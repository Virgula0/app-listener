package update

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	changelogTitleStyle = "212"
	changelogHintStyle  = "240"
	changelogMinWidth   = 40
	changelogMinHeight  = 5
	changelogChrome     = 4
)

// changelogText returns the release notes of r, falling back to a short
// placeholder when GitHub did not publish a body for the release.
func changelogText(r *githubRelease) string {
	if notes := strings.TrimSpace(r.Body); notes != "" {
		return notes
	}
	return fmt.Sprintf("No changelog available for this release.\n\nTag: %s\nPublished: %s", r.TagName, r.PublishedAt)
}

// showChangelog displays release notes in a full-screen read-only TUI
// viewer. Close it with q, Esc or Enter.
func showChangelog(title, notes string) error {
	m := newChangelogModel(title, notes)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type changelogModel struct {
	textarea textarea.Model
	title    string
}

func newChangelogModel(title, notes string) changelogModel {
	ta := textarea.New()
	ta.SetValue(notes)
	ta.SetHeight(24)
	ta.ShowLineNumbers = true
	ta.Focus()
	ta.KeyMap.LineNext.SetKeys("down", "ctrl+n")
	ta.KeyMap.LinePrevious.SetKeys("up", "ctrl+p")
	ta.KeyMap.CharacterBackward.SetKeys("left")
	ta.KeyMap.CharacterForward.SetKeys("right")
	return changelogModel{textarea: ta, title: title}
}

func (m *changelogModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *changelogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.textarea.SetWidth(max(msg.Width-2, changelogMinWidth))
		m.textarea.SetHeight(max(msg.Height-changelogChrome, changelogMinHeight))
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

func (m *changelogModel) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(changelogTitleStyle)).
		Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(changelogHintStyle)).
		Render("q / Esc / Enter: continue"))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	return b.String()
}
