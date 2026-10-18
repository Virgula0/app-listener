//go:build gui

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// navItem is one sidebar entry: icon + label (icon only when compact), a tinted pill and accent bar
// when selected.
type navItem struct {
	widget.BaseWidget
	iconName string
	label    string
	selected bool
	hovered  bool
	compact  bool
	// badge is a dot next to the label (on the icon when compact): something new on that page.
	badge bool
	onTap func()
}

func newNavItem(iconName, label string, onTap func()) *navItem {
	n := &navItem{iconName: iconName, label: label, onTap: onTap}
	n.ExtendBaseWidget(n)
	return n
}

func (n *navItem) setSelected(on bool) {
	if n.selected != on {
		n.selected = on
		n.Refresh()
	}
}

func (n *navItem) setBadge(on bool) {
	if n.badge != on {
		n.badge = on
		n.Refresh()
	}
}

func (n *navItem) setCompact(on bool) {
	if n.compact != on {
		n.compact = on
		n.Refresh()
	}
}

func (n *navItem) Tapped(*fyne.PointEvent) {
	if n.onTap != nil {
		n.onTap()
	}
}

func (n *navItem) MouseIn(*desktop.MouseEvent)    { n.hovered = true; n.Refresh() }
func (n *navItem) MouseMoved(*desktop.MouseEvent) {}
func (n *navItem) MouseOut()                      { n.hovered = false; n.Refresh() }
func (n *navItem) Cursor() desktop.Cursor         { return desktop.PointerCursor }

func (n *navItem) CreateRenderer() fyne.WidgetRenderer {
	r := &navItemRenderer{
		n:   n,
		bg:  canvas.NewRectangle(color.Transparent),
		bar: canvas.NewRectangle(color.Transparent),
		img: canvas.NewImageFromResource(nil),
		txt: canvas.NewText(n.label, color.Transparent),
		dot: canvas.NewCircle(color.Transparent),
	}
	r.img.FillMode = canvas.ImageFillContain
	r.Refresh()
	return r
}

type navItemRenderer struct {
	n       *navItem
	bg, bar *canvas.Rectangle
	img     *canvas.Image
	txt     *canvas.Text
	dot     *canvas.Circle
}

const navIconSize, navHeight, navDot = 20, 40, 8

func (r *navItemRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bar.Move(fyne.NewPos(0, size.Height*0.25))
	r.bar.Resize(fyne.NewSize(3, size.Height*0.5))
	pad := theme.Padding() * 2
	r.dot.Resize(fyne.NewSquareSize(navDot))
	if r.n.compact {
		r.img.Move(fyne.NewPos((size.Width-navIconSize)/2, (size.Height-navIconSize)/2))
		r.img.Resize(fyne.NewSquareSize(navIconSize))
		r.dot.Move(fyne.NewPos((size.Width+navIconSize)/2-navDot/2, (size.Height-navIconSize)/2-navDot/2))
		return
	}
	r.dot.Move(fyne.NewPos(size.Width-pad-navDot, (size.Height-navDot)/2))
	r.img.Move(fyne.NewPos(pad+4, (size.Height-navIconSize)/2))
	r.img.Resize(fyne.NewSquareSize(navIconSize))
	ts := r.txt.MinSize()
	r.txt.Move(fyne.NewPos(pad+4+navIconSize+pad, (size.Height-ts.Height)/2))
	r.txt.Resize(ts)
}

func (r *navItemRenderer) MinSize() fyne.Size {
	pad := theme.Padding() * 2
	if r.n.compact {
		return fyne.NewSize(2*pad+8+navIconSize, navHeight)
	}
	return fyne.NewSize(pad+4+navIconSize+pad+r.txt.MinSize().Width+pad+navDot+pad, navHeight)
}

func (r *navItemRenderer) Refresh() {
	primary := themeColor(theme.ColorNamePrimary)
	r.bg.CornerRadius = 8
	r.bar.CornerRadius = 1.5
	switch {
	case r.n.selected:
		r.bg.FillColor = tint(primary, 0x2A)
		r.bar.FillColor = primary
		r.txt.Color = themeColor(theme.ColorNameForeground)
		r.img.Resource = coloredIcon(r.n.iconName, theme.ColorNamePrimary)
	case r.n.hovered:
		r.bg.FillColor = themeColor(theme.ColorNameHover)
		r.bar.FillColor = color.Transparent
		r.txt.Color = themeColor(theme.ColorNameForeground)
		r.img.Resource = icon(r.n.iconName)
	default:
		r.bg.FillColor = color.Transparent
		r.bar.FillColor = color.Transparent
		r.txt.Color = themeColor(colorMuted)
		r.img.Resource = coloredIcon(r.n.iconName, colorMuted)
	}
	r.txt.TextStyle = fyne.TextStyle{Bold: r.n.selected}
	r.txt.Hidden = r.n.compact
	r.dot.FillColor = themeColor(colorInfo)
	r.dot.Hidden = !r.n.badge
	r.dot.Refresh()
	r.bg.Refresh()
	r.bar.Refresh()
	r.img.Refresh()
	r.txt.Refresh()
}

func (r *navItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.bar, r.img, r.txt, r.dot}
}

func (r *navItemRenderer) Destroy() {}
