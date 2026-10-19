//go:build gui

package gui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// forwardEvents are the event types a temporary rule may name (`edit-protected forward -e`).
var forwardEvents = []string{"OPEN", "READ", "WRITE", "DELETE", "RENAME", "SYMLINK", "HARDLINK", "MKDIR", "MMAP",
	"ATTR", "STAT", "MKNOD"}

// forwardForm is the setup of `edit-protected forward -w/-b [-e]`.
type forwardForm struct {
	p      *protectedPage
	bins   []string
	list   *fyne.Container
	path   *widget.Entry
	rule   *handRadioGroup
	events *checkGrid
}

const (
	ruleAdmit = "Admit (-w): these binaries may reach the directories"
	ruleDeny  = "Deny (-b): these whitelisted binaries may not"
)

func (p *protectedPage) startForward() {
	if p.hello != nil && !p.hello.Live {
		showInformation("Temporary access", "Temporary access needs the daemon running with an "+
			"edit-protected password configured (live mode).", p.s.win)
		return
	}
	f := &forwardForm{p: p, list: container.NewVBox()}
	f.path = widget.NewEntry()
	f.path.SetPlaceHolder("/usr/bin/vim")
	f.path.OnSubmitted = func(string) { f.addTyped() }
	f.rule = newRadioGroup([]string{ruleAdmit, ruleDeny}, nil)
	f.rule.Required = true
	f.rule.SetSelected(ruleAdmit)
	f.events = newCheckGrid(forwardEvents)
	f.refreshList()

	add := newButton("Add", icon("plus"), f.addTyped)
	browse := newButton("Browse…", icon("file"), func() {
		p.s.pickPath(pickFile, "Binary to grant", "", f.add)
	})
	start := newButton("Authenticate and grant", icon("lock-simple-open"), f.submit)
	start.Importance = widget.HighImportance
	back := newButton("Back", nil, func() { p.show(p.over) })
	all := newTextButton("All events", func() { f.setEvents(true) })
	none := newTextButton("None", func() { f.setEvents(false) })

	form := container.NewVBox(
		sectionTitle("BINARIES"),
		newPanel(container.NewVBox(f.list, container.NewBorder(nil, nil, nil, container.NewHBox(add, browse), f.path))),
		sectionTitle("RULE"),
		newPanel(container.NewVBox(f.rule, mutedText("An admit grant applies to every user running the binary, "+
			"and admitting an interpreter (python, node, bash) admits every script it runs."))),
		sectionTitle("EVENTS"),
		newPanel(container.NewVBox(f.events.root, container.NewHBox(all, none),
			mutedText("The event types the rule covers. All checked = every event."))),
		container.NewHBox(back, start),
	)
	p.show(newFadeVScroll(form))
}

func (f *forwardForm) setEvents(on bool) {
	for _, c := range f.events.checks {
		c.SetChecked(on)
	}
}

func (f *forwardForm) addTyped() {
	if t := strings.TrimSpace(f.path.Text); t != "" {
		f.add(t)
		f.path.SetText("")
	}
}

func (f *forwardForm) add(p string) {
	if !filepath.IsAbs(p) {
		showInformation("Binary", p+" is not an absolute path.", f.p.s.win)
		return
	}
	for _, b := range f.bins {
		if b == p {
			return
		}
	}
	f.bins = append(f.bins, p)
	f.refreshList()
}

func (f *forwardForm) refreshList() {
	f.list.Objects = nil
	if len(f.bins) == 0 {
		f.list.Add(mutedText("No binaries yet: type a path or browse for one."))
	}
	for i, b := range f.bins {
		rm := newButton("", icon("x"), func() {
			f.bins = append(f.bins[:i], f.bins[i+1:]...)
			f.refreshList()
		})
		rm.Importance = widget.LowImportance
		f.list.Add(container.NewBorder(nil, nil, newIconTile("terminal-window", colorMuted, 26), rm, selectable(b)))
	}
	f.list.Refresh()
}

// isOwnExecutable: admitting this window's own binary would hand the GUI build, with every
// subcommand in it, the access this design keeps in the headless backend.
func isOwnExecutable(p string) bool {
	sd, si, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return true
	}
	d, i, err := ebpf.StatInode(p)
	return err == nil && d == sd && i == si
}

