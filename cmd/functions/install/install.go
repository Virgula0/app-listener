// Package install implements the `app-listener install` wizard: it builds
// the binary if needed, generates the fscrypt master key, lets the user
// pick which users and critical directories to protect (with curated
// whitelisted binaries), migrates directories to fscrypt encryption with a
// non-destructive backup, installs the systemd services and pacman hook,
// and finally enables and verifies the daemon.
package install

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"syscall"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

func init() {
	InstallCmd.Flags().Bool("restore-backups", false,
		"Undo the fscrypt migration: list the found .app_listener.backup directories in the TUI and, after one confirmation, delete the encrypted copies and move the backups back (aborts while the daemon is running)")
	InstallCmd.Flags().Bool("delete-post-backups", false,
		"Delete the found .app_listener.backup directories (listed in the TUI, confirmed once, with progress)")
	InstallCmd.Flags().Bool("binary-only", false,
		"Non-interactive: move the freshly built binary to the install path (and recreate the PATH symlink), then restart the daemon; no wizard, no config, no fscrypt, no systemd units")
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Interactive installer: protect directories with the daemon (fscrypt + systemd)",
	Long: `Interactive, root-only installer for the daemon mode.

The wizard walks through the whole installation:

  0. stops a running daemon before anything else: an active systemd unit is
     stopped (and re-enabled at the end); a daemon process running outside
     systemd is a fatal error — stop it manually first
  1. builds the binary when build/linux/app-listener is missing
  2. ensures the fscrypt master key exists in /etc/app-listener/fscrypt.key
  3. asks which users to protect and probes a catalog of critical
     directories (SSH, keys, AI agents, IDEs, browsers, VPNs, ...)
  4. lets you add further directories manually (they must exist)
  5. shows the generated daemon.conf in an embedded editor for review
  6. checks the encryption state of every configured directory and
     verifies every already-encrypted directory unlocks with the master
     key (a directory encrypted while declared need_encryption: false is
     a fatal error)
  7. asks per-directory whether fscrypt encryption is required, but ONLY
     for directories that are not yet encrypted and declare
     need_encryption: true — need_encryption: false resources are skipped
     silently, and already-encrypted ones are never asked and never
     migrated — then encrypts the ones that need it (keeping a
     .app_listener.backup copy) while a progress bar shows the copy
     progress
  8. installs the systemd units and pacman reload hook from the embedded
     daemon-samples, copies the binary to /usr/local/sbin/app-listener and
     the config to /etc/app-listener/daemon.conf. Already-installed files
     and an existing config are compared with the bundled ones: identical
     files are left alone, differing ones show the diff in the TUI and
     ask whether to overwrite them
  9. enables the daemon across reboots and ensures it is running. When the
     config changed on a running daemon it is reloaded with SIGHUP instead
     of restarted
 10. cleans up orphaned fscrypt metadata: policies and raw-key protectors
     left behind by directories that were encrypted and no longer exist
     (only metadata carrying the installer's own app-listener-key-*
     signature is removed — existing directories keep their metadata)
 11. asks whether to delete the migration backups

The installer is safe to re-run: an interrupted or completed installation
is detected and resumed (existing master key kept, already encrypted
directories verified, identical files skipped, backups never overwritten).

Aborting (Esc) any step stops the installer; the system is left untouched
except for steps that already completed. A directory whose
.app_listener.backup exists aborts the migration with a fatal error: never
overwrite an older backup.

Use --restore-backups to undo the migration instead of installing: the
found .app_listener.backup directories are shown in a TUI list and, after
a single confirmation, the encrypted directories are deleted and the
backups moved back, restoring the original unencrypted content. It aborts
while the daemon is running.

Use --delete-post-backups to delete the found .app_listener.backup
directories instead of installing (also shown in a TUI list first). Both
options show a TUI progress bar while running.

Use --binary-only for a non-interactive shortcut that only deploys the
binary: it builds build/linux/app-listener when it does not exist, stops
a running daemon, replaces /usr/local/sbin/app-listener atomically,
recreates the /usr/local/bin/app-listener symlink and starts the daemon
again — nothing else is touched (no config, no fscrypt, no services).`,
	Args: cobra.NoArgs,
	RunE: runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("install must be run as root: sudo app-listener install")
	}

	done, maintenanceErr := runMaintenanceMode(cmd)
	if maintenanceErr != nil {
		return maintenanceErr
	}
	if done {
		return nil
	}

	// The installer reconfigures files the daemon has open: a daemon left
	// alive during the migration would keep guarding and unlocking the very
	// directories being moved. Stop it first, fatally refusing when it runs
	// outside systemd (where the installer cannot control it).
	if err := systemd.StopDaemonIfRunning(); err != nil {
		return err
	}

	if err := prepareInstallation(); err != nil {
		return err
	}

	cfgText, cfg, err := selectAndEditConfig()
	if err != nil {
		return err
	}
	if len(cfg.Resources) == 0 {
		return errors.New("config contains no [watch] sections")
	}

	vault := fscrypt.New()
	cfgText, err = secureResources(vault, cfgText, cfg)
	if err != nil {
		return err
	}

	if err := deploy(cfgText); err != nil {
		return err
	}

	if err := cleanOrphanedFscrypt(cfg); err != nil {
		return err
	}

	return cleanupBackups(cfg)
}

