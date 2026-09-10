// Package install implements the `app-listener install` wizard: deploys the
// binary if missing (copying the running executable, never recompiling),
// ensures the fscrypt master key, picks users and critical dirs, migrates to
// fscrypt with backups, installs units/hook, enables the daemon.
package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
		"Non-interactive: move the freshly built binary to the install path (and recreate the PATH symlink), then restart the daemon if its service is installed; no wizard, no config, no fscrypt, no systemd units")
	InstallCmd.Flags().Bool("update-catalog-only", false,
		"Re-scan the catalog whitelists for every guarded directory in the existing daemon.conf: unlocks encrypted vaults, re-expands glob patterns, drops deleted binaries, picks up new ones, and overwrites the config (use --yes to skip confirmation; requires a previous installation)")
	InstallCmd.Flags().Bool("diff-catalog", false,
		"Diff the catalog against the installed daemon.conf: list the critical directories that exist on the host but are not yet guarded, let you pick which to add, then append them to the config and encrypt them like a fresh install (existing sections untouched; stops the daemon for the cycle; requires a previous installation)")
	InstallCmd.Flags().Bool("live", false,
		"With --update-catalog-only: refresh the whitelists WITHOUT stopping the daemon (vaults are already unlocked and guarded by the running daemon; the config change is applied via SIGHUP reload). Requires the daemon to be running")
	InstallCmd.Flags().BoolP("yes", "y", false,
		"Skip all confirmation prompts (use with --update-catalog-only for non-interactive use, e.g. the pacman/apt hooks and the boot-time refresh unit)")
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Interactive installer: protect directories with the daemon (fscrypt + systemd)",
	Long: `Interactive, root-only installer for the daemon mode.

The wizard walks through the whole installation:

  0. stops a running daemon before anything else: an active systemd unit is
     stopped (and re-enabled at the end); a daemon process running outside
     systemd is a fatal error — stop it manually first
  1. deploys the app-listener binary to /usr/local/sbin/app-listener when it
     is not already there, by copying the running executable as-is (no
     recompilation). An existing install is left untouched — upgrade it with
     "install --binary-only" or "app-listener update"
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
  8. installs the systemd units and the package-manager catalog-refresh
     hook from the embedded daemon-samples and writes the config to
     /etc/app-listener/daemon.conf (the binary is left as deployed). The
     PATH symlink /usr/local/bin/app-listener is (re)created. The refresh
     hook is chosen from the
     host's package manager: pacman (/etc/pacman.d/hooks) or apt
     (/etc/apt/apt.conf.d); dnf/zypper are reported and left to the
     boot-time refresh. Already-installed files and an existing config are
     compared with the bundled ones: identical files are left alone,
     differing ones show the diff in the TUI and ask whether to overwrite
  9. enables the daemon across reboots and ensures it is running, and
     enables app-listener-catalog-refresh.service — a boot-time --live
     catalog refresh that catches package changes made while no hook ran
     (offline installs, direct dpkg/pacman -U, a failed hook). When the
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
again — nothing else is touched (no config, no fscrypt, no services).
When the daemon service is not installed yet it just deploys the binary
and the symlink and stops there (there is no unit to restart).

Use --update-catalog-only to refresh the whitelist of every guarded
directory from the catalog: encrypted vaults are unlocked, glob patterns
are re-expanded (picking up new binaries and dropping deleted ones),
and the diff is shown before overwriting. User-added sections (not in
the catalog) are preserved as-is. Requires a previous installation.

Use --diff-catalog to add directories the catalog now discovers but the
installed daemon.conf does not yet guard (a catalog update, or a newly
installed application): the missing critical directories are listed in the
same picker as the wizard, the selected ones are appended to the config
and encrypted like a fresh install, and every existing section is left
byte-for-byte intact. It stops the daemon for the cycle and restarts it on
the merged config. Requires a previous installation. Unlike
--update-catalog-only this never touches existing sections' whitelists —
run both to fully re-sync.`,
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

	// The daemon holds open the files being reconfigured and would keep
	// guarding the dirs being moved: stop it first; outside-systemd daemons
	// are fatally refused (the installer cannot control them).
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

	// Collected now, persisted last (finalizeEditPassword): an aborted
	// install must never leave a dangling password hash.
	editPassword, err := promptEditPassword(cfg)
	if err != nil {
		return err
	}

	if err := deploy(cfgText); err != nil {
		return err
	}

	if err := cleanOrphanedFscrypt(cfg); err != nil {
		return err
	}

	if err := cleanupBackups(cfg); err != nil {
		return err
	}

	// The very last install action.
	return finalizeEditPassword(editPassword)
}

