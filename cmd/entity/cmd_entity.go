package entity

import "github.com/Virgula0/app-listener/internal/infrastructure"

const (
	MonitorKey = "monitor"
)

var (
	Recursive  bool
	Depth      int
	GUI        bool
	Headless   bool
	EventTypes []ebpf.EventType
)