func (f *forwardForm) selectedEvents() []string {
	sel := f.events.selected()
	if len(sel) == len(forwardEvents) {
		return nil
	}
	return sel
}

func (f *forwardForm) submit() {
	p := f.p
	if len(f.bins) == 0 {
		showInformation("Temporary access", "Add at least one binary.", p.s.win)
		return
	}
	events := f.selectedEvents()
	if events != nil && len(events) == 0 {
		showInformation("Temporary access", "Check at least one event type.", p.s.win)
		return
	}
	allow := f.rule.Selected == ruleAdmit
	if allow {
		for _, b := range f.bins {
			if isOwnExecutable(b) {
				showInformation("Temporary access", b+" is this window's own binary: it is never "+
					"admitted. Its file access goes through the installed headless binary instead.", p.s.win)
				return
			}
		}
	}
	bins := append([]string(nil), f.bins...)
	// The backend checks binaries only after AUTH: checking opens them as the daemon's binary.
	p.begin(func(link backendLink) {
		p.authenticate(link, "Temporary access", func(resources []string) {
			var m editbackend.Message
			var err error
			p.s.background(func() {
				m, err = link.Call(&editbackend.Request{Op: editbackend.OpForwardCheck, Binaries: bins})
			}, func() {
				if err != nil {
					p.s.showError("Temporary access", err)
					p.finish()
					return
				}
				f.confirmChecks(m.Binaries, allow, func() {
					p.pickMany("Temporary access", resources, func(picked []string) {
						f.grant(link, picked, bins, allow, events)
					})
				})
			})
		})
	})
}

// confirmChecks mirrors confirmUserPlaced: unusable binaries stop here; admitting one that isn't
// root-placed asks first; launch-injectable runtimes are noted.
func (f *forwardForm) confirmChecks(checks []editbackend.BinaryCheck, allow bool, ok func()) {
	p := f.p
	var bad, risky, notes []string
	for _, c := range checks {
		switch {
		case c.Err != "":
			bad = append(bad, c.Path+": "+c.Err)
		case !c.SystemPlaced:
			risky = append(risky, c.Path)
		}
		if len(c.Runtimes) > 0 {
			notes = append(notes, fmt.Sprintf("%s is a %s runtime: started with code-loading env or flags "+
				"(NODE_OPTIONS, ELECTRON_RUN_AS_NODE, --inspect, JAVA_TOOL_OPTIONS, …) it is refused the "+
				"directories.", c.Path, strings.Join(c.Runtimes, "/")))
		}
	}
	if len(bad) > 0 {
		p.s.showError("Temporary access", fmt.Errorf("%s", strings.Join(bad, "\n")))
		p.finish()
		return
	}
	if !allow || len(risky) == 0 {
		if len(notes) > 0 {
			showInformation("Temporary access", strings.Join(notes, "\n"), p.s.win)
		}
		ok()
		return
	}
	msg := "Not root-placed (a user can rewrite them while the grant is active; it is revoked if their " +
		"content changes):\n  " + strings.Join(risky, "\n  ") +
		"\nThe grant covers every user running them, and an interpreter (python, bash, node) admits " +
		"every script it is given."
	if len(notes) > 0 {
		msg += "\n\n" + strings.Join(notes, "\n")
	}
	d := dialog.NewCustomConfirm("Temporary access", "Admit them", "Cancel",
		scrollable(promptBody("Admit binaries that are not root-placed?", msg)), func(yes bool) {
			if yes {
				ok()
			} else {
				p.finish()
			}
		}, p.s.win)
	p.s.showSized(d, fyne.NewSize(720, 360))
}

func (f *forwardForm) grant(link backendLink, resources, bins []string, allow bool, events []string) {
	p := f.p
	var m editbackend.Message
	var err error
	p.s.background(func() {
		m, err = link.Call(&editbackend.Request{Op: editbackend.OpForward, Resources: resources, Binaries: bins,
			Allow: allow, Events: events})
	}, func() {
		if err != nil {
			p.s.showError("Temporary access", err)
			p.finish()
			return
		}
		v := newForwardView(p, link, resources, bins, allow, events)
		p.leave = v.stop
		p.show(v.root)
		v.start()
		if len(m.Notes) > 0 {
			showInformation("Temporary access", strings.Join(m.Notes, "\n"), p.s.win)
		}
	})
}

