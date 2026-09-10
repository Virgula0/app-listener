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
	// ModeReadOnly allows every process to READ the guarded tree while
	// gating every modifying operation (write, truncate, rename, unlink,
	// chmod, mkdir, mknod, xattr, writable mmap) on the whitelist, exactly
	// like ModeWhitelist. It backs the daemon's self-protection guard on
	// /etc/app-listener, whose daemon.conf must stay world-readable while
	// only the app-listener binary can change it. mount(2) over the guarded
	// dir is NOT blocked in this mode (it needs CAP_SYS_ADMIN and blocking
	// it broke systemd's mount-namespace setup for the daemon unit's own
	// ExecReload helper). Mirrors GUARD_MODE_READONLY in guard.bpf.c.
	ModeReadOnly
)

const (
	GUARD_BLOCK = 1
	GUARD_ALLOW = 2
	// GUARD_ALLOW_ROOT is honored by the BPF layer only for uid 0: the
	// guard owner's own binary must keep working (fscrypt ioctls, inode
	// scans) without other local users executing the same file inheriting
	// the allow. See WithSelfAllowBinary.
	GUARD_ALLOW_ROOT = 3
)

// BinaryEntry aliases the shared infrastructure type, keeping the guard's public API unchanged.
type BinaryEntry = ebpf.BinaryEntry

// ComputeBinaryEntry hashes a binary and derives its comm.
var ComputeBinaryEntry = ebpf.ComputeBinaryEntry

type GuardEvent struct {
	ebpf.FileEvent
	Blocked bool
	// RawDevice is set when the event came from the raw block-device gate
	// (a denied open of a block device backing a guarded filesystem) rather
	// than a decision about a watched path. The gate is device-granular, so
	// the event names the device, not any resource — userspace must not
	// attribute it to a specific watched path.
	RawDevice bool
}

// guardReasonRawDevice mirrors GUARD_REASON_RAW_DEVICE in guard.bpf.c.
const guardReasonRawDevice = 1

// RawDeviceResourceLabel is the resource string the daemon logs for a raw
// block-device denial, in place of a (misleading) specific watched path.
const RawDeviceResourceLabel = "raw-block-device"

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
	// selfBinary is the guard owner's own executable (the daemon), registered as GUARD_ALLOW_ROOT —
	// honored only for uid 0 — with a minimal event mask, so the binary is not a universal key for
	// other local users (WithSelfAllowBinary).
	selfBinary *BinaryEntry
	selfEvents []ebpf.EventType
	// selfKey / selfKeySet pin the inode key the self binary was registered
	// under, so GrantSelfEditAccess can widen (and RevokeSelfEditAccess
	// restore) its guard_exe_events mask at runtime for a live edit-protected
	// session. selfGrantMu serializes the grant/revoke pair.
	selfKey     GuardInodeKey
	selfKeySet  bool
	selfGrantMu sync.Mutex
	selfGranted bool
	// binaryVerifyStates pins the admitted binaries' inode keys and content
	// hashes for the in-place replacement detector; verifyStop shuts the
	// verifier goroutine down (see startBinaryHashVerifier).
	binaryVerifyStates map[string]*binaryVerifyState
	verifyStop         chan struct{}
	// degradeStop shuts the BPF degradation watcher down (startDegradeWatch).
	degradeStop chan struct{}
	// deployed tracks, per symlink-canonicalized whitelisted path, the (dev, ino) currently in the BPF
	// maps. In-place replacements leave a stale inode key denying the binary until ReSyncBinaries
	// rewrites it; keys are never deleted, so a still-running pre-replacement process keeps admission.
	deployed map[string]GuardInodeKey
	// eagerPopulate scans the whole guarded tree into guard_inodes while LSM hooks are detached (see WithEagerPopulate).
	eagerPopulate bool
	// sweepRootKey / sweepRootMtime fingerprint the watch root between
	// SweepInodes ticks: a single-file root that its app deletes and recreates
	// gets a new inode (nothing else re-maps it), and a directory root whose
	// own mtime moved gained or lost a top-level entry. An unchanged
	// fingerprint means the periodic re-scan can be skipped entirely — deeper
	// changes are covered by BPF runtime discovery and the ancestor walk.
	sweepRootKey   GuardInodeKey
	sweepRootMtime time.Time
	sweepLastFull  time.Time
	// pinPrefix, when set (WithPinning), is the bpffs path prefix each LSM
	// link is pinned at (prefix + hook name) after it attaches: the pin holds
	// the link — and through it the program and its maps — alive past process
	// death, so a SIGKILL (or OOM, or power loss between crash and restart)
	// leaves the guarded tree still enforced instead of instantly
	// unprotected. Stop() removes the pins; pins a killed process left behind
	// are retired by CleanupStalePins on the next start. Cleared to "" if
	// pinning fails mid-attach (see pinDegraded).
	pinPrefix string
	// pinDegraded is set when pinning was requested but the kernel refused it
	// (hardened bpffs, unusual mount): the guard still enforces while the
	// process is alive, but it will NOT survive a SIGKILL. Surfaced so the
	// daemon can warn loudly.
	pinDegraded bool
	// rawDevices / rawDevicesSet override the raw block-device gate's default
	// "derive the device from my own watched path". When rawDevicesSet is
	// true the guard writes exactly rawDevices (rdev-encoded major<<20|minor)
	// into guard_fs_devices — an empty slice disables the gate for this
	// guard. The daemon computes the union of every resource's backing device
	// once and hands it to a single guard (the rest opt out with an empty
	// set), so the same device is not stamped — and mis-attributed — N times.
	// See WithBackingDevices.
	rawDevices    []uint32
	rawDevicesSet bool
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

