// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux
// +build linux

package constantfetch

import (
	"runtime"

	"github.com/DataDog/datadog-agent/pkg/security/ebpf/kernel"
)

// FallbackConstantFetcher is a constant fetcher that uses the old fallback
// heuristics to fetch constants
type FallbackConstantFetcher struct {
	kernelVersion *kernel.Version
	res           map[string]uint64
}

// NewFallbackConstantFetcher returns a new FallbackConstantFetcher
func NewFallbackConstantFetcher(kv *kernel.Version) *FallbackConstantFetcher {
	return &FallbackConstantFetcher{
		kernelVersion: kv,
		res:           make(map[string]uint64),
	}
}

func (f *FallbackConstantFetcher) String() string {
	return "fallback"
}

func (f *FallbackConstantFetcher) appendRequest(id string) {
	var value = ErrorSentinel
	switch id {
	case "sizeof_inode":
		value = getSizeOfStructInode(f.kernelVersion)
	case "sb_magic_offset":
		value = getSuperBlockMagicOffset(f.kernelVersion)
	case "tty_offset":
		value = getSignalTTYOffset(f.kernelVersion)
	case "tty_name_offset":
		value = getTTYNameOffset(f.kernelVersion)
	case "creds_uid_offset":
		value = getCredsUIDOffset(f.kernelVersion)
	case "bpf_map_id_offset":
		value = getBpfMapIDOffset(f.kernelVersion)
	case "bpf_map_name_offset":
		value = getBpfMapNameOffset(f.kernelVersion)
	case "bpf_map_type_offset":
		value = getBpfMapTypeOffset(f.kernelVersion)
	case "bpf_prog_aux_offset":
		value = getBpfProgAuxOffset(f.kernelVersion)
	case "bpf_prog_tag_offset":
		value = getBpfProgTagOffset(f.kernelVersion)
	case "bpf_prog_type_offset":
		value = getBpfProgTypeOffset(f.kernelVersion)
	case "bpf_prog_attach_type_offset":
		value = getBpfProgAttachTypeOffset(f.kernelVersion)
	case "bpf_prog_aux_id_offset":
		value = getBpfProgAuxIDOffset(f.kernelVersion)
	case "bpf_prog_aux_name_offset":
		value = getBpfProgAuxNameOffset(f.kernelVersion)
	case "pid_level_offset":
		value = getPIDLevelOffset(f.kernelVersion)
	case "pid_numbers_offset":
		value = getPIDNumbersOffset(f.kernelVersion)
	case "sizeof_upid":
		value = getSizeOfUpid(f.kernelVersion)
	case "dentry_sb_offset":
		value = getDentrySuperBlockOffset(f.kernelVersion)
	case "pipe_inode_info_bufs_offset":
		value = getPipeInodeInfoBufsOffset(f.kernelVersion)
	case "net_device_ifindex_offset":
		value = getNetDeviceIfindexOffset(f.kernelVersion)
	case "net_ns_offset":
		value = getNetNSOffset(f.kernelVersion)
	case "net_proc_inum_offset":
		value = getNetProcINumOffset(f.kernelVersion)
	case "sock_common_skc_net_offset":
		value = getSockCommonSKCNetOffset(f.kernelVersion)
	case "socket_sock_offset":
		value = getSocketSockOffset(f.kernelVersion)
	case "nf_conn_ct_net_offset":
		value = getNFConnCTNetOffset(f.kernelVersion)
	case "sock_common_skc_family_offset":
		value = getSockCommonSKCFamilyOffset(f.kernelVersion)
	case "flowi4_saddr_offset":
		value = getFlowi4SAddrOffset(f.kernelVersion)
	case "flowi6_saddr_offset":
		value = getFlowi6SAddrOffset(f.kernelVersion)
	case "flowi4_uli_offset":
		value = getFlowi4ULIOffset(f.kernelVersion)
	case "flowi6_uli_offset":
		value = getFlowi6ULIOffset(f.kernelVersion)
	}
	f.res[id] = value
}

