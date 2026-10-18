//go:build gui

package gui

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	fadeIn   = 120 * time.Millisecond
	fadeHold = 900 * time.Millisecond
	fadeOut  = 450 * time.Millisecond
	// A glide eases wheel steps and page taps (fyne applies a notch as one 25px jump): each frame
	// closes 1-e^(-dt/glideTau) of the gap, ~95% in 3 tau.
	glideTau     = 40 * time.Millisecond
	barStrip     = float32(12) // hit area along the edge
	barThin      = float32(5)
	barThick     = float32(8) // hovered or dragged
	barInset     = float32(3)
	barMinLength = float32(28)
	pageStep     = float32(0.9)
)

// fadeScroll is the pages' scroll container: it clips and offsets its content like
// container.Scroll, but its bars are thin overlay pills that fade in while scrolling (or hovered)
// and fade out once idle. Tables and lists keep fyne's own scroller.
type fadeScroll struct {
	widget.BaseWidget
	content fyne.CanvasObject
	// both: horizontal scrolling too (the window-level page wrapper); else vertical only.
	both    bool
	minSize fyne.Size
	offset  fyne.Position
	// target is where the running glide ends; equal to offset when none runs.
	target  fyne.Position
	glider  *fyne.Animation
	glideAt time.Time // the glide's last step
	clock   func() time.Time
	v, h    *fadeBar
	alpha   float32
	anim    *fyne.Animation
	render  *fadeScrollRenderer
}

var _ fyne.Scrollable = (*fadeScroll)(nil)

func newFadeVScroll(content fyne.CanvasObject) *fadeScroll { return fadeScrollOf(content, false) }

// newFadeScroll scrolls both ways.
func newFadeScroll(content fyne.CanvasObject) *fadeScroll { return fadeScrollOf(content, true) }

func fadeScrollOf(content fyne.CanvasObject, both bool) *fadeScroll {
	s := &fadeScroll{content: content, both: both, clock: time.Now}
	s.v = &fadeBar{s: s, vertical: true}
	s.h = &fadeBar{s: s}
	s.v.ExtendBaseWidget(s.v)
	s.h.ExtendBaseWidget(s.h)
	s.ExtendBaseWidget(s)
	return s
}

func (s *fadeScroll) CreateRenderer() fyne.WidgetRenderer {
	s.render = &fadeScrollRenderer{s: s}
	return s.render
}

// SetMinSize is container.Scroll's: the least size the scroll asks for, whatever its content.
func (s *fadeScroll) SetMinSize(size fyne.Size) { s.minSize = size }

func (s *fadeScroll) MinSize() fyne.Size {
	m := fyne.NewSize(32, 32).Max(s.minSize)
	if !s.both {
		m.Width = fyne.Max(m.Width, s.content.MinSize().Width)
	}
	return m
}

func (s *fadeScroll) Scrolled(e *fyne.ScrollEvent) {
	dx, dy := e.Scrolled.DX, e.Scrolled.DY
	if !s.both {
		if dy == 0 {
			dy = dx
		}
		dx = 0
	}
	s.glideTo(fyne.NewPos(s.target.X-dx, s.target.Y-dy))
}

// scrollTo moves the view to off (clamped) at once, ending any glide, and shows the bars.
func (s *fadeScroll) scrollTo(off fyne.Position) {
	s.stopGlide()
	off = s.clamped(off)
	s.target = off
	if off == s.offset {
		return
	}
	s.moveTo(off)
	s.wake()
}

// glideTo eases the view to off (clamped). Mid-glide it only moves the end: restarting the animation
// would cost a frame, as fyne starts a new one only after the current frame.
func (s *fadeScroll) glideTo(off fyne.Position) {
	off = s.clamped(off)
	if off == s.target {
		return
	}
	s.target = off
	s.wake()
	if s.glider != nil {
		return
	}
	s.glideAt = s.clock()
	s.glider = fyne.NewAnimation(time.Second, s.glideStep)
	s.glider.RepeatCount = fyne.AnimationRepeatForever
	s.glider.Curve = fyne.AnimationLinear
	s.glider.Start()
}

