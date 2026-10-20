package ebpf

import (
	"encoding/binary"
	"math/bits"
	"net"
	"strconv"
	"strings"
	"time"
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
	// NetPtrace and NetTracedExec are the network guard's process gates: memory access to, or a
	// traced exec of, a whitelisted process. Always reported; not -e event types.
	NetPtrace
	NetTracedExec
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
	case NetPtrace:
		return "PTRACE"
	case NetTracedExec:
		return "TRACED_EXEC"
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

// NetRecord is the ringbuf layout of networkmonitor.bpf.c's struct net_event, which is also the
// prefix of networkguard.bpf.c's struct net_guard_event. Field order must match the C structs.
type NetRecord struct {
	PID      uint32
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
	TID      uint32
	NetNS    uint64
	CgroupID uint64
}

// Event converts r, stamping it once at ingestion: dual TUI consumers (local and browser) must
// display the same time for the same event.
func (r *NetRecord) Event() NetEvent {
	return NetEvent{
		PID:       r.PID,
		TID:       r.TID,
		UID:       r.UID,
		GID:       r.GID,
		Type:      NetEventType(r.Type),
		Protocol:  r.Proto,
		Size:      r.Size,
		FD:        r.FD,
		Comm:      Cstr(r.Comm[:]),
		SrcAddr:   FormatAddr(r.AF, r.Saddr[:], r.Sport),
		DstAddr:   FormatAddr(r.AF, r.Daddr[:], r.Dport),
		NetNS:     r.NetNS,
		CgroupID:  r.CgroupID,
		Timestamp: time.Now().UnixNano(),
	}
}

// Ntohs byte-swaps a port the BPF side copied in network order (both release arches are
// little-endian).
func Ntohs(v uint16) uint16 {
	return bits.ReverseBytes16(v)
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
		buf.WriteString(strconv.FormatUint(uint64(Ntohs(port)), 10))
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
