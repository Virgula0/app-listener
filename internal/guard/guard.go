package guard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type Mode int

const (
	ModeBlacklist Mode = iota
	ModeWhitelist
)

const (
	GUARD_BLOCK = 1
	GUARD_ALLOW = 2
)

// BinaryEntry identifies a binary the guard keys on. It is the shared
// infrastructure type; this alias keeps the guard's public API unchanged.
type BinaryEntry = ebpf.BinaryEntry

// ComputeBinaryEntry hashes a binary and derives its comm.
var ComputeBinaryEntry = ebpf.ComputeBinaryEntry

type GuardEvent struct {
	ebpf.FileEvent
	Blocked bool
}

const maxResolveAttempts = 5

type deferredBinary struct {
	rule     daemonconfig.BinaryRule
	attempts int
}

type Guard struct {
	objs    GuardObjects
	links   []link.Link
	events  chan GuardEvent
	done    chan struct{}
	mu      sync.Mutex
	stopped bool

	path      string
	mode      Mode
	binaries  []BinaryEntry
	exeEvents map[string][]ebpf.EventType
	recursive bool
	depth     int
	// canonicalPaths maps each binary's configured path to its symlink-
	// resolved real path (computed at NewGuard time, when the binary was
	// already readable). It is used to attribute events to guarded
	// binaries whose config entry is a symlink, so that a process
	// launched through the symlink is not misreported as "spoofed comm".
	canonicalPaths map[string]string
	// deferred stores whitelist entries that could not be read while
	// their resource tree was still fscrypt-locked. They stay unlisted
	// (denied in whitelist mode) until ResolvePendingBinaries runs after
	// the unlock.
	deferred []deferredBinary
	// deployed tracks, per whitelisted binary path (symlink-canonicalized),
	// the (dev, ino) currently written into the BPF maps. An application
	// update that replaces a binary in place leaves the map keyed by a stale
	// inode, denying the replaced binary; ReSyncBinaries re-stats the path
	// and rewrites changed keys. Keys are never deleted from the maps, so a
	// pre-replacement process that still runs keeps being admitted.
	deployed map[string]GuardInodeKey
	// eagerPopulate scans the whole guarded tree into guard_inodes while
	// the LSM hooks are not yet attached (see WithEagerPopulate).
	eagerPopulate bool
}

// GuardOption customizes a Guard before its BPF maps are populated.
type GuardOption func(*Guard)

// WithEagerPopulate registers every file and directory under the guarded
// path in guard_inodes BEFORE the LSM hooks attach. In whitelist mode the
// guard blocks every binary that is not allowlisted — including its own
// process — so a post-attach startup walk (the eager population) would be
// blocked by the guard's own file_open hook and fail with EPERM. Guards
// that must scan only after an unlock (the daemon's fscrypt flow) omit
// this option and call PopulateInodes themselves.
func WithEagerPopulate() GuardOption {
	return func(g *Guard) {
		g.eagerPopulate = true
	}
}

// WithBinaryEvents restricts each listed binary to the given event types
// (blocking semantics: unlisted events are denied with EPERM). Binaries
// without an entry keep allowing every event, preserving the plain guard
// behavior.
func WithBinaryEvents(events map[string][]ebpf.EventType) GuardOption {
	return func(g *Guard) {
		g.exeEvents = events
	}
}

// WithPendingBinaries registers whitelist entries that could not be read
// while their resource tree was still locked (see daemonconfig.Resource.
// PendingBinaries). They are left out of the BPF whitelist — denied in
// whitelist mode — until ResolvePendingBinaries succeeds after the unlock,
// so the whitelist never admits a binary ahead of its tree being usable.
func WithPendingBinaries(rules []daemonconfig.BinaryRule) GuardOption {
	return func(g *Guard) {
		g.deferred = make([]deferredBinary, len(rules))
		for i, r := range rules {
			g.deferred[i] = deferredBinary{rule: r}
		}
	}
}

// eventMask converts event types into the BPF bitmask stored in
// guard_exe_events. Listing READ, WRITE or MMAP implicitly allows OPEN: a
// binary cannot perform those operations without opening the file first.
// The returned mask never has the OPEN bit cleared when any of those three
// is present.
func eventMask(types []ebpf.EventType) (uint32, error) {
	var mask uint32
	for _, t := range types {
		if t < 0 || t >= 32 {
			return 0, fmt.Errorf("event type %d out of range", t)
		}
		mask |= 1 << uint(t) //nolint:gosec // t is range-checked to [0, 32) above
		switch t {
		case ebpf.EventRead, ebpf.EventWrite, ebpf.EventMmap:
			mask |= 1 << uint(ebpf.EventOpen)
		}
	}
	return mask, nil
}

