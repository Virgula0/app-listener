package codeedit

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	rw "github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
	"github.com/rivo/uniseg"
)

// Selection is a codeedit textarea's text selection, from an anchor to the cursor. Shift (or
// Ctrl+Shift) with the movement keys, Ctrl+A or a mouse drag make one; Ctrl+C copies it, Esc drops
// it, and in an editor Ctrl+X cuts it while typing, Backspace, Delete or a paste replace it. The
// zero value selects nothing.
type Selection struct {
	anchor   pos
	active   bool
	dragging bool
	note     string
}

type pos struct{ line, col int }

func (a pos) before(b pos) bool { return a.line < b.line || a.line == b.line && a.col < b.col }

func cursorPos(ta *textarea.Model) pos { return pos{ta.Line(), column(ta)} }

// selectMoves maps a selecting key to the movement it extends the selection with.
var selectMoves = map[string]tea.KeyType{
	"shift+left":       tea.KeyLeft,
	"shift+right":      tea.KeyRight,
	"shift+up":         tea.KeyUp,
	"shift+down":       tea.KeyDown,
	"shift+home":       tea.KeyHome,
	"shift+end":        tea.KeyEnd,
	"ctrl+shift+left":  tea.KeyCtrlLeft,
	"ctrl+shift+right": tea.KeyCtrlRight,
	"ctrl+shift+up":    tea.KeyUp,
	"ctrl+shift+down":  tea.KeyDown,
	"ctrl+shift+home":  tea.KeyCtrlHome,
	"ctrl+shift+end":   tea.KeyCtrlEnd,
}

// Clear drops the selection.
func (s *Selection) Clear() { s.active, s.dragging = false, false }

// Dragging reports whether a mouse drag is extending the selection: the host should route mouse
// motion here even once it leaves the editor.
func (s *Selection) Dragging() bool { return s.dragging }

// TakeNote returns, once, what the last key did that the host may show (a copy's size).
func (s *Selection) TakeNote() string {
	n := s.note
	s.note = ""
	return n
}

// bounds is the selection in text order; ok is false when nothing is selected. A nil Selection
// selects nothing.
func (s *Selection) bounds(ta *textarea.Model) (from, to pos, ok bool) {
	if s == nil || !s.active {
		return pos{}, pos{}, false
	}
	from, to = s.anchor, cursorPos(ta)
	if to.before(from) {
		from, to = to, from
	}
	return from, to, from != to
}

// Text is the selected text, "" when nothing is selected.
func (s *Selection) Text(ta *textarea.Model) string {
	from, to, ok := s.bounds(ta)
	if !ok {
		return ""
	}
	rs := []rune(ta.Value())
	text := lines(ta)
	return string(rs[offset(text, from):offset(text, to)])
}

// offset is p's rune index in the buffer whose lines are text.
func offset(text [][]rune, p pos) int {
	n := 0
	for _, l := range text[:p.line] {
		n += len(l) + 1
	}
	return n + min(p.col, len(text[p.line]))
}

// Key handles the selection keys; editable enables cut and replace. It reports whether it consumed
// msg: a key it leaves (typing over a selection, which it first deletes) goes on to the host.
func (s *Selection) Key(ta *textarea.Model, msg tea.KeyMsg, editable bool) (bool, tea.Cmd) {
	k := msg.String()
	if s.extend(ta, k) {
		return true, nil
	}
	if k == "ctrl+v" && editable && clipboardText() != "" {
		// The textarea's own paste reads the system clipboard, which a root process under sudo
		// rarely reaches: paste what was copied here.
		s.delete(ta)
		ta.InsertString(clipboardText())
		reposition(ta)
		return true, nil
	}
	if _, _, ok := s.bounds(ta); !ok {
		s.Clear()
		return false, nil
	}
	switch {
	case k == "ctrl+c":
		return true, s.copy(ta)
	case k == "esc":
		s.Clear()
		return true, nil
	case !editable:
		s.Clear()
		return false, nil
	case k == "ctrl+x":
		cmd := s.copy(ta)
		s.delete(ta)
		return true, cmd
	case replacing[k] || msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace:
		s.delete(ta)
		return deleting[k], nil
	}
	s.Clear()
	return false, nil
}

// replacing are the keys that remove the selection first; deleting ones are then done.
var (
	replacing = map[string]bool{
		"backspace": true, "ctrl+h": true, "delete": true, "enter": true, "ctrl+m": true, "ctrl+v": true,
	}
	deleting = map[string]bool{"backspace": true, "ctrl+h": true, "delete": true}
)

// extend handles the keys that make or grow a selection: Shift+movement and Ctrl+A.
func (s *Selection) extend(ta *textarea.Model, k string) bool {
	if k == "ctrl+a" {
		text := lines(ta)
		moveTo(ta, 0, 0)
		s.anchor, s.active, s.dragging = pos{}, true, false
		moveTo(ta, len(text)-1, len(text[len(text)-1]))
		reposition(ta)
		return true
	}
	t, ok := selectMoves[k]
	if !ok {
		return false
	}
	if !s.active || s.dragging {
		s.anchor, s.active, s.dragging = cursorPos(ta), true, false
	}
	km := tea.KeyMsg{Type: t}
	if !Navigate(ta, km) {
		*ta, _ = ta.Update(km)
	}
	return true
}

