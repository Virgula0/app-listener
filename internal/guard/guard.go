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
	"sync/atomic"
	"syscall"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

type Mode int

const (
	ModeBlacklist Mode = iota
	ModeWhitelist
	// ModeReadOnly lets every process READ the tree but gates all modifying ops (write, truncate,
	// rename, unlink, chmod, mkdir, mknod, xattr, writable mmap) on the whitelist, like
	// ModeWhitelist. Backs the self-protection guard on /etc/app-listener (daemon.conf stays
	// world-readable). mount(2) over the dir is NOT blocked: doing so broke systemd's
	// mount-namespace setup for ExecReload. Mirrors GUARD_MODE_READONLY in guard.bpf.c.
	ModeReadOnly
)

const (
	GUARD_BLOCK = 1
	GUARD_ALLOW = 2
	// GUARD_ALLOW_ROOT is honored only for uid 0, so other users executing the owner's binary don't
	// inherit the allow. See WithSelfAllowBinary.
	GUARD_ALLOW_ROOT = 3
)

// BinaryEntry aliases the shared infrastructure type, keeping the guard's public API unchanged.
type BinaryEntry = ebpf.BinaryEntry

// ComputeBinaryEntry hashes a binary and derives its comm.
var ComputeBinaryEntry = ebpf.ComputeBinaryEntry

type GuardEvent struct {
	ebpf.FileEvent
	Blocked bool
	// RawDevice: the event is a raw block-device gate denial. The gate is device-granular, so the
	// event names the device, not a resource; don't attribute it to a watched path.
	RawDevice bool
	// Process names a process-gate denial ("PTRACE", "TRACED_EXEC", "PROC_MEM"): no file involved;
	// Path is "pid=<n> comm=<name>" of the other task. Empty for path-keyed events.
	Process string
}

// guard_event.reason values — mirror GUARD_REASON_* in guard.bpf.c.
const (
	guardReasonRawDevice  = 1
	guardReasonPtrace     = 2
	guardReasonTracedExec = 3
	guardReasonProcMem    = 4
)

// processGateLabel maps a process-gate reason to its logged op label; "" for path-keyed events.
func processGateLabel(reason uint32) string {
	switch reason {
	case guardReasonPtrace:
		return "PTRACE"
	case guardReasonTracedExec:
		return "TRACED_EXEC"
	case guardReasonProcMem:
		return "PROC_MEM"
	default:
		return ""
	}
}

// RawDeviceResourceLabel is the resource string logged for raw block-device denials.
const RawDeviceResourceLabel = "raw-block-device"

const maxResolveAttempts = 5

type deferredBinary struct {
	rule     daemonconfig.BinaryRule
	attempts int
}

type Guard struct {
	// resID is this resource's slot in the shared engine's kernel tables; every per-resource map
	// key carries it. See engine.
	resID uint32
	// eventsDropped counts events discarded because this guard's consumer stalled (dispatch).
	eventsDropped atomic.Uint64
	events        chan GuardEvent
	done          chan struct{}
	mu            sync.Mutex
	stopped       bool

	path      string
	mode      Mode
	binaries  []BinaryEntry
	exeEvents map[string][]ebpf.EventType
	recursive bool
	depth     int
	// canonicalPaths maps configured binary paths to symlink-resolved real paths (set in NewGuard),
	// so events from symlinked entries aren't reported as spoofed comm.
	canonicalPaths map[string]string
	// deferred: whitelist entries unreadable while the tree was fscrypt-locked; denied until
	// ResolvePendingBinaries runs post-unlock.
	deferred []deferredBinary
	// selfBinary is the owner's own executable (the daemon), registered GUARD_ALLOW_ROOT (uid 0
	// only) with a minimal event mask (WithSelfAllowBinary).
	selfBinary *BinaryEntry
	selfEvents []ebpf.EventType
	// selfKey/selfKeySet: inode key of selfBinary, so Grant/RevokeSelfEditAccess can widen/restore
	// its guard_exe_events mask. selfGrantMu serializes the pair.
	selfKey     GuardInodeKey
	selfKeySet  bool
	selfGrantMu sync.Mutex
	selfGranted bool
	// binaryVerifyStates pins admitted binaries' inode keys and hashes for the in-place replacement
	// detector; verifyStop stops it (startBinaryHashVerifier).
	binaryVerifyStates map[string]*binaryVerifyState
	verifyStop         chan struct{}
	// degradeStop shuts the BPF degradation watcher down (startDegradeWatch).
	degradeStop chan struct{}
	// deployed: per canonical whitelisted path, the (dev, ino) in the BPF maps. Keys are never
	// deleted, so a running pre-replacement process keeps admission; ReSyncBinaries rewrites stale
	// ones.
	deployed map[string]GuardInodeKey
	// eagerPopulate scans the whole guarded tree into guard_inodes while LSM hooks are detached (see WithEagerPopulate).
	eagerPopulate bool
	// sweepRootMtime/sweepLastFull fingerprint a DIRECTORY root between SweepInodes ticks:
	// unchanged mtime = no top-level change, skip the re-scan (deeper changes are covered by BPF
	// discovery and the ancestor walk).
	sweepRootMtime time.Time
	sweepLastFull  time.Time
	// rootKey is the watch root's (dev, ino): sole source for the kernel root-confinement check
	// (guard_config[3]/[4], root_in_chain) and the one key ReconcileInodes never evicts, so a
	// transient stat failure can't unguard the tree. Set in populateMaps; re-anchored by
	// updateRootKey when a single-file root is recreated (a stale anchor stops guarding the real
	// file and guards whatever reuses the freed inode). Guarded by mu.
	rootKey GuardInodeKey
	// pinPrefix (WithPinning) is the bpffs prefix each LSM link is pinned at (prefix+hook) so
	// enforcement survives SIGKILL/OOM. Stop() removes pins; CleanupStalePins retires a killed
	// process's. Cleared if pinning fails mid-attach (pinDegraded).
	pinPrefix string
	// rawDevices/rawDevicesSet override the raw block-device gate's default (device of the own
	// watched path): when set, write exactly rawDevices (major<<20|minor) to guard_fs_devices;
	// empty disables the gate. The daemon passes the cross-resource union to one guard
	// (WithBackingDevices).
	rawDevices    []uint32
	rawDevicesSet bool
	// dropAllowed: no consumer wants allowed events (WithoutAllowedEvents).
	dropAllowed bool
}

// GuardOption customizes a Guard before its BPF maps are populated.
type GuardOption func(*Guard)

// WithEagerPopulate scans the whole tree into guard_inodes BEFORE the hooks attach: in whitelist
// mode a post-attach walk is denied by the guard's own file_open hook. Unlock-time guards omit it
// and call PopulateInodes.
func WithEagerPopulate() GuardOption {
	return func(g *Guard) {
		g.eagerPopulate = true
	}
}

// WithBinaryEvents restricts listed binaries to the given event types (others get EPERM); unlisted
// binaries allow all events.
func WithBinaryEvents(events map[string][]ebpf.EventType) GuardOption {
	return func(g *Guard) {
		g.exeEvents = events
	}
}

