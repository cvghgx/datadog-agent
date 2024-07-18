// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package uprobes

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	manager "github.com/DataDog/ebpf-manager"
	"golang.org/x/exp/maps"

	"github.com/DataDog/datadog-agent/pkg/ebpf"
	"github.com/DataDog/datadog-agent/pkg/network/usm/sharedlibraries"
	"github.com/DataDog/datadog-agent/pkg/network/usm/utils"
	"github.com/DataDog/datadog-agent/pkg/process/monitor"
	"github.com/DataDog/datadog-agent/pkg/util/kernel"
)

// ExcludeMode defines the different optiont to exclude processes from attachment
type ExcludeMode uint8

const (
	// ExcludeSelf excludes the agent's own PID
	ExcludeSelf ExcludeMode = 1 << iota
	// ExcludeInternal excludes internal DataDog processes
	ExcludeInternal
	// ExcludeBuildkit excludes buildkitd processes
	ExcludeBuildkit
	// ExcludeContainerdTmp excludes containerd tmp mounts
	ExcludeContainerdTmp
)

const procFSUpdateTimeout = 10 * time.Millisecond

var (
	// ErrSelfExcluded is returned when the PID is the same as the agent's PID.
	ErrSelfExcluded = errors.New("self-excluded")
	// ErrInternalDDogProcessRejected is returned when the PID is an internal datadog process.
	ErrInternalDDogProcessRejected = errors.New("internal datadog process rejected")
	// ErrNoMatchingRule is returned when no rule matches the shared library path.
	ErrNoMatchingRule = errors.New("no matching rule")
	// regex that defines internal DataDog processes
	internalProcessRegex = regexp.MustCompile("datadog-agent/.*/((process|security|trace)-agent|system-probe|agent)")
)

// AttachTarget defines the target to which we should attach the probes, libraries or executables
type AttachTarget uint8

const (
	// AttachToExecutable attaches to the main executable
	AttachToExecutable AttachTarget = 1 << iota
	// AttachToSharedLibraries attaches to shared libraries
	AttachToSharedLibraries
)

// AttachRule defines how to attach a certain set of probes. Uprobes can be attached
// to shared libraries or executables, this structure tells the attacher which ones to
// select and to which targets to do it.
type AttachRule struct {
	// LibraryNameRegex defines which libraries should be matched by this rule
	LibraryNameRegex *regexp.Regexp
	// Targets defines the targets to which we should attach the probes, shared libraries and/or executables
	Targets AttachTarget
	// ProbesSelectors defines which probes should be attached and how should we validate
	// the attachment (e.g., whether we need all probes active or just one of them, or in a best-effort basis)
	ProbesSelector []manager.ProbesSelector
}

func (r *AttachRule) canTarget(target AttachTarget) bool {
	return r.Targets&target != 0
}

func (r *AttachRule) matchesLibrary(path string) bool {
	return r.canTarget(AttachToSharedLibraries) && r.LibraryNameRegex != nil && r.LibraryNameRegex.MatchString(path)
}

func (r *AttachRule) matchesExecutable(_ string) bool {
	return r.canTarget(AttachToExecutable)
}

// AttacherConfig defines the configuration for the attacher
type AttacherConfig struct {
	// Rules defines a series of rules that tell the attacher how to attach the probes
	Rules []*AttachRule

	// ScanTerminatedProcessesInterval defines the interval at which we scan for terminated processes. Set
	// to zero to disable
	ScanTerminatedProcessesInterval time.Duration

	// ProcRoot is the root directory of the proc filesystem
	ProcRoot string

	// ExcludeTargets defines the targets that should be excluded from the attacher
	ExcludeTargets ExcludeMode

	// EbpfConfig is the configuration for the eBPF program
	EbpfConfig *ebpf.Config

	// PerformInitialScan defines if the attacher should perform an initial scan of the processes before starting the monitor
	PerformInitialScan bool

	// ProcessMonitorEventStream defines whether the process monitor is using the event stream
	ProcessMonitorEventStream bool
}

