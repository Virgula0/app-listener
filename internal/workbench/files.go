package workbench

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// GuardMode is a target's guard state.
type GuardMode int

const (
	GuardOff GuardMode = iota
	GuardWhitelist
	GuardBlacklist
)

func (m GuardMode) String() string {
	switch m {
	case GuardWhitelist:
		return "whitelist"
	case GuardBlacklist:
		return "blacklist"
	default:
		return "off"
	}
}

func (m GuardMode) engineMode() guard.Mode {
	if m == GuardBlacklist {
		return guard.ModeBlacklist
	}
	return guard.ModeWhitelist
}

// Policy is a target's guard configuration.
type Policy struct {
	Mode     GuardMode
	Binaries []ebpf.BinaryEntry
}

// FileSpec is a watched file or directory. Recursive/Depth apply to the monitor and the guard alike.
type FileSpec struct {
	Path      string
	IsDir     bool
	Recursive bool
	Depth     int
	Events    []ebpf.EventType // monitor display filter; empty = all
}

// FileGuard is the slice of the guard engine a target drives.
type FileGuard interface {
	Start() error
	Stop()
	Events() <-chan guard.GuardEvent
	SetBinaries([]ebpf.BinaryEntry) error
}

// FileEngines builds the engines; tests swap in fakes.
type FileEngines struct {
	NewMonitor func(t ebpf.Target, recursive bool, depth int) (repository.MonitorRepository, error)
	NewGuard   func(path string, mode guard.Mode, bins []ebpf.BinaryEntry, recursive bool, depth int) (FileGuard, error)
	Hash       func(path string) (ebpf.BinaryEntry, error)
}

// Files is the Files workspace: several targets side by side, each with its own monitor and guard.
type Files struct {
	engines FileEngines
	exes    *exeResolver
	ownPID  uint32

	mu      sync.Mutex
	targets []*FileTarget
	nextID  int
}

// NewFiles returns an empty workspace.
func NewFiles(engines FileEngines) *Files {
	return &Files{engines: engines, exes: newExeResolver("/proc"), ownPID: uint32(os.Getpid())} //nolint:gosec // pid fits
}

// Targets returns the targets in insertion order.
func (w *Files) Targets() []*FileTarget {
	w.mu.Lock()
	defer w.mu.Unlock()
	return slices.Clone(w.targets)
}

// NormalizeSpec resolves spec.Path to an absolute existing path and fills IsDir.
func NormalizeSpec(spec FileSpec) (FileSpec, error) {
	if strings.TrimSpace(spec.Path) == "" {
		return spec, errors.New("choose a file or folder to watch")
	}
	abs, err := filepath.Abs(spec.Path)
	if err != nil {
		return spec, err
	}
	fi, err := os.Stat(abs)
	if err != nil {
		return spec, fmt.Errorf("cannot watch %s: %w", abs, err)
	}
	if spec.Depth < 0 {
		return spec, errors.New("depth must be 0 (unlimited) or more")
	}
	spec.Path, spec.IsDir = abs, fi.IsDir()
	if !spec.IsDir {
		spec.Recursive, spec.Depth = false, 0
	}
	return spec, nil
}

// overlaps reports whether a and b can't both be targets: the same path, or one inside the other
// (two guards would then judge the same files).
func overlaps(a, b FileSpec) bool {
	inside := func(p, dir string) bool { return strings.HasPrefix(p+"/", dir+"/") }
	switch {
	case a.Path == b.Path:
		return true
	case a.IsDir && inside(b.Path, a.Path):
		return true
	case b.IsDir && inside(a.Path, b.Path):
		return true
	}
	return false
}

// Add validates spec and starts monitoring it; the guard starts off.
func (w *Files) Add(spec FileSpec) (*FileTarget, error) {
	spec, err := NormalizeSpec(spec)
	if err != nil {
		return nil, err
	}
	w.mu.Lock()
	for _, t := range w.targets {
		if overlaps(t.Spec, spec) {
			w.mu.Unlock()
			return nil, fmt.Errorf("%s overlaps the target %s: remove one of them first", spec.Path, t.Spec.Path)
		}
	}
	w.nextID++
	t := newFileTarget(w, w.nextID, spec)
	w.targets = append(w.targets, t)
	w.mu.Unlock()

	if err := t.startMonitor(); err != nil {
		w.Remove(t)
		return nil, err
	}
	return t, nil
}