// WithPendingBinaries registers entries unreadable while the tree was locked
// (daemonconfig.Resource.PendingBinaries); they stay denied until ResolvePendingBinaries succeeds
// post-unlock.
func WithPendingBinaries(rules []daemonconfig.BinaryRule) GuardOption {
	return func(g *Guard) {
		g.deferred = make([]deferredBinary, len(rules))
		for i, r := range rules {
			g.deferred[i] = deferredBinary{rule: r}
		}
	}
}

// WithSelfAllowBinary registers the owner's executable as GUARD_ALLOW_ROOT (uid 0 only, limited to
// the given events) for the fscrypt lifecycle; the uid gate stops other users running the same file
// from inheriting a universal key. Whitelist mode only; kept out of the plain whitelist so re-sync
// can't re-register it as unconditional GUARD_ALLOW.
func WithSelfAllowBinary(entry BinaryEntry, events []ebpf.EventType) GuardOption {
	return func(g *Guard) {
		g.selfBinary = &entry
		g.selfEvents = events
	}
}

// WithBackingDevices sets exactly rdevs (major<<20|minor, from BackingDevice) in guard_fs_devices
// instead of deriving one from the watched path. Pass the cross-resource union to one guard and an
// empty slice to the rest (the gate is device-granular; this avoids N duplicate, misattributed
// stamps). Empty disables the gate.
func WithBackingDevices(rdevs []uint32) GuardOption {
	return func(g *Guard) {
		g.rawDevices = rdevs
		g.rawDevicesSet = true
	}
}

// WithoutAllowedEvents discards allowed events before they are queued, for a consumer that prints
// only denials: a whitelisted app's normal I/O (Steam's CEF cache) otherwise floods the queue.
func WithoutAllowedEvents() GuardOption {
	return func(g *Guard) {
		g.dropAllowed = true
	}
}

// WithPinning pins each LSM link at prefix+<hook> on bpffs so enforcement survives daemon death.
// prefix must be unique per guard and encode the generation (PinPrefix). Stop() unpins;
// CleanupStalePins retires leftovers.
func WithPinning(prefix string) GuardOption {
	return func(g *Guard) {
		g.pinPrefix = prefix
	}
}

// eventMask converts event types to the guard_exe_events bitmask. READ, WRITE or MMAP imply OPEN
// and STAT (else `ssh READ,WRITE` is denied at inode_getattr when probing ~/.ssh/config).
// DELETE/RENAME/HARDLINK stay independent.
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

// VerifyLoad loads every guard and trust eBPF program into this kernel's verifier and releases
// them; it attaches nothing.
//
// A non-nil error means a program was rejected (verifier complexity budget or CO-RE relocation):
// enforcement is impossible, so callers MUST treat it as fatal (daemon refuses to start; install
// aborts before enabling the service). CO-RE fixes offsets but not verifier complexity, so a
// prebuilt binary can fail on a newer kernel. The trust object is best-effort at runtime, but a
// rejection there silently drops every trust protection, so the preflight refuses it too.
func VerifyLoad() error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("removing memlock rlimit (need CAP_SYS_RESOURCE / root): %w", err)
	}
	for _, obj := range []struct {
		name string
		load func() (*cilium.CollectionSpec, error)
	}{{"guard", LoadGuard}, {"trust", LoadGuardTrust}} {
		spec, err := obj.load()
		if err != nil {
			return fmt.Errorf("reading embedded %s objects: %w", obj.name, err)
		}
		if err := verifyCollection(spec); err != nil {
			return err
		}
	}
	return nil
}

