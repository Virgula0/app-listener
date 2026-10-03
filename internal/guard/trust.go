package guard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// guard_trusted_files value flags — must match guard_trust.bpf.c.
const (
	trustedBinary uint8 = 1
	trustedLib    uint8 = 2
)

// trust_event.kind — must match guard_trust.bpf.c.
const (
	trustLibload    uint32 = 0
	trustWriteblock uint32 = 1
	trustPlant      uint32 = 2
	trustSuspect    uint32 = 3
)

// TrustGuard owns guard_trusted_files and the daemon-wide trusted-binary/library protections
// (guard_trust.bpf.c). Created and attached ONCE for the whole daemon, independent of the
// per-resource guards. Both protections (#1 writer attribution, #2 library-load allowlist) are
// always enforced once attached.
type TrustGuard struct {
	objs  GuardTrustObjects
	memfd bool // trust_memfd_alloc's target exists (trustSpec)
	links []link.Link
	rd    *ringbuf.Reader
	done  chan struct{}
	// ownerByPath: UpdaterPlan.Owners, for AllowReplacement.
	ownerMu     sync.Mutex
	ownerByPath map[string]uint64
	// guard_vouched_devs upkeep (trust_mounts.go).
	mountMu   sync.Mutex
	mountFd   int
	mountStop int
	mountDone chan struct{}
	vouched   int
	// mountUntracked: the watch died, so SyncMounts refuses to vouch anything again.
	mountUntracked bool
	// retiredMu guards retired: the rows forgetExe dropped for a freed key, so a replacement its
	// guards re-sync only after that prune can still adopt them (adoptRows). Userspace only: a
	// reused inode number gains nothing from them. Cleared by SetTrusted, which re-derives every row.
	retiredMu sync.Mutex
	retired   map[GuardTrustInodeKey]trustRows
}

// trustRows is one key's rows across guard_trusted_files and the bit maps (trustBitMaps order).
type trustRows struct {
	flags uint8
	bits  [4]uint64
	has   [4]bool
}

// maxRetiredRows bounds TrustGuard.retired; past it the oldest prune's rows are simply not kept
// (a later replacement then waits for a reload: stricter).
const maxRetiredRows = 1024

func (t *TrustGuard) trustBitMaps() [4]*cilium.Map {
	return [4]*cilium.Map{t.objs.GuardBinOwner, t.objs.GuardBinUpdaters, t.objs.GuardGlobWriters,
		t.objs.GuardLibdirUsers}
}

// trustEvent mirrors struct trust_event in guard_trust.bpf.c.
type trustEvent struct {
	PID  uint32
	UID  uint32
	Kind uint32
	Comm [16]byte
	Path [256]byte
}

// NewTrustGuard loads the trust BPF object. It does not attach anything or
// populate the set — call SetTrusted / SetGuardedDirs then Start.
func NewTrustGuard() (*TrustGuard, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("trust guard: removing memlock rlimit: %v", err)
	}
	shared, err := supersedeMaps()
	if err != nil {
		return nil, err
	}
	t := &TrustGuard{done: make(chan struct{}), mountFd: -1, mountStop: -1, vouched: -1}
	spec, memfd, err := trustSpec()
	if err != nil {
		return nil, err
	}
	t.memfd = memfd
	if err := spec.LoadAndAssign(&t.objs, &cilium.CollectionOptions{MapReplacements: shared}); err != nil {
		return nil, fmt.Errorf("loading trust BPF objects: %w", err)
	}
	return t, nil
}

