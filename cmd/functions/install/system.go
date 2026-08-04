package install

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	inst "github.com/Virgula0/app-listener/internal/install"
)

const (
	// buildBinaryPath is the Makefile output path for the linux build.
	buildBinaryPath = "build/linux/app-listener"
	// installBinaryPath is where the service unit expects the binary.
	installBinaryPath = "/usr/local/sbin/app-listener"
	// systemConfigDir and systemConfigPath hold the installed daemon config.
	systemConfigDir  = "/etc/app-listener"
	systemConfigPath = "/etc/app-listener/daemon.conf"
	// systemdDir is the system unit directory.
	systemdDir = "/etc/systemd/system"
	// pacmanHooksDir is Arch's post-transaction hook directory.
	pacmanHooksDir = "/etc/pacman.d/hooks"
	// daemonServiceName is the systemd unit name (without .service).
	daemonServiceName = "app-listener-daemon"
)

const (
	daemonEnabledState = "enabled"
	daemonActiveState  = "active"
)

// selectedUsers remembers the users picked in the selection step so the
// per-user ssh-agent unit can be installed during deployment.
var selectedUsers []inst.User

// mustCwd returns the current working directory (empty on failure).
func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// runCmd runs a command with CGO/GOOS/GOARCH pinned to the Makefile build
// environment and returns a wrapped error including the combined output.
func runCmd(prog string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), prog, args...)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1", "GOOS=linux", "GOARCH=amd64")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", prog, strings.Join(args, " "), err, out)
	}
	return nil
}

// installServices copies the embedded service units and pacman hook from
// daemon-samples into place, skipping anything that already exists.
// ssh-agent.service is installed per selected user.
func installServices() error {
	files, err := inst.SampleFiles()
	if err != nil {
		return err
	}
	for _, name := range files {
		switch {
		case strings.HasSuffix(name, ".hook"):
			if _, err := os.Stat("/etc/pacman.d"); err != nil {
				log.Warnf("no /etc/pacman.d directory: skipping pacman hook %s", name)
				continue
			}
			if err := installFile(name, filepath.Join(pacmanHooksDir, name), 0o644); err != nil {
				return err
			}
		case name == "ssh-agent.service":
			for _, u := range selectedUsers {
				if err := installSSHAgent(u); err != nil {
					return err
				}
			}
		case strings.HasSuffix(name, ".service"):
			if err := installFile(name, filepath.Join(systemdDir, name), 0o644); err != nil {
				return err
			}
		default:
			// daemon.conf and anything else: not installed as a system file.
			log.Debugf("not installing %s", name)
		}
	}
	return nil
}

// installFile writes one embedded sample to dest unless it already exists.
func installFile(name, dest string, mode os.FileMode) error {
	data, err := inst.SampleContent(name)
	if err != nil {
		return err
	}
	return upsertFile(dest, dest, data, mode, -1)
}

// upsertFile writes data to path unless it already matches. An existing
// file with different content is diffed and the user is asked whether to
// overwrite it; identical files are skipped silently. When uid is
// non-negative the file is chowned afterwards (used for per-user units).
func upsertFile(path, label string, data []byte, mode os.FileMode, uid int) error {
	existing, statErr := os.ReadFile(path)
	newFile := false
	switch {
	case statErr == nil && bytes.Equal(existing, data):
		log.Infof("%s already up to date: skipping", label)
		return nil
	case statErr == nil:
		overwrite, err := inst.ConfirmOverwrite(path, existing, data)
		if err != nil {
			return err
		}
		if !overwrite {
			log.Warnf("keeping existing %s", label)
			return nil
		}
		log.Infof("overwrote %s", label)
	case !os.IsNotExist(statErr):
		return fmt.Errorf("reading %s: %w", label, statErr)
	default:
		newFile = true
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return fmt.Errorf("installing %s: %w", label, err)
	}
	if newFile {
		log.Infof("installed %s", label)
	}
	if uid >= 0 {
		if err := os.Chown(path, uid, -1); err != nil {
			return fmt.Errorf("chown %s: %w", label, err)
		}
	}
	return nil
}

// installSSHAgent installs the ssh-agent user unit for one user, hands it
// to the user and enables it for that user's session. An existing unit is
// compared with the bundled one: identical units are skipped, differing
// ones show a diff and ask whether to overwrite. Root is skipped: it
// normally has no interactive user session.
func installSSHAgent(u inst.User) error {
	if u.UID == 0 {
		log.Debugf("skipping ssh-agent unit for root (no interactive user session)")
		return nil
	}
	unitDir := filepath.Join(u.Home, ".config", "systemd", "user")
	unitPath := filepath.Join(unitDir, "ssh-agent.service")
	data, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(unitDir, 0o755); err != nil {
		return err
	}
	if err := upsertFile(unitPath, fmt.Sprintf("ssh-agent unit for %s", u.Name), data, 0o600, int(u.UID)); err != nil {
		return err
	}
	wantsDir := filepath.Join(unitDir, "default.target.wants")
	if err := os.MkdirAll(wantsDir, 0o755); err != nil {
		return err
	}
	// Equivalent of `systemctl --user enable ssh-agent` without requiring
	// the user's session bus: a symlink in default.target.wants.
	if err := os.Symlink(unitPath, filepath.Join(wantsDir, "ssh-agent.service")); err != nil && !os.IsExist(err) {
		return err
	}
	log.Infof("ssh-agent enabled for %s (relogin or start manually: systemctl --user start ssh-agent)", u.Name)
	return nil
}

