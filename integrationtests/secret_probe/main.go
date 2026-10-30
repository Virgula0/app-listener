// secret_probe is a Secret Service client for the integration tests: a whitelisted copy plays the
// app that owns an item, other copies try to reach it. Exit codes: 0 = got the secret, 1 = refused
// (AccessDenied or a refused connect), 3 = setup error, 4 = nothing captured (own: not the owner).
package main

import (
	"bytes"
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

func fail(code int, format string, a ...any) {
	fmt.Printf(format+"\n", a...)
	os.Exit(code)
}

// refusal maps a D-Bus error to exit 1 when the gatekeeper or the bus policy refused it.
func refusal(what string, err error) {
	var de dbus.Error
	if errors.As(err, &de) && (de.Name == "org.freedesktop.DBus.Error.AccessDenied" ||
		strings.Contains(de.Name, "NotAuthorized") || strings.Contains(de.Name, "IsDenied")) {
		fail(1, "REFUSED %s: %v", what, err)
	}
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
		fail(1, "REFUSED %s: %v", what, err)
	}
	fail(3, "ERROR %s: %v", what, err)
}

func busPath() string {
	addr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	p, ok := strings.CutPrefix(addr, "unix:path=")
	if !ok {
		fail(3, "ERROR unsupported DBUS_SESSION_BUS_ADDRESS %q", addr)
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
func lookup(c *dbus.Conn, attr, val string) (dbus.ObjectPath, dbus.ObjectPath) {
	svc := c.Object(service, svcPath)
	var out dbus.Variant
	var session dbus.ObjectPath
	if err := svc.Call(svcIfc+".OpenSession", 0, "plain", dbus.MakeVariant("")).Store(&out, &session); err != nil {
		refusal("OpenSession", err)
	}
	var unlocked, locked []dbus.ObjectPath
	if err := svc.Call(svcIfc+".SearchItems", 0, map[string]string{attr: val}).Store(&unlocked, &locked); err != nil {
		refusal("SearchItems", err)
	}
	if len(unlocked) == 0 {
		fail(3, "ERROR no unlocked item %s=%s (locked: %v)", attr, val, locked)
	}
	return session, unlocked[0]
}

func get(c *dbus.Conn, attr, val string) {
	session, item := lookup(c, attr, val)
	var s secret
	if err := c.Object(service, item).Call("org.freedesktop.Secret.Item.GetSecret", 0, session).Store(&s); err != nil {
		refusal("GetSecret", err)
	}
	fail(0, "SECRET=%s", s.Value)
}

func store(attr, val, value string) {
	c := dial()
	svc := c.Object(service, svcPath)
	var out dbus.Variant
	var session dbus.ObjectPath
	if err := svc.Call(svcIfc+".OpenSession", 0, "plain", dbus.MakeVariant("")).Store(&out, &session); err != nil {
		refusal("OpenSession", err)
	}
	coll := dbus.ObjectPath("/org/freedesktop/secrets/aliases/default")
	props := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label":      dbus.MakeVariant("probe " + val),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{attr: val}),
	}
	sec := secret{Session: session, Value: []byte(value), ContentType: "text/plain"}
	var item, prompt dbus.ObjectPath
	if err := c.Object(service, coll).Call("org.freedesktop.Secret.Collection.CreateItem", 0, props, sec, true).
		Store(&item, &prompt); err != nil {
		refusal("CreateItem", err)
	}
	fail(0, "STORED %s", item)
}

// monitor becomes a bus monitor (falling back to an eavesdrop match) and reports whether needle
// crossed the bus within secs.
func monitor(needle string, secs time.Duration) {
	c := dial()
	err := c.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, []string{}, uint32(0)).Err
	if err != nil {
		fmt.Printf("BecomeMonitor: %v\n", err)
		if err := c.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "eavesdrop='true'").Err; err != nil {
			fail(1, "REFUSED eavesdrop: %v", err)
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
				fail(0, "CAPTURED %v", m)
			}
		case <-deadline:
			fail(4, "NOTHING")
		}
	}
}

