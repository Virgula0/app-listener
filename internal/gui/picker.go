//go:build gui

package gui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type pickMode int

const (
	pickFile pickMode = iota
	pickFolder
)

// picker is the file and folder chooser: fyne's FileDialog has no search, and opens the chosen
// file (as root) just to hand back its path.
type picker struct {
	s      *shell
	mode   pickMode
	onPick func(string)

	dir    string
	all    []pickEntry
	shown  []pickEntry
	cursor int
	chosen string
	opt    listOpts
	deep   bool
	// keying marks a selection made from the search field's arrow keys: it highlights a folder
	// instead of entering it.
	keying bool
	// gen stops a running subfolder search once the query, folder or dialog changes.
	gen atomic.Uint64

	d        *dialog.CustomDialog
	list     *widget.List
	path     *keyEntry
	search   *keyEntry
	selected *widget.Label
	status   *widget.Label
	open     *handButton
}

// pickPath asks for a file or folder, starting in start (the invoking user's home when empty), and
// passes the chosen absolute path to onPick.
func (s *shell) pickPath(mode pickMode, title, start string, onPick func(string)) {
	p := &picker{s: s, mode: mode, onPick: onPick, cursor: -1, opt: listOpts{foldersOnly: mode == pickFolder}}
	p.d = dialog.NewCustomWithoutButtons(title, p.content(), s.win)
	p.d.SetButtons([]fyne.CanvasObject{newTextButton("Cancel", p.close), p.open})
	p.d.SetOnClosed(func() {
		p.gen.Add(1)
		if s.picker == p {
			s.picker = nil
		}
	})
	if start == "" {
		start = invokingHome()
	}
	if fi, err := os.Stat(start); err != nil || !fi.IsDir() {
		start = "/"
	}
	s.picker = p
	p.navigate(start)
	s.showSized(p.d, fyne.NewSize(900, 620))
	s.win.Canvas().Focus(p.search)
}

func (p *picker) content() fyne.CanvasObject {
	p.path = newKeyEntry()
	p.path.OnSubmitted = p.goTo
	p.path.onShortcut = p.findShortcut
	p.path.onKey = func(k fyne.KeyName) bool {
		if k != fyne.KeyEscape {
			return false
		}
		p.close()
		return true
	}
	up := newButton("", theme.MoveUpIcon(), func() { p.navigate(filepath.Dir(p.dir)) })
	hidden := widget.NewCheck("Hidden", func(on bool) { p.opt.hidden = on; p.reload() })

	p.search = newKeyEntry()
	p.search.SetPlaceHolder("Search this folder (Ctrl+F)")
	p.search.ActionItem = widget.NewIcon(icon("magnifying-glass"))
	p.search.OnChanged = func(string) { p.refilter() }
	p.search.onKey = p.searchKey
	p.search.onShortcut = p.findShortcut
	deep := widget.NewCheck("Subfolders", func(on bool) { p.deep = on; p.refilter() })

	p.list = widget.NewList(func() int { return len(p.shown) }, newPickRow, p.updateRow)
	p.list.OnSelected = p.selectRow

	p.selected = widget.NewLabel("")
	p.selected.Truncation = fyne.TextTruncateEllipsis
	p.status = widget.NewLabel("")
	p.status.Importance = widget.LowImportance
	label := "Open"
	if p.mode == pickFolder {
		label = "Choose folder"
	}
	p.open = newButton(label, nil, p.confirm)
	p.open.Importance = widget.HighImportance
	p.open.Disable()

	top := container.NewVBox(hfill(p.path, up, hidden), hfill(p.search, deep))
	bottom := hfill(p.selected, p.status)
	return container.NewBorder(top, bottom, container.NewVScroll(p.places()), nil, p.list)
}

func (p *picker) places() fyne.CanvasObject {
	home := invokingHome()
	box := container.NewVBox()
	for _, pl := range []struct{ name, path, icon string }{
		{"Home", home, "house"},
		{"~/.config", filepath.Join(home, ".config"), "gear"},
		{"~/.local/share", filepath.Join(home, ".local", "share"), "folder"},
		{"/", "/", "folder"},
		{"/usr/bin", "/usr/bin", "terminal-window"},
		{"/etc", "/etc", "folder"},
		{"/opt", "/opt", "folder"},
	} {
		if !filepath.IsAbs(pl.path) {
			continue
		}
		if fi, err := os.Stat(pl.path); err != nil || !fi.IsDir() {
			continue
		}
		b := newButton(pl.name, icon(pl.icon), func() { p.navigate(pl.path) })
		b.Alignment = widget.ButtonAlignLeading
		b.Importance = widget.LowImportance
		box.Add(b)
	}
	return box
}

func newPickRow() fyne.CanvasObject {
	l := widget.NewLabel("")
	l.Truncation = fyne.TextTruncateEllipsis
	return container.NewBorder(nil, nil, widget.NewIcon(icon("file")), nil, l)
}

func (p *picker) updateRow(id widget.ListItemID, o fyne.CanvasObject) {
	if id >= len(p.shown) {
		return
	}
	e := p.shown[id]
	for _, c := range o.(*fyne.Container).Objects {
		switch w := c.(type) {
		case *widget.Label:
			w.SetText(e.rel)
		case *widget.Icon:
			if e.dir {
				w.SetResource(icon("folder"))
			} else {
				w.SetResource(icon("file"))
			}
		}
	}
}

// navigate shows dir and clears the search, which belonged to the folder left.
func (p *picker) navigate(dir string) {
	p.dir = filepath.Clean(dir)
	p.path.SetText(p.dir)
	if p.mode == pickFolder {
		p.setChosen(p.dir)
	} else {
		p.setChosen("")
	}
	p.search.SetText("")
	p.reload()
}

