package update

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/codeedit"
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
	hl       *codeedit.Highlighter
	find     codeedit.Finder
	title    string
}

func newChangelogModel(title, notes string) changelogModel {
	ta := codeedit.New()
	ta.SetHeight(24)
	ta.Focus()
	codeedit.SetText(&ta, notes)
	return changelogModel{
		textarea: ta,
		hl:       codeedit.NewHighlighter("CHANGELOG.md", notes),
		find:     codeedit.NewFinder(false),
		title:    title,
	}
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
		if m.find.Update(&m.textarea, msg) {
			return m, nil
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

func (m *changelogModel) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(changelogTitleStyle)).
		Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(changelogHintStyle)).
		Render("q / Esc / Enter: continue  ·  Ctrl+F find"))
	b.WriteString("\n\n")
	b.WriteString(m.find.View(m.hl, &m.textarea))
	return b.String()
}
