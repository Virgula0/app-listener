package install

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/backups"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// restoreBackups undoes the fscrypt migration: the found .app_listener.backup dirs are listed (all
// preselected) and, after one confirmation, the encrypted copies are deleted and backups moved
// back. Aborts while the daemon runs: it would keep using the directories being deleted and leave
// daemon.conf declaring need_encryption: true for a resource with no policy (issue #53).
// RequireDaemonStopped (not the bare systemd check) also catches a manually started daemon.
func restoreBackups() error {
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
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

// deletePostBackups deletes every found .app_listener.backup: listed (all preselected), one
// confirmation, then removed under a TUI progress bar. Aborts while the daemon runs, like
// restoreBackups: a backup is the only recovery for a resource the daemon may still rely on (e.g. a
// crashed unlock), so never delete it under a live daemon.
func deletePostBackups() error {
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
	}
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
