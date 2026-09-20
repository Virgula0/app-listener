// Package protected holds the daemon-guard and fscrypt-key checks shared by the installer,
// uninstaller and edit-protected: verifying which catalog directories are encrypted with the master
// key, and refusing to run while the daemon is alive.
package protected

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/install"
)

// daemonServiceName is the systemd unit that guards the encrypted
// directories; the guard against running it while churning the vaults.
const daemonServiceName = "app-listener-daemon"

// DaemonRunning reports whether the app-listener-daemon systemd unit is
// active (empty when systemctl is unavailable).
func DaemonRunning() bool {
	return strings.TrimSpace(systemctlActive()) == "active"
}

// systemctlActive returns the trimmed systemctl is-active output for the
// daemon unit (empty on failure).
func systemctlActive() string {
	cmd := exec.CommandContext(context.Background(), "systemctl", "is-active", daemonServiceName)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// FindDaemonProcesses scans procRoot for an app-listener daemon process ("app-listener daemon", any
// --headless or foreground run) and returns PIDs. Exited or unreadable processes are skipped.
func FindDaemonProcesses(procRoot string) ([]int, error) {
	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, convErr := strconv.Atoi(e.Name())
		if convErr != nil {
			continue
		}
		cmdline, readErr := os.ReadFile(filepath.Join(procRoot, e.Name(), "cmdline"))
		if readErr != nil {
			continue
		}
		if isDaemonCmdline(string(cmdline)) {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

// isDaemonCmdline: a NUL-separated /proc cmdline belongs to an app-listener daemon (exe is
// app-listener, "daemon" among the args).
func isDaemonCmdline(cmdline string) bool {
	args := strings.Split(cmdline, "\x00")
	if len(args) == 0 || args[0] == "" {
		return false
	}
	if filepath.Base(args[0]) != "app-listener" {
		return false
	}
	for _, a := range args[1:] {
		if a == "daemon" {
			return true
		}
	}
	return false
}

// RequireDaemonStopped fatally refuses while a daemon runs: an active systemd unit or a daemon
// process outside systemd; the operator must stop them manually.
func RequireDaemonStopped() error {
	if DaemonRunning() {
		return fmt.Errorf("fatal: the daemon is running — stop it first: systemctl stop %s", daemonServiceName)
	}
	pids, err := FindDaemonProcesses("/proc")
	if err != nil {
		return fmt.Errorf("scanning for a running daemon process: %w", err)
	}
	if len(pids) > 0 {
		return fmt.Errorf("fatal: the daemon is running outside systemd (pid(s) %v) — stop it first, e.g. kill %d", pids, pids[0])
	}
	log.Info("daemon is not running")
	return nil
}

// ScanEncryptedCatalogDirs re-scans the catalog for every local user (plus system entries) and
// returns the directories currently fscrypt-encrypted, mirroring what the installer would have
// picked.
func ScanEncryptedCatalogDirs(vault interface {
	IsEncrypted(path string) (bool, error)
}) ([]string, error) {
	users, err := install.ListUsers()
	if err != nil {
		return nil, err
	}
	candidates := install.DiscoverForUsers(users)
	if len(candidates) == 0 {
		log.Warn("no catalog directories found for any user")
		return nil, nil
	}
	paths := make([]string, 0, len(candidates))
	for i := range candidates {
		paths = append(paths, candidates[i].Path)
	}
	encrypted, err := ExistingEncryptedDirs(vault, paths)
	if err != nil {
		return nil, err
	}
	if len(encrypted) == 0 {
		log.Info("no catalog directory is encrypted")
	}
	return encrypted, nil
}

// ExistingEncryptedDirs returns the paths in the set that are currently encrypted: a kernel fscrypt
// policy (directory) or a file-vault ciphertext record (regular file; vault.IsEncrypted dispatches
// on type, see isRegularFileTarget in internal/fscrypt). Plain entries of either kind are skipped.
func ExistingEncryptedDirs(vault interface {
	IsEncrypted(path string) (bool, error)
}, paths []string) ([]string, error) {
	var encrypted []string
	for _, path := range paths {
		ok, err := vault.IsEncrypted(path)
		if err != nil {
			return nil, fmt.Errorf("checking encryption of %s: %w", path, err)
		}
		if ok {
			log.Infof("%s is encrypted", path)
			encrypted = append(encrypted, path)
		}
	}
	return encrypted, nil
}