// SetTrusted (re)populates guard_trusted_files from resolved binary and library paths. A path that
// is both carries both flags. Unresolvable paths are skipped with a warning. Idempotent: the map is
// synced to exactly the new set (drops entries for removed binaries, no stale over-trust). Uses
// syncMap (put-new-then-delete-stale) rather than clearing first, so a reload never opens a window
// where no binary is trusted — during that window the library allowlist and binary write-protection
// would enforce nothing (SIGHUP is the normal way the whitelist changes; see trustManager.reload).
func (t *TrustGuard) SetTrusted(binaries, libs []string) error {
	flags := make(map[GuardInodeKey]uint8)
	add := func(path string, flag uint8) {
		dev, ino, err := ebpf.StatConfined(path)
		if err != nil {
			log.Warnf("trust guard: skipping unresolvable %s: %v", path, err)
			return
		}
		flags[GuardInodeKey{Dev: dev, Ino: ino}] |= flag
	}
	for _, b := range binaries {
		add(b, trustedBinary)
	}
	for _, l := range libs {
		add(l, trustedLib)
	}
	for k := range flags {
		liftSuperseded(k)
	}
	t.retiredMu.Lock()
	t.retired = nil
	t.retiredMu.Unlock()
	if err := t.keepLive(flags); err != nil {
		return fmt.Errorf("keeping live trusted files: %w", err)
	}
	if err := syncMap(t.objs.GuardTrustedFiles, flags); err != nil {
		return fmt.Errorf("syncing trusted files: %w", err)
	}
	log.Infof("trust guard: %d trusted inode(s) loaded (%d binaries, %d libraries requested)",
		len(flags), len(binaries), len(libs))
	return nil
}

// keepLive adds to flags every trusted inode a guard still admits without a config path naming it:
// a system binary admitted live (admitSystemMatches), or a superseded one a process started before
// its update still runs. A reload rebuilds the set from paths, and losing the flag would lift the
// library allowlist from a process that holds its resource's secrets. Bit rows (owners, updaters,
// writers) are not kept: their bits are reassigned per reload.
func (t *TrustGuard) keepLive(flags map[GuardInodeKey]uint8) error {
	var k GuardInodeKey
	var f uint8
	it := t.objs.GuardTrustedFiles.Iterate()
	for it.Next(&k, &f) {
		if _, ok := flags[k]; ok {
			continue
		}
		if mark, ok := supersededMark(k); (ok && mark.Freed == 0) || (!ok && sharedEngine.admitsExe(k)) {
			flags[k] = f
		}
	}
	return it.Err()
}

// trusts reports whether k is a trusted file.
func (t *TrustGuard) trusts(k GuardInodeKey) bool {
	var f uint8
	return t.objs.GuardTrustedFiles.Lookup(GuardTrustInodeKey{Dev: k.Dev, Ino: k.Ino}, &f) == nil
}

