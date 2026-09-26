// Package install implements the `app-listener install` wizard: deploys the binary if missing
// (copies the running executable, never recompiles), ensures the fscrypt master key, picks users
// and dirs, migrates to fscrypt with backups, installs units/hook, enables the daemon.
package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

func init() {
	InstallCmd.Flags().Bool("restore-backups", false,
		"Undo the fscrypt migration: list the found .app_listener.backup directories in the TUI and, after one confirmation, delete the encrypted copies and move the backups back (aborts while the daemon is running)")
	InstallCmd.Flags().Bool("delete-post-backups", false,
		"Delete the found .app_listener.backup directories (listed in the TUI, confirmed once, with progress)")
	InstallCmd.Flags().Bool("binary-only", false,
		"Move the freshly built binary to the install path (and recreate the PATH symlink); no wizard, no config, no fscrypt, no systemd units. If the daemon service is installed, asks (TUI) whether to restart it now to run the new binary or skip the restart for later")
	InstallCmd.Flags().Bool("update-catalog-only", false,
		"Re-scan the catalog whitelists for every guarded directory in the existing daemon.conf: unlocks encrypted vaults, re-expands glob patterns, drops deleted binaries, picks up new ones, and overwrites the config (use --yes to skip confirmation; requires a previous installation)")
	InstallCmd.Flags().Bool("diff-catalog", false,
		"Diff the catalog against the installed daemon.conf: list the critical directories that exist on the host but are not yet guarded, let you pick which to add, then append them to the config and encrypt them like a fresh install (existing sections untouched; stops the daemon for the cycle; requires a previous installation)")
	InstallCmd.Flags().Bool("live", false,
		"With --update-catalog-only: refresh the whitelists WITHOUT stopping the daemon (vaults are already unlocked and guarded by the running daemon; the config change is applied via SIGHUP reload). Requires the daemon to be running")
	InstallCmd.Flags().BoolVar(&allowMetadataOutput, "allow-metadata-output", false,
		"Install the daemon unit WITHOUT --no-log-metadata-blocks, so denied metadata-only process inspections (op=PTRACE mode=READ) are logged; useful to diagnose an app that breaks with nothing logged")
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
  6a. for every user whose ~/.ssh ends up guarded, asks once (naming the
     user and paths) whether to set up ssh-agent: unit + shell SSH_AUTH_SOCK
     block (installed in step 8) and "AddKeysToAgent yes" at the top of
     ~/.ssh/config (created if missing, an existing AddKeysToAgent is kept),
     done before encryption so the file is part of the migrated tree
  6b. checks each backing filesystem is fscrypt-ready; for a fixable gap
     (missing 'encrypt' feature flag, no 'fscrypt setup') it shows the
     exact command and its reason and offers to run it now as root —
     declining aborts the install, just like the old hard error
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
     differing ones show the diff in the TUI and ask whether to overwrite.
     The per-user ssh-agent systemd unit (and the SSH_AUTH_SOCK block in the
     user's shell rc file) is installed for the users who accepted it in
     step 6a
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

Use --binary-only for a shortcut that only deploys the binary: it builds
build/linux/app-listener when it does not exist and replaces
/usr/local/sbin/app-listener atomically, recreating the
/usr/local/bin/app-listener symlink — nothing else is touched (no config,
no fscrypt, no services). When the daemon service is not installed yet it
just deploys the binary and the symlink and stops there (there is no unit
to restart). Otherwise it asks (TUI) whether to restart the daemon now —
stopping it, replacing the binary and starting it again, the historical
default — or skip the restart, which still replaces the binary file but
leaves the running daemon on its old (already-loaded) binary until you
restart it yourself.

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

	if err := checkRunningBinaryMatchesInstalled(); err != nil {
		return err
	}

	// Stop the daemon first: it holds the files being reconfigured and would keep guarding the dirs
	// being moved. Outside-systemd daemons are fatally refused (can't be controlled).
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

	sshUsers, bunUsers, err := gatherPerUserSetup(cfg)
	if err != nil {
		return err
	}

	vault := fscrypt.New()
	cfgText, err = secureResources(vault, cfgText, cfg)
	if err != nil {
		return err
	}

	// Collected now, persisted in deploy() only after every step that could still abort has
	// succeeded, just before the daemon's first start (writeEditPasswordHash).
	editPassword, err := promptEditPassword(cfg)
	if err != nil {
		return err
	}

	if err := deploy(cfgText, sshUsers, bunUsers, editPassword); err != nil {
		return err
	}

	if err := cleanOrphanedFscrypt(cfg); err != nil {
		return err
	}

	return cleanupBackups(cfg)
}

// maintenanceFlags holds the parsed CLI flags of the install command's mutually exclusive
// maintenance modes.
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

// runMaintenanceMode dispatches the mutually exclusive maintenance flags, reporting whether one
// ran so the caller skips the installer body.
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

// validateMaintenanceFlags rejects unsupported combinations: --yes/--live outside their owning
// modes, and more than one mode at once.
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

// installBinaryOnly deploys only the freshly built binary (systemd.DeployInstalledBinary); config,
// fscrypt and units are untouched. If the service is installed, asks whether to restart the daemon
// now or leave it on the old binary until a manual restart.
func installBinaryOnly() error {
	if err := buildBinaryIfNeeded(); err != nil {
		return err
	}

	if !systemd.DaemonUnitInstalled() {
		if err := systemd.DeployInstalledBinary(buildBinaryPath); err != nil {
			return fmt.Errorf("installing binary: %w", err)
		}
		log.Info("binary-only install complete: binary deployed; the daemon service is not installed (run `sudo app-listener install`)")
		return nil
	}

	restart := true
	if formErr := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Restart the daemon now?").
			Description("The daemon service is installed. Restarting picks up the new binary immediately (the current default). Skipping leaves the running daemon on the old binary — restart it yourself later (systemctl restart app-listener-daemon) to run the new one.").
			Affirmative("Yes, restart now").
			Negative("No, skip the restart").
			Value(&restart),
	)).Run(); formErr != nil {
		return formErr
	}

	if !restart {
		if err := systemd.DeployInstalledBinaryNoRestart(buildBinaryPath); err != nil {
			return fmt.Errorf("installing binary: %w", err)
		}
		log.Info("binary-only install complete: binary deployed; the daemon restart was skipped — restart it manually (systemctl restart app-listener-daemon) to run the new binary")
		return nil
	}

	if err := systemd.DeployInstalledBinary(buildBinaryPath); err != nil {
		return fmt.Errorf("installing binary: %w", err)
	}
	log.Info("binary-only install complete: the daemon is running the new binary")
	return nil
}

