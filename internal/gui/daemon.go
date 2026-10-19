//go:build gui

package gui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// daemonPage starts and stops the daemon unit with systemctl. Its runs count as shell.running: an
// install, update or edit-protected session never starts while the daemon is being started or stopped.
type daemonPage struct {
	runHost
	state daemonState
	// loaded: state has been read once; reading: a read is in flight.
	loaded, reading bool

	unitT, stateT, bootT *statusTile
	startBtn, stopBtn    *handButton
	note                 *widget.Label
	box                  *fyne.Container
}

// daemonEUID is the uid systemctl runs as (a test seam).
var daemonEUID = os.Geteuid

// unitActive is systemd's ActiveState of a running unit.
const unitActive = "active"

type daemonState struct {
	installed, enabled bool
	// proc is a daemon process running outside systemd (or an unreadable /proc).
	proc        bool
	active, sub string
}

func readDaemonState() daemonState {
	st := daemonState{installed: systemd.DaemonUnitInstalled(), enabled: systemd.IsDaemonEnabled()}
	st.active, st.sub = systemd.DaemonUnitState()
	pids, err := protected.FindDaemonProcesses("/proc")
	st.proc = err != nil || len(pids) > 0
	return st
}

func newDaemonPage(s *shell) *daemonPage {
	return &daemonPage{runHost: runHost{s: s}}
}

func (d *daemonPage) view() fyne.CanvasObject {
	d.unitT = newStatusTile("gear", "Unit")
	d.stateT = newStatusTile("pulse", "State")
	d.bootT = newStatusTile("arrow-clockwise", "At boot")
	d.startBtn = newButton("Start", icon("play"), d.start)
	d.startBtn.Importance = widget.HighImportance
	d.stopBtn = newButton("Stop", icon("stop"), d.confirmStop)
	d.stopBtn.Importance = widget.DangerImportance
	d.note = warnText()
	d.box = container.NewVBox(boldLabel(systemd.DaemonServiceName),
		mutedText("Start attaches every guard before it unlocks the encrypted folders: systemctl returns once the "+
			"daemon is ready. Stop keeps the guards attached until every folder is locked again, then checks "+
			"they are; until the next start the encrypted folders stay unreadable and the others unguarded."),
		d.note, container.NewHBox(d.startBtn, d.stopBtn))
	d.startBtn.Disable()
	d.stopBtn.Disable()

	over := newFadeVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), newGapVBox(
		newAdaptiveGrid(3, d.unitT.card, d.stateT.card, d.bootT.card),
		container.NewVBox(sectionTitle("CONTROL"), newPanel(d.box)),
		container.NewVBox(sectionTitle("NEEDS THE DAEMON"), newPanel(mutedText(
			"Running: live edits, temporary access, configuration edits and recording (Protected), and a live "+
				"catalog refresh (Install). Stopped: the catalog diff, restoring backups and uninstalling."))),
	)))
	root := d.mount(over, "Back to Daemon")
	d.refresh()
	d.s.every(2*time.Second, func() {
		if d.s.current == pageDaemon && !d.running {
			d.refresh()
		}
	})
	return root
}

// refresh re-reads the unit off the UI goroutine.
func (d *daemonPage) refresh() {
	if d.reading {
		return
	}
	d.reading = true
	var st daemonState
	d.s.background(func() { st = readDaemonState() }, func() {
		d.reading, d.state, d.loaded = false, st, true
		d.render()
	})
}

func (d *daemonPage) render() {
	st := d.state
	if st.installed {
		d.unitT.set(true, "installed", "systemctl manages the daemon.", theme.ColorNameSuccess)
	} else {
		d.unitT.set(false, "not installed", "Set it up on the Install page.", "")
	}
	switch st.active {
	case unitActive:
		d.stateT.set(true, "running", "Guarding its configuration ("+st.sub+").", theme.ColorNameSuccess)
	case "activating", "deactivating", "reloading":
		d.stateT.set(true, st.active, "systemd is "+st.active+" it ("+st.sub+").", colorInfo)
	case "failed":
		d.stateT.set(false, "failed", "It exited with an error: journalctl -u "+systemd.DaemonServiceName+" -e", "")
	default:
		d.stateT.set(true, "stopped", "Encrypted folders are locked.", theme.ColorNameWarning)
	}
	if st.enabled {
		d.bootT.set(true, "enabled", "It starts at boot.", theme.ColorNameSuccess)
	} else {
		d.bootT.set(true, "disabled", "It doesn't start at boot (the installer enables it).", colorMuted)
	}
	d.renderButtons()
}

