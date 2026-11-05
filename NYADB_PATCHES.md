# NyaDB MyGo maintenance branch

This repository is a fork of [egoist/mygo](https://github.com/egoist/mygo). The upstream `main` branch is kept suitable for GitHub **Sync fork**. NyaDB-specific changes belong on version-specific `nyadb/vX.Y.Z` branches based on their corresponding upstream release tags.

## NyaDB v0.7.4 compatibility release

Based on upstream `v0.7.4`, commit
`9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e`, released November 1, 2026.
The immutable Go module tag is `v0.7.4-nyadb.1` on `nyadb/v0.7.4`.
The module path remains `github.com/egoist/mygo`, and Go stays at 1.27.3.
The fork's `main` tracks upstream separately from this fixed release.

Preserves the complete `v0.7.3-nyadb.2` customization: multi-line native
headers and COMMENT tips, discrete column selection, independent header
actions, cell-click callbacks, stripes/grid, outline triangles/full-row
selection/context menus, and Windows compact menu density. Table subtree
themes and input/runtime changes from v0.7.4 remain intact.

Windows menu measurement retains upstream's `itemWidth`/`learnMenuExtra`
algorithm, measured item widths and DPI, and compensation invalidation on
system font changes. Compact sizing supplies the dimensions and font
policy to the new measurement/painting paths; it does not restore the old
`SM_CXMENUCHECK` subtraction that clipped labels at high scaling. A new
deterministic Windows test checks both densities, menu/DPI isolation,
partly remeasured menus, and font invalidation. It is cross-compiled on
Linux; Windows runtime execution still requires its desktop.

Release-size comparison against unmodified v0.7.4, with identical
`CGO_ENABLED=0`, `-buildvcs=false`, `-trimpath`, `mygo_noinspector`, stripped
symbols and `production=1`:

| App and platform | Upstream bytes | Fork bytes | Difference |
| --- | ---: | ---: | ---: |
| linux/amd64/counter-native | 9,576,711 | 9,580,807 | +4,096 |
| linux/amd64/hello | 9,621,767 | 9,621,767 | 0 |
| windows/amd64/counter-native | 9,899,008 | 9,905,664 | +6,656 |
| windows/amd64/hello | 9,980,928 | 9,986,048 | +5,120 |
| darwin/arm64/counter-native | 9,389,762 | 9,389,986 | +224 |
| darwin/arm64/hello | 9,578,370 | 9,578,370 | 0 |

The native counter still links no `pointerCellClick` or grid drawing
machine code. `-dumpdep` and `nm -size` trace the increment to the guarded
row-input callback (150/160/144 bytes on Linux/Windows/Darwin), retained
ListState/TableStyle metadata, and the existing Windows compact menu
policy reachable through shared menu measurement/painting. Linux `.text`
grows 224 bytes, `.gopclntab` 257 and `.go.type` 1,056; the text segment
crosses a page boundary, producing the 4,096-byte file increment.
Windows menu symbols account for most of its shared code growth; PE
alignment and metadata determine the final file size. Darwin's native
counter growth includes the callback, type equality and unwind/metadata
layout. Unused webview examples remain unchanged on Linux and Darwin.
Logs, symbol/link maps and images remain in NyaDB's ignored
`temp/research/mygo-upgrade-20261105/`, outside this fork's published tree.

Full Go tests, Linux/Windows amd64 and Darwin arm64 builds/vet,
production-tag and js/wasm vet pass with cgo disabled. Table/list/outline
and public menu-density focused race checks pass with cgo enabled.
NyaDB's candidate passes full Go tests, Go/MyGo vet, unchanged generated
translations (1,711 messages), three-platform builds, and native Linux
SQLite workflow replay in Chinese and English at 1,280/840 pixels. The
replay closes every owned process and verifies the database is unchanged.
The general MyGo e2e harness cannot start on this host because it creates
a webview while WebKitGTK is absent; its failed attempt is recorded,
and is not described as passing. Windows/macOS desktop GUI checks are
not part of the cross-compilation results.

## Native grid headers: v0.7.3-nyadb.2

Based on `nyadb/v0.7.3` at the immutable `.1` commit
`922da57bd8b2ec6cf1576441afbe7f6565a536ca`. The new immutable Go module tag
is `v0.7.3-nyadb.2` on `nyadb/v0.7.3`. The module path and upstream
version remain unchanged; this is a Go-only fork patch, not an npm release.
The inherited npm release workflow excludes `v*-nyadb.*` tags; publishing
a new immutable Go tag does not publish upstream npm packages.

Adds opt-in whole/discrete column selection, native COMMENT tooltips and a
third header line, an independent hover action menu, and an undragged cell
click callback. Header interactions retain keyed identity, native resize,
reorder and scrolling. Empty tables use the same native headers with no
synthetic record. Column IDs are reconciled after removal and private Shift
state is revoked on direct selection reset. A ListState can be reused by a
new window/Tester, while two lists in the same context are still rejected.
See [native header API and lifecycle](docs/nyadb-grid-headers.md).

The implementation is independently developed from runtime behavior and
the existing MyGo sources. No proprietary code, icons or screenshots enter
the patch. NyaDB retains its own theme and icon set.

Verification with cgo disabled: full Go tests, Linux/Windows amd64 and
Darwin arm64 builds/vet, production-tag vet, and js/wasm vet excluding the
CLI pass. Native header tests advance Tester time rather than sleeping.
Linux runtime checks cover header selection, menus, COMMENTs, empty tables,
width, zoom and both axes of native selection autoscroll. NyaDB's 180-row
nonadjacent-column clipboard output matches the original byte for byte.

Release-size comparison against `.1`, using identical `-buildvcs=false`,
`-trimpath`, `mygo_noinspector`, stripped symbols and `production=1`:

| App and platform | `.1` bytes | Candidate bytes | Difference |
| --- | ---: | ---: | ---: |
| linux/amd64/counter-native | 9,564,423 | 9,564,423 | +0 |
| linux/amd64/hello | 9,613,575 | 9,613,575 | +0 |
| windows/amd64/counter-native | 9,881,088 | 9,882,112 | +1,024 |
| windows/amd64/hello | 9,974,272 | 9,974,272 | +0 |
| darwin/arm64/counter-native | 9,372,866 | 9,372,978 | +112 |
| darwin/arm64/hello | 9,561,634 | 9,561,634 | +0 |

The new pointer-click dispatcher is registered by an opt-in Table; `nm`
and `-dumpdep` confirm that column hit-testing machine code is not linked
into counter-native. Its remaining Windows/Darwin growth is the small
guarded row-input dispatch and callback/type metadata. Windows `.text`
grows 160 virtual bytes and `.rdata` 696, rounded to two 512-byte sections;
Darwin text grows 144 bytes, with metadata/alignment accounting for the
112-byte file difference. Unused webview examples remain unchanged.
An initial measurement included differing VCS build info; the controlled
comparison removes it and verifies identical `go version -m` output.
Logs, source hashes and symbol/section deltas are retained in NyaDB's
ignored `temp/research/grid-alignment-20261031/mygo-validation/`.

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
