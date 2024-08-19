// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package lsof

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/prometheus/procfs"

	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// see the documentation for /proc for more details about files and their format
// https://www.kernel.org/doc/html/latest/filesystems/proc.html

// openFilesLister stores the state needed to list open files
// this is helpful for tests to help with mocking
type openFilesLister struct {
	pid      int
	procPath string

	readlink func(string) (string, error)
	stat     func(string) (os.FileInfo, error)
	lstat    func(string) (os.FileInfo, error)

	proc       procfsProc
	socketInfo map[uint64]socketInfo
}

// procfsProc is an interface to allow mocking of procfs.Proc
type procfsProc interface {
	ProcMaps() ([]*procfs.ProcMap, error)
	FileDescriptors() ([]uintptr, error)
	Cwd() (string, error)
}

type socketInfo struct {
	Description string
	State       string
	Protocol    string
}

func openFiles(_ context.Context, pid int) (Files, error) {
	ofl := &openFilesLister{
		pid: pid,

		readlink: os.Readlink,
		stat:     os.Stat,
		lstat:    os.Lstat,
	}

	ofl.procPath = procPath()

	fs, err := procfs.NewFS(ofl.procPath)
	if err != nil {
		return nil, err
	}

	ofl.proc, err = fs.Proc(pid)
	if err != nil {
		return nil, err
	}

	ofl.socketInfo = readSocketInfo(ofl.procPIDPath())

	return ofl.openFiles()
}

func (ofl *openFilesLister) procPIDPath() string {
	return fmt.Sprintf("%s/%d", ofl.procPath, ofl.pid)
}

func (ofl *openFilesLister) openFiles() (Files, error) {
	var files Files

	// open files, socket, pipe (everything with a file descriptor, from /proc/<pid>/fd)
	openFDFiles, err := ofl.fdMetadata()
	if err != nil {
		log.Debugf("Failed to get open FDs for pid %d: %s", ofl.pid, err)
	} else {
		files = append(files, openFDFiles...)
	}

	// memory mapped files, code, regions (from /proc/<pid>/maps)
	mmapFiles, err := ofl.mmapMetadata()
	if err != nil {
		log.Debugf("Failed to get memory maps for pid %d: %s", ofl.pid, err)
	} else {
		files = append(files, mmapFiles...)
	}

	return files, nil
}

func (ofl *openFilesLister) mmapMetadata() (Files, error) {
	maps, err := ofl.proc.ProcMaps()
	if err != nil {
		return nil, err
	}

	cwd, _ := ofl.proc.Cwd()

	var files Files
	for i, m := range maps {
		if i > 0 && m.Pathname == maps[i-1].Pathname {
			// skip duplicate entries
			continue
		}
		if m.Pathname == "" {
			// anonymous mapping
			continue
		}

		if m.Dev == 0 || m.Inode == 0 {
			// virtual memory region, eg. [heap], [stack], [vvar], etc
			continue
		}

		file := File{
			Name:     m.Pathname,
			OpenPerm: permToString(m.Perms),
		}

		if file.Type, file.FilePerm, file.Size, _ = fileStats(ofl.stat, m.Pathname); file.Type == "" {
			continue
		}
		file.Fd = mmapFD(m.Pathname, file.Type, cwd)

		files = append(files, file)
	}

	return files, nil
}

func permToString(perms *procfs.ProcMapPermissions) string {
	s := ""

	for _, perm := range []struct {
		set       bool
		charSet   string
		charUnset string
	}{
		{perms.Read, "r", "-"},
		{perms.Write, "w", "-"},
		{perms.Execute, "x", "-"},
		{perms.Private, "p", ""},
		{perms.Shared, "s", ""},
	} {
		if perm.set {
			s += perm.charSet
		} else {
			s += perm.charUnset
		}
	}

	return s
}

func mmapFD(path string, ty, cwd string) string {
	/*
		cwd  current working directory;
		ltx  shared library text (code and data);
		mem  memory-mapped file;
		mmap memory-mapped device;
		rtd  root directory;
		txt  program text (code and data);
	*/
	fd := "unknown"
	if ty == "REG" {
		// knowing whether the file is memory mapped or program text would require knowing permissions
		// which aren't parsed by gopsutil, so we just assume it's memory mapped
		fd = "mem"
	} else if ty == "DIR" {
		if path == "/" {
			fd = "rtd"
		} else if cwd != "" && path == cwd {
			fd = "cwd"
		}
	}
	return fd
}

func (ofl *openFilesLister) fdMetadata() (Files, error) {
	openFDs, err := ofl.proc.FileDescriptors()
	if err != nil {
		return nil, err
	}

	var files Files
	for _, openFD := range openFDs {
		file, ok := ofl.fdStat(openFD)
		if ok {
			files = append(files, file)
		}
	}

	return files, nil
}

