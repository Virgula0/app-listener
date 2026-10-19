//go:build gui

package gui

import (
	"os"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	log "github.com/sirupsen/logrus"
)

func TestJournalLine(t *testing.T) {
	cases := []struct {
		raw   string
		level log.Level
		text  string
	}{
		{`time="2026-10-19 17:29:05" level=warning msg="guard created — resource 9"`, log.WarnLevel,
			"guard created — resource 9"},
		{`time="2026-10-19 17:29:05" level=error msg=plain`, log.ErrorLevel, "plain"},
		{"DAEMON DENIED  op=STAT  comm=python3  path=/home/a/\x1b[31mx", log.WarnLevel,
			"DAEMON DENIED  op=STAT  comm=python3  path=/home/a/?[31mx"},
		{"Started app-listener-daemon.service.", log.InfoLevel, "Started app-listener-daemon.service."},
	}
	for _, c := range cases {
		got := journalLine(c.raw)
		if got.level != c.level || got.text != c.text {
			t.Errorf("%q: got %v %q", c.raw, got.level, got.text)
		}
	}
}

func TestDaemonPageGates(t *testing.T) {
	daemonEUID = func() int { return 0 }
	t.Cleanup(func() { daemonEUID = os.Geteuid })
	a := test.NewApp()
	t.Cleanup(a.Quit)
	a.Settings().SetTheme(&appTheme{variant: theme.VariantDark})
	w := a.NewWindow("app-listener")
	s := newShell(a, w, &Options{Version: "dev"})
	s.syncBackground = true
	w.SetContent(s.view())
	s.show(pageDaemon)
	d := s.daemon
	cases := []struct {
		st              daemonState
		note            string
		startOn, stopOn bool
	}{
		{daemonState{}, "No daemon unit", false, false},
		{daemonState{installed: true, active: "inactive", proc: true}, outsideSystemdNote, false, false},
		{daemonState{installed: true, active: "deactivating"}, "wait for it", false, false},
		{daemonState{installed: true, active: "inactive"}, "", true, false},
		{daemonState{installed: true, active: "failed"}, "", true, false},
		{daemonState{installed: true, active: unitActive, proc: true}, "", false, true},
	}
	for _, c := range cases {
		d.state, d.loaded = c.st, true
		d.render()
		if (c.note == "") != (d.note.Text == "") || !strings.Contains(d.note.Text, c.note) {
			t.Errorf("%+v: note %q, want %q", c.st, d.note.Text, c.note)
		}
		if d.startBtn.Disabled() == c.startOn || d.stopBtn.Disabled() == c.stopOn {
			t.Errorf("%+v: start disabled=%v stop disabled=%v", c.st, d.startBtn.Disabled(), d.stopBtn.Disabled())
		}
	}
	daemonEUID = func() int { return 1000 }
	d.render()
	if !strings.Contains(d.note.Text, "needs root") || !d.stopBtn.Disabled() {
		t.Errorf("non-root: note %q, stop disabled=%v", d.note.Text, d.stopBtn.Disabled())
	}
}
