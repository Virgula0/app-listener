package usecase

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/constants"
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
	// inodeGCEvery is the guard_inodes eviction cadence (ReconcileInodes): a
	// full, complete tree walk, much coarser than resyncSweepEvery — it
	// exists to bound how long a stale (dev, ino) entry for a long-deleted
	// file can linger and risk colliding with an unrelated inode reused
	// elsewhere on the same filesystem, not to catch changes quickly.
	inodeGCEvery = time.Hour

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

// concurrencyLimit bounds a resource-indexed fan-out to core count, and never
// higher than there is work to do — avoids a thundering herd of concurrent
// tree walks / BPF syscalls on hosts with many configured resources.
func concurrencyLimit(n int) int {
	if n < 1 {
		return 1
	}
	if lim := runtime.NumCPU(); lim < n {
		return lim
	}
	return n
}

// unlockRoots unlocks every root in roots via unlockOne, fanned out with
// bounded concurrency. Safe because: (1) roots is always already deduplicated
// by the caller (uniqueEncryptionRoots) — a shared `watch:` group root never
// appears twice, so no two goroutines ever race to unlock the SAME vault;
// (2) every guard for every resource — including every member of a
// shared-root group — is attached before startGuards ever calls this (guards
// are built by buildGuards, which returns before Start() runs), so no root
// unlocked here can expose a sibling resource whose own guard is not yet
// live. Each Vault.Unlock call only touches its own root's kernel keyring
// state (or, for a file-vault target, decrypts that one file); the only
// shared mutable state in the Vault (the collateral map) is already
// mutex-protected for concurrent access. On any failure every already-issued
// unlock still runs to completion (no early cancellation) — they are all for
// resources whose guards are already attached and denying, so letting them
// finish is not a protection gap, just possibly-unneeded work that the
// caller's error path locks back down.
func unlockRoots(roots []string, unlockOne func(root string) error) error {
	if len(roots) == 0 {
		return nil
	}
	errs := make([]error, len(roots))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrencyLimit(len(roots)))
	for i, root := range roots {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, root string) {
			defer wg.Done()
			defer func() { <-sem }()
			errs[i] = unlockOne(root)
		}(i, root)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("%w: unlocking %s: %w", constants.ErrCriticalStartup, roots[i], err)
		}
	}
	return nil
}

// prepareOneGuard runs one guard's post-unlock startup: populate its inode
// map, resolve whitelist entries deferred while its resource was locked,
// re-sync replaced binaries, then start its ringbuf reader. Every guard is
// already attached and every vault is already unlocked before this is called
// (startGuards' two barriers above it), and each guard owns its own inode map
// and BPF objects — no shared mutable state across resources — so
// prepareAndStartGuards runs this concurrently across resources.
//
// The one behavior this trades away versus a fully sequential pass: a
// whitelist entry naming a binary that lives inside a DIFFERENT resource's
// tree may not resolve on this call if that other resource's populate hasn't
// landed yet. It stays deferred — still fail-closed, never falls open — and
// is picked up by the periodic re-sync (forwardEvents/periodicSweep) shortly
// after, same as a binary replaced in place after startup already is.
func (d *daemonUseCase) prepareOneGuard(i int) error {
	g := d.guards[i]
	if err := g.PopulateInodes(); err != nil {
		// Critical: every vault is already unlocked and every guard attached
		// here, so a populate failure is a property of the tree on disk (an
		// unreadable subtree, a tree the walk cannot classify) and reproduces
		// identically on every start. Untagged, it exited 1 and systemd's
		// Restart=on-failure crash-looped the daemon — unlocking and
		// re-locking every vault every two seconds.
		return fmt.Errorf("%w: populating guard for %s: %w", constants.ErrCriticalStartup, d.resources[i].Path, err)
	}
	if err := g.ResolvePendingBinaries(); err != nil {
		return fmt.Errorf("resolving deferred binaries for %s: %w", d.resources[i].Path, err)
	}
	if _, err := g.ReSyncBinaries(); err != nil {
		return fmt.Errorf("re-syncing binary whitelist for %s: %w", d.resources[i].Path, err)
	}
	if err := g.Start(); err != nil {
		return fmt.Errorf("starting guard for %s: %w", d.resources[i].Path, err)
	}
	return nil
}

