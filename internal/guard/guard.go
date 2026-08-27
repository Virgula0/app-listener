package guard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// BinaryEntry aliases the shared infrastructure type, keeping the guard's public API unchanged.
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
	// canonicalPaths maps each configured binary path to its symlink-resolved real path (computed at
	// NewGuard time, when the binary was readable), attributing events to binaries launched via a
	// symlinked config entry so they are not misreported as spoofed comm.
	canonicalPaths map[string]string
	// deferred holds whitelist entries unreadable while their resource tree was fscrypt-locked; they
	// stay unlisted (denied in whitelist mode) until ResolvePendingBinaries runs after the unlock.
	deferred []deferredBinary
	// deployed tracks, per symlink-canonicalized whitelisted path, the (dev, ino) currently in the BPF
	// maps. In-place replacements leave a stale inode key denying the binary until ReSyncBinaries
	// rewrites it; keys are never deleted, so a still-running pre-replacement process keeps admission.
	deployed map[string]GuardInodeKey
	// eagerPopulate scans the whole guarded tree into guard_inodes while LSM hooks are detached (see WithEagerPopulate).
	eagerPopulate bool
}

// GuardOption customizes a Guard before its BPF maps are populated.
type GuardOption func(*Guard)

// WithEagerPopulate registers every file and directory under the guarded path in guard_inodes BEFORE the
// LSM hooks attach: in whitelist mode the guard denies any non-allowlisted binary — including its own
// process — so a post-attach startup walk fails EPERM via its own file_open hook. Unlock-time guards omit
// this option and call PopulateInodes themselves.
func WithEagerPopulate() GuardOption {
	return func(g *Guard) {
		g.eagerPopulate = true
	}
}

// WithBinaryEvents restricts each listed binary to the given event types (unlisted events are denied with
// EPERM); binaries without an entry keep allowing every event, preserving plain guard behavior.
func WithBinaryEvents(events map[string][]ebpf.EventType) GuardOption {
	return func(g *Guard) {
		g.exeEvents = events
	}
}

// WithPendingBinaries registers whitelist entries unreadable while their resource tree was still locked
// (see daemonconfig.Resource.PendingBinaries); they stay out of the BPF whitelist — denied in whitelist
// mode — until ResolvePendingBinaries succeeds post-unlock, never admitting a binary ahead of usability.
func WithPendingBinaries(rules []daemonconfig.BinaryRule) GuardOption {
	return func(g *Guard) {
		g.deferred = make([]deferredBinary, len(rules))
		for i, r := range rules {
			g.deferred[i] = deferredBinary{rule: r}
		}
	}
}

// eventMask converts event types into the BPF bitmask stored in guard_exe_events. Listing READ, WRITE or
// MMAP implicitly allows OPEN (those operations require opening the file first); the OPEN bit is never
// cleared when any of the three is present.
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

	// Populate BPF maps first (mode + inode map), then attach LSM hooks, so the guard does not block
	// its own filesystem operations during startup (readdir, stat).
	if err := g.populateMaps(); err != nil {
		g.cleanup()
		return nil, fmt.Errorf("populating BPF maps: %w", err)
	}

	// Register the whole tree in guard_inodes while hooks are still detached: populateMaps records only
	// the root's inode, and a post-attach walk would be blocked by the guard's own file_open hook in
	// whitelist mode (EPERM). Unlock-time guards skip this and call PopulateInodes themselves.
	if g.eagerPopulate {
		if err := g.PopulateInodes(); err != nil {
			g.cleanup()
			return nil, fmt.Errorf("populating inode map for %s (is the path readable?): %w", g.path, err)
		}
	}

	// Required hooks: without file_open and file_permission files can be opened and read without
	// restriction, so the guard refuses to start if they cannot attach. All other hooks are optional —
	// failure yields reduced protection (no mmap, rename, or unlink blocking).
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

// guardModeKey converts a Mode to the uint64 value stored in the BPF config map, rejecting unknown modes.
func guardModeKey(m Mode) (uint64, error) {
	switch m {
	case ModeBlacklist, ModeWhitelist:
		return uint64(m), nil //nolint:gosec // m is validated to one of two small enum values above
	default:
		return 0, fmt.Errorf("invalid guard mode %d", m)
	}
}

// addBinaryActions stores per-binary allow/block flags in guard_exe_actions, keyed by filesystem inode.
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

// addBinaryEvents stores per-binary allowed-event bitmasks in guard_exe_events for binaries with an explicit
// list only; the BPF layer treats a missing entry as "all events allowed".
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

// resolveDeferred computes a BinaryEntry for every deferred rule currently readable, returning the entries,
// their event lists (keyed by canonical path), and rules still unreadable. It never mutates the guard.
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

// ResolvePendingBinaries retries whitelist entries deferred because their resource was still locked; it must run
// only after the unlock (the daemon calls it post-unlock, post-populateInodes). Readable rules are hashed,
// canonicalized and written to the whitelist maps; unreadable ones stay deferred and logged, never weakening protection.
func (g *Guard) ResolvePendingBinaries() error {
	g.mu.Lock()
	hasDeferred := len(g.deferred) > 0
	g.mu.Unlock()

	if !hasDeferred {
		return nil
	}

	// Maps stay writable while attached; inode-keyed entries make the write idempotent by construction.
	resolved, events, stillDeferred := g.resolveDeferred()
	if err := g.retryDeferredBinaries(resolved, events, stillDeferred); err != nil {
		return err
	}
	if len(resolved) == 0 {
		return nil
	}
	if err := g.addBinaryEvents(resolved, events); err != nil {
		return err
	}

	g.mu.Lock()
	for _, b := range resolved {
		g.canonicalPaths[b.Path] = b.Path
	}
	g.mu.Unlock()

	log.Infof("guard %s: resolved %d deferred binary whitelist entries", g.path, len(resolved))
	return nil
}

