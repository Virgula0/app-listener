package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/procstats"
)

const maxGuardEvents = 500

type guardEventMsg guard.GuardEvent
type guardStatsMsg struct{}

var (
	blockedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	allowedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76")).
			Bold(true)

	guardModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	guardBinaryStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("141"))

	guardCountStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	guardAllowedCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("76"))
)

type guardEventLine struct {
	event guard.GuardEvent
	line  string
}

type guardModel struct {
	events   <-chan guard.GuardEvent
	lines    []guardEventLine
	ready    bool
	viewport viewport.Model
	width    int
	height   int

	guardPath string
	mode      guard.Mode
	binaries  []guard.BinaryEntry
	startTime time.Time
	eventID   int
	blocked   int
	allowed   int

	lastStats    *procstats.Stats
	prevEventID  int
	eventsPerSec float64
}

func NewGuardModel(events <-chan guard.GuardEvent, guardPath string, mode guard.Mode, binaries []guard.BinaryEntry) tea.Model {
	return &guardModel{
		events:    events,
		lines:     make([]guardEventLine, 0, maxGuardEvents),
		guardPath: guardPath,
		mode:      mode,
		binaries:  binaries,
		startTime: time.Now(),
	}
}

func (m *guardModel) Init() tea.Cmd {
	return tea.Batch(
		listenForGuardEvents(m.events),
		tickGuardStats(),
	)
}

func listenForGuardEvents(ch <-chan guard.GuardEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return guardEventMsg(ev)
	}
}

func tickGuardStats() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return guardStatsMsg{}
	})
}

func (m *guardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 4
		footerHeight := 3
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			vp := viewport.New(m.width-4, m.height-headerHeight-footerHeight-1)
			m.viewport = vp
			m.ready = true
		} else {
			m.viewport.Width = m.width - 4
			m.viewport.Height = m.height - headerHeight - footerHeight - 1
		}

		m.renderGuardViewport()

	case guardEventMsg:
		ev := guard.GuardEvent(msg)
		m.addGuardEvent(&ev)
		m.renderGuardViewport()
		m.viewport.GotoBottom()

		return m, listenForGuardEvents(m.events)

	case guardStatsMsg:
		m.updateGuardStats()
		m.renderGuardViewport()

		return m, tickGuardStats()

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

func (m *guardModel) updateGuardStats() {
	s, err := procstats.Read()
	if err != nil {
		return
	}

	eventsDelta := m.eventID - m.prevEventID
	m.eventsPerSec = float64(eventsDelta)
	m.prevEventID = m.eventID
	m.lastStats = s
}

func (m *guardModel) addGuardEvent(ev *guard.GuardEvent) {
	m.eventID++

	if ev.Blocked {
		m.blocked++
	} else {
		m.allowed++
	}

	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatGuardType(ev.Type)
	decision := formatDecision(ev.Blocked)

	pathPart := ev.Path
	extra := ""

	switch ev.Type {
	case ebpf.EventRead, ebpf.EventWrite:
		extra = fmt.Sprintf(" fd=%d", ev.FD)
	case ebpf.EventRename:
		extra = fmt.Sprintf(" \u2192 %s", ev.Dest)
	case ebpf.EventSymlink:
		extra = fmt.Sprintf(" \u2192 %s", ev.Dest)
	case ebpf.EventHardlink:
		extra = fmt.Sprintf(" \u2192 %s", ev.Dest)
	}

	line := fmt.Sprintf("%s %s %s %s[%d] %s%s",
		timeStyle.Render(ts),
		t,
		decision,
		commStyle.Render(ev.Comm),
		ev.PID,
		pathPart,
		extra,
	)

	m.lines = append(m.lines, guardEventLine{event: *ev, line: line})
	if len(m.lines) > maxGuardEvents {
		m.lines = m.lines[len(m.lines)-maxGuardEvents:]
	}
}

func formatDecision(blocked bool) string {
	if blocked {
		return blockedStyle.Render("BLOCKED ")
	}
	return allowedStyle.Render("ALLOWED ")
}

func formatGuardType(t ebpf.EventType) string {
	if s, ok := typeStyles[t]; ok {
		return s.Render(fmt.Sprintf("%-8s", t.String()))
	}
	return fmt.Sprintf("%-8s", t.String())
}

func (m *guardModel) renderGuardViewport() {
	lines := make([]string, len(m.lines))
	for i, el := range m.lines {
		lines[i] = el.line
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	m.viewport.SetContent(content)
}

func (m *guardModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	modeStr := "blacklist"
	modeLabel := "Blocking"
	if m.mode == guard.ModeWhitelist {
		modeStr = "whitelist"
		modeLabel = "Allowing"
	}

	binaryNames := make([]string, len(m.binaries))
	for i, b := range m.binaries {
		binaryNames[i] = b.Path
	}

	header := headerStyle.Render("app-listener \u2014 Guard (eBPF LSM)")

	info := infoStyle.Render(fmt.Sprintf(
		"Guarding: %s  |  Mode: %s %s  |  Binaries: %s",
		m.guardPath,
		guardModeStyle.Render(modeStr),
		modeLabel,
		guardBinaryStyle.Render(strings.Join(binaryNames, ", ")),
	))

	statsLine := infoStyle.Render(fmt.Sprintf(
		"Events: %d  |  %s %d  |  %s %d  |  Uptime: %s",
		m.eventID,
		guardCountStyle.Render("BLOCKED"),
		m.blocked,
		guardAllowedCountStyle.Render("ALLOWED"),
		m.allowed,
		time.Since(m.startTime).Round(time.Second),
	))

	resLine := m.renderGuardResourceBar()

	footer := footerStyle.Render("Press q or ctrl+c to exit")

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		info,
		statsLine,
		m.viewport.View(),
		resLine,
		"",
		footer,
	))
}

func (m *guardModel) renderGuardResourceBar() string {
	if m.lastStats == nil {
		return ""
	}

	rssMB := float64(m.lastStats.RSS) / 1024 / 1024
	cpuTotal := m.lastStats.CPUUser + m.lastStats.CPUSys
	cpuSec := cpuTotal.Seconds()
	uptimeSec := time.Since(m.startTime).Seconds()
	cpuPct := 0.0
	if uptimeSec > 0 {
		cpuPct = (cpuSec / uptimeSec) * 100
	}

	rss := memStyle.Render(fmt.Sprintf("%.1f MB", rssMB))
	cpu := cpuStyle.Render(fmt.Sprintf("%.1f%%", cpuPct))
	eps := epsStyle.Render(fmt.Sprintf("%.0f/s", m.eventsPerSec))

	return resourceStyle.Render(
		resourceLabelStyle.Render(" Mem:")+" "+rss+
			resourceLabelStyle.Render("  CPU:")+" "+cpu+
			resourceLabelStyle.Render("  Evt/s:")+" "+eps,
	)
}
