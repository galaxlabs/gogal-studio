# Gogal Studio

> **Laravel productivity · Frappe metadata · Django admin simplicity · Go performance.**

Gogal Studio is a Go-powered, metadata-driven business application framework built by [Galaxy Labs](https://github.com/galaxlabs). Define your data model visually, and the engine provides the REST API, permissions, migrations, naming, validation, and admin UI automatically — no boilerplate, no scattered generators, no manual API wiring.

**The promise:** Describe app → Choose features → Preview migration → Apply → Ship.

---

## Why Gogal Studio?

Most Go frameworks hand you a router and leave the rest to you. Gogal Studio gives you:

- A **visual DocType builder** — define data models in the browser, not in code files
- A **metadata engine** — the entire app runs from JSON-defined DocTypes, no hand-written controllers
- A **safe migration planner** — always preview `ALTER TABLE` plans before they touch the database
- A **hardened CRUD engine** — permission gate, metadata validation, safe column filtering, system field injection, and transaction wrap on every write
- A **permission engine** — role-based per-DocType control over read / write / create / delete
- A **naming series engine** — configurable document naming patterns (series, field, expression, auto-increment)
- **Single binary deployment** — entire app ships as one Go binary with embedded assets

No Laravel-style command explosion. No Frappe Python dependency. No Node.js build pipeline required.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.22+, Fiber v2.52 |
| Database | PostgreSQL 15+ (pgx v5) |
| Migration Engine | Custom safe planner (widening-only, preview before apply) |
| Frontend Admin | Vanilla JS + CSS (zero framework) |
| CLI | Cobra |
| Auth | JWT / session hybrid *(planned)* |
| Queue / Cache | Redis + Asynq *(planned)* |
| AI Providers | OpenAI, Gemini, Claude, Ollama *(planned)* |
| Deployment | Single Go binary + config |

---

## Architecture

```
Gogal Studio
├── Go API Server (Fiber v2)
├── PostgreSQL
├── Core Metadata Engine        ← DocType, DocField, DocPerm, Module Def, etc.
├── Schema Migration Planner    ← safe diff, preview, apply
├── Permission Engine           ← role-based, per-DocType, per-action
├── Naming Series Engine        ← series / field / expression / auto
├── System Fields Engine        ← auto-inject name, owner, creation, modified
├── Document Validator          ← required, select options, hidden/read-only skip
├── Dynamic CRUD Engine         ← hardened POST / PUT / DELETE + safe read
├── Vanilla JS Admin UI         ← Studio, Builder, Resource List/Form
├── DocType Builder UI
├── Feature Pack Engine         ← (planned)
├── AI Agent Engine             ← (planned)
├── Bot Engine                  ← (planned)
├── Jobs / Queue Engine         ← (planned)
└── gogal CLI
```

---

## Getting Started

**Prerequisites:** Go 1.22+, PostgreSQL 15+, Docker (optional)

```bash
# Clone
git clone https://github.com/galaxlabs/gogal-studio
cd gogal-studio

# Install — creates DB, runs migrations, seeds core metadata
go run ./cmd/gogal install

# Start server
go run ./cmd/gogal start
# → http://127.0.0.1:8080
```

Open:
- **Studio UI** → `http://localhost:8080/studio`
- **DocType builder** → `http://localhost:8080/studio-assets/builder.html`
- **Health check** → `http://localhost:8080/health`

---

## CLI Reference

```bash
gogal install                        # Bootstrap DB, tables, core seed data
gogal start                          # Start HTTP server on :8080
gogal doctor                         # Diagnose config and DB health
gogal migrate plan   <doctype>       # Preview schema diff for a DocType
gogal migrate apply  <doctype>       # Apply schema changes (requires --confirm)
gogal migrate status                 # Show pending/applied migrations

# DocType management
gogal doctype list
gogal doctype show  "Naming Series"
gogal doctype export "Naming Series"   # Write JSON to modules/ path

# Naming series
gogal naming list
gogal naming next  "INV-.YYYY.-"

# Slug helpers
gogal slug app    gogal_studio        # → gogal_studio
gogal slug module Core               # → core
gogal slug doctype "Naming Series"   # → naming_series

# Field types
gogal fieldtype list
gogal fieldtype show Link
```

---

## What's Built

### Core Infrastructure
- [x] Go module — `github.com/galaxylabs/gogal-studio`
- [x] Fiber HTTP server — health, version, and full API route tree
- [x] PostgreSQL connection pool (pgx v5)
- [x] Cobra CLI — single `gogal` binary, all subcommands registered
- [x] Bootstrap installer (`gogal install`) — idempotent, safe to re-run
- [x] Environment config loader (`GOGAL_*` env vars + config file)

### Metadata Engine
- [x] Core metadata tables — `tabDocType`, `tabDocField`, `tabDocPerm`, `tabModule Def`, `tabNaming Series`, `tabInstalled App`, `tabInstalled Module`, `tabRole`, `tabUser`, `tabHas Role`
- [x] Core DocType JSON files under `modules/core/doctype/` — source of truth for core schema
- [x] DocType JSON reader/writer — reads JSON, validates, writes slug-pathed files
- [x] DocType JSON importer — imports from file or SQL
- [x] DocType meta loader — loads full DocType + fields + perms from DB at runtime
- [x] Core DocTypes seeded: Module Def, DocType, DocField, DocPerm, DocType Action, DocType Link, DocType State, Naming Series, Role, User, Installed App, Installed Module
- [x] `tabSchema Migration Log` — migration audit trail

### Migration Planner *(safe-only)*
- [x] Schema diff engine — compares DocType definition vs live `information_schema`
- [x] Safe widening detection — allows `varchar(n) → varchar(m)` where `m ≥ n`, `int → bigint`, etc.
- [x] Blocked operations — DROP COLUMN, DROP TABLE, narrowing casts always require explicit confirmation
- [x] Index planner — `CREATE INDEX` for Link / Data fields
- [x] `gogal migrate plan` — prints human-readable diff before any change
- [x] `gogal migrate apply` — executes plan inside a transaction, writes audit log row

### Permission Engine
- [x] `tabDocPerm` — role-based permissions: read, write, create, delete, submit, cancel, amend
- [x] `CanUserRead / CanUserWrite / CanUserCreate / CanUserDelete` — checked on every CRUD call
- [x] `GET /api/core/permissions/check?user=&doctype=&action=` — live permission check endpoint
- [x] Both user-role and direct-user permission modes
- [x] 36 unit tests across permission, migration, CRUD, sysfields, and validator packages

### Naming Series Engine
- [x] `tabNaming Series` — configurable series patterns per DocType
- [x] Supported naming rules: `series:`, `field:`, `expression:`, prompt, auto-increment
- [x] `NextSeries` — atomically increments counter in DB
- [x] CRUD create automatically resolves document name via naming rule → fallback chain

### System Fields Engine (`internal/core/sysfields`)
- [x] `InjectCreate` — sets `name`, `owner`, `creation`, `modified`, `modified_by`, `docstatus`, `idx`
- [x] `InjectUpdate` — sets `modified`, `modified_by`
- [x] `ProtectedFields` — `name`, `owner`, `creation`, `docstatus` — never overwritten by user payload
- [x] 6 unit tests

### Document Validator (`internal/core/validator`)
- [x] Required field enforcement — missing or empty value → `400`
- [x] Select option validation — value must be in declared options list
- [x] Hidden field skip — hidden required fields never block save
- [x] Read-only field skip — read-only required fields never block save
- [x] Layout field skip — Section Break, Column Break, etc. are never validated
- [x] 10 unit tests

### CRUD Engine (Hardened)

All write endpoints enforce the full chain:

```
CanUser{Create|Write|Delete}  →  load meta  →  validate  →  resolve name
  →  inject system fields  →  safe column filter  →  transaction  →  RETURNING *
```

| Endpoint | Method | Gate |
|---|---|---|
| `POST /api/resource/:doctype` | Create | `CanUserCreate` + validate + naming + sysfields + safe cols |
| `PUT /api/resource/:doctype/:name` | Update | `CanUserWrite` + validate + sysfields + safe cols |
| `DELETE /api/resource/:doctype/:name` | Delete | `CanUserDelete` |
| `GET /api/resource/:doctype` | List | `CanUserRead` + safe col filter + list filters + pagination |
| `GET /api/resource/:doctype/:name` | Get | `CanUserRead` + safe col filter |

Read endpoints:

- `?fields=name,series_key` — whitelist columns (meta + DB intersection)
- `?filters={"series_key":"TEST"}` — SQL WHERE with safe column check
- `?limit=20&offset=0` — pagination
- `?order_by=name&order_dir=asc`

`buildAllowedColumns` — only stored, non-layout, non-system fields from DocType metadata pass through to SQL. Missing columns tracked and returned in response as `missing_columns`.

### API Error Helpers (`internal/core/api/errors.go`)

All error responses use a consistent shape:

```json
{
  "error": {
    "code": "forbidden",
    "message": "Permission denied",
    "detail": "user guest cannot create Naming Series"
  }
}
```

Helpers: `badRequest`, `forbidden`, `notFound`, `serverError`.

### Safe Write Payload Filter (`internal/core/crud/write_safety.go`)

`SafeWritablePayload` filters user input through three gates:
1. **Protected field block** — `name`, `owner`, `creation`, `modified`, `modified_by`, `docstatus`, `idx`, `id`
2. **Metadata gate** — only non-hidden, non-read-only, non-layout fields from DocType meta
3. **DB column gate** — only fields that exist as actual columns in `information_schema`

Returns `Values`, `SkippedFields`, and `MissingColumns` for full transparency.

### DocType Builder & Studio UI
- [x] `/studio` — read-only metadata viewer (DocType list, fields / actions / links / perms / states tabs, JSON preview)
- [x] Visual DocType builder — create, configure, preview, save DocType drafts
  - Module selection, Single/Submittable/Child Table/Editable Grid toggles
  - Add / edit / reorder / remove fields, actions, permissions
  - Live JSON preview with copy button
  - Saves JSON to `modules/{module}/doctype/{slug}/` and syncs to DB

### Slug Helpers
- [x] `slug.FromAppName` — `Gogal Studio` → `gogal_studio`
- [x] `slug.FromModuleName` — `Core` → `core`
- [x] `slug.FromDocTypeName` — `Naming Series` → `naming_series`
- [x] `slug.DocTypeFolderPath` — `modules/core/doctype/naming_series`
- [x] `slug.DocTypeJSONPath` — `modules/core/doctype/naming_series/naming_series.json`

### Field Type Registry
- [x] 24 canonical stored types: Data, Link, Table, Select, Currency, Check, Date, Datetime, Int, Float, Text, Long Text, Small Text, Password, Attach, Attach Image, Read Only, HTML, Markdown Editor, Code, Rating, Duration, Barcode, JSON
- [x] Layout types (never stored, never validated): Section Break, Column Break, Tab Break, Fold, Heading

### Validation Tower
- [x] Field name — `^[a-z][a-z0-9_]*$`, reserved name protection
- [x] DocType name — `^[A-Z][A-Za-z0-9]*( [A-Z][A-Za-z0-9]*)*$`
- [x] Module name — `^[A-Z][A-Za-z0-9]*$`
- [x] App name — `^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$`
- [x] All enforced at bootstrap seed time and on API write

### Dead Code Cleanup *(2026-05-08)*
- [x] `cmd/sync-doctypes/` — deleted (targeted removed `core_*` tables)
- [x] `cmd/server/main.go` — deleted (`cmd/gogal start` is the only entry point)
- [x] `internal/core/api/children.go` — deleted (queried non-existent `core_*` tables)
- [x] `internal/core/api/save.go.bak` — deleted (old backup, not compiled)
- [x] `internal/core/doctype/sync.go` — deleted; `slugify`/`hashBytes` extracted to `helpers.go`

### Unit Tests — 64 passing

| Package | Tests | Count |
|---|---|---|
| `internal/core/migration` | `TestQuoteIdent`, `TestBuildIndexName`, `TestIsSafeApplyOperation` (8 cases), `TestIsSafeWidening` (8 cases), `TestPostgresType` (9 cases) | 7 |
| `internal/core/permission` | `TestActionColumnKnownActions` (12 actions), `TestActionColumnUnknown`, `TestActionColumnEmpty` | 3 |
| `internal/core/crud` | `TestQuoteIdent`, `TestJoinQuotedColumns`, `TestIsLayoutFieldtype` (12 cases), `TestReadableColumns`, `TestBuildWhereClause` (3 cases) | 10 |
| `internal/core/sysfields` | `TestSystemFieldsContainsAll`, `TestIsSystemField` (10 cases), `TestIsProtectedField` (8 cases), `TestInjectCreate`, `TestInjectUpdate` | 6 |
| `internal/core/validator` | `TestIsStoredFieldtype`, `TestValidateRequired*`, `TestValidateSelect*`, `TestValidateSkips*` | 10 |
| `internal/core/doctype` | JSON validation, writer, schema, importer | 28 |

---

## What's Next

### CRUD Hardening
- [ ] Replace temporary name generation with full Naming Series resolution
- [ ] `POST /api/resource/:doctype` — child table nested write support
- [ ] Optimistic lock on `PUT` — check `modified` timestamp before update

### Auth & Sessions
- [ ] User login — `POST /api/auth/login` (JWT + secure cookie)
- [ ] `GET /api/auth/me`
- [ ] Role assignment UI
- [ ] All CRUD endpoints require authenticated user (replace `?user=` query param)

### Studio UI Resource Screen
- [ ] Resource list page — paginated table, column sort, search, filters
- [ ] Resource form page — dynamic form built from DocType field definitions
- [ ] Create / edit / delete from UI

### Feature Packs *(builder UI options)*
- [ ] **Admin CRUD** — full list view + form view for any DocType, one click
- [ ] **REST API** — auto-enable public `GET/POST/PUT/DELETE` for a DocType
- [ ] **Import/Export** — CSV import with validation preview, CSV export
- [ ] **Workflow/States** — state machine, transitions, action buttons, audit log
- [ ] **Dashboard/Reports** — saved reports, charts, aggregations, filters

### AI & Bots *(future)*
- [ ] AI builder — describe app in plain text, AI suggests DocTypes, fields, permissions
- [ ] Bot engine — query/create/update records via natural language (permission-aware)
- [ ] Providers: OpenAI, Gemini, Claude, Ollama/local

### Developer Experience
- [ ] `gogal backup` — database backup to file
- [ ] Code generation — Go structs from DocType definitions
- [ ] React/Vue/Next.js frontend client templates

---

## Project Structure

```
gogal-studio/
├── cmd/
│   └── gogal/              # Main CLI entry point (only entry point)
├── internal/
│   ├── bootstrap/          # Installer: DB setup, core seed data
│   ├── cli/                # Cobra subcommands (start, install, migrate, ...)
│   ├── config/             # Environment config
│   ├── core/
│   │   ├── api/            # HTTP handlers (CRUD, meta, migration, permissions, resource)
│   │   ├── crud/           # Reader + write_safety helpers
│   │   ├── doctype/        # DocType validation, JSON writer/reader/importer, schema
│   │   ├── meta/           # Runtime DocType meta loader (from DB)
│   │   ├── migration/      # Migration planner, safe diff, apply
│   │   ├── naming/         # Naming series engine
│   │   ├── permission/     # Permission checker (CanUserCreate/Read/Write/Delete)
│   │   ├── sysfields/      # System field injection (InjectCreate / InjectUpdate)
│   │   ├── system/         # System-level helpers
│   │   └── validator/      # Document field validator
│   ├── db/                 # Database connection pool
│   └── http/               # Fiber server setup and route registration
├── modules/
│   ├── core/
│   │   └── doctype/        # Core DocType JSON definitions (source of truth)
│   └── security/
│       └── doctype/        # Security module DocType definitions
└── public/
    └── studio/             # Admin UI, builder, resource pages (Vanilla JS)
```

---

## DB Tables (current)

```
tabApp
tabDocField
tabDocPerm
tabDocType
tabHas Role
tabInstalled App
tabInstalled Module
tabModule Def
tabNaming Series
tabRole
tabUser
tabSchema Migration Log    ← written on every migrate apply
```

---

## Design Principles

1. **UI-first** — build apps in the browser, not the terminal
2. **Metadata-first** — the engine runs from DocType definitions, not hand-written code
3. **Preview before apply** — no dangerous automatic schema changes, ever
4. **Security from day one** — permissions enforced on every write, audit logs, safe migrations, input validation at every layer
5. **Consistent error responses** — every API error returns `{error:{code, message, detail}}`
6. **Vanilla JS first** — fast admin UI with zero JS framework dependency
7. **Single binary** — deploy the entire app as one Go binary
8. **Dead code is deleted** — no stubs, no `.bak` files, no unreferenced packages in the tree

---

## Contributing

Gogal Studio is in active early development. The core validation tower, migration planner, and feature pack engine are being built now.

Issues, ideas, and pull requests are welcome.

---

## License

MIT — Galaxy Labs, 2026