// SetDefaults configures the AttacherConfig with default values for those fields for which the compiler
// defaults are not enough
func (ac *AttacherConfig) SetDefaults() {
	if ac.ScanTerminatedProcessesInterval == 0 {
		ac.ScanTerminatedProcessesInterval = 30 * time.Minute
	}

	if ac.ProcRoot == "" {
		ac.ProcRoot = "/proc"
	}

	if ac.EbpfConfig == nil {
		ac.EbpfConfig = ebpf.NewConfig()
	}
}

// ProbeManager is an interface that defines the methods that a Manager implements,
// so that we can replace it in tests for a mock object
type ProbeManager interface {
	AddHook(string, *manager.Probe) error
	DetachHook(manager.ProbeIdentificationPair) error
	GetProbe(manager.ProbeIdentificationPair) (*manager.Probe, bool)
}

// FileRegistry is an interface that defines the methods that a FileRegistry implements, so that we can replace it in tests for a mock object
type FileRegistry interface {
	Register(namespacedPath string, pid uint32, activationCB, deactivationCB func(utils.FilePath) error) error
	Unregister(uint32) error
	Clear()
	GetRegisteredProcesses() map[uint32]struct{}
}

// AttachCallback is a callback that is called whenever a probe is attached successfully
type AttachCallback func(*manager.Probe, *utils.FilePath)

// UprobeAttacher is a struct that handles the attachment of uprobes to processes and libraries
type UprobeAttacher struct {
	name         string
	done         chan struct{}
	wg           sync.WaitGroup
	config       *AttacherConfig
	fileRegistry FileRegistry
	manager      ProbeManager
	inspector    BinaryInspector

	// pathToAttachedProbes maps a filesystem path to the probes attached to it. Used to detach them
	// once the path is no longer used.
	pathToAttachedProbes   map[string][]manager.ProbeIdentificationPair
	onAttachCallback       AttachCallback
	soWatcher              *sharedlibraries.EbpfProgram
	handlesLibrariesCached *bool
}

// NewUprobeAttacher creates a new UprobeAttacher. Receives as arguments
// the name of the attacher, the configuration, the probe manage (ebpf.Manager usually), a callback to be called
// whenever a probe is attached (optional, can be nil), and the binary inspector to be used (e.g., to attach to
// Go functions we need to inspect the binary in a different way)
func NewUprobeAttacher(name string, config *AttacherConfig, mgr ProbeManager, onAttachCallback AttachCallback, inspector BinaryInspector) (*UprobeAttacher, error) {
	config.SetDefaults()

	ua := &UprobeAttacher{
		name:                 name,
		config:               config,
		fileRegistry:         utils.NewFileRegistry(name),
		manager:              mgr,
		onAttachCallback:     onAttachCallback,
		pathToAttachedProbes: make(map[string][]manager.ProbeIdentificationPair),
		done:                 make(chan struct{}),
		inspector:            inspector,
	}

	if ua.handlesLibraries() {
		ua.soWatcher = sharedlibraries.NewEBPFProgram(config.EbpfConfig)
	}

	utils.AddAttacher(name, ua)

	return ua, nil
}

func (ua *UprobeAttacher) handlesLibraries() bool {
	if ua.handlesLibrariesCached != nil {
		return *ua.handlesLibrariesCached
	}

	result := false
	for _, rule := range ua.config.Rules {
		if rule.LibraryNameRegex != nil {
			result = true
			break
		}
	}
	ua.handlesLibrariesCached = &result
	return result
}

