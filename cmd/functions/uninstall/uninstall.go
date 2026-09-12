// The `app-listener uninstall` command reverts the installer: refuses while
// the daemon runs, re-scans the catalog (not daemon.conf) for encrypted dirs,
// verifies the master key, asks per dir to decrypt permanently (default: no),
// removes units/hook/binary/symlink/config, reverts sample-matching ssh-agent
// units after a separate confirmation, and deletes the key only with
// --delete-key.
package uninstall

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/backups"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// deleteKeyFlag removes the master key at the end; default false, since
// without it still-encrypted directories can never be unlocked again.
var deleteKeyFlag bool

func init() {
	UninstallCmd.Flags().BoolVar(&deleteKeyFlag, "delete-key", false,
		"Also delete the fscrypt master key at /etc/app-listener/fscrypt.key after the uninstall (kept by default: without it still-encrypted directories can never be unlocked again)")
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Interactive uninstaller: revert the daemon installation (fscrypt + systemd)",
	Long: `Interactive, root-only uninstaller for the daemon mode.

The wizard reverts everything the installer installed:

  0. fatally refuses while the daemon is running: an active systemd unit or
     a daemon process running outside systemd must be stopped manually
     first
  1. re-scans the catalog (the current daemon.conf is deliberately NOT
     used) and detects which directories are actually encrypted with fscrypt
  2. verifies the master key against every encrypted directory: a directory
     that does not unlock with the master key is a fatal error, because
     removing the daemon (and possibly the key with --delete-key) would
     lock it forever
  3. asks per directory whether its fscrypt encryption must be permanently
     removed (default: no) and decrypts the confirmed ones in place with a
     progress bar: an encrypted directory is only removed after the
     plaintext copy completed, so a failure leaves it untouched
  4. cleans the orphaned fscrypt metadata left behind by the decrypted
     directories (still-encrypted directories keep their metadata)
  5. reverts the installed systemd units (daemon + boot-time catalog
     refresh), the pacman/apt catalog-refresh hook, binary, PATH symlink
     and config, and — after a final confirmation (default: no) and only
     for units whose content matches the bundled sample — the per-user
     ssh-agent systemd units
  6. deletes the fscrypt master key ONLY when --delete-key is passed; the
     default keeps it, because without it every still-encrypted directory
     can never be unlocked again
  7. asks whether to delete the migration backups (.app_listener.backup):
     the found ones are listed (all preselected) and, after one
     confirmation, removed — they are plain unencrypted copies, so keeping
     them is harmless but usually pointless once the daemon is gone

To move a backup back over its directory instead of deleting it, use
'app-listener install --restore-backups'.

Aborting (Esc) any step cancels the uninstall; completed steps stay
completed.`,
	Args: cobra.NoArgs,
	RunE: runUninstall,
}

// runUninstall drives the whole uninstall flow.
func runUninstall(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("uninstall must be run as root: sudo app-listener uninstall")
	}

	// The daemon guards the very dirs being decrypted/deleted; refuse while
	// it lives instead of stopping it — nothing gets re-enabled afterwards.
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
	}

	vault := fscrypt.New()

	if err := decryptStep(vault); err != nil {
		return err
	}

	// Scan for migration backups while daemon.conf still exists (it names
	// manually added directories the catalog does not); the prompt runs at
	// the end, after revertSystemFiles.
	backupEntries, err := backups.Find()
	if err != nil {
		return err
	}

	if err := revertSSHAgents(); err != nil {
		return err
	}

	if err := revertSystemFiles(); err != nil {
		return err
	}

	if deleteKeyFlag {
		if err := removeMasterKey(); err != nil {
			return err
		}
	} else {
		log.Infof("keeping the fscrypt master key at %s (pass --delete-key to remove it)", fscrypt.MasterKeyFile)
	}

	if err := offerBackupCleanup(backupEntries); err != nil {
		return err
	}

	log.Info("uninstall complete")
	return nil
}

