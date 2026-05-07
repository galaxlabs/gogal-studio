# Gogal Studio

> **Laravel productivity + Frappe metadata + Django admin simplicity + Go performance.**

Gogal Studio is a Go-powered, UI-first business application framework built by [Galaxy Labs](https://github.com/galaxlabs). It lets you build production-ready business apps by describing your data model in a visual builder — no boilerplate, no scattered generator commands, no manual API wiring.

**The promise:** Describe app → Choose features → Preview migration → Apply → Use app immediately.

---

## Why Gogal Studio?

Most Go frameworks give you a router and leave the rest to you. Gogal Studio gives you:

- A **visual DocType builder** — define your data model in the browser, not in code files
- A **metadata engine** — the entire app runs from JSON-defined DocTypes, not hand-written controllers
- **Feature packs** — attach Admin CRUD, REST API, Permissions, Workflows, Reports, and AI agents to any DocType in one click
- **Safe migration planner** — always preview schema changes before they touch your database
- **Single binary deployment** — ship your entire app as one Go binary with embedded assets

No Laravel-style command explosion. No Frappe Python dependency. No Node.js build pipeline required.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go (Fiber v2) |
| Database | PostgreSQL (pgx v5) |
| Migrations | Goose → custom migration planner |
| Frontend Admin | Vanilla JS + CSS (no framework) |
| Auth | JWT / session hybrid *(planned)* |
| Queue / Cache | Redis + Asynq *(planned)* |
| AI Providers | OpenAI, Gemini, Claude, Ollama *(planned)* |
| Deployment | Single binary + config + migrations |

---

## Architecture

```
Gogal Studio
├── Go API Server (Fiber)
├── PostgreSQL
├── Core Metadata Engine
├── Schema Migration Planner
├── Dynamic CRUD Engine
├── Vanilla JS Desk / Admin UI
├── DocType Builder UI
├── Permission Engine
├── Feature Pack Engine
├── AI Agent Engine          (planned)
├── Bot Engine               (planned)
├── Jobs / Queue Engine      (planned)
├── Reports / Dashboard      (planned)
└── gogal CLI
```

---

## Getting Started

**Prerequisites:** Go 1.22+, PostgreSQL 15+

```bash
# Clone
git clone https://github.com/galaxlabs/gogal-studio
cd gogal-studio

# Install (creates DB, seeds core metadata)
go run ./cmd/gogal install

# Start server
go run ./cmd/gogal start
```

Open:
- **Metadata viewer** → `http://localhost:8080/studio`
- **DocType builder** → `http://localhost:8080/studio-assets/builder.html`
- **Health check** → `http://localhost:8080/health`

---

## CLI Reference

```bash
gogal install          # Bootstrap DB, tables, core seed data
gogal start            # Start the HTTP server
gogal doctor           # Diagnose configuration and DB health

# Slug helpers
gogal slug app gogal_studio              # → gogal_studio
gogal slug module Core                   # → core
gogal slug doctype "Naming Series"       # → naming_series
gogal slug doctype-path Core "Sales Invoice"
# Folder: modules/core/doctype/sales_invoice
# JSON:   modules/core/doctype/sales_invoice/sales_invoice.json

# Field types
gogal fieldtype list
gogal fieldtype show Link

# Naming series
gogal naming list
```

---

## What's Done

### Core Infrastructure
- [x] Go module — `github.com/galaxylabs/gogal-studio`
- [x] Fiber HTTP server with health, version, and API routes
- [x] PostgreSQL connection pool (pgx v5)
- [x] Cobra CLI (`gogal` command)
- [x] Bootstrap installer (`gogal install`) — idempotent, safe to re-run

### Metadata Engine
- [x] Core metadata tables — `tabDocType`, `tabDocField`, `tabDocPerm`, `tabModule Def`, `tabNaming Series`, `tabInstalled App`, `tabInstalled Module`, etc.
- [x] Core DocType JSON files under `modules/core/doctype/`
- [x] DocType sync engine — reads JSON, syncs to DB
- [x] Core DocTypes seeded: Module Def, DocType, DocField, DocPerm, DocType Action, DocType Link, DocType State, Naming Series, Role, User, and more

### Validation Tower *(production-grade)*
- [x] **Field type registry** — 24 canonical types (Data, Link, Table, Select, Currency, Check, Date, etc.)
- [x] **Field name validation** — `^[a-z][a-z0-9_]*$`, reserved name protection (`name`, `owner`, `creation`, etc.)
- [x] **DocType name validation** — `^[A-Z][A-Za-z0-9]*( [A-Z][A-Za-z0-9]*)*$`
- [x] **Module name validation** — `^[A-Z][A-Za-z0-9]*$` (one Title Case word)
- [x] **App name validation** — `^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$` (snake_case)
- [x] All validation enforced at bootstrap seed time — bad data is rejected before it reaches the DB

### Slug Helpers
- [x] `slug.FromAppName` — `gogal_studio` → `gogal_studio`
- [x] `slug.FromModuleName` — `Core` → `core`
- [x] `slug.FromDocTypeName` — `Naming Series` → `naming_series`
- [x] `slug.DocTypeFolderPath` — `modules/core/doctype/naming_series`
- [x] `slug.DocTypeJSONPath` — `modules/core/doctype/naming_series/naming_series.json`

### API Endpoints
- [x] `GET /api/core/doctypes` — list all DocTypes
- [x] `GET /api/core/doctypes/:name` — DocType detail
- [x] `GET /api/core/doctypes/:name/fields`
- [x] `GET /api/core/doctypes/:name/actions`
- [x] `GET /api/core/doctypes/:name/links`
- [x] `GET /api/core/doctypes/:name/permissions`
- [x] `GET /api/core/doctypes/:name/states`
- [x] `GET /fieldtypes` — list all field types
- [x] `GET /fieldtypes/:name` — field type detail

### Builder UI
- [x] `/studio` — read-only metadata viewer (DocType list, field/action/link/perm/state tabs, JSON preview)
- [x] `/studio-assets/builder.html` — visual DocType builder
  - Create DocType draft with module selection
  - Mutually exclusive options (Single / Submittable / Child Table / Editable Grid)
  - Add, edit, reorder, remove fields
  - Add, edit, reorder, remove actions
  - Live JSON preview with copy button
  - Save writes `modules/{module}/doctype/{doctype}/{doctype}.json` and syncs to DB
  - Normal canvas hides internal/meta fields (checkbox to inspect)

---

## What's Next

### Immediate — Core Validation Tower
- [ ] DocType JSON writer (uses slug paths, validated before write)
- [ ] Migration planner — diff current schema vs DocType definition
- [ ] Migration preview UI — show `ALTER TABLE` plan before applying
- [ ] Safe schema apply — `CREATE TABLE`, `ADD COLUMN`, no auto-drop

### Feature Packs *(builder UI options)*
- [ ] **Admin CRUD** — dynamic list view, form view, create, edit, delete, search, filters, pagination
- [ ] **REST API** — `GET/POST/PUT/DELETE /api/resource/:doctype/:id`
- [ ] **Permissions** — role-based read/write/create/delete/submit/cancel per DocType
- [ ] **Import/Export** — CSV import with validation preview, CSV export
- [ ] **Workflow/States** — state machine, allowed transitions, action buttons, audit log
- [ ] **Dashboard/Reports** — saved reports, charts, aggregations, filters

### Auth & Multi-tenancy
- [ ] User login (JWT + session)
- [ ] Role assignment
- [ ] Permission engine enforced on all dynamic CRUD endpoints
- [ ] Multi-site support

### AI & Bots *(future)*
- [ ] AI builder — describe app in plain text, AI suggests modules, DocTypes, fields, permissions
- [ ] App-level AI assistant — suggest fields, validations, reports, help text
- [ ] Bot engine — query/create/update records, run reports, trigger workflows (permission-aware)
- [ ] Providers: OpenAI, Gemini, Claude, Ollama/local

### Developer Experience
- [ ] `gogal migrate` — run pending migrations
- [ ] `gogal sync` — sync all DocType JSON files to DB
- [ ] `gogal backup` — database backup
- [ ] Code generation (optional) — generate Go structs from DocType definitions
- [ ] React/Vue/Next.js frontend client templates

---

## Project Structure

```
gogal-studio/
├── cmd/
│   ├── gogal/          # Main CLI entry point
│   ├── server/         # Standalone server command
│   └── sync-doctypes/  # DocType JSON sync command
├── internal/
│   ├── bootstrap/      # Installer: DB setup, core seed data
│   ├── cli/            # Cobra subcommands
│   ├── config/         # Environment config
│   ├── core/
│   │   ├── app/        # App name validation
│   │   ├── doctype/    # DocType name validation, table naming
│   │   ├── fieldtype/  # Field type registry + field validation
│   │   ├── module/     # Module name validation
│   │   ├── naming/     # Naming series helpers
│   │   ├── slug/       # Slug helpers for folder paths
│   │   └── api/        # Core API route handlers
│   ├── db/             # Database connection pool
│   └── http/           # Fiber server setup
├── modules/
│   └── core/
│       └── doctype/    # Core DocType JSON definitions
└── public/
    └── studio/         # Admin UI and builder assets
```

---

## Design Principles

1. **UI-first** — build apps in the browser, not the terminal
2. **Metadata-first** — the engine runs from DocType definitions, not hand-written code
3. **Preview before apply** — no dangerous automatic schema changes, ever
4. **Security from day one** — permissions, audit logs, safe migrations, input validation at every layer
5. **Feature packs over scattered generators** — one builder, one click, all layers created together
6. **Vanilla JS first** — fast admin UI with zero JS framework dependency
7. **Single binary** — deploy your entire app as one Go binary

---

## Contributing

Gogal Studio is in active early development. The core validation tower, migration planner, and feature pack engine are being built now.

Issues, ideas, and pull requests are welcome.

---

## License

MIT — Galaxy Labs, 2026
