//go:build gui

package gui

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/headless"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// stagedBinary is the headless build the installer deploys, copied or downloaded into the page's
// root-only staging directory so nobody else can swap it before the deploy reads it.
type stagedBinary struct {
	path, sha string
	// tag is set for a signature-verified release download, empty for a local file.
	tag string
}

// binSource is where the headless build comes from.
type binSource int

const (
	srcThisVersion binSource = iota // this window's own release tag
	srcStable
	srcPreRelease
	srcLocal
)

var srcLabels = map[binSource]string{
	srcThisVersion: "This version's release",
	srcStable:      "Latest stable release",
	srcPreRelease:  "Latest pre-release",
	srcLocal:       "Local file",
}

// installedBinary is what systemd.InstallBinaryPath holds now.
type installedBinary struct {
	present bool
	version string
	sha     string
	err     error
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
		"Pick the users and folders to protect, review daemon.conf, encrypt with fscrypt and start the daemon. " +
			"Safe to re-run: a finished or interrupted install is resumed.", true},
	{InstallUpdateCatalog, "arrow-clockwise", theme.ColorNameSuccess, "Refresh catalog whitelists",
		"install --update-catalog-only", "Re-expand the catalog's whitelists for every guarded folder, show the " +
			"diff and rewrite daemon.conf. Your own sections are kept. Live reloads the running daemon instead of " +
			"stopping it.", false},
	{InstallDiffCatalog, "plus", theme.ColorNameSuccess, "Protect new catalog folders", "install --diff-catalog",
		"Find the critical folders on this host that daemon.conf doesn't guard yet, pick some, and encrypt and " +
			"guard them like a fresh install.", false},
	{InstallRestoreBackups, "lock-simple-open", theme.ColorNameWarning, "Restore backups", "install --restore-backups",
		"Undo a migration: delete the encrypted folders and move their .app_listener.backup copies back. " +
			"Refused while the daemon runs.", false},
	{InstallDeleteBackups, "trash", theme.ColorNameError, "Delete backups", "install --delete-post-backups",
		"Permanently delete the .app_listener.backup copies the migration kept.", false},
}

// installPage is the Install workspace: the headless build the daemon will run, the installer's
// options and entry points, and the run view each of them switches to.
type installPage struct {
	s     *shell
	hooks InstallHooks
	stack *fyne.Container
	over  fyne.CanvasObject
	run   *runView

	running  bool
	staging  string
	bin      *stagedBinary
	present  installedBinary
	daemonUp bool

	installedT, daemonT, windowT *statusTile
	binTile                      *iconTile
	binState                     *pill
	binHead, binCard             *fyne.Container
	binInfo, binIntro, binNote   *widget.Label
	binBtn                       *widget.Button
	source                       binSource
	meta, live                   *widget.Check
	binButtons                   []*widget.Button
}

func newInstallPage(s *shell) *installPage {
	return &installPage{s: s, hooks: s.opts.Install}
}

func (p *installPage) view() fyne.CanvasObject {
	p.installedT = newStatusTile("terminal-window", "Installed binary")
	p.daemonT = newStatusTile("pulse", "Daemon")
	p.windowT = newStatusTile("house", "This window")
	p.windowT.set(true, p.s.opts.Version, "The desktop build. It is never deployed as the daemon: the daemon "+
		"runs the headless build below.", colorInfo)

	p.meta = widget.NewCheck("Log denied metadata inspections", nil)
	p.live = widget.NewCheck("Live: keep the daemon running", nil)

	p.over = container.NewVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), container.NewVBox(
		container.NewGridWithColumns(3, p.installedT.card, p.daemonT.card, p.windowT.card),
		sectionTitle("DAEMON BINARY"),
		p.binaryCard(),
		sectionTitle("OPTIONS"),
		newPanel(container.NewVBox(p.meta, mutedText("Installs the daemon unit without "+
			"--no-log-metadata-blocks, so denied process inspections (op=PTRACE mode=READ) are logged. Use it "+
			"to diagnose an app that breaks with nothing logged. Applies to every action that writes the unit."))),
		sectionTitle("ACTIONS"),
		p.actionGrid(),
	)))
	p.run = newRunView(p.s, func() { p.show(p.over) })
	p.stack = container.NewStack(p.over)
	p.renderBinary()
	p.refreshStatus()
	return p.stack
}

func (p *installPage) show(o fyne.CanvasObject) {
	p.stack.Objects = []fyne.CanvasObject{o}
	p.stack.Refresh()
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
	p.binBtn = widget.NewButton("", p.acquire)
	p.binBtn.Importance = widget.HighImportance

	sources := []binSource{srcStable, srcPreRelease, srcLocal}
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
	pick := widget.NewSelect(labels, func(l string) {
		p.source = sources[slices.Index(labels, l)]
		p.renderSource()
	})
	pick.Selected = labels[0] // SetSelected fires the callback before the card exists
	p.source = sources[0]

	p.binHead = container.NewBorder(nil, nil, p.binTile, nil,
		container.NewVBox(boldLabel("Headless build to deploy"), container.NewHBox(p.binState)))
	source := container.NewBorder(nil, nil, boldLabel("Source"), p.binBtn, pick)
	p.binCard = container.NewVBox(p.binHead, source, p.binIntro, p.binInfo, p.binNote)
	p.renderSource()
	return newPanel(p.binCard)
}

