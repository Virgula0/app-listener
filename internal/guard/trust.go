package guard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

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
// cleared first, so a reload that dropped a binary also drops its trust entry (no stale over-trust).
func (t *TrustGuard) SetTrusted(binaries, libs []string) error {
	if err := clearInodeMap(t.objs.GuardTrustedFiles); err != nil {
		return fmt.Errorf("clearing trusted files: %w", err)
	}
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
	for key, f := range flags {
		if err := t.objs.GuardTrustedFiles.Put(key, f); err != nil {
			return fmt.Errorf("populating trusted file %+v: %w", key, err)
		}
	}
	log.Infof("trust guard: %d trusted inode(s) loaded (%d binaries, %d libraries requested)",
		len(flags), len(binaries), len(libs))
	return nil
}

// SetGuardedDirs records the (dev, ino) of every guarded resource root so a library loaded from
// inside a write-protected tree is trusted without allow_lib (under_guarded_tree in
// guard_trust.bpf.c). Unresolvable paths are skipped with a warning. Idempotent (cleared first).
func (t *TrustGuard) SetGuardedDirs(roots []string) error {
	if err := clearInodeMap(t.objs.GuardTrustedDirs); err != nil {
		return fmt.Errorf("clearing guarded dirs: %w", err)
	}
	v := uint8(1)
	n := 0
	for _, r := range roots {
		dev, ino, err := ebpf.StatInode(r)
		if err != nil {
			log.Warnf("trust guard: skipping unresolvable guarded root %s: %v", r, err)
			continue
		}
		if err := t.objs.GuardTrustedDirs.Put(GuardInodeKey{Dev: dev, Ino: ino}, v); err != nil {
			return fmt.Errorf("recording guarded root %s: %w", r, err)
		}
		n++
	}
	log.Infof("trust guard: %d guarded resource root(s) recorded (libraries inside them are trusted)", n)
	return nil
}

// clearInodeMap deletes every entry of a GuardInodeKey-keyed hash map so SetTrusted/SetGuardedDirs
// can be re-applied on reload as a true replace, not an add-only merge.
func clearInodeMap(m *cilium.Map) error {
	var keys []GuardInodeKey
	var k GuardInodeKey
	var v uint8
	it := m.Iterate()
	for it.Next(&k, &v) {
		keys = append(keys, k)
	}
	if err := it.Err(); err != nil {
		return err
	}
	for i := range keys {
		if err := m.Delete(keys[i]); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
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
		default:
			logTrustDenied("LIBLOAD", comm, ev.PID, ev.UID, path,
				"a whitelisted binary tried to exec-map an untrusted file "+
					"(not root-owned, not in a guarded tree, not allow_lib)")
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
