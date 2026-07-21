package gui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

var colLabels = []string{"Time", "Type", "Process", "Path", "Detail"}
var colDefaultWidths = []float32{120, 80, 90, 400, 200}

type columnState struct {
	sync.RWMutex
	widths []float32
}

func newColumnState(widths []float32) *columnState {
	w := make([]float32, len(widths))
	copy(w, widths)
	return &columnState{widths: w}
}

type guiModel struct {
	events    <-chan ebpf.FileEvent
	allEvents []ebpf.FileEvent
	filtered  []ebpf.FileEvent
	filter    string
	mu        sync.RWMutex
	directory string
	recursive bool
	depth     int
	cols      *columnState
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

func (m *guiModel) getFiltered(i int) ebpf.FileEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if i >= len(m.filtered) {
		return ebpf.FileEvent{}
	}
	return m.filtered[i]
}

func (m *guiModel) filteredLen() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.filtered)
}

// ── dataRow ────────────────────────────────────────────────────────────────

type dataRow struct {
	widget.BaseWidget
	cells []*widget.Label
	cols  *columnState
}

func newDataRow(cols *columnState) *dataRow {
	r := &dataRow{cols: cols}
	for range colLabels {
		l := widget.NewLabel("")
		l.Selectable = true
		l.Wrapping = fyne.TextTruncate
		r.cells = append(r.cells, l)
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *dataRow) setEvent(ev ebpf.FileEvent) {
	r.cells[0].SetText(time.Unix(0, ev.Timestamp).Format("15:04:05.000"))
	r.cells[1].SetText(ev.Type.String())
	r.cells[2].SetText(ev.Comm)

	path := ev.Path
	if ev.PID != 0 {
		path = fmt.Sprintf("%s (pid=%d)", ev.Path, ev.PID)
	}
	r.cells[3].SetText(path)

	detail := ""
	switch ev.Type {
	case ebpf.EventRead, ebpf.EventWrite:
		detail = fmt.Sprintf("fd=%d", ev.FD)
	case ebpf.EventRename, ebpf.EventSymlink, ebpf.EventHardlink:
		detail = ev.Dest
	}
	r.cells[4].SetText(detail)
}

func (r *dataRow) CreateRenderer() fyne.WidgetRenderer {
	objs := make([]fyne.CanvasObject, len(r.cells))
	for i, c := range r.cells {
		objs[i] = c
	}
	return &dataRowRenderer{row: r, objs: objs}
}

type dataRowRenderer struct {
	row  *dataRow
	objs []fyne.CanvasObject
}

func (rr *dataRowRenderer) Layout(size fyne.Size) {
	rr.row.cols.RLock()
	widths := rr.row.cols.widths
	rr.row.cols.RUnlock()

	total := float32(0)
	for _, w := range widths {
		total += w
	}
	if total == 0 {
		return
	}

	x := float32(0)
	for i, l := range rr.row.cells {
		w := size.Width * widths[i] / total
		l.Resize(fyne.NewSize(w, size.Height))
		l.Move(fyne.NewPos(x, 0))
		x += w
	}
}

func (rr *dataRowRenderer) MinSize() fyne.Size {
	return fyne.NewSize(800, 36)
}

func (rr *dataRowRenderer) Objects() []fyne.CanvasObject {
	return rr.objs
}

func (rr *dataRowRenderer) Refresh() {
	rr.Layout(rr.row.Size())
}

func (rr *dataRowRenderer) Destroy() {}

// ── headerWidget ───────────────────────────────────────────────────────────

type headerWidget struct {
	widget.BaseWidget
	labels   []*widget.Label
	cols     *columnState
	onResize func()
}

func newHeaderWidget(cols *columnState, onResize func()) *headerWidget {
	h := &headerWidget{
		cols:     cols,
		onResize: onResize,
	}
	for _, label := range colLabels {
		l := widget.NewLabelWithStyle(label, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		h.labels = append(h.labels, l)
	}
	h.ExtendBaseWidget(h)
	return h
}

// dividerIndex returns the column divider index closest to x, or -1 if none.
func (h *headerWidget) dividerIndex(x float32, size fyne.Size) int {
	h.cols.RLock()
	total := float32(0)
	for _, w := range h.cols.widths {
		total += w
	}
	h.cols.RUnlock()
	if total == 0 || size.Width == 0 {
		return -1
	}

	cum := float32(0)
	for i := 0; i < len(h.cols.widths)-1; i++ {
		cum += size.Width * h.cols.widths[i] / total
		if x >= cum-4 && x <= cum+4 {
			return i
		}
	}
	return -1
}

func (h *headerWidget) Dragged(e *fyne.DragEvent) {
	idx := h.dividerIndex(e.Position.X, h.Size())
	if idx < 0 || idx+1 >= len(h.cols.widths) {
		return
	}

	h.cols.Lock()
	h.cols.widths[idx] += e.Dragged.DX
	h.cols.widths[idx+1] -= e.Dragged.DX
	minW := float32(40)
	if h.cols.widths[idx] < minW {
		h.cols.widths[idx] = minW
	}
	if h.cols.widths[idx+1] < minW {
		h.cols.widths[idx+1] = minW
	}
	h.cols.Unlock()

	h.Refresh()
	if h.onResize != nil {
		h.onResize()
	}
}

func (h *headerWidget) DragEnd() {}

func (h *headerWidget) MinSize() fyne.Size {
	return fyne.NewSize(800, 36)
}

func (h *headerWidget) CreateRenderer() fyne.WidgetRenderer {
	objs := make([]fyne.CanvasObject, len(h.labels))
	for i, l := range h.labels {
		objs[i] = l
	}
	return &headerRenderer{h: h, objs: objs}
}

type headerRenderer struct {
	h    *headerWidget
	objs []fyne.CanvasObject
}

func (hr *headerRenderer) Layout(size fyne.Size) {
	hr.h.cols.RLock()
	widths := hr.h.cols.widths
	hr.h.cols.RUnlock()

	total := float32(0)
	for _, w := range widths {
		total += w
	}
	if total == 0 || size.Width == 0 {
		return
	}

	x := float32(0)
	for i, l := range hr.h.labels {
		w := size.Width * widths[i] / total
		l.Resize(fyne.NewSize(w, size.Height))
		l.Move(fyne.NewPos(x, 0))
		x += w
	}
}

func (hr *headerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(800, 36)
}

func (hr *headerRenderer) Objects() []fyne.CanvasObject {
	return hr.objs
}

func (hr *headerRenderer) Refresh() {
	hr.Layout(hr.h.Size())
}

func (hr *headerRenderer) Destroy() {}

// ── Run ────────────────────────────────────────────────────────────────────

func Run(events <-chan ebpf.FileEvent, directory string, recursive bool, depth int) {
	m := &guiModel{
		events:    events,
		allEvents: make([]ebpf.FileEvent, 0, 500),
		filtered:  make([]ebpf.FileEvent, 0, 500),
		directory: directory,
		recursive: recursive,
		depth:     depth,
		cols:      newColumnState(colDefaultWidths),
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

	header := newHeaderWidget(m.cols, func() {
		// called after columns are resized
	})

	list := widget.NewList(
		m.filteredLen,
		func() fyne.CanvasObject {
			return newDataRow(m.cols)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			row := item.(*dataRow)
			ev := m.getFiltered(id)
			row.setEvent(ev)
		},
	)

	header.onResize = func() {
		list.Refresh()
	}

	top := container.NewVBox(infoLabel, filterEntry, header)
	content := container.NewBorder(top, nil, nil, nil, list)
	win.SetContent(content)
	win.Resize(fyne.NewSize(1100, 650))

	startTime := time.Now()
	go m.listenEvents(infoLabel, list, startTime, recursiveStr)
	go m.uptimeTicker(infoLabel, startTime, recursiveStr)
	floatWindow()

	win.ShowAndRun()
}

func (m *guiModel) listenEvents(
	infoLabel *widget.Label, list *widget.List,
	startTime time.Time, recursiveStr string,
) {
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
			list.Refresh()
			infoLabel.SetText(fmt.Sprintf(
				"Directory: %s  |  Recursive: %s  |  Events: %d  |  Uptime: %s",
				m.directory, recursiveStr, count, uptime,
			))
		})
	}
}