// deliverReload applies a patched config to the daemon: SIGHUP reload (atomic; new guards attach
// before old detach) with a restart fallback when running, or a start after the stopped flow. A
// package var so tests can observe delivery.
var deliverReload = systemd.EnableAndVerify

// runUpdateCatalogOnly re-scans the catalog whitelist for every guarded directory in daemon.conf:
// encrypted vaults are unlocked for stat, globs re-expanded, the diff shown for confirmation;
// user-added sections kept verbatim.
//
// Live mode: the daemon keeps running (vaults unlocked, guards attached), so no unlock/lock cycle
// or stop is needed and a changed config is delivered via SIGHUP. A stopped daemon can't support it
// (locked vaults hide plaintext names), so the caller is told to use the stopped flow.
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
		// updateCatalogConfig writes the config only after every section patched, so an error
		// leaves the on-disk config unchanged. The flow stopped the daemon: bring it back on the
		// existing config so a failed refresh never leaves protection off, then surface the
		// original error.
		if wasActive {
			if restartErr := deliverReload(true); restartErr != nil {
				log.Errorf("catalog refresh failed AND the daemon could not be restarted: %v", restartErr)
			}
		}
		return softenAutomatedRefreshErr(err, autoConfirm)
	}
	// The flow stopped the daemon: it must run again whether or not the config changed.
	return deliverReload(true)
}

// softenAutomatedRefreshErr turns "all-manual config, nothing matched the catalog" into a no-op for
// automated callers (--yes: pacman/apt hooks, boot refresh unit): a package transaction or boot
// must not fail because no section is catalog-managed. Interactive callers still see the error.
func softenAutomatedRefreshErr(err error, autoConfirm bool) error {
	if err != nil && autoConfirm && errors.Is(err, errNoCatalogMatch) {
		log.Infof("catalog refresh: %v — nothing to do", err)
		return nil
	}
	return err
}

// applyLiveRefresh delivers a patched config to the running daemon: SIGHUP reload (atomic), restart
// fallback. A result matching the running config is a no-op. Delivery is the mandatory last step of
// live mode (an earlier version patched the file but never delivered it, so the daemon kept the old
// whitelist silently).
func applyLiveRefresh(changed bool) error {
	if !changed {
		log.Info("nothing to reload: the refreshed config matches the running daemon")
		return nil
	}
	return deliverReload(true)
}

// prepareInstallation deploys the binary if missing (copying the running executable, never
// recompiling) and ensures the master key exists.
func prepareInstallation() error {
	if err := ensureInstalledBinary(); err != nil {
		return err
	}
	return ensureMasterKey()
}

// installedBinaryPath is the service path ensureInstalledBinary manages; a package var so tests can
// use a temp file.
var installedBinaryPath = systemd.InstallBinaryPath

