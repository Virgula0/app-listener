package editprotected

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/vaultfs"
)

func newTestBackend(t *testing.T) (*backend, *bytes.Buffer) {
	t.Helper()
	var out bytes.Buffer
	return &backend{out: editbackend.NewWriter(&out), timeout: SessionTimeoutDefault}, &out
}

// editingBackend is a backend past SELECT on a temp tree (root/a.txt, root/sub/), and an outside
// file a confined op must never reach.
func editingBackend(t *testing.T) (b *backend, root, outside string) {
	t.Helper()
	base := t.TempDir()
	root = filepath.Join(base, "vault")
	outside = filepath.Join(base, "outside.txt")
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o700); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "a.txt"), "alpha")
	mustWrite(t, outside, "OUTSIDE")
	b, _ = newTestBackend(t)
	v, err := vaultfs.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(v.Close)
	b.vault, b.resource, b.state = v, root, stEditing
	return b, root, outside
}

func call(b *backend, req editbackend.Request) *editbackend.Message {
	return b.handle(&req)
}

func TestBackendRefusesOpsOutOfState(t *testing.T) {
	b, _ := newTestBackend(t)
	// forward_check opens the named file: before AUTH it would be a probe through this binary's
	// read baseline on guarded trees.
	for _, op := range []string{editbackend.OpList, editbackend.OpRead, editbackend.OpWrite, editbackend.OpSelect,
		editbackend.OpForward, editbackend.OpForwardCheck, editbackend.OpConfigGet, editbackend.OpConfigPut, "nope"} {
		if m := call(b, editbackend.Request{ID: 1, Op: op, Path: "."}); m.Err == "" {
			t.Errorf("%s before auth: accepted, want refused", op)
		}
	}
	b.state = stAuthed
	if m := call(b, editbackend.Request{Op: editbackend.OpRead, Path: "a"}); m.Err == "" {
		t.Error("read after auth but before select: accepted")
	}
	if m := call(b, editbackend.Request{Op: editbackend.OpOfflineOpen, Resource: "/"}); m.Err == "" {
		t.Error("offline_open after auth: accepted")
	}
}

func TestBackendSelectOnlyListedResources(t *testing.T) {
	b, _ := newTestBackend(t)
	b.state, b.resources = stAuthed, []string{"/home/u/.ssh"}
	// A line break would smuggle a second control line to the daemon; only exact list entries pass.
	for _, r := range []string{"/home/u/.gnupg", "/home/u/.ssh\nFORWARD allow * 1 1", "/home/u/.ssh/"} {
		if m := call(b, editbackend.Request{Op: editbackend.OpSelect, Resource: r}); m.Err == "" ||
			!strings.Contains(m.Err, "not a guarded directory") {
			t.Errorf("select %q: err %q, want refused before the daemon is asked", r, m.Err)
		}
	}
}

func TestBackendPathsStayInsideTheVault(t *testing.T) {
	b, _, outside := editingBackend(t)
	for _, p := range []string{"../outside.txt", "sub/../../outside.txt", outside, "a\x00b", "\xff"} {
		if m := call(b, editbackend.Request{Op: editbackend.OpRead, Path: p}); m.Err == "" {
			t.Errorf("read %q: accepted", p)
		}
		if m := call(b, editbackend.Request{Op: editbackend.OpWrite, Path: p, Data: []byte("x")}); m.Err == "" {
			t.Errorf("write %q: accepted", p)
		}
	}
	if got := mustRead(t, outside); got != "OUTSIDE" {
		t.Fatalf("outside file changed: %q", got)
	}
}

func TestBackendSymlinksRefused(t *testing.T) {
	b, root, outside := editingBackend(t)
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Dir(outside), filepath.Join(root, "dirlink")); err != nil {
		t.Fatal(err)
	}
	for _, req := range []editbackend.Request{
		{Op: editbackend.OpRead, Path: "link"},
		{Op: editbackend.OpWrite, Path: "link", Data: []byte("pwned")},
		{Op: editbackend.OpWrite, Path: "dirlink/outside.txt", Data: []byte("pwned")},
		{Op: editbackend.OpChmod, Path: "link", Mode: "0777"},
		{Op: editbackend.OpCreate, Path: "dirlink", Name: "new"},
		{Op: editbackend.OpList, Path: "dirlink"},
	} {
		if m := call(b, req); m.Err == "" {
			t.Errorf("%s %s: accepted through a symlink", req.Op, req.Path)
		}
	}
	if got := mustRead(t, outside); got != "OUTSIDE" {
		t.Fatalf("outside file changed: %q", got)
	}
	if info, err := os.Stat(outside); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("outside file mode changed: %v %v", info.Mode(), err)
	}
}

