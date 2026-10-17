//go:build gui

package gui

import (
	"embed"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Icons are Phosphor Icons 2.1.1 (MIT, assets/icons/LICENSE): fill-based paths, the only kind
// fyne's themed resources can recolor (it rewrites `fill`, never `stroke`).
//
//go:embed assets/icons/*.svg assets/logo.svg
var assets embed.FS

func asset(name string) fyne.Resource {
	data, err := assets.ReadFile(name)
	if err != nil {
		panic("gui: missing embedded asset " + name) // embedded at build: unreachable
	}
	return fyne.NewStaticResource(path.Base(name), data)
}

// icon is a foreground-colored Phosphor icon, e.g. icon("shield-check").
func icon(name string) fyne.Resource {
	return theme.NewThemedResource(asset("assets/icons/" + name + ".svg"))
}

// coloredIcon is a Phosphor icon in a theme color (primary, success, error, ...).
func coloredIcon(name string, c fyne.ThemeColorName) fyne.Resource {
	return theme.NewColoredResource(asset("assets/icons/"+name+".svg"), c)
}

var logoResource = asset("assets/logo.svg")
