//go:build gui

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/editbackend"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// startConfig is `edit-protected --edit-config`: fetch daemon.conf, edit it here, review the diff;
// the daemon writes it only while it still holds what was fetched, and reloads.
func (p *protectedPage) startConfig() {
	if p.hello != nil && !p.hello.Live {
		dialog.ShowInformation("Daemon configuration", "Editing the configuration needs the daemon running with "+
			"an edit-protected password configured (live mode). With the daemon stopped, edit "+
			"/etc/app-listener/daemon.conf directly.", p.s.win)
		return
	}
	p.begin(func(link backendLink) {
		p.authenticate(link, "Daemon configuration", func([]string) {
			var m editbackend.Message
			var err error
			p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: editbackend.OpConfigGet}) }, func() {
				if err != nil {
					p.s.showError("Daemon configuration", err)
					p.finish()
					return
				}
				v := newConfigView(p, link, m.Data)
				p.leave = v.close
				p.show(v.root)
			})
		})
	})
}

// configView edits daemon.conf; base is what the daemon held when it was fetched.
type configView struct {
	p      *protectedPage
	link   backendLink
	base   []byte
	editor *codeEntry
	status *widget.Label
	root   fyne.CanvasObject
}

func newConfigView(p *protectedPage, link backendLink, base []byte) *configView {
	v := &configView{p: p, link: link, base: base, status: mutedText("")}
	v.editor = newCodeEntry()
	v.editor.SetText(string(base))
	v.editor.onSave = v.review
	apply := newButton("Review and apply", icon("check"), v.review)
	apply.Importance = widget.HighImportance
	revert := newButton("Revert", icon("arrow-clockwise"), func() {
		v.editor.SetText(string(v.base))
		v.status.SetText("reverted to the daemon's configuration")
	})
	closeBtn := newButton("Close", icon("x"), v.close)
	head := newPanel(container.NewVBox(
		container.NewBorder(nil, nil, nil, closeBtn, boldLabel("daemon.conf (live)")),
		mutedText("Ctrl+S or Review and apply shows the diff; the daemon then writes and reloads it. If the reload "+
			"fails, the daemon keeps the current configuration and restores the file."),
	))
	v.root = container.NewBorder(container.NewVBox(head, container.NewHBox(apply, revert)), v.status, nil, nil,
		v.editor)
	return v
}

func (v *configView) review() {
	edited := v.editor.Text
	if edited == string(v.base) {
		v.status.SetText("configuration unchanged — nothing to apply")
		return
	}
	if _, err := daemonconfig.Parse([]byte(edited)); err != nil {
		v.p.s.showError("The configuration does not parse", err)
		return
	}
	body := container.NewBorder(promptBody("Apply these changes and reload the daemon?",
		"If the reload fails, the daemon keeps the current configuration and restores daemon.conf."), nil, nil, nil,
		container.NewScroll(diffText(inst.UnifiedDiff(string(v.base), edited))))
	d := dialog.NewCustomConfirm("daemon.conf — changes to apply", "Apply", "Back to the editor", body,
		func(ok bool) {
			if ok {
				v.apply(edited)
			}
		}, v.p.s.win)
	v.p.s.showSized(d, fyne.NewSize(1040, 720))
}

func (v *configView) apply(edited string) {
	base := v.base
	var m editbackend.Message
	var err error
	v.status.SetText("applying — the daemon rebuilds every guard …")
	v.p.s.background(func() {
		m, err = v.link.Call(&editbackend.Request{Op: editbackend.OpConfigPut, Base: base, Data: []byte(edited)})
	}, func() {
		switch {
		case err == nil:
			v.base = []byte(edited)
			v.status.SetText("configuration saved and reloaded by the daemon")
			dialog.ShowInformation("Daemon configuration", "Saved and reloaded by the daemon.", v.p.s.win)
		case m.Stale:
			// Diff the next review against what the daemon now holds; the edit stays in the editor.
			v.base = m.Data
			v.status.SetText("daemon.conf changed while you were editing it")
			v.p.s.showError("Not applied", fmt.Errorf("%w — review the diff against the current file and apply "+
				"again", err))
		default:
			v.status.SetText("not applied: " + err.Error())
			v.p.s.showError("The configuration was not applied", err)
		}
	})
}

func (v *configView) close() {
	if v.editor.Text != string(v.base) {
		dialog.ShowConfirm("Unsaved changes", "Discard your changes to daemon.conf?", func(ok bool) {
			if ok {
				v.p.finish()
			}
		}, v.p.s.win)
		return
	}
	v.p.finish()
}
