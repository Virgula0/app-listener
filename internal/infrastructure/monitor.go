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

type Monitor struct {
	objs    MonitorObjects
	links   []link.Link
	events  chan FileEvent
	done    chan struct{}
	mu      sync.Mutex
	stopped bool

	watchDir  string
	recursive bool
	depth     int
	ownPID    int
}

func NewMonitor(watchDir string, recursive bool, depth int) (*Monitor, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	m := &Monitor{
		events:    make(chan FileEvent, 1024),
		done:      make(chan struct{}),
		watchDir:  watchDir,
		recursive: recursive,
		depth:     depth,
		ownPID:    os.Getpid(),
	}

	var objs MonitorObjects
	if err := LoadMonitorObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("loading BPF objects: %w", err)
	}
	m.objs = objs

	type kprobeDef struct {
		prog *cilium.Program
		name string
	}

	kprobes := []kprobeDef{
		{objs.TraceDoSysOpen, "do_sys_open"},
		{objs.TraceDoSysOpenat2, "do_sys_openat2"},
		{objs.TraceKsysRead, "ksys_read"},
		{objs.TraceKsysWrite, "ksys_write"},
		{objs.TraceUnlinkat, "__x64_sys_unlinkat"},
		{objs.TraceRenameat2, "__x64_sys_renameat2"},
		{objs.TraceSymlinkat, "__x64_sys_symlinkat"},
		{objs.TraceLinkat, "__x64_sys_linkat"},
		{objs.TraceMkdirat, "__x64_sys_mkdirat"},
	}

	for _, k := range kprobes {
		kp, err := link.Kprobe(k.name, k.prog, nil)
		if err != nil {
			m.cleanup()
			return nil, fmt.Errorf("kprobe %s: %w", k.name, err)
		}
		m.links = append(m.links, kp)
	}

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

	log.Infof("monitor started — %d kprobes attached, watching %s", len(m.links), m.watchDir)
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

		m.resolveFDPath(&ev)

		if !m.inWatchDir(ev.Path) {
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

func (m *Monitor) inWatchDir(path string) bool {
	if path == "" {
		return false
	}
	return strings.HasPrefix(path, m.watchDir)
}

func (m *Monitor) matchesFilter(ev *FileEvent) bool {
	path := ev.Path

	if !m.recursive {
		rel, err := filepath.Rel(m.watchDir, path)
		if err != nil {
			return false
		}
		if strings.Contains(rel, string(filepath.Separator)) {
			return false
		}
	}

	if m.recursive && m.depth > 0 {
		rel, err := filepath.Rel(m.watchDir, path)
		if err != nil {
			return false
		}
		depth := 0
		if rel != "." {
			depth = len(strings.Split(rel, string(filepath.Separator)))
		}
		if depth > m.depth {
			return false
		}
	}

	return true
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
