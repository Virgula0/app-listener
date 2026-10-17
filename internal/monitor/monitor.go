package monitor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type Monitor struct {
	objs    MonitorObjects
	links   []link.Link
	events  chan ebpf.FileEvent
	done    chan struct{}
	mu      sync.Mutex
	stopped bool

	targets    []ebpf.Target
	recursive  bool
	depth      int
	ownPID     int
	eventTypes atomic.Value // stores []ebpf.EventType

	watchInodes map[uint64][]watchedInode // by inode number, for hardlink detection

	// paths memoizes EvalSymlinks for event paths not lexically inside a watched target: on the
	// unfiltered VFS firehose the same paths recur constantly.
	paths  *pathCache
	mounts *mountTable
}

type watchedInode struct {
	path  string
	sdev  uint32 // kernel s_dev of its mount (what BPF reports); 0 when unknown
	stDev uint64
}

const (
	pathCacheMax = 16384
	pathCacheTTL = 30 * time.Second
)

type pathCacheEntry struct {
	resolved string // filepath.EvalSymlinks(path); "" when it errored
	at       time.Time
}

type pathCache struct {
	mu sync.Mutex
	m  map[string]pathCacheEntry
}

func newPathCache() *pathCache {
	return &pathCache{m: make(map[string]pathCacheEntry, 1024)}
}

func (c *pathCache) lookup(path string) (pathCacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[path]
	if !ok || time.Since(e.at) > pathCacheTTL {
		return pathCacheEntry{}, false
	}
	return e, true
}

func (c *pathCache) store(path string, e pathCacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.m) >= pathCacheMax {
		c.m = make(map[string]pathCacheEntry, pathCacheMax)
	}
	c.m[path] = e
}

func NewMonitor(targets []ebpf.Target, recursive bool, depth int) (*Monitor, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	m := &Monitor{
		events:    make(chan ebpf.FileEvent, 1024),
		done:      make(chan struct{}),
		targets:   targets,
		recursive: recursive,
		depth:     depth,
		ownPID:    os.Getpid(),
		paths:     newPathCache(),
		mounts:    newMountTable(),
	}
	m.eventTypes.Store(ebpf.EventTypes())

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
		{objs.TraceVfsOpen, "kprobe", "", "vfs_open"},
		{objs.TraceVfsRead, "kprobe", "", "vfs_read"},
		{objs.TraceVfsWrite, "kprobe", "", "vfs_write"},
		{objs.TraceVfsReadv, "kprobe", "", "vfs_readv"},
		{objs.TraceVfsCopyFileRange, "kprobe", "", "vfs_copy_file_range"},
		{objs.TraceDoSplice, "kprobe", "", "do_splice"},
		{objs.TraceDoSplice, "kprobe", "", "__do_splice"},
		{objs.TraceDoSpliceDirect, "kprobe", "", "do_splice_direct"},
		{objs.TraceSpliceFileRange, "kprobe", "", "splice_file_range"},
		{objs.TraceDoSendfile, "kprobe", "", "do_sendfile"},
		{objs.TraceVfsIterRead, "kprobe", "", "vfs_iter_read"},
		{objs.TraceIoRead, "kprobe", "", "io_read"},
		{objs.TraceIoWrite, "kprobe", "", "io_write"},
		{objs.TraceSecurityMmapFile, "kprobe", "", "security_mmap_file"},
		{objs.TraceVfsMkdir, "kprobe", "", "vfs_mkdir"},
		{objs.TraceVfsRmdir, "kprobe", "", "vfs_rmdir"},
		{objs.TraceVfsUnlink, "kprobe", "", "vfs_unlink"},
		{objs.TraceVfsRename, "kprobe", "", "vfs_rename"},
		{objs.TraceVfsSymlink, "kprobe", "", "vfs_symlink"},
		{objs.TraceVfsLink, "kprobe", "", "vfs_link"},
		{objs.TraceNotifyChange, "kprobe", "", "notify_change"},
		{objs.TraceVfsSetxattr, "kprobe", "", "vfs_setxattr"},
		{objs.TraceVfsRemovexattr, "kprobe", "", "vfs_removexattr"},
		{objs.TraceVfsMknod, "kprobe", "", "vfs_mknod"},
		{objs.TraceVfsGetattr, "kprobe", "", "vfs_getattr"},
		{objs.TraceVfsReadlink, "kprobe", "", "vfs_readlink"},
		{objs.TraceDoFaccessat, "kprobe", "", "do_faccessat"},
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

	m.watchInodes = make(map[uint64][]watchedInode)
	m.populateWatchInodes()

	paths := make([]string, len(targets))
	for i, t := range targets {
		paths[i] = t.Path
	}
	log.Infof("monitor created \u2014 %d probes/tracepoints attached, watching: %s",
		len(m.links), strings.Join(paths, ", "))
	return m, nil
}

