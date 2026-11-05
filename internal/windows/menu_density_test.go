//go:build windows && (amd64 || arm64)

package windows

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

func TestMenuSizingPreservesDefaultsAndCompactValues(t *testing.T) {
	old := platform.CurrentMenuDensity()
	defer platform.SetMenuDensity(old)
	platform.SetMenuDensity(platform.MenuDefault)
	def := menuSizing()
	if def.itemHeight != 32 || def.separatorHeight != 9 || def.padding != 12 ||
		def.barHeight != 30 || def.markColumn != 28 {
		t.Fatalf("default Fluent sizing regressed: %+v", def)
	}
	platform.SetMenuDensity(platform.MenuCompact)
	compact := menuSizing()
	if compact.itemHeight != 24 || compact.separatorHeight != 5 || compact.padding != 8 ||
		compact.barHeight != 26 || compact.markColumn != 24 {
		t.Fatalf("unexpected compact menu sizing: %+v", compact)
	}
	if compact.itemHeight >= def.itemHeight || compact.barHeight >= def.barHeight ||
		compact.accelGap >= def.accelGap {
		t.Fatal("compact sizing must reduce item height, bar height and accelerator gap")
	}
}

func TestMenuLookInsetsFollowDensity(t *testing.T) {
	old := platform.CurrentMenuDensity()
	defer platform.SetMenuDensity(old)
	l := menuLook{first: true, last: true}
	platform.SetMenuDensity(platform.MenuDefault)
	top, bottom := l.padding(96)
	if top != 4 || bottom != 4 {
		t.Fatalf("default padding: %d, %d", top, bottom)
	}
	platform.SetMenuDensity(platform.MenuCompact)
	top, bottom = l.padding(96)
	if top != 2 || bottom != 2 {
		t.Fatalf("compact padding: %d, %d", top, bottom)
	}
}

func TestCompactMenuFontUsesSystemHeight(t *testing.T) {
	old := platform.CurrentMenuDensity()
	defer platform.SetMenuDensity(old)
	platform.SetMenuDensity(platform.MenuDefault)
	if got := menuFontHeight(-12); got != -14 {
		t.Fatalf("default font scaling: %d", got)
	}
	platform.SetMenuDensity(platform.MenuCompact)
	for _, fontHeight := range []int32{-12, -15, -18} {
		if got := menuFontHeight(fontHeight); got != fontHeight {
			t.Fatalf("compact menu must retain system font scale: %d => %d", fontHeight, got)
		}
	}
	if got := toPx(menuSizing().itemHeight, 120); got != 30 {
		t.Fatalf("compact 125%% DPI scaling: %d", got)
	}
}

func TestMenuWidthCompensationKeepsDensityDPIAndMenuIsolation(t *testing.T) {
	old := platform.CurrentMenuDensity()
	defer platform.SetMenuDensity(old)
	for _, density := range []platform.MenuDensity{platform.MenuDefault, platform.MenuCompact} {
		platform.SetMenuDensity(density)
		b := &Backend{}
		a := &menuLook{menu: 1, dpi: 192, width: 200}
		b.menus.looks = map[uintptr]*menuLook{
			1: a,
			2: {menu: 1, dpi: 192, width: 100},
			3: {menu: 2, dpi: 96, width: 500},
		}
		b.learnMenuExtra(a, 192, 226)
		if b.menuExtra[192] != 26 {
			t.Fatalf("density %v: width compensation = %v, want 26", density, b.menuExtra)
		}
		b.learnMenuExtra(a, 192, 240)
		b.learnMenuExtra(a, 192, 190)
		if b.menuExtra[192] != 26 {
			t.Fatal("larger or negative observed padding replaced the least known compensation")
		}
		b.menus.looks[2].dpi = 96
		b.learnMenuExtra(a, 192, 210)
		if b.menuExtra[192] != 26 {
			t.Fatal("a partly remeasured menu changed DPI compensation")
		}
		b.menus.looks[2].dpi = 192
		b.learnMenuExtra(a, 192, 220)
		if b.menuExtra[192] != 20 {
			t.Fatal("smaller observed padding was not retained")
		}
		if _, ok := b.menuExtra[96]; ok {
			t.Fatal("compensation leaked to another DPI")
		}
		b.dropMenuFonts()
		if b.menuExtra != nil {
			t.Fatal("system font invalidation retained obsolete menu compensation")
		}
	}
}
