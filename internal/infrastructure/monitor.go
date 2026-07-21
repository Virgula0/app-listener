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
		{objs.TraceDoSysOpen, "kprobe", "", "do_sys_open"},
		{objs.TraceDoSysOpenat2, "kprobe", "", "do_sys_openat2"},
		{objs.TraceKsysRead, "kprobe", "", "ksys_read"},
		{objs.TraceKsysWrite, "kprobe", "", "ksys_write"},
		{objs.TraceOpen, "tracepoint", "syscalls", "sys_enter_open"},
		{objs.TraceOpenat, "tracepoint", "syscalls", "sys_enter_openat"},
		{objs.TraceUnlinkat, "tracepoint", "syscalls", "sys_enter_unlinkat"},
		{objs.TraceRenameat2, "tracepoint", "syscalls", "sys_enter_renameat2"},
		{objs.TraceSymlinkat, "tracepoint", "syscalls", "sys_enter_symlinkat"},
		{objs.TraceLinkat, "tracepoint", "syscalls", "sys_enter_linkat"},
		{objs.TraceMkdirat, "tracepoint", "syscalls", "sys_enter_mkdirat"},
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

		if !m.inWatchPath(ev.Path, ev.Dest) {
			continue
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
	if ev.Type != EventRead && ev.Type != EventWrite {
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
