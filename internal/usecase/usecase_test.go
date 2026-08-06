package usecase

import (
	"errors"
	"sync"
	"testing"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
)

var errBoom = errors.New("boom")

type fakeMonitorRepo struct {
	started  bool
	stopped  bool
	startErr error
	events   chan ebpf.FileEvent
	eventTyp []ebpf.EventType
}

func newFakeMonitorRepo() *fakeMonitorRepo {
	return &fakeMonitorRepo{events: make(chan ebpf.FileEvent)}
}

func (f *fakeMonitorRepo) Start() error {
	f.started = true
	return f.startErr
}

func (f *fakeMonitorRepo) Stop() {
	f.stopped = true
	close(f.events)
}

func (f *fakeMonitorRepo) Events() <-chan ebpf.FileEvent {
	return f.events
}

func (f *fakeMonitorRepo) SetEventTypes(types []ebpf.EventType) {
	f.eventTyp = types
}

type fakeGuardRepo struct {
	mu          sync.Mutex
	started     bool
	stopped     bool
	populated   bool
	startErr    error
	populateErr error
	events      chan guard.GuardEvent
}

func newFakeGuardRepo() *fakeGuardRepo {
	return &fakeGuardRepo{events: make(chan guard.GuardEvent)}
}

func (f *fakeGuardRepo) PopulateInodes() error {
	if f.populateErr != nil {
		return f.populateErr
	}
	f.populated = true
	return nil
}

func (f *fakeGuardRepo) Start() error {
	f.started = true
	return f.startErr
}

func (f *fakeGuardRepo) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stopped {
		return
	}
	f.stopped = true
	close(f.events)
}

// isStopped is the mutex-safe stopped-state accessor for tests that observe
// Stop concurrently.
func (f *fakeGuardRepo) isStopped() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.stopped
}

func (f *fakeGuardRepo) Events() <-chan guard.GuardEvent {
	return f.events
}

type fakeNetworkMonitorRepo struct {
	started  bool
	stopped  bool
	startErr error
	events   chan ebpf.NetEvent
	eventTyp []ebpf.NetEventType
}

func newFakeNetworkMonitorRepo() *fakeNetworkMonitorRepo {
	return &fakeNetworkMonitorRepo{events: make(chan ebpf.NetEvent)}
}

func (f *fakeNetworkMonitorRepo) Start() error {
	f.started = true
	return f.startErr
}

func (f *fakeNetworkMonitorRepo) Stop() {
	f.stopped = true
	close(f.events)
}

func (f *fakeNetworkMonitorRepo) Events() <-chan ebpf.NetEvent {
	return f.events
}

func (f *fakeNetworkMonitorRepo) SetEventTypes(types []ebpf.NetEventType) {
	f.eventTyp = types
}

type fakeNetworkGuardRepo struct {
	started  bool
	stopped  bool
	startErr error
	events   chan networkguard.NetGuardEvent
}

func newFakeNetworkGuardRepo() *fakeNetworkGuardRepo {
	return &fakeNetworkGuardRepo{events: make(chan networkguard.NetGuardEvent)}
}

func (f *fakeNetworkGuardRepo) Start() error {
	f.started = true
	return f.startErr
}

func (f *fakeNetworkGuardRepo) Stop() {
	f.stopped = true
	close(f.events)
}

func (f *fakeNetworkGuardRepo) Events() <-chan networkguard.NetGuardEvent {
	return f.events
}

func TestMonitorUseCaseLifecycle(t *testing.T) {
	repo := newFakeMonitorRepo()
	uc := NewMonitorUseCase(repo)

	if err := uc.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !repo.started {
		t.Fatal("repo not started")
	}

	types := []ebpf.EventType{ebpf.EventRead, ebpf.EventWrite}
	repo.SetEventTypes(types)
	if len(repo.eventTyp) != 2 || repo.eventTyp[0] != ebpf.EventRead {
		t.Fatalf("SetEventTypes not forwarded: %v", repo.eventTyp)
	}

	if ch := uc.Events(); ch != (chan ebpf.FileEvent)(repo.events) {
		t.Fatal("Events did not return the repo channel")
	}

	uc.Stop()
	if !repo.stopped {
		t.Fatal("repo not stopped")
	}
}

func TestMonitorUseCaseStartError(t *testing.T) {
	repo := &fakeMonitorRepo{}
	repo.startErr = errBoom
	uc := NewMonitorUseCase(repo)

	if err := uc.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate repo error, got %v", err)
	}
}

func TestGuardUseCaseLifecycle(t *testing.T) {
	repo := newFakeGuardRepo()
	uc := NewGuardUseCase(repo)

	if err := uc.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !repo.started {
		t.Fatal("repo not started")
	}
	if ch := uc.Events(); ch != (chan guard.GuardEvent)(repo.events) {
		t.Fatal("Events did not return the repo channel")
	}
	uc.Stop()
	if !repo.stopped {
		t.Fatal("repo not stopped")
	}
}

func TestGuardUseCaseStartError(t *testing.T) {
	repo := &fakeGuardRepo{}
	repo.startErr = errBoom
	uc := NewGuardUseCase(repo)

	if err := uc.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate repo error, got %v", err)
	}
}

func TestNetworkMonitorUseCaseLifecycle(t *testing.T) {
	repo := newFakeNetworkMonitorRepo()
	uc := NewNetworkMonitorUseCase(repo)

	if err := uc.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !repo.started {
		t.Fatal("repo not started")
	}

	types := []ebpf.NetEventType{ebpf.NetConnect}
	uc.SetEventTypes(types)
	if len(repo.eventTyp) != 1 || repo.eventTyp[0] != ebpf.NetConnect {
		t.Fatalf("SetEventTypes not forwarded: %v", repo.eventTyp)
	}

	if ch := uc.Events(); ch != (chan ebpf.NetEvent)(repo.events) {
		t.Fatal("Events did not return the repo channel")
	}
	uc.Stop()
	if !repo.stopped {
		t.Fatal("repo not stopped")
	}
}

func TestNetworkGuardUseCaseLifecycle(t *testing.T) {
	repo := newFakeNetworkGuardRepo()
	uc := NewNetworkGuardUseCase(repo)

	if err := uc.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !repo.started {
		t.Fatal("repo not started")
	}
	if ch := uc.Events(); ch != (chan networkguard.NetGuardEvent)(repo.events) {
		t.Fatal("Events did not return the repo channel")
	}
	uc.Stop()
	if !repo.stopped {
		t.Fatal("repo not stopped")
	}
}

func TestNetworkGuardUseCaseStartError(t *testing.T) {
	repo := &fakeNetworkGuardRepo{}
	repo.startErr = errBoom
	uc := NewNetworkGuardUseCase(repo)

	if err := uc.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate repo error, got %v", err)
	}
}

var _ repository.MonitorRepository = (*fakeMonitorRepo)(nil)
var _ repository.GuardRepository = (*fakeGuardRepo)(nil)
var _ repository.NetworkMonitorRepository = (*fakeNetworkMonitorRepo)(nil)
var _ repository.NetworkGuardRepository = (*fakeNetworkGuardRepo)(nil)
