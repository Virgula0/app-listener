package usecase

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

type tempTrace struct {
	mu    sync.Mutex
	steps []string
}

func (t *tempTrace) add(s string) {
	t.mu.Lock()
	t.steps = append(t.steps, s)
	t.mu.Unlock()
}

func (t *tempTrace) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.Join(t.steps, " ")
}

type fakeTempGrant struct {
	name     string
	trace    *tempTrace
	rows     []guard.GuardResInodeKey
	report   []string
	applyErr error
}

func (f *fakeTempGrant) JournalRows() []guard.GuardResInodeKey { return f.rows }
func (f *fakeTempGrant) Report() []string                      { return f.report }
func (f *fakeTempGrant) Apply() error {
	f.trace.add("apply:" + f.name)
	return f.applyErr
}
func (f *fakeTempGrant) Revoke() error {
	f.trace.add("revoke:" + f.name)
	return nil
}

func tempBinary(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "tool")
	if err := os.WriteFile(p, []byte(content), 0o700); err != nil { //nolint:gosec // test fixture must be executable
		t.Fatal(err)
	}
	return p
}

func tempGrantSetup(t *testing.T, grants ...*fakeTempGrant) (DaemonUseCase, []string) {
	t.Helper()
	var resources []daemonconfig.Resource
	var guards []repository.GuardRepository
	var paths []string
	for i, gr := range grants {
		p := filepath.Join("/res", string(rune('a'+i)))
		paths = append(paths, p)
		resources = append(resources, resource(p))
		repo := newFakeGuardRepo()
		repo.tempPlan = func([]*guard.TempBinary, guard.TempRule) (guard.TemporaryGrant, error) { return gr, nil }
		guards = append(guards, repo)
	}
	d, err := NewDaemonUseCase(resources, newFakeVault(), guards)
	if err != nil {
		t.Fatal(err)
	}
	return d, paths
}

func journalTo(trace *tempTrace) func([]guard.GuardResInodeKey) error {
	return func(rows []guard.GuardResInodeKey) error {
		if rows == nil {
			trace.add("journal:clear")
		} else {
			trace.add("journal:write")
		}
		return nil
	}
}

func TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(t *testing.T) {
	trace := &tempTrace{}
	row := []guard.GuardResInodeKey{{ResId: 1}}
	d, res := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace, rows: row},
		&fakeTempGrant{name: "b", trace: trace, rows: row})

	acc, err := d.GrantTemporaryAccess(res, []string{tempBinary(t, "x")}, guard.TempRule{Allow: true}, journalTo(trace))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.GrantTemporaryAccess(res, []string{tempBinary(t, "y")}, guard.TempRule{Allow: true},
		journalTo(trace)); err == nil {
		t.Fatal("a second concurrent grant was accepted")
	}
	if err := acc.Revoke(); err != nil {
		t.Fatal(err)
	}
	<-acc.Done
	want := "journal:write apply:a apply:b revoke:b revoke:a journal:clear"
	if got := trace.String(); got != want {
		t.Fatalf("steps = %q, want %q", got, want)
	}
	again, err := d.GrantTemporaryAccess(res, []string{tempBinary(t, "z")}, guard.TempRule{Allow: true}, journalTo(trace))
	if err != nil {
		t.Fatalf("revoke did not free the session slot: %v", err)
	}
	_ = again.Revoke()
}

func TestGrantTemporaryAccessRollsBackOnApplyFailure(t *testing.T) {
	trace := &tempTrace{}
	row := []guard.GuardResInodeKey{{ResId: 1}}
	d, res := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace, rows: row},
		&fakeTempGrant{name: "b", trace: trace, applyErr: errBoom})

	if _, err := d.GrantTemporaryAccess(res, []string{tempBinary(t, "x")}, guard.TempRule{Allow: true},
		journalTo(trace)); !errors.Is(err, errBoom) {
		t.Fatalf("err = %v, want boom", err)
	}
	want := "journal:write apply:a apply:b revoke:a journal:clear"
	if got := trace.String(); got != want {
		t.Fatalf("steps = %q, want %q", got, want)
	}
}

