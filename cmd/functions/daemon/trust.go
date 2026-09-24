package daemon

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

// buildTrustedSet computes the daemon-wide trusted sets from the config:
//   - binaries: every whitelisted binary (TRUSTED_BINARY);
//   - libs: each binary's static library closure, the allow_lib entries and /etc/ld.so.preload
//     (TRUSTED_LIB);
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

	// Static dependency closure of every whitelisted binary (the auto part the
	// operator never has to list).
	rejected = addLibraryClosures(binSet, libSet)

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

// addLibraryClosures adds every binary's auto-trustable library closure to libSet and returns the
// members it refused.
func addLibraryClosures(binSet, libSet map[string]struct{}) map[string]*libRejection {
	rejected := make(map[string]*libRejection)
	for b := range binSet {
		closure, refused, err := ebpf.LibraryClosure(b)
		if err != nil {
			log.Warnf("trust guard: could not resolve library closure of %s: %v", b, err)
			continue
		}
		for _, l := range closure {
			libSet[l] = struct{}{}
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
// whitelist changes (the pacman/apt catalog-refresh hooks reload rather than restart), so a binary
// added by a reload must be re-applied to guard_trusted_files or it gets no library allowlist while
// still holding access to the secrets.
type trustManager struct {
	tg *guard.TrustGuard
}

// startTrustGuard brings up the daemon-wide trust guard: binary write-protection (#1), the
// library-load allowlist (#2) and reserved glob names (#3), all always enforced. Best-effort: a failure logs and starts
// nothing, never blocking startup or per-resource enforcement.
func startTrustGuard(cfg *daemonconfig.Config) *trustManager {
	binaries, libs, dirs, rejected := buildTrustedSet(cfg)
	if len(binaries) == 0 {
		return &trustManager{} // no whitelisted binaries: nothing to protect
	}
	tg, err := guard.NewTrustGuard()
	if err != nil {
		log.Warnf("trust guard: not started (%v) — binary write-protection and library enforcement "+
			"unavailable; per-resource enforcement is unaffected", err)
		return &trustManager{}
	}
	if err := applyTrustSet(tg, cfg, binaries, libs, dirs, rejected); err != nil {
		log.Warnf("trust guard: %v — not started", err)
		tg.Stop()
		return &trustManager{}
	}
	if err := tg.Start(); err != nil {
		log.Warnf("trust guard: could not attach (%v) — not started", err)
		tg.Stop()
		return &trustManager{}
	}
	guard.SetReplacementCheck(tg.AllowReplacement)
	log.Infof("trust guard: enforcing write-protection + library allowlist for %d whitelisted binary(ies)", len(binaries))
	return &trustManager{tg: tg}
}

// reload re-applies the trusted set from the new config to the running trust guard (SetTrusted /
// SetGuardedDirs clear-then-fill, so it is a true replace). No-op if the trust guard never started.
func (m *trustManager) reload(cfg *daemonconfig.Config) {
	if m == nil || m.tg == nil {
		return
	}
	binaries, libs, dirs, rejected := buildTrustedSet(cfg)
	if err := applyTrustSet(m.tg, cfg, binaries, libs, dirs, rejected); err != nil {
		log.Warnf("trust guard: reload could not re-apply the trusted set (%v) — keeping the previous set", err)
		return
	}
	log.Infof("trust guard: trusted set rebuilt after reload (%d whitelisted binary(ies))", len(binaries))
}

func (m *trustManager) stop() {
	if m != nil && m.tg != nil {
		guard.SetReplacementCheck(nil)
		m.tg.Stop()
	}
}

func applyTrustSet(tg *guard.TrustGuard, cfg *daemonconfig.Config, binaries, libs []string,
	dirs []guard.TrustedDir, rejected map[string]*libRejection) error {
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
	reservations := buildGlobReservations(cfg, users)
	if err := tg.SetGlobReservations(reservations); err != nil {
		return fmt.Errorf("reserving glob names: %w", err)
	}
	warnUntrustedLibs(rejected, dirs, reservations)
	return nil
}
