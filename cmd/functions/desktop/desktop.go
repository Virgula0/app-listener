// Package desktop is the gui build's entry (`app-listener-gui`): it opens the desktop window, which
// configures every operation itself.
package desktop

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Virgula0/app-listener/cmd/functions/install"
	"github.com/Virgula0/app-listener/cmd/functions/uninstall"
	"github.com/Virgula0/app-listener/cmd/functions/update"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// Run opens the GUI. It needs the gui build and root (every workspace loads eBPF programs).
func Run() error {
	if !gui.Available {
		return errors.New("this binary was built without GUI support: use the app-listener-gui release " +
			"asset or rebuild with `make build GUI=1`")
	}
	if os.Geteuid() != 0 {
		return fmt.Errorf("the window needs root to load eBPF programs: run `sudo %s`", os.Args[0])
	}
	return gui.Run(&gui.Options{Version: constants.Version, Release: releaseBuild(), Install: installHooks(),
		Update: gui.UpdateHooks{Run: update.RunWith, Latest: update.Latest},
		Uninstall: gui.UninstallHooks{Run: func(ui wizard.UI, req gui.UninstallRequest) error {
			return uninstall.RunWith(ui, uninstall.Request{DeleteKey: req.DeleteKey, KeepGUI: req.KeepGUI,
				Facultatives: req.Facultatives})
		}}})
}

// releaseBuild reports whether -ldflags stamped constants.Version, as the release pipeline does: the
// source default names the last release, so a local build must not fetch that tag as "itself".
func releaseBuild() bool {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return false
	}
	for _, s := range info.Settings {
		if s.Key == "-ldflags" && strings.Contains(s.Value, "/internal/constants.Version=") {
			return true
		}
	}
	return false
}

var installModes = map[gui.InstallMode]install.Mode{
	gui.InstallWizard:         install.ModeWizard,
	gui.InstallRefreshCatalog: install.ModeRefreshCatalog,
	gui.InstallDiffCatalog:    install.ModeDiffCatalog,
	gui.InstallRestoreBackups: install.ModeRestoreBackups,
	gui.InstallDeleteBackups:  install.ModeDeleteBackups,
	gui.InstallFacultatives:   install.ModeFacultatives,
}

func installHooks() gui.InstallHooks {
	return gui.InstallHooks{
		Run: func(ui wizard.UI, req gui.InstallRequest) error {
			mode, ok := installModes[req.Mode]
			if !ok {
				return fmt.Errorf("unknown installer mode %d", req.Mode)
			}
			return install.RunWith(ui, install.Request{
				Mode: mode, Binary: req.Binary, GUIBinary: req.GUIBinary, Live: req.Live,
				AllowMetadataOutput: req.AllowMetadataOutput,
			})
		},
		Fetch: func(tag, dir string) (gui.StagedRelease, error) {
			return staged(update.FetchVerified(tag, dir))
		},
		FetchLatest: func(pre bool, dir string) (gui.StagedRelease, error) {
			return staged(update.FetchLatest(pre, dir))
		},
		InstalledVersion:    update.InstalledVersion,
		InstalledGUIVersion: update.InstalledGUIVersion,
	}
}

func staged(r update.Release, err error) (gui.StagedRelease, error) {
	return gui.StagedRelease{Tag: r.Tag, Binary: r.Binary, GUI: r.GUI}, err
}
