package ui

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/internal/platform"
)

// Titles deliberately differ from IDs: changing display order must not turn
// the header choice into a choice of a different source field.
type headerGridFixture struct {
	state     ListState
	selection Selection[int]
	selected  int
	rows      int
	columns   []TableColumn
	sort      SortOrder
}

func newHeaderGridFixture(rows int) *headerGridFixture {
	f := &headerGridFixture{rows: rows, selected: -1, sort: SortOrder{Column: "gamma", Descending: true}}
	f.columns = []TableColumn{
		{Title: "#", ID: "row", RowHeader: true, Fixed: true, Width: 35},
		{Title: "Alpha", ID: "alpha", Width: 100, MinWidth: 60, Sortable: true},
		{Title: "Beta", ID: "beta", Width: 100, MinWidth: 60, Sortable: true},
		{Title: "Gamma", ID: "gamma", Width: 100, MinWidth: 60, Sortable: true},
		{Title: "Delta", ID: "delta", Width: 100, MinWidth: 60, Sortable: true},
	}
	f.state = ListState{HeaderSelect: true, CellSelect: true, DragSelect: true, Selection: &f.selection, Selected: &f.selected, Sort: &f.sort}
	return f
}

func (f *headerGridFixture) view(c *Context) {
	Column(c).Fill().Padding(10).Children(func() {
		Table(c, &f.state, f.columns, f.rows, func(row, col int) {
			Text(c, fmt.Sprintf("cell %d %s", row, f.columns[col].ID)).SingleLine()
		}).Grow(1)
	})
}

func headerGridClick(t *testing.T, tt *Tester, mods Modifiers, title string) {
	t.Helper()
	if err := tt.ClickWith(mods, title); err != nil {
		t.Fatal(err)
	}
}

func headerGridColumns(t *testing.T, f *headerGridFixture, want ...string) {
	t.Helper()
	if got := f.state.SelectedColumns(); !slices.Equal(got, want) {
		t.Fatalf("selected columns %q, want %q (cells %+v, header columns %q)", got, want, f.state.Cells, f.state.ColumnSelection)
	}
}

func headerGridActionNode(t *testing.T, tt *Tester, title string) *node {
	t.Helper()
	var find func(*node) *node
	find = func(n *node) *node {
		if n.label == title+" actions" {
			return n
		}
		for child := n.first; child != nil; child = child.next {
			if got := find(child); got != nil {
				return got
			}
		}
		return nil
	}
	if got := find(tt.rt.c.root); got != nil {
		return got
	}
	t.Fatalf("no actions button for %s", title)
	return nil
}

func TestGridHeaderChoosesAllLoadedRowsWithoutSorting(t *testing.T) {
	f := newHeaderGridFixture(6)
	tt := NewTester(f.view, 480, 330)
	tt.Frame()
	before := f.sort
	headerGridClick(t, tt, 0, "Beta")
	headerGridColumns(t, f, "beta")
	if got := chosenRows(&f.selection); !slices.Equal(got, []int{0, 1, 2, 3, 4, 5}) || f.selected != 0 {
		t.Fatalf("header chose rows %v and cursor %d; want every loaded row with cursor 0", got, f.selected)
	}
	if f.sort != before || f.state.Cells.Cursor != "beta" || f.state.ColumnSelection == nil {
		t.Fatalf("header choice changed sort or lost field identity: sort %+v, cells %+v, header choice %q", f.sort, f.state.Cells, f.state.ColumnSelection)
	}
	for row := range f.rows {
		for col := range f.columns {
			if got := f.state.CellSelected(row, col); got != (col == 2) {
				t.Fatalf("row %d source column %d selected %v; want only beta", row, col, got)
			}
		}
	}
	headerGridClick(t, tt, 0, "cell 3 alpha")
	headerGridClick(t, tt, 0, "Gamma")
	if f.selected != 3 {
		t.Fatalf("header replaced valid active row with %d", f.selected)
	}
}

func TestGridHeaderCmdAddsAndRemovesNonadjacentColumns(t *testing.T) {
	f := newHeaderGridFixture(5)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	headerGridColumns(t, f, "alpha", "gamma")
	if f.state.CellSelected(2, 2) || !f.state.CellSelected(2, 1) || !f.state.CellSelected(2, 3) {
		t.Fatal("nonadjacent header choice included the intervening column or omitted a selected field")
	}
	headerGridClick(t, tt, Cmd, "Gamma")
	headerGridColumns(t, f, "alpha")
	headerGridClick(t, tt, Cmd, "Alpha")
	headerGridColumns(t, f, "alpha")
	if f.selection.Len() != f.rows {
		t.Fatal("toggling a header altered the selected loaded rows")
	}
}

