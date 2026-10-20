//go:build gui

package gui

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

var netModeHelp = map[workbench.GuardMode]string{
	workbench.GuardOff: "No network guard: activity is only observed.",
	workbench.GuardWhitelist: "Only the executables below may use the network (IPv4 and IPv6). " +
		"Every other program is cut off, root included.",
	workbench.GuardBlacklist: "The executables below are cut off the network. Every other program is allowed.",
}

// netGuardTab edits the network guard: `network guard`'s -w/-b, --auto-infra, --unsafe,
// --no-throttle and -e.
type netGuardTab struct {
	p    *networkPage
	root fyne.CanvasObject
	ed   *policyEditor

	autoInfra *handCheck
	unsafe    *handCheck
	throttle  *handCheck
	events    []ebpf.NetEventType
	eventsBox *handCheckGroup
}

func newNetGuardTab(p *networkPage) *netGuardTab {
	g := &netGuardTab{p: p}
	g.ed = newPolicyEditor(p.shell, policyHooks{
		help:    netModeHelp,
		status:  g.status,
		changed: g.syncOptions,
		apply:   g.confirmApply,
		revert:  g.refresh,
	})
	onCheck := func(bool) { g.ed.changed() }
	g.autoInfra = newCheck("Keep network daemons online", onCheck)
	g.unsafe = newCheck("Block local sockets too", onCheck)
	g.throttle = newCheck("Rate-limit events", onCheck)
	g.eventsBox = newNetEventChecks(nil, func(t []ebpf.NetEventType) { g.events = t; g.ed.changed() })
	g.eventsBox.Horizontal = false
	opts := newPanel(container.NewVBox(boldLabel("Options"),
		g.autoInfra, optionHint("Whitelist only. Admits the running DNS resolver and network manager, "+
			"without which name resolution fails for every program."),
		g.unsafe, optionHint("Whitelist only. Also guards AF_UNIX sockets (X11, Wayland, D-Bus, systemd): "+
			"the desktop may break."),
		g.throttle, optionHint("At most one event per type, verdict and program every 250 ms."),
		boldLabel("Report events"), g.eventsBox,
		optionHint("Only filters what is shown: every operation is enforced.")))
	g.root = g.ed.layout(opts)
	if !p.shell.status.canGuard() {
		for _, w := range []fyne.Disableable{g.autoInfra, g.unsafe, g.throttle, g.eventsBox} {
			w.Disable()
		}
	}
	g.refresh()
	return g
}

func (g *netGuardTab) draft() workbench.NetPolicy {
	return workbench.NetPolicy{
		Mode: g.ed.mode, Binaries: g.ed.entries(), AutoInfra: g.autoInfra.Checked, Unsafe: g.unsafe.Checked,
		Throttle: g.throttle.Checked, Events: slices.Clone(g.events),
	}
}

// syncOptions enables the whitelist-only options for a whitelist draft and clears them otherwise,
// as `network guard` refuses them with -b.
func (g *netGuardTab) syncOptions() {
	if g.autoInfra == nil {
		return // still being built
	}
	if g.ed.mode == workbench.GuardWhitelist && g.p.shell.status.canGuard() {
		g.autoInfra.Enable()
		g.unsafe.Enable()
		return
	}
	g.autoInfra.Checked, g.unsafe.Checked = false, false
	g.autoInfra.Disable() // Disable refreshes: the cleared boxes repaint
	g.unsafe.Disable()
}

func (g *netGuardTab) status() (string, bool) {
	enf := g.p.n.Policy()
	text := "Enforced now: no network guard."
	if enf.Mode != workbench.GuardOff {
		text = fmt.Sprintf("Enforced now: %s with %d executable(s)", enf.Mode, len(enf.Binaries))
		if x := len(g.p.n.Extras()); x > 0 {
			text += fmt.Sprintf(" and %d added automatically", x)
		}
		text += "."
	}
	if g.autoInfra == nil {
		return text, false
	}
	d := g.draft()
	if d.Mode == workbench.GuardOff {
		return text, enf.Mode != workbench.GuardOff
	}
	dirty := d.Mode != enf.Mode || !sameEntries(d.Binaries, enf.Binaries) || d.AutoInfra != enf.AutoInfra ||
		d.Unsafe != enf.Unsafe || d.Throttle != enf.Throttle || !slices.Equal(d.Events, enf.Events)
	return text, dirty
}

