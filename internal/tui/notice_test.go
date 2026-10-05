package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestRenderBulletsVerbatimAndIndented(t *testing.T) {
	items := []string{
		"file .app_listener/id_ed25519 is readable/writable by group or others (mode -rw-r--r--) — protected content should be 0600",
		"symlink a_b -> *target*",
	}
	out := renderBullets(items, 40)
	for _, want := range []string{".app_listener/id_ed25519", "a_b", "*target*"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q lost from the rendering:\n%s", want, out)
		}
	}
	if strings.Contains(out, "\x1b[") {
		t.Errorf("markup escapes in the rendering:\n%q", out)
	}
	lines := strings.Split(out, "\n")
	bullets := 0
	for _, l := range lines {
		if lipgloss.Width(l) > 40 {
			t.Errorf("line wider than 40 cells: %q", l)
		}
		switch {
		case strings.HasPrefix(l, "• "):
			bullets++
		case !strings.HasPrefix(l, "  "):
			t.Errorf("continuation line not indented under its bullet: %q", l)
		}
	}
	if bullets != len(items) || len(lines) <= len(items) {
		t.Fatalf("expected %d wrapped bullets, got:\n%s", len(items), out)
	}
}

func TestNoticeViewFitsAndPlacesOKLast(t *testing.T) {
	m := &noticeModel{title: "Post-edit audit", intro: "Advisory only.",
		items: []string{strings.Repeat("long_path_segment/", 12), "short"}}
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 30})
	view := m.View()
	lines := strings.Split(view, "\n")
	ok, short := -1, -1
	for i, l := range lines {
		if lipgloss.Width(l) > 60 {
			t.Errorf("line %d wider than the terminal: %q", i, l)
		}
		if strings.Contains(l, " OK ") {
			ok = i
		}
		if strings.Contains(l, "• short") {
			short = i
		}
	}
	if ok < 0 || short < 0 || ok < short {
		t.Fatalf("OK (line %d) must follow the last finding (line %d):\n%s", ok, short, view)
	}
}