func TestGridHeaderDiscreteChoiceDrawingMatchesSelectedCells(t *testing.T) {
	f := newHeaderGridFixture(4)
	cells := make(map[[2]int]*node)
	var accent Color
	tt := NewTester(func(c *Context) {
		accent = c.Theme().Accent
		clear(cells)
		Table(c, &f.state, f.columns, f.rows, func(row, col int) {
			cells[[2]int{row, col}] = Text(c, fmt.Sprintf("cell %d %s", row, f.columns[col].ID)).node().parent
		}).Grow(1)
	}, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	for rc, n := range cells {
		if got, drawn := f.state.CellSelected(rc[0], rc[1]), n.bg == accent; got != drawn {
			t.Fatalf("selected-cell query %v differs from its native background %v at row/source %v", got, drawn, rc)
		}
	}
}

func TestGridHeaderShiftKeepsIndependentColumnsAndShrinksItsRange(t *testing.T) {
	f := newHeaderGridFixture(5)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Delta")
	headerGridClick(t, tt, Shift, "Beta")
	headerGridColumns(t, f, "alpha", "beta", "gamma", "delta")
	headerGridClick(t, tt, Shift, "Gamma")
	headerGridColumns(t, f, "alpha", "gamma", "delta")
	headerGridClick(t, tt, 0, "Beta")
	headerGridClick(t, tt, Shift, "Gamma")
	headerGridColumns(t, f, "beta", "gamma")
}

func TestGridHeaderShiftDoesNotReviveChoiceAfterApplicationReset(t *testing.T) {
	f := newHeaderGridFixture(5)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Delta")
	headerGridClick(t, tt, Shift, "Beta")
	headerGridColumns(t, f, "alpha", "beta", "gamma", "delta")
	// Page replacement and explicit editor activation assign Cells directly.
	// Clearing whole-column mode must revoke its previous Shift history too.
	f.state.ColumnSelection = nil
	f.state.Cells = ColumnRange{Cursor: "beta", Anchor: "beta"}
	f.selection.Clear()
	f.selection.Add(1)
	f.selected = 1
	tt.Frame()
	headerGridClick(t, tt, Shift, "Gamma")
	headerGridColumns(t, f, "gamma")
}

func TestGridHeaderShiftDoesNotReviveHiddenBaseOrPivot(t *testing.T) {
	for _, removed := range []string{"alpha", "delta"} {
		t.Run(removed, func(t *testing.T) {
			f := newHeaderGridFixture(5)
			original := slices.Clone(f.columns)
			tt := NewTester(f.view, 480, 300)
			tt.Frame()
			headerGridClick(t, tt, 0, "Alpha")
			headerGridClick(t, tt, Cmd, "Delta")
			headerGridClick(t, tt, Shift, "Beta")
			index := slices.IndexFunc(f.columns, func(col TableColumn) bool { return col.ID == removed })
			f.columns = slices.Delete(f.columns, index, index+1)
			tt.Frame()
			// Showing the field again does not restore selection history which
			// belonged to a field absent from the intervening frame.
			f.columns = original
			tt.Frame()
			headerGridClick(t, tt, Shift, "Gamma")
			if removed == "alpha" {
				headerGridColumns(t, f, "gamma", "delta")
			} else {
				headerGridColumns(t, f, "gamma")
			}
		})
	}
}

func TestGridHeaderContextMenuKeepsInsideChoiceAndSelectsOutside(t *testing.T) {
	f := newHeaderGridFixture(5)
	built, chosen := -1, -1
	f.state.HeaderMenu = func(m *Menu, col int) {
		built = col
		if m.Item("Header copy").Chosen() {
			chosen = col
		}
	}
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	if err := tt.RightClick("Gamma"); err != nil {
		t.Fatal(err)
	}
	headerGridColumns(t, f, "alpha", "gamma")
	if built != 3 || !slices.Equal(tt.Menu(), []string{"Header copy"}) {
		t.Fatalf("context menu %q built for %d, want original gamma index 3", tt.Menu(), built)
	}
	tt.CloseMenu()
	if err := tt.RightClick("Beta"); err != nil {
		t.Fatal(err)
	}
	headerGridColumns(t, f, "beta")
	if err := tt.ChooseMenuItem("Header copy"); err != nil {
		t.Fatal(err)
	}
	if chosen != 2 {
		t.Fatalf("context item acted on source column %d, want 2", chosen)
	}
}

func TestGridHeaderActionMenuDoesNotSelectItsColumn(t *testing.T) {
	f := newHeaderGridFixture(5)
	built, chosen := -1, -1
	f.state.HeaderMenu = func(m *Menu, col int) { m.Item("Context copy") }
	f.state.HeaderActionMenu = func(m *Menu, col int) {
		built = col
		if m.Item("Sort ascending").Chosen() {
			chosen = col
		}
	}
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	beta := header(t, tt, "Beta")
	if button := headerGridActionNode(t, tt, "Beta"); !button.opacitySet || button.opacity != 0 {
		t.Fatal("unhovered header's action button was visible")
	}
	tt.Move(float32(beta.X+beta.W/2), float32(beta.Y+beta.H/2))
	if button := headerGridActionNode(t, tt, "Beta"); button.opacitySet && button.opacity == 0 {
		t.Fatal("header actions did not appear in the hover frame")
	}
	if after := header(t, tt, "Beta"); beta.W != after.W {
		t.Fatalf("revealing header actions changed column width: before %v after %v", beta, after)
	}
	headerGridClick(t, tt, 0, "Beta actions")
	if !slices.Equal(tt.Menu(), []string{"Sort ascending"}) || built != 2 {
		t.Fatalf("action menu %q built for %d; want independent menu for original beta index", tt.Menu(), built)
	}
	headerGridColumns(t, f, "alpha", "gamma")
	if err := tt.ChooseMenuItem("Sort ascending"); err != nil {
		t.Fatal(err)
	}
	if chosen != 2 {
		t.Fatalf("action item acted on %d, want original beta index", chosen)
	}
	headerGridColumns(t, f, "alpha", "gamma")
}

func TestGridHeaderActionMenuKeyboardReachesHiddenButton(t *testing.T) {
	f := newHeaderGridFixture(5)
	built := -1
	f.state.HeaderActionMenu = func(m *Menu, col int) {
		built = col
		m.Item("Keyboard action")
	}
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	tt.Move(470, 290)
	for count := 0; count < 12 && !headerGridActionNode(t, tt, "Beta").Focused(); count++ {
		tt.Key(0, KeyTab)
	}
	if !headerGridActionNode(t, tt, "Beta").Focused() {
		t.Fatal("keyboard navigation could not reach the hover-hidden header action")
	}
	// Tab applies focus after construction; render the view with that focus
	// before checking the button's presentation.
	tt.Frame()
	if button := headerGridActionNode(t, tt, "Beta"); button.opacitySet && button.opacity == 0 {
		t.Fatal("keyboard-focused header actions stayed hidden")
	}
	tt.Key(0, KeyEnter)
	if built != 2 || !slices.Equal(tt.Menu(), []string{"Keyboard action"}) {
		t.Fatalf("keyboard action menu %q for column %d, want beta's independent menu", tt.Menu(), built)
	}
	headerGridColumns(t, f, "alpha")
}

func TestGridHeaderChoiceFollowsFieldIDsAfterReorderAndRemoval(t *testing.T) {
	f := newHeaderGridFixture(5)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	f.state.Columns.Order = []string{"delta", "gamma", "beta", "alpha"}
	tt.Frame()
	headerGridColumns(t, f, "gamma", "alpha")
	if !f.state.CellSelected(2, 1) || !f.state.CellSelected(2, 3) || f.state.CellSelected(2, 2) {
		t.Fatal("reordering changed the source fields of the selected cells")
	}
	f.columns = slices.Delete(f.columns, 3, 4)
	tt.Frame()
	headerGridColumns(t, f, "alpha")
	if f.state.Cells.Cursor != "alpha" {
		t.Fatalf("removed cursor field was not reconciled to a remaining selected field: %+v", f.state.Cells)
	}
}

func TestGridHeaderEmptyToLoadedAndBackKeepsNativeInteraction(t *testing.T) {
	f := newHeaderGridFixture(0)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Beta")
	headerGridColumns(t, f, "beta")
	if f.selected != -1 || f.selection.Len() != 0 || f.state.CellSelected(0, 2) {
		t.Fatal("choosing an empty-table header invented a loaded row")
	}
	f.rows = 4
	tt.Frame()
	headerGridClick(t, tt, 0, "Gamma")
	headerGridColumns(t, f, "gamma")
	if f.selection.Len() != 4 {
		t.Fatal("loaded rows did not use the same native header choice")
	}
	f.rows = 0
	tt.Frame()
	f.state.HeaderMenu = func(m *Menu, col int) { m.Item("Empty header") }
	tt.Frame()
	if err := tt.RightClick("Alpha"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tt.Menu(), []string{"Empty header"}) || f.selected != -1 {
		t.Fatalf("empty header interaction lost after loaded rows: menu %q, cursor %d", tt.Menu(), f.selected)
	}
}

func TestGridHeaderEmptyStateCanMoveToAnotherTester(t *testing.T) {
	f := newHeaderGridFixture(0)
	first := NewTester(f.view, 480, 300)
	first.Frame()
	headerGridClick(t, first, 0, "Alpha")
	// Each new window begins with the same frame/pass counters. Moving an
	// empty table's state to its new context must not look like two lists
	// registered during one pass in the old window.
	second := NewTester(f.view, 480, 300)
	second.Frame()
	headerGridClick(t, second, 0, "Beta")
	headerGridColumns(t, f, "beta")
	if f.selection.Len() != 0 || f.selected != -1 {
		t.Fatal("moving an empty native table created a data row")
	}
}

func TestGridHeaderStateStillRejectsTwoTablesInOneContext(t *testing.T) {
	f := newHeaderGridFixture(0)
	defer func() {
		if recover() == nil {
			t.Fatal("one ListState accepted two native tables in one context pass")
		}
	}()
	tt := NewTester(func(c *Context) {
		Table(c.Key("first"), &f.state, f.columns, 0, func(int, int) {}).Grow(1)
		Table(c.Key("second"), &f.state, f.columns, 0, func(int, int) {}).Grow(1)
	}, 480, 300)
	tt.Frame()
}

func TestGridHeaderChoiceExcludesStatisticalSection(t *testing.T) {
	f := newHeaderGridFixture(6)
	f.state.Header = func(row int) bool { return row == 0 }
	tt := NewTester(f.view, 480, 330)
	tt.Frame()
	headerGridClick(t, tt, 0, "Gamma")
	if got := chosenRows(&f.selection); !slices.Equal(got, []int{1, 2, 3, 4, 5}) || f.selected != 1 {
		t.Fatalf("whole-column choice included a section or selected no real row: rows %v cursor %d", got, f.selected)
	}
	if f.state.CellSelected(0, 3) || !f.state.CellSelected(1, 3) {
		t.Fatal("section-header selection was not excluded from cells")
	}
}

func TestGridHeaderArrowReturnsToOneCell(t *testing.T) {
	f := newHeaderGridFixture(5)
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Alpha")
	headerGridClick(t, tt, Cmd, "Gamma")
	tt.Key(0, KeyRight)
	if f.state.ColumnSelection != nil || f.selection.Len() != 1 || f.state.Cells.All {
		t.Fatalf("arrow did not leave one cell: rows %v, cells %+v, header choice %q", chosenRows(&f.selection), f.state.Cells, f.state.ColumnSelection)
	}
	headerGridColumns(t, f, "delta")
}

func TestGridHeaderTooltipUsesDelayAndCancelsOnLeaveAndPress(t *testing.T) {
	f := newHeaderGridFixture(4)
	f.columns[2].HeaderTooltip = "Chinese comment 中文\nSecond line"
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	beta := header(t, tt, "Beta")
	x, y := float32(beta.X+20), float32(beta.Y+beta.H/2)
	tt.Move(x, y)
	tt.Advance(590 * time.Millisecond)
	if tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatal("header tooltip appeared before 600ms")
	}
	tt.Advance(20 * time.Millisecond)
	if !tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatalf("nonempty header comment did not appear after the delay: %q", tt.Texts())
	}
	tt.Move(470, 290)
	tt.Frame()
	if tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatal("header tooltip remained after the pointer left")
	}
	tt.Advance(time.Second)
	tt.Move(x, y)
	tt.Advance(610 * time.Millisecond)
	tt.Press(x, y)
	if tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatal("pressing the header kept its tooltip open")
	}
	tt.Release(x, y)
	tt.Advance(time.Second)
	if tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatal("pressed tooltip reappeared before leaving the header")
	}
	tt.Move(470, 290)
	tt.Advance(time.Second)
	alpha := header(t, tt, "Alpha")
	tt.Move(float32(alpha.X+20), float32(alpha.Y+alpha.H/2))
	tt.Advance(time.Second)
	if tt.HasText(f.columns[2].HeaderTooltip) {
		t.Fatal("empty header comment displayed another field's tooltip")
	}
}

