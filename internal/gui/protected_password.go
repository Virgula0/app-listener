//go:build gui

package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// fillPasswordBox offers what `--set-password` / `--clear-password` would accept now.
func (p *protectedPage) fillPasswordBox() {
	p.pwBox.Objects = nil
	if p.pw == nil {
		p.pwBox.Add(mutedText("Password state unknown: see Edit mode above."))
		p.pwBox.Refresh()
		return
	}
	set := newButton("Set password", icon("fingerprint"), func() { p.passwordDialog(false) })
	change := newButton("Change password", icon("fingerprint"), func() { p.passwordDialog(true) })
	remove := newButton("Remove password", icon("trash"), p.clearDialog)
	remove.Importance = widget.DangerImportance
	switch {
	case !p.pw.Configured && p.pw.DaemonRunning:
		set.Disable()
		p.pwBox.Add(hfill(mutedText("The first password is created by the installer while the daemon runs. "+
			"Stop the daemon to set it here."), set))
	case !p.pw.Configured:
		p.pwBox.Add(hfill(mutedText("No password: live mode is off."), set))
	case p.pw.Origin == "install":
		change.Disable()
		p.pwBox.Add(hfill(mutedText("Chosen during installation: re-run the installer to change it."), change, remove))
	default:
		p.pwBox.Add(hfill(mutedText("Set with edit-protected: you can change or remove it."), change, remove))
	}
	p.pwBox.Refresh()
}

// passwordDialog is --set-password: the current password (when rotating) and the new one twice.
func (p *protectedPage) passwordDialog(rotate bool) {
	cur, next, again := widget.NewPasswordEntry(), widget.NewPasswordEntry(), widget.NewPasswordEntry()
	items := []*widget.FormItem{}
	if rotate {
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
			"Password saved. A running daemon was reloaded to pick it up.")
	}, p.s.win)
	p.s.showSized(d, fyne.NewSize(560, 320))
}

// clearDialog is --clear-password: proving the current password removes it, disabling live mode.
func (p *protectedPage) clearDialog() {
	cur := widget.NewPasswordEntry()
	d := dialog.NewForm("Remove the edit-protected password", "Remove", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Current", cur),
		widget.NewFormItem("", mutedText("Live edits, temporary access and configuration edits stop working.")),
	}, func(ok bool) {
		if ok {
			p.passwordCall(&editbackend.Request{Op: editbackend.OpClearPassword, Password: cur.Text},
				"Password removed: live edit-protected is now disabled.")
		}
	}, p.s.win)
	p.s.showSized(d, fyne.NewSize(560, 240))
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
			dialog.ShowInformation("Edit-protected password", done, p.s.win)
		})
	})
}
