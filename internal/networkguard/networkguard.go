package networkguard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

type Mode int

const (
	ModeBlacklist Mode = iota
	ModeWhitelist
)

// BinaryEntry identifies a binary the network guard keys on. It is the
// shared infrastructure type; this alias keeps the public API unchanged.
type BinaryEntry = ebpf.BinaryEntry

// ComputeBinaryEntry hashes a binary and derives its comm.
var ComputeBinaryEntry = ebpf.ComputeBinaryEntry

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

	mode         Mode
	binaries     []BinaryEntry
	unsafe       bool
	eventset     []ebpf.NetEventType
	ownPID       int
	throttle     bool
	hasMulticall bool // a uutils multicall applet is configured: the exec tagger must attach
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

	hooks := len(g.links)
	if err := g.attachMulticallTagger(); err != nil {
		g.cleanup()
		return nil, err
	}

	log.Infof("network guard created (mode=%s) \u2014 %d/%d LSM hooks attached, binaries: %d, events: %v",
		modeLabel(mode, unsafe), hooks, len(attachments), len(binaries), eventsetSummary(es))
	return g, nil
}

func modeLabel(mode Mode, unsafe bool) string {
	l := "blacklist"
	if mode == ModeWhitelist {
		l = "whitelist"
	}
	if unsafe {
		l += " [UNSAFE]"
	}
	return l
}

// attachMulticallTagger attaches the uutils applet attestation programs (mc_tag.h) when a multicall
// applet is configured: exec stamps each multicall process with the applet it ran as, fork inherits
// it, task_free drops it. All three are required then: without the tagger every multicall process
// keys as MC_TAG_NONE (a whitelisted applet silently denied), and without task_free leaked stamps
// fill the map until no exec can be attested.
func (g *NetGuard) attachMulticallTagger() error {
	if !g.hasMulticall {
		return nil
	}
	progs := []struct {
		prog *cilium.Program
		name string
		lsm  bool
	}{
		{g.objs.NetgExecApplet, "sched_process_exec", false},
		{g.objs.NetgSchedFork, "sched_process_fork", false},
		{g.objs.NetgTaskFree, "task_free", true},
	}
	for _, a := range progs {
		var l link.Link
		var err error
		if a.lsm {
			l, err = link.AttachLSM(link.LSMOptions{Program: a.prog})
		} else {
			l, err = link.AttachTracing(link.TracingOptions{Program: a.prog})
		}
		if err != nil {
			return fmt.Errorf("attaching multicall applet tagger (%s): %w — a uutils applet is "+
				"configured but its per-applet identity cannot be attested on this kernel", a.name, err)
		}
		g.links = append(g.links, l)
	}
	return nil
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

	return g.putBinaries()
}

// putBinaries stores every binary's inode key with the mode's action. A uutils multicall applet is
// keyed by its attested applet (guard.ExeKey), so whitelisting one applet does not admit the other
// ~108; its inode's name table is written first (the kernel tagger needs it). An opaque multicall
// (busybox/toybox, or a uutils build whose applets can't be told apart) is refused — one identity
// there would admit every applet. Non-multicall binaries keep their plain inode key, unchanged.
func (g *NetGuard) putBinaries() error {
	if err := g.putBtrfsLayout(); err != nil {
		return err
	}
	bpfVal := bpfBlock
	if g.mode == ModeWhitelist {
		bpfVal = bpfAllow
	}
	for _, b := range g.binaries {
		keys, err := g.exeKeys(b.Path)
		if err != nil {
			return err
		}
		for _, k := range keys {
			if err := g.objs.GuardNetExeActions.Put(k, bpfVal); err != nil {
				return fmt.Errorf("storing exe inode for %s: %w", b.Path, err)
			}
		}
	}

	return nil
}

