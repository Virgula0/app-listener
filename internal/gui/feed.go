//go:build gui

package gui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/workbench"
)

// feedAction is a button acting on the selected event's or process's executable.
type feedAction struct {
	label, icon string
	importance  widget.Importance
	// needs gates the button on the host (default: guards available).
	needs func(sysStatus) bool
	run   func(exe workbench.Exe, comm string)
}

func (a *feedAction) button() *handButton {
	b := newButton(a.label, icon(a.icon), nil)
	b.Importance = a.importance
	b.Disable()
	return b
}

func (a *feedAction) allowed(st sysStatus) bool {
	if a.needs == nil {
		return st.canGuard()
	}
	return a.needs(st)
}

// actionBar is the strip under a table: a hint or the selection, and the actions.
type actionBar struct {
	shell   *shell
	root    fyne.CanvasObject
	info    *widget.Label
	hint    string
	actions []feedAction
	buttons []*handButton
}

func newActionBar(s *shell, hint string, actions []feedAction) *actionBar {
	b := &actionBar{shell: s, info: mutedText(hint), hint: hint, actions: actions}
	objs := make([]fyne.CanvasObject, len(actions))
	for i := range actions {
		b.buttons = append(b.buttons, actions[i].button())
		objs[i] = b.buttons[i]
	}
	p := newPanel(container.NewBorder(nil, nil, nil, container.NewHBox(objs...), b.info))
	p.inset = 10
	b.root = p
	return b
}

// set shows info for a selected executable (nil exe: nothing selected) and gates the actions.
func (b *actionBar) set(exe *workbench.Exe, comm, info string) {
	if exe == nil {
		b.info.SetText(b.hint)
	} else {
		b.info.SetText(info)
	}
	for i, btn := range b.buttons {
		a := b.actions[i]
		if exe == nil || !exe.Known() || !a.allowed(b.shell.status) {
			btn.Disable()
			continue
		}
		e, c := *exe, comm
		btn.OnTapped = func() { a.run(e, c) }
		btn.Enable()
	}
}

// feedSpec configures an eventFeed.
type feedSpec struct {
	cols        []column
	cell        func(r *workbench.Row, col int) cell
	source      func(workbench.Filter) ([]workbench.Row, uint64)
	clear       func()
	placeholder string
	hint        string
	actions     []feedAction
	// exportName names exported files ("files-events").
	exportName string
}

// eventFeed is a filterable, pausable event table with actions on the selected event's executable.
type eventFeed struct {
	spec    feedSpec
	root    fyne.CanvasObject
	search  *widget.Entry
	table   *dataTable
	bar     *actionBar
	rows    []workbench.Row
	filter  workbench.Filter
	version uint64
	paused  bool
	sel     *workbench.Row
}

var verdictFilters = map[string][]workbench.Verdict{
	"All verdicts": nil,
	"Denied":       {workbench.Denied},
	"Allowed":      {workbench.Allowed},
	"Observed":     {workbench.Observed},
}

func newEventFeed(s *shell, spec *feedSpec) *eventFeed {
	f := &eventFeed{spec: *spec}
	f.table = newDataTable(spec.cols, func() int { return len(f.rows) },
		func(row, col int) cell { return spec.cell(&f.rows[row], col) }, f.onSelect)

	f.search = newSearchEntry(spec.placeholder, func(s string) { f.filter.Text = s; f.refresh(true) })

	verdict := widget.NewSelect([]string{"All verdicts", "Denied", "Allowed", "Observed"}, func(s string) {
		f.filter.Verdicts = verdictFilters[s]
		f.refresh(true)
	})
	// Assigned, not SetSelected: the callback would fire before the feed is built.
	verdict.Selected = "All verdicts"

	var pause *handButton
	pause = newButton("Pause", icon("pause"), func() {
		f.paused = !f.paused
		if f.paused {
			pause.SetText("Resume")
			pause.SetIcon(icon("play"))
			return
		}
		pause.SetText("Pause")
		pause.SetIcon(icon("pause"))
		f.refresh(true)
	})
	clearBtn := newButton("Clear", icon("broom"), func() { spec.clear(); f.refresh(true) })
	export := newButton("Export", icon("export"), func() { s.exportTable(spec.exportName, f.table) })
	toolbar := container.NewBorder(nil, nil, nil, container.NewHBox(verdict, pause, clearBtn, export), f.search)

	f.bar = newActionBar(s, spec.hint, spec.actions)
	f.root = container.NewBorder(container.NewPadded(toolbar), f.bar.root, nil, nil, f.table)
	return f
}

func (f *eventFeed) refresh(force bool) {
	if f.paused && !force {
		return
	}
	rows, version := f.spec.source(f.filter)
	if !force && version == f.version {
		return
	}
	f.rows, f.version = rows, version
	if f.sel != nil && !force {
		// Keep the chosen event while new ones scroll in.
		f.table.Refresh()
		return
	}
	f.sel = nil
	f.table.reset()
	f.bar.set(nil, "", "")
}

