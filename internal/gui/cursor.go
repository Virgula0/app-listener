//go:build gui

package gui

import (
	"image/color"
	"reflect"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// buttonCursor is a transparent shim stacked over a fyne-built *widget.Button (dialog buttons)
// to give it the pointing hand. The driver takes the cursor from the last Cursorable under the
// pointer, while taps and hover go to the last Tappable/Hoverable: the shim is only Cursorable,
// so the button below keeps both.
type buttonCursor struct {
	widget.BaseWidget
	b *widget.Button
}

func newButtonCursor(b *widget.Button) *buttonCursor {
	c := &buttonCursor{b: b}
	c.ExtendBaseWidget(c)
	return c
}

func (c *buttonCursor) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func (c *buttonCursor) Cursor() desktop.Cursor {
	if !c.b.Visible() || c.b.Disabled() {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

// handCursors stacks a buttonCursor over every *widget.Button in o's container tree. Idempotent:
// a dialog shown again keeps its shims.
func handCursors(o fyne.CanvasObject) {
	c, ok := o.(*fyne.Container)
	if !ok || slices.ContainsFunc(c.Objects, func(o fyne.CanvasObject) bool {
		_, done := o.(*buttonCursor)
		return done
	}) {
		return
	}
	for i, child := range c.Objects {
		if b, ok := child.(*widget.Button); ok {
			c.Objects[i] = container.NewStack(b, newButtonCursor(b))
			continue
		}
		handCursors(child)
	}
	c.Refresh()
}

// handCursorsOn applies handCursors to the dialogs open on win. A PopUp's overlay is fyne's
// internal OverlayContainer, reachable only through its exported Content field.
func handCursorsOn(win fyne.Window) {
	for _, o := range win.Canvas().Overlays().List() {
		v := reflect.ValueOf(o)
		if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
			continue
		}
		f := v.Elem().FieldByName("Content")
		if !f.IsValid() || !f.CanInterface() {
			continue
		}
		if p, ok := f.Interface().(*widget.PopUp); ok {
			handCursors(p.Content)
		}
	}
}

// showDialog shows d with hand cursors on its buttons. Every dialog goes through it or showSized.
func showDialog(d dialog.Dialog, win fyne.Window) {
	d.Show()
	handCursorsOn(win)
}

func showConfirm(title, message string, callback func(bool), win fyne.Window) {
	showDialog(dialog.NewConfirm(title, message, callback, win), win)
}

func showCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject, callback func(bool),
	win fyne.Window) {
	showDialog(dialog.NewCustomConfirm(title, confirm, dismiss, content, callback, win), win)
}

func showInformation(title, message string, win fyne.Window) {
	showDialog(dialog.NewInformation(title, message, win), win)
}

func showErrorDialog(err error, win fyne.Window) {
	showDialog(dialog.NewError(err, win), win)
}

func showForm(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool), win fyne.Window) {
	showDialog(dialog.NewForm(title, confirm, dismiss, items, callback, win), win)
}
