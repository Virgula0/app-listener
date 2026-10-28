//go:build gui

// Package gui is the desktop window `app-listener-gui` opens, behind the `gui` build tag: fyne
// adds X11/OpenGL/font/image dependencies (~15 MiB) and package-init goroutines that don't belong
// in the daemon sharing this binary. Build with `-tags gui` (or `make build GUI=1`).
package gui

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/wm"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// Available reports whether this binary was built with the fyne desktop GUI
// (this file is only compiled under the `gui` build tag).
const Available = true

const appID = "com.virgula0.app-listener"

type pageID int

const (
	pageHome pageID = iota
	pageFiles
	pageNetwork
	pageProtected
	pageUpdate
	pageInstall
	pageDaemon
	pageTrust
)

type pageDef struct {
	id       pageID
	icon     string
	title    string
	subtitle string
}

var pageDefs = []pageDef{
	{pageHome, "house", "Home", "System checks and workspaces"},
	{pageFiles, "folders", "Files", "Monitor folders and files, then guard them with a whitelist or blacklist"},
	{pageNetwork, "globe-hemisphere-west", "Network", "Watch binaries' network activity and enforce a network guard"},
	{pageProtected, "lock-simple", "Protected", "Edit protected directories, grant temporary access, edit the daemon's configuration"},
	{pageTrust, "shield-check", "Trust", "Confirm whitelisted binaries the daemon refused because their content changed or is new"},
	{pageUpdate, "arrow-clockwise", "Update", "Update the installed daemon to the latest signed release"},
	{pageInstall, "download-simple", "Install", "Set up the protection daemon"},
	{pageDaemon, "pulse", "Daemon", "Start or stop the protection daemon with systemctl"},
}

// shell is the window frame: sidebar navigation, page header and the active page.
type shell struct {
	app    fyne.App
	win    fyne.Window
	opts   Options
	status sysStatus

	nav      map[pageID]*navItem
	views    map[pageID]fyne.CanvasObject
	title    *widget.Label
	subtitle *widget.Label
	body     *fyne.Container
	home     *homePage
	current  pageID
	filesPg  *filesPage
	netPg    *networkPage
	files    *workbench.Files
	network  *workbench.Network
	install  *installPage
	update   *updatePage
	protect  *protectedPage
	trust    *trustPage
	daemon   *daemonPage
	// startBackend runs the installed binary's edit-protected backend (a fake in tests).
	startBackend func() (backendLink, error)

	wm        wm.Session
	collapsed bool
	sidebar   *fyne.Container
	navBox    *fyne.Container
	frame     *fyne.Container
	stats     procStats
	statLines [3]*widget.Label
	// picker is the open file or folder picker, which Ctrl+F searches.
	picker *picker
	// dialogs are the open sized dialogs and the size each asked for (refitDialogs).
	dialogs map[sizedDialog]fyne.Size
	// syncBackground makes background synchronous and every a no-op: fyne's test driver runs fyne.Do
	// inline on the calling goroutine, so a real goroutine would race the test's own rendering.
	syncBackground bool
}

// Run opens the window and blocks until it is closed.
func Run(opts *Options) (err error) {
	wmEnv := restoreSessionEnv()
	defer cleanupSessionEnv()
	defer filterFyneLogs()()
	// fyne logs a failed GLFW init and then panics on the first window call; a panic once the window
	// is up is a bug, reported with its stack.
	var started atomic.Bool
	defer func() {
		r := recover()
		switch {
		case r == nil:
		case started.Load():
			err = fmt.Errorf("gui: internal error: %v\n%s", r, debug.Stack())
		default:
			err = fmt.Errorf("gui: no display reachable (%v): run it from a graphical session "+
				"(under sudo it needs the invoking user's WAYLAND_DISPLAY or DISPLAY)", r)
		}
	}()

	// Opts out of fyne's off-main-goroutine marshaling: every widget call from a goroutine here
	// goes through fyne.Do. Set before the app exists: fyne reads it once, on first use.
	app.SetMetadata(fyne.AppMetadata{
		ID: appID, Name: "app-listener", Version: opts.Version,
		Migrations: map[string]bool{"fyneDo": true},
	})
	a := app.NewWithID(appID)
	a.Settings().SetTheme(&appTheme{variant: savedVariant(a)})
	a.SetIcon(logoResource)

	w := a.NewWindow("app-listener")
	s := newShell(a, w, opts)
	s.wm = detectWM(wmEnv)
	w.SetContent(s.view())
	w.Resize(preferredSize)
	w.CenterOnScreen()
	w.SetMaster()
	w.SetCloseIntercept(s.confirmClose)
	a.Lifecycle().SetOnStarted(func() {
		started.Store(true)
		if s.floatPref() {
			s.setFloating(true)
		}
	})

	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault},
		func(fyne.Shortcut) { s.focusSearch() })

	s.reprobe()
	s.trust.probe()
	s.checkForUpdate()
	s.every(time.Second, s.updateStats)
	s.every(6*time.Hour, s.checkForUpdate)
	w.ShowAndRun()
	s.closeBackends()
	s.network.Close()
	s.files.Close()
	if s.install != nil {
		s.install.cleanup()
	}
	return nil
}

