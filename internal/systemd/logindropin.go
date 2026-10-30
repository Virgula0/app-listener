package systemd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// LoginDropInFile is the drop-in on user-runtime-dir@.service: around each user's /run/user/<uid>
// mount it runs `daemon login|logout <uid>` (cmd/functions/daemon/runtime.go), so runtime sections
// are guarded before user@<uid> starts. A package var so tests can point it at a temp directory.
var LoginDropInFile = filepath.Join(SystemdDir, "user-runtime-dir@.service.d", "app-listener.conf")

// loginDropIn fails the login when the daemon refuses (`daemon login` exits 0 with no daemon
// running, as does a removed binary). After= the daemon: a login during its start waits for READY=1.
var loginDropIn = []byte(`# Installed by app-listener: make and guard the runtime sections of daemon.conf
# (runtime: true) in /run/user/<uid> before any of that user's processes start.
[Unit]
After=` + DaemonServiceName + `.service

[Service]
ExecStartPost=/bin/sh -c '[ ! -x ` + InstallBinaryPath + ` ] || exec ` + InstallBinaryPath + ` daemon login %i'
ExecStopPost=-/bin/sh -c '[ ! -x ` + InstallBinaryPath + ` ] || exec ` + InstallBinaryPath + ` daemon logout %i'
`)

// WriteLoginDropIn installs the drop-in unless it is already current; the caller daemon-reloads.
func WriteLoginDropIn() error {
	if cur, err := os.ReadFile(LoginDropInFile); err == nil && bytes.Equal(cur, loginDropIn) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(LoginDropInFile), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(LoginDropInFile), err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(LoginDropInFile), ".app-listener-dropin-*")
	if err != nil {
		return fmt.Errorf("staging %s: %w", LoginDropInFile, err)
	}
	defer os.Remove(tmp.Name())
	_, werr := tmp.Write(loginDropIn)
	if werr == nil {
		werr = tmp.Chmod(0o644)
	}
	if cerr := tmp.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		return fmt.Errorf("staging %s: %w", LoginDropInFile, werr)
	}
	if err := os.Rename(tmp.Name(), LoginDropInFile); err != nil {
		return fmt.Errorf("installing %s: %w", LoginDropInFile, err)
	}
	log.Infof("installed %s", LoginDropInFile)
	return nil
}

// LoginDropInInstalled reports whether the drop-in exists.
func LoginDropInInstalled() bool {
	_, err := os.Lstat(LoginDropInFile)
	return err == nil
}

// RemoveLoginDropIn deletes the drop-in and its directory once empty; the caller daemon-reloads.
func RemoveLoginDropIn() (bool, error) {
	if err := os.Remove(LoginDropInFile); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("removing %s: %w", LoginDropInFile, err)
	}
	_ = os.Remove(filepath.Dir(LoginDropInFile)) // only succeeds when empty
	log.Infof("removed %s", LoginDropInFile)
	return true, nil
}
