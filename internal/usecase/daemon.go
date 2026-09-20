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
	// maxLockRetries bounds the EBUSY retry loop when force-flushing an fscrypt lock.
	maxLockRetries = 100
	lockRetryDelay = 10 * time.Millisecond
	// resyncMinInterval throttles post-denial re-syncs so a denial storm doesn't re-stat every
	// binary per event.
	resyncMinInterval = 2 * time.Second
	// resyncSweepEvery is the background sweep interval, catching binary replacements that never
	// produced a denial and recreated single-file roots. Fingerprint-gated (SweepInodes), so cheap
	// when nothing changed.
	resyncSweepEvery = 30 * time.Second
	// inodeGCEvery is the ReconcileInodes cadence (full tree walk), much coarser than
	// resyncSweepEvery: it only bounds how long a stale (dev, ino) entry can linger and collide
	// with a reused inode.
	inodeGCEvery = time.Hour
	// fileRootFollowEvery is how quickly a SINGLE-FILE watch root replaced by an atomic save (temp
	// + rename over; Steam's registry.vdf) is re-anchored to its new inode. Until then the new file
	// is unguarded. The kernel side can't follow the rename (guard_path_rename has no verifier
	// budget left, issue #45); for a file root the check is one stat, so it runs far more often
	// than the general sweep.
	fileRootFollowEvery = time.Second

	// Rollback lock-back budget is bounded (unlike Stop's infinite wait): rollback runs on the
	// SIGHUP handler, which must keep serving signals.
	maxRollbackLockRounds  = 3
	rollbackLockRetryDelay = 500 * time.Millisecond
)

// DaemonEvent is one guard event tagged with the resource it belongs to.
type DaemonEvent struct {
	Resource string
	Event    guard.GuardEvent
}

// DaemonUseCase orchestrates the daemon lifecycle: verify encryption, unlock, run guards, atomic
// SIGHUP reloads, secure lockdown. Order is attach → unlock → populate → resolve → re-sync, so
// resources are never readable without live protection.
type DaemonUseCase interface {
	Start() error
	Reload(resources []daemonconfig.Resource, guards []repository.GuardRepository) error
	Stop()
	Events() <-chan DaemonEvent
	Resources() []daemonconfig.Resource
	// GrantEditAccess widens one resource's guard so the app-listener binary (uid 0 only) may
	// modify it for an authenticated live edit-protected session. One grant at a time. The returned
	// revoke restores the read-only baseline; idempotent.
	GrantEditAccess(resourcePath string) (revoke func() error, err error)
}

