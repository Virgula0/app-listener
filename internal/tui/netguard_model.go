package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/procstats"
)

const maxNetGuardEvents = 500

type netGuardEventMsg networkguard.NetGuardEvent
type netGuardStatsMsg struct{}

var (
	netGuardBlockedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)

	netGuardAllowedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("76")).
				Bold(true)

	netGuardModeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214")).
				Bold(true)

	netGuardBlockedCountStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("196"))

	netGuardAllowedCountStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("76"))
)

type netGuardEventLine struct {
	event networkguard.NetGuardEvent
	line  string
}

type netGuardModel struct {
	events   <-chan networkguard.NetGuardEvent
	lines    []netGuardEventLine
	ready    bool
	viewport viewport.Model
	width    int
	height   int

	mode      networkguard.Mode
	binaries  []networkguard.BinaryEntry
	startTime time.Time
	eventID   int
	blocked   int
	allowed   int

	lastStats    *procstats.Stats
	prevEventID  int
	eventsPerSec float64
}

func NewNetGuardModel(events <-chan networkguard.NetGuardEvent, mode networkguard.Mode, binaries []networkguard.BinaryEntry) tea.Model {
	return &netGuardModel{
		events:    events,
		lines:     make([]netGuardEventLine, 0, maxNetGuardEvents),
		mode:      mode,
		binaries:  binaries,
		startTime: time.Now(),
	}
}

func (m *netGuardModel) Init() tea.Cmd {
	return tea.Batch(
		m.waitForEvent(),
		m.waitForStats(),
	)
}

func (m *netGuardModel) waitForEvent() tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-m.events
		if !ok {
			return nil
		}
		return netGuardEventMsg(ev)
	}
}

func (m *netGuardModel) waitForStats() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return netGuardStatsMsg{}
	})
}

func (m *netGuardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}

		m.renderContent()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case quitKey, "q":
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case netGuardEventMsg:
		ev := networkguard.NetGuardEvent(msg)
		m.eventID++
		if ev.Blocked {
			m.blocked++
		} else {
			m.allowed++
		}

		m.lines = append(m.lines, netGuardEventLine{
			event: ev,
			line:  m.formatEvent(&ev),
		})

		if len(m.lines) > maxNetGuardEvents {
			m.lines = m.lines[len(m.lines)-maxNetGuardEvents:]
		}

		m.renderContent()
		m.viewport.GotoBottom()

		cmd := m.waitForEvent()
		return m, cmd

	case netGuardStatsMsg:
		s, err := procstats.Read()
		if err == nil {
			eventsDelta := m.eventID - m.prevEventID
			m.eventsPerSec = float64(eventsDelta) / 3.0
			m.prevEventID = m.eventID
			m.lastStats = s
		}

		cmd := m.waitForStats()
		return m, cmd
	}

	return m, nil
}

func (m *netGuardModel) formatEvent(ev *networkguard.NetGuardEvent) string {
	style := netTypeStyles[ev.Type]
	eventName := style.Render(ev.Type.String())

	status := ""
	if ev.Blocked {
		status = netGuardBlockedStyle.Render(" BLOCKED")
	} else {
		status = netGuardAllowedStyle.Render(" ALLOWED")
	}

	comm := ev.Comm
	protocol := protoString(ev.Protocol)

	addr := ev.DstAddr
	if addr == "" {
		addr = ev.SrcAddr
	}

	details := fmt.Sprintf("pid=%d tid=%d", ev.PID, ev.TID)

	if ev.Size > 0 {
		details += fmt.Sprintf(" size=%d", ev.Size)
	}

	return fmt.Sprintf("[%d] %s %s %s %s %s%s %s",
		m.eventID,
		time.Now().Format("15:04:05.000"),
		eventName,
		comm,
		protocol,
		addr,
		status,
		details,
	)
}

func protoString(proto uint32) string {
	switch proto {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	case 1:
		return "ICMP"
	case 58:
		return "ICMPv6"
	default:
		return "UNKN"
	}
}

func (m *netGuardModel) renderContent() {
	var b strings.Builder

	for i := range m.lines {
		b.WriteString(m.lines[i].line)
		b.WriteString("\n")
	}

	m.viewport.SetContent(b.String())
}

func (m *netGuardModel) View() string {
	if !m.ready {
		return initializingView
	}

	var b strings.Builder

	title := " Network Guard"

	fmt.Fprintf(&b, "%s\n", lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Bold(true).
		Render(title))

	modeLabel := "BLACKLIST"
	if m.mode == networkguard.ModeWhitelist {
		modeLabel = "WHITELIST"
	}
	fmt.Fprintf(&b, " Mode: %s  Binaries: %s\n", netGuardModeStyle.Render(modeLabel), m.listSummary(m.binaries))
	fmt.Fprint(&b, " "+m.renderResourceBar()+"\n\n")
	fmt.Fprint(&b, m.viewport.View())

	return b.String()
}

func (m *netGuardModel) listSummary(entries []networkguard.BinaryEntry) string {
	parts := make([]string, len(entries))
	for i, b := range entries {
		parts[i] = b.Path
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}

func (m *netGuardModel) renderResourceBar() string {
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
	blocked := netGuardBlockedCountStyle.Render(fmt.Sprintf("%d blocked", m.blocked))
	allowed := netGuardAllowedCountStyle.Render(fmt.Sprintf("%d allowed", m.allowed))

	return resourceStyle.Render(
		resourceLabelStyle.Render(" Mem:") + " " + rss +
			resourceLabelStyle.Render("  CPU:") + " " + cpu +
			resourceLabelStyle.Render("  Evt/s:") + " " + eps +
			"  " + blocked +
			"  " + allowed,
	)
}

const (
	headerHeight = 5
	footerHeight = 0
)
