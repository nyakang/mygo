package ui

import "slices"

// selectHeaderColumn runs while the native header is built, before its menu
// collects items. Use the current row count and node, rather than a previous
// frame's listFrame: the rows can have changed while the menu was requested.
func (s *ListState) selectHeaderColumn(table *node, id string, mods Modifiers, n int) {
	if !s.HeaderSelect || !s.choosesCells() {
		return
	}
	at := s.colIndex(id)
	if at < 0 || s.tableCols[at].rowHeader {
		return
	}
	before := s.ColumnSelection
	ids := slices.Clone(before)
	switch {
	case mods&Shift != 0 && s.colIndex(s.headerPivot) >= 0:
		if s.headerBase == nil {
			s.headerBase = slices.Clone(ids)
		}
		ids = slices.Clone(s.headerBase)
		pivot := s.colIndex(s.headerPivot)
		for k := min(pivot, at); k <= max(pivot, at); k++ {
			if col := s.tableCols[k]; !col.rowHeader && !slices.Contains(ids, col.id) {
				ids = append(ids, col.id)
			}
		}
	case mods&Cmd != 0 && len(ids) > 0:
		s.headerBase, s.headerPivot = nil, id
		if k := slices.Index(ids, id); k >= 0 {
			if len(ids) > 1 {
				ids = slices.Delete(ids, k, k+1)
			}
		} else {
			ids = append(ids, id)
		}
	default:
		ids = []string{id}
		s.headerBase, s.headerPivot = nil, id
	}
	// Drop base entries that a hide/removal made stale during Shift selection.
	ids = slices.DeleteFunc(ids, func(id string) bool {
		k := s.colIndex(id)
		return k < 0 || s.tableCols[k].rowHeader
	})
	cursor := id
	if !slices.Contains(ids, cursor) && len(ids) > 0 {
		cursor = ids[len(ids)-1]
	}
	cells := ColumnRange{Cursor: cursor, Anchor: s.headerPivot}
	changed := !slices.Equal(before, ids) || s.Cells != cells
	s.ColumnSelection, s.Cells = ids, cells
	selected := s.cursor()
	if *selected < 0 || *selected >= n || s.Header != nil && s.Header(*selected) {
		*selected = -1
		for row := range n {
			if s.Header == nil || !s.Header(row) {
				*selected = row
				break
			}
		}
		changed = true
	}
	if *selected >= 0 {
		s.pivot, s.span = *selected, 0
	}
	// Only real, currently loaded row keys are selected. No placeholder row
	// is introduced for an empty table or a statistical section.
	changed = s.Selection.clear() || changed
	for row := range n {
		if s.Header != nil && s.Header(row) {
			continue
		}
		var key any = row
		if s.Key != nil {
			key = s.Key(row)
		}
		changed = s.Selection.set(key, true) || changed
	}
	if changed {
		table.st.markChanged()
		table.c.rt.consumed = true
	}
}
