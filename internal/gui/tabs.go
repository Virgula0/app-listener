//go:build gui

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// tabItem is one page of a tabBar.
type tabItem struct {
	label, iconName string
	content         fyne.CanvasObject
}

// tabBar is a top tab strip over one visible page. It replaces container.AppTabs, whose internal
// tab buttons keep the arrow cursor.
type tabBar struct {
	root     *fyne.Container
	buttons  []*tabButton
	pages    []fyne.CanvasObject
	selected int
}

func newTabBar(items ...tabItem) *tabBar {
	t := &tabBar{}
	strip, pages := container.NewHBox(), container.NewStack()
	for i, it := range items {
		b := newTabButton(it.label, it.iconName, func() { t.SelectIndex(i) })
		t.buttons = append(t.buttons, b)
		t.pages = append(t.pages, it.content)
		strip.Add(b)
		pages.Add(it.content)
	}
	top := container.New(layout.NewCustomPaddedVBoxLayout(0), strip, widget.NewSeparator())
	t.root = container.NewBorder(top, nil, nil, nil, pages)
	t.SelectIndex(0)
	return t
}

// SelectIndex shows page i; out of range is ignored.
func (t *tabBar) SelectIndex(i int) {
	if i < 0 || i >= len(t.pages) {
		return
	}
	t.selected = i
	for j, p := range t.pages {
		t.buttons[j].setSelected(j == i)
		if j == i {
			p.Show()
		} else {
			p.Hide()
		}
	}
}

func (t *tabBar) SelectedIndex() int { return t.selected }

// tabButton is one tab: icon + bold label, primary-colored with an underline when selected.
type tabButton struct {
	widget.BaseWidget
	label, iconName   string
	selected, hovered bool
	onTap             func()
}

func newTabButton(label, iconName string, onTap func()) *tabButton {
	b := &tabButton{label: label, iconName: iconName, onTap: onTap}
	b.ExtendBaseWidget(b)
	return b
}

func (b *tabButton) setSelected(on bool) {
	if b.selected != on {
		b.selected = on
		b.Refresh()
	}
}

func (b *tabButton) Tapped(*fyne.PointEvent) {
	if b.onTap != nil {
		b.onTap()
	}
}

func (b *tabButton) MouseIn(*desktop.MouseEvent)    { b.hovered = true; b.Refresh() }
func (b *tabButton) MouseMoved(*desktop.MouseEvent) {}
func (b *tabButton) MouseOut()                      { b.hovered = false; b.Refresh() }
func (b *tabButton) Cursor() desktop.Cursor         { return desktop.PointerCursor }

func (b *tabButton) CreateRenderer() fyne.WidgetRenderer {
	r := &tabButtonRenderer{
		b:   b,
		bg:  canvas.NewRectangle(color.Transparent),
		ind: canvas.NewRectangle(color.Transparent),
		img: canvas.NewImageFromResource(nil),
		txt: canvas.NewText(b.label, color.Transparent),
	}
	r.img.FillMode = canvas.ImageFillContain
	r.txt.TextStyle = fyne.TextStyle{Bold: true}
	r.Refresh()
	return r
}

type tabButtonRenderer struct {
	b       *tabButton
	bg, ind *canvas.Rectangle
	img     *canvas.Image
	txt     *canvas.Text
}

const tabIndicator = 2

func (r *tabButtonRenderer) Layout(size fyne.Size) {
	pad, is := theme.InnerPadding(), theme.IconInlineSize()
	r.bg.Resize(size)
	r.ind.Move(fyne.NewPos(0, size.Height-tabIndicator))
	r.ind.Resize(fyne.NewSize(size.Width, tabIndicator))
	r.img.Move(fyne.NewPos(pad, (size.Height-is)/2))
	r.img.Resize(fyne.NewSquareSize(is))
	ts := r.txt.MinSize()
	r.txt.Move(fyne.NewPos(pad+is+theme.Padding(), (size.Height-ts.Height)/2))
	r.txt.Resize(ts)
}

func (r *tabButtonRenderer) MinSize() fyne.Size {
	pad, is := theme.InnerPadding(), theme.IconInlineSize()
	ts := r.txt.MinSize()
	return fyne.NewSize(2*pad+is+theme.Padding()+ts.Width, fyne.Max(ts.Height, is)+2*pad)
}

func (r *tabButtonRenderer) Refresh() {
	r.bg.CornerRadius = theme.SelectionRadiusSize()
	r.bg.FillColor = color.Transparent
	if r.b.hovered {
		r.bg.FillColor = themeColor(theme.ColorNameHover)
	}
	r.ind.FillColor = color.Transparent
	r.txt.Color = themeColor(theme.ColorNameForeground)
	r.img.Resource = icon(r.b.iconName)
	if r.b.selected {
		r.ind.FillColor = themeColor(theme.ColorNamePrimary)
		r.txt.Color = r.ind.FillColor
		r.img.Resource = coloredIcon(r.b.iconName, theme.ColorNamePrimary)
	}
	r.bg.Refresh()
	r.ind.Refresh()
	r.img.Refresh()
	r.txt.Refresh()
}

func (r *tabButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.ind, r.img, r.txt}
}

func (r *tabButtonRenderer) Destroy() {}
