package secretgate

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/godbus/dbus/v5"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func call(iface, member string, path dbus.ObjectPath, body ...any) *dbus.Message {
	h := map[dbus.HeaderField]dbus.Variant{
		dbus.FieldPath:        dbus.MakeVariant(path),
		dbus.FieldMember:      dbus.MakeVariant(member),
		dbus.FieldDestination: dbus.MakeVariant("org.freedesktop.secrets"),
	}
	if iface != "" {
		h[dbus.FieldInterface] = dbus.MakeVariant(iface)
	}
	if len(body) > 0 {
		h[dbus.FieldSignature] = dbus.MakeVariant(dbus.SignatureOf(body...))
	}
	return &dbus.Message{Type: dbus.TypeMethodCall, Headers: h, Body: body}
}

func TestClassify(t *testing.T) {
	item := dbus.ObjectPath("/org/freedesktop/secrets/collection/login/1")
	session := dbus.ObjectPath("/org/freedesktop/secrets/session/s1")
	attrs := map[string]string{"service": "gh:github.com"}
	props := map[string]dbus.Variant{propAttributes: dbus.MakeVariant(attrs)}
	secret := struct {
		S dbus.ObjectPath
		P []byte
		V []byte
		C string
	}{session, nil, []byte("x"), "text/plain"}
	cases := []struct {
		name  string
		m     *dbus.Message
		kind  verdictKind
		op    ebpf.EventType
		items int
		attrs int
	}{
		{"OpenSession", call(ifaceService, "OpenSession", "/org/freedesktop/secrets", "plain", dbus.MakeVariant("")), allow, 0, 0, 0},
		{"SearchItems", call(ifaceService, "SearchItems", "/org/freedesktop/secrets", attrs), allow, 0, 0, 0},
		{"GetSecrets", call(ifaceService, "GetSecrets", "/org/freedesktop/secrets", []dbus.ObjectPath{item, item}, session),
			gated, ebpf.EventRead, 2, 0},
		{"GetSecrets bad args", call(ifaceService, "GetSecrets", "/org/freedesktop/secrets", "x"), refuse, 0, 0, 0},
		{"GetSecret", call(ifaceItem, "GetSecret", item, session), gated, ebpf.EventRead, 1, 0},
		{"SetSecret", call(ifaceItem, "SetSecret", item, secret), gated, ebpf.EventWrite, 1, 0},
		{"Item.Delete", call(ifaceItem, "Delete", item), gated, ebpf.EventDelete, 1, 0},
		{"CreateItem", call(ifaceCollection, "CreateItem", "/org/freedesktop/secrets/aliases/default", props, secret, true),
			gated, ebpf.EventWrite, 0, 1},
		{"CreateItem no attrs", call(ifaceCollection, "CreateItem", "/org/freedesktop/secrets/aliases/default",
			map[string]dbus.Variant{}, secret, false), gated, ebpf.EventWrite, 0, 1},
		{"CreateItem bad attrs", call(ifaceCollection, "CreateItem", "/c",
			map[string]dbus.Variant{propAttributes: dbus.MakeVariant("x")}, secret, false), refuse, 0, 0, 0},
		{"Collection.Delete", call(ifaceCollection, "Delete", "/c"), needManager, 0, 0, 0},
		{"SetAlias", call(ifaceService, "SetAlias", "/org/freedesktop/secrets", "default", dbus.ObjectPath("/c")), needManager, 0, 0, 0},
		{"Props.Get", call(ifaceProps, "Get", item, ifaceItem, "Attributes"), allow, 0, 0, 0},
		{"Props.Set label", call(ifaceProps, "Set", item, ifaceItem, "Label", dbus.MakeVariant("l")), gated, ebpf.EventWrite, 1, 0},
		{"Props.Set attrs", call(ifaceProps, "Set", item, ifaceItem, "Attributes", dbus.MakeVariant(attrs)),
			gated, ebpf.EventWrite, 1, 1},
		{"Props.Set collection", call(ifaceProps, "Set", "/c", ifaceCollection, "Label", dbus.MakeVariant("l")), needManager, 0, 0, 0},
		{"Props.Set other", call(ifaceProps, "Set", "/c", "org.example.X", "Y", dbus.MakeVariant("l")), refuse, 0, 0, 0},
		{"Prompter callback", call(ifacePrompterCallback, "PromptReady", "/p", "yes", map[string]dbus.Variant{}, "x"),
			needPrompter, 0, 0, 0},
		{"master password", call(ifaceGuiltRidden, "UnlockWithMasterPassword", "/org/freedesktop/secrets",
			dbus.ObjectPath("/c"), secret), needManager, 0, 0, 0},
		{"portal", call("org.freedesktop.impl.portal.Secret", "RetrieveSecret", "/org/freedesktop/portal/desktop"), refuse, 0, 0, 0},
		{"no interface", call("", "GetSecret", item, session), refuse, 0, 0, 0},
		{"unknown member", call(ifaceItem, "Steal", item), refuse, 0, 0, 0},
		{"introspect", call("org.freedesktop.DBus.Introspectable", "Introspect", "/"), allow, 0, 0, 0},
	}
	for _, c := range cases {
		v := classify(c.m)
		if v.kind != c.kind || v.op != c.op || len(v.items) != c.items || len(v.attrs) != c.attrs {
			t.Errorf("%s: got %+v", c.name, v)
		}
	}
}

