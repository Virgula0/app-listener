package tui

import (
	"fmt"
	"time"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// formatGuardEventLine renders one guard event line, shared by the guard
// and daemon views.
func formatGuardEventLine(ev *guard.GuardEvent) string {
	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatGuardType(ev.Type)
	decision := formatDecision(ev.Blocked)

	pathPart := ev.Path
	extra := ""

	switch ev.Type {
	case ebpf.EventRead, ebpf.EventWrite:
		extra = fmt.Sprintf(" fd=%d", ev.FD)
	case ebpf.EventRename, ebpf.EventSymlink, ebpf.EventHardlink:
		extra = fmt.Sprintf(" \u2192 %s", ev.Dest)
	}

	return fmt.Sprintf("%s %s %s %s[%d] %s%s",
		timeStyle.Render(ts),
		t,
		decision,
		commStyle.Render(ev.Comm),
		ev.PID,
		pathPart,
		extra,
	)
}
