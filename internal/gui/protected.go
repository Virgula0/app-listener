//go:build gui

package gui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// backendLink is one running `edit-protected --backend` (editbackend.Client). The window never
// opens a protected path or the daemon's control socket itself: every read, write and grant
// happens in that child, the installed binary the daemon's peer check admits.
type backendLink interface {
	Hello() editbackend.Hello
	Call(req *editbackend.Request) (editbackend.Message, error)
	Events() <-chan editbackend.Message
	Ended() <-chan string
	Done() <-chan struct{}
	Close()
}

func startBackend() (backendLink, error) {
	c, err := editbackend.Start()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// sessionTimeouts are the live idle timeouts offered (`--timeout-session`).
var sessionTimeouts = []struct {
	label string
	d     time.Duration
}{
	{"5 minutes", 5 * time.Minute},
	{"15 minutes", 15 * time.Minute},
	{"30 minutes", 30 * time.Minute},
	{"1 hour", time.Hour},
	{"4 hours", 4 * time.Hour},
}

// protectedPage is the edit-protected workspace: edit a protected directory (live through the
// daemon, or offline with it stopped), grant other binaries temporary access, edit daemon.conf,
// manage the password. One workflow (one backend) at a time.
type protectedPage struct {
	s     *shell
	stack *fyne.Container
	over  fyne.CanvasObject
	// link is the open workflow's backend; nil when idle.
	link backendLink
	// leave ends the open workflow from outside (window close); nil when idle.
	leave func()

	hello   *editbackend.Hello
	pw      *editbackend.PwState
	probing bool
	probErr error

	modeT, pwT, daemonT *statusTile
	timeout             *widget.Select
	actions             []*handButton
	pwBox               *fyne.Container
}

func newProtectedPage(s *shell) *protectedPage { return &protectedPage{s: s} }

func (p *protectedPage) view() fyne.CanvasObject {
	p.modeT = newStatusTile("lock-simple", "Edit mode")
	p.daemonT = newStatusTile("pulse", "Daemon")
	p.pwT = newStatusTile("fingerprint", "Password")
	labels := make([]string, len(sessionTimeouts))
	for i, t := range sessionTimeouts {
		labels[i] = t.label
	}
	p.timeout = widget.NewSelect(labels, nil)
	p.timeout.Selected = "30 minutes"
	p.pwBox = container.NewVBox()

	cards := newAdaptiveGrid(3,
		p.actionCard("lock-simple-open", theme.ColorNamePrimary, "Edit protected files",
			"Open one guarded directory in a two-pane editor. Live: the daemon grants write access to that "+
				"directory for this session only. Offline (daemon stopped): it is unlocked with the master key "+
				"and locked again when you close the editor.", p.startEdit),
		p.actionCard("arrows-left-right", colorInfo, "Temporary access",
			"Let other binaries reach guarded directories (admit) or deny whitelisted ones, until you stop it. "+
				"Live mode only.", p.startForward),
		p.actionCard("sliders-horizontal", colorInfo, "Daemon configuration",
			"Edit the running daemon's daemon.conf: you review the diff, then the daemon writes and reloads it. "+
				"A configuration that fails to reload is not kept. Live mode only.", p.startConfig),
	)
	timeoutCard := newPanel(container.NewVBox(
		hfill(boldLabel("Live session idle timeout"), p.timeout),
		mutedText("The daemon revokes a live session after this long without input in the session's view "+
			"(and, for temporary access, without a granted binary touching the directories)."),
	))
	p.over = container.NewVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), newGapVBox(
		newAdaptiveGrid(3, p.modeT.card, p.daemonT.card, p.pwT.card),
		container.NewVBox(sectionTitle("ACTIONS"), cards),
		container.NewVBox(sectionTitle("SESSION"), timeoutCard),
		container.NewVBox(sectionTitle("PASSWORD"), newPanel(p.pwBox)),
		container.NewVBox(sectionTitle("HOW IT WORKS"), newPanel(mutedText(
			"This window runs the installed headless binary ("+"edit-protected --backend"+") and drives it. "+
				"That process, the one the daemon admits, authenticates, holds the grant and reads and writes "+
				"the files; this window only shows them. The password goes to it over a private pipe and is "+
				"checked by the daemon before any protected path is shown.\n"+
				"On X11 any program in your session can read this window and type into it: prefer a Wayland "+
				"session when editing secrets."))),
	)))
	p.stack = container.NewStack(p.over)
	p.probe()
	return p.stack
}

