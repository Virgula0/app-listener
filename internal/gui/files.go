//go:build gui

package gui

import (
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/workbench"
)

// filesPage is the Files workspace: a target list and the selected target's view.
type filesPage struct {
	shell   *shell
	wb      *workbench.Files
	targets []*workbench.FileTarget
	list    *widget.List
	detail  *fyne.Container
	views   map[int]*targetView
	current *targetView

	root      *fyne.Container
	side      fyne.CanvasObject
	offset    float64
	collapsed bool
}

const prefTargetsCollapsed = "targetsCollapsed"

const pollInterval = 400 * time.Millisecond

func newFilesPage(s *shell, wb *workbench.Files) *filesPage {
	return &filesPage{shell: s, wb: wb, views: map[int]*targetView{}, detail: container.NewStack()}
}

func (p *filesPage) view() fyne.CanvasObject {
	p.list = widget.NewList(
		func() int { return len(p.targets) },
		func() fyne.CanvasObject { return newTargetRow() },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < len(p.targets) {
				o.(*targetRow).set(p.targets[id])
			}
		})
	p.list.OnSelected = func(id widget.ListItemID) {
		if id < len(p.targets) {
			p.showTarget(p.targets[id])
		}
	}

	add := newButton("Add", icon("plus"), p.addTarget)
	add.Importance = widget.HighImportance
	collapse := newButton("", icon("sidebar-simple"), p.toggleTargets)
	collapse.Importance = widget.LowImportance
	header := hfill(boldLabel("Targets"), add, collapse)
	side := newPanel(container.NewBorder(header, nil, nil, nil, p.list))
	side.inset = theme.Padding() * 2
	p.side = side

	p.showEmpty()
	p.offset = 0.24
	p.root = container.NewStack()
	p.collapsed = p.shell.app.Preferences().Bool(prefTargetsCollapsed)
	p.layoutTargets()
	p.shell.every(pollInterval, func() {
		if p.current != nil {
			p.current.refresh(false)
		}
		p.list.Refresh()
	})
	return p.root
}

// layoutTargets shows the target list beside the detail, or a narrow rail when collapsed.
func (p *filesPage) layoutTargets() {
	if p.collapsed {
		expand := newButton("", icon("caret-right"), p.toggleTargets)
		expand.Importance = widget.LowImportance
		add := newButton("", icon("plus"), p.addTarget)
		add.Importance = widget.HighImportance
		rail := newPanel(container.NewVBox(expand, add))
		rail.inset = theme.Padding() * 2
		p.root.Objects = []fyne.CanvasObject{container.NewBorder(nil, nil, rail, nil, p.detail)}
	} else {
		split := container.NewHSplit(p.side, p.detail)
		split.Offset = p.offset
		p.root.Objects = []fyne.CanvasObject{split}
	}
	p.root.Refresh()
}

func (p *filesPage) toggleTargets() {
	if split, ok := p.root.Objects[0].(*container.Split); ok {
		p.offset = split.Offset
	}
	p.collapsed = !p.collapsed
	p.shell.app.Preferences().SetBool(prefTargetsCollapsed, p.collapsed)
	p.layoutTargets()
}

func (p *filesPage) showEmpty() {
	p.current = nil
	addBtn := newButton("Add a folder or file", icon("plus"), p.addTarget)
	addBtn.Importance = widget.HighImportance
	p.detail.Objects = []fyne.CanvasObject{container.NewCenter(container.NewVBox(
		container.NewCenter(newIconTile("folders", theme.ColorNamePrimary, 72)),
		widget.NewLabelWithStyle("Nothing watched yet", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Add the folders or files to watch. Each one gets its own monitor and,\n"+
			"once you have seen who uses it, its own guard.", fyne.TextAlignCenter, fyne.TextStyle{}),
		container.NewCenter(addBtn),
	))}
	p.detail.Refresh()
}

func (p *filesPage) showTarget(t *workbench.FileTarget) {
	v, ok := p.views[t.ID]
	if !ok {
		v = newTargetView(p, t)
		p.views[t.ID] = v
	}
	p.current = v
	p.detail.Objects = []fyne.CanvasObject{v.root}
	p.detail.Refresh()
	v.refresh(true)
}

func (p *filesPage) reloadTargets(selectT *workbench.FileTarget) {
	p.targets = p.wb.Targets()
	p.list.Refresh()
	for i, t := range p.targets {
		if t == selectT {
			p.list.Select(i)
			return
		}
	}
	if len(p.targets) == 0 {
		p.list.UnselectAll()
		p.showEmpty()
	}
}

func (p *filesPage) removeTarget(t *workbench.FileTarget) {
	msg := fmt.Sprintf("Stop monitoring %s?", t.Spec.Path)
	if t.Policy().Mode != workbench.GuardOff {
		msg += "\n\nIts guard stops too: the path is no longer protected by this window."
	}
	dialog.ShowConfirm("Remove target", msg, func(ok bool) {
		if !ok {
			return
		}
		p.shell.busy("Removing target…", func() error {
			p.wb.Remove(t)
			return nil
		}, func() {
			delete(p.views, t.ID)
			p.reloadTargets(nil)
			if len(p.targets) > 0 {
				p.list.Select(0)
			}
		})
	}, p.shell.win)
}

// poll refreshes the visible target and the list's pills while the window lives.
// targetRow is one entry of the target list.
type targetRow struct {
	widget.BaseWidget
	tile  *iconTile
	name  *widget.Label
	path  *widget.Label
	guard *pill
	box   *fyne.Container
}

func newTargetRow() *targetRow {
	r := &targetRow{
		tile:  newIconTile("folder", theme.ColorNamePrimary, 32),
		name:  boldLabel(""),
		path:  mutedText(""),
		guard: newPill("", colorMuted),
	}
	r.name.Truncation = fyne.TextTruncateEllipsis
	r.path.Wrapping = fyne.TextWrapOff
	r.path.Truncation = fyne.TextTruncateEllipsis
	r.path.SizeName = theme.SizeNameCaptionText
	r.ExtendBaseWidget(r)
	return r
}

func (r *targetRow) set(t *workbench.FileTarget) {
	name := filepath.Base(t.Spec.Path)
	iconName := "file"
	if t.Spec.IsDir {
		iconName = "folder"
	}
	r.tile.Set(iconName, theme.ColorNamePrimary)
	r.name.SetText(name)
	r.path.SetText(t.Spec.Path)
	text, c := guardPill(t.Policy())
	r.guard.Set(text, c)
	if r.box != nil {
		r.box.Refresh() // the pill's width changed: lay the row out again
	}
}

func (r *targetRow) CreateRenderer() fyne.WidgetRenderer {
	texts := container.NewVBox(r.name, r.path, container.NewHBox(r.guard))
	r.box = container.NewBorder(nil, nil, container.NewCenter(r.tile), nil, texts)
	return widget.NewSimpleRenderer(r.box)
}

// guardPill is the short guard status shown on a target.
func guardPill(p workbench.Policy) (string, fyne.ThemeColorName) {
	switch p.Mode {
	case workbench.GuardWhitelist:
		return fmt.Sprintf("whitelist · %d", len(p.Binaries)), theme.ColorNameSuccess
	case workbench.GuardBlacklist:
		return fmt.Sprintf("blacklist · %d", len(p.Binaries)), theme.ColorNameWarning
	default:
		return "unguarded", colorMuted
	}
}
