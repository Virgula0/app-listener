//go:build gui

package gui

import (
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/recorder"
)

// fakeLink answers the backend protocol from memory and records every request.
type fakeLink struct {
	mu     sync.Mutex
	live   bool
	files  map[string]string
	reqs   []editbackend.Request
	trust  []editbackend.TrustRow
	events chan editbackend.Message
	ended  chan string
	done   chan struct{}
	closed bool
}

func newFakeLink(live bool) *fakeLink {
	return &fakeLink{live: live, files: map[string]string{
		"config":          "Host *\n  IdentitiesOnly yes\n",
		"id_ed25519":      "-----BEGIN OPENSSH PRIVATE KEY-----\n",
		"authorized_keys": "ssh-ed25519 AAAA… user@host\n",
	}, events: make(chan editbackend.Message, 8), ended: make(chan string, 1), done: make(chan struct{})}
}

func (f *fakeLink) Hello() editbackend.Hello {
	return editbackend.Hello{Proto: editbackend.Proto, Live: f.live, HasPassword: true}
}

func (f *fakeLink) Call(req *editbackend.Request) (editbackend.Message, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reqs = append(f.reqs, *req)
	m := editbackend.Message{}
	switch req.Op {
	case editbackend.OpPasswordState:
		m.Password = &editbackend.PwState{Configured: true, Origin: "cli", DaemonRunning: true}
	case editbackend.OpAuth, editbackend.OpOfflineScan:
		m.Resources = []string{"/home/alice/.ssh", "/home/alice/.gnupg"}
	case editbackend.OpList:
		now := time.Date(2026, 10, 17, 9, 30, 0, 0, time.UTC)
		names := make([]string, 0, len(f.files))
		for n := range f.files {
			names = append(names, n)
		}
		slices.Sort(names)
		for _, n := range names {
			m.Entries = append(m.Entries, editbackend.Entry{Name: n, Size: int64(len(f.files[n])), Mode: 0o600,
				ModeStr: "-rw-------", UID: 1000, GID: 1000, ModTime: now})
		}
		m.Entries = append(m.Entries, editbackend.Entry{Name: "agent", Dir: true, Mode: 0o700, ModeStr: "drwx------",
			UID: 1000, GID: 1000, ModTime: now})
	case editbackend.OpRead:
		m.File = &editbackend.File{Data: []byte(f.files[req.Path]), Size: int64(len(f.files[req.Path]))}
	case editbackend.OpWrite:
		f.files[req.Path] = string(req.Data)
	case editbackend.OpConfigGet:
		m.Data = []byte("[/home/alice/.ssh]\nmode = whitelist\n")
	case editbackend.OpForwardCheck:
		for _, b := range req.Binaries {
			m.Binaries = append(m.Binaries, editbackend.BinaryCheck{Path: b, SystemPlaced: true})
		}
	case editbackend.OpTrustList:
		m.Trust = f.trust
	case editbackend.OpTrustInspect:
		for _, b := range req.Binaries {
			m.Trust = append(m.Trust, editbackend.TrustRow{Line: b, Resolved: b, SHA256: strings.Repeat("c", 64)})
		}
	case editbackend.OpTrustConfirm, editbackend.OpTrustAdd:
		m.TrustDone = &editbackend.TrustDone{Reloaded: !req.NoReload}
		for _, it := range req.Trust {
			m.TrustDone.Confirmed = append(m.TrustDone.Confirmed, it.Line)
		}
	}
	return m, nil
}

func (f *fakeLink) ops() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, 0, len(f.reqs))
	for _, r := range f.reqs {
		out = append(out, r.Op)
	}
	return out
}

func (f *fakeLink) Events() <-chan editbackend.Message { return f.events }
func (f *fakeLink) Ended() <-chan string               { return f.ended }
func (f *fakeLink) Done() <-chan struct{}              { return f.done }
func (f *fakeLink) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.closed {
		f.closed = true
		close(f.done)
	}
}

func protectedShell(t *testing.T, link *fakeLink) (*shell, *protectedPage) {
	t.Helper()
	a := test.NewApp()
	t.Cleanup(a.Quit)
	a.Settings().SetTheme(&appTheme{variant: theme.VariantDark})
	w := a.NewWindow("app-listener")
	s := newShell(a, w, &Options{Version: "dev"})
	s.syncBackground = true
	s.startBackend = func() (backendLink, error) { return link, nil }
	w.SetContent(s.view())
	w.Resize(fyne.NewSize(1280, 820))
	s.show(pageProtected)
	return s, s.protect
}