// Start starts the attacher, attaching to the processes and libraries as needed
func (ua *UprobeAttacher) Start() error {
	procMonitor := monitor.GetProcessMonitor()
	err := procMonitor.Initialize(ua.config.ProcessMonitorEventStream)
	if err != nil {
		return fmt.Errorf("error initializing process monitor: %w", err)
	}

	if ua.config.PerformInitialScan {
		err := ua.initialScan()
		if err != nil {
			return fmt.Errorf("error during initial scan: %w", err)
		}
	}

	cleanupExec := procMonitor.SubscribeExec(ua.handleProcessStart)
	cleanupExit := procMonitor.SubscribeExit(ua.handleProcessExit)

	if ua.soWatcher != nil {
		err := ua.soWatcher.Init()
		if err != nil {
			return fmt.Errorf("error initializing shared library program: %w", err)
		}
		err = ua.soWatcher.Start()
		if err != nil {
			return fmt.Errorf("error starting shared library program: %w", err)
		}
	}

	ua.wg.Add(1)
	go func() {
		processSync := time.NewTicker(ua.config.ScanTerminatedProcessesInterval)

		defer func() {
			processSync.Stop()
			cleanupExec()
			cleanupExit()
			procMonitor.Stop()
			ua.fileRegistry.Clear()
			if ua.soWatcher != nil {
				ua.soWatcher.Stop()
			}
			ua.wg.Done()
		}()

		var sharedLibDataChan <-chan *ebpf.DataEvent
		var sharedLibLostChan <-chan uint64

		if ua.soWatcher != nil {
			sharedLibDataChan = ua.soWatcher.GetPerfHandler().DataChannel()
			sharedLibLostChan = ua.soWatcher.GetPerfHandler().LostChannel()
		}

		for {
			select {
			case <-ua.done:
				return
			case <-processSync.C:
				processSet := ua.fileRegistry.GetRegisteredProcesses()
				deletedPids := monitor.FindDeletedProcesses(processSet)
				for deletedPid := range deletedPids {
					ua.handleProcessExit(deletedPid)
				}
			case event, ok := <-sharedLibDataChan:
				if !ok {
					return
				}
				_ = ua.handleLibraryOpen(event)
			case <-sharedLibLostChan:
				// Nothing to do in this case
				break
			}
		}
	}()

	return nil
}

// Stop stops the attache
func (ua *UprobeAttacher) Stop() {
	close(ua.done)
	ua.wg.Wait()
}

func (ua *UprobeAttacher) initialScan() error {
	thisPID, err := kernel.RootNSPID()
	if err != nil {
		return fmt.Errorf("error getting PID of our own process: %w", err)
	}

	err = kernel.WithAllProcs(ua.config.ProcRoot, func(pid int) error {
		if pid == thisPID { // don't scan ourselves
			return nil
		}

		return ua.AttachPIDWithOptions(uint32(pid), true)
	})

	return err
}

// handleProcessStart is called when a new process is started, wraps AttachPIDWithOptions but ignoring the error
// for API compatibility with processMonitor
func (ua *UprobeAttacher) handleProcessStart(pid uint32) {
	_ = ua.AttachPIDWithOptions(pid, false) // Do not try to attach to libraries on process start, it hasn't loaded them yet
}

// handleProcessExit is called when a process finishes, wraps DetachPID but ignoring the error
// for API compatibility with processMonitor
func (ua *UprobeAttacher) handleProcessExit(pid uint32) {
	_ = ua.DetachPID(pid)
}

func (ua *UprobeAttacher) handleLibraryOpen(event *ebpf.DataEvent) error {
	defer event.Done()

	libpath := sharedlibraries.ToLibPath(event.Data)
	path := sharedlibraries.ToBytes(&libpath)

	err := ua.AttachLibrary(string(path), libpath.Pid)

	return err
}

// AttachLibrary attaches the probes to the given library, opened by a given PID
func (ua *UprobeAttacher) AttachLibrary(path string, pid uint32) error {
	if (ua.config.ExcludeTargets&ExcludeSelf) != 0 && int(pid) == os.Getpid() {
		return ErrSelfExcluded
	}

	matchingRules := ua.getRulesForLibrary(path)
	if len(matchingRules) == 0 {
		return ErrNoMatchingRule
	}

	registerCB := func(path utils.FilePath) error {
		return ua.attachToBinary(path, matchingRules)
	}
	unregisterCB := func(path utils.FilePath) error {
		return ua.detachFromBinary(path)
	}

	return ua.fileRegistry.Register(path, pid, registerCB, unregisterCB)
}

