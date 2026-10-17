package codeedit

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func selKey(t *testing.T, s *Selection, ta *textarea.Model, editable bool, msgs ...tea.KeyMsg) tea.Cmd {
	t.Helper()
	var last tea.Cmd
	for _, m := range msgs {
		ok, cmd := s.Key(ta, m, editable)
		if !ok && !Edit(ta, m) {
			*ta, _ = ta.Update(m)
		}
		last = cmd
	}
	return last
}

func repeat(t tea.KeyType, n int) []tea.KeyMsg {
	out := make([]tea.KeyMsg, n)
	for i := range out {
		out[i] = tea.KeyMsg{Type: t}
	}
	return out
}

func TestSelectionShiftArrowsCopy(t *testing.T) {
	var out bytes.Buffer
	clipboardOut = &out
	t.Cleanup(func() { clipboardOut = nil; setClipboard("") })
	ta := findEditor(t, "hello world\nsecond line\nthird")
	var s Selection
	selKey(t, &s, ta, true, repeat(tea.KeyShiftRight, 5)...)
	if got := s.Text(ta); got != "hello" {
		t.Fatalf("shift+right x5 = %q", got)
	}
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyCtrlShiftDown})
	if got := s.Text(ta); got != "hello world\nsecon" {
		t.Fatalf("ctrl+shift+down = %q", got)
	}
	cmd := selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil || clipboardText() != "hello world\nsecon" || s.TakeNote() != "copied 17 characters" {
		t.Fatalf("copy: cmd %v clip %q", cmd, clipboardText())
	}
	cmd()
	if want := base64.StdEncoding.EncodeToString([]byte("hello world\nsecon")); !strings.Contains(out.String(), "]52;c;"+want) {
		t.Fatalf("no OSC 52 sequence: %q", out.String())
	}
	if s.Text(ta) == "" {
		t.Fatal("copy dropped the selection")
	}
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyEsc})
	if s.Text(ta) != "" {
		t.Fatal("esc kept the selection")
	}
}

func TestSelectionEditing(t *testing.T) {
	t.Cleanup(func() { setClipboard("") })
	ta := findEditor(t, "hello world")
	var s Selection
	selKey(t, &s, ta, true, repeat(tea.KeyShiftRight, 5)...)
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	if ta.Value() != "X world" {
		t.Fatalf("typing over a selection = %q", ta.Value())
	}
	selKey(t, &s, ta, true, repeat(tea.KeyShiftRight, 6)...)
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyCtrlX})
	if ta.Value() != "X" || clipboardText() != " world" {
		t.Fatalf("cut: value %q clip %q", ta.Value(), clipboardText())
	}
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyCtrlV}, tea.KeyMsg{Type: tea.KeyCtrlV})
	if ta.Value() != "X world world" {
		t.Fatalf("paste = %q", ta.Value())
	}
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyCtrlA}, tea.KeyMsg{Type: tea.KeyBackspace})
	if ta.Value() != "" {
		t.Fatalf("ctrl+a, backspace = %q", ta.Value())
	}
}

func TestSelectionReadOnly(t *testing.T) {
	ta := findEditor(t, "hello world")
	var s Selection
	selKey(t, &s, ta, false, repeat(tea.KeyShiftRight, 5)...)
	if ok, _ := s.Key(ta, tea.KeyMsg{Type: tea.KeyCtrlX}, false); ok || ta.Value() != "hello world" {
		t.Fatalf("read-only cut: consumed %v, value %q", ok, ta.Value())
	}
	if s.Text(ta) != "" {
		t.Fatal("another key kept the selection")
	}
}

func TestSelectionLargeDelete(t *testing.T) {
	content := strings.Repeat("0123456789\n", 1000) + "tail"
	ta := findEditor(t, content)
	var s Selection
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyRight},
		tea.KeyMsg{Type: tea.KeyCtrlShiftEnd}, tea.KeyMsg{Type: tea.KeyShiftLeft})
	selKey(t, &s, ta, true, tea.KeyMsg{Type: tea.KeyDelete})
	if ta.Value() != "0l" || ta.Line() != 0 || column(ta) != 1 {
		t.Fatalf("value %q cursor (%d,%d)", ta.Value(), ta.Line(), column(ta))
	}
}

// cellOf is the screen cell of the first rune of needle in the rendered view.
func cellOf(t *testing.T, view, needle string) (x, y int) {
	t.Helper()
	for y, row := range strings.Split(ansi.Strip(view), "\n") {
		if i := strings.Index(row, needle); i >= 0 {
			return ansi.StringWidth(row[:i]), y
		}
	}
	t.Fatalf("%q not on screen:\n%s", needle, ansi.Strip(view))
	return 0, 0
}

func TestSelectionMouseDrag(t *testing.T) {
	withColor(t)
	ta := findEditor(t, "first line\nhello world\nlast")
	h := NewHighlighter("x.txt", ta.Value())
	var s Selection
	view := h.View(ta)
	x, y := cellOf(t, view, "world")
	if !s.Mouse(ta, h, tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft}, x, y) {
		t.Fatal("press not handled")
	}
	if ta.Line() != 1 || column(ta) != 6 {
		t.Fatalf("press put the cursor at (%d,%d)", ta.Line(), column(ta))
	}
	s.Mouse(ta, h, tea.MouseMsg{Action: tea.MouseActionMotion, Button: tea.MouseButtonLeft}, x+5, y)
	s.Mouse(ta, h, tea.MouseMsg{Action: tea.MouseActionRelease}, x+5, y)
	if got := s.Text(ta); got != "world" {
		t.Fatalf("drag selected %q", got)
	}
	if s.Dragging() {
		t.Fatal("still dragging after release")
	}
	sel := (&Finder{}).View(h, ta, &s)
	if !strings.Contains(sel, "48;5;24") && !strings.Contains(sel, "48;5;153") {
		t.Fatalf("selection not painted:\n%q", sel)
	}
	s.Mouse(ta, h, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft}, x, y+5)
	if ta.Line() != 2 || column(ta) != 4 {
		t.Fatalf("press below the text = (%d,%d), want the end of the last line", ta.Line(), column(ta))
	}
}
