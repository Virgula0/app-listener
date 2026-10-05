package editprotected

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Virgula0/app-listener/internal/guard"
)

// ControlSocket is the daemon's local control endpoint, present only while an edit-protected
// password is configured. Created 0600 root-owned; the daemon also checks SO_PEERCRED (uid 0) and
// that the peer's exe is its own binary before honoring a request.
const ControlSocket = "/run/app-listener-daemon.control"

// Control protocol (line based, UTF-8), two phases on one connection:
//
//	client -> AUTH\n<password>\n
//	server -> OK <n>\n<resource-1>\n...<resource-n>\n   (authenticated)
//	       -> ERR <reason>\n                          (refused; conn closed)
//	client -> [TIMEOUT <seconds>\n]                   (idle timeout; default SessionTimeoutDefault)
//	          SELECT <resource>\n                       (edit)
//	       |  FORWARD <allow|block> <events|*> <nres> <nbin>\n
//	          <resource>\n × nres  <binary>\n × nbin   (--forward)
//	server -> OK\n (SELECT) | OK <m>\n<note>\n × m (FORWARD) | ERR <reason>\n
//	server -> EVENT <json guard.GuardEvent>\n …          (FORWARD: the granted binaries' events)
//	client -> PING\n …                                 (activity: resets the idle timeout)
//	client -> END\n                                   (session finished)
//
// or a configuration edit: CONFIG -> OK <n>\n<daemon.conf>; then END, or PUT <m>\n<bytes> -> OK |
// ERR (the daemon reloads it, and restores the file if the reload fails).
//
// The password is checked before any resource path is disclosed. The server revokes the grant on
// EOF or after the idle timeout with neither a PING nor (FORWARD) a granted binary's event, and
// closes the connection when it revokes on its own.
const (
	ctrlAuth    = "AUTH"
	ctrlTimeout = "TIMEOUT"
	ctrlSelect  = "SELECT"
	ctrlForward = "FORWARD"
	ctrlEvent   = "EVENT"
	ctrlPing    = "PING"
	ctrlEnd     = "END"
	ctrlConfig  = "CONFIG"
	ctrlPut     = "PUT"
	ctrlOK      = "OK"
	ctrlErr     = "ERR"
)

// SessionTimeoutDefault is a live session's idle timeout when --timeout-session isn't given;
// SessionTimeoutMax bounds what the daemon accepts.
const (
	SessionTimeoutDefault = 30 * time.Minute
	SessionTimeoutMax     = 24 * time.Hour
)

// configReloadTimeout bounds the wait for a PUT's reload: each guard attach is a verifier pass.
const configReloadTimeout = 5 * time.Minute

// controlIOTimeout bounds each control-socket read/write.
const controlIOTimeout = 10 * time.Second

// liveSession is a control connection. After Authenticate it is proven; after
// Select it holds a write grant on one resource. End() releases it.
type liveSession struct {
	conn     net.Conn
	r        *bufio.Reader
	selected bool
	// timeout is sent before SELECT/FORWARD (0 = the daemon's default).
	timeout time.Duration
	// wmu serializes writes: PINGs come from the TUI goroutine.
	wmu sync.Mutex
}

// dialLiveSession connects to the daemon control socket. It sends nothing —
// call Authenticate next.
func dialLiveSession() (*liveSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), controlIOTimeout)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", ControlSocket)
	if err != nil {
		return nil, fmt.Errorf("connecting to the daemon control socket %s: %w "+
			"(is the daemon running with an edit-protected password configured?)", ControlSocket, err)
	}
	return &liveSession{conn: conn, r: bufio.NewReader(conn)}, nil
}

