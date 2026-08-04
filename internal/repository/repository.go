package repository

import (
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
)

// MonitorRepository is the port implemented by the file system monitor
// engine. It exposes the full lifecycle of the eBPF event pipeline.
type MonitorRepository interface {
	Start() error
	Stop()
	Events() <-chan ebpf.FileEvent
	SetEventTypes(types []ebpf.EventType)
}

// GuardRepository is the port implemented by the file system guard engine.
type GuardRepository interface {
	Start() error
	Stop()
	Events() <-chan guard.GuardEvent
}

// NetworkMonitorRepository is the port implemented by the network monitor
// engine.
type NetworkMonitorRepository interface {
	Start() error
	Stop()
	Events() <-chan ebpf.NetEvent
	SetEventTypes(types []ebpf.NetEventType)
}

// NetworkGuardRepository is the port implemented by the network guard engine.
type NetworkGuardRepository interface {
	Start() error
	Stop()
	Events() <-chan networkguard.NetGuardEvent
}
