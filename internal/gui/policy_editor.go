//go:build gui

package gui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

var modeLabels = []string{"Off", "Whitelist", "Blacklist"}

const sideWidth = 360

// policyHooks is what a policyEditor asks of its workspace.
type policyHooks struct {
	help map[workbench.GuardMode]string
	// status describes the enforced policy and whether the draft differs from it.
	status func() (enforced string, dirty bool)
	// changed runs after any draft edit, before status (optional).
	changed func()
	apply   func()
	revert  func()
}

// policyEditor is a draft guard policy: mode, executables, Apply and Revert.
type policyEditor struct {
	shell *shell
	hooks policyHooks
	mode  workbench.GuardMode

	modes    *handRadioGroup
	help     *widget.Label
	bins     *binaryList
	enforced *widget.Label
	apply    *handButton
	revert   *handButton
	syncing  bool
}

func newPolicyEditor(s *shell, hooks policyHooks) *policyEditor {
	e := &policyEditor{shell: s, hooks: hooks}
	e.modes = newRadioGroup(modeLabels, func(label string) {
		if e.syncing {
			return
		}
		e.mode = workbench.GuardMode(slices.Index(modeLabels, label))
		e.changed()
	})
	e.modes.Horizontal = true
	e.modes.Required = true
	e.help = mutedText("")
	e.bins = newBinaryList(s, e.changed)
	e.revert = newButton("Revert", icon("arrow-clockwise"), hooks.revert)
	e.apply = newButton("Apply", icon("check"), hooks.apply)
	e.apply.Importance = widget.HighImportance
	e.enforced = mutedText("")
	return e
}

// layout stacks the mode card, the executables (with side controls, if any, in a scrolling column on
// their right) and the Apply footer.
func (e *policyEditor) layout(side fyne.CanvasObject) fyne.CanvasObject {
	top := newPanel(container.NewVBox(boldLabel("Mode"), e.modes, e.help))
	footer := newPanel(container.NewBorder(nil, nil, nil, container.NewHBox(e.revert, e.apply), e.enforced))
	footer.inset = 10
	var right fyne.CanvasObject
	if side != nil {
		scroll := newFadeVScroll(side)
		scroll.SetMinSize(fyne.NewSize(sideWidth, 0))
		right = scroll
	}
	body := fyne.CanvasObject(container.NewBorder(top, footer, nil, right, e.bins.card("Executables")))
	if !e.shell.status.canGuard() {
		l := widget.NewLabel("Guards are unavailable on this host (BPF LSM inactive or not root): see Home.")
		l.Importance = widget.WarningImportance
		l.Wrapping = fyne.TextWrapWord
		body = container.NewBorder(l, nil, nil, nil, body)
		e.modes.Disable()
		e.bins.setEditable(false)
	}
	return body
}

func (e *policyEditor) setMode(m workbench.GuardMode) {
	e.mode = m
	e.syncing = true
	e.modes.SetSelected(modeLabels[m])
	e.syncing = false
	e.changed()
}

func (e *policyEditor) entries() []ebpf.BinaryEntry { return slices.Clone(e.bins.entries) }

func (e *policyEditor) changed() {
	e.help.SetText(e.hooks.help[e.mode])
	if e.hooks.changed != nil {
		e.hooks.changed()
	}
	text, dirty := e.hooks.status()
	e.enforced.SetText(text)
	if dirty && e.shell.status.canGuard() {
		e.apply.Enable()
		e.revert.Enable()
	} else {
		e.apply.Disable()
		e.revert.Disable()
	}
}

// promoteBase is the list a promotion into mode extends: the draft's when it already uses mode,
// else enforced's when that does, else none.
func (e *policyEditor) promoteBase(mode, enforcedMode workbench.GuardMode, enforced []ebpf.BinaryEntry,
	entry ebpf.BinaryEntry) []ebpf.BinaryEntry {
	var base []ebpf.BinaryEntry
	switch mode {
	case e.mode:
		base = e.entries()
	case enforcedMode:
		base = slices.Clone(enforced)
	}
	base = slices.DeleteFunc(base, func(b ebpf.BinaryEntry) bool { return b.Path == entry.Path })
	return append(base, entry)
}

func sameEntries(a, b []ebpf.BinaryEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		if !slices.Contains(b, x) {
			return false
		}
	}
	return true
}

// confirmWarnings runs ok at once, or after the user accepts the warnings.
func confirmWarnings(s *shell, title string, warnings []string, ok func()) {
	if len(warnings) == 0 {
		ok()
		return
	}
	msg := ""
	for _, w := range warnings {
		msg += "• " + w + "\n"
	}
	showConfirm(title, msg+"\nApply anyway?", func(yes bool) {
		if yes {
			ok()
		}
	}, s.win)
}
