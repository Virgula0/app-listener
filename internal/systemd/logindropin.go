package systemd

import (
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
	written, err := installFile(LoginDropInFile, loginDropIn)
	if written {
		log.Infof("installed %s", LoginDropInFile)
	}
	return err
}

// LoginDropInInstalled reports whether the drop-in exists.
func LoginDropInInstalled() bool {
	_, err := os.Lstat(LoginDropInFile)
	return err == nil
}

// RemoveLoginDropIn deletes the drop-in and its directory once empty; the caller daemon-reloads.
func RemoveLoginDropIn() (bool, error) {
	removed, err := removeFile(LoginDropInFile, true)
	if removed {
		log.Infof("removed %s", LoginDropInFile)
	}
	return removed, err
}