func (m *guiModel) uptimeTicker(
	infoLabel *widget.Label, startTime time.Time, recursiveStr string,
) {
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

// ── X11 float hint ─────────────────────────────────────────────────────────

func floatWindow() {
	go func() {
		time.Sleep(200 * time.Millisecond)
		setFloatHint()
	}()
}

func setFloatHint() {
	conn, err := xgb.NewConn()
	if err != nil {
		return
	}
	defer conn.Close()

	setup := xproto.Setup(conn)
	root := setup.DefaultScreen(conn).Root

	wid := findTopLevelWindow(conn, root, "app-listener")
	if wid == 0 {
		return
	}

	typeAtom := internAtom(conn, "_NET_WM_WINDOW_TYPE")
	utilityAtom := internAtom(conn, "_NET_WM_WINDOW_TYPE_DIALOG")
	if typeAtom == 0 || utilityAtom == 0 {
		return
	}

	data := make([]byte, 4)
	xgb.Put32(data, uint32(utilityAtom))
	xproto.ChangeProperty(conn, xproto.PropModeReplace, wid,
		typeAtom, xproto.AtomAtom, 32, 1, data)
}

func findTopLevelWindow(conn *xgb.Conn, root xproto.Window, substr string) xproto.Window {
	netClientList := internAtom(conn, "_NET_CLIENT_LIST")
	if netClientList == 0 {
		return 0
	}

	reply, err := xproto.GetProperty(conn, false, root,
		netClientList, xproto.AtomWindow, 0, 4096).Reply()
	if err != nil || reply == nil {
		return 0
	}

	netWmName := internAtom(conn, "_NET_WM_NAME")
	utf8Str := internAtom(conn, "UTF8_STRING")
	if netWmName == 0 || utf8Str == 0 {
		return 0
	}

	for i := 0; i < len(reply.Value)/4; i++ {
		wid := xproto.Window(xgb.Get32(reply.Value[i*4:]))
		nameReply, err := xproto.GetProperty(conn, false, wid,
			netWmName, utf8Str, 0, 4096).Reply()
		if err != nil || nameReply == nil || len(nameReply.Value) == 0 {
			continue
		}
		if strings.Contains(string(nameReply.Value), substr) {
			return wid
		}
	}
	return 0
}

func internAtom(conn *xgb.Conn, name string) xproto.Atom {
	reply, err := xproto.InternAtom(conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		return 0
	}
	return reply.Atom
}