// exeKeys returns the rows path's action is stored under, registering a multicall's applet table
// first. A blacklisted applet also gets its inode's unattested key (fexecve, exec -a, a full stamp
// map): uutils may still run the blacklisted applet from argv[0].
func (g *NetGuard) exeKeys(path string) ([]GuardNetInodeKey, error) {
	ik, err := statInodeKey(path)
	if err != nil {
		return nil, fmt.Errorf("stating binary %s: %w", path, err)
	}
	file := guard.GuardInodeKey{Dev: ik.Dev, Ino: ik.Ino}
	exe, err := guard.ExeKey(path, nil, file)
	if err != nil {
		return nil, fmt.Errorf("resolving binary %s: %w", path, err)
	}
	keys := []GuardNetInodeKey{{Dev: exe.Dev, Ino: exe.Ino}}
	if !exe.Tagged() {
		return keys, nil
	}
	if err := g.syncMulticall(file); err != nil {
		return nil, fmt.Errorf("registering multicall binary %s: %w", path, err)
	}
	if g.mode == ModeBlacklist {
		u := file.Unattested()
		keys = append(keys, GuardNetInodeKey{Dev: u.Dev, Ino: u.Ino})
	}
	return keys, nil
}

// syncMulticall writes the multicall inode's applet name table, then its presence row, before any
// tagged exe-action row of it: a process of a multicall missing from the presence map keys as NONE.
func (g *NetGuard) syncMulticall(file guard.GuardInodeKey) error {
	tags, err := guard.MulticallNameTags(file)
	if err != nil {
		return err
	}
	exe := GuardNetInodeKey{Dev: file.Dev, Ino: file.Ino}
	for name, tag := range tags {
		if err := g.objs.McMulticallNames.Put(mcNameKey(exe, name), tag); err != nil {
			return fmt.Errorf("writing multicall applet %s: %w", name, err)
		}
	}
	g.hasMulticall = true
	return g.objs.McMulticall.Put(exe, uint32(1))
}

// mcNameKey encodes an applet name into the kernel name-table key (mc_tag.h mc_name_key).
func mcNameKey(exe GuardNetInodeKey, name string) GuardNetMcNameKey {
	k := GuardNetMcNameKey{Exe: exe}
	for i := 0; i < len(name) && i < len(k.Name)-1; i++ {
		k.Name[i] = int8(name[i]) //nolint:gosec // applet names are ASCII (validAppletName)
	}
	return k
}

// putBtrfsLayout lets inode_dev key btrfs inodes by subvolume, as stat does. Without it no btrfs
// binary matches its entry, so one configured there refuses the start: a blacklisted one would
// connect freely.
func (g *NetGuard) putBtrfsLayout() error {
	l, err := ebpf.ResolveBtrfsLayout()
	if err != nil {
		paths := make([]string, len(g.binaries))
		for i, b := range g.binaries {
			paths[i] = b.Path
		}
		if p := ebpf.FirstOnBtrfs(paths); p != "" {
			return fmt.Errorf("%s is on btrfs, but the kernel's BTF does not describe btrfs: %w", p, err)
		}
		return nil
	}
	v := GuardNetBtrfsLayout{Valid: l.Valid, InodeRoot: l.InodeRoot, InodeVfs: l.InodeVfs, RootAnonDev: l.RootAnonDev}
	if err := g.objs.BtrfsLayout.Put(uint32(0), v); err != nil {
		return fmt.Errorf("storing the btrfs layout: %w", err)
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
	dev, ino, err := ebpf.StatInode(path)
	if err != nil {
		return nil, err
	}
	return &GuardNetInodeKey{
		Dev: dev,
		Ino: ino,
	}, nil
}

var infraBinaries = []string{
	"/usr/lib/systemd/systemd-resolved",
	"/usr/bin/NetworkManager",
	"/usr/sbin/NetworkManager",
	"/usr/lib/systemd/systemd-networkd",
}

// DiscoverInfraBinaries returns the paths of running essential system networking daemons. In
// whitelist mode (especially --auto-infra) they must be allowed so name resolution and connection
// management keep working for whitelisted apps.
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
			PID:      be.PID,
			TID:      be.TID,
			UID:      be.UID,
			GID:      be.GID,
			Type:     ebpf.NetEventType(be.Type),
			Protocol: be.Proto,
			Size:     be.Size,
			FD:       be.FD,
			Comm:     ebpf.Cstr(be.Comm[:]),
			NetNS:    be.NetNS,
			CgroupID: be.CgroupID,
			// Stamp once at ingestion: dual TUI consumers (local and
			// browser) must display the same time for the same event.
			Timestamp: time.Now().UnixNano(),
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
