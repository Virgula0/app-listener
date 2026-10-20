//go:build gui

package gui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
)

func newTestFadeScroll(t *testing.T, content fyne.Size, both bool) *fadeScroll {
	t.Helper()
	a := test.NewApp()
	t.Cleanup(a.Quit)
	r := canvas.NewRectangle(nil)
	r.SetMinSize(content)
	s := fadeScrollOf(r, both)
	w := a.NewWindow("x")
	w.SetPadded(false)
	w.SetContent(s)
	w.Resize(fyne.NewSize(200, 100))
	return s
}

// settle runs a glide to its end on a fake clock: the test driver ticks an animation once, at once.
func settle(t *testing.T, s *fadeScroll) {
	t.Helper()
	at := time.Now()
	s.clock = func() time.Time { return at }
	for range 1000 {
		if s.glider == nil {
			return
		}
		at = at.Add(16 * time.Millisecond)
		s.glideStep(0)
	}
	t.Fatalf("glide still running at %v, heading for %v", s.offset, s.target)
}

func TestFadeScrollClampsOffset(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(500, 400), false)
	s.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-50, -1000)})
	settle(t, s)
	if s.offset.X != 0 || s.offset.Y != 300 {
		t.Fatalf("offset %v, want (0,300): vertical only, clamped to the overflow", s.offset)
	}
	if got := s.content.Position(); got.Y != -300 {
		t.Fatalf("content at %v, want y=-300", got)
	}
	s.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, 1000)})
	settle(t, s)
	if s.offset.Y != 0 {
		t.Fatalf("offset %v, want 0 after scrolling back past the top", s.offset)
	}
}

func TestFadeScrollBothDirections(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(500, 400), true)
	s.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-120, -40)})
	settle(t, s)
	if want := fyne.NewPos(120*wheelGain, 40*wheelGain); s.offset != want {
		t.Fatalf("offset %v, want %v", s.offset, want)
	}
	if !s.v.Visible() || !s.h.Visible() {
		t.Fatal("both bars should show when the content overflows both ways")
	}
}

func TestFadeScrollBarsHiddenWhenContentFits(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(50, 50), true)
	if s.v.Visible() || s.h.Visible() {
		t.Fatal("a bar shows while the content fits: it would take clicks at the edge for nothing")
	}
}

func TestFadeScrollFadesUnlessHeld(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(200, 400), false)
	// The test driver runs an animation straight to its end: the bars have faded out.
	s.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -10)})
	settle(t, s)
	if s.alpha != 0 {
		t.Fatalf("alpha %v after the fade, want 0", s.alpha)
	}
	s.v.MouseIn(nil)
	if s.alpha != 1 {
		t.Fatalf("alpha %v while hovered, want 1", s.alpha)
	}
	s.v.MouseOut()
	if s.alpha != 0 {
		t.Fatalf("alpha %v after leaving, want 0", s.alpha)
	}
}

func TestFadeBarDragMovesContent(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(200, 400), false)
	s.v.Dragged(&fyne.DragEvent{Dragged: fyne.NewDelta(0, 1000)})
	s.v.DragEnd()
	if s.offset.Y != 300 {
		t.Fatalf("offset %v after dragging past the end, want 300", s.offset)
	}
}

func TestFadeScrollWheelStepsFromGlideEnd(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(200, 400), false)
	s.target = fyne.NewPos(0, 100) // mid-glide: the view hasn't got there yet
	s.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -25)})
	settle(t, s)
	if want := 100 + 25*wheelGain; s.offset.Y != want || s.target.Y != want {
		t.Fatalf("offset %v target %v, want %v: a step adds to the glide's end, not the view", s.offset, s.target, want)
	}
}

func TestFadeBarDragEndsGlide(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(200, 400), false)
	s.target = fyne.NewPos(0, 300)
	s.v.Dragged(&fyne.DragEvent{Dragged: fyne.NewDelta(0, 1)})
	s.v.DragEnd()
	if s.target != s.offset {
		t.Fatalf("target %v, offset %v: a drag must leave no glide running", s.target, s.offset)
	}
}

func TestFadeBarTapPages(t *testing.T) {
	s := newTestFadeScroll(t, fyne.NewSize(200, 400), false)
	s.v.Tapped(&fyne.PointEvent{Position: fyne.NewPos(5, 95)})
	settle(t, s)
	if s.offset.Y != 90 {
		t.Fatalf("offset %v after a tap below the thumb, want one page (90)", s.offset)
	}
}
