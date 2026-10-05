package codeedit

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	rw "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// View renders ta like textarea.View, with the gutter's line numbers and each rune colored by its
// token. The textarea exposes no per-rune styling hook, so this mirrors its View for the rows on
// screen (bubbles v1.0.0; highlight_test.go pins the same text and layout as textarea.View) and
// reads the private viewport's scroll state. Anything unexpected falls back to textarea.View,
// whose prompt still carries the numbers.
func (h *Highlighter) View(ta *textarea.Model) string {
	if h == nil {
		h = &Highlighter{}
	}
	text := ta.Value()
	h.tokenize(text)
	h.fitGutter(ta, strings.Count(text, "\n")+1)
	// textarea.View refreshes the viewport's content, which its own scrolling (repositionView)
	// measures: without it the view stops following the cursor. Only an edit or resize changes it.
	// It runs after wrap: the gutter prompt reads h.labels.
	if width := ta.Width(); h.wraps == nil || width != h.wrapWidth || ta.Height() != h.vpHeight {
		h.wrap(width)
		h.vpHeight = ta.Height()
		_ = ta.View()
	}
	if text == "" {
		return ta.View() // the placeholder, if any
	}
	vp, ok := viewportOf(ta)
	if !ok || vp.lines != h.rowCount()+ta.Height()+1 {
		return ta.View()
	}
	return h.render(ta, vp)
}

// viewportState is what render needs from textarea's private viewport.
type viewportState struct {
	yOffset, xOffset, width, height, lines, longest int
}

func viewportOf(ta *textarea.Model) (viewportState, bool) {
	p := reflect.ValueOf(ta).Elem().FieldByName("viewport")
	if !p.IsValid() || p.Kind() != reflect.Pointer || p.IsNil() {
		return viewportState{}, false
	}
	v := p.Elem()
	ok := true
	num := func(name string) int {
		f := v.FieldByName(name)
		if !f.IsValid() || !f.CanInt() {
			ok = false
			return 0
		}
		return int(f.Int())
	}
	vs := viewportState{
		yOffset: num("YOffset"),
		xOffset: num("xOffset"),
		width:   num("Width"),
		height:  num("Height"),
		longest: num("longestLineWidth"),
	}
	if l := v.FieldByName("lines"); l.IsValid() && l.Kind() == reflect.Slice {
		vs.lines = l.Len()
	} else {
		ok = false
	}
	return vs, ok
}

func (h *Highlighter) wrap(width int) {
	lines := strings.Split(h.text, "\n")
	h.wrapWidth = width
	h.wraps = make([][][]rune, len(lines))
	h.labels = h.labels[:0]
	for i, l := range lines {
		h.wraps[i] = wrap([]rune(l), width)
		h.labels = append(h.labels, strconv.Itoa(i+1))
		for range len(h.wraps[i]) - 1 {
			h.labels = append(h.labels, "")
		}
	}
}

// minGutterDigits is the narrowest number column; a longer buffer widens it for every row.
const minGutterDigits = 3

// gutterLabel is the number cell after the prompt bar: right-aligned in a column digits wide, so
// the text starts on the same column whatever the line number's length.
func gutterLabel(digits int, label string) string {
	return fmt.Sprintf(" %*s ", digits, label)
}

// prompt is the textarea's prompt func: the bar plus displayRow's line number.
func (h *Highlighter) prompt(displayRow int) string {
	label := ""
	if displayRow < len(h.labels) {
		label = h.labels[displayRow]
	}
	return h.bar + gutterLabel(h.digits, label)
}

// fitGutter sizes the number column for lines and binds the prompt to h.
func (h *Highlighter) fitGutter(ta *textarea.Model, lines int) {
	h.bar, h.digits = ta.Prompt, max(minGutterDigits, len(strconv.Itoa(lines)))
	setGutter(ta, h.digits, h.prompt)
}

// setGutter installs prompt as a gutter digits wide. SetWidth subtracts the prompt's width from
// the outer width it is given, so a width change re-applies the current outer width to rewrap.
func setGutter(ta *textarea.Model, digits int, prompt func(int) string) {
	width := uniseg.StringWidth(ta.Prompt) + digits + len(gutterLabel(0, ""))
	vp, ok := viewportOf(ta)
	ta.SetPromptFunc(width, prompt)
	if ok && vp.width-ta.Width() != width {
		ta.SetWidth(vp.width + baseStyle(ta).GetHorizontalFrameSize())
	}
}