// Remove stops t's engines and drops it.
func (w *Files) Remove(t *FileTarget) {
	t.close()
	w.mu.Lock()
	w.targets = slices.DeleteFunc(w.targets, func(x *FileTarget) bool { return x == t })
	w.mu.Unlock()
}

// Close stops every target (window closed).
func (w *Files) Close() {
	for _, t := range w.Targets() {
		w.Remove(t)
	}
}

// FileTarget is one watched path with its monitor, optional guard and event buffer.
type FileTarget struct {
	ID   int
	Spec FileSpec

	w *Files
	*activity

	mu       sync.Mutex
	mon      repository.MonitorRepository
	monStop  chan struct{}
	grd      FileGuard
	grdStop  chan struct{}
	policy   Policy
	closed   bool
	lastErr  string
	policyMu sync.Mutex
}

func newFileTarget(w *Files, id int, spec FileSpec) *FileTarget {
	return &FileTarget{ID: id, Spec: spec, w: w, activity: newActivity(w.ownPID)}
}

func (t *FileTarget) target() ebpf.Target {
	tg := ebpf.Target{Path: t.Spec.Path, IsDir: t.Spec.IsDir}
	if t.Spec.IsDir {
		tg.Dir = t.Spec.Path
	} else {
		tg.Dir, tg.File = filepath.Dir(t.Spec.Path), filepath.Base(t.Spec.Path)
	}
	return tg
}

func (t *FileTarget) startMonitor() error {
	mon, err := t.w.engines.NewMonitor(t.target(), t.Spec.Recursive, t.Spec.Depth)
	if err != nil {
		return fmt.Errorf("creating monitor: %w", err)
	}
	if len(t.Spec.Events) > 0 {
		mon.SetEventTypes(t.Spec.Events)
	}
	if err := mon.Start(); err != nil {
		return fmt.Errorf("starting monitor: %w", err)
	}
	stop := make(chan struct{})
	t.mu.Lock()
	t.mon, t.monStop = mon, stop
	t.mu.Unlock()
	go t.pumpMonitor(mon.Events(), stop)
	return nil
}

func (t *FileTarget) pumpMonitor(events <-chan ebpf.FileEvent, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			t.ingest(&Row{
				Time: time.Unix(0, ev.Timestamp), Source: "monitor", Op: ev.Type.String(), Comm: ev.Comm,
				PID: ev.PID, UID: ev.UID, Path: ev.Path, Dest: ev.Dest, Verdict: Observed,
				Exe: t.w.exes.resolve(ev.PID),
			})
		}
	}
}

func (t *FileTarget) pumpGuard(events <-chan guard.GuardEvent, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			row := guardRow(&ev, t.w.exes)
			t.ingest(&row)
		}
	}
}

// guardRow converts a guard event; the exe path is taken from /proc only when it names the inode
// the kernel judged.
func guardRow(ev *guard.GuardEvent, exes *exeResolver) Row {
	v := Allowed
	if ev.Blocked {
		v = Denied
	}
	detail := ev.FsGate
	if ev.Scope != "" {
		detail = strings.TrimSpace(detail + " scope=" + ev.Scope)
	}
	r := Row{
		Time: time.Unix(0, ev.Timestamp), Source: "guard", Op: ev.Op(), Comm: ev.Comm, PID: ev.PID,
		UID: ev.UID, Path: ev.Path, Dest: ev.Dest, Verdict: v, Detail: detail,
	}
	if ev.Exe != (guard.GuardInodeKey{}) {
		r.Exe.Key = ev.Exe.Real()
		if seen := exes.resolve(ev.PID); seen.Key == r.Exe.Key {
			r.Exe.Path = seen.Path
		}
	}
	return r
}

