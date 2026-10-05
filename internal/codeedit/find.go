package codeedit

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// span is a match: runes [start, end) of logical line.
type span struct{ line, start, end int }

type findButton int

const (
	btnNext findButton = iota // Search / Replace
	btnAll                    // Search all / Replace all
	btnClose
	buttonCount
)

// Finder is the Ctrl+F (find) / Ctrl+R (find & replace) box drawn over a codeedit textarea.
// Matching is literal and smart-case (case-sensitive only when the query has an upper-case
// letter); matches never span lines. The zero value is a closed, find-only Finder.
type Finder struct {
	editable bool // Ctrl+R allowed
	open     bool
	replace  bool
	showAll  bool // Search all: mark every match, not just the current one

	query, with textinput.Model
	focus       int // 0 query, 1 replacement (replace mode), then the buttons

	matches []span
	cur     int // index into matches, -1 none
	text    string
	pattern string
	status  string
}

// NewFinder returns a closed Finder; editable enables Ctrl+R (replace).
func NewFinder(editable bool) Finder {
	return Finder{editable: editable, cur: -1}
}

// Active reports whether the box is open (it then takes the keyboard).
func (f *Finder) Active() bool { return f.open }

// Close hides the box and drops every match mark.
func (f *Finder) Close() {
	f.open, f.showAll, f.matches, f.cur, f.status = false, false, nil, -1, ""
	f.query.Blur()
	f.with.Blur()
}

// Update handles Ctrl+F / Ctrl+R and, while the box is open, every key but Ctrl+S / Ctrl+C (left
// to the host). It reports whether it consumed msg.
func (f *Finder) Update(ta *textarea.Model, msg tea.KeyMsg) bool {
	switch k := msg.String(); k {
	case "ctrl+f":
		f.start(ta, false)
		return true
	case "ctrl+r":
		if !f.editable {
			return f.open
		}
		f.start(ta, true)
		return true
	case "ctrl+s", "ctrl+c":
		return false
	}
	if !f.open {
		return false
	}
	switch msg.String() {
	case "esc":
		f.Close()
	case "tab", "down":
		f.setFocus(f.focus + 1)
	case "shift+tab", "up":
		f.setFocus(f.focus - 1)
	case "left", "right":
		switch {
		case f.focus < f.fields():
			f.typeKey(ta, msg)
		case msg.String() == "left":
			f.setFocus(max(f.focus-1, f.fields()))
		default:
			f.setFocus(min(f.focus+1, f.fields()+int(buttonCount)-1))
		}
	case "enter":
		f.activate(ta)
	default:
		f.typeKey(ta, msg)
	}
	return true
}

func (f *Finder) start(ta *textarea.Model, replace bool) {
	if !f.open {
		f.query = newFindInput("text to find")
		f.with = newFindInput("replacement")
		f.cur, f.showAll, f.status = -1, false, ""
		if f.pattern != "" {
			f.query.SetValue(f.pattern)
			f.query.CursorEnd()
		}
	}
	f.open, f.replace = true, replace
	f.setFocus(0)
	f.refresh(ta)
	f.selectFrom(ta, ta.Line(), column(ta))
}

func newFindInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = placeholder
	ti.Cursor.SetMode(cursor.CursorStatic)
	return ti
}

func (f *Finder) fields() int {
	if f.replace {
		return 2
	}
	return 1
}

func (f *Finder) setFocus(i int) {
	n := f.fields() + int(buttonCount)
	f.focus = (i%n + n) % n
	f.query.Blur()
	f.with.Blur()
	switch f.focus {
	case 0:
		f.query.Focus()
	case 1:
		if f.replace {
			f.with.Focus()
		}
	}
}

func (f *Finder) typeKey(ta *textarea.Model, msg tea.KeyMsg) {
	switch {
	case f.focus == 0:
		f.query, _ = f.query.Update(msg)
		if f.query.Value() != f.pattern {
			f.refresh(ta)
			f.selectFrom(ta, ta.Line(), column(ta))
		}
	case f.focus == 1 && f.replace:
		f.with, _ = f.with.Update(msg)
	}
}

// activate is Enter: on the query, Search; on the replacement, Replace; on a button, that button.
func (f *Finder) activate(ta *textarea.Model) {
	btn := btnNext
	switch {
	case f.focus == 0:
		f.next(ta)
		return
	case f.focus >= f.fields():
		btn = findButton(f.focus - f.fields())
	}
	switch btn {
	case btnNext:
		if f.replace {
			f.replaceOne(ta)
		} else {
			f.next(ta)
		}
	case btnAll:
		if f.replace {
			f.replaceAll(ta)
		} else {
			f.showAll = true
			f.report()
		}
	case btnClose:
		f.Close()
	}
}

