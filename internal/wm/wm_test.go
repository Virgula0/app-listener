package wm

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// shortDir is a temp dir short enough for a unix socket path (108 bytes).
func shortDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "wm")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(d) })
	return d
}

// serve answers each connection on path with reply(request) until the test ends.
func serve(t *testing.T, path string, handle func(net.Conn)) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = l.Close() })
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go func() { defer c.Close(); handle(c) }()
		}
	}()
}

func envOf(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func uid() uint32 { return uint32(os.Getuid()) } //nolint:gosec // uids fit

func TestFitSize(t *testing.T) {
	area := Rect{W: 1920, H: 1080}
	cases := []struct {
		pref, minSize, want Size
	}{
		{Size{1280, 820}, Size{600, 400}, Size{1280, 820}},
		{Size{1280, 820}, Size{1400, 900}, Size{1400, 900}},
		{Size{2560, 1600}, Size{0, 0}, Size{1728, 972}},
		{Size{1280, 820}, Size{3000, 3000}, Size{1728, 972}},
	}
	for _, c := range cases {
		if got := FitSize(c.pref, c.minSize, area); got != c.want {
			t.Errorf("FitSize(%v, %v) = %v, want %v", c.pref, c.minSize, got, c.want)
		}
	}
}

func TestDetect(t *testing.T) {
	dir := shortDir(t)
	sig := "abc_123"
	hypr := filepath.Join(dir, "hypr", sig, ".socket.sock")
	serve(t, hypr, func(net.Conn) {})
	sway := filepath.Join(dir, "sway.sock")
	serve(t, sway, func(net.Conn) {})
	plain := filepath.Join(dir, "plain")
	if err := os.WriteFile(plain, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		env  map[string]string
		uid  uint32
		want Kind
	}{
		{"hyprland", map[string]string{"XDG_RUNTIME_DIR": dir, "HYPRLAND_INSTANCE_SIGNATURE": sig}, uid(), Hyprland},
		{"hyprland signature with a path", map[string]string{"XDG_RUNTIME_DIR": dir,
			"HYPRLAND_INSTANCE_SIGNATURE": "../hypr/" + sig}, uid(), None},
		{"socket of another user", map[string]string{"XDG_RUNTIME_DIR": dir, "HYPRLAND_INSTANCE_SIGNATURE": sig},
			uid() + 1, None},
		{"sway", map[string]string{"SWAYSOCK": sway}, uid(), Sway},
		{"i3sock not a socket", map[string]string{"I3SOCK": plain}, uid(), None},
		{"relative socket path", map[string]string{"SWAYSOCK": "sway.sock"}, uid(), None},
		{"nothing", map[string]string{"XDG_RUNTIME_DIR": dir}, uid(), None},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Detect(envOf(c.env), c.uid).Kind; got != c.want {
				t.Fatalf("Detect = %v, want %v", got, c.want)
			}
		})
	}
}

const hyprClients = `[{"address":"0xaaa","mapped":true,"floating":false,"monitor":1,"pid":4242},
 {"address":"0xbbb","mapped":true,"floating":false,"monitor":0,"pid":1}]`

const hyprMonitors = `[{"id":0,"x":0,"y":0,"width":1920,"height":1080,"scale":1,"transform":0,"reserved":[0,0,0,0]},
 {"id":1,"x":1920,"y":0,"width":5120,"height":2880,"scale":2,"transform":0,"reserved":[0,40,0,0]}]`

