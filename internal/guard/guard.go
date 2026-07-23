package guard

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
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

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type Mode int

const (
	ModeBlacklist Mode = iota
	ModeWhitelist
)

type BinaryEntry struct {
	Path string
	Hash [sha256.Size]byte
	Comm string
}

type GuardEvent struct {
	ebpf.FileEvent
	Blocked bool
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
	recursive bool
	depth     int
}

func NewGuard(path string, mode Mode, binaries []BinaryEntry, recursive bool, depth int) (*Guard, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	entries := make([]BinaryEntry, len(binaries))
	copy(entries, binaries)

	g := &Guard{
		events:    make(chan GuardEvent, 1024),
		done:      make(chan struct{}),
		path:      path,
		mode:      mode,
		binaries:  entries,
		recursive: recursive,
		depth:     depth,
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

	attachments := []struct {
		prog *cilium.Program
		hook string
	}{
		{g.objs.GuardFileOpen, "file_open"},
		{g.objs.GuardFilePermission, "file_permission"},
		{g.objs.GuardMmapFile, "mmap_file"},
		{g.objs.GuardPathUnlink, "path_unlink"},
		{g.objs.GuardPathRename, "path_rename"},
		{g.objs.GuardPathSymlink, "path_symlink"},
		{g.objs.GuardPathLink, "path_link"},
		{g.objs.GuardPathMkdir, "path_mkdir"},
	}

	for _, a := range attachments {
		l, err := link.AttachLSM(link.LSMOptions{
			Program: a.prog,
		})
		if err != nil {
			log.Warnf("skipping LSM hook %s: %v", a.hook, err)
			continue
		}
		g.links = append(g.links, l)
	}

	log.Infof("guard created \u2014 %d LSM hooks attached, watching: %s (%s)",
		len(g.links), path, modeString(mode))
	return g, nil
}

func modeString(mode Mode) string {
	if mode == ModeWhitelist {
		return "whitelist"
	}
	return "blacklist"
}

func (g *Guard) populateMaps() error {
	// Set the operating mode
	if err := g.objs.GuardConfig.Put(uint32(0), uint64(g.mode)); err != nil {
		return fmt.Errorf("setting mode in config: %w", err)
	}

	info, err := os.Stat(g.path)
	if err != nil {
		return fmt.Errorf("stating guarded path %s: %w", g.path, err)
	}

	if info.IsDir() {
		if err := g.scanDirInodes(g.path, 0); err != nil {
			return err
		}
	} else {
		if err := g.addInode(g.path); err != nil {
			return err
		}
	}

	for _, b := range g.binaries {
		var key struct {
			Comm [16]byte
		}
		copy(key.Comm[:], b.Comm)

		var val uint8 = 1
		if err := g.objs.GuardComms.Put(key, val); err != nil {
			return fmt.Errorf("storing comm %q in map: %w", b.Comm, err)
		}

		// Store the binary's exe inode for anti-spoof verification.
		// The BPF program reads current->mm->exe_file->f_inode and
		// checks it against guard_exe_inodes to prevent prctl-based
		// comm spoofing.
		var es syscall.Stat_t
		if err := syscall.Stat(b.Path, &es); err != nil {
			log.Warnf("cannot stat binary %s for exe inode: %v", b.Path, err)
			continue
		}
		inodeKey := GuardInodeKey{
			Dev: uint64((unix.Major(es.Dev) << 20) | unix.Minor(es.Dev)),
			Ino: es.Ino,
		}
		if err := g.objs.GuardExeInodes.Put(inodeKey, val); err != nil {
			return fmt.Errorf("storing exe inode for %s: %w", b.Path, err)
		}
	}

	// Store the guarded path for symlink target matching
	var pathBuf [256]byte
	copy(pathBuf[:], g.path)
	if err := g.objs.GuardPath.Put(uint32(0), pathBuf); err != nil {
		return fmt.Errorf("storing guarded path: %w", err)
	}

	return nil
}

func (g *Guard) addInode(path string) error {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return err
	}

	major := unix.Major(s.Dev)
	minor := unix.Minor(s.Dev)
	kernelDev := uint64((major << 20) | minor)

	key := GuardInodeKey{
		Dev: kernelDev,
		Ino: s.Ino,
	}

	var val uint8 = 1
	if err := g.objs.GuardInodes.Put(key, val); err != nil {
		return fmt.Errorf("adding inode %s to map: %w", path, err)
	}
	return nil
}

func (g *Guard) scanDirInodes(dir string, currentDepth int) error {
	if err := g.addInode(dir); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !g.recursive {
				continue
			}
			if g.depth > 0 && currentDepth+1 >= g.depth {
				continue
			}
			if err := g.scanDirInodes(fullPath, currentDepth+1); err != nil {
				return err
			}
		} else {
			if err := g.addInode(fullPath); err != nil {
				return err
			}
		}
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

	if !ge.Blocked {
		g.verifyHash(ge)
	}

	return ge, true
}

func (g *Guard) verifyHash(ge *GuardEvent) {
	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ge.PID))
	if err != nil {
		return
	}

	if !g.isGuardedBinary(exePath) {
		log.Warnf("process %d (%s) spoofed comm \u2014 actual binary: %s",
			ge.PID, ge.Comm, exePath)
	}
}

func (g *Guard) isGuardedBinary(exePath string) bool {
	for _, b := range g.binaries {
		if b.Path == exePath {
			return true
		}
		abs, err := filepath.EvalSymlinks(exePath)
		if err == nil && abs == b.Path {
			return true
		}
	}
	return false
}

func ComputeBinaryEntry(path string) (BinaryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}

	hash := sha256.Sum256(data)
	comm := filepath.Base(path)

	if len(comm) > 15 {
		comm = comm[:15]
	}

	return BinaryEntry{
		Path: path,
		Hash: hash,
		Comm: comm,
	}, nil
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
