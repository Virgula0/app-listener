package install

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/Virgula0/app-listener/internal/codeedit"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// UnifiedDiff renders a unified diff between two file contents, labeled
// as "existing" and "new" with three lines of context.
func UnifiedDiff(existing, desired string) string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(existing),
		B:        difflib.SplitLines(desired),
		FromFile: "existing",
		ToFile:   "new",
		Context:  3,
	}
	out, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return fmt.Sprintf("--- existing\n+++ new\n(diff failed: %v)", err)
	}
	return out
}

// ConfirmOverwrite compares an existing file with the desired content: identical is skipped without
// asking; otherwise the diff shows in a full-screen viewer, then the user is asked to overwrite.
// True only if the file was (or will be) replaced.
func ConfirmOverwrite(path string, existing, desired []byte) (bool, error) {
	if bytes.Equal(existing, desired) {
		return false, nil
	}
	if err := ShowDiff(path, UnifiedDiff(string(existing), string(desired))); err != nil {
		return false, err
	}
	return wizard.Confirm(wizard.Question{
		Title:       fmt.Sprintf("Overwrite %s with the new version?", path),
		Description: "The existing file differs from the one this installer ships. Keeping the existing file leaves the previous version in place.",
		Affirmative: "Overwrite",
		Negative:    "Keep existing",
	})
}

// ShowDiff displays a diff in a full-screen read-only viewer. Close it
// with q, Esc or Enter.
func ShowDiff(title, diff string) error {
	if ui := wizard.Active(); ui != nil {
		return ui.ShowDiff(title, diff)
	}
	m := newDiffModel(title, diff)
	p := tea.NewProgram(&m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

type diffModel struct {
	textarea textarea.Model
	hl       *codeedit.Highlighter
	find     codeedit.Finder
	sel      codeedit.Selection
	title    string
	// top is the viewer's first screen row in the last View (mouse coordinates).
	top int
}

func newDiffModel(title, diff string) diffModel {
	ta := codeedit.New()
	ta.SetHeight(24)
	ta.Focus()
	codeedit.SetText(&ta, diff)
	return diffModel{
		textarea: ta,
		hl:       codeedit.NewHighlighter("changes.diff", diff),
		find:     codeedit.NewFinder(false),
		title:    title,
	}
}

func (m *diffModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *diffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.textarea.SetWidth(max(msg.Width-editorSidePadding, editMinWidth))
		m.textarea.SetHeight(max(msg.Height-editorChromeLines, editMinHeight))
	}
	if msg, ok := msg.(tea.MouseMsg); ok {
		m.sel.Mouse(&m.textarea, m.hl, msg, msg.X, msg.Y-m.top)
		return m, nil
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		if m.find.Update(&m.textarea, msg) {
			return m, nil
		}
		if ok, cmd := m.sel.Key(&m.textarea, msg, false); ok {
			return m, cmd
		}
		switch msg.String() {
		case "q", "esc", "enter", "ctrl+q":
			return m, tea.Quit
		}
		if codeedit.Navigate(&m.textarea, msg) {
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *diffModel) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("q / Esc / Enter: continue  ·  Ctrl+F find  ·  Shift+arrows/mouse select  ·  Ctrl+C copy"))
	b.WriteString("\n\n")
	m.top = lipgloss.Height(b.String()) - 1
	b.WriteString(m.find.View(m.hl, &m.textarea, &m.sel))
	return b.String()
}
