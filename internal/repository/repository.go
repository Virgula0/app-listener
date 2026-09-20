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
	// PopulateInodes fills the inode map with the guarded tree. Hooks are already attached
	// (protection live); runs after the resource is unlocked. Tolerant: vanished entries are
	// skipped, a full map degrades coverage but isn't fatal.
	PopulateInodes() error
	// ResolvePendingBinaries retries whitelist entries deferred because their tree was locked at
	// guard build. Call only after the resource is accessible (post-unlock, post-PopulateInodes).
	// Resolved binaries join the running whitelist; still-unreadable ones stay deferred and logged
	// (fail-closed).
	ResolvePendingBinaries() error
	// ReSyncBinaries re-stats every whitelisted binary and rewrites its inode-keyed entry after an
	// in-place replacement (app update). Stale keys are never deleted, so pre-replacement processes
	// keep working. Safe to repeat; also retries still-deferred binaries. Returns how many entries
	// changed.
	ReSyncBinaries() (int, error)
	// SweepInodes is the cheap periodic guard_inodes refresh: re-scans only when the watch root's
	// fingerprint moved (recreated single-file root, or a directory that gained/lost a top-level
	// entry). Deeper changes are covered by BPF runtime discovery and the ancestor walk.
	SweepInodes() error
	// ReconcileInodes is the coarse-cadence counterpart of SweepInodes: deletes guard_inodes
	// entries whose (dev, ino) no longer exists on disk. Everything else only adds, and a stale
	// entry can collide with an unrelated file via inode reuse (false denial). Deletes only when
	// backed by a fresh, complete tree walk, and never the watch root's key.
	ReconcileInodes() error
	// GrantSelfEditAccess widens the root-gated self binary mask for an authenticated live
	// edit-protected session; RevokeSelfEditAccess restores the read-only baseline. Both touch only
	// the uid-0-gated self inode, never the whitelist, and live only in the BPF map (gone on
	// restart).
	GrantSelfEditAccess() error
	RevokeSelfEditAccess() error
	// WithSelfVaultAccess widens this guard's self-access for the fscrypt file-vault's in-place
	// unlock/lock (internal/fscrypt/filevault.go) for the duration of fn, then always restores the
	// baseline. Only single-file resources need it (directory unlock/lock is pure keyring work).
	// Distinct from GrantSelfEditAccess, the broader session-scoped grant for interactive edits.
	WithSelfVaultAccess(fn func() error) error
	// SnapshotTaintedPIDs returns the tgids marked tainted (holding guarded content in memory);
	// RestoreTaintedPIDs re-stamps them into a fresh guard. Together they carry the memory-read
	// taint across a reload: new guards start with empty maps, so without this a SIGHUP would drop
	// process_vm_readv/ptrace protection for processes that already read a guarded file.
	SnapshotTaintedPIDs() ([]uint32, error)
	RestoreTaintedPIDs(pids []uint32) error
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