// runMaintenanceMode reads the non-installing switches and runs the
// requested maintenance mode. It reports whether one was requested (done)
// so the caller knows the installer body must be skipped. The modes are
// mutually exclusive.
func runMaintenanceMode(cmd *cobra.Command) (bool, error) {
	restore, flagErr := cmd.Flags().GetBool("restore-backups")
	if flagErr != nil {
		return false, flagErr
	}
	deleteBackups, flagErr := cmd.Flags().GetBool("delete-post-backups")
	if flagErr != nil {
		return false, flagErr
	}
	binaryOnly, flagErr := cmd.Flags().GetBool("binary-only")
	if flagErr != nil {
		return false, flagErr
	}
	modes := 0
	for _, on := range []bool{restore, deleteBackups, binaryOnly} {
		if on {
			modes++
		}
	}
	if modes > 1 {
		return false, errors.New("--restore-backups, --delete-post-backups and --binary-only are mutually exclusive")
	}
	switch {
	case restore:
		return true, restoreBackups()
	case deleteBackups:
		return true, deletePostBackups()
	case binaryOnly:
		return true, installBinaryOnly()
	}
	return false, nil
}

// installBinaryOnly deploys just the binary: it builds build/linux/app-listener
// when it does not exist, stops a running daemon, replaces the installed
// binary atomically, recreates the PATH symlink and brings the daemon back
// to enabled-and-running. Everything else (config, fscrypt migration,
// systemd units) is deliberately left untouched.
func installBinaryOnly() error {
	if err := buildBinaryIfNeeded(); err != nil {
		return err
	}

	// Same contract as the full installer and `update`: replacing the
	// running binary requires the daemon to be stopped (a daemon outside
	// systemd is a fatal error).
	if err := systemd.StopDaemonIfRunning(); err != nil {
		return err
	}

	if err := systemd.ReplaceInstalledBinary(buildBinaryPath, systemd.InstallBinaryPath); err != nil {
		return fmt.Errorf("installing binary: %w", err)
	}
	log.Infof("installed binary at %s", systemd.InstallBinaryPath)

	if err := systemd.EnsureBinSymlink(); err != nil {
		return err
	}

	if err := systemd.EnableAndVerify(false); err != nil {
		return err
	}
	log.Info("binary-only install complete: the daemon is running the new binary")
	return nil
}

// prepareInstallation covers the two steps that precede any user input:
// building the binary (unless it already exists) and ensuring the master
// key file exists.
func prepareInstallation() error {
	if err := buildBinaryIfNeeded(); err != nil {
		return err
	}
	return ensureMasterKey()
}

