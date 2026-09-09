package editprotected

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// ControlSocket is the daemon's local control endpoint. The daemon creates it
// 0600 root-owned and additionally checks SO_PEERCRED (uid 0) and that the
// peer's executable is the daemon's own binary before honoring a request. It
// exists only while an edit-protected password is configured.
const ControlSocket = "/run/app-listener-daemon.control"

// Control protocol (line based, UTF-8):
//
//	client -> BEGIN <resource-path>\n
//	client -> <password>\n
//	server -> OK\n                     (edit window open; keep the conn)
//	server -> ERR <reason>\n           (refused; conn closed)
//	client -> END\n                    (edit finished; server revokes)
//
// The server also revokes the write grant on EOF (client crash) or when
// EditSessionMaxDuration elapses, whichever comes first.
const (
	ctrlBegin = "BEGIN"
	ctrlEnd   = "END"
	ctrlOK    = "OK"
	ctrlErr   = "ERR"
)

// EditSessionMaxDuration is the hard cap the daemon puts on a single live
// edit window regardless of client behavior.
const EditSessionMaxDuration = 30 * time.Minute

// controlDialTimeout bounds the connect + handshake.
const controlDialTimeout = 10 * time.Second

// liveSession is an open control connection holding a write grant on one
// resource. End() releases it.
type liveSession struct {
	conn net.Conn
}

// beginLiveSession authenticates to the daemon and, on success, holds a
// write grant on resource open until End is called. A non-nil error means no
// grant was issued (the message is safe to show the user).
func beginLiveSession(resource, password string) (*liveSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), controlDialTimeout)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", ControlSocket)
	if err != nil {
		return nil, fmt.Errorf("connecting to the daemon control socket %s: %w "+
			"(is the daemon running with an edit-protected password configured?)", ControlSocket, err)
	}
	ok := false
	defer func() {
		if !ok {
			_ = conn.Close()
		}
	}()

	_ = conn.SetDeadline(time.Now().Add(controlDialTimeout))
	if _, wErr := fmt.Fprintf(conn, "%s %s\n%s\n", ctrlBegin, resource, password); wErr != nil {
		return nil, fmt.Errorf("sending the control request: %w", wErr)
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading the control response: %w", err)
	}
	line = strings.TrimRight(line, "\r\n")
	switch {
	case line == ctrlOK:
		_ = conn.SetDeadline(time.Time{})
		ok = true
		return &liveSession{conn: conn}, nil
	case strings.HasPrefix(line, ctrlErr+" "):
		return nil, errors.New(strings.TrimPrefix(line, ctrlErr+" "))
	default:
		return nil, fmt.Errorf("unexpected control response %q", line)
	}
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
	_ = s.conn.SetDeadline(time.Now().Add(controlDialTimeout))
	_, _ = fmt.Fprintf(s.conn, "%s\n", ctrlEnd)
	_ = s.conn.Close()
	s.conn = nil
}
