// Package codeedit gives the bubbles textarea the code-editor behavior the TUIs share: a steady
// cursor, content opening at its first line, Enter keeping the indentation, Tab/Shift+Tab indent,
// PgUp/PgDn page, Ctrl+Left/Right stopping at word and punctuation boundaries across lines, and
// (Highlighter.View) a fixed-width line-number gutter plus syntax highlighting.
package codeedit

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// indent is what Tab inserts: the textarea can't hold a tab character (its sanitizer expands
// tabs to four spaces), so indentation is soft.
const indent = "    "

// New returns a focused-ready textarea with line numbers and the shared key bindings.
func New() textarea.Model {
	ta := textarea.New()
	// The textarea sizes its number column by MaxHeight's digits (1 here), so each extra digit
	// shifted the text right: Highlighter.View draws the numbers in a fixed-width prompt instead.
	ta.ShowLineNumbers = false
	blank := ta.Prompt + gutterLabel(minGutterDigits, "")
	setGutter(&ta, minGutterDigits, func(int) string { return blank })
	ta.Cursor.SetMode(cursor.CursorStatic)
	// The default caps the buffer at 99 lines: Enter does nothing past it.
	ta.MaxHeight = 0
	ta.KeyMap.LineNext.SetKeys("down", "ctrl+n")
	ta.KeyMap.LinePrevious.SetKeys("up", "ctrl+p")
	ta.KeyMap.CharacterBackward.SetKeys("left")
	ta.KeyMap.CharacterForward.SetKeys("right")
	// Word jumps are handled by Navigate; keep only the emacs-style alternates here.
	ta.KeyMap.WordBackward.SetKeys("alt+b")
	ta.KeyMap.WordForward.SetKeys("alt+f")
	return ta
}

// SetText replaces the content and puts the cursor at the start of the first line: SetValue
// leaves it on the last, and the view follows the cursor, so long content opened scrolled to its
// end and looked empty.
func SetText(ta *textarea.Model, s string) {
	ta.SetValue(s)
	moveTo(ta, 0, 0)
}

// Navigate applies the movement keys the textarea lacks (PgUp/PgDn, word jumps). It reports
// whether it consumed msg; otherwise pass msg to the textarea's Update.
func Navigate(ta *textarea.Model, msg tea.KeyMsg) bool {
	switch msg.String() {
	case "pgdown":
		for range max(ta.Height(), 1) {
			ta.CursorDown()
		}
	case "pgup":
		for range max(ta.Height(), 1) {
			ta.CursorUp()
		}
	case "ctrl+right", "alt+right":
		row, col := wordRight(lines(ta), ta.Line(), column(ta))
		moveTo(ta, row, col)
	case "ctrl+left", "alt+left":
		row, col := wordLeft(lines(ta), ta.Line(), column(ta))
		moveTo(ta, row, col)
	default:
		return false
	}
	reposition(ta)
	return true
}

// Edit is Navigate plus the editing keys: Enter keeps the line's indentation, Tab indents at the
// cursor, Shift+Tab removes one indent level from the start of the line, and Backspace/Delete
// inside the indentation remove a whole indent step.
func Edit(ta *textarea.Model, msg tea.KeyMsg) bool {
	switch msg.String() {
	case "backspace", "ctrl+h":
		return repeatKey(ta, tea.KeyBackspace, softTabBack(lines(ta)[ta.Line()], column(ta)))
	case "delete":
		return repeatKey(ta, tea.KeyDelete, softTabForward(lines(ta)[ta.Line()], column(ta)))
	case "enter", "ctrl+m":
		ta.InsertString("\n" + leadingIndent(lines(ta)[ta.Line()], column(ta)))
	case "tab":
		ta.InsertString(indent)
	case "shift+tab":
		dedent(ta)
	default:
		return Navigate(ta, msg)
	}
	reposition(ta)
	return true
}

// reposition runs the textarea's update with no input: it scrolls the view to the cursor.
func reposition(ta *textarea.Model) {
	*ta, _ = ta.Update(nil)
}

func lines(ta *textarea.Model) [][]rune {
	parts := strings.Split(ta.Value(), "\n")
	out := make([][]rune, len(parts))
	for i, p := range parts {
		out[i] = []rune(p)
	}
	return out
}