func NewGuard(path string, mode Mode, binaries []BinaryEntry, recursive bool, depth int, opts ...GuardOption) (*Guard, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	entries := make([]BinaryEntry, len(binaries))
	copy(entries, binaries)

	canonical := make(map[string]string, len(entries))
	for _, b := range entries {
		resolved, err := filepath.EvalSymlinks(b.Path)
		if err != nil {
			resolved = b.Path // unresolvable now (stays equal, fail-closed)
		}
		canonical[b.Path] = resolved
	}

	g := &Guard{
		events:         make(chan GuardEvent, 1024),
		done:           make(chan struct{}),
		path:           path,
		mode:           mode,
		binaries:       entries,
		canonicalPaths: canonical,
		recursive:      recursive,
		depth:          depth,
		deployed:       make(map[string]GuardInodeKey),
	}
	for _, opt := range opts {
		opt(g)
	}

	var objs GuardObjects
	if err := LoadGuardObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("loading guard BPF objects: %w", err)
	}
	g.objs = objs

	// Populate BPF maps first (sets operating mode and inode map), then
	// attach LSM hooks. This avoids the guard blocking its own filesystem
	// operations during startup (readdir, stat, etc.).
	if err := g.populateMaps(); err != nil {
		g.cleanup()
		return nil, fmt.Errorf("populating BPF maps: %w", err)
	}

	// Eagerly register the whole tree in guard_inodes while the hooks are
	// still detached. populateMaps only records the guarded path's own
	// inode; the full recursive scan must run pre-attach, otherwise the
	// guard's own walk is blocked by its own file_open hook in whitelist
	// mode (the guard's exe is not allowlisted) and startup fails with
	// EPERM. Guards that need to scan after an unlock (daemon fscrypt
	// flow) skip this and call PopulateInodes themselves.
	if g.eagerPopulate {
		if err := g.PopulateInodes(); err != nil {
			g.cleanup()
			return nil, fmt.Errorf("populating inode map for %s (is the path readable?): %w", g.path, err)
		}
	}

	// For each hook, track which hooks are critical (must attach for
	// the guard to provide meaningful protection):
	//
	//   file_open, file_permission – required.  Without these, files
	//   can be opened and read without restriction.  The guard will
	//   fail if they cannot attach.
	//
	//   All other hooks – optional.  If they fail, a reduced set of
	//   protection is available (e.g. no mmap, rename, or unlink
	//   blocking).
	required := map[string]bool{
		"file_open":       true,
		"file_permission": true,
	}

	attachments := guardLSMHooks(g)

	var failedRequired []string
	for _, a := range attachments {
		l, err := link.AttachLSM(link.LSMOptions{
			Program: a.prog,
		})
		if err != nil {
			if required[a.hook] {
				failedRequired = append(failedRequired, a.hook)
				log.Errorf("CRITICAL: required LSM hook %s failed to attach: %v", a.hook, err)
			} else {
				log.Warnf("skipping optional LSM hook %s: %v", a.hook, err)
			}
			continue
		}
		g.links = append(g.links, l)
	}

	if len(failedRequired) > 0 {
		g.cleanup()
		return nil, fmt.Errorf(
			"required LSM hooks failed to attach: %v — read/write/open protection unavailable. "+
				"Ensure your kernel supports BPF LSM (CONFIG_BPF_LSM=y) and LSM=bpf is in the "+
				"boot command line (/sys/kernel/security/lsm). "+
				"The guard REQUIRES the file_open and file_permission hooks; if they cannot attach, "+
				"the guard cannot provide meaningful protection and will refuse to start",
			failedRequired)
	}

	log.Infof("guard created \u2014 %d/%d LSM hooks attached, watching: %s (%s)",
		len(g.links), len(attachments), path, modeString(mode))
	return g, nil
}

