package platform

import "sync/atomic"

// MenuDensity selects the native menu metrics. The default preserves the
// platform's standard menu rendering; compact is currently supported by the
// Windows 11 owner-drawn menus.
type MenuDensity uint8

const (
	MenuDefault MenuDensity = iota
	MenuCompact
)

var currentMenuDensity atomic.Uint32

func SetMenuDensity(d MenuDensity) {
	currentMenuDensity.Store(uint32(d))
}

func CurrentMenuDensity() MenuDensity {
	return MenuDensity(currentMenuDensity.Load())
}
