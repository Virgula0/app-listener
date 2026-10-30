// Package secretgate is the Secret Service gatekeeper: the keyring daemon (the provider, gnome-keyring)
// reaches the user's session bus only through a relay the daemon runs as root, and every method call
// from the bus to the keyring is judged there before the keyring sees it. Sessions, prompts and
// signals pass through unchanged: the gate filters one connection's stream, it never re-encrypts or
// re-exports. Items belong to [secret] scopes by their attributes; a call touching a scoped item must
// come from a process its whitelist admits (attested by the pidfd the bus took at connect,
// guard.AttestBusClient). Unscoped items stay as open as the keyring keeps them.
package secretgate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Role is a [secret_service] role beyond item scopes.
type Role int

const (
	RoleManager Role = iota
	RolePrompter
)

// Caller is an attested bus client, opaque to the gate.
type Caller any

// Scope is a [secret] scope as the gate matches items against it.
type Scope struct {
	Name    string
	Matches []daemonconfig.SecretMatch
}

// Judge is the daemon's side of every decision. Admits and InRole log the decision.
type Judge interface {
	// Attest resolves the pidfd a caller's bus reported; the gate keeps ownership of pidfd.
	Attest(pidfd int) (Caller, error)
	Scopes() []Scope
	Admits(scope string, c Caller, op ebpf.EventType, what string) bool
	InRole(r Role, c Caller, what string) bool
	// RelayGuarded reports a live guard on the relay sockets: only it narrows who connects.
	RelayGuarded() bool
}

const (
	callTimeout = 10 * time.Second
	// maxQueued bounds the calls waiting behind one caller's pending decision: beyond it, refused.
	maxQueued = 256
	// injectedSerial starts the serials of the errors the gate sends as the keyring: far above any
	// the keyring reaches itself.
	injectedSerial = 0x80000000
)

// Gate serves one user's relay socket.
type Gate struct {
	uid   uint32
	bus   string
	judge Judge
	ln    *net.UnixListener

	mu      sync.Mutex
	current *session
	closed  bool
}

// Listen binds the relay socket at path for uid's provider (only uid may connect; the guard narrows it
// to the provider binary) and relays to the bus socket at bus.
func Listen(uid, gid uint32, path, bus string, judge Judge) (*Gate, error) {
	ln, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		return nil, err
	}
	ln.SetUnlinkOnClose(true)
	if err := os.Chown(path, int(uid), int(gid)); err == nil {
		err = os.Chmod(path, 0o600)
	}
	if err != nil {
		ln.Close()
		return nil, err
	}
	return &Gate{uid: uid, bus: bus, judge: judge, ln: ln}, nil
}

// Serve accepts the provider's connections, one live at a time, until Close.
func (g *Gate) Serve() {
	for {
		c, err := g.ln.AcceptUnix()
		if err != nil {
			g.mu.Lock()
			closed := g.closed
			g.mu.Unlock()
			if !closed {
				log.Errorf("secret gate %d: accept: %v", g.uid, err)
			}
			return
		}
		if !g.judge.RelayGuarded() {
			log.Errorf("secret gate %d: refusing a keyring connection: no guard holds the relay socket", g.uid)
			c.Close()
			continue
		}
		g.mu.Lock()
		busy := g.current != nil || g.closed
		var s *session
		if !busy {
			s = &session{g: g, down: c, queues: map[string]*senderQueue{}, pidfds: map[string]int{}}
			s.serial.Store(injectedSerial)
			g.current = s
		}
		g.mu.Unlock()
		if busy {
			log.Warnf("secret gate %d: refusing a second keyring connection", g.uid)
			c.Close()
			continue
		}
		go func() {
			if err := s.run(); err != nil {
				log.Warnf("secret gate %d: keyring session ended: %v", g.uid, err)
			}
			g.mu.Lock()
			if g.current == s {
				g.current = nil
			}
			g.mu.Unlock()
		}()
	}
}

