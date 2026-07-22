package ebpf

import (
	"bytes"
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
)

type Target struct {
	Path  string // original path (for display)
	Dir   string // parent dir (for files) or the dir itself
	File  string // basename (empty for directories)
	IsDir bool
}

type Monitor struct {
	objs    MonitorObjects
	links   []link.Link
	events  chan FileEvent
	done    chan struct{}
	mu      sync.Mutex
	stopped bool

	targets   []Target
	recursive bool
	depth     int
	ownPID    int

	watchInodes map[string]string // "dev:ino" → watched path for hardlink detection
}

func NewMonitor(targets []Target, recursive bool, depth int) (*Monitor, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	m := &Monitor{
		events:    make(chan FileEvent, 1024),
		done:      make(chan struct{}),
		targets:   targets,
		recursive: recursive,
		depth:     depth,
		ownPID:    os.Getpid(),
	}

	var objs MonitorObjects
	if err := LoadMonitorObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("loading BPF objects: %w", err)
	}
	m.objs = objs

	type attachDef struct {
		prog  *cilium.Program
		typ   string
		group string
		name  string
	}

	attachments := []attachDef{
		// VFS-level kprobes catch ALL open/read/write (syscalls, io_uring, execve, etc.)
		// Path is resolved from dentry manually instead of using bpf_d_path.
		{objs.TraceVfsOpen, "kprobe", "", "vfs_open"},
		{objs.TraceVfsRead, "kprobe", "", "vfs_read"},
		{objs.TraceVfsWrite, "kprobe", "", "vfs_write"},
		// Additional read-variant kprobes (readv, copy_file_range, splice, sendfile)
		{objs.TraceVfsReadv, "kprobe", "", "vfs_readv"},
		{objs.TraceVfsCopyFileRange, "kprobe", "", "vfs_copy_file_range"},
		{objs.TraceDoSplice, "kprobe", "", "do_splice"},
		{objs.TraceDoSpliceDirect, "kprobe", "", "do_splice_direct"},
		{objs.TraceSpliceFileRange, "kprobe", "", "splice_file_range"},
		{objs.TraceDoSendfile, "kprobe", "", "do_sendfile"},
		{objs.TraceVfsIterRead, "kprobe", "", "vfs_iter_read"},
		// Memory-mapped I/O via kprobe (works without tracefs)
		{objs.TraceSecurityMmapFile, "kprobe", "", "security_mmap_file"},
		// Tracepoints for filesystem metadata operations
		{objs.TraceUnlinkat, "tracepoint", "syscalls", "sys_enter_unlinkat"},
		{objs.TraceRenameat2, "tracepoint", "syscalls", "sys_enter_renameat2"},
		{objs.TraceSymlinkat, "tracepoint", "syscalls", "sys_enter_symlinkat"},
		{objs.TraceLinkat, "tracepoint", "syscalls", "sys_enter_linkat"},
		{objs.TraceMkdirat, "tracepoint", "syscalls", "sys_enter_mkdirat"},
		// Legacy tracepoint for mmap (kept for tracefs-enabled systems)
		{objs.TraceMmap, "tracepoint", "syscalls", "sys_enter_mmap"},
	}

	for _, a := range attachments {
		var l link.Link
		var err error

		switch a.typ {
		case "kprobe":
			l, err = link.Kprobe(a.name, a.prog, nil)
		case "tracepoint":
			l, err = link.Tracepoint(a.group, a.name, a.prog, nil)
		}

		if err != nil {
			log.Warnf("skipping %s/%s: %v", a.typ, a.name, err)
			continue
		}
		m.links = append(m.links, l)
	}

	m.watchInodes = make(map[string]string)
	m.populateWatchInodes()

	paths := make([]string, len(targets))
	for i, t := range targets {
		paths[i] = t.Path
	}
	log.Infof("monitor created — %d probes/tracepoints attached, watching: %s",
		len(m.links), strings.Join(paths, ", "))
	return m, nil
}

func (m *Monitor) Events() <-chan FileEvent {
	return m.events
}

func (m *Monitor) Start() error {
	rd, err := ringbuf.NewReader(m.objs.Rb)
	if err != nil {
		m.cleanup()
		return fmt.Errorf("ringbuf reader: %w", err)
	}

	paths := make([]string, len(m.targets))
	for i, t := range m.targets {
		paths[i] = t.Path
	}
	log.Infof("monitor started — watching: %s", strings.Join(paths, ", "))
	go m.readLoop(rd)
	return nil
}

