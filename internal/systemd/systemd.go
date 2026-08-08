// Package systemd centralizes the systemd integration and the installed
// binary lifecycle shared by the install, uninstall and update commands:
// where the daemon files live, how the daemon service is stopped, enabled,
// started and verified, and how the /usr/local/bin PATH symlink is managed.
package systemd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/protected"
)

const (
	// InstallBinaryPath is where the service unit expects the binary.
	InstallBinaryPath = "/usr/local/sbin/app-listener"
	// BinSymlinkPath makes the installed binary reachable as `app-listener`
	// from any terminal: /usr/local/bin is in the default PATH of every
	// user, while /usr/local/sbin normally is not.
	BinSymlinkPath = "/usr/local/bin/app-listener"
	// SystemConfigDir and SystemConfigPath hold the installed daemon config.
	SystemConfigDir  = "/etc/app-listener"
	SystemConfigPath = "/etc/app-listener/daemon.conf"
	// SystemdDir is the system unit directory.
	SystemdDir = "/etc/systemd/system"
	// PacmanHooksDir is Arch's post-transaction hook directory.
	PacmanHooksDir = "/etc/pacman.d/hooks"
	// DaemonServiceName is the systemd unit name (without .service).
	DaemonServiceName = "app-listener-daemon"
)

const (
	daemonEnabledState = "enabled"
	daemonActiveState  = "active"
)

// StopDaemonIfRunning makes sure no app-listener daemon is left alive while
// the installation reconfigures it: an active systemd unit is stopped
// cleanly (the caller re-enables and starts it at the end), a daemon process
// started outside systemd (e.g. a manual --headless run) is a fatal error
// because it cannot be controlled, and a daemon that is not running at all
// needs no action. After a systemd stop, the installed config's resources
// are verified to be locked: a daemon killed without its lockdown would
// leave the watched directories unlocked and unguarded.
func StopDaemonIfRunning() error {
	if strings.TrimSpace(SystemctlOutput("is-active", DaemonServiceName)) == daemonActiveState {
		log.Infof("daemon is running under systemd: stopping it ...")
		if err := RunCmd("systemctl", "stop", DaemonServiceName); err != nil {
			return fmt.Errorf("stopping %s: %w", DaemonServiceName, err)
		}
		log.Infof("daemon stopped")
		if err := VerifyInstalledResourcesLocked(); err != nil {
			return err
		}
		return nil
	}
	pids, err := protected.FindDaemonProcesses("/proc")
	if err != nil {
		return fmt.Errorf("scanning for a running daemon process: %w", err)
	}
	if len(pids) > 0 {
		return fmt.Errorf("fatal: the daemon is running outside systemd (pid(s) %v) — stop it first, e.g. kill %d", pids, pids[0])
	}
	log.Info("daemon is not running: nothing to stop")
	return nil
}

// VerifyInstalledResourcesLocked loads the installed daemon config — the one
// the stopped daemon was running with — and checks that every encrypted
// resource is keyless. The daemon's own shutdown only returns once all
// resources are locked, so a violation here means the daemon was killed
// without running its lockdown (crash, SIGKILL, older binary): continuing
// would leave the trees unlocked and unguarded.
func VerifyInstalledResourcesLocked() error {
	if _, err := os.Stat(SystemConfigPath); os.IsNotExist(err) {
		log.Debug("no installed config yet: nothing to verify")
		return nil
	}
	cfg, err := daemonconfig.Load(SystemConfigPath)
	if err != nil {
		return fmt.Errorf("cannot verify the stopped daemon's lockdown: parsing %s: %w — stop any stray daemon process and close files in the watched directories first", SystemConfigPath, err)
	}
	return VerifyResourcesLocked(cfg.Resources, fscrypt.New().IsProvisioned)
}

// VerifyResourcesLocked reports an error when any encrypted resource is
// still provisioned (unlocked).
func VerifyResourcesLocked(resources []daemonconfig.Resource, provisioned func(string) (bool, error)) error {
	var unlocked []string
	for _, r := range resources {
		if !r.NeedEncryption {
			continue
		}
		provisionedNow, err := provisioned(r.Path)
		if err != nil {
			return fmt.Errorf("checking lock state of %s: %w", r.Path, err)
		}
		if provisionedNow {
			unlocked = append(unlocked, r.Path)
		}
	}
	if len(unlocked) > 0 {
		return fmt.Errorf("fatal: the daemon exited with %d directory(ies) still unlocked: %v — its lockdown did not run or could not finish. Close every file in these directories (lsof +D <dir>, fuser -v <dir>) and re-run the installer",
			len(unlocked), unlocked)
	}
	return nil
}

// RunCmd runs a command and returns a wrapped error including the combined
// output.
func RunCmd(prog string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), prog, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", prog, strings.Join(args, " "), err, out)
	}
	return nil
}

