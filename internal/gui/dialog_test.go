//go:build gui

package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/test"
)

func TestShowFileSizesAfterShow(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("x")
	w.Resize(fyne.NewSize(700, 500))
	s := newShell(a, w, Options{})
	d := dialog.NewFileOpen(func(fyne.URIReadCloser, error) {}, w)
	s.showFile(d, t.TempDir()) // a Resize before Show dereferenced the unbuilt dialog
	if got := s.dialogs[d]; got.Width == 0 {
		t.Fatal("file dialog not tracked for refitting")
	}
}
