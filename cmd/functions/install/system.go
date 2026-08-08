package install

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

const (
	// buildBinaryPath is the Makefile output path for the linux build.
	buildBinaryPath = "build/linux/app-listener"
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
			if err := installFile(name, filepath.Join(systemd.PacmanHooksDir, name), 0o644); err != nil {
				return err
			}
		case name == "ssh-agent.service":
			for _, u := range selectedUsers {
				if err := installSSHAgent(u); err != nil {
					return err
				}
			}
		case strings.HasSuffix(name, ".service"):
			if err := installFile(name, filepath.Join(systemd.SystemdDir, name), 0o644); err != nil {
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
	if err := inst.CopyTree(buildBinaryPath, systemd.InstallBinaryPath); err != nil {
		return false, fmt.Errorf("installing binary: %w", err)
	}
	if err := os.Chmod(systemd.InstallBinaryPath, 0o700); err != nil {
		return false, err
	}
	log.Infof("installed binary at %s", systemd.InstallBinaryPath)

	if err := systemd.EnsureBinSymlink(); err != nil {
		return false, err
	}

	if err := os.MkdirAll(systemd.SystemConfigDir, 0o700); err != nil {
		return false, err
	}
	desired := []byte(cfgText)
	existing, statErr := os.ReadFile(systemd.SystemConfigPath)
	if statErr == nil {
		if bytes.Equal(existing, desired) {
			log.Infof("config already up to date: %s", systemd.SystemConfigPath)
			return false, nil
		}
		overwrite, err := inst.ConfirmOverwrite(systemd.SystemConfigPath, existing, desired)
		if err != nil {
			return false, err
		}
		if !overwrite {
			log.Warnf("keeping existing config %s: the daemon keeps running with it", systemd.SystemConfigPath)
			return false, nil
		}
		if err := os.WriteFile(systemd.SystemConfigPath, desired, 0o600); err != nil {
			return false, fmt.Errorf("installing config: %w", err)
		}
		log.Infof("overwrote config at %s", systemd.SystemConfigPath)
		return true, nil
	}
	if !os.IsNotExist(statErr) {
		return false, fmt.Errorf("reading %s: %w", systemd.SystemConfigPath, statErr)
	}
	if err := os.WriteFile(systemd.SystemConfigPath, desired, 0o600); err != nil {
		return false, fmt.Errorf("installing config: %w", err)
	}
	log.Infof("installed config at %s", systemd.SystemConfigPath)
	return true, nil
}