func guardLSMHooks(g *Guard) []struct {
	prog *cilium.Program
	hook string
} {
	return []struct {
		prog *cilium.Program
		hook string
	}{
		{g.objs.GuardFileOpen, "file_open"},
		{g.objs.GuardFilePermission, "file_permission"},
		{g.objs.GuardFileTruncate, "file_truncate"},
		{g.objs.GuardMmapFile, "mmap_file"},
		{g.objs.GuardPathUnlink, "path_unlink"},
		{g.objs.GuardPathRename, "path_rename"},
		{g.objs.GuardPathSymlink, "path_symlink"},
		{g.objs.GuardPathLink, "path_link"},
		{g.objs.GuardPathMkdir, "path_mkdir"},
		{g.objs.GuardPathTruncate, "path_truncate"},
		{g.objs.GuardInodeSetattr, "inode_setattr"},
		{g.objs.GuardInodeSetxattr, "inode_setxattr"},
		{g.objs.GuardInodeRemovexattr, "inode_removexattr"},
		{g.objs.GuardPathMknod, "path_mknod"},
		{g.objs.GuardPathRmdir, "path_rmdir"},
		{g.objs.GuardInodePermission, "inode_permission"},
		{g.objs.GuardInodeGetattr, "inode_getattr"},
		{g.objs.GuardInodeReadlink, "inode_readlink"},
		{g.objs.GuardSbMount, "sb_mount"},
		{g.objs.GuardPtraceAccessCheck, "ptrace_access_check"},
		{g.objs.GuardTaskAlloc, "task_alloc"},
		{g.objs.GuardTaskFree, "task_free"},
	}
}

func modeString(mode Mode) string {
	if mode == ModeWhitelist {
		return "whitelist"
	}
	return "blacklist"
}

// guardModeKey converts a Mode to the uint64 value stored in the BPF config
// map, rejecting unknown or negative modes.
func guardModeKey(m Mode) (uint64, error) {
	switch m {
	case ModeBlacklist, ModeWhitelist:
		return uint64(m), nil //nolint:gosec // m is validated to one of two small enum values above
	default:
		return 0, fmt.Errorf("invalid guard mode %d", m)
	}
}

// addBinaryActions stores the per-binary allow/block flags in
// guard_exe_actions, keyed by the binary's filesystem inode.
func (g *Guard) addBinaryActions(binaries []BinaryEntry) error {
	for _, b := range binaries {
		dev, ino, err := ebpf.StatInode(b.Path)
		if err != nil {
			return fmt.Errorf("storing exe action for %s: %w", b.Path, err)
		}
		inodeKey := GuardInodeKey{
			Dev: dev,
			Ino: ino,
		}
		action := uint8(GUARD_BLOCK)
		if g.mode == ModeWhitelist {
			action = uint8(GUARD_ALLOW)
		}
		if err := g.objs.GuardExeActions.Put(inodeKey, action); err != nil {
			return fmt.Errorf("storing exe action for %s: %w", b.Path, err)
		}
		g.mu.Lock()
		g.deployed[canonicalBinaryPath(b.Path)] = inodeKey
		g.mu.Unlock()
	}
	return nil
}

// addBinaryEvents stores the per-binary allowed-event bitmasks in
// guard_exe_events. Only binaries with an explicit event list are stored;
// the BPF layer treats a missing entry as "all events allowed".
func (g *Guard) addBinaryEvents(binaries []BinaryEntry, events map[string][]ebpf.EventType) error {
	for _, b := range binaries {
		types, ok := events[b.Path]
		if !ok {
			continue
		}
		if g.mode != ModeWhitelist {
			return fmt.Errorf("per-binary event restrictions are only supported in whitelist mode (binary %s)", b.Path)
		}
		if len(types) == 0 {
			continue // empty list = all events allowed = no mask entry
		}
		mask, err := eventMask(types)
		if err != nil {
			return fmt.Errorf("invalid event mask for binary %s: %w", b.Path, err)
		}
		dev, ino, err := ebpf.StatInode(b.Path)
		if err != nil {
			return fmt.Errorf("cannot stat binary %s for event mask: %w", b.Path, err)
		}
		if err := g.objs.GuardExeEvents.Put(GuardInodeKey{Dev: dev, Ino: ino}, mask); err != nil {
			return fmt.Errorf("storing exe events for %s: %w", b.Path, err)
		}
	}
	return nil
}