func (f *eventFeed) onSelect(row int) {
	if row < 0 || row >= len(f.rows) {
		return
	}
	r := f.rows[row]
	f.sel = &r
	f.bar.set(&r.Exe, r.Comm, fmt.Sprintf("%s (pid %d) · %s · %s", r.Comm, r.PID, exeText(r.Exe),
		r.Time.Format(time.TimeOnly)))
}

// newSearchEntry is a filter field; Ctrl+F focuses the visible one (shell.focusSearch).
func newSearchEntry(placeholder string, onChanged func(string)) *widget.Entry {
	e := widget.NewEntry()
	e.SetPlaceHolder(placeholder)
	e.ActionItem = widget.NewIcon(icon("magnifying-glass"))
	e.OnChanged = onChanged
	return e
}

// procSpec configures a processView.
type procSpec struct {
	cols       []column
	cell       func(p *workbench.Process, col int) cell
	source     func() []workbench.Process
	hint       string
	actions    []feedAction
	exportName string
}

// processView lists the executables seen, one row per inode, with actions on the selected one.
type processView struct {
	spec   procSpec
	root   fyne.CanvasObject
	search *widget.Entry
	table  *dataTable
	bar    *actionBar
	procs  []workbench.Process
	filter string
}

func newProcessView(s *shell, spec *procSpec) *processView {
	p := &processView{spec: *spec}
	p.table = newDataTable(spec.cols, func() int { return len(p.procs) },
		func(row, col int) cell { return spec.cell(&p.procs[row], col) },
		func(row int) {
			if row >= 0 && row < len(p.procs) {
				pr := p.procs[row]
				p.bar.set(&pr.Exe, pr.Comm, fmt.Sprintf("%s · %s", pr.Comm, exeText(pr.Exe)))
			}
		})
	p.bar = newActionBar(s, spec.hint, spec.actions)
	p.search = newSearchEntry("Search processes", func(text string) {
		p.filter = strings.ToLower(text)
		p.refresh()
	})
	export := newButton("Export", icon("export"), func() { s.exportTable(spec.exportName, p.table) })
	toolbar := container.NewBorder(nil, nil, nil, export, p.search)
	p.root = container.NewBorder(container.NewPadded(toolbar), p.bar.root, nil, nil, p.table)
	return p
}

func (p *processView) refresh() {
	all := p.spec.source()
	if p.filter == "" {
		p.procs = all
	} else {
		p.procs = nil
		for _, pr := range all {
			if strings.Contains(strings.ToLower(pr.Comm), p.filter) ||
				strings.Contains(strings.ToLower(exeText(pr.Exe)), p.filter) {
				p.procs = append(p.procs, pr)
			}
		}
	}
	p.table.Refresh()
}

// procColumns are the leading columns every process table shares (procCell).
var procColumns = []column{
	{"Process", 140}, {"Executable", 360}, {"Inode", 110}, {"Events", 80}, {"Denied", 80}, {"Last seen", 100},
}

// procCell renders procColumns; ok is false past them.
func procCell(pr *workbench.Process, col int) (c cell, ok bool) {
	switch col {
	case 0:
		return cell{text: pr.Comm, bold: true}, true
	case 1:
		return cell{text: exeText(pr.Exe)}, true
	case 2:
		if pr.Exe.Known() {
			return cell{text: strconv.FormatUint(pr.Exe.Key.Ino, 10), importance: widget.LowImportance}, true
		}
		return cell{text: "?", importance: widget.LowImportance}, true
	case 3:
		return cell{text: strconv.FormatUint(pr.Events, 10)}, true
	case 4:
		if pr.Denied > 0 {
			return cell{text: strconv.FormatUint(pr.Denied, 10), importance: widget.DangerImportance, bold: true}, true
		}
		return cell{text: "0", importance: widget.LowImportance}, true
	case 5:
		return cell{text: pr.LastSeen.Format(time.TimeOnly)}, true
	}
	return cell{}, false
}

func timeCell(r *workbench.Row) cell { return cell{text: r.Time.Format("15:04:05.000")} }

func pidCell(r *workbench.Row) cell { return cell{text: strconv.FormatUint(uint64(r.PID), 10)} }

func verdictCell(v workbench.Verdict) cell {
	switch v {
	case workbench.Denied:
		return cell{text: "DENIED", importance: widget.DangerImportance, bold: true}
	case workbench.Allowed:
		return cell{text: "allowed", importance: widget.SuccessImportance}
	default:
		return cell{text: "observed", importance: widget.LowImportance}
	}
}

func exeText(e workbench.Exe) string {
	switch {
	case e.Path != "":
		return e.Path
	case e.Known():
		return fmt.Sprintf("inode %d", e.Key.Ino)
	default:
		return "unknown"
	}
}

// policyMark is how a process table shows an executable's place in a policy.
func policyMark(mode workbench.GuardMode, bins []string, exe workbench.Exe) cell {
	for _, b := range bins {
		if b == exe.Path && exe.Path != "" {
			if mode == workbench.GuardWhitelist {
				return cell{text: "whitelisted", importance: widget.SuccessImportance, bold: true}
			}
			return cell{text: "blacklisted", importance: widget.DangerImportance, bold: true}
		}
	}
	return cell{text: "—", importance: widget.LowImportance}
}
