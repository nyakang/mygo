package ui

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

func TestGridCellCanceledPressDoesNotActivate(t *testing.T) {
	for _, drag := range []bool{false, true} {
		f := newHeaderGridFixture(8)
		clicks := 0
		f.state.OnCellClick = func(int, int, Modifiers) { clicks++ }
		tt := NewTester(f.view, 480, 300)
		tt.Frame()
		x, y := rowCenter(t, tt, "cell 1 beta")
		tt.Press(x, y)
		if drag {
			x2, y2 := rowCenter(t, tt, "cell 3 gamma")
			tt.Move(x2, y2)
		}
		tt.CancelPress()
		tt.Frame()
		tt.Release(x, y)
		if clicks != 0 {
			t.Fatalf("canceled press (drag %v) activated %d cells", drag, clicks)
		}
		headerGridClick(t, tt, 0, "cell 1 beta")
		if clicks != 1 {
			t.Fatalf("click after cancel (drag %v) activated %d cells", drag, clicks)
		}
	}
}

func TestGridModifierOnlyKeysPreserveChoiceAndRemainUntaken(t *testing.T) {
	f := newHeaderGridFixture(4)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Beta")
	before := f.state.Cells
	for _, code := range []platform.KeyCode{platform.CodeShiftLeft, platform.CodeControlRight, platform.CodeAltLeft, platform.CodeSuperRight} {
		for _, kind := range []platform.SurfaceEventKind{platform.RawKeyPressed, platform.RawKeyReleased} {
			if tt.rt.event(platform.SurfaceEvent{Kind: kind, Code: code}) {
				t.Fatalf("grid swallowed modifier-only event %v/%v", code, kind)
			}
			tt.Frame()
		}
	}
	headerGridColumns(t, f, "beta")
	if f.state.Cells != before || f.selected != 0 || f.selection.Len() != 4 {
		t.Fatalf("modifier-only event changed the grid: %+v, row %d", f.state.Cells, f.selected)
	}
}

func TestGridHeaderExtensionsAndCellsStyleCoexist(t *testing.T) {
	f := newHeaderGridFixture(6)
	f.state.CellsStyle = CellsStyle{Cell: RGB(30, 60, 90), Row: RGB(70, 80, 90), Stripe: RGB(100, 110, 120)}
	f.state.TableStyle = &TableStyle{AlternateRowColor: RGB(10, 20, 30), GridColor: RGB(40, 50, 60), RowLines: true, ColumnLines: true}
	var built []int
	rows := map[int]*node{}
	tt := NewTester(func(c *Context) {
		f.state.HeaderCell = func(cell Element, col int) {
			built = append(built, col)
			Text(c, "header-extension").SingleLine()
		}
		Column(c).Fill().Padding(10).Children(func() {
			Table(c, &f.state, f.columns, f.rows, func(row, col int) {
				text := Textf(c, "cell %d %s", row, f.columns[col].ID).SingleLine()
				rows[row] = text.node().parent.parent.parent
			}).Grow(1)
		})
	}, 480, 330)
	tt.Frame()
	headerGridClick(t, tt, 0, "Beta")
	headerGridClick(t, tt, Cmd, "Delta")
	for i := 0; i < f.rows; i++ {
		row := rows[i]
		if row == nil || row.bg != f.state.CellsStyle.Row {
			t.Fatalf("row %d lost its explicit CellsStyle background", i)
		}
	}
	f.state.Columns.Order = []string{"row", "delta", "gamma", "beta", "alpha"}
	built = nil
	tt.Frame()
	for start := 0; start < len(built); start += len(f.columns) {
		if len(built)-start < len(f.columns) || !slices.Equal(built[start:start+len(f.columns)], []int{0, 4, 3, 2, 1}) {
			t.Fatalf("HeaderCell source columns after reorder: %v", built)
		}
	}
	headerGridColumns(t, f, "delta", "beta")
	if !f.state.CellSelected(2, 4) || f.state.CellSelected(2, 3) {
		t.Fatal("CellsStyle painting and discrete column selection disagree")
	}
	headerGridClick(t, tt, 0, "cell 1 alpha")
	if rows[1].bg != f.state.CellsStyle.Row || rows[3].bg != f.state.CellsStyle.Stripe {
		t.Fatal("TableStyle replaced explicit selected-row or unselected-row colors")
	}
}
