# Core Status

## Removed

- `cmd/sync-doctypes` removed.
- Reason: it targeted old `core_*` tables.
- Current database structure uses Frappe-style `tab*` tables:
  - `tabDocType`
  - `tabDocField`
  - `tabDocPerm`
  - `tabModule Def`
  - `tabInstalled App`
  - `tabInstalled Module`

- `internal/core/api/children.go` removed.
- Reason: `GetDocTypeActions`, `GetDocTypeLinks`, `GetDocTypeStates` all queried
  `core_doctype_action`, `core_doctype_link`, `core_doctype_state` — none of which exist.
  They were never registered in routes so caused no live errors.

## Removed old DocType syncer

- Removed `internal/core/doctype/sync.go`.
- Confirmed `cmd/sync-doctypes` already removed.
- Reason: old syncer targeted removed `core_*` tables.
- `slugify` and `hashBytes` helpers that `writer.go` depended on were extracted to `internal/core/doctype/helpers.go` before deletion.
- Current Core uses `tab*` metadata tables:
  - `tabDocType`
  - `tabDocField`
  - `tabDocPerm`
  - `tabModule Def`
  - `tabInstalled App`
  - `tabInstalled Module`

Future DocType save/sync should use the current `tab*` metadata path, not old `core_*` tables.

## Removed dead API handlers

- Removed `internal/core/api/children.go`.
- Reason: it referenced old `core_*` tables that no longer exist.
- Removed/unregistered:
  - `GetDocTypeActions`
  - `GetDocTypeLinks`
  - `GetDocTypeStates`

Future actions/links/states should be reintroduced only after we define current `tab*` metadata tables or child DocTypes for them.

## Removed orphan server binary

- Removed `cmd/server`.
- Reason: `cmd/gogal` is the main entry point.
- Server should be started with:

```powershell
go run .\cmd\gogal start
```

- This avoids duplicate server entry points and keeps CLI operations centralized.

## Current Core Direction

DocType metadata is managed through:

- installer seed
- Core API
- migration planner
- future Studio UI builder

Do not use old `core_*` sync pipeline.
