package ui

// TableStyle decorates native table rows and cells without changing their
// measured size, column identities, selection, or scrolling behavior.
// Its zero value leaves the table unchanged.
type TableStyle struct {
	// AlternateRowColor fills odd result rows; transparent disables striping.
	AlternateRowColor Color
	// GridColor colors the enabled grid lines; transparent draws no lines.
	GridColor             Color
	RowLines, ColumnLines bool
	// RowOffset keeps striping continuous when a table displays one page.
	RowOffset int64
}

func tableGridPaint(style TableStyle, columnHeader, rtl bool) func(*Painter, Rect) {
	return func(p *Painter, r Rect) {
		bottom := float32(0)
		if style.RowLines && !columnHeader {
			bottom = min(float32(1), r.H)
		}
		if style.ColumnLines {
			x := r.X + r.W - 1
			if rtl {
				x = r.X
			}
			// The bottom edge owns the corner too, so translucent lines
			// do not become darker where the two edges meet.
			p.Fill(Rect{x, r.Y, 1, max(float32(0), r.H-bottom)}, style.GridColor, 0)
		}
		if style.RowLines && !columnHeader {
			p.Fill(Rect{r.X, r.Y + r.H - 1, r.W, 1}, style.GridColor, 0)
		}
	}
}

func decorateTableCell(cell *node, style *TableStyle, columnHeader bool) {
	if style == nil || style.GridColor.A == 0 || (!style.RowLines && !style.ColumnLines) {
		return
	}
	grid := tableGridPaint(*style, columnHeader, cell.rtl())
	if columnHeader {
		// Resize handles remain above the normal grid when hovered or dragged.
		cell.Draw(grid)
		return
	}
	// Cursor painting is installed by native selection. Compose rather than
	// replacing it: the focus outline must stay above the ordinary grid.
	cursor := cell.paintAfterFn
	cell.DrawOver(func(p *Painter, r Rect) {
		grid(p, r)
		if cursor != nil {
			cursor(p, r)
		}
	})
}