// glideStep ignores the animation's progress: it steps by the time since the last step and stops
// once the view is within half a pixel of the end.
func (s *fadeScroll) glideStep(float32) {
	now := s.clock()
	k := float32(1 - math.Exp(-float64(now.Sub(s.glideAt))/float64(glideTau)))
	s.glideAt = now
	gap := s.target.Subtract(s.offset)
	if math.Abs(float64(gap.X)) < 0.5 && math.Abs(float64(gap.Y)) < 0.5 {
		s.moveTo(s.target)
		s.stopGlide()
		return
	}
	s.moveTo(s.clamped(s.offset.Add(fyne.NewDelta(gap.X*k, gap.Y*k))))
}

func (s *fadeScroll) stopGlide() {
	if s.glider != nil {
		s.glider.Stop()
		s.glider = nil
	}
}

func (s *fadeScroll) moveTo(off fyne.Position) {
	if off == s.offset {
		return
	}
	s.offset = off
	if s.render != nil {
		s.render.place()
	}
}

// overflow is how far the content extends past the viewport on each axis.
func (s *fadeScroll) overflow() fyne.Size {
	c, v := s.content.Size(), s.Size()
	return fyne.NewSize(fyne.Max(0, c.Width-v.Width), fyne.Max(0, c.Height-v.Height))
}

func (s *fadeScroll) clamped(p fyne.Position) fyne.Position {
	o := s.overflow()
	return fyne.NewPos(fyne.Min(fyne.Max(0, p.X), o.Width), fyne.Min(fyne.Max(0, p.Y), o.Height))
}

// clamp keeps the view and the glide's end inside the content after a resize.
func (s *fadeScroll) clamp() {
	s.offset = s.clamped(s.offset)
	s.target = s.clamped(s.target)
}

func (s *fadeScroll) held() bool {
	return s.v.hovered || s.v.dragging || s.h.hovered || s.h.dragging
}

// wake fades the bars in, holds them, then fades them out unless a bar is hovered or dragged
// (its MouseOut / DragEnd wakes again).
func (s *fadeScroll) wake() {
	if s.anim != nil {
		s.anim.Stop()
	}
	from := s.alpha
	total := fadeIn + fadeHold + fadeOut
	s.anim = fyne.NewAnimation(total, func(p float32) {
		at := time.Duration(p * float32(total))
		switch {
		case at < fadeIn:
			s.setAlpha(from + (1-from)*float32(at)/float32(fadeIn))
		case at < fadeIn+fadeHold || s.held():
			s.setAlpha(1)
		default:
			s.setAlpha(1 - float32(at-fadeIn-fadeHold)/float32(fadeOut))
		}
	})
	s.anim.Curve = fyne.AnimationLinear
	s.anim.Start()
}

func (s *fadeScroll) setAlpha(a float32) {
	s.alpha = fyne.Min(1, fyne.Max(0, a))
	s.v.paint()
	s.h.paint()
}

type fadeScrollRenderer struct {
	s *fadeScroll
}

func (r *fadeScrollRenderer) Layout(size fyne.Size) {
	s := r.s
	cmin := s.content.MinSize()
	cs := fyne.NewSize(size.Width, fyne.Max(size.Height, cmin.Height))
	if s.both {
		cs.Width = fyne.Max(size.Width, cmin.Width)
	}
	s.content.Resize(cs)
	s.clamp()
	r.place()
}

// place moves the content to the offset and the bars to the edges.
func (r *fadeScrollRenderer) place() {
	s := r.s
	size := s.Size()
	s.content.Move(fyne.NewPos(-s.offset.X, -s.offset.Y))
	o := s.overflow()
	s.v.Hidden = o.Height <= 0
	s.h.Hidden = !s.both || o.Width <= 0
	s.v.Move(fyne.NewPos(size.Width-barStrip, 0))
	s.v.Resize(fyne.NewSize(barStrip, size.Height))
	s.h.Move(fyne.NewPos(0, size.Height-barStrip))
	s.h.Resize(fyne.NewSize(size.Width, barStrip))
	s.v.Refresh()
	s.h.Refresh()
	canvas.Refresh(s.content)
}

func (r *fadeScrollRenderer) MinSize() fyne.Size { return r.s.MinSize() }

func (r *fadeScrollRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.s.content, r.s.v, r.s.h}
}

func (r *fadeScrollRenderer) Refresh() {
	r.Layout(r.s.Size())
	r.s.content.Refresh()
}

func (r *fadeScrollRenderer) Destroy() {
	r.s.stopGlide()
	if r.s.anim != nil {
		r.s.anim.Stop()
	}
}