// Policy returns the policy the guard enforces now.
func (t *FileTarget) Policy() Policy {
	t.mu.Lock()
	defer t.mu.Unlock()
	return Policy{Mode: t.policy.Mode, Binaries: slices.Clone(t.policy.Binaries)}
}

// LastError is the most recent engine failure, for the view.
func (t *FileTarget) LastError() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastErr
}

// ApplyLeavesGap reports whether applying p unguards the path for a moment: a mode change away from
// a whitelist must stop it before the new guard's tree walk, which the whitelist would deny.
func (t *FileTarget) ApplyLeavesGap(p Policy) bool {
	cur := t.Policy().Mode
	return cur == GuardWhitelist && p.Mode != GuardOff && p.Mode != cur
}

// Apply makes p the enforced policy. Same mode: rows change in place on the running guard
// (Guard.SetBinaries). Turning on: a fresh guard. Mode change: from a blacklist the new guard is
// up before the old one goes (no gap); from a whitelist the old one must go first (ApplyLeavesGap).
// On error the previous guard keeps running, except after a gap step, where the target is left
// unguarded and LastError says so. A whitelist binary changed since it was added is refused
// (verifyHashes); a blacklist blocks whatever the path names now, so it needs no check.
func (t *FileTarget) Apply(p Policy) error {
	t.policyMu.Lock()
	defer t.policyMu.Unlock()
	p.Binaries = slices.Clone(p.Binaries)

	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return errors.New("target removed")
	}
	cur, old := t.policy, t.grd
	t.mu.Unlock()

	var err error
	if p.Mode == GuardWhitelist {
		err = verifyHashes(t.w.engines.Hash, p.Binaries)
	}
	switch {
	case err != nil:
	case p.Mode == GuardOff:
		t.swapGuard(nil, p)
	case cur.Mode == p.Mode:
		if err = old.SetBinaries(p.Binaries); err == nil {
			t.setPolicy(p)
		}
	case cur.Mode == GuardWhitelist:
		t.swapGuard(nil, Policy{Mode: GuardOff})
		err = t.startGuard(p)
	default: // off or blacklist: the new guard comes up first
		err = t.startGuard(p)
	}
	t.mu.Lock()
	t.lastErr = ""
	if err != nil {
		t.lastErr = err.Error()
	}
	t.mu.Unlock()
	return err
}

func (t *FileTarget) startGuard(p Policy) error {
	g, err := t.w.engines.NewGuard(t.Spec.Path, p.Mode.engineMode(), p.Binaries, t.Spec.Recursive || !t.Spec.IsDir, t.Spec.Depth)
	if err != nil {
		return fmt.Errorf("creating guard: %w", err)
	}
	if err := g.Start(); err != nil {
		g.Stop()
		return fmt.Errorf("starting guard: %w", err)
	}
	t.swapGuard(g, p)
	return nil
}

// swapGuard installs next (nil = none) and only then stops the previous guard.
func (t *FileTarget) swapGuard(next FileGuard, p Policy) {
	var stop chan struct{}
	if next != nil {
		stop = make(chan struct{})
	}
	t.mu.Lock()
	prev, prevStop := t.grd, t.grdStop
	t.grd, t.grdStop, t.policy = next, stop, p
	t.mu.Unlock()
	if next != nil {
		go t.pumpGuard(next.Events(), stop)
	}
	stopPumped(prev, prevStop)
}

func (t *FileTarget) setPolicy(p Policy) {
	t.mu.Lock()
	t.policy = p
	t.mu.Unlock()
}

// stopPumped closes a pumped engine's stop channel, then stops the engine; nil (none) is a no-op.
func stopPumped(e interface{ Stop() }, stop chan struct{}) {
	if e == nil {
		return
	}
	close(stop)
	e.Stop()
}

func (t *FileTarget) close() {
	t.policyMu.Lock()
	defer t.policyMu.Unlock()
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
	}
	t.closed = true
	mon, monStop, grd, grdStop := t.mon, t.monStop, t.grd, t.grdStop
	t.mon, t.grd = nil, nil
	t.mu.Unlock()
	stopPumped(grd, grdStop)
	stopPumped(mon, monStop)
}
