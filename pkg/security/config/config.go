// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package config holds config related files
package config

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	sysconfig "github.com/DataDog/datadog-agent/cmd/system-probe/config"
	logsconfig "github.com/DataDog/datadog-agent/comp/logs/agent/config"
	coreconfig "github.com/DataDog/datadog-agent/pkg/config"
	configUtils "github.com/DataDog/datadog-agent/pkg/config/utils"
	logshttp "github.com/DataDog/datadog-agent/pkg/logs/client/http"
	pconfig "github.com/DataDog/datadog-agent/pkg/security/probe/config"
	"github.com/DataDog/datadog-agent/pkg/security/secl/compiler/eval"
	"github.com/DataDog/datadog-agent/pkg/security/secl/model"
	"github.com/DataDog/datadog-agent/pkg/security/seclog"
	"github.com/DataDog/datadog-agent/pkg/security/utils"
)

const (
	// ADMinMaxDumSize represents the minimum value for runtime_security_config.activity_dump.max_dump_size
	ADMinMaxDumSize = 100
)

// Policy represents a policy file in the configuration file
type Policy struct {
	Name  string   `mapstructure:"name"`
	Files []string `mapstructure:"files"`
	Tags  []string `mapstructure:"tags"`
}

// RuntimeSecurityConfig holds the configuration for the runtime security agent
type RuntimeSecurityConfig struct {
	// RuntimeEnabled defines if the runtime security module should be enabled
	RuntimeEnabled bool
	// PoliciesDir defines the folder in which the policy files are located
	PoliciesDir string
	// WatchPoliciesDir activate policy dir inotify
	WatchPoliciesDir bool
	// PolicyMonitorEnabled enable policy monitoring
	PolicyMonitorEnabled bool
	// PolicyMonitorPerRuleEnabled enabled per-rule policy monitoring
	PolicyMonitorPerRuleEnabled bool
	// PolicyMonitorReportInternalPolicies enable internal policies monitoring
	PolicyMonitorReportInternalPolicies bool
	// SocketPath is the path to the socket that is used to communicate with the security agent
	SocketPath string
	// EventServerBurst defines the maximum burst of events that can be sent over the grpc server
	EventServerBurst int
	// EventServerRate defines the grpc server rate at which events can be sent
	EventServerRate int
	// EventServerRetention defines an event retention period so that some fields can be resolved
	EventServerRetention time.Duration
	// FIMEnabled determines whether fim rules will be loaded
	FIMEnabled bool
	// SelfTestEnabled defines if the self tests should be executed at startup or not
	SelfTestEnabled bool
	// SelfTestSendReport defines if a self test event will be emitted
	SelfTestSendReport bool
	// RemoteConfigurationEnabled defines whether to use remote monitoring
	RemoteConfigurationEnabled bool
	// RemoteConfigurationDumpPolicies defines whether to dump remote config policy
	RemoteConfigurationDumpPolicies bool
	// LogPatterns pattern to be used by the logger for trace level
	LogPatterns []string
	// LogTags tags to be used by the logger for trace level
	LogTags []string
	// HostServiceName string
	HostServiceName string
	// OnDemandEnabled defines whether the on-demand probes should be enabled
	OnDemandEnabled bool
	// OnDemandRateLimiterEnabled defines whether the on-demand probes rate limit getting hit disabled the on demand probes
	OnDemandRateLimiterEnabled bool

	// InternalMonitoringEnabled determines if the monitoring events of the agent should be sent to Datadog
	InternalMonitoringEnabled bool

	// ActivityDumpEnabled defines if the activity dump manager should be enabled
	ActivityDumpEnabled bool
	// ActivityDumpCleanupPeriod defines the period at which the activity dump manager should perform its cleanup
	// operation.
	ActivityDumpCleanupPeriod time.Duration
	// ActivityDumpTagsResolutionPeriod defines the period at which the activity dump manager should try to resolve
	// missing container tags.
	ActivityDumpTagsResolutionPeriod time.Duration
	// ActivityDumpLoadControlPeriod defines the period at which the activity dump manager should trigger the load controller
	ActivityDumpLoadControlPeriod time.Duration
	// ActivityDumpLoadControlMinDumpTimeout defines minimal duration of a activity dump recording
	ActivityDumpLoadControlMinDumpTimeout time.Duration

	// ActivityDumpTracedCgroupsCount defines the maximum count of cgroups that should be monitored concurrently. Leave this parameter to 0 to prevent the generation
	// of activity dumps based on cgroups.
	ActivityDumpTracedCgroupsCount int
	// ActivityDumpTracedEventTypes defines the list of events that should be captured in an activity dump. Leave this
	// parameter empty to monitor all event types. If not already present, the `exec` event will automatically be added
	// to this list.
	ActivityDumpTracedEventTypes []model.EventType
	// ActivityDumpCgroupDumpTimeout defines the cgroup activity dumps timeout.
	ActivityDumpCgroupDumpTimeout time.Duration
	// ActivityDumpRateLimiter defines the kernel rate of max events per sec for activity dumps.
	ActivityDumpRateLimiter int
	// ActivityDumpCgroupWaitListTimeout defines the time to wait before a cgroup can be dumped again.
	ActivityDumpCgroupWaitListTimeout time.Duration
	// ActivityDumpCgroupDifferentiateArgs defines if system-probe should differentiate process nodes using process
	// arguments for dumps.
	ActivityDumpCgroupDifferentiateArgs bool
	// ActivityDumpLocalStorageDirectory defines the output directory for the activity dumps and graphs. Leave
	// this field empty to prevent writing any output to disk.
	ActivityDumpLocalStorageDirectory string
	// ActivityDumpLocalStorageFormats defines the formats that should be used to persist the activity dumps locally.
	ActivityDumpLocalStorageFormats []StorageFormat
	// ActivityDumpLocalStorageCompression defines if the local storage should compress the persisted data.
	ActivityDumpLocalStorageCompression bool
	// ActivityDumpLocalStorageMaxDumpsCount defines the maximum count of activity dumps that should be kept locally.
	// When the limit is reached, the oldest dumps will be deleted first.
	ActivityDumpLocalStorageMaxDumpsCount int
	// ActivityDumpSyscallMonitorPeriod defines the minimum amount of time to wait between 2 syscalls event for the same
	// process.
	ActivityDumpSyscallMonitorPeriod time.Duration
	// ActivityDumpMaxDumpCountPerWorkload defines the maximum amount of dumps that the agent should send for a workload
	ActivityDumpMaxDumpCountPerWorkload int
	// ActivityDumpWorkloadDenyList defines the list of workloads for which we shouldn't generate dumps. Workloads should
	// be provided as strings in the following format "{image_name}:[{image_tag}|*]". If "*" is provided instead of a
	// specific image tag, then the entry will match any workload with the input {image_name} regardless of their tag.
	ActivityDumpWorkloadDenyList []string
	// ActivityDumpTagRulesEnabled enable the tagging of nodes with matched rules
	ActivityDumpTagRulesEnabled bool
	// ActivityDumpSilentWorkloadsDelay defines the minimum amount of time to wait before the activity dump manager will start tracing silent workloads
	ActivityDumpSilentWorkloadsDelay time.Duration
	// ActivityDumpSilentWorkloadsTicker configures ticker that will check if a workload is silent and should be traced
	ActivityDumpSilentWorkloadsTicker time.Duration
	// ActivityDumpAutoSuppressionEnabled bool do not send event if part of a dump
	ActivityDumpAutoSuppressionEnabled bool

	// # Dynamic configuration fields:
	// ActivityDumpMaxDumpSize defines the maximum size of a dump
	ActivityDumpMaxDumpSize func() int

	// SecurityProfileEnabled defines if the Security Profile manager should be enabled
	SecurityProfileEnabled bool
	// SecurityProfileMaxImageTags defines the maximum number of profile versions to maintain
	SecurityProfileMaxImageTags int
	// SecurityProfileDir defines the directory in which Security Profiles are stored
	SecurityProfileDir string
	// SecurityProfileWatchDir defines if the Security Profiles directory should be monitored
	SecurityProfileWatchDir bool
	// SecurityProfileCacheSize defines the count of Security Profiles held in cache
	SecurityProfileCacheSize int
	// SecurityProfileMaxCount defines the maximum number of Security Profiles that may be evaluated concurrently
	SecurityProfileMaxCount int
	// SecurityProfileRCEnabled defines if remote-configuration is enabled
	SecurityProfileRCEnabled bool
	// SecurityProfileDNSMatchMaxDepth defines the max depth of subdomain to be matched for DNS anomaly detection (0 to match everything)
	SecurityProfileDNSMatchMaxDepth int

	// SecurityProfileAutoSuppressionEnabled do not send event if part of a profile
	SecurityProfileAutoSuppressionEnabled bool
	// SecurityProfileAutoSuppressionEventTypes defines the list of event types the can be auto suppressed using security profiles
	SecurityProfileAutoSuppressionEventTypes []model.EventType

	// AnomalyDetectionEventTypes defines the list of events that should be allowed to generate anomaly detections
	AnomalyDetectionEventTypes []model.EventType
	// AnomalyDetectionDefaultMinimumStablePeriod defines the default minimum amount of time during which the events
	// that diverge from their profiles are automatically added in their profiles without triggering an anomaly detection
	// event.
	AnomalyDetectionDefaultMinimumStablePeriod time.Duration
	// AnomalyDetectionMinimumStablePeriods defines the minimum amount of time per event type during which the events
	// that diverge from their profiles are automatically added in their profiles without triggering an anomaly detection
	// event.
	AnomalyDetectionMinimumStablePeriods map[model.EventType]time.Duration
	// AnomalyDetectionUnstableProfileTimeThreshold defines the maximum amount of time to wait until a profile that
	// hasn't reached a stable state is considered as unstable.
	AnomalyDetectionUnstableProfileTimeThreshold time.Duration
	// AnomalyDetectionUnstableProfileSizeThreshold defines the maximum size a profile can reach past which it is
	// considered unstable
	AnomalyDetectionUnstableProfileSizeThreshold int64
	// AnomalyDetectionWorkloadWarmupPeriod defines the duration we ignore the anomaly detections for
	// because of workload warm up
	AnomalyDetectionWorkloadWarmupPeriod time.Duration
	// AnomalyDetectionRateLimiterPeriod is the duration during which a limited number of anomaly detection events are allowed
	AnomalyDetectionRateLimiterPeriod time.Duration
	// AnomalyDetectionRateLimiterNumEventsAllowed is the number of anomaly detection events allowed per duration by the rate limiter
	AnomalyDetectionRateLimiterNumEventsAllowed int
	// AnomalyDetectionRateLimiterNumKeys is the number of keys in the rate limiter
	AnomalyDetectionRateLimiterNumKeys int
	// AnomalyDetectionTagRulesEnabled defines if the events that triggered anomaly detections should be tagged with the
	// rules they might have matched.
	AnomalyDetectionTagRulesEnabled bool
	// AnomalyDetectionSilentRuleEventsEnabled do not send rule event if also part of an anomaly event
	AnomalyDetectionSilentRuleEventsEnabled bool
	// AnomalyDetectionEnabled defines if we should send anomaly detection events
	AnomalyDetectionEnabled bool

	// SBOMResolverEnabled defines if the SBOM resolver should be enabled
	SBOMResolverEnabled bool
	// SBOMResolverWorkloadsCacheSize defines the count of SBOMs to keep in memory in order to prevent re-computing
	// the SBOMs of short-lived and periodical workloads
	SBOMResolverWorkloadsCacheSize int
	// SBOMResolverHostEnabled defines if the SBOM resolver should compute the host's SBOM
	SBOMResolverHostEnabled bool

	// HashResolverEnabled defines if the hash resolver should be enabled
	HashResolverEnabled bool
	// HashResolverMaxFileSize defines the maximum size of the files that the hash resolver is allowed to hash
	HashResolverMaxFileSize int64
	// HashResolverMaxHashRate defines the rate at which the hash resolver may compute hashes
	HashResolverMaxHashRate int
	// HashResolverMaxHashBurst defines the burst of files for which the hash resolver may compute a hash
	HashResolverMaxHashBurst int
	// HashResolverHashAlgorithms defines the hashes that hash resolver needs to compute
	HashResolverHashAlgorithms []model.HashAlgorithm
	// HashResolverEventTypes defines the list of event which files may be hashed
	HashResolverEventTypes []model.EventType
	// HashResolverCacheSize defines the number of hashes to keep in cache
	HashResolverCacheSize int
	// HashResolverReplace is used to apply specific hash to specific file path
	HashResolverReplace map[string]string

	// UserSessionsCacheSize defines the size of the User Sessions cache size
	UserSessionsCacheSize int

	// EBPFLessEnabled enables the ebpfless probe
	EBPFLessEnabled bool
	// EBPFLessSocket defines the socket used for the communication between system-probe and the ebpfless source
	EBPFLessSocket string

	// Enforcement capabilities
	EnforcementEnabled           bool
	EnforcementRawSyscallEnabled bool
	EnforcementRuleSourceAllowed []string

	//WindowsFilenameCacheSize is the max number of filenames to cache
	WindowsFilenameCacheSize int
	//WindowsRegistryCacheSize is the max number of registry paths to cache
	WindowsRegistryCacheSize int

	// ETWEventsChannelSize windows specific ETW channel buffer size
	ETWEventsChannelSize int

	//ETWEventsMaxBuffers sets the maximumbuffers argument to ETW
	ETWEventsMaxBuffers int

	// WindowsProbeChannelUnbuffered defines if the windows probe channel should be unbuffered
	WindowsProbeBlockOnChannelSend bool

	// IMDSIPv4 is used to provide a custom IP address for the IMDS endpoint
	IMDSIPv4 uint32
}