// installBinaryAndConfig copies the freshly built binary to the service
// path (always, so upgrades are deployed) and writes the final config.
// An existing config with different content is diffed and the user is
// asked whether to overwrite it; identical configs are left alone. It
// reports whether the installed config differs from what the daemon was
// running with.
func installBinaryAndConfig(cfgText string) (configChanged bool, err error) {
	if err := inst.CopyTree(buildBinaryPath, installBinaryPath); err != nil {
		return false, fmt.Errorf("installing binary: %w", err)
	}
	if err := os.Chmod(installBinaryPath, 0o700); err != nil {
		return false, err
	}
	log.Infof("installed binary at %s", installBinaryPath)

	if err := os.MkdirAll(systemConfigDir, 0o700); err != nil {
		return false, err
	}
	desired := []byte(cfgText)
	existing, statErr := os.ReadFile(systemConfigPath)
	if statErr == nil {
		if bytes.Equal(existing, desired) {
			log.Infof("config already up to date: %s", systemConfigPath)
			return false, nil
		}
		overwrite, err := inst.ConfirmOverwrite(systemConfigPath, existing, desired)
		if err != nil {
			return false, err
		}
		if !overwrite {
			log.Warnf("keeping existing config %s: the daemon keeps running with it", systemConfigPath)
			return false, nil
		}
		if err := os.WriteFile(systemConfigPath, desired, 0o600); err != nil {
			return false, fmt.Errorf("installing config: %w", err)
		}
		log.Infof("overwrote config at %s", systemConfigPath)
		return true, nil
	}
	if !os.IsNotExist(statErr) {
		return false, fmt.Errorf("reading %s: %w", systemConfigPath, statErr)
	}
	if err := os.WriteFile(systemConfigPath, desired, 0o600); err != nil {
		return false, fmt.Errorf("installing config: %w", err)
	}
	log.Infof("installed config at %s", systemConfigPath)
	return true, nil
}

// enableAndVerify brings the daemon to the desired state no matter how a
// previous or interrupted installation left it: when the config changed
// and the daemon is already running, it is reloaded with SIGHUP instead of
// restarted; when it is not running it is started; when it is not enabled
// across reboots it is enabled. Both states are verified before success.
func enableAndVerify(configChanged bool) error {
	if err := runCmd("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("systemctl daemon-reload: %w", err)
	}

	active := strings.TrimSpace(systemctlOutput("is-active", daemonServiceName))
	switch {
	case active == daemonActiveState && configChanged:
		log.Infof("daemon is running with a changed config: reloading (SIGHUP) ...")
		if err := runCmd("systemctl", "reload", daemonServiceName); err != nil {
			log.Warnf("systemctl reload failed (%v): restarting instead", err)
			if err := runCmd("systemctl", "restart", daemonServiceName); err != nil {
				return fmt.Errorf("systemctl restart %s: %w", daemonServiceName, err)
			}
		}
		log.Infof("daemon reloaded with the new config")
	case active != daemonActiveState:
		log.Infof("daemon is not running: starting ...")
		if err := runCmd("systemctl", "start", daemonServiceName); err != nil {
			return fmt.Errorf("systemctl start %s: %w — inspect with: journalctl -u %s -e", daemonServiceName, err, daemonServiceName)
		}
	}

	enabled := strings.TrimSpace(systemctlOutput("is-enabled", daemonServiceName))
	if enabled != daemonEnabledState {
		log.Infof("daemon is %q: enabling across reboots ...", enabled)
		if err := runCmd("systemctl", "enable", daemonServiceName); err != nil {
			return fmt.Errorf("systemctl enable %s: %w", daemonServiceName, err)
		}
	}
	enabled = strings.TrimSpace(systemctlOutput("is-enabled", daemonServiceName))
	if enabled != daemonEnabledState {
		return fmt.Errorf("daemon is not enabled across reboots (is-enabled: %q)", enabled)
	}
	log.Infof("daemon is enabled across reboots (%s)", enabled)

	active = strings.TrimSpace(systemctlOutput("is-active", daemonServiceName))
	if active != daemonActiveState {
		return fmt.Errorf("daemon is not running (is-active: %q) — inspect with: journalctl -u %s -e", active, daemonServiceName)
	}
	log.Infof("daemon is running (%s)", active)
	return nil
}

// systemctlOutput runs a read-only systemctl query, returning its trimmed
// output (empty on failure).
func systemctlOutput(args ...string) string {
	cmd := exec.CommandContext(context.Background(), "systemctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}
