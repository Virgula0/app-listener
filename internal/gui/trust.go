//go:build gui

package gui

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/logging"
)

// trustPage is `app-listener trust`: the whitelisted binaries the daemon refused. binaries.db is
// sealed to the daemon's own binary, so every read and write goes through the installed binary's
// edit-protected backend; this window only shows the rows and sends back the lines picked, each
// with the content shown.
type trustPage struct {
	s *shell
	// link is the open "add by hand" workflow's backend; nil when idle.
	link    backendLink
	probing bool

	hello *editbackend.Hello
	pw    *editbackend.PwState
	rows  []editbackend.TrustRow
	err   error

	refusedT, daemonT *statusTile
	list              *fyne.Container
	// checks parallels rows: nil for a row that can't be confirmed now.
	checks   []*handCheck
	reload   *handCheck
	trustBtn *handButton
	allBtn   *handButton

	added    []string
	addList  *fyne.Container
	addEntry *widget.Entry
	add      actionGate
}

func newTrustPage(s *shell) *trustPage { return &trustPage{s: s} }

func (t *trustPage) view() fyne.CanvasObject {
	t.refusedT = newStatusTile("shield-check", "Refused binaries")
	t.daemonT = newStatusTile("pulse", "Daemon")
	t.list = container.NewVBox()
	t.reload = newCheck("Reload a running daemon afterwards, so the confirmed binaries are admitted now", nil)
	t.reload.Checked = true
	t.trustBtn = newButton("Trust selected", icon("check"), t.confirmSelected)
	t.trustBtn.Importance = widget.HighImportance
	t.allBtn = newTextButton("Select all", t.selectAll)
	refresh := newButton("Refresh", icon("arrow-clockwise"), t.probe)

	t.addList = container.NewVBox()
	t.addEntry = widget.NewEntry()
	t.addEntry.SetPlaceHolder("/home/alice/.local/bin/tool")
	t.addEntry.OnSubmitted = func(string) { t.addTyped() }
	t.refreshAdded()
	addBtn := newButton("Add", icon("plus"), t.addTyped)
	browse := newButton("Browse…", icon("file"), func() {
		t.s.pickPath(pickFile, "Whitelist line to trust", "", t.addLine)
	})
	review := newButton("Review and trust", icon("caret-right"), t.reviewAdded)
	review.IconPlacement = widget.ButtonIconTrailingText
	t.add.btn, t.add.note = review, warnText()
	t.add.box = container.NewVBox(
		mutedText("A line you added to daemon.conf yourself (trust add): its content as it is now is "+
			"recorded, so the next reload admits it. While the daemon runs this asks for the "+
			"edit-protected password first."),
		t.addList,
		container.NewBorder(nil, nil, nil, container.NewHBox(addBtn, browse), t.addEntry),
		t.add.note, container.NewHBox(review))

	over := newFadeVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), newGapVBox(
		newAdaptiveGrid(2, t.refusedT.card, t.daemonT.card),
		container.NewVBox(sectionTitle("REFUSED BINARIES"), newPanel(container.NewVBox(
			mutedText("The daemon records the content of every whitelisted binary it admits. One whose "+
				"content differs, or a new whitelist line, is refused unless only root could have placed it "+
				"or the app's own updater wrote it while the daemon ran. Trust one only if you installed or "+
				"updated it yourself: only the exact content that was refused is recorded."),
			t.list, t.reload,
			container.NewHBox(t.trustBtn, t.allBtn, layout.NewSpacer(), refresh)))),
		container.NewVBox(sectionTitle("TRUST A LINE BY HAND"), newPanel(t.add.box)),
		container.NewVBox(sectionTitle("HOW IT WORKS"), newPanel(mutedText(
			"The binary ledger is readable only by the daemon's own binary. This window runs the installed "+
				"headless binary (edit-protected --backend), which reads the ledger, hashes each file again "+
				"and records a binary only while it still holds the content shown here."))),
	)))
	t.render()
	t.probe()
	return over
}

