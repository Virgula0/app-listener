package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkmonitor"
	"github.com/Virgula0/app-listener/internal/procstats"
)

const maxNetEvents = 500

type netEventMsg ebpf.NetEvent
type netStatsMsg struct{}

var (
	netConnectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	netAcceptStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	netSendStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	netRecvStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	netCloseStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	netDNSStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	netBindStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	netListenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("44"))

	netTypeStyles = map[ebpf.NetEventType]lipgloss.Style{
		ebpf.NetConnect: netConnectStyle,
		ebpf.NetAccept:  netAcceptStyle,
		ebpf.NetSend:    netSendStyle,
		ebpf.NetRecv:    netRecvStyle,
		ebpf.NetClose:   netCloseStyle,
		ebpf.NetDNS:     netDNSStyle,
		ebpf.NetBind:    netBindStyle,
		ebpf.NetListen:  netListenStyle,
	}

	protoTCP  = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("TCP")
	protoUDP  = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("UDP")
	protoICMP = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("ICMP")
)

type netEventLine struct {
	event ebpf.NetEvent
	line  string
}

type netModel struct {
	events   <-chan ebpf.NetEvent
	lines    []netEventLine
	ready    bool
	viewport viewport.Model
	width    int
	height   int

	binaries  []networkmonitor.BinaryEntry
	startTime time.Time
	eventID   int

	lastStats    *procstats.Stats
	prevEventID  int
	eventsPerSec float64
}

func NewNetModel(events <-chan ebpf.NetEvent, binaries []networkmonitor.BinaryEntry) tea.Model {
	return &netModel{
		events:    events,
		lines:     make([]netEventLine, 0, maxNetEvents),
		binaries:  binaries,
		startTime: time.Now(),
	}
}

func (m *netModel) Init() tea.Cmd {
	return tea.Batch(
		listenForNetEvents(m.events),
		tickNetStats(),
	)
}

func listenForNetEvents(ch <-chan ebpf.NetEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return netEventMsg(ev)
	}
}

func tickNetStats() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return netStatsMsg{}
	})
}

//nolint:dupl // bubbletea update-loop skeleton shared by the three views
func (m *netModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		syncViewport(&m.viewport, &m.ready, m.width, m.height, 4, m.renderNetViewport)

	case netEventMsg:
		ev := ebpf.NetEvent(msg)
		m.addNetEvent(&ev)
		m.renderNetViewport()
		m.viewport.GotoBottom()

		return m, listenForNetEvents(m.events)

	case netStatsMsg:
		m.updateNetStats()
		m.renderNetViewport()

		return m, tickNetStats()

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

func (m *netModel) updateNetStats() {
	s, err := procstats.Read()
	if err != nil {
		return
	}

	eventsDelta := m.eventID - m.prevEventID
	m.eventsPerSec = float64(eventsDelta)
	m.prevEventID = m.eventID
	m.lastStats = s
}

func (m *netModel) addNetEvent(ev *ebpf.NetEvent) {
	m.eventID++

	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatNetType(ev.Type)
	proto := formatNetProto(ev.Protocol)

	src := ev.SrcAddr
	if src == "" {
		src = "-"
	}
	dst := ev.DstAddr
	if dst == "" {
		dst = "-"
	}
	sizeStr := ""
	if ev.Size > 0 {
		sizeStr = fmt.Sprintf(" %d", ev.Size)
	}

	line := fmt.Sprintf("%s %s %s %s %s %s[%d]%s",
		timeStyle.Render(ts),
		t,
		proto,
		commStyle.Render(ev.Comm),
		formatAddr(src),
		formatAddr(dst),
		ev.PID,
		sizeStr,
	)

	m.lines = append(m.lines, netEventLine{event: *ev, line: line})
	if len(m.lines) > maxNetEvents {
		m.lines = m.lines[len(m.lines)-maxNetEvents:]
	}
}

func formatNetType(t ebpf.NetEventType) string {
	if s, ok := netTypeStyles[t]; ok {
		return s.Render(fmt.Sprintf("%-8s", t.String()))
	}
	return fmt.Sprintf("%-8s", t.String())
}

func formatNetProto(proto uint32) string {
	switch proto {
	case 6:
		return protoTCP
	case 17:
		return protoUDP
	case 1:
		return protoICMP
	case 58:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("ICMPv6")
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render("?")
	}
}

func formatAddr(addr string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(addr)
}

func (m *netModel) renderNetViewport() {
	lines := make([]string, len(m.lines))
	for i := range m.lines {
		lines[i] = m.lines[i].line
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	m.viewport.SetContent(content)
}

func (m *netModel) View() string {
	if !m.ready {
		return initializingView
	}

	binaryNames := make([]string, len(m.binaries))
	for i, b := range m.binaries {
		binaryNames[i] = b.Path
	}

	header := headerStyle.Render("app-listener \u2014 Network Monitor (eBPF)")

	info := infoStyle.Render(fmt.Sprintf(
		"Watching: %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Render(strings.Join(binaryNames, ", ")),
	))

	statsLine := infoStyle.Render(fmt.Sprintf(
		"Events: %d  |  Uptime: %s",
		m.eventID,
		time.Since(m.startTime).Round(time.Second),
	))

	resLine := m.renderNetResourceBar()

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

func (m *netModel) renderNetResourceBar() string {
	return formatResourceBar(m.lastStats, m.startTime, m.eventsPerSec)
}