// WithSelfAllowBinary registers the guard owner's own executable with a root-gated allow action
// (GUARD_ALLOW_ROOT): the BPF layer honors it only when the caller's uid is 0, and the event mask
// restricts it to the listed event types. The daemon needs to open and read its guarded resources
// for the fscrypt lifecycle; without the uid gate any local user executing the same binary file
// would inherit that allow — a universal key over every guarded tree. Whitelist mode only; the
// entry is kept out of the plain whitelist so re-sync/deferred resolution can never re-register
// it as an unconditional GUARD_ALLOW.
func WithSelfAllowBinary(entry BinaryEntry, events []ebpf.EventType) GuardOption {
	return func(g *Guard) {
		g.selfBinary = &entry
		g.selfEvents = events
	}
}

// WithBackingDevices overrides the raw block-device gate for this guard: it
// writes exactly rdevs (each major<<20|minor, from BackingDevice) into
// guard_fs_devices instead of deriving one device from the watched path. Pass
// the union of every guarded resource's backing device to a single guard and
// an empty slice to the others — the gate is device-granular (a block device
// is the whole filesystem), so one guard enforcing the union covers the whole
// daemon, and stamping it once keeps the denial attributable to the daemon's
// raw-device protection rather than a random resource. A nil/empty slice
// disables the gate for the guard it is passed to.
func WithBackingDevices(rdevs []uint32) GuardOption {
	return func(g *Guard) {
		g.rawDevices = rdevs
		g.rawDevicesSet = true
	}
}

// WithPinning pins every attached LSM link at prefix+<hook> on bpffs, so the
// guard keeps enforcing after the daemon process dies (SIGKILL, OOM, power
// loss). prefix must be unique per guard instance and encode the daemon
// generation — see guard.PinPrefix. Stop() unpins; CleanupStalePins retires
// pins a killed daemon left behind.
func WithPinning(prefix string) GuardOption {
	return func(g *Guard) {
		g.pinPrefix = prefix
	}
}

// eventMask converts event types into the BPF bitmask stored in
// guard_exe_events. Listing READ, WRITE or MMAP implicitly allows OPEN and
// STAT: a binary trusted to read, write or map a guarded file's contents is
// necessarily trusted to open it and to see its metadata, and every real
// consumer stat()s a file before it uses it. Without the implied STAT bit a
// restricted mask like `ssh READ,WRITE` is denied at the inode_getattr hook
// (op=STAT) the moment ssh probes ~/.ssh/config, breaking the whitelisted
// binary. The implied bits are never cleared when any of the three is present.
// DELETE / RENAME / HARDLINK stay independent (an unlink needs no open).
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
			mask |= 1 << uint(ebpf.EventStat)
		}
	}
	return mask, nil
}

