//go:build gui

package gui

import (
	"errors"
	"github.com/Virgula0/app-listener/internal/editbackend"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
	"github.com/Virgula0/app-listener/internal/wizard"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// TestSnapshots renders the window off-screen into $APPLISTENER_GUI_SNAPSHOTS (a directory) for
// design review; skipped when unset.
func TestSnapshots(t *testing.T) {
	dir := os.Getenv("APPLISTENER_GUI_SNAPSHOTS")
	if dir == "" {
		t.Skip("APPLISTENER_GUI_SNAPSHOTS not set")
	}
	for _, v := range []struct {
		name    string
		variant fyne.ThemeVariant
	}{{"dark", theme.VariantDark}, {"light", theme.VariantLight}} {
		a := test.NewApp()
		a.Settings().SetTheme(&appTheme{variant: v.variant})
		if v.variant == theme.VariantLight {
			a.Preferences().SetString("theme", "light")
		}
		w := a.NewWindow("app-listener")
		s := newShell(a, w, &Options{Version: "dev"})
		s.syncBackground = true
		fe := &snapEngines{}
		s.files = workbench.NewFiles(fe.engines())
		ne := &snapNet{}
		s.network = workbench.NewNetwork(ne.engines())
		snapLink := newFakeLink(true)
		s.startBackend = func() (backendLink, error) { return snapLink, nil }
		w.SetContent(s.view())
		w.Resize(fyne.NewSize(1280, 820))
		s.setStatus(sysStatus{root: true, lsmErr: errors.New("bpf not in /sys/kernel/security/lsm")})
		for _, p := range pageDefs {
			s.show(p.id)
			save(t, filepath.Join(dir, v.name+"-"+p.title+".png"), w)
		}

		s.setStatus(sysStatus{root: true})
		fillFiles(t, s, fe)
		fp := newFilesPage(s, s.files)
		s.views[pageFiles] = container.NewScroll(fp.view())
		s.show(pageFiles)
		fp.reloadTargets(s.files.Targets()[0])
		for i, tab := range []string{"Activity", "Processes", "Guard"} {
			fp.current.tabs.SelectIndex(i)
			fp.current.refresh(true)
			save(t, filepath.Join(dir, v.name+"-Files-"+tab+".png"), w)
		}
		fp.current.tabs.SelectIndex(0)
		fp.current.activity.table.Select(widget.TableCellID{Row: 1, Col: 2})
		save(t, filepath.Join(dir, v.name+"-Files-RowSelected.png"), w)
		fp.toggleTargets()
		save(t, filepath.Join(dir, v.name+"-Files-TargetsCollapsed.png"), w)
		fp.toggleTargets()

		fillNetwork(t, s, ne)
		np := newNetworkPage(s, s.network)
		s.views[pageNetwork] = container.NewScroll(np.view())
		s.show(pageNetwork)
		for i, tab := range []string{"Watched", "Activity", "Processes", "Guard"} {
			np.tabs.SelectIndex(i)
			np.refresh(true)
			save(t, filepath.Join(dir, v.name+"-Network-"+tab+".png"), w)
		}

		s.toggleSidebar()
		save(t, filepath.Join(dir, v.name+"-Collapsed.png"), w)
		s.toggleSidebar()

		w.Resize(fyne.NewSize(860, 600))
		s.show(pageHome)
		save(t, filepath.Join(dir, v.name+"-Narrow-Home.png"), w)
		s.show(pageNetwork)
		np.tabs.SelectIndex(1)
		s.exportTable("network-events", np.feed.table)
		save(t, filepath.Join(dir, v.name+"-Narrow-Export.png"), w)
		w.Canvas().Overlays().Top().Hide()
		np.tabs.SelectIndex(0)
		np.watch.bins.pick()
		save(t, filepath.Join(dir, v.name+"-Narrow-AddExecutable.png"), w)
		w.Canvas().Overlays().Top().Hide()
		s.show(pageFiles)
		fp.addTarget()
		save(t, filepath.Join(dir, v.name+"-Narrow-AddTarget.png"), w)
		w.Canvas().Overlays().Top().Hide()
		w.Resize(fyne.NewSize(1280, 820))

		s.pickPath(pickFile, "File to monitor", pickerTree(t), func(string) {})
		s.picker.deep = true
		test.Type(s.picker.search, "alpha")
		s.picker.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
		save(t, filepath.Join(dir, v.name+"-Picker.png"), w)
		s.picker.close()

		snapInstall(t, s, filepath.Join(dir, v.name+"-Install-"))
		snapUpdate(t, s, filepath.Join(dir, v.name+"-Update-"))
		snapProtected(t, s, snapLink, filepath.Join(dir, v.name+"-Protected-"))
		w.Close()
	}
}

// TestIconGallery renders every embedded icon (and the logo) for a rasterizer check.
func TestIconGallery(t *testing.T) {
	dir := os.Getenv("APPLISTENER_GUI_SNAPSHOTS")
	if dir == "" {
		t.Skip("APPLISTENER_GUI_SNAPSHOTS not set")
	}
	a := test.NewApp()
	a.Settings().SetTheme(&appTheme{variant: theme.VariantDark})
	entries, _ := assets.ReadDir("assets/icons")
	grid := container.NewGridWithColumns(6)
	for _, e := range entries {
		n := strings.TrimSuffix(e.Name(), ".svg")
		img := canvas.NewImageFromResource(icon(n))
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSquareSize(32))
		grid.Add(container.NewVBox(img, widget.NewLabel(n)))
	}
	logo := newLogo(96)
	w := a.NewWindow("g")
	w.SetContent(container.NewBorder(logo, nil, nil, nil, grid))
	w.Resize(fyne.NewSize(1200, 1100))
	save(t, filepath.Join(dir, "gallery.png"), w)
}