func verifyCollection(spec *cilium.CollectionSpec) error {
	coll, err := cilium.NewCollection(spec)
	if err != nil {
		// Dump the FULL verifier log (%+v) to a file; the wrapped error is truncated.
		var ve *cilium.VerifierError
		if errors.As(err, &ve) {
			const logPath = "/tmp/app-listener-verifier.log"
			if werr := os.WriteFile(logPath, []byte(fmt.Sprintf("%+v\n", ve)), 0o600); werr == nil {
				log.Errorf("full verifier log written to %s (read with sudo)", logPath)
			}
			// Echo the log tail (where the rejection happens) to the terminal/journal.
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

	// Join the shared engine: the LSM programs attach once for the whole process (see engine).
	// The slot starts INACTIVE, so this resource denies nothing until its inodes are registered —
	// the same window the old attach-after-populate ordering gave, while every other resource
	// stays enforced throughout.
	resID, err := sharedEngine.acquire(g)
	if err != nil {
		return nil, err
	}
	g.resID = resID

	if err := g.populateMaps(); err != nil {
		g.cleanup()
		return nil, fmt.Errorf("populating BPF maps: %w", err)
	}

	// Register the whole tree before the resource goes live (see WithEagerPopulate); populateMaps
	// only records the root.
	if g.eagerPopulate {
		if err := g.PopulateInodes(); err != nil {
			g.cleanup()
			return nil, fmt.Errorf("populating inode map for %s (is the path readable?): %w", g.path, err)
		}
	}

	// Enforcement for this resource begins here.
	if err := g.activate(); err != nil {
		g.cleanup()
		return nil, fmt.Errorf("activating guard for %s: %w", g.path, err)
	}

	g.pinSelfMaps()

	log.Infof("guard created \u2014 resource %d watching: %s (%s)", g.resID, path, modeString(mode))
	return g, nil
}

// objs is the shared engine's loaded BPF objects; every Guard keys into them by resID.
func (g *Guard) objs() *GuardObjects { return &sharedEngine.objs }

// activate flips this resource's slot live, so the kernel starts resolving its inodes to a real
// policy. Called only after the inode map and whitelist are populated.
func (g *Guard) activate() error {
	cfg, err := g.resConfig()
	if err != nil {
		return err
	}
	cfg.Active = 1
	return g.objs().GuardResConfig.Put(g.resID, cfg)
}

// resConfig reads this resource's current kernel config slot.
func (g *Guard) resConfig() (GuardResConfig, error) {
	var cfg GuardResConfig
	if err := g.objs().GuardResConfig.Lookup(g.resID, &cfg); err != nil {
		return cfg, fmt.Errorf("reading resource slot %d: %w", g.resID, err)
	}
	return cfg, nil
}

// PinDegraded reports that requested pinning was refused: enforcing, but won't survive SIGKILL.
func (g *Guard) PinDegraded() bool { return SharedPinDegraded() }

// requiredHooks: without these, files could be opened and read unchecked.
var requiredHooks = map[string]bool{"file_open": true, "file_permission": true}

// ExeActionsPinName/ExeEventsPinName are the pin suffixes of guard_exe_actions/guard_exe_events
// (g.pinPrefix+suffix); exported so `daemon --lockdown` (no live Guard) can reopen just those maps.
const (
	ExeActionsPinName = "exe-actions"
	ExeEventsPinName  = "exe-events"
	// ResIDPinName holds one u32: which resource slot this guard owns. The whitelist maps are
	// shared and keyed by (res_id, inode), so a process with no live Guard needs the id to find
	// this resource's rows in them.
	ResIDPinName = "res-id"
)

// pinSelfMaps pins guard_exe_actions/guard_exe_events when a self binary is registered and pinning
// is on, solely for `daemon --lockdown`, which widens a file-vault's self access from a fresh
// process to lock it after a crash. Failure is logged, not fatal (only that recovery path is lost).
func (g *Guard) pinSelfMaps() {
	if g.pinPrefix == "" || !g.selfKeySet {
		return
	}
	// The whitelist maps are shared, so they are pinned once by the engine; what is per-resource is
	// only which slot this guard owns.
	sharedEngine.pinSharedMaps()

	m, err := cilium.NewMap(&cilium.MapSpec{
		Type:       cilium.Array,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1,
	})
	if err != nil {
		log.Errorf("guard %s: creating resource-id map failed (%v) — 'daemon --lockdown' will not be "+
			"able to widen self-access here if the daemon crashes while this resource is unlocked",
			g.path, err)
		return
	}
	defer m.Close()
	if err := m.Put(uint32(0), g.resID); err != nil {
		log.Errorf("guard %s: recording resource id failed: %v", g.path, err)
		return
	}
	if err := m.Pin(g.pinPrefix + ResIDPinName); err != nil {
		log.Errorf("guard %s: pinning %s failed (%v) — 'daemon --lockdown' will not be able to widen "+
			"self-access here if the daemon crashes while this resource is unlocked",
			g.path, ResIDPinName, err)
	}
}

// unpinSelfMaps removes the map pins on clean shutdown (pins outlive only unclean death).
// Best-effort; CleanupStalePins retires leftovers.
func (g *Guard) unpinSelfMaps() {
	if g.pinPrefix == "" || !g.selfKeySet {
		return
	}
	for _, name := range []string{ResIDPinName} {
		if err := os.Remove(g.pinPrefix + name); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Errorf("guard %s: removing pinned %s failed (%v) — remove it manually under %s*",
				g.path, name, err, g.pinPrefix)
		}
	}
}

func guardLSMHooks(o *GuardObjects) []struct {
	prog *cilium.Program
	hook string
} {
	return []struct {
		prog *cilium.Program
		hook string
	}{
		{o.GuardFileOpen, "file_open"},
		{o.GuardFilePermission, "file_permission"},
		{o.GuardFileTruncate, "file_truncate"},
		{o.GuardMmapFile, "mmap_file"},
		{o.GuardPathUnlink, "path_unlink"},
		{o.GuardPathRename, "path_rename"},
		{o.GuardPathSymlink, "path_symlink"},
		{o.GuardPathLink, "path_link"},
		{o.GuardPathMkdir, "path_mkdir"},
		{o.GuardPathTruncate, "path_truncate"},
		{o.GuardInodeSetattr, "inode_setattr"},
		{o.GuardInodeSetxattr, "inode_setxattr"},
		{o.GuardInodeRemovexattr, "inode_removexattr"},
		{o.GuardPathMknod, "path_mknod"},
		{o.GuardPathRmdir, "path_rmdir"},
		{o.GuardInodePermission, "inode_permission"},
		{o.GuardInodeGetattr, "inode_getattr"},
		{o.GuardInodeGetxattr, "inode_getxattr"},
		{o.GuardInodeListxattr, "inode_listxattr"},
		{o.GuardInodeReadlink, "inode_readlink"},
		{o.GuardSbMount, "sb_mount"},
		{o.GuardPtraceAccessCheck, "ptrace_access_check"},
		{o.GuardBprmCheckSecurity, "bprm_check_security"},
		{o.GuardTaskFree, "task_free"},
		{o.GuardInodeFree, "inode_free_security"},
		{o.GuardBprmCommitted, "bprm_committed_creds"},
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
		if err := g.putExeAction(inodeKey, action); err != nil {
			return fmt.Errorf("storing exe action for %s: %w", b.Path, err)
		}
		g.mu.Lock()
		g.deployed[canonicalBinaryPath(b.Path)] = inodeKey
		g.mu.Unlock()
	}
	return nil
}

// addAllowRootBinary stores the owner's executable as GUARD_ALLOW_ROOT (uid-0 gated in BPF) with an
// optional event mask. Never touches g.binaries, so re-sync/deferred resolution can't re-register
// it as a plain GUARD_ALLOW.
func (g *Guard) addAllowRootBinary(b BinaryEntry, events []ebpf.EventType) error {
	// ModeReadOnly reuses the whitelist path for modifying ops, so the root-gated self allow
	// applies. The event mask is only consulted by the whitelist branch, so register the self
	// binary maskless there (as the daemon does).
	if g.mode != ModeWhitelist && g.mode != ModeReadOnly {
		return fmt.Errorf("self allow is only supported in whitelist / read-only mode")
	}

	dev, ino, err := ebpf.StatInode(b.Path)
	if err != nil {
		return fmt.Errorf("storing self exe action for %s: %w", b.Path, err)
	}
	key := GuardInodeKey{Dev: dev, Ino: ino}

	if err := g.putExeAction(key, uint8(GUARD_ALLOW_ROOT)); err != nil {
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
		if err := g.objs().GuardExeEvents.Put(g.resKey(key), mask); err != nil {
			return fmt.Errorf("storing self exe events for %s: %w", b.Path, err)
		}
	}
	return nil
}

// editGrantEvents is the widened self mask for a live edit-protected session: the read set plus
// what the editor's atomic save does (create temp, write, rename, new file/dir, delete, chmod).
// SYMLINK, HARDLINK, MMAP are excluded; the post-edit audit flags any.
var editGrantEvents = []ebpf.EventType{
	ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat,
	ebpf.EventWrite, ebpf.EventDelete, ebpf.EventRename,
	ebpf.EventMkdir, ebpf.EventMknod, ebpf.EventAttr,
}

// GrantSelfEditAccess widens the root-gated self binary's mask (editGrantEvents) for an
// authenticated live edit session. Touches only the self key, never the whitelist; BPF-map only
// (gone on restart). Pair with RevokeSelfEditAccess. Whitelist/read-only mode only.
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
	if err := g.objs().GuardExeEvents.Put(g.resKey(key), mask); err != nil {
		return fmt.Errorf("guard %s: widening self edit mask: %w", g.path, err)
	}
	g.selfGranted = true
	log.Warnf("guard %s: edit-protected write access GRANTED to the app-listener binary (uid 0 only) for this resource", g.path)
	return nil
}

// RevokeSelfEditAccess restores the baseline self mask. Idempotent; safe if the grant never took.
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
		if err := g.objs().GuardExeEvents.Delete(g.resKey(key)); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			restoreErr = fmt.Errorf("guard %s: clearing self edit mask: %w", g.path, err)
		}
	} else {
		mask, err := eventMask(base)
		if err != nil {
			return err
		}
		if err := g.objs().GuardExeEvents.Put(g.resKey(key), mask); err != nil {
			restoreErr = fmt.Errorf("guard %s: restoring self edit mask: %w", g.path, err)
		}
	}
	if restoreErr == nil && g.selfGranted {
		log.Infof("guard %s: edit-protected write access REVOKED — self binary back to read-only", g.path)
	}
	g.selfGranted = false
	return restoreErr
}

// vaultAccessEvents is the minimal self mask for the file-vault's in-place unlock/lock on its own
// path: open, read, stat, write (file_truncate checks EVENT_WRITE). Narrower than editGrantEvents:
// no delete/rename/mkdir/mknod/attr, since the transform only truncates and rewrites the same file.
var vaultAccessEvents = []ebpf.EventType{
	ebpf.EventOpen, ebpf.EventRead, ebpf.EventWrite, ebpf.EventStat,
}

