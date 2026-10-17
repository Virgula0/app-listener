//go:build gui

package gui

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// startEdit is `edit-protected` without flags: live when the daemon has a password, else offline.
func (p *protectedPage) startEdit() {
	p.begin(func(link backendLink) {
		if link.Hello().Live {
			p.authenticate(link, "Edit protected files", func(resources []string) {
				if len(resources) == 0 {
					p.s.showError("Edit protected files", fmt.Errorf("the daemon reports no guarded directories"))
					p.finish()
					return
				}
				p.pickOne("Edit protected files", "Which guarded directory do you want to edit?", resources,
					func(r string) { p.openEdit(link, editbackend.OpSelect, r, true) })
			})
			return
		}
		var m editbackend.Message
		var err error
		p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: editbackend.OpOfflineScan}) }, func() {
			if err != nil {
				p.s.showError("Edit protected files (offline)", err)
				p.finish()
				return
			}
			p.pickOne("Edit protected files (offline)", "Which protected directory do you want to open? It is "+
				"unlocked until you end the session, then locked again.", m.Resources,
				func(r string) { p.openEdit(link, editbackend.OpOfflineOpen, r, false) })
		})
	})
}

func (p *protectedPage) openEdit(link backendLink, op, resource string, live bool) {
	var m editbackend.Message
	var err error
	p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: op, Resource: resource}) }, func() {
		if err != nil {
			p.s.showError("Opening "+resource, err)
			p.finish()
			return
		}
		v := newEditView(p, link, resource, live, m.SingleFile)
		p.leave = v.end
		p.show(v.root)
		v.start()
	})
}

// kindFile and kindFolder name entries in messages and are also their icons' names.
const kindFile, kindFolder = "file", "folder"

// codeEntry is the editor: Ctrl+S saves, every key is session activity.
type codeEntry struct {
	widget.Entry
	onSave     func()
	onActivity func()
}

func newCodeEntry() *codeEntry {
	e := &codeEntry{}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapOff
	e.TextStyle = fyne.TextStyle{Monospace: true}
	e.ExtendBaseWidget(e)
	return e
}

func (e *codeEntry) TypedShortcut(sc fyne.Shortcut) {
	if cs, ok := sc.(*desktop.CustomShortcut); ok && cs.KeyName == fyne.KeyS &&
		cs.Modifier == fyne.KeyModifierShortcutDefault && e.onSave != nil {
		e.onSave()
		return
	}
	e.Entry.TypedShortcut(sc)
}

func (e *codeEntry) TypedKey(k *fyne.KeyEvent) {
	if e.onActivity != nil {
		e.onActivity()
	}
	e.Entry.TypedKey(k)
}

func (e *codeEntry) TypedRune(r rune) {
	if e.onActivity != nil {
		e.onActivity()
	}
	e.Entry.TypedRune(r)
}

// editView is one open resource: tree left, editor right; every access is a backend call.
type editView struct {
	p        *protectedPage
	link     backendLink
	resource string
	live     bool
	single   bool

	children map[string][]editbackend.Entry
	entries  map[string]editbackend.Entry
	selected string

	editPath string
	original string
	dirty    bool
	loading  bool
	ended    bool
	lastPing time.Time

	root    fyne.CanvasObject
	tree    *widget.Tree
	editor  *codeEntry
	title   *widget.Label
	info    *widget.Label
	status  *widget.Label
	banner  *fyne.Container
	buttons []*handButton
	save    *handButton
}

func newEditView(p *protectedPage, link backendLink, resource string, live, single bool) *editView {
	v := &editView{p: p, link: link, resource: resource, live: live, single: single,
		children: map[string][]editbackend.Entry{}, entries: map[string]editbackend.Entry{}}
	v.build()
	return v
}

