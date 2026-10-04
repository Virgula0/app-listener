package daemon

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

// buildTrustedSet computes the daemon-wide trusted sets from the config:
//   - binaries: every whitelisted binary (TRUSTED_BINARY);
//   - libs: the allow_lib entries and /etc/ld.so.preload (TRUSTED_LIB). Never the static library
//     closure: its members that pass rootOwnedSafe are system-trusted by the kernel anyway, and
//     SetTrusted re-stats the unresolved path, so a symlink in a user-writable RPATH dir re-pointed
//     after the check would make a user-owned inode TRUSTED_LIB;
//   - dirs: every guarded resource root, so libraries inside those write-protected trees are
//     trusted without allow_lib (a read-only lib_dir only for its own writers).
//
// A currently-unresolvable binary/library is skipped and picked up by the periodic re-sync: a
// coverage gap, never a protection gap. rejected maps each closure library that is not safe to
// auto-trust to why and to the binaries loading it (warnUntrustedLibs).
func buildTrustedSet(cfg *daemonconfig.Config) (binaries, libs []string, dirs []guard.TrustedDir,
	rejected map[string]*libRejection) {
	binSet := make(map[string]struct{})
	libSet := make(map[string]struct{})

	for i := range cfg.Resources {
		dirs = append(dirs, trustedDirOf(&cfg.Resources[i]))
		for _, b := range cfg.Resources[i].Binaries {
			binSet[b.Path] = struct{}{}
		}
		// PendingBinaries: whitelisted binaries unreadable at parse time (a still-locked vault — every
		// binary living inside an encrypted watch tree lands here). The guard admits them post-unlock
		// via ResolvePendingBinaries; without them here they get NO library allowlist (trust_mmap
		// bails unless the exe is TRUSTED_BINARY) while still holding full access to the secrets, so
		// LD_PRELOAD works against them. SetTrusted re-stats every path after unlock.
		for _, b := range cfg.Resources[i].PendingBinaries {
			binSet[b.Path] = struct{}{}
		}
		for _, l := range cfg.Resources[i].AllowLibs {
			libSet[l] = struct{}{}
		}
		// Unreadable at parse time (typically a still-locked vault). The trust set is built after
		// unlock and SetTrusted re-stats every path, so these resolve now or are skipped with a
		// warning (previously they were parked and never re-read, leaving such libraries
		// untrusted).
		for _, l := range cfg.Resources[i].PendingLibs {
			libSet[l] = struct{}{}
		}
	}
	// [libraries] blocks: library trust is daemon-wide, not per resource.
	for _, l := range cfg.SharedAllowLibs {
		libSet[l] = struct{}{}
	}
	// An inspector reads every protected process's metadata: no preloaded code may ride it.
	for _, p := range inspectorPaths(cfg, systemPlaced) {
		binSet[p] = struct{}{}
	}

	rejected = closureRejections(binSet)

	// Global force-preloads: everything in /etc/ld.so.preload is loaded into
	// every dynamically linked process and must be trusted.
	if preloads, err := ebpf.LdSoPreloadPaths(); err != nil {
		log.Warnf("trust guard: reading /etc/ld.so.preload: %v", err)
	} else {
		for _, p := range preloads {
			libSet[p] = struct{}{}
		}
	}

	return keys(binSet), keys(libSet), dirs, rejected
}

// trustedDirOf is res's root for the library allowlist: a read-only lib_dir is loadable only by
// its writers, since its contents may predate the guard.
func trustedDirOf(res *daemonconfig.Resource) guard.TrustedDir {
	d := guard.TrustedDir{Path: res.Path}
	if res.ReadOnly {
		d.Loaders = []string{}
		for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
			for _, b := range list {
				d.Loaders = append(d.Loaders, b.Path)
			}
		}
	}
	return d
}

// closureRejections maps each closure library the kernel won't system-trust to why and to the
// binaries loading it. Only for warnUntrustedLibs: no closure member enters the trusted set.
func closureRejections(binSet map[string]struct{}) map[string]*libRejection {
	rejected := make(map[string]*libRejection)
	for b := range binSet {
		_, refused, err := libraryClosure(b)
		if err != nil {
			log.Warnf("trust guard: could not resolve library closure of %s: %v", b, err)
			continue
		}
		for l, why := range refused {
			if rejected[l] == nil {
				rejected[l] = &libRejection{why: why}
			}
			rejected[l].bins = append(rejected[l].bins, b)
		}
	}
	return rejected
}

var libraryClosure = ebpf.LibraryClosure

