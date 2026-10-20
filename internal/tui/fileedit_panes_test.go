package tui

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

// openAlpha opens <root>/alpha.txt in the editor of a 120x30 model.
func openAlpha(t *testing.T) (*fileEditModel, string) {
	t.Helper()
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	for _, k := range []string{"j", "l", "j", "j", "e"} {
		updateKey(m, k)
	}
	if m.mode != modeEdit || m.editPath != filepath.Join(root, "alpha.txt") {
		t.Fatalf("mode %v editing %q", m.mode, m.editPath)
	}
	return m, root
}

// screen is the frame as the terminal shows it: bubbletea drops the top rows of a frame taller
// than the terminal.
func screen(m *fileEditModel) []string {
	rows := strings.Split(ansi.Strip(m.View()), "\n")
	return rows[max(0, len(rows)-m.height):]
}

// screenCell is where needle first appears on the screen.
func screenCell(t *testing.T, m *fileEditModel, needle string) (x, y int) {
	t.Helper()
	for y, row := range screen(m) {
		if before, _, ok := strings.Cut(row, needle); ok {
			return ansi.StringWidth(before), y
		}
	}
	t.Fatalf("%q not on screen:\n%s", needle, ansi.Strip(m.View()))
	return 0, 0
}

func click(m *fileEditModel, x, y int) {
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
}

func TestFileEditPaneFocusKeepsEditorOpen(t *testing.T) {
	m, _ := openAlpha(t)
	editW := m.rightW
	updateKey(m, "alt+left")
	if m.mode != modeNav || !m.editorOpen || m.rightW != editW {
		t.Fatalf("alt+left: mode %v open %v rightW %d (was %d)", m.mode, m.editorOpen, m.rightW, editW)
	}
	if !strings.Contains(m.View(), "hello alpha") {
		t.Fatal("the editor left the screen when the tree took the focus")
	}
	updateKey(m, "k")
	if m.selected().name != "tool.sh" {
		t.Fatalf("tree keys after alt+left moved to %q", m.selected().name)
	}
	updateKey(m, "alt+right")
	if m.mode != modeEdit || m.editPath == "" || !strings.HasSuffix(m.editPath, "alpha.txt") {
		t.Fatalf("alt+right: mode %v editing %q", m.mode, m.editPath)
	}
	updateKey(m, "alt+left")
	updateKey(m, "esc")
	if m.editorOpen || m.mode != modeNav {
		t.Fatal("esc in the tree did not close the open editor")
	}
}

func TestFileEditHideTree(t *testing.T) {
	m, _ := openAlpha(t)
	usable := m.width - appStyle.GetHorizontalMargins()
	updateKey(m, "ctrl+b")
	if !m.treeHidden || m.leftW != 0 || m.rightW != usable || m.mode != modeEdit {
		t.Fatalf("hidden: leftW %d rightW %d mode %v", m.leftW, m.rightW, m.mode)
	}
	if got := maxLineWidth(m.body()); got != usable {
		t.Fatalf("hidden: body width %d, want %d", got, usable)
	}
	if strings.Contains(ansi.Strip(m.body()), "tool.sh") {
		t.Fatal("hidden tree still drawn")
	}
	updateKey(m, "alt+left")
	if m.treeHidden || m.mode != modeNav {
		t.Fatal("alt+left did not bring the tree back")
	}
	updateKey(m, "esc")
	updateKey(m, "ctrl+b")
	if m.treeHidden {
		t.Fatal("tree hidden with no editor open")
	}
}

func TestFileEditSwitchFileAsksBeforeDiscarding(t *testing.T) {
	m, root := openAlpha(t)
	m.editor.SetValue("changed")
	m.dirty = true
	updateKey(m, "alt+left")
	updateKey(m, "j") // zebra.txt
	updateKey(m, "e")
	if m.mode != modeConfirm {
		t.Fatalf("opening another file over unsaved changes: mode %v", m.mode)
	}
	updateKey(m, "n")
	if m.editor.Value() != "changed" {
		t.Fatal("declining discarded the buffer")
	}
	updateKey(m, "e")
	updateKey(m, "y")
	if m.editPath != filepath.Join(root, "zebra.txt") || m.editor.Value() != "zebra" || m.mode != modeEdit {
		t.Fatalf("after confirming: editing %q = %q", m.editPath, m.editor.Value())
	}
}

func TestFileEditMouse(t *testing.T) {
	m, root := openAlpha(t)
	x, y := screenCell(t, m, "hello alpha")
	x += len("hello ")
	m.Update(tea.MouseMsg{X: x + 2, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	m.Update(tea.MouseMsg{X: x + 5, Y: y, Action: tea.MouseActionMotion, Button: tea.MouseButtonLeft})
	m.Update(tea.MouseMsg{X: x + 5, Y: y, Action: tea.MouseActionRelease})
	if got := m.sel.Text(&m.editor); got != "pha" {
		t.Fatalf("drag in the editor selected %q", got)
	}

	zx, zy := screenCell(t, m, "zebra.txt")
	click(m, zx, zy)
	if m.mode != modeNav || m.selected().name != "zebra.txt" || !m.editorOpen {
		t.Fatalf("tree click: mode %v row %q open %v", m.mode, m.selected().name, m.editorOpen)
	}
	click(m, zx, zy)
	if m.editPath != filepath.Join(root, "zebra.txt") || m.mode != modeEdit {
		t.Fatalf("second click opened %q (mode %v)", m.editPath, m.mode)
	}
	updateKey(m, "alt+left")
	ex, ey := screenCell(t, m, "1 zebra")
	click(m, ex+2, ey)
	if m.mode != modeEdit {
		t.Fatal("a click in the editor did not focus it")
	}
}

func TestFileEditFrameFitsTerminal(t *testing.T) {
	m, _ := openAlpha(t)
	check := func(state string) {
		t.Helper()
		if got := len(strings.Split(m.View(), "\n")); got != m.height {
			t.Fatalf("%s %dx%d: frame is %d rows", state, m.width, m.height, got)
		}
	}
	for _, sz := range [][2]int{{120, 30}, {100, 16}, {200, 50}, {80, 24}, {60, 12}} {
		m.Update(tea.WindowSizeMsg{Width: sz[0], Height: sz[1]})
		check("editor")
		updateKey(m, "ctrl+b")
		check("tree hidden")
		updateKey(m, "alt+left")
		check("tree focused")
		updateKey(m, "esc")
		check("details")
		updateKey(m, "g")
		check("directory details")
		for i, r := range m.rows {
			if r.name == "alpha.txt" {
				m.cur = i
			}
		}
		updateKey(m, "e")
	}
}
