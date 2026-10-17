package wizard

import (
	"errors"
	"sync"

	"github.com/charmbracelet/huh"
)

// Question is a yes/no prompt.
type Question struct {
	Title, Description    string
	Affirmative, Negative string
	Default               bool
}

// Choice is one option of a multi-select.
type Choice struct {
	Label    string
	Selected bool
}

// Field is a one-line text prompt.
type Field struct {
	Title, Description, Placeholder string
	Secret                          bool
	// Path offers a file and folder picker beside the entry (non-terminal UIs).
	Path bool
}

// UI is a front end for the interactive commands' prompts other than the terminal (the desktop
// GUI). Every method blocks until the user answers; dismissing a prompt returns ErrCanceled.
type UI interface {
	Confirm(q Question) (bool, error)
	MultiSelect(title, description string, choices []Choice) ([]int, error)
	Input(f Field) (string, error)
	EditText(title, name, text string) (string, error)
	ShowDiff(title, diff string) error
	// ShowText shows read-only plain text (a changelog) until the user closes it.
	ShowText(title, text string) error
	// Progress runs run with set reporting its progress (fraction in [0, 1]).
	Progress(run func(set func(label string, fraction float64)) error) error
	// Stage names the step the command is at.
	Stage(name string)
}

// ErrCanceled is returned when the user dismisses a prompt of a non-terminal UI.
var ErrCanceled = errors.New("canceled by the user")

var (
	uiMu   sync.Mutex
	active UI
)

// Use routes every prompt to ui until restore is called. One UI at a time: the commands it drives
// keep package state.
func Use(ui UI) (restore func()) {
	uiMu.Lock()
	prev := active
	active = ui
	uiMu.Unlock()
	return func() {
		uiMu.Lock()
		active = prev
		uiMu.Unlock()
	}
}

// Active returns the UI set by Use, nil for the terminal.
func Active() UI {
	uiMu.Lock()
	defer uiMu.Unlock()
	return active
}

// Stage names the step a command is at; the terminal shows nothing (the logs say it).
func Stage(name string) {
	if ui := Active(); ui != nil {
		ui.Stage(name)
	}
}

// Confirm asks a yes/no question.
func Confirm(q Question) (bool, error) {
	if ui := Active(); ui != nil {
		return ui.Confirm(q)
	}
	answer := q.Default
	c := huh.NewConfirm().Title(q.Title).Affirmative(q.Affirmative).Negative(q.Negative).Value(&answer)
	if q.Description != "" {
		c = c.Description(q.Description)
	}
	if err := huh.NewForm(huh.NewGroup(c)).Run(); err != nil {
		return false, err
	}
	return answer, nil
}

// MultiSelect asks which choices to keep and returns their indexes; height is the terminal list's.
func MultiSelect(title, description string, choices []Choice, height int) ([]int, error) {
	if ui := Active(); ui != nil {
		return ui.MultiSelect(title, description, choices)
	}
	opts := make([]huh.Option[int], 0, len(choices))
	for i, c := range choices {
		opts = append(opts, huh.NewOption(c.Label, i).Selected(c.Selected))
	}
	var picked []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title(title).
			Description(description).
			Options(opts...).
			Height(height).
			Value(&picked),
	))
	if err := form.WithKeyMap(MultiSelectKeymap()).Run(); err != nil {
		return nil, err
	}
	return picked, nil
}

// Input asks for one line of text.
func Input(f Field) (string, error) {
	if ui := Active(); ui != nil {
		return ui.Input(f)
	}
	var v string
	in := huh.NewInput().Title(f.Title).Value(&v)
	if f.Description != "" {
		in = in.Description(f.Description)
	}
	if f.Placeholder != "" {
		in = in.Prompt("> ").Placeholder(f.Placeholder)
	}
	if f.Secret {
		in = in.EchoMode(huh.EchoModePassword)
	}
	if err := huh.NewForm(huh.NewGroup(in)).Run(); err != nil {
		return "", err
	}
	return v, nil
}