// gateNote is why neither button may run now ("" = none).
func (d *daemonPage) gateNote() string {
	st := d.state
	switch {
	case daemonEUID() != 0:
		return "systemctl needs root: restart the window with sudo."
	case !d.loaded:
		return "Checking the daemon…"
	case !st.installed:
		return "No daemon unit is installed: set it up on the Install page."
	case st.active != unitActive && st.proc:
		return outsideSystemdNote
	case st.active == "activating" || st.active == "deactivating" || st.active == "reloading":
		return "The daemon is " + st.active + ": wait for it."
	case d.s.running() && !d.running:
		return "An installer or updater run is in progress."
	case d.s.sessionOpen():
		return "An edit-protected session is open: end it first (Protected workspace)."
	}
	if active := d.s.activeGuards(); len(active) > 0 {
		return "Turn off " + strings.Join(active, ", ") + " first (Files and Network workspaces): " +
			"the daemon's lockdown and unlock aren't admitted by a guard this window runs."
	}
	return ""
}

func (d *daemonPage) renderButtons() {
	note := d.gateNote()
	setWarn(d.note, note)
	up := d.state.active == unitActive
	enable := func(b *handButton, on bool) {
		if on && note == "" && !d.running {
			b.Enable()
		} else {
			b.Disable()
		}
	}
	enable(d.startBtn, !up)
	enable(d.stopBtn, up)
	d.box.Refresh()
}

func (d *daemonPage) confirmStop() {
	showConfirm("Stop the daemon?", "Every encrypted folder is locked again and stays unreadable until the "+
		"daemon starts; folders without encryption are unguarded meanwhile.", func(ok bool) {
		if ok {
			d.launch("Stop the daemon", "stop", stopProgress, systemd.StopDaemonIfRunning)
		}
	}, d.s.win)
}

func (d *daemonPage) start() {
	d.launch("Start the daemon", "play", startProgress, systemd.StartDaemon)
}

func (d *daemonPage) launch(name, iconName string, progress map[string]float64, op func() error) {
	if d.begin(name, iconName, func(ui wizard.UI) error {
		ui.Stage(name)
		return d.watch(ui, progress, op)
	}, nil) {
		d.renderButtons()
	}
}

// startProgress and stopProgress place each systemd SubState on the bar: the bar follows the unit,
// it never advances on time alone.
var (
	startProgress = map[string]float64{"dead": 0.05, "failed": 0.05, "start-pre": 0.25, "start": 0.5,
		"start-post": 0.8, "running": 1}
	stopProgress = map[string]float64{"running": 0.05, "stop": 0.2, "stop-sigterm": 0.35, "stop-sigkill": 0.5,
		"stop-post": 0.7, "final-sigterm": 0.85, "final-sigkill": 0.9, "dead": 1, "failed": 1}
)

// watch runs op with the bar on the unit's state and the unit's journal in the run log.
func (d *daemonPage) watch(ui wizard.UI, progress map[string]float64, op func() error) error {
	ctx, cancel := context.WithCancel(context.Background())
	// The journal lags the unit: keep tailing a moment after op returns.
	defer time.AfterFunc(time.Second, cancel)
	go d.tailJournal(ctx, time.Now())
	return ui.Progress(func(set func(string, float64)) error {
		done := make(chan error, 1)
		go func() { done <- op() }()
		tick := time.NewTicker(250 * time.Millisecond)
		defer tick.Stop()
		var at float64
		for {
			select {
			case err := <-done:
				return err
			case <-tick.C:
				active, sub := systemd.DaemonUnitState()
				at = max(at, progress[sub])
				set(fmt.Sprintf("%s (%s)", orDefault(active, "unknown"), orDefault(sub, "unknown")), at)
			}
		}
	})
}

// tailJournal mirrors the unit's journal lines since into the run log until ctx ends.
func (d *daemonPage) tailJournal(ctx context.Context, since time.Time) {
	//nolint:gosec // fixed program and arguments
	cmd := exec.CommandContext(ctx, "journalctl", "-u", systemd.DaemonServiceName, "-f", "--lines=all",
		"--since=@"+strconv.FormatInt(since.Unix(), 10), "-o", "cat", "--no-pager")
	out, err := cmd.StdoutPipe()
	if err == nil {
		err = cmd.Start()
	}
	if err != nil {
		log.Warnf("daemon journal unavailable: %v", err)
		return
	}
	sc := bufio.NewScanner(out)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := journalLine(sc.Text())
		fyne.Do(func() { d.run.appendLog(line) })
	}
	_ = cmd.Wait()
}

// journalLine reads a daemon journal line: logrus logfmt (level=, msg=) or a headless event line.
func journalLine(raw string) logLine {
	l := logLine{level: log.InfoLevel, at: time.Now(), text: raw}
	if _, rest, ok := strings.Cut(raw, "level="); ok {
		lvl, _, _ := strings.Cut(rest, " ")
		if parsed, err := log.ParseLevel(lvl); err == nil {
			l.level = parsed
		}
		if _, msg, ok := strings.Cut(rest, "msg="); ok {
			if q, err := strconv.QuotedPrefix(msg); err == nil {
				msg, _ = strconv.Unquote(q)
			}
			l.text = msg
		}
	} else if strings.Contains(raw, " DENIED ") {
		l.level = log.WarnLevel
	}
	l.text = logging.SanitizeText(l.text)
	return l
}
