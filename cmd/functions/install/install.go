// Package install implements the `app-listener install` wizard: deploys the binary if missing
// (copies the running executable, never recompiles), ensures the fscrypt master key, picks users
// and dirs, migrates to fscrypt with backups, installs units/hook, enables the daemon.
package install

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the daemon: pick directories, encrypt them, enable the service",
	Long: `Interactive installer for the daemon. It deploys the binary to
/usr/local/sbin/app-listener, creates the fscrypt master key, lets you pick the
directories to protect (a catalog of SSH, GPG, AI agents, browsers, ...) and
review daemon.conf, encrypts them (keeping a .app_listener.backup copy), asks
for the edit-protected password, installs and starts the systemd service.

Re-run it to resume an interrupted install. A completed install is recorded and
the wizard then refuses: use the subcommands, 'update', or 'uninstall' first.`,
	Example: `  sudo app-listener install
  sudo app-listener install binary
  sudo app-listener install catalog diff
  sudo app-listener install facultatives`,
	Args: cobra.NoArgs,
	RunE: runInstall,
}

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Deploy only this build's binary, then offer to restart the daemon",
	Long: `Replace /usr/local/sbin/app-listener with build/linux/app-listener (built when
missing) and recreate the /usr/local/bin symlink. Nothing else changes. With the
service installed, a differing daemon unit is diffed and offered, then you pick
whether to restart the daemon now or leave it on the old binary.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error { return asRoot(installBinaryOnly) },
}

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Sync an existing install with the catalog of sensitive directories",
	Args:  cobra.NoArgs,
}

var catalogRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Re-expand the catalog whitelists of every guarded directory",
	Long: `Re-expand the catalog's whitelist of every guarded directory (new binaries in,
deleted ones out) and show the diff before writing daemon.conf. Sections not
from the catalog are kept as they are. The daemon does this by itself; this is
for manual or debug use. Without --live the daemon is stopped for the cycle.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error {
		return asRoot(func() error { return runRefreshCatalogOnly(refreshYes, refreshLive) })
	},
}

var catalogDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Add catalog directories the installed daemon.conf doesn't guard yet",
	Long: `List the catalog directories that exist on this host but aren't guarded yet,
append the ones you pick to daemon.conf and encrypt them like a fresh install.
Existing sections stay byte-for-byte intact. The daemon must be stopped first;
at the end you are asked whether to start it. Run 'catalog refresh' too to
update existing whitelists.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error { return asRoot(runDiffCatalog) },
}

var backupsCmd = &cobra.Command{
	Use:   "backups",
	Short: "Restore or delete the .app_listener.backup copies the encryption left",
	Args:  cobra.NoArgs,
}

var backupsRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Undo the encryption: move the backups back over the encrypted directories",
	Long: `Delete the encrypted directories and move their .app_listener.backup copies
back, after a confirmation naming each one. Refused while the daemon runs.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error { return asRoot(restoreBackups) },
}

var backupsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the backups (the daemon may keep running)",
	Args:  cobra.NoArgs,
	RunE:  func(*cobra.Command, []string) error { return asRoot(deletePostBackups) },
}

var (
	refreshYes   bool
	refreshLive  bool
	markComplete bool
)

func init() {
	InstallCmd.Flags().BoolVar(&allowMetadataOutput, "allow-metadata-output", false,
		"Install the daemon unit without --no-log-metadata-blocks (logs denied op=PTRACE mode=READ inspections)")
	// The installer runs it on the installed binary, which may be another version: never rename it.
	// The GUI window isn't the app-listener binary, so the daemon's guard on /etc/app-listener denies
	// it the completed-install record; the installed binary writes it.
	InstallCmd.Flags().BoolVar(&markComplete, "mark-complete", false,
		"Record a completed installation (the desktop installer runs it on the installed binary)")

	catalogRefreshCmd.Flags().BoolVar(&refreshLive, "live", false,
		"Keep the daemon running: apply the new config with a reload (needs the daemon running)")
	catalogRefreshCmd.Flags().BoolVarP(&refreshYes, "yes", "y", false, "Write without asking (scripts)")

	catalogCmd.AddCommand(catalogRefreshCmd, catalogDiffCmd)
	backupsCmd.AddCommand(backupsRestoreCmd, backupsDeleteCmd)
	InstallCmd.AddCommand(binaryCmd, catalogCmd, backupsCmd, facultativesCmd)
}