type snapEngines struct {
	monitors []chan ebpf.FileEvent
	guards   []chan guard.GuardEvent
}

type snapMonitor struct{ ch chan ebpf.FileEvent }

func (m *snapMonitor) Start() error                   { return nil }
func (m *snapMonitor) Stop()                          {}
func (m *snapMonitor) Events() <-chan ebpf.FileEvent  { return m.ch }
func (m *snapMonitor) SetEventTypes([]ebpf.EventType) {}

type snapGuard struct{ ch chan guard.GuardEvent }

func (g *snapGuard) Start() error                         { return nil }
func (g *snapGuard) Stop()                                {}
func (g *snapGuard) Events() <-chan guard.GuardEvent      { return g.ch }
func (g *snapGuard) SetBinaries([]ebpf.BinaryEntry) error { return nil }

func (f *snapEngines) engines() workbench.FileEngines {
	return workbench.FileEngines{
		NewMonitor: func(ebpf.Target, bool, int) (repository.MonitorRepository, error) {
			ch := make(chan ebpf.FileEvent, 64)
			f.monitors = append(f.monitors, ch)
			return &snapMonitor{ch: ch}, nil
		},
		NewGuard: func(string, guard.Mode, []ebpf.BinaryEntry, bool, int) (workbench.FileGuard, error) {
			ch := make(chan guard.GuardEvent, 64)
			f.guards = append(f.guards, ch)
			return &snapGuard{ch: ch}, nil
		},
	}
}