// toolbar is the editing actions and the end-session button.
func (v *editView) toolbar() (bar *fyne.Container, end *handButton) {
	v.save = newButton("Save", icon("check"), v.saveFile)
	v.save.Importance = widget.HighImportance
	v.save.Disable()
	newFile := newButton("New file", icon(kindFile), func() { v.create(false) })
	newDir := newButton("New folder", icon(kindFolder), func() { v.create(true) })
	perms := newButton("Permissions", icon("sliders-horizontal"), v.chmod)
	owner := newButton("Owner", icon("user"), v.chown)
	del := newButton("Delete", icon("trash"), v.remove)
	del.Importance = widget.DangerImportance
	end = newButton("End session", icon("lock-simple"), v.end)
	v.buttons = []*handButton{newFile, newDir, perms, owner, del}
	if v.single {
		newFile.Hide()
		newDir.Hide()
		del.Hide()
	}
	return container.NewHBox(v.save, widget.NewSeparator(), newFile, newDir, perms, owner, del), end
}

func (v *editView) build() {
	v.editor = newCodeEntry()
	v.editor.onSave = v.saveFile
	v.editor.onActivity = v.ping
	v.editor.OnChanged = func(string) { v.setDirty(v.editPath != "" && v.editor.Text != v.original) }
	v.editor.Disable()
	v.title = boldLabel("")
	v.title.Truncation = fyne.TextTruncateEllipsis
	v.info = mutedText("")
	v.status = mutedText("")
	v.status.Truncation = fyne.TextTruncateEllipsis
	v.banner = container.NewVBox()

	mode, modeColor, note := "LIVE", theme.ColorNameSuccess,
		"the daemon keeps guarding it; write access is granted only for this session and revoked after "+
			v.p.timeout.Selected+" without input"
	if !v.live {
		mode, modeColor, note = "OFFLINE", theme.ColorNameWarning,
			"it is unlocked and accessible while this session is open, and locked again when you end it"
	}
	toolbar, end := v.toolbar()
	head := newPanel(container.NewVBox(
		container.NewBorder(nil, nil, newPill(mode, modeColor), end, boldLabel(v.resource)),
		mutedText("Editing "+v.resource+": "+note+"."),
	))

	right := container.NewBorder(container.NewVBox(v.title, v.info), nil, nil, nil, v.editor)
	var body fyne.CanvasObject = right
	if !v.single {
		v.tree = widget.NewTree(v.childUIDs, v.isBranch, v.createNode, v.updateNode)
		v.tree.OnSelected = v.onSelect
		v.tree.OnBranchOpened = func(uid string) {
			if _, ok := v.children[relOf(uid)]; !ok {
				v.load(relOf(uid), nil)
			}
		}
		split := container.NewHSplit(newPanel(v.tree), right)
		split.Offset = 0.3
		body = split
	}
	v.root = container.NewBorder(container.NewVBox(head, v.banner, toolbar), v.status, nil, nil, body)
}

// start loads the tree (or the single file) and watches for the daemon ending the session.
func (v *editView) start() {
	v.p.watchEnded(v.link, v.onEnded)
	if v.single {
		v.selected = "."
		v.openFile(".")
		v.load(".", nil)
		return
	}
	v.load(".", nil)
}

func relOf(uid string) string {
	if uid == "" {
		return "."
	}
	return uid
}

func joinRel(dir, name string) string {
	if dir == "." {
		return name
	}
	return dir + "/" + name
}

func (v *editView) childUIDs(uid string) []string {
	dir := relOf(uid)
	out := make([]string, 0, len(v.children[dir]))
	for _, e := range v.children[dir] {
		out = append(out, joinRel(dir, e.Name))
	}
	return out
}

func (v *editView) isBranch(uid string) bool {
	if uid == "" {
		return true
	}
	e, ok := v.entries[uid]
	return ok && e.Dir && !e.Symlink && !e.Invalid
}

func (v *editView) createNode(bool) fyne.CanvasObject {
	l := widget.NewLabel("")
	l.Truncation = fyne.TextTruncateEllipsis
	return container.NewBorder(nil, nil, widget.NewIcon(nil), nil, l)
}

func (v *editView) updateNode(uid string, branch bool, o fyne.CanvasObject) {
	box := o.(*fyne.Container)
	l, ic := box.Objects[0].(*widget.Label), box.Objects[1].(*widget.Icon)
	e := v.entries[uid]
	name := path.Base(uid)
	switch {
	case e.Symlink:
		name += "  →  (symlink)"
	case e.Invalid:
		name += "  (name is not UTF-8)"
	}
	l.SetText(name)
	l.Importance = widget.MediumImportance
	if e.Symlink || e.Invalid {
		l.Importance = widget.LowImportance
	}
	if branch {
		ic.SetResource(icon(kindFolder))
	} else {
		ic.SetResource(icon(kindFile))
	}
}

