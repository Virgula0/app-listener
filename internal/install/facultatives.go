package install

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// SSHAgentUnitPath is the per-user ssh-agent unit the installer deploys under home.
func SSHAgentUnitPath(home string) string {
	return filepath.Join(home, ".config", "systemd", "user", "ssh-agent.service")
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
	return out
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