// fillFiles adds two targets with a whitelist and a little traffic.
func fillFiles(t *testing.T, s *shell, fe *snapEngines) {
	t.Helper()
	home := t.TempDir()
	ssh := filepath.Join(home, ".ssh")
	gnupg := filepath.Join(home, ".gnupg")
	for _, d := range []string{ssh, gnupg} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	tg, err := s.files.Add(workbench.FileSpec{Path: ssh, Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.files.Add(workbench.FileSpec{Path: gnupg, Recursive: true}); err != nil {
		t.Fatal(err)
	}
	agent := ebpf.BinaryEntry{Path: "/usr/bin/ssh-agent", Comm: "ssh-agent", Hash: [32]byte{0xab, 0xcd, 0x12}}
	if err := tg.Apply(workbench.Policy{Mode: workbench.GuardWhitelist, Binaries: []ebpf.BinaryEntry{agent}}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixNano()
	for i, ev := range []struct {
		comm, path string
		typ        ebpf.EventType
	}{
		{"ssh", "id_ed25519", ebpf.EventOpen}, {"ssh", "id_ed25519", ebpf.EventRead},
		{"ssh-agent", "id_ed25519", ebpf.EventRead}, {"code", "config", ebpf.EventStat},
		{"python3", "known_hosts", ebpf.EventOpen},
	} {
		fe.monitors[0] <- ebpf.FileEvent{PID: uint32(4000 + i), Comm: ev.comm, Path: filepath.Join(ssh, ev.path), //nolint:gosec // test
			Type: ev.typ, Timestamp: now + int64(i)*int64(time.Millisecond)}
	}
	fe.guards[0] <- guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: 4242, Comm: "stealer", Type: ebpf.EventOpen,
		Path: filepath.Join(ssh, "id_ed25519"), Timestamp: now + int64(10*time.Millisecond)}, Blocked: true,
		Exe: guard.GuardInodeKey{Dev: 2049, Ino: 1311}}
	fe.guards[0] <- guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: 4100, Comm: "ssh-agent", Type: ebpf.EventRead,
		Path: filepath.Join(ssh, "id_ed25519"), Timestamp: now + int64(11*time.Millisecond)}}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if c, _ := tg.Counts(); c.Observed+c.Allowed+c.Denied == 7 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("events never reached the target")
}

func save(t *testing.T, path string, w fyne.Window) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, w.Canvas().Capture()); err != nil {
		t.Fatal(err)
	}
}

type snapNet struct {
	mon   chan ebpf.NetEvent
	guard chan networkguard.NetGuardEvent
}

type snapNetMonitor struct{ ch chan ebpf.NetEvent }

func (m *snapNetMonitor) Start() error                      { return nil }
func (m *snapNetMonitor) Stop()                             {}
func (m *snapNetMonitor) Events() <-chan ebpf.NetEvent      { return m.ch }
func (m *snapNetMonitor) SetEventTypes([]ebpf.NetEventType) {}

type snapNetGuard struct {
	ch chan networkguard.NetGuardEvent
}

func (g *snapNetGuard) Start() error                              { return nil }
func (g *snapNetGuard) Stop()                                     {}
func (g *snapNetGuard) Events() <-chan networkguard.NetGuardEvent { return g.ch }

func (f *snapNet) engines() workbench.NetEngines {
	return workbench.NetEngines{
		NewMonitor: func([]ebpf.BinaryEntry) (repository.NetworkMonitorRepository, error) {
			f.mon = make(chan ebpf.NetEvent, 64)
			return &snapNetMonitor{ch: f.mon}, nil
		},
		NewGuard: func(workbench.NetGuardConfig, repository.NetworkGuardRepository) (repository.NetworkGuardRepository, error) {
			f.guard = make(chan networkguard.NetGuardEvent, 64)
			return &snapNetGuard{ch: f.guard}, nil
		},
		Infra: func() ([]string, error) { return []string{"/usr/lib/systemd/systemd-resolved"}, nil },
		Hash:  func(p string) (ebpf.BinaryEntry, error) { return snapEntry(p), nil },
		Self:  func() (string, error) { return "/usr/local/bin/app-listener-gui", nil },
	}
}

func snapEntry(path string) ebpf.BinaryEntry {
	return ebpf.EntryWithHash(path, [32]byte{byte(len(path)), 0x5e, 0x77})
}

