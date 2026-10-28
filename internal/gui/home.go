//go:build gui

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// statusTile is one system-check card: icon tile, caption, value pill, detail line.
type statusTile struct {
	tile   *iconTile
	value  *pill
	detail *widget.Label
	card   *panel
}

func newStatusTile(iconName, caption string) *statusTile {
	s := &statusTile{
		tile:   newIconTile(iconName, colorMuted, 40),
		value:  newPill("checking…", colorMuted),
		detail: mutedText(""),
	}
	head := container.NewBorder(nil, nil, s.tile, nil,
		container.NewVBox(boldLabel(caption), container.NewHBox(s.value)))
	s.card = newPanel(container.NewVBox(head, s.detail))
	return s
}

func (s *statusTile) set(ok bool, value, detail string, okColor fyne.ThemeColorName) {
	c := okColor
	if !ok {
		c = theme.ColorNameError
	}
	s.tile.Set(s.tile.name, c)
	s.value.Set(value, c)
	s.detail.SetText(detail)
}

// homePage is the landing view: system checks and the workspace launcher.
type homePage struct {
	shell                     *shell
	privs, ebpfT, lsm, daemon *statusTile
}

func newHomePage(s *shell) *homePage {
	return &homePage{
		shell:  s,
		privs:  newStatusTile("fingerprint", "Privileges"),
		ebpfT:  newStatusTile("cpu", "eBPF"),
		lsm:    newStatusTile("shield-check", "BPF LSM"),
		daemon: newStatusTile("pulse", "Daemon"),
	}
}

func (h *homePage) view() fyne.CanvasObject {
	hero := container.NewBorder(nil, nil, newLogo(64), nil, container.NewVBox(
		widget.NewLabelWithStyle("Welcome to app-listener", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mutedText("Watch and lock down who touches your secrets. eBPF observes every file and network "+
			"operation; the LSM guard denies the ones you don't allow, keyed on the executable's inode."),
	))

	checks := newAdaptiveGrid(4, h.privs.card, h.ebpfT.card, h.lsm.card, h.daemon.card)

	cards := newAdaptiveGrid(3,
		h.workspaceCard("folders", theme.ColorNamePrimary, "Files",
			"Monitor folders and files, then guard them: whitelist or blacklist the processes you see, or any binary.",
			pageFiles),
		h.workspaceCard("globe-hemisphere-west", colorInfo, "Network",
			"Monitor the network activity of chosen binaries and enforce a system-wide allow or deny list.",
			pageNetwork),
		h.workspaceCard("lock-simple", theme.ColorNameWarning, "Protected",
			"Edit the daemon's protected directories, grant a binary temporary access, or edit its configuration.",
			pageProtected),
		h.workspaceCard("shield-check", theme.ColorNamePrimary, "Trust",
			"Confirm whitelisted binaries the daemon refused because their content changed or is new.",
			pageTrust),
		h.workspaceCard("arrow-clockwise", colorInfo, "Update",
			"Update the installed daemon to the latest signed stable or pre-release build.",
			pageUpdate),
		h.workspaceCard("download-simple", theme.ColorNameSuccess, "Install",
			"Set up the protection daemon: pick what to guard, encrypt it with fscrypt and enable the service.",
			pageInstall),
	)

	// The right inset keeps the cards clear of the scroll bar.
	return newFadeVScroll(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 14), newGapVBox(
		newPanel(hero),
		container.NewVBox(sectionTitle("SYSTEM"), checks),
		container.NewVBox(sectionTitle("WORKSPACES"), cards),
		container.NewVBox(sectionTitle("WINDOW"), h.windowCard()),
	)))
}

// windowCard holds the window settings: floating on tiling window managers.
func (h *homePage) windowCard() fyne.CanvasObject {
	s := h.shell
	detail := "No tiling window manager detected (Hyprland, sway and i3 are supported): the setting " +
		"applies when the window opens on one."
	if s.wm.Tiling() {
		detail = "Detected " + s.wm.Kind.String() + ". Unchecking it tiles the window again."
	}
	float := newCheck("Open as a floating window on tiling window managers", func(on bool) {
		s.app.Preferences().SetBool(prefFloatTiling, on)
		s.setFloating(on)
	})
	float.Checked = s.floatPref()
	return newPanel(container.NewVBox(float, mutedText(detail)))
}

func (h *homePage) workspaceCard(iconName string, c fyne.ThemeColorName, title, desc string, target pageID) fyne.CanvasObject {
	open := newButton("Open", icon("caret-right"), func() { h.shell.show(target) })
	open.IconPlacement = widget.ButtonIconTrailingText
	open.Importance = widget.LowImportance
	body := container.NewVBox(
		newIconTile(iconName, c, 48),
		widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mutedText(desc),
		container.NewHBox(open),
	)
	return newClickPanel(body, func() { h.shell.show(target) })
}

// applyStatus fills the check tiles once probeSystem returns.
func (h *homePage) applyStatus(st sysStatus) {
	if st.root {
		h.privs.set(true, "root", "Running as uid 0: eBPF programs can load.", theme.ColorNameSuccess)
	} else {
		h.privs.set(false, "not root", "Restart with sudo: every workspace loads eBPF programs.", "")
	}
	if st.ebpfErr == nil {
		h.ebpfT.set(true, "available", "Kernel accepts BPF programs: monitors can run.", theme.ColorNameSuccess)
	} else {
		h.ebpfT.set(false, "unavailable", st.ebpfErr.Error(), "")
	}
	if st.lsmErr == nil {
		h.lsm.set(true, "active", "bpf is in the active LSM list: guards can deny.", theme.ColorNameSuccess)
	} else {
		h.lsm.set(false, "inactive", "Guards are disabled (monitors still work). "+st.lsmErr.Error(), "")
	}
	if st.daemon {
		h.daemon.set(true, "running", "The protection daemon is enforcing its configuration.", theme.ColorNameSuccess)
	} else {
		h.daemon.set(true, "stopped", "No daemon is running: start it on the Daemon page, or set one up with Install.", theme.ColorNameWarning)
	}
}