// Config defines a security config
type Config struct {
	// Probe Config
	Probe *pconfig.Config

	// CWS specific parameters
	RuntimeSecurity *RuntimeSecurityConfig
}

// NewConfig returns a new Config object
func NewConfig() (*Config, error) {
	probeConfig, err := pconfig.NewConfig()
	if err != nil {
		return nil, err
	}

	rsConfig, err := NewRuntimeSecurityConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Probe:           probeConfig,
		RuntimeSecurity: rsConfig,
	}, nil
}

// NewRuntimeSecurityConfig returns the runtime security (CWS) config, build from the system probe one
func NewRuntimeSecurityConfig() (*RuntimeSecurityConfig, error) {
	sysconfig.Adjust(coreconfig.SystemProbe())

	eventTypeStrings := map[string]model.EventType{}

	var eventType model.EventType
	for i := uint64(0); i != uint64(model.MaxKernelEventType); i++ {
		eventType = model.EventType(i)
		eventTypeStrings[eventType.String()] = eventType
	}

	// parseEventTypeStringSlice converts a string list to a list of event types
	parseEventTypeStringSlice := func(eventTypes []string) []model.EventType {
		var output []model.EventType
		for _, eventTypeStr := range eventTypes {
			if eventType := eventTypeStrings[eventTypeStr]; eventType != model.UnknownEventType {
				output = append(output, eventType)
			}
		}
		return output
	}

	rsConfig := &RuntimeSecurityConfig{
		RuntimeEnabled:                 coreconfig.SystemProbe().GetBool("runtime_security_config.enabled"),
		FIMEnabled:                     coreconfig.SystemProbe().GetBool("runtime_security_config.fim_enabled"),
		WindowsFilenameCacheSize:       coreconfig.SystemProbe().GetInt("runtime_security_config.windows_filename_cache_max"),
		WindowsRegistryCacheSize:       coreconfig.SystemProbe().GetInt("runtime_security_config.windows_registry_cache_max"),
		ETWEventsChannelSize:           coreconfig.SystemProbe().GetInt("runtime_security_config.etw_events_channel_size"),
		ETWEventsMaxBuffers:            coreconfig.SystemProbe().GetInt("runtime_security_config.etw_events_max_buffers"),
		WindowsProbeBlockOnChannelSend: coreconfig.SystemProbe().GetBool("runtime_security_config.windows_probe_block_on_channel_send"),

		SocketPath:           coreconfig.SystemProbe().GetString("runtime_security_config.socket"),
		EventServerBurst:     coreconfig.SystemProbe().GetInt("runtime_security_config.event_server.burst"),
		EventServerRate:      coreconfig.SystemProbe().GetInt("runtime_security_config.event_server.rate"),
		EventServerRetention: coreconfig.SystemProbe().GetDuration("runtime_security_config.event_server.retention"),

		SelfTestEnabled:                 coreconfig.SystemProbe().GetBool("runtime_security_config.self_test.enabled"),
		SelfTestSendReport:              coreconfig.SystemProbe().GetBool("runtime_security_config.self_test.send_report"),
		RemoteConfigurationEnabled:      isRemoteConfigEnabled(),
		RemoteConfigurationDumpPolicies: coreconfig.SystemProbe().GetBool("runtime_security_config.remote_configuration.dump_policies"),

		OnDemandEnabled:            coreconfig.SystemProbe().GetBool("runtime_security_config.on_demand.enabled"),
		OnDemandRateLimiterEnabled: coreconfig.SystemProbe().GetBool("runtime_security_config.on_demand.rate_limiter.enabled"),

		// policy & ruleset
		PoliciesDir:                         coreconfig.SystemProbe().GetString("runtime_security_config.policies.dir"),
		WatchPoliciesDir:                    coreconfig.SystemProbe().GetBool("runtime_security_config.policies.watch_dir"),
		PolicyMonitorEnabled:                coreconfig.SystemProbe().GetBool("runtime_security_config.policies.monitor.enabled"),
		PolicyMonitorPerRuleEnabled:         coreconfig.SystemProbe().GetBool("runtime_security_config.policies.monitor.per_rule_enabled"),
		PolicyMonitorReportInternalPolicies: coreconfig.SystemProbe().GetBool("runtime_security_config.policies.monitor.report_internal_policies"),

		LogPatterns: coreconfig.SystemProbe().GetStringSlice("runtime_security_config.log_patterns"),
		LogTags:     coreconfig.SystemProbe().GetStringSlice("runtime_security_config.log_tags"),

		// custom events
		InternalMonitoringEnabled: coreconfig.SystemProbe().GetBool("runtime_security_config.internal_monitoring.enabled"),

		// activity dump
		ActivityDumpEnabled:                   coreconfig.SystemProbe().GetBool("runtime_security_config.activity_dump.enabled"),
		ActivityDumpCleanupPeriod:             coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.cleanup_period"),
		ActivityDumpTagsResolutionPeriod:      coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.tags_resolution_period"),
		ActivityDumpLoadControlPeriod:         coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.load_controller_period"),
		ActivityDumpLoadControlMinDumpTimeout: coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.min_timeout"),
		ActivityDumpTracedCgroupsCount:        coreconfig.SystemProbe().GetInt("runtime_security_config.activity_dump.traced_cgroups_count"),
		ActivityDumpTracedEventTypes:          parseEventTypeStringSlice(coreconfig.SystemProbe().GetStringSlice("runtime_security_config.activity_dump.traced_event_types")),
		ActivityDumpCgroupDumpTimeout:         coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.dump_duration"),
		ActivityDumpRateLimiter:               coreconfig.SystemProbe().GetInt("runtime_security_config.activity_dump.rate_limiter"),
		ActivityDumpCgroupWaitListTimeout:     coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.cgroup_wait_list_timeout"),
		ActivityDumpCgroupDifferentiateArgs:   coreconfig.SystemProbe().GetBool("runtime_security_config.activity_dump.cgroup_differentiate_args"),
		ActivityDumpLocalStorageDirectory:     coreconfig.SystemProbe().GetString("runtime_security_config.activity_dump.local_storage.output_directory"),
		ActivityDumpLocalStorageMaxDumpsCount: coreconfig.SystemProbe().GetInt("runtime_security_config.activity_dump.local_storage.max_dumps_count"),
		ActivityDumpLocalStorageCompression:   coreconfig.SystemProbe().GetBool("runtime_security_config.activity_dump.local_storage.compression"),
		ActivityDumpSyscallMonitorPeriod:      coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.syscall_monitor.period"),
		ActivityDumpMaxDumpCountPerWorkload:   coreconfig.SystemProbe().GetInt("runtime_security_config.activity_dump.max_dump_count_per_workload"),
		ActivityDumpTagRulesEnabled:           coreconfig.SystemProbe().GetBool("runtime_security_config.activity_dump.tag_rules.enabled"),
		ActivityDumpSilentWorkloadsDelay:      coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.silent_workloads.delay"),
		ActivityDumpSilentWorkloadsTicker:     coreconfig.SystemProbe().GetDuration("runtime_security_config.activity_dump.silent_workloads.ticker"),
		ActivityDumpWorkloadDenyList:          coreconfig.SystemProbe().GetStringSlice("runtime_security_config.activity_dump.workload_deny_list"),
		ActivityDumpAutoSuppressionEnabled:    coreconfig.SystemProbe().GetBool("runtime_security_config.activity_dump.auto_suppression.enabled"),
		// activity dump dynamic fields
		ActivityDumpMaxDumpSize: func() int {
			mds := coreconfig.SystemProbe().GetInt("runtime_security_config.activity_dump.max_dump_size")
			if mds < ADMinMaxDumSize {
				mds = ADMinMaxDumSize
			}
			return mds * (1 << 10)
		},

		// SBOM resolver
		SBOMResolverEnabled:            coreconfig.SystemProbe().GetBool("runtime_security_config.sbom.enabled"),
		SBOMResolverWorkloadsCacheSize: coreconfig.SystemProbe().GetInt("runtime_security_config.sbom.workloads_cache_size"),
		SBOMResolverHostEnabled:        coreconfig.SystemProbe().GetBool("runtime_security_config.sbom.host.enabled"),

		// Hash resolver
		HashResolverEnabled:        coreconfig.SystemProbe().GetBool("runtime_security_config.hash_resolver.enabled"),
		HashResolverEventTypes:     parseEventTypeStringSlice(coreconfig.SystemProbe().GetStringSlice("runtime_security_config.hash_resolver.event_types")),
		HashResolverMaxFileSize:    coreconfig.SystemProbe().GetInt64("runtime_security_config.hash_resolver.max_file_size"),
		HashResolverHashAlgorithms: parseHashAlgorithmStringSlice(coreconfig.SystemProbe().GetStringSlice("runtime_security_config.hash_resolver.hash_algorithms")),
		HashResolverMaxHashBurst:   coreconfig.SystemProbe().GetInt("runtime_security_config.hash_resolver.max_hash_burst"),
		HashResolverMaxHashRate:    coreconfig.SystemProbe().GetInt("runtime_security_config.hash_resolver.max_hash_rate"),
		HashResolverCacheSize:      coreconfig.SystemProbe().GetInt("runtime_security_config.hash_resolver.cache_size"),
		HashResolverReplace:        coreconfig.SystemProbe().GetStringMapString("runtime_security_config.hash_resolver.replace"),

		// security profiles
		SecurityProfileEnabled:          coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.enabled"),
		SecurityProfileMaxImageTags:     coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.max_image_tags"),
		SecurityProfileDir:              coreconfig.SystemProbe().GetString("runtime_security_config.security_profile.dir"),
		SecurityProfileWatchDir:         coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.watch_dir"),
		SecurityProfileCacheSize:        coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.cache_size"),
		SecurityProfileMaxCount:         coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.max_count"),
		SecurityProfileRCEnabled:        coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.remote_configuration.enabled"),
		SecurityProfileDNSMatchMaxDepth: coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.dns_match_max_depth"),

		// auto suppression
		SecurityProfileAutoSuppressionEnabled:    coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.auto_suppression.enabled"),
		SecurityProfileAutoSuppressionEventTypes: parseEventTypeStringSlice(coreconfig.SystemProbe().GetStringSlice("runtime_security_config.security_profile.auto_suppression.event_types")),

		// anomaly detection
		AnomalyDetectionEventTypes:                   parseEventTypeStringSlice(coreconfig.SystemProbe().GetStringSlice("runtime_security_config.security_profile.anomaly_detection.event_types")),
		AnomalyDetectionDefaultMinimumStablePeriod:   coreconfig.SystemProbe().GetDuration("runtime_security_config.security_profile.anomaly_detection.default_minimum_stable_period"),
		AnomalyDetectionMinimumStablePeriods:         parseEventTypeDurations(coreconfig.SystemProbe(), "runtime_security_config.security_profile.anomaly_detection.minimum_stable_period"),
		AnomalyDetectionWorkloadWarmupPeriod:         coreconfig.SystemProbe().GetDuration("runtime_security_config.security_profile.anomaly_detection.workload_warmup_period"),
		AnomalyDetectionUnstableProfileTimeThreshold: coreconfig.SystemProbe().GetDuration("runtime_security_config.security_profile.anomaly_detection.unstable_profile_time_threshold"),
		AnomalyDetectionUnstableProfileSizeThreshold: coreconfig.SystemProbe().GetInt64("runtime_security_config.security_profile.anomaly_detection.unstable_profile_size_threshold"),
		AnomalyDetectionRateLimiterPeriod:            coreconfig.SystemProbe().GetDuration("runtime_security_config.security_profile.anomaly_detection.rate_limiter.period"),
		AnomalyDetectionRateLimiterNumKeys:           coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.anomaly_detection.rate_limiter.num_keys"),
		AnomalyDetectionRateLimiterNumEventsAllowed:  coreconfig.SystemProbe().GetInt("runtime_security_config.security_profile.anomaly_detection.rate_limiter.num_events_allowed"),
		AnomalyDetectionTagRulesEnabled:              coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.anomaly_detection.tag_rules.enabled"),
		AnomalyDetectionSilentRuleEventsEnabled:      coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.anomaly_detection.silent_rule_events.enabled"),
		AnomalyDetectionEnabled:                      coreconfig.SystemProbe().GetBool("runtime_security_config.security_profile.anomaly_detection.enabled"),

		// enforcement
		EnforcementEnabled:           coreconfig.SystemProbe().GetBool("runtime_security_config.enforcement.enabled"),
		EnforcementRawSyscallEnabled: coreconfig.SystemProbe().GetBool("runtime_security_config.enforcement.raw_syscall.enabled"),
		EnforcementRuleSourceAllowed: coreconfig.SystemProbe().GetStringSlice("runtime_security_config.enforcement.rule_source_allowed"),

		// User Sessions
		UserSessionsCacheSize: coreconfig.SystemProbe().GetInt("runtime_security_config.user_sessions.cache_size"),

		// ebpf less
		EBPFLessEnabled: coreconfig.SystemProbe().GetBool("runtime_security_config.ebpfless.enabled"),
		EBPFLessSocket:  coreconfig.SystemProbe().GetString("runtime_security_config.ebpfless.socket"),

		// IMDS
		IMDSIPv4: parseIMDSIPv4(),
	}

	if err := rsConfig.sanitize(); err != nil {
		return nil, err
	}

	return rsConfig, nil
}