func TestBackendEditFlow(t *testing.T) {
	b, root, _ := editingBackend(t)
	m := call(b, editbackend.Request{Op: editbackend.OpList, Path: "."})
	if m.Err != "" || len(m.Entries) != 2 {
		t.Fatalf("list: %+v", m)
	}
	m = call(b, editbackend.Request{Op: editbackend.OpRead, Path: "a.txt"})
	if m.Err != "" || string(m.File.Data) != "alpha" {
		t.Fatalf("read: %+v", m)
	}
	if m = call(b, editbackend.Request{Op: editbackend.OpWrite, Path: "a.txt", Data: []byte("beta")}); m.Err != "" {
		t.Fatalf("write: %s", m.Err)
	}
	if got := mustRead(t, filepath.Join(root, "a.txt")); got != "beta" {
		t.Fatalf("content = %q", got)
	}
	if m = call(b, editbackend.Request{Op: editbackend.OpCreate, Path: "sub", Name: "n.txt"}); m.Err != "" {
		t.Fatalf("create: %s", m.Err)
	}
	if m = call(b, editbackend.Request{Op: editbackend.OpChmod, Path: "sub/n.txt", Mode: "0600"}); m.Err != "" {
		t.Fatalf("chmod: %s", m.Err)
	}
	if m = call(b, editbackend.Request{Op: editbackend.OpRemove, Path: "."}); m.Err == "" {
		t.Fatal("removing the vault root was accepted")
	}
	if m = call(b, editbackend.Request{Op: editbackend.OpRemove, Path: "sub"}); m.Err != "" {
		t.Fatalf("remove: %s", m.Err)
	}
	if _, err := os.Lstat(filepath.Join(root, "sub")); !os.IsNotExist(err) {
		t.Fatalf("sub still there: %v", err)
	}
}

func TestBackendReadWithholdsBinaryAndOversized(t *testing.T) {
	b, root, _ := editingBackend(t)
	mustWrite(t, filepath.Join(root, "bin"), "a\x00b")
	m := call(b, editbackend.Request{Op: editbackend.OpRead, Path: "bin"})
	if m.Err != "" || !m.File.Binary || len(m.File.Data) != 0 {
		t.Fatalf("binary read: %+v", m.File)
	}
	big := filepath.Join(root, "big")
	if err := os.WriteFile(big, bytes.Repeat([]byte("x"), vaultfs.MaxEditableBytes+1), 0o600); err != nil {
		t.Fatal(err)
	}
	m = call(b, editbackend.Request{Op: editbackend.OpRead, Path: "big"})
	if m.Err != "" || !m.File.TooLarge || len(m.File.Data) != 0 {
		t.Fatalf("oversized read: too_large=%v len=%d", m.File.TooLarge, len(m.File.Data))
	}
	tooBig := make([]byte, editbackend.MaxWriteBytes+1)
	if m = call(b, editbackend.Request{Op: editbackend.OpWrite, Path: "a.txt", Data: tooBig}); m.Err == "" {
		t.Fatal("oversized write accepted")
	}
}

func TestBackendChownOnlyLoginUsers(t *testing.T) {
	b, _, _ := editingBackend(t)
	m := call(b, editbackend.Request{Op: editbackend.OpChown, Path: "a.txt", UID: 4242424, GID: 4242424})
	if m.Err == "" || !strings.Contains(m.Err, "not a login user") {
		t.Fatalf("chown to an arbitrary uid: err %q", m.Err)
	}
}

func TestBackendRefusesEverythingAfterTheDaemonEnded(t *testing.T) {
	b, out := newTestBackend(t)
	b2, _, _ := editingBackend(t)
	b2.out = b.out
	b2.markEnded()
	if !strings.Contains(out.String(), `"kind":"ended"`) {
		t.Fatalf("no ended notice sent: %s", out.String())
	}
	if m := call(b2, editbackend.Request{Op: editbackend.OpWrite, Path: "a.txt", Data: []byte("late")}); m.Err == "" {
		t.Fatal("write after the daemon ended the session was accepted")
	}
}

