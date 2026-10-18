//go:build gui

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// column is one dataTable column.
type column struct {
	title string
	width float32
}

// cell is what a dataTable row returns for a column: text and an optional emphasis.
type cell struct {
	text       string
	importance widget.Importance
	bold       bool
}

// dataTable is a header-row table over a row count and a cell function, with whole-row selection:
// fyne selects a single cell, so the selection is kept here and painted behind every cell of the row.
type dataTable struct {
	*widget.Table
	cols     []column
	rows     func() int
	cell     func(row, col int) cell
	selected int
	onSelect func(row int)
}

func newDataTable(cols []column, rows func() int, cellFn func(row, col int) cell, onSelect func(row int)) *dataTable {
	d := &dataTable{cols: cols, rows: rows, cell: cellFn, selected: -1, onSelect: onSelect}
	d.Table = widget.NewTableWithHeaders(
		func() (int, int) { return d.rows(), len(d.cols) },
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.Truncation = fyne.TextTruncateEllipsis
			return container.New(bleedLayout{}, canvas.NewRectangle(color.Transparent), l)
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			bg, l := box.Objects[0].(*canvas.Rectangle), box.Objects[1].(*widget.Label)
			fill := color.Color(color.Transparent)
			if id.Row == d.selected {
				fill = themeColor(theme.ColorNameSelection)
			}
			if bg.FillColor != fill {
				bg.FillColor = fill
				bg.Refresh()
			}
			if id.Row >= d.rows() {
				l.SetText("")
				return
			}
			c := d.cell(id.Row, id.Col)
			l.Importance = c.importance
			l.TextStyle = fyne.TextStyle{Bold: c.bold}
			l.SetText(c.text)
		})
	d.ShowHeaderColumn = false
	// Per-cell separators would split the selection band into one block per column.
	d.HideSeparators = true
	d.CreateHeader = func() fyne.CanvasObject {
		l := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		l.Importance = widget.LowImportance
		return l
	}
	d.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(d.cols) {
			o.(*widget.Label).SetText(d.cols[id.Col].title)
		}
	}
	for i, c := range cols {
		d.SetColumnWidth(i, c.width)
	}
	d.OnSelected = func(id widget.TableCellID) {
		d.selected = id.Row
		// Drop fyne's single-cell highlight: the row background marks the selection.
		d.UnselectAll()
		d.Refresh()
		if d.onSelect != nil {
			d.onSelect(id.Row)
		}
	}
	return d
}

// reset clears the selection after the rows were replaced.
func (d *dataTable) reset() {
	d.selected = -1
	d.UnselectAll()
	d.Refresh()
}

// texts is the table as header titles and cell texts, for export.
func (d *dataTable) texts() (headers []string, rows [][]string) {
	for _, c := range d.cols {
		headers = append(headers, c.title)
	}
	n := d.rows()
	rows = make([][]string, n)
	for r := range n {
		rows[r] = make([]string, len(d.cols))
		for c := range d.cols {
			rows[r][c] = d.cell(r, c).text
		}
	}
	return headers, rows
}

// bleedLayout stretches its first object (the selection background) over the padding fyne's Table
// leaves between cells, so a selected row reads as one band; the rest fill the cell.
type bleedLayout struct{}

func (bleedLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	pad := theme.Padding()
	for i, o := range objs {
		if i == 0 {
			o.Move(fyne.NewPos(-pad/2, -pad/2))
			o.Resize(size.AddWidthHeight(pad, pad))
			continue
		}
		o.Move(fyne.NewPos(0, 0))
		o.Resize(size)
	}
}

func (bleedLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	var m fyne.Size
	for _, o := range objs[1:] {
		m = m.Max(o.MinSize())
	}
	return m
}
