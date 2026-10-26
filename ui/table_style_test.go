package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"testing"
)

// Inspect actual cell rectangles and raster output, including padding and
// empty cells, rather than assuming a row height or text width.
func gridFixture(s *ListState, cols *[]TableColumn, count int) (*Tester, map[[2]int]*node) {
	cells := map[[2]int]*node{}
	tt := NewTester(func(c *Context) {
		clear(cells)
		Table(c, s, *cols, count, func(row, col int) {
			text := fmt.Sprintf("r%dc%d", row, col)
			if col == 2 {
				text = "" // Empty text still has full cell boundaries.
			}
			cells[[2]int{row, col}] = Text(c, text).SingleLine().node().parent
		}).Grow(1)
	}, 400, 300)
	return tt, cells
}

func gridTestColumns() []TableColumn {
	return []TableColumn{{Title: "#", Width: 40, Fixed: true, RowHeader: true}, {Title: "A", Width: 180}, {Title: "B", Width: 180}, {Title: "C", Width: 180}}
}

func gridRGBA(c Color) color.RGBA { return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A} }

func TestTableStyleZeroPreservesRenderingAndLayout(t *testing.T) {
	s := &ListState{}
	cols := gridTestColumns()
	tt, cells := gridFixture(s, &cols, 30)
	before := bytes.Clone(tt.Image().Pix)
	old := header(t, tt, "A")
	s.TableStyle = &TableStyle{}
	tt.Frame()
	if !bytes.Equal(before, tt.Image().Pix) || header(t, tt, "A") != old {
		t.Fatal("zero style changed the default rendering or header layout")
	}
	cell := cells[[2]int{1, 1}].st
	w, h := cell.w, cell.h
	s.TableStyle = &TableStyle{AlternateRowColor: RGB(70, 80, 90), GridColor: RGB(230, 10, 20), RowLines: true, ColumnLines: true}
	tt.Frame()
	cell = cells[[2]int{1, 1}].st
	if header(t, tt, "A") != old || cell.w != w || cell.h != h {
		t.Fatal("decorations participated in table measurement")
	}
}

func TestTableStyleParityAndVirtualScrolling(t *testing.T) {
	stripe := RGB(70, 80, 90)
	s := &ListState{TableStyle: &TableStyle{AlternateRowColor: stripe}}
	cols := gridTestColumns()
	tt, cells := gridFixture(s, &cols, 100000)
	check := func() {
		t.Helper()
		img := tt.Image()
		if len(cells) > 200 {
			t.Fatalf("style disabled virtualization: %d cells built", len(cells))
		}
		checked := 0
		for pos, n := range cells {
			r := n.st
			if r.x < 0 || r.x+3 >= 390 || r.y < 35 || r.y+r.h > 290 {
				continue
			}
			got := img.RGBAAt(int(r.x)+2, int(r.y)+2)
			odd := (int64(pos[0])+s.TableStyle.RowOffset)%2 != 0
			if (got == gridRGBA(stripe)) != odd {
				t.Fatalf("cell %v offset %d: stripe=%v, pixel=%v", pos, s.TableStyle.RowOffset, odd, got)
			}
			checked++
		}
		if checked < 3 {
			t.Fatal("no visible cells checked")
		}
	}
	check()
	s.TableStyle.RowOffset = 101
	tt.Frame()
	check()
	tt.Scroll(200, 150, 0, 800)
	tt.Move(399, 299) // Keep hover out of sampled cells.
	if s.first == 0 {
		t.Fatal("vertical scroll did not move rows")
	}
	check()
	tt.Scroll(200, 150, 230, 0)
	tt.Move(399, 299)
	check()
}

