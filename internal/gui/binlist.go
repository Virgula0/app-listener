//go:build gui

package gui

import (
	"encoding/hex"
	"fmt"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/workbench"
)

// binaryList edits a list of executables; extras are shown read-only after them.
type binaryList struct {
	shell    *shell
	entries  []ebpf.BinaryEntry
	extras   []workbench.NetExtra
	list     *widget.List
	add      *widget.Button
	onChange func()
	editable bool
}

func newBinaryList(s *shell, onChange func()) *binaryList {
	b := &binaryList{shell: s, onChange: onChange, editable: true}
	b.list = widget.NewList(
		func() int { return len(b.entries) + len(b.extras) },
		func() fyne.CanvasObject { return newBinaryRow() },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*binaryRow)
			switch {
			case id < len(b.entries):
				var remove func()
				if b.editable {
					remove = func() { b.remove(id) }
				}
				row.set(b.entries[id], "", remove)
			case id-len(b.entries) < len(b.extras):
				x := b.extras[id-len(b.entries)]
				row.set(x.Entry, x.Why, nil)
			}
		})
	b.add = widget.NewButtonWithIcon("Add executable…", icon("plus"), b.pick)
	return b
}

// card is the list in a panel under a title row with the Add button.
func (b *binaryList) card(title string) fyne.CanvasObject {
	return newPanel(container.NewBorder(hfill(boldLabel(title), b.add), nil, nil, nil, b.list))
}

// set replaces the entries without calling onChange.
func (b *binaryList) set(entries []ebpf.BinaryEntry, extras []workbench.NetExtra) {
	b.entries, b.extras = slices.Clone(entries), slices.Clone(extras)
	b.list.Refresh()
}

func (b *binaryList) setEditable(on bool) {
	b.editable = on
	if on {
		b.add.Enable()
	} else {
		b.add.Disable()
	}
	b.list.Refresh()
}

// put adds e, replacing an entry with the same path.
func (b *binaryList) put(e ebpf.BinaryEntry) {
	b.entries = slices.DeleteFunc(slices.Clone(b.entries), func(x ebpf.BinaryEntry) bool { return x.Path == e.Path })
	b.entries = append(b.entries, e)
	b.changed()
}

func (b *binaryList) remove(i int) {
	if i < len(b.entries) {
		b.entries = slices.Delete(slices.Clone(b.entries), i, i+1)
		b.changed()
	}
}

func (b *binaryList) changed() {
	b.list.Refresh()
	if b.onChange != nil {
		b.onChange()
	}
}

func (b *binaryList) pick() {
	d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		path := r.URI().Path()
		_ = r.Close()
		var entry ebpf.BinaryEntry
		b.shell.busy("Hashing "+path+"…", func() error {
			var hashErr error
			entry, hashErr = ebpf.ComputeBinaryEntry(path)
			return hashErr
		}, func() { b.put(entry) })
	}, b.shell.win)
	b.shell.startAt(d, "/usr/bin")
	d.Show()
}

// binaryRow is one executable of a list; a caption replaces the remove button on read-only rows.
type binaryRow struct {
	widget.BaseWidget
	path    *widget.Label
	detail  *widget.Label
	caption *pill
	remove  *widget.Button
}

func newBinaryRow() *binaryRow {
	r := &binaryRow{path: boldLabel(""), detail: mutedText(""), caption: newPill("", colorInfo)}
	r.path.Truncation = fyne.TextTruncateEllipsis
	r.detail.Wrapping = fyne.TextWrapOff
	r.detail.SizeName = theme.SizeNameCaptionText
	r.remove = widget.NewButtonWithIcon("", icon("x"), nil)
	r.remove.Importance = widget.LowImportance
	r.ExtendBaseWidget(r)
	return r
}

func (r *binaryRow) set(e ebpf.BinaryEntry, caption string, onRemove func()) {
	r.path.SetText(e.Path)
	r.detail.SetText(fmt.Sprintf("%s · sha256 %s…", e.Comm, hex.EncodeToString(e.Hash[:6])))
	r.remove.OnTapped = onRemove
	if onRemove == nil {
		r.remove.Hide()
	} else {
		r.remove.Show()
	}
	if caption == "" {
		r.caption.Hide()
	} else {
		r.caption.Set(caption, colorInfo)
		r.caption.Show()
	}
	r.Refresh()
}

func (r *binaryRow) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil,
		container.NewCenter(newIconTile("terminal-window", colorMuted, 30)),
		container.NewHBox(container.NewCenter(r.caption), r.remove),
		container.NewVBox(r.path, r.detail)))
}

// verifyExe hashes the executable a process ran (workbench.EntryFor) and passes the entry to next.
func (s *shell) verifyExe(exe workbench.Exe, next func(ebpf.BinaryEntry)) {
	var entry ebpf.BinaryEntry
	s.busy("Verifying "+exeText(exe)+"…", func() error {
		var err error
		entry, err = workbench.EntryFor(exe)
		return err
	}, func() { next(entry) })
}

// confirmExecutable shows what is about to be trusted or blocked before ok runs.
func (s *shell) confirmExecutable(verb string, entry ebpf.BinaryEntry, exe workbench.Exe, comm, note string, ok func()) {
	form := widget.NewForm(
		widget.NewFormItem("Process", widget.NewLabel(comm)),
		widget.NewFormItem("Executable", selectable(entry.Path)),
		widget.NewFormItem("Inode", widget.NewLabel(fmt.Sprintf("%d (device %d)", exe.Key.Ino, exe.Key.Dev))),
		widget.NewFormItem("SHA-256", selectable(hex.EncodeToString(entry.Hash[:]))),
	)
	body := mutedText(note)
	dialog.ShowCustomConfirm(verb+" executable", verb, "Cancel", container.NewVBox(form, body), func(yes bool) {
		if yes {
			ok()
		}
	}, s.win)
}

func selectable(s string) *widget.Label {
	l := widget.NewLabel(s)
	l.Selectable = true
	l.Wrapping = fyne.TextWrapBreak
	return l
}
