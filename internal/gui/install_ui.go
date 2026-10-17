//go:build gui

package gui

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/wizard"
)

// installUI is the wizard.UI of one run: every prompt is a dialog over the window, asked from the
// run's goroutine; stages and progress go to the run view.
type installUI struct {
	s   *shell
	run *runView
}

var _ wizard.UI = (*installUI)(nil)

// ask shows a dialog on the UI goroutine and blocks the caller until reply is called (only the
// first call counts: a dialog's buttons and its Escape dismissal can both fire).
func ask[T any](show func(reply func(T))) T {
	ch := make(chan T, 1)
	var once sync.Once
	fyne.Do(func() { show(func(v T) { once.Do(func() { ch <- v }) }) })
	return <-ch
}

// promptBody is a prompt's bold question over its muted explanation.
func promptBody(title, description string, rest ...fyne.CanvasObject) *fyne.Container {
	q := boldLabel(title)
	q.Wrapping = fyne.TextWrapWord
	box := container.NewVBox(q)
	if description != "" {
		box.Add(mutedText(description))
	}
	for _, o := range rest {
		box.Add(o)
	}
	return box
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func (u *installUI) Confirm(q wizard.Question) (bool, error) {
	return ask(func(reply func(bool)) {
		d := dialog.NewCustomConfirm(u.run.name, orDefault(q.Affirmative, "Yes"), orDefault(q.Negative, "No"),
			scrollable(promptBody(q.Title, q.Description)), reply, u.s.win)
		u.s.showSized(d, fyne.NewSize(640, 280))
	}), nil
}

type multiAnswer struct {
	picked []int
	ok     bool
}

func (u *installUI) MultiSelect(title, description string, choices []wizard.Choice) ([]int, error) {
	a := ask(func(reply func(multiAnswer)) { u.showMultiSelect(title, description, choices, reply) })
	if !a.ok {
		return nil, wizard.ErrCanceled
	}
	return a.picked, nil
}

func (u *installUI) showMultiSelect(title, description string, choices []wizard.Choice, reply func(multiAnswer)) {
	checks := make([]*widget.Check, len(choices))
	list := container.NewVBox()
	count := mutedText("")
	recount := func() {
		n := 0
		for _, c := range checks {
			if c.Checked {
				n++
			}
		}
		count.SetText(fmt.Sprintf("%d of %d selected", n, len(checks)))
	}
	for i, c := range choices {
		checks[i] = widget.NewCheck(c.Label, func(bool) { recount() })
		checks[i].Checked = c.Selected
		list.Add(checks[i])
	}
	recount()
	setAll := func(on bool) {
		for _, c := range checks {
			c.SetChecked(on)
		}
	}
	all := newTextButton("Select all", func() { setAll(true) })
	none := newTextButton("Select none", func() { setAll(false) })
	top := promptBody(title, description, hfill(count, all, none))
	body := container.NewBorder(top, nil, nil, nil, container.NewVScroll(list))
	d := dialog.NewCustomConfirm(u.run.name, "Continue", "Cancel", body, func(ok bool) {
		var picked []int
		for i, c := range checks {
			if c.Checked {
				picked = append(picked, i)
			}
		}
		reply(multiAnswer{picked: picked, ok: ok})
	}, u.s.win)
	u.s.showSized(d, fyne.NewSize(760, 600))
}

type textAnswer struct {
	text string
	ok   bool
}

func (u *installUI) Input(f wizard.Field) (string, error) {
	a := ask(func(reply func(textAnswer)) {
		entry := widget.NewEntry()
		if f.Secret {
			entry = widget.NewPasswordEntry()
		}
		entry.SetPlaceHolder(f.Placeholder)
		field := fyne.CanvasObject(entry)
		if f.Path {
			field = container.NewBorder(nil, nil, nil, u.pathButtons(entry), entry)
		}
		d := dialog.NewCustomConfirm(u.run.name, "OK", "Cancel", scrollable(promptBody(f.Title, f.Description, field)),
			func(ok bool) { reply(textAnswer{text: entry.Text, ok: ok}) }, u.s.win)
		entry.OnSubmitted = func(string) { d.Confirm() }
		u.s.showSized(d, fyne.NewSize(680, 260))
		u.s.win.Canvas().Focus(entry)
	})
	if !a.ok {
		return "", wizard.ErrCanceled
	}
	return a.text, nil
}

// pathButtons fill entry from a folder or file picker opened in the invoking user's home.
func (u *installUI) pathButtons(entry *widget.Entry) fyne.CanvasObject {
	folder := newButton("Folder…", icon("folder"), func() {
		u.s.pickPath(pickFolder, "Choose a folder", "", entry.SetText)
	})
	file := newButton("File…", icon("file"), func() {
		u.s.pickPath(pickFile, "Choose a file", "", entry.SetText)
	})
	return container.NewHBox(folder, file)
}

func (u *installUI) EditText(title, name, text string) (string, error) {
	a := ask(func(reply func(textAnswer)) {
		ed := widget.NewMultiLineEntry()
		ed.TextStyle = fyne.TextStyle{Monospace: true}
		ed.SetText(text)
		top := promptBody(title, "Editing "+name+". Save to continue the installer with this text.")
		d := dialog.NewCustomConfirm(u.run.name, "Save", "Cancel", container.NewBorder(top, nil, nil, nil, ed),
			func(ok bool) { reply(textAnswer{text: ed.Text, ok: ok}) }, u.s.win)
		u.s.showSized(d, fyne.NewSize(1040, 720))
	})
	if !a.ok {
		return "", wizard.ErrCanceled
	}
	return a.text, nil
}

func (u *installUI) ShowDiff(title, diff string) error {
	ask(func(reply func(struct{})) {
		d := dialog.NewCustom(u.run.name, "Close", container.NewBorder(promptBody(title, ""), nil, nil, nil,
			container.NewScroll(diffText(diff))), u.s.win)
		d.SetOnClosed(func() { reply(struct{}{}) })
		u.s.showSized(d, fyne.NewSize(1040, 720))
	})
	return nil
}

// diffText colors a unified diff's added and removed lines.
func diffText(diff string) *widget.RichText {
	lines := strings.Split(strings.TrimRight(diff, "\n"), "\n")
	segs := make([]widget.RichTextSegment, 0, len(lines))
	for _, l := range lines {
		style := widget.RichTextStyle{TextStyle: fyne.TextStyle{Monospace: true}}
		switch {
		case strings.HasPrefix(l, "+++"), strings.HasPrefix(l, "---"), strings.HasPrefix(l, "@@"):
			style.ColorName = colorInfo
		case strings.HasPrefix(l, "+"):
			style.ColorName = theme.ColorNameSuccess
		case strings.HasPrefix(l, "-"):
			style.ColorName = theme.ColorNameError
		}
		segs = append(segs, &widget.TextSegment{Text: l, Style: style})
	}
	return widget.NewRichText(segs...)
}

func (u *installUI) Progress(run func(set func(label string, fraction float64)) error) error {
	fyne.Do(func() { u.run.progress("", 0) })
	defer fyne.Do(u.run.progressDone)
	return run(func(label string, fraction float64) {
		fyne.Do(func() { u.run.progress(label, fraction) })
	})
}

func (u *installUI) Stage(name string) {
	fyne.Do(func() { u.run.stage(name) })
}

// logHook mirrors logrus entries into the run view while a run is active.
type logHook struct{ run *runView }

func (h *logHook) Levels() []log.Level { return log.AllLevels }

func (h *logHook) Fire(e *log.Entry) error {
	text := e.Message
	if len(e.Data) > 0 {
		keys := make([]string, 0, len(e.Data))
		for k := range e.Data {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			text += fmt.Sprintf(" %s=%v", k, e.Data[k])
		}
	}
	line := logLine{level: e.Level, at: e.Time, text: text}
	fyne.Do(func() { h.run.appendLog(line) })
	return nil
}

// captureLogs adds a logHook to the standard logger until restore is called.
func captureLogs(run *runView) (restore func()) {
	std := log.StandardLogger()
	prev := std.ReplaceHooks(make(log.LevelHooks))
	next := make(log.LevelHooks)
	for l, hs := range prev {
		next[l] = slices.Clone(hs)
	}
	next.Add(&logHook{run: run})
	std.ReplaceHooks(next)
	return func() { std.ReplaceHooks(prev) }
}