func asRoot(run func() error) error {
	if os.Geteuid() != 0 {
		return errors.New("install must be run as root: sudo app-listener install")
	}
	return run()
}

func runInstall(*cobra.Command, []string) error {
	if markComplete {
		return asRoot(markInstallComplete)
	}
	return asRoot(runWizard)
}

// runWizard is the installer body (docs/installation.md#the-install-wizard).
func runWizard() error {
	if err := protected.RequireNoInstallCert(); err != nil {
		return err
	}
	if err := checkRunningBinaryMatchesInstalled(); err != nil {
		return err
	}

	// Stop the daemon first: it holds the files being reconfigured and would keep guarding the dirs
	// being moved. Outside-systemd daemons are fatally refused (can't be controlled).
	wizard.Stage("Stopping the daemon")
	if err := systemd.StopDaemonIfRunning(); err != nil {
		return err
	}

	wizard.Stage("Deploying the binary and the master key")
	if err := prepareInstallation(); err != nil {
		return err
	}

	wizard.Stage("Choosing what to protect")
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

	wizard.Stage("Encrypting")
	vault := fscrypt.New()
	cfgText, err = secureResources(vault, cfgText, cfg)
	if err != nil {
		return err
	}

	// Collected now, persisted in deploy() only after every step that could still abort has
	// succeeded, just before the daemon's first start (writeEditPasswordHash).
	editPassword, err := promptEditPassword()
	if err != nil {
		return err
	}

	wizard.Stage("Installing services and starting the daemon")
	if err := deploy(cfgText, sshUsers, bunUsers, editPassword); err != nil {
		return err
	}

	wizard.Stage("Cleaning up")
	if err := cleanOrphanedFscrypt(cfg); err != nil {
		return err
	}
	// Before the backup prompts: keeping or skipping a backup doesn't make the install incomplete.
	certifyInstall()

	return cleanupBackups(cfg)
}

// certifyInstall records the completed install, after which the wizard refuses to run again. Under
// the GUI the installed binary writes it (see --mark-complete). A failure only loses that refusal,
// never protection, so it warns instead of failing a finished install.
func certifyInstall() {
	var err error
	if deploySource == "" {
		err = markInstallComplete()
	} else {
		err = systemd.RunCmd(systemd.InstallBinaryPath, "install", "--mark-complete")
	}
	if err != nil {
		log.Warnf("the installation completed, but recording it in %s failed: %v — the install wizard "+
			"won't refuse a re-run", protected.InstallCertPath(), err)
		return
	}
	log.Infof("installation recorded in %s: run `app-listener uninstall` before installing again",
		protected.InstallCertPath())
}

// markInstallComplete is --mark-complete: only a running, enabled daemon counts as installed.
func markInstallComplete() error {
	if !systemd.IsDaemonActive() {
		return errors.New("the daemon is not running: the installation is not complete")
	}
	if !systemd.IsDaemonEnabled() {
		return errors.New("the daemon is not enabled across reboots: the installation is not complete")
	}
	return protected.WriteInstallCert(constants.Version)
}

// installBinaryOnly deploys only the freshly built binary (systemd.DeployInstalledBinary); config
// and fscrypt are untouched. If the service is installed, a daemon unit differing from the one the
// build ships is diffed and asked, then whether to restart the daemon now or leave it on the old
// binary until a manual restart.
func installBinaryOnly() error {
	src := deploySource
	if src == "" {
		if err := buildBinaryIfNeeded(); err != nil {
			return err
		}
		src = buildBinaryPath
	}

	if !systemd.DaemonUnitInstalled() {
		if err := systemd.DeployInstalledBinary(src, nil); err != nil {
			return fmt.Errorf("installing binary: %w", err)
		}
		log.Info("binary-only install complete: binary deployed; the daemon service is not installed (run `sudo app-listener install`)")
		return nil
	}

	unit, err := systemd.PlanDaemonUnitUpdate(src, inst.ConfirmOverwrite)
	if err != nil {
		return err
	}

	restart, formErr := wizard.Confirm(wizard.Question{
		Title:       "Restart the daemon now?",
		Description: "The daemon service is installed. Restarting picks up the new binary immediately (the current default). Skipping leaves the running daemon on the old binary — restart it yourself later (systemctl restart app-listener-daemon) to run the new one.",
		Affirmative: "Yes, restart now",
		Negative:    "No, skip the restart",
		Default:     true,
	})
	if formErr != nil {
		return formErr
	}

	if !restart {
		if err := systemd.DeployInstalledBinaryNoRestart(src, unit); err != nil {
			return fmt.Errorf("installing binary: %w", err)
		}
		log.Info("binary-only install complete: binary deployed; the daemon restart was skipped — restart it manually (systemctl restart app-listener-daemon) to run the new binary")
		return nil
	}

	if err := systemd.DeployInstalledBinary(src, unit); err != nil {
		return fmt.Errorf("installing binary: %w", err)
	}
	log.Info("binary-only install complete: the daemon is running the new binary")
	return nil
}

