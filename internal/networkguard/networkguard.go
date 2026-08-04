package networkguard

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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

type NetGuardEvent struct {
	ebpf.NetEvent
	Blocked bool
}

type NetGuard struct {
	objs   GuardNetObjects
	links  []link.Link
	events chan NetGuardEvent
	done   chan struct{}
	mu     sync.Mutex
	stop   bool

	mode     Mode
	binaries []BinaryEntry
	unsafe   bool
	eventset []ebpf.NetEventType
	ownPID   int
	throttle bool
}

func NewNetGuard(mode Mode, binaries []BinaryEntry, eventset []ebpf.NetEventType, unsafe, throttle bool) (*NetGuard, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Warnf("failed to remove memlock rlimit: %v", err)
	}

	bins := make([]BinaryEntry, len(binaries))
	copy(bins, binaries)

	es := eventset
	if len(es) == 0 {
		es = ebpf.NetEventTypes()
	}

	g := &NetGuard{
		events:   make(chan NetGuardEvent, 1024),
		done:     make(chan struct{}),
		mode:     mode,
		binaries: bins,
		unsafe:   unsafe,
		eventset: es,
		ownPID:   os.Getpid(),
		throttle: throttle,
	}

	var objs GuardNetObjects
	if err := LoadGuardNetObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("loading network guard BPF objects: %w", err)
	}
	g.objs = objs

	if err := g.populateMaps(); err != nil {
		g.cleanup()
		return nil, fmt.Errorf("populating BPF maps: %w", err)
	}

	required := map[string]bool{
		"socket_connect": true,
		"socket_bind":    true,
		"socket_listen":  true,
	}

	attachments := []struct {
		prog *cilium.Program
		hook string
	}{
		{g.objs.GuardNetSocketConnect, "socket_connect"},
		{g.objs.GuardNetSocketBind, "socket_bind"},
		{g.objs.GuardNetSocketListen, "socket_listen"},
		{g.objs.GuardNetSocketSendmsg, "socket_sendmsg"},
		{g.objs.GuardNetSocketRecvmsg, "socket_recvmsg"},
	}

	var failedRequired []string
	for _, a := range attachments {
		l, err := link.AttachLSM(link.LSMOptions{
			Program: a.prog,
		})
		if err != nil {
			if required[a.hook] {
				failedRequired = append(failedRequired, a.hook)
				log.Errorf("CRITICAL: required LSM hook %s failed to attach: %v", a.hook, err)
			} else {
				log.Warnf("skipping optional LSM hook %s: %v", a.hook, err)
			}
			continue
		}
		g.links = append(g.links, l)
	}

	if len(failedRequired) > 0 {
		g.cleanup()
		return nil, fmt.Errorf(
			"required LSM hooks failed to attach: %v — network blocking unavailable; "+
				"ensure your kernel supports BPF LSM (CONFIG_BPF_LSM=y) and LSM=bpf is in the "+
				"boot command line (/sys/kernel/security/lsm)",
			failedRequired)
	}

	modeLabel := "blacklist"
	if mode == ModeWhitelist {
		modeLabel = "whitelist"
	}
	unsafeLabel := ""
	if unsafe {
		unsafeLabel = " [UNSAFE]"
	}
	log.Infof("network guard created (mode=%s%s) \u2014 %d/%d LSM hooks attached, binaries: %d, events: %v",
		modeLabel, unsafeLabel, len(g.links), len(attachments), len(binaries), eventsetSummary(es))
	return g, nil
}

func eventsetSummary(es []ebpf.NetEventType) string {
	parts := make([]string, len(es))
	for i, et := range es {
		parts[i] = et.String()
	}
	return strings.Join(parts, ",")
}

const (
	bpfBlock uint8 = 1
	bpfAllow uint8 = 2

	// Default actions for guard_net_config[0]
	defaultAllow uint64 = 0
	defaultBlock uint64 = 1
)

func (g *NetGuard) populateMaps() error {
	for _, et := range g.eventset {
		key, err := eventTypeKey(et)
		if err != nil {
			return err
		}
		if err := g.objs.GuardNetEvents.Put(key, uint64(1)); err != nil {
			return fmt.Errorf("setting event type %d in filter: %w", et, err)
		}
	}

	// Set default action: whitelist mode → block by default, blacklist → allow by default
	defaultAction := defaultAllow
	if g.mode == ModeWhitelist {
		defaultAction = defaultBlock
	}
	if err := g.objs.GuardNetConfig.Put(uint32(0), defaultAction); err != nil {
		return fmt.Errorf("setting default action: %w", err)
	}

	// Config[1]: blocking_enabled — always on for both blacklist and whitelist modes
	if err := g.objs.GuardNetConfig.Put(uint32(1), uint64(1)); err != nil {
		return fmt.Errorf("setting blocking enabled: %w", err)
	}

	// Config[2]: unsafe_families — guard AF_UNIX too (only with --unsafe in whitelist mode)
	unsafeFamilies := uint64(0)
	if g.unsafe {
		unsafeFamilies = 1
	}
	if err := g.objs.GuardNetConfig.Put(uint32(2), unsafeFamilies); err != nil {
		return fmt.Errorf("setting unsafe families: %w", err)
	}

	// Config[3]: throttle_enabled — rate-limit events per (type, comm) to
	// protect the ring buffer from flooding; disabled with --no-throttle.
	throttleEnabled := uint64(0)
	if g.throttle {
		throttleEnabled = 1
	}
	if err := g.objs.GuardNetConfig.Put(uint32(3), throttleEnabled); err != nil {
		return fmt.Errorf("setting throttle enabled: %w", err)
	}

	bpfVal := bpfBlock
	if g.mode == ModeWhitelist {
		bpfVal = bpfAllow
	}
	for _, b := range g.binaries {
		ik, err := statInodeKey(b.Path)
		if err != nil {
			return fmt.Errorf("stating binary %s: %w", b.Path, err)
		}
		if err := g.objs.GuardNetExeActions.Put(ik, bpfVal); err != nil {
			return fmt.Errorf("storing exe inode for %s: %w", b.Path, err)
		}
	}

	return nil
}