// canonicalBinaryPath resolves symlinks for whitelist bookkeeping; falls back to the literal path if broken.
func canonicalBinaryPath(path string) string {
	if target, err := filepath.EvalSymlinks(path); err == nil {
		return target
	}
	return path
}

// putBinaryKey writes one binary's allow/block action into guard_exe_actions and, in whitelist mode, its event
// mask into guard_exe_events; masks are skipped in blacklist mode, where blocklisted binaries are denied outright.
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

// retryDeferredBinaries admits deferred rules that became readable: writes their map entries and merges them
// into the guard's bookkeeping. Shared by ResolvePendingBinaries and ReSyncBinaries.
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

// ReSyncBinaries re-stats whitelisted binaries and rewrites inode-keyed map entries after an in-place replacement;
// stale keys are never deleted, so still-running pre-replacement processes keep being admitted and paths vanished
// mid-update keep their prior key until they reappear. Still-deferred unreadable rules are retried here too.
func (g *Guard) ReSyncBinaries() (int, error) {
	// Retry deferred rules first; addBinaryActions records newly resolved binaries in deployed, so the pass below skips them.
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
			// Vanished mid-update (rename, then removal); the previously deployed key stays valid.
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

// writeGuardConfig stores mode, recursion flag and depth limit in guard_config[0..2].
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

	// The watch root's (dev, ino) feeds the BPF root-confinement check (guard_config[3..4]): an inode
	// entry only guards an access whose dentry chain reaches this root, so inode reuse outside the tree
	// can never deny an unrelated file.
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
	// The guarded path's own inode is stored before the hooks attach, so the BPF ancestor walk recognizes
	// the whole tree from the instant hooks go live; the root inode stays statable even while an encrypted
	// tree is locked.
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

	// Detect and block the backing block device: tools like debugfs(8) can open a real block device
	// (e.g. /dev/sda1, /dev/mapper/cryptlvm) directly and read contents without ever open()ing the guarded
	// path — bypassing VFS access control; the detected device goes into a BPF map so any open() on it is blocked.
	if err := g.addBackingBlockDevice(); err != nil {
		log.Warnf("backing block device detection: %v", err)
	}
	if err := g.addFsDeviceGate(); err != nil {
		log.Warnf("filesystem device gate: %v", err)
	}

	return nil
}

// addFsDeviceGate records the guarded path's fs device in guard_fs_sbdevs so the BPF ancestor walk skips other
// filesystems with one lookup; st_dev equals i_sb->s_dev on all fs types (major-0 tmpfs/overlayfs anon devices,
// shared-superblock btrfs). Must be unconditional: an empty map silently disables the walk there, exposing deeper content.
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

// addBackingBlockDevice adds the block device backing the guarded path to guard_fs_devices so raw opens of it are blocked.
func (g *Guard) addBackingBlockDevice() error {
	var s syscall.Stat_t
	if err := syscall.Stat(g.path, &s); err != nil {
		return fmt.Errorf("stating guarded path: %w", err)
	}

	major := unix.Major(s.Dev)
	if major == 0 {
		// Pseudo-filesystem (tmpfs, overlay, procfs, etc.) — no backing block device.
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

// PopulateInodes fills guard_inodes with the inodes of every file and directory under the guarded path; it must run
// once the path is readable (after the unlock — NewGuard attaches the hooks before the caller unlocks, and the root
// inode is already mapped, so the whole subtree is protected via the BPF ancestor walk from hook attach). The scan is
// tolerant (see walkInodes); degraded coverage weakens rename/unlink/mmap precision but keeps open/read enforcement.
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

// walkInodes walks dir depth-first, invoking add for dir and every entry below it: ENOENT/ENOTDIR failures
// (vanished mid-walk, dangling symlinks) and any other per-entry stat error are skipped with a warning;
// E2BIG stops the walk with degraded coverage (open/read enforcement survives via the BPF ancestor walk);
// a directory whose every entry fails to stat is almost certainly fscrypt-encrypted and locked (error).
//
// addErrKind classifies a failed addInode stat so the walk can decide what the failure means for coverage.
type addErrKind int

const (
	addErrOther   addErrKind = iota // unclassified failure: entry skipped
	addErrMissing                   // ENOENT/ENOTDIR: entry vanished mid-walk
	addErrFull                      // E2BIG: inode map full
)

// classifyAddErr maps an add error to its walk semantics: E2BIG degrades coverage (open/read enforcement
// survives via the BPF ancestor walk), anything else is logged and skipped.
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
			// Root vanished mid-walk (nothing left to scan) or map already full: stop without failing.
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

// walkEntries visits one directory level, erroring only when a nested recursion fails or every entry failed
// to stat (the locked-encrypted-tree signature).
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

	// comm is telemetry: the spoof warning is diagnostic only, never enforcement (BPF decisions key on exe inode).
	g.checkCommSpoof(ge)

	return ge, true
}

// checkCommSpoof warns when an event's comm claims a guarded binary's name while the real binary differs;
// running for blocked events too catches impersonation that would mislead logs. Guarded binaries whose threads
// rename themselves (Chromium's "libuv-worker", Bun's "Bun Pool N") are normal and produce no warning.
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

// commMatchesGuardedBinary reports whether comm claims a guarded binary's name; the kernel truncates comm to
// 15 bytes (TASK_COMM_LEN), so both sides are compared truncated.
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

// BinariesSummary formats the whitelist entries for log lines.
func BinariesSummary(binaries []BinaryEntry) string {
	return ebpf.BinariesSummary(binaries)
}
