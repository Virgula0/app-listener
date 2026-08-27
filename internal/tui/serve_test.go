package tui

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/usecase"
)

type serveTestModel struct {
	width  int
	height int
}

func (m *serveTestModel) Init() tea.Cmd { return nil }

func (m *serveTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
	}
	return m, nil
}

func (m *serveTestModel) View() string {
	return fmt.Sprintf("%dx%d", m.width, m.height)
}

func TestParseResize(t *testing.T) {
	tests := []struct {
		name string
		body string
		ok   bool
	}{
		{name: "valid", body: `{"type":"resize","cols":120,"rows":40}`, ok: true},
		{name: "tiny terminal", body: `{"type":"resize","cols":2,"rows":1}`, ok: true},
		{name: "huge terminal", body: `{"type":"resize","cols":1000,"rows":500}`, ok: true},
		{name: "one column", body: `{"type":"resize","cols":1,"rows":40}`},
		{name: "zero rows", body: `{"type":"resize","cols":120,"rows":0}`},
		{name: "overwide", body: `{"type":"resize","cols":1001,"rows":40}`},
		{name: "too tall", body: `{"type":"resize","cols":120,"rows":501}`},
		{name: "wrong type", body: `{"type":"input","cols":120,"rows":40}`},
		{name: "unknown field", body: `{"type":"resize","cols":120,"rows":40,"key":"q"}`},
		{name: "trailing value", body: `{"type":"resize","cols":120,"rows":40} {}`},
		{name: "invalid json", body: `{`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			size, err := parseResize([]byte(test.body))
			if test.ok && err != nil {
				t.Fatalf("parseResize() error = %v", err)
			}
			if !test.ok && err == nil {
				t.Fatal("parseResize() unexpectedly succeeded")
			}
			if test.ok && (size.Width == 0 || size.Height == 0) {
				t.Fatalf("parseResize() returned invalid size: %#v", size)
			}
		})
	}
}

func TestBasicAuth(t *testing.T) {
	protected := basicAuth("alice", "secret", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, test := range []struct {
		name     string
		user     string
		password string
		want     int
	}{
		{name: "missing", want: http.StatusUnauthorized},
		{name: "wrong user", user: "mallory", password: "secret", want: http.StatusUnauthorized},
		{name: "wrong password", user: "alice", password: "wrong", want: http.StatusUnauthorized},
		{name: "valid", user: "alice", password: "secret", want: http.StatusNoContent},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:9999/", nil)
			if test.user != "" {
				req.SetBasicAuth(test.user, test.password)
			}
			response := httptest.NewRecorder()
			protected.ServeHTTP(response, req)
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
			if test.want == http.StatusUnauthorized && response.Header().Get("WWW-Authenticate") == "" {
				t.Fatal("unauthorized response has no WWW-Authenticate header")
			}
		})
	}
}

