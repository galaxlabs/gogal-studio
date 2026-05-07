# Gogal Studio Roadmap v2 — Current Build Status

Updated: 2026-05-06  
Project path: `E:\dev\gogal-studio`  
Brand/company: `Galaxy Labs`  
Product: `Gogal Studio`  
Engine: `Gogal Engine`  
CLI: `gogal`  
Go module: `github.com/galaxylabs/gogal-studio`

---

## 1. Product Vision

Gogal Studio is a Go-based full-stack business application framework for building business apps with a UI-first, metadata-first, safe-migration-first workflow.

It is inspired by Laravel, Django, Flask, Frappe, Nova, Filament, and modern admin builders, but it must not become a clone of any of them.

The core direction remains:

```text
Less commands
More UI actions
More automation
High performance
Security from day one
Vanilla JS modern admin first
Metadata-driven app creation
Preview before apply
Optional AI-assisted app creation later
```

Main promise:

```text
Describe app → Choose features → Preview → Apply → Use app
```

The system should create working features, not just code files.

---

## 2. Current Implementation Status

### 2.1 Repository foundation — Done

```text
E:\dev\gogal-studio
```

Completed:

```text
git init
go mod init github.com/galaxylabs/gogal-studio
Fiber server
Clean internal server split
PostgreSQL config
pgxpool database connection
Health route
Version route
```

Current server files:

```text
cmd/server/main.go
internal/config/config.go
internal/db/db.go
internal/http/server.go
internal/http/routes.go
```

Working routes:

```text
GET /
GET /health
GET /api
GET /api/version
```

Current development database:

```text
postgres://door:door@localhost:5432/door_app?sslmode=disable
```

---

### 2.2 Core metadata foundation — Done

Core metadata tables are created through Goose migrations.

Current core metadata tables:

```text
core_module_def
core_doctype
core_docfield
core_doctype_action
core_doctype_link
core_docperm
core_doctype_state
core_migration_log
```

Core DocType JSON files are created:

```text
modules/core/doctype/module_def/module_def.json
modules/core/doctype/doctype/doctype.json
modules/core/doctype/docfield/docfield.json
modules/core/doctype/doctype_action/doctype_action.json
modules/core/doctype/doctype_link/doctype_link.json
modules/core/doctype/docperm/docperm.json
modules/core/doctype/doctype_state/doctype_state.json
```

Core DocType schema structs are created:

```text
internal/core/doctype/schema.go
internal/core/doctype/docfield.go
internal/core/doctype/doctype_action.go
internal/core/doctype/doctype_link.go
internal/core/doctype/docperm.go
internal/core/doctype/doctype_state.go
```

DocType sync engine is created:

```text
internal/core/doctype/sync.go
cmd/sync-doctypes/main.go
```

Sync result:

```text
Success: 7
Failed: 0
```

Synced Core DocTypes:

```text
Module Def
DocType
DocField
DocType Action
DocType Link
DocPerm
DocType State
```

---

### 2.3 Core metadata API — Done

Working API endpoints:

```text
GET /api/core/doctypes
GET /api/core/doctypes/:name
GET /api/core/doctypes/:name/fields
GET /api/core/doctypes/:name/actions
GET /api/core/doctypes/:name/links
GET /api/core/doctypes/:name/permissions
GET /api/core/doctypes/:name/states
```

Important route rule:

```text
DocType names with spaces must be URL-decoded before DB lookup.
Example: Module%20Def → Module Def
```

---

### 2.4 Studio metadata viewer — Done

Current read-only metadata UI:

```text
GET /studio
/static assets: /studio-assets/*
```

Files:

```text
public/studio/index.html
public/studio/styles.css
public/studio/app.js
```

Working features:

```text
Sidebar DocType list
DocType detail panel
Fields tab
Actions tab
Links tab
Permissions tab
States tab
JSON Preview tab
Counts summary
```

Important rule:

```text
/studio is the stable metadata viewer.
Do not break it while building the new builder UI.
```

---

### 2.5 Separate DocType Builder UI — In progress

Separate builder page:

```text
/studio-assets/builder.html
```

Files:

```text
public/studio/builder.html
public/studio/builder.css
public/studio/builder.js
public/studio/builder-create.js
public/studio/builder-policy.js
```

Purpose:

```text
Visual DocType builder
Draft UI only for now
No physical DB table creation
No migration apply yet
```

Current / planned builder layout:

