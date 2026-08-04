package ebpf

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type NetEventType int

const (
	NetConnect NetEventType = 0 + iota
	NetAccept
	NetSend
	NetRecv
	NetClose
	NetDNS
	NetBind
	NetListen
)

func (t NetEventType) String() string {
	switch t {
	case NetConnect:
		return "CONNECT"
	case NetAccept:
		return "ACCEPT"
	case NetSend:
		return "SEND"
	case NetRecv:
		return "RECV"
	case NetClose:
		return "CLOSE"
	case NetDNS:
		return "DNS"
	case NetBind:
		return "BIND"
	case NetListen:
		return "LISTEN"
	default:
		return unknownLabel
	}
}

type NetEvent struct {
	PID       uint32
	TID       uint32
	UID       uint32
	GID       uint32
	Type      NetEventType
	Protocol  uint32 // IPPROTO_* (6=TCP, 17=UDP, 1=ICMP) or 0=unknown
	Size      uint32 // data size for SEND/RECV, backlog for LISTEN
	FD        uint32
	Comm      string
	SrcAddr   string // "IP:port" or empty
	DstAddr   string // "IP:port"
	NetNS     uint64
	CgroupID  uint64
	Timestamp int64
}

func ParseNetEventType(s string) (NetEventType, bool) {
	switch strings.ToUpper(s) {
	case "CONNECT":
		return NetConnect, true
	case "ACCEPT":
		return NetAccept, true
	case "SEND":
		return NetSend, true
	case "RECV":
		return NetRecv, true
	case "CLOSE":
		return NetClose, true
	case "DNS":
		return NetDNS, true
	case "BIND":
		return NetBind, true
	case "LISTEN":
		return NetListen, true
	default:
		return 0, false
	}
}

func NetEventTypes() []NetEventType {
	return []NetEventType{
		NetConnect, NetAccept, NetSend, NetRecv, NetClose, NetDNS, NetBind, NetListen,
	}
}

type NetBpfEvent struct {
	PID      uint32
	TID      uint32
	UID      uint32
	GID      uint32
	Type     uint32
	Proto    uint32
	Size     uint32
	FD       uint32
	AF       uint32
	Saddr    [4]uint32
	Daddr    [4]uint32
	Sport    uint16
	Dport    uint16
	Comm     [16]byte
	NetNS    uint64
	CgroupID uint64
}

func (e *NetBpfEvent) ToNetEvent() NetEvent {
	return NetEvent{
		PID:      e.PID,
		TID:      e.TID,
		UID:      e.UID,
		GID:      e.GID,
		Type:     NetEventType(e.Type),
		Protocol: e.Proto,
		Size:     e.Size,
		FD:       e.FD,
		Comm:     Cstr(e.Comm[:]),
		SrcAddr:  FormatAddr(e.AF, e.Saddr[:], e.Sport),
		DstAddr:  FormatAddr(e.AF, e.Daddr[:], e.Dport),
		NetNS:    e.NetNS,
		CgroupID: e.CgroupID,
	}
}

func Ntohs(v uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return binary.LittleEndian.Uint16(b)
}

func FormatAddr(af uint32, ip []uint32, port uint16) string {
	if af == 0 || len(ip) == 0 {
		return ""
	}

	var buf strings.Builder

	switch {
	case af == 2:
		addr := make(net.IP, 4)
		binary.BigEndian.PutUint32(addr, ip[0])
		buf.WriteString(addr.String())
	case af == 10 && len(ip) >= 4:
		addr := make(net.IP, 16)
		binary.BigEndian.PutUint32(addr[0:4], ip[0])
		binary.BigEndian.PutUint32(addr[4:8], ip[1])
		binary.BigEndian.PutUint32(addr[8:12], ip[2])
		binary.BigEndian.PutUint32(addr[12:16], ip[3])
		buf.WriteString("[" + addr.String() + "]")
	default:
		return ""
	}

	if port != 0 {
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%d", Ntohs(port)))
	}
	return buf.String()
}

func ProtocolString(proto uint32) string {
	switch proto {
	case 1:
		return "ICMP"
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	case 58:
		return "ICMPv6"
	default:
		return unknownLabel
	}
}