func (p *protectedPage) actionCard(iconName string, c fyne.ThemeColorName, title, desc string, run func()) fyne.CanvasObject {
	open := newButton("Start", icon("caret-right"), run)
	open.IconPlacement = widget.ButtonIconTrailingText
	open.Importance = widget.LowImportance
	p.actions = append(p.actions, open)
	return newPanel(container.NewVBox(newIconTile(iconName, c, 48),
		widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mutedText(desc), container.NewHBox(open)))
}

// show swaps the page body; session views get the overview's right inset (clear of the scroll bar).
func (p *protectedPage) show(o fyne.CanvasObject) {
	if o != p.over {
		o = container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), o)
	}
	p.stack.Objects = []fyne.CanvasObject{o}
	p.stack.Refresh()
}

func (p *protectedPage) timeoutSec() int64 {
	for _, t := range sessionTimeouts {
		if t.label == p.timeout.Selected {
			return int64(t.d / time.Second)
		}
	}
	return int64((30 * time.Minute) / time.Second)
}

// probe asks a short-lived backend which mode applies and the password's state.
func (p *protectedPage) probe() {
	if p.probing || p.link != nil {
		return
	}
	p.probing = true
	var hello *editbackend.Hello
	var pw *editbackend.PwState
	var err error
	p.s.background(func() {
		link, serr := p.s.startBackend()
		if serr != nil {
			err = serr
			return
		}
		defer link.Close()
		h := link.Hello()
		hello = &h
		m, cerr := link.Call(&editbackend.Request{Op: editbackend.OpPasswordState})
		if cerr != nil {
			err = cerr
			return
		}
		pw = m.Password
	}, func() {
		p.probing = false
		p.hello, p.pw, p.probErr = hello, pw, err
		p.applyStatus()
	})
}

func (p *protectedPage) applyStatus() {
	switch {
	case p.probErr != nil:
		p.modeT.set(false, "unavailable", p.probErr.Error(), "")
	case p.hello.Live:
		p.modeT.set(true, "live", "The daemon is running with a password: edits go through its control "+
			"socket while it keeps guarding.", theme.ColorNameSuccess)
	case p.hello.HasPassword:
		p.modeT.set(true, "offline", "A password is configured, but the daemon control socket is not present: "+
			"edits use the offline flow (the daemon must be stopped).", colorInfo)
	default:
		p.modeT.set(true, "offline", "No password configured: edits use the offline flow (the daemon must be "+
			"stopped). Set a password to edit live.", colorInfo)
	}
	if p.pw == nil {
		p.daemonT.set(false, "unknown", "", "")
		p.pwT.set(false, "unknown", "", "")
	} else {
		if p.pw.DaemonRunning {
			p.daemonT.set(true, "running", "Live edits, temporary access and configuration edits are available "+
				"with the password.", theme.ColorNameSuccess)
		} else {
			p.daemonT.set(true, "stopped", "Offline edits only: nothing guards the directories right now.", colorInfo)
		}
		if p.pw.Configured {
			p.pwT.set(true, "set", "Chosen "+originLabel(p.pw.Origin)+".", theme.ColorNameSuccess)
		} else {
			p.pwT.set(true, "none", "Live mode is off until a password is set.", colorMuted)
		}
	}
	p.fillPasswordBox()
}

func originLabel(origin string) string {
	if origin == "install" {
		return "during installation"
	}
	return "with edit-protected"
}

// begin starts a workflow's backend unless another run or session is open; then runs on the UI
// goroutine with it.
func (p *protectedPage) begin(then func(link backendLink)) {
	if p.link != nil {
		dialog.ShowInformation("A session is open", "Finish the open edit-protected session first.", p.s.win)
		return
	}
	if p.s.running() {
		dialog.ShowInformation("A run is in progress", "Wait for the installer or updater run to finish.", p.s.win)
		return
	}
	if active := p.s.activeGuards(); len(active) > 0 {
		dialog.ShowInformation("Guards are running", "Turn off "+strings.Join(active, ", ")+
			" first (Files and Network workspaces): a guard running in this window doesn't admit the "+
			"edit-protected backend, so it would deny its reads and writes.", p.s.win)
		return
	}
	for _, b := range p.actions {
		b.Disable()
	}
	var link backendLink
	var err error
	p.s.background(func() { link, err = p.s.startBackend() }, func() {
		for _, b := range p.actions {
			b.Enable()
		}
		if err != nil {
			p.s.showError("edit-protected", err)
			return
		}
		p.link = link
		then(link)
	})
}