// resolveDeferred computes a BinaryEntry for every deferred rule that is
// currently readable, returning the resolved entries, their event lists
// (keyed by the binary's canonical path) and the rules that remain
// unreadable. It never mutates the guard.
func (g *Guard) resolveDeferred() (resolved []BinaryEntry, events map[string][]ebpf.EventType, stillDeferred []deferredBinary) {
	g.mu.Lock()
	deferredList := append([]deferredBinary(nil), g.deferred...)
	g.mu.Unlock()

	resolved = make([]BinaryEntry, 0, len(deferredList))
	events = make(map[string][]ebpf.EventType, len(deferredList))
	for _, deferred := range deferredList {
		rule := deferred.rule
		path := rule.Path
		if realPath, err := filepath.EvalSymlinks(path); err == nil {
			path = realPath
		}
		entry, err := ComputeBinaryEntry(path)
		if err != nil {
			deferred.attempts++
			if deferred.attempts >= maxResolveAttempts {
				log.Warnf("guard %s: aborting deferred binary, still unreadable after %d attempts: %s", g.path, deferred.attempts, rule.Path)
			} else {
				log.Warnf("guard %s: binary still unreadable, keeping it deferred: %s", g.path, rule.Path)
				stillDeferred = append(stillDeferred, deferred)
			}
			continue
		}
		resolved = append(resolved, entry)
		if len(rule.Events) > 0 {
			events[path] = rule.Events
		}
	}
	return resolved, events, stillDeferred
}

// ResolvePendingBinaries retries the whitelist entries that were deferred
// because their resource was still locked. It must run only after the
// resource is accessible (the daemon calls it post-unlock and post-
// populateInodes). Each now-readable rule is hashed, canonicalized and
// written into the whitelist maps; a rule that is still unreadable (or
// genuinely gone) stays deferred and logged — the whitelist keeps
// excluding it, so protection is never weakened by a failed resolve.
func (g *Guard) ResolvePendingBinaries() error {
	g.mu.Lock()
	hasDeferred := len(g.deferred) > 0
	g.mu.Unlock()

	if !hasDeferred {
		return nil
	}

	// The maps are writable while the program is attached. Entries are
	// keyed by inode, so the write is idempotent by construction.
	resolved, events, stillDeferred := g.resolveDeferred()

	g.mu.Lock()
	g.deferred = stillDeferred
	g.mu.Unlock()

	if len(resolved) == 0 {
		return nil
	}
	if err := g.addBinaryActions(resolved); err != nil {
		return err
	}
	if err := g.addBinaryEvents(resolved, events); err != nil {
		return err
	}

	g.mu.Lock()
	g.binaries = append(g.binaries, resolved...)
	for _, b := range resolved {
		g.canonicalPaths[b.Path] = b.Path
	}
	if g.exeEvents == nil {
		g.exeEvents = make(map[string][]ebpf.EventType)
	}
	for path, types := range events {
		g.exeEvents[path] = types
	}
	g.mu.Unlock()

	log.Infof("guard %s: resolved %d deferred binary whitelist entries", g.path, len(resolved))
	return nil
}

// canonicalBinaryPath resolves symlinks for whitelist bookkeeping, falling
// back to the literal path when the link chain is temporarily broken.
func canonicalBinaryPath(path string) string {
	if target, err := filepath.EvalSymlinks(path); err == nil {
		return target
	}
	return path
}

// putBinaryKey writes the allow/block action for one binary into
// guard_exe_actions and, in whitelist mode, any event mask into
// guard_exe_events. Event masks are skipped in blacklist mode, where a
// blocklisted binary is denied outright.
func (g *Guard) putBinaryKey(key GuardInodeKey, path string, events map[string][]ebpf.EventType) error {
	action := uint8(GUARD_BLOCK)
	if g.mode == ModeWhitelist {
		action = uint8(GUARD_ALLOW)
	}
	if err := g.objs.GuardExeActions.Put(key, action); err != nil {
		return fmt.Errorf("storing exe action for %s: %w", path, err)
	}
	if g.mode == ModeWhitelist {
		if types, ok := events[path]; ok && len(types) > 0 {
			mask, err := eventMask(types)
			if err != nil {
				return fmt.Errorf("invalid event mask for binary %s: %w", path, err)
			}
			if err := g.objs.GuardExeEvents.Put(key, mask); err != nil {
				return fmt.Errorf("storing exe events for %s: %w", path, err)
			}
		}
	}
	return nil
}

// retryDeferredBinaries retries the deferred whitelist rules once,
// admitting every rule that became readable: it writes the map entries
// and merges the resolved entries into the guard's bookkeeping. It is the
// shared bookkeeping of ResolvePendingBinaries and ReSyncBinaries.
func (g *Guard) retryDeferredBinaries(resolved []BinaryEntry, resolvedEvents map[string][]ebpf.EventType, stillDeferred []deferredBinary) error {
	g.mu.Lock()
	g.deferred = stillDeferred
	g.mu.Unlock()

	if len(resolved) == 0 {
		return nil
	}
	if err := g.addBinaryActions(resolved); err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.binaries = append(g.binaries, resolved...)
	if g.exeEvents == nil {
		g.exeEvents = make(map[string][]ebpf.EventType)
	}
	for path, types := range resolvedEvents {
		g.exeEvents[path] = types
	}
	return nil
}

