package install

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/backups"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// restoreBackups undoes the fscrypt migration: the found
// .app_listener.backup directories are shown in a TUI list (all
// preselected) and, after a single confirmation, the encrypted copies are
// deleted and the backups moved back to the original locations. It aborts
// when the daemon is running — restoring while the daemon is active would
// let it keep unlocking and using the very directories being deleted.
func restoreBackups() error {
	if protected.DaemonRunning() {
		return errors.New("fatal: the daemon is running — stop it before restoring backups: systemctl stop app-listener-daemon")
	}
	entries, err := backups.Find()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no migration backups found: nothing to restore")
		return nil
	}
	entries, err = backups.Select(entries,
		"Migration backups found — select the ones to restore",
		"All are preselected. Restoring DELETES the encrypted directory and moves the backup back to the original location.")
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no backups selected: nothing to restore")
		return nil
	}
	ok, err := wizard.ConfirmOnce(fmt.Sprintf("Restore %d selected backup(s)? The encrypted directories will be DELETED.", len(entries)), "Restore all")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("restore canceled")
		return nil
	}
	if err := backups.Restore(entries); err != nil {
		return err
	}
	log.Infof("restored %d backup(s)", len(entries))
	return nil
}

// deletePostBackups deletes every found .app_listener.backup: the backups
// are shown in a TUI list (all preselected) and, after a single
// confirmation, removed with a TUI progress bar showing the progress.
func deletePostBackups() error {
	entries, err := backups.Find()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no migration backups found: nothing to delete")
		return nil
	}
	entries, err = backups.Select(entries,
		"Migration backups found — select the ones to delete",
		"All are preselected. The backups are plain, unencrypted copies of the migrated directories.")
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("no backups selected: nothing to delete")
		return nil
	}
	ok, err := wizard.ConfirmOnce(fmt.Sprintf("Delete %d selected backup(s) permanently?", len(entries)), "Delete all")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("delete canceled")
		return nil
	}
	if err := backups.Delete(entries); err != nil {
		return err
	}
	log.Infof("deleted %d backup(s)", len(entries))
	return nil
}