// checkRunningBinaryMatchesInstalled fails fast, before the wizard, when this executable isn't the
// one already deployed at installedBinaryPath (inode identity, as the daemon's self-guard checks;
// see ensureInstalledBinary and guard.WithSelfAllowBinary).
//
// Why: with a binary already installed the wizard leaves it untouched and deploy() restarts the
// daemon on that SAME binary; every self-guard decision (config, fscrypt key, edit-auth hash)
// trusts the exact on-disk file the daemon exec'd, never a same-name/same-content binary elsewhere
// (that would bypass the "identity is dev:ino, never path/name" invariant). Running the wizard from
// a freshly rebuilt standalone binary is an easy-to-hit mismatch that used to surface only at the
// end as "operation not permitted" writing the password hash, after backups and migration had run.
// That write now happens before first start, but the mismatch is still a footgun for other
// self-guarded writes. os.SameFile compares dev:ino like the guard.
func checkRunningBinaryMatchesInstalled() error {
	installedInfo, err := os.Stat(installedBinaryPath)
	if err != nil {
		return nil //nolint:nilerr // nothing installed yet: ensureInstalledBinary deploys this exact binary next
	}

	self, err := os.Executable()
	if err != nil {
		return nil //nolint:nilerr // best-effort: let the wizard proceed and surface any real error itself
	}
	if resolved, resolveErr := filepath.EvalSymlinks(self); resolveErr == nil {
		self = resolved
	}
	selfInfo, err := os.Stat(self)
	if err != nil {
		return nil //nolint:nilerr // same reasoning
	}

	if os.SameFile(selfInfo, installedInfo) {
		return nil
	}
	return fmt.Errorf(
		"this binary (%s) is not the one already installed at %s — `install` leaves an existing binary "+
			"untouched (see `install --binary-only` / `app-listener update`), so the daemon that comes back "+
			"up during this wizard would still be running the OLD binary. Its self-protection only trusts "+
			"its own on-disk file, so the wizard's last step (writing the edit-protected password hash) "+
			"would fail with a confusing \"operation not permitted\" after everything else already ran.\n"+
			"Run 'sudo %s install --binary-only' (or 'app-listener update') first to deploy this binary and "+
			"restart the daemon on it, then re-run 'sudo app-listener install'",
		self, installedBinaryPath, self)
}

// ensureInstalledBinary makes sure the binary is deployed at its service path. An existing one is
// left untouched (upgrades go through `install --binary-only` or `app-listener update`, which
// restart the daemon). A missing one (first run after `make build` or the one-line installer) is
// filled by copying the running executable as-is: no recompilation, no separate --binary-only step.
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

// secureResources verifies encryption state (fatal: encrypted dir declared need_encryption: false,
// or master-key mismatch), asks only unencrypted need_encryption: true dirs and migrates them; the
// text reflects decisions.
func secureResources(vault *fscrypt.Vault, cfgText string, cfg *daemonconfig.Config) (string, error) {
	if err := verifyEncryptionState(vault, cfg); err != nil {
		return "", err
	}
	if err := resolveFilesystemPrereqs(vault, cfg); err != nil {
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

// deploy installs services/hook, copies binary+config, enables the daemon. Existing files are
// diffed: identical stay, differing ones show the diff and ask; a changed config reaches a running
// daemon via SIGHUP. sshUsers (asked before encryption) get the ssh-agent unit. editPassword (maybe "") is
// persisted via writeEditPasswordHash right before EnableAndVerify's first start on a fresh
// install: everything that could abort has succeeded and the daemon isn't running, so that first
// start already self-guards the hash and opens the control socket (no follow-up reload).
func deploy(cfgText string, sshUsers []inst.User, bunUsers []bunUser, editPassword string) error {
	if err := installServices(sshUsers); err != nil {
		return err
	}
	configChanged, err := installConfig(cfgText)
	if err != nil {
		return err
	}
	if err := preflightDeployedBinary(); err != nil {
		return err
	}
	// Before the daemon (re)starts: the reserved .bun-* pattern activates only once each opted-in
	// user's tmp dir exists.
	if err := setupBunTmpdir(bunUsers); err != nil {
		return err
	}
	if err := writeEditPasswordHash(editPassword); err != nil {
		return err
	}
	if err := systemd.EnableAndVerify(configChanged); err != nil {
		return err
	}
	return systemd.EnableCatalogRefresh()
}

// preflightDeployedBinary runs `<installed binary> daemon --check` on the just-deployed binary and
// config BEFORE enabling the service: BPF-LSM active and every guard eBPF program accepted by this
// kernel's verifier (nothing attached). A failure means the daemon would crash-loop (or panic, with
// a prebuilt object too complex for a newer kernel), so the install aborts having changed only the
// on-disk binary/config.
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

// buildBinaryIfNeeded compiles build/linux/app-listener if missing; refuses without a source tree.
// Only the --binary-only path uses it (the full wizard never builds).
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

// ensureMasterKey keeps an existing master key (with a warning), else creates one.
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

// verifyEncryptionState checks every resource: encrypted while need_encryption: false is fatal
// (unmanaged encryption), and encrypted dirs must unlock with the current master key. Grouped
// sections are verified once per encryption root (the policy is per vault root).
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

// encryptDirectories migrates the selected dirs behind a progress bar, fatal on the first error
// (prior backups kept).
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

// cleanupBackups asks per backup whether to delete it (final step, post-verify). The backup is made
// at the encryption root (the whole vault is renamed aside), so grouped sections are handled once,
// at that root.
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