func TestGridHeaderTooltipBubbleContainsSingleAndMultilineComments(t *testing.T) {
	for _, comment := range []string{"单行字段注释", "中文字段注释\nSecond line"} {
		t.Run(comment, func(t *testing.T) {
			f := newHeaderGridFixture(4)
			f.columns[2].HeaderTooltip = comment
			tt := NewTester(func(c *Context) {
				theme := *c.Theme()
				theme.FontSize, theme.Spacing = 12, 3.4
				c.SetTheme(&theme)
				Column(c).Fill().Padding(10).Children(func() {
					Text(c, "Grid toolbar").SingleLine()
					Table(c, &f.state, f.columns, f.rows, func(row, col int) {
						Text(c, fmt.Sprintf("cell %d %s", row, f.columns[col].ID)).SingleLine()
					}).Grow(1)
				})
			}, 480, 300)
			tt.Frame()
			beta := header(t, tt, "Beta")
			tt.Move(float32(beta.X+20), float32(beta.Y+beta.H/2))
			tt.Advance(610 * time.Millisecond)
			tt.Frame()
			var tip *node
			var find func(*node)
			find = func(n *node) {
				if n.role == RoleTooltip {
					tip = n
				}
				for child := n.first; child != nil; child = child.next {
					find(child)
				}
			}
			find(tt.rt.c.root)
			if tip == nil {
				t.Fatal("header comment has no native tooltip bubble")
			}
			bubble := tip.Bounds()
			text, ok := tt.Find(comment)
			if !ok || bubble.W <= 0 || bubble.H <= 0 {
				t.Fatalf("comment tooltip has no measured bubble or text: bubble %+v text %+v found %v", bubble, text, ok)
			}
			if text.X < bubble.X-0.5 || text.Y < bubble.Y-0.5 || text.X+text.W > bubble.X+bubble.W+0.5 || text.Y+text.H > bubble.Y+bubble.H+0.5 {
				t.Fatalf("header comment paints outside its native tooltip background: bubble %+v text %+v", bubble, text)
			}
		})
	}
}

