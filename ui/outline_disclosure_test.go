package ui

import "testing"

func TestOutlineTrianglePathDirections(t *testing.T) {
	r := Rect{X: 2, Y: 3, W: 10, H: 10}
	for _, tc := range []struct {
		name        string
		deg         float32
		rtl         bool
		right, down bool
	}{
		{"closed-ltr", 0, false, true, false},
		{"open-ltr", 90, false, false, true},
		{"closed-rtl", 0, true, false, false},
		{"open-rtl", 90, true, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := outlineTrianglePath(r, tc.deg, tc.rtl)
			if len(p.cmds) != 4 || p.cmds[3].op != 4 {
				t.Fatalf("triangle must be a closed three-vertex path: %+v", p.cmds)
			}
			pts := [3][2]float32{p.cmds[0].pts[0], p.cmds[1].pts[0], p.cmds[2].pts[0]}
			baseX, tipX := (pts[0][0]+pts[2][0])/2, pts[1][0]
			baseY, tipY := (pts[0][1]+pts[2][1])/2, pts[1][1]
			switch {
			case tc.right && !(tipX > baseX):
				t.Fatalf("collapsed LTR triangle must point right: %v", pts)
			case tc.down && !(tipY > baseY):
				t.Fatalf("expanded triangle must point down: %v", pts)
			case tc.rtl && !tc.down && !(tipX < baseX):
				t.Fatalf("collapsed RTL triangle must point left: %v", pts)
			}
		})
	}
}

func TestOutlineTriangleClickKeyboardAndAccess(t *testing.T) {
	selected := -1
	s := OutlineState[string]{Disclosure: OutlineTriangle, DisclosureSize: 10, Indent: 17}
	s.List.Selected = &selected
	tt := coreNewTester(func(c *context) {
		coreOutline(c, &s, []string{"src"}, children, func(path string) {
			coreText(c, baseName(path))
		}).Grow(1)
	}, 400, 200)
	tt.SetPreferences(Preferences{ReduceMotion: true})
	src, ok := tt.Find("src")
	if !ok {
		t.Fatal("root item missing")
	}
	tt.ClickAt(src.X-9, src.Y+src.H/2)
	if !s.Open.Has("src") || selected != -1 || !tt.HasText("main.go") {
		t.Fatalf("triangle click: open=%v selected=%d rows=%v", s.Open.Has("src"), selected, tt.Texts())
	}
	tt.Click("src")
	tt.Key(0, KeyLeft)
	if s.Open.Has("src") {
		t.Fatal("keyboard Left should close triangle node")
	}
	tt.Key(0, KeyRight)
	if !s.Open.Has("src") {
		t.Fatal("keyboard Right should open triangle node")
	}
	tt.ClickAt(src.X-9, src.Y+src.H/2)
	if s.Open.Has("src") {
		t.Fatal("triangle click should collapse expanded node")
	}
}

func TestOutlineDisclosureSizeValidation(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("negative disclosure size should panic")
		}
	}()
	s := OutlineState[string]{DisclosureSize: -1}
	coreNewTester(func(c *context) { coreOutline(c, &s, []string{"src"}, children, func(string) {}) }, 400, 200)
}