// fillNetwork watches two programs, starts a whitelist and feeds a little traffic.
func fillNetwork(t *testing.T, s *shell, f *snapNet) {
	t.Helper()
	curl, firefox := snapEntry("/usr/bin/curl"), snapEntry("/usr/lib/firefox/firefox")
	if err := s.network.Watch([]ebpf.BinaryEntry{curl, firefox}); err != nil {
		t.Fatal(err)
	}
	if err := s.network.Apply(workbench.NetPolicy{Mode: workbench.GuardWhitelist, Binaries: []ebpf.BinaryEntry{firefox},
		AutoInfra: true, Throttle: true}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixNano()
	at := func(i int) int64 { return now + int64(i)*int64(time.Millisecond) }
	f.mon <- ebpf.NetEvent{PID: 5100, Comm: "firefox", Type: ebpf.NetDNS, Protocol: 17,
		SrcAddr: "10.0.0.7:51122", DstAddr: "127.0.0.53:53", Timestamp: at(0)}
	f.mon <- ebpf.NetEvent{PID: 5100, Comm: "firefox", Type: ebpf.NetConnect, Protocol: 6,
		SrcAddr: "10.0.0.7:40112", DstAddr: "140.82.121.4:443", Timestamp: at(1)}
	f.mon <- ebpf.NetEvent{PID: 5100, Comm: "firefox", Type: ebpf.NetSend, Protocol: 6, Size: 517,
		SrcAddr: "10.0.0.7:40112", DstAddr: "140.82.121.4:443", Timestamp: at(2)}
	f.mon <- ebpf.NetEvent{PID: 5201, Comm: "curl", Type: ebpf.NetConnect, Protocol: 6,
		SrcAddr: "10.0.0.7:40200", DstAddr: "93.184.215.14:80", Timestamp: at(3)}
	f.guard <- networkguard.NetGuardEvent{NetEvent: ebpf.NetEvent{PID: 5201, Comm: "curl", Type: ebpf.NetConnect,
		Protocol: 6, DstAddr: "93.184.215.14:80", Timestamp: at(4)}, Blocked: true}
	f.guard <- networkguard.NetGuardEvent{NetEvent: ebpf.NetEvent{PID: 5300, Comm: "node", Type: ebpf.NetConnect,
		Protocol: 6, DstAddr: "104.16.0.35:443", Timestamp: at(5)}, Blocked: true, Reason: "ld-preload"}
	f.guard <- networkguard.NetGuardEvent{NetEvent: ebpf.NetEvent{PID: 5100, Comm: "firefox", Type: ebpf.NetConnect,
		Protocol: 6, DstAddr: "140.82.121.4:443", Timestamp: at(6)}}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if c, _ := s.network.Counts(); c.Observed+c.Allowed+c.Denied == 7 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("network events never arrived")
}

// snapInstall renders the Install page with a staged build, a run in progress, a failed run and a
// prompt.
func snapInstall(t *testing.T, s *shell, prefix string) {
	t.Helper()
	ip := newInstallPage(s)
	s.install = ip
	s.views[pageInstall] = ip.view()
	s.show(pageInstall)
	ip.bin = &stagedBinary{path: "/tmp/app-listener-gui-81723/app-listener", sha: strings.Repeat("3fa9c0d1", 8),
		tag: "v0.3.0", guiPath: "/tmp/app-listener-gui-81723/app-listener-gui", guiSha: strings.Repeat("9be0a412", 8)}
	ip.present = installedBinary{present: true, version: "v0.3.0", sha: "ff", guiVersion: "v0.2.9"}
	ip.daemonUp = true
	ip.renderBinary()
	s.win.Resize(fyne.NewSize(1280, 2450)) // the whole overview, actions and uninstall included
	save(t, prefix+"Overview.png", s.win)
	s.win.Resize(fyne.NewSize(1280, 820))

	ip.run.start("Install wizard", "shield-check")
	ip.show(ip.run.root)
	for _, st := range []string{"Stopping the daemon", "Deploying the binary and the master key",
		"Choosing what to protect", "Encrypting"} {
		ip.run.stage(st)
	}
	now := time.Now()
	for i, l := range []logLine{
		{level: log.InfoLevel, text: "daemon stopped"},
		{level: log.InfoLevel, text: "binary is installed: /usr/local/sbin/app-listener"},
		{level: log.WarnLevel, text: "~/.config/Code: general-purpose tools whitelisted (python3)"},
		{level: log.InfoLevel, text: "migrating /home/user/.ssh into an encrypted directory"},
	} {
		l.at = now.Add(time.Duration(i) * time.Second)
		ip.run.appendLog(l)
	}
	ip.run.progress("Encrypting /home/user/.ssh", 0.42)
	save(t, prefix+"Run.png", s.win)

	ip.run.finish(errors.New("filesystem /home is not fscrypt-ready: run `fscrypt setup /home`"), false)
	save(t, prefix+"Failed.png", s.win)

	ui := &installUI{s: s, run: ip.run}
	ui.showMultiSelect("Which users should be protected?", "Their catalog folders are probed next.",
		[]wizard.Choice{{Label: "alice (/home/alice)", Selected: true}, {Label: "bob (/home/bob)"}},
		func(multiAnswer) {})
	save(t, prefix+"Prompt.png", s.win)
	s.win.Canvas().Overlays().Top().Hide()
}

// snapUpdate renders the Update page with a newer release, a run and the changelog prompt.
func snapUpdate(t *testing.T, s *shell, prefix string) {
	t.Helper()
	up := newUpdatePage(s)
	s.update = up
	s.views[pageUpdate] = up.view()
	s.show(pageUpdate)
	up.present = installedBinary{present: true, version: "v0.3.0"}
	up.daemonUp, up.loaded = true, true
	up.latest[false] = &latestRelease{tag: "v0.4.0", newer: true}
	up.renderStatus()
	s.setUpdateAvailable("v0.4.0", true)
	save(t, prefix+"Overview.png", s.win)
	s.setUpdateAvailable("", false)

	up.run.start("Update: Latest stable release", "arrow-clockwise")
	up.show(up.run.root)
	for _, st := range []string{"Checking the installed version", "Looking up the latest stable release",
		"Downloading and verifying v0.4.0"} {
		up.run.stage(st)
	}
	up.run.appendLog(logLine{level: log.InfoLevel, at: time.Now(), text: "latest stable release: v0.4.0"})
	up.run.progress("Downloading app-listener", 0.63)
	save(t, prefix+"Run.png", s.win)

	ui := &installUI{s: s, run: up.run}
	ui.showText("Changelog up to v0.4.0", "## v0.4.0\n\n- gui: updater page\n- guard: fixes", func() {})
	save(t, prefix+"Changelog.png", s.win)
	s.win.Canvas().Overlays().Top().Hide()
}

func snapProtected(t *testing.T, s *shell, link *fakeLink, prefix string) {
	t.Helper()
	s.show(pageProtected)
	p := s.protect
	save(t, prefix+"Overview.png", s.win)

	p.link = link
	p.authenticate(link, "Edit protected files", func([]string) {})
	save(t, prefix+"Password.png", s.win)
	s.win.Canvas().Overlays().Top().Hide()

	v := newEditView(p, link, "/home/alice/.ssh", true, false)
	p.show(v.root)
	v.start()
	v.tree.Select("config")
	test.Type(v.editor, "# edited\n")
	save(t, prefix+"Edit.png", s.win)
	v.onEnded("the daemon ended the live session (idle timeout or reload) — saves made before it are kept, " +
		"later ones were not written")
	save(t, prefix+"Edit-Ended.png", s.win)

	off := newEditView(p, link, "/home/alice/.gnupg", false, false)
	p.show(off.root)
	off.start()
	save(t, prefix+"Edit-Offline.png", s.win)

	p.show(p.over)
	p.link = nil
	p.startForward()
	save(t, prefix+"Forward-Form.png", s.win)
	fv := newForwardView(p, link, []string{"/home/alice/.ssh"}, []string{"/usr/bin/vim"}, true, []string{"OPEN", "READ"})
	p.show(fv.root)
	for i, op := range []string{"OPEN", "READ", "STAT", "WRITE"} {
		fv.add(&editbackend.Event{Time: time.Date(2026, 10, 17, 9, 30, i, 0, time.UTC), PID: 4242, UID: 1000,
			Comm: "vim", Op: op, Path: "/home/alice/.ssh/config", Allowed: op != "WRITE"})
	}
	save(t, prefix+"Forward-Active.png", s.win)

	cv := newConfigView(p, link, []byte("[/home/alice/.ssh]\nmode = whitelist\nbinaries = /usr/bin/ssh\n"))
	p.show(cv.root)
	save(t, prefix+"Config.png", s.win)
	p.show(p.over)
}
