package usecase

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/repository"
)

const (
	// maxLockRetries bounds the EBUSY retry loop when force-flushing an
	// fscrypt lock during key deprovisioning.
	maxLockRetries = 100
	lockRetryDelay = 10 * time.Millisecond
	// resyncMinInterval throttles post-denial re-syncs: a denial storm must
	// not re-stat every whitelisted binary on each event.
	resyncMinInterval = 2 * time.Second
	// resyncSweepEvery is the background re-sync interval, catching
	// replacements that never produced a denied access (swapped while the
	// old binary was still running).
	resyncSweepEvery = 30 * time.Second
)

// DaemonEvent is one guard event tagged with the resource it belongs to.
type DaemonEvent struct {
	Resource string
	Event    guard.GuardEvent
}

// DaemonUseCase orchestrates the daemon lifecycle: verify encryption, unlock, run guards,
// apply atomic SIGHUP reloads and a secure lockdown. Ordering is attach → unlock → populate
// → resolve → re-sync, so resources are never readable without live protection (no TOCTOU window).
type DaemonUseCase interface {
	Start() error
	Reload(resources []daemonconfig.Resource, guards []repository.GuardRepository) error
	Stop()
	Events() <-chan DaemonEvent
	Resources() []daemonconfig.Resource
}

type daemonUseCase struct {
	mu        sync.RWMutex
	resources []daemonconfig.Resource
	vault     repository.Vault
	guards    []repository.GuardRepository
	// stops holds one close-channel per guard to retire its event forwarder
	// on reload (guard event channels themselves are never closed).
	stops []chan struct{}

	events chan DaemonEvent
	done   chan struct{}

	start   sync.Once
	stop    sync.Once
	started bool
	// stopping is set once lockdown begins; Reload refuses when set, so a
	// SIGHUP mid-shutdown can never attach guards or unlock resources.
	stopping bool
}

func NewDaemonUseCase(resources []daemonconfig.Resource, vault repository.Vault, guards []repository.GuardRepository) (DaemonUseCase, error) {
	if len(resources) != len(guards) {
		return nil, fmt.Errorf("daemon: %d resources but %d guard engines", len(resources), len(guards))
	}
	return &daemonUseCase{
		resources: resources,
		vault:     vault,
		guards:    guards,
		events:    make(chan DaemonEvent, 1024),
		done:      make(chan struct{}),
	}, nil
}

// populateInodes scans the already-attached guards' resource trees into their inode maps.
// It runs after the resources are unlocked (a locked fscrypt tree cannot be stat'ed); guards
// are attached by the caller before Start, preserving attach → unlock → populate ordering.
func (d *daemonUseCase) populateInodes() error {
	for i := range d.guards {
		if err := d.guards[i].PopulateInodes(); err != nil {
			return fmt.Errorf("populating guard for %s: %w", d.resources[i].Path, err)
		}
	}
	return nil
}

func (d *daemonUseCase) Start() error {
	var startErr error
	d.start.Do(func() {
		startErr = d.startGuards()
	})
	return startErr
}

// verifyEncryptionStates checks every resource's encryption state against its
// need_encryption setting before anything is unlocked or protected: a misconfigured
// resource aborts the start.
func (d *daemonUseCase) verifyEncryptionStates() error {
	for _, r := range d.resources {
		encrypted, err := d.vault.IsEncrypted(r.Path)
		if err != nil {
			return fmt.Errorf("checking encryption on %s: %w", r.Path, err)
		}
		if r.NeedEncryption && !encrypted {
			return fmt.Errorf("directory %s is NOT encrypted: run the fscrypt migration first or set need_encryption: false", r.Path)
		}
		if !r.NeedEncryption && encrypted {
			log.Warnf("resource %s is encrypted but need_encryption: false \u2014 leaving it locked", r.Path)
		}
	}
	return nil
}

