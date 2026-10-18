package install

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/backups"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// restoreBackups undoes the fscrypt migration: the found .app_listener.backup dirs are listed (all
// preselected) and, after a confirmation naming them, the encrypted copies are deleted and backups
// moved back. Aborts while the daemon runs: it would keep using the directories being deleted and leave
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
	ok, err := wizard.Confirm(wizard.Question{
		Title: fmt.Sprintf("Restore %d selected backup(s)?", len(entries)),
		Description: "Each encrypted directory is DELETED and its backup moved back in its place:\n" +
			listBackups(entries, func(b backups.Backup) string { return b.BackupPath + " -> " + b.Path }),
		Affirmative: "Restore",
		Negative:    "Cancel",
	})
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

// deletePostBackups deletes every found .app_listener.backup: listed (all preselected), a
// confirmation naming them, then removed under a TUI progress bar. The daemon may keep running: a
// backup lies outside every guarded tree and the daemon never reads it.
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
	ok, err := wizard.Confirm(wizard.Question{
		Title: fmt.Sprintf("Delete %d selected backup(s) permanently?", len(entries)),
		Description: "These unencrypted copies are removed for good; the encrypted directories keep your data:\n" +
			listBackups(entries, func(b backups.Backup) string { return b.BackupPath }),
		Affirmative: "Delete",
		Negative:    "Cancel",
	})
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

func listBackups(entries []backups.Backup, line func(backups.Backup) string) string {
	lines := make([]string, len(entries))
	for i := range entries {
		lines[i] = "  " + line(entries[i])
	}
	return strings.Join(lines, "\n")
}
