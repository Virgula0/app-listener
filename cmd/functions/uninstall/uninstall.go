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
  5. reverts the installed systemd service, pacman reload hook, binary,
     PATH symlink and config, and — after a final confirmation (default:
     no) and only for units whose content matches the bundled sample — the
     per-user ssh-agent systemd units
  6. deletes the fscrypt master key ONLY when --delete-key is passed; the
     default keeps it, because without it every still-encrypted directory
     can never be unlocked again

Migration backups (the .app_listener.backup directories) are NOT touched:
use 'app-listener install --restore-backups' / '--delete-post-backups'
for those.

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

	encrypted, err := protected.ScanEncryptedCatalogDirs(vault)
	if err != nil {
		return err
	}

	var decrypted []string
	if len(encrypted) > 0 {
		if err := protected.VerifyEncryptedKeys(vault, encrypted); err != nil {
			return err
		}
		toDecrypt, err := pickDirsToDecrypt(encrypted)
		if err != nil {
			return err
		}
		if err := decryptDirectories(vault, toDecrypt); err != nil {
			return err
		}
		decrypted = toDecrypt
		log.Infof("permanently decrypted %d directory(ies)", len(decrypted))
	}

	// Decrypted dirs left orphaned /.fscrypt metadata behind; remove it
	// while keeping the still-encrypted directories' metadata.
	if len(decrypted) > 0 {
		if err := cleanOrphanedMetadata(); err != nil {
			return err
		}
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

	log.Info("uninstall complete")
	return nil
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
