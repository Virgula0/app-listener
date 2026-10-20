//go:build gui

package gui

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/headless"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// stagedBinary is the signature-verified release the installer deploys (its headless build as the
// daemon, its desktop build at systemd.GUIBinaryPath), downloaded into the page's root-only staging
// directory so nobody else can swap it before the deploy reads it. guiPath "" = the release has no
// desktop build.
type stagedBinary struct {
	path, sha, tag  string
	guiPath, guiSha string
}

// binSource is the release channel the headless build is downloaded from. Only signed releases:
// an unverified local file is never offered for deploy.
type binSource int

const (
	srcThisVersion binSource = iota // this window's own release tag
	srcStable
	srcPreRelease
)

var srcLabels = map[binSource]string{
	srcThisVersion: "This version's release",
	srcStable:      "Latest stable release",
	srcPreRelease:  "Latest pre-release",
}

// installedBinary is what systemd.InstallBinaryPath holds now; guiVersion is the desktop build's
// at systemd.GUIBinaryPath ("" = not installed).
type installedBinary struct {
	present    bool
	version    string
	sha        string
	err        error
	guiVersion string
}

type installAction struct {
	mode     InstallMode
	icon     string
	color    fyne.ThemeColorName
	title    string
	flag     string
	desc     string
	needsBin bool
}

var installActions = []installAction{
	{InstallWizard, "shield-check", theme.ColorNamePrimary, "Install wizard", "install",
		"Pick the users and folders to protect, review daemon.conf, encrypt with fscrypt and start the daemon, " +
			"then install the release's desktop build as app-listener-gui. An interrupted install is resumed; " +
			"after a completed one, uninstall first to run it again.", true},
	{InstallRefreshCatalog, "arrow-clockwise", theme.ColorNameSuccess, "Refresh catalog whitelists",
		"install catalog refresh", "Re-expand the catalog's whitelists for every guarded folder, show the " +
			"diff and rewrite daemon.conf. Your own sections are kept. Live reloads the running daemon instead of " +
			"stopping it.", false},
	{InstallDiffCatalog, "plus", theme.ColorNameSuccess, "Protect new catalog folders", "install catalog diff",
		"Find the critical folders on this host that daemon.conf doesn't guard yet, pick some, and encrypt and " +
			"guard them like a fresh install. Needs the daemon stopped; asks at the end whether to start it.", false},
	{InstallRestoreBackups, "lock-simple-open", theme.ColorNameWarning, "Restore backups", "install backups restore",
		"Undo a migration: delete the encrypted folders and move their .app_listener.backup copies back. " +
			"Refused while the daemon runs; asks before touching anything.", false},
	{InstallDeleteBackups, "trash", theme.ColorNameError, "Delete backups", "install backups delete",
		"Permanently delete the .app_listener.backup copies the migration kept. Works with the daemon running; " +
			"asks before deleting.", false},
	{InstallFacultatives, "plugs-connected", theme.ColorNamePrimary, "Install facultatives",
		"install facultatives", "Pick the optional per-user setup to install: the ssh-agent unit, the Bun " +
			"$TMPDIR redirect. Offered for every user, whether or not daemon.conf guards what it serves. " +
			"The daemon keeps running.", false},
}

// installPage is the Install workspace: the headless build the daemon will run, the installer's
// options and entry points, and the run view each of them switches to.
type installPage struct {
	runHost
	hooks InstallHooks

	staging  string
	bin      *stagedBinary
	present  installedBinary
	daemonUp bool
	// daemonProc is a daemon process running outside systemd (uninstall refuses it too).
	daemonProc bool
	// loaded: the daemon state above has been read once.
	loaded bool
	// cert is the completed-install record (nil = none); certErr is why it couldn't be read.
	cert    *protected.InstallCert
	certErr error

	installedT, daemonT, windowT *statusTile
	binTile                      *iconTile
	binState                     *pill
	binHead, binCard             *fyne.Container
	binInfo, binIntro, binNote   *widget.Label
	binBtn                       *handButton
	source                       binSource
	meta, live, withGUI          *handCheck
	gates                        []actionGate
	uninst                       uninstallCard
}

// actionGate is an action card's Start button and the note that says why it is refused.
type actionGate struct {
	a    *installAction
	btn  *handButton
	note *widget.Label
	box  *fyne.Container
}

