package workbench

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// journal records engine calls in order across fakes.
type journal struct {
	mu  sync.Mutex
	log []string
}

func (j *journal) add(s string) {
	j.mu.Lock()
	j.log = append(j.log, s)
	j.mu.Unlock()
}

func (j *journal) entries() []string {
	j.mu.Lock()
	defer j.mu.Unlock()
	return slices.Clone(j.log)
}

type fakeMonitor struct {
	j      *journal
	events chan ebpf.FileEvent
}

func (m *fakeMonitor) Start() error                     { m.j.add("monitor start"); return nil }
func (m *fakeMonitor) Stop()                            { m.j.add("monitor stop") }
func (m *fakeMonitor) Events() <-chan ebpf.FileEvent    { return m.events }
func (m *fakeMonitor) SetEventTypes(_ []ebpf.EventType) {}

type fakeGuard struct {
	j      *journal
	name   string
	events chan guard.GuardEvent
	bins   []ebpf.BinaryEntry
	setErr error
}

func (g *fakeGuard) Start() error                    { g.j.add(g.name + " start"); return nil }
func (g *fakeGuard) Stop()                           { g.j.add(g.name + " stop") }
func (g *fakeGuard) Events() <-chan guard.GuardEvent { return g.events }
func (g *fakeGuard) SetBinaries(b []ebpf.BinaryEntry) error {
	if g.setErr != nil {
		return g.setErr
	}
	g.j.add(fmt.Sprintf("%s set %d", g.name, len(b)))
	g.bins = b
	return nil
}

type fakeEngines struct {
	j        *journal
	n        int
	guards   []*fakeGuard
	monitor  *fakeMonitor
	guardErr error
}

func (f *fakeEngines) engines() FileEngines {
	return FileEngines{
		NewMonitor: func(ebpf.Target, bool, int) (repository.MonitorRepository, error) {
			f.monitor = &fakeMonitor{j: f.j, events: make(chan ebpf.FileEvent, 16)}
			return f.monitor, nil
		},
		NewGuard: func(_ string, mode guard.Mode, bins []ebpf.BinaryEntry, _ bool, _ int) (FileGuard, error) {
			if f.guardErr != nil {
				f.j.add("guard create failed")
				return nil, f.guardErr
			}
			f.n++
			g := &fakeGuard{j: f.j, name: fmt.Sprintf("guard%d(%d)", f.n, mode), events: make(chan guard.GuardEvent), bins: bins}
			f.j.add(g.name + " create")
			f.guards = append(f.guards, g)
			return g, nil
		},
	}
}

func newTestTarget(t *testing.T) (*FileTarget, *fakeEngines) {
	t.Helper()
	fe := &fakeEngines{j: &journal{}}
	w := NewFiles(fe.engines())
	t.Cleanup(w.Close)
	tg, err := w.Add(FileSpec{Path: t.TempDir(), Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	return tg, fe
}

var binA = ebpf.BinaryEntry{Path: "/a", Hash: [32]byte{1}}

func since(fe *fakeEngines, from int) []string { return fe.j.entries()[from:] }

func TestApplyOffToWhitelistStartsGuard(t *testing.T) {
	tg, fe := newTestTarget(t)
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{binA}}); err != nil {
		t.Fatal(err)
	}
	want := []string{"guard1(1) create", "guard1(1) start"}
	if got := since(fe, mark); !slices.Equal(got, want) {
		t.Fatalf("calls = %q, want %q", got, want)
	}
	if tg.Policy().Mode != GuardWhitelist {
		t.Fatal("policy not recorded")
	}
}

func TestApplySameModeEditsInPlace(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{binA}}); err != nil {
		t.Fatal(err)
	}
	if got := since(fe, mark); !slices.Equal(got, []string{"guard1(1) set 1"}) {
		t.Fatalf("a same-mode edit must change rows in place, got %q", got)
	}
}