// SystemctlOutput runs a read-only systemctl query, returning its trimmed
// output (empty on failure).
func SystemctlOutput(args ...string) string {
	cmd := exec.CommandContext(context.Background(), "systemctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// EnableAndVerify brings the daemon to the desired state no matter how a
// previous or interrupted installation left it: when the config changed
// and the daemon is already running, it is reloaded with SIGHUP instead of
// restarted; when it is not running it is started; when it is not enabled
// across reboots it is enabled. Both states are verified before success.
func EnableAndVerify(configChanged bool) error {
	if err := RunCmd("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("systemctl daemon-reload: %w", err)
	}

	active := strings.TrimSpace(SystemctlOutput("is-active", DaemonServiceName))
	switch {
	case active == daemonActiveState && configChanged:
		log.Infof("daemon is running with a changed config: reloading (SIGHUP) ...")
		if err := RunCmd("systemctl", "reload", DaemonServiceName); err != nil {
			log.Warnf("systemctl reload failed (%v): restarting instead", err)
			if err := RunCmd("systemctl", "restart", DaemonServiceName); err != nil {
				return fmt.Errorf("systemctl restart %s: %w", DaemonServiceName, err)
			}
		}
		log.Infof("daemon reloaded with the new config")
	case active != daemonActiveState:
		log.Infof("daemon is not running: starting ...")
		if err := RunCmd("systemctl", "start", DaemonServiceName); err != nil {
			return fmt.Errorf("systemctl start %s: %w — inspect with: journalctl -u %s -e", DaemonServiceName, err, DaemonServiceName)
		}
	}

	enabled := strings.TrimSpace(SystemctlOutput("is-enabled", DaemonServiceName))
	if enabled != daemonEnabledState {
		log.Infof("daemon is %q: enabling across reboots ...", enabled)
		if err := RunCmd("systemctl", "enable", DaemonServiceName); err != nil {
			return fmt.Errorf("systemctl enable %s: %w", DaemonServiceName, err)
		}
	}
	enabled = strings.TrimSpace(SystemctlOutput("is-enabled", DaemonServiceName))
	if enabled != daemonEnabledState {
		return fmt.Errorf("daemon is not enabled across reboots (is-enabled: %q)", enabled)
	}
	log.Infof("daemon is enabled across reboots (%s)", enabled)

	active = strings.TrimSpace(SystemctlOutput("is-active", DaemonServiceName))
	if active != daemonActiveState {
		return fmt.Errorf("daemon is not running (is-active: %q) — inspect with: journalctl -u %s -e", active, DaemonServiceName)
	}
	log.Infof("daemon is running (%s)", active)
	return nil
}

// EnsureBinSymlink makes the installed binary reachable as `app-listener`
// from any terminal by symlinking /usr/local/bin/app-listener to the
// installed binary.
func EnsureBinSymlink() error {
	return EnsureBinSymlinkAt(BinSymlinkPath, InstallBinaryPath)
}

// EnsureBinSymlinkAt creates linkPath -> targetPath. It is idempotent: an
// existing symlink already pointing at targetPath is left alone, a regular
// file or a foreign symlink is never touched (a warning is logged instead),
// and a missing parent directory is only a warning — the daemon itself does
// not need the symlink.
func EnsureBinSymlinkAt(linkPath, targetPath string) error {
	info, err := os.Lstat(linkPath)
	switch {
	case err == nil && info.Mode()&os.ModeSymlink == 0:
		log.Warnf("refusing to replace %s: it is a regular file", linkPath)
		return nil
	case err == nil:
		target, readErr := os.Readlink(linkPath)
		if readErr != nil {
			return fmt.Errorf("reading symlink %s: %w", linkPath, readErr)
		}
		if target == targetPath {
			log.Infof("symlink %s already points at %s: keeping it", linkPath, targetPath)
			return nil
		}
		log.Warnf("refusing to replace %s: it already points at %s (not %s)", linkPath, target, targetPath)
		return nil
	case !os.IsNotExist(err):
		return fmt.Errorf("reading %s: %w", linkPath, err)
	}
	if _, statErr := os.Stat(filepath.Dir(linkPath)); statErr != nil {
		log.Warnf("cannot create the %s symlink: %s does not exist (%v)", linkPath, filepath.Dir(linkPath), statErr)
		return nil
	}
	if err := os.Symlink(targetPath, linkPath); err != nil {
		return fmt.Errorf("creating symlink %s -> %s: %w", linkPath, targetPath, err)
	}
	log.Infof("created symlink %s -> %s", linkPath, targetPath)
	return nil
}

// RemoveBinSymlink deletes the PATH symlink when it points at the installed
// binary. A missing symlink, a regular file or a foreign symlink are left
// untouched.
func RemoveBinSymlink() {
	RemoveBinSymlinkAt(BinSymlinkPath, InstallBinaryPath)
}

// RemoveBinSymlinkAt removes linkPath only when it is a symlink pointing at
// targetPath. A missing symlink, a regular file or a foreign symlink are
// left untouched.
func RemoveBinSymlinkAt(linkPath, targetPath string) {
	target, err := os.Readlink(linkPath)
	if err != nil {
		return
	}
	if target != targetPath {
		log.Warnf("not removing %s: it points at %s, not %s", linkPath, target, targetPath)
		return
	}
	if err := os.Remove(linkPath); err != nil {
		log.Warnf("removing symlink %s: %v", linkPath, err)
		return
	}
	log.Infof("removed symlink %s", linkPath)
}
