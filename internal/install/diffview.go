package install

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pmezard/go-difflib/difflib"
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

// ConfirmOverwrite compares an existing file with the desired content.
// Identical files are skipped without asking. When they differ, the diff
// is shown in a full-screen TUI viewer, then the user is asked whether to
// overwrite. It returns true only when the existing file was (or will be)
// replaced.
func ConfirmOverwrite(path string, existing, desired []byte) (bool, error) {
	if bytes.Equal(existing, desired) {
		return false, nil
	}
	if err := ShowDiff(path, UnifiedDiff(string(existing), string(desired))); err != nil {
		return false, err
	}
	overwrite := false
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Overwrite %s with the new version?", path)).
			Description("The existing file differs from the one this installer ships. Keeping the existing file leaves the previous version in place.").
			Affirmative("Overwrite").
			Negative("Keep existing").
			Value(&overwrite),
	)).Run(); err != nil {
		return false, err
	}
	return overwrite, nil
}

// ShowDiff displays a diff in a full-screen read-only viewer. Close it
// with q, Esc or Enter.
func ShowDiff(title, diff string) error {
	m := newDiffModel(title, diff)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type diffModel struct {
	textarea textarea.Model
	title    string
}

func newDiffModel(title, diff string) diffModel {
	ta := textarea.New()
	ta.SetValue(diff)
	ta.SetHeight(24)
	ta.ShowLineNumbers = true
	ta.Focus()
	ta.KeyMap.LineNext.SetKeys("down", "ctrl+n")
	ta.KeyMap.LinePrevious.SetKeys("up", "ctrl+p")
	ta.KeyMap.CharacterBackward.SetKeys("left")
	ta.KeyMap.CharacterForward.SetKeys("right")
	return diffModel{textarea: ta, title: title}
}

func (m *diffModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *diffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.textarea.SetWidth(max(msg.Width-editorSidePadding, editMinWidth))
		m.textarea.SetHeight(max(msg.Height-editorChromeLines, editMinHeight))
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "esc", "enter", "ctrl+q":
			return m, tea.Quit
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
		Render("q / Esc / Enter: continue"))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	return b.String()
}
