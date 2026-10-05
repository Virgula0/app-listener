package codeedit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

func withColor(t *testing.T) {
	t.Helper()
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.ANSI256)
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })
}

// editSession drives ta through cursor moves, typing, wrapping and scrolling, calling check after
// every step.
func editSession(t *testing.T, content string, width, height int, check func(step string, ta *textarea.Model)) {
	t.Helper()
	ta := New()
	ta.SetWidth(width)
	ta.SetHeight(height)
	ta.Focus()
	SetText(&ta, content)
	key := func(k string) tea.KeyMsg {
		switch k {
		case "down":
			return tea.KeyMsg{Type: tea.KeyDown}
		case "up":
			return tea.KeyMsg{Type: tea.KeyUp}
		case "end":
			return tea.KeyMsg{Type: tea.KeyEnd}
		case "right":
			return tea.KeyMsg{Type: tea.KeyRight}
		}
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
	}
	check("open", &ta)
	steps := []string{"right", "right", "end", "down", "down", "end", "x", "pgdown", "down", "end", "pgdown",
		"pgdown", "end", "up", "pgup", "y", "pgup"}
	for i, s := range steps {
		msg := key(s)
		if !Edit(&ta, msg) {
			ta, _ = ta.Update(msg)
		}
		check(fmt.Sprintf("step %d (%s)", i, s), &ta)
	}
	ta.Blur()
	check("blurred", &ta)
}

func sampleGo(lines int) string {
	var b strings.Builder
	b.WriteString("package main\n\nimport \"fmt\"\n\n// main prints a very long line that has to soft-wrap across rows\n")
	for i := range lines {
		fmt.Fprintf(&b, "func f%d() { fmt.Println(\"value\", %d) } // trailing comment that wraps around\n", i, i)
	}
	b.WriteString("    \n世界 wide runes 世界世界世界世界世界世界世界世界世界世界")
	return b.String()
}

// textarea.View renders the gutter prompt as one styled run where View styles the bar and the
// number apart, so the two agree on text and layout, not bytes.
func TestHighlighterUnstyledMatchesTextareaView(t *testing.T) {
	withColor(t)
	for _, tc := range []struct {
		lines, width, height int
	}{{2, 40, 6}, {12, 30, 5}, {120, 50, 10}, {1100, 80, 24}} {
		hl := NewHighlighter("notes.txt", "")
		editSession(t, sampleGo(tc.lines), tc.width, tc.height, func(step string, ta *textarea.Model) {
			got, want := hl.View(ta), ta.View()
			if ansi.Strip(got) != ansi.Strip(want) {
				t.Fatalf("%dx%d/%d lines, %s: render differs from textarea.View\ngot:\n%s\nwant:\n%s",
					tc.width, tc.height, tc.lines, step, ansi.Strip(got), ansi.Strip(want))
			}
			if step != "blurred" && !strings.Contains(got, "\x1b[7m") {
				t.Fatalf("%s: no reverse-video cursor in\n%q", step, got)
			}
		})
	}
}

func TestGutterKeepsTextColumnAligned(t *testing.T) {
	var b strings.Builder
	for i := 1; i <= 1000; i++ {
		fmt.Fprintf(&b, "L%d\n", i)
	}
	ta := New()
	ta.SetWidth(40)
	ta.SetHeight(1200)
	SetText(&ta, b.String())
	view := ansi.Strip(NewHighlighter("x.txt", "").View(&ta))
	col := -1
	for _, row := range strings.Split(view, "\n") {
		i := strings.Index(row, "L")
		if i < 0 {
			continue
		}
		if col < 0 {
			col = i
		}
		if i != col {
			t.Fatalf("text starts at column %d, not %d: %q", i, col, row)
		}
		if w := ansi.StringWidth(row); w != 40 {
			t.Fatalf("row is %d columns wide, want the 40 it was sized to: %q", w, row)
		}
	}
	if !strings.Contains(view, "   9 L9 ") || !strings.Contains(view, " 1000 L1000 ") {
		t.Fatalf("numbers not right-aligned in a 4-digit column:\n%s", view[:200])
	}
}

func TestHighlighterOnlyChangesColors(t *testing.T) {
	withColor(t)
	hl := NewHighlighter("main.go", "")
	editSession(t, sampleGo(40), 50, 8, func(step string, ta *textarea.Model) {
		got, want := hl.View(ta), ta.View()
		if ansi.Strip(got) != ansi.Strip(want) {
			t.Fatalf("%s: highlighted text differs from textarea.View\ngot:\n%s\nwant:\n%s",
				step, ansi.Strip(got), ansi.Strip(want))
		}
	})
}

func TestHighlighterColorsTokens(t *testing.T) {
	withColor(t)
	ta := New()
	ta.SetWidth(60)
	ta.SetHeight(5)
	ta.Focus()
	SetText(&ta, "package main\n\nfunc main() {}\n")
	got := NewHighlighter("main.go", "").View(&ta)
	// Line 3: off the cursor line, so the keyword's color is not layered over its background.
	if !strings.Contains(got, tokenStyles[clsKeyword].Render("func")) {
		t.Errorf("keyword not colored:\n%q", got)
	}
	if got == ta.View() {
		t.Error("highlighted view equals the plain one")
	}
}

