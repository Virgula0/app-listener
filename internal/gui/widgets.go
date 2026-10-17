//go:build gui

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// panel is a rounded surface card around content. Its colors are re-read from the theme on every
// Refresh, so a light/dark switch repaints it.
type panel struct {
	widget.BaseWidget
	content fyne.CanvasObject
	inset   float32
	// fill overrides the surface color (nil = colorSurface); hover lightens it.
	fill    fyne.ThemeColorName
	hovered bool
}

func newPanel(content fyne.CanvasObject) *panel {
	p := &panel{content: content, inset: theme.Padding() * 3}
	p.ExtendBaseWidget(p)
	return p
}

func (p *panel) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.Transparent)
	hover := canvas.NewRectangle(color.Transparent)
	r := &panelRenderer{p: p, bg: bg, hover: hover}
	r.Refresh()
	return r
}

type panelRenderer struct {
	p         *panel
	bg, hover *canvas.Rectangle
}

func (r *panelRenderer) Layout(size fyne.Size) {
	// A tiling WM can impose a window smaller than the content's MinSize: never lay out a negative
	// size (the software painter panics on one).
	size = size.Max(fyne.Size{})
	r.bg.Resize(size)
	r.hover.Resize(size)
	in := r.p.inset
	r.p.content.Move(fyne.NewPos(in, in))
	r.p.content.Resize(size.SubtractWidthHeight(2*in, 2*in).Max(fyne.Size{}))
}

func (r *panelRenderer) MinSize() fyne.Size {
	return r.p.content.MinSize().AddWidthHeight(2*r.p.inset, 2*r.p.inset)
}

func (r *panelRenderer) Refresh() {
	fill := r.p.fill
	if fill == "" {
		fill = colorSurface
	}
	radius := theme.Size(theme.SizeNameCardRadius)
	r.bg.FillColor = themeColor(fill)
	r.bg.StrokeColor = themeColor(colorSurfaceBorder)
	r.bg.StrokeWidth = 1
	r.bg.CornerRadius = radius
	r.hover.CornerRadius = radius
	r.hover.FillColor = color.Transparent
	if r.p.hovered {
		r.hover.FillColor = themeColor(theme.ColorNameHover)
	}
	r.bg.Refresh()
	r.hover.Refresh()
	r.p.content.Refresh()
}

func (r *panelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.hover, r.p.content}
}

func (r *panelRenderer) Destroy() {}

// clickPanel is a panel that reacts to hover and taps (workspace cards).
type clickPanel struct {
	panel
	onTap func()
}

func newClickPanel(content fyne.CanvasObject, onTap func()) *clickPanel {
	c := &clickPanel{panel: panel{content: content, inset: theme.Padding() * 3}, onTap: onTap}
	c.ExtendBaseWidget(c)
	return c
}

func (c *clickPanel) Tapped(*fyne.PointEvent) {
	if c.onTap != nil {
		c.onTap()
	}
}

func (c *clickPanel) MouseIn(*desktop.MouseEvent)    { c.hovered = true; c.Refresh() }
func (c *clickPanel) MouseMoved(*desktop.MouseEvent) {}
func (c *clickPanel) MouseOut()                      { c.hovered = false; c.Refresh() }
func (c *clickPanel) Cursor() desktop.Cursor         { return desktop.PointerCursor }

// pill is a small rounded status label tinted by a theme color.
type pill struct {
	widget.BaseWidget
	text  string
	color fyne.ThemeColorName
}

func newPill(text string, c fyne.ThemeColorName) *pill {
	p := &pill{text: text, color: c}
	p.ExtendBaseWidget(p)
	return p
}

// Set changes the pill's text and color.
func (p *pill) Set(text string, c fyne.ThemeColorName) {
	p.text, p.color = text, c
	p.Refresh()
}

func (p *pill) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.Transparent)
	txt := canvas.NewText("", color.Transparent)
	txt.TextStyle = fyne.TextStyle{Bold: true}
	r := &pillRenderer{p: p, bg: bg, txt: txt}
	r.Refresh()
	return r
}

