package protected

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

// installCertPath records an install wizard that ran to completion: the wizard refuses while it
// exists, and uninstall removes it before anything else. Not a secret: the daemon's read-only guard
// on /etc/app-listener lets every process read it and only the app-listener binary write it.
var installCertPath = "/etc/app-listener/install-complete"

// InstallCertPath is where the completed-install record lives.
func InstallCertPath() string { return installCertPath }

// InstallCert is the completed-install record.
type InstallCert struct {
	Version     string    `json:"version"`
	CompletedAt time.Time `json:"completed_at"`
}

// ReadInstallCert returns the record, nil when none. A record that can't be read or parsed is an
// error: the wizard then refuses rather than guess the install never completed.
func ReadInstallCert() (*InstallCert, error) {
	f, err := safeio.OpenRegularNoFollow(installCertPath, os.O_RDONLY, 0)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", installCertPath, err)
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 4096))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", installCertPath, err)
	}
	var c InstallCert
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", installCertPath, err)
	}
	return &c, nil
}

// WriteInstallCert atomically records a completed install of version.
func WriteInstallCert(version string) error {
	data, err := json.Marshal(InstallCert{Version: version, CompletedAt: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		return err
	}
	dir, name := filepath.Split(installCertPath)
	dirFD, err := unix.Open(dir, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", dir, err)
	}
	defer unix.Close(dirFD)
	return safeio.AtomicWriteAt(dirFD, name, "."+name+".tmp", append(data, '\n'), 0o644, -1, -1)
}

// RemoveInstallCert deletes the record; none is fine.
func RemoveInstallCert() error {
	if err := os.Remove(installCertPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", installCertPath, err)
	}
	return nil
}

// RequireNoInstallCert refuses the install wizard once an install has completed: from then on only
// the maintenance modes (or an uninstall first) apply. An interrupted install has no record and is
// resumed.
func RequireNoInstallCert() error {
	c, err := ReadInstallCert()
	if err != nil {
		return fmt.Errorf("%w — refusing the install wizard; `app-listener uninstall` removes it", err)
	}
	if c == nil {
		return nil
	}
	return fmt.Errorf("app-listener %s is already installed (completed %s): the wizard only runs on a fresh "+
		"or interrupted installation. Use `install --diff-catalog` to protect more folders, "+
		"`install --refresh-catalog-only` to refresh the whitelists, `edit-protected --edit-config` to change "+
		"daemon.conf or `app-listener update` for a new version, or run `app-listener uninstall` first to start over",
		c.Version, c.CompletedAt.Local().Format("2006-01-02 15:04"))
}
