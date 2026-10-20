//go:build gui

package gui

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Virgula0/app-listener/internal/editbackend"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/recorder"
)

// recordForm is the setup of `edit-protected config record`: the paths to learn.
type recordForm struct {
	p     *protectedPage
	paths []string
	list  *fyne.Container
	entry *widget.Entry
}

func (p *protectedPage) startRecord() {
	if p.hello != nil && !p.hello.Live {
		showInformation("Record access", "Recording needs the daemon running with an edit-protected password "+
			"configured (live mode): the new sections are added to its daemon.conf.", p.s.win)
		return
	}
	f := &recordForm{p: p, list: container.NewVBox()}
	f.entry = widget.NewEntry()
	f.entry.SetPlaceHolder("/home/alice/.config/app")
	f.entry.OnSubmitted = func(string) { f.addTyped() }
	f.refreshList()
	add := newButton("Add", icon("plus"), f.addTyped)
	folder := newButton("Folder…", icon("folder"), func() {
		p.s.pickPath(pickFolder, "Directory to record", "", f.add)
	})
	file := newButton("File…", icon("file"), func() {
		p.s.pickPath(pickFile, "File to record", "", f.add)
	})
	start := newButton("Authenticate and record", icon("circle-fill"), f.submit)
	start.Importance = widget.HighImportance
	back := newButton("Back", nil, func() { p.show(p.over) })
	form := container.NewVBox(
		sectionTitle("PATHS TO GUARD"),
		newPanel(container.NewVBox(f.list,
			container.NewBorder(nil, nil, nil, container.NewHBox(add, folder, file), f.entry))),
		sectionTitle("HOW IT WORKS"),
		newPanel(mutedText("Every access to these paths is recorded, by the executable the kernel reports, "+
			"until you stop it. Use the applications that need them meanwhile. You then pick the binaries to "+
			"whitelist and the libraries to trust, can encrypt the paths, and the daemon adds one section "+
			"per path and reloads. A path daemon.conf already guards is extended instead: its whitelist is "+
			"kept and shown, only new lines are added. Anything that touched the paths is listed, malware included: approve only "+
			"what you ran on purpose.")),
		container.NewHBox(back, start),
	)
	p.show(newFadeVScroll(form))
}

func (f *recordForm) addTyped() {
	if t := strings.TrimSpace(f.entry.Text); t != "" {
		f.add(t)
		f.entry.SetText("")
	}
}

func (f *recordForm) add(path string) {
	if !filepath.IsAbs(path) {
		showInformation("Record access", path+" is not an absolute path.", f.p.s.win)
		return
	}
	path = filepath.Clean(path)
	if !slices.Contains(f.paths, path) {
		f.paths = append(f.paths, path)
	}
	f.refreshList()
}

func (f *recordForm) refreshList() {
	fillRemovableList(f.list, &f.paths, "folder", "No paths yet: type one or browse for a folder or file.",
		f.refreshList)
}

// submit sends the paths only after AUTH: the backend validates and watches them as the daemon's
// binary.
func (f *recordForm) submit() {
	p := f.p
	if len(f.paths) == 0 {
		showInformation("Record access", "Add at least one path.", p.s.win)
		return
	}
	paths := append([]string(nil), f.paths...)
	p.begin(func(link backendLink) {
		p.authenticate(link, "Record access", func([]string) {
			var m editbackend.Message
			var err error
			p.s.background(func() {
				m, err = link.Call(&editbackend.Request{Op: editbackend.OpRecordStart, Resources: paths})
			}, func() {
				if err != nil {
					p.s.showError("Record access", err)
					p.finish()
					return
				}
				v := newRecordView(p, link, m.Resources, m.RecordTargets)
				p.leave = v.discard
				p.show(v.root)
				v.start()
			})
		})
	})
}

// recordView shows a running recording, then its review.
type recordView struct {
	p       *protectedPage
	link    backendLink
	paths   []string
	prog    recorder.Progress
	elapsed *widget.Label
	count   *widget.Label
	table   *dataTable
	root    fyne.CanvasObject
	stopped chan struct{}
	done    bool
}

func newRecordView(p *protectedPage, link backendLink, paths []string, targets []recorder.Target) *recordView {
	v := &recordView{p: p, link: link, paths: paths, elapsed: boldLabel("00:00"), count: mutedText("0 accesses"),
		stopped: make(chan struct{})}
	stop := newButton("Stop recording", icon("stop"), v.stop)
	stop.Importance = widget.DangerImportance
	head := container.NewVBox(
		container.NewBorder(nil, nil, container.NewHBox(newPill("REC", theme.ColorNameError), v.elapsed), stop,
			nil),
		selectable(strings.Join(paths, "\n")),
		mutedText("Use the applications that need these paths now. Everything that touches them is listed for "+
			"your review when you stop."),
	)
	for _, o := range existingNotices(targets) {
		head.Add(o)
	}
	v.table = newDataTable([]column{{"Last", 90}, {"Accesses", 90}, {"Events", 220}, {"Executable", 520}},
		func() int { return len(v.prog.Recent) }, v.cell, nil)
	v.root = container.NewBorder(newPanel(head), v.count, nil, nil, v.table)
	return v
}