type pillRenderer struct {
	p   *pill
	bg  *canvas.Rectangle
	txt *canvas.Text
}

const pillPadH, pillPadV = 10, 3

func (r *pillRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	ts := r.txt.MinSize()
	r.txt.Move(fyne.NewPos((size.Width-ts.Width)/2, (size.Height-ts.Height)/2))
	r.txt.Resize(ts)
}

func (r *pillRenderer) MinSize() fyne.Size {
	return r.txt.MinSize().AddWidthHeight(2*pillPadH, 2*pillPadV)
}

func (r *pillRenderer) Refresh() {
	c := themeColor(r.p.color)
	r.txt.Text = r.p.text
	r.txt.Color = c
	r.txt.TextSize = theme.Size(theme.SizeNameCaptionText)
	r.bg.FillColor = tint(c, 0x26)
	r.bg.CornerRadius = r.MinSize().Height / 2
	r.bg.Refresh()
	r.txt.Refresh()
}

func (r *pillRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.bg, r.txt} }
func (r *pillRenderer) Destroy()                     {}

// iconTile is an icon centered on a rounded square tinted with the icon's color.
type iconTile struct {
	widget.BaseWidget
	name  string
	color fyne.ThemeColorName
	side  float32
}

func newIconTile(name string, c fyne.ThemeColorName, side float32) *iconTile {
	t := &iconTile{name: name, color: c, side: side}
	t.ExtendBaseWidget(t)
	return t
}

// Set swaps the tile's icon and color.
func (t *iconTile) Set(name string, c fyne.ThemeColorName) {
	t.name, t.color = name, c
	t.Refresh()
}

func (t *iconTile) CreateRenderer() fyne.WidgetRenderer {
	r := &iconTileRenderer{t: t, bg: canvas.NewRectangle(color.Transparent), img: canvas.NewImageFromResource(nil)}
	r.img.FillMode = canvas.ImageFillContain
	r.Refresh()
	return r
}

type iconTileRenderer struct {
	t   *iconTile
	bg  *canvas.Rectangle
	img *canvas.Image
}

func (r *iconTileRenderer) Layout(size fyne.Size) {
	side := fyne.Min(size.Width, size.Height)
	r.bg.Resize(fyne.NewSquareSize(side))
	in := side * 0.24
	r.img.Move(fyne.NewPos(in, in))
	r.img.Resize(fyne.NewSquareSize(side - 2*in))
}

func (r *iconTileRenderer) MinSize() fyne.Size { return fyne.NewSquareSize(r.t.side) }

func (r *iconTileRenderer) Refresh() {
	r.bg.FillColor = tint(themeColor(r.t.color), 0x24)
	r.bg.CornerRadius = r.t.side * 0.28
	r.img.Resource = coloredIcon(r.t.name, r.t.color)
	r.bg.Refresh()
	r.img.Refresh()
}

func (r *iconTileRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.bg, r.img} }
func (r *iconTileRenderer) Destroy()                     {}

// mutedText is secondary text in colorMuted, wrapping.
func mutedText(s string) *widget.Label {
	l := widget.NewLabel(s)
	l.Wrapping = fyne.TextWrapWord
	l.Importance = widget.LowImportance
	return l
}

// sectionTitle is an uppercase caption heading a group of cards.
func sectionTitle(s string) fyne.CanvasObject {
	l := widget.NewLabelWithStyle(s, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	l.Importance = widget.LowImportance
	l.SizeName = theme.SizeNameCaptionText
	return l
}

// boldLabel is a single-line bold label.
func boldLabel(s string) *widget.Label {
	return widget.NewLabelWithStyle(s, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

// hfill is an HBox whose first object takes the remaining width.
func hfill(grow fyne.CanvasObject, rest ...fyne.CanvasObject) *fyne.Container {
	return container.NewBorder(nil, nil, nil, container.NewHBox(rest...), grow)
}
