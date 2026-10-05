package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const (
	// controlHandshakeTimeout bounds the AUTH exchange. The client sends AUTH right after
	// connecting (password prompted first), so it only needs headroom for scheduling and a slow
	// hash-file read, not a human typing.
	controlHandshakeTimeout = 30 * time.Second
	// controlSelectTimeout is how long the client has, after AUTH, to pick a
	// resource and send SELECT (it is navigating the picker in that window).
	controlSelectTimeout = 5 * time.Minute
	// maxAuthFailures consecutive bad passwords trip a cooldown.
	maxAuthFailures = 5
	// authLockout is how long the socket refuses every request after the
	// failure budget is exhausted.
	authLockout = 1 * time.Minute
)

// Control protocol, two phases on one connection:
//
//	client -> AUTH\n<password>\n
//	server -> OK <n>\n<resource-1>\n...<resource-n>\n   (authenticated; the configured watch paths)
//	       -> ERR <reason>\n                          (refused; closed)
//	client -> [TIMEOUT <seconds>\n]                   (idle timeout)
//	          SELECT <resource>\n                       (edit: widen the self mask)
//	       |  FORWARD <allow|block> <events|*> <nres> <nbin>\n
//	          <resource>\n × nres  <binary>\n × nbin   (temporary -w/-b rows)
//	server -> OK\n (SELECT) | OK <m>\n<note>\n × m (FORWARD) | ERR <reason>\n
//	server -> EVENT <json>\n …                        (FORWARD: the granted binaries' events)
//	client -> PING\n …                                 (activity)
//	client -> END\n                                   (or EOF / idle timeout)
//
// The password is verified BEFORE any resource path is disclosed, so an unauthenticated caller
// learns nothing about the protected directories.
const (
	ctrlAuth    = "AUTH"
	ctrlTimeout = "TIMEOUT"
	ctrlSelect  = "SELECT"
	ctrlForward = "FORWARD"
	ctrlEvent   = "EVENT"
	ctrlPing    = "PING"
	ctrlEnd     = "END"
	// maxForwardResources / maxForwardBinaries bound one FORWARD request.
	maxForwardResources = 64
	maxForwardBinaries  = 32
)

// controlServer answers the local edit-protected control socket: authenticates a live edit request
// (peer uid 0, peer exe == this daemon's binary, password matches the hash), discloses the watch
// paths only after that, and on SELECT asks the use case to widen the target resource's guard for
// the session. One session at a time.
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
	// timeout: revoke after this long without activity (a client PING or, FORWARD, an event).
	timeout  time.Duration
	activity chan struct{}
	// events: FORWARD's stream for the client; nil for SELECT.
	events <-chan guard.GuardEvent
}