// IsRuntimeEnabled returns true if any feature is enabled. Has to be applied in config package too
func (c *RuntimeSecurityConfig) IsRuntimeEnabled() bool {
	return c.RuntimeEnabled || c.FIMEnabled
}

// parseIMDSIPv4 returns the uint32 representation of the IMDS IP set by the configuration
func parseIMDSIPv4() uint32 {
	ip := coreconfig.SystemProbe().GetString("runtime_security_config.imds_ipv4")
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0
	}
	return binary.LittleEndian.Uint32(parsedIP.To4())
}

// If RC is globally enabled, RC is enabled for CWS, unless the CWS-specific RC value is explicitly set to false
func isRemoteConfigEnabled() bool {
	// This value defaults to true
	rcEnabledInSysprobeConfig := coreconfig.SystemProbe().GetBool("runtime_security_config.remote_configuration.enabled")

	if !rcEnabledInSysprobeConfig {
		return false
	}

	if coreconfig.IsRemoteConfigEnabled(coreconfig.Datadog()) {
		return true
	}

	return false
}

// GetAnomalyDetectionMinimumStablePeriod returns the minimum stable period for a given event type
func (c *RuntimeSecurityConfig) GetAnomalyDetectionMinimumStablePeriod(eventType model.EventType) time.Duration {
	if minimumStablePeriod, found := c.AnomalyDetectionMinimumStablePeriods[eventType]; found {
		return minimumStablePeriod
	}
	return c.AnomalyDetectionDefaultMinimumStablePeriod
}

