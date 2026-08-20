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
	// PopulateInodes fills the inode map with the guarded tree's contents.
	// The hooks are already attached (protection is live) when it is
	// called; it runs after the resource was unlocked. Tolerant scanning:
	// entries vanishing mid-scan are skipped, a full map degrades coverage
	// but is not fatal.
	PopulateInodes() error
	// ResolvePendingBinaries retries the whitelist entries that were
	// deferred because their resource tree was still locked when the
	// guard was built. It must be called only after the resource is
	// accessible (post-unlock, post-PopulateInodes). On success the
	// deferred binaries are added to the running whitelist; a binary
	// that stays unreadable is kept deferred and logged — protection
	// remains fail-closed.
	ResolvePendingBinaries() error
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
