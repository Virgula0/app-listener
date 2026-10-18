package editbackend

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

// TestHelperBackend is the fake backend the client tests launch (re-exec of the test binary).
func TestHelperBackend(t *testing.T) {
	if os.Getenv("EDITBACKEND_HELPER") == "" {
		t.Skip("helper process")
	}
	in, out := NewReader(os.Stdin), NewWriter(os.Stdout)
	proto := Proto
	if os.Getenv("EDITBACKEND_HELPER") == "oldproto" {
		proto = Proto + 1
	}
	for {
		var req Request
		if err := in.Read(&req); err != nil {
			os.Exit(0)
		}
		switch req.Op {
		case OpHello:
			_ = out.Write(&Message{Kind: KindReply, ID: req.ID, Hello: &Hello{Proto: proto, Live: true}})
		case OpPing:
			_ = out.Write(&Message{Kind: KindEvent, Event: &Event{Op: "OPEN", Path: "/x"}})
			_ = out.Write(&Message{Kind: KindEnded, Err: "gone"})
			_ = out.Write(&Message{Kind: KindReply, ID: req.ID})
		case OpRead:
			_ = out.Write(&Message{Kind: KindReply, ID: req.ID, File: &File{Data: []byte(req.Path)}})
		default:
			_ = out.Write(&Message{Kind: KindReply, ID: req.ID, Err: "refused " + req.Op})
		}
	}
}

func helper(t *testing.T, mode string) (*Client, error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestHelperBackend$")
	cmd.Env = append(os.Environ(), "EDITBACKEND_HELPER="+mode)
	return launch(cmd, func(int) error { return nil })
}

func TestClientCallsAndEvents(t *testing.T) {
	c, err := helper(t, "ok")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	m, err := c.Call(&Request{Op: OpRead, Path: "a/b"})
	if err != nil || string(m.File.Data) != "a/b" {
		t.Fatalf("read: %+v %v", m, err)
	}
	if _, err := c.Call(&Request{Op: OpWrite}); err == nil || !strings.Contains(err.Error(), "refused write") {
		t.Fatalf("an error reply must surface as err, got %v", err)
	}
	if _, err := c.Call(&Request{Op: OpPing}); err != nil {
		t.Fatal(err)
	}
	select {
	case ev := <-c.Events():
		if ev.Event == nil || ev.Event.Op != "OPEN" {
			t.Fatalf("event: %+v", ev)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no event")
	}
	select {
	case why := <-c.Ended():
		if why != "gone" {
			t.Fatalf("ended: %q", why)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no ended notice")
	}
}

func TestClientCloseEndsTheBackend(t *testing.T) {
	c, err := helper(t, "ok")
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
	select {
	case <-c.Done():
	default:
		t.Fatal("Close returned with the backend still running")
	}
	if _, err := c.Call(&Request{Op: OpRead}); !errors.Is(err, ErrClosed) {
		t.Fatalf("call after close: %v", err)
	}
}

func TestClientRefusesAnotherProtocol(t *testing.T) {
	if _, err := helper(t, "oldproto"); err == nil || !strings.Contains(err.Error(), "protocol") {
		t.Fatalf("protocol mismatch: %v", err)
	}
}

func TestClientVerifyRunsBeforeAnythingIsSent(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestHelperBackend$")
	cmd.Env = append(os.Environ(), "EDITBACKEND_HELPER=ok")
	refused := errors.New("not the daemon's binary")
	if _, err := launch(cmd, func(int) error { return refused }); !errors.Is(err, refused) {
		t.Fatalf("launch with a failing verify: %v", err)
	}
}

func TestCleanRel(t *testing.T) {
	for in, want := range map[string]string{"": ".", ".": ".", "a/b": "a/b", "a/../b": "b", "./a/": "a"} {
		if got, err := CleanRel(in); err != nil || got != want {
			t.Errorf("CleanRel(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"..", "../a", "a/../../b", "/etc/shadow", "a\x00b", "\xff", strings.Repeat("a", 5000)} {
		if got, err := CleanRel(in); err == nil {
			t.Errorf("CleanRel(%q) = %q, want refused", in, got)
		}
	}
}

func TestChildEnvIsAllowlisted(t *testing.T) {
	t.Setenv("LD_PRELOAD", "/tmp/evil.so")
	t.Setenv("PATH", "/tmp/evil:/usr/bin")
	t.Setenv("SUDO_UID", "1000")
	t.Setenv("SUDO_GID", "1000;rm -rf")
	env := childEnv()
	for _, kv := range env {
		if strings.HasPrefix(kv, "LD_PRELOAD=") || strings.Contains(kv, "/tmp/evil") || strings.HasPrefix(kv, "SUDO_GID=") {
			t.Errorf("leaked into the backend: %s", kv)
		}
	}
	if !slices.Contains(env, "SUDO_UID=1000") {
		t.Errorf("SUDO_UID not passed: %v", env)
	}
}

func TestReaderCapsLines(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_, _ = w.Write([]byte(`{"op":"` + strings.Repeat("a", MaxLine) + `"}` + "\n"))
		_ = w.Close()
	}()
	var req Request
	if err := NewReader(r).Read(&req); err == nil {
		t.Fatal("a line over MaxLine was accepted")
	}
	_ = r.Close()
}
