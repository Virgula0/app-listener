// The `app-listener uninstall` wizard is the counterpart of the installer:
// it reverts everything `app-listener install` deployed. It fatally refuses
// while the daemon is running, re-scans the catalog (NOT the installed
// daemon.conf) to find which guarded directories are currently encrypted,
// verifies the master key against every one, asks per directory whether the
// fscrypt encryption may be permanently removed (default: no) and decrypts
// the confirmed ones with a progress bar, then removes the installed
// systemd units, pacman hook, binary and config. The per-user ssh-agent
// systemd units are only reverted after a separate confirmation (default:
// no) and only when their content matches the bundled sample. The fscrypt
// master key is never deleted unless --delete-key is passed.
package install

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// deleteKeyFlag controls whether the fscrypt master key is removed at the
// end of the uninstall. It defaults to false: without the key, every
// still-encrypted directory can never be unlocked again.
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
  5. reverts the installed systemd service, pacman reload hook, binary and
     config, and — after a final confirmation (default: no) and only for
     units whose content matches the bundled sample — the per-user
     ssh-agent systemd units
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

	// The daemon guards the very directories being decrypted and deleted:
	// reverting it while alive would let it keep unlocking and denying
	// access. Fatally refuse instead of stopping it — there is nothing to
	// re-enable at the end of an uninstall.
	if err := requireDaemonStopped(); err != nil {
		return err
	}

	vault := fscrypt.New()

	encrypted, err := scanEncryptedCatalogDirs(vault)
	if err != nil {
		return err
	}

	var decrypted []string
	if len(encrypted) > 0 {
		if err := verifyEncryptedKeys(vault, encrypted); err != nil {
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

	// The decrypted directories left their policy/protector metadata
	// behind in /.fscrypt: remove the orphaned metadata while keeping the
	// still-encrypted directories' metadata.
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

// requireDaemonStopped fatally refuses the uninstall while an
// app-listener-daemon is running: an active systemd unit and a daemon
// process running outside systemd are both fatal, and the operator must
// stop them manually.
func requireDaemonStopped() error {
	if daemonRunning() {
		return errors.New("fatal: the daemon is running — stop it first: systemctl stop app-listener-daemon")
	}
	pids, err := findDaemonProcesses("/proc")
	if err != nil {
		return fmt.Errorf("scanning for a running daemon process: %w", err)
	}
	if len(pids) > 0 {
		return fmt.Errorf("fatal: the daemon is running outside systemd (pid(s) %v) — stop it first, e.g. kill %d", pids, pids[0])
	}
	log.Info("daemon is not running")
	return nil
}

// scanEncryptedCatalogDirs re-scans the catalog for every local user (plus
// the system-level entries) and returns the directories that are currently
// encrypted with fscrypt. The installed daemon.conf is deliberately not
// consulted: the uninstaller protects whatever is actually protected now.
func scanEncryptedCatalogDirs(vault *fscrypt.Vault) ([]string, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	candidates := inst.DiscoverForUsers(users)
	if len(candidates) == 0 {
		log.Warn("no catalog directories found for any user")
		return nil, nil
	}
	paths := make([]string, 0, len(candidates))
	for i := range candidates {
		paths = append(paths, candidates[i].Path)
	}
	encrypted, err := existingEncryptedDirs(vault, paths)
	if err != nil {
		return nil, err
	}
	if len(encrypted) == 0 {
		log.Info("no catalog directory is encrypted: nothing to decrypt")
	}
	return encrypted, nil
}

// existingEncryptedDirs returns the paths from the given set that currently
// carry an fscrypt policy. Regular files and plain directories are skipped
// (a file can never carry a policy).
func existingEncryptedDirs(vault *fscrypt.Vault, paths []string) ([]string, error) {
	var encrypted []string
	for _, path := range paths {
		ok, err := vault.IsEncrypted(path)
		if err != nil {
			return nil, fmt.Errorf("checking encryption of %s: %w", path, err)
		}
		if ok {
			log.Infof("%s is encrypted", path)
			encrypted = append(encrypted, path)
		}
	}
	return encrypted, nil
}

// verifyEncryptedKeys checks the master key against every encrypted
// directory. A directory that does not unlock with the master key is a
// fatal error: removing the daemon (and possibly the key with --delete-key)
// would lock that directory forever since no usable key remains.
func verifyEncryptedKeys(vault *fscrypt.Vault, encrypted []string) error {
	for _, path := range encrypted {
		log.Infof("verifying master key against %s ...", path)
		if err := vault.VerifyKey(path); err != nil {
			return fmt.Errorf("fatal: %v — the master key %s does not match the policy of %s; restore the correct key before uninstalling", err, fscrypt.MasterKeyFile, path)
		}
	}
	return nil
}

// pickDirsToDecrypt shows, per encrypted directory, the confirmation
// (default: no) whether its fscrypt encryption must be permanently removed
// with the master key. Only the confirmed directories are returned.
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

// askDecrypt shows the per-directory confirmation (default: no) whether the
// fscrypt encryption of path must be permanently removed using the master
// key.
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

// decryptDirectories permanently decrypts every directory in the given
// list. While a directory is being decrypted a TUI progress bar shows the
// copy progress. It fails fatally and leaves the directory encrypted on the
// first error: the plaintext copy is discarded, never silently swapped in.
func decryptDirectories(vault *fscrypt.Vault, toDecrypt []string) error {
	if len(toDecrypt) == 0 {
		return nil
	}
	for _, path := range toDecrypt {
		log.Infof("decrypting %s ...", path)
		err := withBottomBar(func(bar *bottomBar) error {
			return vault.DecryptWithProgress(path, func(copied, total int64) {
				fraction := 0.0
				if total > 0 {
					fraction = float64(copied) / float64(total)
				}
				bar.set(fmt.Sprintf("Decrypting %s", path), fraction)
			})
		})
		if err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		log.Infof("%s is no longer encrypted", path)
	}
	return nil
}

// cleanOrphanedMetadata removes the fscrypt policy and protector metadata
// left behind by directories that were decrypted moments ago. The live set
// is derived from the catalog exactly like the installer scopes itself:
// every discoverable watch directory for every local user plus the
// system-level entries. Metadata of directories that still exist and are
// still encrypted is kept; only orphaned app-listener-key-* raw-key
// protector/policy pairs are deleted.
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

	live, err := fscrypt.LivePolicyDescriptors(paths)
	if err != nil {
		return err
	}

	anchor := "/"
	for _, p := range paths {
		if _, statErr := os.Stat(p); statErr == nil {
			anchor = p
			break
		}
	}
	removedProtectors, removedPolicies, err := fscrypt.CleanOrphans(anchor, live)
	if err != nil {
		return fmt.Errorf("cleaning orphaned fscrypt metadata: %w", err)
	}
	if removedProtectors == 0 && removedPolicies == 0 {
		log.Info("no orphaned fscrypt metadata to clean")
		return nil
	}
	log.Infof("cleaned orphaned fscrypt metadata: removed %d protector(s) and %d policy(ies)", removedProtectors, removedPolicies)
	return nil
}
