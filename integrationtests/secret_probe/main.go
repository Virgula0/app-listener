// secret_probe is a Secret Service client for the integration tests: a whitelisted copy plays the
// app that owns an item, other copies try to reach it. Exit codes: 0 = got the secret, 1 = refused
// (AccessDenied or a refused connect), 3 = setup error, 4 = nothing captured (own: not the owner).
package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
	"golang.org/x/sys/unix"
)

const (
	service = "org.freedesktop.secrets"
	svcPath = dbus.ObjectPath("/org/freedesktop/secrets")
	svcIfc  = "org.freedesktop.Secret.Service"
)

type secret struct {
	Session     dbus.ObjectPath
	Params      []byte
	Value       []byte
	ContentType string
}

func failf(code int, format string, a ...any) {
	fmt.Printf(format+"\n", a...)
	os.Exit(code)
}

// refusal maps a D-Bus error to exit 1 when the gatekeeper or the bus policy refused it.
func refusal(what string, err error) {
	var de dbus.Error
	if errors.As(err, &de) && (de.Name == "org.freedesktop.DBus.Error.AccessDenied" ||
		strings.Contains(de.Name, "NotAuthorized") || strings.Contains(de.Name, "IsDenied")) {
		failf(1, "REFUSED %s: %v", what, err)
	}
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
		failf(1, "REFUSED %s: %v", what, err)
	}
	failf(3, "ERROR %s: %v", what, err)
}

func busPath() string {
	addr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	p, ok := strings.CutPrefix(addr, "unix:path=")
	if !ok {
		failf(3, "ERROR unsupported DBUS_SESSION_BUS_ADDRESS %q", addr)
	}
	if i := strings.IndexByte(p, ','); i >= 0 {
		p = p[:i]
	}
	return p
}

func dial() *dbus.Conn {
	c, err := dbus.Dial(os.Getenv("DBUS_SESSION_BUS_ADDRESS"))
	if err != nil {
		refusal("dial", err)
	}
	return handshake(c)
}

func handshake(c *dbus.Conn) *dbus.Conn {
	if err := c.Auth(nil); err != nil {
		refusal("auth", err)
	}
	if err := c.Hello(); err != nil {
		refusal("hello", err)
	}
	return c
}

// lookup opens a plain session (as go-keyring does: the secret crosses the bus in clear) and
// returns the first unlocked item matching attr=val.
func lookup(c *dbus.Conn, attr, val string) (session, item dbus.ObjectPath) {
	svc := c.Object(service, svcPath)
	var out dbus.Variant
	if err := svc.Call(svcIfc+".OpenSession", 0, "plain", dbus.MakeVariant("")).Store(&out, &session); err != nil {
		refusal("OpenSession", err)
	}
	var unlocked, locked []dbus.ObjectPath
	if err := svc.Call(svcIfc+".SearchItems", 0, map[string]string{attr: val}).Store(&unlocked, &locked); err != nil {
		refusal("SearchItems", err)
	}
	if len(unlocked) == 0 {
		failf(3, "ERROR no unlocked item %s=%s (locked: %v)", attr, val, locked)
	}
	return session, unlocked[0]
}

func get(c *dbus.Conn, attr, val string) {
	session, item := lookup(c, attr, val)
	var s secret
	if err := c.Object(service, item).Call("org.freedesktop.Secret.Item.GetSecret", 0, session).Store(&s); err != nil {
		refusal("GetSecret", err)
	}
	failf(0, "SECRET=%s", s.Value)
}

func store(attr, val, value string) {
	createItem("probe "+val, map[string]string{attr: val, "user": "alice"}, &secret{Value: []byte(value), ContentType: "text/plain"})
}

// replace re-creates, replace=true, whatever item includes attrs (pairs): the keyring picks it by a
// subset match. "fail" sends a parameter a plain session rejects: the keyring has rewritten the
// item's attributes by then and keeps them.
func replace(mode string, pairs []string) {
	attrs := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		attrs[pairs[i]] = pairs[i+1]
	}
	sec := secret{Value: []byte("SECRET-evil"), ContentType: "text/plain"}
	if mode == "fail" {
		sec.Params = []byte("x")
	}
	createItem("probe replace", attrs, &sec)
}