```text
Left sidebar:
  DocType list
  Create DocType button
  Back to Metadata button

Center canvas:
  Form canvas
  Field cards
  Drag/drop field ordering
  Action tab/panel
  JSON Preview

Right properties panel:
  DocType properties
  Field properties
  Action properties
```

---

## 3. Updated Builder UI Design Decisions

### 3.1 DocType create flow

When creating a new DocType, ask:

```text
DocType Name
Module Def
Create New Module Def option
Is Single
Is Submittable
Is Child Table / istable
Editable Grid
```

Module Def rules:

```text
Module Def is selected at DocType creation.
Module Def can be changed later in DocType properties.
Module Def is a DocType property.
Module Def must NOT appear on the main field canvas.
```

---

### 3.2 DocType option interaction rules

#### Is Submittable

Explanation:

```text
Once submitted, submittable documents cannot be changed.
They can only be Cancelled and Amended.
```

When `Is Submittable = true`:

```text
Disable/hide Is Single
Disable/hide Is Child Table / istable
Disable/hide Editable Grid
```

#### Is Child Table / istable

Explanation:

```text
Child Tables are shown as a Grid in other DocTypes.
```

When `Is Child Table / istable = true`:

```text
Disable/hide Is Single
Disable/hide Is Submittable
Show/enable Editable Grid
Default Editable Grid = true
```

#### Editable Grid

When `Editable Grid = true`:

```text
Force Is Child Table / istable = true
Disable/hide Is Single
Disable/hide Is Submittable
```

#### Is Single

When `Is Single = true`:

```text
Disable/hide Is Submittable
Disable/hide Is Child Table / istable
Disable/hide Editable Grid
```

---

## 4. Builder Canvas Policy

The canvas is for real business/user fields only.

### 4.1 Show by default on canvas

Examples:

```text
customer_name
mobile
email
status
description
posting_date
amount
remarks
```

### 4.2 Hide from normal canvas

System/meta fields must stay hidden unless inspection mode is enabled.

Hide these by default:

```text
name
owner
creation
modified
modified_by
docstatus
idx
parent
parenttype
parentfield
doctype
module
is_submittable
istable
is_child_table
issingle
is_single
editable_grid
quick_entry
track_changes
track_seen
track_views
custom
beta
is_virtual
queue_in_background
engine
migration_hash
row_format
permissions
actions
links
states
fields
form_builder
settings_tab
connections_tab
advanced
json_hash
source_path
created_at
updated_at
status
oldfieldname
oldfieldtype
```

Add UI checkbox:

```text
Show meta/internal fields
```

Default:

```text
Unchecked
```

---

## 5. Field Builder Requirements

The field builder must support:

```text
Add Field
Select field card
Edit label
Edit fieldname for draft/new field
Edit fieldtype
Edit options
Toggle required
Toggle unique
Toggle hidden
Toggle read-only
Toggle in list view
Set default value
Set description
Set depends_on
Drag/drop reorder
Remove non-system field
```

Field types:

```text
Data
Small Text
Text
Long Text
Int
Float
Currency
Check
Date
Datetime
Time
Select
Link
Table
Attach
Attach Image
JSON
Code
Section Break
Column Break
Tab Break
Button
HTML
```

Field type behavior:

```text
Link      → Options = target DocType
Table     → Options = child table DocType
Select    → Options textarea, one option per line
Breaks    → Layout-only cards, not normal input fields
Button    → Action/handler configured from Actions panel later
```

---

## 6. Action UI Requirements

Actions are separate from fields.

Actions must not appear as field cards.

Action builder must support:

```text
Add Action
Select Action
Edit Label
Edit Action Name
Action Type
Method
Handler
Route
Permission
Visible When
Enabled toggle
Drag/drop reorder
Remove Action
```

Action types:

```text
server
client
route
modal
external
```

HTTP methods:

```text
GET
POST
PUT
PATCH
DELETE
```

Action JSON target:

```json
"actions": []
```

Later sync target:

```text
core_doctype_action
```

---

## 7. JSON Preview Requirements

Builder must generate formatted JSON from current draft state.

JSON must include:

```text
name
module
label
table_name
is_core
is_single
is_submittable
is_child_table
editable_grid
quick_entry
allow_import
allow_export
track_changes
naming_rule
title_field
sort_field
sort_order
fields
actions
links
permissions
states
```

Required UI:

```text
JSON Preview tab/panel
Copy JSON button
Live update after field/action edits
```

---

## 8. Current Save/Migration Rule

For current builder phase:

