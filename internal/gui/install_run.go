//go:build gui

package gui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

const maxLogLines = 5000

type logLine struct {
	level log.Level
	at    time.Time
	text  string
}

type stageState int

const (
	stageActive stageState = iota
	stageDone
	stageFailed
)

type stageItem struct {
	name  string
	state stageState
}

// runView is the panel of one installer or updater run: its stages as a stepper, the current progress bar and
// the run's log. Every method runs on the UI goroutine.
type runView struct {
	name   string
	root   fyne.CanvasObject
	tile   *iconTile
	title  *widget.Label
	state  *pill
	head   *fyne.Container
	steps  *fyne.Container
	bar    *widget.ProgressBar
	barLbl *widget.Label
	barBox *fyne.Container
	logs   *logGrid
	result *widget.Label
	// logBox holds barBox and result: showing or hiding them needs its relayout.
	logBox *fyne.Container
	back   *handButton

	stages []stageItem
	lines  []logLine
}

func newRunView(s *shell, back string, onBack func()) *runView {
	r := &runView{
		tile:   newIconTile("terminal-window", theme.ColorNamePrimary, 44),
		title:  boldLabel(""),
		state:  newPill("", colorMuted),
		steps:  container.NewVBox(),
		bar:    widget.NewProgressBar(),
		barLbl: mutedText(""),
		result: widget.NewLabel(""),
	}
	r.result.Wrapping = fyne.TextWrapWord
	r.barBox = container.NewVBox(r.barLbl, r.bar)
	r.barBox.Hide()
	r.logs = newLogGrid(func(row int) {
		if row >= 0 && row < len(r.lines) {
			s.app.Clipboard().SetContent(r.lines[row].text)
		}
	})

	r.back = newTextButton(back, onBack)
	copyLog := newButton("Copy log", icon("copy"), func() {
		s.app.Clipboard().SetContent(r.logText())
	})
	copyLog.Importance = widget.LowImportance
	r.head = container.NewBorder(nil, nil, r.tile, container.NewHBox(copyLog, r.back),
		container.NewVBox(r.title, container.NewHBox(r.state)))

	stepsCard := newPanel(container.NewBorder(boldLabel("Steps"), nil, nil, nil, newFadeVScroll(r.steps)))
	hint := widget.NewLabel("click a line to copy it")
	hint.Importance = widget.LowImportance
	hint.SizeName = theme.SizeNameCaptionText
	r.logBox = container.NewBorder(container.NewVBox(hfill(boldLabel("Log"), hint), r.barBox),
		r.result, nil, nil, r.logs)
	logCard := newPanel(r.logBox)
	split := container.NewHSplit(stepsCard, logCard)
	split.Offset = 0.28
	r.root = container.NewBorder(container.NewVBox(newPanel(r.head)), nil, nil, nil, split)
	return r
}

// logGrid is the run log: one row per line, scrolling both ways so a long line is never cut. A
// tapped row is handed to onTap.
type logGrid struct {
	widget.TextGrid
	onTap    func(row int)
	rendered bool
}

func newLogGrid(onTap func(row int)) *logGrid {
	g := &logGrid{onTap: onTap}
	g.Scroll = fyne.ScrollBoth
	g.ExtendBaseWidget(g)
	return g
}

func (g *logGrid) CreateRenderer() fyne.WidgetRenderer {
	r := g.TextGrid.CreateRenderer()
	g.rendered = true
	return r
}

func (g *logGrid) Tapped(e *fyne.PointEvent) {
	row, _ := g.CursorLocationForPosition(e.Position)
	g.onTap(row)
}

// scrollToBottom is TextGrid.ScrollToBottom, whose scroller exists only once rendered.
func (g *logGrid) scrollToBottom() {
	if g.rendered {
		g.ScrollToBottom()
	}
}

func logStyle(l log.Level) widget.TextGridStyle {
	var c color.Color
	switch {
	case l <= log.ErrorLevel:
		c = themeColor(theme.ColorNameError)
	case l == log.WarnLevel:
		c = themeColor(theme.ColorNameWarning)
	case l >= log.DebugLevel:
		c = themeColor(colorMuted)
	default:
		return nil
	}
	return &widget.CustomTextGridStyle{FGColor: c}
}