func (d *daemonUseCase) startGuards() error {
	if err := d.verifyEncryptionStates(); err != nil {
		return err
	}

	for _, r := range d.resources {
		if !r.NeedEncryption {
			continue
		}
		if err := d.vault.Unlock(r.Path); err != nil {
			return fmt.Errorf("unlocking %s: %w", r.Path, err)
		}
	}

	// Guards are already attached (built before Start): the scan runs with
	// protection live and resources readable — no unlocked-and-unprotected window.
	if err := d.populateInodes(); err != nil {
		return err
	}

	// Fail-closed contract: deferred whitelist entries stay absent from the BPF
	// whitelist — denied — until resolved here, so unlock never precedes protection
	// (attach → unlock → populate → resolve). The re-sync right after admits binaries
	// replaced on disk while the daemon was down.
	for i := range d.guards {
		if err := d.guards[i].ResolvePendingBinaries(); err != nil {
			return fmt.Errorf("resolving deferred binaries for %s: %w", d.resources[i].Path, err)
		}
		if _, err := d.guards[i].ReSyncBinaries(); err != nil {
			return fmt.Errorf("re-syncing binary whitelist for %s: %w", d.resources[i].Path, err)
		}
	}

	for i := range d.guards {
		if err := d.guards[i].Start(); err != nil {
			return fmt.Errorf("starting guard for %s: %w", d.resources[i].Path, err)
		}
	}

	for i := range d.guards {
		stop := make(chan struct{})
		d.stops = append(d.stops, stop)
		go d.forwardEvents(d.resources[i].Path, d.guards[i], stop)
	}
	d.started = true
	return nil
}

func (d *daemonUseCase) forwardEvents(resource string, g repository.GuardRepository, stop <-chan struct{}) {
	sweep := time.NewTicker(resyncSweepEvery)
	defer sweep.Stop()
	var lastResync time.Time
	for {
		select {
		case ev, ok := <-g.Events():
			if !ok {
				return
			}
			// A denial usually means the binary was replaced in place;
			// re-sync (throttled by resyncMinInterval) to admit the new
			// inode instead of re-statting the whitelist per event.
			if ev.Blocked && time.Since(lastResync) >= resyncMinInterval {
				if _, err := g.ReSyncBinaries(); err != nil {
					log.Errorf("daemon: re-syncing binary whitelist for %s: %v", resource, err)
				}
				lastResync = time.Now()
			}
			select {
			case d.events <- DaemonEvent{Resource: resource, Event: ev}:
			case <-d.done:
				return
			case <-stop:
				return
			}
		case <-sweep.C:
			// Periodic sweep: catches replacements that never produced a denial.
			if _, err := g.ReSyncBinaries(); err != nil {
				log.Errorf("daemon: periodic binary re-sync for %s: %v", resource, err)
			}
			lastResync = time.Now()
		case <-d.done:
			return
		case <-stop:
			return
		}
	}
}

