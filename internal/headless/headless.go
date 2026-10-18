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
	gui, why, err := inspect(path)
	if err != nil {
		return err
	}
	if gui {
		return fmt.Errorf("%s %s: the daemon needs the headless build", path, why)
	}
	return nil
}

// CheckGUI is Check's inverse: path must be the app-listener desktop build (gui tag).
func CheckGUI(path string) error {
	gui, _, err := inspect(path)
	if err != nil {
		return err
	}
	if !gui {
		return fmt.Errorf("%s is the headless build, not the desktop (gui) one", path)
	}
	return nil
}

// inspect reports whether path is the GUI build of app-listener, and why.
func inspect(path string) (gui bool, why string, err error) {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return false, "", fmt.Errorf("reading the build info of %s: %w", path, err)
	}
	if info.Path != mainPath {
		return false, "", fmt.Errorf("%s is %q, not an app-listener build", path, info.Path)
	}
	for _, s := range info.Settings {
		if s.Key == "-tags" && slices.Contains(strings.Split(s.Value, ","), "gui") {
			return true, "is the GUI build (-tags " + s.Value + ")", nil
		}
	}
	for _, d := range info.Deps {
		if strings.HasPrefix(d.Path, "fyne.io/") {
			return true, "links " + d.Path + " (the GUI build)", nil
		}
	}
	return false, "", nil
}