func TestSASL(t *testing.T) {
	for name, client := range map[string]string{
		"initial response": "\x00AUTH EXTERNAL 31303030\r\nNEGOTIATE_UNIX_FD\r\nBEGIN\r\nlnext",
		"data exchange":    "\x00AUTH EXTERNAL\r\nDATA\r\nBEGIN\r\nlnext",
		"retry":            "\x00AUTH EXTERNAL 30\r\nAUTH EXTERNAL 31303030\r\nBEGIN\r\nlnext",
	} {
		var out bytes.Buffer
		br, err := serverAuth(struct {
			io.Reader
			io.Writer
		}{strings.NewReader(client), &out}, 1000)
		if err != nil {
			t.Fatalf("%s: %v (sent %q)", name, err, out.String())
		}
		rest, _ := io.ReadAll(br)
		if string(rest) != "lnext" || !strings.Contains(out.String(), "OK ") {
			t.Errorf("%s: rest %q, replies %q", name, rest, out.String())
		}
		if strings.Contains(out.String(), "AGREE_UNIX_FD") {
			t.Errorf("%s: agreed to fd passing", name)
		}
	}
	for name, client := range map[string]string{
		"other uid":     "\x00AUTH EXTERNAL 30\r\nBEGIN\r\n",
		"no NUL":        "AUTH EXTERNAL 31303030\r\nBEGIN\r\n",
		"other mech":    "\x00AUTH ANONYMOUS\r\nBEGIN\r\n",
		"wrong data":    "\x00AUTH EXTERNAL\r\nDATA 30\r\nBEGIN\r\n",
		"cancel":        "\x00AUTH EXTERNAL 31303030\r\nCANCEL\r\nBEGIN\r\n",
		"never begins":  strings.Repeat("\x00", 1) + strings.Repeat("NOOP\r\n", 20),
		"hangs up auth": "\x00AUTH EXTERNAL",
	} {
		_, err := serverAuth(struct {
			io.Reader
			io.Writer
		}{strings.NewReader(client), io.Discard}, 1000)
		if err == nil {
			t.Errorf("%s: authenticated", name)
		}
	}
}

func TestClientAuth(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	go func() {
		br, err := serverAuth(b, 0)
		if err == nil {
			_, _ = br.ReadByte()
		}
		b.Close()
	}()
	br, err := clientAuth(a, 0)
	if err != nil {
		t.Fatal(err)
	}
	_ = br
	_, _ = a.Write([]byte("l"))
}

func TestReadFrame(t *testing.T) {
	m := call(ifaceItem, "GetSecret", "/org/freedesktop/secrets/collection/login/1", dbus.ObjectPath("/s"))
	frame, err := encodeWithSerial(m, 7)
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(append(append([]byte{}, frame...), frame...))
	for range 2 {
		got, err := readFrame(r)
		if err != nil || !bytes.Equal(got, frame) {
			t.Fatalf("got %d bytes (%v), want %d", len(got), err, len(frame))
		}
		d, err := decodeFrame(got)
		if err != nil || d.Serial() != 7 || header[string](d, dbus.FieldMember) != "GetSecret" {
			t.Fatalf("decoded %v (%v)", d, err)
		}
	}
	huge := make([]byte, 16)
	huge[0] = 'l'
	binary.LittleEndian.PutUint32(huge[4:8], maxMessage)
	if _, err := readFrame(bytes.NewReader(huge)); err == nil {
		t.Error("accepted a message over the D-Bus limit")
	}
	huge[0] = 'x'
	if _, err := readFrame(bytes.NewReader(huge)); err == nil {
		t.Error("accepted a bad endianness byte")
	}
}

func TestClassifyOpenCallsExact(t *testing.T) {
	for _, m := range []*dbus.Message{
		call(ifacePrompt+".Prompt", "GetSecret", "/org/freedesktop/secrets/collection/login/1"),
		call(ifaceSession, "Steal", "/s"),
		call("org.freedesktop.DBus.Introspectable", "", "/"),
	} {
		if v := classify(m); v.kind != refuse {
			t.Errorf("%v: got %+v", m, v)
		}
	}
}
