package guard

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

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
	updaters := resolveBits(singleBinaryPaths(p.Updaters, "updater"), t.statBinary())
	if err := syncMap(t.objs.GuardBinUpdaters, updaters); err != nil {
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

// SetBinaryResolver resolves updater and reserved-name writer paths to inodes through fn instead
// of a fresh stat: the daemon passes the inode it hashed and approved for each path, so a binary
// swapped in after the check, or one it refused, never gains write rights.
func (t *TrustGuard) SetBinaryResolver(fn func(path string) (dev, ino uint64, err error)) {
	t.ownerMu.Lock()
	t.binaryStat = fn
	t.ownerMu.Unlock()
}

func (t *TrustGuard) statBinary() statFunc {
	t.ownerMu.Lock()
	defer t.ownerMu.Unlock()
	if t.binaryStat != nil {
		return t.binaryStat
	}
	return ebpf.StatConfined
}

// CreatedByUpdater reports whether nk, the inode f holds, was created since this trust guard
// started by an updater of the resources owning ownerPaths and never write-opened by another exe
// (guard_bin_origin): how a new whitelist line with no previous inode (a fresh version directory)
// proves its provenance. On the updater's own filesystem only, never FUSE: a user FUSE mount serves
// chosen inode numbers and bytes without any write-open.
func (t *TrustGuard) CreatedByUpdater(f *os.File, nk GuardInodeKey, ownerPaths []string) bool {
	t.ownerMu.Lock()
	var owner uint64
	for _, p := range ownerPaths {
		owner |= t.ownerByPath[p]
	}
	t.ownerMu.Unlock()
	if owner == 0 || f == nil {
		return false
	}
	var o GuardTrustBinOrigin
	k := GuardTrustInodeKey{Dev: nk.Dev, Ino: nk.Ino}
	if err := t.objs.GuardBinOrigin.Lookup(k, &o); err != nil || o.Tainted != 0 || o.Exe.Dev != nk.Dev {
		return false
	}
	var fs unix.Statfs_t
	if unix.Fstatfs(int(f.Fd()), &fs) != nil || fs.Type == unix.FUSE_SUPER_MAGIC {
		return false
	}
	var upd uint64
	return t.objs.GuardBinUpdaters.Lookup(o.Exe, &upd) == nil && upd&owner != 0
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
	if err := t.objs.GuardTrustedFiles.Put(nk, trustedBinary|runtimeClass(f, path)); err != nil {
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

// adoptRows copies every exe/target-keyed trust row of old to newKey. old must be a trusted file,
// or one whose rows the prune retired: admitting an inode without TRUSTED_BINARY would leave it
// outside protections #1 and #2. With neither, newKey qualifies only if a sibling guard's
// admission of the same path already adopted them. Any failed copy removes the rows written.
func (t *TrustGuard) adoptRows(old, newKey GuardTrustInodeKey) error {
	r, ok := t.rowsOf(old)
	if !ok {
		var flags uint8
		if t.objs.GuardTrustedFiles.Lookup(newKey, &flags) == nil && flags&trustedBinary != 0 {
			return nil
		}
		return errors.New("old inode is not a trusted file")
	}
	// A temporary grant's row carries runtime bits only (tempgrant.go): not a whitelisted binary's.
	if r.flags&trustedBinary == 0 {
		return errors.New("old inode is not a trusted binary")
	}
	bitMaps := t.trustBitMaps()
	rollback := func() {
		_ = t.objs.GuardTrustedFiles.Delete(newKey)
		for _, m := range bitMaps {
			_ = m.Delete(newKey)
		}
	}
	// Bit rows before the TRUSTED flag: a flagged inode with no owner row is already fail-closed.
	for i, m := range bitMaps {
		if !r.has[i] {
			continue
		}
		if err := m.Put(newKey, r.bits[i]); err != nil {
			rollback()
			return fmt.Errorf("copying %s row: %w", m, err)
		}
	}
	if err := t.objs.GuardTrustedFiles.Put(newKey, r.flags); err != nil {
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