func TestNewHighlighterLanguage(t *testing.T) {
	for name, want := range map[string]string{
		"/etc/app-listener/daemon.conf": "daemon.conf",
		"x.py":                          "Python",
		"x.c":                           "C",
		"x.go":                          "Go",
		"settings.json":                 "JSON",
		"a.yaml":                        "YAML",
		"a.yml":                         "YAML",
		"changes.diff":                  "Diff",
		"notes.txt":                     "",
		"":                              "",
		"id_ed25519":                    "",
	} {
		if got := NewHighlighter(name, "").Language(); got != want {
			t.Errorf("%s: language %q, want %q", name, got, want)
		}
	}
	for shebang, want := range map[string]string{
		"#!/usr/bin/env python3\nimport os\n": "Python",
		"#!/bin/bash\necho hi\n":              "Bash",
		"#!/usr/bin/env -S bash -e\n":         "Bash",
		"plain text\n":                        "",
	} {
		if got := NewHighlighter("script", shebang).Language(); got != want {
			t.Errorf("%q: language %q, want %q", shebang, got, want)
		}
	}
}

func TestDaemonConfTokens(t *testing.T) {
	conf := "# comment READ\n[watch \"/home/a/.ssh\"]\n/usr/bin/ssh READ,WRITE\n/opt/APP/bin\n" +
		"need_encryption: true\nwatch: /home/a/x\n[libraries \"Steam\"]\nallow_lib /a/b.so\n"
	h := NewHighlighter("daemon.conf", conf)
	h.tokenize(conf)
	lines := strings.Split(conf, "\n")
	classAt := func(line int, sub string) uint8 {
		t.Helper()
		i := strings.Index(lines[line], sub)
		if i < 0 {
			t.Fatalf("%q not in line %d", sub, line)
		}
		return h.classes[line][len([]rune(lines[line][:i]))]
	}
	for _, c := range []struct {
		line int
		sub  string
		want uint8
	}{
		{0, "READ", clsComment},
		{1, "watch", clsKeyword},
		{1, "/home", clsString},
		{2, "/usr", clsNone},
		{2, "READ", clsConstant},
		{2, "WRITE", clsConstant},
		{3, "APP", clsNone},
		{4, "need_encryption", clsTag},
		{4, "true", clsConstant},
		{5, "watch", clsTag},
		{6, "libraries", clsKeyword},
		{7, "allow_lib", clsTag},
	} {
		if got := classAt(c.line, c.sub); got != c.want {
			t.Errorf("line %d %q: class %d, want %d", c.line, c.sub, got, c.want)
		}
	}
}

func TestTokenizeKeepsLineCount(t *testing.T) {
	h := &Highlighter{lexer: chroma.Coalesce(lexers.Get("go"))}
	for _, s := range []string{"package main", "package main\n", "a\n\nb", ""} {
		h.classes = nil
		h.tokenize(s)
		if got, want := len(h.classes), strings.Count(s, "\n")+1; got != want {
			t.Errorf("%q: %d class lines, want %d", s, got, want)
		}
	}
}

// Real editors call only Highlighter.View, never textarea.View: the view must still follow the
// cursor exactly like a plain textarea's, which needs the viewport content View refreshes on edits.
func TestHighlighterFollowsCursorAlone(t *testing.T) {
	hl := NewHighlighter("main.go", "")
	newTA := func() textarea.Model {
		ta := New()
		ta.SetWidth(60)
		ta.SetHeight(6)
		ta.Focus()
		SetText(&ta, sampleGo(80))
		return ta
	}
	plain, lit := newTA(), newTA()
	step := func(i int, msg tea.KeyMsg) {
		for _, ta := range []*textarea.Model{&plain, &lit} {
			if !Edit(ta, msg) {
				*ta, _ = ta.Update(msg)
			}
		}
		got := ansi.Strip(hl.View(&lit))
		setGutter(&plain, hl.digits, hl.prompt) // same text and width: lit's line numbers
		if want := ansi.Strip(plain.View()); got != want {
			t.Fatalf("step %d: highlighted view scrolled differently\ngot:\n%s\nwant:\n%s", i, got, want)
		}
	}
	for i := range 300 {
		switch {
		case i%37 == 0:
			step(i, tea.KeyMsg{Type: tea.KeyEnter})
		case i < 200:
			step(i, tea.KeyMsg{Type: tea.KeyDown})
		default:
			step(i, tea.KeyMsg{Type: tea.KeyPgUp})
		}
	}
	step(300, tea.KeyMsg{Type: tea.KeyCtrlEnd})
	if got := ansi.Strip(hl.View(&lit)); !strings.Contains(got, "wide runes") {
		t.Fatalf("cursor on the last line but it is not on screen:\n%s", got)
	}
}