// VerifyLoad loads every guard eBPF program into THIS kernel's verifier and
// immediately releases them. It attaches nothing and changes no kernel state.
//
// A non-nil error means at least one program was rejected by the running
// kernel — either it exceeds the verifier's complexity budget or it fails a
// CO-RE relocation against this kernel's BTF. Enforcement is then impossible
// (and forcing a partially-loaded guard would be a bypass), so every caller
// MUST treat a failure as fatal: the daemon refuses to start, and `install`
// aborts before enabling the service instead of leaving a crash-looping unit.
//
// A prebuilt binary carries eBPF objects compiled against whatever kernel the
// release runner had; CO-RE fixes field offsets at load time but not verifier
// complexity, so this check is the only thing standing between "runs fine" and
// "the LSM hook corrupts kernel memory" on a kernel newer than the build host.
func VerifyLoad() error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("removing memlock rlimit (need CAP_SYS_RESOURCE / root): %w", err)
	}
	spec, err := LoadGuard()
	if err != nil {
		return fmt.Errorf("reading embedded guard objects: %w", err)
	}
	coll, err := cilium.NewCollection(spec)
	if err != nil {
		// Dump the FULL verifier log (%+v) to a file — the wrapped error is
		// truncated, and the instruction-by-instruction log is what shows
		// where a program blows the complexity budget.
		var ve *cilium.VerifierError
		if errors.As(err, &ve) {
			const logPath = "/tmp/app-listener-verifier.log"
			if werr := os.WriteFile(logPath, []byte(fmt.Sprintf("%+v\n", ve)), 0o600); werr == nil {
				log.Errorf("full verifier log written to %s (read with sudo)", logPath)
			}
			// Echo the tail (where the rejection actually happens) straight to
			// the terminal / journal so it needs no file chasing.
			log.Errorf("verifier log tail:\n%v", ve)
		}
		return fmt.Errorf("this kernel's verifier rejected a guard eBPF program: %w", err)
	}
	coll.Close()
	return nil
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

	failedRequired, total := g.attachHooks()

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
		len(g.links), total, path, modeString(mode))
	return g, nil
}

// requiredHooks are the LSM hooks without which the guard cannot provide
// meaningful protection: files could be opened and read unchecked.
var requiredHooks = map[string]bool{"file_open": true, "file_permission": true}

// attachHooks attaches every guard LSM program and, when pinning is enabled,
// pins each link at g.pinPrefix+<hook> so it survives process death. A
// required hook that fails to attach is collected (the caller aborts); an
// optional one is warned and skipped. A pin failure does NOT abort: the guard
// still enforces while the process runs, so it degrades (pinDegraded) with a
// CRITICAL log and drops the already-made pins rather than leave the tree
// unprotected entirely.
func (g *Guard) attachHooks() (failedRequired []string, total int) {
	attachments := guardLSMHooks(g)
	for _, a := range attachments {
		l, attachErr := link.AttachLSM(link.LSMOptions{Program: a.prog})
		if attachErr != nil {
			if requiredHooks[a.hook] {
				failedRequired = append(failedRequired, a.hook)
				log.Errorf("CRITICAL: required LSM hook %s failed to attach: %v", a.hook, attachErr)
			} else {
				log.Warnf("skipping optional LSM hook %s: %v", a.hook, attachErr)
			}
			continue
		}
		g.links = append(g.links, l)
		if g.pinPrefix != "" {
			// Hook names carry underscores; some hardened bpffs implementations
			// only accept [a-z0-9-] in pin names, so normalise here.
			pinPath := g.pinPrefix + strings.ReplaceAll(a.hook, "_", "-")
			if pinErr := l.Pin(pinPath); pinErr != nil {
				log.Errorf("guard %s: CRITICAL: LSM link pinning failed at %s (%v) \u2014 this guard will NOT "+
					"survive a SIGKILL. The daemon keeps running with live enforcement; investigate bpffs "+
					"(kernel hardening, mount options).", g.path, pinPath, pinErr)
				for _, prev := range g.links {
					_ = prev.Unpin()
				}
				g.pinPrefix = ""
				g.pinDegraded = true
			}
		}
	}
	return failedRequired, len(attachments)
}