// sanitize ensures that the configuration is properly setup
func (c *RuntimeSecurityConfig) sanitize() error {
	serviceName := utils.GetTagValue("service", configUtils.GetConfiguredTags(coreconfig.Datadog(), true))
	if len(serviceName) > 0 {
		c.HostServiceName = serviceName
	}

	if c.IMDSIPv4 == 0 {
		return fmt.Errorf("invalid IPv4 address: got %v", coreconfig.SystemProbe().GetString("runtime_security_config.imds_ipv4"))
	}

	return c.sanitizeRuntimeSecurityConfigActivityDump()
}

// sanitizeRuntimeSecurityConfigActivityDump ensures that runtime_security_config.activity_dump is properly configured
func (c *RuntimeSecurityConfig) sanitizeRuntimeSecurityConfigActivityDump() error {
	var execFound bool
	for _, evtType := range c.ActivityDumpTracedEventTypes {
		if evtType == model.ExecEventType {
			execFound = true
			break
		}
	}
	if !execFound {
		c.ActivityDumpTracedEventTypes = append(c.ActivityDumpTracedEventTypes, model.ExecEventType)
	}

	if formats := coreconfig.SystemProbe().GetStringSlice("runtime_security_config.activity_dump.local_storage.formats"); len(formats) > 0 {
		var err error
		c.ActivityDumpLocalStorageFormats, err = ParseStorageFormats(formats)
		if err != nil {
			return fmt.Errorf("invalid value for runtime_security_config.activity_dump.local_storage.formats: %w", err)
		}
	}

	if c.ActivityDumpTracedCgroupsCount > model.MaxTracedCgroupsCount {
		c.ActivityDumpTracedCgroupsCount = model.MaxTracedCgroupsCount
	}

	if c.SecurityProfileEnabled && c.ActivityDumpEnabled && c.ActivityDumpLocalStorageDirectory != c.SecurityProfileDir {
		return fmt.Errorf("activity dumps storage directory '%s' has to be the same than security profile storage directory '%s'", c.ActivityDumpLocalStorageDirectory, c.SecurityProfileDir)
	}

	return nil
}