// existingNotices warn about each target daemon.conf already guards and show its whitelist.
func existingNotices(targets []recorder.Target) []fyne.CanvasObject {
	var out []fyne.CanvasObject
	for _, t := range targets {
		if t.Existing == nil {
			continue
		}
		w := warnText()
		setWarn(w, logging.SanitizeText(t.Existing.Warning(t.Path)))
		allow := "Already whitelisted: nothing (every access is denied)."
		if lines := t.Existing.Lines(); len(lines) > 0 {
			allow = "Already whitelisted:\n" + logging.SanitizeText(strings.Join(lines, "\n"))
		}
		out = append(out, w, selectable(allow))
	}
	return out
}

func (v *recordView) cell(row, col int) cell {
	s := v.prog.Recent[row]
	switch col {
	case 0:
		return cell{text: formatAgo(time.Since(s.Last))}
	case 1:
		return cell{text: strconv.FormatUint(s.Count, 10)}
	case 2:
		return cell{text: strings.Join(s.Ops, ",")}
	default:
		return cell{text: logging.SanitizeText(s.Exe)}
	}
}

func formatAgo(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm ago", int(d.Minutes()))
}

// start polls the backend once a second until the recording stops or the backend exits.
func (v *recordView) start() {
	if v.p.s.syncBackground {
		return
	}
	link := v.link
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-v.stopped:
				return
			case <-link.Done():
				fyne.Do(func() {
					if v.p.link == link && !v.done {
						v.p.s.showError("Record access", editbackend.ErrClosed)
						v.p.finish()
					}
				})
				return
			case <-tick.C:
				m, err := link.Call(&editbackend.Request{Op: editbackend.OpRecordStatus})
				if err != nil || m.RecordProgress == nil {
					continue
				}
				prog := *m.RecordProgress
				fyne.Do(func() { v.show(prog) })
			}
		}
	}()
}

func (v *recordView) show(prog recorder.Progress) {
	if v.done {
		return
	}
	v.prog = prog
	v.elapsed.SetText(formatClock(prog.Elapsed))
	v.count.SetText(fmt.Sprintf("%d accesses by %d executables", prog.Events, prog.Binaries))
	v.table.Refresh()
	if prog.Elapsed >= recorder.MaxDuration {
		v.stop()
	}
}

func formatClock(d time.Duration) string {
	d = d.Truncate(time.Second)
	if h := int(d.Hours()); h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, int(d.Minutes())%60, int(d.Seconds())%60)
	}
	return fmt.Sprintf("%02d:%02d", int(d.Minutes()), int(d.Seconds())%60)
}

func (v *recordView) stop() {
	if v.done {
		return
	}
	v.done = true
	close(v.stopped)
	link := v.link
	var m editbackend.Message
	var err error
	v.p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: editbackend.OpRecordStop}) }, func() {
		if err != nil || m.Recording == nil {
			if err == nil {
				err = errors.New("the backend returned no recording")
			}
			v.p.s.showError("Record access", err)
			v.p.finish()
			return
		}
		r := newRecordReview(v.p, link, m.Recording, m.EncryptBlockers)
		v.p.leave = r.discard
		v.p.show(r.root)
	})
}

func (v *recordView) discard() {
	if !v.done {
		v.done = true
		close(v.stopped)
	}
	v.p.finish()
}

// recordReview is the approval of one recording.
type recordReview struct {
	p       *protectedPage
	link    backendLink
	rec     *recorder.Recording
	bins    map[string]*handCheck
	libs    map[string]*handCheck
	encrypt map[string]*handCheck
	masks   *handCheck
	status  *widget.Label
	root    fyne.CanvasObject
	applied bool
}

