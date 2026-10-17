//go:build gui

package gui

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// targetView is one target's detail: header, counters and the Activity / Processes / Guard tabs.
type targetView struct {
	page *filesPage
	t    *workbench.FileTarget
	root fyne.CanvasObject

	guardPill *pill
	headerBox *fyne.Container
	errLabel  *widget.Label
	counters  *counterStrip

	activity *eventFeed
	procs    *processView
	policy   *policyTab
	tabs     *container.AppTabs

	countsVersion uint64
}

func newTargetView(p *filesPage, t *workbench.FileTarget) *targetView {
	v := &targetView{page: p, t: t, guardPill: newPill("", colorMuted)}
	v.errLabel = widget.NewLabel("")
	v.errLabel.Importance = widget.DangerImportance
	v.errLabel.Wrapping = fyne.TextWrapWord
	v.errLabel.Hide()

	v.policy = newPolicyTab(v)
	v.activity = newFileFeed(v)
	v.procs = newFileProcesses(v)
	v.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Activity", icon("pulse"), v.activity.root),
		container.NewTabItemWithIcon("Processes", icon("cpu"), v.procs.root),
		container.NewTabItemWithIcon("Guard", icon("shield-check"), v.policy.root),
	)
	v.counters = newCounterStrip("Denied")
	v.root = container.NewBorder(newGapVBox(v.header(), v.counters.root, v.errLabel), nil, nil, nil, v.tabs)
	return v
}

// searchEntry is the filter of the visible table tab, nil on the others.
func (v *targetView) searchEntry() *widget.Entry {
	switch v.tabs.SelectedIndex() {
	case 0:
		return v.activity.search
	case 1:
		return v.procs.search
	}
	return nil
}

func (v *targetView) header() fyne.CanvasObject {
	iconName := "file"
	if v.t.Spec.IsDir {
		iconName = "folder"
	}
	title := boldLabel(v.t.Spec.Path)
	title.Truncation = fyne.TextTruncateEllipsis
	remove := widget.NewButtonWithIcon("Remove", icon("trash"), func() { v.page.removeTarget(v.t) })
	remove.Importance = widget.DangerImportance
	summary := widget.NewLabel(specSummary(v.t.Spec))
	summary.Importance = widget.LowImportance
	summary.Truncation = fyne.TextTruncateEllipsis
	pills := container.NewHBox(newPill("monitoring", colorInfo), v.guardPill)
	v.headerBox = container.NewBorder(nil, nil, container.NewCenter(pills), nil, summary)
	body := container.NewBorder(nil, nil, container.NewCenter(newIconTile(iconName, theme.ColorNamePrimary, 44)),
		container.NewCenter(remove), container.NewVBox(title, v.headerBox))
	return newPanel(body)
}

func specSummary(s workbench.FileSpec) string {
	parts := []string{"single file"}
	if s.IsDir {
		switch {
		case !s.Recursive:
			parts = []string{"top level only"}
		case s.Depth == 0:
			parts = []string{"all subfolders"}
		default:
			parts = []string{fmt.Sprintf("subfolders to depth %d", s.Depth)}
		}
	}
	if len(s.Events) == 0 || len(s.Events) == len(ebpf.EventTypes()) {
		parts = append(parts, "every operation")
	} else {
		names := make([]string, len(s.Events))
		for i, e := range s.Events {
			names[i] = e.String()
		}
		parts = append(parts, strings.Join(names, ", "))
	}
	return strings.Join(parts, " · ")
}

// counterStrip is a row of four stat cards: events, allowed, denied (deniedCaption) and processes.
type counterStrip struct {
	root   fyne.CanvasObject
	labels [4]*widget.Label
}

func newCounterStrip(deniedCaption string) *counterStrip {
	specs := []struct {
		icon, caption string
		c             fyne.ThemeColorName
	}{
		{"pulse", "Events", colorInfo},
		{"check-circle", "Allowed", theme.ColorNameSuccess},
		{"prohibit", deniedCaption, theme.ColorNameError},
		{"cpu", "Processes", theme.ColorNamePrimary},
	}
	cs := &counterStrip{}
	cards := make([]fyne.CanvasObject, len(specs))
	for i, sp := range specs {
		cs.labels[i] = widget.NewLabelWithStyle("0", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		cs.labels[i].SizeName = theme.SizeNameSubHeadingText
		caption := mutedText(sp.caption)
		caption.Wrapping = fyne.TextWrapOff
		caption.SizeName = theme.SizeNameCaptionText
		cards[i] = newPanel(container.NewBorder(nil, nil, container.NewCenter(newIconTile(sp.icon, sp.c, 34)), nil,
			container.NewVBox(cs.labels[i], caption)))
	}
	cs.root = newAdaptiveGrid(len(cards), cards...)
	return cs
}

func (cs *counterStrip) set(c workbench.Counts, processes int) {
	cs.labels[0].SetText(strconv.FormatUint(c.Observed+c.Allowed+c.Denied, 10))
	cs.labels[1].SetText(strconv.FormatUint(c.Allowed, 10))
	cs.labels[2].SetText(strconv.FormatUint(c.Denied, 10))
	cs.labels[3].SetText(strconv.Itoa(processes))
}

// refresh redraws what changed since the last poll (force: everything).
func (v *targetView) refresh(force bool) {
	if text, c := guardPill(v.t.Policy()); text != v.guardPill.text {
		v.guardPill.Set(text, c)
		v.headerBox.Refresh() // the pill's width changed: lay the line out again
	}
	if e := v.t.LastError(); e != "" {
		v.errLabel.SetText("Guard error: " + e)
		v.errLabel.Show()
	} else {
		v.errLabel.Hide()
	}

	counts, version := v.t.Counts()
	if force || version != v.countsVersion {
		v.countsVersion = version
		v.counters.set(counts, len(v.t.Processes()))
		v.activity.refresh(force)
		v.procs.refresh()
	}
	if force {
		v.policy.refresh()
	}
}
