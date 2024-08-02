// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build windows

package network

import (
	"net/netip"
	"syscall"

	"github.com/DataDog/datadog-agent/pkg/network/driver"
	"github.com/DataDog/datadog-agent/pkg/network/protocols"
	"github.com/DataDog/datadog-agent/pkg/process/util"
)

func connFamily(addressFamily uint16) ConnectionFamily {
	if addressFamily == syscall.AF_INET {
		return AFINET
	}
	return AFINET6
}

func connType(protocol uint16) ConnectionType {
	if protocol == syscall.IPPROTO_TCP {
		return TCP
	}
	return UDP
}

func connDirection(flags uint32) ConnectionDirection {
	direction := (flags & driver.FlowDirectionMask) >> driver.FlowDirectionBits
	if (direction & driver.FlowDirectionInbound) == driver.FlowDirectionInbound {
		return INCOMING
	}
	if (direction & driver.FlowDirectionOutbound) == driver.FlowDirectionOutbound {
		return OUTGOING
	}
	return OUTGOING
}

func isFlowClosed(flags uint32) bool {
	// Connection is closed
	return (flags & driver.FlowClosedMask) == driver.FlowClosedMask
}

func isTCPFlowEstablished(flow *driver.PerFlowData) bool {
	//return (flags & driver.TCPFlowEstablishedMask) == driver.TCPFlowEstablishedMask
	tcpdata := flow.TCPFlow()
	if nil != tcpdata {
		if tcpdata.ConnectionStatus == driver.TcpStatusEstablished {
			return true
		}
	}
	return false
}

func convertV4Addr(addr [16]uint8) util.Address {
	// We only read the first 4 bytes for v4 address
	return util.Address{Addr: netip.AddrFrom4([4]byte(addr))}
}

func convertV6Addr(addr [16]uint8) util.Address {
	// We read all 16 bytes for v6 address
	return util.Address{Addr: netip.AddrFrom16(addr)}
}

// Monotonic values include retransmits and headers, while transport does not. We default to using transport
// values and must explicitly enable using monotonic counts in the config. This is consistent with the Linux probe
func monotonicOrTransportBytes(useMonotonicCounts bool, monotonic uint64, transport uint64) uint64 {
	if useMonotonicCounts {
		return monotonic
	}
	return transport
}

// FlowToConnStat converts a driver.PerFlowData into a ConnectionStats struct for use with the tracer
func FlowToConnStat(cs *ConnectionStats, flow *driver.PerFlowData, enableMonotonicCounts bool) {
	var (
		family         ConnectionFamily
		srcAddr        util.Address
		dstAddr        util.Address
		connectionType ConnectionType
	)
	family = connFamily(flow.AddressFamily)
	connectionType = connType(flow.Protocol)

	// V4 Address
	if family == AFINET {
		srcAddr, dstAddr = convertV4Addr(flow.LocalAddress), convertV4Addr(flow.RemoteAddress)
	} else {
		// V6 Address
		srcAddr, dstAddr = convertV6Addr(flow.LocalAddress), convertV6Addr(flow.RemoteAddress)
	}

	*cs = ConnectionStats{}
	cs.Source = srcAddr
	cs.Dest = dstAddr
	// after lengthy discussion, use the transport bytes in/out.  monotonic
	// RecvBytes/SentBytes includes the size of the IP header and transport
	// header, transportBytes is the raw transport data.  At present,
	// the linux probe only reports the raw transport data.  So do that by default.
	cs.Monotonic.SentBytes = monotonicOrTransportBytes(enableMonotonicCounts, flow.MonotonicSentBytes, flow.TransportBytesOut)
	cs.Monotonic.RecvBytes = monotonicOrTransportBytes(enableMonotonicCounts, flow.MonotonicRecvBytes, flow.TransportBytesIn)
	cs.Monotonic.SentPackets = flow.PacketsOut
	cs.Monotonic.RecvPackets = flow.PacketsIn
	cs.LastUpdateEpoch = flow.Timestamp
	cs.Pid = uint32(flow.ProcessId)
	cs.SPort = flow.LocalPort
	cs.DPort = flow.RemotePort
	cs.Type = connectionType
	cs.Family = family
	cs.Direction = connDirection(flow.Flags)
	cs.SPortIsEphemeral = IsPortInEphemeralRange(cs.Family, cs.Type, cs.SPort)
	cs.Cookie = flow.FlowCookie
	if connectionType == TCP {

		tf := flow.TCPFlow()
		if tf != nil {
			cs.TCPFailures = make(map[uint32]uint32)
			cs.Monotonic.Retransmits = uint32(tf.RetransmitCount)
			cs.RTT = uint32(tf.SRTT)
			cs.RTTVar = uint32(tf.RttVariance)

			switch driver.ConnectionStatus(tf.ConnectionStatus) {
			case driver.ConnectionStatusACKRST:
				cs.TCPFailures[111] = 1 // ECONNREFUSED in posix is 111
			case driver.ConnectionStatusTimeout:
				cs.TCPFailures[110] = 1 // ETIMEDOUT in posix is 110
			case driver.ConnectionStatusSentRst:
				fallthrough
			case driver.ConnectionStatusRecvRst:
				cs.TCPFailures[104] = 1 // ECONNRESET in posix is 104
			}

		}

		if isTCPFlowEstablished(flow) {
			cs.Monotonic.TCPEstablished = 1
		}
		if isFlowClosed(flow.Flags) {
			cs.Monotonic.TCPClosed = 1
		}
		if flow.ClassificationStatus == driver.ClassificationClassified {
			switch crq := flow.ClassifyRequest; {
			default:
				// this is unexpected.  The various case statements should
				// encompass all of the available values.

			case crq == driver.ClassificationRequestUnclassified:
				// do nothing because it may be classified in the response if
				// the request portion of the flow was missed.

			case crq >= driver.ClassificationRequestHTTPUnknown && crq < driver.ClassificationRequestHTTPLast:
				cs.ProtocolStack = protocols.Stack{Application: protocols.HTTP}
			case crq == driver.ClassificationRequestHTTP2:
				cs.ProtocolStack = protocols.Stack{Application: protocols.HTTP2}
			case crq == driver.ClassificationRequestTLS:
				cs.ProtocolStack = protocols.Stack{Encryption: protocols.TLS}
			}

			switch crsp := flow.ClassifyResponse; {
			default:
				// this is unexpected.  The various case statements should
				// encompass all of the available values.

			case crsp == driver.ClassificationRequestUnclassified:
				// do nothing because it will have been classified in the request

			case crsp == driver.ClassificationResponseHTTP:
				if flow.HttpUpgradeToH2Accepted == 1 {
					cs.ProtocolStack = protocols.Stack{Application: protocols.HTTP2}
				} else {
					// could have missed the request.  Most likely this is just
					// resetting the existing value
					cs.ProtocolStack = protocols.Stack{Application: protocols.HTTP}
				}
			}
		} else {
			// one of
			// ClassificationUnableInsufficientData, ClassificationUnknown, ClassificationUnclassified
			cs.ProtocolStack = protocols.Stack{}
		}
	}
}
