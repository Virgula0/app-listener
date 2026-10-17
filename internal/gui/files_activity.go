//go:build gui

package gui

import (
	"strings"

	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

func (v *targetView) promoteActions() []feedAction {
	return []feedAction{
		{label: "Whitelist", icon: "shield-check", importance: widget.SuccessImportance,
			run: func(e workbench.Exe, comm string) { v.policy.promote(e, comm, workbench.GuardWhitelist) }},
		{label: "Blacklist", icon: "prohibit", importance: widget.DangerImportance,
			run: func(e workbench.Exe, comm string) { v.policy.promote(e, comm, workbench.GuardBlacklist) }},
	}
}

// newFileFeed is the merged monitor + guard event table of a target.
func newFileFeed(v *targetView) *eventFeed {
	root := v.t.Spec.Path
	return newEventFeed(v.page.shell, &feedSpec{
		cols: []column{
			{"Time", 104}, {"Verdict", 86}, {"Operation", 104}, {"Process", 130}, {"PID", 72},
			{"Path", 380}, {"Executable", 300},
		},
		cell: func(r *workbench.Row, col int) cell {
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
				if r.Dest != "" {
					return cell{text: relTo(root, r.Path) + " → " + relTo(root, r.Dest)}
				}
				return cell{text: relTo(root, r.Path)}
			default:
				return cell{text: exeText(r.Exe), importance: widget.LowImportance}
			}
		},
		source:      v.t.Rows,
		clear:       v.t.Clear,
		placeholder: "Filter by process, path, operation or executable",
		hint:        "Select an event to whitelist or blacklist the executable behind it.",
		actions:     v.promoteActions(),
	})
}

// relTo shows p relative to the target root ("./id_ed25519"); anything outside stays absolute.
func relTo(root, p string) string {
	if p == root {
		return "."
	}
	if rest, ok := strings.CutPrefix(p, root+"/"); ok {
		return "./" + rest
	}
	return p
}

// newFileProcesses lists the executables seen on a target.
func newFileProcesses(v *targetView) *processView {
	return newProcessView(v.page.shell, &procSpec{
		cols: append(append([]column(nil), procColumns...), column{"Policy", 110}),
		cell: func(pr *workbench.Process, col int) cell {
			if c, ok := procCell(pr, col); ok {
				return c
			}
			pol := v.t.Policy()
			return policyMark(pol.Mode, entryPaths(pol.Binaries), pr.Exe)
		},
		source: v.t.Processes,
		hint: "Every executable that touched this target, identified by its inode. " +
			"Whitelist or blacklist one to put it in the guard policy.",
		actions: v.promoteActions(),
	})
}

func entryPaths(bins []ebpf.BinaryEntry) []string {
	out := make([]string, len(bins))
	for i, b := range bins {
		out[i] = b.Path
	}
	return out
}