// ReSyncBinaries re-stats every whitelisted binary and rewrites its
// inode-keyed action when the file was replaced in place — an application
// update that swaps a binary file (Discord's updater) changes the inode
// while the path stays the same, denying the replaced binary. Old map
// keys are never deleted: a pre-replacement process still running keeps
// being admitted, and a path that vanished in the middle of an update
// (rename over, then deletion) keeps its previous key until it appears
// again. Rules still deferred because they were unreadable are retried
// here too. It returns the number of binaries whose map entries changed.
func (g *Guard) ReSyncBinaries() (int, error) {
	// Retry deferred rules first; newly resolved binaries are recorded in
	// deployed by addBinaryActions and the pass below leaves them alone.
	resolved, resolvedEvents, stillDeferred := g.resolveDeferred()
	if err := g.retryDeferredBinaries(resolved, resolvedEvents, stillDeferred); err != nil {
		return 0, err
	}

	g.mu.Lock()
	binaries := append([]BinaryEntry(nil), g.binaries...)
	exeEvents := make(map[string][]ebpf.EventType, len(g.exeEvents))
	for path, types := range g.exeEvents {
		exeEvents[path] = types
	}
	g.mu.Unlock()

	var changed int
	for _, b := range binaries {
		path := canonicalBinaryPath(b.Path)
		dev, ino, err := ebpf.StatInode(path)
		if err != nil {
			// The binary vanished mid-update (rename, then removal); the
			// previously deployed key stays valid in the maps.
			continue
		}
		key := GuardInodeKey{Dev: dev, Ino: ino}

		g.mu.Lock()
		upToDate := g.deployed[path] == key
		g.mu.Unlock()
		if upToDate {
			continue
		}
		if err := g.putBinaryKey(key, b.Path, exeEvents); err != nil {
			return changed, err
		}
		g.mu.Lock()
		g.deployed[path] = key
		g.mu.Unlock()
		changed++
	}
	if changed > 0 {
		log.Infof("guard %s: re-synced %d binary inode(s) after replacement", g.path, changed)
	}
	return changed, nil
}

// writeGuardConfig stores the scalar guard configuration in
// guard_config[0..2]: operating mode, recursion flag, and depth limit.
func (g *Guard) writeGuardConfig(modeKey uint64) error {
	if putErr := g.objs.GuardConfig.Put(uint32(0), modeKey); putErr != nil {
		return fmt.Errorf("setting mode in config: %w", putErr)
	}

	recursiveVal := uint64(0)
	if g.recursive {
		recursiveVal = 1
	}
	if putErr := g.objs.GuardConfig.Put(uint32(1), recursiveVal); putErr != nil {
		return fmt.Errorf("setting recursive in config: %w", putErr)
	}

	if g.depth < 0 {
		return fmt.Errorf("invalid depth %d", g.depth)
	}
	if putErr := g.objs.GuardConfig.Put(uint32(2), uint64(g.depth)); putErr != nil {
		return fmt.Errorf("setting depth in config: %w", putErr)
	}
	return nil
}