// WithSelfVaultAccess widens the self mask to vaultAccessEvents for the duration of fn, then always
// restores the baseline. Needed for single-file resources, whose vault must rewrite its own bytes
// (the baseline open/read/stat mask would deny the daemon by its own guard).
//
// Deliberately separate from GrantSelfEditAccess so a boot-time vault unlock is never mistaken for,
// or widens, a live-edit grant. Whitelist/read-only mode only; returns fn's error, or an error if
// no self binary is registered.
func (g *Guard) WithSelfVaultAccess(fn func() error) error {
	g.selfGrantMu.Lock()
	defer g.selfGrantMu.Unlock()

	g.mu.Lock()
	key, ok := g.selfKey, g.selfKeySet
	base := append([]ebpf.EventType(nil), g.selfEvents...)
	g.mu.Unlock()
	if !ok {
		return fmt.Errorf("guard %s: no self binary registered — cannot grant vault access", g.path)
	}

	mask, err := eventMask(vaultAccessEvents)
	if err != nil {
		return err
	}
	if err := g.objs().GuardExeEvents.Put(g.resKey(key), mask); err != nil {
		return fmt.Errorf("guard %s: widening self vault mask: %w", g.path, err)
	}
	defer func() {
		var restoreErr error
		if len(base) == 0 {
			if derr := g.objs().GuardExeEvents.Delete(g.resKey(key)); derr != nil && !errors.Is(derr, cilium.ErrKeyNotExist) {
				restoreErr = derr
			}
		} else if baseMask, merr := eventMask(base); merr != nil {
			restoreErr = merr
		} else if perr := g.objs().GuardExeEvents.Put(g.resKey(key), baseMask); perr != nil {
			restoreErr = perr
		}
		if restoreErr != nil {
			log.Errorf("guard %s: restoring self baseline mask after vault access: %v — self access may be left WIDENED until the next restart", g.path, restoreErr)
		}
	}()

	return fn()
}

// WithPinnedSelfVaultAccess is WithSelfVaultAccess for a process with no live *Guard: `daemon
// --lockdown` (ExecStopPost) locking a file-vault after its daemon died. It works on the pinned
// guard_exe_actions/guard_exe_events maps left by pinSelfMaps. pinPrefix must be
// base+gen+resource-hash (PinPrefix); the caller recovers gen (see pinstate.go).
//
// Security invariants:
//   - The self key is always computed internally from /proc/self/exe, never a parameter.
//   - It only WIDENS an existing GUARD_ALLOW_ROOT entry for that key in the pinned actions map; if
//     the pin is missing/unreadable or the key isn't GUARD_ALLOW_ROOT it errors without touching
//     guard_exe_events.
//   - The mask is read first and restored via defer, even if fn fails.
//   - Worst case (restore fails): pinned maps are per-generation and retired by CleanupStalePins on
//     next start (RestartSec=2s).
func WithPinnedSelfVaultAccess(pinPrefix, sharedPrefix string, fn func() error) error {
	if pinPrefix == "" || sharedPrefix == "" {
		return errors.New("no pin prefix given — cannot widen a self grant that was never pinned")
	}

	dev, ino, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return fmt.Errorf("resolving own executable: %w", err)
	}

	// The whitelist maps are shared by every resource and keyed by (res_id, inode), so the
	// resource's own slot, pinned beside it, selects this resource's rows.
	resID, err := readPinnedResID(pinPrefix)
	if err != nil {
		return err
	}
	key := GuardResInodeKey{ResId: resID, Ino: GuardInodeKey{Dev: dev, Ino: ino}}

	if allowErr := requirePinnedSelfAllowRoot(sharedPrefix, pinPrefix, key); allowErr != nil {
		return allowErr
	}

	events, err := cilium.LoadPinnedMap(sharedPrefix+ExeEventsPinName, nil)
	if err != nil {
		return fmt.Errorf("loading pinned %s: %w", ExeEventsPinName, err)
	}
	defer events.Close()

	var baseMask uint32
	hadBase := events.Lookup(key, &baseMask) == nil // absent baseline means "all events allowed" (whitelist mode's convention)

	widened, err := eventMask(vaultAccessEvents)
	if err != nil {
		return err
	}
	if putErr := events.Put(key, widened); putErr != nil {
		return fmt.Errorf("widening pinned self vault mask at %s: %w", pinPrefix, putErr)
	}
	defer func() {
		var restoreErr error
		if !hadBase {
			if derr := events.Delete(key); derr != nil && !errors.Is(derr, cilium.ErrKeyNotExist) {
				restoreErr = derr
			}
		} else if perr := events.Put(key, baseMask); perr != nil {
			restoreErr = perr
		}
		if restoreErr != nil {
			log.Errorf("lockdown: restoring pinned self baseline mask at %s failed (%v) — left WIDENED; "+
				"the next daemon start's CleanupStalePins will retire this pin regardless", pinPrefix, restoreErr)
		}
	}()

	return fn()
}

// addBinaryEvents stores per-binary event bitmasks in guard_exe_events (explicit lists only; a
// missing entry means all events allowed).
func (g *Guard) addBinaryEvents(binaries []BinaryEntry, events map[string][]ebpf.EventType) error {
	for _, b := range binaries {
		types, ok := events[b.Path]
		// Empty list = no restriction, valid in every mode.
		if !ok || len(types) == 0 {
			continue
		}
		if g.mode != ModeWhitelist {
			return fmt.Errorf("per-binary event restrictions are only supported in whitelist mode (binary %s)", b.Path)
		}
		mask, err := eventMask(types)
		if err != nil {
			return fmt.Errorf("invalid event mask for binary %s: %w", b.Path, err)
		}
		dev, ino, err := ebpf.StatInode(b.Path)
		if err != nil {
			return fmt.Errorf("cannot stat binary %s for event mask: %w", b.Path, err)
		}
		if err := g.objs().GuardExeEvents.Put(g.resKey(GuardInodeKey{Dev: dev, Ino: ino}), mask); err != nil {
			return fmt.Errorf("storing exe events for %s: %w", b.Path, err)
		}
	}
	return nil
}

// resolveDeferred hashes every deferred rule now readable; returns entries, event lists (by
// canonical path) and still-unreadable rules. Never mutates the guard.
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

// ResolvePendingBinaries retries entries deferred while their resource was locked; run only after
// unlock + populateInodes. Readable rules are hashed, canonicalized and written to the whitelist
// maps; the rest stay deferred (logged), never weakening protection.
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

// putBinaryKey writes a binary's allow/block action to guard_exe_actions and, in whitelist mode,
// its event mask to guard_exe_events.
func (g *Guard) putBinaryKey(key GuardInodeKey, path string, events map[string][]ebpf.EventType) error {
	action := uint8(GUARD_BLOCK)
	if g.mode == ModeWhitelist {
		action = uint8(GUARD_ALLOW)
	}
	if err := g.putExeAction(key, action); err != nil {
		return fmt.Errorf("storing exe action for %s: %w", path, err)
	}
	if g.mode == ModeWhitelist {
		if types, ok := events[path]; ok && len(types) > 0 {
			mask, err := eventMask(types)
			if err != nil {
				return fmt.Errorf("invalid event mask for binary %s: %w", path, err)
			}
			if err := g.objs().GuardExeEvents.Put(g.resKey(key), mask); err != nil {
				return fmt.Errorf("storing exe events for %s: %w", path, err)
			}
		}
	}
	return nil
}