// prepareAndStartGuards fans prepareOneGuard out across every guard, bounded
// by concurrencyLimit, and waits for all of them before returning.
func (d *daemonUseCase) prepareAndStartGuards() error {
	n := len(d.guards)
	errs := make([]error, n)
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrencyLimit(n))
	for i := range n {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			errs[i] = d.prepareOneGuard(i)
		}(i)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
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
			return fmt.Errorf("%w: directory %s is NOT encrypted: run the fscrypt migration first or set need_encryption: false", constants.ErrCriticalStartup, r.Path)
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
	// unlocked-and-unprotected window for either. Within each pass, the
	// per-root unlocks are independent of one another (see unlockRoots) and
	// run concurrently. A file root's guard self access is widened only for
	// the call itself (vaultOpOnRoot -> guard.WithSelfVaultAccess): the
	// in-place unlock reads+rewrites the file's own bytes, which the guard's
	// baseline self mask (open/read/stat only) does not permit.
	// An Unlock failure here means the master key/policy pairing itself is
	// wrong (missing key file, wrong key, unsupported filesystem) — the same
	// key, same policy and same filesystem will fail identically on every
	// restart, so this is the same "administrator action required" class as
	// verifyEncryptionStates' mismatch, not a transient condition.
	dirRoots, fileRoots := partitionEncryptionRoots(uniqueEncryptionRoots(d.resources))
	if err := unlockRoots(dirRoots, d.vault.Unlock); err != nil {
		return err
	}
	if err := unlockRoots(fileRoots, func(root string) error {
		return vaultOpOnRoot(d.resources, d.guards, root, func() error { return d.vault.Unlock(root) })
	}); err != nil {
		return err
	}

	// Every vault is now unlocked and every guard is already attached
	// (built before Start): for each resource, populate its inode map, resolve
	// deferred whitelist entries, re-sync replaced binaries and start its
	// ringbuf reader — fail-closed contract preserved (attach → unlock →
	// populate → resolve), just fanned out across resources instead of run
	// one at a time (see prepareOneGuard for why this is safe to parallelize).
	if err := d.prepareAndStartGuards(); err != nil {
		return err
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
	inodeGC := time.NewTicker(inodeGCEvery)
	defer inodeGC.Stop()
	var lastResync time.Time
	for {
		select {
		case ev, ok := <-g.Events():
			if !ok {
				return
			}
			if !d.dispatchGuardEvent(resource, g, &ev, stop, &lastResync) {
				return
			}
		case <-sweep.C:
			d.periodicSweep(resource, g)
			lastResync = time.Now()
		case <-inodeGC.C:
			if err := g.ReconcileInodes(); err != nil {
				log.Warnf("daemon: periodic inode GC for %s: %v", resource, err)
			}
		case <-d.done:
			return
		case <-stop:
			return
		}
	}
}

