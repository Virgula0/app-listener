package guard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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

// TrustGuard owns guard_trusted_files and the daemon-wide trusted-binary /
// trusted-library protections (see guard_trust.bpf.c). It is created and
// attached ONCE for the whole daemon, independent of the per-resource guards.
//
// Both protections (#1 writer attribution and #2 the library-load allowlist)
// are always enforced once the programs attach.
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

// SetTrusted (re)populates guard_trusted_files from the resolved binary and
// library paths. A path present as both a binary and a library carries both
// flags. Unresolvable paths are skipped with a warning.
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
	for key, f := range flags {
		if err := t.objs.GuardTrustedFiles.Put(key, f); err != nil {
			return fmt.Errorf("populating trusted file %+v: %w", key, err)
		}
	}
	log.Infof("trust guard: %d trusted inode(s) loaded (%d binaries, %d libraries requested)",
		len(flags), len(binaries), len(libs))
	return nil
}

// SetGuardedDirs records the (dev, ino) of every guarded resource root, so a
// library loaded from inside one of those write-protected trees is trusted
// without an explicit allow_lib entry (see under_guarded_tree in
// guard_trust.bpf.c). Unresolvable paths are skipped with a warning.
func (t *TrustGuard) SetGuardedDirs(roots []string) error {
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

// trustHook pairs an LSM program with a human name for logging.
type trustHook struct {
	prog *cilium.Program
	name string
}

func (t *TrustGuard) hooks() []trustHook {
	return []trustHook{
		{t.objs.TrustMmap, "mmap_file"},
		{t.objs.TrustFileOpen, "file_open"},
		{t.objs.TrustPathUnlink, "path_unlink"},
		{t.objs.TrustPathRename, "path_rename"},
		{t.objs.TrustPathTruncate, "path_truncate"},
	}
}

// Start attaches every trust LSM program and begins draining events. Each
// attach is best-effort: a failure logs and that one protection is skipped,
// never blocking the daemon. Enforcement of the per-resource guards is
// unaffected regardless.
func (t *TrustGuard) Start() error {
	for _, h := range t.hooks() {
		l, err := link.AttachLSM(link.LSMOptions{Program: h.prog})
		if err != nil {
			log.Warnf("trust guard: skipping %s (%v) — that protection is unavailable; "+
				"per-resource enforcement is unaffected", h.name, err)
			continue
		}
		t.links = append(t.links, l)
	}
	if len(t.links) == 0 {
		return errors.New("no trust programs attached")
	}

	rd, err := ringbuf.NewReader(t.objs.TrustRb)
	if err != nil {
		return fmt.Errorf("opening trust ringbuf: %w", err)
	}
	t.rd = rd
	go t.readLoop()
	return nil
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
		switch ev.Kind {
		case trustWriteblock:
			log.Warnf("TRUST DENIED  op=WRITE  comm=%s  pid=%d  uid=%d  path=%s — a non-whitelisted process "+
				"tried to modify a protected binary", cStr(ev.Comm[:]), ev.PID, ev.UID, cStr(ev.Path[:]))
		default:
			log.Warnf("TRUST DENIED  op=LIBLOAD  comm=%s  pid=%d  uid=%d  path=%s — a whitelisted binary tried to "+
				"exec-map an untrusted file (not root-owned, not in a guarded tree, not allow_lib)",
				cStr(ev.Comm[:]), ev.PID, ev.UID, cStr(ev.Path[:]))
		}
	}
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
	for _, l := range t.links {
		_ = l.Close()
	}
	t.links = nil
	_ = t.objs.Close()
}

func cStr(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		return string(b[:i])
	}
	return string(b)
}