// probe asks a short-lived backend for the mode, the daemon's state and the refused binaries; the
// sidebar badge follows the count even while the page isn't built.
func (t *trustPage) probe() {
	if t.probing || t.link != nil {
		return
	}
	t.probing = true
	var hello *editbackend.Hello
	var pw *editbackend.PwState
	var rows []editbackend.TrustRow
	var err error
	t.s.background(func() {
		link, serr := t.s.startBackend()
		if serr != nil {
			err = serr
			return
		}
		defer link.Close()
		h := link.Hello()
		hello = &h
		var m editbackend.Message
		if m, err = link.Call(&editbackend.Request{Op: editbackend.OpPasswordState}); err != nil {
			return
		}
		pw = m.Password
		if m, err = link.Call(&editbackend.Request{Op: editbackend.OpTrustList}); err == nil {
			rows = m.Trust
		}
	}, func() {
		t.probing = false
		t.hello, t.pw, t.rows, t.err = hello, pw, rows, err
		if n := t.s.nav[pageTrust]; n != nil {
			n.setBadge(err == nil && len(rows) > 0)
		}
		t.render()
	})
}

func (t *trustPage) render() {
	if t.list == nil {
		return // not built yet
	}
	t.renderTiles()
	t.list.Objects = nil
	t.checks = make([]*handCheck, len(t.rows))
	confirmable := 0
	for i := range t.rows {
		if t.rows[i].Skip == "" {
			t.checks[i] = newCheck("", nil)
			confirmable++
		}
		t.list.Add(trustRowView(&t.rows[i], t.checks[i], true))
	}
	switch {
	case t.err != nil:
		t.list.Add(warnLabel("Unavailable: " + t.err.Error()))
	case t.hello == nil:
		t.list.Add(mutedText("Reading the binary ledger…"))
	case len(t.rows) == 0:
		t.list.Add(mutedText("Nothing to confirm: the daemon admitted every whitelisted binary."))
	}
	t.list.Refresh()
	if confirmable == 0 {
		t.trustBtn.Disable()
		t.allBtn.Disable()
	} else {
		t.trustBtn.Enable()
		t.allBtn.Enable()
	}
	t.add.set(t.addNote())
}

func (t *trustPage) renderTiles() {
	switch {
	case t.err != nil:
		t.refusedT.set(false, "unavailable", t.err.Error(), "")
	case t.hello == nil:
		t.refusedT.set(true, "checking…", "", colorMuted)
	case len(t.rows) == 0:
		t.refusedT.set(true, "none", "Every whitelisted binary was admitted.", theme.ColorNameSuccess)
	default:
		t.refusedT.set(true, fmt.Sprintf("%d refused", len(t.rows)), "Their whitelist lines stay denied until "+
			"you trust them or the content goes back to the one recorded.", theme.ColorNameWarning)
	}
	switch {
	case t.pw == nil:
		t.daemonT.set(false, "unknown", "", "")
	case t.pw.DaemonRunning:
		t.daemonT.set(true, "running", "Confirmed binaries are admitted when it reloads.", theme.ColorNameSuccess)
	default:
		t.daemonT.set(true, "stopped", "Confirmed binaries are admitted when it starts.", colorInfo)
	}
}

// addNote is why trusting a line by hand isn't available now ("" = it is).
func (t *trustPage) addNote() string {
	switch {
	case t.err != nil:
		return "Unavailable: see Refused binaries above."
	case t.hello == nil || t.pw == nil:
		return checkingDaemon
	case !t.hello.Live && t.pw.DaemonRunning:
		return "The daemon runs without its control socket, so the password can't be checked: " +
			"stop it on the Daemon page, or use sudo app-listener trust add."
	}
	return ""
}

func warnLabel(s string) *widget.Label {
	l := warnText()
	setWarn(l, s)
	return l
}

// trustRowView shows one binary: check (nil = not confirmable, with why), the line and its details.
func trustRowView(r *editbackend.TrustRow, check *handCheck, refused bool) fyne.CanvasObject {
	line := selectable(logging.SanitizeText(r.Line))
	line.TextStyle = fyne.TextStyle{Bold: true}
	var lead fyne.CanvasObject = newIconTile("warning", theme.ColorNameWarning, 26)
	if check != nil {
		lead = check
	}
	box := container.NewVBox(container.NewBorder(nil, nil, container.NewCenter(lead), nil, line),
		mutedText(trustDetails(r, refused)))
	if r.Skip != "" {
		box.Add(warnLabel(logging.SanitizeText(r.Skip)))
	}
	return newPanel(box)
}

