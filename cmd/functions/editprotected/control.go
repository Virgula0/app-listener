package editprotected

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// ControlSocket is the daemon's local control endpoint. The daemon creates it
// 0600 root-owned and additionally checks SO_PEERCRED (uid 0) and that the
// peer's executable is the daemon's own binary before honoring a request. It
// exists only while an edit-protected password is configured.
const ControlSocket = "/run/app-listener-daemon.control"

// Control protocol (line based, UTF-8), two phases on one connection:
//
//	client -> AUTH\n<password>\n
//	server -> OK <n>\n<resource-1>\n…<resource-n>\n   (authenticated)
//	       -> ERR <reason>\n                          (refused; conn closed)
//	client -> SELECT <resource>\n
//	server -> OK\n | ERR <reason>\n                   (grant activated)
//	client -> END\n                                   (edit finished)
//
// The password is checked before any resource path is disclosed. The server
// also revokes the grant on EOF or when EditSessionMaxDuration elapses.
const (
	ctrlAuth   = "AUTH"
	ctrlSelect = "SELECT"
	ctrlEnd    = "END"
	ctrlOK     = "OK"
	ctrlErr    = "ERR"
)

// EditSessionMaxDuration is the hard cap the daemon puts on a single live
// edit window regardless of client behavior.
const EditSessionMaxDuration = 30 * time.Minute

// controlIOTimeout bounds each control-socket read/write.
const controlIOTimeout = 10 * time.Second

// liveSession is a control connection. After Authenticate it is proven; after
// Select it holds a write grant on one resource. End() releases it.
type liveSession struct {
	conn     net.Conn
	r        *bufio.Reader
	selected bool
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

// Authenticate proves the password to the daemon and returns the configured
// watch paths. A non-nil error means authentication failed (message safe to
// show).
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
	n, err := strconv.Atoi(strings.TrimSpace(rest))
	if err != nil || n < 0 {
		return nil, fmt.Errorf("bad resource count in auth response %q", line)
	}

	resources := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r, rErr := s.readLine()
		if rErr != nil {
			return nil, fmt.Errorf("reading resource list: %w", rErr)
		}
		resources = append(resources, r)
	}
	_ = s.conn.SetDeadline(time.Time{}) // the user now navigates the picker
	return resources, nil
}

// Select activates the write grant on resource for the session.
func (s *liveSession) Select(resource string) error {
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if _, err := fmt.Fprintf(s.conn, "%s %s\n", ctrlSelect, resource); err != nil {
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

// LiveModeAvailable reports whether the daemon's control socket is present —
// the precise signal that a running daemon has an edit-protected password
// configured (the socket exists only then). Used to pick live vs. offline
// without depending on systemd being the daemon's supervisor.
func LiveModeAvailable() bool {
	_, err := os.Stat(ControlSocket)
	return err == nil
}

// End releases the write grant. Best effort: the daemon also revokes on the
// connection closing.
func (s *liveSession) End() {
	if s == nil || s.conn == nil {
		return
	}
	_ = s.conn.SetDeadline(time.Now().Add(controlIOTimeout))
	if s.selected {
		_, _ = fmt.Fprintf(s.conn, "%s\n", ctrlEnd)
	}
	_ = s.conn.Close()
	s.conn = nil
}