// AppendSizeofRequest appends a sizeof request
func (f *FallbackConstantFetcher) AppendSizeofRequest(id, typeName, headerName string) {
	f.appendRequest(id)
}

// AppendOffsetofRequest appends an offset request
func (f *FallbackConstantFetcher) AppendOffsetofRequest(id, typeName, fieldName, headerName string) {
	f.appendRequest(id)
}

// FinishAndGetResults returns the results
func (f *FallbackConstantFetcher) FinishAndGetResults() (map[string]uint64, error) {
	return f.res, nil
}

func getSizeOfStructInode(kv *kernel.Version) uint64 {
	sizeOf := uint64(600)

	switch {
	case kv.IsRH7Kernel():
		sizeOf = 584
	case kv.IsRH8Kernel():
		sizeOf = 648
	case kv.IsSLES12Kernel():
		sizeOf = 560
	case kv.IsSLES15Kernel():
		sizeOf = 592
	case kv.IsOracleUEKKernel():
		sizeOf = 632
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel4_20):
		sizeOf = 712
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		sizeOf = 704
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		sizeOf = 704
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		sizeOf = 584
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		sizeOf = 584
	case kv.Code != 0 && kv.Code < kernel.Kernel4_16:
		sizeOf = 608
	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_1):
		sizeOf = 584
	case kv.IsInRangeCloseOpen(kernel.Kernel5_13, kernel.Kernel5_15):
		sizeOf = 592
	case kv.Code >= kernel.Kernel5_15:
		sizeOf = 632
	}

	return sizeOf
}

func getSuperBlockMagicOffset(kv *kernel.Version) uint64 {
	sizeOf := uint64(96)

	if kv.IsRH7Kernel() {
		sizeOf = 88
	}

	return sizeOf
}

func getSignalTTYOffset(kv *kernel.Version) uint64 {
	ttyOffset := uint64(400)

	switch {
	case kv.IsRH7Kernel():
		ttyOffset = 416
	case kv.IsRH8Kernel():
		ttyOffset = 392
	case kv.IsSLES12Kernel():
		ttyOffset = 376
	case kv.IsSLES15Kernel():
		ttyOffset = 408
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel4_20):
		ttyOffset = 416
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		ttyOffset = 416
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		ttyOffset = 416
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_14, kernel.Kernel4_15):
		ttyOffset = 368
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		ttyOffset = 400
	case kv.IsInRangeCloseOpen(kernel.Kernel4_13, kernel.Kernel4_19):
		ttyOffset = 376
	case kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel5_0):
		ttyOffset = 400
	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_1):
		ttyOffset = 408
	case kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		ttyOffset = 408
	case kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_7):
		if runtime.GOARCH == "arm64" || kv.IsOracleUEKKernel() {
			ttyOffset = 408
		} else {
			ttyOffset = 400
		}
	case kv.Code != 0 && kv.Code == kernel.Kernel5_10:
		if runtime.GOARCH == "arm64" {
			ttyOffset = 408
		} else {
			ttyOffset = 400
		}
	case kv.IsInRangeCloseOpen(kernel.Kernel5_7, kernel.Kernel5_9) || kv.IsInRangeCloseOpen(kernel.Kernel5_11, kernel.Kernel5_14):
		if runtime.GOARCH == "arm64" {
			ttyOffset = 400
		} else {
			ttyOffset = 408
		}
	case kv.Code != 0 && kv.Code < kernel.Kernel5_3:
		ttyOffset = 368
	case kv.Code >= kernel.Kernel5_16:
		ttyOffset = 416
	}

	return ttyOffset
}

