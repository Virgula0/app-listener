//go:build gui

// Package gui is the desktop window `app-listener --gui` opens, behind the `gui` build tag: fyne
// adds X11/OpenGL/font/image dependencies (~15 MiB) and package-init goroutines that don't belong
// in the daemon sharing this binary. Build with `-tags gui` (or `make build GUI=1`).
package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

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
	pageInstall
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
	{pageInstall, "download-simple", "Install", "Set up the protection daemon"},
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
	files    *workbench.Files
	network  *workbench.Network
	install  *installPage
	// syncBackground makes background synchronous and every a no-op: fyne's test driver runs fyne.Do
	// inline on the calling goroutine, so a real goroutine would race the test's own rendering.
	syncBackground bool
}

// Run opens the window and blocks until it is closed.
func Run(opts Options) (err error) {
	restoreSessionEnv()
	defer cleanupSessionEnv()
	// fyne logs a failed GLFW init and then panics on the first window call.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("--gui: no display reachable (%v): run it from a graphical session "+
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
	w.SetContent(s.view())
	w.Resize(fyne.NewSize(1280, 820))
	w.CenterOnScreen()
	w.SetMaster()
	w.SetCloseIntercept(s.confirmClose)

	s.reprobe()
	w.ShowAndRun()
	s.network.Close()
	s.files.Close()
	if s.install != nil {
		s.install.cleanup()
	}
	return nil
}

// confirmClose asks before quitting mid-install: the installer would stop between two steps.
func (s *shell) confirmClose() {
	if s.install == nil || !s.install.running {
		s.win.Close()
		return
	}
	dialog.ShowConfirm("Quit during an installer run?", "Closing the window stops the installer midway. "+
		"Steps already completed stay applied; re-run the wizard afterwards to resume.", func(ok bool) {
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

func newShell(a fyne.App, w fyne.Window, opts Options) *shell {
	s := &shell{
		app:      a,
		win:      w,
		opts:     opts,
		nav:      map[pageID]*navItem{},
		views:    map[pageID]fyne.CanvasObject{},
		title:    widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subtitle: mutedText(""),
		body:     container.NewStack(),
	}
	s.title.SizeName = theme.SizeNameHeadingText
	s.home = newHomePage(s)
	s.files = workbench.NewFiles(workbench.RealFileEngines())
	s.network = workbench.NewNetwork(workbench.RealNetEngines())
	return s
}

func (s *shell) view() fyne.CanvasObject {
	items := container.NewVBox()
	for _, d := range pageDefs {
		n := newNavItem(d.icon, d.title, func() { s.show(d.id) })
		s.nav[d.id] = n
		items.Add(n)
	}

	brand := container.NewHBox(newLogo(34), container.NewVBox(
		widget.NewLabelWithStyle("app-listener", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	))
	version := mutedText(s.opts.Version)
	version.SizeName = theme.SizeNameCaptionText

	themeBtn := widget.NewButtonWithIcon("", s.themeIcon(), nil)
	themeBtn.Importance = widget.LowImportance
	themeBtn.OnTapped = func() {
		s.toggleTheme()
		themeBtn.SetIcon(s.themeIcon())
	}

	sidebar := container.NewBorder(
		container.NewVBox(container.NewPadded(brand), widget.NewSeparator(), items),
		container.NewVBox(widget.NewSeparator(), hfill(version, themeBtn)),
		nil, nil,
	)
	sidebarBox := newPanel(sidebar)
	sidebarBox.fill = colorSidebar
	sidebarBox.inset = theme.Padding() * 2

	header := container.New(layout.NewCustomPaddedVBoxLayout(-10), s.title, s.subtitle)
	main := container.NewBorder(container.New(layout.NewCustomPaddedLayout(4, 12, 0, 0), header), nil, nil, nil, s.body)

	s.show(pageHome)
	return container.NewBorder(nil, nil, container.NewPadded(sidebarBox), nil,
		container.New(layout.NewCustomPaddedLayout(12, 16, 8, 24), main))
}

func (s *shell) show(id pageID) {
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
	if !ok {
		v = s.build(id)
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
		return newFilesPage(s, s.files).view()
	case pageNetwork:
		return newNetworkPage(s, s.network).view()
	default:
		s.install = newInstallPage(s)
		return s.install.view()
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
				dialog.ShowError(err, s.win)
				return
			}
			if done != nil {
				done()
			}
		})
	}()
}

// startAt opens d in dir, or the invoking user's home when empty: the GUI runs as root under sudo,
// and the secrets worth watching live in that user's home.
func (s *shell) startAt(d *dialog.FileDialog, dir string) {
	if dir == "" {
		dir = invokingHome()
	}
	if dir == "" {
		return
	}
	if l, err := storage.ListerForURI(storage.NewFileURI(dir)); err == nil {
		d.SetLocation(l)
	}
}

func newLogo(side float32) fyne.CanvasObject {
	img := canvas.NewImageFromResource(logoResource)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSquareSize(side))
	return img
}