```text
No physical DB table creation
No schema migration apply
No dangerous backend changes
No automatic ALTER TABLE
```

Allowed now:

```text
Draft UI state
JSON preview
Optional safe metadata save endpoint later
```

Later safe flow:

```text
Create/Edit DocType draft
Preview JSON
Save JSON to module folder
Sync metadata
Preview migration plan
User clicks Apply
Create/update physical table safely
Log migration
```

---

## 9. Updated Repo Structure Target

Current product path must be Gogal Studio, not First Brick.

```text
gogal-studio/
  cmd/
    gogal/
    server/
    sync-doctypes/
  internal/
    app/
    auth/
    config/
    db/
    http/
    studio/
    core/
      api/
      ai/
      bots/
      crud/
      doctype/
      generator/
      migration/
      module/
      orm/
      permissions/
      reports/
  modules/
    core/
      module.json
      doctype/
        module_def/
        doctype/
        docfield/
        doctype_action/
        doctype_link/
        docperm/
        doctype_state/
  migrations/
  public/
    studio/
      index.html
      styles.css
      app.js
      builder.html
      builder.css
      builder.js
      builder-create.js
      builder-policy.js
  docs/
  templates/
  scripts/
  README.md
  .gitignore
  go.mod
  go.sum
```

---

## 10. Updated Phase Roadmap

### Phase 0: Vision and package direction — Done

Completed:

```text
Product name: Gogal Studio
Brand: Galaxy Labs
Engine: Gogal Engine
CLI: gogal
UI-first product direction
Feature-pack approach
Metadata-first design
Preview-before-apply rule
```

---

### Phase 1: Core Engine Stabilization — Mostly done

Completed:

```text
Go server foundation
PostgreSQL connection
Core metadata migrations
7 Core DocType JSON files
Core schema structs
DocType JSON sync engine
Core metadata API read endpoints
/studio metadata viewer
JSON preview in viewer
```

Still pending in Phase 1:

```text
Safe metadata save endpoint
README update with current state
More validation around DocType JSON payloads
Basic automated tests for Core API/sync
```

---

### Phase 2: Builder UI Draft Interaction — Current active phase

Goal:

```text
Create/edit DocType visually and generate JSON preview without saving physical schema.
```

Scope:

```text
Separate builder page
DocType create flow
Module Def select/create
DocType option rules
Field add/edit/remove
Field drag/drop
Field properties panel
Action add/edit/remove
Action drag/drop
Action properties panel
Canvas field visibility policy
Show meta/internal fields toggle
JSON preview
Copy JSON button
No DB alter
No migration apply
```

Status:

```text
In progress
```

---

### Phase 3: Builder Metadata Save

Goal:

```text
Save created/edited DocType metadata safely.
```

Scope:

```text
POST /api/core/doctypes
Validate DocType payload
Validate field payload
Validate action payload
Log source as ui://doctype/{DocTypeName}
Sync metadata using existing SyncOne()
Do not create physical table yet
```

---

### Phase 4: JSON File Writer

Goal:

```text
Builder can write DocType JSON into the correct module folder.
```

Scope:

```text
Validate module name
Create module folder if needed
Generate module slug
Generate doctype slug
Create modules/{module}/doctype/{doctype_slug}/{doctype_slug}.json
Sync metadata after write
Log builder history
```

---

### Phase 5: Schema Migration Planner

Goal:

```text
Show database changes before applying.
```

Scope:

```text
POST /api/core/migration/preview
Compare DocType fields with PostgreSQL table
Plan create table
Plan add column
Plan indexes
Block drop column
Block unsafe type change
Protect core tables
Return human-readable migration plan
```

---

### Phase 6: Apply Migration From UI

Goal:

```text
User applies safe migration after preview.
```

Scope:

```text
POST /api/core/migration/apply
Create table safely
Add columns safely
Add indexes safely
Record migration log
Show success/failure in UI
```

---

### Phase 7: Dynamic CRUD API

Goal:

```text
Any migrated DocType gets runtime CRUD API.
```

Routes:

```text
GET    /api/resource/:doctype
POST   /api/resource/:doctype
GET    /api/resource/:doctype/:name
PUT    /api/resource/:doctype/:name
DELETE /api/resource/:doctype/:name
```

Scope:

```text
Required field validation
Field type validation
Pagination
Search
Filters
Permission checks later
```

---

### Phase 8: Dynamic Admin UI

Goal:

```text
Any migrated DocType becomes usable from Gogal admin UI.
```

Scope:

