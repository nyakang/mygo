package platform

import "testing"

func TestMenuDensity(t *testing.T) {
	old := CurrentMenuDensity()
	defer SetMenuDensity(old)
	SetMenuDensity(MenuCompact)
	if got := CurrentMenuDensity(); got != MenuCompact {
		t.Fatalf("compact density: got %d", got)
	}
	SetMenuDensity(MenuDefault)
	if got := CurrentMenuDensity(); got != MenuDefault {
		t.Fatalf("default density: got %d", got)
	}
}