func (m *Monitor) Events() <-chan ebpf.FileEvent {
	return m.events
}

func (m *Monitor) SetEventTypes(types []ebpf.EventType) {
	if len(types) > 0 {
		m.eventTypes.Store(types)
	}
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
	log.Infof("monitor started \u2014 watching: %s", strings.Join(paths, ", "))
	go m.readLoop(rd)
	return nil
}

func (m *Monitor) readLoop(rd *ringbuf.Reader) {
	defer rd.Close()

	for {
		ev, ok := m.readEvent(rd)
		if !ok {
			return
		}
		if ev == nil {
			continue
		}

		select {
		case m.events <- *ev:
		case <-m.done:
			return
		}
	}
}

func (m *Monitor) readEvent(rd *ringbuf.Reader) (*ebpf.FileEvent, bool) {
	record, err := rd.Read()
	if err != nil {
		if errors.Is(err, ringbuf.ErrClosed) {
			return nil, false
		}
		log.Errorf("ringbuf read error: %v", err)
		return nil, true
	}

	var be ebpf.BpfEvent
	if !ebpf.DecodeBpfEvent(record.RawSample, &be) {
		log.Errorf("decode event: short record (%d bytes, want %d)", len(record.RawSample), ebpf.BpfEventSize)
		return nil, true
	}

	ev := be.ToFileEvent()
	ev.Timestamp = time.Now().UnixNano()

	if int(ev.PID) == m.ownPID || !m.eventTypeAllowed(ev.Type) {
		return nil, true
	}
	if !m.filterEvent(&ev, &be) {
		return nil, true
	}
	return &ev, true
}

// filterEvent completes ev's paths and reports whether it concerns a watch target. Every event on
// the system reaches it, so it must not touch the event's file: stat-ing other programs' files
// both costs a syscall per event and trips the daemon's guards on their protected ones.
func (m *Monitor) filterEvent(ev *ebpf.FileEvent, be *ebpf.BpfEvent) bool {
	kernelPath := be.Dev != 0
	if kernelPath {
		m.mapKernelPaths(ev, be.Dev)
	} else {
		ev.Path = m.resolveRelativePath(ev.PID, ev.Path)
		ev.Dest = m.resolveRelativePath(ev.PID, ev.Dest)
		m.resolveFDPath(ev)
	}

	// A kernel path names the object itself: only a symlink can lead into a target.
	if !kernelPath || be.Mode&syscall.S_IFMT == syscall.S_IFLNK {
		ev.Path = m.resolveSymlinkToTarget(ev.Path)
	}

	if m.inWatchPath(ev.Path, ev.Dest) {
		return m.matchesFilter(ev)
	}

	hardlinkPath := m.hardlinkByInode(be.Dev, be.Ino, ev.Path)
	if hardlinkPath == "" {
		return false
	}
	ev.Path = hardlinkPath
	return m.matchesFilter(ev)
}

