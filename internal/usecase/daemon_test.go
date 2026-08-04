package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/repository"
)

type fakeVault struct {
	encrypted map[string]bool
	// unlocked is the fscrypt lock state: Unlock provisions, Lock
	// deprovisions.
	unlocked map[string]bool
	// unlockCalls records every Unlock invocation, for assertions.
	unlockCalls []string

	// lockErrs maps path -> queue of errors returned by successive Lock
	// calls; the queue is consumed from the front.
	lockErrs map[string][]error

	unlockErr error
	checkErr  error
}

func newFakeVault(paths ...string) *fakeVault {
	v := &fakeVault{
		encrypted: make(map[string]bool),
		unlocked:  make(map[string]bool),
		lockErrs:  make(map[string][]error),
	}
	for _, p := range paths {
		v.encrypted[p] = true
	}
	return v
}

func (f *fakeVault) IsEncrypted(path string) (bool, error) {
	if f.checkErr != nil {
		return false, f.checkErr
	}
	return f.encrypted[path], nil
}

func (f *fakeVault) IsProvisioned(path string) (bool, error) {
	if f.checkErr != nil {
		return false, f.checkErr
	}
	return f.unlocked[path], nil
}

func (f *fakeVault) Unlock(path string) error {
	if f.unlockErr != nil {
		return f.unlockErr
	}
	f.unlockCalls = append(f.unlockCalls, path)
	f.unlocked[path] = true
	return nil
}

func (f *fakeVault) Lock(path string, forceFlush bool) error {
	queue := f.lockErrs[path]
	if len(queue) > 0 {
		err := queue[0]
		f.lockErrs[path] = queue[1:]
		if err != nil {
			return err
		}
	}
	delete(f.unlocked, path)
	return nil
}

func resource(path string) daemonconfig.Resource {
	return daemonconfig.Resource{Path: path, NeedEncryption: true}
}

func TestDaemonUseCaseConstructorMismatch(t *testing.T) {
	_, err := NewDaemonUseCase([]daemonconfig.Resource{resource("/a")}, newFakeVault(), []repository.GuardRepository{})
	if err == nil {
		t.Fatal("expected error on resource/guard count mismatch")
	}
}

func TestDaemonUseCaseStartLifecycle(t *testing.T) {
	vault := newFakeVault("/vault", "/plain")
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault"), {Path: "/plain", NeedEncryption: false}},
		vault,
		[]repository.GuardRepository{repo, newFakeGuardRepo()},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !vault.unlocked["/vault"] {
		t.Error("encrypted resource should be unlocked")
	}
	if vault.unlocked["/plain"] {
		t.Error("need_encryption:false resource must not be unlocked")
	}
	if !repo.started {
		t.Error("guard repo should be started")
	}
}

func TestDaemonUseCaseStartNotEncryptedFails(t *testing.T) {
	vault := newFakeVault() // nothing encrypted
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	err = d.Start()
	if err == nil {
		t.Fatal("expected error for non-encrypted resource")
	}
	if vault.unlocked["/vault"] {
		t.Error("must not unlock before verifying every resource")
	}
	if repo.started {
		t.Error("must not start guards when verification fails")
	}
}

func TestDaemonUseCaseStartFailsBeforeAnyUnlock(t *testing.T) {
	vault := newFakeVault("/vault")
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault"), resource("/missing")},
		vault,
		[]repository.GuardRepository{newFakeGuardRepo(), newFakeGuardRepo()},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	err = d.Start()
	if err == nil {
		t.Fatal("expected error for non-encrypted second resource")
	}
	if vault.unlocked["/vault"] {
		t.Error("unlock must not start until every resource is verified")
	}
}

func TestDaemonUseCaseStartGuardFailure(t *testing.T) {
	vault := newFakeVault("/vault")
	repo := newFakeGuardRepo()
	repo.startErr = errBoom
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	if err := d.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate guard error, got %v", err)
	}
}

func TestDaemonUseCaseEventsMerged(t *testing.T) {
	vault := newFakeVault("/vault")
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	repo.events <- guard.GuardEvent{Blocked: true}
	ev := <-d.Events()
	if ev.Resource != "/vault" {
		t.Errorf("resource tag = %q, want /vault", ev.Resource)
	}
	if !ev.Event.Blocked {
		t.Error("event payload not forwarded")
	}
	d.Stop()
}