func getTTYNameOffset(kv *kernel.Version) uint64 {
	nameOffset := uint64(368)

	switch {
	case kv.IsRH7Kernel():
		nameOffset = 312
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel4_20):
		nameOffset = 552
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		nameOffset = 552
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		nameOffset = 544
	case kv.IsInRangeCloseOpen(kernel.Kernel4_13, kernel.Kernel5_8):
		nameOffset = 368
	case kv.IsInRangeCloseOpen(kernel.Kernel5_8, kernel.Kernel5_9) && runtime.GOARCH == "arm64":
		nameOffset = 368
	case kv.IsInRangeCloseOpen(kernel.Kernel5_8, kernel.Kernel5_14):
		nameOffset = 360
	case kv.Code >= kernel.Kernel5_14:
		nameOffset = 352
	}

	return nameOffset
}

func getCredsUIDOffset(kv *kernel.Version) uint64 {
	size := uint64(4)

	switch {
	case kv.IsCOSKernel():
		size += 16
	}

	return size
}

func getBpfMapIDOffset(kv *kernel.Version) uint64 {
	switch {
	case kv.IsInRangeCloseOpen(kernel.Kernel5_15, kernel.Kernel5_16):
		return 52
	case kv.Code >= kernel.Kernel5_16:
		return 60
	default:
		return 48
	}
}

func getBpfMapNameOffset(kv *kernel.Version) uint64 {
	nameOffset := uint64(168)

	switch {
	case kv.IsRH7Kernel():
		nameOffset = 112
	case kv.IsRH8Kernel():
		nameOffset = 80
	case kv.IsSLES15Kernel():
		nameOffset = 88
	case kv.IsSLES12Kernel():
		nameOffset = 176

	case kv.IsInRangeCloseOpen(kernel.Kernel4_18, kernel.Kernel5_1):
		nameOffset = 176
	case kv.IsInRangeCloseOpen(kernel.Kernel5_1, kernel.Kernel5_3):
		nameOffset = 200
	case kv.IsInRangeCloseOpen(kernel.Kernel5_3, kernel.Kernel5_5):
		if kv.IsOracleUEKKernel() {
			nameOffset = 200
		} else {
			nameOffset = 168
		}
	case kv.IsInRangeCloseOpen(kernel.Kernel5_5, kernel.Kernel5_11):
		nameOffset = 88
	case kv.IsInRangeCloseOpen(kernel.Kernel5_11, kernel.Kernel5_13):
		nameOffset = 80
	case kv.IsInRangeCloseOpen(kernel.Kernel5_13, kernel.Kernel5_15):
		nameOffset = 80
	case kv.IsInRangeCloseOpen(kernel.Kernel5_15, kernel.Kernel5_16):
		nameOffset = 88
	case kv.Code >= kernel.Kernel5_16:
		nameOffset = 96
	case kv.Code != 0 && kv.Code < kernel.Kernel4_15:
		return ErrorSentinel
	}

	return nameOffset
}

func getBpfMapTypeOffset(kv *kernel.Version) uint64 {
	return uint64(24)
}

func getBpfProgAuxOffset(kv *kernel.Version) uint64 {
	auxOffset := uint64(32)

	switch {
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_14, kernel.Kernel4_15):
		auxOffset = 24
	case kv.Code >= kernel.Kernel5_13:
		auxOffset = 56
	}

	return auxOffset
}

func getBpfProgTagOffset(kv *kernel.Version) uint64 {
	progTagOffset := uint64(20)
	switch {
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_14, kernel.Kernel4_15):
		progTagOffset = 16
	default:
	}

	return progTagOffset
}

func getBpfProgTypeOffset(kv *kernel.Version) uint64 {
	return uint64(4)
}

func getBpfProgAttachTypeOffset(kv *kernel.Version) uint64 {
	return uint64(8)
}

func getBpfProgAuxIDOffset(kv *kernel.Version) uint64 {
	idOffset := uint64(24)

	switch {
	case kv.IsRH7Kernel():
		idOffset = 8
	case kv.IsRH8Kernel():
		idOffset = 32
	case kv.IsSLES15Kernel():
		idOffset = 28
	case kv.IsSLES12Kernel():
		idOffset = 16
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_14, kernel.Kernel4_15):
		idOffset = 16

	case kv.IsInRangeCloseOpen(kernel.Kernel4_18, kernel.Kernel5_0):
		idOffset = 16
	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_4):
		idOffset = 20
	case kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_8):
		idOffset = 24
	case kv.IsInRangeCloseOpen(kernel.Kernel5_8, kernel.Kernel5_13):
		idOffset = 28
	case kv.Code != 0 && kv.Code >= kernel.Kernel5_13:
		idOffset = 32
	}

	return idOffset
}