// closeBackends ends the open edit-protected backends: a live grant is narrowed, an offline vault
// re-locked.
func (s *shell) closeBackends() {
	if s.protect != nil && s.protect.link != nil {
		s.protect.link.Close()
	}
	if s.trust.link != nil {
		s.trust.link.Close()
	}
}

// confirmClose asks before quitting mid-run: the installer or updater would stop between two steps,
// and an open edit-protected session ends (unsaved edits are lost).
func (s *shell) confirmClose() {
	if s.protect != nil && s.protect.link != nil {
		showConfirm("End the edit-protected session?", "Closing the window ends the open session: a "+
			"live grant is revoked, an offline directory is locked again, and unsaved edits are lost.",
			func(ok bool) {
				if ok {
					s.win.Close()
				}
			}, s.win)
		return
	}
	if !s.running() {
		s.win.Close()
		return
	}
	showConfirm("Quit during a run?", "Closing the window stops the installer or updater midway. "+
		"Steps already completed stay applied; re-run it afterwards to resume.", func(ok bool) {
		if ok {
			s.win.Close()
		}
	}, s.win)
}

// reprobe re-runs the system checks off the UI goroutine.
func (s *shell) reprobe() {
	var st sysStatus
	s.background(func() { st = probeSystem() }, func() { s.setStatus(st) })
}

// every runs fn on the UI goroutine at each interval for the window's lifetime (never under
// syncBackground: tests refresh by hand).
func (s *shell) every(interval time.Duration, fn func()) {
	if s.syncBackground {
		return
	}
	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()
		for range tick.C {
			fyne.Do(fn)
		}
	}()
}

// background runs work off the UI goroutine, then done on it.
func (s *shell) background(work, done func()) {
	if s.syncBackground {
		work()
		done()
		return
	}
	go func() {
		work()
		fyne.Do(done)
	}()
}

func savedVariant(a fyne.App) fyne.ThemeVariant {
	if a.Preferences().StringWithFallback("theme", "dark") == "light" {
		return theme.VariantLight
	}
	return theme.VariantDark
}

func newShell(a fyne.App, w fyne.Window, opts *Options) *shell {
	s := &shell{
		app:      a,
		win:      w,
		opts:     *opts,
		nav:      map[pageID]*navItem{},
		views:    map[pageID]fyne.CanvasObject{},
		dialogs:  map[sizedDialog]fyne.Size{},
		title:    widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subtitle: mutedText(""),
		body:     container.NewStack(),
	}
	s.title.SizeName = theme.SizeNameHeadingText
	s.subtitle.Wrapping = fyne.TextWrapOff
	s.subtitle.Truncation = fyne.TextTruncateEllipsis
	s.home = newHomePage(s)
	s.trust = newTrustPage(s)
	s.startBackend = startBackend
	s.files = workbench.NewFiles(workbench.RealFileEngines())
	s.network = workbench.NewNetwork(workbench.RealNetEngines())
	return s
}

func (s *shell) view() fyne.CanvasObject {
	s.navBox = container.NewVBox()
	for _, d := range pageDefs {
		n := newNavItem(d.icon, d.title, func() { s.show(d.id) })
		s.nav[d.id] = n
		s.navBox.Add(n)
	}
	for i := range s.statLines {
		s.statLines[i] = mutedText("")
		s.statLines[i].Wrapping = fyne.TextWrapOff
		s.statLines[i].SizeName = theme.SizeNameCaptionText
	}
	s.updateStats()

	s.sidebar = container.NewStack()
	s.collapsed = s.app.Preferences().Bool(prefSidebarCollapsed)
	s.layoutSidebar()
	sidebarBox := newPanel(s.sidebar)
	sidebarBox.fill = colorSidebar
	sidebarBox.inset = theme.Padding() * 2

	header := container.New(layout.NewCustomPaddedVBoxLayout(-10), s.title, s.subtitle)
	main := container.NewBorder(container.New(layout.NewCustomPaddedLayout(4, 12, 0, 0), header), nil, nil, nil, s.body)

	s.show(pageHome)
	s.frame = container.NewBorder(nil, nil, container.NewPadded(sidebarBox), nil,
		container.New(layout.NewCustomPaddedLayout(12, 16, 8, 24), main))
	return container.New(&resizeWatch{onResize: s.refitDialogs}, s.frame)
}