// load lists dir; then (if set) runs on the UI goroutine.
func (v *editView) load(dir string, then func()) {
	var m editbackend.Message
	var err error
	v.p.s.background(func() { m, err = v.link.Call(&editbackend.Request{Op: editbackend.OpList, Path: dir}) }, func() {
		if err != nil {
			v.setStatus(err.Error())
			return
		}
		if v.single {
			if len(m.Entries) == 1 {
				v.entries["."] = m.Entries[0]
				v.showInfo(".")
			}
			return
		}
		list := m.Entries
		sort.Slice(list, func(i, j int) bool {
			if list[i].Dir != list[j].Dir {
				return list[i].Dir
			}
			return list[i].Name < list[j].Name
		})
		v.children[dir] = list
		for _, e := range list {
			v.entries[joinRel(dir, e.Name)] = e
		}
		v.tree.Refresh()
		if then != nil {
			then()
		}
	})
}

// reload re-lists dir, keeping it open in the tree.
func (v *editView) reload(dir string) {
	v.load(dir, func() {
		if dir != "." {
			v.tree.OpenBranch(dir)
		}
	})
}

func (v *editView) onSelect(uid string) {
	v.ping()
	v.selected = uid
	v.showInfo(uid)
	e := v.entries[uid]
	if e.Dir || e.Invalid || uid == v.editPath {
		return
	}
	if e.Symlink {
		v.setStatus(path.Join(v.resource, uid) + " is a symlink — editing it is refused")
		return
	}
	if v.dirty {
		dialog.ShowConfirm("Unsaved changes", "Discard unsaved changes to "+path.Base(v.editPath)+"?",
			func(ok bool) {
				if ok {
					v.openFile(uid)
				}
			}, v.p.s.win)
		return
	}
	v.openFile(uid)
}

func (v *editView) showInfo(uid string) {
	e, ok := v.entries[uid]
	if !ok {
		v.info.SetText("")
		return
	}
	kind := kindFile
	switch {
	case e.Symlink:
		kind = "symlink"
	case e.Dir:
		kind = "directory"
	}
	v.info.SetText(fmt.Sprintf("%s · %s · mode %04o (%s) · uid %d gid %d · modified %s", kind,
		fileSize(e.Size), e.Mode, e.ModeStr, e.UID, e.GID, e.ModTime.Format("2006-01-02 15:04")))
}

func fileSize(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.1f KiB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	}
}

func (v *editView) openFile(rel string) {
	if v.loading {
		return
	}
	v.loading = true
	var m editbackend.Message
	var err error
	v.p.s.background(func() { m, err = v.link.Call(&editbackend.Request{Op: editbackend.OpRead, Path: rel}) }, func() {
		v.loading = false
		full := path.Join(v.resource, rel)
		switch {
		case err != nil:
			v.setStatus(err.Error())
			return
		case m.File.TooLarge:
			v.setStatus(fmt.Sprintf("%s is %s — only files up to 2 MiB can be edited", full, fileSize(m.File.Size)))
			return
		case m.File.Binary:
			v.setStatus(full + " looks like a binary file — it cannot be edited")
			return
		}
		v.editPath, v.original = rel, string(m.File.Data)
		v.title.SetText(full)
		v.editor.SetText(v.original)
		v.setDirty(false)
		if !v.ended {
			v.editor.Enable()
		}
		v.setStatus("editing " + full)
	})
}

func (v *editView) setDirty(on bool) {
	v.dirty = on
	if on && !v.ended {
		v.save.Enable()
	} else {
		v.save.Disable()
	}
	if v.editPath != "" {
		mark := ""
		if on {
			mark = "  •  unsaved"
		}
		v.title.SetText(path.Join(v.resource, v.editPath) + mark)
	}
}

func (v *editView) setStatus(s string) { v.status.SetText(s) }

// ping is session activity (the daemon's idle timeout), at most every few seconds.
func (v *editView) ping() {
	if !v.live || v.ended || time.Since(v.lastPing) < 3*time.Second {
		return
	}
	v.lastPing = time.Now()
	link := v.link
	v.p.s.background(func() { _, _ = link.Call(&editbackend.Request{Op: editbackend.OpPing}) }, func() {})
}