func createItem(label string, attrs map[string]string, sec *secret) {
	c := dial()
	svc := c.Object(service, svcPath)
	var out dbus.Variant
	var session dbus.ObjectPath
	if err := svc.Call(svcIfc+".OpenSession", 0, "plain", dbus.MakeVariant("")).Store(&out, &session); err != nil {
		refusal("OpenSession", err)
	}
	coll := dbus.ObjectPath("/org/freedesktop/secrets/aliases/default")
	props := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label":      dbus.MakeVariant(label),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(attrs),
	}
	sec.Session = session
	var item, prompt dbus.ObjectPath
	if err := c.Object(service, coll).Call("org.freedesktop.Secret.Collection.CreateItem", 0, props, *sec, true).
		Store(&item, &prompt); err != nil {
		refusal("CreateItem", err)
	}
	failf(0, "STORED %s", item)
}

// monitor becomes a bus monitor (falling back to an eavesdrop match) and reports whether needle
// crossed the bus within secs.
func monitor(needle string, secs time.Duration) {
	c := dial()
	err := c.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, []string{}, uint32(0)).Err
	if err != nil {
		fmt.Printf("BecomeMonitor: %v\n", err)
		if err := c.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "eavesdrop='true'").Err; err != nil {
			failf(1, "REFUSED eavesdrop: %v", err)
		}
	}
	// Set only now: an eavesdrop channel also swallows method replies.
	ch := make(chan *dbus.Message, 256)
	c.Eavesdrop(ch)
	fmt.Println("MONITORING")
	deadline := time.After(secs)
	for {
		select {
		case m := <-ch:
			var b bytes.Buffer // the secret is an ay: search the wire form, not fmt's numbers
			if m.EncodeTo(&b, binary.LittleEndian) == nil && bytes.Contains(b.Bytes(), []byte(needle)) {
				failf(0, "CAPTURED %v", m)
			}
		case <-deadline:
			failf(4, "NOTHING")
		}
	}
}

// connexec connects to the bus, hands the unauthenticated socket to readerBin, then execs into
// argv: the bus credits the connection to this pid, now running argv's exe.
func connexec(attr, val, readerBin string, argv []string) {
	fd, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		failf(3, "ERROR socket: %v", err)
	}
	if err := unix.Connect(fd, &unix.SockaddrUnix{Name: busPath()}); err != nil {
		refusal("connect", err)
	}
	reader := exec.CommandContext(context.Background(), readerBin, "reader", attr, val)
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(fd), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		failf(3, "ERROR reader: %v", err)
	}
	execInto(argv)
}

// sendexec queues GetSecret on its own connection, hands a dup to readerBin, and execs into argv
// before the gatekeeper answers.
func sendexec(attr, val, readerBin string, argv []string) {
	c := dial()
	session, item := lookup(c, attr, val)
	fd := connFD()
	dup, err := unix.Dup(fd)
	if err != nil {
		failf(3, "ERROR dup: %v", err)
	}
	reader := exec.CommandContext(context.Background(), readerBin, "rawreader", "SECRET-")
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(dup), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		failf(3, "ERROR reader: %v", err)
	}
	writeGetSecret(fd, item, session)
	execInto(argv)
}