// ActivityDumpRemoteStorageEndpoints returns the list of activity dump remote storage endpoints parsed from the agent config
func ActivityDumpRemoteStorageEndpoints(endpointPrefix string, intakeTrackType logsconfig.IntakeTrackType, intakeProtocol logsconfig.IntakeProtocol, intakeOrigin logsconfig.IntakeOrigin) (*logsconfig.Endpoints, error) {
	logsConfig := logsconfig.NewLogsConfigKeys("runtime_security_config.activity_dump.remote_storage.endpoints.", coreconfig.Datadog())
	endpoints, err := logsconfig.BuildHTTPEndpointsWithConfig(coreconfig.Datadog(), logsConfig, endpointPrefix, intakeTrackType, intakeProtocol, intakeOrigin)
	if err != nil {
		endpoints, err = logsconfig.BuildHTTPEndpoints(coreconfig.Datadog(), intakeTrackType, intakeProtocol, intakeOrigin)
		if err == nil {
			httpConnectivity := logshttp.CheckConnectivity(endpoints.Main, coreconfig.Datadog())
			endpoints, err = logsconfig.BuildEndpoints(coreconfig.Datadog(), httpConnectivity, intakeTrackType, intakeProtocol, intakeOrigin)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("invalid endpoints: %w", err)
	}

	for _, status := range endpoints.GetStatus() {
		seclog.Infof("activity dump remote storage endpoint: %v\n", status)
	}
	return endpoints, nil
}

// ParseEvalEventType convert a eval.EventType (string) to its uint64 representation
// the current algorithm is not efficient but allows us to reduce the number of conversion functions
func ParseEvalEventType(eventType eval.EventType) model.EventType {
	for i := uint64(0); i != uint64(model.MaxAllEventType); i++ {
		if model.EventType(i).String() == eventType {
			return model.EventType(i)
		}
	}

	return model.UnknownEventType
}

// parseEventTypeDurations converts a map of durations indexed by event types
func parseEventTypeDurations(cfg coreconfig.Config, prefix string) map[model.EventType]time.Duration {
	eventTypeMap := cfg.GetStringMap(prefix)
	eventTypeDurations := make(map[model.EventType]time.Duration, len(eventTypeMap))
	for eventType := range eventTypeMap {
		eventTypeDurations[ParseEvalEventType(eventType)] = cfg.GetDuration(prefix + "." + eventType)
	}
	return eventTypeDurations
}

// parseHashAlgorithmStringSlice converts a string list to a list of hash algorithms
func parseHashAlgorithmStringSlice(algorithms []string) []model.HashAlgorithm {
	var output []model.HashAlgorithm
	for _, hashAlgorithm := range algorithms {
		for i := model.HashAlgorithm(0); i < model.MaxHashAlgorithm; i++ {
			if i.String() == hashAlgorithm {
				output = append(output, i)
				break
			}
		}
	}
	return output
}

// GetFamilyAddress returns the address famility to use for system-probe <-> security-agent communication
func GetFamilyAddress(path string) (string, string) {
	if strings.HasPrefix(path, "/") {
		return "unix", path
	}
	return "tcp", path
}