// refresh recomputes the matches when the buffer or the query changed.
func (f *Finder) refresh(ta *textarea.Model) {
	text, pattern := ta.Value(), f.query.Value()
	if text == f.text && pattern == f.pattern && f.matches != nil {
		return
	}
	f.text, f.pattern = text, pattern
	f.matches = findAll(text, pattern)
	f.cur = -1
}

// selectFrom makes the first match at or after (row, col) current, wrapping to the top, and puts
// the cursor on it.
func (f *Finder) selectFrom(ta *textarea.Model, row, col int) {
	f.cur = -1
	for i, m := range f.matches {
		if m.line > row || (m.line == row && m.start >= col) {
			f.cur = i
			break
		}
	}
	if f.cur < 0 && len(f.matches) > 0 {
		f.cur = 0
	}
	f.jump(ta)
}

// next makes the match after the current one current (wrapping).
func (f *Finder) next(ta *textarea.Model) {
	f.refresh(ta)
	if len(f.matches) == 0 {
		f.cur = -1
		f.report()
		return
	}
	if f.cur < 0 {
		f.selectFrom(ta, ta.Line(), column(ta))
		return
	}
	f.cur = (f.cur + 1) % len(f.matches)
	f.jump(ta)
}

func (f *Finder) jump(ta *textarea.Model) {
	if f.cur >= 0 {
		m := f.matches[f.cur]
		moveTo(ta, m.line, m.start)
		reposition(ta)
	}
	f.report()
}

func (f *Finder) report() {
	switch {
	case f.pattern == "":
		f.status = ""
	case len(f.matches) == 0:
		f.status = "no matches"
	case f.cur >= 0:
		f.status = fmt.Sprintf("%d of %d", f.cur+1, len(f.matches))
	default:
		f.status = fmt.Sprintf("%d matches", len(f.matches))
	}
}

// replaceOne replaces the current match and moves on to the next one after it.
func (f *Finder) replaceOne(ta *textarea.Model) {
	f.refresh(ta)
	if f.cur < 0 {
		f.next(ta)
		return
	}
	m := f.matches[f.cur]
	lines := strings.Split(f.text, "\n")
	r := []rune(lines[m.line])
	repl := f.with.Value()
	lines[m.line] = string(r[:m.start]) + repl + string(r[m.end:])
	ta.SetValue(strings.Join(lines, "\n"))
	after := m.start + len([]rune(repl))
	moveTo(ta, m.line, after) // SetValue leaves the cursor at the end
	reposition(ta)
	f.refresh(ta)
	f.selectFrom(ta, m.line, after)
}

// replaceAll replaces every match; the cursor stays on its line.
func (f *Finder) replaceAll(ta *textarea.Model) {
	f.refresh(ta)
	n := len(f.matches)
	if n == 0 {
		f.report()
		return
	}
	row, col := ta.Line(), column(ta)
	lines := strings.Split(f.text, "\n")
	repl := f.with.Value()
	// Back to front, so earlier offsets on the same line stay valid.
	for i := n - 1; i >= 0; i-- {
		m := f.matches[i]
		r := []rune(lines[m.line])
		lines[m.line] = string(r[:m.start]) + repl + string(r[m.end:])
	}
	ta.SetValue(strings.Join(lines, "\n"))
	moveTo(ta, row, min(col, len([]rune(lines[row]))))
	reposition(ta)
	f.refresh(ta)
	f.status = fmt.Sprintf("replaced %d", n)
}

// findAll is every non-overlapping, smart-case match of pattern, line by line.
func findAll(text, pattern string) []span {
	needle := []rune(pattern)
	if len(needle) == 0 {
		return []span{}
	}
	fold := !strings.ContainsFunc(pattern, unicode.IsUpper)
	if fold {
		needle = lowerRunes(needle)
	}
	out := []span{}
	for l, line := range strings.Split(text, "\n") {
		hay := []rune(line)
		if fold {
			hay = lowerRunes(hay)
		}
		for i := 0; i+len(needle) <= len(hay); {
			if runesEqual(hay[i:i+len(needle)], needle) {
				out = append(out, span{l, i, i + len(needle)})
				i += len(needle)
				continue
			}
			i++
		}
	}
	return out
}

