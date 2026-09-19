package daemon

import (
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// buildTrustedSet computes the daemon-wide trusted sets from the config:
//   - binaries: every whitelisted binary (TRUSTED_BINARY);
//   - libs: each binary's statically resolved library closure, the operator's
//     allow_lib entries, and /etc/ld.so.preload (TRUSTED_LIB);
//   - dirs: every guarded resource root, so a library loaded from inside one
//     of those write-protected trees is trusted without an allow_lib entry.
//
// A currently-unresolvable binary/library is skipped; the periodic re-sync
// picks it up later. It is only ever a coverage gap, never a protection gap.
func buildTrustedSet(cfg *daemonconfig.Config) (binaries, libs, dirs []string) {
	binSet := make(map[string]struct{})
	libSet := make(map[string]struct{})
	dirSet := make(map[string]struct{})

	for i := range cfg.Resources {
		dirSet[cfg.Resources[i].Path] = struct{}{}
		for _, b := range cfg.Resources[i].Binaries {
			binSet[b.Path] = struct{}{}
		}
		for _, l := range cfg.Resources[i].AllowLibs {
			libSet[l] = struct{}{}
		}
	}

	// Static dependency closure of every whitelisted binary (the auto part the
	// operator never has to list).
	for b := range binSet {
		closure, err := ebpf.ResolveLibraryClosure(b)
		if err != nil {
			log.Warnf("trust guard: could not resolve library closure of %s: %v", b, err)
			continue
		}
		for _, l := range closure {
			libSet[l] = struct{}{}
		}
	}

	// Global force-preloads: everything in /etc/ld.so.preload is loaded into
	// every dynamically linked process and must be trusted.
	if preloads, err := ebpf.LdSoPreloadPaths(); err != nil {
		log.Warnf("trust guard: reading /etc/ld.so.preload: %v", err)
	} else {
		for _, p := range preloads {
			libSet[p] = struct{}{}
		}
	}

	return keys(binSet), keys(libSet), keys(dirSet)
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// startTrustGuard brings up the daemon-wide trust guard: binary
// write-protection (#1) and the library-load allowlist (#2), both always
// enforced. It returns a cleanup func (a no-op when nothing started) so the
// caller defers a single statement. Best-effort — any failure logs and starts
// nothing, never blocking daemon startup; per-resource enforcement is
// unaffected regardless.
func startTrustGuard(cfg *daemonconfig.Config) func() {
	noop := func() {}
	binaries, libs, dirs := buildTrustedSet(cfg)
	if len(binaries) == 0 {
		return noop // no whitelisted binaries: nothing to protect
	}
	tg, err := guard.NewTrustGuard()
	if err != nil {
		log.Warnf("trust guard: not started (%v) — binary write-protection and library enforcement "+
			"unavailable; per-resource enforcement is unaffected", err)
		return noop
	}
	if err := tg.SetGuardedDirs(dirs); err != nil {
		log.Warnf("trust guard: could not record guarded roots (%v) — not started", err)
		tg.Stop()
		return noop
	}
	if err := tg.SetTrusted(binaries, libs); err != nil {
		log.Warnf("trust guard: could not load the trusted set (%v) — not started", err)
		tg.Stop()
		return noop
	}
	if err := tg.Start(); err != nil {
		log.Warnf("trust guard: could not attach (%v) — not started", err)
		tg.Stop()
		return noop
	}
	log.Infof("trust guard: enforcing write-protection + library allowlist for %d whitelisted binary(ies)", len(binaries))
	return tg.Stop
}