func TestBackendEndClosesTheSession(t *testing.T) {
	b, _, _ := editingBackend(t)
	m := b.end(7)
	if m.ID != 7 || m.Kind != editbackend.KindReply {
		t.Fatalf("end reply: %+v", m)
	}
	if b.vault != nil || b.state != stDone {
		t.Fatalf("end left vault=%v state=%v", b.vault, b.state)
	}
	if r := call(b, editbackend.Request{Op: editbackend.OpRead, Path: "a.txt"}); r.Err == "" {
		t.Fatal("read after end accepted")
	}
	// Our own close is not the daemon ending the session.
	b.markEnded()
	if b.ended.Load() {
		t.Fatal("cleanup's own disconnect was reported as the daemon ending the session")
	}
}

func TestCheckSecretRefusesLineBreaks(t *testing.T) {
	for _, pw := range []string{"", "a\nSELECT /x", "a\rb", "a\x00b", "\xff\xfe"} {
		if checkSecret(pw) == nil {
			t.Errorf("checkSecret(%q) accepted", pw)
		}
	}
	if err := checkSecret("Corr3ct-Horse!"); err != nil {
		t.Errorf("valid password refused: %v", err)
	}
}

func TestBackendServeNeedsHelloFirst(t *testing.T) {
	b, _ := newTestBackend(t)
	in := strings.NewReader(`{"id":1,"op":"list","path":"."}` + "\n")
	if err := b.serve(editbackend.NewReader(in)); err == nil {
		t.Fatal("a stream not starting with hello was served")
	}
}

func TestBackendAuthRefusedOffline(t *testing.T) {
	b, _ := newTestBackend(t)
	if m := call(b, editbackend.Request{Op: editbackend.OpAuth, Password: "x"}); m.Err == "" ||
		!strings.Contains(m.Err, "live mode") {
		t.Fatalf("auth without the control socket: err %q", m.Err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBackendAuthTimeoutBounds(t *testing.T) {
	b, _ := newTestBackend(t)
	b.live = true
	for _, sec := range []int64{-1, 86401, 1 << 62} {
		m := call(b, editbackend.Request{Op: editbackend.OpAuth, Password: "x", TimeoutSec: sec})
		if !strings.Contains(m.Err, "timeout") {
			t.Errorf("timeout %d: err %q, want refused before dialing", sec, m.Err)
		}
	}
}

// record_start resolves and watches caller-named paths as this binary: never before AUTH. The
// later ops only follow their predecessor, so no step (plan, apply, delete) can be skipped to.
func TestBackendRecordOpsFollowTheWorkflow(t *testing.T) {
	b, _ := newTestBackend(t)
	ops := []string{editbackend.OpRecordStart, editbackend.OpRecordStatus, editbackend.OpRecordStop,
		editbackend.OpRecordPlan, editbackend.OpRecordApply, editbackend.OpRecordDeleteBackup}
	for _, op := range ops {
		if m := call(b, editbackend.Request{Op: op, Resources: []string{"/home/u/.ssh"}, Path: "/x"}); m.Err == "" {
			t.Errorf("%s before auth: accepted", op)
		}
	}
	b.state = stAuthed
	for _, op := range ops[1:] {
		if m := call(b, editbackend.Request{Op: op, Path: "/x"}); m.Err == "" {
			t.Errorf("%s before record_start: accepted", op)
		}
	}
	b.state = stRecording
	for _, op := range []string{editbackend.OpRecordPlan, editbackend.OpRecordApply, editbackend.OpRecordDeleteBackup,
		editbackend.OpConfigPut, editbackend.OpSelect} {
		if m := call(b, editbackend.Request{Op: op, Path: "/x"}); m.Err == "" {
			t.Errorf("%s while recording: accepted", op)
		}
	}
	b.state = stRecorded
	if m := call(b, editbackend.Request{Op: editbackend.OpRecordDeleteBackup, Path: "/x"}); m.Err == "" {
		t.Error("record_delete_backup before an apply: accepted")
	}
}
