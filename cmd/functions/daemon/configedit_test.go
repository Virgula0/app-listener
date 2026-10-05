package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

const (
	oldConf = "[watch /a]\nneed_encryption: false\n/usr/bin/cat\n"
	newConf = "[watch /a]\nneed_encryption: false\n/usr/bin/cat\n/usr/bin/head\n"
)

func newTestConfigEditor(t *testing.T, reloadErr error) (*configEditor, *int) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(path, []byte(oldConf), 0o600); err != nil {
		t.Fatal(err)
	}
	reloads := 0
	return &configEditor{path: path, reload: func() (*daemonconfig.Config, error) {
		reloads++
		got, _ := os.ReadFile(path)
		if string(got) != newConf {
			return nil, fmt.Errorf("reload saw %q", got)
		}
		return nil, reloadErr
	}}, &reloads
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestConfigEditorApply(t *testing.T) {
	ed, reloads := newTestConfigEditor(t, nil)
	if err := ed.apply([]byte(oldConf), []byte(newConf)); err != nil {
		t.Fatal(err)
	}
	if *reloads != 1 || readFile(t, ed.path) != newConf {
		t.Fatalf("reloads=%d file=%q", *reloads, readFile(t, ed.path))
	}
}

func TestConfigEditorRestoresOnFailedReload(t *testing.T) {
	ed, _ := newTestConfigEditor(t, errors.New("no free slots"))
	err := ed.apply([]byte(oldConf), []byte(newConf))
	if err == nil || !strings.Contains(err.Error(), "no free slots") || !strings.Contains(err.Error(), "restored") {
		t.Fatalf("err = %v", err)
	}
	if readFile(t, ed.path) != oldConf {
		t.Fatalf("a failed reload left %q", readFile(t, ed.path))
	}
}

func TestConfigEditorRefusals(t *testing.T) {
	ed, reloads := newTestConfigEditor(t, nil)
	if err := ed.apply([]byte(oldConf), []byte("[watch\n")); err == nil {
		t.Error("an unparsable configuration was applied")
	}
	if err := ed.apply([]byte("something else"), []byte(newConf)); err == nil {
		t.Error("a stale base was applied")
	}
	if *reloads != 0 || readFile(t, ed.path) != oldConf {
		t.Fatalf("a refused edit touched state: reloads=%d file=%q", *reloads, readFile(t, ed.path))
	}
}

func TestRunConfigEditOverTheSocket(t *testing.T) {
	ed, _ := newTestConfigEditor(t, nil)
	cs := &controlServer{configEditor: func() *configEditor { return ed }}
	server, client := net.Pipe()
	defer client.Close()
	go func() {
		cs.runConfigEdit(server, bufio.NewReader(server))
		server.Close()
	}()

	r := bufio.NewReader(client)
	head, err := r.ReadString('\n')
	if err != nil || head != fmt.Sprintf("OK %d\n", len(oldConf)) {
		t.Fatalf("head %q, %v", head, err)
	}
	body := make([]byte, len(oldConf))
	if _, err := io.ReadFull(r, body); err != nil || string(body) != oldConf {
		t.Fatalf("body %q, %v", body, err)
	}
	if _, err := fmt.Fprintf(client, "PUT %d\n%s", len(newConf), newConf); err != nil {
		t.Fatal(err)
	}
	if ans, _ := r.ReadString('\n'); ans != "OK\n" {
		t.Fatalf("answer %q", ans)
	}
	if readFile(t, ed.path) != newConf {
		t.Fatalf("file %q", readFile(t, ed.path))
	}
}

func TestRunConfigEditRefusedDuringALiveSession(t *testing.T) {
	ed, reloads := newTestConfigEditor(t, nil)
	cs := &controlServer{configEditor: func() *configEditor { return ed }, active: &editControlSession{}}
	server, client := net.Pipe()
	defer client.Close()
	go func() {
		cs.runConfigEdit(server, bufio.NewReader(server))
		server.Close()
	}()
	r := bufio.NewReader(client)
	_, _ = r.ReadString('\n')
	_, _ = io.ReadFull(r, make([]byte, len(oldConf)))
	_, _ = fmt.Fprintf(client, "PUT %d\n%s", len(newConf), newConf)
	if ans, _ := r.ReadString('\n'); !strings.HasPrefix(ans, "ERR ") {
		t.Fatalf("answer %q, want a refusal", ans)
	}
	if *reloads != 0 {
		t.Fatal("reloaded during a live session")
	}
}

func TestReadGrantRequestConfig(t *testing.T) {
	req, err := readGrantRequest(bufio.NewReader(strings.NewReader("CONFIG\n")))
	if err != nil || !req.config {
		t.Fatalf("%+v, %v", req, err)
	}
	if _, err := readGrantRequest(bufio.NewReader(strings.NewReader("CONFIG /x\n"))); err == nil {
		t.Fatal("CONFIG with an argument accepted")
	}
}
