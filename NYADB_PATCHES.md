# NyaDB MyGo maintenance branch

This repository is a fork of [egoist/mygo](https://github.com/egoist/mygo). The upstream `main` branch is kept suitable for GitHub **Sync fork**. NyaDB-specific changes belong on a separate `nyadb/v0.6.3` branch based on the upstream `v0.6.3` tag.

## Patch on `nyadb/v0.6.3`

Only `ui/table.go` differs from the corresponding upstream source. It extends `TableColumn` with optional `Subtitle`, `SubtitleIcon` and `HeaderBadge` fields to display database field metadata inside the same native table header cell. With a blank subtitle, the original single-line header is preserved. Column identity, resizing, dragging, sorting, context menus and horizontal scrolling continue using the existing Table layout and state.

This is a native, independently developed UI extension for NyaDB, not code copied from Navicat or any other proprietary application. The original MyGo license remains unchanged.

## Keeping up with upstream

1. Sync the fork's `main` branch from `egoist/mygo:main` (GitHub Sync fork or a reviewed workflow).
2. Review a new upstream release tag before upgrading NyaDB. Never silently merge moving `main` into a version-pinned release.
3. Create a new `nyadb/vX.Y.Z` branch based on that exact upstream release tag, then reapply or adapt the minimal `ui/table.go` patch.
4. Run MyGo tests and NyaDB Go tests, vet checks, builds, and native UI regression tests.
5. Publish an immutable annotated/custom SemVer tag (for example `v0.6.3-nyadb.1`) and update NyaDB's version-pinned `replace` directive. Never move an existing version tag.
6. Drop this patch if upstream adopts an equivalent, stable API.

The Go module path in the fork remains `github.com/egoist/mygo`; NyaDB uses a Go module `replace` directive to resolve that module from `github.com/nyakang/mygo`.