func TestGridHeaderDetailAndWidthRemainNative(t *testing.T) {
	f := newHeaderGridFixture(3)
	f.columns[2].Subtitle = "varchar(255)"
	f.columns[2].HeaderDetail = "Comment detail"
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	if !tt.HasText("Comment detail") || !tt.HasText("varchar(255)") {
		t.Fatalf("native header lost detail or subtitle: %q", tt.Texts())
	}
	beta := header(t, tt, "Beta")
	x, y := float32(beta.X+beta.W-2), float32(beta.Y+beta.H/2)
	tt.Press(x, y)
	tt.Move(x+35, y)
	tt.Release(x+35, y)
	if got := f.state.Columns.Widths["beta"]; got != 135 {
		t.Fatalf("three-line header did not resize its original field: width %v", got)
	}
	if f.state.ColumnSelection != nil {
		t.Fatal("resizing selected a column")
	}
}

func TestGridHeaderDetailScalesWithoutChangingFieldChoice(t *testing.T) {
	f := newHeaderGridFixture(3)
	f.columns[2].Subtitle = "varchar(255)"
	f.columns[2].HeaderDetail = "Comment detail"
	zoom := float32(1)
	tt := NewTester(func(c *Context) {
		theme := *c.Theme()
		theme.Spacing, theme.FontSize = theme.Space(1)*zoom, theme.FontSize*zoom
		c.SetTheme(&theme)
		f.view(c)
	}, 500, 400)
	tt.Frame()
	headerGridClick(t, tt, 0, "Beta")
	before := header(t, tt, "Beta")
	zoom = 2
	tt.Frame()
	after := header(t, tt, "Beta")
	headerGridColumns(t, f, "beta")
	if after.H < before.H*1.5 {
		t.Fatalf("native multi-line header did not scale with its theme: before %v after %v", before, after)
	}
	if detail, ok := tt.Find("Comment detail"); !ok || detail.Y+detail.H > float32(after.Y+after.H)+0.5 {
		t.Fatalf("zoom clipped the third header line: detail %+v inside %+v", detail, after)
	}
}