// Reload atomically applies a new configuration: old and new LSM programs run concurrently
// and ANY deny wins, so protection is never weaker than either configuration; on error the
// old config keeps running and the new guards are detached. Removing a resource is refused —
// it would be left unlocked and unguarded; dropping protection requires stop/edit/start.
func (d *daemonUseCase) Reload(resources []daemonconfig.Resource, guards []repository.GuardRepository) error {
	if len(resources) != len(guards) {
		return fmt.Errorf("daemon: %d resources but %d guard engines", len(resources), len(guards))
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.started {
		return fmt.Errorf("daemon: reload before start")
	}

	if d.stopping {
		d.rollbackReload(guards, nil)
		return fmt.Errorf("daemon: reload refused: shutdown in progress — restart the daemon to apply the new configuration")
	}

	// Phase 0 — refuse to drop any protected resource (see Reload). Nothing
	// is committed yet, so the new guards are detached and the old config
	// keeps running.
	newByPath := make(map[string]bool, len(resources))
	for _, r := range resources {
		newByPath[r.Path] = true
	}
	var removed []string
	for _, r := range d.resources {
		if !newByPath[r.Path] {
			removed = append(removed, r.Path)
		}
	}
	if len(removed) > 0 {
		d.rollbackReload(guards, nil)
		return fmt.Errorf("reload refused: %v would be left unlocked and unprotected — stop the daemon, edit the config, and start it again to drop them", removed)
	}

	oldByPath := make(map[string]int, len(d.resources))
	for i, r := range d.resources {
		oldByPath[r.Path] = i
	}

	// Phase 1 — validate and unlock resources new to the config; nothing is
	// committed yet, so errors roll back (detach new guards, re-lock unlocks).
	// Guards were attached by the caller: unlocks always have live protection.
	var unlocked []string
	if err := d.prepareNewResources(resources, oldByPath, &unlocked); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 2 — populate the new guards' inode maps (hooks live, plaintext
	// readable — no protection gap), then resolve deferred entries and
	// re-sync: same fail-closed ordering as startup.
	if err := d.prepareGuards(resources, guards); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 3 — start new guards' ringbuf readers; old guards stay attached.
	if err := d.startNewGuards(guards); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 4 — commit: swap references, retire old forwarders, detach old LSM programs.
	d.commitReload(resources, guards)
	log.Info("daemon: configuration reloaded without dropping protection")
	return nil
}

// prepareNewResources validates and unlocks resources new to the configuration; kept
// resources are untouched. Each unlock is recorded so rollback can lock it back.
func (d *daemonUseCase) prepareNewResources(resources []daemonconfig.Resource, oldByPath map[string]int, unlocked *[]string) error {
	for _, r := range resources {
		if _, existed := oldByPath[r.Path]; existed {
			continue
		}
		if err := d.prepareAddedResource(r, unlocked); err != nil {
			return err
		}
	}
	return nil
}

func (d *daemonUseCase) prepareAddedResource(r daemonconfig.Resource, unlocked *[]string) error {
	encrypted, err := d.vault.IsEncrypted(r.Path)
	if err != nil {
		return fmt.Errorf("reload: checking encryption on %s: %w", r.Path, err)
	}
	if r.NeedEncryption && !encrypted {
		return fmt.Errorf("reload: directory %s is NOT encrypted: run the fscrypt migration first or set need_encryption: false", r.Path)
	}
	if !r.NeedEncryption && encrypted {
		log.Warnf("daemon: reload: resource %s is encrypted but need_encryption: false \u2014 leaving it locked", r.Path)
		return nil
	}
	provisioned, err := d.vault.IsProvisioned(r.Path)
	if err != nil {
		return fmt.Errorf("reload: checking lock state of %s: %w", r.Path, err)
	}
	if provisioned {
		return nil
	}
	if err := d.vault.Unlock(r.Path); err != nil {
		return fmt.Errorf("reload: unlocking %s: %w", r.Path, err)
	}
	*unlocked = append(*unlocked, r.Path)
	return nil
}

func (d *daemonUseCase) startNewGuards(guards []repository.GuardRepository) error {
	for _, g := range guards {
		if err := g.Start(); err != nil {
			return fmt.Errorf("reload: starting new guard: %w", err)
		}
	}
	return nil
}

// prepareGuards fills the freshly attached guards' inode maps, resolves deferred
// whitelist entries and re-syncs replaced binaries; it runs post-unlock, so scans
// see plaintext while hooks are live.
func (d *daemonUseCase) prepareGuards(resources []daemonconfig.Resource, guards []repository.GuardRepository) error {
	for i := range guards {
		if err := guards[i].PopulateInodes(); err != nil {
			return fmt.Errorf("reload: populating guard inodes for %s: %w", resources[i].Path, err)
		}
		if err := guards[i].ResolvePendingBinaries(); err != nil {
			return fmt.Errorf("reload: resolving deferred binaries for %s: %w", resources[i].Path, err)
		}
		if _, err := guards[i].ReSyncBinaries(); err != nil {
			return fmt.Errorf("reload: re-syncing binary whitelist for %s: %w", resources[i].Path, err)
		}
	}
	return nil
}

// rollbackReload detaches the not-yet-committed new guards and locks back
// any resource the aborted reload unlocked.
func (d *daemonUseCase) rollbackReload(guards []repository.GuardRepository, unlocked []string) {
	for _, g := range guards {
		g.Stop()
	}
	for _, path := range unlocked {
		if err := d.lockWithRetry(path); err != nil {
			log.Errorf("daemon: rollback: failed to re-lock %s: %v", path, err)
		}
	}
}

// commitReload swaps in the new configuration; new forwarders spawn before old ones
// retire and old LSM programs detach, so the event stream never goes silent either.
func (d *daemonUseCase) commitReload(resources []daemonconfig.Resource, guards []repository.GuardRepository) {
	newStops := make([]chan struct{}, len(guards))
	for i := range guards {
		stop := make(chan struct{})
		newStops[i] = stop
		go d.forwardEvents(resources[i].Path, guards[i], stop)
	}

	for _, oldStop := range d.stops {
		close(oldStop)
	}
	for _, g := range d.guards {
		g.Stop()
	}

	d.resources = resources
	d.guards = guards
	d.stops = newStops
}

// Stop performs the secure lockdown: guards stay attached (denying all non-whitelisted
// access) until every resource is keyless; only then are LSM hooks detached. A resource
// whose key cannot be deprovisioned blocks shutdown indefinitely — the tree is never
// left unlocked and unguarded, and the daemon logs how to find the pinning process.
func (d *daemonUseCase) Stop() {
	d.stop.Do(func() {
		log.Info("daemon: initiating secure lockdown")

		// Hold the write lock through lockdown: Reload serializes behind it and
		// refuses once stopping is set — no SIGHUP attaches/unlocks mid-shutdown.
		d.mu.Lock()
		defer d.mu.Unlock()

		d.stopping = true
		close(d.done)

		resources := d.resources
		guards := d.guards

		// First pass: remove the key where possible.
		for _, r := range resources {
			if !r.NeedEncryption {
				continue
			}
			if err := d.vault.Lock(r.Path, false); err != nil && !errors.Is(err, repository.ErrKeyMissing) {
				log.Warnf("daemon: first-pass lock of %s: %v", r.Path, err)
			}
		}

		// Second pass: force-flush with EBUSY retries; ENOKEY (ErrKeyMissing) counts
		// as success. Busy resources retry forever — detaching guards here would leave
		// the tree unlocked and unguarded, so shutdown never gives up.
		pending := make([]string, 0, len(resources))
		for _, r := range resources {
			if r.NeedEncryption {
				pending = append(pending, r.Path)
			}
		}
		d.lockUntilAllKeyless(pending)

		for _, g := range guards {
			g.Stop()
		}
		log.Info("daemon: shutdown complete, all vaults locked")
	})
}

// lockUntilAllKeyless retries the force-flush lock until every resource is keyless; it
// never returns while one is unlocked — guards stay attached denying all access, so the
// tree is never unprotected and shutdown just waits for file pins to close.
func (d *daemonUseCase) lockUntilAllKeyless(pending []string) {
	for len(pending) > 0 {
		var still []string
		for _, path := range pending {
			if err := d.lockWithRetry(path); err != nil {
				log.Errorf("daemon: %s is still unlocked: a process holds open files in it (investigate with: lsof +D %s, fuser -v %s): %v",
					path, path, path, err)
				still = append(still, path)
			}
		}
		if len(still) == 0 {
			break
		}
		log.Errorf("daemon: shutdown blocked on %d resource(s) still unlocked — guards stay attached and deny all access until they lock", len(still))
		pending = still
		time.Sleep(time.Second)
	}
}

func (d *daemonUseCase) lockWithRetry(path string) error {
	for i := 0; i < maxLockRetries; i++ {
		err := d.vault.Lock(path, true)
		if err == nil || errors.Is(err, repository.ErrKeyMissing) {
			log.Infof("daemon: locked fscrypt directory fully: %s", path)
			return nil
		}
		if !errors.Is(err, repository.ErrKeyBusy) {
			return err
		}
		time.Sleep(lockRetryDelay)
	}
	return fmt.Errorf("could not fully lock %s: still busy after %d retries", path, maxLockRetries)
}

// Events returns the merged, per-resource tagged event stream.
func (d *daemonUseCase) Events() <-chan DaemonEvent {
	return d.events
}

// Resources exposes the configured resources for presentation.
func (d *daemonUseCase) Resources() []daemonconfig.Resource {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]daemonconfig.Resource, len(d.resources))
	copy(out, d.resources)
	return out
}