// Authenticate proves the password to the daemon and returns the configured watch paths. A non-nil
// error means authentication failed (message safe to show).
func (s *liveSession) Authenticate(password string) ([]string, error) {
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if _, err := fmt.Fprintf(s.conn, "%s\n%s\n", ctrlAuth, password); err != nil {
		return nil, fmt.Errorf("sending the auth request: %w", err)
	}

	line, err := s.readLine()
	if err != nil {
		return nil, fmt.Errorf("reading the auth response: %w", err)
	}
	if msg, isErr := parseErr(line); isErr {
		return nil, errors.New(msg)
	}
	rest, ok := strings.CutPrefix(line, ctrlOK+" ")
	if !ok {
		return nil, fmt.Errorf("unexpected auth response %q", line)
	}
	resources, err := s.readList(rest)
	if err != nil {
		return nil, fmt.Errorf("reading resource list: %w", err)
	}
	_ = s.conn.SetDeadline(time.Time{}) // the user now navigates the picker
	return resources, nil
}

// readList reads the n lines announced by an "OK <n>" response (count is what followed "OK ").
func (s *liveSession) readList(count string) ([]string, error) {
	n, err := strconv.Atoi(strings.TrimSpace(count))
	if err != nil || n < 0 {
		return nil, fmt.Errorf("bad line count %q", count)
	}
	out := make([]string, 0, n)
	for range n {
		l, err := s.readLine()
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

// Forward activates a temporary rule (allow = -w, else -b; events nil = all) for binaries on
// resources, returning the daemon's notes on binaries it left untouched.
func (s *liveSession) Forward(resources, binaries []string, allow bool, events []string) ([]string, error) {
	action := "block"
	if allow {
		action = "allow"
	}
	evs := "*"
	if len(events) > 0 {
		evs = strings.Join(events, ",")
	}
	var b strings.Builder
	b.WriteString(s.timeoutLine())
	fmt.Fprintf(&b, "%s %s %s %d %d\n", ctrlForward, action, evs, len(resources), len(binaries))
	for _, p := range append(slices.Clone(resources), binaries...) {
		if strings.ContainsAny(p, "\r\n") {
			return nil, fmt.Errorf("path %q contains a line break", p)
		}
		b.WriteString(p)
		b.WriteByte('\n')
	}

	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if _, err := io.WriteString(s.conn, b.String()); err != nil {
		return nil, fmt.Errorf("sending the forward request: %w", err)
	}
	line, err := s.readLine()
	if err != nil {
		return nil, fmt.Errorf("reading the forward response: %w", err)
	}
	if msg, isErr := parseErr(line); isErr {
		return nil, errors.New(msg)
	}
	rest, ok := strings.CutPrefix(line, ctrlOK+" ")
	if !ok {
		return nil, fmt.Errorf("unexpected forward response %q", line)
	}
	notes, err := s.readList(rest)
	if err != nil {
		return nil, fmt.Errorf("reading the forward notes: %w", err)
	}
	_ = s.conn.SetDeadline(time.Time{})
	s.selected = true
	return notes, nil
}

// Select activates the write grant on resource for the session.
func (s *liveSession) Select(resource string) error {
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if _, err := fmt.Fprintf(s.conn, "%s%s %s\n", s.timeoutLine(), ctrlSelect, resource); err != nil {
		return fmt.Errorf("sending the select request: %w", err)
	}
	line, err := s.readLine()
	if err != nil {
		return fmt.Errorf("reading the select response: %w", err)
	}
	if line == ctrlOK {
		_ = s.conn.SetDeadline(time.Time{})
		s.selected = true
		return nil
	}
	if msg, isErr := parseErr(line); isErr {
		return errors.New(msg)
	}
	return fmt.Errorf("unexpected select response %q", line)
}

// Config fetches the daemon's daemon.conf. The daemon then waits for PutConfig or End.
func (s *liveSession) Config() ([]byte, error) {
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if _, err := fmt.Fprintf(s.conn, "%s\n", ctrlConfig); err != nil {
		return nil, fmt.Errorf("sending the config request: %w", err)
	}
	line, err := s.readLine()
	if err != nil {
		return nil, fmt.Errorf("reading the config response: %w", err)
	}
	if msg, isErr := parseErr(line); isErr {
		return nil, errors.New(msg)
	}
	rest, ok := strings.CutPrefix(line, ctrlOK+" ")
	n, err := strconv.Atoi(rest)
	if !ok || err != nil || n < 0 {
		return nil, fmt.Errorf("unexpected config response %q", line)
	}
	body := make([]byte, n)
	if _, err := io.ReadFull(s.r, body); err != nil {
		return nil, fmt.Errorf("reading the configuration: %w", err)
	}
	s.selected = true // End tells the daemon nothing is coming
	return body, nil
}

// PutConfig sends the new daemon.conf; a nil error means the daemon reloaded it. On a refusal
// the daemon's file and running configuration are unchanged.
func (s *liveSession) PutConfig(next []byte) error {
	// The daemon reloads before it answers: every guard is rebuilt.
	_ = s.conn.SetDeadline(time.Now().Add(configReloadTimeout))
	if _, err := fmt.Fprintf(s.conn, "%s %d\n%s", ctrlPut, len(next), next); err != nil {
		return fmt.Errorf("sending the configuration: %w", err)
	}
	s.selected = false
	line, err := s.readLine()
	if err != nil {
		return fmt.Errorf("reading the daemon's answer: %w", err)
	}
	if line == ctrlOK {
		return nil
	}
	if msg, isErr := parseErr(line); isErr {
		return errors.New(msg)
	}
	return fmt.Errorf("unexpected answer %q", line)
}

func (s *liveSession) timeoutLine() string {
	if s.timeout <= 0 {
		return ""
	}
	return fmt.Sprintf("%s %d\n", ctrlTimeout, int(s.timeout/time.Second))
}

// Watch reads the daemon's side of an active session: events streams FORWARD's EVENT lines, and
// ended closes once the daemon closes the connection (it revoked the grant: idle timeout, reload,
// a granted binary changed). Call it once, after Select/Forward.
func (s *liveSession) Watch() (events <-chan guard.GuardEvent, ended <-chan struct{}) {
	evs := make(chan guard.GuardEvent, 64)
	done := make(chan struct{})
	r := s.r
	go func() {
		defer close(done)
		for {
			line, err := r.ReadString('\n')
			if payload, ok := strings.CutPrefix(strings.TrimRight(line, "\r\n"), ctrlEvent+" "); ok {
				var ev guard.GuardEvent
				if json.Unmarshal([]byte(payload), &ev) == nil {
					select {
					case evs <- ev:
					default: // a slow display drops, never stalls the reader
					}
				}
			}
			if err != nil {
				return
			}
		}
	}()
	return evs, done
}

// Pinger returns an activity callback sending PING at most once per interval, so the daemon's
// idle timeout counts from the operator's last action.
func (s *liveSession) Pinger(interval time.Duration) func() {
	var mu sync.Mutex
	var last time.Time
	return func() {
		mu.Lock()
		if time.Since(last) < interval {
			mu.Unlock()
			return
		}
		last = time.Now()
		mu.Unlock()
		s.wmu.Lock()
		defer s.wmu.Unlock()
		if s.conn != nil {
			_ = s.conn.SetWriteDeadline(time.Now().Add(controlIOTimeout))
			_, _ = fmt.Fprintf(s.conn, "%s\n", ctrlPing)
		}
	}
}

func (s *liveSession) readLine() (string, error) {
	line, err := s.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func parseErr(line string) (msg string, isErr bool) {
	if line == ctrlErr {
		return "refused", true
	}
	if rest, ok := strings.CutPrefix(line, ctrlErr+" "); ok {
		return rest, true
	}
	return "", false
}

// LiveModeAvailable reports whether the daemon's control socket exists: the precise signal that a
// running daemon has an edit-protected password (works without systemd supervising it). Picks live
// vs offline.
func LiveModeAvailable() bool {
	_, err := os.Stat(ControlSocket)
	return err == nil
}

// End releases the write grant. Best effort: the daemon also revokes on the
// connection closing.
func (s *liveSession) End() {
	if s == nil {
		return
	}
	s.wmu.Lock()
	defer s.wmu.Unlock()
	if s.conn == nil {
		return
	}
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if s.selected {
		_, _ = fmt.Fprintf(s.conn, "%s\n", ctrlEnd)
	}
	_ = s.conn.Close()
	s.conn = nil
}
