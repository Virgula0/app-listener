package install

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/headless"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// Mode is an installer entry point the desktop GUI offers: the wizard or one maintenance flag.
type Mode int

const (
	ModeWizard         Mode = iota
	ModeRefreshCatalog      // install catalog refresh (--live: Request.Live)
	ModeDiffCatalog         // install catalog diff
	ModeRestoreBackups      // install backups restore
	ModeDeleteBackups       // install backups delete
	ModeFacultatives        // install facultatives
)

// Request is one GUI run.
type Request struct {
	Mode Mode
	// Binary is the verified headless build deployed in place of the running executable: the GUI
	// build must never become the daemon. Required by the modes that deploy (needsBinary).
	Binary string
	// GUIBinary is the verified desktop build of the same release, installed at
	// systemd.GUIBinaryPath once the wizard completes ("" = none: not chosen, or the release has none).
	GUIBinary           string
	Live                bool
	AllowMetadataOutput bool
}

// noDeploySource is deploySource for a GUI mode that deploys nothing: a path that can't exist, so a
// code path that would deploy fails instead of falling back to the running (GUI) executable.
const noDeploySource = "/nonexistent/app-listener-gui-deploys-nothing"

// needsBinary reports whether mode deploys the binary.
func needsBinary(mode Mode) bool { return mode == ModeWizard }

// deploySource replaces the running executable as what the installer deploys; "" = itself (CLI).
var deploySource string

var runMu sync.Mutex

// RunWith runs req with every prompt on ui. One run at a time: the installer keeps package state.
// TODO(self-binary-guard): runs in the GUI process, see update.RunWith.
func RunWith(ui wizard.UI, req Request) error {
	if !runMu.TryLock() {
		return errors.New("an installer run is already in progress")
	}
	defer runMu.Unlock()
	if needsBinary(req.Mode) {
		if req.Binary == "" {
			return errors.New("no headless binary to deploy")
		}
		// Re-checked here, not trusted from the caller: the daemon must never be the GUI build.
		if err := headless.Check(req.Binary); err != nil {
			return fmt.Errorf("refusing to deploy %s: %w", req.Binary, err)
		}
		if req.GUIBinary != "" {
			if err := headless.CheckGUI(req.GUIBinary); err != nil {
				return fmt.Errorf("refusing to install %s as the desktop build: %w", req.GUIBinary, err)
			}
		}
	}
	restoreUI := wizard.Use(ui)
	prevSrc, prevMeta := deploySource, allowMetadataOutput
	deploySource, allowMetadataOutput = req.Binary, req.AllowMetadataOutput
	if deploySource == "" {
		// Never "" under the GUI: that means "copy the running executable", the fyne build.
		deploySource = noDeploySource
	}
	defer func() {
		deploySource, allowMetadataOutput = prevSrc, prevMeta
		restoreUI()
	}()

	switch req.Mode {
	case ModeWizard:
		if err := runWizard(); err != nil {
			return err
		}
		return installGUI(req.GUIBinary)
	case ModeRefreshCatalog:
		return runRefreshCatalogOnly(false, req.Live)
	case ModeDiffCatalog:
		return runDiffCatalog()
	case ModeRestoreBackups:
		return restoreBackups()
	case ModeDeleteBackups:
		return deletePostBackups()
	case ModeFacultatives:
		return runFacultatives()
	}
	return fmt.Errorf("unknown installer mode %d", req.Mode)
}

// installGUI installs the desktop build next to the daemon after a completed wizard.
func installGUI(path string) error {
	if path == "" {
		log.Info("the desktop build was not installed: only the daemon was")
		return nil
	}
	wizard.Stage("Installing the desktop build")
	if err := systemd.DeployGUIBinary(path); err != nil {
		return fmt.Errorf("the daemon is installed, but its desktop build was not: %w", err)
	}
	log.Infof("desktop build installed: run it with `sudo %s`", filepath.Base(systemd.GUISymlinkPath))
	return nil
}
