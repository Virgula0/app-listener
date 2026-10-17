// Package headless tells the daemon-capable app-listener build from the desktop one: the daemon
// must never run the fyne (`gui` tag) build.
package headless

import (
	"debug/buildinfo"
	"fmt"
	"slices"
	"strings"
)

const mainPath = "github.com/Virgula0/app-listener"

// Check refuses path unless it is an app-listener build without the gui tag or a fyne dependency.
// Unreadable build info refuses too.
func Check(path string) error {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading the build info of %s: %w", path, err)
	}
	if info.Path != mainPath {
		return fmt.Errorf("%s is %q, not an app-listener build", path, info.Path)
	}
	for _, s := range info.Settings {
		if s.Key == "-tags" && slices.Contains(strings.Split(s.Value, ","), "gui") {
			return fmt.Errorf("%s is the GUI build (-tags %s): the daemon needs the headless one", path, s.Value)
		}
	}
	for _, d := range info.Deps {
		if strings.HasPrefix(d.Path, "fyne.io/") {
			return fmt.Errorf("%s links %s: it is the GUI build, the daemon needs the headless one", path, d.Path)
		}
	}
	return nil
}