func TestApplySameModeFailureKeepsPolicy(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	fe.guards[0].setErr = errors.New("boom")
	if err := tg.Apply(Policy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{binA}}); err == nil {
		t.Fatal("want the engine error")
	}
	if p := tg.Policy(); len(p.Binaries) != 0 || p.Mode != GuardWhitelist {
		t.Fatalf("a failed edit must leave the recorded policy as enforced, got %+v", p)
	}
	if tg.LastError() == "" {
		t.Fatal("LastError must report the failure")
	}
}

// From a blacklist the replacement attaches before the old guard goes: no unguarded moment.
func TestApplyFromBlacklistOverlaps(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardBlacklist}); err != nil {
		t.Fatal(err)
	}
	if tg.ApplyLeavesGap(Policy{Mode: GuardWhitelist}) {
		t.Fatal("blacklist -> whitelist must not report a gap")
	}
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	want := []string{"guard2(1) create", "guard2(1) start", "guard1(0) stop"}
	if got := since(fe, mark); !slices.Equal(got, want) {
		t.Fatalf("calls = %q, want %q", got, want)
	}
}

func TestApplyFromBlacklistFailureKeepsOldGuard(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardBlacklist, Binaries: []ebpf.BinaryEntry{binA}}); err != nil {
		t.Fatal(err)
	}
	fe.guardErr = errors.New("verifier said no")
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err == nil {
		t.Fatal("want the engine error")
	}
	if got := since(fe, mark); slices.Contains(got, "guard1(0) stop") {
		t.Fatalf("the running blacklist must survive a failed replacement, got %q", got)
	}
	if p := tg.Policy(); p.Mode != GuardBlacklist || len(p.Binaries) != 1 {
		t.Fatalf("policy must still be the running blacklist, got %+v", p)
	}
}

// From a whitelist the old guard must go first (its hooks deny the new guard's walk): a gap,
// reported up front.
func TestApplyFromWhitelistStopsFirst(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	if !tg.ApplyLeavesGap(Policy{Mode: GuardBlacklist}) {
		t.Fatal("whitelist -> blacklist must report the gap")
	}
	if tg.ApplyLeavesGap(Policy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{binA}}) {
		t.Fatal("a same-mode edit has no gap")
	}
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardBlacklist}); err != nil {
		t.Fatal(err)
	}
	want := []string{"guard1(1) stop", "guard2(0) create", "guard2(0) start"}
	if got := since(fe, mark); !slices.Equal(got, want) {
		t.Fatalf("calls = %q, want %q", got, want)
	}
}

func TestApplyFromWhitelistFailureReportsUnguarded(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	fe.guardErr = errors.New("verifier said no")
	if err := tg.Apply(Policy{Mode: GuardBlacklist}); err == nil {
		t.Fatal("want the engine error")
	}
	if tg.Policy().Mode != GuardOff {
		t.Fatal("after the gap step failed the target is unguarded and must say so")
	}
}

func TestApplyOffStopsGuard(t *testing.T) {
	tg, fe := newTestTarget(t)
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err != nil {
		t.Fatal(err)
	}
	mark := len(fe.j.entries())
	if err := tg.Apply(Policy{Mode: GuardOff}); err != nil {
		t.Fatal(err)
	}
	if got := since(fe, mark); !slices.Equal(got, []string{"guard1(1) stop"}) {
		t.Fatalf("calls = %q", got)
	}
}