func (s *Selection) copy(ta *textarea.Model) tea.Cmd {
	text := s.Text(ta)
	setClipboard(text)
	s.note = fmt.Sprintf("copied %d characters", len([]rune(text)))
	return func() tea.Msg {
		termenv.NewOutput(clipboardOut).Copy(text)
		return nil
	}
}

// deleteKeysMax bounds deleting a selection key by key, which keeps the view where it is; a larger
// one is cut out of the value, which scrolls the view to bring the cursor back.
const deleteKeysMax = 4096

func (s *Selection) delete(ta *textarea.Model) {
	from, to, ok := s.bounds(ta)
	s.Clear()
	if !ok {
		return
	}
	text := lines(ta)
	n := offset(text, to) - offset(text, from)
	if n <= deleteKeysMax {
		moveTo(ta, from.line, from.col)
		for range n {
			*ta, _ = ta.Update(tea.KeyMsg{Type: tea.KeyDelete})
		}
	} else {
		rs := []rune(ta.Value())
		ta.SetValue(string(rs[:offset(text, from)]) + string(rs[offset(text, to):]))
		moveTo(ta, from.line, from.col)
	}
	reposition(ta)
}

// clipboardOut receives the OSC 52 sequence that puts copied text on the terminal's clipboard:
// it reaches the user's clipboard through sudo, ssh and tmux, where no clipboard tool would.
var clipboardOut io.Writer = os.Stdout

var (
	clipMu   sync.Mutex
	clipText string
)

func setClipboard(text string) {
	clipMu.Lock()
	clipText = text
	clipMu.Unlock()
}

func clipboardText() string {
	clipMu.Lock()
	defer clipMu.Unlock()
	return clipText
}

// wheelStep is how many lines one mouse-wheel notch scrolls.
const wheelStep = 3

// Mouse handles a mouse event at (x, y), in cells from the editor's top-left corner as h last
// drew it: a left press puts the cursor there and starts a selection, dragging extends it (past
// the edges, scrolling), the wheel scrolls. It reports whether it used msg.
func (s *Selection) Mouse(ta *textarea.Model, h *Highlighter, msg tea.MouseMsg, x, y int) bool {
	switch {
	case msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown:
		if !s.dragging {
			s.Clear()
		}
		for range wheelStep {
			if msg.Button == tea.MouseButtonWheelUp {
				ta.CursorUp()
			} else {
				ta.CursorDown()
			}
		}
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
		p, ok := h.posAt(ta, x, y)
		if !ok {
			return false
		}
		moveTo(ta, p.line, p.col)
		s.anchor, s.active, s.dragging = p, true, true
	case msg.Action == tea.MouseActionMotion && s.dragging:
		if p, ok := h.posAt(ta, x, y); ok {
			moveTo(ta, p.line, p.col)
		}
	case msg.Action == tea.MouseActionRelease && s.dragging:
		s.dragging = false
	default:
		return false
	}
	reposition(ta)
	return true
}

// posAt is the text position drawn at (x, y) in the last View; a row above or below the view, or
// a column in the gutter, lands on the nearest text.
func (h *Highlighter) posAt(ta *textarea.Model, x, y int) (pos, bool) {
	vp, ok := viewportOf(ta)
	if h == nil || !ok || len(h.wraps) == 0 {
		return pos{}, false
	}
	st := baseStyle(ta)
	x -= st.GetMarginLeft() + st.GetBorderLeftSize() + st.GetPaddingLeft()
	y -= st.GetMarginTop() + st.GetBorderTopSize() + st.GetPaddingTop()
	x -= uniseg.StringWidth(h.bar) + h.digits + len(gutterLabel(0, ""))
	d := min(max(y+vp.yOffset, 0), h.rowCount()-1)
	for l, rows := range h.wraps {
		if d >= len(rows) {
			d -= len(rows)
			continue
		}
		start := 0
		for _, r := range rows[:d] {
			start += len(r)
		}
		col := start + cellToRune(rows[d], x+vp.xOffset)
		return pos{l, min(col, lineLen(rows))}, true
	}
	return pos{}, false
}

// cellToRune is the index of the rune drawn at cell x of row.
func cellToRune(row []rune, x int) int {
	w := 0
	for i, r := range row {
		w += rw.RuneWidth(r)
		if x < w {
			return i
		}
	}
	return len(row)
}

// selectedStyle is a selected cell's background, under its token color.
var selectedStyle = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "153", Dark: "24"})

// selected reports whether rune col of line l is in the selection; a line's end cell stands for
// its newline.
func (h *Highlighter) selected(l, col int) bool {
	if !h.hasSel {
		return false
	}
	from, to := h.selFrom, h.selTo
	return (l > from.line || l == from.line && col >= from.col) &&
		(l < to.line || l == to.line && col < to.col)
}

// markSelection records sel's range for the next render.
func (h *Highlighter) markSelection(ta *textarea.Model, sel *Selection) {
	h.selFrom, h.selTo, h.hasSel = sel.bounds(ta)
}
