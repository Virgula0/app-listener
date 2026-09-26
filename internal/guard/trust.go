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
)

// TrustGuard owns guard_trusted_files and the daemon-wide trusted-binary/library protections
// (guard_trust.bpf.c). Created and attached ONCE for the whole daemon, independent of the
// per-resource guards. Both protections (#1 writer attribution, #2 library-load allowlist) are
// always enforced once attached.
type TrustGuard struct {
	objs  GuardTrustObjects
	links []link.Link
	rd    *ringbuf.Reader
	done  chan struct{}
	// ownerByPath: UpdaterPlan.Owners, for AllowReplacement.
	ownerMu     sync.Mutex
	ownerByPath map[string]uint64
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
	t := &TrustGuard{done: make(chan struct{})}
	if err := LoadGuardTrustObjects(&t.objs, nil); err != nil {
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
		dev, ino, err := ebpf.StatInode(path)
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
	if err := syncMap(t.objs.GuardTrustedFiles, flags); err != nil {
		return fmt.Errorf("syncing trusted files: %w", err)
	}
	log.Infof("trust guard: %d trusted inode(s) loaded (%d binaries, %d libraries requested)",
		len(flags), len(binaries), len(libs))
	return nil
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
		dev, ino, err := ebpf.StatInode(path)
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
	}
}

// Start attaches every trust LSM program and drains events. A required hook (mmap_file) that cannot
// attach fails the start; other hooks are best-effort. It never blocks the per-resource guards.
func (t *TrustGuard) Start() error {
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
	t.attachMemfdProvenance()

	rd, err := ringbuf.NewReader(t.objs.TrustRb)
	if err != nil {
		return fmt.Errorf("opening trust ringbuf: %w", err)
	}
	t.rd = rd
	go t.readLoop()
	return nil
}

// attachMemfdProvenance attaches the one non-LSM trust program: an fexit on the kernel's memfd
// allocation, recording that a whitelisted process created this anonymous inode (trust_memfd_alloc
// in guard_trust.bpf.c). A memfd never reaches security_file_open, so without it the first step of
// a GPU driver's JIT fallback chain has no provenance.
//
// Best-effort and fail-closed: without the record a memfd exec-map stays denied and the driver
// falls through to O_TMPFILE/mkstemp, which file_open sees.
func (t *TrustGuard) attachMemfdProvenance() {
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
	if t.rd != nil {
		_ = t.rd.Close()
	}
	t.detachLinks()
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