// retryDeferredBinaries admits deferred rules that became readable (map entries + bookkeeping).
// Shared by ResolvePendingBinaries and ReSyncBinaries.
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
	// Pin the freshly resolved binaries so in-place replacement detection covers them.
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

// ReSyncBinaries re-stats whitelisted binaries and rewrites inode-keyed entries after in-place
// replacement. Stale keys are never deleted (running pre-replacement processes stay admitted;
// vanished paths keep their key). Still-deferred rules are retried.
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

// writeGuardConfig stores this resource's mode, recursion flag and depth limit in its slot,
// leaving it inactive: NewGuard activates it once the inode map is populated.
func (g *Guard) writeGuardConfig(modeKey uint64) error {
	if g.depth < 0 {
		return fmt.Errorf("invalid depth %d", g.depth)
	}
	recursiveVal := uint64(0)
	if g.recursive {
		recursiveVal = 1
	}
	cfg := GuardResConfig{
		Mode:      modeKey,
		Recursive: recursiveVal,
		Depth:     uint64(g.depth),
	}
	if putErr := g.objs().GuardResConfig.Put(g.resID, cfg); putErr != nil {
		return fmt.Errorf("writing resource slot %d: %w", g.resID, putErr)
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

	// Watch root (dev, ino) feeds the BPF root-confinement check (guard_config[3..4]): an entry
	// only guards an access whose dentry chain reaches this root, so inode reuse elsewhere can't
	// deny unrelated files.
	rootDev, rootIno, rootStatErr := ebpf.StatInode(g.path)
	if rootStatErr != nil {
		return fmt.Errorf("stating guard root %s: %w", g.path, rootStatErr)
	}
	if rootErr := g.updateRootKey(GuardInodeKey{Dev: rootDev, Ino: rootIno}); rootErr != nil {
		return rootErr
	}

	_, statErr := os.Stat(g.path)
	if statErr != nil {
		return fmt.Errorf("stating guarded path %s: %w", g.path, statErr)
	}
	// Store the path's own inode before attach so the ancestor walk covers the tree from the moment
	// hooks go live (statable even while an encrypted tree is locked).
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

	// Register the watch root for path_symlink's target check (guard_path: path -> resource).
	if putErr := g.objs().GuardPath.Put(g.pathKey(), g.resID); putErr != nil {
		return fmt.Errorf("storing guarded path: %w", putErr)
	}

	// Block the backing block device: debugfs(8)/dd can read /dev/sda1 etc. without open()ing the
	// guarded path, bypassing VFS access control.
	if err := g.addBackingBlockDevice(); err != nil {
		log.Warnf("backing block device detection: %v", err)
	}
	if err := g.addFsDeviceGate(); err != nil {
		log.Warnf("filesystem device gate: %v", err)
	}

	return nil
}

// addFsDeviceGate records the path's fs device in guard_fs_sbdevs so the ancestor walk skips other
// filesystems in one lookup (st_dev == i_sb->s_dev on all fs types, incl. anon
// tmpfs/overlayfs/btrfs). Must be unconditional: an empty map disables the walk and exposes deeper
// content.
func (g *Guard) addFsDeviceGate() error {
	var s syscall.Stat_t
	if err := syscall.Stat(g.path, &s); err != nil {
		return fmt.Errorf("stating guarded path: %w", err)
	}

	major := unix.Major(s.Dev)
	dev := uint64(major)<<20 | uint64(unix.Minor(s.Dev))
	var val uint8 = 1
	if err := g.objs().GuardFsSbdevs.Put(dev, val); err != nil {
		return fmt.Errorf("storing filesystem device %d:%d in map: %w", major, unix.Minor(s.Dev), err)
	}
	return nil
}

// BackingDevice returns the guard_fs_devices key (dev_t, major<<20|minor) of path's backing block
// device; false for tmpfs, overlayfs, procfs and other major-0 pseudo-filesystems.
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

// addBackingBlockDevice populates guard_fs_devices so raw opens of the hosting block device(s) are
// blocked: the WithBackingDevices set if given (empty = opt out), else the device of the watched
// path.
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
	if len(rdevs) > 0 {
		sharedEngine.claimRawDevices(g)
	}
	var val uint8 = 1
	for _, rdev := range rdevs {
		if err := g.objs().GuardFsDevices.Put(rdev, val); err != nil {
			return fmt.Errorf("storing backing block device %d:%d in map: %w",
				rdev>>20, rdev&0xFFFFF, err)
		}
		log.Infof("blocking raw access to backing block device %d:%d", rdev>>20, rdev&0xFFFFF)
	}
	return nil
}

// updateRootKey re-anchors the watch-root identity to newKey in guard_config[3..4] (root_in_chain)
// and g.rootKey. Called from populateMaps and from SweepInodes when a single-file root is
// recreated.
func (g *Guard) updateRootKey(newKey GuardInodeKey) error {
	cfg, err := g.resConfig()
	if err != nil {
		return err
	}
	cfg.RootDev = newKey.Dev
	cfg.RootIno = newKey.Ino
	if putErr := g.objs().GuardResConfig.Put(g.resID, cfg); putErr != nil {
		return fmt.Errorf("setting watch root in resource slot %d: %w", g.resID, putErr)
	}
	g.mu.Lock()
	g.rootKey = newKey
	g.mu.Unlock()
	return nil
}

// treeInode returns the inode a tree walk records for path. Whitelist guards follow symlinks (also
// guarding an in-tree link's target). Read-only guards (`lib_dir`) must not: trees like Steam's
// pressure-vessel link into host /usr/lib (would guard system libs, or fail on dangling links and
// trip the locked-vault heuristic); the link's own entry belongs to the tree.
func (g *Guard) treeInode(path string) (dev, ino uint64, err error) {
	if g.mode == ModeReadOnly {
		return ebpf.LstatInode(path)
	}
	return ebpf.StatInode(path)
}

func (g *Guard) addInode(path string) error {
	dev, ino, err := g.treeInode(path)
	if err != nil {
		return err
	}

	key := GuardInodeKey{
		Dev: dev,
		Ino: ino,
	}

	if err := g.objs().GuardInodes.Put(key, g.resID); err != nil {
		return fmt.Errorf("adding inode %s to map: %w", path, err)
	}
	return nil
}

// PopulateInodes fills guard_inodes with every file/dir under the path; run once the path is
// readable (after unlock; the root inode is already mapped, so the ancestor walk protects the tree
// meanwhile). Tolerant (see walkInodes): degraded coverage weakens rename/unlink/mmap precision,
// not open/read.
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

