package guard

import (
	"fmt"
	"maps"
	"os"
	"sync"
	"sync/atomic"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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
	if err := syncMap(t.objs.GuardBinUpdaters, resolveBits(p.Updaters, ebpf.StatConfined)); err != nil {
		return fmt.Errorf("binary updaters: %w", err)
	}
	if err := syncMap(t.objs.GuardBinOwner, resolveBits(p.Owners, ebpf.StatConfined)); err != nil {
		return fmt.Errorf("binary owners: %w", err)
	}
	t.ownerMu.Lock()
	t.ownerByPath = maps.Clone(p.Owners)
	t.ownerMu.Unlock()
	return nil
}

// AllowReplacement reports whether newKey, the inode f holds at the whitelisted path, may be
// admitted in place of old: created by one of its updaters (updaterCreated) or a system file at a
// name only root could have placed (systemFile). On approval it inherits old's trust rows, so it keeps its
// write-protection, library allowlist and writer bits until the next reload; a row that can't be
// copied refuses it.
func (t *TrustGuard) AllowReplacement(path string, f *os.File, old, newKey GuardInodeKey) bool {
	nk := GuardTrustInodeKey{Dev: newKey.Dev, Ino: newKey.Ino}
	if t.updaterCreated(path, old, nk) {
		liftSuperseded(newKey)
		if err := t.adoptRows(GuardTrustInodeKey{Dev: old.Dev, Ino: old.Ino}, nk); err != nil {
			log.Warnf("trust guard: replacement of %s not re-admitted: %v", path, err)
			return false
		}
		return true
	}
	if !t.systemFile(path, f, newKey) {
		return false
	}
	liftSuperseded(newKey)
	if old != (GuardInodeKey{}) && t.trusts(old) {
		if err := t.adoptRows(GuardTrustInodeKey{Dev: old.Dev, Ino: old.Ino}, nk); err != nil {
			log.Warnf("trust guard: system binary %s not admitted: %v", path, err)
			return false
		}
		return true
	}
	// A new install: nothing to inherit. Protection #1 exempts system files, so only the library
	// allowlist (TRUSTED_BINARY) applies.
	if err := t.objs.GuardTrustedFiles.Put(nk, trustedBinary); err != nil {
		log.Warnf("trust guard: system binary %s not admitted: %v", path, err)
		return false
	}
	return true
}

// updaterCreated: newKey was created by an updater of path's resources and never written by another
// exe since (guard_bin_origin), on old's filesystem (a user FUSE mount serves chosen inode numbers
// and bytes without any write-open).
func (t *TrustGuard) updaterCreated(path string, old GuardInodeKey, nk GuardTrustInodeKey) bool {
	if old == (GuardInodeKey{}) || old.Dev != nk.Dev {
		return false
	}
	t.ownerMu.Lock()
	owner := t.ownerByPath[path]
	t.ownerMu.Unlock()
	if owner == 0 {
		return false
	}
	var o GuardTrustBinOrigin
	if err := t.objs.GuardBinOrigin.Lookup(nk, &o); err != nil || o.Tainted != 0 {
		return false
	}
	var upd uint64
	return t.objs.GuardBinUpdaters.Lookup(o.Exe, &upd) == nil && upd&owner != 0
}

// systemFile: path is root-placed (ebpf.OpenSystemPlaced: every directory and symlink hop leading
// to it is root's) and ends on k, the inode f holds open so its number can't be reused meanwhile, on
// a superblock root vouches for (mount_vouches_ownership: by superblock device, which on btrfs is
// not k.Dev). The name is judged, not only the inode: a user can point a name it controls at any
// root-owned binary.
func (t *TrustGuard) systemFile(path string, f *os.File, k GuardInodeKey) bool {
	if f == nil {
		return false
	}
	sp, err := ebpf.OpenSystemPlaced(path)
	if err != nil {
		log.Debugf("trust guard: %s is not a system file: %v", path, err)
		return false
	}
	defer sp.Close()
	if dev, ino, statErr := ebpf.StatFile(sp); statErr != nil || dev != k.Dev || ino != k.Ino {
		return false
	}
	sbdev, err := ebpf.SuperblockDev(int(sp.Fd()))
	if err != nil {
		return false
	}
	var synced, died uint64
	if t.objs.GuardVouchedDevs.Lookup(sbdev, &synced) != nil {
		return false
	}
	return t.objs.GuardDeadDevs.Lookup(sbdev, &died) != nil || died <= synced
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

// ReplacementCheck approves admitting newKey, the inode f holds, at a whitelisted path whose
// admitted inode was old (zero for a new catalog match).
type ReplacementCheck func(path string, f *os.File, old, newKey GuardInodeKey) bool

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

func replacementAllowed(path string, f *os.File, old, newKey GuardInodeKey) bool {
	fn := replacementCheck.Load()
	return fn != nil && (*fn)(path, f, old, newKey)
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