const prefSidebarCollapsed = "sidebarCollapsed"

// layoutSidebar fills the sidebar for the current collapsed state: brand, navigation, process stats,
// version and theme toggle, or only the icons.
func (s *shell) layoutSidebar() {
	for _, n := range s.nav {
		n.setCompact(s.collapsed)
	}
	collapseIcon := "sidebar-simple"
	if s.collapsed {
		collapseIcon = "caret-right"
	}
	toggle := newButton("", icon(collapseIcon), s.toggleSidebar)
	toggle.Importance = widget.LowImportance

	themeBtn := newButton("", s.themeIcon(), nil)
	themeBtn.Importance = widget.LowImportance
	themeBtn.OnTapped = func() {
		s.toggleTheme()
		themeBtn.SetIcon(s.themeIcon())
	}

	var top, bottom fyne.CanvasObject
	if s.collapsed {
		top = container.NewVBox(container.NewCenter(newLogo(30)), container.NewCenter(toggle), widget.NewSeparator(), s.navBox)
		bottom = container.NewVBox(widget.NewSeparator(), container.NewCenter(themeBtn))
	} else {
		brand := container.NewHBox(newLogo(34), container.NewCenter(
			widget.NewLabelWithStyle("app-listener", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})))
		version := mutedText(s.opts.Version)
		version.SizeName = theme.SizeNameCaptionText
		stats := container.New(layout.NewCustomPaddedVBoxLayout(-12), s.statLines[0], s.statLines[1], s.statLines[2])
		top = container.NewVBox(hfill(container.NewPadded(brand), container.NewCenter(toggle)), widget.NewSeparator(), s.navBox)
		bottom = container.NewVBox(widget.NewSeparator(), sectionTitle("USAGE"), stats, widget.NewSeparator(),
			hfill(version, themeBtn))
	}
	s.sidebar.Objects = []fyne.CanvasObject{container.NewBorder(top, bottom, nil, nil)}
	s.sidebar.Refresh()
	if s.frame != nil {
		s.frame.Refresh() // the sidebar's width changed: lay the page out again
	}
}

func (s *shell) toggleSidebar() {
	s.collapsed = !s.collapsed
	s.app.Preferences().SetBool(prefSidebarCollapsed, s.collapsed)
	s.layoutSidebar()
}

func (s *shell) updateStats() {
	cpu, mem, rate := s.stats.sample().lines()
	for i, text := range []string{cpu, mem, rate} {
		if s.statLines[i] != nil && s.statLines[i].Text != text {
			s.statLines[i].SetText(text)
		}
	}
}

func (s *shell) show(id pageID) {
	s.current = id
	for pid, n := range s.nav {
		n.setSelected(pid == id)
	}
	for _, d := range pageDefs {
		if d.id == id {
			s.title.SetText(d.title)
			s.subtitle.SetText(d.subtitle)
		}
	}
	v, ok := s.views[id]
	if ok {
		s.refreshPage(id)
	} else {
		// The page fills the window when it fits and scrolls when the window is smaller than it
		// (fyne grows scroll content to the viewport): a tiling WM may impose any size.
		v = newFadeScroll(s.build(id))
		s.views[id] = v
	}
	s.body.Objects = []fyne.CanvasObject{v}
	s.body.Refresh()
}

func (s *shell) build(id pageID) fyne.CanvasObject {
	switch id {
	case pageHome:
		return s.home.view()
	case pageFiles:
		s.filesPg = newFilesPage(s, s.files)
		return s.filesPg.view()
	case pageNetwork:
		s.netPg = newNetworkPage(s, s.network)
		return s.netPg.view()
	case pageProtected:
		s.protect = newProtectedPage(s)
		return s.protect.view()
	case pageUpdate:
		s.update = newUpdatePage(s)
		return s.update.view()
	case pageDaemon:
		s.daemon = newDaemonPage(s)
		return s.daemon.view()
	case pageTrust:
		return s.trust.view()
	default:
		s.install = newInstallPage(s)
		return s.install.view()
	}
}

// refreshPage re-reads the daemon state a page shows when it comes back into view: the daemon may
// have been started or stopped meanwhile.
func (s *shell) refreshPage(id pageID) {
	switch {
	case id == pageInstall && s.install != nil && !s.install.running:
		s.install.refreshStatus()
	case id == pageProtected && s.protect != nil:
		s.protect.probe()
	case id == pageTrust:
		s.trust.probe()
	case id == pageDaemon && s.daemon != nil && !s.daemon.running:
		s.daemon.refresh()
	}
}