func (g *Guard) populateMaps() error {
	modeKey, err := guardModeKey(g.mode)
	if err != nil {
		return err
	}
	if cfgErr := g.writeGuardConfig(modeKey); cfgErr != nil {
		return cfgErr
	}

	// The watch root's (dev, ino) is stored for the BPF root-confinement
	// check (guard_config[3..4]): an inode map entry only guards an access
	// when the accessed file's dentry chain reaches this root, so stale
	// entries whose inode number was reused outside the tree (ext4 inode
	// reuse) can never deny an unrelated file.
	rootDev, rootIno, rootStatErr := ebpf.StatInode(g.path)
	if rootStatErr != nil {
		return fmt.Errorf("stating guard root %s: %w", g.path, rootStatErr)
	}
	if putErr := g.objs.GuardConfig.Put(uint32(3), rootDev); putErr != nil {
		return fmt.Errorf("setting root dev in config: %w", putErr)
	}
	if putErr := g.objs.GuardConfig.Put(uint32(4), rootIno); putErr != nil {
		return fmt.Errorf("setting root ino in config: %w", putErr)
	}

	_, statErr := os.Stat(g.path)
	if statErr != nil {
		return fmt.Errorf("stating guarded path %s: %w", g.path, statErr)
	}
	// The guarded path's own inode is stored before the LSM hooks are
	// attached: the BPF ancestor walk then recognizes the whole tree as
	// guarded from the first instant the hooks are live. The inode is
	// statable regardless of encryption state — even while an encrypted
	// tree is locked, its own inode stays visible.
	if addErr := g.addInode(g.path); addErr != nil {
		return addErr
	}

	if binErr := g.addBinaryActions(g.binaries); binErr != nil {
		return binErr
	}

	if eventsErr := g.addBinaryEvents(g.binaries, g.exeEvents); eventsErr != nil {
		return eventsErr
	}

	// Store the guarded path for symlink target matching
	var pathBuf [256]byte
	copy(pathBuf[:], g.path)
	if putErr := g.objs.GuardPath.Put(uint32(0), pathBuf); putErr != nil {
		return fmt.Errorf("storing guarded path: %w", putErr)
	}

	// Detect and block the backing block device
	//
	// When a guarded path is on a real block device (e.g. /dev/sda1,
	// /dev/mapper/cryptlvm), tools like debugfs(8) can open the block
	// device directly and read the file contents without ever calling
	// open() on the guarded file path — bypassing VFS access control.
	//
	// We automatically detect the backing device and add it to a BPF
	// map so that any open() on that block device is blocked.
	if err := g.addBackingBlockDevice(); err != nil {
		log.Warnf("backing block device detection: %v", err)
	}
	if err := g.addFsDeviceGate(); err != nil {
		log.Warnf("filesystem device gate: %v", err)
	}

	return nil
}

// addFsDeviceGate records the filesystem device of the guarded path in
// guard_fs_sbdevs so the BPF ancestor walk can skip every access on other
// filesystems with a single lookup. st_dev is i_sb->s_dev by construction,
// so the recorded device matches the superblock the BPF side reads on
// every filesystem — real block devices, tmpfs and overlayfs (anonymous
// devices), and btrfs subvolumes (which share one superblock and only
// widen the gate conservatively). The gate must be populated
// unconditionally — including major-0 filesystems (tmpfs, overlayfs, anon
// devices), where st_dev carries the minor device number with a zero
// major. Leaving the map empty would make every lookup fail and silently
// disable the ancestor walk there, exposing runtime-created content
// deeper than the discovery bound.
func (g *Guard) addFsDeviceGate() error {
	var s syscall.Stat_t
	if err := syscall.Stat(g.path, &s); err != nil {
		return fmt.Errorf("stating guarded path: %w", err)
	}

	major := unix.Major(s.Dev)
	dev := uint64(major)<<20 | uint64(unix.Minor(s.Dev))
	var val uint8 = 1
	if err := g.objs.GuardFsSbdevs.Put(dev, val); err != nil {
		return fmt.Errorf("storing filesystem device %d:%d in map: %w", major, unix.Minor(s.Dev), err)
	}
	return nil
}

// addBackingBlockDevice resolves the block device that backs the guarded
// path and adds it to guard_fs_devices so that raw opens of that block
// device are blocked.
func (g *Guard) addBackingBlockDevice() error {
	var s syscall.Stat_t
	if err := syscall.Stat(g.path, &s); err != nil {
		return fmt.Errorf("stating guarded path: %w", err)
	}

	major := unix.Major(s.Dev)
	if major == 0 {
		// Pseudo-filesystem (tmpfs, overlay, procfs, etc.) — no backing
		// block device to block.
		return nil
	}

	// Encode as the BPF map key (32-bit dev_t).
	minor := unix.Minor(s.Dev)
	rdev := major<<20 | minor

	var val uint8 = 1
	if err := g.objs.GuardFsDevices.Put(rdev, val); err != nil {
		return fmt.Errorf("storing backing block device %d:%d in map: %w", major, minor, err)
	}

	log.Infof("blocking raw access to backing block device %d:%d for guarded path %s",
		major, minor, g.path)
	return nil
}

func (g *Guard) addInode(path string) error {
	dev, ino, err := ebpf.StatInode(path)
	if err != nil {
		return err
	}

	key := GuardInodeKey{
		Dev: dev,
		Ino: ino,
	}

	var val uint8 = 1
	if err := g.objs.GuardInodes.Put(key, val); err != nil {
		return fmt.Errorf("adding inode %s to map: %w", path, err)
	}
	return nil
}