func TestGridHeaderDisabledOptInKeepsSorting(t *testing.T) {
	f := newHeaderGridFixture(4)
	f.state.HeaderSelect = false
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "Beta")
	if f.sort != (SortOrder{Column: "beta"}) || f.state.ColumnSelection != nil || f.selection.Len() != 0 {
		t.Fatalf("default table behavior changed: sort %+v, header choice %q, rows %v", f.sort, f.state.ColumnSelection, chosenRows(&f.selection))
	}
}

func TestGridCellClickCallbackDoesNotRunForDragOrRowHeader(t *testing.T) {
	f := newHeaderGridFixture(100)
	type click struct {
		row, col int
		mods     Modifiers
	}
	var clicks []click
	f.state.OnCellClick = func(row, col int, mods Modifiers) { clicks = append(clicks, click{row, col, mods}) }
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	headerGridClick(t, tt, 0, "cell 1 beta")
	if !slices.Equal(clicks, []click{{1, 2, 0}}) {
		t.Fatalf("plain cell click callbacks %+v; want original source row and column", clicks)
	}
	clicks = nil
	x, y := rowCenter(t, tt, "cell 1 beta")
	x2, y2 := rowCenter(t, tt, "cell 3 gamma")
	tt.Press(x, y)
	tt.Move(x2, y2)
	tt.Release(x2, y2)
	if len(clicks) != 0 {
		t.Fatalf("drag activated a cell editor: %+v", clicks)
	}
	headerGridClick(t, tt, 0, "cell 2 row")
	if len(clicks) != 0 {
		t.Fatalf("row-header click activated a cell editor: %+v", clicks)
	}
	x, y = rowCenter(t, tt, "cell 1 beta")
	tt.Press(x, y)
	tt.Move(x, y+4)
	tt.Move(x, 299)
	tt.Advance(2 * time.Second)
	if f.state.frame.e.st.scrollY <= 0 || len(chosenRows(&f.selection)) <= 10 {
		t.Fatalf("held drag did not scroll beyond initially visible rows: scroll=%v, rows=%v", f.state.frame.e.st.scrollY, chosenRows(&f.selection))
	}
	tt.Release(x, 299)
	if len(clicks) != 0 {
		t.Fatalf("edge drag activated a cell editor: %+v", clicks)
	}
}

