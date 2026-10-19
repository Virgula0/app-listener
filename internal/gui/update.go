//go:build gui

package gui

import (
	"errors"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// updateChannels are the Update page's sources: `update --channel` values, no local file.
var updateChannels = []struct {
	label      string
	preRelease bool
}{
	{"Latest stable release", false},
	{"Latest pre-release", true},
}

// latestRelease is the newest release of one channel, as last checked.
type latestRelease struct {
	checking bool
	tag      string
	newer    bool
	err      error
}

// updatePage is the Update workspace: `app-listener update` run in-process (update.RunWith) with its
// changelog, confirmations and unit diff as dialogs. It deploys the release's headless asset.
type updatePage struct {
	runHost
	hooks                   UpdateHooks
	versionOf, guiVersionOf func() string

	preRelease bool
	// picked is set once the user chose a channel; until then it follows the installed version's.
	picked bool
	// loaded is set once the installed binary and the daemon were read: Latest compares against it.
	loaded   bool
	present  installedBinary
	daemonUp bool
	latest   map[bool]*latestRelease

	installedT, daemonT, latestT *statusTile
	intro                        *widget.Label
	card                         *fyne.Container
	updateBtn, checkBtn          *handButton
	pick                         *handSelect
	withGUI                      *handCheck
	banner, content              *fyne.Container
	bannerText                   *widget.Label
}

func newUpdatePage(s *shell) *updatePage {
	p := &updatePage{runHost: runHost{s: s}, hooks: s.opts.Update, versionOf: s.opts.Install.InstalledVersion,
		guiVersionOf: s.opts.Install.InstalledGUIVersion, latest: map[bool]*latestRelease{}}
	p.settled = func() {
		clear(p.latest)
		p.refreshStatus()
	}
	return p
}

func (p *updatePage) view() fyne.CanvasObject {
	p.installedT = newStatusTile("terminal-window", "Installed binary")
	p.daemonT = newStatusTile("pulse", "Daemon")
	p.latestT = newStatusTile("download-simple", "Latest release")

	p.content = container.NewVBox(
		p.bannerCard(),
		container.NewGridWithColumns(3, p.installedT.card, p.daemonT.card, p.latestT.card),
		sectionTitle("SOURCE"),
		p.sourceCard(),
		sectionTitle("WHAT HAPPENS"),
		newPanel(mutedText("1. The release's headless asset for this architecture is downloaded with its sha256 "+
			"checksum and Ed25519 signature. The signature, the checksum and the GitHub digest are verified, "+
			"and a failed check stops the update.\n"+
			"2. The changelog since the installed version is shown, and you confirm.\n"+
			"3. If the release ships a different daemon unit, its diff is shown and you choose whether to replace "+
			"the installed one. Local edits to the unit are not carried over.\n"+
			"4. The daemon is stopped, its folders are verified locked, the binary is replaced atomically and the "+
			"daemon is started again.\n"+
			"5. The release's desktop build, downloaded and verified with the headless one in step 1, replaces "+
			systemd.GUIBinaryPath+" (installed there if missing), when the desktop build option above is "+
			"checked. Unchecked, that file is left as is. An open window keeps running its old version until "+
			"you restart it.\n\n"+
			"Nothing changes before step 4. When the daemon is already the newest, Update only brings the desktop "+
			"build to the daemon's version. This window's own executable is never copied.")))
	over := newFadeVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), p.content))
	root := p.mount(over, "Back to Update")
	p.updateBtn.Disable()
	p.renderLatest()
	p.refreshStatus()
	return root
}

func (p *updatePage) sourceCard() fyne.CanvasObject {
	labels := make([]string, len(updateChannels))
	for i, c := range updateChannels {
		labels[i] = c.label
	}
	p.pick = newSelect(labels, func(l string) {
		p.picked = true
		p.setChannel(updateChannels[slices.Index(labels, l)].preRelease)
	})
	p.pick.Selected = labels[0] // SetSelected fires the callback before the card exists
	p.intro = mutedText("")
	p.checkBtn = newButton("Check", icon("arrow-clockwise"), p.checkLatest)
	p.updateBtn = newButton("Update", icon("download-simple"), p.start)
	p.updateBtn.Importance = widget.HighImportance
	p.withGUI = newCheck("Also install or update the desktop build (app-listener-gui)",
		func(bool) { p.renderLatest() })
	p.withGUI.Checked = true
	p.card = container.NewVBox(
		container.NewBorder(nil, nil, boldLabel("Channel"), container.NewHBox(p.checkBtn, p.updateBtn), p.pick),
		p.intro, p.withGUI)
	p.renderSource()
	return newPanel(p.card)
}

func (p *updatePage) setChannel(pre bool) {
	p.preRelease = pre
	p.renderSource()
	p.renderLatest()
	if p.latest[p.preRelease] == nil {
		p.checkLatest()
	}
}

// bannerCard is the "update available" strip, shown while the selected channel has a newer release.
func (p *updatePage) bannerCard() fyne.CanvasObject {
	p.bannerText = mutedText("")
	btn := newButton("Update now", icon("download-simple"), p.start)
	btn.Importance = widget.HighImportance
	p.banner = container.NewVBox(newPanel(container.NewBorder(nil, nil,
		newIconTile("download-simple", colorInfo, 40), container.NewCenter(btn),
		container.NewVBox(boldLabel("A new version is available"), p.bannerText))))
	p.banner.Hide()
	return p.banner
}