// connexec connects to the bus, hands the unauthenticated socket to readerBin, then execs into
// argv: the bus credits the connection to this pid, now running argv's exe.
func connexec(attr, val, readerBin string, argv []string) {
	fd, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		fail(3, "ERROR socket: %v", err)
	}
	if err := unix.Connect(fd, &unix.SockaddrUnix{Name: busPath()}); err != nil {
		refusal("connect", err)
	}
	reader := exec.Command(readerBin, "reader", attr, val)
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(fd), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		fail(3, "ERROR reader: %v", err)
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
		fail(3, "ERROR dup: %v", err)
	}
	reader := exec.Command(readerBin, "rawreader", "SECRET-")
	reader.ExtraFiles = []*os.File{os.NewFile(uintptr(dup), "bus")}
	reader.Stdout, reader.Stderr = os.Stdout, os.Stderr
	if err := reader.Start(); err != nil {
		fail(3, "ERROR reader: %v", err)
	}
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
		fail(3, "ERROR encode: %v", err)
	}
	raw := buf.Bytes()
	binary.LittleEndian.PutUint32(raw[8:12], 0x7fff0001) // header serial: godbus keeps it unexported, 0 is invalid
	if _, err := unix.Write(fd, raw); err != nil {
		fail(3, "ERROR write: %v", err)
	}
	execInto(argv)
}

// own asks the bus for the Secret Service name: whoever owns it receives every app's secrets.
func own() {
	reply, err := dial().RequestName(service, dbus.NameFlagDoNotQueue)
	if err != nil {
		refusal("RequestName", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fail(4, "NOT OWNER (reply %d)", reply)
	}
	fail(0, "OWNER")
}

func execInto(argv []string) {
	err := syscall.Exec(argv[0], argv, os.Environ())
	fail(3, "ERROR exec %s: %v", argv[0], err)
}

// reader speaks D-Bus on an inherited, unauthenticated socket (fd 3).
func reader(attr, val string) {
	time.Sleep(500 * time.Millisecond) // let the parent exec first
	nc, err := net.FileConn(os.NewFile(3, "bus"))
	if err != nil {
		fail(3, "ERROR fileconn: %v", err)
	}
	c, err := dbus.NewConn(nc)
	if err != nil {
		fail(3, "ERROR newconn: %v", err)
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
				fail(0, "SECRET=%s", strings.TrimRight(string(seen[i:min(len(seen), i+64)]), "\x00"))
			}
			if bytes.Contains(seen, []byte("AccessDenied")) {
				fail(1, "REFUSED GetSecret (AccessDenied)")
			}
		}
		if err != nil && !errors.Is(err, unix.EAGAIN) {
			fail(3, "ERROR read: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}
	fail(4, "NOTHING")
}

// connFD finds the bus socket: godbus keeps its transport private.
func connFD() int {
	ents, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		fail(3, "ERROR fds: %v", err)
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
	fail(3, "ERROR bus fd not found")
	return -1
}

func usage() {
	fail(3, "usage: secret_probe get|store|own|monitor|connexec|sendexec|reader|rawreader|sleep ...")
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	a := os.Args[2:]
	switch os.Args[1] {
	case "get":
		get(dial(), a[0], a[1])
	case "store":
		store(a[0], a[1], a[2])
	case "own":
		own()
	case "monitor":
		d, _ := time.ParseDuration(a[1])
		monitor(a[0], d)
	case "connexec":
		connexec(a[0], a[1], a[2], a[3:])
	case "sendexec":
		sendexec(a[0], a[1], a[2], a[3:])
	case "reader":
		reader(a[0], a[1])
	case "rawreader":
		rawreader(a[0])
	case "sleep":
		d, _ := time.ParseDuration(a[0])
		time.Sleep(d)
	default:
		usage()
	}
}