// selectAndEditConfig runs the interactive selection flow: users, catalog
// discovery, manual additions, then the embedded editor with strict
// validation of the result.
func selectAndEditConfig() (string, *daemonconfig.Config, error) {
	users, err := pickUsers()
	if err != nil {
		return "", nil, err
	}
	candidates, err := pickDirectories(users)
	if err != nil {
		return "", nil, err
	}
	candidates, err = addManualDirectories(candidates)
	if err != nil {
		return "", nil, err
	}
	if len(candidates) == 0 {
		return "", nil, errors.New("no directories selected: nothing to protect")
	}
	return editConfig(candidates)
}

// secureResources verifies the encryption state of every directory (fatal
// on an encrypted directory declared need_encryption: false, and on a key
// mismatch), asks only for the directories that are NOT yet encrypted AND
// declare need_encryption: true whether encryption is required, and
// migrates the ones that need it. The returned text reflects the final
// need_encryption decisions.
func secureResources(vault *fscrypt.Vault, cfgText string, cfg *daemonconfig.Config) (string, error) {
	if err := verifyEncryptionState(vault, cfg); err != nil {
		return "", err
	}
	if err := askFilesystemsReady(vault, cfg); err != nil {
		return "", err
	}
	updated, toEncrypt, err := askEncryption(vault, cfgText, cfg)
	if err != nil {
		return "", err
	}
	if err := encryptDirectories(vault, toEncrypt); err != nil {
		return "", err
	}
	return updated, nil
}

// askFilesystemsReady fails fast when a target filesystem was never
// initialized for fscrypt (`fscrypt setup`): without it every migration
// aborts deep inside the library, after the user already answered the
// prompts. Filesystems are deduplicated by device so a multi-fs setup is
// each verified exactly once.
func askFilesystemsReady(vault *fscrypt.Vault, cfg *daemonconfig.Config) error {
	var checkedDevs []uint64
	for _, r := range cfg.Resources {
		if !r.NeedEncryption {
			continue
		}
		info, statErr := os.Stat(r.Path)
		if statErr != nil {
			return fmt.Errorf("stat %s: %w", r.Path, statErr)
		}
		dev := info.Sys().(*syscall.Stat_t).Dev
		if slices.Contains(checkedDevs, dev) {
			continue
		}
		checkedDevs = append(checkedDevs, dev)
		if readyErr := vault.CheckFilesystemReady(r.Path); readyErr != nil {
			return readyErr
		}
	}
	return nil
}

// deploy installs the services and hook, copies the binary and the config
// into place, then brings the daemon to the enabled-and-running state.
// Existing files (config and services) are diffed against the bundled
// versions: identical files are left alone, differing ones show the diff
// and ask whether to overwrite. A changed config is delivered to a running
// daemon with SIGHUP.
func deploy(cfgText string) error {
	if err := installServices(); err != nil {
		return err
	}
	configChanged, err := installBinaryAndConfig(cfgText)
	if err != nil {
		return err
	}
	return systemd.EnableAndVerify(configChanged)
}

// buildBinaryIfNeeded compiles the binary with the Makefile build flags
// (CGO_ENABLED=1 GOOS=linux GOARCH=amd64) when build/linux/app-listener
// does not exist yet. Without a source tree in the current directory the
// installer refuses: the binary must exist beforehand in that case.
func buildBinaryIfNeeded() error {
	if _, err := os.Stat(buildBinaryPath); err == nil {
		log.Infof("binary already built: %s", buildBinaryPath)
		return nil
	}
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("source tree not found (no go.mod in %s) and %s does not exist: build the binary with `make build` first",
			mustCwd(), buildBinaryPath)
	}
	log.Infof("building %s ...", buildBinaryPath)
	if err := systemd.RunCmd("go", "build", "-o", buildBinaryPath, "."); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}
	log.Infof("built %s", buildBinaryPath)
	return nil
}

