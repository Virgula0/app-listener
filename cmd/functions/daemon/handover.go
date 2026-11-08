package daemon

import (
	"errors"
	"io/fs"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
)

// takeHandover returns the generation a clean stop left guarding the runtime sections
// (guard.Guard.HandOver), "" if none, once this binary may read them as its guards' populate
// will: after an update the stopped run's guards know only the old inode. A handover that still
// denies it would fail every start with the sections on the old policy and every vault locked, so
// on any failure it is retired first, as every stop's pins were before handovers: the sections
// are then unguarded until this start's guards hold them.
func takeHandover(cfg *daemonconfig.Config, base string) string {
	gen, err := guard.AdmitSelfInHandover(base, daemonSelfBaselineEvents)
	for i := 0; err == nil && gen != "" && i < len(cfg.Resources); i++ {
		if cfg.Resources[i].Runtime {
			err = walkReadable(cfg.Resources[i].Path)
		}
	}
	if err == nil {
		return gen
	}
	log.Errorf("daemon: CRITICAL: cannot start under the runtime-section guards the stopped daemon left (%v) — "+
		"retiring them: those sections are unguarded until this start's guards attach", err)
	if rErr := guard.RetireHandover(base); rErr != nil {
		log.Errorf("daemon: retiring the stopped daemon's runtime-section pins: %v", rErr)
	}
	return ""
}

func walkReadable(root string) error {
	return filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err == nil {
			_, err = d.Info()
		}
		if errors.Is(err, fs.ErrNotExist) {
			return nil // an agent removed its socket meanwhile, or the user logged out
		}
		return err
	})
}
