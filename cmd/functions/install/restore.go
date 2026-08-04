package install

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// daemonRunning reports whether the app-listener-daemon systemd unit is
// active (empty when systemctl is unavailable).
func daemonRunning() bool {
	return strings.TrimSpace(systemctlOutput("is-active", daemonServiceName)) == daemonActiveState
}

// backupEntry is one discovered migration backup and the directory it
// belongs to.
type backupEntry struct {
	path   string
	backup string
}

// restoreBackups undoes the fscrypt migration: the found
// .app_listener.backup directories are shown in a TUI list (all
// preselected) and, after a single confirmation, the encrypted copies are
// deleted and the backups moved back to the original locations. It aborts
// when the daemon is running — restoring while the daemon is active would
// let it keep unlocking and using the very directories being deleted. A
// TUI progress bar shows the operation progress.
func restoreBackups() error {
	if daemonRunning() {
		return errors.New("fatal: the daemon is running — stop it before restoring backups: systemctl stop app-listener-daemon")
	}
	entries, err := findBackups()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no migration backups found: nothing to restore")
		return nil
	}
	entries, err = pickBackups(entries,
		"Migration backups found — select the ones to restore",
		"All are preselected. Restoring DELETES the encrypted directory and moves the backup back to the original location.")
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no backups selected: nothing to restore")
		return nil
	}
	ok, err := confirmOnce(fmt.Sprintf("Restore %d selected backup(s)? The encrypted directories will be DELETED.", len(entries)), "Restore all")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("restore canceled")
		return nil
	}
	vault := fscrypt.New()
	if err := runWithProgress("Restoring", entries, func(path string) error {
		if err := vault.RestoreBackup(path); err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		return nil
	}); err != nil {
		return err
	}
	log.Infof("restored %d backup(s)", len(entries))
	return nil
}

// deletePostBackups deletes every found .app_listener.backup: the backups
// are shown in a TUI list (all preselected) and, after a single
// confirmation, removed with a TUI progress bar showing the progress.
func deletePostBackups() error {
	entries, err := findBackups()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no migration backups found: nothing to delete")
		return nil
	}
	entries, err = pickBackups(entries,
		"Migration backups found — select the ones to delete",
		"All are preselected. The backups are plain, unencrypted copies of the migrated directories.")
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no backups selected: nothing to delete")
		return nil
	}
	ok, err := confirmOnce(fmt.Sprintf("Delete %d selected backup(s) permanently?", len(entries)), "Delete all")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("delete canceled")
		return nil
	}
	if err := runWithProgress("Deleting", entries, func(path string) error {
		backup := path + fscrypt.BackupSuffix
		if err := os.RemoveAll(backup); err != nil {
			return fmt.Errorf("removing backup %s: %w", backup, err)
		}
		return nil
	}); err != nil {
		return err
	}
	log.Infof("deleted %d backup(s)", len(entries))
	return nil
}

// findBackups lists every catalog/config path that still carries a
// .app_listener.backup.
func findBackups() ([]backupEntry, error) {
	candidates, err := restoreCandidates()
	if err != nil {
		return nil, err
	}
	var out []backupEntry
	for _, path := range candidates {
		backup := path + fscrypt.BackupSuffix
		if _, err := os.Lstat(backup); err != nil {
			continue
		}
		out = append(out, backupEntry{path: path, backup: backup})
	}
	return out, nil
}

// pickBackups shows the found backups in a TUI multi-select (all
// preselected) and returns the selected ones.
func pickBackups(entries []backupEntry, title, description string) ([]backupEntry, error) {
	opts := make([]huh.Option[int], 0, len(entries))
	for i := range entries {
		opts = append(opts, huh.NewOption(entries[i].backup, i).Selected(true))
	}
	var pickedIdx []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title(title).
			Description(description).
			Options(opts...).
			Height(10).
			Value(&pickedIdx),
	))
	if err := form.WithKeyMap(selectionKeymap()).Run(); err != nil {
		return nil, err
	}
	picked := make([]backupEntry, 0, len(pickedIdx))
	for _, i := range pickedIdx {
		picked = append(picked, entries[i])
	}
	return picked, nil
}

// confirmOnce asks a single yes/no for the whole batch (default: no).
func confirmOnce(question, affirmative string) (bool, error) {
	answer := false
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(question).
			Affirmative(affirmative).
			Negative("Cancel").
			Value(&answer),
	)).Run(); err != nil {
		return false, err
	}
	return answer, nil
}

// runWithProgress executes op for every entry while a single-line progress
// bar at the bottom of the terminal shows the operation progress (one step
// per directory); the fscrypt logs scroll normally above it.
func runWithProgress(verb string, entries []backupEntry, op func(path string) error) error {
	total := len(entries)
	return withBottomBar(func(bar *bottomBar) error {
		for i := range entries {
			label := fmt.Sprintf("%s %d/%d: %s", verb, i+1, total, entries[i].backup)
			bar.set(label, float64(i)/float64(total))
			if err := op(entries[i].path); err != nil {
				return err
			}
			bar.set(label+" (done)", float64(i+1)/float64(total))
		}
		return nil
	})
}

// restoreCandidates lists every path that could carry a migration backup:
// the resources of /etc/app-listener/daemon.conf when it exists (covering
// manually added directories) merged with every catalog path discovered
// for all local users plus the system-level entries. Deduplicated by path.
func restoreCandidates() ([]string, error) {
	seen := make(map[string]bool)
	var out []string
	add := func(path string) {
		if !seen[path] {
			seen[path] = true
			out = append(out, path)
		}
	}

	if data, err := os.ReadFile(systemConfigPath); err == nil {
		cfg, loadErr := validateConfigText(string(data))
		if loadErr != nil {
			return nil, fmt.Errorf("reading %s: %w", systemConfigPath, loadErr)
		}
		for _, r := range cfg.Resources {
			add(r.Path)
		}
	}

	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	found := inst.DiscoverForUsers(users)
	for i := range found {
		add(found[i].Path)
	}
	return out, nil
}