func (m *Monitor) readLoop(rd *ringbuf.Reader) {
	defer rd.Close()

	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				return
			}
			log.Errorf("ringbuf read error: %v", err)
			continue
		}

		var be bpfEvent
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &be); err != nil {
			log.Errorf("decode event: %v", err)
			continue
		}

		ev := be.toFileEvent()
		ev.Timestamp = time.Now().UnixNano()

		if int(ev.PID) == m.ownPID {
			continue
		}

		ev.Path = m.resolveRelativePath(ev.PID, ev.Path)
		ev.Dest = m.resolveRelativePath(ev.PID, ev.Dest)

		m.resolveFDPath(&ev)

		// Resolve symlinks so that cat on a symlink to a watched file is detected.
		ev.Path = m.resolveSymlinkToTarget(ev.Path)

		if !m.inWatchPath(ev.Path, ev.Dest) {
			// If still no match, try hardlink detection via inode.
			if hardlinkPath := m.checkHardlinkByInode(ev.Path); hardlinkPath != "" {
				ev.Path = hardlinkPath
			} else {
				continue
			}
		}

		if !m.matchesFilter(&ev) {
			continue
		}

		select {
		case m.events <- ev:
		case <-m.done:
			return
		}
	}
}

func (m *Monitor) resolveRelativePath(pid uint32, path string) string {
	if path == "" || path[0] == '/' {
		return path
	}
	cwd, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	if err != nil {
		return path
	}
	return filepath.Join(cwd, path)
}

func (m *Monitor) resolveFDPath(ev *FileEvent) {
	if ev.Type != EventMmap {
		return
	}
	if ev.Path != "" {
		return
	}

	path := fmt.Sprintf("/proc/%d/fd/%d", ev.PID, ev.FD)
	target, err := os.Readlink(path)
	if err != nil {
		return
	}
	ev.Path = target
}

func (m *Monitor) inWatchPath(path, dest string) bool {
	if path != "" && m.matchesAnyTarget(path) {
		return true
	}
	return dest != "" && m.matchesAnyTarget(dest)
}

func (m *Monitor) matchesAnyTarget(p string) bool {
	for _, t := range m.targets {
		if t.File != "" {
			if p == t.File || strings.HasSuffix(p, "/"+t.File) {
				return true
			}
		} else {
			if p == t.Dir || strings.HasPrefix(p, t.Dir+"/") {
				return true
			}
		}
	}
	return false
}

func (m *Monitor) matchesFilter(ev *FileEvent) bool {
	path := ev.Path
	if path == "" {
		return false
	}

	for _, t := range m.targets {
		if !t.IsDir {
			if path == t.File || strings.HasSuffix(path, "/"+t.File) {
				return true
			}
			continue
		}
		if !m.isWithinDir(path, t.Dir) {
			continue
		}
		if m.recursive && m.depth == 0 {
			return true
		}
		if !m.recursive {
			rel, err := filepath.Rel(t.Dir, path)
			if err != nil {
				continue
			}
			if rel == "." || !strings.Contains(rel, string(filepath.Separator)) {
				return true
			}
		}
		if m.recursive && m.depth > 0 {
			rel, err := filepath.Rel(t.Dir, path)
			if err != nil {
				continue
			}
			depth := 0
			if rel != "." {
				depth = len(strings.Split(rel, string(filepath.Separator)))
			}
			if depth <= m.depth {
				return true
			}
		}
	}
	return false
}

func (m *Monitor) isWithinDir(path, dir string) bool {
	return path == dir || strings.HasPrefix(path, dir+"/")
}

// populateWatchInodes scans watched targets and records their (dev, inode)
// pairs so we can detect hardlink access by inode match.
func (m *Monitor) populateWatchInodes() {
	for _, t := range m.targets {
		if !t.IsDir {
			m.addInode(t.Path)
		} else {
			m.scanDirInodes(t.Dir, 0)
		}
	}
}

func (m *Monitor) addInode(path string) {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return
	}
	key := fmt.Sprintf("%d:%d", s.Dev, s.Ino)
	m.watchInodes[key] = path
}

func (m *Monitor) scanDirInodes(dir string, currentDepth int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !m.recursive {
				continue
			}
			if m.depth > 0 && currentDepth+1 > m.depth {
				continue
			}
			m.scanDirInodes(fullPath, currentDepth+1)
			continue
		}
		m.addInode(fullPath)
	}
}

// resolveSymlinkToTarget resolves all symlinks in path and, if the resolved
// path falls within a watched target, returns the resolved path.
func (m *Monitor) resolveSymlinkToTarget(path string) string {
	if path == "" || m.matchesAnyTarget(path) {
		return path
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	if m.matchesAnyTarget(resolved) {
		return resolved
	}
	return path
}

// checkHardlinkByInode stats the given path and, if its (dev, inode) matches
// a pre-recorded watched file, returns the corresponding watched path.
func (m *Monitor) checkHardlinkByInode(path string) string {
	if path == "" || len(m.watchInodes) == 0 {
		return ""
	}
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return ""
	}
	key := fmt.Sprintf("%d:%d", s.Dev, s.Ino)
	if watchedPath, ok := m.watchInodes[key]; ok {
		return watchedPath
	}
	return ""
}

func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return
	}
	m.stopped = true
	close(m.done)
	m.cleanup()
}

func (m *Monitor) cleanup() {
	for _, l := range m.links {
		l.Close()
	}
	m.links = nil
	m.objs.Close()
}
