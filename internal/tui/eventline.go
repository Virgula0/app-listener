package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// formatGuardEventLine renders one guard event line, shared by the guard
// and daemon views.
func formatGuardEventLine(ev *guard.GuardEvent) string {
	safeEvent := *ev
	safeEvent.Dest = sanitizeTerminalText(safeEvent.Dest)
	ev = &safeEvent
	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatGuardType(ev.Type)
	decision := formatDecision(ev.Blocked)

	pathPart := sanitizeTerminalText(ev.Path)
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
		commStyle.Render(sanitizeTerminalText(ev.Comm)),
		ev.PID,
		pathPart,
		extra,
	)
}

func sanitizeTerminalText(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return '?'
		}
		return r
	}, value)
}

func sanitizeTerminalTexts(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = sanitizeTerminalText(value)
	}
	return result
}
