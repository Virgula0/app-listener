package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)
)

type model struct {
	directory string
	recursive bool
	width     int
	height    int
}

func NewModel(directory string, recursive bool) tea.Model {
	return model{
		directory: directory,
		recursive: recursive,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	recursiveStr := "no"
	if m.recursive {
		recursiveStr = "yes"
	}

	header := titleStyle.Render("app-listener — File System Monitor (eBPF)")

	info := infoStyle.Render(fmt.Sprintf(
		"Directory : %s\nRecursive : %s",
		m.directory, recursiveStr,
	))

	status := statusStyle.Render("\n  ⏳ Waiting for file system events...")

	footer := warnStyle.Render("\nPress q or ctrl+c to exit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		info,
		status,
		footer,
	)

	return borderStyle.Render(content)
}
