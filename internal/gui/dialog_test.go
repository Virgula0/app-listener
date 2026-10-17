//go:build gui

package gui

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

func pickerTree(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, d := range []string{"sub/deeper", ".hidden", "alpine"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"alpha.txt", "Beta.conf", "sub/alphabet", "sub/deeper/ALPHA.key", ".hidden/alpha.secret"} {
		if err := os.WriteFile(filepath.Join(dir, f), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func openPicker(t *testing.T, mode pickMode, dir string) (*shell, *picker, *string) {
	t.Helper()
	a := test.NewApp()
	t.Cleanup(a.Quit)
	w := a.NewWindow("x")
	w.Resize(fyne.NewSize(1000, 700))
	s := newShell(a, w, Options{})
	s.syncBackground = true
	got := new(string)
	s.pickPath(mode, "Pick", dir, func(p string) { *got = p })
	if s.picker == nil {
		t.Fatal("picker not open")
	}
	return s, s.picker, got
}

func shownNames(p *picker) []string {
	var out []string
	for _, e := range p.shown {
		out = append(out, e.rel)
	}
	return out
}

func TestPickerTrackedUntilClosed(t *testing.T) {
	s, p, _ := openPicker(t, pickFile, t.TempDir())
	if got := s.dialogs[p.d]; got.Width == 0 {
		t.Fatal("picker not tracked for refitting")
	}
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	if s.picker != nil || len(s.dialogs) != 0 {
		t.Fatal("closed picker still tracked")
	}
}

func TestPickerSearchFiltersFolder(t *testing.T) {
	dir := pickerTree(t)
	_, p, _ := openPicker(t, pickFile, dir)
	if want := []string{"alpine", "sub", "alpha.txt", "Beta.conf"}; !reflect.DeepEqual(shownNames(p), want) {
		t.Fatalf("listing = %v, want %v", shownNames(p), want)
	}
	test.Type(p.search, "ALP")
	if want := []string{"alpine", "alpha.txt"}; !reflect.DeepEqual(shownNames(p), want) {
		t.Fatalf("filtered = %v, want %v", shownNames(p), want)
	}
	p.search.SetText("al txt")
	if want := []string{"alpha.txt"}; !reflect.DeepEqual(shownNames(p), want) {
		t.Fatalf("two words = %v, want %v", shownNames(p), want)
	}
}

func TestPickerSubfolderSearchPicksFile(t *testing.T) {
	dir := pickerTree(t)
	_, p, got := openPicker(t, pickFile, dir)
	p.deep = true
	test.Type(p.search, "alpha")
	want := []string{"alpha.txt", "sub/alphabet", "sub/deeper/ALPHA.key"}
	if !reflect.DeepEqual(shownNames(p), want) {
		t.Fatalf("subfolder search = %v, want %v", shownNames(p), want)
	}
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if *got != filepath.Join(dir, "sub/alphabet") {
		t.Fatalf("picked %q", *got)
	}
}

func TestPickerArrowHighlightsFolderEnterOpensIt(t *testing.T) {
	dir := pickerTree(t)
	_, p, got := openPicker(t, pickFile, dir)
	test.Type(p.search, "sub")
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	if p.dir != dir || p.chosen != "" {
		t.Fatalf("arrow entered %q / chose %q", p.dir, p.chosen)
	}
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if p.dir != filepath.Join(dir, "sub") || p.search.Text != "" || *got != "" {
		t.Fatalf("enter: dir %q search %q picked %q", p.dir, p.search.Text, *got)
	}
}

func TestPickerFolderModeChoosesShownFolder(t *testing.T) {
	dir := pickerTree(t)
	_, p, got := openPicker(t, pickFolder, dir)
	if want := []string{"alpine", "sub"}; !reflect.DeepEqual(shownNames(p), want) {
		t.Fatalf("folder listing = %v, want %v", shownNames(p), want)
	}
	p.search.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if *got != dir {
		t.Fatalf("picked %q, want %q", *got, dir)
	}
}

func TestPickerCtrlFFocusesSearch(t *testing.T) {
	s, p, _ := openPicker(t, pickFile, t.TempDir())
	canvas := s.win.Canvas()
	canvas.Focus(p.path)
	p.path.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault})
	if canvas.Focused() != p.search {
		t.Fatal("Ctrl+F in the path field did not focus the search")
	}
	canvas.Unfocus()
	s.focusSearch()
	if canvas.Focused() != p.search {
		t.Fatal("window Ctrl+F did not focus the picker's search")
	}
}

func TestPickerTypedFilePathSelectsIt(t *testing.T) {
	dir := pickerTree(t)
	_, p, _ := openPicker(t, pickFile, "/")
	p.goTo(filepath.Join(dir, "sub", "alphabet"))
	if p.dir != filepath.Join(dir, "sub") || p.chosen != filepath.Join(dir, "sub", "alphabet") {
		t.Fatalf("dir %q chosen %q", p.dir, p.chosen)
	}
}

func TestSearchTreeFollowsNoSymlink(t *testing.T) {
	dir := pickerTree(t)
	if err := os.Symlink(dir, filepath.Join(dir, "sub", "loop")); err != nil {
		t.Fatal(err)
	}
	res, truncated := searchTree(dir, "alpha", listOpts{}, func() bool { return false })
	if truncated || len(res) != 3 {
		t.Fatalf("results %v truncated %v", res, truncated)
	}
	loop, _ := listDir(filepath.Join(dir, "sub"), listOpts{foldersOnly: true})
	if len(loop) != 2 || loop[1].rel != "loop" || !loop[1].dir {
		t.Fatalf("symlink to a folder not listed as one: %v", loop)
	}
	hidden, _ := searchTree(dir, "secret", listOpts{hidden: true}, func() bool { return false })
	if len(hidden) != 1 {
		t.Fatalf("hidden search = %v", hidden)
	}
	if res, _ := searchTree(dir, "alpha", listOpts{}, func() bool { return true }); len(res) != 0 {
		t.Fatalf("stopped search returned %v", res)
	}
}