func (p *installPage) renderSource() {
	const verified = " Its Ed25519 signature, checksum and GitHub digest are verified before anything is " +
		"installed. A network whitelist active in this window would block the download."
	switch p.source {
	case srcThisVersion:
		p.binIntro.SetText("The headless asset of this window's own release, " + p.s.opts.Version + "." + verified)
	case srcStable:
		p.binIntro.SetText("The newest stable release (vX.Y.Z), what update installs by default." + verified)
	case srcPreRelease:
		p.binIntro.SetText("The newest pre-release (the pre-YYYYMMDD builds of main), what " +
			"update --channel pre-release installs." + verified)
	case srcLocal:
		p.binIntro.SetText("A headless binary you built (make build, without GUI=1), copied into a " +
			"root-only folder. It is not signature-verified.")
	}
	if p.source == srcLocal {
		p.binBtn.SetText("Choose binary…")
		p.binBtn.SetIcon(icon("file"))
	} else {
		p.binBtn.SetText("Download and verify")
		p.binBtn.SetIcon(icon("download-simple"))
	}
	p.binCard.Refresh()
}

func (p *installPage) acquire() {
	if p.source == srcLocal {
		p.chooseLocal()
		return
	}
	p.download(p.source)
}

func (p *installPage) actionGrid() fyne.CanvasObject {
	grid := container.NewGridWithColumns(3)
	for i := range installActions {
		a := &installActions[i]
		btn := widget.NewButtonWithIcon("Start", icon("play"), func() { p.startAction(a) })
		btn.Importance = widget.MediumImportance
		if a.mode == InstallWizard {
			btn.Importance = widget.HighImportance
		}
		if a.needsBin {
			p.binButtons = append(p.binButtons, btn)
		}
		flag := widget.NewLabel(a.flag)
		flag.TextStyle = fyne.TextStyle{Monospace: true}
		flag.Importance = widget.LowImportance
		flag.SizeName = theme.SizeNameCaptionText
		top := container.NewVBox(newIconTile(a.icon, a.color, 40), boldLabel(a.title), flag, mutedText(a.desc))
		bottom := container.NewHBox(btn)
		if a.mode == InstallUpdateCatalog {
			bottom = container.NewVBox(p.live, container.NewHBox(btn))
		}
		if a.needsBin {
			bottom = container.NewVBox(mutedText("Needs the headless build above."), bottom)
		}
		grid.Add(newPanel(container.NewBorder(top, bottom, nil, nil)))
	}
	return grid
}

// renderBinary shows the staged build and enables the actions that deploy it.
func (p *installPage) renderBinary() {
	switch {
	case p.bin == nil:
		p.binTile.Set("terminal-window", colorMuted)
		p.binState.Set("none yet", colorMuted)
		p.binInfo.SetText("")
		p.binInfo.Hide()
	case p.bin.tag != "":
		p.binTile.Set("shield-check", theme.ColorNameSuccess)
		p.binState.Set("release "+p.bin.tag+" · signature verified", theme.ColorNameSuccess)
		p.binInfo.SetText("sha256 " + p.bin.sha + "\n" + p.bin.path)
		p.binInfo.Show()
	default:
		p.binTile.Set("warning", theme.ColorNameWarning)
		p.binState.Set("local build · not signature-verified", theme.ColorNameWarning)
		p.binInfo.SetText("sha256 " + p.bin.sha + "\n" + p.bin.path)
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
	for _, b := range p.binButtons {
		if p.bin == nil {
			b.Disable()
		} else {
			b.Enable()
		}
	}
	p.renderInstalled()
}

// versionNote warns when the staged build is another version than this window, whose installer code
// writes daemon.conf.
func (p *installPage) versionNote() string {
	if p.s.opts.Release && p.bin.tag == p.s.opts.Version {
		return ""
	}
	name := p.s.opts.Version
	if !p.s.opts.Release {
		name = "a development build"
	}
	return "Not this window's version: this window is " + name + ", and its installer writes daemon.conf. " +
		"Before the service is enabled the deployed build checks that config (daemon --check); the install " +
		"stops if it can't load it."
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
			"update it from a terminal (app-listener update, or install --binary-only) first.", theme.ColorNameWarning)
	}
	if p.daemonUp {
		p.daemonT.set(true, "running", "The installer stops it where a step needs that and starts it again.",
			theme.ColorNameSuccess)
	} else {
		p.daemonT.set(true, "stopped", "No daemon is running.", theme.ColorNameWarning)
	}
}

