//go:build gui

package gui

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// addTarget asks for a path and its options (the CLI's -w, -r, -d, -e) and starts monitoring it.
func (p *filesPage) addTarget() {
	win := p.shell.win
	path := widget.NewEntry()
	path.SetPlaceHolder("/home/user/.ssh")
	browseDir := newButton("Folder…", icon("folder"), func() {
		p.shell.pickPath(pickFolder, "Folder to monitor", "", path.SetText)
	})
	browseFile := newButton("File…", icon("file"), func() {
		p.shell.pickPath(pickFile, "File to monitor", "", path.SetText)
	})

	depth := widget.NewEntry()
	depth.SetText("0")
	depthHint := mutedText("0 = every level")
	depthHint.Wrapping = fyne.TextWrapOff
	recursive := widget.NewCheck("Include subfolders", func(on bool) {
		if on {
			depth.Enable()
		} else {
			depth.Disable()
		}
	})
	recursive.SetChecked(true)

	names := make([]string, 0, len(ebpf.EventTypes()))
	for _, e := range ebpf.EventTypes() {
		names = append(names, e.String())
	}
	events := newCheckGrid(names)

	form := widget.NewForm(
		widget.NewFormItem("Path", container.NewBorder(nil, nil, nil, container.NewHBox(browseDir, browseFile), path)),
		widget.NewFormItem("Scope", recursive),
		widget.NewFormItem("Max depth", container.NewBorder(nil, nil, nil, depthHint, depth)),
		widget.NewFormItem("Show events", events.root),
	)
	note := mutedText("The monitor starts right away. The guard starts off: configure it in the Guard tab, " +
		"or whitelist and blacklist processes from the activity you observe.")

	d := dialog.NewCustomConfirm("Watch a folder or file", "Start watching", "Cancel",
		scrollable(container.NewVBox(form, note)), func(ok bool) {
			if !ok {
				return
			}
			spec, err := specFromForm(path.Text, recursive.Checked, depth.Text, events.selected())
			if err != nil {
				showErrorDialog(err, win)
				return
			}
			var t *workbench.FileTarget
			p.shell.busy("Starting the monitor…", func() error {
				var addErr error
				t, addErr = p.wb.Add(spec)
				return addErr
			}, func() { p.reloadTargets(t) })
		}, win)
	p.shell.showSized(d, fyne.NewSize(760, 440))
}

func specFromForm(path string, recursive bool, depthText string, selected []string) (workbench.FileSpec, error) {
	spec := workbench.FileSpec{Path: strings.TrimSpace(path), Recursive: recursive}
	if recursive {
		n, err := strconv.Atoi(strings.TrimSpace(depthText))
		if err != nil || n < 0 {
			return spec, errBadDepth
		}
		spec.Depth = n
	}
	if len(selected) == 0 {
		return spec, errNoEvents
	}
	if len(selected) < len(ebpf.EventTypes()) {
		for _, s := range selected {
			if et, ok := ebpf.ParseEventType(s); ok {
				spec.Events = append(spec.Events, et)
			}
		}
	}
	return spec, nil
}

var (
	errBadDepth = errors.New("max depth must be a whole number, 0 for every level")
	errNoEvents = errors.New("pick at least one event type to show")
)