func TestHyprlandFloat(t *testing.T) {
	dir := shortDir(t)
	sock := filepath.Join(dir, "hypr", "s", ".socket.sock")
	var mu sync.Mutex
	var dispatched []string
	serve(t, sock, func(c net.Conn) {
		buf := make([]byte, 4096)
		n, _ := c.Read(buf)
		req := string(buf[:n])
		switch {
		case req == "j/clients":
			_, _ = io.WriteString(c, hyprClients)
		case req == "j/monitors":
			_, _ = io.WriteString(c, hyprMonitors)
		case strings.HasPrefix(req, "[[BATCH]]"):
			mu.Lock()
			dispatched = append(dispatched, req)
			mu.Unlock()
			_, _ = io.WriteString(c, strings.Repeat("ok\n\n", strings.Count(req, "dispatch ")))
		}
	})
	s := Detect(envOf(map[string]string{"XDG_RUNTIME_DIR": dir, "HYPRLAND_INSTANCE_SIGNATURE": "s"}), uid())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Float(ctx, Window{PID: 4242}, Size{1280, 820}, Size{800, 500}); err != nil {
		t.Fatal(err)
	}
	// Monitor 1: 2560x1440 logical at x=1920, minus a 40px top bar -> area 1920,40 2560x1400.
	want := "[[BATCH]]dispatch setfloating address:0xaaa;dispatch resizewindowpixel exact 1280 820,address:0xaaa;" +
		"dispatch movewindowpixel exact 2560 330,address:0xaaa"
	mu.Lock()
	defer mu.Unlock()
	if len(dispatched) != 1 || dispatched[0] != want {
		t.Fatalf("dispatched %q, want %q", dispatched, want)
	}
}

func TestHyprlandFloatTimesOutWithoutWindow(t *testing.T) {
	dir := shortDir(t)
	sock := filepath.Join(dir, "hypr", "s", ".socket.sock")
	serve(t, sock, func(c net.Conn) { _, _ = io.WriteString(c, "[]") })
	s := Detect(envOf(map[string]string{"XDG_RUNTIME_DIR": dir, "HYPRLAND_INSTANCE_SIGNATURE": "s"}), uid())
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	if err := s.Float(ctx, Window{PID: 4242}, Size{1280, 820}, Size{}); err == nil {
		t.Fatal("Float succeeded without the window in the client list")
	}
}

func TestI3Float(t *testing.T) {
	dir := shortDir(t)
	sock := filepath.Join(dir, "i3", "ipc-socket.1")
	var mu sync.Mutex
	var commands []string
	serve(t, sock, func(c net.Conn) {
		head := make([]byte, 14)
		if _, err := io.ReadFull(c, head); err != nil {
			return
		}
		body := make([]byte, binary.NativeEndian.Uint32(head[6:]))
		_, _ = io.ReadFull(c, body)
		var reply string
		switch binary.NativeEndian.Uint32(head[10:]) {
		case i3GetTree:
			reply = `{"nodes":[{"window":77,"nodes":[]}]}`
		case i3GetWorkspaces:
			reply = `[{"focused":false,"rect":{"x":0,"y":0,"width":800,"height":600}},
				{"focused":true,"rect":{"x":1920,"y":20,"width":1920,"height":1060}}]`
		case i3RunCommand:
			mu.Lock()
			commands = append(commands, string(body))
			mu.Unlock()
			reply = `[{"success":true}]`
		}
		out := append([]byte("i3-ipc"), binary.NativeEndian.AppendUint32(nil, uint32(len(reply)))...) //nolint:gosec // test
		out = binary.NativeEndian.AppendUint32(out, 0)
		_, _ = c.Write(append(out, reply...))
	})
	env := map[string]string{"XDG_RUNTIME_DIR": dir, "DISPLAY": ":0"}
	s := Detect(envOf(env), uid())
	if s.Kind != I3 {
		t.Fatalf("Detect = %v, want i3 from the runtime-dir socket", s.Kind)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Float(ctx, Window{PID: 1, XID: 77}, Size{1280, 820}, Size{}); err != nil {
		t.Fatal(err)
	}
	want := "[id=77] floating enable, resize set 1280 px 820 px, move position center"
	mu.Lock()
	defer mu.Unlock()
	if len(commands) != 1 || commands[0] != want {
		t.Fatalf("commands %q, want %q", commands, want)
	}
	if err := s.Float(ctx, Window{PID: 1}, Size{1280, 820}, Size{}); err == nil {
		t.Fatal("i3 Float without an X11 id succeeded")
	}
}
