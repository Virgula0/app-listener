package daemon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const (
	// controlHandshakeTimeout bounds the BEGIN + password exchange.
	controlHandshakeTimeout = 10 * time.Second
	// maxAuthFailures consecutive bad handshakes trip a cooldown.
	maxAuthFailures = 5
	// authLockout is how long the socket refuses every request after the
	// failure budget is exhausted.
	authLockout = 1 * time.Minute
)

// controlServer answers the local edit-protected control socket: it
// authenticates a live edit request (peer uid 0, peer executable == this
// daemon's binary, password matches the configured hash) and, on success,
// asks the use case to widen the target resource's guard for the duration of
// the session. Only one session is served at a time.
type controlServer struct {
	uc       usecase.DaemonUseCase
	ln       net.Listener
	selfDev  uint64
	selfIno  uint64
	hashPath string

	mu        sync.Mutex
	active    *editControlSession
	failures  int
	lockUntil time.Time

	closeOnce sync.Once
}

type editControlSession struct {
	resource string
	revoke   func() error
	done     chan struct{}
	once     sync.Once
}

// end tears the session down: revoke the write grant, then unblock the
// handler goroutine. Idempotent.
func (s *editControlSession) end() {
	s.once.Do(func() {
		if err := s.revoke(); err != nil {
			log.Errorf("daemon: control: revoking edit access for %s: %v", s.resource, err)
		}
		close(s.done)
	})
}

// controlManager owns the (optional) control socket across the daemon's life:
// it is created when an edit-protected password exists and torn down when the
// password is removed. The reload handler calls refresh() so a password added
// or cleared via `edit-protected --set-password/--clear-password` on a
// running daemon takes effect on the next SIGHUP without a restart.
type controlManager struct {
	uc usecase.DaemonUseCase
	mu sync.Mutex
	cs *controlServer
}

func newControlManager(uc usecase.DaemonUseCase) *controlManager {
	return &controlManager{uc: uc}
}

// refresh brings the control socket in line with whether a password is
// configured: start it if it should run and does not, stop it otherwise.
func (m *controlManager) refresh() {
	if m == nil {
		return
	}
	exists, err := editprotected.HashFileExists()
	if err != nil {
		log.Errorf("daemon: cannot check the edit-protected password file (%v) — leaving the control socket as-is", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	switch {
	case exists && m.cs == nil:
		cs, startErr := startControlServer(m.uc)
		if startErr != nil {
			log.Errorf("daemon: edit-protected control socket unavailable (live editing disabled): %v", startErr)
			return
		}
		m.cs = cs
	case !exists && m.cs != nil:
		log.Info("daemon: edit-protected password removed — closing the control socket")
		m.cs.close()
		m.cs = nil
	}
}

func (m *controlManager) endActiveSession(reason string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	cs := m.cs
	m.mu.Unlock()
	cs.endActiveSession(reason)
}

func (m *controlManager) close() {
	if m == nil {
		return
	}
	m.mu.Lock()
	cs := m.cs
	m.cs = nil
	m.mu.Unlock()
	cs.close()
}

// startControlServer opens the control socket. The socket is 0600 and the
// handler still verifies SO_PEERCRED and the peer executable per connection.
func startControlServer(uc usecase.DaemonUseCase) (*controlServer, error) {
	selfDev, selfIno, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving own executable for the control socket: %w", err)
	}

	// A stale socket from a killed predecessor would make Listen fail.
	if rmErr := os.Remove(editprotected.ControlSocket); rmErr != nil && !os.IsNotExist(rmErr) {
		return nil, fmt.Errorf("removing stale control socket %s: %w", editprotected.ControlSocket, rmErr)
	}
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "unix", editprotected.ControlSocket)
	if err != nil {
		return nil, fmt.Errorf("listening on control socket %s: %w", editprotected.ControlSocket, err)
	}
	if chErr := os.Chmod(editprotected.ControlSocket, 0o600); chErr != nil {
		_ = ln.Close()
		_ = os.Remove(editprotected.ControlSocket)
		return nil, fmt.Errorf("securing control socket %s: %w", editprotected.ControlSocket, chErr)
	}

	cs := &controlServer{uc: uc, ln: ln, selfDev: selfDev, selfIno: selfIno, hashPath: editprotected.HashFile}
	go cs.acceptLoop()
	log.Infof("daemon: edit-protected control socket ready at %s", editprotected.ControlSocket)
	return cs, nil
}

