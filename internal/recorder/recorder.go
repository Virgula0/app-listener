// Package recorder learns a whitelist for paths the daemon doesn't guard yet: while it records,
// every access to the targets is grouped by the executable inode the kernel reports for the acting
// process (never its comm, argv or a path it names), together with the code that executable maps,
// so `edit-protected --record` can propose daemon.conf sections. A proposal is only a suggestion:
// anything that touched a target during the window, malware included, is in it, so the operator
// reviews every line and Classify marks what deserves suspicion.
package recorder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/monitor"
)

// MaxDuration ends a recording nobody stopped.
const MaxDuration = 2 * time.Hour

// Bounds on what one recording keeps: past them new executables or libraries are dropped (counted
// in Recording.Dropped), never an unbounded map.
const (
	maxAccessors   = 1024
	maxMappers     = 4096
	maxLibsPerExe  = 1024
	maxRecentSeens = 64
)

// Target is one path being recorded.
type Target struct {
	Path  string
	IsDir bool
}

type exeKey struct{ dev, ino uint64 }

// accessKey keys an accessor by executable and the name it was exec'd by: a uutils multicall runs
// many applets from one inode, and its applet is part of its identity (guard.ExeKey).
type accessKey struct {
	exe  exeKey
	name string
}

type accessor struct {
	exe     string
	comms   map[string]struct{}
	ops     map[ebpf.EventType]struct{}
	targets map[string]struct{}
	uids    map[uint32]struct{}
	count   uint64
	first   time.Time
	last    time.Time
}

// source is the event stream (monitor.Monitor; a fake in tests).
type source interface {
	Events() <-chan ebpf.FileEvent
	Stop()
}

// Recorder aggregates one recording.
type Recorder struct {
	src     source
	targets []Target
	self    exeKey
	started time.Time

	mu      sync.Mutex
	events  uint64
	dropped uint64
	access  map[accessKey]*accessor
	maps    map[exeKey]map[string]struct{}

	stop     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
	stopped  time.Time
}

// Start records accesses to targets (validated by ValidateTargets) until Stop or MaxDuration.
func Start(targets []Target) (*Recorder, error) {
	if len(targets) == 0 {
		return nil, errors.New("nothing to record")
	}
	self, err := selfExe()
	if err != nil {
		return nil, err
	}
	mon, err := monitor.NewMonitor(monitorTargets(targets), true, 0)
	if err != nil {
		return nil, fmt.Errorf("loading the monitor: %w", err)
	}
	mon.RecordExecMappings(true)
	if err := mon.AttachRecordProbes(); err != nil {
		mon.Stop()
		return nil, fmt.Errorf("loading the monitor: %w", err)
	}
	if err := mon.Start(); err != nil {
		mon.Stop()
		return nil, fmt.Errorf("starting the monitor: %w", err)
	}
	r := newRecorder(mon, targets, self)
	go r.loop()
	return r, nil
}

func newRecorder(src source, targets []Target, self exeKey) *Recorder {
	return &Recorder{
		src: src, targets: targets, self: self, started: time.Now(),
		access: map[accessKey]*accessor{}, maps: map[exeKey]map[string]struct{}{},
		stop: make(chan struct{}), done: make(chan struct{}),
	}
}

// monitorTargets names a file target by its absolute path: the monitor matches a bare basename
// anywhere on the system.
func monitorTargets(targets []Target) []ebpf.Target {
	out := make([]ebpf.Target, len(targets))
	for i, t := range targets {
		if t.IsDir {
			out[i] = ebpf.Target{Path: t.Path, Dir: t.Path, IsDir: true}
		} else {
			out[i] = ebpf.Target{Path: t.Path, Dir: filepath.Dir(t.Path), File: t.Path}
		}
	}
	return out
}

// selfExe is this binary's inode: the daemon and this process run it, and neither is recorded.
func selfExe() (exeKey, error) {
	dev, ino, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return exeKey{}, fmt.Errorf("reading this binary's identity: %w", err)
	}
	return exeKey{dev, ino}, nil
}

func (r *Recorder) loop() {
	defer close(r.done)
	limit := time.NewTimer(MaxDuration)
	defer limit.Stop()
	events := r.src.Events()
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			r.observe(&ev)
		case <-limit.C:
			return
		case <-r.stop:
			return
		}
	}
}

