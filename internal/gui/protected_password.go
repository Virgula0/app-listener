//go:build gui

package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// fillPasswordBox offers what `edit-protected password` would accept now: a change, which asks the current
// password, or the first password of an install older than the mandatory one. Never a removal.
func (p *protectedPage) fillPasswordBox() {
	p.pwBox.Objects = nil
	if p.pw == nil {
		p.pwBox.Add(mutedText("Password state unknown: see Edit mode above."))
		p.pwBox.Refresh()
		return
	}
	set := newButton("Set password", icon("fingerprint"), func() { p.passwordDialog(false) })
	change := newButton("Change password", icon("fingerprint"), func() { p.passwordDialog(true) })
	switch {
	case !p.pw.Configured && p.pw.DaemonRunning:
		set.Disable()
		note := warnText()
		setWarn(note, daemonStopNote(true, false)+" Or re-run the installer.")
		p.pwBox.Add(hfill(container.NewVBox(mutedText("No password (an install older than the mandatory one): "+
			"a first one is set only with the daemon stopped."), note), set))
	case !p.pw.Configured:
		p.pwBox.Add(hfill(mutedText("No password: live mode is off until one is set."), set))
	default:
		p.pwBox.Add(hfill(mutedText("Changing it asks the current password first. It can't be removed."), change))
	}
	p.pwBox.Refresh()
}

// passwordDialog is edit-protected password: the current password (when changing) and the new one twice.
func (p *protectedPage) passwordDialog(change bool) {
	cur, next, again := widget.NewPasswordEntry(), widget.NewPasswordEntry(), widget.NewPasswordEntry()
	items := []*widget.FormItem{}
	if change {
		items = append(items, widget.NewFormItem("Current", cur))
	}
	items = append(items, widget.NewFormItem("New", next), widget.NewFormItem("Again", again),
		widget.NewFormItem("", mutedText("At least 12 characters, no whitespace, three of: lowercase, uppercase, "+
			"digits, symbols.")))
	d := dialog.NewForm("Edit-protected password", "Save", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}
		if next.Text != again.Text {
			p.s.showError("Edit-protected password", errors.New("the two entries do not match"))
			return
		}
		p.passwordCall(&editbackend.Request{Op: editbackend.OpSetPassword, Password: cur.Text, NewPassword: next.Text},
			"Password saved.")
	}, p.s.win)
	p.s.showSized(d, fyne.NewSize(560, 320))
}

func (p *protectedPage) passwordCall(req *editbackend.Request, done string) {
	p.begin(func(link backendLink) {
		var err error
		p.s.background(func() { _, err = link.Call(req) }, func() {
			p.finish()
			if err != nil {
				p.s.showError("Edit-protected password", err)
				return
			}
			showInformation("Edit-protected password", done, p.s.win)
		})
	})
}
