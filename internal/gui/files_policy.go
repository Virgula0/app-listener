//go:build gui

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

var fileModeHelp = map[workbench.GuardMode]string{
	workbench.GuardOff:       "No guard: activity is only observed.",
	workbench.GuardWhitelist: "Only the executables below may touch this target. Every other process is denied, root included.",
	workbench.GuardBlacklist: "The executables below are denied. Every other process is allowed.",
}

// policyTab is a target's Guard tab.
type policyTab struct {
	v    *targetView
	root fyne.CanvasObject
	ed   *policyEditor
}

func newPolicyTab(v *targetView) *policyTab {
	p := &policyTab{v: v}
	p.ed = newPolicyEditor(v.page.shell, policyHooks{
		help:   fileModeHelp,
		status: p.status,
		apply:  p.confirmApply,
		revert: p.refresh,
	})
	p.root = p.ed.layout(nil)
	p.refresh()
	return p
}

func (p *policyTab) draft() workbench.Policy {
	return workbench.Policy{Mode: p.ed.mode, Binaries: p.ed.entries()}
}

func (p *policyTab) status() (string, bool) {
	enf := p.v.t.Policy()
	text := "Enforced now: no guard."
	if enf.Mode != workbench.GuardOff {
		text = fmt.Sprintf("Enforced now: %s with %d executable(s).", enf.Mode, len(enf.Binaries))
	}
	d := p.draft()
	return text, d.Mode != enf.Mode || !sameEntries(d.Binaries, enf.Binaries)
}

// refresh resets the draft to the enforced policy.
func (p *policyTab) refresh() {
	enf := p.v.t.Policy()
	p.ed.bins.set(enf.Binaries, nil)
	p.ed.setMode(enf.Mode)
}

func (p *policyTab) confirmApply() {
	draft := p.draft()
	var warnings []string
	if p.v.t.ApplyLeavesGap(draft) {
		warnings = append(warnings, "Leaving a whitelist stops it before the new guard can scan the tree, so "+
			p.v.t.Spec.Path+" is unprotected for a moment.")
	}
	if draft.Mode == workbench.GuardWhitelist && len(draft.Binaries) == 0 {
		warnings = append(warnings, "The whitelist is empty: every process will be denied, root included.")
	}
	if draft.Mode == workbench.GuardOff {
		warnings = append(warnings, "The guard stops: "+p.v.t.Spec.Path+" is no longer protected by this window.")
	}
	confirmWarnings(p.v.page.shell, "Apply guard policy", warnings, func() { p.applyDraft(draft) })
}

func (p *policyTab) applyDraft(draft workbench.Policy) {
	p.v.page.shell.busy("Applying guard policy…", func() error {
		return p.v.t.Apply(draft)
	}, func() {
		p.refresh()
		p.v.refresh(true)
	})
}

// promote adds the executable a process ran to a mode's list, after showing what will be trusted.
func (p *policyTab) promote(exe workbench.Exe, comm string, mode workbench.GuardMode) {
	p.v.page.shell.verifyExe(exe, func(entry ebpf.BinaryEntry) {
		enf := p.v.t.Policy()
		next := workbench.Policy{Mode: mode, Binaries: p.ed.promoteBase(mode, enf.Mode, enf.Binaries, entry)}
		p.v.page.shell.confirmExecutable(verbFor(mode), entry, exe, comm, p.promoteNote(next), func() {
			p.ed.bins.set(next.Binaries, nil)
			p.ed.setMode(next.Mode)
			p.applyDraft(next)
		})
	})
}

func verbFor(mode workbench.GuardMode) string {
	if mode == workbench.GuardBlacklist {
		return "Blacklist"
	}
	return "Whitelist"
}

func (p *policyTab) promoteNote(next workbench.Policy) string {
	effect := "It may touch " + p.v.t.Spec.Path + "."
	if next.Mode == workbench.GuardBlacklist {
		effect = "It is denied on " + p.v.t.Spec.Path + "."
	}
	enf := p.v.t.Policy()
	var how string
	switch {
	case enf.Mode == next.Mode:
		how = "The running guard is updated in place."
	case enf.Mode == workbench.GuardOff && next.Mode == workbench.GuardWhitelist:
		how = fmt.Sprintf("This starts a whitelist guard with %d executable(s): every other process is denied.", len(next.Binaries))
	case enf.Mode == workbench.GuardOff:
		how = "This starts a blacklist guard."
	default:
		how = fmt.Sprintf("This switches the guard from %s to %s.", enf.Mode, next.Mode)
		if p.v.t.ApplyLeavesGap(next) {
			how += " The path is unprotected for a moment while the whitelist is replaced."
		}
	}
	return effect + " " + how
}
