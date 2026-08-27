package tui

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	minServeCols = 2
	maxServeCols = 1000
	minServeRows = 1
	maxServeRows = 500

	// fanoutBuffer is the per-subscriber event queue of the dual-TUI fanout.
	// A subscriber that cannot keep up drops its oldest queued event instead
	// of blocking the engine reader (which would overflow the kernel ring
	// buffer and lose events for both subscribers).
	fanoutBuffer = 256
)

//go:embed embeds/serve.html
var serveHTML []byte

type ServeOptions struct {
	Address  string
	Username string
	Password string
	Reload   func()
}

type resizeMessage struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// EventFanout turns the single point-to-point engine event channel into two
// ordered streams so the local TUI and the browser TUI can run as two
// independent Bubble Tea programs (separate dimensions, scroll state and
// counters) without competing for events. Stop closes both output channels
// after the dispatcher exits, releasing any pending event receive command.
type EventFanout[T any] struct {
	source    <-chan T
	local     chan T
	browser   chan T
	stop      chan struct{}
	stopOnce  sync.Once
	dispatchd chan struct{}
}

// NewEventFanout starts duplicating source into two independent subscribers.
func NewEventFanout[T any](source <-chan T) *EventFanout[T] {
	f := &EventFanout[T]{
		source:    source,
		local:     make(chan T, fanoutBuffer),
		browser:   make(chan T, fanoutBuffer),
		stop:      make(chan struct{}),
		dispatchd: make(chan struct{}),
	}
	go f.dispatch()
	return f
}

// Local returns the event stream consumed by the local terminal TUI.
func (f *EventFanout[T]) Local() <-chan T { return f.local }

// Browser returns the event stream consumed by the browser TUI.
func (f *EventFanout[T]) Browser() <-chan T { return f.browser }

// Stop stops the dispatcher and closes both output channels. Safe to call
// once; it blocks until the outputs are closed.
func (f *EventFanout[T]) Stop() {
	f.stopOnce.Do(func() { close(f.stop) })
	<-f.dispatchd
}

func (f *EventFanout[T]) dispatch() {
	defer close(f.dispatchd)
	defer close(f.browser)
	defer close(f.local)

	for {
		select {
		case ev, ok := <-f.source:
			if !ok {
				return
			}
			f.deliver(f.local, ev)
			f.deliver(f.browser, ev)
		case <-f.stop:
			return
		}
	}
}

// deliver queues ev without ever blocking: on a full queue the oldest queued
// event is dropped first, so one slow subscriber cannot stall the other or
// back-pressure the engine into kernel ring-buffer loss.
func (f *EventFanout[T]) deliver(queue chan T, ev T) {
	select {
	case queue <- ev:
		return
	default:
	}
	select {
	case <-queue:
	default:
	}
	select {
	case queue <- ev:
	default:
	}
}

type snapshotHub struct {
	mu       sync.Mutex
	latest   string
	reserved bool
	closed   bool
	updates  chan string
	conn     *websocket.Conn
}

func (h *snapshotHub) publish(snapshot string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.latest = snapshot
	if h.updates == nil {
		return
	}
	select {
	case <-h.updates:
	default:
	}
	select {
	case h.updates <- snapshot:
	default:
	}
}

func (h *snapshotHub) reserve() (updates chan string, latest string, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.reserved || h.closed {
		return nil, "", false
	}
	h.reserved = true
	h.updates = make(chan string, 1)
	return h.updates, h.latest, true
}

func (h *snapshotHub) attach(conn *websocket.Conn) {
	h.mu.Lock()
	h.conn = conn
	h.mu.Unlock()
}

func (h *snapshotHub) release(conn *websocket.Conn) {
	h.mu.Lock()
	if h.conn == conn || conn == nil {
		h.conn = nil
		h.updates = nil
		h.reserved = false
	}
	h.mu.Unlock()
}

func (h *snapshotHub) close() {
	h.mu.Lock()
	h.closed = true
	conn := h.conn
	h.mu.Unlock()
	if conn != nil {
		_ = conn.Close()
	}
}