// forgetExe drops every trust row keyed by k (supersede.go). guard_bin_origin is left to the
// kernel: every creation rewrites or drops its row (origin_record).
func (t *TrustGuard) forgetExe(k GuardInodeKey) error {
	key := GuardTrustInodeKey{Dev: k.Dev, Ino: k.Ino}
	t.retire(key)
	for _, m := range []*cilium.Map{t.objs.GuardTrustedFiles, t.objs.GuardBinOwner, t.objs.GuardBinUpdaters,
		t.objs.GuardGlobWriters, t.objs.GuardLibdirUsers} {
		if err := m.Delete(key); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

// retire stashes key's rows before forgetExe drops them, if it is a trusted file.
func (t *TrustGuard) retire(key GuardTrustInodeKey) {
	var r trustRows
	if t.objs.GuardTrustedFiles.Lookup(key, &r.flags) != nil {
		return
	}
	for i, m := range t.trustBitMaps() {
		r.has[i] = m.Lookup(key, &r.bits[i]) == nil
	}
	t.retiredMu.Lock()
	defer t.retiredMu.Unlock()
	if t.retired == nil {
		t.retired = make(map[GuardTrustInodeKey]trustRows)
	}
	if len(t.retired) < maxRetiredRows {
		t.retired[key] = r
	}
}

// rowsOf reads key's live rows, else the ones retire stashed. ok false: neither.
func (t *TrustGuard) rowsOf(key GuardTrustInodeKey) (trustRows, bool) {
	var r trustRows
	if t.objs.GuardTrustedFiles.Lookup(key, &r.flags) == nil {
		for i, m := range t.trustBitMaps() {
			r.has[i] = m.Lookup(key, &r.bits[i]) == nil
		}
		return r, true
	}
	t.retiredMu.Lock()
	defer t.retiredMu.Unlock()
	r, ok := t.retired[key]
	return r, ok
}

// trustedDirAny mirrors TRUSTED_DIR_ANY in guard_trust.bpf.c.
const trustedDirAny = uint64(1) << 63

// TrustedDir is a guarded resource root for the library allowlist. Loaders nil: every whitelisted
// binary may load from it (a whitelist-mode root). Otherwise only those binaries may (a read-only
// lib_dir, whose Loaders are its writers).
type TrustedDir struct {
	Path    string
	Loaders []string
}

// planTrustedDirs assigns one guard_trusted_dirs bit per distinct loader set, returning each root's
// mask and each loader's bits. Past 63 sets the rest get mask 0 (nobody loads from them: fail
// closed) and are returned in refused.
func planTrustedDirs(dirs []TrustedDir) (roots, loaders map[string]uint64, refused []string) {
	roots = make(map[string]uint64, len(dirs))
	loaders = make(map[string]uint64)
	bitOf := make(map[string]uint64)
	for _, d := range dirs {
		if d.Loaders == nil {
			roots[d.Path] |= trustedDirAny
			continue
		}
		set := append([]string(nil), d.Loaders...)
		sort.Strings(set)
		key := strings.Join(slices.Compact(set), "\x00")
		bit, ok := bitOf[key]
		if !ok {
			if len(bitOf) >= 63 {
				refused = append(refused, d.Path)
				if _, seen := roots[d.Path]; !seen {
					roots[d.Path] = 0
				}
				continue
			}
			bit = uint64(1) << len(bitOf)
			bitOf[key] = bit
			for _, l := range set {
				loaders[l] |= bit
			}
		}
		roots[d.Path] |= bit
	}
	return roots, loaders, refused
}

// SetGuardedDirs records every guarded resource root and who may load libraries from below it
// (under_guarded_tree in guard_trust.bpf.c). Unresolvable paths are skipped with a warning.
// Idempotent: synced to exactly the new set via syncMap (put-then-delete-stale), never cleared
// first, so a reload never opens a window where the library allowlist trusts no one.
func (t *TrustGuard) SetGuardedDirs(dirs []TrustedDir) error {
	roots, loaders, refused := planTrustedDirs(dirs)
	for _, p := range refused {
		log.Warnf("trust guard: too many distinct lib_dir writer sets, libraries under %s are "+
			"trusted for no one", p)
	}
	userBits := make(map[GuardInodeKey]uint64)
	for path, bits := range loaders {
		dev, ino, err := ebpf.StatConfined(path)
		if err != nil {
			log.Warnf("trust guard: skipping unresolvable lib_dir writer %s: %v", path, err)
			continue
		}
		userBits[GuardInodeKey{Dev: dev, Ino: ino}] |= bits
	}
	if err := syncMap(t.objs.GuardLibdirUsers, userBits); err != nil {
		return fmt.Errorf("syncing lib_dir loaders: %w", err)
	}
	rootKeys := make(map[GuardInodeKey]uint64, len(roots))
	for r, mask := range roots {
		dev, ino, err := ebpf.StatInode(r)
		if err != nil {
			log.Warnf("trust guard: skipping unresolvable guarded root %s: %v", r, err)
			continue
		}
		rootKeys[GuardInodeKey{Dev: dev, Ino: ino}] = mask
	}
	if err := syncMap(t.objs.GuardTrustedDirs, rootKeys); err != nil {
		return fmt.Errorf("syncing guarded dirs: %w", err)
	}
	log.Infof("trust guard: %d guarded resource root(s) recorded (libraries inside them are trusted)", len(rootKeys))
	return nil
}

// trustHook pairs an LSM program with a human name for logging; required marks a hook whose attach
// failure must fail the whole trust guard.
type trustHook struct {
	prog     *cilium.Program
	name     string
	required bool
}

func (t *TrustGuard) hooks() []trustHook {
	return []trustHook{
		// mmap_file is protection #2 — the library-load allowlist that stops LD_PRELOAD of attacker
		// code into a whitelisted process. It is REQUIRED: without it the trust guard denies nothing
		// meaningful, so if it can't attach (e.g. the per-attach-point trampoline cap
		// BPF_MAX_TRAMP_LINKS is exhausted) the guard must fail rather than report success.
		{t.objs.TrustMmap, "mmap_file", true},
		{t.objs.TrustFileOpen, "file_open", false},
		{t.objs.TrustPathUnlink, "path_unlink", false},
		{t.objs.TrustPathRename, "path_rename", false},
		{t.objs.TrustPathTruncate, "path_truncate", false},
		{t.objs.TrustPathMknod, "path_mknod", false},
		{t.objs.TrustPathMkdir, "path_mkdir", false},
		{t.objs.TrustPathSymlink, "path_symlink", false},
		{t.objs.TrustPathLink, "path_link", false},
		// The code-suspect lifecycle (trust_code_suspect). bprm_committed_creds and task_free clear
		// marks; without them marks go stale. Required together with the fork tracepoint
		// (attachSuspectFork), or a suspect process's children would escape the mark.
		{t.objs.TrustBprmCommitted, "bprm_committed_creds", true},
		{t.objs.TrustTaskFree, "task_free", true},
		// Superseded keys for trusted libraries (exe_supersede.h): without them a freed library's
		// number, reused by any file, would stay TRUSTED_LIB.
		{t.objs.TrustInodeUnlink, "inode_unlink", true},
		{t.objs.TrustInodeRename, "inode_rename", true},
		{t.objs.TrustInodeFree, "inode_free_security", true},
	}
}

// Start attaches every trust program and drains events. A required hook (mmap_file, the
// code-suspect lifecycle, sb_delete) that cannot attach fails the start; other hooks are
// best-effort. It never blocks the per-resource guards.
func (t *TrustGuard) Start() error {
	if err := t.startMountVouch(); err != nil {
		t.detachLinks()
		return err
	}
	for _, h := range t.hooks() {
		l, err := link.AttachLSM(link.LSMOptions{Program: h.prog})
		if err != nil {
			if h.required {
				t.detachLinks()
				return fmt.Errorf("required trust hook %s failed to attach: %w — the LD_PRELOAD "+
					"library allowlist would be unenforced", h.name, err)
			}
			log.Warnf("trust guard: skipping %s (%v) — that protection is unavailable; "+
				"per-resource enforcement is unaffected", h.name, err)
			continue
		}
		t.links = append(t.links, l)
	}
	if len(t.links) == 0 {
		return errors.New("no trust programs attached")
	}
	if err := t.attachSuspectFork(); err != nil {
		t.detachLinks()
		return err
	}
	t.attachMemfdProvenance()
	if err := t.startMountWatch(); err != nil {
		t.detachLinks()
		return err
	}

	rd, err := ringbuf.NewReader(t.objs.TrustRb)
	if err != nil {
		return fmt.Errorf("opening trust ringbuf: %w", err)
	}
	t.rd = rd
	go t.readLoop()
	setSupersedeTrust(t)
	return nil
}

// attachSuspectFork attaches the fork tracepoint that carries a code-suspect mark to the child,
// which shares the parent's mapped code. Required: without it a suspect process escapes the mark
// by forking.
func (t *TrustGuard) attachSuspectFork() error {
	l, err := link.AttachTracing(link.TracingOptions{
		Program:    t.objs.TrustSchedProcessFork,
		AttachType: cilium.AttachTraceRawTp,
	})
	if err != nil {
		return fmt.Errorf("required trust hook sched_process_fork failed to attach: %w — code "+
			"preloaded before a binary was whitelisted would survive a fork", err)
	}
	t.links = append(t.links, l)
	return nil
}

// attachMemfdProvenance attaches an fexit on the kernel's memfd allocation, recording that a
// whitelisted process created this anonymous inode (trust_memfd_alloc in guard_trust.bpf.c). A
// memfd never reaches security_file_open, so without it the first step of a GPU driver's JIT
// fallback chain has no provenance.
//
// Best-effort and fail-closed: without the record a memfd exec-map stays denied and the driver
// falls through to O_TMPFILE/mkstemp, which file_open sees.
func (t *TrustGuard) attachMemfdProvenance() {
	if !t.memfd {
		log.Warnf("trust guard: skipping memfd provenance (this kernel has no %s) — a whitelisted "+
			"process's memfd-backed runtime code stays denied (GPU drivers fall back to their "+
			"file-backed JIT paths, which are covered); every other protection is unaffected", memfdTarget)
		return
	}
	l, err := link.AttachTracing(link.TracingOptions{
		Program:    t.objs.TrustMemfdAlloc,
		AttachType: cilium.AttachTraceFExit,
	})
	if err != nil {
		log.Warnf("trust guard: skipping memfd provenance (%v) — a whitelisted process's "+
			"memfd-backed runtime code stays denied (GPU drivers fall back to their "+
			"file-backed JIT paths, which are covered); every other protection is unaffected", err)
		return
	}
	t.links = append(t.links, l)
}

func (t *TrustGuard) readLoop() {
	for {
		rec, err := t.rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				return
			}
			continue
		}
		var ev trustEvent
		if err := binary.Read(bytes.NewReader(rec.RawSample), binary.LittleEndian, &ev); err != nil {
			continue
		}
		comm, path := logging.SanitizeText(cStr(ev.Comm[:])), logging.SanitizeText(cStr(ev.Path[:]))
		switch ev.Kind {
		case trustWriteblock:
			logTrustDenied("WRITE", comm, ev.PID, ev.UID, path,
				"a non-whitelisted process tried to modify a protected binary")
		case trustPlant:
			logTrustDenied("PLANT", comm, ev.PID, ev.UID, path,
				"only the owning app's binaries may create or modify a file a catalog glob reserves")
		case trustSuspect:
			logTrustDenied("PRELOADED", comm, ev.PID, ev.UID, path,
				"this process mapped untrusted code before its binary was whitelisted; restart it")
		default:
			logTrustDenied("LIBLOAD", comm, ev.PID, ev.UID, path,
				"a whitelisted binary tried to exec-map an untrusted file "+
					"(not root-owned, not in a guarded tree, not allow_lib, not its app's reserved library)")
		}
	}
}

// logTrustDenied prints a denial with the syslog <4> (warning) marker so journald colors it
// yellow like DAEMON DENIED. Callers must pass sanitized comm/path (event-controlled).
func logTrustDenied(op, comm string, pid, uid uint32, path, reason string) {
	fmt.Fprintf(os.Stderr, "<4>TRUST DENIED  op=%s  comm=%s  pid=%d  uid=%d  path=%s — %s\n",
		op, comm, pid, uid, path, reason)
}

// Stop detaches the trust programs and stops the reader. Safe to call once.
func (t *TrustGuard) Stop() {
	select {
	case <-t.done:
		return
	default:
		close(t.done)
	}
	supersede.mu.Lock()
	if supersede.trust == t {
		supersede.trust = nil
	}
	supersede.mu.Unlock()
	if t.rd != nil {
		_ = t.rd.Close()
	}
	t.detachLinks()
	t.stopMountWatch()
	_ = t.objs.Close()
}

// detachLinks closes every attached LSM link (used both by Stop and by a failed Start that must not
// leave a partial attach behind).
func (t *TrustGuard) detachLinks() {
	for _, l := range t.links {
		_ = l.Close()
	}
	t.links = nil
}

func cStr(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		return string(b[:i])
	}
	return string(b)
}
