package daemon

import (
	"errors"
	"io/fs"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/logging"
)

// inspectorAdmitter is the TrustGuard's AdmitInspector.
type inspectorAdmitter func(path string) (guard.GuardInodeKey, error)

// resolveInspectors admits each configured inspector that root placed (TrustGuard.AdmitInspector).
// A refused entry grants nothing and is logged; a default not installed here is expected.
func resolveInspectors(cfg *daemonconfig.Config, admit inspectorAdmitter) []guard.Inspector {
	var out []guard.Inspector
	for _, p := range cfg.Inspectors {
		k, err := admit(p)
		switch {
		case err == nil:
			out = append(out, guard.Inspector{Path: p, Key: k})
			log.Infof("daemon: inspector %s (inode %d) may read protected processes' /proc metadata",
				logging.SanitizeText(p), k.Ino)
		case cfg.InspectorsDefault && errors.Is(err, fs.ErrNotExist):
			log.Debugf("daemon: default inspector %s not installed", p)
		default:
			log.Errorf("daemon: inspector %s refused (%v) — it may not inspect protected processes; "+
				"only a root-owned program in root-owned directories qualifies", logging.SanitizeText(p), err)
		}
	}
	return out
}

// inspectorPaths are the configured inspectors root placed, for the trusted-binary set: only those
// may become inspectors, and TRUSTED_BINARY would write-protect a user's own file.
func inspectorPaths(cfg *daemonconfig.Config, placed func(string) bool) []string {
	var out []string
	for _, p := range cfg.Inspectors {
		if placed(p) {
			out = append(out, p)
		}
	}
	return out
}
