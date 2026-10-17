//go:build gui

package gui

import (
	"fmt"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// networkPage is the Network workspace: the network monitor's watched programs and the host-wide
// network guard, over one merged event feed.
type networkPage struct {
	shell *shell
	n     *workbench.Network

	monPill   *pill
	guardPill *pill
	pills     *fyne.Container
	errLabel  *widget.Label
	counters  *counterStrip
	feed      *eventFeed
	procs     *processView
	watch     *watchTab
	guard     *netGuardTab
	tabs      *container.AppTabs

	countsVersion uint64
}

func newNetworkPage(s *shell, n *workbench.Network) *networkPage {
	return &networkPage{shell: s, n: n, monPill: newPill("", colorMuted), guardPill: newPill("", colorMuted)}
}

func (p *networkPage) view() fyne.CanvasObject {
	p.errLabel = widget.NewLabel("")
	p.errLabel.Importance = widget.DangerImportance
	p.errLabel.Wrapping = fyne.TextWrapWord
	p.errLabel.Hide()
	p.counters = newCounterStrip("Blocked")
	p.guard = newNetGuardTab(p)
	p.watch = newWatchTab(p)
	p.feed = p.newFeed()
	p.procs = p.newProcesses()
	p.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Watched", icon("eye"), p.watch.root),
		container.NewTabItemWithIcon("Activity", icon("pulse"), p.feed.root),
		container.NewTabItemWithIcon("Processes", icon("cpu"), p.procs.root),
		container.NewTabItemWithIcon("Guard", icon("shield-check"), p.guard.root),
	)
	p.pills = container.NewHBox(p.monPill, p.guardPill)
	top := newGapVBox(p.counters.root, p.pills, p.errLabel)
	p.refresh(true)
	p.shell.every(pollInterval, func() { p.refresh(false) })
	return container.NewBorder(top, nil, nil, nil, p.tabs)
}

// searchEntry is the filter of the visible table tab, nil on the others.
func (p *networkPage) searchEntry() *widget.Entry {
	switch p.tabs.SelectedIndex() {
	case 1:
		return p.feed.search
	case 2:
		return p.procs.search
	}
	return nil
}

func (p *networkPage) refresh(force bool) {
	mon, mc := "not monitoring", colorMuted
	if n := len(p.n.Watched()); n > 0 {
		mon, mc = fmt.Sprintf("monitoring %d program(s)", n), colorInfo
	}
	text, c := netGuardPill(p.n.Policy())
	if mon != p.monPill.text || text != p.guardPill.text {
		p.monPill.Set(mon, mc)
		p.guardPill.Set(text, c)
		p.pills.Refresh() // the pills' widths changed: lay the row out again
	}
	if e := p.n.LastError(); e != "" {
		p.errLabel.SetText("Network guard error: " + e)
		p.errLabel.Show()
	} else {
		p.errLabel.Hide()
	}
	counts, version := p.n.Counts()
	if force || version != p.countsVersion {
		p.countsVersion = version
		p.counters.set(counts, len(p.n.Processes()))
		p.feed.refresh(force)
		p.procs.refresh()
	}
}

func netGuardPill(pol workbench.NetPolicy) (string, fyne.ThemeColorName) {
	switch pol.Mode {
	case workbench.GuardWhitelist:
		s := fmt.Sprintf("network whitelist · %d", len(pol.Binaries))
		if pol.Unsafe {
			s += " · local sockets"
		}
		return s, theme.ColorNameSuccess
	case workbench.GuardBlacklist:
		return fmt.Sprintf("network blacklist · %d", len(pol.Binaries)), theme.ColorNameWarning
	default:
		return "network unguarded", colorMuted
	}
}

func (p *networkPage) actions() []feedAction {
	return []feedAction{
		{label: "Watch", icon: "eye", importance: widget.MediumImportance, needs: sysStatus.canMonitor,
			run: func(e workbench.Exe, _ string) { p.watch.watchExe(e) }},
		{label: "Whitelist", icon: "shield-check", importance: widget.SuccessImportance,
			run: func(e workbench.Exe, comm string) { p.guard.promote(e, comm, workbench.GuardWhitelist) }},
		{label: "Blacklist", icon: "prohibit", importance: widget.DangerImportance,
			run: func(e workbench.Exe, comm string) { p.guard.promote(e, comm, workbench.GuardBlacklist) }},
	}
}

func (p *networkPage) newFeed() *eventFeed {
	return newEventFeed(p.shell, &feedSpec{
		cols: []column{
			{"Time", 104}, {"Verdict", 86}, {"Event", 96}, {"Process", 130}, {"PID", 72}, {"Proto", 64},
			{"Local", 170}, {"Remote", 200}, {"Size", 80}, {"Executable", 280},
		},
		cell:        netCell,
		source:      p.n.Rows,
		clear:       p.n.Clear,
		placeholder: "Search events",
		hint:        "Select an event to watch, whitelist or blacklist the executable behind it.",
		actions:     p.actions(),
		exportName:  "network-events",
	})
}