// PopulateInodes fills guard_inodes with the inodes of every file and
// directory under the guarded path. It must be called once the path is
// readable — after the resource was unlocked — but the LSM hooks are
// already attached by then (NewGuard attaches them before the caller
// unlocks), and the guarded root inode is already in the map, so the
// entire subtree is protected from the moment the hooks attach (the BPF
// ancestor walk). The scan mirrors the original ssh-guard's "marks
// attached, then unlock, then walk" ordering. Scanning is tolerant:
// entries that vanish mid-walk (a live application renaming files) and
// dangling symlinks are skipped with a warning; a full inode map degrades
// the rename/unlink/mmap precision but keeps open/read enforcement (the
// BPF ancestor walk) intact. A tree whose entries all fail to stat is
// almost certainly fscrypt-encrypted and locked, which is reported as an
// error.
func (g *Guard) PopulateInodes() error {
	info, statErr := os.Stat(g.path)
	if statErr != nil {
		return fmt.Errorf("stating guarded path %s: %w", g.path, statErr)
	}
	if info.IsDir() {
		if scanErr := g.scanDirInodes(g.path, 0); scanErr != nil {
			return scanErr
		}
	} else {
		if addErr := g.addInode(g.path); addErr != nil {
			return addErr
		}
	}
	return nil
}

func (g *Guard) scanDirInodes(dir string, currentDepth int) error {
	return walkInodes(dir, g.recursive, g.depth, currentDepth, g.addInode)
}

// walkInodes walks dir depth-first, invoking add for dir and every entry
// below it. Semantics:
//
//   - an entry that fails to stat with ENOENT/ENOTDIR (it vanished during
//     the walk, or is a dangling symlink) is skipped with a warning and the
//     walk continues;
//   - an entry that fails with E2BIG means the inode map is full: the walk
//     stops and reports degraded coverage (open/read enforcement survives
//     through the BPF ancestor walk);
//   - any other per-entry failure is skipped with a warning;
//   - when every entry under a directory fails to stat, the tree is almost
//     certainly encrypted and locked: an error is returned with a hint.
//
// addErrKind classifies a failed addInode stat so the walk can decide what
// the failure means for coverage.
type addErrKind int

const (
	addErrOther   addErrKind = iota // unclassified failure: entry skipped
	addErrMissing                   // ENOENT/ENOTDIR: entry vanished mid-walk
	addErrFull                      // E2BIG: inode map full
)

// classifyAddErr maps an add error to its walk semantics. E2BIG degrades
// coverage (open/read enforcement survives through the BPF ancestor walk),
// anything else is logged and skipped.
func classifyAddErr(path string, err error) addErrKind {
	switch {
	case errors.Is(err, unix.ENOENT), errors.Is(err, unix.ENOTDIR):
		return addErrMissing
	case errors.Is(err, unix.E2BIG):
		log.Warnf("guard inode map full while scanning %s: continuing with degraded coverage", path)
		return addErrFull
	default:
		log.Warnf("guard: skipping unreadable entry %s: %v", path, err)
		return addErrOther
	}
}

func walkInodes(dir string, recursive bool, depthLimit, currentDepth int, add func(string) error) error {
	if err := add(dir); err != nil {
		switch classifyAddErr(dir, err) {
		case addErrMissing, addErrFull:
			// The root vanished mid-walk (nothing left to scan), or the
			// inode map is already full: stop without failing.
			return nil
		default:
			return fmt.Errorf("adding inode %s: %w", dir, err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("reading %s: %w", dir, err)
	}

	return walkEntries(dir, entries, recursive, depthLimit, currentDepth, add)
}

// walkEntries visits the entries of one directory level. It reports an
// error only when a nested recursion fails, or when every entry failed to
// stat (the locked-encrypted-tree signature).
func walkEntries(dir string, entries []os.DirEntry, recursive bool, depthLimit, currentDepth int, add func(string) error) error {
	total := 0
	missing := 0
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !recursive {
				continue
			}
			if depthLimit > 0 && currentDepth+1 >= depthLimit {
				continue
			}
			if err := walkInodes(fullPath, recursive, depthLimit, currentDepth+1, add); err != nil {
				return err
			}
			continue
		}
		total++
		if err := add(fullPath); err != nil {
			switch classifyAddErr(fullPath, err) {
			case addErrMissing:
				missing++
			case addErrFull:
				return nil
			}
		}
	}

	if total > 0 && missing == total {
		return fmt.Errorf("all %d entries under %s failed to stat: the directory is probably fscrypt-encrypted and locked (unlock it before building the guard)",
			total, dir)
	}
	return nil
}

func (g *Guard) Events() <-chan GuardEvent {
	return g.events
}

