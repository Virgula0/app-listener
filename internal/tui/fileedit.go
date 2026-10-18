package tui

import (
	"fmt"
	"time"

	"github.com/Virgula0/app-listener/internal/vaultfs"
)

// humanSize renders a byte count with a binary unit suffix.
func humanSize(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.1f KiB", float64(n)/(1<<10))
	case n < 1<<30:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	default:
		return fmt.Sprintf("%.1f GiB", float64(n)/(1<<30))
	}
}

// clipLabel truncates s to w runes, appending an ellipsis when shortened.
func clipLabel(s string, w int) string {
	if w <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w == 1 {
		return "…"
	}
	return string(r[:w-1]) + "…"
}

// dirEntry is one row of the right-pane directory listing.
type dirEntry struct {
	name  string
	isDir bool
	size  int64
	mtime time.Time
}

// listEntries stats up to limit children, skipping vanished entries.
func listEntries(v *vaultfs.Vault, n *fileNode, limit int) []dirEntry {
	out := make([]dirEntry, 0, min(limit, len(n.children)))
	for _, c := range n.children {
		if len(out) >= limit {
			break
		}
		info, err := v.Lstat(c.path)
		if err != nil {
			continue
		}
		out = append(out, dirEntry{name: c.name, isDir: c.isDir, size: info.Size(), mtime: info.ModTime()})
	}
	return out
}