type servedModel struct {
	model tea.Model
	hub   *snapshotHub
}

func (m *servedModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m *servedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.model.Update(msg)
	m.model = model
	return m, cmd
}

func (m *servedModel) View() string {
	m.hub.publish(m.model.View())
	return ""
}

// Serve runs the local TUI and the browser TUI side by side. Both models must
// be independent instances fed by an EventFanout: they own their dimensions
// and state separately. Quitting the local TUI (q, ctrl+c or SIGINT/SIGTERM)
// shuts down the browser endpoint; the process keeps its secure shutdown
// ordering because the caller's deferred use case stop runs after Serve
// returns.
func Serve(local, browser tea.Model, options ServeOptions) error {
	if !isInteractiveTerminal() {
		return errors.New("--serve needs an interactive terminal for the local TUI; use --headless without --serve for services")
	}

	listener, err := newServeListener(options.Address)
	if err != nil {
		return err
	}

	hub := &snapshotHub{}
	browserProgram := tea.NewProgram(
		&servedModel{model: browser, hub: hub},
		tea.WithInput(nil),
		tea.WithOutput(io.Discard),
		tea.WithoutRenderer(),
		tea.WithoutSignalHandler(),
	)
	localProgram := tea.NewProgram(local, tea.WithAltScreen(), tea.WithoutSignalHandler())
	server := newServeServer(serveHandler(options, browserProgram, hub))

	return runServe(localProgram, browserProgram, server, listener, hub, options)
}

func isInteractiveTerminal() bool {
	return term.IsTerminal(os.Stdin.Fd())
}

func newServeListener(address string) (net.Listener, error) {
	var config net.ListenConfig
	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("starting browser TUI listener: %w", err)
	}
	return listener, nil
}

func newServeServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
}

func runServe(localProgram, browserProgram *tea.Program, server *http.Server, listener net.Listener, hub *snapshotHub, options ServeOptions) error {
	serveSignals := make(chan os.Signal, 1)
	signal.Notify(serveSignals, syscall.SIGINT, syscall.SIGTERM)
	if options.Reload != nil {
		signal.Notify(serveSignals, syscall.SIGHUP)
	}
	stopRelay := make(chan struct{})
	relayDone := make(chan struct{})
	go relayServeSignals(serveSignals, stopRelay, relayDone, localProgram, options.Reload)

	localDone := make(chan error, 1)
	go func() { localDone <- waitProgram(localProgram) }()
	browserDone := make(chan error, 1)
	go func() { browserDone <- waitProgram(browserProgram) }()
	serverDone := make(chan error, 1)
	go func() { serverDone <- serveHTTP(server, listener, localProgram) }()

	log.Infof("browser TUI available at http://%s", options.Address)
	runErr := <-localDone

	signal.Stop(serveSignals)
	close(stopRelay)
	<-relayDone // waits for an in-flight daemon reload before teardown

	browserProgram.Quit()
	<-browserDone
	hub.close()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdownErr := server.Shutdown(shutdownCtx)
	serveErr := <-serverDone

	if serveErr != nil {
		return fmt.Errorf("serving browser TUI: %w", serveErr)
	}
	if shutdownErr != nil {
		return fmt.Errorf("stopping browser TUI server: %w", shutdownErr)
	}
	if runErr != nil && !errors.Is(runErr, tea.ErrProgramKilled) {
		return runErr
	}
	return nil
}

func relayServeSignals(signals chan os.Signal, stop, done chan struct{}, localProgram *tea.Program, reload func()) {
	defer close(done)
	for {
		select {
		case s := <-signals:
			if reload != nil && s == syscall.SIGHUP {
				log.Info("daemon: SIGHUP received, reloading configuration")
				reload()
				continue
			}
			localProgram.Quit()
		case <-stop:
			return
		}
	}
}

func waitProgram(program *tea.Program) error {
	_, err := program.Run()
	return err
}

