package install

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

const (
	// buildBinaryPath is the Makefile output path for the linux build.
	buildBinaryPath = "build/linux/app-listener"
)

// mustCwd returns the current working directory (empty on failure).
func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// installServices copies the embedded system unit files from daemon-samples
// into place, skipping anything that already exists, offers the per-user
// ssh-agent units, and drops the package-manager catalog-refresh hooks.
func installServices(cfg *daemonconfig.Config) error {
	files, err := inst.SampleFiles()
	if err != nil {
		return err
	}
	for _, name := range files {
		switch {
		case name == systemd.PacmanHookName || name == systemd.AptHookSample:
			// Package-manager hooks are installed together, after the loop,
			// once per detected manager (see installReloadHooks).
			continue
		case name == "ssh-agent.service":
			if err := offerSSHAgentUnits(cfg); err != nil {
				return err
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
	return installReloadHooks()
}

// installReloadHooks drops the catalog-refresh post-transaction hook for every
// package manager present on the host. Package managers without a shipped hook
// (dnf, zypper) are reported: the boot-time app-listener-catalog-refresh unit
// still covers reboots, and `app-listener install --update-catalog-only` can
// be run by hand otherwise.
func installReloadHooks() error {
	managers := systemd.DetectPackageManagers()
	if len(managers) == 0 {
		log.Warn("no package manager detected: the catalog whitelist will not refresh " +
			"automatically after package changes — the boot-time refresh still runs, or run " +
			"`app-listener install --update-catalog-only` manually")
		return nil
	}
	hooked := false
	for _, pm := range managers {
		switch pm {
		case systemd.PkgPacman:
			if err := os.MkdirAll(systemd.PacmanHooksDir, 0o755); err != nil {
				return fmt.Errorf("creating %s: %w", systemd.PacmanHooksDir, err)
			}
			if err := installFile(systemd.PacmanHookName,
				filepath.Join(systemd.PacmanHooksDir, systemd.PacmanHookName), 0o644); err != nil {
				return err
			}
			hooked = true
		case systemd.PkgApt:
			if err := os.MkdirAll(systemd.AptHooksDir, 0o755); err != nil {
				return fmt.Errorf("creating %s: %w", systemd.AptHooksDir, err)
			}
			if err := installFileAs(systemd.AptHookSample,
				filepath.Join(systemd.AptHooksDir, systemd.AptHookName), 0o644); err != nil {
				return err
			}
			hooked = true
		default:
			log.Warnf("%s has no bundled catalog-refresh hook: the boot-time refresh covers "+
				"reboots; run `app-listener install --update-catalog-only` after package changes",
				pm)
		}
	}
	if !hooked {
		log.Warn("no package manager with a bundled catalog-refresh hook was found: relying on " +
			"the boot-time refresh and manual `app-listener install --update-catalog-only`")
	}
	return nil
}

// installFile writes one embedded sample to dest unless it already exists.
func installFile(name, dest string, mode os.FileMode) error {
	return installFileAs(name, dest, mode)
}

// installFileAs writes the embedded sample "sample" to "dest" (whose basename
// may differ from the sample name), unless dest already matches.
func installFileAs(sample, dest string, mode os.FileMode) error {
	data, err := inst.SampleContent(sample)
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

// offerSSHAgentUnits installs the per-user ssh-agent systemd unit — but only
// for a user whose ~/.ssh is guarded by this config, and only after asking.
// The unit is per-user (it lives under ~/.config/systemd/user and runs in
// that user's session, not system-wide), so there is one question per such
// user and it names the user and the exact path. Users without a guarded
// ~/.ssh are skipped silently; root is always skipped (no interactive
// session); an already-installed matching unit is kept without a prompt.
func offerSSHAgentUnits(cfg *daemonconfig.Config) error {
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return err
	}
	covered := configPaths(cfg)
	for i := range users {
		u := users[i]
		if u.UID == 0 {
			continue
		}
		sshDir := filepath.Join(u.Home, ".ssh")
		if !pathCovered(sshDir, covered) {
			continue
		}
		unitPath := filepath.Join(u.Home, ".config", "systemd", "user", "ssh-agent.service")
		if cur, rerr := os.ReadFile(unitPath); rerr == nil && bytes.Equal(cur, sample) {
			log.Infof("ssh-agent unit for %s already installed — keeping it", u.Name)
			continue
		}

		install := true
		if ferr := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Install the ssh-agent unit for user %s?", u.Name)).
				Description(fmt.Sprintf(
					"%s is guarded, so ssh-agent must run as a known service to keep reading the keys.\n"+
						"This installs and enables a per-user systemd unit (that user's session only) at:\n    %s",
					sshDir, unitPath)).
				Affirmative("Install it").
				Negative("Skip").
				Value(&install),
		)).Run(); ferr != nil {
			return ferr
		}
		if !install {
			log.Infof("skipping the ssh-agent unit for %s — start ssh-agent another way, or guarded ~/.ssh access from it will be denied", u.Name)
			continue
		}
		if err := installSSHAgent(u); err != nil {
			return err
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

// installConfig writes the final daemon config to its system path and
// ensures the PATH symlink. The binary is deployed separately (the
// one-line installer or `install --binary-only`); the wizard only checks
// it is present (see ensureInstalledBinary). An existing config with
// different content is diffed and the user is asked whether to overwrite
// it; identical configs are left alone. It reports whether the installed
// config differs from what the daemon was running with.
func installConfig(cfgText string) (configChanged bool, err error) {
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