// daemonChanged refreshes every built page whose actions depend on whether the daemon runs, after
// a run (any of them may have started or stopped it).
func (s *shell) daemonChanged() {
	if s.install != nil && !s.install.running {
		s.install.refreshStatus()
	}
	if s.protect != nil {
		s.protect.probe()
	}
	s.trust.probe()
	if s.daemon != nil && !s.daemon.running {
		s.daemon.refresh()
	}
}

// preferredSize is the window's size when nothing smaller is imposed (fyne units).
var preferredSize = fyne.NewSize(1280, 820)

// sizedDialog is a dialog showSized can resize.
type sizedDialog interface {
	Show()
	Resize(fyne.Size)
	SetOnClosed(func())
}

const dialogMargin = 24

// showSized shows d at want, capped to the window less a margin; refitDialogs keeps it capped as the
// window is resized (a tiling WM can shrink the window under an open dialog). Content larger than
// the cap must scroll: see scrollable.
func (s *shell) showSized(d sizedDialog, want fyne.Size) {
	s.dialogs[d] = want
	d.SetOnClosed(func() { delete(s.dialogs, d) })
	d.Show()
	handCursorsOn(s.win)
	s.fitDialog(d, want)
}

func (s *shell) fitDialog(d sizedDialog, want fyne.Size) {
	c := s.win.Canvas().Size()
	if c.Width <= 0 || c.Height <= 0 {
		d.Resize(want)
		return
	}
	d.Resize(fyne.NewSize(
		fyne.Max(1, fyne.Min(want.Width, c.Width-2*dialogMargin)),
		fyne.Max(1, fyne.Min(want.Height, c.Height-2*dialogMargin))))
	// A modal's overlay centers it only on layout: re-center at the new size.
	for _, o := range s.win.Canvas().Overlays().List() {
		o.Refresh()
	}
}

func (s *shell) refitDialogs(fyne.Size) {
	for d, want := range s.dialogs {
		s.fitDialog(d, want)
	}
}

// scrollable lets a dialog's content shrink below its natural size and scroll instead.
func scrollable(content fyne.CanvasObject) fyne.CanvasObject {
	return newFadeScroll(content)
}

// resizeWatch is a stack layout that reports size changes of the window content.
type resizeWatch struct {
	last     fyne.Size
	onResize func(fyne.Size)
}

func (r *resizeWatch) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objs {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(size)
	}
	if size != r.last {
		r.last = size
		if r.onResize != nil {
			r.onResize(size)
		}
	}
}

func (r *resizeWatch) MinSize(objs []fyne.CanvasObject) fyne.Size {
	var m fyne.Size
	for _, o := range objs {
		m = m.Max(o.MinSize())
	}
	return m
}

// focusSearch focuses the open picker's search, else the filter of the table on screen (Ctrl+F).
func (s *shell) focusSearch() {
	var e *widget.Entry
	switch {
	case s.picker != nil:
		s.picker.focusSearch()
		return
	case s.current == pageFiles && s.filesPg != nil && s.filesPg.current != nil:
		e = s.filesPg.current.searchEntry()
	case s.current == pageNetwork && s.netPg != nil:
		e = s.netPg.searchEntry()
	}
	if e != nil {
		s.win.Canvas().Focus(e)
	}
}

func (s *shell) setStatus(st sysStatus) {
	s.status = st
	s.home.applyStatus(st)
}

// themeIcon names the variant the toggle switches to.
func (s *shell) themeIcon() fyne.Resource {
	if savedVariant(s.app) == theme.VariantLight {
		return icon("moon")
	}
	return icon("sun")
}

func (s *shell) toggleTheme() {
	next := theme.VariantLight
	name := "light"
	if savedVariant(s.app) == theme.VariantLight {
		next, name = theme.VariantDark, "dark"
	}
	s.app.Preferences().SetString("theme", name)
	s.app.Settings().SetTheme(&appTheme{variant: next})
}

// busy shows a progress dialog while work runs off the UI goroutine, then reports its error or
// runs done.
func (s *shell) busy(title string, work func() error, done func()) {
	d := dialog.NewCustomWithoutButtons(title, widget.NewProgressBarInfinite(), s.win)
	d.Show()
	go func() {
		err := work()
		fyne.Do(func() {
			d.Hide()
			if err != nil {
				showErrorDialog(err, s.win)
				return
			}
			if done != nil {
				done()
			}
		})
	}()
}

func newLogo(side float32) fyne.CanvasObject {
	img := canvas.NewImageFromResource(logoResource)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSquareSize(side))
	return img
}
