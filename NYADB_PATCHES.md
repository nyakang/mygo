# NyaDB MyGo maintenance branch

This repository is a fork of [egoist/mygo](https://github.com/egoist/mygo). The upstream `main` branch is kept suitable for GitHub **Sync fork**. NyaDB-specific changes belong on a separate `nyadb/v0.6.3` branch based on the upstream `v0.6.3` tag.

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
