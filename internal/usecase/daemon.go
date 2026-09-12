package usecase

import (
	"errors"
	"fmt"
	"os"
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
	// resyncSweepEvery is the background sweep interval, catching binary
	// replacements that never produced a denied access (swapped while the old
	// binary was still running) and recreated single-file watch roots. The
	// sweep is fingerprint-gated (SweepInodes) so each tick is cheap when
	// nothing changed.
	resyncSweepEvery = 30 * time.Second

	// Rollback lock-back budget: bounded (unlike Stop's infinite wait)
	// because rollback runs on the SIGHUP handler, which must keep serving
	// signals.
	maxRollbackLockRounds  = 3
	rollbackLockRetryDelay = 500 * time.Millisecond
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
	// GrantEditAccess widens the guard of one configured resource so the
	// app-listener binary (uid 0 only) may modify it for an authenticated
	// live edit-protected session. Only one grant may be active at a time.
	// The returned revoke restores the read-only baseline; it is idempotent.
	GrantEditAccess(resourcePath string) (revoke func() error, err error)
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
	// orphans holds not-yet-committed guards from an aborted reload whose
	// freshly unlocked resources could not be locked back within the
	// bounded rollback budget: they stay attached denying access and are
	// retired by Stop after its own lockdown (see rollbackReload).
	orphans []repository.GuardRepository
	// editGrantActive is set while a live edit-protected session holds a
	// write grant on one resource (see GrantEditAccess); only one at a time.
	editGrantActive bool
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

// encryptionRootOf returns the vault root whose fscrypt key lifecycle governs
// resource r: the `watch:`-group's section path when set, the resource path
// itself otherwise (ungrouped resources are their own encryption root).
func encryptionRootOf(r *daemonconfig.Resource) string {
	if r.EncryptionRoot != "" {
		return r.EncryptionRoot
	}
	return r.Path
}

// uniqueEncryptionRoots deduplicates the encryption roots of the encrypted
// resources: grouped watch paths share one vault root, and the fscrypt key
// lifecycle must run exactly once per root.
func uniqueEncryptionRoots(resources []daemonconfig.Resource) []string {
	seen := make(map[string]bool)
	var roots []string
	for i := range resources {
		r := &resources[i]
		if !r.NeedEncryption {
			continue
		}
		root := encryptionRootOf(r)
		if !seen[root] {
			seen[root] = true
			roots = append(roots, root)
		}
	}
	return roots
}

// partitionEncryptionRoots splits roots into fscrypt directories and
// file-vault regular files. Used only to sequence startGuards' unlock pass
// into two explicit phases — directories (unchanged kernel-keyring cycle)
// then files (in-place, same-inode cycle, see internal/fscrypt/filevault.go)
// — so the two different unlock mechanisms, and that neither can affect the
// other's guard state, is visible by inspection. Each Vault.Unlock call is
// independently safe regardless of interleaving; this split is a clarity
// and auditability choice, not a correctness requirement. A stat failure
// (e.g. the root does not exist yet) falls into dirs, matching the
// pre-partition behavior of just calling Unlock and letting it report the
// real error.
func partitionEncryptionRoots(roots []string) (dirs, files []string) {
	for _, root := range roots {
		if info, err := os.Stat(root); err == nil && info.Mode().IsRegular() {
			files = append(files, root)
			continue
		}
		dirs = append(dirs, root)
	}
	return dirs, files
}

// vaultOpForGuard runs a Vault.Unlock/Lock call on root, temporarily
// widening g's self-access when root is a file-vault target (see
// guard.WithSelfVaultAccess) — a directory root's fscrypt lifecycle never
// touches file content (pure kernel-keyring operations), so it always runs
// unwidened. g may be nil (no guard found for root): op still runs, so the
// real (denied) error surfaces instead of a silent skip.
func vaultOpForGuard(g repository.GuardRepository, root string, op func() error) error {
	info, statErr := os.Stat(root)
	if statErr != nil || !info.Mode().IsRegular() || g == nil {
		return op()
	}
	return g.WithSelfVaultAccess(op)
}

// vaultOpOnRoot is vaultOpForGuard, looking the guard up by encryption root
// in resources/guards (index-aligned — Reload enforces this, and
// d.resources/d.guards always are). A regular-file resource is never
// grouped (WatchRelPaths targets directories only), so exactly one
// resource/guard matches a file root — no ambiguity to resolve.
func vaultOpOnRoot(resources []daemonconfig.Resource, guards []repository.GuardRepository, root string, op func() error) error {
	for i := range resources {
		if encryptionRootOf(&resources[i]) == root {
			return vaultOpForGuard(guards[i], root, op)
		}
	}
	return vaultOpForGuard(nil, root, op)
}

// verifyEncryptionStates checks every resource's encryption state against its
// need_encryption setting before anything is unlocked or protected: a misconfigured
// resource aborts the start.
func (d *daemonUseCase) verifyEncryptionStates() error {
	for i := range d.resources {
		r := &d.resources[i]
		encrypted, err := d.vault.IsEncrypted(encryptionRootOf(r))
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

	// Grouped watch paths share one encryption root: the key is provisioned
	// exactly once per vault (unlock is vault-wide by fscrypt semantics).
	// Directories and file-vault regular files are unlocked in two explicit
	// passes (see partitionEncryptionRoots) — guards for BOTH kinds are
	// already attached at this point (built before Start), so there is no
	// unlocked-and-unprotected window for either. A file root's guard self
	// access is widened only for the call itself (vaultOpOnRoot ->
	// guard.WithSelfVaultAccess): the in-place unlock reads+rewrites the
	// file's own bytes, which the guard's baseline self mask (open/read/stat
	// only) does not permit.
	dirRoots, fileRoots := partitionEncryptionRoots(uniqueEncryptionRoots(d.resources))
	for _, root := range dirRoots {
		if err := d.vault.Unlock(root); err != nil {
			return fmt.Errorf("unlocking %s: %w", root, err)
		}
	}
	for _, root := range fileRoots {
		if err := vaultOpOnRoot(d.resources, d.guards, root, func() error { return d.vault.Unlock(root) }); err != nil {
			return fmt.Errorf("unlocking %s: %w", root, err)
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
			// The raw block-device gate names a device, not this resource:
			// label it as such and skip the re-sync (it is never an
			// in-place binary replacement).
			label := resource
			if ev.RawDevice {
				label = guard.RawDeviceResourceLabel
			}
			// A denial usually means the binary was replaced in place;
			// re-sync (throttled by resyncMinInterval) to admit the new
			// inode instead of re-statting the whitelist per event.
			if ev.Blocked && !ev.RawDevice && time.Since(lastResync) >= resyncMinInterval {
				if _, err := g.ReSyncBinaries(); err != nil {
					log.Errorf("daemon: re-syncing binary whitelist for %s: %v", resource, err)
				}
				lastResync = time.Now()
			}
			select {
			case d.events <- DaemonEvent{Resource: label, Event: ev}:
			case <-d.done:
				return
			case <-stop:
				return
			}
		case <-sweep.C:
			d.periodicSweep(resource, g)
			lastResync = time.Now()
		case <-d.done:
			return
		case <-stop:
			return
		}
	}
}

// periodicSweep re-syncs the binary whitelist (catching replacements that
// never produced a denial) and refreshes the inode map. SweepInodes only does
// real work when the watch root's fingerprint moved (a recreated single-file
// root — sqlite journals — or a directory that gained/lost a top-level
// entry); an unconditional full re-walk every tick was a large share of the
// daemon's steady-state CPU.
func (d *daemonUseCase) periodicSweep(resource string, g repository.GuardRepository) {
	if _, err := g.ReSyncBinaries(); err != nil {
		log.Errorf("daemon: periodic binary re-sync for %s: %v", resource, err)
	}
	if err := g.SweepInodes(); err != nil {
		log.Errorf("daemon: periodic inode sweep for %s: %v", resource, err)
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
		d.rollbackReload(resources, guards, nil)
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
		d.rollbackReload(resources, guards, nil)
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
	if err := d.prepareNewResources(resources, guards, oldByPath, &unlocked); err != nil {
		d.rollbackReload(resources, guards, unlocked)
		return err
	}

	// Phase 2 — populate the new guards' inode maps (hooks live, plaintext
	// readable — no protection gap), then resolve deferred entries and
	// re-sync: same fail-closed ordering as startup.
	if err := d.prepareGuards(resources, guards); err != nil {
		d.rollbackReload(resources, guards, unlocked)
		return err
	}

	// Phase 3 — start new guards' ringbuf readers; old guards stay attached.
	if err := d.startNewGuards(guards); err != nil {
		d.rollbackReload(resources, guards, unlocked)
		return err
	}

	// Phase 4 — commit: swap references, retire old forwarders, detach old LSM programs.
	d.commitReload(resources, guards)
	log.Info("daemon: configuration reloaded without dropping protection")
	return nil
}

// prepareNewResources validates and unlocks resources new to the configuration; kept
// resources are untouched. Each unlock is recorded so rollback can lock it back.
// resources/guards are index-aligned (Reload enforces this).
func (d *daemonUseCase) prepareNewResources(resources []daemonconfig.Resource, guards []repository.GuardRepository, oldByPath map[string]int, unlocked *[]string) error {
	for i := range resources {
		if _, existed := oldByPath[resources[i].Path]; existed {
			continue
		}
		if err := d.prepareAddedResource(&resources[i], guards[i], unlocked); err != nil {
			return err
		}
	}
	return nil
}

func (d *daemonUseCase) prepareAddedResource(r *daemonconfig.Resource, g repository.GuardRepository, unlocked *[]string) error {
	// Grouped watch paths share the group's encryption root: the vault-level
	// checks and the unlock target the root, and the rollback re-locks roots.
	root := encryptionRootOf(r)
	encrypted, err := d.vault.IsEncrypted(root)
	if err != nil {
		return fmt.Errorf("reload: checking encryption on %s: %w", root, err)
	}
	if r.NeedEncryption && !encrypted {
		return fmt.Errorf("reload: directory %s is NOT encrypted: run the fscrypt migration first or set need_encryption: false", root)
	}
	if !r.NeedEncryption && encrypted {
		log.Warnf("daemon: reload: resource %s is encrypted but need_encryption: false \u2014 leaving it locked", root)
		return nil
	}
	provisioned, err := d.vault.IsProvisioned(root)
	if err != nil {
		return fmt.Errorf("reload: checking lock state of %s: %w", root, err)
	}
	if provisioned {
		return nil
	}
	if err := vaultOpForGuard(g, root, func() error { return d.vault.Unlock(root) }); err != nil {
		return fmt.Errorf("reload: unlocking %s: %w", root, err)
	}
	*unlocked = append(*unlocked, root)
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

// rollbackReload aborts a reload without ever leaving a resource unlocked
// and unguarded: while the not-yet-committed guards are still attached
// (denying every non-whitelisted access), the freshly unlocked resources
// are locked back with bounded retries, and the guards detach only once
// every vault is keyless — the same ordering discipline as Stop. If a pin
// outlasts the retry budget, the new guards are kept attached as orphans
// (registered for the daemon's Stop) so the resource remains guarded;
// an unlocked-and-unguarded state is unreachable by construction.
//
// Called only from Reload while holding d.mu, so the orphan registration
// needs no extra locking. The retry budget is bounded — unlike Stop —
// because rollback runs on the SIGHUP handler, which must keep serving
// signals.
func (d *daemonUseCase) rollbackReload(resources []daemonconfig.Resource, guards []repository.GuardRepository, unlocked []string) {
	if len(unlocked) > 0 {
		pending := append([]string(nil), unlocked...)
		for round := 0; round < maxRollbackLockRounds && len(pending) > 0; round++ {
			var still []string
			for _, path := range pending {
				if err := d.lockWithRetry(resources, guards, path); err != nil {
					log.Errorf("daemon: rollback: %s is still unlocked: %v", path, err)
					still = append(still, path)
				}
			}
			pending = still
			if len(pending) > 0 {
				time.Sleep(rollbackLockRetryDelay)
			}
		}
		if len(pending) > 0 {
			log.Errorf("daemon: rollback gave up locking %v \u2014 the new guards STAY ATTACHED and deny access until the daemon stops", pending)
			d.orphans = append(d.orphans, guards...)
			return
		}
	}
	for _, g := range guards {
		g.Stop()
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
		// Grouped watch paths share one encryption root: lock each unique
		// vault once (the lock is vault-wide by fscrypt semantics).
		roots := uniqueEncryptionRoots(resources)

		// First pass: remove the key where possible.
		for _, root := range roots {
			if err := vaultOpOnRoot(resources, guards, root, func() error { return d.vault.Lock(root, false) }); err != nil && !errors.Is(err, repository.ErrKeyMissing) {
				log.Warnf("daemon: first-pass lock of %s: %v", root, err)
			}
		}

		// Second pass: force-flush with EBUSY retries; ENOKEY (ErrKeyMissing) counts
		// as success. Busy resources retry forever — detaching guards here would leave
		// the tree unlocked and unguarded, so shutdown never gives up.
		d.lockUntilAllKeyless(resources, guards, roots)

		for _, g := range guards {
			g.Stop()
		}
		// Guards orphaned by an aborted reload (their resource could not
		// be locked back in time) kept denying access throughout this
		// lockdown; retire them last, after every configured vault is
		// keyless.
		for _, g := range d.orphans {
			g.Stop()
		}
		d.orphans = nil
		log.Info("daemon: shutdown complete, all vaults locked")
	})
}

// lockUntilAllKeyless retries the force-flush lock until every resource is keyless; it
// never returns while one is unlocked — guards stay attached denying all access, so the
// tree is never unprotected and shutdown just waits for file pins to close.
func (d *daemonUseCase) lockUntilAllKeyless(resources []daemonconfig.Resource, guards []repository.GuardRepository, pending []string) {
	for len(pending) > 0 {
		var still []string
		for _, path := range pending {
			if err := d.lockWithRetry(resources, guards, path); err != nil {
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

func (d *daemonUseCase) lockWithRetry(resources []daemonconfig.Resource, guards []repository.GuardRepository, path string) error {
	for i := 0; i < maxLockRetries; i++ {
		err := vaultOpOnRoot(resources, guards, path, func() error { return d.vault.Lock(path, true) })
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

// GrantEditAccess widens the guard for resourcePath to admit interactive
// edits by the app-listener binary (uid 0 only), for the duration of an
// authenticated live edit-protected session. Only one grant is active at a
// time. The returned revoke restores the read-only baseline and clears the
// active flag; it is safe to call more than once and safe to call after a
// reload swapped the guard out (the new guard is built read-only anyway).
func (d *daemonUseCase) GrantEditAccess(resourcePath string) (func() error, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.stopping {
		return nil, errors.New("daemon: shutdown in progress")
	}
	if d.editGrantActive {
		return nil, errors.New("another live edit session is already active")
	}

	idx := -1
	for i := range d.resources {
		if d.resources[i].Path == resourcePath {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, fmt.Errorf("%s is not a guarded resource in the running configuration", resourcePath)
	}

	g := d.guards[idx]
	if err := g.GrantSelfEditAccess(); err != nil {
		return nil, err
	}
	d.editGrantActive = true

	var once sync.Once
	revoke := func() error {
		var rerr error
		once.Do(func() {
			rerr = g.RevokeSelfEditAccess()
			d.mu.Lock()
			d.editGrantActive = false
			d.mu.Unlock()
		})
		return rerr
	}
	return revoke, nil
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