// set shows note and disables the button while it is non-empty.
func (g *actionGate) set(note string) {
	setWarn(g.note, note)
	if note == "" {
		g.btn.Enable()
	} else {
		g.btn.Disable()
	}
	g.box.Refresh()
}

func newInstallPage(s *shell) *installPage {
	return &installPage{runHost: runHost{s: s}, hooks: s.opts.Install}
}

func (p *installPage) view() fyne.CanvasObject {
	p.installedT = newStatusTile("terminal-window", "Installed binary")
	p.daemonT = newStatusTile("pulse", "Daemon")
	p.windowT = newStatusTile("sidebar-simple", "Desktop build")

	p.meta = newCheck("Log denied metadata inspections", nil)
	p.live = newCheck("Live: keep the daemon running", func(bool) { p.renderActions() })
	p.live.Checked = true
	p.withGUI = newCheck("Also install the desktop build (app-listener-gui)", nil)
	p.withGUI.Checked = true

	over := newFadeVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), container.NewVBox(
		container.NewGridWithColumns(3, p.installedT.card, p.daemonT.card, p.windowT.card),
		sectionTitle("RELEASE"),
		p.binaryCard(),
		sectionTitle("OPTIONS"),
		newPanel(container.NewVBox(
			p.withGUI, mutedText("Once the wizard completes, the staged release's desktop build is installed at "+
				systemd.GUIBinaryPath+" (symlinked in /usr/local/bin): run it with sudo app-listener-gui. "+
				"It is verified like the headless build, and Update keeps it on the daemon's release."),
			p.meta, mutedText("Installs the daemon unit without "+
				"--no-log-metadata-blocks, so denied process inspections (op=PTRACE mode=READ) are logged. Use it "+
				"to diagnose an app that breaks with nothing logged. Applies to every action that writes the unit."))),
		sectionTitle("ACTIONS"),
		p.actionGrid(),
		sectionTitle("UNINSTALL"),
		container.NewGridWithColumns(2, p.uninstallCard(), p.facultativesCard()),
	)))
	root := p.mount(over, "Back to Install")
	p.renderBinary()
	p.refreshStatus()
	return root
}

func (p *installPage) binaryCard() fyne.CanvasObject {
	p.binTile = newIconTile("terminal-window", colorMuted, 44)
	p.binState = newPill("", colorMuted)
	p.binInfo = selectable("")
	p.binInfo.TextStyle = fyne.TextStyle{Monospace: true}
	p.binIntro = mutedText("")
	p.binNote = widget.NewLabel("")
	p.binNote.Wrapping = fyne.TextWrapWord
	p.binNote.Importance = widget.WarningImportance
	p.binBtn = newButton("Download and verify", icon("download-simple"), func() { p.download(p.source) })
	p.binBtn.Importance = widget.HighImportance

	sources := []binSource{srcStable, srcPreRelease}
	if p.s.opts.Release {
		sources = append([]binSource{srcThisVersion}, sources...)
	}
	labels := make([]string, len(sources))
	for i, src := range sources {
		labels[i] = srcLabels[src]
		if src == srcThisVersion {
			labels[i] += " (" + p.s.opts.Version + ")"
		}
	}
	pick := newSelect(labels, func(l string) {
		p.source = sources[slices.Index(labels, l)]
		p.renderSource()
	})
	pick.Selected = labels[0] // SetSelected fires the callback before the card exists
	p.source = sources[0]

	p.binHead = container.NewBorder(nil, nil, p.binTile, nil,
		container.NewVBox(boldLabel("Release to deploy"), container.NewHBox(p.binState)))
	source := container.NewBorder(nil, nil, boldLabel("Source"), p.binBtn, pick)
	p.binCard = container.NewVBox(p.binHead, source, p.binIntro, p.binInfo, p.binNote)
	p.renderSource()
	return newPanel(p.binCard)
}

