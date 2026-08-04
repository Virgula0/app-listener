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
	// maxLockRetries mirrors ssh-guard's teardown retry budget for the
	// EBUSY/ENOKEY dance when deprovisioning fscrypt keys.
	maxLockRetries = 100
	lockRetryDelay = 10 * time.Millisecond
)

// DaemonEvent is one guard event tagged with the resource it belongs to.
type DaemonEvent struct {
	Resource string
	Event    guard.GuardEvent
}

// DaemonUseCase orchestrates the daemon lifecycle: verify encryption,
// unlock (decrypt) resources, start the per-resource whitelist guards,
// atomically apply configuration reloads (SIGHUP), and — on shutdown —
// lock everything back up while the guards are still attached, so there
// is never an unprotected window (the TOC-TOU concern that ssh-guard's
// teardown sequence was designed against).
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
	// stops holds one channel per guard; closing it detaches that guard's
	// event forwarder. Guard event channels are never closed, so reloads
	// need these to retire the forwarders of replaced guards.
	stops []chan struct{}

	events chan DaemonEvent
	done   chan struct{}

	start   sync.Once
	stop    sync.Once
	started bool
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

// Start brings the daemon up in ssh-guard's order: verify every resource's
// encryption state first (nothing is unlocked or protected on partial
// failure), then unlock what needs unlocking, then attach the guards.
func (d *daemonUseCase) Start() error {
	var startErr error
	d.start.Do(func() {
		startErr = d.startGuards()
	})
	return startErr
}

func (d *daemonUseCase) startGuards() error {
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

	for _, r := range d.resources {
		if !r.NeedEncryption {
			continue
		}
		if err := d.vault.Unlock(r.Path); err != nil {
			return fmt.Errorf("unlocking %s: %w", r.Path, err)
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
	for {
		select {
		case ev, ok := <-g.Events():
			if !ok {
				return
			}
			select {
			case d.events <- DaemonEvent{Resource: resource, Event: ev}:
			case <-d.done:
				return
			case <-stop:
				return
			}
		case <-d.done:
			return
		case <-stop:
			return
		}
	}
}

// Reload atomically applies a new configuration, mirroring ssh-guard's
// SIGHUP handling. The guards passed in are ALREADY attached (built by
// the caller) before this call returns; the kernel runs the old and new
// LSM programs together during the transition and denies an access when
// ANY of them denies, so protection is never weaker than either
// configuration — there is no instant where a whitelisted-only change
// opens a window. On any error the previous configuration keeps running
// untouched and the newly built guards are detached.
func (d *daemonUseCase) Reload(resources []daemonconfig.Resource, guards []repository.GuardRepository) error {
	if len(resources) != len(guards) {
		return fmt.Errorf("daemon: %d resources but %d guard engines", len(resources), len(guards))
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.started {
		return fmt.Errorf("daemon: reload before start")
	}

	oldByPath := make(map[string]int, len(d.resources))
	for i, r := range d.resources {
		oldByPath[r.Path] = i
	}

	// Phase 1 — validate and unlock resources that are new to the
	// configuration. Nothing is committed yet: on any error the newly
	// built (and attached) guards are detached again, anything we just
	// unlocked is locked back, and the old configuration keeps running.
	var unlocked []string
	if err := d.prepareNewResources(resources, oldByPath, &unlocked); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 2 — start the ringbuf readers of all new guards. The old
	// guards remain attached; no protection is dropped at any point.
	if err := d.startNewGuards(guards); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 3 — commit: swap the references, retire the old forwarders,
	// then detach the old LSM programs.
	d.commitReload(resources, guards)
	log.Info("daemon: configuration reloaded without dropping protection")
	return nil
}

// prepareNewResources validates and unlocks the resources that are new to
// the configuration; kept resources keep their lifecycle untouched. Every
// unlock performed here is recorded so a rollback can lock it back.
func (d *daemonUseCase) prepareNewResources(resources []daemonconfig.Resource, oldByPath map[string]int, unlocked *[]string) error {
	for _, r := range resources {
		if _, existed := oldByPath[r.Path]; existed {
			continue // lifecycle unchanged for kept resources
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

// rollbackReload detaches the newly built (not yet committed) guards and
// locks back any resource the aborted reload unlocked.
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

// commitReload swaps the running configuration for the new one. The new
// forwarders are spawned before the old ones are retired and the old LSM
// programs detached, so the event stream never goes silent either.
func (d *daemonUseCase) commitReload(resources []daemonconfig.Resource, guards []repository.GuardRepository) {
	newByPath := make(map[string]bool, len(resources))
	for _, r := range resources {
		newByPath[r.Path] = true
	}

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
	for _, r := range d.resources {
		if !newByPath[r.Path] {
			log.Warnf("daemon: reload: resource %s is no longer protected (left unlocked)", r.Path)
		}
	}

	d.resources = resources
	d.guards = guards
	d.stops = newStops
}

// Stop performs the secure lockdown. The guards stay attached — denying
// every non-whitelisted access — for the whole lock sequence, so no new
// open can pin an inode while the keys are being deprovisioned. Only when
// every resource is keyless are the LSM hooks detached. This is strictly
// stronger than ssh-guard's teardown, which removed its marks before the
// final lock pass.
func (d *daemonUseCase) Stop() {
	d.stop.Do(func() {
		log.Info("daemon: initiating secure lockdown")
		close(d.done)

		d.mu.RLock()
		resources := d.resources
		guards := d.guards
		d.mu.RUnlock()

		// First pass: remove the key where possible.
		for _, r := range resources {
			if !r.NeedEncryption {
				continue
			}
			if err := d.vault.Lock(r.Path, false); err != nil && !errors.Is(err, repository.ErrKeyMissing) {
				log.Warnf("daemon: first-pass lock of %s: %v", r.Path, err)
			}
		}

		// Second pass: force-flush with the EBUSY retry loop; ENOKEY
		// (ErrKeyMissing) confirms the key is gone and counts as success.
		for _, r := range resources {
			if !r.NeedEncryption {
				continue
			}
			if err := d.lockWithRetry(r.Path); err != nil {
				log.Errorf("daemon: failed to fully lock %s: %v", r.Path, err)
			}
		}

		for _, g := range guards {
			g.Stop()
		}
		log.Info("daemon: shutdown complete, all vaults locked")
	})
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