// Close stops accepting and ends the live session.
func (g *Gate) Close() {
	g.mu.Lock()
	g.closed = true
	s := g.current
	g.mu.Unlock()
	g.ln.Close()
	if s != nil {
		s.close()
	}
}

type session struct {
	g      *Gate
	down   *net.UnixConn
	downR  *bufio.Reader
	downMu sync.Mutex
	up     net.Conn
	upR    *bufio.Reader
	upMu   sync.Mutex
	ctl    *dbus.Conn
	serial atomic.Uint32

	mu     sync.Mutex
	queues map[string]*senderQueue
	pidfds map[string]int
	ended  bool
}

func (s *session) run() error {
	defer s.close()
	ucred, err := peerCred(s.down)
	if err != nil {
		return err
	}
	if ucred.Uid != s.g.uid {
		return fmt.Errorf("peer uid %d is not %d", ucred.Uid, s.g.uid)
	}
	if s.downR, err = serverAuth(s.down, s.g.uid); err != nil {
		return err
	}
	if s.up, err = net.Dial("unix", s.g.bus); err != nil {
		return err
	}
	if s.upR, err = clientAuth(s.up, uint32(os.Geteuid())); err != nil { //nolint:gosec // euid >= 0
		return err
	}
	if err := s.dialControl(); err != nil {
		return err
	}
	log.Infof("secret gate %d: keyring (pid %d) relayed to %s", s.g.uid, ucred.Pid, s.g.bus)
	errc := make(chan error, 2)
	go func() { errc <- s.pumpUp() }()
	go func() { errc <- s.pumpDown() }()
	return <-errc
}

// dialControl opens the gate's own bus connection: credentials of callers, item attributes, and
// NameOwnerChanged to drop a gone caller's pidfd.
func (s *session) dialControl() error {
	ctl, err := dbus.Dial("unix:path=" + s.g.bus)
	if err != nil {
		return err
	}
	if err := ctl.Auth(nil); err == nil {
		err = ctl.Hello()
	}
	if err == nil {
		err = ctl.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.DBus"),
			dbus.WithMatchMember("NameOwnerChanged"), dbus.WithMatchSender("org.freedesktop.DBus"))
	}
	if err != nil {
		ctl.Close()
		return err
	}
	s.ctl = ctl
	sigs := make(chan *dbus.Signal, 64)
	ctl.Signal(sigs)
	go func() {
		for sig := range sigs {
			if len(sig.Body) == 3 {
				name, _ := sig.Body[0].(string)
				owner, _ := sig.Body[2].(string)
				if owner == "" && len(name) > 0 && name[0] == ':' {
					s.forget(name)
				}
			}
		}
	}()
	return nil
}

func (s *session) close() {
	s.mu.Lock()
	if s.ended {
		s.mu.Unlock()
		return
	}
	s.ended = true
	for name, q := range s.queues {
		q.stop()
		delete(s.queues, name)
	}
	for name, fd := range s.pidfds {
		unix.Close(fd)
		delete(s.pidfds, name)
	}
	s.mu.Unlock()
	s.down.Close()
	if s.up != nil {
		s.up.Close()
	}
	if s.ctl != nil {
		s.ctl.Close()
	}
}

// pumpUp passes the keyring's own traffic to the bus unchanged.
func (s *session) pumpUp() error {
	for {
		frame, err := readFrame(s.downR)
		if err != nil {
			return fmt.Errorf("from keyring: %w", err)
		}
		if err := s.writeUp(frame); err != nil {
			return err
		}
	}
}

// pumpDown routes the bus's traffic to the keyring: method calls through their sender's queue
// (judged in order), everything else straight through. It never blocks on a decision.
func (s *session) pumpDown() error {
	for {
		frame, err := readFrame(s.upR)
		if err != nil {
			return fmt.Errorf("from bus: %w", err)
		}
		if frame[1] != byte(dbus.TypeMethodCall) {
			if err := s.writeDown(frame); err != nil {
				return err
			}
			continue
		}
		m, err := decodeFrame(frame)
		if err != nil {
			log.Warnf("secret gate %d: dropping an undecodable call: %v", s.g.uid, err)
			continue
		}
		s.enqueue(header[string](m, dbus.FieldSender), m, frame)
	}
}