// writeGetSecret sends Item.GetSecret raw on fd.
func writeGetSecret(fd int, item, session dbus.ObjectPath) {
	msg := &dbus.Message{
		Type: dbus.TypeMethodCall,
		Headers: map[dbus.HeaderField]dbus.Variant{
			dbus.FieldPath:        dbus.MakeVariant(item),
			dbus.FieldInterface:   dbus.MakeVariant("org.freedesktop.Secret.Item"),
			dbus.FieldMember:      dbus.MakeVariant("GetSecret"),
			dbus.FieldDestination: dbus.MakeVariant(service),
			dbus.FieldSignature:   dbus.MakeVariant(dbus.SignatureOf(session)),
		},
		Body: []any{session},
	}
	var buf bytes.Buffer
	if err := msg.EncodeTo(&buf, binary.LittleEndian); err != nil {
		failf(3, "ERROR encode: %v", err)
	}
	raw := buf.Bytes()
	binary.LittleEndian.PutUint32(raw[8:12], 0x7fff0001) // header serial: godbus keeps it unexported, 0 is invalid
	if _, err := unix.Write(fd, raw); err != nil {
		failf(3, "ERROR write: %v", err)
	}
}

// lend opens a session and finds the item on its own connection, hands a dup of the socket to
// readerBin, drops its own handle and stays alive: the bus credits the reader's sends to this process.
func lend(attr, val, readerBin string) {
	c := dial()
	session, item := lookup(c, attr, val)
	dup, err := unix.Dup(connFD())
	if err != nil {
		failf(3, "ERROR dup: %v", err)
	}
	reader := exec.CommandContext(context.Background(), readerBin, "rawget", string(item), string(session)) //nolint:gosec // the test names the reader
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(dup), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		failf(3, "ERROR reader: %v", err)
	}
	c.Close() // its reader would take the reply; the socket stays open through the dup
	_ = reader.Wait()
	os.Exit(reader.ProcessState.ExitCode())
}

// heldexec holds a bus connection until waitfile appears (the daemon restarts meanwhile), looks the
// item up on it, hands a dup to readerBin (rawget) and execs into argv, which connects anew: the bus
// still credits the held connection to this pid.
func heldexec(attr, val, waitfile, readerBin string, argv []string) {
	c := dial()
	fmt.Println("CONNECTED")
	for {
		if _, err := os.Stat(waitfile); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	session, item := lookup(c, attr, val)
	dup, err := unix.Dup(connFD())
	if err != nil {
		failf(3, "ERROR dup: %v", err)
	}
	reader := exec.CommandContext(context.Background(), readerBin, "rawget", string(item), string(session), "1500ms") //nolint:gosec // the test names the reader
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(dup), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		failf(3, "ERROR reader: %v", err)
	}
	execInto(argv)
}

// rawget sends GetSecret on an inherited bus socket (fd 3) after delay (default: until the lender
// dropped its handle) and reads the answer.
func rawget(item, session string, delay time.Duration) {
	time.Sleep(delay)
	writeGetSecret(3, dbus.ObjectPath(item), dbus.ObjectPath(session))
	rawreader("SECRET-")
}

// own asks the bus for the Secret Service name: whoever owns it receives every app's secrets.
func own() {
	reply, err := dial().RequestName(service, dbus.NameFlagDoNotQueue)
	if err != nil {
		refusal("RequestName", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		failf(4, "NOT OWNER (reply %d)", reply)
	}
	failf(0, "OWNER")
}

// fakebus binds a listener and renames it over path (the real bus linked aside first): every
// relay connect then reaches it. It answers the first connection's SASL: 0 = the write reached the
// relay, 1 = refused, 4 = nobody connected within secs.
func fakebus(path string, secs time.Duration) {
	tmp := path + ".fake"
	ln, err := net.ListenUnix("unix", &net.UnixAddr{Name: tmp, Net: "unix"})
	if err != nil {
		failf(3, "ERROR listen: %v", err)
	}
	ln.SetUnlinkOnClose(false)
	if rerr := os.Rename(tmp, path); rerr != nil {
		failf(3, "ERROR rename over %s: %v", path, rerr)
	}
	fmt.Println("LISTENING")
	_ = ln.SetDeadline(time.Now().Add(secs))
	c, err := ln.AcceptUnix()
	if err != nil {
		failf(4, "NOTHING: %v", err)
	}
	_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 256)
	n, _ := c.Read(buf)
	if _, err := c.Write([]byte("OK 00112233445566778899aabbccddeeff\r\n")); err != nil {
		refusal("write to the relay", err)
	}
	failf(0, "WROTE after %q", buf[:n])
}

