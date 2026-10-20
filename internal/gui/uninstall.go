//go:build gui

package gui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// uninstallCard is the Install page's `uninstall` entry. It runs on the page's run view, so it
// excludes installer and updater runs like they exclude each other.
type uninstallCard struct {
	keepGUI, deleteKey *handCheck
	btn                *handButton
	note               *widget.Label
	box                *fyne.Container
}

func (p *installPage) uninstallCard() fyne.CanvasObject {
	u := &p.uninst
	u.keepGUI = newCheck("Keep the desktop build (app-listener-gui)", nil)
	u.deleteKey = newCheck("Also delete the fscrypt master key", nil)
	u.btn = newButton("Uninstall", icon("trash"), p.confirmUninstall)
	u.btn.Importance = widget.DangerImportance
	u.note = warnText()

	flag := widget.NewLabel("uninstall")
	flag.TextStyle = fyne.TextStyle{Monospace: true}
	flag.Importance = widget.LowImportance
	flag.SizeName = theme.SizeNameCaptionText
	head := container.NewBorder(nil, nil, newIconTile("shield-slash", theme.ColorNameError, 40), nil,
		container.NewVBox(boldLabel("Uninstall app-listener"), flag))
	u.box = container.NewVBox(head,
		mutedText("Ends the protection. Each encrypted folder is asked about: decrypt it permanently or keep it "+
			"encrypted. Then the daemon unit, the binary, daemon.conf, the binary ledger and password, the "+
			"per-user ssh-agent and Bun setup (asked first) are removed, and the migration backups are offered "+
			"for deletion. Refused while the daemon runs."),
		u.keepGUI, mutedText("Removing it deletes "+systemd.GUIBinaryPath+"; this window keeps running until "+
			"you close it."),
		u.deleteKey, mutedText("--delete-key. Without the key, any folder left encrypted can never be unlocked "+
			"again: keeping a folder encrypted asks once more before the key goes."),
		u.note, container.NewHBox(u.btn))
	return newPanel(u.box)
}

// facultativesCard is `uninstall facultatives`, which works with the daemon running.
func (p *installPage) facultativesCard() fyne.CanvasObject {
	btn := newButton("Remove facultatives", icon("broom"),
		func() { p.startUninstall(UninstallRequest{Facultatives: true}) })
	flag := widget.NewLabel("uninstall facultatives")
	flag.TextStyle = fyne.TextStyle{Monospace: true}
	flag.Importance = widget.LowImportance
	flag.SizeName = theme.SizeNameCaptionText
	head := container.NewBorder(nil, nil, newIconTile("broom", theme.ColorNameWarning, 40), nil,
		container.NewVBox(boldLabel("Remove facultatives"), flag))
	return newPanel(container.NewBorder(container.NewVBox(head,
		mutedText("Lists the optional per-user setup the installer added (the ssh-agent unit and its "+
			"SSH_AUTH_SOCK block, the Bun launcher wrappers) and removes only the ones you pick. The rest "+
			"of the install stays, and the daemon keeps running.")), container.NewHBox(btn), nil, nil))
}

const outsideSystemdNote = "A daemon process runs outside systemd: stop it first."

// daemonStopNote is why an action that needs the daemon stopped is refused ("" = none runs).
func daemonStopNote(unit, proc bool) string {
	switch {
	case unit:
		return "The daemon is running: stop it first on the Daemon page (or sudo systemctl stop " +
			systemd.DaemonServiceName + "). Stopping it locks its folders again."
	case proc:
		return outsideSystemdNote
	}
	return ""
}

// renderUninstall enables the uninstall only with no daemon running, as the uninstaller requires.
func (p *installPage) renderUninstall() {
	u := &p.uninst
	if u.btn == nil {
		return
	}
	note := daemonStopNote(p.daemonUp, p.daemonProc)
	setWarn(u.note, note)
	if note == "" && !p.running {
		u.btn.Enable()
	} else {
		u.btn.Disable()
	}
	u.box.Refresh()
}

func (p *installPage) confirmUninstall() {
	u := &p.uninst
	req := UninstallRequest{DeleteKey: u.deleteKey.Checked, KeepGUI: u.keepGUI.Checked}
	lines := []string{"Every protected folder loses its guard. Encrypted folders are asked about one by one."}
	if req.DeleteKey {
		lines = append(lines, "The fscrypt master key is DELETED at the end.")
	} else {
		lines = append(lines, "The fscrypt master key is kept.")
	}
	if req.KeepGUI {
		lines = append(lines, "The desktop build is kept.")
	} else {
		lines = append(lines, "The desktop build is removed.")
	}
	showConfirm("Uninstall app-listener?", strings.Join(lines, "\n"), func(ok bool) {
		if ok {
			p.startUninstall(req)
		}
	}, p.s.win)
}

func (p *installPage) startUninstall(req UninstallRequest) {
	run := p.s.opts.Uninstall.Run
	if p.begin("Uninstall", "shield-slash", func(ui wizard.UI) error {
		if run == nil {
			return errors.New("the uninstaller is not available in this build")
		}
		return run(ui, req)
	}, nil) && !req.Facultatives {
		p.uninst.btn.Disable()
	}
}
