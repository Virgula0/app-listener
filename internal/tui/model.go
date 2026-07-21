package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

const maxEvents = 500

type eventMsg ebpf.FileEvent

var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	typeStyles = map[ebpf.EventType]lipgloss.Style{
		ebpf.EventOpen:     lipgloss.NewStyle().Foreground(lipgloss.Color("76")),
		ebpf.EventRead:     lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		ebpf.EventWrite:    lipgloss.NewStyle().Foreground(lipgloss.Color("214")),
		ebpf.EventDelete:   lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		ebpf.EventRename:   lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
		ebpf.EventSymlink:  lipgloss.NewStyle().Foreground(lipgloss.Color("45")),
		ebpf.EventHardlink: lipgloss.NewStyle().Foreground(lipgloss.Color("45")),
		ebpf.EventMkdir:    lipgloss.NewStyle().Foreground(lipgloss.Color("184")),
	}

	timeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	commStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
)

type eventLine struct {
	event ebpf.FileEvent
	line  string
}

type model struct {
	events   <-chan ebpf.FileEvent
	lines    []eventLine
	ready    bool
	viewport viewport.Model
	width    int
	height   int

	paths     []string
	recursive bool
	depth     int
	startTime time.Time
	eventID   int
}

func NewModel(events <-chan ebpf.FileEvent, paths []string, recursive bool, depth int) tea.Model {
	return &model{
		events:    events,
		lines:     make([]eventLine, 0, maxEvents),
		paths:     paths,
		recursive: recursive,
		depth:     depth,
		startTime: time.Now(),
	}
}

func (m *model) Init() tea.Cmd {
	return listenForEvents(m.events)
}

func listenForEvents(ch <-chan ebpf.FileEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return eventMsg(ev)
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 5
		footerHeight := 1
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			vp := viewport.New(m.width-4, m.height-headerHeight-footerHeight)
			m.viewport = vp
			m.ready = true
		} else {
			m.viewport.Width = m.width - 4
			m.viewport.Height = m.height - headerHeight - footerHeight
		}

		m.renderViewport()

	case eventMsg:
		ev := ebpf.FileEvent(msg)
		m.addEvent(&ev)
		m.renderViewport()
		m.viewport.GotoBottom()

		return m, listenForEvents(m.events)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) addEvent(ev *ebpf.FileEvent) {
	m.eventID++

	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatType(ev.Type)
	pathPart := ev.Path
	extra := ""

	switch ev.Type {
	case ebpf.EventRead, ebpf.EventWrite:
		extra = fmt.Sprintf(" fd=%d", ev.FD)
	case ebpf.EventRename:
		extra = fmt.Sprintf(" → %s", ev.Dest)
	case ebpf.EventSymlink:
		extra = fmt.Sprintf(" → %s", ev.Dest)
	case ebpf.EventHardlink:
		extra = fmt.Sprintf(" → %s", ev.Dest)
	}

	line := fmt.Sprintf("%s %s %s[%d] %s%s",
		timeStyle.Render(ts),
		t,
		commStyle.Render(ev.Comm),
		ev.PID,
		pathPart,
		extra,
	)

	m.lines = append(m.lines, eventLine{event: *ev, line: line})
	if len(m.lines) > maxEvents {
		m.lines = m.lines[len(m.lines)-maxEvents:]
	}
}

func (m *model) renderViewport() {
	lines := make([]string, len(m.lines))
	for i, el := range m.lines {
		lines[i] = el.line
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	m.viewport.SetContent(content)
}

func formatType(t ebpf.EventType) string {
	if s, ok := typeStyles[t]; ok {
		return s.Render(fmt.Sprintf("%-8s", t.String()))
	}
	return fmt.Sprintf("%-8s", t.String())
}

func (m *model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	recursiveStr := "no"
	if m.recursive {
		recursiveStr = fmt.Sprintf("yes (depth:%d)", m.depth)
	}

	header := headerStyle.Render("app-listener — File System Monitor (eBPF)")

	info := infoStyle.Render(fmt.Sprintf(
		"Watching: %s  |  Recursive: %s  |  Events: %d  |  Uptime: %s",
		strings.Join(m.paths, ", "),
		recursiveStr,
		m.eventID,
		time.Since(m.startTime).Round(time.Second),
	))

	footer := footerStyle.Render("Press q or ctrl+c to exit")

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		info,
		m.viewport.View(),
		"",
		footer,
	))
}