func execInto(argv []string) {
	err := syscall.Exec(argv[0], argv, os.Environ()) //nolint:gosec // the test names what to exec
	failf(3, "ERROR exec %s: %v", argv[0], err)
}

// reader speaks D-Bus on an inherited, unauthenticated socket (fd 3).
func reader(attr, val string) {
	time.Sleep(500 * time.Millisecond) // let the parent exec first
	nc, err := net.FileConn(os.NewFile(3, "bus"))
	if err != nil {
		failf(3, "ERROR fileconn: %v", err)
	}
	c, err := dbus.NewConn(nc)
	if err != nil {
		failf(3, "ERROR newconn: %v", err)
	}
	get(handshake(c), attr, val)
}

// rawreader reads an inherited bus socket (fd 3) for 5s and reports whether a reply carried needle.
func rawreader(needle string) {
	f := os.NewFile(3, "bus")
	_ = unix.SetNonblock(3, true)
	var seen []byte
	buf := make([]byte, 65536)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		n, err := unix.Read(int(f.Fd()), buf)
		if n > 0 {
			seen = append(seen, buf[:n]...)
			if i := bytes.Index(seen, []byte(needle)); i >= 0 {
				failf(0, "SECRET=%s", strings.TrimRight(string(seen[i:min(len(seen), i+64)]), "\x00"))
			}
			if bytes.Contains(seen, []byte("AccessDenied")) {
				failf(1, "REFUSED GetSecret (AccessDenied)")
			}
		}
		if err != nil && !errors.Is(err, unix.EAGAIN) {
			failf(3, "ERROR read: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}
	failf(4, "NOTHING")
}

// connFD finds the bus socket: godbus keeps its transport private.
func connFD() int {
	ents, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		failf(3, "ERROR fds: %v", err)
	}
	for _, e := range ents {
		var fd int
		if _, err := fmt.Sscan(e.Name(), &fd); err != nil || fd < 3 {
			continue
		}
		sa, err := unix.Getpeername(fd)
		if err != nil {
			continue
		}
		if u, ok := sa.(*unix.SockaddrUnix); ok && u.Name == busPath() {
			return fd
		}
	}
	failf(3, "ERROR bus fd not found")
	return -1
}

func usage() {
	failf(3, "usage: secret_probe get|store|replace|own|monitor|connexec|sendexec|lend|heldexec|reader|rawget|rawreader|fakebus|dialsleep|sleep ...")
}

func duration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

// modes run with the arguments after the mode name.
var modes = map[string]func(a []string){
	"get":       func(a []string) { get(dial(), a[0], a[1]) },
	"store":     func(a []string) { store(a[0], a[1], a[2]) },
	"replace":   func(a []string) { replace(a[0], a[1:]) },
	"own":       func([]string) { own() },
	"monitor":   func(a []string) { monitor(a[0], duration(a[1])) },
	"connexec":  func(a []string) { connexec(a[0], a[1], a[2], a[3:]) },
	"sendexec":  func(a []string) { sendexec(a[0], a[1], a[2], a[3:]) },
	"reader":    func(a []string) { reader(a[0], a[1]) },
	"lend":      func(a []string) { lend(a[0], a[1], a[2]) },
	"heldexec":  func(a []string) { heldexec(a[0], a[1], a[2], a[3], a[4:]) },
	"rawreader": func(a []string) { rawreader(a[0]) },
	"fakebus":   func(a []string) { fakebus(a[0], duration(a[1])) },
	"sleep":     func(a []string) { time.Sleep(duration(a[0])) },
	"rawget": func(a []string) {
		delay := 500 * time.Millisecond
		if len(a) > 2 {
			delay = duration(a[2])
		}
		rawget(a[0], a[1], delay)
	},
	"dialsleep": func(a []string) {
		c := dial()
		time.Sleep(duration(a[0]))
		c.Close()
	},
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	run, ok := modes[os.Args[1]]
	if !ok {
		usage()
	}
	run(os.Args[2:])
}