// eventTypeKey converts a parsed event type to the uint32 key used in the BPF
// filter map, rejecting out-of-range values.
func eventTypeKey(et ebpf.NetEventType) (uint32, error) {
	if et < 0 || et >= 0x100 {
		return 0, fmt.Errorf("event type %d out of range", et)
	}
	return uint32(et), nil //nolint:gosec // et is range-checked to [0, 256) above
}

func statInodeKey(path string) (*GuardNetInodeKey, error) {
	var es syscall.Stat_t
	if err := syscall.Stat(path, &es); err != nil {
		return nil, err
	}
	return &GuardNetInodeKey{
		Dev: uint64(mkdev(unix.Major(es.Dev), unix.Minor(es.Dev))),
		Ino: es.Ino,
	}, nil
}

func mkdev(major, minor uint32) uint32 {
	return (major << 20) | minor
}

var infraBinaries = []string{
	"/usr/lib/systemd/systemd-resolved",
	"/usr/bin/NetworkManager",
	"/usr/sbin/NetworkManager",
	"/usr/lib/systemd/systemd-networkd",
}

// DiscoverInfraBinaries returns the paths of essential system networking daemons
// that are currently running. In whitelist mode (especially with --auto-infra)
// these must be allowed so that name resolution and connection management keep
// working on behalf of whitelisted applications.
func DiscoverInfraBinaries() ([]string, error) {
	running, err := runningExecutables()
	if err != nil {
		return nil, fmt.Errorf("scanning running processes: %w", err)
	}
	return infraFromRunning(running), nil
}

func infraFromRunning(running map[string]bool) []string {
	var found []string
	seen := make(map[string]bool)
	for _, path := range infraBinaries {
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			continue
		}
		if !running[resolved] || seen[resolved] {
			continue
		}
		seen[resolved] = true
		found = append(found, path)
	}
	return found
}

func runningExecutables() (map[string]bool, error) {
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	running := make(map[string]bool)
	for _, dir := range procs {
		if !dir.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(dir.Name()); err != nil {
			continue
		}
		exe, err := os.Readlink(fmt.Sprintf("/proc/%s/exe", dir.Name()))
		if err != nil {
			continue
		}
		resolved, err := filepath.EvalSymlinks(exe)
		if err != nil {
			continue
		}
		running[resolved] = true
	}
	return running, nil
}

func (g *NetGuard) Events() <-chan NetGuardEvent {
	return g.events
}

func (g *NetGuard) Start() error {
	rd, err := ringbuf.NewReader(g.objs.GuardNetRb)
	if err != nil {
		g.cleanup()
		return fmt.Errorf("ringbuf reader: %w", err)
	}

	modeLabel := "blacklist"
	if g.mode == ModeWhitelist {
		modeLabel = "whitelist"
	}
	unsafeLabel := ""
	if g.unsafe {
		unsafeLabel = " [UNSAFE]"
	}
	log.Infof("network guard started (mode=%s%s) \u2014 %d binaries", modeLabel, unsafeLabel, len(g.binaries))
	go g.readLoop(rd)
	return nil
}

func (g *NetGuard) readLoop(rd *ringbuf.Reader) {
	defer rd.Close()

	for {
		ev, ok := g.readEvent(rd)
		if !ok {
			return
		}
		if ev == nil {
			continue
		}

		select {
		case g.events <- *ev:
		case <-g.done:
			return
		}
	}
}

func (g *NetGuard) readEvent(rd *ringbuf.Reader) (*NetGuardEvent, bool) {
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
		Blocked  uint32
	}
	if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &be); err != nil {
		log.Errorf("decode net guard event: %v", err)
		return nil, true
	}

	if int(be.PID) == g.ownPID {
		return nil, true
	}

	ev := NetGuardEvent{
		NetEvent: ebpf.NetEvent{
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
			Timestamp: 0,
		},
		Blocked: be.Blocked != 0,
	}

	ev.SrcAddr = ebpf.FormatAddr(be.AF, be.Saddr[:], be.Sport)
	ev.DstAddr = ebpf.FormatAddr(be.AF, be.Daddr[:], be.Dport)

	return &ev, true
}

func (g *NetGuard) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.stop {
		return
	}
	g.stop = true
	close(g.done)
	g.cleanup()
}

func (g *NetGuard) cleanup() {
	for _, l := range g.links {
		l.Close()
	}
	g.links = nil
	g.objs.Close()
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