func TestDaemonUseCaseStopLocksThenStopsGuards(t *testing.T) {
	vault := newFakeVault("/vault")
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop()
	if vault.unlocked["/vault"] {
		t.Error("resource should be locked on stop")
	}
	if !repo.stopped {
		t.Error("guard must stop after the vault is locked")
	}
}

func TestDaemonUseCaseStopRetriesOnBusy(t *testing.T) {
	vault := newFakeVault("/vault")
	// First pass Lock(false) succeeds; second pass Lock(true) returns
	// ErrKeyBusy twice, then succeeds.
	vault.lockErrs["/vault"] = []error{
		nil,
		repository.ErrKeyBusy,
		repository.ErrKeyBusy,
		nil,
	}
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop()
	if vault.unlocked["/vault"] {
		t.Error("resource should be locked after retries")
	}
	if len(vault.lockErrs["/vault"]) != 0 {
		t.Errorf("lock error queue not consumed: %v", vault.lockErrs["/vault"])
	}
	if !repo.stopped {
		t.Error("guard must be stopped after locking")
	}
}

func TestDaemonUseCaseStopKeyMissingIsSuccess(t *testing.T) {
	vault := newFakeVault("/vault")
	// First pass: already gone (ErrKeyMissing) — treated as success.
	// Second pass: also gone.
	vault.lockErrs["/vault"] = []error{repository.ErrKeyMissing, repository.ErrKeyMissing}
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop()
	if !repo.stopped {
		t.Error("guard must be stopped even when the key was already gone")
	}
}

func TestDaemonUseCaseStopRetryExhausted(t *testing.T) {
	vault := newFakeVault("/vault")
	// First pass: nil. Second pass: busy forever.
	busy := make([]error, maxLockRetries+1)
	for i := range busy {
		busy[i] = repository.ErrKeyBusy
	}
	vault.lockErrs["/vault"] = busy
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop() // must not hang; logs failure, still stops the guard
	if !repo.stopped {
		t.Error("guard must be stopped even when locking is exhausted")
	}
}

func TestDaemonUseCaseStopNonEncryptedResourceNotLocked(t *testing.T) {
	vault := newFakeVault()
	repo := newFakeGuardRepo()
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{{Path: "/plain", NeedEncryption: false}},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop()
	if len(vault.lockErrs) != 0 {
		t.Errorf("plain resource must not be locked: %v", vault.lockErrs)
	}
	if !repo.stopped {
		t.Error("guard must still stop")
	}
}

func TestDaemonUseCaseStopIdempotent(t *testing.T) {
	vault := newFakeVault("/vault")
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{newFakeGuardRepo()},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	d.Stop()
	// Second Stop must not double-lock (the queue would empty out and
	// panic the fake otherwise).
	d.Stop()
}

func TestDaemonUseCaseResources(t *testing.T) {
	res := []daemonconfig.Resource{resource("/vault")}
	d, err := NewDaemonUseCase(res, newFakeVault("/vault"), []repository.GuardRepository{newFakeGuardRepo()})
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if got := d.Resources(); len(got) != 1 || got[0].Path != "/vault" {
		t.Fatalf("Resources() = %v", got)
	}
}

// startDaemon is a helper that boots a daemon with the given vault,
// resources and guards.
func startDaemon(t *testing.T, vault *fakeVault, resources []daemonconfig.Resource, guards []repository.GuardRepository) DaemonUseCase {
	t.Helper()
	d, err := NewDaemonUseCase(resources, vault, guards)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return d
}

func TestDaemonUseCaseReloadLenMismatch(t *testing.T) {
	d := startDaemon(t, newFakeVault("/a"), []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{newFakeGuardRepo()})
	err := d.Reload([]daemonconfig.Resource{resource("/a")}, nil)
	if err == nil {
		t.Fatal("expected error on resource/guard count mismatch")
	}
}