func baseStyle(ta *textarea.Model) lipgloss.Style {
	if ta.Focused() {
		return ta.FocusedStyle.Base
	}
	return ta.BlurredStyle.Base
}

func (h *Highlighter) rowCount() int {
	n := 0
	for _, w := range h.wraps {
		n += len(w)
	}
	return n
}

// lineStyles are textarea's computed styles (Style.computed*) for the focus state.
type lineStyles struct {
	cursorLine, text, cursorLineNumber, lineNumber, endOfBuffer, prompt lipgloss.Style
}

func stylesOf(ta *textarea.Model) lineStyles {
	s := ta.BlurredStyle
	if ta.Focused() {
		s = ta.FocusedStyle
	}
	return lineStyles{
		cursorLine:       s.CursorLine.Inherit(s.Base).Inline(true),
		text:             s.Text.Inherit(s.Base).Inline(true),
		cursorLineNumber: s.CursorLineNumber.Inherit(s.CursorLine).Inherit(s.Base).Inline(true),
		lineNumber:       s.LineNumber.Inherit(s.Base).Inline(true),
		endOfBuffer:      s.EndOfBuffer.Inherit(s.Base).Inline(true),
		prompt:           s.Prompt.Inherit(s.Base).Inline(true),
	}
}

func (h *Highlighter) render(ta *textarea.Model, vp viewportState) string {
	st := stylesOf(ta)
	total := h.rowCount()
	top, bottom := max(0, vp.yOffset), min(vp.yOffset+vp.height, vp.lines)

	rows := make([]string, 0, max(0, bottom-top))
	display := 0
	for l := range h.wraps {
		if display >= bottom {
			break
		}
		if display+len(h.wraps[l]) <= top {
			display += len(h.wraps[l])
			continue
		}
		start := 0
		for wl, wrapped := range h.wraps[l] {
			if display >= top && display < bottom {
				rows = append(rows, h.renderRow(ta, &st, l, wl, wrapped, start))
			}
			start += len(wrapped)
			display++
		}
	}
	for d := max(display, top); d < bottom; d++ {
		if d >= total+ta.Height() {
			rows = append(rows, "") // the split's empty tail after the final newline
			continue
		}
		gap := strings.Repeat(" ", max(0, ta.Width()-1))
		rows = append(rows, st.prompt.Render(h.prompt(d))+
			st.endOfBuffer.Render(string(ta.EndOfBufferCharacter)+gap))
	}
	return viewportView(ta, vp, rows)
}

// renderRow is one display row of textarea.View: prompt bar, line number, text (with the
// cursor), padding. start is wrapped's rune offset in its logical line l.
func (h *Highlighter) renderRow(
	ta *textarea.Model, st *lineStyles, l, wl int, wrapped []rune, start int,
) string {
	row, li := ta.Line(), ta.LineInfo()
	style, lnStyle := st.text, st.lineNumber
	if l == row {
		style, lnStyle = st.cursorLine, st.cursorLineNumber
	}
	var b strings.Builder
	b.WriteString(style.Render(st.prompt.Render(h.bar)))
	label := ""
	if wl == 0 {
		label = strconv.Itoa(l + 1)
	}
	b.WriteString(style.Render(lnStyle.Render(gutterLabel(h.digits, label))))
	width := ta.Width()
	strwidth := uniseg.StringWidth(string(wrapped))
	padding := width - strwidth
	if strwidth > width {
		wrapped = []rune(strings.TrimSuffix(string(wrapped), " "))
		padding -= width - strwidth
	}
	if l == row && li.RowOffset == wl {
		h.paint(&b, &style, wrapped[:li.ColumnOffset], l, start)
		cur := ta.Cursor
		cur.TextStyle = st.cursorLine
		if li.StartColumn+li.ColumnOffset >= lineLen(h.wraps[l]) && li.CharOffset >= width {
			cur.SetChar(" ")
			b.WriteString(cur.View())
		} else {
			cur.SetChar(string(wrapped[li.ColumnOffset]))
			b.WriteString(style.Render(cur.View()))
			h.paint(&b, &style, wrapped[li.ColumnOffset+1:], l, start+li.ColumnOffset+1)
		}
	} else {
		h.paint(&b, &style, wrapped, l, start)
	}
	b.WriteString(style.Render(strings.Repeat(" ", max(0, padding))))
	return b.String()
}

// lineLen is the logical line's rune count: its wrapped rows minus wrap's one trailing space.
func lineLen(rows [][]rune) int {
	n := -1
	for _, r := range rows {
		n += len(r)
	}
	return n
}

