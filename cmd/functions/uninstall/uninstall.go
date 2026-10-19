// The `app-listener uninstall` command reverts the installer: refuses while the daemon runs,
// re-scans the catalog (not daemon.conf) for encrypted dirs, verifies the master key, asks per dir
// to decrypt permanently (default no), removes units/hook/binary/symlink/config, reverts
// sample-matching ssh-agent units after a separate confirmation, and deletes the key only with
// --delete-key.
package uninstall

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

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

// keepGUIFlag leaves the desktop build (systemd.GUIBinaryPath) installed.
var keepGUIFlag bool

func init() {
	UninstallCmd.Flags().BoolVar(&deleteKeyFlag, "delete-key", false,
		"Also delete the fscrypt master key (without it, directories left encrypted can never be unlocked)")
	UninstallCmd.Flags().BoolVar(&keepGUIFlag, "keep-gui", false,
		"Keep the desktop build (/usr/local/sbin/app-listener-gui) and its PATH symlink")
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the install: decrypt the directories you pick, remove the service",
	Long: `Revert what install did. Refused while the daemon runs. It re-scans the catalog
for encrypted directories, checks each unlocks with the master key, decrypts
the ones you pick, then removes the service, binaries and config, and offers
to remove the ssh-agent setup and the migration backups. The master key is
kept unless --delete-key.`,
	Args: cobra.NoArgs,
	RunE: runUninstall,
}

// runUninstall drives the whole uninstall flow.
func runUninstall(cmd *cobra.Command, args []string) error {
	return run(Request{DeleteKey: deleteKeyFlag, KeepGUI: keepGUIFlag})
}

// Request is one uninstall: the --delete-key and --keep-gui choices.
type Request struct {
	DeleteKey, KeepGUI bool
}

var runMu sync.Mutex

// RunWith uninstalls with every prompt on ui: the desktop GUI's uninstaller. One run at a time:
// wizard.Use is process-wide.
func RunWith(ui wizard.UI, req Request) error {
	if !runMu.TryLock() {
		return errors.New("an uninstall is already in progress")
	}
	defer runMu.Unlock()
	defer wizard.Use(ui)()
	return run(req)
}

func run(req Request) error {
	if os.Geteuid() != 0 {
		return errors.New("uninstall must be run as root: sudo app-listener uninstall")
	}

	// The daemon guards the very dirs being decrypted/deleted; refuse while
	// it lives instead of stopping it — nothing gets re-enabled afterwards.
	wizard.Stage("Checking the daemon is stopped")
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
	}
	// First: an uninstall interrupted past this point leaves a system the wizard may resume.
	if err := protected.RemoveInstallCert(); err != nil {
		return err
	}

	vault := fscrypt.New()

	wizard.Stage("Decrypting")
	stillEncrypted, err := decryptStep(vault)
	if err != nil {
		return err
	}
	deleteKey := req.DeleteKey
	if deleteKey && len(stillEncrypted) > 0 {
		if deleteKey, err = confirmKeyDeletion(stillEncrypted); err != nil {
			return err
		}
	}

	// Scan for migration backups while daemon.conf still exists (it names manual dirs the catalog
	// doesn't); the prompt runs at the end, after revertSystemFiles.
	backupEntries, err := backups.Find()
	if err != nil {
		return err
	}

	wizard.Stage("Reverting the per-user setup")
	if err := revertSSHAgents(); err != nil {
		return err
	}

	if err := revertBunLaunchers(); err != nil {
		return err
	}

	wizard.Stage("Removing the installed files")
	if err := revertSystemFiles(req.KeepGUI); err != nil {
		return err
	}

	if deleteKey {
		if err := removeMasterKey(); err != nil {
			return err
		}
	} else {
		log.Infof("keeping the fscrypt master key at %s (pass --delete-key to remove it)", fscrypt.MasterKeyFile)
	}

	wizard.Stage("Cleaning up backups")
	if err := offerBackupCleanup(backupEntries); err != nil {
		return err
	}

	log.Info("uninstall complete")
	return nil
}

// offerBackupCleanup asks at the end whether to delete the .app_listener.backup dirs the migration
// left. They're plain unencrypted copies, usually pointless but harmless to keep, so the user
// chooses (all preselected, one confirmation).
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

// confirmKeyDeletion asks before --delete-key locks the still-encrypted dirs forever (default: keep).
func confirmKeyDeletion(stillEncrypted []string) (bool, error) {
	del, err := wizard.Confirm(wizard.Question{
		Title: fmt.Sprintf("Delete the master key although %d directory(ies) stay encrypted?", len(stillEncrypted)),
		Description: "Without the key these can never be unlocked again:\n  " +
			strings.Join(stillEncrypted, "\n  "),
		Affirmative: "Delete the key",
		Negative:    "Keep the key",
	})
	if err != nil {
		return false, err
	}
	if !del {
		log.Infof("keeping the fscrypt master key at %s: directories stay encrypted", fscrypt.MasterKeyFile)
	}
	return del, nil
}

// decryptStep scans for fscrypt-encrypted catalog dirs, verifies the master key against each, asks
// per dir whether to decrypt permanently, decrypts the confirmed ones and cleans the fscrypt
// metadata they orphan; it returns the dirs left encrypted. No-op when nothing is encrypted.
func decryptStep(vault *fscrypt.Vault) ([]string, error) {
	encrypted, scanErr := protected.ScanEncryptedCatalogDirs(vault)
	if scanErr != nil {
		return nil, scanErr
	}
	if len(encrypted) == 0 {
		return nil, nil
	}
	if verifyErr := protected.VerifyEncryptedKeys(vault, encrypted); verifyErr != nil {
		return nil, verifyErr
	}
	toDecrypt, pickErr := pickDirsToDecrypt(encrypted)
	if pickErr != nil {
		return nil, pickErr
	}
	kept := slices.DeleteFunc(slices.Clone(encrypted), func(p string) bool { return slices.Contains(toDecrypt, p) })
	if decErr := decryptDirectories(vault, toDecrypt); decErr != nil {
		return nil, decErr
	}
	if len(toDecrypt) == 0 {
		return kept, nil
	}
	log.Infof("permanently decrypted %d directory(ies)", len(toDecrypt))
	// Decrypted dirs left orphaned /.fscrypt metadata behind; remove it
	// while keeping the still-encrypted directories' metadata.
	return kept, cleanOrphanedMetadata()
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
	return wizard.Confirm(wizard.Question{
		Title:       fmt.Sprintf("Permanently remove fscrypt encryption from %s?", path),
		Description: "The directory remains readable, but the encryption — and its confidentiality — is permanently removed.",
		Affirmative: "Decrypt & unprotect",
		Negative:    "Keep encrypted",
	})
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

// cleanOrphanedMetadata removes policy/protector metadata orphaned by the just-decrypted dirs,
// scoped like the installer (catalog watch dirs for all users plus system entries); only
// app-listener-key-* pairs are deleted.
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