// PinDegraded reports that link pinning was requested for this guard but the
// kernel refused it: enforcement is live but will not survive a SIGKILL.
func (g *Guard) PinDegraded() bool { return g.pinDegraded }

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
		{g.objs.GuardBprmCheckSecurity, "bprm_check_security"},
		{g.objs.GuardTaskAlloc, "task_alloc"},
		{g.objs.GuardTaskFree, "task_free"},
	}
}

func modeString(mode Mode) string {
	switch mode {
	case ModeWhitelist:
		return "whitelist"
	case ModeReadOnly:
		return "readonly"
	default:
		return "blacklist"
	}
}

// guardModeKey converts a Mode to the uint64 value stored in the BPF config map, rejecting unknown modes.
func guardModeKey(m Mode) (uint64, error) {
	switch m {
	case ModeBlacklist, ModeWhitelist, ModeReadOnly:
		return uint64(m), nil //nolint:gosec // m is validated to one of three small enum values above
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
		if g.mode == ModeWhitelist || g.mode == ModeReadOnly {
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

// addAllowRootBinary stores the guard owner's own executable as GUARD_ALLOW_ROOT (uid-0 gated in
// the BPF layer) with an optional event mask. It never touches g.binaries, so re-sync and deferred
// resolution cannot re-register it as a plain GUARD_ALLOW.
func (g *Guard) addAllowRootBinary(b BinaryEntry, events []ebpf.EventType) error {
	// ModeReadOnly reuses the whitelist decision path for every modifying
	// operation, so the root-gated self allow applies there too. The
	// per-binary event mask is only consulted by the BPF whitelist branch;
	// in ModeReadOnly the self binary is expected to be registered maskless
	// (all events), which is what the daemon does.
	if g.mode != ModeWhitelist && g.mode != ModeReadOnly {
		return fmt.Errorf("self allow is only supported in whitelist / read-only mode")
	}

	dev, ino, err := ebpf.StatInode(b.Path)
	if err != nil {
		return fmt.Errorf("storing self exe action for %s: %w", b.Path, err)
	}
	key := GuardInodeKey{Dev: dev, Ino: ino}

	if err := g.objs.GuardExeActions.Put(key, uint8(GUARD_ALLOW_ROOT)); err != nil {
		return fmt.Errorf("storing self exe action for %s: %w", b.Path, err)
	}
	g.mu.Lock()
	g.selfKey = key
	g.selfKeySet = true
	g.mu.Unlock()
	if len(events) > 0 {
		mask, err := eventMask(events)
		if err != nil {
			return fmt.Errorf("invalid self event mask for %s: %w", b.Path, err)
		}
		if err := g.objs.GuardExeEvents.Put(key, mask); err != nil {
			return fmt.Errorf("storing self exe events for %s: %w", b.Path, err)
		}
	}
	return nil
}

// editGrantEvents is the widened self event mask a live edit-protected
// session needs: the read set plus every operation the embedded editor's
// atomic save performs (temp file create + write + rename, plus new
// file/dir creation and deletion and mode preservation). SYMLINK, HARDLINK
// and MMAP are deliberately excluded — the editor never needs them and the
// post-edit audit flags any that appeared.
var editGrantEvents = []ebpf.EventType{
	ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat,
	ebpf.EventWrite, ebpf.EventDelete, ebpf.EventRename,
	ebpf.EventMkdir, ebpf.EventMknod, ebpf.EventAttr,
}

// GrantSelfEditAccess widens this guard's root-gated self binary mask to
// cover the operations an interactive edit performs, for the duration of an
// authenticated live edit-protected session. It touches only the self
// (GUARD_ALLOW_ROOT, uid-0-gated) inode key — never the whitelist — and the
// widened mask lives in the BPF map only: it is gone on daemon restart and is
// meant to be paired with RevokeSelfEditAccess. Whitelist / read-only mode
// only (the only modes that register a self binary).
func (g *Guard) GrantSelfEditAccess() error {
	g.selfGrantMu.Lock()
	defer g.selfGrantMu.Unlock()

	g.mu.Lock()
	key, ok := g.selfKey, g.selfKeySet
	g.mu.Unlock()
	if !ok {
		return fmt.Errorf("guard %s: no self binary registered — cannot grant edit access", g.path)
	}

	mask, err := eventMask(editGrantEvents)
	if err != nil {
		return err
	}
	if err := g.objs.GuardExeEvents.Put(key, mask); err != nil {
		return fmt.Errorf("guard %s: widening self edit mask: %w", g.path, err)
	}
	g.selfGranted = true
	log.Warnf("guard %s: edit-protected write access GRANTED to the app-listener binary (uid 0 only) for this resource", g.path)
	return nil
}

// RevokeSelfEditAccess restores the self binary's baseline (read-only) mask
// after a live edit-protected session. Idempotent; safe to call even if the
// grant never took (best-effort teardown path).
func (g *Guard) RevokeSelfEditAccess() error {
	g.selfGrantMu.Lock()
	defer g.selfGrantMu.Unlock()

	g.mu.Lock()
	key, ok := g.selfKey, g.selfKeySet
	base := append([]ebpf.EventType(nil), g.selfEvents...)
	g.mu.Unlock()
	if !ok {
		return nil
	}

	var restoreErr error
	if len(base) == 0 {
		// No baseline mask means "all events allowed": drop the entry.
		if err := g.objs.GuardExeEvents.Delete(key); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			restoreErr = fmt.Errorf("guard %s: clearing self edit mask: %w", g.path, err)
		}
	} else {
		mask, err := eventMask(base)
		if err != nil {
			return err
		}
		if err := g.objs.GuardExeEvents.Put(key, mask); err != nil {
			restoreErr = fmt.Errorf("guard %s: restoring self edit mask: %w", g.path, err)
		}
	}
	if restoreErr == nil && g.selfGranted {
		log.Infof("guard %s: edit-protected write access REVOKED — self binary back to read-only", g.path)
	}
	g.selfGranted = false
	return restoreErr
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
	// Pin the freshly resolved binaries so in-place replacement detection
	// covers them too (their hash was computed from the unlocked content).
	if g.binaryVerifyStates != nil {
		for i := range resolved {
			entry := resolved[i]
			canonical := canonicalBinaryPath(entry.Path)
			if key, ok := g.deployed[canonical]; ok {
				g.binaryVerifyStates[canonical] = &binaryVerifyState{key: key, hash: entry.Hash}
			}
		}
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

	if g.selfBinary != nil {
		if selfErr := g.addAllowRootBinary(*g.selfBinary, g.selfEvents); selfErr != nil {
			return selfErr
		}
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

// BackingDevice returns the guard_fs_devices map key (kernel dev_t, encoded
// major<<20|minor) for the block device backing path, and whether path has one
// at all — false for tmpfs, overlayfs, procfs and every other major-0
// pseudo-filesystem, which have no block device to raw-read.
func BackingDevice(path string) (rdev uint32, hasDevice bool, err error) {
	var s syscall.Stat_t
	if statErr := syscall.Stat(path, &s); statErr != nil {
		return 0, false, fmt.Errorf("stating %s: %w", path, statErr)
	}
	major := unix.Major(s.Dev)
	if major == 0 {
		return 0, false, nil
	}
	return major<<20 | unix.Minor(s.Dev), true, nil
}

// addBackingBlockDevice populates guard_fs_devices so raw opens of the block
// device(s) hosting guarded content are blocked. With WithBackingDevices the
// caller supplies the exact set (the daemon's cross-resource union, or an
// empty set to opt this guard out); otherwise the device is derived from the
// watched path (standalone `guard`).
func (g *Guard) addBackingBlockDevice() error {
	if g.rawDevicesSet {
		return g.putBackingDevices(g.rawDevices)
	}
	rdev, hasDevice, err := BackingDevice(g.path)
	if err != nil {
		return err
	}
	if !hasDevice {
		// Pseudo-filesystem (tmpfs, overlay, procfs, etc.) — no backing block device.
		return nil
	}
	return g.putBackingDevices([]uint32{rdev})
}

func (g *Guard) putBackingDevices(rdevs []uint32) error {
	var val uint8 = 1
	for _, rdev := range rdevs {
		if err := g.objs.GuardFsDevices.Put(rdev, val); err != nil {
			return fmt.Errorf("storing backing block device %d:%d in map: %w",
				rdev>>20, rdev&0xFFFFF, err)
		}
		log.Infof("blocking raw access to backing block device %d:%d", rdev>>20, rdev&0xFFFFF)
	}
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

// SweepInodes is the cheap periodic refresh of guard_inodes. It replaces an
// unconditional full re-walk on every tick (which, across many guards over
// large trees, was ~14% of the daemon's CPU):
//
//   - single-file watch root: re-map it only when its inode changed (an app
//     that deletes and recreates the file — sqlite journals — would otherwise
//     leave the new inode out of the fast-path map);
//   - directory watch root: re-scan only when the directory's own mtime moved
//     since the last sweep (a top-level entry was added or removed). Deeper
//     additions are still guarded — new dirs are mapped by the BPF path_mkdir
//     hook, new files fall through to the ancestor walk — so a full recursive
//     re-scan every tick is wasted work.
func (g *Guard) SweepInodes() error {
	info, err := os.Stat(g.path)
	if err != nil {
		return fmt.Errorf("stating guarded path %s: %w", g.path, err)
	}

	if !info.IsDir() {
		dev, ino, statErr := ebpf.StatInode(g.path)
		if statErr != nil {
			return statErr
		}
		key := GuardInodeKey{Dev: dev, Ino: ino}
		g.mu.Lock()
		unchanged := key == g.sweepRootKey
		g.mu.Unlock()
		if unchanged {
			return nil
		}
		if addErr := g.addInode(g.path); addErr != nil {
			return addErr
		}
		g.mu.Lock()
		g.sweepRootKey = key
		g.mu.Unlock()
		return nil
	}

	mtime := info.ModTime()
	g.mu.Lock()
	// Re-walk only when a top-level entry changed, and never more often than
	// dirRescanMinGap: a chatty root (an app writing lockfiles beside its
	// data) must not force a full recursive walk every tick. New subtrees in
	// the meantime are still mapped by the BPF path_mkdir hook and covered by
	// the ancestor walk.
	skip := mtime.Equal(g.sweepRootMtime) || time.Since(g.sweepLastFull) < dirRescanMinGap
	g.mu.Unlock()
	if skip {
		return nil
	}
	if scanErr := g.scanDirInodes(g.path, 0); scanErr != nil {
		return scanErr
	}
	now := time.Now()
	g.mu.Lock()
	g.sweepRootMtime = mtime
	g.sweepLastFull = now
	g.mu.Unlock()
	return nil
}

const dirRescanMinGap = 5 * time.Minute

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

	g.startDegradeWatch()

	if g.mode == ModeWhitelist && len(g.binaries) > 0 {
		g.startBinaryHashVerifier()
	}
	return nil
}

// BPF degradation counter slots (guard_degrade).
const (
	degradeInodesFull = 0 // guard_inodes map full: discovery adds dropped
	degradeTaintFull  = 1 // guard_tainted_pids map full: taint stamps lost
)

const degradeWatchInterval = 30 * time.Second

// startDegradeWatch logs BPF-side degradation that is otherwise silent: map
// allocations failing inside the kernel reduce coverage (inode discovery
// falls back to the ancestor walk) or disable per-process ptrace protection
// without any Go-side error. Purely observability.
func (g *Guard) startDegradeWatch() {
	g.degradeStop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(degradeWatchInterval)
		defer ticker.Stop()
		seen := map[uint32]uint64{
			degradeInodesFull: 0,
			degradeTaintFull:  0,
		}
		labels := map[uint32]string{
			degradeInodesFull: "guard_inodes map full: runtime inode discovery dropped entries, coverage degrades to the ancestor walk",
			degradeTaintFull:  "tainted-pids map full: ptrace/process_vm_readv protection was not stamped for some processes",
		}
		for {
			select {
			case <-g.degradeStop:
				return
			case <-ticker.C:
				for slot, label := range labels {
					var count uint64
					if err := g.objs.GuardDegrade.Lookup(slot, &count); err != nil {
						continue
					}
					if count > seen[slot] {
						log.Warnf("guard %s: degraded coverage: %s (%d time(s) so far)", g.path, label, count)
						seen[slot] = count
					}
				}
			}
		}
	}()
}

// binaryVerifyState is the verifier's pinned identity of one admitted
// whitelist binary: the inode key it was admitted under, the content hash,
// and the stat fingerprint that hash was computed from (so an unchanged
// binary is never re-hashed).
type binaryVerifyState struct {
	key     GuardInodeKey
	stat    ebpf.BinaryStat
	hash    [32]byte
	hashed  bool // a real content hash has been taken at least once
	demoted bool
}

// binaryHashVerifyInterval is short because the check is now cheap: an
// unchanged binary costs one stat and no read (see verifyBinaryHashesOnce),
// so a fast tick keeps in-place-tamper detection responsive without cost.
const binaryHashVerifyInterval = 10 * time.Second

// sharedBinaryHashes dedupes content hashing across every guard in the
// process. The same binary is admitted by many guards (the ssh tool set, and
// especially Discord's ~170 MiB app blob shared by 10 grouped watch paths);
// re-reading and re-hashing it once per guard per tick was the daemon's
// single largest CPU and allocation cost. Keyed by inode; the stored
// fingerprint gates reuse. Realistic size is a few dozen entries.
var (
	sharedHashMu sync.Mutex
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
)

type sharedHashEntry struct {
	stat ebpf.BinaryStat
	hash [32]byte
}

// hashBinaryShared returns path's content hash, reusing one another guard (or
// an earlier tick) already computed for the same inode+fingerprint and only
// reading the file when nothing cached matches.
func hashBinaryShared(path string, fp ebpf.BinaryStat) ([32]byte, error) {
	key := GuardInodeKey{Dev: fp.Dev, Ino: fp.Ino}

	sharedHashMu.Lock()
	if e, ok := sharedHashes[key]; ok && e.stat == fp {
		sharedHashMu.Unlock()
		return e.hash, nil
	}
	sharedHashMu.Unlock()

	entry, err := ComputeBinaryEntry(path)
	if err != nil {
		return [32]byte{}, err
	}

	sharedHashMu.Lock()
	if len(sharedHashes) > 8192 { // pathological only; keeps the map bounded
		sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	}
	sharedHashes[key] = sharedHashEntry{stat: fp, hash: entry.Hash}
	sharedHashMu.Unlock()
	return entry.Hash, nil
}

// startBinaryHashVerifier pins every admitted whitelist binary's content
// hash and periodically re-verifies it: whitelist identity is the exe inode
// alone, so anyone who can write to an admitted binary OUTSIDE the guarded
// tree can replace its content in place (same inode) and inherit the allow.
// A same-inode hash change demotes the BPF map entry to GUARD_BLOCK —
// permanently, until the inode itself changes: a legitimate package-style
// replacement is re-admitted by ReSyncBinaries/reload and the verifier
// re-pins the new identity on its next tick.
func (g *Guard) startBinaryHashVerifier() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.binaryVerifyStates = make(map[string]*binaryVerifyState, len(g.binaries))
	for _, b := range g.binaries {
		canonical := canonicalBinaryPath(b.Path)
		key, ok := g.deployed[canonical]
		if !ok {
			continue // deferred and not yet resolved: nothing admitted yet
		}
		g.binaryVerifyStates[canonical] = &binaryVerifyState{key: key, hash: b.Hash}
	}
	g.verifyStop = make(chan struct{})
	go g.verifyBinaryHashLoop()
}

func (g *Guard) verifyBinaryHashLoop() {
	ticker := time.NewTicker(binaryHashVerifyInterval)
	defer ticker.Stop()
	for {
		select {
		case <-g.verifyStop:
			return
		case <-ticker.C:
			g.verifyBinaryHashesOnce()
		}
	}
}

// verifyBinaryHashesOnce re-hashes admitted binaries whose stat fingerprint
// moved and demotes same-inode replacements. An unchanged binary costs one
// stat and no read. States are owned by this goroutine (created in
// startBinaryHashVerifier, extended by ReSyncBinaries under g.mu), so field
// access needs no extra locking.
func (g *Guard) verifyBinaryHashesOnce() {
	g.mu.Lock()
	states := make(map[string]*binaryVerifyState, len(g.binaryVerifyStates))
	for canonical, st := range g.binaryVerifyStates {
		states[canonical] = st
	}
	g.mu.Unlock()

	for canonical, st := range states {
		fp, err := ebpf.StatBinary(canonical)
		if err != nil {
			continue // vanished mid-update: ReSyncBinaries handles it
		}
		freshKey := GuardInodeKey{Dev: fp.Dev, Ino: fp.Ino}

		// Common path: same inode and same size/mtime/ctime as the last
		// hash — an in-place overwrite always bumps mtime and ctime, so
		// nothing changed. No read, no hash.
		if st.hashed && freshKey == st.key && fp == st.stat {
			continue
		}

		hash, err := hashBinaryShared(canonical, fp)
		if err != nil {
			continue // unreadable right now: keep the current decision
		}

		if freshKey != st.key {
			// Inode changed: a replacement flow (ReSyncBinaries, SIGHUP
			// reload) owns re-admission — adopt the new identity so a
			// legitimate upgrade is never mistaken for tampering.
			st.key = freshKey
			st.stat = fp
			st.hash = hash
			st.hashed = true
			st.demoted = false
			continue
		}
		st.stat = fp
		st.hashed = true
		if st.demoted {
			continue // once tampering is detected the entry stays blocked
		}
		if hash != st.hash {
			// Same inode, different content: the admitted binary was
			// replaced in place. Demote to GUARD_BLOCK — fail-closed.
			if putErr := g.objs.GuardExeActions.Put(st.key, uint8(GUARD_BLOCK)); putErr != nil {
				log.Errorf("guard %s: demoting in-place replaced binary %s: %v", g.path, canonical, putErr)
				continue
			}
			st.demoted = true
			log.Errorf("guard %s: whitelisted binary %s was modified in place (inode unchanged, hash changed) \u2014 whitelist entry demoted to BLOCK", g.path, canonical)
		}
	}
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
		RawDevice: be.Reason == guardReasonRawDevice,
	}

	// comm is telemetry: the spoof warning is diagnostic only, never enforcement (BPF decisions key on exe inode).
	g.checkCommSpoof(ge)

	return ge, true
}

// checkCommSpoof warns when an event's comm claims a guarded binary's name while the real binary differs;
// running for blocked events too catches impersonation that would mislead logs. Guarded binaries whose threads
// rename themselves (Chromium's "libuv-worker", Bun's "Bun Pool N") are normal and produce no warning.
func (g *Guard) checkCommSpoof(ge *GuardEvent) {
	// Cheap string gate first \u2014 this runs once per guard event, and the vast
	// majority of events come from processes whose comm does not collide with
	// any guarded binary name. Only a collision is worth a /proc lookup.
	if !commMatchesGuardedBinary(ge.Comm, g.binaries) {
		return
	}

	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ge.PID))
	if err != nil {
		return
	}
	if g.isGuardedBinary(exePath) {
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
	// Resolve the process's exe once, not once per whitelist entry.
	absExe, absErr := filepath.EvalSymlinks(exePath)
	for _, b := range g.binaries {
		canonical := g.canonicalPaths[b.Path]
		if canonical == "" {
			canonical = b.Path // guard created without canonical map
		}
		if exePath == b.Path || exePath == canonical {
			return true
		}
		if absErr == nil && (absExe == b.Path || absExe == canonical) {
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
	if g.verifyStop != nil {
		close(g.verifyStop)
		g.verifyStop = nil
	}
	if g.degradeStop != nil {
		close(g.degradeStop)
		g.degradeStop = nil
	}
	close(g.done)
	g.cleanup()
}

func (g *Guard) cleanup() {
	for _, l := range g.links {
		// A pinned link outlives Close() until the pin is removed, so unpin
		// first: Stop() means "this guard is going away for good" (clean
		// shutdown, reload swap, or rollback), and leaving the bpffs entry
		// would keep the LSM program attached with no owner. NewGuard pins
		// all-or-nothing, so pinPrefix != "" implies every link here is pinned.
		if g.pinPrefix != "" {
			if err := l.Unpin(); err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Errorf("guard %s: unpinning LSM link failed (%v) — the program stays attached; "+
					"remove its pin file under %s* manually", g.path, err, g.pinPrefix)
			}
		}
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
	Reason  uint32
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