// lowerRunes lowers rune by rune, keeping offsets (strings.ToLower may change the rune count).
func lowerRunes(rs []rune) []rune {
	out := make([]rune, len(rs))
	for i, r := range rs {
		out[i] = unicode.ToLower(r)
	}
	return out
}

func runesEqual(a, b []rune) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// marks is what the highlighter underlines: the current match, plus every match after Search all.
func (f *Finder) marks() (map[int][]span, span, bool) {
	if !f.open {
		return nil, span{}, false
	}
	out := map[int][]span{}
	if f.showAll {
		for _, m := range f.matches {
			out[m.line] = append(out[m.line], m)
		}
	}
	var cur span
	if f.cur >= 0 {
		cur = f.matches[f.cur]
		if !f.showAll {
			out[cur.line] = append(out[cur.line], cur)
		}
	}
	return out, cur, f.cur >= 0
}

var (
	findBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "97", Dark: "141"}).
			Padding(0, 1)
	findTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "162", Dark: "212"})
	findDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "245"})
	findButtonStyle = lipgloss.NewStyle().Padding(0, 1).
			Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "252"})
	findButtonFocus = findButtonStyle.Bold(true).
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.AdaptiveColor{Light: "97", Dark: "62"})
)

// findBoxMaxWidth caps the box; it shrinks to fit narrow panes.
const findBoxMaxWidth = 46

// View renders ta through h (as Highlighter.View, with the matches marked) and, while open, the box
// over its top-right corner, or the bottom-right one when the cursor is under the top.
func (f *Finder) View(h *Highlighter, ta *textarea.Model) string {
	if h == nil {
		h = &Highlighter{}
	}
	if f.open {
		f.refresh(ta)
	}
	h.marks, h.current, h.hasCurrent = f.marks()
	base := h.View(ta)
	if !f.open {
		return base
	}
	rows := strings.Split(base, "\n")
	width := 0
	for _, r := range rows {
		width = max(width, ansi.StringWidth(r))
	}
	box := strings.Split(f.box(min(findBoxMaxWidth, width-2)), "\n")
	boxW := ansi.StringWidth(box[0])
	y := 0
	if c := h.cursorScreenRow(ta); c >= 0 && c < len(box) && c < len(rows)-len(box) {
		y = len(rows) - len(box)
	}
	x := max(0, width-boxW-1)
	for i, b := range box {
		if y+i < len(rows) {
			rows[y+i] = overlay(rows[y+i], b, x, boxW)
		}
	}
	return strings.Join(rows, "\n")
}

// overlay puts box over line at column x, keeping line's cells on both sides.
func overlay(line, box string, x, boxW int) string {
	left := ansi.Truncate(line, x, "")
	if w := ansi.StringWidth(left); w < x {
		left += strings.Repeat(" ", x-w)
	}
	return left + "\x1b[0m" + box + "\x1b[0m" + ansi.TruncateLeft(line, x+boxW, "")
}

func (f *Finder) box(width int) string {
	inner := max(width-findBoxStyle.GetHorizontalFrameSize(), 30)
	title := "Find"
	labels := []string{"Find"}
	buttons := []string{"Search", "Search all", "Esc"}
	if f.replace {
		title = "Find & replace"
		labels = append(labels, "Replace")
		buttons = []string{"Replace", "Replace all", "Esc"}
	}
	labelW := 0
	for _, l := range labels {
		labelW = max(labelW, len(l))
	}
	var b strings.Builder
	gap := max(1, inner-lipgloss.Width(title)-lipgloss.Width(f.status))
	b.WriteString(findTitleStyle.Render(title) + strings.Repeat(" ", gap) + findDimStyle.Render(f.status))
	inputs := []*textinput.Model{&f.query, &f.with}
	for i, l := range labels {
		in := inputs[i]
		in.Width = max(inner-labelW-3, 1)
		b.WriteString("\n" + findDimStyle.Render(fmt.Sprintf("%-*s ", labelW, l)) + in.View())
	}
	b.WriteString("\n")
	for i, name := range buttons {
		style := findButtonStyle
		if f.focus == f.fields()+i {
			style = findButtonFocus
		}
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(style.Render(name))
	}
	b.WriteString("\n" + findDimStyle.Render(ansi.Truncate("tab next · ↵ activate · esc close", inner, "…")))
	return findBoxStyle.Width(inner + 2).Render(b.String())
}
