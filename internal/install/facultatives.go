package install

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// SSHAgentUnitPath is the per-user ssh-agent unit the installer deploys under home.
func SSHAgentUnitPath(home string) string {
	return filepath.Join(home, ".config", "systemd", "user", "ssh-agent.service")
}

// legacySSHAgentUnit is the ssh-agent.service earlier releases deployed (socket in %t, unguarded).
const legacySSHAgentUnit = `[Unit]
Description=SSH key agent
Documentation=man:ssh-agent(1)

[Service]
Type=simple
Environment=SSH_AUTH_SOCK=%t/ssh-agent.socket
# The minus sign (-) tells systemd to ignore errors if the file doesn't exist yet
ExecStartPre=-/usr/bin/rm -f %t/ssh-agent.socket
ExecStart=/usr/bin/ssh-agent -D -a %t/ssh-agent.socket
Restart=on-failure
RestartSec=2s

[Install]
WantedBy=default.target
`

// SSHAgentRuntimeDir is the /run/user/<uid> directory the bundled ssh-agent unit binds in.
const SSHAgentRuntimeDir = "ssh-agent"

// SSHAgentSocketGuarded reports whether cfg guards u's SSHAgentRuntimeDir (a runtime section).
func SSHAgentSocketGuarded(cfg *daemonconfig.Config, u User) bool {
	return slices.Contains(cfg.RuntimePaths, fmt.Sprintf("/run/user/%d/%s", u.UID, SSHAgentRuntimeDir))
}

// IsInstallerSSHAgentUnit reports whether data is an ssh-agent unit the installer deployed: the
// bundled one or an earlier release's.
func IsInstallerSSHAgentUnit(data []byte) bool {
	sample, err := SampleContent("ssh-agent.service")
	return (err == nil && bytes.Equal(data, sample)) || string(data) == legacySSHAgentUnit
}

// LoadConfigIfPresent parses the daemon.conf at path; a missing file (nothing installed) is an
// empty config.
func LoadConfigIfPresent(path string) (*daemonconfig.Config, error) {
	cfg, err := daemonconfig.Load(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &daemonconfig.Config{}, nil
	}
	return cfg, err
}

// CoveredPaths returns every path a config resource occupies: the watch path and its encryption
// root (equal unless the section is a grouped watch).
func CoveredPaths(cfg *daemonconfig.Config) []string {
	out := make([]string, 0, len(cfg.Resources)*2)
	for i := range cfg.Resources {
		out = append(out, cfg.Resources[i].Path, cfg.Resources[i].EncryptionRootOrPath())
	}
	return append(out, cfg.RuntimePaths...) // a logged-out user's are no resource
}

// PathCovered: p is already protected by a config path: equal, inside, or a parent of it (a new
// parent watch would overlap an existing guarded tree).
func PathCovered(p string, covered []string) bool {
	for _, cv := range covered {
		if cv == "" {
			continue
		}
		if p == cv || isInsidePath(p, cv) || isInsidePath(cv, p) {
			return true
		}
	}
	return false
}