func TestGridCellClickCallbackUsesClickedSourceAfterPointerMoves(t *testing.T) {
	f := newHeaderGridFixture(4)
	var got [][3]int
	f.state.OnCellClick = func(row, col int, mods Modifiers) {
		got = append(got, [3]int{row, col, int(mods)})
	}
	tt := NewTester(f.view, 480, 300)
	tt.Frame()
	x, y := rowCenter(t, tt, "cell 1 alpha")
	x2, _ := rowCenter(t, tt, "cell 1 gamma")
	tt.rt.event(platform.SurfaceEvent{Kind: platform.PointerDown, X: float64(x), Y: float64(y)})
	tt.rt.event(platform.SurfaceEvent{Kind: platform.PointerUp, X: float64(x), Y: float64(y)})
	tt.rt.event(platform.SurfaceEvent{Kind: platform.PointerMove, X: float64(x2), Y: float64(y)})
	tt.Frame()
	if !slices.Equal(got, [][3]int{{1, 1, 0}}) {
		t.Fatalf("cell activation followed the later pointer instead of the clicked source: %v", got)
	}
	got = nil
	if err := tt.RightClick("cell 2 beta"); err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("secondary click activated a cell editor: %v", got)
	}
	headerGridClick(t, tt, Cmd, "cell 3 gamma")
	if !slices.Equal(got, [][3]int{{3, 3, int(Cmd)}}) {
		t.Fatalf("modified primary click lost its source or modifiers: %v", got)
	}
}