func getBpfProgAuxNameOffset(kv *kernel.Version) uint64 {
	nameOffset := uint64(176)

	switch {
	case kv.IsRH7Kernel():
		nameOffset = 144
	case kv.IsRH8Kernel():
		nameOffset = 520
	case kv.IsSLES15Kernel():
		nameOffset = 256
	case kv.IsSLES12Kernel():
		nameOffset = 160
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		nameOffset = 544

	case kv.IsInRangeCloseOpen(kernel.Kernel4_18, kernel.Kernel4_19):
		nameOffset = 152
	case kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel5_0):
		nameOffset = 160
	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_8):
		nameOffset = 176
	case kv.IsInRangeCloseOpen(kernel.Kernel5_8, kernel.Kernel5_10):
		nameOffset = 416
	case kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		nameOffset = 496
	case kv.IsInRangeCloseOpen(kernel.Kernel5_11, kernel.Kernel5_13):
		nameOffset = 504
	case kv.IsInRangeCloseOpen(kernel.Kernel5_13, kernel.Kernel5_16):
		nameOffset = 528
	case kv.Code != 0 && kv.Code >= kernel.Kernel5_16:
		nameOffset = 544
	}

	return nameOffset
}

func getPIDLevelOffset(kv *kernel.Version) uint64 {
	return uint64(4)
}

func getPIDNumbersOffset(kv *kernel.Version) uint64 {
	pidNumbersOffset := uint64(48)

	switch {
	case kv.IsRH7Kernel():
		pidNumbersOffset = 48
	case kv.IsRH8Kernel():
		pidNumbersOffset = 56
	case kv.IsSLES12Kernel():
		pidNumbersOffset = 48
	case kv.IsSLES15Kernel():
		pidNumbersOffset = 80
	case kv.IsDebianKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel4_20):
		pidNumbersOffset = 56
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel4_19, kernel.Kernel4_20):
		pidNumbersOffset = 56
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_4, kernel.Kernel5_5):
		pidNumbersOffset = 96
	case kv.IsCOSKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		pidNumbersOffset = 128

	case kv.IsInRangeCloseOpen(kernel.Kernel4_15, kernel.Kernel5_0):
		pidNumbersOffset = 48
	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_3):
		pidNumbersOffset = 56
	case kv.IsInRangeCloseOpen(kernel.Kernel5_3, kernel.Kernel5_7):
		pidNumbersOffset = 80
	case kv.Code != 0 && kv.Code >= kernel.Kernel5_7:
		pidNumbersOffset = 96
	}
	return pidNumbersOffset
}

func getSizeOfUpid(kv *kernel.Version) uint64 {
	sizeOfUpid := uint64(16)

	switch {
	case kv.IsRH7Kernel():
		sizeOfUpid = 32
	case kv.IsRH8Kernel():
		sizeOfUpid = 16
	case kv.IsSLES12Kernel():
		sizeOfUpid = 16
	case kv.IsSLES15Kernel():
		sizeOfUpid = 32
	case kv.IsAmazonLinuxKernel() && kv.Code != 0 && kv.Code < kernel.Kernel4_15:
		sizeOfUpid = 32
	}
	return sizeOfUpid
}

func getDentrySuperBlockOffset(kv *kernel.Version) uint64 {
	offset := uint64(104)

	switch {
	case kv.IsCOSKernel():
		offset = 128
	}

	return offset
}

