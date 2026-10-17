//go:build gui

package gui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/wizard"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// runHost is a page that runs installer or updater commands in its run view. Only one host runs at
// a time across the window (shell.running): both stop the daemon and swap its binary, and
// wizard.Use is process-wide.
type runHost struct {
	s     *shell
	stack *fyne.Container
	over  fyne.CanvasObject
	run   *runView
	// running is set while this host's command runs.
	running bool
	// settled runs on the UI goroutine after every run.
	settled func()
}

// mount builds the run view (its back button labeled back) over the page's overview.
func (h *runHost) mount(over fyne.CanvasObject, back string) fyne.CanvasObject {
	h.over = over
	h.run = newRunView(h.s, back, func() { h.show(h.over) })
	h.stack = container.NewStack(over)
	return h.stack
}

func (h *runHost) show(o fyne.CanvasObject) {
	h.stack.Objects = []fyne.CanvasObject{o}
	h.stack.Refresh()
}

// begin switches to the run view and runs work off the UI goroutine with its prompts, stages and
// logs on that view; after runs on success.
func (h *runHost) begin(name, iconName string, work func(ui wizard.UI) error, after func()) bool {
	if h.s.running() {
		dialog.ShowInformation("A run is in progress", "Wait for the installer or updater run to finish.", h.s.win)
		return false
	}
	if why := h.s.guardsBlocking(); why != "" {
		dialog.ShowInformation("Guards are running", why, h.s.win)
		return false
	}
	h.running = true
	h.run.start(name, iconName)
	h.show(h.run.root)
	ui := &installUI{s: h.s, run: h.run}
	restoreLogs := captureLogs(h.run)
	go func() {
		err := work(ui)
		restoreLogs()
		fyne.Do(func() {
			h.running = false
			canceled := errors.Is(err, wizard.ErrCanceled) || errors.Is(err, inst.ErrEditCanceled)
			h.run.finish(err, canceled)
			if err == nil && after != nil {
				after()
			}
			if h.settled != nil {
				h.settled()
			}
			h.s.reprobe()
		})
	}()
	return true
}

// running reports whether an installer or updater run is in progress.
func (s *shell) running() bool {
	return (s.install != nil && s.install.running) || (s.update != nil && s.update.running)
}

// guardsBlocking names why no installer or updater run may start now: a guard this window runs
// would deny it (they move and encrypt guarded folders, restart the daemon, download from GitHub)
// since none of them admits this window.
func (s *shell) guardsBlocking() string {
	var active []string
	for _, t := range s.files.Targets() {
		if t.Policy().Mode != workbench.GuardOff {
			active = append(active, "the guard on "+t.Spec.Path)
		}
	}
	if s.network.Policy().Mode != workbench.GuardOff {
		active = append(active, "the network guard")
	}
	if len(active) == 0 {
		return ""
	}
	return "Turn off " + strings.Join(active, ", ") + " first (Files and Network workspaces).\n\n" +
		"The installer and the updater move and encrypt protected folders, restart the daemon and " +
		"download from GitHub. A guard running in this window doesn't admit them, so it would deny those steps."
}