func TestGrantTemporaryAccessRefusals(t *testing.T) {
	trace := &tempTrace{}
	bin := tempBinary(t, "x")
	t.Run("unknown resource", func(t *testing.T) {
		d, _ := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace})
		if _, err := d.GrantTemporaryAccess([]string{"/elsewhere"}, []string{bin}, guard.TempRule{Allow: true},
			journalTo(trace)); err == nil {
			t.Fatal("unknown resource accepted")
		}
	})
	t.Run("nothing to change", func(t *testing.T) {
		d, res := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace, report: []string{"already whitelisted"}})
		if _, err := d.GrantTemporaryAccess(res, []string{bin}, guard.TempRule{Allow: true},
			journalTo(trace)); err == nil || !strings.Contains(err.Error(), "already whitelisted") {
			t.Fatalf("err = %v, want a nothing-to-change refusal", err)
		}
	})
	t.Run("not executable", func(t *testing.T) {
		d, res := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace})
		p := filepath.Join(t.TempDir(), "data")
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := d.GrantTemporaryAccess(res, []string{p}, guard.TempRule{Allow: true}, journalTo(trace)); err == nil {
			t.Fatal("non-executable binary accepted")
		}
	})
	if got := trace.String(); got != "" {
		t.Fatalf("a refused grant touched state: %q", got)
	}
}

func TestGrantTemporaryAccessRevokesOnInPlaceRewrite(t *testing.T) {
	old := tempBinaryCheckEvery
	tempBinaryCheckEvery = 10 * time.Millisecond
	t.Cleanup(func() { tempBinaryCheckEvery = old })

	trace := &tempTrace{}
	d, res := tempGrantSetup(t, &fakeTempGrant{name: "a", trace: trace, rows: []guard.GuardResInodeKey{{ResId: 1}}})
	bin := tempBinary(t, "original")
	acc, err := d.GrantTemporaryAccess(res, []string{bin}, guard.TempRule{Allow: true}, journalTo(trace))
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(bin, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("rewritten"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	select {
	case <-acc.Done:
	case <-time.After(5 * time.Second):
		t.Fatal("an in-place rewrite of a granted binary did not revoke the grant")
	}
	if got := trace.String(); !strings.HasSuffix(got, "revoke:a journal:clear") {
		t.Fatalf("steps = %q", got)
	}
}

func TestGrantTemporaryAccessStreamsOnlyGrantedExeEvents(t *testing.T) {
	trace := &tempTrace{}
	gr := &fakeTempGrant{name: "a", trace: trace, rows: []guard.GuardResInodeKey{{ResId: 1}}}
	repo := newFakeGuardRepo()
	repo.tempPlan = func([]*guard.TempBinary, guard.TempRule) (guard.TemporaryGrant, error) { return gr, nil }
	d, err := NewDaemonUseCase([]daemonconfig.Resource{resource("/res/a")}, newFakeVault(),
		[]repository.GuardRepository{repo})
	if err != nil {
		t.Fatal(err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	// The granted binary is this test binary, so this process's events match by exe inode.
	acc, err := d.GrantTemporaryAccess([]string{"/res/a"}, []string{self}, guard.TempRule{Allow: true}, journalTo(trace))
	if err != nil {
		t.Fatal(err)
	}
	repo.mu.Lock()
	tap := repo.tap
	repo.mu.Unlock()
	if tap == nil {
		t.Fatal("no event tap installed on the granted resource's guard")
	}

	other := exec.Command("sleep", "5")
	if err := other.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = other.Process.Kill(); _ = other.Wait() }()
	tap(guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: uint32(other.Process.Pid), Path: "/res/a/x"}}) //nolint:gosec // test pid
	tap(guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: uint32(os.Getpid()), Path: "/res/a/mine"}})    //nolint:gosec // test pid

	select {
	case ev := <-acc.Events:
		if ev.Path != "/res/a/mine" {
			t.Fatalf("streamed %q, want only the granted exe's event", ev.Path)
		}
	case <-time.After(time.Second):
		t.Fatal("the granted exe's event was not streamed")
	}
	select {
	case ev := <-acc.Events:
		t.Fatalf("an unrelated process's event was streamed: %q", ev.Path)
	default:
	}

	if err := acc.Revoke(); err != nil {
		t.Fatal(err)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.tap != nil {
		t.Fatal("revoke left the event tap installed")
	}
}
