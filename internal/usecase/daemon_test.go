package usecase

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/repository"
)

type fakeVault struct {
	mu sync.Mutex

	encrypted map[string]bool
	// unlocked is the fscrypt lock state: Unlock provisions, Lock
	// deprovisions.
	unlocked map[string]bool
	// unlockCalls records every Unlock invocation, for assertions.
	unlockCalls []string

	// lockErrs maps path -> queue of errors returned by successive Lock
	// calls; the queue is consumed from the front.
	lockErrs map[string][]error

	// busyForever: Lock on these paths always returns ErrKeyBusy (no
	// queue consumption) until cleared — simulates a process that keeps
	// holding files open for an unbounded time.
	busyForever map[string]bool

	unlockErr error
	checkErr  error

	// lockCalls counts every Lock invocation, for tests that must observe
	// the lockdown starting.
	lockCalls int
}

func newFakeVault(paths ...string) *fakeVault {
	v := &fakeVault{
		encrypted:   make(map[string]bool),
		unlocked:    make(map[string]bool),
		lockErrs:    make(map[string][]error),
		busyForever: make(map[string]bool),
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
	f.mu.Lock()
	defer f.mu.Unlock()

	f.lockCalls++

	if f.busyForever[path] {
		return repository.ErrKeyBusy
	}
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

// clearLockErrs makes subsequent Lock calls for path succeed, unblocking a
// lockdown that is stuck on a busy key.
func (f *fakeVault) clearLockErrs(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.lockErrs, path)
	delete(f.busyForever, path)
}

// isUnlocked is the mutex-safe lock-state accessor for concurrent tests.
func (f *fakeVault) isUnlocked(path string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.unlocked[path]
}

// lockCallCount returns the number of Lock invocations observed, for
// tests that must wait until the lockdown has started.
func (f *fakeVault) lockCallCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lockCalls
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
	if !repo.populated {
		t.Error("guard inode map should be populated after unlock and before start")
	}
	if !repo.resolved {
		t.Error("deferred binaries should be resolved after populate and before start")
	}
}

