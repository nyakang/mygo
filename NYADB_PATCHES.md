# NyaDB MyGo maintenance branch

This repository is a fork of [egoist/mygo](https://github.com/egoist/mygo). The upstream `main` branch is kept suitable for GitHub **Sync fork**. NyaDB-specific changes belong on a separate `nyadb/v0.6.3` branch based on the upstream `v0.6.3` tag.

## Patch on `nyadb/v0.6.3`

The `ui/table.go` patch extends the upstream table implementation. It extends `TableColumn` with optional `Subtitle`, `SubtitleIcon` and `HeaderBadge` fields to display database field metadata inside the same native table header cell. With a blank subtitle, the original single-line header is preserved. Column identity, resizing, dragging, sorting, context menus and horizontal scrolling continue using the existing Table layout and state.

This is a native, independently developed UI extension for NyaDB, not code copied from Navicat or any other proprietary application. The original MyGo license remains unchanged.

## Opt-in compact menu density (Windows 11)

The fork adds `mygo.App.SetMenuDensity(mygo.MenuCompact)`, to be called before `App.Run`. Defaults are unchanged unless explicitly enabled. Windows 11 continues using the native HMENU / MFT_OWNERDRAW backend, with smaller font amplification, 24-DIP rows, 5-DIP separators and reduced horizontal padding. Windows 10 keeps Win32 native menu metrics. Linux GTK menu styling remains application-owned and macOS continues using native AppKit menus.

The setting preserves DPI scaling and honors the system menu font size, including user accessibility choices. It is deliberately immutable once the application starts so native menu handles and cached fonts remain internally consistent. Tests cover default/compact geometry, font policy and public API behavior. Windows pixel-perfect appearance still requires desktop validation; cross-compilation is not GUI verification.

## Keeping up with upstream

1. Sync the fork's `main` branch from `egoist/mygo:main` (GitHub Sync fork or a reviewed workflow).
2. Review a new upstream release tag before upgrading NyaDB. Never silently merge moving `main` into a version-pinned release.
3. Create a new `nyadb/vX.Y.Z` branch based on that exact upstream release tag, then reapply or adapt the minimal `ui/table.go` patch.
4. Run MyGo tests and NyaDB Go tests, vet checks, builds, and native UI regression tests.
5. Publish an immutable annotated/custom SemVer tag (for example `v0.6.3-nyadb.2`) and update NyaDB's version-pinned `replace` directive. Never move an existing version tag.
6. Drop this patch if upstream adopts an equivalent, stable API.

The Go module path in the fork remains `github.com/egoist/mygo`; NyaDB uses a Go module `replace` directive to resolve that module from `github.com/nyakang/mygo`.
