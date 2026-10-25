package mygo

import "github.com/egoist/mygo/internal/platform"

// MenuDensity controls native menu spacing and font sizing.
type MenuDensity = platform.MenuDensity

const (
	// MenuDefault preserves MyGo's standard native menu appearance.
	MenuDefault = platform.MenuDefault
	// MenuCompact opts into smaller Windows 11 native menu metrics.
	// Windows 10 and non-Windows platforms retain their existing native style.
	MenuCompact = platform.MenuCompact
)

// SetMenuDensity configures the app's native menus. Set it before App.Run;
// existing native menu handles and font caches are not dynamically rebuilt.
// This setting does not affect ordinary UI widgets or OS text scaling.
func (a *Application) SetMenuDensity(d MenuDensity) {
	if d != MenuDefault && d != MenuCompact {
		panic("mygo: unsupported menu density")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		panic("mygo: SetMenuDensity must be called before App.Run")
	}
	platform.SetMenuDensity(d)
}

// MenuDensity reports the configured menu density.
func (a *Application) MenuDensity() MenuDensity {
	return platform.CurrentMenuDensity()
}
