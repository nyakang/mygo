package ui

import "testing"

func TestOutlineLeadingInsetPreservesHierarchy(t *testing.T) {
	measure := func(inset float32) (float32, float32) {
		t.Helper()
		s := OutlineState[string]{Indent: 14, LeadingInset: inset}
		s.Open.Add("src")
		tt := coreNewTester(func(c *context) {
			coreOutline(c, &s, []string{"src"}, children, func(path string) {
				coreText(c, baseName(path))
			}).Grow(1)
		}, 400, 200)
		root, rootOK := tt.Find("src")
		child, childOK := tt.Find("ui")
		if !rootOK || !childOK {
			t.Fatal("missing root/child text")
		}
		return root.X, child.X
	}
	root0, child0 := measure(0)
	root14, child14 := measure(14)
	if root14-root0 != 14 || child14-child0 != 14 {
		t.Fatalf("leading inset must shift every level equally: root %v child %v", root14-root0, child14-child0)
	}
	if child14-root14 != child0-root0 {
		t.Fatalf("leading inset changed the existing depth step: before %v after %v", child0-root0, child14-root14)
	}
}

func TestOutlineLeadingInsetRejectsNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("negative LeadingInset did not panic")
		}
	}()
	s := OutlineState[string]{LeadingInset: -1}
	coreNewTester(func(c *context) {
		coreOutline(c, &s, []string{"src"}, children, func(string) {})
	}, 400, 200)
}

func TestOutlineRowContextMenuIncludesGutterAndTrailingArea(t *testing.T) {
	s := OutlineState[string]{LeadingInset: 14, Indent: 14, Disclosure: OutlineTriangle, DisclosureSize: 10}
	s.RowContextMenu = func(item string, m *Menu) { m.Item("Open " + item) }
	tt := coreNewTester(func(c *context) {
		coreOutline(c, &s, []string{"src"}, children, func(item string) {
			coreText(c, baseName(item))
		}).Grow(1)
	}, 320, 180)
	root, ok := tt.Find("src")
	if !ok {
		t.Fatal("missing root row")
	}
	for _, x := range []float32{8, 20, 250} {
		tt.RightClickAt(x, root.Y+root.H/2)
		if menu := tt.Menu(); len(menu) == 0 || menu[0] != "Open src" {
			t.Fatalf("right click x=%.1f had no row menu: %v", x, menu)
		}
		tt.CloseMenu()
		tt.Frame()
	}
}
