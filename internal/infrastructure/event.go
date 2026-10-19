package ebpf

import (
	"encoding/binary"
	"strings"
)

// unknownLabel is the display name used when an event type cannot be mapped.
const unknownLabel = "UNKNOWN"

type EventType int

const (
	EventOpen EventType = iota
	EventRead
	EventWrite
	EventDelete
	EventRename
	EventSymlink
	EventHardlink
	EventMkdir
	EventMmap
	EventAttr
	EventStat
	EventMknod
)

func (t EventType) String() string {
	switch t {
	case EventOpen:
		return "OPEN"
	case EventRead:
		return "READ"
	case EventWrite:
		return "WRITE"
	case EventDelete:
		return "DELETE"
	case EventRename:
		return "RENAME"
	case EventSymlink:
		return "SYMLINK"
	case EventHardlink:
		return "HARDLINK"
	case EventMkdir:
		return "MKDIR"
	case EventMmap:
		return "MMAP"
	case EventAttr:
		return "ATTR"
	case EventStat:
		return "STAT"
	case EventMknod:
		return "MKNOD"
	default:
		return unknownLabel
	}
}

type FileEvent struct {
	PID       uint32
	UID       uint32
	GID       uint32
	Type      EventType
	FD        uint32
	Comm      string
	Path      string
	Dest      string
	Timestamp int64
	// ExeDev/ExeIno is the acting process's executable as the kernel saw it (dev encoded like a
	// guard key), 0/0 for kernel threads; Exe is the path the monitor resolved for that inode, ""
	// when it could not. Identity is ExeDev/ExeIno: Exe and Comm are display only.
	ExeDev uint64
	ExeIno uint64
	Exe    string
	// ExecName is the path the process was exec'd by (bprm->filename) when the monitor saw that
	// exec, else "": it names a multicall's applet.
	ExecName string
	// Prot is an MMAP event's protection (unix.PROT_*); 0 for every other event.
	Prot uint32
}

type BpfEvent struct {
	PID  uint32
	UID  uint32
	GID  uint32
	Type uint32
	FD   uint32
	Pad  uint32
	Comm [16]byte
	Path [256]byte
	Dest [256]byte
	// Dev is the walked dentry's kernel s_dev (MAJOR<<20|MINOR, as mountinfo prints it); 0 when
	// Path/Dest are user strings. Kernel paths are relative to that superblock's root.
	Dev    uint32
	Mode   uint32
	Ino    uint64
	ExeDev uint64
	ExeIno uint64
	Prot   uint32
	Pad2   uint32
}

// BpfEventSize is BpfEvent's wire size as laid out by the BPF side: 6 x u32 header, the three
// fixed byte arrays, dev/mode/ino, exe_dev/exe_ino, prot/pad2. Keep in sync with C struct event {}
// in monitor.bpf.c.
const BpfEventSize = 6*4 + 16 + 256 + 256 + 4 + 4 + 8 + 8 + 8 + 4 + 4

// DecodeBpfEvent fills e from a raw ring-buffer record without reflection or a per-record
// bytes.Reader (binary.Read does both, which adds up on the monitor's unfiltered VFS firehose).
// False for a short record.
func DecodeBpfEvent(raw []byte, e *BpfEvent) bool {
	if len(raw) < BpfEventSize {
		return false
	}
	e.PID = binary.LittleEndian.Uint32(raw[0:4])
	e.UID = binary.LittleEndian.Uint32(raw[4:8])
	e.GID = binary.LittleEndian.Uint32(raw[8:12])
	e.Type = binary.LittleEndian.Uint32(raw[12:16])
	e.FD = binary.LittleEndian.Uint32(raw[16:20])
	e.Pad = binary.LittleEndian.Uint32(raw[20:24])
	copy(e.Comm[:], raw[24:40])
	copy(e.Path[:], raw[40:296])
	copy(e.Dest[:], raw[296:552])
	e.Dev = binary.LittleEndian.Uint32(raw[552:556])
	e.Mode = binary.LittleEndian.Uint32(raw[556:560])
	e.Ino = binary.LittleEndian.Uint64(raw[560:568])
	e.ExeDev = binary.LittleEndian.Uint64(raw[568:576])
	e.ExeIno = binary.LittleEndian.Uint64(raw[576:584])
	e.Prot = binary.LittleEndian.Uint32(raw[584:588])
	e.Pad2 = binary.LittleEndian.Uint32(raw[588:592])
	return true
}

func (e *BpfEvent) ToFileEvent() FileEvent {
	return FileEvent{
		PID:  e.PID,
		UID:  e.UID,
		GID:  e.GID,
		Type: EventType(e.Type),
		FD:   e.FD,
		Comm: Cstr(e.Comm[:]),
		Path: Cstr(e.Path[:]),
		Dest: Cstr(e.Dest[:]),

		ExeDev: e.ExeDev,
		ExeIno: e.ExeIno,
		Prot:   e.Prot,
	}
}

func Cstr(b []byte) string {
	i := 0
	for i < len(b) && b[i] != 0 {
		i++
	}
	return string(b[:i])
}

func ParseEventType(s string) (EventType, bool) {
	switch strings.ToUpper(s) {
	case "OPEN":
		return EventOpen, true
	case "READ":
		return EventRead, true
	case "WRITE":
		return EventWrite, true
	case "DELETE":
		return EventDelete, true
	case "RENAME":
		return EventRename, true
	case "SYMLINK":
		return EventSymlink, true
	case "HARDLINK":
		return EventHardlink, true
	case "MKDIR":
		return EventMkdir, true
	case "MMAP":
		return EventMmap, true
	case "ATTR":
		return EventAttr, true
	case "STAT":
		return EventStat, true
	case "MKNOD":
		return EventMknod, true
	default:
		return 0, false
	}
}

func EventTypes() []EventType {
	return []EventType{
		EventOpen, EventRead, EventWrite, EventDelete,
		EventRename, EventSymlink, EventHardlink, EventMkdir, EventMmap,
		EventAttr, EventStat, EventMknod,
	}
}
