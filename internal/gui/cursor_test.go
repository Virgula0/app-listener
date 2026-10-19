//go:build gui

package gui

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// shimmedButtons counts the *widget.Buttons under o and how many carry a buttonCursor.
func shimmedButtons(o fyne.CanvasObject) (buttons, shims int) {
	c, ok := o.(*fyne.Container)
	if !ok {
		return 0, 0
	}
	for _, child := range c.Objects {
		switch v := child.(type) {
		case *widget.Button:
			buttons++
		case *buttonCursor:
			shims++
		default:
			b, s := shimmedButtons(v)
			buttons, shims = buttons+b, shims+s
		}
	}
	return buttons, shims
}

func TestDialogButtonsGetHandCursor(t *testing.T) {
	a := test.NewApp()
	t.Cleanup(a.Quit)
	w := a.NewWindow("x")
	w.Resize(fyne.NewSize(800, 600))
	showConfirm("Title", "Message", func(bool) {}, w)
	handCursorsOn(w) // a second pass must not stack another shim
	var b, s int
	for _, o := range w.Canvas().Overlays().List() {
		ob, os := shimmedButtons(reflect.ValueOf(o).Elem().FieldByName("Content").Interface().(*widget.PopUp).Content)
		b, s = b+ob, s+os
	}
	if b != 2 || s != 2 {
		t.Fatalf("buttons=%d shims=%d, want 2 and 2", b, s)
	}
}

func TestButtonCursorFollowsState(t *testing.T) {
	test.NewTempApp(t)
	b := widget.NewButton("x", nil)
	c := newButtonCursor(b)
	if c.Cursor() != desktop.PointerCursor {
		t.Fatal("enabled button: want the pointer")
	}
	b.Disable()
	if c.Cursor() != desktop.DefaultCursor {
		t.Fatal("disabled button: want the arrow")
	}
}

func TestTabBarShowsOnlySelectedPage(t *testing.T) {
	test.NewTempApp(t)
	pages := []fyne.CanvasObject{widget.NewLabel("a"), widget.NewLabel("b"), widget.NewLabel("c")}
	tb := newTabBar(tabItem{"A", "eye", pages[0]}, tabItem{"B", "pulse", pages[1]}, tabItem{"C", "cpu", pages[2]})
	test.Tap(tb.buttons[2])
	if tb.SelectedIndex() != 2 {
		t.Fatalf("selected %d, want 2", tb.SelectedIndex())
	}
	for i, p := range pages {
		if p.Visible() != (i == 2) || tb.buttons[i].selected != (i == 2) {
			t.Fatalf("page %d visible=%v selected=%v", i, p.Visible(), tb.buttons[i].selected)
		}
	}
	tb.SelectIndex(9)
	if tb.SelectedIndex() != 2 {
		t.Fatal("out-of-range index changed the selection")
	}
}

func TestChoiceWidgetsPointingHand(t *testing.T) {
	a := test.NewApp()
	t.Cleanup(a.Quit)
	opts := []string{"a", "b"}
	type choice interface {
		desktop.Cursorable
		fyne.Disableable
	}
	for _, w := range []choice{newCheck("x", nil), newRadioGroup(opts, nil), newCheckGroup(opts, nil),
		newSelect(opts, nil)} {
		if w.Cursor() != desktop.PointerCursor {
			t.Errorf("%T: enabled cursor %v", w, w.Cursor())
		}
		w.Disable()
		if w.Cursor() != desktop.DefaultCursor {
			t.Errorf("%T: disabled cursor %v", w, w.Cursor())
		}
	}
}