// ensureMasterKey keeps an existing master key (with a warning) and
// creates a fresh one otherwise.
func ensureMasterKey() error {
	exists, err := fscrypt.MasterKeyExists()
	if err != nil {
		return fmt.Errorf("checking master key: %w", err)
	}
	if exists {
		log.Warnf("master key already exists at %s: keeping it (existing fscrypt directories stay usable)", fscrypt.MasterKeyFile)
		return nil
	}
	if err := fscrypt.GenerateMasterKey(false); err != nil {
		return fmt.Errorf("generating master key: %w", err)
	}
	log.Infof("generated new master key at %s", fscrypt.MasterKeyFile)
	return nil
}

// verifyEncryptionState checks every resource's encryption state, always:
// a directory that is encrypted while the config explicitly declares
// need_encryption: false is a fatal error (the daemon would leave an
// encrypted directory unmanaged), and every already-encrypted directory
// must unlock with the current master key (fatal on the first mismatch).
func verifyEncryptionState(vault *fscrypt.Vault, cfg *daemonconfig.Config) error {
	for _, r := range cfg.Resources {
		encrypted, err := vault.IsEncrypted(r.Path)
		if err != nil {
			return fmt.Errorf("checking encryption of %s: %w", r.Path, err)
		}
		if encrypted && !r.NeedEncryption {
			return fmt.Errorf("fatal: %s is already encrypted with fscrypt but the config declares need_encryption: false — set need_encryption: true in the editor, or decrypt the directory first (an encrypted directory must never be left unmanaged)", r.Path)
		}
		if !encrypted {
			continue
		}
		log.Infof("verifying master key against %s ...", r.Path)
		if err := vault.VerifyKey(r.Path); err != nil {
			return fmt.Errorf("fatal: %v — the master key %s does not match the policy of %s; fix the key before installing", err, fscrypt.MasterKeyFile, r.Path)
		}
	}
	return nil
}

// encryptDirectories migrates every not-yet-encrypted directory that was
// selected for encryption. It fails fatally on the first error: backups
// created before the failure are deliberately left in place. While a
// directory is being migrated a TUI progress bar shows the copy progress.
func encryptDirectories(vault *fscrypt.Vault, toEncrypt []string) error {
	if len(toEncrypt) == 0 {
		return nil
	}
	if err := fscrypt.EnsureSystemSetup(); err != nil {
		return fmt.Errorf("setting up fscrypt system-wide: %w", err)
	}
	for _, path := range toEncrypt {
		log.Infof("encrypting %s (backup: %s%s) ...", path, path, fscrypt.BackupSuffix)
		err := wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
			return vault.EncryptWithProgress(path, func(copied, total int64) {
				fraction := 0.0
				if total > 0 {
					fraction = float64(copied) / float64(total)
				}
				bar.Set(fmt.Sprintf("Encrypting %s", path), fraction)
			})
		})
		if err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		log.Infof("%s is now encrypted", path)
	}
	return nil
}

// cleanupBackups asks for each remaining migration backup whether it may
// be deleted. This is the very last step, after the daemon is verified.
func cleanupBackups(cfg *daemonconfig.Config) error {
	for _, r := range cfg.Resources {
		backup := r.Path + fscrypt.BackupSuffix
		if _, err := os.Lstat(backup); err != nil {
			continue
		}
		remove := false
		err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Remove migration backup %s?", backup)).
				Description("The installation succeeded; the backup is no longer needed. The original unencrypted data was already copied into the encrypted directory.").
				Affirmative("Remove").
				Negative("Keep").
				Value(&remove),
		)).Run()
		if err != nil {
			return err
		}
		if !remove {
			log.Infof("keeping %s", backup)
			continue
		}
		if err := os.RemoveAll(backup); err != nil {
			return fmt.Errorf("removing backup %s: %w", backup, err)
		}
		log.Infof("removed backup %s", backup)
	}
	return nil
}