func (s *session) writeUp(frame []byte) error {
	s.upMu.Lock()
	defer s.upMu.Unlock()
	_, err := s.up.Write(frame)
	return err
}

func (s *session) writeDown(frame []byte) error {
	s.downMu.Lock()
	defer s.downMu.Unlock()
	_, err := s.down.Write(frame)
	return err
}

// senderQueue judges one caller's calls in order, so a held call is never overtaken by the next.
type senderQueue struct {
	mu      sync.Mutex
	cond    *sync.Cond
	pending []queued
	stopped bool
}

type queued struct {
	m     *dbus.Message
	frame []byte
}

func (q *senderQueue) stop() {
	q.mu.Lock()
	q.stopped = true
	q.mu.Unlock()
	q.cond.Broadcast()
}

func (s *session) enqueue(sender string, m *dbus.Message, frame []byte) {
	s.mu.Lock()
	if s.ended {
		s.mu.Unlock()
		return
	}
	q, ok := s.queues[sender]
	if !ok {
		q = &senderQueue{}
		q.cond = sync.NewCond(&q.mu)
		s.queues[sender] = q
		go s.drain(sender, q)
	}
	s.mu.Unlock()
	q.mu.Lock()
	full := len(q.pending) >= maxQueued
	if !full {
		q.pending = append(q.pending, queued{m, frame})
	}
	q.mu.Unlock()
	if full {
		s.deny(m, "too many pending calls")
		return
	}
	q.cond.Signal()
}

func (s *session) drain(sender string, q *senderQueue) {
	for {
		q.mu.Lock()
		for len(q.pending) == 0 && !q.stopped {
			q.cond.Wait()
		}
		if q.stopped {
			q.mu.Unlock()
			return
		}
		next := q.pending[0]
		q.pending = q.pending[1:]
		q.mu.Unlock()
		s.handle(sender, next.m, next.frame)
	}
}

// forget drops a caller that left the bus: its pidfd and its queue.
func (s *session) forget(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if fd, ok := s.pidfds[name]; ok {
		unix.Close(fd)
		delete(s.pidfds, name)
	}
	if q, ok := s.queues[name]; ok {
		q.stop()
		delete(s.queues, name)
	}
}

func (s *session) handle(sender string, m *dbus.Message, frame []byte) {
	v := classify(m)
	what := fmt.Sprintf("secret-service %s.%s %s", header[string](m, dbus.FieldInterface),
		header[string](m, dbus.FieldMember), header[dbus.ObjectPath](m, dbus.FieldPath))
	ok := false
	switch v.kind {
	case allow:
		ok = true
	case needManager, needPrompter:
		role := RoleManager
		if v.kind == needPrompter {
			role = RolePrompter
		}
		if c, err := s.attest(sender); err == nil {
			ok = s.g.judge.InRole(role, c, what)
		} else {
			log.Warnf("secret gate %d: refused %s from %s: %v", s.g.uid, what, sender, err)
		}
	case gated:
		ok = s.admitGated(sender, v, what)
	default:
		log.Warnf("secret gate %d: refused %s from %s: not a call the gate relays", s.g.uid, what, sender)
	}
	if !ok {
		s.deny(m, "refused by app-listener")
		return
	}
	if err := s.writeDown(frame); err != nil {
		log.Debugf("secret gate %d: forwarding: %v", s.g.uid, err)
	}
}