// mapKernelPaths rewrites the superblock-relative paths the BPF walk yields to absolute ones.
// Unmapped paths (a mount outside this namespace) are left as they are.
func (m *Monitor) mapKernelPaths(ev *ebpf.FileEvent, dev uint32) {
	if m.mounts == nil {
		return
	}
	if p, ok := m.mounts.resolve(dev, ev.Path, m.matchesAnyTarget); ok {
		ev.Path = p
	}
	switch ev.Type {
	case ebpf.EventRename, ebpf.EventHardlink:
		if p, ok := m.mounts.resolve(dev, ev.Dest, m.matchesAnyTarget); ok {
			ev.Dest = p
		}
	case ebpf.EventSymlink:
		// Dest is the link body; a relative one is relative to the link's directory.
		if ev.Dest != "" && !filepath.IsAbs(ev.Dest) {
			ev.Dest = filepath.Join(filepath.Dir(ev.Path), ev.Dest)
		}
	}
}

// pathInfo returns the cached EvalSymlinks resolution for path. Nil-safe: a Monitor without a
// cache computes each time.
func (m *Monitor) pathInfo(path string) pathCacheEntry {
	if m.paths != nil {
		if e, ok := m.paths.lookup(path); ok {
			return e
		}
	}
	e := pathCacheEntry{resolved: evalSymlinksOrEmpty(path), at: time.Now()}
	if m.paths != nil {
		m.paths.store(path, e)
	}
	return e
}

func evalSymlinksOrEmpty(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return ""
	}
	return resolved
}

// hardlinkByInode returns the watched path sharing the event's inode, or "". Only an inode the
// kernel reported as a watched one is lstat'd, to confirm it: inode numbers repeat across
// filesystems, and across btrfs subvolumes of one superblock.
func (m *Monitor) hardlinkByInode(dev uint32, ino uint64, path string) string {
	if ino == 0 || path == "" {
		return ""
	}
	for _, w := range m.watchInodes[ino] {
		if w.sdev != 0 && w.sdev != dev {
			continue
		}
		var s syscall.Stat_t
		if syscall.Lstat(path, &s) == nil && s.Ino == ino && s.Dev == w.stDev {
			return w.path
		}
	}
	return ""
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

func (m *Monitor) resolveFDPath(ev *ebpf.FileEvent) {
	if ev.Type != ebpf.EventMmap {
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

func (m *Monitor) matchesFilter(ev *ebpf.FileEvent) bool {
	if ev.Path == "" {
		return false
	}

	for _, t := range m.targets {
		if m.matchTarget(ev.Path, t) {
			return true
		}
	}
	return false
}

func (m *Monitor) matchTarget(path string, t ebpf.Target) bool {
	if !t.IsDir {
		return path == t.File || strings.HasSuffix(path, "/"+t.File)
	}
	if !m.isWithinDir(path, t.Dir) {
		return false
	}
	if m.recursive && m.depth == 0 {
		return true
	}
	if !m.recursive {
		return m.matchDirectChildren(path, t.Dir)
	}
	return m.matchRecursiveDepth(path, t.Dir)
}

func (m *Monitor) matchDirectChildren(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return rel == "." || !strings.Contains(rel, string(filepath.Separator))
}

func (m *Monitor) matchRecursiveDepth(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	depth := 0
	if rel != "." {
		depth = len(strings.Split(rel, string(filepath.Separator)))
	}
	return depth <= m.depth
}

func (m *Monitor) eventTypeAllowed(et ebpf.EventType) bool {
	types := m.eventTypes.Load().([]ebpf.EventType)
	for _, allowed := range types {
		if et == allowed {
			return true
		}
	}
	return false
}

func (m *Monitor) isWithinDir(path, dir string) bool {
	return path == dir || strings.HasPrefix(path, dir+"/")
}

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
	w := watchedInode{path: path, stDev: s.Dev}
	if m.mounts != nil {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			w.sdev = m.mounts.devOf(resolved)
		}
	}
	m.watchInodes[s.Ino] = append(m.watchInodes[s.Ino], w)
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

// resolveSymlinkToTarget rewrites an event path that is a symlink resolving into a watched target;
// others are returned unchanged. EvalSymlinks is cached.
func (m *Monitor) resolveSymlinkToTarget(path string) string {
	if path == "" || m.matchesAnyTarget(path) {
		return path
	}
	if resolved := m.pathInfo(path).resolved; resolved != "" && m.matchesAnyTarget(resolved) {
		return resolved
	}
	return path
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