// runMaintenanceMode dispatches the mutually exclusive maintenance flags,
// reporting whether one ran so the caller skips the installer body.
// maintenanceFlags holds the parsed CLI flags for the install command's
// mutually exclusive maintenance modes.
type maintenanceFlags struct {
	restore       bool
	deleteBackup  bool
	binaryOnly    bool
	updateCatalog bool
	diffCatalog   bool
	live          bool
	autoConfirm   bool
}

func parseMaintenanceFlags(cmd *cobra.Command) (maintenanceFlags, error) {
	var f maintenanceFlags
	var err error
	if f.restore, err = cmd.Flags().GetBool("restore-backups"); err != nil {
		return f, err
	}
	if f.deleteBackup, err = cmd.Flags().GetBool("delete-post-backups"); err != nil {
		return f, err
	}
	if f.binaryOnly, err = cmd.Flags().GetBool("binary-only"); err != nil {
		return f, err
	}
	if f.updateCatalog, err = cmd.Flags().GetBool("update-catalog-only"); err != nil {
		return f, err
	}
	if f.diffCatalog, err = cmd.Flags().GetBool("diff-catalog"); err != nil {
		return f, err
	}
	if f.live, err = cmd.Flags().GetBool("live"); err != nil {
		return f, err
	}
	if f.autoConfirm, err = cmd.Flags().GetBool("yes"); err != nil {
		return f, err
	}
	return f, nil
}

func runMaintenanceMode(cmd *cobra.Command) (bool, error) {
	f, err := parseMaintenanceFlags(cmd)
	if err != nil {
		return false, err
	}
	if err := validateMaintenanceFlags(f); err != nil {
		return false, err
	}
	switch {
	case f.restore:
		return true, restoreBackups()
	case f.deleteBackup:
		return true, deletePostBackups()
	case f.binaryOnly:
		return true, installBinaryOnly()
	case f.updateCatalog:
		return true, runUpdateCatalogOnly(f.autoConfirm, f.live)
	case f.diffCatalog:
		return true, runDiffCatalog()
	}
	return false, nil
}

// validateMaintenanceFlags rejects flag combinations the maintenance modes
// do not support: --yes / --live outside their owning modes, and more than
// one mutually exclusive mode at once.
func validateMaintenanceFlags(f maintenanceFlags) error {
	if f.autoConfirm && !f.updateCatalog && !f.restore && !f.deleteBackup {
		return errors.New("--yes can only be used with --update-catalog-only, --restore-backups, or --delete-post-backups")
	}
	if f.live && !f.updateCatalog {
		return errors.New("--live can only be used with --update-catalog-only")
	}
	modes := 0
	for _, on := range []bool{f.restore, f.deleteBackup, f.binaryOnly, f.updateCatalog, f.diffCatalog} {
		if on {
			modes++
		}
	}
	if modes > 1 {
		return errors.New("--restore-backups, --delete-post-backups, --binary-only, --update-catalog-only and --diff-catalog are mutually exclusive")
	}
	return nil
}

// installBinaryOnly deploys only the freshly built binary (see
// systemd.DeployInstalledBinary); config, fscrypt and units stay untouched.
func installBinaryOnly() error {
	if err := buildBinaryIfNeeded(); err != nil {
		return err
	}

	if err := systemd.DeployInstalledBinary(buildBinaryPath); err != nil {
		return fmt.Errorf("installing binary: %w", err)
	}
	if !systemd.DaemonUnitInstalled() {
		log.Info("binary-only install complete: binary deployed; the daemon service is not installed (run `sudo app-listener install`)")
		return nil
	}
	log.Info("binary-only install complete: the daemon is running the new binary")
	return nil
}

// deliverReload applies a patched config to the daemon: SIGHUP reload
// (atomic — new guards attach before old detach) with a restart fallback
// when running, or a start after the stopped flow. Package-level
// indirection so tests can observe the delivery.
var deliverReload = systemd.EnableAndVerify