// reanchorRoot moves the root identity from oldKey to newKey: rescan via the caller-chosen scan
// (addInode for a file, scanDirInodes for a dir), move the guard_config[3..4] anchor and g.rootKey,
// then evict oldKey. A stale anchor would stop guarding the recreated resource and deny whatever
// reuses the freed inode number.
func (g *Guard) reanchorRoot(oldKey, newKey GuardInodeKey, rescan func() error) error {
	if err := rescan(); err != nil {
		return err
	}
	if err := g.updateRootKey(newKey); err != nil {
		return err
	}
	if delErr := g.objs().GuardInodes.Delete(oldKey); delErr != nil && !errors.Is(delErr, cilium.ErrKeyNotExist) {
		log.Warnf("guard %s: evicting the stale root inode %+v: %v", g.path, oldKey, delErr)
	}
	return nil
}

// sweepDirRootRecreated checks whether a directory root itself was deleted and recreated (new
// inode: fscrypt migration, backup restore, app rebuild), not just a top-level entry change. Runs
// ahead of the mtime/dirRescanMinGap throttle, which must never delay it. Reports whether it
// handled this tick (repaired), so SweepInodes skips the mtime re-scan.
func (g *Guard) sweepDirRootRecreated(info os.FileInfo) (bool, error) {
	dev, ino, statErr := ebpf.StatInode(g.path)
	if statErr != nil {
		return false, statErr
	}
	newKey := GuardInodeKey{Dev: dev, Ino: ino}
	g.mu.Lock()
	oldKey := g.rootKey
	g.mu.Unlock()
	if newKey == oldKey {
		return false, nil
	}

	// The whole subtree is new (new inodes even for same names): re-scan fully, including the
	// root's entry, before moving the anchor. Stale non-root entries age out via inode GC.
	if err := g.reanchorRoot(oldKey, newKey, func() error { return g.scanDirInodes(g.path, 0) }); err != nil {
		return false, err
	}
	now := time.Now()
	g.mu.Lock()
	g.sweepRootMtime = info.ModTime()
	g.sweepLastFull = now
	g.mu.Unlock()
	return true, nil
}

// SweepInodes is the cheap periodic refresh of guard_inodes (a full re-walk per tick cost ~14%
// daemon CPU):
//   - single-file root: re-map only when its inode changed;
//   - directory root: first, unthrottled, check whether the directory itself was recreated and
//     re-anchor; then re-scan the top level only if the dir mtime moved. Deeper additions are
//     covered by the path_mkdir hook and the ancestor walk.
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
		newKey := GuardInodeKey{Dev: dev, Ino: ino}
		g.mu.Lock()
		oldKey := g.rootKey
		g.mu.Unlock()
		if newKey == oldKey {
			return nil
		}
		// Root was deleted and recreated with a new inode: reanchorRoot moves protection (a stale
		// anchor is dangerous, see its doc).
		return g.reanchorRoot(oldKey, newKey, func() error { return g.addInode(g.path) })
	}

	if handled, recreateErr := g.sweepDirRootRecreated(info); recreateErr != nil || handled {
		return recreateErr
	}

	mtime := info.ModTime()
	g.mu.Lock()
	// Re-walk only when a top-level entry changed, and no more often than dirRescanMinGap; new
	// subtrees meanwhile are covered by path_mkdir and the ancestor walk.
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

// addErrKind classifies a failed addInode stat for the walk.
type addErrKind int

const (
	addErrOther   addErrKind = iota // unclassified failure: entry skipped
	addErrMissing                   // ENOENT/ENOTDIR: entry vanished mid-walk
	addErrFull                      // E2BIG: inode map full
)

// classifyAddErr maps an add error to walk semantics: E2BIG degrades coverage (open/read survive
// via the ancestor walk); anything else is logged and skipped.
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

// walkInodes walks dir depth-first, calling add for dir and every entry: ENOENT/ENOTDIR and other
// per-entry stat errors are skipped with a warning; E2BIG stops the walk with degraded coverage
// (open/read survive via the ancestor walk); a dir whose every entry fails to stat is almost
// certainly locked fscrypt (error).
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

// walkEntries visits one directory level; errors only if a nested recursion fails or every entry
// failed to stat (locked-encrypted-tree signature).
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

