package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noticeTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	noticeIntroStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	noticeButtonStyle = lipgloss.NewStyle().Bold(true).Padding(0, 3).
				Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62"))
	noticeHintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

const (
	noticeMargin = 2
	// noticeChrome: title, blank, intro gap, blank, button, hint, around the scrolled list.
	noticeChrome = 8
)

// ShowNotice shows title, intro and a bulleted list full-screen until acknowledged (Enter, Esc,
// q). Text is shown verbatim: paths are full of '_' and '*', which markdown renderers eat.
func ShowNotice(title, intro string, items []string) error {
	_, err := tea.NewProgram(&noticeModel{title: title, intro: intro, items: items}, tea.WithAltScreen()).Run()
	return err
}

type noticeModel struct {
	title, intro string
	items        []string
	vp           viewport.Model
	width        int
	ready        bool
}

func (m *noticeModel) Init() tea.Cmd { return nil }

func (m *noticeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = max(msg.Width-2*noticeMargin, 20)
		introRows := lipgloss.Height(m.renderIntro())
		list := renderBullets(m.items, m.width)
		// Fit the list, scrolling only when the screen is shorter: OK stays right under it.
		m.vp = viewport.New(m.width, max(min(lipgloss.Height(list), msg.Height-noticeChrome-introRows), 3))
		m.vp.SetContent(list)
		m.ready = true
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", "q", "ctrl+c", " ":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *noticeModel) renderIntro() string {
	return noticeIntroStyle.Width(m.width).Render(m.intro)
}

func (m *noticeModel) View() string {
	if !m.ready {
		return ""
	}
	hint := "Enter / Esc: OK"
	if m.vp.TotalLineCount() > m.vp.Height {
		hint = "↑/↓ PgUp/PgDn: scroll  ·  " + hint
	}
	body := lipgloss.JoinVertical(lipgloss.Left,
		noticeTitleStyle.Width(m.width).Render(m.title),
		"",
		m.renderIntro(),
		"",
		m.vp.View(),
		"",
		noticeButtonStyle.Render("OK"),
		noticeHintStyle.Render(hint),
	)
	return lipgloss.NewStyle().Margin(1, noticeMargin).Render(body)
}

// renderBullets wraps each item to width with a hanging indent under its bullet.
func renderBullets(items []string, width int) string {
	const bullet, hang = "• ", "  "
	wrapper := lipgloss.NewStyle().Width(max(width-len(hang), 10))
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString("\n")
		}
		for j, line := range strings.Split(wrapper.Render(item), "\n") {
			if j == 0 {
				b.WriteString(bullet)
			} else {
				b.WriteString("\n" + hang)
			}
			b.WriteString(strings.TrimRight(line, " "))
		}
	}
	return b.String()
}