// runUpdateCatalogOnly re-scans the catalog whitelist for every guarded
// directory in the existing daemon.conf. Encrypted vaults are unlocked
// for stat access, glob patterns are re-expanded, and the diff is shown
// for confirmation. User-added sections are preserved verbatim.
//
// In live mode the daemon keeps running: its vaults are already unlocked
// and its guards attached, so the re-scan needs no unlock/lock cycle and
// no daemon stop — a changed config is delivered via SIGHUP reload. A
// stopped daemon cannot support live mode (locked vaults hide plaintext
// names), so the caller is told to use the stopped flow instead.
func runUpdateCatalogOnly(autoConfirm, live bool) error {
	if live {
		if !systemd.IsDaemonActive() {
			return errors.New("live catalog refresh requires the daemon to be running: " +
				"locked vaults hide plaintext names — use --update-catalog-only without --live, or start the daemon first")
		}
		vault := fscrypt.New()
		changed, err := updateCatalogConfig(vault, autoConfirm, true)
		if err != nil {
			return softenAutomatedRefreshErr(err, autoConfirm)
		}
		return applyLiveRefresh(changed)
	}
	wasActive := systemd.IsDaemonActive()
	if err := systemd.StopDaemonIfRunning(); err != nil {
		return err
	}
	vault := fscrypt.New()
	if _, err := updateCatalogConfig(vault, autoConfirm, false); err != nil {
		// updateCatalogConfig writes the config only after every section
		// patched successfully, so on error the on-disk config is unchanged.
		// The daemon was stopped by this flow (StopDaemonIfRunning) — bring
		// it back on the existing config so a failed refresh never silently
		// leaves protection off, then surface the original error.
		if wasActive {
			if restartErr := deliverReload(true); restartErr != nil {
				log.Errorf("catalog refresh failed AND the daemon could not be restarted: %v", restartErr)
			}
		}
		return softenAutomatedRefreshErr(err, autoConfirm)
	}
	// The daemon was stopped by this flow: it must run again regardless of
	// whether the config changed.
	return deliverReload(true)
}

// softenAutomatedRefreshErr turns the "all-manual config, nothing matched the
// catalog" outcome into a no-op for automated callers (--yes: the pacman/apt
// hooks and the boot-time refresh unit). A package transaction or a boot must
// not be reported as failed just because the config has no catalog-managed
// section. Interactive callers still see the error.
func softenAutomatedRefreshErr(err error, autoConfirm bool) error {
	if err != nil && autoConfirm && errors.Is(err, errNoCatalogMatch) {
		log.Infof("catalog refresh: %v — nothing to do", err)
		return nil
	}
	return err
}

// applyLiveRefresh delivers a patched config to the running daemon: SIGHUP
// reload (atomic — new guards attach before old detach), restart fallback.
// A refresh whose result matches the running config is a no-op.
//
// Regression guard for the first live implementation: the config was
// patched on disk but never delivered to the running daemon, which kept
// enforcing the old whitelist while journalctl stayed silent — the
// delivery is the mandatory last step of live mode.
func applyLiveRefresh(changed bool) error {
	if !changed {
		log.Info("nothing to reload: the refreshed config matches the running daemon")
		return nil
	}
	return deliverReload(true)
}

// prepareInstallation deploys the binary if it is missing (copying the
// running executable — never recompiling) and ensures the master key exists.
func prepareInstallation() error {
	if err := ensureInstalledBinary(); err != nil {
		return err
	}
	return ensureMasterKey()
}

// installedBinaryPath is the service path ensureInstalledBinary manages; a
// package var so tests can point it at a temp file.
var installedBinaryPath = systemd.InstallBinaryPath

// ensureInstalledBinary makes sure the app-listener binary is deployed at its
// service path. An existing binary is left untouched — upgrades go through
// `install --binary-only` or `app-listener update`, which also restart the
// daemon. A missing binary (the common first run straight after `make build`,
// or after the one-line installer) is filled in by copying the currently
// running executable into place as-is: no recompilation, and no separate
// --binary-only step needed for a one-time install.
func ensureInstalledBinary() error {
	if _, err := os.Stat(installedBinaryPath); err == nil {
		log.Infof("binary is installed: %s", installedBinaryPath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking %s: %w", installedBinaryPath, err)
	}

	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locating the running binary to deploy: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(self); resolveErr == nil {
		self = resolved
	}
	log.Infof("binary not installed yet: deploying the running executable %s -> %s (no recompilation)", self, installedBinaryPath)
	if err := os.MkdirAll(filepath.Dir(installedBinaryPath), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(installedBinaryPath), err)
	}
	if err := systemd.ReplaceInstalledBinary(self, installedBinaryPath); err != nil {
		return fmt.Errorf("deploying the binary: %w", err)
	}
	return nil
}

// selectAndEditConfig runs users, catalog, manual additions, then the editor.
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

// secureResources verifies encryption state (fatal: encrypted dir declared
// need_encryption: false, or master-key mismatch), asks only unencrypted
// dirs with need_encryption: true and migrates them; text reflects decisions.
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