// admitGated: every scope of every touched item (its current attributes, and those the call would
// give it) must admit the caller. No scope: the item is not the gate's to protect. A lookup that fails
// refuses (fail closed): the keyring would answer the call for an item the gate couldn't place.
func (s *session) admitGated(sender string, v verdict, what string) bool {
	scopes := s.g.judge.Scopes()
	var hit []string
	add := func(attrs map[string]string) {
		for _, sc := range scopes {
			if !slices.Contains(hit, sc.Name) && slices.ContainsFunc(sc.Matches, func(m daemonconfig.SecretMatch) bool {
				return m.Matches(attrs)
			}) {
				hit = append(hit, sc.Name)
			}
		}
	}
	for _, a := range v.attrs {
		add(a)
	}
	for _, item := range v.items {
		attrs, err := s.attributes(item)
		if err != nil {
			log.Warnf("secret gate %d: refused %s from %s: attributes of %s: %v", s.g.uid, what, sender, item, err)
			return false
		}
		add(attrs)
	}
	if len(hit) == 0 {
		return true
	}
	c, err := s.attest(sender)
	if err != nil {
		log.Warnf("secret gate %d: refused %s from %s (scopes %v): %v", s.g.uid, what, sender, hit, err)
		return false
	}
	for _, name := range hit {
		if !s.g.judge.Admits(name, c, v.op, what) {
			return false
		}
	}
	return true
}

func (s *session) attributes(item dbus.ObjectPath) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()
	var attrs map[string]string
	err := s.ctl.Object("org.freedesktop.secrets", item).CallWithContext(ctx, ifaceProps+".Get", 0,
		ifaceItem, "Attributes").Store(&attrs)
	return attrs, err
}

// attest resolves sender through the pidfd its bus took at connect (cached for the connection).
func (s *session) attest(sender string) (Caller, error) {
	s.mu.Lock()
	fd, ok := s.pidfds[sender]
	s.mu.Unlock()
	if !ok {
		var err error
		if fd, err = s.processFD(sender); err != nil {
			return nil, err
		}
		s.mu.Lock()
		if s.ended {
			s.mu.Unlock()
			unix.Close(fd)
			return nil, errors.New("session ended")
		}
		if old, dup := s.pidfds[sender]; dup {
			unix.Close(fd)
			fd = old
		} else {
			s.pidfds[sender] = fd
		}
		s.mu.Unlock()
	}
	return s.g.judge.Attest(fd)
}

func (s *session) processFD(sender string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()
	var creds map[string]dbus.Variant
	if err := s.ctl.BusObject().CallWithContext(ctx, "org.freedesktop.DBus.GetConnectionCredentials", 0,
		sender).Store(&creds); err != nil {
		return -1, err
	}
	v, ok := creds["ProcessFD"]
	if !ok {
		return -1, errors.New("the bus reports no ProcessFD (dbus-broker 35+ or dbus-daemon 1.15.6+ needed)")
	}
	fd, ok := v.Value().(dbus.UnixFD)
	if !ok {
		return -1, fmt.Errorf("ProcessFD has type %T", v.Value())
	}
	return int(fd), nil
}

// deny answers m with AccessDenied as if from the keyring (the bus routes it to the caller by
// reply serial). A call expecting no reply is dropped.
func (s *session) deny(m *dbus.Message, text string) {
	if m.Flags&dbus.FlagNoReplyExpected != 0 {
		return
	}
	e := &dbus.Message{
		Type:  dbus.TypeError,
		Flags: dbus.FlagNoReplyExpected,
		Headers: map[dbus.HeaderField]dbus.Variant{
			dbus.FieldReplySerial: dbus.MakeVariant(m.Serial()),
			dbus.FieldDestination: dbus.MakeVariant(header[string](m, dbus.FieldSender)),
			dbus.FieldErrorName:   dbus.MakeVariant("org.freedesktop.DBus.Error.AccessDenied"),
			dbus.FieldSignature:   dbus.MakeVariant(dbus.SignatureOf(text)),
		},
		Body: []any{text},
	}
	frame, err := encodeWithSerial(e, s.serial.Add(1))
	if err == nil {
		err = s.writeUp(frame)
	}
	if err != nil {
		log.Debugf("secret gate %d: answering a refused call: %v", s.g.uid, err)
	}
}

func peerCred(c *net.UnixConn) (*unix.Ucred, error) {
	raw, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}
	var cred *unix.Ucred
	var credErr error
	if err := raw.Control(func(fd uintptr) {
		cred, credErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
	}); err != nil {
		return nil, err
	}
	return cred, credErr
}