func netCell(r *workbench.Row, col int) cell {
	switch col {
	case 0:
		return timeCell(r)
	case 1:
		return verdictCell(r.Verdict)
	case 2:
		return cell{text: r.Op}
	case 3:
		return cell{text: r.Comm, bold: true}
	case 4:
		return pidCell(r)
	case 5:
		return cell{text: r.Proto, importance: widget.LowImportance}
	case 6:
		return cell{text: r.Local}
	case 7:
		return cell{text: r.Remote}
	case 8:
		return sizeCell(r)
	default:
		if r.Detail != "" {
			return cell{text: r.Detail, importance: widget.WarningImportance}
		}
		return cell{text: exeText(r.Exe), importance: widget.LowImportance}
	}
}

// sizeCell is the byte count of SEND/RECV and the backlog of LISTEN.
func sizeCell(r *workbench.Row) cell {
	switch r.Op {
	case ebpf.NetSend.String(), ebpf.NetRecv.String():
		return cell{text: strconv.FormatUint(uint64(r.Size), 10) + " B"}
	case ebpf.NetListen.String():
		return cell{text: "backlog " + strconv.FormatUint(uint64(r.Size), 10), importance: widget.LowImportance}
	}
	return cell{}
}

func (p *networkPage) newProcesses() *processView {
	cols := append(append([]column(nil), procColumns...), column{"Watched", 90}, column{"Guard", 110})
	return newProcessView(p.shell, &procSpec{
		cols: cols,
		cell: func(pr *workbench.Process, col int) cell {
			if c, ok := procCell(pr, col); ok {
				return c
			}
			if col == len(procColumns) {
				if pr.Exe.Path != "" && slices.Contains(entryPaths(p.n.Watched()), pr.Exe.Path) {
					return cell{text: "watched", importance: widget.HighImportance}
				}
				return cell{text: "—", importance: widget.LowImportance}
			}
			pol := p.n.Policy()
			return policyMark(pol.Mode, entryPaths(pol.Binaries), pr.Exe)
		},
		source: p.n.Processes,
		hint: "Every program seen on the network, identified by its inode. Watch one to trace all of its " +
			"connections, or put it in the network guard's policy.",
		actions:    p.actions(),
		exportName: "network-processes",
	})
}

// watchTab edits the programs the network monitor follows (network-monitor's binaries and -e).
type watchTab struct {
	p    *networkPage
	root fyne.CanvasObject
	bins *binaryList
}

func newWatchTab(p *networkPage) *watchTab {
	w := &watchTab{p: p}
	w.bins = newBinaryList(p.shell, w.apply)
	w.bins.set(p.n.Watched(), nil)
	if !p.shell.status.canMonitor() {
		w.bins.setEditable(false)
	}
	events := newNetEventChecks(p.n.MonitorEvents(), func(types []ebpf.NetEventType) { p.n.SetMonitorEvents(types) })
	opts := newPanel(container.NewVBox(boldLabel("Show events"), events,
		mutedText("Only the programs listed below are traced. Changes apply at once.")))
	w.root = container.NewBorder(opts, nil, nil, nil, w.bins.card("Watched programs"))
	return w
}

// apply restarts the monitor on the edited list; a failure puts the running list back.
func (w *watchTab) apply() {
	next := w.bins.entries
	w.p.shell.busy("Updating the network monitor…", func() error {
		return w.p.n.Watch(next)
	}, func() {
		w.bins.set(w.p.n.Watched(), nil)
		w.p.refresh(true)
	})
}

func (w *watchTab) watchExe(e workbench.Exe) {
	w.p.shell.verifyExe(e, func(entry ebpf.BinaryEntry) { w.bins.put(entry) })
}

// newNetEventChecks is a horizontal check group of the network event types (nil = all, as -e). The
// last checked box can't be cleared: no event type at all is not a filter -e can express.
func newNetEventChecks(selected []ebpf.NetEventType, onChange func([]ebpf.NetEventType)) *widget.CheckGroup {
	all := ebpf.NetEventTypes()
	names := make([]string, len(all))
	for i, t := range all {
		names[i] = t.String()
	}
	g := widget.NewCheckGroup(names, nil)
	g.Horizontal = true
	if len(selected) == 0 {
		selected = all
	}
	for _, t := range selected {
		g.Selected = append(g.Selected, t.String())
	}
	prev := slices.Clone(g.Selected)
	g.OnChanged = func(sel []string) {
		if len(sel) == 0 {
			g.SetSelected(prev)
			return
		}
		prev = slices.Clone(sel)
		var types []ebpf.NetEventType
		for _, s := range sel {
			if t, ok := ebpf.ParseNetEventType(s); ok {
				types = append(types, t)
			}
		}
		if len(types) == len(all) {
			types = nil
		}
		onChange(types)
	}
	return g
}
