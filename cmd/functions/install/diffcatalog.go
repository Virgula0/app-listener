package install

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// runDiffCatalog is the incremental counterpart of the full wizard: it
// compares the catalog against the installed daemon.conf and offers to add —
// protect and encrypt — every discoverable critical directory that is not
// already guarded. Existing sections are left byte-for-byte intact; the
// selected new sections are appended and the normal install pipeline
// (verify encryption state → ask encryption → migrate → deploy → (re)start)
// runs on the merged config. Requires a previous installation; interactive.
func runDiffCatalog() error {
	mergedText, mergedCfg, proceed, err := collectDiffAdditions()
	if err != nil || !proceed {
		return err
	}
	return applyDiffAdditions(mergedText, mergedCfg)
}

// collectDiffAdditions reads the installed config, discovers the catalog
// directories missing from it, runs the picker and the editor, and returns
// the merged config text + parsed config. proceed is false (with nil error)
// when there is nothing to add or the user selected nothing.
func collectDiffAdditions() (mergedText string, mergedCfg *daemonconfig.Config, proceed bool, err error) {
	oldBytes, readErr := os.ReadFile(systemd.SystemConfigPath)
	if readErr != nil {
		return "", nil, false, fmt.Errorf("no previous installation found — run `sudo app-listener install` first: %w", readErr)
	}
	cfg, loadErr := daemonconfig.Load(systemd.SystemConfigPath)
	if loadErr != nil {
		return "", nil, false, fmt.Errorf("parsing the installed configuration: %w", loadErr)
	}

	users, usersErr := inst.ListUsers()
	if usersErr != nil {
		return "", nil, false, fmt.Errorf("listing users: %w", usersErr)
	}

	fresh := uncoveredCandidates(cfg, users)
	if len(fresh) == 0 {
		log.Info("catalog diff: every discoverable critical directory is already in daemon.conf — nothing to add")
		return "", nil, false, nil
	}
	log.Infof("catalog diff: %d critical path(s) discovered that are not in daemon.conf", len(fresh))

	picked, pickErr := pickFromCandidates(fresh,
		"New critical directories found — select the ones to add and protect",
		"These exist on the host but are not in the current daemon.conf. All are preselected. "+
			"Selected paths are appended to the config and encrypted like a fresh install.")
	if pickErr != nil {
		return "", nil, false, pickErr
	}
	if len(picked) == 0 {
		log.Info("catalog diff: nothing selected — daemon.conf unchanged")
		return "", nil, false, nil
	}

	text, parsed, editErr := appendSectionsAndEdit(string(oldBytes), picked)
	if editErr != nil {
		return "", nil, false, editErr
	}
	return text, parsed, true, nil
}

// applyDiffAdditions runs the encrypt + deploy half of the install pipeline
// on the merged config. The daemon is stopped for the whole cycle (a new
// section may need an in-place fscrypt migration) and brought back on the
// existing config if the user aborts before deploy.
func applyDiffAdditions(mergedText string, mergedCfg *daemonconfig.Config) error {
	wasActive := systemd.IsDaemonActive()
	if stopErr := systemd.StopDaemonIfRunning(); stopErr != nil {
		return stopErr
	}

	if keyErr := ensureMasterKey(); keyErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, keyErr)
	}

	vault := fscrypt.New()
	securedText, secErr := secureResources(vault, mergedText, mergedCfg)
	if secErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, secErr)
	}

	// deploy writes the config only after ConfirmOverwrite, runs the eBPF
	// preflight, and (re)starts the daemon — its own failure modes already
	// leave the box in a documented state, so surface them as-is.
	if deployErr := deploy(securedText); deployErr != nil {
		return deployErr
	}
	if cleanErr := cleanOrphanedFscrypt(mergedCfg); cleanErr != nil {
		return cleanErr
	}
	return cleanupBackups(mergedCfg)
}

// uncoveredCandidates returns the catalog candidates whose path is neither a
// configured watch path / encryption root nor nested with one (in either
// direction) — i.e. the directories a diff would propose adding.
func uncoveredCandidates(cfg *daemonconfig.Config, users []inst.User) []inst.Candidate {
	covered := make([]string, 0, len(cfg.Resources)*2)
	for i := range cfg.Resources {
		covered = append(covered, cfg.Resources[i].Path, cfg.Resources[i].EncryptionRootOrPath())
	}
	all := inst.DiscoverForUsers(users)
	var fresh []inst.Candidate
	for i := range all {
		if !pathCovered(all[i].Path, covered) {
			fresh = append(fresh, all[i])
		}
	}
	return fresh
}

// pathCovered reports whether p is already protected by an existing config
// path: equal to it, inside it, or a parent of it (a new parent watch would
// overlap an existing guarded tree).
func pathCovered(p string, covered []string) bool {
	for _, cv := range covered {
		if cv == "" {
			continue
		}
		if p == cv || isInsidePath(p, cv) || isInsidePath(cv, p) {
			return true
		}
	}
	return false
}

// appendSectionsAndEdit appends the generated sections for picked to the
// existing config text verbatim and opens the editor on the result, so the
// diff the user reviews (here and again at ConfirmOverwrite) is exactly the
// new sections.
func appendSectionsAndEdit(oldText string, picked []inst.Candidate) (string, *daemonconfig.Config, error) {
	merged := strings.TrimRight(oldText, "\n") + "\n" +
		inst.GenerateSections(sectionsFromCandidates(picked)) + "\n"
	return runConfigEditor(
		"app-listener daemon.conf — new sections appended, review and save (Ctrl+S)",
		merged)
}

// restoreDaemonAfterDiffAbort brings the daemon back on the existing config
// when a diff-catalog run is aborted after the daemon was stopped, then
// returns the original cause.
func restoreDaemonAfterDiffAbort(wasActive bool, cause error) error {
	if wasActive {
		log.Warn("catalog diff aborted — restarting the daemon on the existing config ...")
		if e := deliverReload(false); e != nil {
			log.Errorf("could not restart the daemon (%v) — run: sudo systemctl start %s", e, systemd.DaemonServiceName)
		}
	}
	return cause
}