// refreshStatus re-reads the installed binary and the daemon off the UI goroutine.
func (p *installPage) refreshStatus() {
	versionOf := p.hooks.InstalledVersion
	var ib installedBinary
	var up bool
	p.s.background(func() {
		ib = readInstalled(versionOf)
		up = protected.DaemonRunning()
	}, func() {
		p.present, p.daemonUp = ib, up
		p.renderInstalled()
	})
}

func readInstalled(versionOf func() string) installedBinary {
	if _, err := os.Stat(systemd.InstallBinaryPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return installedBinary{}
		}
		return installedBinary{err: err}
	}
	ib := installedBinary{present: true, version: "unknown"}
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

func hashReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// blocked names why no installer run may start now: a guard this window runs would deny the
// installer (it moves and encrypts guarded folders, restarts the daemon, downloads from GitHub)
// since none of them admits this window.
func (p *installPage) blocked() string {
	var active []string
	for _, t := range p.s.files.Targets() {
		if t.Policy().Mode != workbench.GuardOff {
			active = append(active, "the guard on "+t.Spec.Path)
		}
	}
	if p.s.network.Policy().Mode != workbench.GuardOff {
		active = append(active, "the network guard")
	}
	if len(active) == 0 {
		return ""
	}
	return "Turn off " + strings.Join(active, ", ") + " first (Files and Network workspaces).\n\n" +
		"The installer moves and encrypts protected folders, restarts the daemon and downloads from " +
		"GitHub. A guard running in this window doesn't admit the installer, so it would deny those steps."
}

// begin switches to the run view and runs work off the UI goroutine with its prompts, stages and
// logs on that view; after runs on success.
func (p *installPage) begin(name, iconName string, work func(ui wizard.UI) error, after func()) bool {
	if p.running {
		return false
	}
	if why := p.blocked(); why != "" {
		dialog.ShowInformation("Guards are running", why, p.s.win)
		return false
	}
	p.running = true
	p.run.start(name, iconName)
	p.show(p.run.root)
	ui := &installUI{s: p.s, run: p.run}
	restoreLogs := captureLogs(p.run)
	go func() {
		err := work(ui)
		restoreLogs()
		fyne.Do(func() {
			p.running = false
			canceled := errors.Is(err, wizard.ErrCanceled) || errors.Is(err, inst.ErrEditCanceled)
			p.run.finish(err, canceled)
			if err == nil && after != nil {
				after()
			}
			p.refreshStatus()
			p.s.reprobe()
		})
	}()
	return true
}

func (p *installPage) startAction(a *installAction) {
	if a.needsBin && p.bin == nil {
		return
	}
	req := InstallRequest{Mode: a.mode, AllowMetadataOutput: p.meta.Checked,
		Live: a.mode == InstallUpdateCatalog && p.live.Checked}
	if p.bin != nil {
		req.Binary = p.bin.path
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
		dialog.ShowError(err, p.s.win)
		return
	}
	var staged *stagedBinary
	started := p.begin("Download: "+srcLabels[src], "download-simple", func(ui wizard.UI) error {
		if fetch == nil || fetchLatest == nil {
			return errors.New("release download is not available in this build")
		}
		restore := wizard.Use(ui)
		defer restore()
		var path, tag string
		var err error
		if src == srcThisVersion {
			ui.Stage("Downloading and verifying " + thisTag)
			tag = thisTag
			path, err = fetch(tag, dir)
		} else {
			ui.Stage("Downloading and verifying the " + strings.ToLower(srcLabels[src]))
			path, tag, err = fetchLatest(src == srcPreRelease, dir)
		}
		if err != nil {
			return err
		}
		ui.Stage("Checking it is the headless build")
		if checkErr := headless.Check(path); checkErr != nil {
			return checkErr
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		sum, err := hashReader(f)
		if err != nil {
			return err
		}
		staged = &stagedBinary{path: path, sha: sum, tag: tag}
		return nil
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

func (p *installPage) chooseLocal() {
	d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		dir, err := p.ensureStaging()
		if err != nil {
			_ = r.Close()
			dialog.ShowError(err, p.s.win)
			return
		}
		p.bin = nil
		p.renderBinary()
		var staged *stagedBinary
		p.s.busy("Checking the binary…", func() error {
			defer r.Close()
			var stageErr error
			staged, stageErr = stageLocal(r, dir, headless.Check)
			return stageErr
		}, func() {
			p.bin = staged
			p.renderBinary()
		})
	}, p.s.win)
	start, _ := os.Getwd()
	p.s.showFile(d, start)
}

// stageLocal copies src into dir (a fresh 0700 file) and refuses it unless it is the headless build.
func stageLocal(src io.Reader, dir string, check func(string) error) (*stagedBinary, error) {
	dst := filepath.Join(dir, "app-listener-local")
	if err := os.Remove(dst); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o700)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	_, err = io.Copy(io.MultiWriter(f, h), src)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = check(dst)
	}
	if err != nil {
		_ = os.Remove(dst)
		return nil, err
	}
	return &stagedBinary{path: dst, sha: hex.EncodeToString(h.Sum(nil))}, nil
}
