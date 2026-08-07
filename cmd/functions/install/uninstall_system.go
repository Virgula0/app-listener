package install

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// revertSystemFiles removes every file the installer deployed: the systemd
// daemon unit, the pacman reload hook, the binary at /usr/local/sbin and the
// config at /etc/app-listener/daemon.conf. The daemon unit is disabled first
// (best effort — it may already be disabled), and systemd is reloaded so the
// removal is visible to the manager. The /etc/app-listener directory itself
// is kept: the fscrypt master key lives there and is removed only by
// removeMasterKey.
func revertSystemFiles() error {
	if err := runCmd("systemctl", "disable", daemonServiceName); err != nil {
		log.Warnf("systemctl disable %s failed: %v", daemonServiceName, err)
	}

	removed := 0
	for _, path := range []string{
		filepath.Join(systemdDir, daemonServiceName+".service"),
		filepath.Join(pacmanHooksDir, "50-app-listener-reload.hook"),
		installBinaryPath,
		systemConfigPath,
	} {
		if _, err := os.Lstat(path); os.IsNotExist(err) {
			continue
		}
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("removing %s: %w", path, err)
		}
		log.Infof("removed %s", path)
		removed++
	}

	// A stale pid file from a crash that never cleaned it up.
	if err := os.Remove("/run/app-listener-daemon.pid"); err == nil {
		log.Info("removed stale /run/app-listener-daemon.pid")
	}

	if err := runCmd("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("systemctl daemon-reload: %w", err)
	}
	if removed == 0 {
		log.Info("no installed daemon files found")
	}
	log.Infof("reverted %d installed daemon file(s)", removed)
	return nil
}

// userSSHAgentUnit is one installed per-user ssh-agent unit (owned by the
// user it was installed for).
type userSSHAgentUnit struct {
	User inst.User
	Path string
}

// detectSSHAgentUnits finds the per-user ssh-agent systemd units the
// installer deployed: the unit file at ~/.config/systemd/user/ssh-agent.service
// whose content matches the bundled sample. Root is skipped (the installer
// never installed one for root), and a unit whose content differs from the
// bundled sample is not ours — the user's own unit is never touched.
func detectSSHAgentUnits() ([]userSSHAgentUnit, error) {
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return nil, err
	}
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	var units []userSSHAgentUnit
	for i := range users {
		u := &users[i]
		if u.UID == 0 {
			continue
		}
		path := filepath.Join(u.Home, ".config", "systemd", "user", "ssh-agent.service")
		if isInstallerSSHAgentUnit(sample, path) {
			units = append(units, userSSHAgentUnit{User: *u, Path: path})
		}
	}
	return units, nil
}

// isInstallerSSHAgentUnit reports whether the ssh-agent unit at path was
// installed by the installer, i.e. its content matches the bundled sample.
// A missing file or a modified unit is not ours; the user's own unit is
// never touched.
func isInstallerSSHAgentUnit(sample []byte, path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return bytes.Equal(data, sample)
}

// revertSSHAgents reverts the per-user ssh-agent units installed by the
// installer. A single TUI confirmation (default: no) precedes the removal
// of every detected unit and its default.target.wants symlink.
func revertSSHAgents() error {
	units, err := detectSSHAgentUnits()
	if err != nil {
		return err
	}
	if len(units) == 0 {
		log.Info("no installer-provided ssh-agent units found")
		return nil
	}

	ok, err := confirmOnce(
		fmt.Sprintf("Remove the per-user ssh-agent units installed by the installer for %d user(s)?", len(units)),
		"Remove units")
	if err != nil {
		return err
	}
	if !ok {
		log.Info("ssh-agent units kept")
		return nil
	}
	for _, u := range units {
		if err := removeSSHAgentUnit(u); err != nil {
			return err
		}
	}
	log.Infof("reverted %d ssh-agent unit(s)", len(units))
	return nil
}

// removeSSHAgentUnit reverts one per-user ssh-agent unit: the
// default.target.wants symlink (created by the installer instead of
// `systemctl --user enable`) and the unit file itself.
func removeSSHAgentUnit(u userSSHAgentUnit) error {
	wants := filepath.Join(filepath.Dir(u.Path), "default.target.wants", "ssh-agent.service")
	if err := os.Remove(wants); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing symlink %s: %w", wants, err)
	}
	if err := os.Remove(u.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing unit %s: %w", u.Path, err)
	}
	log.Infof("reverted ssh-agent unit for %s", u.User.Name)
	return nil
}

// removeMasterKey deletes the fscrypt master key and, when the parent
// directory becomes empty, the directory itself. The key must only be
// deleted after every intended decryption completed and only when
// --delete-key was passed.
func removeMasterKey() error {
	return removeKeyAndEmptyDir(fscrypt.MasterKeyFile, systemConfigDir)
}

// removeKeyAndEmptyDir deletes the key file and, when dir becomes empty,
// dir itself. A key that is already gone is a success (nothing left to
// clean); a dir that is not empty (e.g. a stray file the operator keeps) is
// left alone.
func removeKeyAndEmptyDir(keyFile, dir string) error {
	if err := os.Remove(keyFile); err != nil {
		if os.IsNotExist(err) {
			log.Info("master key already gone")
			return nil
		}
		return fmt.Errorf("removing master key %s: %w", keyFile, err)
	}
	log.Infof("removed master key %s", keyFile)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // directory already gone: nothing left to clean
		}
		return fmt.Errorf("reading %s: %w", dir, err)
	}
	if len(entries) == 0 {
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("removing empty %s: %w", dir, err)
		}
		log.Infof("removed empty %s", dir)
	}
	return nil
}