func TestTableStyleLinesFollowColumns(t *testing.T) {
	grid := RGB(230, 10, 20)
	s := &ListState{TableStyle: &TableStyle{GridColor: grid, RowLines: true, ColumnLines: true}}
	cols := gridTestColumns()
	tt, cells := gridFixture(s, &cols, 40)
	check := func(col int, title string) {
		t.Helper()
		n := cells[[2]int{1, col}].st
		h := header(t, tt, title)
		if float64(n.x) != h.X || float64(n.w) != h.W {
			t.Fatal("header and cell diverged")
		}
		img := tt.Image()
		for _, xy := range [][2]int{{int(n.x+n.w) - 1, int(n.y) + 3}, {int(n.x) + 3, int(n.y+n.h) - 1}, {int(h.X+h.W) - 1, int(h.Y) + 3}} {
			if got := img.RGBAAt(xy[0], xy[1]); got != gridRGBA(grid) {
				t.Fatalf("grid at %v: %v", xy, got)
			}
		}
	}
	check(1, "A")
	h := header(t, tt, "A")
	x, y := float32(h.X+h.W)-2, float32(h.Y+h.H/2)
	tt.Press(x, y)
	tt.Move(x+25, y)
	tt.Release(x+25, y)
	tt.Move(399, 299)
	check(1, "A")
	s.Columns.Order = []string{"#", "B", "A", "C"}
	tt.Frame()
	check(2, "B")
	cols = []TableColumn{cols[0], cols[1], cols[2]}
	tt.Frame()
	check(2, "B")
	tt.Scroll(200, 150, 100, 0)
	tt.Move(399, 299)
	check(1, "A")
}

func TestTableStyleIndependentEdgesAndSinglePaintCorners(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		for _, lines := range [][2]bool{{true, false}, {false, true}, {true, true}} {
			style := TableStyle{GridColor: RGBA(255, 0, 0, 0.5), RowLines: lines[0], ColumnLines: lines[1]}
			tt := NewTester(func(c *Context) {
				Box(c).Size(100, 100).Background(RGB(255, 255, 255)).DrawOver(tableGridPaint(style, false, rtl))
			}, 120, 120)
			img := tt.Image()
			x := 99
			if rtl {
				x = 0
			}
			bottom, side, corner := img.RGBAAt(50, 99), img.RGBAAt(x, 50), img.RGBAAt(x, 99)
			white := color.RGBA{255, 255, 255, 255}
			if (bottom != white) != lines[0] || (side != white) != lines[1] {
				t.Fatal("line switches are not independent")
			}
			if lines[0] && corner != bottom || !lines[0] && corner != side {
				t.Fatal("shared corner painted more than once")
			}
		}
	}
}

func TestTableStyleSelectionAndFocusStayAboveGrid(t *testing.T) {
	sel := &Selection[int]{}
	s := &ListState{CellSelect: true, Selection: sel, TableStyle: &TableStyle{AlternateRowColor: RGB(70, 80, 90), GridColor: RGB(230, 10, 20), RowLines: true, ColumnLines: true}}
	cols := gridTestColumns()
	tt, cells := gridFixture(s, &cols, 20)
	tt.Click("r1c1")
	n := cells[[2]int{1, 1}]
	if n.bg != tt.rt.c.theme.Accent || n.paintAfterFn == nil {
		t.Fatal("selection or cursor lost")
	}
	img := tt.Image()
	// The selected cursor's bottom edge is AccentText, above the red grid.
	if got := img.RGBAAt(int(n.st.x)+15, int(n.st.y+n.st.h)-1); got == gridRGBA(s.TableStyle.GridColor) {
		t.Fatal("grid obscures cursor border")
	}
	tt.ClickWith(Shift, "r3c1")
	if len(chosenRows(sel)) != 3 {
		t.Fatal("block selection changed")
	}
	tt.Click("r1c0")
	if !s.Cells.All || cells[[2]int{1, 0}].bg != tt.rt.c.theme.Accent {
		t.Fatal("row header selection lost")
	}
	tt.Key(0, KeyDown)
	if *s.cursor() != 2 {
		t.Fatal("keyboard navigation changed")
	}
}