// column is the cursor's rune index in its (logical) line.
func column(ta *textarea.Model) int {
	li := ta.LineInfo()
	return li.StartColumn + li.ColumnOffset
}

// moveTo puts the cursor on (row, col). Rows are reached with CursorUp/Down, which step through
// soft-wrapped sub-lines, so each loop stops as soon as a step makes no progress.
func moveTo(ta *textarea.Model, row, col int) {
	for ta.Line() < row {
		before := ta.LineInfo().RowOffset + ta.Line()<<16
		ta.CursorDown()
		if ta.LineInfo().RowOffset+ta.Line()<<16 == before {
			break
		}
	}
	for ta.Line() > row {
		before := ta.LineInfo().RowOffset + ta.Line()<<16
		ta.CursorUp()
		if ta.LineInfo().RowOffset+ta.Line()<<16 == before {
			break
		}
	}
	ta.SetCursor(col)
}

// leadingIndent is line's leading whitespace, cut at col when the cursor is inside it.
func leadingIndent(line []rune, col int) string {
	n := 0
	for n < len(line) && n < col && (line[n] == ' ' || line[n] == '\t') {
		n++
	}
	return string(line[:n])
}

// repeatKey sends key n times; n < 2 is an ordinary keypress, left to the caller (false).
func repeatKey(ta *textarea.Model, key tea.KeyType, n int) bool {
	if n < 2 {
		return false
	}
	for range n {
		*ta, _ = ta.Update(tea.KeyMsg{Type: key})
	}
	reposition(ta)
	return true
}

// softTabBack is how many spaces Backspace removes at col: back to the previous indent stop while
// only spaces precede the cursor, else 0.
func softTabBack(line []rune, col int) int {
	if col == 0 || col > len(line) || !onlySpaces(line[:col]) {
		return 0
	}
	if n := col % len(indent); n != 0 {
		return n
	}
	return len(indent)
}

// softTabForward is how many spaces Delete removes at col: up to the next indent stop while the
// cursor is in the leading spaces and only spaces reach that stop, else 0.
func softTabForward(line []rune, col int) int {
	n := len(indent) - col%len(indent)
	if col > len(line) || col+n > len(line) || !onlySpaces(line[:col+n]) {
		return 0
	}
	return n
}

func onlySpaces(rs []rune) bool {
	for _, r := range rs {
		if r != ' ' {
			return false
		}
	}
	return true
}

func dedent(ta *textarea.Model) {
	row, col := ta.Line(), column(ta)
	line := lines(ta)[row]
	n := 0
	for n < len(indent) && n < len(line) && line[n] == ' ' {
		n++
	}
	if n == 0 {
		return
	}
	ta.CursorStart()
	for range n {
		*ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyDelete})
	}
	ta.SetCursor(max(col-n, 0))
}

type class int

const (
	classSpace class = iota
	classWord
	classPunct
)

func classOf(r rune) class {
	switch {
	case unicode.IsSpace(r):
		return classSpace
	case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
		return classWord
	}
	return classPunct
}

// wordRight is where Ctrl+Right lands from (row, col): past the spaces, then to the end of the
// run of word characters or of punctuation; at a line's end, the start of the next line.
func wordRight(text [][]rune, row, col int) (newRow, newCol int) {
	line := text[row]
	if col >= len(line) {
		if row+1 < len(text) {
			return row + 1, 0
		}
		return row, len(line)
	}
	for col < len(line) && classOf(line[col]) == classSpace {
		col++
	}
	if col < len(line) {
		c := classOf(line[col])
		for col < len(line) && classOf(line[col]) == c {
			col++
		}
	}
	return row, col
}

// wordLeft mirrors wordRight: to the start of the previous run; at a line's start, the end of the
// previous line.
func wordLeft(text [][]rune, row, col int) (newRow, newCol int) {
	line := text[row]
	col = min(col, len(line))
	if col == 0 {
		if row > 0 {
			return row - 1, len(text[row-1])
		}
		return 0, 0
	}
	for col > 0 && classOf(line[col-1]) == classSpace {
		col--
	}
	if col > 0 {
		c := classOf(line[col-1])
		for col > 0 && classOf(line[col-1]) == c {
			col--
		}
	}
	return row, col
}