func TestDaemonUseCaseReloadBeforeStart(t *testing.T) {
	vault := newFakeVault("/a")
	d, err := NewDaemonUseCase([]daemonconfig.Resource{resource("/a")}, vault, []repository.GuardRepository{newFakeGuardRepo()})
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}
	if err := d.Reload([]daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{newFakeGuardRepo()}); err == nil {
		t.Fatal("expected error when reloading before start")
	}
}

func TestDaemonUseCaseReloadKeptResources(t *testing.T) {
	vault := newFakeVault("/a")
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	next := newFakeGuardRepo()
	if err := d.Reload([]daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{next}); err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if !old.stopped {
		t.Error("the replaced guard must be stopped after the new one is attached")
	}
	if !next.started {
		t.Error("the new guard must be started")
	}
	if got := d.Resources(); len(got) != 1 || got[0].Path != "/a" {
		t.Fatalf("Resources() = %v", got)
	}

	// Events from the new guard must flow into the merged stream.
	next.events <- guard.GuardEvent{Blocked: true}
	select {
	case ev := <-d.Events():
		if ev.Resource != "/a" || !ev.Event.Blocked {
			t.Fatalf("unexpected event: %+v", ev)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no event from the new guard after reload")
	}
}

func TestDaemonUseCaseReloadAddedResourceUnlocked(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if !old.stopped {
		t.Error("the replaced guard for the kept resource must be stopped")
	}
	if !newA.started || !newB.started {
		t.Error("all new guards must be started")
	}
	if !vault.unlocked["/b"] {
		t.Error("newly added encrypted resource must be unlocked")
	}
}

func TestDaemonUseCaseReloadAddedAlreadyProvisionedNotUnlocked(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	vault.unlocked["/b"] = true // already unlocked by someone else
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	if err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	for _, p := range vault.unlockCalls {
		if p == "/b" {
			t.Error("already-provisioned resource must not be unlocked")
		}
	}
}

func TestDaemonUseCaseReloadAddedUnencryptedFailsKeepsOld(t *testing.T) {
	vault := newFakeVault("/a") // /b has no fscrypt policy
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	)
	if err == nil {
		t.Fatal("expected error for unencrypted added resource")
	}

	if old.stopped {
		t.Error("the old configuration must keep running after a failed reload")
	}
	if !newA.stopped || !newB.stopped {
		t.Error("newly built guards must be detached on failed reload")
	}
	if got := d.Resources(); len(got) != 1 || got[0].Path != "/a" {
		t.Fatalf("Resources() must still expose the old configuration: %v", got)
	}
}

func TestDaemonUseCaseReloadRollbackLocksUnlocked(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	newB.startErr = errBoom
	err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	)
	if err == nil {
		t.Fatal("expected error when a new guard fails to start")
	}

	if old.stopped {
		t.Error("a failed reload must not stop the old guards")
	}
	if vault.unlocked["/b"] {
		t.Error("a resource unlocked during the reload must be locked back on rollback")
	}
	// /a was unlocked at start; /b must have been unlocked (and locked
	// back) by the reload.
	if len(vault.unlockCalls) != 2 || vault.unlockCalls[1] != "/b" {
		t.Errorf("unexpected unlock calls: %v", vault.unlockCalls)
	}
}

func TestDaemonUseCaseReloadRemovedResource(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	oldA, oldB := newFakeGuardRepo(), newFakeGuardRepo()
	d := startDaemon(t, vault,
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{oldA, oldB},
	)

	newA := newFakeGuardRepo()
	if err := d.Reload([]daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{newA}); err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if !oldB.stopped {
		t.Error("guard of the removed resource must be stopped")
	}
	if got := d.Resources(); len(got) != 1 || got[0].Path != "/a" {
		t.Fatalf("Resources() = %v", got)
	}
	// The removed resource stays unlocked, matching ssh-guard: reload
	// only reconciles protection, never the fscrypt lifecycle.
	if !vault.unlocked["/b"] {
		t.Error("removed resource must not be locked by a reload")
	}
}

func TestDaemonUseCaseReloadThenStop(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	if err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	); err != nil {
		t.Fatalf("Reload: %v", err)
	}

	d.Stop()
	if !newA.stopped || !newB.stopped {
		t.Error("current guards must be stopped on shutdown after reload")
	}
	if vault.unlocked["/b"] {
		t.Error("resource unlocked by reload must be locked on shutdown")
	}
}