func TestDaemonUseCaseResolveFailureFailsStart(t *testing.T) {
	// Set after unlock+populate, so a failed resolution must leave the
	// resource in the pre-unlock state once the caller runs Stop().
	vault := newFakeVault("/vault")
	repo := newFakeGuardRepo()
	repo.resolveErr = errBoom
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	if err := d.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate resolution error, got %v", err)
	}
	if !vault.unlocked["/vault"] {
		t.Error("resource should have been unlocked before the resolution attempt")
	}
	if !repo.populated {
		t.Error("inode population should precede deferred-binary resolution")
	}
	if repo.started {
		t.Error("guard must not be started after resolution failure")
	}

	d.Stop()
	if vault.unlocked["/vault"] {
		t.Error("resource must be locked back after a failed start")
	}
	if !repo.stopped {
		t.Error("guard must be stopped after a failed start")
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

func TestDaemonUseCaseStopBlocksUntilAllLocked(t *testing.T) {
	vault := newFakeVault("/vault")
	// A process holds files open in the vault for an unbounded time: every
	// Lock returns ErrKeyBusy until the test "closes" the files.
	vault.busyForever["/vault"] = true
	repo := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/vault")}, []repository.GuardRepository{repo})

	stopDone := make(chan struct{})
	go func() {
		d.Stop()
		close(stopDone)
	}()

	// Stop has exhausted its retry budget by now; it must NOT give up and
	// detach the guard while the vault is still unlocked.
	time.Sleep(lockRetryDelay * time.Duration(maxLockRetries+10))
	select {
	case <-stopDone:
		t.Fatal("Stop returned while the vault was still locked: fail-open regression")
	default:
	}
	if repo.isStopped() {
		t.Error("guard must stay attached while the vault is still unlocked")
	}
	if !vault.isUnlocked("/vault") {
		t.Error("vault should still be provisioned (busy) at this point")
	}

	// The pinning process finally closes its files: the vault locks and
	// Stop completes, detaching the guards only afterwards.
	vault.clearLockErrs("/vault")
	select {
	case <-stopDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not complete once the vault could be locked")
	}
	if vault.isUnlocked("/vault") {
		t.Error("vault must be locked once Stop returned")
	}
	if !repo.isStopped() {
		t.Error("guard must be stopped after the vault is locked")
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

func TestDaemonUseCaseStartPopulateFailure(t *testing.T) {
	vault := newFakeVault("/vault")
	repo := newFakeGuardRepo()
	repo.populateErr = errBoom
	d, err := NewDaemonUseCase(
		[]daemonconfig.Resource{resource("/vault")},
		vault,
		[]repository.GuardRepository{repo},
	)
	if err != nil {
		t.Fatalf("NewDaemonUseCase: %v", err)
	}

	if err := d.Start(); !errors.Is(err, errBoom) {
		t.Fatalf("Start should propagate populate error, got %v", err)
	}
	if !vault.unlocked["/vault"] {
		t.Error("resource should have been unlocked before the populate attempt")
	}

	// The caller (runDaemon) responds to a failed Start with Stop(): the
	// already-unlocked resource must be locked back and the guard stopped.
	d.Stop()
	if vault.unlocked["/vault"] {
		t.Error("resource must be locked back after a failed start")
	}
	if !repo.stopped {
		t.Error("guard must be stopped after a failed start")
	}
}

func TestDaemonUseCaseReloadPopulateFailure(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	newA, newB := newFakeGuardRepo(), newFakeGuardRepo()
	newB.populateErr = errBoom
	err := d.Reload(
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{newA, newB},
	)
	if err == nil {
		t.Fatal("expected error when a new guard fails to populate")
	}

	if old.stopped {
		t.Error("a failed reload must not stop the old guards")
	}
	if !newA.stopped || !newB.stopped {
		t.Error("newly built guards must be detached on failed reload")
	}
	if vault.unlocked["/b"] {
		t.Error("a resource unlocked during the reload must be locked back on rollback")
	}
	if got := d.Resources(); len(got) != 1 || got[0].Path != "/a" {
		t.Fatalf("Resources() must still expose the old configuration: %v", got)
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

func TestDaemonUseCaseReloadRemovedResourceRefused(t *testing.T) {
	vault := newFakeVault("/a", "/b")
	oldA, oldB := newFakeGuardRepo(), newFakeGuardRepo()
	d := startDaemon(t, vault,
		[]daemonconfig.Resource{resource("/a"), resource("/b")},
		[]repository.GuardRepository{oldA, oldB},
	)

	newA := newFakeGuardRepo()
	err := d.Reload([]daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{newA})
	if err == nil {
		t.Fatal("reload that removes a protected resource must be refused")
	}
	if !strings.Contains(err.Error(), "/b") {
		t.Errorf("refusal must name the dropped resource, got: %v", err)
	}

	// The old configuration is untouched and keeps running.
	if oldA.stopped || oldB.stopped {
		t.Error("the running guards must be kept after a refused reload")
	}
	if !newA.isStopped() {
		t.Error("the newly built guard must be detached after a refused reload")
	}
	if got := d.Resources(); len(got) != 2 {
		t.Fatalf("Resources() must still expose the full configuration: %v", got)
	}
	// /b stays protected and unlocked (it was unlocked at start) — the
	// refusal leaves the fscrypt lifecycle untouched.
	if !vault.unlocked["/b"] {
		t.Error("/b must remain untouched by a refused reload")
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

func TestDaemonUseCaseReloadRefusedDuringShutdown(t *testing.T) {
	vault := newFakeVault("/a")
	// A process holds files open in /a, so the lockdown blocks in its
	// force-flush loop. A SIGHUP-driven reload arrives in that window: it
	// must be refused — attaching new guards or unlocking resources during
	// shutdown could leave plaintext behind after the daemon exits.
	vault.busyForever["/a"] = true
	old := newFakeGuardRepo()
	d := startDaemon(t, vault, []daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{old})

	stopDone := make(chan struct{})
	go func() {
		d.Stop()
		close(stopDone)
	}()

	// Wait until the lockdown has started (it holds the daemon's write
	// lock for its whole duration).
	deadline := time.Now().Add(5 * time.Second)
	for vault.lockCallCount() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("Stop never started the lockdown")
		}
		time.Sleep(10 * time.Millisecond)
	}

	unlockCallsAfterStart := len(vault.unlockCalls)

	// The reload is serialized behind the lockdown's write lock and must
	// come back refused once the lockdown completes.
	next := newFakeGuardRepo()
	reloadDone := make(chan error, 1)
	go func() {
		reloadDone <- d.Reload([]daemonconfig.Resource{resource("/a")}, []repository.GuardRepository{next})
	}()

	// While the lockdown is still blocked, the reload must not have made
	// any progress: no unlock, no commit.
	select {
	case err := <-reloadDone:
		t.Fatalf("reload completed while the lockdown was still running: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	if len(vault.unlockCalls) != unlockCallsAfterStart {
		t.Fatalf("reload unlocked resources during shutdown: %v", vault.unlockCalls)
	}

	// The pinning process closes its files: the lockdown finishes and the
	// waiting reload is refused.
	vault.clearLockErrs("/a")
	select {
	case <-stopDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not complete once the vault could be locked")
	}
	select {
	case err := <-reloadDone:
		if err == nil {
			t.Fatal("reload must be refused once shutdown is in progress")
		}
		if !strings.Contains(err.Error(), "shutdown") {
			t.Errorf("refusal should mention the shutdown, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("reload never returned after the lockdown completed")
	}

	if !next.isStopped() {
		t.Error("the new guard must be detached after a refused reload")
	}
	if vault.isUnlocked("/a") {
		t.Error("the resource must stay locked after the shutdown")
	}
}