var systemPlaced = func(path string) bool {
	f, err := ebpf.OpenSystemPlaced(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

type libRejection struct {
	why  error
	bins []string
}

// warnUntrustedLibs reports closure libraries trust_mmap will refuse: not root-owned, and neither
// inside a guarded tree nor a reserved name loaded by one of its writers. One line per library.
func warnUntrustedLibs(rejected map[string]*libRejection, dirs []guard.TrustedDir, r guard.GlobReservations) {
	libs := make([]string, 0, len(rejected))
	for l := range rejected {
		libs = append(libs, l)
	}
	sort.Strings(libs)
	for _, l := range libs {
		rej := rejected[l]
		if slices.ContainsFunc(rej.bins, func(b string) bool { return inGuardedTree(l, b, dirs) }) ||
			slices.ContainsFunc(rej.bins, func(b string) bool { return r.LibTrusted(b, l) }) {
			continue
		}
		log.Warnf("library closure: %s will be refused (%v; not in a guarded tree, not a reserved "+
			"library of its app) — loaded by %d whitelisted binary(ies); add it via allow_lib or a "+
			"lib_dir if a load is denied", l, rej.why, len(rej.bins))
	}
}

// inGuardedTree mirrors under_guarded_tree: the innermost root above path admits bin.
func inGuardedTree(path, bin string, dirs []guard.TrustedDir) bool {
	var inner *guard.TrustedDir
	for i := range dirs {
		d := &dirs[i]
		if (path == d.Path || strings.HasPrefix(path, d.Path+"/")) && (inner == nil || len(d.Path) > len(inner.Path)) {
			inner = d
		}
	}
	return inner != nil && (inner.Loaders == nil || slices.Contains(inner.Loaders, bin))
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// trustManager owns the daemon-wide trust guard across its lifetime, including SIGHUP reloads. The
// trusted set is built once at start AND rebuilt on every reload: SIGHUP is the normal way the
// whitelist changes (the in-daemon catalog refresh reloads rather than restarts), so a binary
// added by a reload must be re-applied to guard_trusted_files or it gets no library allowlist while
// still holding access to the secrets.
type trustManager struct {
	tg       *guard.TrustGuard
	vet      *binaryVetter
	binaries int
	bunRoots map[inodeID]bool // Bun tmp dirs the last apply reserved (reserveGlobs)
	// inspectMu guards inspectCfg, the config whose [inspectors] were last applied; nil until the
	// trust guard is started (AdmitInspector needs its vouched superblocks).
	inspectMu  sync.Mutex
	inspectCfg *daemonconfig.Config
}

// startTrustGuard loads and attaches the daemon-wide trust guard — binary write-protection (#1), the
// library-load allowlist (#2) and reserved glob names (#3) — with every path resolvable now. Called
// before any vault is unlocked, and the daemon refuses to start if it fails: the per-resource guards
// judge a process by its exe inode alone, so without this guard code injected into a whitelisted
// process (a preloaded library) or a whitelisted binary rewritten in place reads the secrets. Started
// even with no whitelisted binary, so a reload that adds one has a guard to update.
func startTrustGuard(cfg *daemonconfig.Config) (m *trustManager, err error) {
	vet, err := openBinaryVetter(cfg)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			vet.close()
		}
	}()
	if fault := injectedTrustFault("start"); fault != nil {
		return nil, trustStartupError(fault)
	}
	tg, err := guard.NewTrustGuard()
	if err != nil {
		return nil, trustStartupError(err)
	}
	m = &trustManager{tg: tg, vet: vet}
	// Before the first apply: until the guards' build approves a binary, it updates nothing and
	// binds no reserved name.
	tg.SetBinaryResolver(vet.resolve)
	vet.setTrust(tg)
	if err := m.apply(cfg); err != nil {
		tg.Stop()
		return nil, trustStartupError(err)
	}
	if err := tg.Start(); err != nil {
		tg.Stop()
		return nil, trustStartupError(err)
	}
	guard.SetReplacementCheck(func(path string, f *os.File, old, newKey guard.GuardInodeKey) bool {
		if !tg.AllowReplacement(path, f, old, newKey) {
			return false
		}
		vet.recordLive(path, f, newKey)
		return true
	})
	return m, nil
}

// afterUnlock re-applies the trusted set once vaults are unlocked: binaries and libraries inside
// them only resolve now.
func (m *trustManager) afterUnlock(cfg *daemonconfig.Config) error {
	if err := injectedTrustFault("after-unlock"); err != nil {
		return trustStartupError(err)
	}
	if err := m.apply(cfg); err != nil {
		return trustStartupError(err)
	}
	log.Infof("trust guard: enforcing write-protection + library allowlist for %d whitelisted binary(ies)",
		m.binaries)
	install.WarnGeneralTools(cfg)
	if err := m.applyInspectors(cfg); err != nil {
		return trustStartupError(err)
	}
	if err := m.tg.SetElectronApps(cfg.ElectronApps); err != nil {
		return trustStartupError(err)
	}
	m.vet.endBootstrap()
	return nil
}

// applyInspectors grants cfg's [inspectors] after the trusted set made them TRUSTED_BINARY.
func (m *trustManager) applyInspectors(cfg *daemonconfig.Config) error {
	m.inspectMu.Lock()
	defer m.inspectMu.Unlock()
	if err := guard.SetInspectors(resolveInspectors(cfg, m.tg.AdmitInspector), cfg.Inspectors); err != nil {
		return fmt.Errorf("granting inspectors: %w", err)
	}
	m.inspectCfg = cfg
	return nil
}

// resyncInspectors re-resolves the inspectors (a package upgrade replaced one), on the catalog
// watch's re-sync.
func (m *trustManager) resyncInspectors() {
	m.inspectMu.Lock()
	cfg := m.inspectCfg
	m.inspectMu.Unlock()
	if cfg == nil {
		return
	}
	if err := m.applyInspectors(cfg); err != nil {
		log.Errorf("daemon: re-syncing inspectors: %v", err)
	}
}

func trustStartupError(err error) error {
	return fmt.Errorf("%w: trust guard unavailable (%w) — refusing to start: without it code injected into a "+
		"whitelisted process could read the guarded secrets", constants.ErrCriticalStartup, err)
}

// reload re-applies the trusted set from the new config (SetTrusted / SetGuardedDirs clear-then-fill,
// so it is a true replace). On failure the previous set stays.
func (m *trustManager) reload(cfg *daemonconfig.Config) {
	if err := m.tg.SyncMounts(); err != nil {
		log.Errorf("trust guard: CRITICAL: mount re-sync on reload failed (%v) — system libraries are "+
			"refused to whitelisted binaries until it succeeds", err)
	}
	if err := m.apply(cfg); err != nil {
		log.Errorf("trust guard: CRITICAL: reload could not re-apply the trusted set (%v) — keeping the previous "+
			"set; binaries added by this reload have no library allowlist until this is fixed", err)
		return
	}
	log.Infof("trust guard: trusted set rebuilt after reload (%d whitelisted binary(ies))", m.binaries)
	install.WarnGeneralTools(cfg)
	if err := m.applyInspectors(cfg); err != nil {
		log.Errorf("daemon: CRITICAL: reload could not apply [inspectors] (%v) — inspector grants may be "+
			"stale until this is fixed", err)
	}
	if err := m.tg.SetElectronApps(cfg.ElectronApps); err != nil {
		log.Errorf("daemon: CRITICAL: reload could not apply [electron_apps] (%v) — a removed app may still "+
			"be admitted until this is fixed", err)
	}
}

func (m *trustManager) apply(cfg *daemonconfig.Config) error {
	binaries, libs, dirs, rejected := buildTrustedSet(cfg)
	if err := m.applyTrustSet(cfg, binaries, libs, dirs, rejected); err != nil {
		return err
	}
	m.binaries = len(binaries)
	return nil
}

func (m *trustManager) stop() {
	if m != nil && m.tg != nil {
		guard.SetReplacementCheck(nil)
		m.tg.Stop()
	}
	if m != nil && m.vet != nil {
		m.vet.close()
	}
}

func (m *trustManager) applyTrustSet(cfg *daemonconfig.Config, binaries, libs []string,
	dirs []guard.TrustedDir, rejected map[string]*libRejection) error {
	tg := m.tg
	if err := tg.SetGuardedDirs(dirs); err != nil {
		return fmt.Errorf("recording guarded roots: %w", err)
	}
	if err := tg.SetTrusted(binaries, libs); err != nil {
		return fmt.Errorf("loading the trusted set: %w", err)
	}
	if err := tg.SetUpdaters(planUpdaters(cfg)); err != nil {
		return fmt.Errorf("scoping binary updaters: %w", err)
	}
	users, err := install.ListUsers()
	if err != nil {
		return fmt.Errorf("listing users for reserved glob names: %w", err)
	}
	reservations, bun := buildGlobReservations(cfg, users)
	if err := m.reserveGlobs(reservations, bun); err != nil {
		return fmt.Errorf("reserving glob names: %w", err)
	}
	warnUntrustedLibs(rejected, dirs, reservations)
	return nil
}
