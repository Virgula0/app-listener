package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// runDiffCatalog is the incremental counterpart of the full wizard: compares the catalog with the
// installed daemon.conf and offers to protect and encrypt every discoverable critical directory not
// already guarded. Existing sections stay byte-for-byte intact; selected sections are appended and
// the normal pipeline (verify encryption -> ask -> migrate -> deploy) runs on the merged config.
// It needs the daemon stopped (a new section may need an in-place fscrypt migration) and asks at
// the end whether to start it. Requires a previous installation; interactive.
func runDiffCatalog() error {
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
	}
	prevCfg, mergedText, mergedCfg, proceed, err := collectDiffAdditions()
	if err == nil && proceed {
		err = applyDiffAdditions(prevCfg, mergedText, mergedCfg)
	}
	if err != nil {
		return err
	}
	return offerDaemonStart()
}

// offerDaemonStart asks to start (and enable) a stopped daemon on the installed config.
func offerDaemonStart() error {
	if systemd.IsDaemonActive() || !systemd.DaemonUnitInstalled() {
		return nil
	}
	start, err := wizard.ConfirmOnce("The daemon is stopped: its folders stay locked until it runs. Start it now?",
		"Start the daemon")
	if err != nil {
		return err
	}
	if !start {
		log.Infof("daemon left stopped — start it with: sudo systemctl start %s", systemd.DaemonServiceName)
		return nil
	}
	return systemd.EnableAndVerify(false)
}

// collectDiffAdditions reads the installed config, discovers catalog directories missing from it,
// runs the picker and editor, and returns the installed config plus the merged text + parsed config.
// proceed=false (nil error) when there's nothing to add or nothing selected.
func collectDiffAdditions() (prevCfg *daemonconfig.Config, mergedText string, mergedCfg *daemonconfig.Config,
	proceed bool, err error) {
	oldBytes, readErr := os.ReadFile(systemd.SystemConfigPath)
	if readErr != nil {
		return nil, "", nil, false, fmt.Errorf("no previous installation found — run `sudo app-listener install` first: %w", readErr)
	}
	cfg, loadErr := daemonconfig.Load(systemd.SystemConfigPath)
	if loadErr != nil {
		return nil, "", nil, false, fmt.Errorf("parsing the installed configuration: %w", loadErr)
	}

	users, usersErr := inst.ListUsers()
	if usersErr != nil {
		return nil, "", nil, false, fmt.Errorf("listing users: %w", usersErr)
	}

	fresh := uncoveredCandidates(cfg, users)
	if len(fresh) == 0 {
		log.Info("catalog diff: every discoverable critical directory is already in daemon.conf — nothing to add")
		return nil, "", nil, false, nil
	}
	log.Infof("catalog diff: %d critical path(s) discovered that are not in daemon.conf", len(fresh))

	picked, pickErr := pickFromCandidates(fresh,
		"New critical directories found — select the ones to add and protect",
		"These exist on the host but are not in the current daemon.conf. All are preselected. "+
			"Selected paths are appended to the config and encrypted like a fresh install.")
	if pickErr != nil {
		return nil, "", nil, false, pickErr
	}
	if len(picked) == 0 {
		log.Info("catalog diff: nothing selected — daemon.conf unchanged")
		return nil, "", nil, false, nil
	}

	text, parsed, editErr := appendSectionsAndEdit(string(oldBytes), picked)
	if editErr != nil {
		return nil, "", nil, false, editErr
	}
	inst.WarnGeneralTools(parsed)
	return cfg, text, parsed, true, nil
}

// applyDiffAdditions runs the encrypt + deploy half on the merged config, with the daemon stopped
// (runDiffCatalog): an abort before deploy leaves daemon.conf untouched.
func applyDiffAdditions(prevCfg *daemonconfig.Config, mergedText string, mergedCfg *daemonconfig.Config) error {
	if err := ensureMasterKey(); err != nil {
		return err
	}
	sshUsers, err := askSSHAgentUsers(addedResources(prevCfg, mergedCfg))
	if err == nil {
		err = addKeysToAgent(sshUsers)
	}
	if err != nil {
		return err
	}
	bunUsers, err := askBunTmpdirUsers(mergedCfg)
	if err != nil {
		return err
	}
	securedText, err := secureResources(fscrypt.New(), mergedText, mergedCfg)
	if err != nil {
		return err
	}
	// deployFiles writes the config only after ConfirmOverwrite and runs the eBPF preflight; it never
	// touches the edit-protected password (an existing hash stays, sealed again once the daemon runs).
	if _, err := deployFiles(securedText, sshUsers, bunUsers, ""); err != nil {
		return err
	}
	if err := cleanOrphanedFscrypt(mergedCfg); err != nil {
		return err
	}
	return cleanupBackups(mergedCfg)
}

// addedResources returns merged's resources prev doesn't cover. The ssh-agent setup must only see
// a ~/.ssh this run guards for the first time: an already-guarded one is encrypted and, with the
// daemon stopped, locked, so editing its ~/.ssh/config fails with ENOKEY and aborts the run.
func addedResources(prev, merged *daemonconfig.Config) *daemonconfig.Config {
	covered := inst.CoveredPaths(prev)
	out := &daemonconfig.Config{}
	for i := range merged.Resources {
		if !inst.PathCovered(merged.Resources[i].Path, covered) {
			out.Resources = append(out.Resources, merged.Resources[i])
		}
	}
	return out
}

// uncoveredCandidates returns catalog candidates whose path is neither a configured watch
// path/encryption root nor nested with one (either direction): what a diff would propose adding.
func uncoveredCandidates(cfg *daemonconfig.Config, users []inst.User) []inst.Candidate {
	covered := inst.CoveredPaths(cfg)
	all := inst.DiscoverForUsers(users)
	var fresh []inst.Candidate
	for i := range all {
		if !inst.PathCovered(all[i].Path, covered) {
			fresh = append(fresh, all[i])
		}
	}
	return fresh
}

// appendSectionsAndEdit appends the generated sections verbatim to the existing config and opens
// the editor, so the diff the user reviews (here and at ConfirmOverwrite) is exactly the new
// sections.
func appendSectionsAndEdit(oldText string, picked []inst.Candidate) (string, *daemonconfig.Config, error) {
	merged := inst.InsertSectionsBeforeLibraries(oldText, inst.GenerateSections(sectionsFromCandidates(picked)))
	for _, block := range libraryBlocksFromCandidates(picked) {
		merged = inst.EnsureLibraryBlock(merged, &block)
	}
	prev, err := daemonconfig.Parse([]byte(oldText))
	if err != nil {
		prev = nil // every line is checked
	}
	return runConfigEditor(
		"app-listener daemon.conf — new sections appended, review and save (Ctrl+S)",
		merged, prev)
}