// observe files one event: a code mapping under its executable, an access to a target under its
// accessor. Kernel threads and this binary are never recorded.
func (r *Recorder) observe(ev *ebpf.FileEvent) {
	k := exeKey{ev.ExeDev, ev.ExeIno}
	if k.ino == 0 || k == r.self {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if ev.Type == ebpf.EventMmap && ev.Prot&unix.PROT_EXEC != 0 && filepath.IsAbs(ev.Path) {
		r.noteMapping(k, ev.Path)
	}
	target := r.targetOf(ev.Path)
	if target == "" {
		target = r.targetOf(ev.Dest)
	}
	if target == "" {
		return
	}
	r.events++
	ak := accessKey{exe: k, name: ev.ExecName}
	a := r.access[ak]
	if a == nil {
		if len(r.access) >= maxAccessors {
			r.dropped++
			return
		}
		a = &accessor{comms: map[string]struct{}{}, ops: map[ebpf.EventType]struct{}{},
			targets: map[string]struct{}{}, uids: map[uint32]struct{}{}, first: time.Now()}
		r.access[ak] = a
	}
	if ev.Exe != "" {
		a.exe = ev.Exe
	}
	a.comms[ev.Comm] = struct{}{}
	a.ops[ev.Type] = struct{}{}
	a.targets[target] = struct{}{}
	a.uids[ev.UID] = struct{}{}
	a.count++
	a.last = time.Now()
}

func (r *Recorder) noteMapping(k exeKey, lib string) {
	libs := r.maps[k]
	if libs == nil {
		if len(r.maps) >= maxMappers {
			r.dropped++
			return
		}
		libs = map[string]struct{}{}
		r.maps[k] = libs
	}
	if _, ok := libs[lib]; !ok && len(libs) >= maxLibsPerExe {
		r.dropped++
		return
	}
	libs[lib] = struct{}{}
}

// targetOf is the target p lies in, or "".
func (r *Recorder) targetOf(p string) string {
	if p == "" {
		return ""
	}
	for _, t := range r.targets {
		if p == t.Path || (t.IsDir && strings.HasPrefix(p, t.Path+"/")) {
			return t.Path
		}
	}
	return ""
}

// Seen is one executable in a live progress view (display only).
type Seen struct {
	Exe   string
	Comm  string
	Ops   []string
	Count uint64
	Last  time.Time
}

// Progress is a live view of a running recording.
type Progress struct {
	Started  time.Time
	Elapsed  time.Duration
	Events   uint64
	Binaries int
	Recent   []Seen
}

// Progress reports the recording so far, most recently active executables first.
func (r *Recorder) Progress() Progress {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := Progress{Started: r.started, Elapsed: time.Since(r.started), Events: r.events, Binaries: len(r.access)}
	for k, a := range r.access {
		exe := a.exe
		if exe == "" {
			exe = fmt.Sprintf("(unnamed executable %d:%d)", k.exe.dev, k.exe.ino)
		}
		if k.name != "" && filepath.Base(k.name) != filepath.Base(exe) {
			exe += " as " + filepath.Base(k.name)
		}
		p.Recent = append(p.Recent, Seen{Exe: exe, Comm: strings.Join(sortedKeys(a.comms), ","),
			Ops: opNames(a.ops), Count: a.count, Last: a.last})
	}
	sort.Slice(p.Recent, func(i, j int) bool { return p.Recent[i].Last.After(p.Recent[j].Last) })
	if len(p.Recent) > maxRecentSeens {
		p.Recent = p.Recent[:maxRecentSeens]
	}
	return p
}

// Stop ends the recording and returns what it saw, unclassified. Idempotent.
func (r *Recorder) Stop() *Raw {
	r.stopOnce.Do(func() {
		close(r.stop)
		<-r.done
		r.src.Stop()
		r.stopped = time.Now()
	})
	r.mu.Lock()
	defer r.mu.Unlock()
	raw := &Raw{Targets: r.targets, Started: r.started, Stopped: r.stopped, Events: r.events, Dropped: r.dropped}
	for k, a := range r.access {
		raw.Accessors = append(raw.Accessors, RawAccessor{
			Dev: k.exe.dev, Ino: k.exe.ino, ExecName: k.name, Exe: a.exe, Comms: sortedKeys(a.comms),
			Ops: opsOf(a.ops), Targets: sortedKeys(a.targets), UIDs: sortedUIDs(a.uids), Count: a.count,
			First: a.first, Last: a.last, Libs: sortedKeys(r.maps[k.exe]),
		})
	}
	sort.Slice(raw.Accessors, func(i, j int) bool {
		a, b := raw.Accessors[i], raw.Accessors[j]
		if a.Exe != b.Exe {
			return a.Exe < b.Exe
		}
		return a.ExecName < b.ExecName
	})
	return raw
}

// Raw is a stopped recording before Classify.
type Raw struct {
	Targets   []Target
	Started   time.Time
	Stopped   time.Time
	Events    uint64
	Dropped   uint64
	Accessors []RawAccessor
}

// RawAccessor is one (executable inode, exec name) that touched a target.
type RawAccessor struct {
	Dev, Ino uint64
	ExecName string
	Exe      string
	Comms    []string
	Ops      []ebpf.EventType
	Targets  []string
	UIDs     []uint32
	Count    uint64
	First    time.Time
	Last     time.Time
	// Libs are the files this executable mapped as code during the recording.
	Libs []string
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedUIDs(m map[uint32]struct{}) []uint32 {
	out := make([]uint32, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func opsOf(m map[ebpf.EventType]struct{}) []ebpf.EventType {
	out := make([]ebpf.EventType, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func opNames(m map[ebpf.EventType]struct{}) []string {
	ops := opsOf(m)
	out := make([]string, len(ops))
	for i, o := range ops {
		out[i] = o.String()
	}
	return out
}

// statKey is the (dev, ino) of f, dev encoded as the kernel keys it.
func statKey(f *os.File) (exeKey, error) {
	dev, ino, err := ebpf.StatFile(f)
	return exeKey{dev, ino}, err
}
