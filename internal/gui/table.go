//go:build gui

package gui

import (
	"fyne.io/fyne/v2"
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

// dataTable is a header-row table over a row count and a cell function, with single-row selection.
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
			return l
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
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