// offerBackupCleanup asks, at the end of the uninstall, whether to delete the
// .app_listener.backup directories the install-time fscrypt migration left
// behind. They are plain, unencrypted copies of the original directories —
// keeping them after an uninstall is usually pointless but harmless, so the
// user chooses (all preselected, one confirmation).
func offerBackupCleanup(entries []backups.Backup) error {
	if len(entries) == 0 {
		return nil
	}
	entries, err := backups.Select(entries,
		"Migration backups found — select the ones to delete",
		"Left by the install-time fscrypt migration: plain, unencrypted copies of the original directories. All are preselected; unselect any you want to keep.")
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		log.Info("keeping all migration backups")
		return nil
	}
	ok, err := wizard.ConfirmOnce(fmt.Sprintf("Delete %d selected backup(s) permanently?", len(entries)), "Delete all")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("keeping all migration backups")
		return nil
	}
	if err := backups.Delete(entries); err != nil {
		return err
	}
	log.Infof("deleted %d migration backup(s)", len(entries))
	return nil
}

// decryptStep scans for fscrypt-encrypted catalog directories, verifies the
// master key against each, asks per directory whether to permanently
// decrypt, decrypts the confirmed ones, and cleans the fscrypt metadata
// they orphan. A no-op when nothing is encrypted.
func decryptStep(vault *fscrypt.Vault) error {
	encrypted, scanErr := protected.ScanEncryptedCatalogDirs(vault)
	if scanErr != nil {
		return scanErr
	}
	if len(encrypted) == 0 {
		return nil
	}
	if verifyErr := protected.VerifyEncryptedKeys(vault, encrypted); verifyErr != nil {
		return verifyErr
	}
	toDecrypt, pickErr := pickDirsToDecrypt(encrypted)
	if pickErr != nil {
		return pickErr
	}
	if decErr := decryptDirectories(vault, toDecrypt); decErr != nil {
		return decErr
	}
	if len(toDecrypt) == 0 {
		return nil
	}
	log.Infof("permanently decrypted %d directory(ies)", len(toDecrypt))
	// Decrypted dirs left orphaned /.fscrypt metadata behind; remove it
	// while keeping the still-encrypted directories' metadata.
	return cleanOrphanedMetadata()
}

// pickDirsToDecrypt confirms permanent decryption per encrypted dir (default: no).
func pickDirsToDecrypt(encrypted []string) ([]string, error) {
	var toDecrypt []string
	for _, path := range encrypted {
		answer, err := askDecrypt(path)
		if err != nil {
			return nil, err
		}
		if !answer {
			log.Infof("%s stays encrypted", path)
			continue
		}
		toDecrypt = append(toDecrypt, path)
	}
	return toDecrypt, nil
}

// askDecrypt shows path's permanent-decryption confirmation (default: no).
func askDecrypt(path string) (bool, error) {
	answer := false
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Permanently remove fscrypt encryption from %s?", path)).
			Description("The directory remains readable, but the encryption — and its confidentiality — is permanently removed.").
			Affirmative("Decrypt & unprotect").
			Negative("Keep encrypted").
			Value(&answer),
	)).Run(); err != nil {
		return false, err
	}
	return answer, nil
}

// decryptDirectories decrypts each dir behind a progress bar, failing
// fatally on the first error (dir stays encrypted; plaintext copy discarded).
func decryptDirectories(vault *fscrypt.Vault, toDecrypt []string) error {
	if len(toDecrypt) == 0 {
		return nil
	}
	for _, path := range toDecrypt {
		log.Infof("decrypting %s ...", path)
		err := wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
			return vault.DecryptWithProgress(path, func(copied, total int64) {
				fraction := 0.0
				if total > 0 {
					fraction = float64(copied) / float64(total)
				}
				bar.Set(fmt.Sprintf("Decrypting %s", path), fraction)
			})
		})
		if err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		log.Infof("%s is no longer encrypted", path)
	}
	return nil
}

// cleanOrphanedMetadata removes policy/protector metadata orphaned by the
// just-decrypted dirs, scoped like the installer (catalog watch dirs for all
// local users plus system entries); only app-listener-key-* pairs are deleted.
func cleanOrphanedMetadata() error {
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	var paths []string
	catalog := inst.DiscoverForUsers(users)
	for i := range catalog {
		paths = append(paths, catalog[i].Path)
	}

	return fscrypt.CleanOrphanedMetadata(paths)
}