// dispatchGuardEvent re-syncs the binary whitelist on a throttled denial and
// forwards ev to the daemon's event channel, reporting whether the forward
// loop should keep running (false means the daemon or this resource's guard
// is shutting down).
func (d *daemonUseCase) dispatchGuardEvent(resource string, g repository.GuardRepository, ev *guard.GuardEvent, stop <-chan struct{}, lastResync *time.Time) bool {
	// The raw block-device gate names a device, not this resource: label it
	// as such and skip the re-sync (it is never an in-place binary
	// replacement).
	label := resource
	if ev.RawDevice {
		label = guard.RawDeviceResourceLabel
	}
	// A denial usually means the binary was replaced in place; re-sync
	// (throttled by resyncMinInterval) to admit the new inode instead of
	// re-statting the whitelist per event.
	if ev.Blocked && !ev.RawDevice && time.Since(*lastResync) >= resyncMinInterval {
		if _, err := g.ReSyncBinaries(); err != nil {
			log.Errorf("daemon: re-syncing binary whitelist for %s: %v", resource, err)
		}
		*lastResync = time.Now()
	}
	select {
	case d.events <- DaemonEvent{Resource: label, Event: *ev}:
		return true
	case <-d.done:
		return false
	case <-stop:
		return false
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
	for i := range resources {
		newByPath[resources[i].Path] = true
	}
	var removed []string
	for i := range d.resources {
		if !newByPath[d.resources[i].Path] {
			removed = append(removed, d.resources[i].Path)
		}
	}
	if len(removed) > 0 {
		d.rollbackReload(resources, guards, nil)
		return fmt.Errorf("reload refused: %v would be left unlocked and unprotected — stop the daemon, edit the config, and start it again to drop them", removed)
	}

	oldByPath := make(map[string]int, len(d.resources))
	for i := range d.resources {
		oldByPath[d.resources[i].Path] = i
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
	// Carry the memory-read taint from each old guard to the new guard for
	// the same resource BEFORE the old ones are stopped. The new guards were
	// built with empty guard_tainted_pids maps, so without this a reload
	// would silently drop process_vm_readv / ptrace protection for every
	// process that had already read a guarded file (they keep running across
	// the reload with the secret still in memory).
	d.carryTaintAcrossReload(resources, guards)

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

// carryTaintAcrossReload copies each old guard's tainted-pid set into the new
// guard for the same resource path. Best-effort: a transfer failure only
// weakens memory-read protection for already-tainted processes until they next
// touch the guarded tree, never enforcement of any direct access.
func (d *daemonUseCase) carryTaintAcrossReload(resources []daemonconfig.Resource, guards []repository.GuardRepository) {
	oldByPath := make(map[string]repository.GuardRepository, len(d.guards))
	for i := range d.resources {
		oldByPath[d.resources[i].Path] = d.guards[i]
	}
	for i := range resources {
		r := &resources[i]
		old, ok := oldByPath[r.Path]
		if !ok {
			continue // resource new to this config: nothing to carry
		}
		pids, err := old.SnapshotTaintedPIDs()
		if err != nil {
			log.Warnf("daemon: reading tainted pids for %s during reload: %v", r.Path, err)
			continue
		}
		if len(pids) == 0 {
			continue
		}
		if err := guards[i].RestoreTaintedPIDs(pids); err != nil {
			log.Warnf("daemon: carrying taint across reload for %s: %v", r.Path, err)
		}
	}
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

// lockUntilAllKeyless retries the force-flush lock until every resource is keyless or
// permanently un-lockable; it never returns while one is genuinely still unlocked — guards
// stay attached denying all access, so the tree is never unprotected and shutdown just
// waits for file pins to close. repository.ErrNotEncrypted is the one exception: it means
// path carries no fscrypt policy at all (a permanent condition — see lockWithRetry), so
// retrying it can never succeed, unlike a busy key. Giving up on it here is safe ONLY
// because the caller (Stop) detaches every guard together right after this function
// returns, once every OTHER root is genuinely keyless — there is no per-resource window
// where this one is unguarded while others remain protected. rollbackReload, which detaches
// guards per-resource, must not reuse this leniency (see lockWithRetry).
func (d *daemonUseCase) lockUntilAllKeyless(resources []daemonconfig.Resource, guards []repository.GuardRepository, pending []string) {
	for len(pending) > 0 {
		var still []string
		for _, path := range pending {
			err := d.lockWithRetry(resources, guards, path)
			if err == nil {
				continue
			}
			if errors.Is(err, repository.ErrNotEncrypted) {
				log.Errorf("daemon: %s has no fscrypt policy — it was decrypted, restored from a backup, or otherwise modified outside the daemon while the config still expects it encrypted; nothing left to lock, but this resource needs investigation before it is trusted again: %v", path, err)
				continue
			}
			log.Errorf("daemon: %s is still unlocked: a process holds open files in it (investigate with: lsof +D %s, fuser -v %s): %v",
				path, path, path, err)
			still = append(still, path)
		}
		if len(still) == 0 {
			break
		}
		log.Errorf("daemon: shutdown blocked on %d resource(s) still unlocked — guards stay attached and deny all access until they lock", len(still))
		pending = still
		time.Sleep(time.Second)
	}
}

// lockWithRetry retries Lock until it succeeds, the key is confirmed already gone
// (ErrKeyMissing), or the retry budget for a busy key is exhausted. repository.ErrNotEncrypted
// (path carries no fscrypt policy at all — e.g. a backup restore replaced the encrypted tree
// with plaintext) is returned as-is rather than swallowed into success: unlike ErrKeyMissing,
// there was never a key here to confirm gone, and the two callers must react to it
// differently. lockUntilAllKeyless (Stop, full shutdown) may give up on this one path and
// still finish, because every guard detaches together regardless. rollbackReload (a live
// SIGHUP mid-run) must NOT — a newly-unlocked resource whose policy vanished before rollback
// could lock it back has no old guard to fall back on, so treating this as "resolved" there
// would detach its only remaining protection over an unencrypted resource; returning the
// error keeps it in rollbackReload's pending set so the bounded retry budget's exhaustion
// orphans (keeps attached) rather than stops that guard.
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