// finish closes the workflow's backend (it narrows the grant or re-locks on EOF) and returns to
// the overview.
func (p *protectedPage) finish() {
	link := p.link
	p.link, p.leave = nil, nil
	p.show(p.over)
	if link == nil {
		return
	}
	p.s.background(link.Close, p.probe)
}

// authenticate asks for the password (before any protected path is shown) and AUTHs; a refusal asks
// again, as many times as the operator retries (the daemon rate-limits and locks out).
func (p *protectedPage) authenticate(link backendLink, title string, ok func(resources []string)) {
	entry := widget.NewPasswordEntry()
	body := promptBody("Edit-protected password", "The daemon checks it before it shows any protected path. "+
		"Five wrong passwords in a row lock live mode for a minute.", entry)
	d := dialog.NewCustomConfirm(title, "Unlock", "Cancel", body, func(confirm bool) {
		if !confirm {
			p.finish()
			return
		}
		pw := entry.Text
		var m editbackend.Message
		var err error
		p.s.background(func() {
			m, err = link.Call(&editbackend.Request{Op: editbackend.OpAuth, Password: pw, TimeoutSec: p.timeoutSec()})
		}, func() {
			if errors.Is(err, editbackend.ErrClosed) {
				p.s.showError(title, err)
				p.finish()
				return
			}
			if err != nil {
				dialog.ShowCustomConfirm(title, "Try again", "Cancel", mutedText(err.Error()), func(again bool) {
					if again {
						p.authenticate(link, title, ok)
					} else {
						p.finish()
					}
				}, p.s.win)
				return
			}
			ok(m.Resources)
		})
	}, p.s.win)
	entry.OnSubmitted = func(string) { d.Confirm() }
	p.s.showSized(d, fyne.NewSize(560, 240))
	p.s.win.Canvas().Focus(entry)
}

// pickOne asks which resource to open; a single one is taken as is.
func (p *protectedPage) pickOne(title, question string, resources []string, ok func(string)) {
	if len(resources) == 1 {
		ok(resources[0])
		return
	}
	radio := widget.NewRadioGroup(resources, nil)
	radio.Required = true
	radio.SetSelected(resources[0])
	d := dialog.NewCustomConfirm(title, "Open", "Cancel", scrollable(promptBody(question, "", radio)),
		func(confirm bool) {
			if !confirm || radio.Selected == "" {
				p.finish()
				return
			}
			ok(radio.Selected)
		}, p.s.win)
	p.s.showSized(d, fyne.NewSize(640, 480))
}

// pickMany asks which resources a forward grant covers.
func (p *protectedPage) pickMany(title string, resources []string, ok func([]string)) {
	checks := make([]*widget.Check, len(resources))
	list := container.NewVBox()
	for i, r := range resources {
		checks[i] = widget.NewCheck(r, nil)
		list.Add(checks[i])
	}
	if len(checks) == 1 {
		checks[0].Checked = true
	}
	body := container.NewBorder(promptBody("Which guarded directories should the binaries reach?",
		"The access is revoked when you stop it."), nil, nil, nil, container.NewVScroll(list))
	d := dialog.NewCustomConfirm(title, "Grant", "Cancel", body, func(confirm bool) {
		var picked []string
		for i, c := range checks {
			if c.Checked {
				picked = append(picked, resources[i])
			}
		}
		if !confirm || len(picked) == 0 {
			p.finish()
			return
		}
		ok(picked)
	}, p.s.win)
	p.s.showSized(d, fyne.NewSize(680, 520))
}

// endedBanner is the red line a live view shows once the daemon ended its session.
func endedBanner(why string) fyne.CanvasObject {
	l := widget.NewLabelWithStyle(why, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	l.Wrapping = fyne.TextWrapWord
	l.Importance = widget.DangerImportance
	return newPanel(l)
}

// watchEnded runs onEnded on the UI goroutine when the daemon ends link's session or the backend
// exits.
func (p *protectedPage) watchEnded(link backendLink, onEnded func(why string)) {
	if p.s.syncBackground {
		return
	}
	go func() {
		var why string
		select {
		case why = <-link.Ended():
		case <-link.Done():
			why = "the edit-protected backend exited — the session is over"
		}
		fyne.Do(func() {
			if p.link == link {
				onEnded(why)
			}
		})
	}()
}

// showError is a dialog naming the action that failed.
func (s *shell) showError(title string, err error) {
	d := dialog.NewError(fmt.Errorf("%s: %w", title, err), s.win)
	d.Show()
}