// getRulesForLibrary returns the rules that match the given library path
func (ua *UprobeAttacher) getRulesForLibrary(path string) []*AttachRule {
	var matchedRules []*AttachRule

	for _, rule := range ua.config.Rules {
		if rule.matchesLibrary(path) {
			matchedRules = append(matchedRules, rule)
		}
	}
	return matchedRules
}

// getRulesForExecutable returns the rules that match the given executable
func (ua *UprobeAttacher) getRulesForExecutable(path string) []*AttachRule {
	var matchedRules []*AttachRule

	for _, rule := range ua.config.Rules {
		if rule.matchesExecutable(path) {
			matchedRules = append(matchedRules, rule)
		}
	}
	return matchedRules
}

// getExecutablePath resolves the executable of the given PID looking in procfs. Automatically
// handles delays in procfs updates. Will return an error if the path cannot be resolved
func (ua *UprobeAttacher) getExecutablePath(pid uint32) (string, error) {
	pidAsStr := strconv.FormatUint(uint64(pid), 10)
	exePath := filepath.Join(ua.config.ProcRoot, pidAsStr, "exe")

	var binPath string
	err := errors.New("iteration start")
	end := time.Now().Add(procFSUpdateTimeout)

	for err != nil && end.After(time.Now()) {
		binPath, err = os.Readlink(exePath)
		if err != nil {
			time.Sleep(time.Millisecond)
		}
	}

	if err != nil {
		return "", err
	}

	return binPath, nil
}

// AttachPID attaches the corresponding probes to a given pid
func (ua *UprobeAttacher) AttachPID(pid uint32) error {
	return ua.AttachPIDWithOptions(pid, true)
}

// AttachPIDWithOptions attaches the corresponding probes to a given pid
func (ua *UprobeAttacher) AttachPIDWithOptions(pid uint32, attachToLibs bool) error {
	if (ua.config.ExcludeTargets&ExcludeSelf) != 0 && int(pid) == os.Getpid() {
		return ErrSelfExcluded
	}

	binPath, err := ua.getExecutablePath(pid)
	if err != nil {
		return err
	}

	if (ua.config.ExcludeTargets&ExcludeInternal) != 0 && internalProcessRegex.MatchString(binPath) {
		return ErrInternalDDogProcessRejected
	}

	matchingRules := ua.getRulesForExecutable(binPath)
	registerCB := func(path utils.FilePath) error {
		return ua.attachToBinary(path, matchingRules)
	}
	unregisterCB := func(path utils.FilePath) error {
		return ua.detachFromBinary(path)
	}

	if len(matchingRules) != 0 {
		err = ua.fileRegistry.Register(binPath, pid, registerCB, unregisterCB)
		if err != nil {
			return err
		}
	}

	if attachToLibs && ua.handlesLibraries() {
		err = ua.attachToLibrariesOfPID(pid)
		if err != nil {
			return err
		}
	}

	return nil
}

// DetachPID detaches the uprobes attached to a PID
func (ua *UprobeAttacher) DetachPID(pid uint32) error {
	return ua.fileRegistry.Unregister(pid)
}

const (
	// Defined in https://man7.org/linux/man-pages/man5/proc.5.html.
	taskCommLen = 16
)

var (
	taskCommLenBufferPool = sync.Pool{
		New: func() any {
			buf := make([]byte, taskCommLen)
			return &buf
		},
	}
	buildKitProcessName = []byte("buildkitd")
)

func isBuildKit(procRoot string, pid uint32) bool {
	filePath := filepath.Join(procRoot, strconv.Itoa(int(pid)), "comm")

	var file *os.File
	err := errors.New("iteration start")
	for i := 0; err != nil && i < 30; i++ {
		file, err = os.Open(filePath)
		if err != nil {
			time.Sleep(1 * time.Millisecond)
		}
	}

	buf := taskCommLenBufferPool.Get().(*[]byte)
	defer taskCommLenBufferPool.Put(buf)
	n, err := file.Read(*buf)
	if err != nil {
		// short living process can hit here, or slow start of another process.
		return false
	}
	return bytes.Equal(bytes.TrimSpace((*buf)[:n]), buildKitProcessName)
}