func (v *editView) saveFile() {
	if v.editPath == "" || !v.dirty || v.ended {
		return
	}
	v.ping()
	rel, text := v.editPath, v.editor.Text
	var err error
	v.p.s.background(func() {
		_, err = v.link.Call(&editbackend.Request{Op: editbackend.OpWrite, Path: rel, Data: []byte(text)})
	}, func() {
		if err != nil {
			v.setStatus(err.Error())
			v.p.s.showError("Saving "+path.Join(v.resource, rel), err)
			return
		}
		if v.editPath == rel {
			v.original = text
			v.setDirty(v.editor.Text != text)
		}
		v.setStatus("saved " + path.Join(v.resource, rel))
		v.reload(path.Dir(rel))
	})
}

// targetDir is where a new entry goes: the selected directory, or the selected file's parent.
func (v *editView) targetDir() string {
	if v.selected == "" {
		return "."
	}
	if v.isBranch(v.selected) {
		return v.selected
	}
	return path.Dir(v.selected)
}

func (v *editView) create(isDir bool) {
	v.ping()
	dir := v.targetDir()
	kind := kindFile
	if isDir {
		kind = kindFolder
	}
	entry := widget.NewEntry()
	entry.SetPlaceHolder("name")
	d := dialog.NewForm("New "+kind, "Create", "Cancel",
		[]*widget.FormItem{widget.NewFormItem("In", selectable(path.Join(v.resource, dir))),
			widget.NewFormItem("Name", entry)},
		func(ok bool) {
			if !ok {
				return
			}
			name := strings.TrimSpace(entry.Text)
			var err error
			v.p.s.background(func() {
				_, err = v.link.Call(&editbackend.Request{Op: editbackend.OpCreate, Path: dir, Name: name, Dir: isDir})
			}, func() {
				if err != nil {
					v.setStatus(err.Error())
				} else {
					v.setStatus("created " + path.Join(v.resource, dir, name))
				}
				v.reload(dir)
			})
		}, v.p.s.win)
	entry.OnSubmitted = func(string) { d.Submit() }
	v.p.s.showSized(d, fyne.NewSize(560, 220))
	v.p.s.win.Canvas().Focus(entry)
}

// pickedEntry is the selection a chmod/chown/delete acts on, refusing symlinks like the TUI.
func (v *editView) pickedEntry(verb string) (string, editbackend.Entry, bool) {
	uid := v.selected
	if v.single {
		uid = "."
	}
	e, ok := v.entries[uid]
	switch {
	case uid == "" || !ok:
		v.setStatus("select an entry first")
		return "", e, false
	case e.Invalid:
		v.setStatus("this entry's name is not UTF-8: use the terminal editor")
		return "", e, false
	case e.Symlink:
		v.setStatus(path.Join(v.resource, uid) + " is a symlink — " + verb + " is refused")
		return "", e, false
	}
	return uid, e, true
}

func (v *editView) chmod() {
	v.ping()
	uid, e, ok := v.pickedEntry("chmod")
	if !ok {
		return
	}
	entry := widget.NewEntry()
	// The full mode, setuid/setgid/sticky included: confirming it untouched must not strip them.
	entry.SetText(fmt.Sprintf("%04o", e.Mode))
	d := dialog.NewForm("Permissions", "Apply", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Entry", selectable(path.Join(v.resource, uid))),
		widget.NewFormItem("Mode", entry),
		widget.NewFormItem("", mutedText("1-4 octal digits, e.g. 0600 or 4755 (setuid)")),
	}, func(ok bool) {
		if !ok {
			return
		}
		mode := strings.TrimSpace(entry.Text)
		v.apply(&editbackend.Request{Op: editbackend.OpChmod, Path: uid, Mode: mode}, "chmod "+uid+" = "+mode, uid)
	}, v.p.s.win)
	entry.OnSubmitted = func(string) { d.Submit() }
	v.p.s.showSized(d, fyne.NewSize(560, 260))
}

