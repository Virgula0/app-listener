package ebpf

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
	default:
		return "UNKNOWN"
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
}

type bpfEvent struct {
	PID  uint32
	UID  uint32
	GID  uint32
	Type uint32
	FD   uint32
	Pad  uint32
	Comm [16]byte
	Path [256]byte
	Dest [256]byte
}

func (e *bpfEvent) toFileEvent() FileEvent {
	return FileEvent{
		PID:  e.PID,
		UID:  e.UID,
		GID:  e.GID,
		Type: EventType(e.Type),
		FD:   e.FD,
		Comm: cstr(e.Comm[:]),
		Path: cstr(e.Path[:]),
		Dest: cstr(e.Dest[:]),
	}
}

func cstr(b []byte) string {
	i := 0
	for i < len(b) && b[i] != 0 {
		i++
	}
	return string(b[:i])
}