func isContainerdTmpMount(path string) bool {
	return strings.Contains(path, "tmpmounts/containerd-mount")
}

// getUID() return a key of length 5 as the kernel uprobe registration path is limited to a length of 64
// ebpf-manager/utils.go:GenerateEventName() MaxEventNameLen = 64
// MAX_EVENT_NAME_LEN (linux/kernel/trace/trace.h)
//
// Length 5 is arbitrary value as the full string of the eventName format is
//
//	fmt.Sprintf("%s_%.*s_%s_%s", probeType, maxFuncNameLen, functionName, UID, attachPIDstr)
//
// functionName is variable but with a minimum guarantee of 10 chars
func getUID(lib utils.PathIdentifier) string {
	return lib.Key()[:5]
}

func parseSymbolFromEBPFProbeName(probeName string) (symbol string, isManualReturn bool, err error) {
	parts := strings.Split(probeName, "__")
	if len(parts) < 2 {
		err = fmt.Errorf("invalid probe name %s, no double underscore (__) separating probe type and function name", probeName)
		return
	}

	symbol = parts[1]
	if len(parts) > 2 {
		if parts[2] == "return" {
			isManualReturn = true
		} else {
			err = fmt.Errorf("invalid probe name %s, unexpected third part %s. Format should be probeType__funcName[__return]", probeName, parts[2])
			return
		}
	}

	return
}

func (ua *UprobeAttacher) attachToBinary(fpath utils.FilePath, matchingRules []*AttachRule) error {
	// TODO: Retrieve this information once and reuse it
	if isBuildKit(ua.config.ProcRoot, fpath.PID) {
		return fmt.Errorf("process %d is buildkitd, skipping", fpath.PID)
	} else if isContainerdTmpMount(fpath.HostPath) {
		return fmt.Errorf("path %s from process %d is tempmount of containerd, skipping", fpath.HostPath, fpath.PID)
	}

	symbolsToRequest, err := ua.computeSymbolsToRequest(matchingRules)
	if err != nil {
		return fmt.Errorf("error computing symbols to request for rules %+v: %w", matchingRules, err)
	}

	inspectResult, isAttachable, err := ua.inspector.Inspect(fpath, symbolsToRequest)
	if err != nil {
		return fmt.Errorf("error inspecting %s: %w", fpath.HostPath, err)
	}
	if !isAttachable {
		return fmt.Errorf("incompatible binary %s", fpath.HostPath)
	}

	uid := getUID(fpath.ID)

	for _, rule := range matchingRules {
		for _, selector := range rule.ProbesSelector {
			_, isBestEffort := selector.(*manager.BestEffort)
			for _, probeID := range selector.GetProbesIdentificationPairList() {
				symbol, isManualReturn, err := parseSymbolFromEBPFProbeName(probeID.EBPFFuncName)
				if err != nil {
					return fmt.Errorf("error parsing probe name %s: %w", probeID.EBPFFuncName, err)
				}
				data, found := inspectResult[symbol]
				if !found {
					if isBestEffort {
						continue
					}
					// This should not happen, as getAvailableRequestedSymbols should have already
					// returned an error if mandatory symbols weren't found. However and for safety,
					// we'll check again and return an error if the symbol is not found.
					return fmt.Errorf("symbol %s not found in %s", symbol, fpath.HostPath)
				}

				var locationsToAttach []uint64
				if isManualReturn {
					locationsToAttach = data.ReturnLocations
				} else {
					locationsToAttach = []uint64{data.EntryLocation}
				}

				for i, location := range locationsToAttach {
					newProbeID := manager.ProbeIdentificationPair{
						EBPFFuncName: probeID.EBPFFuncName,
						UID:          fmt.Sprintf("%s_%d", uid, i), // Make UID unique even if we have multiple locations
					}

					probe, found := ua.manager.GetProbe(newProbeID)
					if found {
						// We have already probed this process, just ensure it's running and skip it
						if !probe.IsRunning() {
							err := probe.Attach()
							if err != nil {
								return fmt.Errorf("cannot attach running probe %v: %w", newProbeID, err)
							}
						}
						continue
					}

					newProbe := &manager.Probe{
						ProbeIdentificationPair: newProbeID,
						BinaryPath:              fpath.HostPath,
						UprobeOffset:            location,
						HookFuncName:            symbol,
					}
					err = ua.manager.AddHook("", newProbe)
					if err != nil {
						return fmt.Errorf("error attaching probe %+v: %w", newProbe, err)
					}

					ebpf.AddProgramNameMapping(newProbe.ID(), newProbe.EBPFFuncName, ua.name)
					ua.pathToAttachedProbes[fpath.HostPath] = append(ua.pathToAttachedProbes[fpath.HostPath], newProbeID)

					if ua.onAttachCallback != nil {
						ua.onAttachCallback(newProbe, &fpath)
					}
				}

			}

			manager, ok := ua.manager.(*manager.Manager)
			if ok {
				if err := selector.RunValidator(manager); err != nil {
					return fmt.Errorf("error validating probes: %w", err)
				}
			}
		}
	}

	return nil
}

