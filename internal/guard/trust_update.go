package guard

import (
	"fmt"
	"maps"
	"sync"
	"sync/atomic"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"
)

// UpdaterPlan scopes protection #1 (guard_bin_owner/guard_bin_updaters in guard_trust.bpf.c): a
// protected file may be modified only by an exe sharing one of its bits.
type UpdaterPlan struct {
	// Owners maps each protected binary/library path to the bits of the resources trusting it.
	Owners map[string]uint64
	// Updaters maps each exe to the bits of the resources it may update.
	Updaters map[string]uint64
}

// SetUpdaters (re)applies the updater scoping. Updaters are written before owners: an owner row
// landing first only refuses its app's update for that instant, never grants a foreign one.
func (t *TrustGuard) SetUpdaters(p UpdaterPlan) error {
	if err := syncMap(t.objs.GuardBinUpdaters, resolveBits(p.Updaters)); err != nil {
		return fmt.Errorf("binary updaters: %w", err)
	}
	if err := syncMap(t.objs.GuardBinOwner, resolveBits(p.Owners)); err != nil {
		return fmt.Errorf("binary owners: %w", err)
	}
	t.ownerMu.Lock()
	t.ownerByPath = maps.Clone(p.Owners)
	t.ownerMu.Unlock()
	return nil
}

// AllowReplacement reports whether newKey, now at the whitelisted path, was created by an updater
// of path's resources and never written by another exe since (guard_bin_origin). It also requires
// the old inode's filesystem (a user FUSE mount serves chosen inode numbers and bytes without any
// write-open) and a trusted old inode. On approval the new inode inherits old's trust rows, so the
// updated binary keeps its write-protection, library allowlist and writer bits until the next
// reload; if any row can't be copied the replacement is refused.
func (t *TrustGuard) AllowReplacement(path string, old, newKey GuardInodeKey) bool {
	if old == (GuardInodeKey{}) || old.Dev != newKey.Dev {
		return false
	}
	t.ownerMu.Lock()
	owner := t.ownerByPath[path]
	t.ownerMu.Unlock()
	if owner == 0 {
		return false
	}
	nk := GuardTrustInodeKey{Dev: newKey.Dev, Ino: newKey.Ino}
	var o GuardTrustBinOrigin
	if err := t.objs.GuardBinOrigin.Lookup(nk, &o); err != nil || o.Tainted != 0 {
		return false
	}
	var upd uint64
	if err := t.objs.GuardBinUpdaters.Lookup(o.Exe, &upd); err != nil || upd&owner == 0 {
		return false
	}
	if err := t.adoptRows(GuardTrustInodeKey{Dev: old.Dev, Ino: old.Ino}, nk); err != nil {
		log.Warnf("trust guard: replacement of %s not re-admitted: %v", path, err)
		return false
	}
	return true
}

// adoptRows copies every exe/target-keyed trust row of old to newKey. old must be a trusted file:
// admitting an inode without TRUSTED_BINARY would leave it outside protections #1 and #2. Any
// failed copy removes the rows already written.
func (t *TrustGuard) adoptRows(old, newKey GuardTrustInodeKey) error {
	var flags uint8
	if err := t.objs.GuardTrustedFiles.Lookup(old, &flags); err != nil {
		return fmt.Errorf("old inode is not a trusted file: %w", err)
	}
	bitMaps := []*cilium.Map{t.objs.GuardBinOwner, t.objs.GuardBinUpdaters,
		t.objs.GuardGlobWriters, t.objs.GuardLibdirUsers}
	rollback := func() {
		for _, m := range append([]*cilium.Map{t.objs.GuardTrustedFiles}, bitMaps...) {
			_ = m.Delete(newKey)
		}
	}
	// Bit rows before the TRUSTED flag: a flagged inode with no owner row is already fail-closed.
	for _, m := range bitMaps {
		var bits uint64
		if err := m.Lookup(old, &bits); err != nil {
			continue
		}
		if err := m.Put(newKey, bits); err != nil {
			rollback()
			return fmt.Errorf("copying %s row: %w", m, err)
		}
	}
	if err := t.objs.GuardTrustedFiles.Put(newKey, flags); err != nil {
		rollback()
		return fmt.Errorf("trusting the new inode: %w", err)
	}
	return nil
}

// ReplacementCheck approves re-admitting newKey at a whitelisted path whose admitted inode was old.
type ReplacementCheck func(path string, old, newKey GuardInodeKey) bool

var replacementCheck atomic.Pointer[ReplacementCheck]

// SetReplacementCheck installs the check ReSyncBinaries consults; nil refuses every replacement
// (no trust guard means no provenance, so nothing proves who created the new inode).
func SetReplacementCheck(fn ReplacementCheck) {
	if fn == nil {
		replacementCheck.Store(nil)
		return
	}
	replacementCheck.Store(&fn)
}

func replacementAllowed(path string, old, newKey GuardInodeKey) bool {
	fn := replacementCheck.Load()
	return fn != nil && (*fn)(path, old, newKey)
}

// refusedReplacements remembers (path, inode) pairs already reported, so a denial-driven re-sync
// every 2s does not repeat the warning.
type refusedReplacements struct {
	mu   sync.Mutex
	seen map[string]GuardInodeKey
}

func (r *refusedReplacements) firstTime(path string, key GuardInodeKey) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.seen == nil {
		r.seen = make(map[string]GuardInodeKey)
	}
	if prev, ok := r.seen[path]; ok && prev == key {
		return false
	}
	r.seen[path] = key
	return true
}
