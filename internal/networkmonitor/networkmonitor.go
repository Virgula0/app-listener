package networkmonitor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// BinaryEntry identifies a binary the network monitor keys on. It is the
// shared infrastructure type; this alias keeps the public API unchanged.
type BinaryEntry = ebpf.BinaryEntry

// ComputeBinaryEntry hashes a binary and derives its comm.
var ComputeBinaryEntry = ebpf.ComputeBinaryEntry

type NetworkMonitor struct {
	objs   NetMonObjects
	links  []link.Link
	events chan ebpf.NetEvent
	done   chan struct{}
	mu     sync.Mutex
	stop   bool

	binaries   []BinaryEntry
	ownPID     int
	eventTypes atomic.Value
}

func NewNetworkMonitor(binaries []BinaryEntry) (*NetworkMonitor, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	m := &NetworkMonitor{
		events:   make(chan ebpf.NetEvent, 1024),
		done:     make(chan struct{}),
		binaries: binaries,
		ownPID:   os.Getpid(),
	}
	m.eventTypes.Store(ebpf.NetEventTypes())

	var objs NetMonObjects
	if err := LoadNetMonObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("loading BPF objects: %w", err)
	}
	m.objs = objs

	if err := m.populateMaps(); err != nil {
		m.cleanup()
		return nil, fmt.Errorf("populating BPF maps: %w", err)
	}

	attachments := []struct {
		prog *cilium.Program
		typ  string
		opts any
	}{
		{objs.TraceConnect, "tracepoint", "syscalls/sys_enter_connect"},
		{objs.TraceBind, "tracepoint", "syscalls/sys_enter_bind"},
		{objs.TraceListen, "tracepoint", "syscalls/sys_enter_listen"},
		{objs.TraceAccept, "tracepoint", "syscalls/sys_enter_accept"},
		{objs.TraceAccept4, "tracepoint", "syscalls/sys_enter_accept4"},
		{objs.TraceSendto, "tracepoint", "syscalls/sys_enter_sendto"},
		{objs.TraceRecvfrom, "tracepoint", "syscalls/sys_enter_recvfrom"},
		{objs.TraceSendmsg, "tracepoint", "syscalls/sys_enter_sendmsg"},
		{objs.TraceRecvmsg, "tracepoint", "syscalls/sys_enter_recvmsg"},
		{objs.TraceClose, "tracepoint", "syscalls/sys_enter_close"},
		{objs.TraceInetCskAccept, "kretprobe", "inet_csk_accept"},
	}

	for _, a := range attachments {
		var l link.Link
		var err error

		switch a.typ {
		case "tracepoint":
			parts := strings.SplitN(a.opts.(string), "/", 2)
			l, err = link.Tracepoint(parts[0], parts[1], a.prog, nil)
		case "kretprobe":
			l, err = link.Kretprobe(a.opts.(string), a.prog, nil)
		}

		if err != nil {
			log.Warnf("skipping %s/%v: %v", a.typ, a.opts, err)
			continue
		}
		m.links = append(m.links, l)
	}

	log.Infof("network monitor created \u2014 %d/%d probes/tracepoints attached, watching: %s",
		len(m.links), len(attachments), binariesSummary(binaries))
	return m, nil
}

func (m *NetworkMonitor) populateMaps() error {
	for _, b := range m.binaries {
		var val uint8 = 1

		ck := &NetMonCommKey{}
		commBytes := []byte(b.Comm)
		for i := 0; i < len(ck.Comm) && i < len(commBytes); i++ {
			ck.Comm[i] = int8(commBytes[i])
		}
		if err := m.objs.WatchComms.Put(ck, val); err != nil {
			return fmt.Errorf("storing comm for %s: %w", b.Path, err)
		}

		ik, err := statInodeKey(b.Path)
		if err != nil {
			log.Warnf("cannot stat binary %s: %v \u2014 using comm-only matching for %q", b.Path, err, b.Comm)
			continue
		}
		if err := m.objs.WatchExeInodes.Put(ik, val); err != nil {
			return fmt.Errorf("storing exe inode for %s: %w", b.Path, err)
		}
	}
	return nil
}

func statInodeKey(path string) (*NetMonInodeKey, error) {
	dev, ino, err := ebpf.StatInode(path)
	if err != nil {
		return nil, err
	}
	return &NetMonInodeKey{
		Dev: dev,
		Ino: ino,
	}, nil
}

func (m *NetworkMonitor) Events() <-chan ebpf.NetEvent {
	return m.events
}

func (m *NetworkMonitor) SetEventTypes(types []ebpf.NetEventType) {
	if len(types) > 0 {
		m.eventTypes.Store(types)
	}
}

func (m *NetworkMonitor) Start() error {
	rd, err := ringbuf.NewReader(m.objs.Rb)
	if err != nil {
		m.cleanup()
		return fmt.Errorf("ringbuf reader: %w", err)
	}

	log.Infof("network monitor started \u2014 watching %d binaries", len(m.binaries))
	go m.readLoop(rd)
	return nil
}

func (m *NetworkMonitor) readLoop(rd *ringbuf.Reader) {
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

func (m *NetworkMonitor) readEvent(rd *ringbuf.Reader) (*ebpf.NetEvent, bool) {
	record, err := rd.Read()
	if err != nil {
		if errors.Is(err, ringbuf.ErrClosed) {
			return nil, false
		}
		log.Errorf("ringbuf read error: %v", err)
		return nil, true
	}

	var be struct {
		PID      uint32
		UID      uint32
		GID      uint32
		Type     uint32
		Proto    uint32
		Size     uint32
		FD       uint32
		AF       uint32
		Saddr    [4]uint32
		Daddr    [4]uint32
		Sport    uint16
		Dport    uint16
		Comm     [16]byte
		TID      uint32
		NetNS    uint64
		CgroupID uint64
	}
	if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &be); err != nil {
		log.Errorf("decode net event: %v", err)
		return nil, true
	}

	if int(be.PID) == m.ownPID {
		return nil, true
	}

	ev := ebpf.NetEvent{
		PID:       be.PID,
		TID:       be.TID,
		UID:       be.UID,
		GID:       be.GID,
		Type:      ebpf.NetEventType(be.Type),
		Protocol:  be.Proto,
		Size:      be.Size,
		FD:        be.FD,
		Comm:      ebpf.Cstr(be.Comm[:]),
		NetNS:     be.NetNS,
		CgroupID:  be.CgroupID,
		Timestamp: time.Now().UnixNano(),
	}

	ev.SrcAddr = ebpf.FormatAddr(be.AF, be.Saddr[:], be.Sport)
	ev.DstAddr = ebpf.FormatAddr(be.AF, be.Daddr[:], be.Dport)

	if !m.eventTypeAllowed(ev.Type) {
		return nil, true
	}

	return &ev, true
}

func (m *NetworkMonitor) eventTypeAllowed(et ebpf.NetEventType) bool {
	types := m.eventTypes.Load().([]ebpf.NetEventType)
	for _, allowed := range types {
		if et == allowed {
			return true
		}
	}
	return false
}

func (m *NetworkMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stop {
		return
	}
	m.stop = true
	close(m.done)
	m.cleanup()
}

func (m *NetworkMonitor) cleanup() {
	for _, l := range m.links {
		l.Close()
	}
	m.links = nil
	m.objs.Close()
}

func binariesSummary(binaries []BinaryEntry) string {
	parts := make([]string, len(binaries))
	for i, b := range binaries {
		parts[i] = fmt.Sprintf("%s [sha256:%x..%x]", b.Path, b.Hash[:4], b.Hash[28:])
	}
	return strings.Join(parts, ", ")
}
