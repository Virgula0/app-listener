package entities

type FileEventType int

const (
	EventRead FileEventType = iota
	EventWrite
	EventCreate
	EventDelete
	EventRename
	EventOpen
)

type FileEvent struct {
	Type      FileEventType
	PID       uint32
	Path      string
	Timestamp int64
}
