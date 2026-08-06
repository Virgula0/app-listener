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
// failure), then unlock what needs unlocking, then populate the inode maps
// of the already-attached guards, then start the event readers. The guards
// are built (LSM hooks attached) by the caller BEFORE Start is invoked, so
// a resource is never unlocked while its hooks are not yet live — the
// window the original ssh-guard closed by attaching fanotify marks before
// injecting keys.
// populateInodes scans the already-attached guards' resource trees into
// their inode maps. It must run only after the resources are unlocked: a
// locked fscrypt tree cannot be stat'ed.
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

	// The guards are already attached (the builder ran before Start), so
	// the inode scan below happens with protection live and the resources
	// readable: no unlocked-and-unprotected window.
	if err := d.populateInodes(); err != nil {
		return err
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
//
// Removing a resource from the configuration is refused: its directory
// would be left unlocked (the key is only deprovisioned at shutdown) with
// no guard attached. Dropping protection requires a stop/edit/start
// cycle, so the operator deliberately accepts the exposure.
func (d *daemonUseCase) Reload(resources []daemonconfig.Resource, guards []repository.GuardRepository) error {
	if len(resources) != len(guards) {
		return fmt.Errorf("daemon: %d resources but %d guard engines", len(resources), len(guards))
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.started {
		return fmt.Errorf("daemon: reload before start")
	}

	// Phase 0 — refuse to drop any protected resource (see the comment on
	// Reload). Nothing has been committed yet, so the newly built guards
	// are detached and the old configuration keeps running.
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

	// Phase 1 — validate and unlock resources that are new to the
	// configuration. Nothing is committed yet: on any error the newly
	// built (and attached) guards are detached again, anything we just
	// unlocked is locked back, and the old configuration keeps running.
	// The guards passed in were built (hooks attached) by the caller
	// before this call, so the unlocks below never happen without live
	// protection — the same ordering ssh-guard guarantees at startup.
	var unlocked []string
	if err := d.prepareNewResources(resources, oldByPath, &unlocked); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 2 — populate the inode maps of the new guards. Every new
	// resource is unlocked (or was already readable) and the new hooks are
	// attached, so the scan sees the plaintext without a protection gap.
	for i := range guards {
		if err := guards[i].PopulateInodes(); err != nil {
			d.rollbackReload(guards, unlocked)
			return fmt.Errorf("reload: populating guard inodes for %s: %w", resources[i].Path, err)
		}
	}

	// Phase 3 — start the ringbuf readers of all new guards. The old
	// guards remain attached; no protection is dropped at any point.
	if err := d.startNewGuards(guards); err != nil {
		d.rollbackReload(guards, unlocked)
		return err
	}

	// Phase 4 — commit: swap the references, retire the old forwarders,
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

// Stop performs the secure lockdown. The guards stay attached — denying
// every non-whitelisted access — for the whole lock sequence, so no new
// open can pin an inode while the keys are being deprovisioned. Only when
// every resource is keyless are the LSM hooks detached. This is strictly
// stronger than ssh-guard's teardown, which removed its marks before the
// final lock pass.
//
// A resource whose key cannot be deprovisioned (a process still holding
// open files in it) BLOCKS the shutdown: the guards remain attached and
// deny every non-whitelisted access, so the tree is never left unlocked
// and unguarded. Shutdown only completes when every resource is keyless —
// the daemon keeps retrying and logs how to find the pinning process.
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
		// Resources that stay busy are retried indefinitely — detaching the
		// guards here would leave the tree unlocked and unguarded, so
		// shutdown never gives up.
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

// lockUntilAllKeyless keeps retrying the force-flush lock on every resource
// until each one is keyless. It never returns while a resource is still
// unlocked: the guards are still attached at that point and deny all access,
// so the tree is never unprotected, and shutdown simply waits for whoever
// pins the files to close them.
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
