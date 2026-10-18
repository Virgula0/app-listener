//go:build gui && e2e

package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// reel saves numbered frames plus an ffmpeg concat list giving each one's duration.
type reel struct {
	t    *testing.T
	w    fyne.Window
	dir  string
	n    int
	list *os.File
}

func (r *reel) shot(d time.Duration) {
	r.t.Helper()
	name := fmt.Sprintf("frame-%04d.png", r.n)
	r.n++
	save(r.t, filepath.Join(r.dir, name), r.w)
	if _, err := fmt.Fprintf(r.list, "file '%s'\nduration %.3f\n", name, d.Seconds()); err != nil {
		r.t.Fatal(err)
	}
}

// close repeats the last frame: the concat demuxer ignores the last entry's duration.
func (r *reel) close() {
	_, _ = fmt.Fprintf(r.list, "file 'frame-%04d.png'\n", r.n-1)
	_ = r.list.Close()
}

func waitTotal(t *testing.T, total func() int, n int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for total() < n {
		if time.Now().After(deadline) {
			t.Fatalf("only %d of %d events arrived", total(), n)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// TestDemoReel renders the README's GUI walkthrough as frames into $APPLISTENER_GUI_DEMO,
// assembled with ffmpeg (docs/development.md).
func TestDemoReel(t *testing.T) {
	dir := e2eDir(t, "APPLISTENER_GUI_DEMO")
	a := test.NewApp()
	a.Settings().SetTheme(&appTheme{variant: theme.VariantDark})
	w := a.NewWindow("app-listener")
	s := newShell(a, w, &Options{Version: "v0.4.0"})
	s.syncBackground = true
	fe := &snapEngines{}
	s.files = workbench.NewFiles(fe.engines())
	ne := &snapNet{}
	s.network = workbench.NewNetwork(ne.engines())
	link := newFakeLink(true)
	s.startBackend = func() (backendLink, error) { return link, nil }
	w.SetContent(s.view())
	w.Resize(fyne.NewSize(1280, 820))
	s.setStatus(sysStatus{root: true, daemon: true})

	list, err := os.Create(filepath.Join(dir, "frames.txt"))
	if err != nil {
		t.Fatal(err)
	}
	r := &reel{t: t, w: w, dir: dir, list: list}
	defer r.close()

	s.show(pageHome)
	r.shot(2500 * time.Millisecond)
	demoFiles(t, s, fe, r)
	demoNetwork(t, s, ne, r)
	demoProtected(t, s, link, r)
	demoInstall(s, r)
}

// demoFiles: a package's node reads the SSH key unnoticed, then a whitelist denies it.
func demoFiles(t *testing.T, s *shell, fe *snapEngines, r *reel) {
	t.Helper()
	home := t.TempDir()
	ssh, gnupg := filepath.Join(home, ".ssh"), filepath.Join(home, ".gnupg")
	for _, d := range []string{ssh, gnupg} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	tg, err := s.files.Add(workbench.FileSpec{Path: ssh, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	tg2, err := s.files.Add(workbench.FileSpec{Path: gnupg, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	tg.Spec.Path, tg2.Spec.Path = "/home/alice/.ssh", "/home/alice/.gnupg" // display-only once running
	total := func() int { c, _ := tg.Counts(); return int(c.Observed + c.Allowed + c.Denied) }

	fp := newFilesPage(s, s.files)
	s.views[pageFiles] = container.NewScroll(fp.view())
	s.show(pageFiles)
	fp.reloadTargets(tg)
	v := fp.current
	r.shot(1200 * time.Millisecond)

	start := time.Now()
	at := func(i int) int64 { return start.Add(time.Duration(i) * 700 * time.Millisecond).UnixNano() }
	key := "/home/alice/.ssh/id_ed25519"
	for i, ev := range []struct {
		comm, path string
		typ        ebpf.EventType
	}{
		{"ssh", key, ebpf.EventOpen}, {"ssh", key, ebpf.EventRead}, {"ssh-agent", key, ebpf.EventRead},
		{"code", "/home/alice/.ssh/config", ebpf.EventStat}, {"node", key, ebpf.EventOpen},
		{"node", key, ebpf.EventRead},
	} {
		fe.monitors[0] <- ebpf.FileEvent{PID: uint32(4000 + i), Comm: ev.comm, Path: ev.path, //nolint:gosec // test
			Type: ev.typ, Timestamp: at(i)}
		waitTotal(t, total, i+1)
		v.refresh(true)
		r.shot(650 * time.Millisecond)
	}
	r.shot(1600 * time.Millisecond)

	v.tabs.SelectIndex(2)
	r.shot(1000 * time.Millisecond)
	v.policy.ed.setMode(workbench.GuardWhitelist)
	r.shot(1200 * time.Millisecond)
	bins := []ebpf.BinaryEntry{
		{Path: "/usr/bin/ssh", Comm: "ssh", Hash: [32]byte{0x5c, 0x1e, 0x09}},
		{Path: "/usr/bin/ssh-agent", Comm: "ssh-agent", Hash: [32]byte{0xab, 0xcd, 0x12}},
	}
	v.policy.ed.bins.set(bins, nil)
	v.policy.ed.changed()
	r.shot(1500 * time.Millisecond)
	if err := tg.Apply(v.policy.draft()); err != nil {
		t.Fatal(err)
	}
	v.policy.refresh()
	v.refresh(true)
	r.shot(1800 * time.Millisecond)

	v.tabs.SelectIndex(0)
	g := fe.guards[len(fe.guards)-1]
	for i, ev := range []struct {
		comm    string
		typ     ebpf.EventType
		blocked bool
	}{
		{"node", ebpf.EventOpen, true}, {"ssh", ebpf.EventOpen, false}, {"ssh", ebpf.EventRead, false},
		{"node", ebpf.EventOpen, true},
	} {
		g <- guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: uint32(4100 + i), Comm: ev.comm, Path: key, //nolint:gosec // test
			Type: ev.typ, Timestamp: at(8 + i)}, Blocked: ev.blocked, Exe: guard.GuardInodeKey{Dev: 2049, Ino: uint64(1311 + i)}}
		waitTotal(t, total, 7+i)
		v.refresh(true)
		r.shot(800 * time.Millisecond)
	}
	r.shot(2200 * time.Millisecond)
}

// demoNetwork streams connections while firefox alone is whitelisted.
func demoNetwork(t *testing.T, s *shell, f *snapNet, r *reel) {
	t.Helper()
	curl, firefox := snapEntry("/usr/bin/curl"), snapEntry("/usr/lib/firefox/firefox")
	if err := s.network.Watch([]ebpf.BinaryEntry{curl, firefox}); err != nil {
		t.Fatal(err)
	}
	if err := s.network.Apply(workbench.NetPolicy{Mode: workbench.GuardWhitelist, Binaries: []ebpf.BinaryEntry{firefox},
		AutoInfra: true, Throttle: true}); err != nil {
		t.Fatal(err)
	}
	total := func() int { c, _ := s.network.Counts(); return int(c.Observed + c.Allowed + c.Denied) }
	np := newNetworkPage(s, s.network)
	s.views[pageNetwork] = container.NewScroll(np.view())
	s.show(pageNetwork)
	np.tabs.SelectIndex(1)
	np.refresh(true)
	r.shot(1000 * time.Millisecond)

	start := time.Now()
	at := func(i int) int64 { return start.Add(time.Duration(i) * 600 * time.Millisecond).UnixNano() }
	mon := func(i int, comm string, typ ebpf.NetEventType, proto uint32, src, dst string, size uint32) {
		f.mon <- ebpf.NetEvent{PID: uint32(5100 + len(comm)), Comm: comm, Type: typ, Protocol: proto, //nolint:gosec // test
			SrcAddr: src, DstAddr: dst, Size: size, Timestamp: at(i)}
	}
	deny := func(i int, comm, dst, reason string) {
		f.guard <- networkguard.NetGuardEvent{NetEvent: ebpf.NetEvent{PID: uint32(5300 + i), Comm: comm, //nolint:gosec // test
			Type: ebpf.NetConnect, Protocol: 6, DstAddr: dst, Timestamp: at(i)}, Blocked: true, Reason: reason}
	}
	steps := []func(int){
		func(i int) { mon(i, "firefox", ebpf.NetDNS, 17, "10.0.0.7:51122", "127.0.0.53:53", 0) },
		func(i int) { mon(i, "firefox", ebpf.NetConnect, 6, "10.0.0.7:40112", "140.82.121.4:443", 0) },
		func(i int) { mon(i, "firefox", ebpf.NetSend, 6, "10.0.0.7:40112", "140.82.121.4:443", 517) },
		func(i int) { mon(i, "curl", ebpf.NetConnect, 6, "10.0.0.7:40200", "93.184.215.14:80", 0) },
		func(i int) { deny(i, "curl", "93.184.215.14:80", "") },
		func(i int) { deny(i, "node", "104.16.0.35:443", "ld-preload") },
	}
	for i, step := range steps {
		step(i)
		waitTotal(t, total, i+1)
		np.refresh(true)
		r.shot(700 * time.Millisecond)
	}
	r.shot(2000 * time.Millisecond)
}

// demoProtected edits ~/.ssh/config live through the daemon's grant.
func demoProtected(t *testing.T, s *shell, link *fakeLink, r *reel) {
	t.Helper()
	s.show(pageProtected)
	p := s.protect
	r.shot(1800 * time.Millisecond)

	p.link = link
	v := newEditView(p, link, "/home/alice/.ssh", true, false)
	p.show(v.root)
	v.start()
	r.shot(900 * time.Millisecond)
	v.tree.Select("config")
	r.shot(900 * time.Millisecond)
	for _, chunk := range []string{"Host git", "hub.com\n", "  Identity", "File ~/.ssh/", "id_ed25519\n\n"} {
		test.Type(v.editor, chunk)
		r.shot(220 * time.Millisecond)
	}
	r.shot(900 * time.Millisecond)
	v.saveFile()
	r.shot(2000 * time.Millisecond)
}

// demoInstall runs the wizard's steps to completion.
func demoInstall(s *shell, r *reel) {
	ip := newInstallPage(s)
	s.install = ip
	s.views[pageInstall] = ip.view()
	s.show(pageInstall)
	ip.run.start("Install wizard", "shield-check")
	ip.show(ip.run.root)
	logAt := time.Now()
	logf := func(level log.Level, text string) {
		logAt = logAt.Add(1300 * time.Millisecond)
		ip.run.appendLog(logLine{level: level, at: logAt, text: text})
	}
	ip.run.stage("Checking the system")
	logf(log.InfoLevel, "kernel 7.0.11: BPF LSM active, fscrypt ready on /home")
	r.shot(900 * time.Millisecond)
	ip.run.stage("Choosing what to protect")
	logf(log.InfoLevel, "alice: 9 catalog folders found (SSH, GnuPG, AWS, Claude Code, Firefox, ...)")
	r.shot(900 * time.Millisecond)
	ip.run.stage("Encrypting")
	for i, d := range []string{".ssh", ".gnupg", ".aws", ".claude", ".mozilla/firefox"} {
		logf(log.InfoLevel, "migrating /home/alice/"+d+" into an encrypted directory")
		ip.run.progress("Encrypting /home/alice/"+d, float64(i+1)/5)
		r.shot(600 * time.Millisecond)
	}
	ip.run.progressDone()
	ip.run.stage("Deploying the daemon")
	logf(log.InfoLevel, "daemon --check: every guard program verified on this kernel")
	logf(log.InfoLevel, "app-listener-daemon enabled and started")
	r.shot(1000 * time.Millisecond)
	ip.run.finish(nil, false)
	r.shot(3500 * time.Millisecond)
}
