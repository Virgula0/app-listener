package codeedit

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func findEditor(t *testing.T, content string) *textarea.Model {
	t.Helper()
	ta := New()
	ta.SetWidth(60)
	ta.SetHeight(12)
	ta.Focus()
	SetText(&ta, content)
	return &ta
}

func keys(t *testing.T, f *Finder, ta *textarea.Model, ks ...string) {
	t.Helper()
	for _, k := range ks {
		var msg tea.KeyMsg
		switch k {
		case "ctrl+f":
			msg = tea.KeyMsg{Type: tea.KeyCtrlF}
		case "ctrl+r":
			msg = tea.KeyMsg{Type: tea.KeyCtrlR}
		case "ctrl+s":
			msg = tea.KeyMsg{Type: tea.KeyCtrlS}
		case "tab":
			msg = tea.KeyMsg{Type: tea.KeyTab}
		case "enter":
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		case "esc":
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		default:
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		}
		f.Update(ta, msg)
	}
}

func TestFindAllSmartCase(t *testing.T) {
	text := "Foo foo FOO\nfoofoo\nbar"
	if got := findAll(text, "foo"); len(got) != 5 {
		t.Fatalf("lower-case query should fold case: %v", got)
	}
	if got := findAll(text, "Foo"); len(got) != 1 || got[0] != (span{0, 0, 3}) {
		t.Fatalf("upper-case query should be exact: %v", got)
	}
	if got := findAll(text, "oo"); len(got) != 5 {
		t.Fatalf("matches must not overlap: %v", got)
	}
	if got := findAll("世界 x 世界", "世界"); len(got) != 2 || got[1] != (span{0, 5, 7}) {
		t.Fatalf("offsets are runes: %v", got)
	}
}

func TestFinderSearchMovesCursorAndWraps(t *testing.T) {
	ta := findEditor(t, "alpha\nbeta alpha\ngamma\nalpha")
	f := NewFinder(true)
	keys(t, &f, ta, "ctrl+f", "a", "l", "p")
	if ta.Line() != 0 || column(ta) != 0 || f.status != "1 of 3" {
		t.Fatalf("incremental search: line %d col %d status %q", ta.Line(), column(ta), f.status)
	}
	keys(t, &f, ta, "enter")
	if ta.Line() != 1 || column(ta) != 5 || f.status != "2 of 3" {
		t.Fatalf("next: line %d col %d status %q", ta.Line(), column(ta), f.status)
	}
	keys(t, &f, ta, "enter", "enter")
	if ta.Line() != 0 || f.status != "1 of 3" {
		t.Fatalf("wrap: line %d status %q", ta.Line(), f.status)
	}
	keys(t, &f, ta, "esc")
	if f.Active() || f.matches != nil {
		t.Fatal("esc must close the box and drop the marks")
	}
	if ta.Value() != "alpha\nbeta alpha\ngamma\nalpha" {
		t.Fatal("find changed the text")
	}
}

func TestFinderSearchAllMarksEveryMatch(t *testing.T) {
	ta := findEditor(t, "x y x\nx")
	f := NewFinder(false)
	keys(t, &f, ta, "ctrl+f", "x")
	marks, _, ok := f.marks()
	if !ok || len(marks[0])+len(marks[1]) != 1 {
		t.Fatalf("Search marks only the current match: %v", marks)
	}
	keys(t, &f, ta, "tab", "tab", "enter") // query → Search → Search all
	marks, _, _ = f.marks()
	if len(marks[0]) != 2 || len(marks[1]) != 1 {
		t.Fatalf("Search all marks every match: %v", marks)
	}
}

func TestFinderReplace(t *testing.T) {
	ta := findEditor(t, "cat cat\ndog cat")
	f := NewFinder(true)
	keys(t, &f, ta, "ctrl+r", "c", "a", "t", "tab", "c", "o", "w", "enter")
	if got := ta.Value(); got != "cow cat\ndog cat" {
		t.Fatalf("Replace (enter in the replacement field): %q", got)
	}
	if f.status != "1 of 2" || ta.Line() != 0 || column(ta) != 4 {
		t.Fatalf("Replace moves to the next match: status %q line %d col %d", f.status, ta.Line(), column(ta))
	}
	keys(t, &f, ta, "tab", "tab", "enter") // → Replace all
	if got := ta.Value(); got != "cow cow\ndog cow" {
		t.Fatalf("Replace all: %q", got)
	}
	if f.status != "replaced 2" {
		t.Fatalf("status %q", f.status)
	}
}

func TestFinderReplaceWithSupersetTerminates(t *testing.T) {
	ta := findEditor(t, "a a")
	f := NewFinder(true)
	keys(t, &f, ta, "ctrl+r", "a", "tab", "a", "a", "enter", "enter")
	if got := ta.Value(); got != "aa aa" {
		t.Fatalf("each Replace must move past its own replacement: %q", got)
	}
}

func TestFinderReadOnlyIgnoresReplace(t *testing.T) {
	ta := findEditor(t, "text")
	f := NewFinder(false)
	if f.Update(ta, tea.KeyMsg{Type: tea.KeyCtrlR}) || f.Active() {
		t.Fatal("a read-only viewer must not open replace")
	}
	keys(t, &f, ta, "ctrl+f", "ctrl+r")
	if f.replace {
		t.Fatal("ctrl+r must not switch a read-only finder to replace")
	}
}

func TestFinderLeavesSaveToHost(t *testing.T) {
	ta := findEditor(t, "text")
	f := NewFinder(true)
	keys(t, &f, ta, "ctrl+f")
	if f.Update(ta, tea.KeyMsg{Type: tea.KeyCtrlS}) {
		t.Fatal("ctrl+s must reach the host while the box is open")
	}
	if !f.Update(ta, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}) {
		t.Fatal("typing must go to the box, not the host (q quits viewers)")
	}
}

func TestFinderViewOverlaysBoxAndMarks(t *testing.T) {
	withColor(t)
	content := strings.Repeat("package main // needle here\n", 20)
	ta := findEditor(t, content)
	hl := NewHighlighter("main.go", content)
	f := NewFinder(true)
	plain := hl.View(ta)
	keys(t, &f, ta, "ctrl+r", "n", "e", "e", "d", "l", "e")
	got := f.View(hl, ta)
	rows, base := strings.Split(got, "\n"), strings.Split(plain, "\n")
	if len(rows) != len(base) {
		t.Fatalf("overlay changed the height: %d vs %d", len(rows), len(base))
	}
	for i := range rows {
		if w, bw := ansi.StringWidth(rows[i]), ansi.StringWidth(base[i]); w != bw {
			t.Fatalf("row %d width %d, want %d:\n%s", i, w, bw, got)
		}
	}
	if s := ansi.Strip(got); !strings.Contains(s, "Find & replace") || !strings.Contains(s, "Replace all") {
		t.Fatalf("box missing:\n%s", s)
	}
	// The box sits on the top rows unless the cursor is under it; the match is on row 0, so the
	// box moves to the bottom and the underlined, reversed match stays visible.
	if !strings.Contains(rows[0], "\x1b[1;4;7") && !strings.Contains(rows[0], ";4;7") {
		t.Fatalf("current match not underlined+reversed:\n%q", rows[0])
	}
	f.Close()
	if f.View(hl, ta) != NewHighlighter("main.go", content).View(ta) {
		t.Fatal("a closed finder must render exactly as Highlighter.View")
	}
}