func (cs *controlServer) acceptLoop() {
	for {
		conn, err := cs.ln.Accept()
		if err != nil {
			// Accept returns an error once the listener is closed on shutdown.
			return
		}
		go cs.handle(conn)
	}
}

// close stops accepting, ends any active session and removes the socket.
func (cs *controlServer) close() {
	if cs == nil {
		return
	}
	cs.closeOnce.Do(func() {
		_ = cs.ln.Close()
		cs.endActiveSession("daemon shutting down")
		_ = os.Remove(editprotected.ControlSocket)
	})
}

// endActiveSession revokes the current write grant, if any. Called on reload
// (guards get swapped) and on shutdown.
func (cs *controlServer) endActiveSession(reason string) {
	if cs == nil {
		return
	}
	cs.mu.Lock()
	s := cs.active
	cs.mu.Unlock()
	if s == nil {
		return
	}
	log.Warnf("daemon: control: ending the active edit session on %s (%s)", s.resource, reason)
	s.end()
}

func (cs *controlServer) handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(controlHandshakeTimeout))

	if err := cs.authPeer(conn); err != nil {
		_ = cs.reply(conn, false, err.Error())
		log.Warnf("daemon: control: rejected connection: %v", err)
		return
	}

	reader := bufio.NewReader(conn)
	resource, password, err := readHandshake(reader)
	if err != nil {
		_ = cs.reply(conn, false, "malformed request")
		log.Warnf("daemon: control: malformed request: %v", err)
		return
	}

	session, err := cs.authorize(resource, password)
	if err != nil {
		_ = cs.reply(conn, false, err.Error())
		log.Warnf("daemon: control: denied edit request for %s: %v",
			logging.SanitizeText(resource), err)
		return
	}

	if err := cs.reply(conn, true, ""); err != nil {
		// Could not confirm the grant to the client — do not leave it open.
		session.end()
		cs.clearActive(session)
		return
	}
	log.Warnf("daemon: control: live edit session OPEN for %s (uid 0, cap %s)",
		logging.SanitizeText(resource), editprotected.EditSessionMaxDuration)

	cs.runSession(conn, reader, session)
	cs.clearActive(session)
	log.Infof("daemon: control: live edit session CLOSED for %s", logging.SanitizeText(resource))
}

// runSession blocks until the client sends END, disconnects, or the hard cap
// elapses — then the session's grant is revoked (session.end, also reachable
// from reload/shutdown).
func (cs *controlServer) runSession(conn net.Conn, reader *bufio.Reader, session *editControlSession) {
	_ = conn.SetDeadline(time.Time{})
	deadline := time.NewTimer(editprotected.EditSessionMaxDuration)
	defer deadline.Stop()

	clientGone := make(chan struct{})
	go func() {
		defer close(clientGone)
		for {
			line, err := reader.ReadString('\n')
			if strings.TrimRight(line, "\r\n") == "END" {
				return
			}
			if err != nil {
				return
			}
		}
	}()

	select {
	case <-session.done: // reload or shutdown forced it
	case <-clientGone:
	case <-deadline.C:
		log.Warnf("daemon: control: edit session on %s hit the %s cap — revoking",
			logging.SanitizeText(session.resource), editprotected.EditSessionMaxDuration)
	}
	session.end()
}

func (cs *controlServer) clearActive(session *editControlSession) {
	cs.mu.Lock()
	if cs.active == session {
		cs.active = nil
	}
	cs.mu.Unlock()
}