func (v *editView) chown() {
	v.ping()
	uid, e, ok := v.pickedEntry("chown")
	if !ok {
		return
	}
	var m editbackend.Message
	var err error
	v.p.s.background(func() { m, err = v.link.Call(&editbackend.Request{Op: editbackend.OpUsers}) }, func() {
		if err != nil {
			v.setStatus(err.Error())
			return
		}
		keep := fmt.Sprintf("keep current (uid %d, gid %d)", e.UID, e.GID)
		labels := []string{keep}
		for _, u := range m.Users {
			labels = append(labels, fmt.Sprintf("%s (uid %d, gid %d)", u.Name, u.UID, u.GID))
		}
		pick := widget.NewSelect(labels, nil)
		pick.Selected = keep
		dialog.ShowForm("Owner", "Apply", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Entry", selectable(path.Join(v.resource, uid))),
			widget.NewFormItem("Owner", pick),
		}, func(ok bool) {
			i := pick.SelectedIndex()
			if !ok || i <= 0 {
				return
			}
			u := m.Users[i-1]
			v.apply(&editbackend.Request{Op: editbackend.OpChown, Path: uid, UID: u.UID, GID: u.GID},
				"chown "+uid+" → "+u.Name, uid)
		}, v.p.s.win)
	})
}

func (v *editView) remove() {
	v.ping()
	uid, e, ok := v.pickedEntry("delete")
	if !ok {
		return
	}
	if uid == "." {
		v.setStatus("cannot delete the root of the vault")
		return
	}
	msg := "Delete " + path.Join(v.resource, uid) + "?"
	if e.Dir {
		msg += " (recursively)"
	}
	dialog.ShowConfirm("Delete", msg, func(ok bool) {
		if !ok {
			return
		}
		if v.editPath == uid || strings.HasPrefix(v.editPath, uid+"/") {
			v.closeFile()
		}
		v.selected = ""
		v.apply(&editbackend.Request{Op: editbackend.OpRemove, Path: uid}, "deleted "+uid, uid)
	}, v.p.s.win)
}

func (v *editView) closeFile() {
	v.editPath, v.original = "", ""
	v.editor.SetText("")
	v.editor.Disable()
	v.title.SetText("")
	v.setDirty(false)
}

// apply runs one change and re-lists the parent of uid.
func (v *editView) apply(req *editbackend.Request, done, uid string) {
	var err error
	v.p.s.background(func() { _, err = v.link.Call(req) }, func() {
		if err != nil {
			v.setStatus(err.Error())
		} else {
			v.setStatus(done)
		}
		if v.single {
			v.load(".", nil)
		} else {
			v.reload(path.Dir(uid))
		}
	})
}

// onEnded: the daemon revoked the grant; the editor keeps the text, nothing more is written.
func (v *editView) onEnded(why string) {
	if v.ended {
		return
	}
	v.ended = true
	v.editor.Disable()
	v.save.Disable()
	for _, b := range v.buttons {
		b.Disable()
	}
	v.banner.Objects = []fyne.CanvasObject{endedBanner(why)}
	v.banner.Refresh()
}

// end narrows the grant (or re-locks) and shows the post-edit audit, like the CLI on exit.
func (v *editView) end() {
	if v.dirty && !v.ended {
		dialog.ShowConfirm("Unsaved changes", "End the session and discard unsaved changes to "+
			path.Base(v.editPath)+"?", func(ok bool) {
			if ok {
				v.dirty = false
				v.end()
			}
		}, v.p.s.win)
		return
	}
	var m editbackend.Message
	var err error
	link := v.link
	v.p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: editbackend.OpEnd}) }, func() {
		v.p.finish()
		if err != nil {
			v.p.s.showError("Ending the session", err)
		}
		showAudit(v.p.s, m.Findings)
	})
}

// showAudit is the post-edit audit notice: advisory, the edit is already written.
func showAudit(s *shell, findings []string) {
	if len(findings) == 0 {
		return
	}
	list := container.NewVBox()
	for _, f := range findings {
		l := selectable("• " + f)
		l.Wrapping = fyne.TextWrapWord
		list.Add(l)
	}
	d := dialog.NewCustom("Post-edit audit — potential issues", "Close", container.NewBorder(
		promptBody("Advisory only", "The edit was saved and nothing here is changed or blocked."), nil, nil, nil,
		container.NewVScroll(list)), s.win)
	s.showSized(d, fyne.NewSize(860, 520))
}