func newRecordReview(p *protectedPage, link backendLink, rec *recorder.Recording, blockers map[string]string) *recordReview {
	r := &recordReview{p: p, link: link, rec: rec, bins: map[string]*handCheck{}, libs: map[string]*handCheck{},
		encrypt: map[string]*handCheck{}, status: mutedText("")}
	binBox, refused, covered, libBox := r.binaryChecks()
	if len(r.bins) == 0 {
		binBox.Add(mutedText("No binary to add: applying guards the new paths with an empty whitelist " +
			"(every access denied); the existing sections stay as they are."))
	}
	encBox := r.encryptChecks(blockers)
	r.masks = newCheck("Restrict each binary to the events it was seen doing", nil)
	review := newButton("Review changes", icon("check"), r.review)
	review.Importance = widget.HighImportance
	discard := newButton("Discard", icon("x"), r.discard)

	parts := []fyne.CanvasObject{
		newPanel(container.NewVBox(boldLabel(fmt.Sprintf("Recorded %s, %d accesses",
			formatClock(rec.Stopped.Sub(rec.Started)), rec.Events)),
			mutedText("Only root-placed binaries that aren't general tools start checked. Check only what you ran "+
				"on purpose."))),
	}
	if notices := existingNotices(rec.Targets); len(notices) > 0 {
		parts = append(parts, sectionTitle("ALREADY IN DAEMON.CONF"), newPanel(container.NewVBox(notices...)))
	}
	parts = append(parts, sectionTitle("BINARIES TO WHITELIST"), newPanel(binBox))
	if len(covered.Objects) > 0 {
		parts = append(parts, sectionTitle("ALREADY WHITELISTED"), newPanel(covered))
	}
	if len(refused.Objects) > 0 {
		parts = append(parts, sectionTitle("NOT WHITELISTABLE"), newPanel(refused))
	}
	if len(r.libs) > 0 {
		parts = append(parts, sectionTitle("LIBRARIES TO TRUST (allow_lib)"), newPanel(container.NewVBox(
			mutedText("Trust is daemon-wide and write-protects the file: approve only libraries the application "+
				"itself ships. None start checked."), libBox)))
	}
	parts = append(parts, sectionTitle("OPTIONS"), newPanel(container.NewVBox(r.masks, mutedText("Least "+
		"privilege, but an operation the recording missed is then denied."))))
	if len(encBox.Objects) > 0 {
		parts = append(parts, sectionTitle("ENCRYPTION"), newPanel(encBox))
	}
	parts = append(parts, container.NewHBox(discard, review), r.status)
	r.root = newFadeVScroll(newGapVBox(parts...))
	return r
}

func (r *recordReview) binaryChecks() (binBox, refused, covered, libBox *fyne.Container) {
	binBox, refused, covered, libBox = container.NewVBox(), container.NewVBox(), container.NewVBox(), container.NewVBox()
	for i := range r.rec.Binaries {
		b := &r.rec.Binaries[i]
		if !b.Admissible() {
			refused.Add(binaryCard(b, nil))
			continue
		}
		if b.Covered() {
			covered.Add(binaryCard(b, nil))
			continue
		}
		c := newCheck(logging.SanitizeText(b.Line), nil)
		c.Checked = b.Suggested
		r.bins[b.Line] = c
		binBox.Add(binaryCard(b, c))
		for _, l := range b.Libs {
			if l.Problem != "" || r.libs[l.Path] != nil {
				continue
			}
			lc := newCheck(logging.SanitizeText(l.Path), nil)
			r.libs[l.Path] = lc
			libBox.Add(container.NewVBox(lc, mutedText("loaded by "+logging.SanitizeText(b.Line))))
		}
	}
	return binBox, refused, covered, libBox
}

func (r *recordReview) encryptChecks(blockers map[string]string) *fyne.Container {
	box := container.NewVBox()
	for _, t := range r.rec.Targets {
		if t.Existing != nil {
			continue
		}
		c := newCheck("Encrypt "+t.Path+" at rest", nil)
		if why := blockers[t.Path]; why != "" {
			c.Disable()
			box.Add(container.NewVBox(c, mutedText("Guarded unencrypted: "+why)))
		} else {
			c.Checked = true
			box.Add(container.NewVBox(c, mutedText("Moved to "+t.Path+fscrypt.BackupSuffix+", re-created "+
				"encrypted and locked; the reload attaches its guard, then unlocks it.")))
		}
		r.encrypt[t.Path] = c
	}
	return box
}

// binaryCard describes one recorded executable; check is nil for one that can't be added.
func binaryCard(b *recorder.Binary, check *handCheck) fyne.CanvasObject {
	box := container.NewVBox()
	if check != nil {
		box.Add(check)
	} else {
		name := b.Line
		if name == "" {
			name = b.Exe
		}
		if name == "" {
			name = "(unnamed executable)"
		}
		box.Add(boldLabel(logging.SanitizeText(name)))
		if b.Problem != "" {
			box.Add(mutedText(logging.SanitizeText(b.Problem)))
		}
	}
	box.Add(mutedText(fmt.Sprintf("%d accesses (%s) by %s, uid %v, on %s", b.Count, strings.Join(b.Ops, ","),
		logging.SanitizeText(strings.Join(b.Comms, ",")), b.UIDs, logging.SanitizeText(strings.Join(b.Targets, ", ")))))
	if b.SHA256 != "" {
		box.Add(mutedText("sha256 " + b.SHA256))
	}
	for _, n := range b.Notes {
		w := warnText()
		setWarn(w, n)
		box.Add(w)
	}
	return newPanel(box)
}