// askFilesystemsReady fails fast on filesystems lacking `fscrypt setup`,
// before any prompting; dedup by device verifies each fs exactly once.
// Grouped sections are checked once, at their encryption root (the fscrypt
// lifecycle is per vault root, and every watch sub-path lives on it).
func askFilesystemsReady(vault *fscrypt.Vault, cfg *daemonconfig.Config) error {
	var checkedDevs []uint64
	for _, r := range cfg.EncryptionGroups() {
		if !r.NeedEncryption {
			continue
		}
		root := r.EncryptionRootOrPath()
		info, statErr := os.Stat(root)
		if statErr != nil {
			return fmt.Errorf("stat %s: %w", root, statErr)
		}
		dev := info.Sys().(*syscall.Stat_t).Dev
		if slices.Contains(checkedDevs, dev) {
			continue
		}
		checkedDevs = append(checkedDevs, dev)
		if readyErr := vault.CheckFilesystemReady(root); readyErr != nil {
			return readyErr
		}
	}
	return nil
}

// deploy installs services/hook, copies binary+config, enables the daemon.
// Existing files are diffed: identical ones stay, differing ones show the
// diff and ask; a changed config reaches a running daemon via SIGHUP.
func deploy(cfgText string) error {
	if err := installServices(); err != nil {
		return err
	}
	configChanged, err := installConfig(cfgText)
	if err != nil {
		return err
	}
	if err := preflightDeployedBinary(); err != nil {
		return err
	}
	if err := systemd.EnableAndVerify(configChanged); err != nil {
		return err
	}
	return systemd.EnableCatalogRefresh()
}

// preflightDeployedBinary runs `<installed binary> daemon --check` against the
// just-deployed binary and config BEFORE the service is enabled. It verifies
// the BPF-LSM stack is active and that every guard eBPF program is accepted by
// this kernel's verifier (attaching nothing). A failure here means the daemon
// would crash-loop — or, with a prebuilt object too complex for a newer
// kernel, panic — so the install aborts now, having changed only the on-disk
// binary/config, rather than leaving an enabled unit that never runs.
func preflightDeployedBinary() error {
	log.Info("preflight: verifying the guard eBPF loads on this kernel ...")
	if err := systemd.RunCmd(systemd.InstallBinaryPath, "daemon", "--check", "--config", systemd.SystemConfigPath); err != nil {
		log.Error("if you installed a prebuilt release binary, rebuild from source on this host so " +
			"the eBPF is compiled against this kernel: " +
			"git clone https://github.com/Virgula0/app-listener && cd app-listener && " +
			"make build && sudo ./build/linux/app-listener install")
		return fmt.Errorf("guard eBPF preflight failed — the daemon was NOT enabled "+
			"(binary and config are in place, but no service is running): %w", err)
	}
	log.Info("preflight OK")
	return nil
}

// buildBinaryIfNeeded compiles build/linux/app-listener when it is missing;
// without a source tree it refuses. Only the --binary-only path uses it —
// the full wizard never builds (see ensureInstalledBinary).
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

// verifyEncryptionState checks every resource: encrypted while declared
// need_encryption: false is fatal (unmanaged encryption), and encrypted
// dirs must unlock with the current master key (fatal on mismatch).
// Grouped sections are verified once per encryption root — the fscrypt
// policy is per vault root and every watch sub-path inherits it.
func verifyEncryptionState(vault *fscrypt.Vault, cfg *daemonconfig.Config) error {
	for _, r := range cfg.EncryptionGroups() {
		root := r.EncryptionRootOrPath()
		encrypted, err := vault.IsEncrypted(root)
		if err != nil {
			return fmt.Errorf("checking encryption of %s: %w", root, err)
		}
		if encrypted && !r.NeedEncryption {
			return fmt.Errorf("fatal: %s is already encrypted with fscrypt but the config declares need_encryption: false — set need_encryption: true in the editor, or decrypt the directory first (an encrypted directory must never be left unmanaged)", root)
		}
		if !encrypted {
			continue
		}
		log.Infof("verifying master key against %s ...", root)
		if err := vault.VerifyKey(root); err != nil {
			return fmt.Errorf("fatal: %v — the master key %s does not match the policy of %s; fix the key before installing", err, fscrypt.MasterKeyFile, root)
		}
	}
	return nil
}

// encryptDirectories migrates the selected dirs behind a progress bar,
// failing fatally on the first error (prior backups deliberately kept).
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

// cleanupBackups asks per backup whether to delete it (final step, post-verify).
// The migration backup is created at the encryption root (the whole vault is
// renamed aside, not each watch sub-path), so grouped sections are handled
// once, at that root.
func cleanupBackups(cfg *daemonconfig.Config) error {
	for _, r := range cfg.EncryptionGroups() {
		backup := r.EncryptionRootOrPath() + fscrypt.BackupSuffix
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