// forwardView shows an active grant's events until it is stopped or the daemon ends it.
type forwardView struct {
	p      *protectedPage
	link   backendLink
	rows   []editbackend.Event
	table  *dataTable
	count  *widget.Label
	banner *fyne.Container
	root   fyne.CanvasObject
	ended  bool
	last   time.Time
}

const maxForwardRows = 5000

func newForwardView(p *protectedPage, link backendLink, resources, bins []string, allow bool, events []string) *forwardView {
	v := &forwardView{p: p, link: link, count: mutedText("0 events"), banner: container.NewVBox()}
	verb, c := "ALLOWED", theme.ColorNameSuccess
	if !allow {
		verb, c = "DENIED", theme.ColorNameError
	}
	evs := "ALL"
	if len(events) > 0 {
		evs = strings.Join(events, ",")
	}
	stop := newButton("Stop and revoke", icon("stop"), v.stop)
	stop.Importance = widget.DangerImportance
	head := newPanel(container.NewVBox(
		container.NewBorder(nil, nil, newPill(verb, c), stop, boldLabel(strings.Join(bins, ", "))),
		mutedText(fmt.Sprintf("on %s (events %s) — revoked when you stop it, or after %s without activity "+
			"(a granted binary touching the directories, or your input here).", strings.Join(resources, ", "),
			evs, p.timeout.Selected)),
	))
	v.table = newDataTable([]column{{"Time", 110}, {"PID", 70}, {"UID", 60}, {"Process", 140}, {"Op", 90},
		{"Path", 420}, {"Verdict", 90}}, func() int { return len(v.rows) }, v.cell, func(int) { v.ping() })
	v.root = container.NewBorder(container.NewVBox(head, v.banner), v.count, nil, nil, v.table)
	return v
}

func (v *forwardView) cell(row, col int) cell {
	e := v.rows[len(v.rows)-1-row]
	switch col {
	case 0:
		return cell{text: e.Time.Format("15:04:05.000")}
	case 1:
		return cell{text: strconv.FormatUint(uint64(e.PID), 10)}
	case 2:
		return cell{text: strconv.FormatUint(uint64(e.UID), 10)}
	case 3:
		return cell{text: e.Comm}
	case 4:
		return cell{text: e.Op}
	case 5:
		if e.Dest != "" {
			return cell{text: e.Path + " → " + e.Dest}
		}
		return cell{text: e.Path}
	default:
		if e.Allowed {
			return cell{text: "allowed", importance: widget.SuccessImportance}
		}
		return cell{text: "denied", importance: widget.DangerImportance, bold: true}
	}
}

func (v *forwardView) start() {
	v.p.watchEnded(v.link, v.onEnded)
	if v.p.s.syncBackground {
		return
	}
	link := v.link
	go func() {
		for {
			select {
			case m := <-link.Events():
				if m.Event == nil {
					continue
				}
				ev := *m.Event
				fyne.Do(func() { v.add(&ev) })
			case <-link.Done():
				return
			}
		}
	}()
}

func (v *forwardView) add(ev *editbackend.Event) {
	v.rows = append(v.rows, *ev)
	if len(v.rows) > maxForwardRows {
		v.rows = v.rows[len(v.rows)-maxForwardRows:]
	}
	v.count.SetText(fmt.Sprintf("%d events", len(v.rows)))
	v.table.Refresh()
}

func (v *forwardView) ping() {
	if v.ended || time.Since(v.last) < 3*time.Second {
		return
	}
	v.last = time.Now()
	link := v.link
	v.p.s.background(func() { _, _ = link.Call(&editbackend.Request{Op: editbackend.OpPing}) }, func() {})
}

func (v *forwardView) onEnded(why string) {
	v.ended = true
	v.banner.Objects = []fyne.CanvasObject{endedBanner(why)}
	v.banner.Refresh()
}

func (v *forwardView) stop() {
	link := v.link
	var err error
	v.p.s.background(func() { _, err = link.Call(&editbackend.Request{Op: editbackend.OpEnd}) }, func() {
		v.p.finish()
		if err != nil {
			v.p.s.showError("Revoking the temporary access", err)
		}
	})
}