// authPeer enforces that the connecting process runs as uid 0 and executes
// this very daemon binary (same dev/ino) — a grant on the shared exe inode is
// only useful to a caller that is that binary.
func (cs *controlServer) authPeer(conn net.Conn) error {
	uc, ok := conn.(*net.UnixConn)
	if !ok {
		return errors.New("not a unix socket")
	}
	raw, err := uc.SyscallConn()
	if err != nil {
		return fmt.Errorf("inspecting peer: %w", err)
	}
	var cred *unix.Ucred
	var credErr error
	if ctrlErr := raw.Control(func(fd uintptr) {
		cred, credErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
	}); ctrlErr != nil {
		return fmt.Errorf("inspecting peer: %w", ctrlErr)
	}
	if credErr != nil {
		return fmt.Errorf("reading peer credentials: %w", credErr)
	}
	if cred.Uid != 0 {
		return fmt.Errorf("peer uid %d is not root", cred.Uid)
	}
	dev, ino, err := ebpf.StatInode(fmt.Sprintf("/proc/%d/exe", cred.Pid))
	if err != nil {
		return fmt.Errorf("resolving peer executable: %w", err)
	}
	if dev != cs.selfDev || ino != cs.selfIno {
		return errors.New("peer is not the installed app-listener binary " +
			"(live edit-protected requires /usr/local/sbin/app-listener)")
	}
	return nil
}

// readHandshake parses the two request lines: "BEGIN <resource>" then the
// password on its own line.
func readHandshake(reader *bufio.Reader) (resource, password string, err error) {
	first, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fields := strings.SplitN(strings.TrimRight(first, "\r\n"), " ", 2)
	if len(fields) != 2 || fields[0] != "BEGIN" || strings.TrimSpace(fields[1]) == "" {
		return "", "", errors.New(`expected "BEGIN <resource-path>"`)
	}
	resource = strings.TrimSpace(fields[1])

	pwLine, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	password = strings.TrimRight(pwLine, "\r\n")
	if password == "" {
		return "", "", errors.New("empty password")
	}
	return resource, password, nil
}

// authorize runs the lockout check, verifies the password, enforces
// single-session, and asks the use case for the write grant.
func (cs *controlServer) authorize(resource, password string) (*editControlSession, error) {
	cs.mu.Lock()
	if time.Now().Before(cs.lockUntil) {
		wait := time.Until(cs.lockUntil).Round(time.Second)
		cs.mu.Unlock()
		return nil, fmt.Errorf("too many failed attempts — locked for another %s", wait)
	}
	if cs.active != nil {
		cs.mu.Unlock()
		return nil, errors.New("another live edit session is already active")
	}
	cs.mu.Unlock()

	if err := cs.checkPassword(password); err != nil {
		cs.recordFailure()
		return nil, err
	}

	revoke, err := cs.uc.GrantEditAccess(resource)
	if err != nil {
		// A bad resource path is a caller mistake, not an auth failure, but
		// it should not be free to probe either — count it.
		cs.recordFailure()
		return nil, err
	}

	session := &editControlSession{resource: resource, revoke: revoke, done: make(chan struct{})}
	cs.mu.Lock()
	if cs.active != nil { // lost a race
		cs.mu.Unlock()
		_ = revoke()
		return nil, errors.New("another live edit session is already active")
	}
	cs.active = session
	cs.failures = 0
	cs.mu.Unlock()
	return session, nil
}

func (cs *controlServer) checkPassword(password string) error {
	encoded, err := editprotected.LoadHashFile()
	if err != nil {
		return fmt.Errorf("reading the edit-protected password file: %w", err)
	}
	ok, err := editprotected.Verify(encoded, password)
	if err != nil {
		return fmt.Errorf("verifying the password: %w", err)
	}
	if !ok {
		return errors.New("authentication failed")
	}
	return nil
}

func (cs *controlServer) recordFailure() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.failures++
	if cs.failures >= maxAuthFailures {
		cs.lockUntil = time.Now().Add(authLockout)
		cs.failures = 0
		log.Warnf("daemon: control: %d failed edit-protected auth attempts — socket locked for %s", maxAuthFailures, authLockout)
	}
}

// reply writes the single-line protocol response.
func (cs *controlServer) reply(conn net.Conn, ok bool, msg string) error {
	_ = conn.SetWriteDeadline(time.Now().Add(controlHandshakeTimeout))
	if ok {
		_, err := fmt.Fprint(conn, "OK\n")
		return err
	}
	_, err := fmt.Fprintf(conn, "ERR %s\n", strings.ReplaceAll(msg, "\n", " "))
	return err
}
