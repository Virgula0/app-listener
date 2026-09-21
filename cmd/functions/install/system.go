package install

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/safeio"
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

// installServices copies the embedded unit files from daemon-samples into place (skipping
// existing), installs the per-user ssh-agent units for sshUsers, and drops the package-manager
// catalog-refresh hooks.
func installServices(sshUsers []inst.User) error {
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
			for i := range sshUsers {
				if err := installSSHAgent(sshUsers[i]); err != nil {
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
	return installReloadHooks()
}

// installReloadHooks drops the catalog-refresh post-transaction hook for every package manager
// present. Managers without a hook (dnf, zypper) are reported: the boot-time catalog-refresh unit
// still covers reboots, and `install --update-catalog-only` works by hand.
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

// upsertFile writes data to path unless identical. Differing content is diffed and the user asked
// before overwriting. A non-negative uid chowns the file afterwards (per-user units).
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

// askSSHAgentUsers returns the users to set the ssh-agent up for: non-root users whose ~/.ssh is
// guarded by cfg, after one question per user (naming user and unit path). A user whose matching
// unit is already installed is included without asking. Asked before encryption, because the
// ~/.ssh/config edit must land in the tree that gets migrated.
func askSSHAgentUsers(cfg *daemonconfig.Config) ([]inst.User, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return nil, err
	}
	covered := configPaths(cfg)
	var accepted []inst.User
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
			accepted = append(accepted, u)
			continue
		}

		install := true
		if ferr := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Set up ssh-agent for user %s?", u.Name)).
				Description(fmt.Sprintf(
					"%s is guarded, so ssh-agent must run as a known service to keep reading the keys.\n"+
						"This will:\n"+
						"  - install and enable a per-user systemd unit (that user's session only) at %s\n"+
						"  - add an SSH_AUTH_SOCK block to the user's shell startup file (.zshrc / .bashrc / fish conf.d)\n"+
						"  - put `AddKeysToAgent yes` at the top of %s so ssh loads a key into the agent on first use\n"+
						"    (a loaded key can be used through the agent socket by any process of that user)",
					sshDir, unitPath, filepath.Join(sshDir, "config"))).
				Affirmative("Set it up").
				Negative("Skip").
				Value(&install),
		)).Run(); ferr != nil {
			return nil, ferr
		}
		if !install {
			log.Infof("skipping the ssh-agent setup for %s — start ssh-agent another way, or guarded ~/.ssh access from it will be denied", u.Name)
			continue
		}
		accepted = append(accepted, u)
	}
	return accepted, nil
}

// addKeysToAgent edits each accepted user's ~/.ssh/config; runs before that ~/.ssh is encrypted.
func addKeysToAgent(users []inst.User) error {
	for i := range users {
		path, changed, err := inst.EnsureAddKeysToAgent(users[i])
		if err != nil {
			return fmt.Errorf("setting AddKeysToAgent for %s: %w", users[i].Name, err)
		}
		if changed {
			log.Infof("added AddKeysToAgent yes to %s", path)
		}
	}
	return nil
}

// installSSHAgent installs the ssh-agent user unit for one user, hands it to them and enables it
// for their session. An existing unit is compared with the bundled one: identical skipped,
// differing shown as a diff and asked. Root is skipped.
func installSSHAgent(u inst.User) error {
	if u.UID == 0 {
		log.Debugf("skipping ssh-agent unit for root (no interactive user session)")
		return nil
	}
	unitPath := filepath.Join(u.Home, ".config", "systemd", "user", "ssh-agent.service")
	data, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return err
	}

	// Symlink-safe: root writes into the user's home, so the whole ~/.config/systemd/user chain is
	// resolved O_NOFOLLOW anchored at home. A parent the user has swapped for a symlink (e.g. `user`
	// -> /etc/systemd/system) is refused, instead of root creating and chowning a unit outside the
	// home and systemd later running it as root.
	homeFD, err := unix.Open(u.Home, unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", u.Home, err)
	}
	defer unix.Close(homeFD)
	unitDirFD, err := safeio.DescendCreateOwned(homeFD, []string{".config", "systemd", "user"}, int(u.UID), int(u.GID))
	if err != nil {
		return fmt.Errorf("resolving the ssh-agent unit dir for %s: %w", u.Name, err)
	}
	defer unix.Close(unitDirFD)

	if uErr := upsertUnitAt(unitDirFD, "ssh-agent.service", fmt.Sprintf("ssh-agent unit for %s", u.Name),
		data, 0o600, int(u.UID), int(u.GID)); uErr != nil {
		return uErr
	}

	// Equivalent of `systemctl --user enable ssh-agent` without the user's session bus: a symlink in
	// default.target.wants, created relative to the no-follow-resolved unit dir fd.
	wantsFD, err := safeio.DescendCreateOwned(dupFD(unitDirFD), []string{"default.target.wants"}, int(u.UID), int(u.GID))
	if err != nil {
		return fmt.Errorf("resolving default.target.wants for %s: %w", u.Name, err)
	}
	defer unix.Close(wantsFD)
	if err := unix.Symlinkat(unitPath, wantsFD, "ssh-agent.service"); err != nil && !errors.Is(err, unix.EEXIST) {
		return fmt.Errorf("enabling ssh-agent for %s: %w", u.Name, err)
	}
	log.Infof("ssh-agent enabled for %s (relogin or start manually: systemctl --user start ssh-agent)", u.Name)
	return installSSHAgentEnv(u)
}