func serveHTTP(server *http.Server, listener net.Listener, localProgram *tea.Program) error {
	err := server.Serve(listener)
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	localProgram.Kill()
	return err
}

func serveHandler(options ServeOptions, program *tea.Program, hub *snapshotHub) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(serveHTML)
	})
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(w, r, program, hub)
	})

	var handler http.Handler = mux
	handler = validateHost(options.Address, handler)
	if options.Username != "" {
		handler = basicAuth(options.Username, options.Password, handler)
	}
	return securityHeaders(options.Address, handler)
}

func serveWebSocket(w http.ResponseWriter, r *http.Request, program *tea.Program, hub *snapshotHub) {
	if !sameOrigin(r) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}
	updates, latest, ok := hub.reserve()
	if !ok {
		http.Error(w, "browser TUI already has an active viewer", http.StatusConflict)
		return
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin:      sameOrigin,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.release(nil)
		return
	}
	hub.attach(conn)
	defer func() {
		hub.release(conn)
		_ = conn.Close()
	}()

	configureWebSocket(conn)
	readDone := make(chan struct{})
	go readResizeLoop(conn, program, readDone)

	if latest != "" {
		if err := writeSnapshot(conn, latest); err != nil {
			return
		}
	}
	writeSnapshotLoop(conn, updates, readDone)
}

func configureWebSocket(conn *websocket.Conn) {
	conn.SetReadLimit(1024)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
}

func readResizeLoop(conn *websocket.Conn, program *tea.Program, readDone chan struct{}) {
	defer close(readDone)
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil || messageType != websocket.TextMessage {
			return
		}
		windowSize, err := parseResize(data)
		if err != nil {
			return
		}
		program.Send(windowSize)
	}
}

func writeSnapshotLoop(conn *websocket.Conn, updates chan string, readDone chan struct{}) {
	ping := time.NewTicker(30 * time.Second)
	defer ping.Stop()
	for {
		select {
		case snapshot := <-updates:
			if err := writeSnapshot(conn, snapshot); err != nil {
				return
			}
		case <-ping.C:
			_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-readDone:
			return
		}
	}
}

func parseResize(data []byte) (tea.WindowSizeMsg, error) {
	var resize resizeMessage
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&resize); err != nil || !decodeEndsAtEOF(decoder) {
		return tea.WindowSizeMsg{}, errors.New("invalid resize message")
	}
	if resize.Type != "resize" || resize.Cols < minServeCols || resize.Cols > maxServeCols ||
		resize.Rows < minServeRows || resize.Rows > maxServeRows {
		return tea.WindowSizeMsg{}, errors.New("resize dimensions out of range")
	}
	return tea.WindowSizeMsg{Width: resize.Cols, Height: resize.Rows}, nil
}

func decodeEndsAtEOF(decoder *json.Decoder) bool {
	var extra any
	return errors.Is(decoder.Decode(&extra), io.EOF)
}

func writeSnapshot(conn *websocket.Conn, snapshot string) error {
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, []byte(snapshot))
}

func sameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return false
	}
	return strings.EqualFold(u.Host, r.Host)
}

func validateHost(address string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.Host, address) {
			http.Error(w, "invalid host", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func basicAuth(username, password string, next http.Handler) http.Handler {
	wantUser := sha256.Sum256([]byte(username))
	wantPassword := sha256.Sum256([]byte(password))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		gotUser := sha256.Sum256([]byte(user))
		gotPassword := sha256.Sum256([]byte(pass))
		valid := subtle.ConstantTimeCompare(gotUser[:], wantUser[:]) &
			subtle.ConstantTimeCompare(gotPassword[:], wantPassword[:])
		if !ok || valid != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="app-listener", charset="UTF-8"`)
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func securityHeaders(address string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := fmt.Sprintf("default-src 'none'; script-src 'unsafe-inline' https://cdn.jsdelivr.net; style-src 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; font-src https://fonts.gstatic.com; connect-src 'self' ws://%s wss://%s; img-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'", address, address)
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