// fadeBar is one axis' bar: an edge strip that takes hover, drag and taps, drawing a rounded thumb
// whose opacity follows the scroll's fade.
type fadeBar struct {
	widget.BaseWidget
	s                 *fadeScroll
	vertical          bool
	hovered, dragging bool
	thumb             *canvas.Rectangle
}

var (
	_ desktop.Hoverable = (*fadeBar)(nil)
	_ fyne.Draggable    = (*fadeBar)(nil)
	_ fyne.Tappable     = (*fadeBar)(nil)
)

func (b *fadeBar) CreateRenderer() fyne.WidgetRenderer {
	b.thumb = canvas.NewRectangle(color.Transparent)
	b.paint()
	return &fadeBarRenderer{b: b}
}

// axis returns the viewport length, content length and offset along the bar.
func (b *fadeBar) axis() (view, content, off float32) {
	if b.vertical {
		return b.s.Size().Height, b.s.content.Size().Height, b.s.offset.Y
	}
	return b.s.Size().Width, b.s.content.Size().Width, b.s.offset.X
}

// span is the thumb's leading edge and length along the bar.
func (b *fadeBar) span() (lead, length float32) {
	view, content, off := b.axis()
	track := view - 2*barInset
	if content <= view || track <= 0 {
		return barInset, fyne.Max(0, track)
	}
	length = fyne.Min(track, fyne.Max(barMinLength, track*view/content))
	return barInset + (track-length)*off/(content-view), length
}

func (b *fadeBar) paint() {
	if b.thumb == nil {
		return
	}
	base := float32(0x70)
	if b.hovered || b.dragging {
		base = 0xB0
	}
	b.thumb.FillColor = tint(themeColor(theme.ColorNameForeground), uint8(base*b.s.alpha))
	canvas.Refresh(b.thumb)
}

// at is the position from p with this bar's axis set to along.
func (b *fadeBar) at(p fyne.Position, along float32) fyne.Position {
	if b.vertical {
		p.Y = along
	} else {
		p.X = along
	}
	return p
}

func (b *fadeBar) Dragged(e *fyne.DragEvent) {
	b.dragging = true
	view, content, off := b.axis()
	_, length := b.span()
	free := view - 2*barInset - length
	if free <= 0 {
		return
	}
	d := e.Dragged.DX
	if b.vertical {
		d = e.Dragged.DY
	}
	b.s.scrollTo(b.at(b.s.offset, off+d*(content-view)/free))
	b.Refresh()
}

func (b *fadeBar) DragEnd() {
	b.dragging = false
	b.Refresh()
	b.s.wake()
}

// Tapped beside the thumb pages toward the tap, from the running glide's end.
func (b *fadeBar) Tapped(e *fyne.PointEvent) {
	view, _, _ := b.axis()
	lead, length := b.span()
	tap, from := e.Position.X, b.s.target.X
	if b.vertical {
		tap, from = e.Position.Y, b.s.target.Y
	}
	switch {
	case tap < lead:
		b.s.glideTo(b.at(b.s.target, from-view*pageStep))
	case tap > lead+length:
		b.s.glideTo(b.at(b.s.target, from+view*pageStep))
	}
}

func (b *fadeBar) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
	b.s.wake()
}

func (b *fadeBar) MouseMoved(*desktop.MouseEvent) {}

func (b *fadeBar) MouseOut() {
	b.hovered = false
	b.Refresh()
	b.s.wake()
}

type fadeBarRenderer struct {
	b *fadeBar
}

func (r *fadeBarRenderer) Layout(size fyne.Size) {
	b := r.b
	thick := barThin
	if b.hovered || b.dragging {
		thick = barThick
	}
	lead, length := b.span()
	b.thumb.CornerRadius = thick / 2
	if b.vertical {
		b.thumb.Move(fyne.NewPos(size.Width-barInset-thick, lead))
		b.thumb.Resize(fyne.NewSize(thick, length))
	} else {
		b.thumb.Move(fyne.NewPos(lead, size.Height-barInset-thick))
		b.thumb.Resize(fyne.NewSize(length, thick))
	}
}

func (r *fadeBarRenderer) MinSize() fyne.Size { return fyne.NewSize(barStrip, barStrip) }

func (r *fadeBarRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.b.thumb} }

func (r *fadeBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	r.b.paint()
}

func (r *fadeBarRenderer) Destroy() {}
