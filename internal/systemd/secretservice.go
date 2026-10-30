package systemd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// The Secret Service gate's system files, written by the daemon to match its config ([secret_service]
// with a provider) and left in place when it stops: without the relay the keyring can't own
// org.freedesktop.secrets, so a stopped daemon turns the secret store off rather than open.
var (
	// SecretBusPolicyFile is the session bus policy: only root (the relay) owns the secrets name, only
	// through it are the Secret interfaces reachable (never by a stray provider's unique name), and
	// no one but root monitors or eavesdrops. dbus-daemon needs the per-type eavesdrop denies: a bare
	// eavesdrop="true" deny also refuses ordinary receives.
	SecretBusPolicyFile = "/etc/dbus-1/session.d/app-listener-secrets.conf" //nolint:gosec // a path, not a credential
	// SecretKeyringDropIn points gnome-keyring's user unit at the relay socket (%U: the user's uid).
	SecretKeyringDropIn = "/etc/systemd/user/gnome-keyring-daemon.service.d/app-listener.conf" //nolint:gosec // a path, not a credential
)

var secretBusPolicy = []byte(`<!DOCTYPE busconfig PUBLIC "-//freedesktop//DTD D-BUS Bus Configuration 1.0//EN"
 "http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd">
<!-- Installed by app-listener's daemon for its Secret Service gate ([secret_service] in daemon.conf). -->
<busconfig>
  <policy context="default">
    <allow user="root"/>
    <deny own="org.freedesktop.secrets"/>
    <deny receive_type="method_call" eavesdrop="true"/>
    <deny receive_type="method_return" eavesdrop="true"/>
    <deny receive_type="error" eavesdrop="true"/>
    <deny receive_type="signal" eavesdrop="true"/>
    <deny send_destination="org.freedesktop.DBus" send_interface="org.freedesktop.DBus.Monitoring"/>
    <deny send_interface="org.freedesktop.Secret.Service"/>
    <deny send_interface="org.freedesktop.Secret.Collection"/>
    <deny send_interface="org.freedesktop.Secret.Item"/>
    <deny send_interface="org.freedesktop.Secret.Session"/>
    <deny send_interface="org.freedesktop.Secret.Prompt"/>
    <allow send_destination="org.freedesktop.secrets"/>
  </policy>
  <policy user="root">
    <allow own="org.freedesktop.secrets"/>
    <allow send_destination="org.freedesktop.DBus" send_interface="org.freedesktop.DBus.Monitoring"/>
  </policy>
</busconfig>
`)

var secretKeyringDropIn = []byte(`# Installed by app-listener's daemon: gnome-keyring reaches the session bus only through the
# daemon's Secret Service gate ([secret_service] in daemon.conf).
[Service]
Environment=DBUS_SESSION_BUS_ADDRESS=unix:path=` + daemonconfig.SecretRelayDir + `/%U.bus
`)

// SyncSecretServiceFiles installs (on) or removes the gate's files; changed reports whether any
// file was written or removed.
func SyncSecretServiceFiles(on bool) (changed bool, err error) {
	for _, f := range []struct {
		path    string
		body    []byte
		ownsDir bool
	}{{SecretBusPolicyFile, secretBusPolicy, false}, {SecretKeyringDropIn, secretKeyringDropIn, true}} {
		var did bool
		if on {
			did, err = installFile(f.path, f.body)
		} else {
			did, err = removeFile(f.path, f.ownsDir)
		}
		if err != nil {
			return changed, err
		}
		changed = changed || did
	}
	return changed, nil
}

// installFile writes body to path atomically (0644) unless it already holds it.
func installFile(path string, body []byte) (bool, error) {
	if cur, err := os.ReadFile(path); err == nil && bytes.Equal(cur, body) {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".app-listener-*")
	if err != nil {
		return false, fmt.Errorf("staging %s: %w", path, err)
	}
	defer os.Remove(tmp.Name())
	_, werr := tmp.Write(body)
	if werr == nil {
		werr = tmp.Chmod(0o644)
	}
	if cerr := tmp.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		return false, fmt.Errorf("staging %s: %w", path, werr)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return false, fmt.Errorf("installing %s: %w", path, err)
	}
	return true, nil
}

// removeFile deletes path and, with ownsDir, its directory once empty.
func removeFile(path string, ownsDir bool) (bool, error) {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("removing %s: %w", path, err)
	}
	if ownsDir {
		_ = os.Remove(filepath.Dir(path)) // only succeeds when empty
	}
	return true, nil
}
