package ui

import "testing"

func TestTesterRepeatedLabelChoosesReachableControl(t *testing.T) {
	clicked := 0
	tt := coreNewTester(func(c *context) {
		coreColumn(c).Children(func() {
			coreBox(c).Label("Cancel").Size(0, 0)
			b := coreButton(c, "Cancel").Label("Cancel")
			if b.Clicked() {
				clicked++
			}
		})
	}, 300, 200)
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if clicked != 1 {
		t.Fatalf("reachable control clicked %d times", clicked)
	}
}

func TestTesterRepeatedLabelsStillRejectCoveredControls(t *testing.T) {
	clicked := 0
	tt := coreNewTester(func(c *context) {
		coreColumn(c).Children(func() {
			for i := 0; i < 2; i++ {
				b := coreButton(c, "Cancel").Label("Cancel")
				if b.Clicked() {
					clicked++
				}
			}
		})
		coreButton(c, "Cover").Absolute().Top(0).Left(0).Size(300, 200)
	}, 300, 200)
	if err := tt.Click("Cancel"); err == nil {
		t.Fatal("a repeated label bypassed hit testing")
	}
	if clicked != 0 {
		t.Fatal("covered control received a click")
	}
}
