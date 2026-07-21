package gui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type guiModel struct {
	events    <-chan ebpf.FileEvent
	allEvents []ebpf.FileEvent
	filtered  []ebpf.FileEvent
	filter    string
	mu        sync.RWMutex
	directory string
	recursive bool
	depth     int
}

func Run(events <-chan ebpf.FileEvent, directory string, recursive bool, depth int) {
	m := &guiModel{
		events:    events,
		allEvents: make([]ebpf.FileEvent, 0, 500),
		filtered:  make([]ebpf.FileEvent, 0, 500),
		directory: directory,
		recursive: recursive,
		depth:     depth,
	}

	os.Setenv("FYNE_THEME", "dark")

	fyneApp := app.NewWithID("com.virgula0.app-listener")
	win := fyneApp.NewWindow("app-listener — File System Monitor (eBPF)")

	recursiveStr := "no"
	if recursive {
		recursiveStr = fmt.Sprintf("yes (depth:%d)", depth)
	}

	infoLabel := widget.NewLabel(
		fmt.Sprintf("Directory: %s  |  Recursive: %s  |  Events: 0  |  Uptime: 0s",
			directory, recursiveStr))

	filterEntry := widget.NewEntry()
	filterEntry.SetPlaceHolder("Filter by path, process, or event type...")
	filterEntry.OnChanged = func(s string) {
		m.mu.Lock()
		m.filter = strings.ToLower(s)
		m.applyFilter()
		m.mu.Unlock()
	}

	colWidths := []float32{120, 80, 90, 400, 200}

	table := widget.NewTable(
		func() (int, int) {
			m.mu.RLock()
			n := len(m.filtered)
			m.mu.RUnlock()
			return n, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			m.mu.RLock()
			if id.Row >= len(m.filtered) {
				m.mu.RUnlock()
				return
			}
			ev := m.filtered[id.Row]
			m.mu.RUnlock()

			switch id.Col {
			case 0:
				cell.(*widget.Label).SetText(
					time.Unix(0, ev.Timestamp).Format("15:04:05.000"))
			case 1:
				cell.(*widget.Label).SetText(ev.Type.String())
			case 2:
				cell.(*widget.Label).SetText(ev.Comm)
			case 3:
				if ev.PID != 0 {
					cell.(*widget.Label).SetText(fmt.Sprintf("%s (pid=%d)", ev.Path, ev.PID))
				} else {
					cell.(*widget.Label).SetText(ev.Path)
				}
			case 4:
				detail := ""
				switch ev.Type {
				case ebpf.EventRead, ebpf.EventWrite:
					detail = fmt.Sprintf("fd=%d", ev.FD)
				case ebpf.EventRename, ebpf.EventSymlink, ebpf.EventHardlink:
					detail = ev.Dest
				}
				cell.(*widget.Label).SetText(detail)
			}
		},
	)

	for i, w := range colWidths {
		table.SetColumnWidth(i, w)
	}

	top := container.NewVBox(infoLabel, filterEntry)
	content := container.NewBorder(top, nil, nil, nil, table)
	win.SetContent(content)
	win.Resize(fyne.NewSize(1100, 650))

	startTime := time.Now()

	go m.listenEvents(infoLabel, table, startTime, recursiveStr)
	go m.uptimeTicker(infoLabel, startTime, recursiveStr)

	floatWindow()

	win.ShowAndRun()
}

func floatWindow() {
	go func() {
		time.Sleep(150 * time.Millisecond)
		exec.Command("xprop",
			"-name", "app-listener",
			"-f", "_NET_WM_WINDOW_TYPE", "32a",
			"-set", "_NET_WM_WINDOW_TYPE", "_NET_WM_WINDOW_TYPE_UTILITY",
		).Run()
	}()
}

func (m *guiModel) listenEvents(infoLabel *widget.Label, table *widget.Table, startTime time.Time, recursiveStr string) {
	for ev := range m.events {
		m.mu.Lock()
		m.allEvents = append(m.allEvents, ev)
		if len(m.allEvents) > 500 {
			m.allEvents = m.allEvents[len(m.allEvents)-500:]
		}
		m.applyFilter()
		count := len(m.allEvents)
		uptime := time.Since(startTime).Round(time.Second)
		m.mu.Unlock()

		fyne.Do(func() {
			table.Refresh()
			infoLabel.SetText(fmt.Sprintf(
				"Directory: %s  |  Recursive: %s  |  Events: %d  |  Uptime: %s",
				m.directory, recursiveStr, count, uptime,
			))
		})
	}
}

func (m *guiModel) uptimeTicker(infoLabel *widget.Label, startTime time.Time, recursiveStr string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.RLock()
		count := len(m.allEvents)
		m.mu.RUnlock()

		uptime := time.Since(startTime).Round(time.Second)
		fyne.Do(func() {
			infoLabel.SetText(fmt.Sprintf(
				"Directory: %s  |  Recursive: %s  |  Events: %d  |  Uptime: %s",
				m.directory, recursiveStr, count, uptime,
			))
		})
	}
}

func (m *guiModel) applyFilter() {
	if m.filter == "" {
		m.filtered = make([]ebpf.FileEvent, len(m.allEvents))
		copy(m.filtered, m.allEvents)
		return
	}

	m.filtered = m.filtered[:0]
	for _, ev := range m.allEvents {
		if strings.Contains(strings.ToLower(ev.Path), m.filter) ||
			strings.Contains(strings.ToLower(ev.Comm), m.filter) ||
			strings.Contains(strings.ToLower(ev.Type.String()), m.filter) ||
			strings.Contains(strings.ToLower(ev.Dest), m.filter) {
			m.filtered = append(m.filtered, ev)
		}
	}
}
