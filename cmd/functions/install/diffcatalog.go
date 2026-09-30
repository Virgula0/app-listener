package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// runDiffCatalog is the incremental counterpart of the full wizard: compares the catalog with the
// installed daemon.conf and offers to protect and encrypt every discoverable critical directory not
// already guarded. Existing sections stay byte-for-byte intact; selected sections are appended and
// the normal pipeline (verify encryption -> ask -> migrate -> deploy -> (re)start) runs on the
// merged config. Requires a previous installation; interactive.
func runDiffCatalog() error {
	mergedText, mergedCfg, proceed, err := collectDiffAdditions()
	if err != nil || !proceed {
		return err
	}
	return applyDiffAdditions(mergedText, mergedCfg)
}

// collectDiffAdditions reads the installed config, discovers catalog directories missing from it,
// runs the picker and editor, and returns the merged text + parsed config. proceed=false (nil
// error) when there's nothing to add or nothing selected.
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

// applyDiffAdditions runs the encrypt + deploy half on the merged config. The daemon is stopped for
// the whole cycle (a new section may need an in-place fscrypt migration) and restored on the
// existing config if the user aborts before deploy.
func applyDiffAdditions(mergedText string, mergedCfg *daemonconfig.Config) error {
	wasActive := systemd.IsDaemonActive()
	if stopErr := systemd.StopDaemonIfRunning(); stopErr != nil {
		return stopErr
	}

	if keyErr := ensureMasterKey(); keyErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, keyErr)
	}

	sshUsers, sshErr := askSSHAgentUsers(mergedCfg)
	if sshErr == nil {
		sshErr = addKeysToAgent(sshUsers)
	}
	if sshErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, sshErr)
	}

	bunUsers, bunErr := askBunTmpdirUsers(mergedCfg)
	if bunErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, bunErr)
	}

	vault := fscrypt.New()
	securedText, secErr := secureResources(vault, mergedText, mergedCfg)
	if secErr != nil {
		return restoreDaemonAfterDiffAbort(wasActive, secErr)
	}

	// deploy writes the config only after ConfirmOverwrite, runs the eBPF preflight, and (re)starts
	// the daemon; its failure modes already leave a documented state, so they surface as-is.
	// --diff-catalog never touches the edit-protected password (an existing hash stays and is
	// self-guarded again by the restarted daemon).
	if deployErr := deploy(securedText, sshUsers, bunUsers, ""); deployErr != nil {
		return deployErr
	}
	if cleanErr := cleanOrphanedFscrypt(mergedCfg); cleanErr != nil {
		return cleanErr
	}
	return cleanupBackups(mergedCfg)
}

// configPaths returns every path a config resource occupies: the watch path
// and its encryption root (equal unless the section is a grouped watch).
func configPaths(cfg *daemonconfig.Config) []string {
	out := make([]string, 0, len(cfg.Resources)*2)
	for i := range cfg.Resources {
		out = append(out, cfg.Resources[i].Path, cfg.Resources[i].EncryptionRootOrPath())
	}
	return out
}

// uncoveredCandidates returns catalog candidates whose path is neither a configured watch
// path/encryption root nor nested with one (either direction): what a diff would propose adding.
func uncoveredCandidates(cfg *daemonconfig.Config, users []inst.User) []inst.Candidate {
	covered := configPaths(cfg)
	all := inst.DiscoverForUsers(users)
	var fresh []inst.Candidate
	for i := range all {
		if !pathCovered(all[i].Path, covered) {
			fresh = append(fresh, all[i])
		}
	}
	return fresh
}

// pathCovered: p is already protected by a config path: equal, inside, or a parent of it (a new
// parent watch would overlap an existing guarded tree).
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

// appendSectionsAndEdit appends the generated sections verbatim to the existing config and opens
// the editor, so the diff the user reviews (here and at ConfirmOverwrite) is exactly the new
// sections.
func appendSectionsAndEdit(oldText string, picked []inst.Candidate) (string, *daemonconfig.Config, error) {
	merged := inst.InsertSectionsBeforeLibraries(oldText, inst.GenerateSections(sectionsFromCandidates(picked)))
	for _, block := range libraryBlocksFromCandidates(picked) {
		merged = inst.EnsureLibraryBlock(merged, &block)
	}
	return runConfigEditor(
		"app-listener daemon.conf — new sections appended, review and save (Ctrl+S)",
		merged)
}

// restoreDaemonAfterDiffAbort restarts the daemon on the existing config when the run is aborted
// after the daemon was stopped, then returns the original cause.
func restoreDaemonAfterDiffAbort(wasActive bool, cause error) error {
	if wasActive {
		log.Warn("catalog diff aborted — restarting the daemon on the existing config ...")
		if e := deliverReload(false); e != nil {
			log.Errorf("could not restart the daemon (%v) — run: sudo systemctl start %s", e, systemd.DaemonServiceName)
		}
	}
	return cause
}
