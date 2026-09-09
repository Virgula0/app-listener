package daemon

import (
	"errors"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/repository"
	"github.com/Virgula0/app-listener/internal/usecase"
)

// fakeDaemon is a minimal usecase.DaemonUseCase for the startup-abort tests:
// only Stop is exercised.
type fakeDaemon struct{ stopped atomic.Bool }

func (f *fakeDaemon) Start() error { return nil }
func (f *fakeDaemon) Reload([]daemonconfig.Resource, []repository.GuardRepository) error {
	return nil
}
func (f *fakeDaemon) Stop()                              { f.stopped.Store(true) }
func (f *fakeDaemon) Events() <-chan usecase.DaemonEvent { return nil }
func (f *fakeDaemon) Resources() []daemonconfig.Resource { return nil }
func (f *fakeDaemon) GrantEditAccess(string) (func() error, error) {
	return func() error { return nil }, nil
}

// TestAwaitStartupOrSignalCompletes: no signal — the started daemon is passed
// through untouched.
func TestAwaitStartupOrSignalCompletes(t *testing.T) {
	fd := &fakeDaemon{}
	term := make(chan os.Signal, 1)

	d, err := awaitStartupOrSignal(term, func() (usecase.DaemonUseCase, error) {
		return fd, nil
	})
	if err != nil {
		t.Fatalf("awaitStartupOrSignal: %v", err)
	}
	if d == nil {
		t.Fatal("want the started daemon, got nil")
	}
	if fd.stopped.Load() {
		t.Error("Stop must not be called when startup completes cleanly")
	}
}

// TestAwaitStartupOrSignalAbortsAndLocksDown: a termination signal mid-startup
// waits for startup to finish, then runs the secure lockdown (Stop) and
// returns (nil, nil) so the caller does NOT proceed to the run loop with
// vaults unlocked.
func TestAwaitStartupOrSignalAbortsAndLocksDown(t *testing.T) {
	fd := &fakeDaemon{}
	term := make(chan os.Signal, 1)
	started := make(chan struct{})

	go func() {
		<-started
		term <- syscall.SIGTERM
	}()

	d, err := awaitStartupOrSignal(term, func() (usecase.DaemonUseCase, error) {
		close(started)
		time.Sleep(20 * time.Millisecond) // startup still in flight when the signal lands
		return fd, nil
	})
	if err != nil {
		t.Fatalf("awaitStartupOrSignal: %v", err)
	}
	if d != nil {
		t.Fatal("an aborted startup must return a nil daemon so the caller stops")
	}
	if !fd.stopped.Load() {
		t.Error("the secure lockdown (Stop) must run when startup is aborted by a signal")
	}
}

// TestAwaitStartupOrSignalAbortWithFailedStartup: startup failed on its own
// (already cleaned up) and a signal also arrived — no nil-pointer Stop, still
// returns the abort sentinel.
func TestAwaitStartupOrSignalAbortWithFailedStartup(t *testing.T) {
	term := make(chan os.Signal, 1)
	term <- syscall.SIGTERM

	d, err := awaitStartupOrSignal(term, func() (usecase.DaemonUseCase, error) {
		return nil, errors.New("startup failed and rolled back")
	})
	if err != nil {
		t.Fatalf("awaitStartupOrSignal: %v", err)
	}
	if d != nil {
		t.Fatal("want nil daemon after an aborted/failed startup")
	}
}