func (p *installPage) renderSource() {
	const verified = " Its headless build becomes the daemon and its desktop build is installed as " +
		"app-listener-gui. Each one's Ed25519 signature, checksum and GitHub digest are verified before " +
		"anything is installed. A network whitelist active in this window would block the download."
	switch p.source {
	case srcThisVersion:
		p.binIntro.SetText("This window's own release, " + p.s.opts.Version + "." + verified)
	case srcStable:
		p.binIntro.SetText("The newest stable release (vX.Y.Z), what update installs by default." + verified)
	case srcPreRelease:
		p.binIntro.SetText("The newest pre-release (the pre-YYYYMMDD builds of main), what " +
			"update --channel pre-release installs." + verified)
	}
	p.binCard.Refresh()
}

func (p *installPage) actionGrid() fyne.CanvasObject {
	grid := container.NewGridWithColumns(3)
	for i := range installActions {
		a := &installActions[i]
		btn := newButton("Start", icon("play"), func() { p.startAction(a) })
		btn.Importance = widget.MediumImportance
		if a.mode == InstallWizard {
			btn.Importance = widget.HighImportance
		}
		flag := widget.NewLabel(a.flag)
		flag.TextStyle = fyne.TextStyle{Monospace: true}
		flag.Importance = widget.LowImportance
		flag.SizeName = theme.SizeNameCaptionText
		top := container.NewVBox(newIconTile(a.icon, a.color, 40), boldLabel(a.title), flag, mutedText(a.desc))
		bottom := container.NewVBox()
		if a.needsBin {
			bottom.Add(mutedText("Needs the headless build above."))
		}
		if a.mode == InstallRefreshCatalog {
			bottom.Add(p.live)
		}
		note := warnText()
		bottom.Add(note)
		bottom.Add(container.NewHBox(btn))
		p.gates = append(p.gates, actionGate{a: a, btn: btn, note: note, box: bottom})
		grid.Add(newPanel(container.NewBorder(top, bottom, nil, nil)))
	}
	return grid
}

// renderBinary shows the staged build and enables the actions that deploy it.
func (p *installPage) renderBinary() {
	if p.bin == nil {
		p.binTile.Set("terminal-window", colorMuted)
		p.binState.Set("none yet", colorMuted)
		p.binInfo.SetText("")
		p.binInfo.Hide()
	} else {
		p.binTile.Set("shield-check", theme.ColorNameSuccess)
		p.binState.Set("release "+p.bin.tag+" · signatures verified", theme.ColorNameSuccess)
		info := "daemon   sha256 " + p.bin.sha + "\n         " + p.bin.path
		if p.bin.guiPath != "" {
			info += "\ndesktop  sha256 " + p.bin.guiSha + "\n         " + p.bin.guiPath
		}
		p.binInfo.SetText(info)
		p.binInfo.Show()
	}
	note := ""
	if p.bin != nil {
		note = p.versionNote()
	}
	p.binNote.SetText(note)
	p.binNote.Hidden = note == ""
	p.binHead.Refresh()
	p.binCard.Refresh() // binInfo's and binNote's visibility changed
	if p.bin != nil && p.bin.guiPath == "" {
		p.withGUI.Disable()
	} else {
		p.withGUI.Enable()
	}
	p.renderInstalled()
}

// versionNote warns when the staged build is another version than this window, whose installer code
// writes daemon.conf, or ships no desktop build.
func (p *installPage) versionNote() string {
	var notes []string
	if !p.s.opts.Release || p.bin.tag != p.s.opts.Version {
		name := p.s.opts.Version
		if !p.s.opts.Release {
			name = "a development build"
		}
		notes = append(notes, "Not this window's version: this window is "+name+", and its installer writes "+
			"daemon.conf. Before the service is enabled the deployed build checks that config (daemon check); "+
			"the install stops if it can't load it.")
	}
	if p.bin.guiPath == "" {
		notes = append(notes, "This release ships no desktop build for this architecture: the wizard installs "+
			"only the daemon.")
	}
	return strings.Join(notes, "\n")
}