func (p *picker) reload() {
	all, err := listDir(p.dir, p.opt)
	p.all = all
	p.refilter()
	if err != nil {
		p.status.SetText("Cannot read: " + errText(err))
	}
}

// goTo opens a typed path: a folder, or in file mode a file, which it also selects.
func (p *picker) goTo(text string) {
	path := filepath.Clean(strings.TrimSpace(text))
	fi, err := os.Stat(path)
	switch {
	case !filepath.IsAbs(path):
		p.status.SetText("Type an absolute path")
	case err != nil:
		p.status.SetText(errText(err))
	case fi.IsDir():
		p.navigate(path)
	case p.mode == pickFile:
		p.navigate(filepath.Dir(path))
		p.setChosen(path)
	default:
		p.status.SetText("Not a folder")
	}
}

func errText(err error) string {
	var pe *os.PathError
	if errors.As(err, &pe) {
		return pe.Err.Error()
	}
	return err.Error()
}

func (p *picker) refilter() {
	g := p.gen.Add(1)
	p.cursor = -1
	p.list.UnselectAll()
	query := p.search.Text
	if !p.deep || len(searchTerms(query)) == 0 {
		p.show(filterEntries(p.all, query), "")
		return
	}
	p.show(nil, "Searching subfolders…")
	dir, opt := p.dir, p.opt
	var res []pickEntry
	var truncated bool
	p.s.background(func() {
		res, truncated = searchTree(dir, query, opt, func() bool { return p.gen.Load() != g })
	}, func() {
		if p.gen.Load() != g {
			return
		}
		note := ""
		if truncated {
			note = "search stopped early: refine it"
		}
		p.show(res, note)
	})
}

func (p *picker) show(es []pickEntry, note string) {
	p.shown = es
	p.list.Refresh()
	p.list.ScrollToTop()
	switch {
	case note != "" && es == nil:
		p.status.SetText(note)
	case len(searchTerms(p.search.Text)) == 0:
		p.status.SetText(plural(len(es), "item"))
	case note != "":
		p.status.SetText(fmt.Sprintf("%s, %s", plural(len(es), "match"), note))
	default:
		p.status.SetText(plural(len(es), "match"))
	}
}

func plural(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	if strings.HasSuffix(word, "h") {
		return fmt.Sprintf("%d %ses", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// selectRow enters a tapped folder; a file (or a folder reached with the arrow keys) is highlighted.
func (p *picker) selectRow(id widget.ListItemID) {
	if id >= len(p.shown) {
		return
	}
	p.cursor = id
	e := p.shown[id]
	switch {
	case e.dir && !p.keying:
		p.navigate(e.path)
	case !e.dir:
		p.setChosen(e.path)
	case p.mode == pickFile:
		p.setChosen("")
	}
}

func (p *picker) setChosen(path string) {
	p.chosen = path
	if path == "" {
		p.selected.SetText("")
		p.open.Disable()
		return
	}
	p.selected.SetText(path)
	p.open.Enable()
}

// searchKey makes the search field drive the list: arrows move, Enter opens, Escape cancels.
func (p *picker) searchKey(k fyne.KeyName) bool {
	switch k {
	case fyne.KeyDown:
		p.move(1)
	case fyne.KeyUp:
		p.move(-1)
	case fyne.KeyReturn, fyne.KeyEnter:
		p.activate()
	case fyne.KeyEscape:
		p.close()
	default:
		return false
	}
	return true
}

func (p *picker) move(delta int) {
	if len(p.shown) == 0 {
		return
	}
	next := min(max(p.cursor+delta, 0), len(p.shown)-1)
	p.keying = true
	p.list.Select(next)
	p.keying = false
	p.list.ScrollTo(next)
}

// activate enters the highlighted (else first) folder or picks the file; with nothing to act on, a
// folder picker takes the folder shown.
func (p *picker) activate() {
	if len(p.shown) == 0 || (p.cursor < 0 && len(searchTerms(p.search.Text)) == 0) {
		if p.mode == pickFolder {
			p.confirm()
		}
		return
	}
	e := p.shown[max(p.cursor, 0)]
	if e.dir {
		p.navigate(e.path)
		return
	}
	p.setChosen(e.path)
	p.confirm()
}

func (p *picker) confirm() {
	path := p.chosen
	if path == "" {
		return
	}
	p.close()
	p.onPick(path)
}

func (p *picker) close() { p.d.Hide() }

// findShortcut keeps Ctrl+F working while a picker field has the focus (it would swallow it).
func (p *picker) findShortcut(sc fyne.Shortcut) bool {
	if !isFindShortcut(sc) {
		return false
	}
	p.focusSearch()
	return true
}

func (p *picker) focusSearch() {
	p.s.win.Canvas().Focus(p.search)
	p.search.TypedShortcut(&fyne.ShortcutSelectAll{})
}

func isFindShortcut(sc fyne.Shortcut) bool {
	c, ok := sc.(*desktop.CustomShortcut)
	return ok && c.KeyName == fyne.KeyF && c.Modifier == fyne.KeyModifierShortcutDefault
}

// keyEntry is an Entry whose keys and shortcuts the picker sees first.
type keyEntry struct {
	widget.Entry
	onKey      func(fyne.KeyName) bool
	onShortcut func(fyne.Shortcut) bool
}

func newKeyEntry() *keyEntry {
	e := &keyEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *keyEntry) TypedKey(k *fyne.KeyEvent) {
	if e.onKey != nil && e.onKey(k.Name) {
		return
	}
	e.Entry.TypedKey(k)
}

func (e *keyEntry) TypedShortcut(sc fyne.Shortcut) {
	if e.onShortcut != nil && e.onShortcut(sc) {
		return
	}
	e.Entry.TypedShortcut(sc)
}
