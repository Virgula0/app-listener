package codeedit

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

func TestWordJumps(t *testing.T) {
	text := [][]rune{[]rune("  /home/alice/.ssh = yes"), []rune(""), []rune("next")}
	right := [][2]int{}
	row, col := 0, 0
	for range 9 {
		row, col = wordRight(text, row, col)
		right = append(right, [2]int{row, col})
	}
	want := [][2]int{{0, 3}, {0, 7}, {0, 8}, {0, 13}, {0, 15}, {0, 18}, {0, 20}, {0, 24}, {1, 0}}
	for i := range want {
		if right[i] != want[i] {
			t.Fatalf("ctrl+right stops = %v, want %v", right, want)
		}
	}
	if r, c := wordRight(text, 1, 0); r != 2 || c != 0 {
		t.Errorf("ctrl+right from an empty line = (%d,%d), want (2,0)", r, c)
	}
	if r, c := wordLeft(text, 2, 0); r != 1 || c != 0 {
		t.Errorf("ctrl+left from a line start = (%d,%d), want the previous line's end (1,0)", r, c)
	}
	if r, c := wordLeft(text, 0, 13); r != 0 || c != 8 {
		t.Errorf("ctrl+left from after alice = (%d,%d), want (0,8)", r, c)
	}
	if r, c := wordLeft(text, 0, 0); r != 0 || c != 0 {
		t.Errorf("ctrl+left at the start = (%d,%d)", r, c)
	}
}

func TestSetTextStartsAtFirstLine(t *testing.T) {
	ta := New()
	ta.SetHeight(5)
	ta.Focus()
	SetText(&ta, strings.Repeat("line\n", 100)+"last")
	if ta.Line() != 0 || column(&ta) != 0 {
		t.Fatalf("cursor at (%d,%d), want (0,0)", ta.Line(), column(&ta))
	}
	ta, _ = ta.Update(nil)
	if !strings.Contains(ta.View(), "line") || strings.Contains(ta.View(), "last") {
		t.Fatalf("view does not show the first lines:\n%s", ta.View())
	}
}

func TestEditKeys(t *testing.T) {
	ta := New()
	ta.SetHeight(3)
	ta.Focus()
	SetText(&ta, "a\nb\nc\nd\ne\nf\ng")
	if !Edit(&ta, tea.KeyMsg{Type: tea.KeyTab}) || ta.Value()[:5] != "    a" {
		t.Fatalf("tab: %q", ta.Value())
	}
	if !Edit(&ta, tea.KeyMsg{Type: tea.KeyShiftTab}) || !strings.HasPrefix(ta.Value(), "a\n") {
		t.Fatalf("shift+tab: %q", ta.Value())
	}
	Navigate(&ta, tea.KeyMsg{Type: tea.KeyPgDown})
	if ta.Line() != 3 {
		t.Fatalf("pgdown landed on line %d, want 3", ta.Line())
	}
	Navigate(&ta, tea.KeyMsg{Type: tea.KeyPgUp})
	if ta.Line() != 0 {
		t.Fatalf("pgup landed on line %d, want 0", ta.Line())
	}
	if Navigate(&ta, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}) {
		t.Fatal("a plain key was consumed")
	}
}

func TestCursorDoesNotBlink(t *testing.T) {
	ta := New()
	ta.Focus()
	if ta.Cursor.Mode() != cursor.CursorStatic {
		t.Fatalf("cursor mode = %v, want static", ta.Cursor.Mode())
	}
}

func TestEnterKeepsIndentation(t *testing.T) {
	ta := New()
	ta.Focus()
	SetText(&ta, "[watch /a]\n    /usr/bin/cat")
	moveTo(&ta, 1, len("    /usr/bin/cat"))
	Edit(&ta, tea.KeyMsg{Type: tea.KeyEnter})
	ta.InsertString("/usr/bin/head")
	if want := "[watch /a]\n    /usr/bin/cat\n    /usr/bin/head"; ta.Value() != want {
		t.Fatalf("value = %q, want %q", ta.Value(), want)
	}

	// Inside the indentation only what's left of the cursor carries over; an unindented line none.
	SetText(&ta, "    x")
	moveTo(&ta, 0, 2)
	Edit(&ta, tea.KeyMsg{Type: tea.KeyEnter})
	if want := "  \n    x"; ta.Value() != want {
		t.Fatalf("split inside the indent = %q, want %q", ta.Value(), want)
	}
	SetText(&ta, "x")
	moveTo(&ta, 0, 1)
	Edit(&ta, tea.KeyMsg{Type: tea.KeyEnter})
	if ta.Value() != "x\n" {
		t.Fatalf("unindented = %q", ta.Value())
	}
}

func TestEnterPastNinetyNineLines(t *testing.T) {
	ta := New()
	ta.Focus()
	SetText(&ta, strings.Repeat("l\n", 150)+"end")
	moveTo(&ta, 150, 3)
	Edit(&ta, tea.KeyMsg{Type: tea.KeyEnter})
	if ta.LineCount() != 152 {
		t.Fatalf("lines = %d, want 152: Enter was refused", ta.LineCount())
	}
}

func TestSoftTabStops(t *testing.T) {
	line := []rune("        x = 1")
	cases := []struct {
		col, back, fwd int
	}{
		{8, 4, 0}, // before x: back one stop; Delete would eat the x
		{4, 4, 4},
		{6, 2, 2}, // off-stop: back to 4, forward to 8
		{0, 0, 4},
		{10, 0, 0}, // past the indentation: ordinary keys
	}
	for _, c := range cases {
		if got := softTabBack(line, c.col); got != c.back {
			t.Errorf("softTabBack(col %d) = %d, want %d", c.col, got, c.back)
		}
		if got := softTabForward(line, c.col); got != c.fwd {
			t.Errorf("softTabForward(col %d) = %d, want %d", c.col, got, c.fwd)
		}
	}
}

func TestBackspaceAndDeleteRemoveAnIndentStep(t *testing.T) {
	ta := New()
	ta.Focus()
	SetText(&ta, "a\n        b")
	moveTo(&ta, 1, 8)
	if !Edit(&ta, tea.KeyMsg{Type: tea.KeyBackspace}) || ta.Value() != "a\n    b" || column(&ta) != 4 {
		t.Fatalf("backspace: %q col %d", ta.Value(), column(&ta))
	}
	moveTo(&ta, 1, 0)
	if !Edit(&ta, tea.KeyMsg{Type: tea.KeyDelete}) || ta.Value() != "a\nb" {
		t.Fatalf("delete: %q", ta.Value())
	}
	// Outside the indentation the keys stay ordinary (left to the textarea).
	moveTo(&ta, 1, 1)
	if Edit(&ta, tea.KeyMsg{Type: tea.KeyBackspace}) {
		t.Fatal("a plain backspace was consumed")
	}
}
