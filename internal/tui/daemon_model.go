package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/procstats"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const maxDaemonEvents = 500

type daemonEventMsg usecase.DaemonEvent
type daemonStatsMsg struct{}

var (
	daemonPathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141"))

	daemonStateStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214")).
				Bold(true)
)

// DaemonResourceInfo is the presentation-level view of one watched
// resource, built by the daemon command from the configuration.
type DaemonResourceInfo struct {
	Path           string
	NeedEncryption bool
	Binaries       int
}

type daemonEventLine struct {
	event usecase.DaemonEvent
	line  string
}

type daemonModel struct {
	events   <-chan usecase.DaemonEvent
	lines    []daemonEventLine
	ready    bool
	viewport viewport.Model
	width    int
	height   int

	resources []DaemonResourceInfo
	startTime time.Time
	eventID   int
	blocked   int
	allowed   int

	lastStats    *procstats.Stats
	prevEventID  int
	eventsPerSec float64
}

func NewDaemonModel(events <-chan usecase.DaemonEvent, resources []DaemonResourceInfo) tea.Model {
	return &daemonModel{
		events:    events,
		lines:     make([]daemonEventLine, 0, maxDaemonEvents),
		resources: resources,
		startTime: time.Now(),
	}
}

func (m *daemonModel) Init() tea.Cmd {
	return tea.Batch(
		listenForDaemonEvents(m.events),
		tickDaemonStats(),
	)
}

func tickDaemonStats() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return daemonStatsMsg{}
	})
}

func listenForDaemonEvents(ch <-chan usecase.DaemonEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return daemonEventMsg(ev)
	}
}

func (m *daemonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// One rendered row per watched resource sits between the header
		// block and the stats line; reserve it in the viewport budget.
		syncViewport(&m.viewport, &m.ready, m.width, m.height,
			daemonViewFixedRows+len(m.resources), m.renderDaemonViewport)

	case daemonEventMsg:
		ev := usecase.DaemonEvent(msg)
		m.addDaemonEvent(&ev)
		m.renderDaemonViewport()
		m.viewport.GotoBottom()

		return m, listenForDaemonEvents(m.events)

	case daemonStatsMsg:
		m.updateDaemonStats()
		m.renderDaemonViewport()

		return m, tickDaemonStats()

	case tea.KeyMsg:
		switch msg.String() {
		case quitKey, "q":
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *daemonModel) updateDaemonStats() {
	s, err := procstats.Read()
	if err != nil {
		return
	}

	eventsDelta := m.eventID - m.prevEventID
	m.eventsPerSec = float64(eventsDelta)
	m.prevEventID = m.eventID
	m.lastStats = s
}

func (m *daemonModel) addDaemonEvent(ev *usecase.DaemonEvent) {
	m.eventID++

	if ev.Event.Blocked {
		m.blocked++
	} else {
		m.allowed++
	}

	line := fmt.Sprintf("%s %s",
		daemonPathStyle.Render(sanitizeTerminalText(ev.Resource)),
		formatGuardEventLine(&ev.Event))

	m.lines = append(m.lines, daemonEventLine{event: *ev, line: line})
	if len(m.lines) > maxDaemonEvents {
		m.lines = m.lines[len(m.lines)-maxDaemonEvents:]
	}
}

func (m *daemonModel) renderDaemonViewport() {
	lines := make([]string, len(m.lines))
	for i := range m.lines {
		lines[i] = m.lines[i].line
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	m.viewport.SetContent(content)
}

func (m *daemonModel) View() string {
	if !m.ready {
		return initializingView
	}

	header := headerStyle.Render("app-listener \u2014 Daemon (eBPF LSM + fscrypt)")

	var resourceLines []string
	for _, r := range m.resources {
		state := "encrypted"
		if !r.NeedEncryption {
			state = "plain"
		}
		resourceLines = append(resourceLines, infoStyle.Render(fmt.Sprintf(
			"  %s  |  %s  |  whitelist: %d",
			daemonPathStyle.Render(sanitizeTerminalText(r.Path)),
			daemonStateStyle.Render(state),
			r.Binaries,
		)))
	}

	statsLine := infoStyle.Render(fmt.Sprintf(
		"Events: %d  |  %s %d  |  %s %d  |  Uptime: %s",
		m.eventID,
		guardCountStyle.Render("BLOCKED"),
		m.blocked,
		guardAllowedCountStyle.Render("ALLOWED"),
		m.allowed,
		time.Since(m.startTime).Round(time.Second),
	))

	resLine := formatResourceBar(m.lastStats, m.startTime, m.eventsPerSec)

	footer := footerStyle.Render("Press q or ctrl+c to exit")

	return clipToWidth(appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		infoStyle.Render("Watched resources:"),
		lipgloss.JoinVertical(lipgloss.Left, resourceLines...),
		statsLine,
		m.viewport.View(),
		resLine,
		"",
		footer,
	)), m.width)
}