func (p *installPage) renderInstalled() {
	ib := p.present
	switch {
	case ib.err != nil:
		p.installedT.set(false, "unreadable", ib.err.Error(), "")
	case !ib.present:
		p.installedT.set(true, "not installed", "Nothing at "+systemd.InstallBinaryPath+
			" yet: the wizard deploys the headless build.", colorMuted)
	case p.bin == nil:
		p.installedT.set(true, ib.version, systemd.InstallBinaryPath, theme.ColorNameSuccess)
	case ib.sha == p.bin.sha:
		p.installedT.set(true, ib.version, "Identical to the staged headless build.", theme.ColorNameSuccess)
	default:
		p.installedT.set(true, ib.version, "Differs from the staged build, and the wizard needs them identical: "+
			"update it first (the Update workspace, or install binary from a terminal).", theme.ColorNameWarning)
	}
	if p.daemonUp {
		p.daemonT.set(true, "running", "The installer stops it where a step needs that and starts it again.",
			theme.ColorNameSuccess)
	} else {
		p.daemonT.set(true, "stopped", "No daemon is running.", theme.ColorNameWarning)
	}
	p.renderDesktop()
	p.renderActions()
	p.renderUninstall()
}

// renderActions refuses up front what the installer would: the wizard after a completed install
// (protected.RequireNoInstallCert); the wizard and a stopped-flow refresh stop a systemd daemon
// themselves but not one outside systemd; a live refresh needs the unit active; the catalog diff and
// the restore need no daemon at all.
func (p *installPage) renderActions() {
	for _, g := range p.gates {
		note := ""
		if p.loaded {
			note = p.actionNote(g.a)
		}
		setWarn(g.note, note)
		if note == "" && (!g.a.needsBin || p.bin != nil) {
			g.btn.Enable()
		} else {
			g.btn.Disable()
		}
		g.box.Refresh()
	}
}

func (p *installPage) actionNote(a *installAction) string {
	switch a.mode {
	case InstallRestoreBackups:
		return daemonStopNote(p.daemonUp, p.daemonProc)
	case InstallRefreshCatalog:
		if p.live.Checked {
			if !p.daemonUp {
				return "Live needs the daemon running under systemd (locked folders hide their file names): " +
					"start it on the Daemon page, or untick Live to refresh with it stopped."
			}
			return ""
		}
		fallthrough
	case InstallWizard:
		switch {
		case p.certErr != nil:
			return p.certErr.Error() + ": the wizard refuses to run. Uninstall removes it."
		case p.cert != nil:
			return "Already installed (" + p.cert.Version + ", completed " +
				p.cert.CompletedAt.Local().Format("2006-01-02 15:04") + "): uninstall first to run the wizard " +
				"again. The other actions keep working on this installation."
		case !p.daemonUp && p.daemonProc:
			return outsideSystemdNote
		}
	case InstallDiffCatalog:
		return daemonStopNote(p.daemonUp, p.daemonProc)
	}
	return ""
}

// renderDesktop shows the desktop build at systemd.GUIBinaryPath. This window is never deployed as
// the daemon, nor copied there: the wizard installs the staged release's own desktop build.
func (p *installPage) renderDesktop() {
	ib := p.present
	switch {
	case ib.guiVersion == "":
		p.windowT.set(true, "not installed", "The wizard installs the staged release's desktop build at "+
			systemd.GUIBinaryPath+". This window is "+p.s.opts.Version+".", colorMuted)
	case ib.present && ib.guiVersion != ib.version:
		p.windowT.set(true, ib.guiVersion, "Not the daemon's version ("+ib.version+"): the wizard and Update "+
			"install the matching one. This window is "+p.s.opts.Version+".", theme.ColorNameWarning)
	default:
		p.windowT.set(true, ib.guiVersion, systemd.GUIBinaryPath+". This window is "+p.s.opts.Version+".",
			theme.ColorNameSuccess)
	}
}

// refreshStatus re-reads the installed binary and the daemon off the UI goroutine.
func (p *installPage) refreshStatus() {
	versionOf, guiVersionOf := p.hooks.InstalledVersion, p.hooks.InstalledGUIVersion
	var ib installedBinary
	var up, proc bool
	var cert *protected.InstallCert
	var certErr error
	p.s.background(func() {
		ib = readInstalled(versionOf, guiVersionOf)
		up = protected.DaemonRunning()
		pids, err := protected.FindDaemonProcesses("/proc")
		proc = err != nil || len(pids) > 0
		cert, certErr = protected.ReadInstallCert()
	}, func() {
		p.present, p.daemonUp, p.daemonProc, p.loaded = ib, up, proc, true
		p.cert, p.certErr = cert, certErr
		p.renderInstalled()
	})
}

