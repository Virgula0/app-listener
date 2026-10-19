//go:build gui

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// cardGap separates cards in a grid and stacked sections of a page.
const cardGap = 14

// adaptiveGrid lays cards out in up to maxCols equal columns separated by cardGap, dropping columns
// when the widest card no longer fits. Its MinSize is one column wide, so a narrow window wraps
// cards instead of clipping them. The height follows the last laid-out width: fyne re-checks
// MinSize every frame and relayouts the parent when it changes.
type adaptiveGrid struct {
	maxCols int
	cols    int
	// gap overrides cardGap when set.
	gap float32
}

func newAdaptiveGrid(maxCols int, objs ...fyne.CanvasObject) *fyne.Container {
	return container.New(&adaptiveGrid{maxCols: maxCols, cols: maxCols}, objs...)
}

func (g *adaptiveGrid) spacing() float32 {
	if g.gap > 0 {
		return g.gap
	}
	return cardGap
}

func visible(objs []fyne.CanvasObject) []fyne.CanvasObject {
	out := make([]fyne.CanvasObject, 0, len(objs))
	for _, o := range objs {
		if o.Visible() {
			out = append(out, o)
		}
	}
	return out
}

func cellMin(objs []fyne.CanvasObject) fyne.Size {
	var m fyne.Size
	for _, o := range objs {
		m = m.Max(o.MinSize())
	}
	return m
}

func (g *adaptiveGrid) columnsFor(width float32, cell fyne.Size, n int) int {
	cols := g.maxCols
	for cols > 1 && float32(cols)*cell.Width+float32(cols-1)*g.spacing() > width {
		cols--
	}
	return max(1, min(cols, n))
}

// rowHeights is each row's tallest card at cols columns.
func rowHeights(objs []fyne.CanvasObject, cols int) []float32 {
	var hs []float32
	for i, o := range objs {
		if i%cols == 0 {
			hs = append(hs, 0)
		}
		hs[len(hs)-1] = max(hs[len(hs)-1], o.MinSize().Height)
	}
	return hs
}

func (g *adaptiveGrid) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	objs = visible(objs)
	if len(objs) == 0 {
		return
	}
	g.cols = g.columnsFor(size.Width, cellMin(objs), len(objs))
	gap := g.spacing()
	w := (size.Width - float32(g.cols-1)*gap) / float32(g.cols)
	y := float32(0)
	for r, h := range rowHeights(objs, g.cols) {
		for c := 0; c < g.cols; c++ {
			i := r*g.cols + c
			if i >= len(objs) {
				break
			}
			objs[i].Move(fyne.NewPos(float32(c)*(w+gap), y))
			objs[i].Resize(fyne.NewSize(w, h))
		}
		y += h + gap
	}
}

func (g *adaptiveGrid) MinSize(objs []fyne.CanvasObject) fyne.Size {
	objs = visible(objs)
	if len(objs) == 0 {
		return fyne.Size{}
	}
	cols := max(1, min(g.cols, len(objs)))
	hs := rowHeights(objs, cols)
	h := float32(len(hs)-1) * g.spacing()
	for _, rh := range hs {
		h += rh
	}
	return fyne.NewSize(cellMin(objs).Width, h)
}

// gapVBox stacks objects vertically with cardGap between them.
type gapVBox struct{}

func newGapVBox(objs ...fyne.CanvasObject) *fyne.Container {
	return container.New(gapVBox{}, objs...)
}

func (gapVBox) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for _, o := range visible(objs) {
		h := o.MinSize().Height
		o.Move(fyne.NewPos(0, y))
		o.Resize(fyne.NewSize(size.Width, h))
		y += h + cardGap
	}
}

func (gapVBox) MinSize(objs []fyne.CanvasObject) fyne.Size {
	var s fyne.Size
	objs = visible(objs)
	for i, o := range objs {
		m := o.MinSize()
		s.Width = max(s.Width, m.Width)
		s.Height += m.Height
		if i > 0 {
			s.Height += cardGap
		}
	}
	return s
}

// checkGrid is a set of checkboxes, all checked at first, that wraps into rows on narrow windows
// (a horizontal CheckGroup never wraps and sets its container's minimum width).
type checkGrid struct {
	root   *fyne.Container
	checks []*handCheck
}

func newCheckGrid(names []string) *checkGrid {
	g := &checkGrid{}
	objs := make([]fyne.CanvasObject, len(names))
	for i, n := range names {
		c := newCheck(n, nil)
		c.Checked = true
		g.checks = append(g.checks, c)
		objs[i] = c
	}
	g.root = container.New(&adaptiveGrid{maxCols: 6, cols: 6, gap: theme.Padding()}, objs...)
	return g
}

func (g *checkGrid) selected() []string {
	var out []string
	for _, c := range g.checks {
		if c.Checked {
			out = append(out, c.Text)
		}
	}
	return out
}
