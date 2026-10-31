# Opt-in native data-grid headers

The `v0.7.3-nyadb.2` release extends the existing keyed `ui.Table` header.
It does not change the zero-value behavior of Tables or OutlineTables.
No platform backend or browser bridge changes are required.

```go
var selected int = -1
var selection ui.Selection
grid := ui.ListState{
    Selected: &selected,
    Selection: &selection,
    CellSelect: true,
    DragSelect: true,
    HeaderSelect: true,
}
columns := []ui.TableColumn{
    {ID: "id", Title: "id", Width: 100, Sortable: true},
    {ID: "name", Title: "name", Width: 180,
        Subtitle: "varchar(200)", HeaderTooltip: "Customer name",
        HeaderDetail: "Customer name"},
}
grid.HeaderActionMenu = func(m *ui.Menu, col int) {
    // col is the original index into columns, even after a reorder.
    m.Item("Sort ascending")
}
ui.Table(c, &grid, columns, len(rows), func(row, col int) {
    ui.Text(c, rows[row][col])
})
```

`HeaderSelect` requires `CellSelect` and a `Selection`. A plain click chooses
one entire column, without changing `Sort` or `SortBy`. Ctrl (Cmd on macOS)
toggles columns, retaining the final chosen column. Shift extends from the
last plain or Ctrl/Cmd header click. Independent columns outside that range
remain chosen; a subsequent Shift click replaces the previous extension.
`ColumnRange` remains comparable. The independent ID set lives in
`ColumnSelection`; `SelectedColumns()` returns its members in visible order.
`CellSelected()` and native selection painting use the same membership.

Only current real row keys enter `Selection`; section headers are excluded.
A zero-row table can choose fields while its row selection stays empty and
its row cursor stays -1. Clicking cells or using navigation keys returns to
the ordinary cell range. When an application assigns `Cells` directly or
replaces a page, it must also set `ColumnSelection = nil`; the next table
build clears the associated private Shift state. Hiding/removing fields
reconciles the chosen IDs and the cursor; reordering preserves IDs.

`HeaderTooltip` uses native tooltip behavior: initial 600 ms hover delay,
empty text has no tip, and leaving/pressing/menu opening dismisses it.
Native tooltips also retain their existing short warm transition between
adjacent tipped elements. `HeaderDetail` is a separate final line, independent
of `Subtitle`. CR/LF is flattened only for that line; tooltip text is intact.
Applications own placeholder text and persistence, never the framework.

`HeaderActionMenu` creates a separate hover-revealed `…` button inside the
same keyed header. It shares width, clipping and horizontal scrolling with
the header, reserves its space even while hidden, remains keyboard reachable,
and does not select, sort, resize or drag a column when opened. The existing
`HeaderMenu` is still the right-click menu. With `HeaderSelect`, right-click
inside the selected set preserves it; outside selects the clicked column
before the application builds that menu.

`OnCellClick(row, col, mods)` receives original cell-builder indices after
an undragged primary pointer click selects its cell. It excludes row headers,
section headers, key activation and drag completion. Editors can use it to
activate after selection while leaving native `DragSelect` responsible for
the selection rectangle and viewport-edge autoscroll. Application callbacks
continue to run on the UI thread.

`ui/grid_header_test.go` covers selection, menus, tooltip clock advancement,
reorder/resize, zero-row transitions, cross-Tester reuse, state cleanup,
cell callbacks, keyboard access and horizontal/vertical edge scrolling.
Runtime rendering has been checked on Linux; Windows/macOS builds are
compilation evidence, not desktop GUI validation.