// paint writes seg (runes from offset start of line l) in runs of one token class and find mark,
// each run's color layered over the row's style. An all-unstyled segment is one Render, as
// textarea does.
func (h *Highlighter) paint(b *strings.Builder, row *lipgloss.Style, seg []rune, l, start int) {
	cls := h.classes[l]
	classAt := func(i int) uint8 {
		if j := start + i; j < len(cls) {
			return cls[j]
		}
		return clsNone
	}
	if len(seg) == 0 {
		b.WriteString(row.Render(""))
		return
	}
	for i := 0; i < len(seg); {
		c, mk := classAt(i), h.markAt(l, start+i)
		j := i + 1
		for j < len(seg) && classAt(j) == c && h.markAt(l, start+j) == mk {
			j++
		}
		b.WriteString(markStyle(c, mk).Inherit(*row).Render(string(seg[i:j])))
		i = j
	}
}

const (
	markNone = iota
	markMatch
	markCurrent
)

func (h *Highlighter) markAt(l, col int) int {
	if h.hasCurrent && h.current.line == l && col >= h.current.start && col < h.current.end {
		return markCurrent
	}
	for _, m := range h.marks[l] {
		if col >= m.start && col < m.end {
			return markMatch
		}
	}
	return markNone
}

// markStyle is token class c's color; a find match keeps it (plain text takes the keyword color),
// underlined and bold, and the current match is reversed so the color becomes its background.
func markStyle(c uint8, mk int) lipgloss.Style {
	if mk == markNone {
		return tokenStyles[c]
	}
	if c == clsNone || c >= clsStrong && c <= clsEmph {
		c = clsKeyword
	}
	s := tokenStyles[c].Underline(true).Bold(true)
	if mk == markCurrent {
		s = s.Reverse(true)
	}
	return s
}

// cursorScreenRow is the cursor's row in the last View, -1 when unknown.
func (h *Highlighter) cursorScreenRow(ta *textarea.Model) int {
	vp, ok := viewportOf(ta)
	if !ok || ta.Line() >= len(h.wraps) {
		return -1
	}
	d := ta.LineInfo().RowOffset
	for _, w := range h.wraps[:ta.Line()] {
		d += len(w)
	}
	return d - vp.yOffset
}

// viewportView is viewport.Model.View for rows (the visible slice) with its default, frameless
// style, then textarea's Base style around it.
func viewportView(ta *textarea.Model, vp viewportState, rows []string) string {
	w, hgt := vp.width, vp.height
	if (vp.xOffset != 0 || vp.longest > w) && w != 0 {
		for i := range rows {
			rows[i] = ansi.Cut(rows[i], vp.xOffset, vp.xOffset+w)
		}
	}
	contents := lipgloss.NewStyle().
		Width(w).
		Height(hgt).
		MaxHeight(hgt).
		MaxWidth(w).
		Render(strings.Join(rows, "\n"))
	return baseStyle(ta).Render(lipgloss.NewStyle().Render(contents))
}

// wrap is bubbles v1.0.0 textarea's (unexported) soft-wrap, copied verbatim (MIT,
// github.com/charmbracelet/bubbles): LineInfo is computed from it, so rows must split identically.
func wrap(runes []rune, width int) [][]rune {
	var (
		lines  = [][]rune{{}}
		word   = []rune{}
		row    int
		spaces int
	)

	for _, r := range runes {
		if unicode.IsSpace(r) {
			spaces++
		} else {
			word = append(word, r)
		}

		if spaces > 0 { //nolint:nestif // verbatim upstream
			if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces > width {
				row++
				lines = append(lines, []rune{})
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			} else {
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			}
		} else {
			lastCharLen := rw.RuneWidth(word[len(word)-1])
			if uniseg.StringWidth(string(word))+lastCharLen > width {
				if len(lines[row]) > 0 {
					row++
					lines = append(lines, []rune{})
				}
				lines[row] = append(lines[row], word...)
				word = nil
			}
		}
	}

	if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces >= width {
		lines = append(lines, []rune{})
		lines[row+1] = append(lines[row+1], word...)
		spaces++
		lines[row+1] = append(lines[row+1], repeatSpaces(spaces)...)
	} else {
		lines[row] = append(lines[row], word...)
		spaces++
		lines[row] = append(lines[row], repeatSpaces(spaces)...)
	}

	return lines
}

func repeatSpaces(n int) []rune {
	return []rune(strings.Repeat(" ", n))
}