// walkLiveEntries mirrors walkInodes but is strict at the root: reconcile must not mistake a
// momentarily unreadable root (e.g. mid fscrypt lock/unlock) for an empty tree and evict
// everything, so a root add() failure is a hard error. Below the root, per-entry tolerance is
// unchanged (a vanished file is simply not live).
func walkLiveEntries(root string, recursive bool, depthLimit int, add func(string) error) error {
	if err := add(root); err != nil {
		return fmt.Errorf("collecting root inode %s: %w", root, err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("reading %s: %w", root, err)
	}
	return walkEntries(root, entries, recursive, depthLimit, 0, add)
}

// liveInodeKeys stats the guarded tree now, for ReconcileInodes to diff against guard_inodes.
// Errors whenever the walk can't be trusted as complete, collecting nothing.
func (g *Guard) liveInodeKeys() (map[GuardInodeKey]struct{}, error) {
	live := make(map[GuardInodeKey]struct{})
	collect := func(path string) error {
		// Same stat flavor as addInode: populate and reconcile must agree on each path's inode.
		dev, ino, err := g.treeInode(path)
		if err != nil {
			return err
		}
		live[GuardInodeKey{Dev: dev, Ino: ino}] = struct{}{}
		return nil
	}

	info, err := os.Stat(g.path)
	if err != nil {
		return nil, fmt.Errorf("stating guarded path %s: %w", g.path, err)
	}
	if !info.IsDir() {
		if err := collect(g.path); err != nil {
			return nil, fmt.Errorf("collecting root inode %s: %w", g.path, err)
		}
		return live, nil
	}
	if err := walkLiveEntries(g.path, g.recursive, g.depth, collect); err != nil {
		return nil, err
	}
	return live, nil
}

// SnapshotTaintedPIDs returns every tgid in guard_tainted_pids, to carry taint across a reload
// (RestoreTaintedPIDs). Returns an error rather than a partial set.
func (g *Guard) SnapshotTaintedPIDs() ([]uint32, error) {
	var pids []uint32
	var key uint32
	var val uint32
	it := g.objs().GuardTaintedPids.Iterate()
	for it.Next(&key, &val) {
		// Only this resource's own tainted processes carry over; resGlobal marks a process holding
		// several resources' content, which the surviving guards re-stamp themselves.
		if val == g.resID {
			pids = append(pids, key)
		}
	}
	if err := it.Err(); err != nil {
		return nil, fmt.Errorf("iterating tainted pids for %s: %w", g.path, err)
	}
	return pids, nil
}

// RestoreTaintedPIDs stamps pids into this (fresh) guard's guard_tainted_pids, keeping memory-read
// protection for processes tainted under the guard it replaces on reload. Dead pids are harmless
// (guard_task_free removes them).
func (g *Guard) RestoreTaintedPIDs(pids []uint32) error {
	for _, p := range pids {
		if err := g.objs().GuardTaintedPids.Put(p, g.resID); err != nil {
			return fmt.Errorf("restoring tainted pid %d for %s: %w", p, g.path, err)
		}
	}
	return nil
}

// ReconcileInodes deletes guard_inodes entries whose (dev, ino) no longer belongs to anything under
// the tree. Everything else only adds, so stale entries accumulate; freed inode numbers get reused,
// and a stale entry colliding with an unrelated file can false-DENY it (root_in_chain can fail
// closed beyond its walk bound).
//
// Over-evicting is unsafe (a transiently unreadable file would be unguarded), so it deletes nothing
// unless liveInodeKeys returns a complete, error-free walk, and never the root key.
func (g *Guard) ReconcileInodes() error {
	live, err := g.liveInodeKeys()
	if err != nil {
		return fmt.Errorf("inode GC walk incomplete for %s, skipping this cycle: %w", g.path, err)
	}

	g.mu.Lock()
	rootKey := g.rootKey
	g.mu.Unlock()

	var stale []GuardInodeKey
	var key GuardInodeKey
	var val uint32
	it := g.objs().GuardInodes.Iterate()
	for it.Next(&key, &val) {
		// The map is shared by every resource; another resource's rows are not ours to judge.
		if val != g.resID {
			continue
		}
		if key == rootKey {
			continue
		}
		if _, ok := live[key]; ok {
			continue
		}
		stale = append(stale, key)
	}
	if iterErr := it.Err(); iterErr != nil {
		return fmt.Errorf("iterating guard_inodes for %s: %w", g.path, iterErr)
	}

	for _, k := range stale {
		if delErr := g.objs().GuardInodes.Delete(k); delErr != nil && !errors.Is(delErr, cilium.ErrKeyNotExist) {
			return fmt.Errorf("evicting stale inode from guard_inodes for %s: %w", g.path, delErr)
		}
	}
	if len(stale) > 0 {
		log.Infof("guard %s: evicted %d stale inode(s) no longer present on disk", g.path, len(stale))
	}
	return nil
}

func (g *Guard) Events() <-chan GuardEvent {
	return g.events
}

func (g *Guard) Start() error {
	// Events arrive from the shared engine's reader, which routes each one by resource id.
	log.Infof("guard started \u2014 guarding: %s", g.path)

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

// startDegradeWatch logs otherwise-silent BPF-side degradation (kernel map allocation failures
// reduce inode-discovery coverage or disable ptrace protection). Observability only.
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
					if err := g.objs().GuardDegrade.Lookup(slot, &count); err != nil {
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

// binaryVerifyState is the verifier's pinned identity of an admitted binary: inode key, content
// hash, and the stat fingerprint it was computed from (unchanged binaries aren't re-hashed).
type binaryVerifyState struct {
	key     GuardInodeKey
	stat    ebpf.BinaryStat
	hash    [32]byte
	hashed  bool // a real content hash has been taken at least once
	demoted bool
}

// binaryHashVerifyInterval is short because an unchanged binary costs one stat and no read
// (verifyBinaryHashesOnce).
const binaryHashVerifyInterval = 10 * time.Second

// sharedBinaryHashes dedupes content hashing across all guards in the process (one large app blob
// is admitted by ~10 guards). Keyed by inode; the stored fingerprint gates reuse.
var (
	sharedHashMu sync.Mutex
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
)

type sharedHashEntry struct {
	stat ebpf.BinaryStat
	hash [32]byte
}

// hashBinaryShared returns path's hash, reusing one already computed for the same
// inode+fingerprint.
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

// startBinaryHashVerifier pins each admitted binary's content hash and periodically re-verifies it:
// identity is the exe inode alone, so anyone who can write an admitted binary outside the guarded
// tree could replace it in place and inherit the allow. A same-inode hash change demotes the entry
// to GUARD_BLOCK until the inode changes (a package-style replacement is re-admitted by
// ReSyncBinaries/reload and re-pinned next tick).
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

// verifyBinaryHashesOnce re-hashes binaries whose stat fingerprint moved and demotes same-inode
// replacements. States are owned by this goroutine (extended by ReSyncBinaries under g.mu), so no
// extra locking.
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

		// Same inode and size/mtime/ctime: an in-place overwrite always bumps mtime/ctime, so
		// unchanged. No read.
		if st.hashed && freshKey == st.key && fp == st.stat {
			continue
		}

		hash, err := hashBinaryShared(canonical, fp)
		if err != nil {
			continue // unreadable right now: keep the current decision
		}

		if freshKey != st.key {
			// Inode changed: a replacement flow (ReSyncBinaries, reload) owns re-admission; adopt
			// the new identity.
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
			// Same inode, different content: replaced in place. Demote to GUARD_BLOCK (fail
			// closed).
			if putErr := g.putExeAction(st.key, uint8(GUARD_BLOCK)); putErr != nil {
				log.Errorf("guard %s: demoting in-place replaced binary %s: %v", g.path, canonical, putErr)
				continue
			}
			st.demoted = true
			log.Errorf("guard %s: whitelisted binary %s was modified in place (inode unchanged, hash changed) \u2014 whitelist entry demoted to BLOCK", g.path, canonical)
		}
	}
}

// parseGuardEvent decodes one ringbuf record, returning the event and the resource that produced
// it. The engine calls it once per record and routes the result.
func parseGuardEvent(raw []byte) (*GuardEvent, uint32, bool) {
	var be bpfGuardEvent
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, &be); err != nil {
		log.Errorf("decode guard event: %v", err)
		return nil, 0, false
	}

	fe := be.toFileEvent()
	fe.Timestamp = time.Now().UnixNano()

	ge := &GuardEvent{
		FileEvent: fe,
		Blocked:   be.Blocked != 0,
		RawDevice: be.Reason == guardReasonRawDevice,
		Process:   processGateLabel(be.Reason),
	}
	if ge.Process != "" {
		// The kernel sends the other task's comm in path and its tgid in fd.
		ge.Path = fmt.Sprintf("pid=%d comm=%s", be.FD, ge.Path)
		if ge.Dest != "" {
			// READ = /proc metadata only; ATTACH = memory access.
			ge.Path += " mode=" + ge.Dest
			ge.Dest = ""
		}
	}
	return ge, be.ResID, true
}

// dispatch delivers one decoded event to this guard's consumer.
//
// Never blocks: one reader now serves every resource, so waiting on a slow or absent consumer
// would stall event delivery for ALL of them (the daemon runs dozens). A full channel drops the
// event and counts it, which costs only telemetry — enforcement already happened in the kernel,
// synchronously, and the BPF ringbuf itself drops on overflow for the same reason.
func (g *Guard) dispatch(ev *GuardEvent) {
	// comm is telemetry: the spoof warning is diagnostic only, never enforcement (BPF decisions
	// key on exe inode).
	if g.dropAllowed && !ev.Blocked {
		return
	}
	local := *ev
	g.checkCommSpoof(&local)

	select {
	case g.events <- local:
	case <-g.done:
	default:
		if local.Blocked {
			logBacklogDenial(g.path, &local) // audit data: never lost behind a flood of allowed events
			return
		}
		if dropped := g.eventsDropped.Add(1); dropped == 1 || dropped%1000 == 0 {
			log.Warnf("guard %s: event consumer is not keeping up, %d allowed event(s) dropped — "+
				"enforcement is unaffected, only reporting", g.path, dropped)
		}
	}
}

