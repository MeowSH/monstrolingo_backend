# simbuildcore alias mapping

This module contains the translation-only runtime for `/linkbuild/translate` and
the private alias mapping refresh tooling used to resolve Set Skill / Group
Skill effect labels.

## Manual refresh command

Run from the repository root:

```bash
go run ./internal/simbuildcore/localbridge/cmd/refresh \
  -languages all \
  -cache-mode refresh \
  -out internal/simbuildcore/localbridge/data/local_alias_snapshot.json
```

### Flags

- `-languages`: comma-separated language codes or `all` (default).
  - Supported: `en,ja,ko,zh-hans,zh-hant`
- `-cache-mode`: `live`, `offline`, or `refresh`.
  - `live`: fetch remote first, fallback to cache.
  - `offline`: read only local cache files.
  - `refresh`: force remote fetch and rewrite cache.
- `-cache-root`: local cache directory.
  - Default: `internal/simbuildcore/localbridge/cache`
- `-out`: output snapshot path.
  - Default: `internal/simbuildcore/localbridge/data/local_alias_snapshot.json`
- `-timeout-minutes`: global command timeout.
- `-base-url`: optional source URL override.

## Runtime usage

- The API runtime can load a local snapshot from `SIMBUILD_ALIAS_SNAPSHOT_PATH`.
- Alias resolution remains DB-first.
- Snapshot aliases are used only as fallback when direct skill/effect matching
  fails.