func (ua *UprobeAttacher) computeSymbolsToRequest(rules []*AttachRule) ([]SymbolRequest, error) {
	var requests []SymbolRequest
	for _, rule := range rules {
		for _, selector := range rule.ProbesSelector {
			_, isBestEffort := selector.(*manager.BestEffort)
			for _, selector := range selector.GetProbesIdentificationPairList() {
				symbol, isManualReturn, err := parseSymbolFromEBPFProbeName(selector.EBPFFuncName)
				if err != nil {
					return nil, fmt.Errorf("error parsing probe name %s: %w", selector.EBPFFuncName, err)
				}

				requests = append(requests, SymbolRequest{
					Name:                   symbol,
					IncludeReturnLocations: isManualReturn,
					BestEffort:             isBestEffort,
				})
			}
		}
	}

	return requests, nil
}

func (ua *UprobeAttacher) detachFromBinary(fpath utils.FilePath) error {
	for _, probeID := range ua.pathToAttachedProbes[fpath.HostPath] {
		err := ua.manager.DetachHook(probeID)
		if err != nil {
			return fmt.Errorf("error detaching probe %+v: %w", probeID, err)
		}
	}

	ua.inspector.Cleanup(fpath)

	return nil
}

func (ua *UprobeAttacher) getLibrariesFromMapsFile(pid int) ([]string, error) {
	mapsPath := filepath.Join(ua.config.ProcRoot, strconv.Itoa(int(pid)), "maps")
	mapsFile, err := os.Open(mapsPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open maps file at %s: %w", mapsPath, err)
	}
	defer mapsFile.Close()

	scanner := bufio.NewScanner(bufio.NewReader(mapsFile))
	libs := make(map[string]struct{})
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Fields(line)
		// ensuring we have exactly 6 elements (skip '(deleted)' entries) in the line, and the 4th element (inode) is
		// not zero (indicates it is a path, and not an anonymous path).
		if len(cols) == 6 && cols[4] != "0" {
			libs[cols[5]] = struct{}{}
		}
	}

	return maps.Keys(libs), nil
}

func (ua *UprobeAttacher) attachToLibrariesOfPID(pid uint32) error {
	registerErrors := make([]error, 0)
	successfulMatches := make([]string, 0)
	libs, err := ua.getLibrariesFromMapsFile(int(pid))
	if err != nil {
		return err
	}
	for _, libpath := range libs {
		err := ua.AttachLibrary(libpath, pid)

		if err == nil {
			successfulMatches = append(successfulMatches, libpath)
		} else if !errors.Is(err, ErrNoMatchingRule) {
			registerErrors = append(registerErrors, err)
		}
	}

	if len(successfulMatches) == 0 {
		if len(registerErrors) == 0 {
			return nil // No libraries found to attach
		}
		return fmt.Errorf("no rules matched for pid %d, errors: %v", pid, registerErrors)
	}
	if len(registerErrors) > 0 {
		return fmt.Errorf("partially hooked (%v), errors while attaching pid %d: %v", successfulMatches, pid, registerErrors)
	}
	return nil
}