func getPipeInodeInfoBufsOffset(kv *kernel.Version) uint64 {
	offset := uint64(120)

	switch {
	case kv.IsRH7Kernel():
		offset = 128
	case kv.IsRH8Kernel():
		offset = 120
	case kv.IsAmazonLinuxKernel() && kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		offset = 152

	case kv.IsInRangeCloseOpen(kernel.Kernel4_13, kernel.Kernel5_6):
		offset = 120
	case kv.IsInRangeCloseOpen(kernel.Kernel5_6, kernel.Kernel5_8) ||
		kv.IsInRangeCloseOpen(kernel.Kernel5_10, kernel.Kernel5_11):
		offset = 144
	case kv.Code != 0 && kv.Code >= kernel.Kernel5_8:
		offset = 152
	}

	return offset
}

func getNetDeviceIfindexOffset(kv *kernel.Version) uint64 {
	offset := uint64(260)

	switch {
	case kv.IsRH7Kernel():
		offset = 192
	case kv.IsRH8Kernel():
		offset = 264
	case kv.IsSLES12Kernel():
		offset = 264
	case kv.IsSLES15Kernel():
		offset = 256

	case kv.Code >= kernel.Kernel4_14 && kv.Code < kernel.Kernel5_8:
		offset = 264
	case kv.Code >= kernel.Kernel5_8 && kv.Code < kernel.Kernel5_12:
		offset = 256
	case kv.Code >= kernel.Kernel5_12:
		offset = 208
	}

	return offset
}

func getNetNSOffset(kv *kernel.Version) uint64 {
	switch {
	// Commit 355b98553789b646ed97ad801a619ff898471b92 introduces a hashmix field for security
	// purposes. This commit was cherry-picked in stable releases 4.9.168, 4.14.111, 4.19.34 and 5.0.7
	// and is part of master since 5.1
	case kv.IsRH8Kernel():
		fallthrough
	case (kv.IsInRangeCloseOpen(kernel.Kernel4_9, kernel.Kernel4_10) && kv.Code.Patch() >= 168) ||
		(kv.IsInRangeCloseOpen(kernel.Kernel4_14, kernel.Kernel4_15) && kv.Code.Patch() >= 111) ||
		kv.Code >= kernel.Kernel5_1:
		return 120
	default:
		return 112
	}
}

func getNetProcINumOffset(kv *kernel.Version) uint64 {
	return uint64(72)
}

func getSockCommonSKCNetOffset(kv *kernel.Version) uint64 {
	return uint64(48)
}

func getSocketSockOffset(kv *kernel.Version) uint64 {
	offset := uint64(32)

	switch {
	case kv.IsRH7Kernel():
		offset = 32
	case kv.IsRH8Kernel():
		offset = 32
	case kv.IsSLES12Kernel():
		offset = 32
	case kv.IsSLES15Kernel():
		offset = 24

	case kv.Code >= kernel.Kernel5_3:
		offset = 24
	}

	return offset
}

func getNFConnCTNetOffset(kv *kernel.Version) uint64 {
	offset := uint64(144)

	switch {
	case kv.IsRH7Kernel():
		offset = 240
	}

	return offset
}

func getSockCommonSKCFamilyOffset(kv *kernel.Version) uint64 {
	return 16
}

func getFlowi4SAddrOffset(kv *kernel.Version) uint64 {
	offset := uint64(40)

	switch {
	case kv.IsRH7Kernel():
		offset = 20
	case kv.IsRH8Kernel():
		offset = 56

	case kv.IsInRangeCloseOpen(kernel.Kernel5_0, kernel.Kernel5_1):
		offset = 32
	case kv.Code >= kernel.Kernel5_1:
		offset = 40
	}

	return offset
}

func getFlowi4ULIOffset(kv *kernel.Version) uint64 {
	return getFlowi4SAddrOffset(kv) + 8
}

func getFlowi6SAddrOffset(kv *kernel.Version) uint64 {
	return getFlowi4ULIOffset(kv) + 8
}

func getFlowi6ULIOffset(kv *kernel.Version) uint64 {
	return getFlowi6SAddrOffset(kv) + 20
}
