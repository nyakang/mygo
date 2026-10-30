# NyaDB MyGo maintenance branch

This repository is a fork of [egoist/mygo](https://github.com/egoist/mygo). The upstream `main` branch is kept suitable for GitHub **Sync fork**. NyaDB-specific changes belong on version-specific `nyadb/vX.Y.Z` branches based on their corresponding upstream release tags.

## NyaDB v0.7.3 compatibility release

Based on upstream `v0.7.3` (commit
`6f28d7bb07b31177cf218e51586f6a5e99b38016`), released October 29, 2026.
The immutable Go module tag is `v0.7.3-nyadb.1` on `nyadb/v0.7.3`.
The upstream module path and minimum Go version (1.27.3) remain unchanged.

Reapplied the complete `v0.6.4-nyadb.4` customization: table metadata headers,
opt-in stripes/grid, Windows compact menus, configurable outline triangles,
full-width outline selection and disclosure-gutter context menus. The table
and outline patches apply without code conflicts. Windows menu painting is
adapted to upstream's bitmap callbacks, retaining the new menu-icon rendering
while using the existing default/compact density metrics. The `main` branch
tracks upstream separately from this release baseline.

Verification on Linux: `GOWORK=off CGO_ENABLED=0 go test ./...`, full
Linux/Windows/Darwin builds and vet, production-tag vet, and js/wasm vet pass.
NyaDB's candidate dependency passes its full tests, Go vet, MyGo UI vet,
localization validation, and Linux/Windows/Darwin builds. Desktop GUI checks
on Windows and macOS are not part of these compilation results.

Release-size comparison against the unmodified v0.7.3 baseline, with cgo off,
`-trimpath`, `mygo_noinspector`, stripped symbols and `production=1`:

| App and platform | Upstream bytes | Fork bytes | Difference |
| --- | ---: | ---: | ---: |
| linux/amd64/counter-native | 9,564,423 | 9,564,423 | +0 |
| linux/amd64/hello | 9,613,575 | 9,613,575 | +0 |
| windows/amd64/counter-native | 9,876,480 | 9,881,088 | +4,608 |
| windows/amd64/hello | 9,970,688 | 9,974,784 | +4,096 |
| darwin/arm64/counter-native | 9,372,770 | 9,372,866 | +96 |
| darwin/arm64/hello | 9,561,634 | 9,561,634 | +0 |

Windows growth comes from the retained compact-menu policy in shared native
menu measurement/painting and `menuFontHeight`, plus metadata/alignment.
The Windows counter's `-dumpdep` and `nm -size` show `menuFace` reaching
`menuFontHeight` and the existing default menu paths retaining density
branches. `ListState` also retains `TableStyle` type metadata, as before;
unused table-grid drawing is not linked. The Darwin counter's 96-byte growth
is type metadata. Linux example sizes remain unchanged. Logs and symbol/size
comparisons live in NyaDB's ignored `temp/research/mygo-upgrade-20261030/`.

## NyaDB v0.6.4 compatibility release

Based on upstream `v0.6.4` (commit `6255ea40139d06e0762a7c4a32626148ba410478`).
The `nyadb/v0.6.4` branch independently reapplies three existing NyaDB changes:
native double-line table headers, opt-in Windows compact menu density, and opt-in
table stripes/grid. All three cherry-picked without conflicts.
The upstream `main` branch stays free of NyaDB-specific changes.

Previous patch commits: `5c761105`, `005838d6`, `0a7d0b1d`.
Upstream v0.6.4 requires Go 1.27.3.

## Patch on `nyadb/v0.6.3`

The `ui/table.go` patch extends the upstream table implementation. It extends `TableColumn` with optional `Subtitle`, `SubtitleIcon` and `HeaderBadge` fields to display database field metadata inside the same native table header cell. With a blank subtitle, the original single-line header is preserved. Column identity, resizing, dragging, sorting, context menus and horizontal scrolling continue using the existing Table layout and state.

This is a native, independently developed UI extension for NyaDB, not code copied from Navicat or any other proprietary application. The original MyGo license remains unchanged.

## Opt-in compact menu density (Windows 11)

The fork adds `mygo.App.SetMenuDensity(mygo.MenuCompact)`, to be called before `App.Run`. Defaults are unchanged unless explicitly enabled. Windows 11 continues using the native HMENU / MFT_OWNERDRAW backend, with smaller font amplification, 24-DIP rows, 5-DIP separators and reduced horizontal padding. Windows 10 keeps Win32 native menu metrics. Linux GTK menu styling remains application-owned and macOS continues using native AppKit menus.

The setting preserves DPI scaling and honors the system menu font size, including user accessibility choices. It is deliberately immutable once the application starts so native menu handles and cached fonts remain internally consistent. Tests cover default/compact geometry, font policy and public API behavior. Windows pixel-perfect appearance still requires desktop validation; cross-compilation is not GUI verification.

## Keeping up with upstream

### Native table grid (`v0.6.3-nyadb.3`)

Based on `v0.6.3-nyadb.2` (`005838d6835c7953e2e51f027a81f84b49f7dc5e`).
`ListState.TableStyle` is optional; nil and the zero style preserve rendering
and measurement, and ordinary Lists ignore it. `AlternateRowColor` fills odd
display rows (including row headers and padding), with `RowOffset` preserving
parity across pages. `GridColor`, `RowLines` and `ColumnLines` independently
enable bottom and trailing edges. Edges follow the native cell layout, including
RTL, resize, reorder and scroll, and never affect measurement. Shared edges and
corners paint once. Selection covers stripes; grid painting precedes the native
cursor outline. No changes to native virtualization or keyboard behavior.

Verification: native raster/interaction tests cover default compatibility,
100,000 virtual rows, page parity, bidirectional scrolling, resize/reorder/hide,
independent edges, translucent corners and selection/focus layering. Linux,
Windows, Darwin and browser compilation/vet are checked with CGO disabled.
Runtime GUI validation is on Linux; Windows/macOS rendering still needs their
desktops.

Release-size comparison (CGO=0, trimpath, mygo_noinspector, stripped,
production=1): counter-native Linux 9,359,623 → 9,359,623 bytes; Windows
9,660,416 → 9,661,440; Darwin arm64 9,167,282 → 9,167,378. The unused-feature
growth is type metadata: linker `-dumpdep` shows `ListState → *TableStyle →
TableStyle` and its field names, but no table grid drawing code in the counter.
Section alignment accounts for the platform difference. The webview hello
binary is unchanged on all three platforms (9,433,351 / 9,781,760 / 9,372,834).

1. Sync the fork's `main` branch from `egoist/mygo:main` (GitHub Sync fork or a reviewed workflow).
2. Review a new upstream release tag before upgrading NyaDB. Never silently merge moving `main` into a version-pinned release.
3. Create a new `nyadb/vX.Y.Z` branch based on that exact upstream release tag, then reapply or adapt the minimal `ui/table.go` patch.
4. Run MyGo tests and NyaDB Go tests, vet checks, builds, and native UI regression tests.
5. Publish an immutable annotated/custom SemVer tag (for example `v0.6.3-nyadb.2`) and update NyaDB's version-pinned `replace` directive. Never move an existing version tag.
6. Drop this patch if upstream adopts an equivalent, stable API.

The Go module path in the fork remains `github.com/egoist/mygo`; NyaDB uses a Go module `replace` directive to resolve that module from `github.com/nyakang/mygo`.