func (p *updatePage) renderBanner(show bool) {
	if show {
		st := p.latest[p.preRelease]
		p.bannerText.SetText("The " + strings.ToLower(updateChannels[boolIndex(p.preRelease)].label) + " is " +
			st.tag + ". Installed: " + p.present.version + ".")
	}
	if show == p.banner.Visible() {
		return
	}
	p.banner.Hidden = !show
	p.content.Refresh() // a child's visibility change doesn't relayout its parent
}

func (p *updatePage) renderSource() {
	const tail = " A network whitelist active in this window would block the download."
	if p.preRelease {
		p.intro.SetText("The newest pre-release (the pre-YYYYMMDD builds of main), what " +
			"update --channel pre-release installs. An installed stable version is always replaced." + tail)
	} else {
		p.intro.SetText("The newest stable release (vX.Y.Z), what update installs by default. It is " +
			"installed only when newer than the installed version." + tail)
	}
	p.card.Refresh()
}

// refreshStatus re-reads the installed binary and the daemon off the UI goroutine, then checks the
// channel's latest release against the installed version.
func (p *updatePage) refreshStatus() {
	versionOf, guiVersionOf := p.versionOf, p.guiVersionOf
	var ib installedBinary
	var up bool
	p.s.background(func() {
		ib = readInstalled(versionOf, guiVersionOf)
		up = protected.DaemonRunning()
	}, func() {
		p.present, p.daemonUp, p.loaded = ib, up, true
		if !p.picked && ib.present {
			p.preRelease = isPreRelease(ib.version)
			p.pick.Selected = updateChannels[0].label
			if p.preRelease {
				p.pick.Selected = updateChannels[1].label
			}
			p.pick.Refresh()
			p.renderSource()
		}
		p.renderStatus()
		if p.present.present && p.latest[p.preRelease] == nil {
			p.checkLatest()
		}
	})
}

// checkLatest looks up the selected channel's newest release (refreshStatus does once loaded).
func (p *updatePage) checkLatest() {
	if !p.loaded || !p.present.present {
		return
	}
	latest := p.hooks.Latest
	pre, installed := p.preRelease, p.present.version
	st := &latestRelease{checking: true}
	p.latest[pre] = st
	p.renderLatest()
	var tag string
	var newer bool
	var err error
	p.s.background(func() {
		if latest == nil {
			err = errors.New("release lookup is not available in this build")
			return
		}
		tag, newer, err = latest(pre, installed)
	}, func() {
		if p.latest[pre] != st {
			return // superseded by a newer check or a finished run
		}
		*st = latestRelease{tag: tag, newer: newer, err: err}
		p.renderLatest()
		if err == nil && pre == isPreRelease(installed) {
			p.s.setUpdateAvailable(tag, newer)
		}
	})
}

func (p *updatePage) renderStatus() {
	ib := p.present
	switch {
	case ib.err != nil:
		p.installedT.set(false, "unreadable", ib.err.Error(), "")
	case !ib.present:
		p.installedT.set(true, "not installed", "Nothing at "+systemd.InstallBinaryPath+
			": set the daemon up from Install first.", colorMuted)
	default:
		desktop := "desktop build not installed"
		if ib.guiVersion != "" {
			desktop = "desktop build " + ib.guiVersion
		}
		p.installedT.set(true, ib.version, systemd.InstallBinaryPath+" · "+desktop, theme.ColorNameSuccess)
	}
	if p.daemonUp {
		p.daemonT.set(true, "running", "The update stops it before replacing the binary and starts it again.",
			theme.ColorNameSuccess)
	} else {
		p.daemonT.set(true, "stopped", "No daemon is running: only the binary is replaced.", theme.ColorNameWarning)
	}
	p.renderLatest()
	if ib.present && ib.err == nil && !p.running {
		p.updateBtn.Enable()
	} else {
		p.updateBtn.Disable()
	}
}

func (p *updatePage) renderLatest() {
	st := p.latest[p.preRelease]
	p.renderBanner(st != nil && !st.checking && st.err == nil && st.newer && !p.running)
	switch {
	case st == nil && p.loaded && !p.present.present:
		p.latestT.set(true, "not checked", "Nothing is installed to update.", colorMuted)
	case st == nil:
		p.latestT.set(true, "not checked", "Check looks it up on GitHub.", colorMuted)
	case st.checking:
		p.latestT.set(true, "checking…", "Asking GitHub for the newest release.", colorMuted)
	case st.err != nil:
		p.latestT.set(false, "unknown", st.err.Error(), "")
	case st.newer:
		p.latestT.set(true, st.tag, "Newer than the installed version: Update installs it.", colorInfo)
	case p.present.guiVersion != p.present.version && p.withGUI.Checked:
		p.latestT.set(true, st.tag, "The daemon is up to date: Update installs its matching desktop build.",
			colorInfo)
	default:
		p.latestT.set(true, st.tag, "The installed version is up to date: Update changes nothing.",
			theme.ColorNameSuccess)
	}
}

func (p *updatePage) start() {
	run := p.hooks.Run
	pre := p.preRelease
	withGUI := p.withGUI.Checked
	name := updateChannels[boolIndex(pre)].label
	if p.begin("Update: "+name, "arrow-clockwise", func(ui wizard.UI) error {
		if run == nil {
			return errors.New("the updater is not available in this build")
		}
		return run(ui, pre, withGUI)
	}, nil) {
		p.updateBtn.Disable()
		p.renderBanner(false)
	}
}
