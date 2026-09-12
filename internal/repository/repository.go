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
	// ReSyncBinaries re-stats every whitelisted binary and rewrites its
	// inode-keyed entry when the file was replaced in place (an
	// application update). Stale map keys are never deleted, so
	// pre-replacement processes keep working. It is safe to call
	// repeatedly; it also retries binaries still deferred from load. It
	// returns the number of binaries whose map entries changed.
	ReSyncBinaries() (int, error)
	// SweepInodes is the cheap periodic guard_inodes refresh: it re-scans
	// only when the watch root's fingerprint moved (a recreated single-file
	// root, or a directory root that gained/lost a top-level entry).
	// Unconditional full re-walks are avoided — deeper changes are covered by
	// BPF runtime discovery and the ancestor walk.
	SweepInodes() error
	// GrantSelfEditAccess widens the root-gated self binary mask to cover the
	// operations an interactive edit performs, for an authenticated live
	// edit-protected session; RevokeSelfEditAccess restores the read-only
	// baseline. Both touch only the uid-0-gated self inode, never the
	// whitelist, and the change lives in the BPF map only (gone on restart).
	GrantSelfEditAccess() error
	RevokeSelfEditAccess() error
	// WithSelfVaultAccess widens this resource's guard self-access to cover
	// the fscrypt file-vault's in-place unlock/lock (internal/fscrypt/
	// filevault.go) for the duration of fn, then unconditionally restores
	// the baseline mask. Directory resources never need this (their fscrypt
	// unlock/lock is a pure kernel-keyring operation that never touches
	// file content); only a single-file resource's vault reads and rewrites
	// its own bytes on the guarded path itself. Deliberately distinct from
	// GrantSelfEditAccess/RevokeSelfEditAccess, which is a broader,
	// session-scoped grant for interactive edit-protected use.
	WithSelfVaultAccess(fn func() error) error
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