// logBacklogDenial prints ev in the daemon's DAEMON DENIED layout (scripts/trace-app-libs.sh
// parses it); commFullPath is best-effort telemetry, "~" when unresolved.
func logBacklogDenial(resource string, ev *GuardEvent) {
	op := ev.Type.String()
	if ev.Process != "" {
		op = ev.Process
	}
	exe := "~"
	if target, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ev.PID)); err == nil {
		exe = logging.SanitizeText(target)
	}
	fmt.Fprintf(os.Stderr, "<4>DAEMON DENIED  op=%s  comm=%s  commFullPath=%s  pid=%d  uid=%d  resource=%s  path=%s\n",
		op, logging.SanitizeText(ev.Comm), exe, ev.PID, ev.UID, logging.SanitizeText(resource),
		logging.SanitizeText(ev.Path))
}

// checkCommSpoof warns when an event's comm claims a guarded binary's name but the real binary
// differs (also for blocked events). Binaries that rename their threads (Chromium, Bun) are normal
// and don't warn.
func (g *Guard) checkCommSpoof(ge *GuardEvent) {
	// Cheap string gate first: this runs per event and most comms don't collide; only a collision
	// warrants a /proc lookup.
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

// commMatchesGuardedBinary reports whether comm claims a guarded binary's name; the kernel
// truncates comm to 15 bytes (TASK_COMM_LEN), so both sides are truncated.
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

// cleanup retires this resource. The shared LSM links stay attached while any other guard holds
// the engine, so a reload that swaps resources never drops enforcement.
func (g *Guard) cleanup() {
	sharedEngine.forgetResource(g.resID)
	g.dropResourceState()
	g.unpinSelfMaps()
	sharedEngine.release(g.resID)
}

// dropResourceState removes this resource's rows from the shared maps, so a slot reused later
// starts empty and no stale inode can be judged by another resource's policy.
func (g *Guard) dropResourceState() {
	if sharedEngine.objs.GuardInodes == nil {
		return // engine already torn down
	}
	if err := g.deleteInodesOfResource(); err != nil {
		log.Warnf("guard %s: clearing inode entries: %v", g.path, err)
	}
	if err := deleteResKeys(g.objs().GuardExeActions, g.resID); err != nil {
		log.Warnf("guard %s: clearing whitelist entries: %v", g.path, err)
	}
	if err := deleteResKeys(g.objs().GuardExeEvents, g.resID); err != nil {
		log.Warnf("guard %s: clearing event masks: %v", g.path, err)
	}
	// Only if still ours: a reload registers the same path for the replacement guard before this
	// one stops, and deleting it then would drop that guard's symlink-target check.
	var owner uint32
	if g.objs().GuardPath.Lookup(g.pathKey(), &owner) == nil && owner == g.resID {
		if err := g.objs().GuardPath.Delete(g.pathKey()); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			log.Warnf("guard %s: clearing watch path: %v", g.path, err)
		}
	}
}

// deleteInodesOfResource drops every guard_inodes row owned by this resource.
func (g *Guard) deleteInodesOfResource() error {
	var (
		key   GuardInodeKey
		val   uint32
		stale []GuardInodeKey
	)
	it := g.objs().GuardInodes.Iterate()
	for it.Next(&key, &val) {
		if val == g.resID {
			stale = append(stale, key)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	for i := range stale {
		if err := g.objs().GuardInodes.Delete(stale[i]); err != nil &&
			!errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

// deleteResKeys drops every row of a (res_id, inode)-keyed map belonging to res. It walks keys only
// (NextKey): the maps it serves have different value types, and a mismatched value makes Iterate
// fail on the first row — leaving the rows behind for whichever resource reuses the slot id.
func deleteResKeys(m *cilium.Map, res uint32) error {
	var stale []GuardResInodeKey
	var cur, next GuardResInodeKey
	var prev any // nil starts from the first key
	for {
		err := m.NextKey(prev, &next)
		if errors.Is(err, cilium.ErrKeyNotExist) {
			break
		}
		if err != nil {
			return err
		}
		if next.ResId == res {
			stale = append(stale, next)
		}
		cur = next
		prev = cur
	}
	for i := range stale {
		if err := m.Delete(stale[i]); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

type bpfGuardEvent struct {
	PID     uint32
	UID     uint32
	GID     uint32
	Type    uint32
	FD      uint32
	Blocked uint32
	Reason  uint32
	ResID   uint32
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

// resKey scopes an executable's inode to this resource: a binary whitelisted for one resource must
// not inherit that allow on another.
func (g *Guard) resKey(ik GuardInodeKey) GuardResInodeKey {
	return GuardResInodeKey{ResId: g.resID, Ino: ik}
}

// putExeAction records a whitelist/blacklist decision for this resource, keeping the shared
// cross-resource views in step when it is an allow.
func (g *Guard) putExeAction(ik GuardInodeKey, action uint8) error {
	if err := g.objs().GuardExeActions.Put(g.resKey(ik), action); err != nil {
		return err
	}
	if action == GUARD_ALLOW || action == GUARD_ALLOW_ROOT {
		return sharedEngine.noteAllow(g.resID, ik, action)
	}
	return sharedEngine.noteDeny(g.resID, ik)
}

// readPinnedResID recovers which resource slot the guard pinned at pinPrefix owned.
func readPinnedResID(pinPrefix string) (uint32, error) {
	m, err := cilium.LoadPinnedMap(pinPrefix+ResIDPinName, nil)
	if err != nil {
		return 0, fmt.Errorf("loading pinned %s (was this resource's guard ever pinned?): %w", ResIDPinName, err)
	}
	defer m.Close()

	var resID uint32
	if err := m.Lookup(uint32(0), &resID); err != nil {
		return 0, fmt.Errorf("reading pinned resource id at %s: %w", pinPrefix, err)
	}
	if resID == resGlobal || resID >= GuardMaxRes {
		return 0, fmt.Errorf("pinned resource id %d at %s is not a real resource slot", resID, pinPrefix)
	}
	return resID, nil
}

// requirePinnedSelfAllowRoot refuses unless the pinned whitelist already grants this exact key
// GUARD_ALLOW_ROOT: the lockdown path may only WIDEN an existing self grant, never create one.
func requirePinnedSelfAllowRoot(sharedPrefix, pinPrefix string, key GuardResInodeKey) error {
	actions, err := cilium.LoadPinnedMap(sharedPrefix+ExeActionsPinName, nil)
	if err != nil {
		return fmt.Errorf("loading pinned %s (was this daemon's guard ever pinned?): %w", ExeActionsPinName, err)
	}
	defer actions.Close()

	var action uint8
	if lookupErr := actions.Lookup(key, &action); lookupErr != nil {
		return fmt.Errorf("this process is not the registered self binary for the pinned guard at %s: %w", pinPrefix, lookupErr)
	}
	if action != GUARD_ALLOW_ROOT {
		return fmt.Errorf("refusing to widen: self key is not GUARD_ALLOW_ROOT in the pinned guard at %s (got action %d)", pinPrefix, action)
	}
	return nil
}

// pathKey is this watch root as a guard_path key: no trailing slash (path_symlink matches at '/'
// boundaries), zero-padded to MAX_PATH.
func (g *Guard) pathKey() [256]byte {
	var key [256]byte
	p := g.path
	if len(p) > 1 {
		p = strings.TrimRight(p, "/")
	}
	copy(key[:], p)
	return key
}