// dupFD duplicates fd (DescendCreateOwned consumes/closes the fd it is handed, but the caller still
// needs the unit-dir fd afterwards). Returns -1 on failure, which DescendCreateOwned reports as an
// open error.
func dupFD(fd int) int {
	n, err := unix.Dup(fd)
	if err != nil {
		return -1
	}
	return n
}

// upsertUnitAt writes data to name under dirFD unless identical; a differing existing unit is diffed
// and the user asked. Every access is O_NOFOLLOW: an existing symlink at name is refused, never
// followed. uid/gid >= 0 chowns the file.
func upsertUnitAt(dirFD int, name, label string, data []byte, mode os.FileMode, uid, gid int) error {
	existing, present, err := readAtNoFollow(dirFD, name)
	if err != nil {
		return fmt.Errorf("reading %s: %w", label, err)
	}
	switch {
	case present && bytes.Equal(existing, data):
		log.Infof("%s already up to date: skipping", label)
		return nil
	case present:
		overwrite, cErr := inst.ConfirmOverwrite(name, existing, data)
		if cErr != nil {
			return cErr
		}
		if !overwrite {
			log.Warnf("keeping existing %s", label)
			return nil
		}
		log.Infof("overwrote %s", label)
	default:
		log.Infof("installed %s", label)
	}
	return writeAtNoFollow(dirFD, name, data, mode, uid, gid)
}

// readAtNoFollow reads name under dirFD without following a symlink; a symlink at name is refused.
func readAtNoFollow(dirFD int, name string) (content []byte, present bool, err error) {
	fd, err := unix.Openat(dirFD, name, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		switch {
		case errors.Is(err, unix.ENOENT):
			return nil, false, nil
		case errors.Is(err, unix.ELOOP):
			return nil, false, fmt.Errorf("%w: %s", safeio.ErrSymlink, name)
		default:
			return nil, false, err
		}
	}
	f := os.NewFile(uintptr(fd), name)
	defer f.Close()
	b, err := io.ReadAll(f)
	return b, true, err
}

// writeAtNoFollow creates or truncates name under dirFD (never following a symlink) and chowns it.
func writeAtNoFollow(dirFD int, name string, data []byte, mode os.FileMode, uid, gid int) error {
	fd, err := unix.Openat(dirFD, name,
		unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, uint32(mode))
	if errors.Is(err, unix.EEXIST) {
		var st unix.Stat_t
		if serr := unix.Fstatat(dirFD, name, &st, unix.AT_SYMLINK_NOFOLLOW); serr != nil {
			return serr
		}
		if st.Mode&unix.S_IFMT == unix.S_IFLNK {
			return fmt.Errorf("%w: %s", safeio.ErrSymlink, name)
		}
		fd, err = unix.Openat(dirFD, name, unix.O_WRONLY|unix.O_TRUNC|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	}
	if err != nil {
		return fmt.Errorf("installing %s: %w", name, err)
	}
	f := os.NewFile(uintptr(fd), name)
	if _, werr := f.Write(data); werr != nil {
		_ = f.Close()
		return werr
	}
	if uid >= 0 && gid >= 0 {
		if cerr := f.Chown(uid, gid); cerr != nil {
			_ = f.Close()
			return fmt.Errorf("chown %s: %w", name, cerr)
		}
	}
	return f.Close()
}

// installSSHAgentEnv points u's shells at the unit's socket; without it git/ssh never see the agent.
func installSSHAgentEnv(u inst.User) error {
	files, err := inst.EnsureSSHAgentEnv(u)
	for _, f := range files {
		log.Infof("added SSH_AUTH_SOCK to %s (open a new shell to pick it up)", f)
	}
	if err != nil {
		return fmt.Errorf("setting SSH_AUTH_SOCK for %s: %w", u.Name, err)
	}
	return nil
}

// installConfig writes the final config to its system path and ensures the PATH symlink. The binary
// is deployed separately (one-line installer or `install --binary-only`); the wizard only checks
// it's present (ensureInstalledBinary). A differing existing config is diffed and the user asked;
// identical is left alone. Reports whether the installed config differs from what the daemon was
// running with.
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