func TestSameOrigin(t *testing.T) {
	for _, test := range []struct {
		name   string
		origin string
		want   bool
	}{
		{name: "http", origin: "http://127.0.0.1:9999", want: true},
		{name: "https proxy", origin: "https://127.0.0.1:9999", want: true},
		{name: "missing"},
		{name: "foreign", origin: "https://example.com"},
		{name: "credentialed", origin: "http://user@127.0.0.1:9999"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:9999/ws", nil)
			req.Host = "127.0.0.1:9999"
			if test.origin != "" {
				req.Header.Set("Origin", test.origin)
			}
			if got := sameOrigin(req); got != test.want {
				t.Fatalf("sameOrigin() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestSnapshotHubSingleViewerAndCoalescing(t *testing.T) {
	hub := &snapshotHub{}
	updates, _, ok := hub.reserve()
	if !ok {
		t.Fatal("first viewer was rejected")
	}
	if _, _, ok := hub.reserve(); ok {
		t.Fatal("second viewer was accepted")
	}

	hub.publish("first")
	hub.publish("latest")
	if got := <-updates; got != "latest" {
		t.Fatalf("snapshot = %q, want latest", got)
	}
	hub.release(nil)
	if _, latest, ok := hub.reserve(); !ok || latest != "latest" {
		t.Fatalf("viewer could not reconnect or latest snapshot was lost: ok=%t latest=%q", ok, latest)
	}
}

func TestServeHandlerProtectsPageAndWebSocket(t *testing.T) {
	handler := serveHandler(ServeOptions{
		Address: "127.0.0.1:9999", Username: "alice", Password: "secret",
	}, nil, &snapshotHub{})

	for _, path := range []string{"/", "/ws"} {
		req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:9999"+path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != http.StatusUnauthorized {
			t.Errorf("GET %s status = %d, want %d", path, response.Code, http.StatusUnauthorized)
		}
	}
}

func TestServeWebSocketSnapshotsResizeAndSingleViewer(t *testing.T) {
	hub := &snapshotHub{}
	model := &servedModel{model: &serveTestModel{}, hub: hub}
	program := tea.NewProgram(model, tea.WithInput(nil), tea.WithOutput(io.Discard), tea.WithoutRenderer())
	programDone := make(chan error, 1)
	go func() {
		_, err := program.Run()
		programDone <- err
	}()
	t.Cleanup(func() {
		program.Quit()
		select {
		case <-programDone:
		case <-time.After(time.Second):
			t.Error("Bubble Tea program did not stop")
		}
	})

	server := httptest.NewUnstartedServer(nil)
	address := server.Listener.Addr().String()
	server.Config.Handler = serveHandler(ServeOptions{Address: address}, program, hub)
	server.Start()
	t.Cleanup(server.Close)

	header := http.Header{"Origin": []string{"http://" + address}}
	conn, response, err := websocket.DefaultDialer.Dial("ws://"+address+"/ws", header)
	if err != nil {
		t.Fatalf("dialing WebSocket: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if response != nil {
		_ = response.Body.Close()
	}

	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	if _, snapshot, err := conn.ReadMessage(); err != nil || string(snapshot) != "0x0" {
		t.Fatalf("initial snapshot = %q, err = %v", snapshot, err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"resize","cols":120,"rows":40}`)); err != nil {
		t.Fatalf("sending resize: %v", err)
	}
	if _, snapshot, err := conn.ReadMessage(); err != nil || string(snapshot) != "120x40" {
		t.Fatalf("resized snapshot = %q, err = %v", snapshot, err)
	}

	second, response, err := websocket.DefaultDialer.Dial("ws://"+address+"/ws", header)
	if second != nil {
		_ = second.Close()
	}
	if err == nil {
		t.Fatal("second WebSocket viewer was accepted")
	}
	if response == nil || response.StatusCode != http.StatusConflict {
		t.Fatalf("second viewer status = %#v, want %d", response, http.StatusConflict)
	}
	_ = response.Body.Close()
}

func TestSecurityHeadersAndExactRoot(t *testing.T) {
	handler := serveHandler(ServeOptions{Address: "127.0.0.1:9999"}, nil, &snapshotHub{})

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:9999/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("root status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'") {
		t.Fatal("root response is missing restrictive CSP")
	}

	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1:9999/not-found", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusNotFound {
		t.Fatalf("unknown path status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestSanitizeTerminalText(t *testing.T) {
	got := sanitizeTerminalText("safe\x1b]52;c;payload\a\npath")
	if strings.ContainsAny(got, "\x1b\a\n") {
		t.Fatalf("control characters were not sanitized: %q", got)
	}
}

func TestEventFanoutDuplicatesInOrderAndCloses(t *testing.T) {
	source := make(chan int, 8)
	fan := NewEventFanout(source)

	for i := 0; i < 5; i++ {
		source <- i
	}
	close(source)

	for i := 0; i < 5; i++ {
		local, ok := <-fan.Local()
		if !ok || local != i {
			t.Fatalf("local[%d] = %d, ok = %t", i, local, ok)
		}
		browser, ok := <-fan.Browser()
		if !ok || browser != i {
			t.Fatalf("browser[%d] = %d, ok = %t", i, browser, ok)
		}
	}

	fan.Stop()
	if _, ok := <-fan.Local(); ok {
		t.Error("local stream was not closed after the source closed")
	}
	if _, ok := <-fan.Browser(); ok {
		t.Error("browser stream was not closed after the source closed")
	}
}

func TestEventFanoutStopClosesOutputs(t *testing.T) {
	source := make(chan int)
	fan := NewEventFanout(source)
	fan.Stop()
	if _, ok := <-fan.Local(); ok {
		t.Error("local stream was not closed by Stop")
	}
	if _, ok := <-fan.Browser(); ok {
		t.Error("browser stream was not closed by Stop")
	}
}

func TestEventFanoutDropOldestUnderPressure(t *testing.T) {
	const total = fanoutBuffer + 10
	source := make(chan int, total)
	for i := 0; i < total; i++ {
		source <- i
	}
	close(source)

	fan := NewEventFanout(source)
	<-fan.dispatchd // dispatcher drained the source and closed both outputs

	// Neither subscriber consumed while the dispatcher drained: both queues
	// overflowed and must hold exactly the newest fanoutBuffer events in
	// order instead of blocking the source reader.
	for _, stream := range []struct {
		name   string
		events <-chan int
	}{{"local", fan.Local()}, {"browser", fan.Browser()}} {
		first, ok := <-stream.events
		if !ok {
			t.Fatalf("%s stream closed without events", stream.name)
		}
		if want := total - fanoutBuffer; first != want {
			t.Errorf("%s stream starts at %d, want %d (oldest dropped)", stream.name, first, want)
		}
		count := 1
		for range stream.events {
			count++
		}
		if count != fanoutBuffer {
			t.Errorf("%s stream delivered %d events, want %d", stream.name, count, fanoutBuffer)
		}
	}
}

// assertViewFits verifies a rendered view never exceeds the terminal size
// the browser or terminal reported: overflow silently wraps in xterm and
// pushes the resource bar and footer out of the window.
func assertViewFits(t *testing.T, name, view string, width, height int) {
	t.Helper()
	if viewHeight := lipgloss.Height(view); viewHeight > height {
		t.Errorf("%s: view height %d exceeds terminal height %d", name, viewHeight, height)
	}
	if viewWidth := lipgloss.Width(view); viewWidth > width {
		t.Errorf("%s: view width %d exceeds terminal width %d", name, viewWidth, width)
	}
}

func TestViewsFitReportedTerminalSize(t *testing.T) {
	const (
		width  = 80
		height = 24
	)
	longPath := strings.Repeat("/long/path/", 20)
	longBinary := strings.Repeat("binary-", 20)

	fileEvents := make(chan ebpf.FileEvent, 8)
	fileModel := NewModel(fileEvents, []string{longPath}, true, 3).(*model)
	fileModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
	for i := 0; i < 30; i++ {
		ev := ebpf.FileEvent{
			Timestamp: time.Now().UnixNano(), Type: ebpf.EventRename, PID: 1,
			Comm: longBinary, Path: longPath, Dest: longPath,
		}
		fileModel.addEvent(&ev)
	}
	fileModel.renderViewport()
	assertViewFits(t, "monitor", fileModel.View(), width, height)

	guardEvents := make(chan guard.GuardEvent, 8)
	guardModel := NewGuardModel(guardEvents, longPath, guard.ModeWhitelist,
		[]guard.BinaryEntry{{Path: longBinary}}).(*guardModel)
	guardModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
	for i := 0; i < 30; i++ {
		ev := guard.GuardEvent{
			FileEvent: ebpf.FileEvent{
				Timestamp: time.Now().UnixNano(), Type: ebpf.EventOpen, PID: 1,
				Comm: longBinary, Path: longPath,
			},
			Blocked: i%2 == 0,
		}
		guardModel.addGuardEvent(&ev)
	}
	guardModel.renderGuardViewport()
	assertViewFits(t, "guard", guardModel.View(), width, height)

	netEvents := make(chan ebpf.NetEvent, 8)
	netModel := NewNetModel(netEvents, []ebpf.BinaryEntry{{Path: longBinary}}).(*netModel)
	netModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
	for i := 0; i < 30; i++ {
		ev := ebpf.NetEvent{
			Timestamp: time.Now().UnixNano(), Type: ebpf.NetConnect, PID: 1,
			Comm: longBinary, SrcAddr: "10.0.0.1:1", DstAddr: strings.Repeat("a", 90),
		}
		netModel.addNetEvent(&ev)
	}
	netModel.renderNetViewport()
	assertViewFits(t, "network-monitor", netModel.View(), width, height)

	netGuardEvents := make(chan networkguard.NetGuardEvent, 8)
	netGuardModel := NewNetGuardModel(netGuardEvents, networkguard.ModeWhitelist,
		[]networkguard.BinaryEntry{{Path: longBinary}}).(*netGuardModel)
	netGuardModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
	for i := 0; i < 40; i++ {
		ev := networkguard.NetGuardEvent{
			NetEvent: ebpf.NetEvent{
				Timestamp: time.Now().UnixNano(), Type: ebpf.NetConnect, PID: 1,
				Comm: longBinary, DstAddr: strings.Repeat("a", 90),
			},
			Blocked: true,
		}
		netGuardModel.lines = append(netGuardModel.lines, netGuardEventLine{event: ev, line: netGuardModel.formatEvent(&ev)})
	}
	netGuardModel.renderContent()
	assertViewFits(t, "network-guard", netGuardModel.View(), width, height)

	daemonEvents := make(chan usecase.DaemonEvent, 8)
	daemonModel := NewDaemonModel(daemonEvents, []DaemonResourceInfo{
		{Path: longPath, NeedEncryption: true, Binaries: 2},
		{Path: "/r2", Binaries: 1},
	}).(*daemonModel)
	daemonModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
	for i := 0; i < 20; i++ {
		ev := usecase.DaemonEvent{
			Resource: longPath,
			Event: guard.GuardEvent{
				FileEvent: ebpf.FileEvent{
					Timestamp: time.Now().UnixNano(), Type: ebpf.EventOpen, PID: 1,
					Comm: "ssh", Path: longPath,
				},
				Blocked: i%2 == 0,
			},
		}
		daemonModel.addDaemonEvent(&ev)
	}
	daemonModel.renderDaemonViewport()
	assertViewFits(t, "daemon", daemonModel.View(), width, height)
}