// deliverReload applies a patched config to the daemon: SIGHUP reload (atomic; new guards attach
// before old detach) with a restart fallback when running, or a start after the stopped flow. A
// package var so tests can observe delivery.
var deliverReload = systemd.EnableAndVerify

// runRefreshCatalogOnly re-scans the catalog whitelist for every guarded directory in daemon.conf:
// encrypted vaults are unlocked for stat, globs re-expanded, the diff shown for confirmation;
// user-added sections kept verbatim.
//
// Live mode: the daemon keeps running (vaults unlocked, guards attached), so no unlock/lock cycle
// or stop is needed and a changed config is delivered via SIGHUP. A stopped daemon can't support it
// (locked vaults hide plaintext names), so the caller is told to use the stopped flow.
func runRefreshCatalogOnly(autoConfirm, live bool) error {
	if live {
		if !systemd.IsDaemonActive() {
			return errors.New("live catalog refresh requires the daemon to be running: " +
				"locked vaults hide plaintext names — use `install catalog refresh` without --live, or start the daemon first")
		}
		vault := fscrypt.New()
		changed, err := refreshCatalogConfig(vault, autoConfirm, true)
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
	if _, err := refreshCatalogConfig(vault, autoConfirm, false); err != nil {
		// refreshCatalogConfig writes the config only after every section patched, so an error
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
// automated callers (--yes): a script must not fail because no section is catalog-managed.
// Interactive callers still see the error.
func softenAutomatedRefreshErr(err error, autoConfirm bool) error {
	if err != nil && autoConfirm && errors.Is(err, inst.ErrNoCatalogMatch) {
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
// daemon on that SAME binary; every self-guard decision (config, fscrypt key, binaries.db)
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
	if deploySource != "" {
		return checkInstalledMatchesSource()
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
			"untouched (see `install binary` / `app-listener update`), so the daemon that comes back "+
			"up during this wizard would still be running the OLD binary. Its self-protection only trusts "+
			"its own on-disk file, so the wizard's last step (writing the edit-protected password hash) "+
			"would fail with a confusing \"operation not permitted\" after everything else already ran.\n"+
			"Run 'sudo %s install binary' (or 'app-listener update') first to deploy this binary and "+
			"restart the daemon on it, then re-run 'sudo app-listener install'",
		self, installedBinaryPath, self)
}

// checkInstalledMatchesSource is checkRunningBinaryMatchesInstalled for a GUI run: the window is
// never the deployed binary, so the installed one must be byte-identical to the headless build this
// run would deploy, or the daemon that comes back would be another version than the one configured.
func checkInstalledMatchesSource() error {
	same, err := sameContent(installedBinaryPath, deploySource)
	if err != nil {
		return err
	}
	if !same {
		return fmt.Errorf("the binary installed at %s is not this version's headless build: deploy it "+
			"first from a terminal (`app-listener update`, or `install binary`), then run the wizard again", installedBinaryPath)
	}
	return nil
}

func sameContent(a, b string) (bool, error) {
	ha, err := fileSHA256(a)
	if err != nil {
		return false, err
	}
	hb, err := fileSHA256(b)
	if err != nil {
		return false, err
	}
	return ha == hb, nil
}

func fileSHA256(path string) ([sha256.Size]byte, error) {
	var sum [sha256.Size]byte
	f, err := os.Open(path)
	if err != nil {
		return sum, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return sum, fmt.Errorf("hashing %s: %w", path, err)
	}
	h.Sum(sum[:0])
	return sum, nil
}

// ensureInstalledBinary makes sure the binary is deployed at its service path. An existing one is
// left untouched (upgrades go through `install binary` or `app-listener update`, which
// restart the daemon). A missing one (first run after `make build` or the one-line installer) is
// filled by copying the running executable as-is: no recompilation, no separate `install binary` step.
func ensureInstalledBinary() error {
	if _, err := os.Stat(installedBinaryPath); err == nil {
		log.Infof("binary is installed: %s", installedBinaryPath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking %s: %w", installedBinaryPath, err)
	}

	self := deploySource
	if self == "" {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locating the running binary to deploy: %w", err)
		}
		if resolved, resolveErr := filepath.EvalSymlinks(exe); resolveErr == nil {
			exe = resolved
		}
		self = exe
	}
	log.Infof("binary not installed yet: deploying %s -> %s (no recompilation)", self, installedBinaryPath)
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
	text, cfg, err := editConfig(candidates)
	if err == nil {
		inst.WarnGeneralTools(cfg)
	}
	return text, cfg, err
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
// daemon via SIGHUP. sshUsers (asked before encryption) get the ssh-agent unit. editPassword ("" =
// keep the existing one) is persisted via writeEditPasswordHash right before EnableAndVerify's
// first start: everything that could abort has succeeded and the daemon isn't running, so that
// first start already seals the hash and opens the control socket (no follow-up reload).
func deploy(cfgText string, sshUsers []inst.User, bunUsers []bunUser, editPassword string) error {
	configChanged, err := deployFiles(cfgText, sshUsers, bunUsers, editPassword)
	if err != nil {
		return err
	}
	return systemd.EnableAndVerify(configChanged)
}

// deployFiles is deploy without starting the daemon: units, config, eBPF preflight, Bun dirs,
// password hash and ledger rows. configChanged reports whether daemon.conf was rewritten.
func deployFiles(cfgText string, sshUsers []inst.User, bunUsers []bunUser, editPassword string) (bool, error) {
	if err := installServices(sshUsers); err != nil {
		return false, err
	}
	configChanged, err := installConfig(cfgText)
	if err != nil {
		return false, err
	}
	if err := preflightDeployedBinary(); err != nil {
		return false, err
	}
	// Before the daemon (re)starts: the reserved .bun-* pattern activates only once each opted-in
	// user's tmp dir exists.
	if err := setupBunTmpdir(bunUsers); err != nil {
		return false, err
	}
	if err := writeEditPasswordHash(editPassword); err != nil {
		return false, err
	}
	return configChanged, grantNewBinaries(cfgText)
}

// grantNewBinaries lets the daemon's next start record every whitelist line, and every catalog
// maintainer script the config makes it run, that the ledger doesn't know yet (binledger.Grant):
// choosing a section in this wizard is the confirmation. Lines already recorded keep their hash, so
// a binary that changed while the daemon was stopped for this install still needs `trust`.
func grantNewBinaries(cfgText string) error {
	cfg, err := daemonconfig.Parse([]byte(cfgText))
	if err != nil {
		return err
	}
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	l, err := binledger.Open(binledger.DefaultPath)
	if err != nil {
		return fmt.Errorf("opening the binary ledger: %w", err)
	}
	defer l.Close()
	lines := binledger.Lines(cfg)
	for _, m := range inst.Maintainers(cfg, users) {
		lines = append(lines, m.Script)
	}
	for _, line := range lines {
		if _, found, err := l.Lookup(line); err != nil {
			return err
		} else if found {
			continue
		}
		if err := l.Grant(line); err != nil {
			return fmt.Errorf("granting %s in the binary ledger: %w", line, err)
		}
	}
	return nil
}

// preflightDeployedBinary runs `<installed binary> daemon --check` (the spelling every version
// accepts) on the just-deployed binary and
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
// Only `install binary` uses it (the full wizard never builds).
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
		remove, err := wizard.Confirm(wizard.Question{
			Title:       fmt.Sprintf("Remove migration backup %s?", backup),
			Description: "The installation succeeded; the backup is no longer needed. The original unencrypted data was already copied into the encrypted directory.",
			Affirmative: "Remove",
			Negative:    "Keep",
		})
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