func (r *recordReview) selection() recorder.Selection {
	sel := recorder.Selection{Masks: r.masks.Checked}
	for line, c := range r.bins {
		if c.Checked {
			sel.Binaries = append(sel.Binaries, line)
		}
	}
	for path, c := range r.libs {
		if c.Checked {
			sel.Libs = append(sel.Libs, path)
		}
	}
	for path, c := range r.encrypt {
		if c.Checked && !c.Disabled() {
			sel.Encrypt = append(sel.Encrypt, path)
		}
	}
	slices.Sort(sel.Binaries)
	slices.Sort(sel.Libs)
	slices.Sort(sel.Encrypt)
	return sel
}

func (r *recordReview) review() {
	sel := r.selection()
	link := r.link
	var m editbackend.Message
	var err error
	r.p.s.background(func() {
		m, err = link.Call(&editbackend.Request{Op: editbackend.OpRecordPlan, Selection: &sel})
	}, func() {
		if err != nil {
			r.p.s.showError("Not applicable", err)
			return
		}
		body := container.NewBorder(promptBody("Apply and reload the daemon?",
			"If anything fails, the encrypted paths are restored from their backups and daemon.conf keeps its "+
				"current content."), nil, nil, nil, newFadeScroll(diffText(inst.UnifiedDiff(string(m.Base), string(m.Data)))))
		d := dialog.NewCustomConfirm("daemon.conf — changes", "Apply", "Back", body, func(ok bool) {
			if ok {
				r.apply()
			}
		}, r.p.s.win)
		r.p.s.showSized(d, fyne.NewSize(1040, 720))
	})
}

func (r *recordReview) apply() {
	link := r.link
	r.status.SetText("applying — the daemon rebuilds every guard …")
	stopProgress := r.followProgress()
	var m editbackend.Message
	var err error
	r.p.s.background(func() { m, err = link.Call(&editbackend.Request{Op: editbackend.OpRecordApply}) }, func() {
		stopProgress()
		if err != nil {
			r.status.SetText("not applied: " + err.Error())
			r.p.s.showError("Not applied", err)
			return
		}
		r.applied = true
		r.status.SetText("daemon.conf updated, the daemon guards the paths")
		msg := "daemon.conf is updated with the recorded paths and the daemon guards them."
		if len(m.Notes) > 0 {
			msg += "\n\n" + strings.Join(m.Notes, "\n")
		}
		d := dialog.NewCustomWithoutButtons("Record access", scrollable(mutedText(msg)), r.p.s.win)
		d.SetButtons([]fyne.CanvasObject{newTextButton("OK", func() {
			d.Hide()
			r.offerBackups(m.Backups)
		})})
		r.p.s.showSized(d, fyne.NewSize(640, 320))
	})
}

// followProgress shows record_apply's progress until the returned stop is called.
func (r *recordReview) followProgress() func() {
	if r.p.s.syncBackground {
		return func() {}
	}
	stop := make(chan struct{})
	link := r.link
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-link.Done():
				return
			case m := <-link.Events():
				if m.Progress == nil {
					continue
				}
				pr := *m.Progress
				fyne.Do(func() {
					r.status.SetText(fmt.Sprintf("%s — %d%%", logging.SanitizeText(pr.Step), int(pr.Fraction*100)))
				})
			}
		}
	}()
	return func() { close(stop) }
}

// offerBackups asks, one backup at a time, whether to delete it; then the workflow ends.
func (r *recordReview) offerBackups(backups []string) {
	if len(backups) == 0 {
		r.p.finish()
		return
	}
	b, rest := backups[0], backups[1:]
	showCustomConfirm("Remove the unencrypted backup?", "Remove", "Keep", mutedText(b+" is a plaintext copy of "+
		"what is now encrypted, outside every guarded tree: anything running as its owner can read it."),
		func(remove bool) {
			if !remove {
				r.offerBackups(rest)
				return
			}
			link := r.link
			var err error
			r.p.s.background(func() {
				_, err = link.Call(&editbackend.Request{Op: editbackend.OpRecordDeleteBackup, Path: b})
			}, func() {
				if err != nil {
					r.p.s.showError("Removing "+b, err)
				}
				r.offerBackups(rest)
			})
		}, r.p.s.win)
}

func (r *recordReview) discard() {
	if r.applied {
		r.p.finish()
		return
	}
	showConfirm("Discard the recording?", "Nothing is added to daemon.conf.", func(ok bool) {
		if ok {
			r.p.finish()
		}
	}, r.p.s.win)
}