func TestRemoveStopsGuardThenMonitor(t *testing.T) {
	fe := &fakeEngines{j: &journal{}}
	w := NewFiles(fe.engines())
	tg, err := w.Add(FileSpec{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	if err := tg.Apply(Policy{Mode: GuardBlacklist}); err != nil {
		t.Fatal(err)
	}
	mark := len(fe.j.entries())
	w.Remove(tg)
	if got := since(fe, mark); !slices.Equal(got, []string{"guard1(0) stop", "monitor stop"}) {
		t.Fatalf("calls = %q", got)
	}
	if err := tg.Apply(Policy{Mode: GuardWhitelist}); err == nil {
		t.Fatal("a removed target must refuse a policy")
	}
	if len(w.Targets()) != 0 {
		t.Fatal("target not removed")
	}
}

func TestAddRefusesOverlap(t *testing.T) {
	fe := &fakeEngines{j: &journal{}}
	w := NewFiles(fe.engines())
	defer w.Close()
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	file := filepath.Join(root, "f")
	if err := os.Mkdir(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Add(FileSpec{Path: sub}); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{sub, root} {
		if _, err := w.Add(FileSpec{Path: p}); err == nil {
			t.Fatalf("%s overlaps %s and must be refused", p, sub)
		}
	}
	if _, err := w.Add(FileSpec{Path: file}); err != nil {
		t.Fatalf("a sibling file must be accepted: %v", err)
	}
	if _, err := w.Add(FileSpec{Path: filepath.Join(root, "missing")}); err == nil {
		t.Fatal("a missing path must be refused")
	}
}

func TestIngestAggregatesAndSkipsOwnPID(t *testing.T) {
	tg, fe := newTestTarget(t)
	tg.ingest(&Row{PID: tg.w.ownPID, Comm: "self"})
	key := guard.GuardInodeKey{Dev: 1, Ino: 2}
	tg.ingest(&Row{PID: 10, Comm: "cat", Verdict: Observed, Exe: Exe{Key: key}})
	tg.ingest(&Row{PID: 11, Comm: "cat", Verdict: Denied, Exe: Exe{Path: "/bin/cat", Key: key}})
	tg.ingest(&Row{PID: 12, Comm: "ls", Verdict: Allowed})
	_ = fe

	rows, _ := tg.Rows(Filter{})
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (own pid skipped)", len(rows))
	}
	if rows[0].Comm != "ls" {
		t.Fatal("rows must be newest first")
	}
	denied, _ := tg.Rows(Filter{Verdicts: []Verdict{Denied}})
	if len(denied) != 1 {
		t.Fatalf("denied filter = %d", len(denied))
	}
	if text, _ := tg.Rows(Filter{Text: "/BIN/CAT"}); len(text) != 1 {
		t.Fatalf("text filter must be case-insensitive and match the exe path, got %d", len(text))
	}
	procs := tg.Processes()
	if len(procs) != 2 {
		t.Fatalf("processes = %d, want 2 (one per inode or comm)", len(procs))
	}
	for _, p := range procs {
		if p.Comm == "cat" && (p.Events != 2 || p.Denied != 1 || p.Exe.Path != "/bin/cat") {
			t.Fatalf("cat aggregate = %+v", p)
		}
	}
	tg.Clear()
	if rows, _ := tg.Rows(Filter{}); len(rows) != 0 || len(tg.Processes()) != 0 {
		t.Fatal("Clear must empty rows and processes")
	}
}

func TestRingKeepsNewest(t *testing.T) {
	r := newRing()
	for i := range ringCap + 5 {
		r.add(&Row{PID: uint32(i)}) //nolint:gosec // test
	}
	rows, _ := r.snapshot(Filter{})
	if len(rows) != ringCap || rows[0].PID != ringCap+4 || rows[len(rows)-1].PID != 5 {
		t.Fatalf("len=%d first=%d last=%d", len(rows), rows[0].PID, rows[len(rows)-1].PID)
	}
	c, _ := r.totals()
	if c.Observed != ringCap+5 {
		t.Fatalf("totals count every row ever added, got %d", c.Observed)
	}
}

func TestEntryForRefusesReplacedPath(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tool")
	if err := os.WriteFile(p, []byte("v1"), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, err := ebpf.StatInode(p)
	if err != nil {
		t.Fatal(err)
	}
	seen := Exe{Path: p, Key: guard.GuardInodeKey{Dev: dev, Ino: ino}}
	entry, err := EntryFor(seen)
	if err != nil {
		t.Fatal(err)
	}
	if entry.Path != p {
		t.Fatalf("entry path = %s", entry.Path)
	}

	// The link keeps the old inode allocated, so the replacement can't reuse its number.
	if err := os.Link(p, p+".keep"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(p); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("v2"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := EntryFor(seen); err == nil {
		t.Fatal("a path that now names another inode must be refused")
	}
	if _, err := EntryFor(Exe{Path: p}); !errors.Is(err, ErrExeUnknown) {
		t.Fatalf("an unresolved exe must be ErrExeUnknown, got %v", err)
	}
}