func (g *Guard) Start() error {
	rd, err := ringbuf.NewReader(g.objs.Rb)
	if err != nil {
		return fmt.Errorf("ringbuf reader: %w", err)
	}

	log.Infof("guard started \u2014 guarding: %s", g.path)
	go g.readLoop(rd)
	return nil
}

func (g *Guard) readLoop(rd *ringbuf.Reader) {
	defer rd.Close()

	for {
		ge, ok := g.readEvent(rd)
		if !ok {
			return
		}
		if ge == nil {
			continue
		}

		select {
		case g.events <- *ge:
		case <-g.done:
			return
		}
	}
}

func (g *Guard) readEvent(rd *ringbuf.Reader) (*GuardEvent, bool) {
	record, err := rd.Read()
	if err != nil {
		if errors.Is(err, ringbuf.ErrClosed) {
			return nil, false
		}
		log.Errorf("ringbuf read error: %v", err)
		return nil, true
	}

	var be bpfGuardEvent
	if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &be); err != nil {
		log.Errorf("decode guard event: %v", err)
		return nil, true
	}

	fe := be.toFileEvent()
	fe.Timestamp = time.Now().UnixNano()

	ge := &GuardEvent{
		FileEvent: fe,
		Blocked:   be.Blocked != 0,
	}

	// comm is telemetry: the warning is diagnostic only, never
	// enforcement (access decisions in the BPF program key on the
	// caller's exe inode).
	g.checkCommSpoof(ge)

	return ge, true
}

// checkCommSpoof warns when an event's comm claims the name of a
// guarded binary while the real binary is something else: blocked
// events carrying a guarded name would mislead the logs, so running
// for blocked events too catches that impersonation. A guarded binary
// whose threads rename themselves (Chromium's "libuv-worker", Bun's
// "Bun Pool N") is normal and produces no warning.
func (g *Guard) checkCommSpoof(ge *GuardEvent) {
	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ge.PID))
	if err != nil {
		return
	}

	if g.isGuardedBinary(exePath) {
		return
	}
	if !commMatchesGuardedBinary(ge.Comm, g.binaries) {
		return
	}

	log.Warnf("process %d (%s) spoofed comm \u2014 actual binary: %s",
		ge.PID, ge.Comm, exePath)
}

// commMatchesGuardedBinary reports whether comm claims the name of one
// of the guarded binaries. The kernel truncates comm to 15 bytes
// (TASK_COMM_LEN), so names are compared truncated on both sides.
func commMatchesGuardedBinary(comm string, binaries []BinaryEntry) bool {
	for _, b := range binaries {
		base := filepath.Base(b.Path)
		if len(base) > 15 {
			base = base[:15]
		}
		if comm == base {
			return true
		}
	}
	return false
}

func (g *Guard) isGuardedBinary(exePath string) bool {
	for _, b := range g.binaries {
		canonical := g.canonicalPaths[b.Path]
		if canonical == "" {
			canonical = b.Path // guard created without canonical map
		}
		if b.Path == exePath || canonical == exePath {
			return true
		}
		abs, err := filepath.EvalSymlinks(exePath)
		if err == nil && (abs == b.Path || abs == canonical) {
			return true
		}
	}
	return false
}

func (g *Guard) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.stopped {
		return
	}
	g.stopped = true
	close(g.done)
	g.cleanup()
}

func (g *Guard) cleanup() {
	for _, l := range g.links {
		l.Close()
	}
	g.links = nil
	g.objs.Close()
}

type bpfGuardEvent struct {
	PID     uint32
	UID     uint32
	GID     uint32
	Type    uint32
	FD      uint32
	Blocked uint32
	Comm    [16]byte
	Path    [256]byte
	Dest    [256]byte
}

func (e *bpfGuardEvent) toFileEvent() ebpf.FileEvent {
	return ebpf.FileEvent{
		PID:  e.PID,
		UID:  e.UID,
		GID:  e.GID,
		Type: ebpf.EventType(e.Type),
		FD:   e.FD,
		Comm: ebpf.Cstr(e.Comm[:]),
		Path: ebpf.Cstr(e.Path[:]),
		Dest: ebpf.Cstr(e.Dest[:]),
	}
}

func BinariesSummary(binaries []BinaryEntry) string {
	parts := make([]string, len(binaries))
	for i, b := range binaries {
		parts[i] = fmt.Sprintf("%s [sha256:%x..%x]", b.Path, b.Hash[:4], b.Hash[28:])
	}
	return strings.Join(parts, ", ")
}
