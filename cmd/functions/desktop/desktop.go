// Package desktop is the `--gui` entry: it opens the desktop window when no subcommand is given
// (cmd.Execute lets a subcommand win). The window configures every operation itself, so every
// other argument is ignored (and named in a warning).
package desktop

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/functions/install"
	"github.com/Virgula0/app-listener/cmd/functions/update"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/wizard"
)

const flagName = "--gui"

// Requested reports whether args (os.Args[1:]) ask for the GUI, and the arguments it ignores. A
// `--gui=<bool>` value is honored; an unparsable one leaves the decision to cobra, which rejects it.
func Requested(args []string) (on bool, ignored []string) {
	for i, a := range args {
		if a == "--" {
			ignored = append(ignored, args[i:]...)
			break
		}
		name, val, hasVal := strings.Cut(a, "=")
		if name == flagName {
			if !hasVal {
				on = true
				continue
			}
			if b, err := strconv.ParseBool(val); err == nil {
				on = b
				continue
			}
		}
		ignored = append(ignored, a)
	}
	return on, ignored
}

// Run opens the GUI. It needs the gui build and root (every workspace loads eBPF programs).
func Run(ignored []string) error {
	if !gui.Available {
		return errors.New("--gui: this app-listener binary was built without GUI support: " +
			"use the app-listener-gui release asset or rebuild with `make build GUI=1`")
	}
	if len(ignored) > 0 {
		quoted := make([]string, len(ignored))
		for i, a := range ignored {
			quoted[i] = strconv.Quote(a)
		}
		log.Warnf("--gui: ignoring %s: the GUI configures every operation itself", strings.Join(quoted, " "))
	}
	if os.Geteuid() != 0 {
		return fmt.Errorf("--gui needs root to load eBPF programs: run `sudo %s --gui`", os.Args[0])
	}
	return gui.Run(gui.Options{Version: constants.Version, Release: releaseBuild(), Install: installHooks(),
		Update: gui.UpdateHooks{Run: update.RunWith, Latest: update.Latest}})
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
}

func installHooks() gui.InstallHooks {
	return gui.InstallHooks{
		Run: func(ui wizard.UI, req gui.InstallRequest) error {
			mode, ok := installModes[req.Mode]
			if !ok {
				return fmt.Errorf("unknown installer mode %d", req.Mode)
			}
			return install.RunWith(ui, install.Request{
				Mode: mode, Binary: req.Binary, Live: req.Live, AllowMetadataOutput: req.AllowMetadataOutput,
			})
		},
		Fetch:            update.FetchVerified,
		FetchLatest:      update.FetchLatest,
		InstalledVersion: update.InstalledVersion,
	}
}