func (s *editControlSession) touch() {
	select {
	case s.activity <- struct{}{}:
	default:
	}
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

// controlManager owns the optional control socket across the daemon's life: created when an
// edit-protected password exists, torn down when removed. The reload handler calls refresh(), so
// `edit-protected --set-password/--clear-password` on a running daemon takes effect on the next
// SIGHUP.
type controlManager struct {
	uc usecase.DaemonUseCase
	mu sync.Mutex
	cs *controlServer
}

func newControlManager(uc usecase.DaemonUseCase) *controlManager {
	return &controlManager{uc: uc}
}

// startControlManager creates the control manager and immediately reconciles the socket with the
// current password state.
func startControlManager(uc usecase.DaemonUseCase) *controlManager {
	m := newControlManager(uc)
	m.refresh()
	return m
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

// sessionActive reports a live edit-protected session holding a write grant.
func (m *controlManager) sessionActive() bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	cs := m.cs
	m.mu.Unlock()
	if cs == nil {
		return false
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.active != nil
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

	// Phase 1 — AUTH. Nothing about the protected directories is disclosed
	// until the password checks out.
	password, err := readAuthRequest(reader)
	if err != nil {
		_ = cs.reply(conn, false, "malformed request")
		log.Warnf("daemon: control: malformed AUTH: %v", err)
		return
	}
	if authErr := cs.authenticate(password); authErr != nil {
		_ = cs.reply(conn, false, authErr.Error())
		log.Warnf("daemon: control: authentication refused: %v", authErr)
		return
	}
	resources := cs.resourcePaths()
	if listErr := cs.replyList(conn, resources); listErr != nil {
		return
	}

	// Phase 2 — SELECT or FORWARD; the client has been navigating a picker.
	_ = conn.SetDeadline(time.Now().Add(controlSelectTimeout))
	req, err := readGrantRequest(reader)
	if err != nil {
		_ = cs.reply(conn, false, "malformed request")
		log.Warnf("daemon: control: malformed grant request: %v", err)
		return
	}

	session, notes, err := cs.grantFor(req)
	if err != nil {
		_ = cs.reply(conn, false, err.Error())
		log.Warnf("daemon: control: denied grant for %s: %v", logging.SanitizeText(req.label()), err)
		return
	}
	if err := cs.replyGranted(conn, req, notes); err != nil {
		// Could not confirm the grant to the client — do not leave it open.
		session.end()
		cs.clearActive(session)
		return
	}
	log.Warnf("daemon: control: %s session OPEN for %s (uid 0, idle timeout %s)", req.kind(),
		logging.SanitizeText(session.resource), session.timeout)

	cs.runSession(conn, reader, session)
	cs.clearActive(session)
	log.Infof("daemon: control: %s session CLOSED for %s", req.kind(), logging.SanitizeText(session.resource))
}

// resourcePaths is the configured watch paths, disclosed only post-auth.
func (cs *controlServer) resourcePaths() []string {
	resources := cs.uc.Resources()
	out := make([]string, 0, len(resources))
	for i := range resources {
		out = append(out, resources[i].Path)
	}
	return out
}

// runSession blocks until the client sends END, disconnects, or the idle timeout passes without
// activity, then revokes the grant (session.end, also reachable from reload/shutdown).
func (cs *controlServer) runSession(conn net.Conn, reader *bufio.Reader, session *editControlSession) {
	_ = conn.SetDeadline(time.Time{})
	idle := time.NewTimer(session.timeout)
	defer idle.Stop()

	clientGone := make(chan struct{})
	go readClientLines(reader, session, clientGone)
	if session.events != nil {
		go streamEvents(conn, session)
	}

	for {
		select {
		case <-session.done: // reload, shutdown, or the use case revoked it
		case <-clientGone:
		case <-session.activity:
			idle.Reset(session.timeout)
			continue
		case <-idle.C:
			log.Warnf("daemon: control: session on %s idle for %s — revoking",
				logging.SanitizeText(session.resource), session.timeout)
		}
		session.end()
		return
	}
}

// readClientLines turns PINGs into activity and returns (closing gone) on END or EOF.
func readClientLines(reader *bufio.Reader, session *editControlSession, gone chan<- struct{}) {
	defer close(gone)
	for {
		line, err := reader.ReadString('\n')
		switch strings.TrimRight(line, "\r\n") {
		case ctrlEnd:
			return
		case ctrlPing:
			session.touch()
		}
		if err != nil {
			return
		}
	}
}

// streamEvents writes FORWARD's events to the client as EVENT lines; each one is activity. A
// client that stops reading only loses the stream: the session ends on its own terms.
func streamEvents(conn net.Conn, session *editControlSession) {
	for {
		select {
		case <-session.done:
			return
		case ev := <-session.events:
			session.touch()
			data, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			_ = conn.SetWriteDeadline(time.Now().Add(controlHandshakeTimeout))
			if _, err := fmt.Fprintf(conn, "%s %s\n", ctrlEvent, data); err != nil {
				return
			}
		}
	}
}

func (cs *controlServer) clearActive(session *editControlSession) {
	cs.mu.Lock()
	if cs.active == session {
		cs.active = nil
	}
	cs.mu.Unlock()
}

// authPeer requires the peer to run as uid 0 and execute this very daemon binary (same dev/ino): a
// grant on the shared exe inode is only useful to a caller that is that binary.
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

// readAuthRequest parses "AUTH\n<password>\n".
func readAuthRequest(reader *bufio.Reader) (password string, err error) {
	first, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if strings.TrimRight(first, "\r\n") != ctrlAuth {
		return "", fmt.Errorf("expected %q", ctrlAuth)
	}
	pwLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	password = strings.TrimRight(pwLine, "\r\n")
	if password == "" {
		return "", errors.New("empty password")
	}
	return password, nil
}

// grantRequest is phase 2: SELECT (resource) or FORWARD (forward).
type grantRequest struct {
	resource string
	forward  *forwardRequest
	timeout  time.Duration
}

type forwardRequest struct {
	resources []string
	binaries  []string
	rule      guard.TempRule
}

func (r *grantRequest) kind() string {
	if r.forward != nil {
		return "temporary access"
	}
	return "live edit"
}

func (r *grantRequest) label() string {
	if r.forward != nil {
		return strings.Join(r.forward.resources, ", ")
	}
	return r.resource
}

// readGrantRequest parses an optional "TIMEOUT <seconds>" line, then "SELECT <resource-path>" or a
// FORWARD block.
func readGrantRequest(reader *bufio.Reader) (*grantRequest, error) {
	req := &grantRequest{timeout: editprotected.SessionTimeoutDefault}
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if rest, ok := strings.CutPrefix(strings.TrimRight(line, "\r\n"), ctrlTimeout+" "); ok {
		secs, perr := strconv.Atoi(rest)
		if perr != nil || secs < 1 || time.Duration(secs)*time.Second > editprotected.SessionTimeoutMax {
			return nil, fmt.Errorf("timeout must be 1..%d seconds", int(editprotected.SessionTimeoutMax/time.Second))
		}
		req.timeout = time.Duration(secs) * time.Second
		if line, err = reader.ReadString('\n'); err != nil {
			return nil, err
		}
	}
	fields := strings.SplitN(strings.TrimRight(line, "\r\n"), " ", 2)
	switch fields[0] {
	case ctrlSelect:
		if len(fields) != 2 || strings.TrimSpace(fields[1]) == "" {
			return nil, fmt.Errorf("expected %q <resource-path>", ctrlSelect)
		}
		req.resource = strings.TrimSpace(fields[1])
		return req, nil
	case ctrlForward:
		if len(fields) != 2 {
			return nil, fmt.Errorf("expected %q <allow|block> <events|*> <nres> <nbin>", ctrlForward)
		}
		if req.forward, err = readForwardRequest(reader, fields[1]); err != nil {
			return nil, err
		}
		return req, nil
	}
	return nil, fmt.Errorf("expected %q or %q", ctrlSelect, ctrlForward)
}

func readForwardRequest(reader *bufio.Reader, header string) (*forwardRequest, error) {
	f := strings.Fields(header)
	if len(f) != 4 {
		return nil, fmt.Errorf("expected %q <allow|block> <events|*> <nres> <nbin>", ctrlForward)
	}
	fwd := &forwardRequest{}
	switch f[0] {
	case "allow":
		fwd.rule.Allow = true
	case "block":
	default:
		return nil, fmt.Errorf("unknown FORWARD action %q", f[0])
	}
	if f[1] != "*" {
		for _, name := range strings.Split(f[1], ",") {
			et, ok := ebpf.ParseEventType(name)
			if !ok {
				return nil, fmt.Errorf("unknown event type %q", name)
			}
			fwd.rule.Events = append(fwd.rule.Events, et)
		}
	}
	nres, err := strconv.Atoi(f[2])
	if err != nil || nres < 1 || nres > maxForwardResources {
		return nil, fmt.Errorf("resource count must be 1..%d", maxForwardResources)
	}
	nbin, err := strconv.Atoi(f[3])
	if err != nil || nbin < 1 || nbin > maxForwardBinaries {
		return nil, fmt.Errorf("binary count must be 1..%d", maxForwardBinaries)
	}
	if fwd.resources, err = readPathLines(reader, nres); err != nil {
		return nil, err
	}
	if fwd.binaries, err = readPathLines(reader, nbin); err != nil {
		return nil, err
	}
	return fwd, nil
}

func readPathLines(reader *bufio.Reader, n int) ([]string, error) {
	out := make([]string, 0, n)
	for range n {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		p := strings.TrimRight(line, "\r\n")
		if !filepath.IsAbs(p) {
			return nil, fmt.Errorf("%q is not an absolute path", p)
		}
		out = append(out, p)
	}
	return out, nil
}

// authenticate runs the lockout check and verifies the password. On failure
// it records a strike toward the cooldown; on success it clears the counter.
func (cs *controlServer) authenticate(password string) error {
	cs.mu.Lock()
	if time.Now().Before(cs.lockUntil) {
		wait := time.Until(cs.lockUntil).Round(time.Second)
		cs.mu.Unlock()
		return fmt.Errorf("too many failed attempts — locked for another %s", wait)
	}
	cs.mu.Unlock()

	if err := cs.checkPassword(password); err != nil {
		cs.recordFailure()
		return err
	}
	cs.mu.Lock()
	cs.failures = 0
	cs.mu.Unlock()
	return nil
}

// grantFor enforces single-session and asks the use case for the grant. The caller is already
// authenticated. notes are FORWARD's per-binary remarks.
func (cs *controlServer) grantFor(req *grantRequest) (*editControlSession, []string, error) {
	cs.mu.Lock()
	busy := cs.active != nil
	cs.mu.Unlock()
	if busy {
		return nil, nil, errControlBusy
	}

	if req.forward == nil {
		revoke, err := cs.uc.GrantEditAccess(req.resource)
		if err != nil {
			return nil, nil, err
		}
		session, err := cs.activate(req, revoke, nil)
		return session, nil, err
	}

	fwd := req.forward
	access, err := cs.uc.GrantTemporaryAccess(fwd.resources, fwd.binaries, fwd.rule, writeTempJournal)
	if err != nil {
		return nil, nil, err
	}
	session, err := cs.activate(req, access.Revoke, access.Events)
	if err != nil {
		return nil, nil, err
	}
	// The use case revokes on its own when a granted binary changes in place: end the session too.
	go func() {
		select {
		case <-access.Done:
			session.end()
		case <-session.done:
		}
	}()
	return session, access.Report, nil
}

var errControlBusy = errors.New("another live edit session is already active")

func (cs *controlServer) activate(req *grantRequest, revoke func() error,
	events <-chan guard.GuardEvent) (*editControlSession, error) {
	session := &editControlSession{resource: req.label(), revoke: revoke, done: make(chan struct{}),
		timeout: req.timeout, activity: make(chan struct{}, 1), events: events}
	cs.mu.Lock()
	if cs.active != nil { // lost a race
		cs.mu.Unlock()
		_ = revoke()
		return nil, errControlBusy
	}
	cs.active = session
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

// replyGranted confirms a grant: "OK" for SELECT, "OK <m>" plus m note lines for FORWARD.
func (cs *controlServer) replyGranted(conn net.Conn, req *grantRequest, notes []string) error {
	if req.forward == nil {
		return cs.reply(conn, true, "")
	}
	return cs.replyList(conn, notes)
}

// replyList writes "OK <n>" and n lines: the watch paths after AUTH, FORWARD's notes.
func (cs *controlServer) replyList(conn net.Conn, lines []string) error {
	_ = conn.SetWriteDeadline(time.Now().Add(controlHandshakeTimeout))
	var b strings.Builder
	fmt.Fprintf(&b, "OK %d\n", len(lines))
	for _, r := range lines {
		b.WriteString(strings.ReplaceAll(r, "\n", " "))
		b.WriteByte('\n')
	}
	_, err := io.WriteString(conn, b.String())
	return err
}