// readInstalled execs and hashes the installed binary from this (non-daemon) process, and asks the
// desktop build's version. TODO(self-binary-guard): see update.RunWith.
func readInstalled(versionOf, guiVersionOf func() string) installedBinary {
	var guiVersion string
	if guiVersionOf != nil {
		guiVersion = guiVersionOf()
	}
	if _, err := os.Stat(systemd.InstallBinaryPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return installedBinary{guiVersion: guiVersion}
		}
		return installedBinary{err: err, guiVersion: guiVersion}
	}
	ib := installedBinary{present: true, version: "unknown", guiVersion: guiVersion}
	if versionOf != nil {
		ib.version = versionOf()
	}
	f, err := os.Open(systemd.InstallBinaryPath)
	if err != nil {
		ib.err = err
		return ib
	}
	defer f.Close()
	ib.sha, ib.err = hashReader(f)
	return ib
}

// stageRelease checks that rel's builds are what they claim (headless, desktop) and hashes them.
func stageRelease(rel StagedRelease) (*stagedBinary, error) {
	if err := headless.Check(rel.Binary); err != nil {
		return nil, err
	}
	sb := &stagedBinary{path: rel.Binary, tag: rel.Tag, guiPath: rel.GUI}
	var err error
	if sb.sha, err = hashFile(rel.Binary); err != nil {
		return nil, err
	}
	if rel.GUI == "" {
		return sb, nil
	}
	if checkErr := headless.CheckGUI(rel.GUI); checkErr != nil {
		return nil, checkErr
	}
	if sb.guiSha, err = hashFile(rel.GUI); err != nil {
		return nil, err
	}
	return sb, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return hashReader(f)
}

func hashReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (p *installPage) startAction(a *installAction) {
	if a.needsBin && p.bin == nil {
		return
	}
	req := InstallRequest{Mode: a.mode, AllowMetadataOutput: p.meta.Checked,
		Live: a.mode == InstallRefreshCatalog && p.live.Checked}
	if p.bin != nil {
		req.Binary = p.bin.path
		if p.withGUI.Checked {
			req.GUIBinary = p.bin.guiPath
		}
	}
	run := p.hooks.Run
	p.begin(a.title, a.icon, func(ui wizard.UI) error {
		if run == nil {
			return errors.New("the installer is not available in this build")
		}
		return run(ui, req)
	}, nil)
}

// ensureStaging creates the root-only (0700) directory staged builds live in.
func (p *installPage) ensureStaging() (string, error) {
	if p.staging != "" {
		return p.staging, nil
	}
	dir, err := os.MkdirTemp("", "app-listener-gui-")
	if err != nil {
		return "", fmt.Errorf("creating a staging directory: %w", err)
	}
	p.staging = dir
	return dir, nil
}

func (p *installPage) cleanup() {
	if p.staging != "" {
		_ = os.RemoveAll(p.staging)
	}
}

func (p *installPage) download(src binSource) {
	fetch, fetchLatest := p.hooks.Fetch, p.hooks.FetchLatest
	thisTag := p.s.opts.Version
	dir, err := p.ensureStaging()
	if err != nil {
		showErrorDialog(err, p.s.win)
		return
	}
	var staged *stagedBinary
	started := p.begin("Download: "+srcLabels[src], "download-simple", func(ui wizard.UI) error {
		if fetch == nil || fetchLatest == nil {
			return errors.New("release download is not available in this build")
		}
		restore := wizard.Use(ui)
		defer restore()
		var rel StagedRelease
		var err error
		if src == srcThisVersion {
			ui.Stage("Downloading and verifying " + thisTag)
			rel, err = fetch(thisTag, dir)
		} else {
			ui.Stage("Downloading and verifying the " + strings.ToLower(srcLabels[src]))
			rel, err = fetchLatest(src == srcPreRelease, dir)
		}
		if err != nil {
			return err
		}
		ui.Stage("Checking the headless and desktop builds")
		staged, err = stageRelease(rel)
		return err
	}, func() {
		p.bin = staged
		p.renderBinary()
	})
	if started {
		// The download overwrites the staged file before verifying it: nothing may deploy it meanwhile.
		p.bin = nil
		p.renderBinary()
	}
}