type daemonUseCase struct {
	mu        sync.RWMutex
	resources []daemonconfig.Resource
	vault     repository.Vault
	guards    []repository.GuardRepository
	// stops: one close-channel per guard to retire its event forwarder on reload (guard event
	// channels are never closed).
	stops []chan struct{}

	events chan DaemonEvent
	done   chan struct{}

	start   sync.Once
	stop    sync.Once
	started bool
	// stopping is set once lockdown begins; Reload refuses when set, so a SIGHUP mid-shutdown never
	// attaches or unlocks.
	stopping bool
	// orphans: not-yet-committed guards from an aborted reload whose freshly unlocked resources
	// couldn't be locked back within the rollback budget. They stay attached (denying) and are
	// retired by Stop after its lockdown (rollbackReload).
	orphans []repository.GuardRepository
	// editGrantActive is set while a live edit-protected session holds a write grant
	// (GrantEditAccess); one at a time.
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

// concurrencyLimit bounds a resource-indexed fan-out to core count and never more than the work
// (avoids a thundering herd of tree walks/BPF syscalls).
func concurrencyLimit(n int) int {
	if n < 1 {
		return 1
	}
	if lim := runtime.NumCPU(); lim < n {
		return lim
	}
	return n
}

// unlockRoots unlocks every root via unlockOne with bounded concurrency. Safe because: (1) roots is
// already deduplicated (uniqueEncryptionRoots), so no two goroutines unlock the SAME vault; (2)
// every guard of every resource, including all members of a shared-root group, is attached before
// this runs (buildGuards returns before Start()), so no unlock exposes a sibling whose guard isn't
// live. Each Vault.Unlock touches only its own root's keyring state (or decrypts that one file);
// the collateral map is mutex-protected. On failure already-issued unlocks still finish (their
// guards are attached and denying, so this is no protection gap); the caller's error path locks
// them back.
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

// prepareOneGuard runs one guard's post-unlock startup: populate its inode map, resolve entries
// deferred while locked, re-sync replaced binaries, start its ringbuf reader. All guards are
// attached and vaults unlocked before this (startGuards' two barriers) and each guard owns its
// maps, so prepareAndStartGuards runs it concurrently.
//
// Trade-off vs a sequential pass: an entry naming a binary inside a DIFFERENT resource's tree may
// not resolve if that resource's populate hasn't landed. It stays deferred (fail-closed) and is
// picked up by the periodic re-sync (forwardEvents/periodicSweep).
func (d *daemonUseCase) prepareOneGuard(i int) error {
	g := d.guards[i]
	if err := g.PopulateInodes(); err != nil {
		// Critical: vaults are unlocked and guards attached, so a populate failure is a property of
		// the tree on disk and reproduces on every start. Untagged, Restart=on-failure would
		// crash-loop the daemon, re-locking every vault every two seconds.
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

// prepareAndStartGuards fans prepareOneGuard out across guards (bounded by concurrencyLimit) and
// waits for all.
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

// encryptionRootOf returns the vault root governing r: the `watch:` group's section path, else r's
// own path.
func encryptionRootOf(r *daemonconfig.Resource) string {
	if r.EncryptionRoot != "" {
		return r.EncryptionRoot
	}
	return r.Path
}

// uniqueEncryptionRoots deduplicates encryption roots of encrypted resources (grouped watch paths
// share one vault); the key lifecycle must run once per root.
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

// partitionEncryptionRoots splits roots into fscrypt directories and file-vault regular files, so
// startGuards' unlock runs as two visible phases: directories (kernel-keyring cycle), then files
// (in-place same-inode cycle, filevault.go). A clarity/auditability split, not correctness: each
// Unlock is independently safe. A stat failure (e.g. root missing) falls into dirs so Unlock
// reports the real error.
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

// vaultOpForGuard runs a Vault.Unlock/Lock on root, widening g's self-access
// (guard.WithSelfVaultAccess) only when root is a file-vault (directory roots touch only the
// keyring, so run unwidened). g may be nil (no guard for root): op still runs so the real (denied)
// error surfaces instead of a silent skip.
func vaultOpForGuard(g repository.GuardRepository, root string, op func() error) error {
	info, statErr := os.Stat(root)
	if statErr != nil || !info.Mode().IsRegular() || g == nil {
		return op()
	}
	return g.WithSelfVaultAccess(op)
}

// vaultOpOnRoot is vaultOpForGuard, finding the guard by encryption root in resources/guards
// (index-aligned; Reload enforces this). A regular-file resource is never grouped, so exactly one
// resource/guard matches a file root.
func vaultOpOnRoot(resources []daemonconfig.Resource, guards []repository.GuardRepository, root string, op func() error) error {
	for i := range resources {
		if encryptionRootOf(&resources[i]) == root {
			return vaultOpForGuard(guards[i], root, op)
		}
	}
	return vaultOpForGuard(nil, root, op)
}

// verifyEncryptionStates checks each resource's encryption state against need_encryption before
// anything is unlocked or protected; a mismatch aborts the start.
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

	// Grouped watch paths share one encryption root: the key is provisioned once per vault.
	// Directories and file-vault files unlock in two passes (partitionEncryptionRoots); guards for
	// BOTH kinds are already attached, so there is no unlocked-and-unprotected window. Within a
	// pass, unlocks are independent and concurrent (unlockRoots). A file root's self access is
	// widened only for the call itself (vaultOpOnRoot -> guard.WithSelfVaultAccess), since the
	// in-place unlock rewrites the file's bytes, beyond the baseline open/read/stat mask.
	//
	// An Unlock failure means the key/policy pairing is wrong (missing key, wrong key, unsupported
	// fs) and fails identically on every restart: the same "administrator action required" class as
	// verifyEncryptionStates, not transient.
	dirRoots, fileRoots := partitionEncryptionRoots(uniqueEncryptionRoots(d.resources))
	if err := unlockRoots(dirRoots, d.vault.Unlock); err != nil {
		return err
	}
	if err := unlockRoots(fileRoots, func(root string) error {
		return vaultOpOnRoot(d.resources, d.guards, root, func() error { return d.vault.Unlock(root) })
	}); err != nil {
		return err
	}

	// Vaults unlocked, guards attached (built before Start): per resource, populate the inode map,
	// resolve deferred entries, re-sync replaced binaries, start the ringbuf reader. Fail-closed
	// order preserved, just fanned out (see prepareOneGuard).
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
	// nil (never fires) unless the resource is a single file; a watch root doesn't change type
	// while guarded.
	var followRoot <-chan time.Time
	if info, err := os.Lstat(resource); err == nil && !info.IsDir() {
		follow := time.NewTicker(fileRootFollowEvery)
		defer follow.Stop()
		followRoot = follow.C
	}
	var lastResync time.Time
	for {
		select {
		case <-followRoot:
			// Errors are left to the periodic sweep, which reports them at a sane rate.
			_ = g.SweepInodes()
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

// dispatchGuardEvent re-syncs the binary whitelist on a throttled denial and forwards ev to the
// event channel; false = the daemon or this guard is shutting down.
func (d *daemonUseCase) dispatchGuardEvent(resource string, g repository.GuardRepository, ev *guard.GuardEvent, stop <-chan struct{}, lastResync *time.Time) bool {
	// The raw block-device gate names a device, not this resource: label it so and skip the re-sync
	// (never an in-place binary replacement).
	label := resource
	if ev.RawDevice {
		label = guard.RawDeviceResourceLabel
	}
	// A denial usually means an in-place binary replacement: re-sync (throttled by
	// resyncMinInterval) to admit the new inode. A process-gate denial is not one.
	if ev.Blocked && !ev.RawDevice && ev.Process == "" && time.Since(*lastResync) >= resyncMinInterval {
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

// periodicSweep re-syncs the binary whitelist (catching replacements that never produced a denial)
// and refreshes the inode map. SweepInodes only works when the root's fingerprint moved (recreated
// single-file root, or a directory that gained/lost a top-level entry); an unconditional re-walk
// per tick was a large share of steady-state CPU.
func (d *daemonUseCase) periodicSweep(resource string, g repository.GuardRepository) {
	if _, err := g.ReSyncBinaries(); err != nil {
		log.Errorf("daemon: periodic binary re-sync for %s: %v", resource, err)
	}
	if err := g.SweepInodes(); err != nil {
		log.Errorf("daemon: periodic inode sweep for %s: %v", resource, err)
	}
}

// Reload atomically applies a new configuration: old and new LSM programs run concurrently and ANY
// deny wins, so protection is never weaker than either; on error the old config keeps running and
// new guards are detached. Removing a resource is refused (it would be left unlocked and
// unguarded); dropping protection requires stop/edit/start.
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

	// Phase 0: refuse to drop any protected resource (see Reload). Nothing is committed: new guards
	// detach, old config keeps running.
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

	// Phase 1: validate and unlock resources new to the config. Nothing is committed, so errors
	// roll back (detach new guards, re-lock). Guards are already attached, so unlocks always have
	// live protection.
	var unlocked []string
	if err := d.prepareNewResources(resources, guards, oldByPath, &unlocked); err != nil {
		d.rollbackReload(resources, guards, unlocked)
		return err
	}

	// Phase 2: populate the new guards' inode maps (hooks live, plaintext readable, no gap), then
	// resolve deferred entries and re-sync: same fail-closed order as startup.
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

// prepareNewResources validates and unlocks resources new to the config (kept ones untouched),
// recording each unlock so rollback can lock it back. resources/guards are index-aligned (Reload
// enforces this).
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
	// Grouped watch paths share the group's encryption root: vault-level checks, unlock and
	// rollback all target the root.
	root := encryptionRootOf(r)
	encrypted, err := d.vault.IsEncrypted(root)
	if err != nil {
		return fmt.Errorf("reload: checking encryption on %s: %w", root, err)
	}
	if r.NeedEncryption && !encrypted {
		return fmt.Errorf("reload: directory %s is NOT encrypted: run the fscrypt migration first or set need_encryption: false", root)
	}
	if !r.NeedEncryption {
		// Nothing to unlock. Stop here for the unencrypted case: falling through to IsProvisioned
		// asked fscrypt for the policy of a directory that has none and failed the WHOLE reload
		// (any lib_dir/need_encryption:false resource added by a catalog refresh), while a restart
		// accepted the same config.
		if encrypted {
			log.Warnf("daemon: reload: resource %s is encrypted but need_encryption: false \u2014 leaving it locked", root)
		}
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

// prepareGuards fills the new guards' inode maps, resolves deferred entries and re-syncs replaced
// binaries; runs post-unlock, so scans see plaintext while hooks are live.
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

// rollbackReload aborts a reload without leaving a resource unlocked and unguarded: with the
// uncommitted guards still attached (denying), the freshly unlocked resources are locked back with
// bounded retries, and guards detach only once every vault is keyless (same discipline as Stop). If
// a pin outlasts the budget, the new guards stay attached as orphans (retired by Stop), so
// unlocked-and-unguarded is unreachable.
//
// Called only from Reload holding d.mu, so orphan registration needs no extra locking. The budget
// is bounded (unlike Stop) because rollback runs on the SIGHUP handler, which must keep serving
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

// commitReload swaps in the new config; new forwarders spawn before old ones retire and old LSM
// programs detach, so the event stream never goes silent.
func (d *daemonUseCase) commitReload(resources []daemonconfig.Resource, guards []repository.GuardRepository) {
	// Carry the memory-read taint from each old guard to the new one for the same resource BEFORE
	// the old ones stop: new guards start with empty guard_tainted_pids, so otherwise a reload
	// would drop process_vm_readv/ptrace protection for processes already holding a secret in
	// memory.
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

// carryTaintAcrossReload copies each old guard's tainted-pid set into the new guard for the same
// resource path. Best-effort: a failure only weakens memory-read protection for already-tainted
// processes until they next touch the tree, never direct-access enforcement.
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

// Stop performs the secure lockdown: guards stay attached (denying all non-whitelisted access)
// until every resource is keyless; only then do LSM hooks detach. A resource whose key can't be
// deprovisioned blocks shutdown indefinitely (the daemon logs how to find the pinning process):
// never unlocked and unguarded.
func (d *daemonUseCase) Stop() {
	d.stop.Do(func() {
		log.Info("daemon: initiating secure lockdown")

		// Hold the write lock through lockdown: Reload serializes behind it and refuses once
		// stopping is set.
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