```text
Dynamic list view
Dynamic form view
Create/edit/save/delete
Child table grid
Filters
Bulk actions
Row actions
Status pills
Column settings
```

---

### Phase 9: Permission Engine

Goal:

```text
UI and API become permission-aware.
```

Scope:

```text
Roles
DocPerm enforcement
Owner rules
Create/read/write/delete
Import/export/report
Action visibility
Field-level permission later
```

---

### Phase 10: Feature Pack Engine

Goal:

```text
Feature checkboxes control generated behavior.
```

Feature packs:

```text
Database Table
Admin CRUD
REST API
Permissions
Import/Export
Workflow/States
Dashboard Cards
Reports
Public API
TypeScript Client later
AI Assistant later
Bot later
Webhook/Jobs later
```

---

### Phase 11: AI Agent Foundation

Goal:

```text
AI suggests app blueprints but cannot apply dangerous changes directly.
```

Rule:

```text
Suggest → Preview → User Approval → Apply → Audit Log
```

---

### Phase 12: Bot Builder

Goal:

```text
Create permission-aware bots for each app/module.
```

---

### Phase 13: Reports and Dashboards

Goal:

```text
Business apps get report/dashboard builder.
```

---

### Phase 14: Jobs, Webhooks, Integrations

Goal:

```text
Production automation support.
```

---

### Phase 15: Multi-site / SaaS

Goal:

```text
Each site can have its own database and installed modules.
```

---

### Phase 16: External Frontend Generators

Goal:

```text
Generate clients/templates for React, Next, Vue, Angular, Flutter, React Native.
```

---

## 11. Updated Next Implementation Order

Current immediate order:

```text
1. Finish builder UI interaction
2. Make field cards editable from UI
3. Add field add/remove/reorder
4. Add action add/remove/reorder
5. Add JSON Preview + Copy JSON
6. Update README with current builder state
7. Add safe POST /api/core/doctypes metadata save
8. Add JSON file writer
9. Add migration preview endpoint
10. Add migration apply endpoint
```

Do not jump to:

```text
React/Vue/npm UI libraries
Physical table creation before preview
AI before migration planner
Dynamic CRUD before migration planner
Multi-site before core app flow works
```

---

## 12. Updated Codex Rules

Every Codex task should include:

```text
Read README.md first.
Read roadmap/docs first.
Do not use First Brick paths.
Use E:\dev\gogal-studio.
Use imports github.com/galaxylabs/gogal-studio/...
Do not break /studio metadata viewer.
Keep builder separate until stable.
Keep Vanilla JS first.
Do not add React/Vue/npm unless explicitly requested.
Do not create physical DB tables from builder yet.
Do not apply migrations automatically.
Do not drop columns.
Do not alter core tables except planned migrations.
Run go fmt ./...
Run go test ./...
Report files changed.
Report remaining limitations.
```

For current builder tasks, Codex must focus on:

```text
DocType UI interaction
Field editing
Action UI setup
Canvas visibility policy
JSON preview
README update
```

---

## 13. Updated README Section To Keep

README should contain:

```md
## Current Focus: DocType Builder UI

Gogal Studio currently has:

- Go + Fiber backend
- PostgreSQL connection
- Core metadata tables
- Core DocType JSON sync
- Core metadata API
- Read-only /studio metadata viewer
- Separate /studio-assets/builder.html builder page

The current builder page is draft-only.
It supports or is being updated to support:

- Create DocType draft
- Select existing Module Def or create a new Module Def name
- Mutually exclusive DocType options
- Field add/edit/reorder
- Action add/edit/reorder
- JSON Preview

The builder does not yet apply schema migrations or create physical database tables.
```

---

## 14. MVP Definition — Updated

The first MVP is complete when this works:

```text
User opens builder
Creates Module Def
Creates DocType
Adds fields by UI
Adds actions by UI
Selects feature options
Reviews JSON preview
Saves JSON to module folder
Syncs metadata
Shows migration preview
Applies migration safely
Opens dynamic list/form
Creates records
Uses REST API
```

AI is not required for MVP.

---

## 15. Final Direction

Continue with the advanced UI-first architecture, but keep implementation small and safe.

Build now:

```text
DocType Builder UI
Field editing
Action editing
JSON preview
Safe metadata save
JSON file writer
Migration preview
```

Build later:

```text
Physical migration apply
Dynamic CRUD
Dynamic Admin UI
Permissions
Feature packs
AI
Bots
Reports
Jobs
Multi-site
Frontend generators
```