func (ofl *openFilesLister) fdStat(fd uintptr) (File, bool) {
	var err error
	file := File{
		Fd: fmt.Sprintf("%d", fd),
	}

	fdLinkPath := fmt.Sprintf("%s/fd/%d", ofl.procPIDPath(), fd)

	var inode uint64
	if file.Type, file.OpenPerm, _, _ = fileStats(ofl.lstat, fdLinkPath); file.Type == "" {
		return File{}, false
	}

	if file.Type, file.FilePerm, file.Size, inode = fileStats(ofl.stat, fdLinkPath); file.Type == "" {
		return File{}, false
	}

	if info, ok := ofl.socketInfo[inode]; ok {
		file.Name = info.Description
		file.FilePerm = info.State
		file.Type = info.Protocol
	} else {
		file.Name, err = ofl.readlink(fdLinkPath)
		if err != nil {
			return File{}, false
		}
	}

	return file, true
}

// readSocketInfo reads the socket information from /proc/<pid>/net/{tcp,tcp6,udp,udp6,unix}
// returns a map of inode to socketInfo
// see https://www.kernel.org/doc/Documentation/networking/proc_net_tcp.txt
func readSocketInfo(procPIDPath string) map[uint64]socketInfo {
	si := make(map[uint64]socketInfo)

	fs, err := procfs.NewFS(procPIDPath)
	if err != nil {
		log.Debugf("Failed to read %s: %s", procPIDPath, err)
		return si
	}

	for protocol, parser := range map[string]func() (procfs.NetTCP, error){
		"tcp":  fs.NetTCP,
		"tcp6": fs.NetTCP6,
	} {
		addrs, err := parser()
		if err != nil {
			log.Debugf("Failed to read %s socket info in %s: %s", protocol, procPIDPath, err)
			continue
		}
		for _, entry := range addrs {
			si[entry.Inode] = socketInfo{
				fmt.Sprintf("%s:%d->%s:%d", entry.LocalAddr, entry.LocalPort, entry.RemAddr, entry.RemPort),
				// see https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/net/tcp_states.h
				// for the values of states
				fmt.Sprintf("%d", entry.St),
				protocol,
			}
		}
	}

	for protocol, parser := range map[string]func() (procfs.NetUDP, error){
		"udp":  fs.NetUDP,
		"udp6": fs.NetUDP6,
	} {
		addrs, err := parser()
		if err != nil {
			log.Debugf("Failed to read %s socket info in %s: %s", protocol, procPIDPath, err)
			continue
		}
		for _, entry := range addrs {
			si[entry.Inode] = socketInfo{
				fmt.Sprintf("%s:%d->%s:%d", entry.LocalAddr, entry.LocalPort, entry.RemAddr, entry.RemPort),
				fmt.Sprintf("%d", entry.St),
				protocol,
			}
		}
	}

	unix, err := fs.NetUNIX()
	if err == nil {
		for _, entry := range unix.Rows {
			si[entry.Inode] = socketInfo{
				fmt.Sprintf("%s:%s", entry.Type, entry.Path),
				fmt.Sprintf("%s:%s", entry.State, entry.Flags),
				"unix",
			}
		}
	} else {
		log.Debugf("Failed to read unix socket info in %s: %s", procPIDPath, err)
	}

	return si
}

// fileStats returns the type, permission, size, and inode of a file
func fileStats(statf func(string) (os.FileInfo, error), path string) (string, string, int64, uint64) {
	stat, err := statf(path)
	if err != nil {
		log.Debugf("stat failed for %s: %v", path, err)
		return "", "", 0, 0
	}

	ty := "REG"
	if stat.Mode()&os.ModeSocket != 0 {
		ty = "SOCKET"
	} else if stat.Mode()&os.ModeNamedPipe != 0 {
		ty = "PIPE"
	} else if stat.Mode()&os.ModeDevice != 0 {
		ty = "DEV"
	} else if stat.Mode()&os.ModeDir != 0 {
		ty = "DIR"
	} else if stat.Mode()&os.ModeCharDevice != 0 {
		ty = "CHAR"
	} else if stat.Mode()&os.ModeSymlink != 0 {
		ty = "LINK"
	} else if stat.Mode()&os.ModeIrregular != 0 {
		ty = "?"
	}

	size := stat.Size()
	perm := stat.Mode().Perm().String()

	var ino uint64
	if sys, ok := stat.Sys().(*syscall.Stat_t); ok {
		ino = sys.Ino
	}

	return ty, perm, size, ino
}

func procPath() string {
	if procPath, ok := os.LookupEnv("HOST_PROC"); ok {
		return procPath
	}
	return "/proc"
}
