package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/procstats"
)

const maxEvents = 500

// quitKey is the keybinding that exits the TUI views.
const quitKey = "ctrl+c"

// initializingView is shown while the TUI is not ready yet.
const initializingView = "\n  Initializing..."

// Rows each view spends outside the viewport: header/info/footer blocks,
// blank separators and the two appStyle vertical margins. The viewport must
// be sized as height minus these rows, otherwise the bottom of the view
// (resource bar, footer) overflows the reported terminal height.
const (
	fileViewFixedRows    = 8 // header, blank, info | resource, blank, footer | margins
	networkViewFixedRows = 9 // header, blank, info, stats | resource, blank, footer | margins
	guardViewFixedRows   = 9 // header, blank, info, stats | resource, blank, footer | margins
	daemonViewFixedRows  = 9 // header, blank, resources label, stats | resource, blank, footer | margins
)

type eventMsg ebpf.FileEvent
type statsMsg struct{}

var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	resourceStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1)

	resourceLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))

	cpuStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	memStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("76"))

	epsStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("141"))

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
		ebpf.EventMmap:     lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
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

	lastStats    *procstats.Stats
	prevEventID  int
	eventsPerSec float64
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
	return tea.Batch(
		listenForEvents(m.events),
		tickStats(),
	)
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

func tickStats() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return statsMsg{}
	})
}

//nolint:dupl // bubbletea update-loop skeleton shared by the three views
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		syncViewport(&m.viewport, &m.ready, m.width, m.height, fileViewFixedRows, m.renderViewport)

	case eventMsg:
		ev := ebpf.FileEvent(msg)
		m.addEvent(&ev)
		m.renderViewport()
		m.viewport.GotoBottom()

		return m, listenForEvents(m.events)

	case statsMsg:
		m.updateStats()
		m.renderViewport()

		return m, tickStats()

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

func (m *model) updateStats() {
	s, err := procstats.Read()
	if err != nil {
		return
	}

	eventsDelta := m.eventID - m.prevEventID
	m.eventsPerSec = float64(eventsDelta)
	m.prevEventID = m.eventID
	m.lastStats = s
}

func (m *model) addEvent(ev *ebpf.FileEvent) {
	m.eventID++

	ts := time.Unix(0, ev.Timestamp).Format("15:04:05.000")
	t := formatType(ev.Type)
	pathPart := sanitizeTerminalText(ev.Path)
	extra := ""

	switch ev.Type {
	case ebpf.EventRead, ebpf.EventWrite:
		extra = fmt.Sprintf(" fd=%d", ev.FD)
	case ebpf.EventRename:
		extra = fmt.Sprintf(" → %s", sanitizeTerminalText(ev.Dest))
	case ebpf.EventSymlink:
		extra = fmt.Sprintf(" → %s", sanitizeTerminalText(ev.Dest))
	case ebpf.EventHardlink:
		extra = fmt.Sprintf(" → %s", sanitizeTerminalText(ev.Dest))
	}

	line := fmt.Sprintf("%s %s %s[%d] %s%s",
		timeStyle.Render(ts),
		t,
		commStyle.Render(sanitizeTerminalText(ev.Comm)),
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
		return initializingView
	}

	recursiveStr := "no"
	if m.recursive {
		recursiveStr = fmt.Sprintf("yes (depth:%d)", m.depth)
	}

	header := headerStyle.Render("app-listener — File System Monitor (eBPF)")

	info := infoStyle.Render(fmt.Sprintf(
		"Watching: %s  |  Recursive: %s  |  Events: %d  |  Uptime: %s",
		strings.Join(sanitizeTerminalTexts(m.paths), ", "),
		recursiveStr,
		m.eventID,
		time.Since(m.startTime).Round(time.Second),
	))

	resLine := m.renderResourceBar()

	footer := footerStyle.Render("Press q or ctrl+c to exit")

	return clipToWidth(appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		info,
		m.viewport.View(),
		resLine,
		"",
		footer,
	)), m.width)
}

// clipToWidth truncates every rendered line to width. Chrome lines and event
// lines wider than the terminal would otherwise wrap in xterm and silently
// consume extra rows, pushing the resource bar and footer out of the window.
func clipToWidth(view string, width int) string {
	if width <= 0 {
		return view
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(view)
}

func (m *model) renderResourceBar() string {
	return formatResourceBar(m.lastStats, m.startTime, m.eventsPerSec)
}

// syncViewport sizes the viewport, creating it on first use, then re-renders
// the viewport content through render. fixedRows is every row the view spends
// outside the viewport (headers, footers, separators and margins); the
// viewport gets exactly the remaining height so the whole view fits the
// reported terminal size.
func syncViewport(v *viewport.Model, ready *bool, width, height, fixedRows int, render func()) {
	viewportWidth := max(1, width-appStyle.GetHorizontalMargins())
	viewportHeight := max(1, height-fixedRows)
	if !*ready {
		*v = viewport.New(viewportWidth, viewportHeight)
		*ready = true
	} else {
		v.Width = viewportWidth
		v.Height = viewportHeight
	}
	render()
}

// formatResourceBar renders the shared CPU/RSS/event-rate status line.
func formatResourceBar(s *procstats.Stats, start time.Time, eventsPerSec float64) string {
	if s == nil {
		return ""
	}

	rssMB := float64(s.RSS) / 1024 / 1024
	cpuTotal := s.CPUUser + s.CPUSys
	cpuSec := cpuTotal.Seconds()
	uptimeSec := time.Since(start).Seconds()
	cpuPct := 0.0
	if uptimeSec > 0 {
		cpuPct = (cpuSec / uptimeSec) * 100
	}

	rss := memStyle.Render(fmt.Sprintf("%.1f MB", rssMB))
	cpu := cpuStyle.Render(fmt.Sprintf("%.1f%%", cpuPct))
	eps := epsStyle.Render(fmt.Sprintf("%.0f/s", eventsPerSec))

	return resourceStyle.Render(
		resourceLabelStyle.Render(" Mem:") + " " + rss +
			resourceLabelStyle.Render("  CPU:") + " " + cpu +
			resourceLabelStyle.Render("  Evt/s:") + " " + eps,
	)
}