func TestProtectedEditFlow(t *testing.T) {
	link := newFakeLink(true)
	_, p := protectedShell(t, link)
	if p.pw == nil || !p.pw.Configured || !p.hello.Live {
		t.Fatalf("probe: hello=%+v pw=%+v", p.hello, p.pw)
	}
	link.reqs = nil
	p.link = link
	p.openEdit(link, editbackend.OpSelect, "/home/alice/.ssh", true)
	if p.leave == nil || p.stack.Objects[0] == p.over {
		t.Fatal("the edit view is not shown")
	}
	if ops := link.ops(); len(ops) < 2 || ops[0] != editbackend.OpSelect || ops[1] != editbackend.OpList {
		t.Fatalf("open: %v", ops)
	}
	view := newEditView(p, link, "/home/alice/.ssh", true, false)
	view.start()
	view.onSelect("config")
	if view.editor.Text != "Host *\n  IdentitiesOnly yes\n" || view.editor.Disabled() {
		t.Fatalf("editor = %q disabled=%v", view.editor.Text, view.editor.Disabled())
	}
	test.Type(view.editor, "# edited\n")
	if !view.dirty {
		t.Fatal("typing did not mark the file dirty")
	}
	view.editor.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierShortcutDefault})
	if !strings.HasPrefix(link.files["config"], "# edited\n") || view.dirty {
		t.Fatalf("Ctrl+S did not save: %q dirty=%v", link.files["config"], view.dirty)
	}
	view.onEnded("the daemon ended the live session")
	if !view.editor.Disabled() || !view.save.Disabled() {
		t.Fatal("editing still enabled after the daemon ended the session")
	}
	view.end()
	// After end the overview probes again (password_state on a fresh backend).
	ops := link.ops()
	if i := slices.Index(ops, editbackend.OpEnd); i < 0 || slices.Contains(ops[i+1:], editbackend.OpWrite) {
		t.Fatalf("end not sent after the edits: %v", ops)
	}
	if p.link != nil {
		t.Fatal("the session is still open after end")
	}
}

func TestForwardRefusesOwnExecutable(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if !isOwnExecutable(exe) {
		t.Fatalf("%s not recognized as this process's executable", exe)
	}
	if isOwnExecutable("/bin/sh") {
		t.Fatal("/bin/sh taken for this process's executable")
	}
	link := newFakeLink(true)
	_, p := protectedShell(t, link)
	link.reqs = nil
	p.startForward()
	f := &forwardForm{p: p, bins: []string{exe}, events: newCheckGrid(forwardEvents)}
	f.rule = newRadioGroup([]string{ruleAdmit, ruleDeny}, nil)
	f.rule.SetSelected(ruleAdmit)
	f.submit()
	if slices.Contains(link.ops(), editbackend.OpForwardCheck) || p.link != nil {
		t.Fatalf("own executable reached the backend: %v", link.ops())
	}
}

func TestRecordReviewExtendsAnExistingSection(t *testing.T) {
	link := newFakeLink(true)
	_, p := protectedShell(t, link)
	ssh, fresh := "/home/alice/.ssh", "/home/alice/.config/tool"
	rec := &recorder.Recording{
		Targets: []recorder.Target{{Path: ssh, IsDir: true, Existing: &recorder.Existing{Section: ssh,
			Allow: []recorder.ExistingRule{{Line: "/usr/bin/ssh"}}}}, {Path: fresh, IsDir: true}},
		Binaries: []recorder.Binary{
			{Line: "/usr/bin/ssh", Exe: "/usr/bin/ssh", Targets: []string{ssh}, Whitelisted: []string{ssh}},
			{Line: "/usr/bin/git", Exe: "/usr/bin/git", Targets: []string{ssh, fresh}, Suggested: true},
		},
	}
	r := newRecordReview(p, link, rec, map[string]string{ssh: "already in daemon.conf"})
	if r.bins["/usr/bin/ssh"] != nil || r.bins["/usr/bin/git"] == nil {
		t.Fatalf("only a binary with a line to add is selectable: %v", r.bins)
	}
	if r.encrypt[ssh] != nil || r.encrypt[fresh] == nil {
		t.Fatalf("encryption is offered for new paths only: %v", r.encrypt)
	}
	if sel := r.selection(); !slices.Equal(sel.Binaries, []string{"/usr/bin/git"}) {
		t.Fatalf("selection = %+v", sel)
	}
	if notices := existingNotices(rec.Targets); len(notices) != 2 {
		t.Fatalf("one warning and the current whitelist per existing target: %d objects", len(notices))
	}
}

func TestProtectedLiveCardsGated(t *testing.T) {
	for _, live := range []bool{true, false} {
		_, p := protectedShell(t, newFakeLink(live))
		if len(p.live) != 3 {
			t.Fatalf("live cards: %d", len(p.live))
		}
		for _, g := range p.live {
			if g.btn.Disabled() == live || (g.note.Text == "") != live {
				t.Errorf("live=%v: disabled=%v note=%q", live, g.btn.Disabled(), g.note.Text)
			}
		}
	}
}