func trustDetails(r *editbackend.TrustRow, refused bool) string {
	var b strings.Builder
	if r.Resolved != r.Line {
		fmt.Fprintf(&b, "Resolves to: %s\n", logging.SanitizeText(r.Resolved))
	}
	if refused {
		fmt.Fprintf(&b, "Resource: %s\nRefused at: %s\n", logging.SanitizeText(r.Resources),
			r.SeenAt.Local().Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintf(&b, "sha256: %s", r.SHA256)
	if r.HasStat {
		fmt.Fprintf(&b, "\nOwner uid %d · mode %04o · modified %s", r.UID, r.Mode,
			r.ModTime.Local().Format("2006-01-02 15:04:05"))
	}
	return b.String()
}

func (t *trustPage) selectAll() {
	for _, c := range t.checks {
		if c != nil {
			c.SetChecked(true)
		}
	}
}

// blocked is why no trust workflow may start now ("" = one may): a confirmation reloads the
// daemon, which ends an open edit-protected session; a run may be swapping the binary.
func (t *trustPage) blocked() string {
	switch {
	case t.s.running():
		return "Wait for the running installer, updater or daemon run to finish."
	case t.s.protect != nil && t.s.protect.link != nil:
		return "Finish the open edit-protected session first (Protected page): reloading the daemon would end it."
	case t.link != nil:
		return "A review is already open."
	}
	return ""
}

func (t *trustPage) confirmSelected() {
	var items []editbackend.TrustItem
	var lines []string
	for i, c := range t.checks {
		if c != nil && c.Checked {
			items = append(items, editbackend.TrustItem{Line: t.rows[i].Line, SHA256: t.rows[i].SHA256})
			lines = append(lines, logging.SanitizeText(t.rows[i].Line))
		}
	}
	if len(items) == 0 {
		showInformation("Trust", "Check the binaries to trust first.", t.s.win)
		return
	}
	if why := t.blocked(); why != "" {
		showInformation("Trust", why, t.s.win)
		return
	}
	body := promptBody(fmt.Sprintf("Trust %d binary(ies)?", len(items)),
		"Each one is recorded with the content shown, and only while the file still holds it.",
		selectable(strings.Join(lines, "\n")))
	d := dialog.NewCustomConfirm("Trust", "Trust", "Cancel", scrollable(body), func(ok bool) {
		if !ok {
			return
		}
		noReload := !t.reload.Checked
		var m editbackend.Message
		var err error
		t.trustBtn.Disable()
		t.s.background(func() {
			link, serr := t.s.startBackend()
			if serr != nil {
				err = serr
				return
			}
			defer link.Close()
			m, err = link.Call(&editbackend.Request{Op: editbackend.OpTrustConfirm, Trust: items, NoReload: noReload})
		}, func() {
			if err != nil {
				t.s.showError("Trust", err)
			} else {
				t.showDone(m.TrustDone, noReload)
			}
			t.probe()
		})
	}, t.s.win)
	t.s.showSized(d, fyne.NewSize(640, 420))
}

// showDone reports what a trust_confirm / trust_add recorded and whether the daemon reloaded.
func (t *trustPage) showDone(done *editbackend.TrustDone, noReload bool) {
	if done == nil {
		done = &editbackend.TrustDone{}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d binary(ies) confirmed.", len(done.Confirmed))
	skipped := make([]string, 0, len(done.Skipped))
	for line := range done.Skipped {
		skipped = append(skipped, line)
	}
	slices.Sort(skipped)
	for _, line := range skipped {
		fmt.Fprintf(&b, "\n\nNot recorded: %s\n%s", logging.SanitizeText(line), logging.SanitizeText(done.Skipped[line]))
	}
	switch {
	case len(done.Confirmed) == 0:
	case done.ReloadErr != "":
		fmt.Fprintf(&b, "\n\nReloading the daemon failed: %s", logging.SanitizeText(done.ReloadErr))
	case done.Reloaded:
		b.WriteString("\n\nA running daemon was reloaded to admit them.")
	case noReload:
		b.WriteString("\n\nThey are admitted at the daemon's next reload or start.")
	}
	d := dialog.NewCustom("Trust", "Close", scrollable(selectable(b.String())), t.s.win)
	t.s.showSized(d, fyne.NewSize(640, 360))
}

func (t *trustPage) addTyped() {
	if p := strings.TrimSpace(t.addEntry.Text); p != "" {
		t.addLine(p)
		t.addEntry.SetText("")
	}
}

func (t *trustPage) addLine(p string) {
	if !strings.HasPrefix(p, "/") {
		showInformation("Trust", p+" is not an absolute path.", t.s.win)
		return
	}
	if !slices.Contains(t.added, p) {
		t.added = append(t.added, p)
		t.refreshAdded()
	}
}

func (t *trustPage) refreshAdded() {
	fillRemovableList(t.addList, &t.added, "terminal-window", "No lines yet: type a path or browse for one.",
		t.refreshAdded)
}

// reviewAdded starts the add workflow: AUTH first while the daemon runs (the backend hashes the
// paths with a binary that reads guarded trees), then inspect, review, and record.
func (t *trustPage) reviewAdded() {
	if len(t.added) == 0 {
		showInformation("Trust", "Add the whitelist lines to trust first.", t.s.win)
		return
	}
	if why := t.blocked(); why != "" {
		showInformation("Trust", why, t.s.win)
		return
	}
	lines := slices.Clone(t.added)
	t.add.btn.Disable()
	var link backendLink
	var err error
	t.s.background(func() { link, err = t.s.startBackend() }, func() {
		t.add.set(t.addNote())
		if err != nil {
			t.s.showError("Trust", err)
			return
		}
		t.link = link
		if link.Hello().Live {
			t.s.authenticate(link, "Trust a line", 0, t.endAdd, func([]string) { t.inspect(link, lines) })
			return
		}
		t.inspect(link, lines)
	})
}

func (t *trustPage) inspect(link backendLink, lines []string) {
	var m editbackend.Message
	var err error
	t.s.background(func() {
		m, err = link.Call(&editbackend.Request{Op: editbackend.OpTrustInspect, Binaries: lines})
	}, func() {
		if err != nil {
			t.s.showError("Trust", err)
			t.endAdd()
			return
		}
		t.reviewInspected(link, m.Trust)
	})
}

func (t *trustPage) reviewInspected(link backendLink, rows []editbackend.TrustRow) {
	var items []editbackend.TrustItem
	list := container.NewVBox()
	for i := range rows {
		r := &rows[i]
		list.Add(trustRowView(r, nil, false))
		if r.Skip == "" {
			items = append(items, editbackend.TrustItem{Line: r.Line, SHA256: r.SHA256})
		}
	}
	if len(items) == 0 {
		d := dialog.NewCustom("Trust", "Close", scrollable(container.NewVBox(
			boldLabel("None of these lines can be trusted now."), list)), t.s.win)
		d.SetOnClosed(t.endAdd)
		t.s.showSized(d, fyne.NewSize(680, 480))
		return
	}
	body := container.NewBorder(promptBody(fmt.Sprintf("Trust %d line(s) with this content?", len(items)),
		"Each is recorded only while the file still holds the content shown."), nil, nil, nil, scrollable(list))
	d := dialog.NewCustomConfirm("Trust", "Trust", "Cancel", body, func(ok bool) {
		if !ok {
			t.endAdd()
			return
		}
		noReload := !t.reload.Checked
		var m editbackend.Message
		var err error
		t.s.background(func() {
			m, err = link.Call(&editbackend.Request{Op: editbackend.OpTrustAdd, Trust: items, NoReload: noReload})
		}, func() {
			t.endAdd()
			if err != nil {
				t.s.showError("Trust", err)
				return
			}
			if m.TrustDone == nil {
				m.TrustDone = &editbackend.TrustDone{}
			}
			for _, line := range m.TrustDone.Confirmed {
				t.added = slices.DeleteFunc(t.added, func(p string) bool { return p == line })
			}
			t.refreshAdded()
			t.showDone(m.TrustDone, noReload)
		})
	}, t.s.win)
	t.s.showSized(d, fyne.NewSize(680, 520))
}

// endAdd closes the add workflow's backend, then reads the ledger again.
func (t *trustPage) endAdd() {
	link := t.link
	t.link = nil
	if link == nil {
		return
	}
	t.s.background(link.Close, t.probe)
}