// refresh resets the draft to the enforced policy.
func (g *netGuardTab) refresh() {
	enf := g.p.n.Policy()
	g.ed.bins.set(enf.Binaries, g.p.n.Extras())
	g.autoInfra.Checked, g.unsafe.Checked, g.throttle.Checked = enf.AutoInfra, enf.Unsafe, enf.Throttle
	g.events = slices.Clone(enf.Events)
	g.eventsBox.Selected = nil
	for _, t := range eventsOrAll(enf.Events) {
		g.eventsBox.Selected = append(g.eventsBox.Selected, t.String())
	}
	for _, w := range []fyne.Widget{g.autoInfra, g.unsafe, g.throttle, g.eventsBox} {
		w.Refresh()
	}
	g.ed.setMode(enf.Mode)
}

func eventsOrAll(t []ebpf.NetEventType) []ebpf.NetEventType {
	if len(t) == 0 {
		return ebpf.NetEventTypes()
	}
	return t
}

func (g *netGuardTab) confirmApply() {
	d := g.draft()
	confirmWarnings(g.p.shell, "Apply network guard", netWarnings(d), func() { g.applyDraft(d) })
}

func netWarnings(d workbench.NetPolicy) []string {
	var w []string
	switch d.Mode {
	case workbench.GuardOff:
		return []string{"The network guard stops: every program may use the network again."}
	case workbench.GuardWhitelist:
		if len(d.Binaries) == 0 {
			w = append(w, "The whitelist is empty: no program may use the network.")
		} else {
			w = append(w, "Every program not listed loses network access, root included.")
		}
		if !d.AutoInfra {
			w = append(w, "System network daemons are not kept online: name resolution through "+
				"systemd-resolved fails for every program, listed ones included.")
		}
		w = append(w, "Listed programs already running are judged once from /proc: restart them for full "+
			"code-integrity checks.")
	}
	if d.Unsafe {
		w = append(w, "Local sockets are blocked too: X11, Wayland, D-Bus and systemd connections of every "+
			"program not listed fail, and the desktop may break. This window is added to the whitelist to keep "+
			"its display. If it still stops responding, press Ctrl+C in the terminal that started it: the "+
			"guard stops with it.")
	}
	return w
}

func (g *netGuardTab) applyDraft(d workbench.NetPolicy) {
	g.p.shell.busy("Applying network guard…", func() error {
		return g.p.n.Apply(d)
	}, func() {
		g.refresh()
		g.p.refresh(true)
	})
}

// promote adds the executable a process ran to a mode's list, after showing what changes.
func (g *netGuardTab) promote(exe workbench.Exe, comm string, mode workbench.GuardMode) {
	g.p.shell.verifyExe(exe, func(entry ebpf.BinaryEntry) {
		enf := g.p.n.Policy()
		next := g.draft()
		if next.Mode != mode {
			next = workbench.NetPolicy{Mode: mode, Throttle: true}
			if enf.Mode == mode {
				next = enf
			}
		}
		next.Binaries = g.ed.promoteBase(mode, enf.Mode, enf.Binaries, entry)
		effect := "It may use the network."
		if mode == workbench.GuardBlacklist {
			effect = "It loses network access."
		}
		how := "The network guard is replaced: the new one starts before the old one stops, so nothing " +
			"either policy blocks gets through meanwhile."
		if enf.Mode == workbench.GuardOff {
			how = "This starts the network guard."
			if mode == workbench.GuardWhitelist {
				how += " Every program not listed loses network access."
			}
		}
		note := effect + " " + how
		if ws := netWarnings(next); len(ws) > 0 {
			note += "\n\n• " + strings.Join(ws, "\n\n• ")
		}
		g.p.shell.confirmExecutable(verbFor(mode), entry, exe, comm, note, func() {
			g.applyDraft(next)
		})
	})
}

func optionHint(s string) *widget.Label {
	l := mutedText(s)
	l.SizeName = theme.SizeNameCaptionText
	return l
}
