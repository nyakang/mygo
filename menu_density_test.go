package mygo

import "testing"

func TestMenuDensitySetting(t *testing.T) {
	app := &Application{}
	old := app.MenuDensity()
	defer app.SetMenuDensity(old)
	app.SetMenuDensity(MenuCompact)
	if got := app.MenuDensity(); got != MenuCompact {
		t.Fatalf("got %d, want compact", got)
	}
	app.SetMenuDensity(MenuDefault)
	if got := app.MenuDensity(); got != MenuDefault {
		t.Fatalf("got %d, want default", got)
	}
}

func TestInvalidMenuDensityPanics(t *testing.T) {
	app := &Application{}
	defer func() {
		if recover() == nil {
			t.Fatal("invalid menu density must panic")
		}
	}()
	app.SetMenuDensity(MenuDensity(200))
}