func logRow(l logLine) widget.TextGridRow {
	text := l.at.Format("15:04:05") + "  " + strings.ReplaceAll(l.text, "\n", " ")
	cells := make([]widget.TextGridCell, 0, len(text))
	for _, c := range text {
		cells = append(cells, widget.TextGridCell{Rune: c})
	}
	return widget.TextGridRow{Cells: cells, Style: logStyle(l.level)}
}

// start resets the view for a run named name.
func (r *runView) start(name, iconName string) {
	r.name = name
	r.stages, r.lines = nil, nil
	r.tile.Set(iconName, theme.ColorNamePrimary)
	r.title.SetText(name)
	r.setState("running", theme.ColorNamePrimary)
	r.result.SetText("")
	r.result.Hide()
	r.logBox.Refresh()
	r.back.Disable()
	r.progressDone()
	r.renderSteps()
	r.logs.Rows = nil
	r.logs.Refresh()
}

func (r *runView) setState(text string, c fyne.ThemeColorName) {
	r.state.Set(text, c)
	r.head.Refresh() // the pill's width changed
}

// stage closes the active step and opens name.
func (r *runView) stage(name string) {
	if n := len(r.stages); n > 0 && r.stages[n-1].state == stageActive {
		r.stages[n-1].state = stageDone
	}
	r.stages = append(r.stages, stageItem{name: name})
	r.renderSteps()
}

// finish closes the run: err nil = success, else the active step failed.
func (r *runView) finish(err error, canceled bool) {
	r.progressDone()
	state := stageDone
	switch {
	case err == nil:
		r.setState("done", theme.ColorNameSuccess)
		r.tile.Set("check-circle", theme.ColorNameSuccess)
		r.result.Importance = widget.SuccessImportance
		r.result.SetText("Finished.")
	case canceled:
		state = stageFailed
		r.setState("canceled", theme.ColorNameWarning)
		r.tile.Set("warning-circle", theme.ColorNameWarning)
		r.result.Importance = widget.WarningImportance
		r.result.SetText("Stopped: a prompt was canceled. Steps already completed stay applied.")
	default:
		state = stageFailed
		r.setState("failed", theme.ColorNameError)
		r.tile.Set("x-circle", theme.ColorNameError)
		r.result.Importance = widget.DangerImportance
		r.result.SetText("Failed: " + err.Error())
	}
	if n := len(r.stages); n > 0 && r.stages[n-1].state == stageActive {
		r.stages[n-1].state = state
	}
	r.result.Show()
	r.logBox.Refresh()
	r.renderSteps()
	r.back.Enable()
}

func (r *runView) renderSteps() {
	r.steps.Objects = r.steps.Objects[:0]
	if len(r.stages) == 0 {
		r.steps.Add(mutedText("Starting…"))
	}
	for _, st := range r.stages {
		name, c := "circle-fill", theme.ColorNamePrimary
		switch st.state {
		case stageDone:
			name, c = "check-circle-fill", theme.ColorNameSuccess
		case stageFailed:
			name, c = "x-circle-fill", theme.ColorNameError
		}
		lbl := widget.NewLabel(st.name)
		lbl.Wrapping = fyne.TextWrapWord
		if st.state == stageActive {
			lbl.TextStyle = fyne.TextStyle{Bold: true}
		}
		ic := widget.NewIcon(coloredIcon(name, c))
		r.steps.Add(container.NewBorder(nil, nil, container.NewCenter(ic), nil, lbl))
	}
	r.steps.Refresh()
}

func (r *runView) progress(label string, fraction float64) {
	if label == "" {
		label = "Working…"
	}
	r.barLbl.SetText(label)
	r.bar.SetValue(fraction)
	if !r.barBox.Visible() {
		r.barBox.Show()
		r.logBox.Refresh()
	}
}

func (r *runView) progressDone() {
	if r.barBox.Visible() {
		r.barBox.Hide()
		r.logBox.Refresh()
	}
}

func (r *runView) appendLog(l logLine) {
	r.lines = append(r.lines, l)
	r.logs.Rows = append(r.logs.Rows, logRow(l))
	if over := len(r.lines) - maxLogLines; over > 0 {
		r.lines = r.lines[over:]
		r.logs.Rows = r.logs.Rows[over:]
	}
	r.logs.Refresh()
	r.logs.scrollToBottom()
}

func (r *runView) logText() string {
	var b strings.Builder
	for _, l := range r.lines {
		fmt.Fprintf(&b, "%s %-7s %s\n", l.at.Format(time.TimeOnly), strings.ToUpper(l.level.String()), l.text)
	}
	return b.String()
}
