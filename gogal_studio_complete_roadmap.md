# Gogal Studio Roadmap

## 1. Product Vision

Gogal Studio is a Go-based full-stack business application framework designed for developers who like Laravel, Django, Flask, Frappe, Nova, Filament, and modern admin builders, but want Go performance, simple deployment, and a UI-first low-code experience.

The product direction is:

```text
Less commands
More UI actions
More automation
High performance
Security from day one
AI-assisted app creation
Vanilla JS modern admin UI first
Optional code generation later
```

The framework should not become a Laravel clone, a Frappe clone, or a Django clone. It should copy the best productivity ideas and implement them in a Go-native way.

## 2. Main Promise

A user should be able to:

```text
Describe app
Choose features
Preview model and migration
Apply safely
Use app immediately
```

The system should create working features, not just code files.

Example:

```text
Create Customer DocType
Select Admin CRUD + API + Permissions + Import/Export
Preview migration
Apply
Open /desk/customer
Use REST API immediately
```

## 3. Why Not Many Laravel-Style Commands

Laravel generators usually create separate files using separate commands:

```text
make:model
make:controller
make:migration
make:crud
make:api
make:request
make:resource
```

Gogal Studio should avoid that as the main workflow.

Better approach:

```text
Open Builder UI
Create Module / DocType
Choose feature packs
Generate preview
Apply migration
App ready
```

CLI should exist, but only as an admin/developer helper, not the main user experience.

## 4. Product Positioning

Recommended positioning:

```text
A Go full-stack framework for building business apps with UI-first CRUD, DocTypes, migrations, APIs, admin panels, and AI agents.
```

Do not market it only as:

```text
Laravel for Go
```

Better message:

```text
Laravel productivity + Frappe metadata + Django admin simplicity + Go performance.
```

## 5. Core Principles

### 5.1 UI-first

Most app-building actions should happen from the web UI.

### 5.2 Feature packs instead of many packages

A single builder should offer feature options:

```text
Admin CRUD
REST API
Permissions
Import/Export
Workflow/States
Dashboard Cards
Reports
AI Assistant
Bot
Webhook
Background Jobs
Frontend Client
```

### 5.3 Preview before apply

No dangerous automatic schema changes. Always preview migration first.

### 5.4 Security from day one

Permissions, audit logs, migration logs, API key encryption, rate limits, and safe schema changes are required foundation features.

### 5.5 Vanilla JS first

Use a fast, dependency-light Vanilla JS admin UI first. React/Vue/Next can come later through generated clients or templates.

### 5.6 Metadata first, code optional

The engine should run from metadata. Code generation can be optional for advanced developers.

## 6. High-Level Architecture

```text
Gogal Studio
├─ Go API Server
├─ PostgreSQL
├─ Core Metadata Engine
├─ Schema Migration Planner
├─ Dynamic CRUD Engine
├─ Vanilla JS Desk/Admin UI
├─ Builder UI
├─ Permission Engine
├─ Feature Pack Engine
├─ AI Agent Engine
├─ Bot Engine
├─ Jobs/Queue Engine
├─ Reports/Dashboard Engine
└─ CLI for setup/admin tasks
```

## 7. Recommended Tech Stack

```text
Backend: Go
HTTP: Fiber currently, later abstract router if needed
Database: PostgreSQL
DB Driver: pgx
Migrations: Goose now, custom migration planner later
Cache/Queue: Redis + Asynq later
Frontend Admin: Vanilla JS + CSS
Auth: JWT/session hybrid later
AI Providers: OpenAI, OpenRouter, Gemini, Claude, Ollama/local
Deployment: Single binary + config + migrations
```

## 8. Current Foundation Already Started

Current direction already includes:

```text
Go backend
PostgreSQL
Core DocType metadata
JSON sync
Vanilla JS /desk UI
Migration planner direction
Feature-based builder idea
AI agent/bot future
```

Existing plan already says Go is strong for fast backend APIs, multi-tenant servers, auth, permissions, dynamic DocType-like systems, queues, WebSockets, CLI tools, and high-performance services.

## 9. Main User Flow

### 9.1 Create app from UI

```text
Open /desk/builder
Enter module name
Enter DocType name
Add fields
Drag/drop fields
Choose feature packs
Generate preview
Review JSON
Review migration plan
Apply
Use app
```

### 9.2 Create app using AI

```text
Open AI Builder
Write: Create a transport company app
AI suggests modules, DocTypes, fields, permissions, workflows, reports
User reviews
System generates blueprint
Preview migration
Apply
```

### 9.3 Developer CLI flow

```text
gogal serve
gogal migrate
gogal sync
gogal backup
gogal status
```

CLI should be minimal.

## 10. Feature Pack System

Instead of separate generator packages, create one feature engine.

### 10.1 Feature Pack: Database Table

Creates or updates physical PostgreSQL table from DocType fields.

Must include:

```text
Migration preview
Safe create table
Safe add column
No drop by default
No type change by default
Migration log
Rollback metadata only initially
```

### 10.2 Feature Pack: Admin CRUD

Creates dynamic admin UI:

```text
List view
Form view
Create
Edit
Delete
Search
Filters
Pagination
Row actions
Bulk actions
Status pills
```

### 10.3 Feature Pack: REST API

Creates dynamic endpoints:

```text
GET /api/resource/:doctype
POST /api/resource/:doctype
GET /api/resource/:doctype/:id
PUT /api/resource/:doctype/:id
DELETE /api/resource/:doctype/:id
```

### 10.4 Feature Pack: Permissions

Adds role/permission rules:

```text
Read
Write
Create
Delete
Submit
Cancel
Import
Export
Report
Print
Email
Owner-only
Permission level
```

### 10.5 Feature Pack: Import/Export

Adds:

```text
CSV import
CSV export
Excel later
Import validation
Import preview
Error report
Background import job
```

### 10.6 Feature Pack: Workflow/States

Adds:

```text
State definitions
Status colors
Allowed transitions
Transition permissions
Action buttons
Audit log
```

### 10.7 Feature Pack: Dashboard/Reports

Adds:

```text
Dashboard cards
Saved reports
Charts
Aggregations
Filters
Export
```

### 10.8 Feature Pack: AI Assistant

Adds app-level AI helper:

```text
Suggest fields
Suggest validations
Suggest reports
Explain schema
Generate sample data
Suggest permissions
Create help text
```

### 10.9 Feature Pack: Bot

Adds bot/agent tools connected to the app metadata:

```text
Query records
Create records
Update records
Run reports
Explain dashboard
Trigger workflows
```

Bot must obey permissions.

### 10.10 Feature Pack: Webhook/Jobs

Adds:

```text
Outgoing webhooks
Incoming webhooks
Scheduled jobs
Background jobs
Retry queue
Job logs
```

## 11. Admin UI Design Direction

The UI should feel like a mix of:

```text
Laravel Nova
Filament
Django Admin
Flask Admin
Frappe Desk
```

But built in our own Vanilla JS style.

### 11.1 Main Desk Layout

```text
Sidebar
Topbar
Module cards
DocType list
Recent records
Global search
User menu
Notifications
```

### 11.2 Builder Layout

```text
Left: Module/DocType tree
Center: Visual field builder
Right: Properties panel
Bottom/Side: JSON preview and warnings
Top: Generate, Preview, Apply buttons
```

### 11.3 Field Builder

Features:

```text
Add field
Drag/drop order
Field type selector
Required toggle
Unique toggle
Hidden toggle
Read-only toggle
List view toggle
Default value
Options
Depends on
Validation rule
Child table selector
Link target selector
```

### 11.4 Record UI

Features:

```text
Dynamic form renderer
Dynamic list renderer
Child grid
Quick create
Bulk actions
Row actions
Filters
Saved views
Kanban later
Report view later
Dashboard view later
```

## 12. AI Agent Architecture

### 12.1 AI providers

Users can add their own provider/API key:

```text
OpenAI
OpenRouter
Gemini
Claude
Ollama/local
```

### 12.2 Provider storage

AI keys must be encrypted in database.

Suggested tables:

```text
core_ai_provider
core_ai_key
core_ai_agent
core_ai_agent_tool
core_ai_log
```

### 12.3 AI safety rules

AI must never directly apply migrations or delete data without preview and approval.

AI actions must follow:

```text
Suggest
Preview
User approval
Apply
Audit log
```

### 12.4 AI tools

Possible tools:

```text
read_doctype_schema
suggest_doctype_fields
generate_blueprint
preview_migration
create_report
query_records
create_record
update_record
explain_error
```

### 12.5 Bot examples

```text
Sales Bot
Support Bot
Transport Bot
Inventory Bot
HR Bot
Finance Bot
Report Bot
```

Each bot is linked to modules and permissions.

## 13. Security Design From Day One

### 13.1 Authentication

Start with:

```text
Email/password
Password hash
JWT or secure session
Refresh token later
```

Later:

```text
Google login
OTP login
Two-factor auth
API tokens
Personal access tokens
```

### 13.2 Authorization

Every API must pass through permission engine.

Rules:

```text
No direct table access from API without permission check
DocType permissions control CRUD
Field-level permission later
Owner-only rules supported
Bot respects same permission engine
AI respects same permission engine
```

### 13.3 Schema safety

Migration planner must block dangerous changes by default:

```text
Drop column blocked
Rename column requires explicit mapping
Change type blocked unless safe
System tables protected
Core tables protected
Foreign keys delayed until stable
```

### 13.4 Audit logs

Log:

```text
Login attempts
Schema changes
Migration apply
Record create/update/delete
AI actions
Bot actions
API key changes
Permission changes
```

### 13.5 AI key security

Rules:

```text
Encrypt keys
Never expose keys in UI after save
Do not log API keys
Allow key rotation
Allow provider disable
```

## 14. Database Design Roadmap

### 14.1 Core metadata tables

Already planned/started:

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

### 14.2 Builder tables

Add later:

```text
core_app_blueprint
core_generation_request
core_builder_history
core_feature_pack
core_feature_option
```

### 14.3 AI/bot tables

Add later:

```text
core_ai_provider
core_ai_agent
core_ai_agent_tool
core_ai_log
core_bot
core_bot_tool
core_bot_run
```

### 14.4 Site/multi-tenant tables

Add later:

```text
core_site
core_site_database
core_site_installation
core_installed_module
core_site_migration_log
```

## 15. Repo File Structure Target

```text
first-brick/
  cmd/
    server/
    sync-doctypes/
    first-brick/
  internal/
    auth/
    db/
    config/
    http/
    core/
      api/
      doctype/
      module/
      migration/
      generator/
      permissions/
      ai/
      bots/
      jobs/
      reports/
    desk/
    cli/
  modules/
    core/
      module.json
      doctype/
        doctype/
        docfield/
        docperm/
        doctype_action/
        doctype_link/
        doctype_state/
  public/
    index.html
    app.js
    styles.css
    builder.html
    builder.js
    builder.css
  migrations/
  docs/
```

## 16. Should We Use Old Method First Or Advanced From Day One?

Best answer: use the advanced architecture from day one, but implement it in small safe layers.

Do not build the old command-heavy method first and then rewrite. That wastes time.

But also do not build all advanced features at once.

Correct approach:

```text
Advanced architecture
Small implementation steps
Security rules early
Preview-first migrations
UI-first builder
CLI only as helper
```

## 17. Phase Roadmap

### Phase 0: Vision Freeze

Output files:

```text
docs/product-vision.md
docs/architecture.md
docs/security-rules.md
docs/phase-roadmap.md
```

Goal:

```text
Freeze product direction before more coding.
```

### Phase 1: Core Engine Stabilization

Goal:

```text
Core metadata works cleanly.
```

Scope:

```text
Module Def
DocType
DocField
DocType Action
DocType Link
DocPerm
DocType State
JSON sync
Full API read endpoints
Metadata save endpoint
Clean docs
```

Do not yet create business physical tables automatically.

### Phase 2: Builder UI Preview

Goal:

```text
Create Module/DocType visually and generate JSON preview.
```

Scope:

```text
/desk/builder
Field drag/drop
Field properties panel
Feature options panel
JSON preview
Generator preview API
No file write yet
No DB alter yet
```

### Phase 3: JSON File Writer

Goal:

```text
Builder can save DocType JSON in correct module folder.
```

Scope:

```text
Validate module name
Validate DocType name
Generate slug
Create folder
Write modules/{module}/doctype/{slug}/{slug}.json
Sync metadata
Log builder history
```

### Phase 4: Schema Migration Planner

Goal:

```text
Show database changes before applying.
```

Scope:

```text
Compare DocType fields with physical table
Plan create table
Plan add column
Plan indexes
Warn about unsafe changes
No drop column
No type change
Protected core tables
```

### Phase 5: Apply Migration From UI

Goal:

```text
User clicks Apply and table is created safely.
```

Scope:

```text
Apply create table
Apply add column
Add indexes
Record migration log
Show success/failure
```

### Phase 6: Dynamic CRUD API

Goal:

```text
Any migrated DocType automatically gets CRUD API.
```

Scope:

```text
List records
Create record
Read record
Update record
Delete record
Validate required fields
Validate field types
Permission checks
Pagination
Search
Filters
```

### Phase 7: Dynamic Admin UI

Goal:

```text
Any DocType becomes usable in /desk without generated JS files.
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
```

### Phase 8: Permission Engine

Goal:

```text
All UI and API actions become permission-aware.
```

Scope:

```text
Roles
DocPerm
Owner rules
Create/read/write/delete
Import/export/report
Action visibility
Bot/AI permission integration
```

### Phase 9: Feature Pack Engine

Goal:

```text
Feature checkboxes control generated behavior.
```

Scope:

```text
Feature pack registry
Feature dependencies
Feature defaults
Feature warnings
Feature apply hooks
```

### Phase 10: AI Agent Foundation

Goal:

```text
AI can suggest blueprints but cannot directly apply dangerous changes.
```

Scope:

```text
AI provider setup
Encrypted API keys
Prompt templates
Blueprint generation
Schema explanation
Migration explanation
AI logs
```

### Phase 11: Bot Builder

Goal:

```text
Create permission-aware bots for each app/module.
```

Scope:

```text
Bot config
Bot tools
Query records
Create records
Run reports
Bot logs
Permission checks
```

### Phase 12: Reports and Dashboards

Goal:

```text
Business apps get dashboard/report builder.
```

Scope:

```text
Report builder
Chart widgets
Dashboard cards
Saved filters
Export
Scheduled reports
```

### Phase 13: Jobs, Webhooks, Integrations

Goal:

```text
Production automation support.
```

Scope:

```text
Background jobs
Scheduled jobs
Webhook registry
Retry logs
Email/SMS/WhatsApp hooks
```

### Phase 14: Multi-site / SaaS

Goal:

```text
Each site can have its own database and installed modules.
```

Scope:

```text
Control database
Site database creation
Site install
Module install
Backup/restore
Subdomain routing
```

### Phase 15: External Frontend Generators

Goal:

```text
Generate clients for React, Next, Vue, Angular, Flutter, React Native.
```

Scope:

```text
OpenAPI export
TypeScript client
Flutter client
Frontend auth template
Module page template
```

## 18. Performance Design

### 18.1 Use metadata cache

Cache DocType schemas in memory.

### 18.2 Avoid heavy reflection per request

Pre-compile field metadata into runtime structures.

### 18.3 Use pgx pool

Use connection pooling properly.

### 18.4 Add query limits

All dynamic list APIs must have pagination and max limit.

### 18.5 Use background jobs

Long imports, exports, AI generation, report generation should run as jobs.

### 18.6 Avoid npm in core first

Keep Vanilla JS core light and fast.

## 19. Exact Documentation Files To Create

Create these docs before next major coding:

```text
docs/product-vision.md
docs/architecture.md
docs/security-rules.md
docs/phase-roadmap.md
docs/feature-packs.md
docs/ai-agent-bot-plan.md
docs/builder-ui-spec.md
docs/migration-planner-spec.md
docs/repo-structure.md
docs/development-rules.md
```

## 20. Manual File Creation Order

If creating manually, use this order:

```text
1. docs/product-vision.md
2. docs/architecture.md
3. docs/security-rules.md
4. docs/phase-roadmap.md
5. docs/feature-packs.md
6. docs/builder-ui-spec.md
7. docs/migration-planner-spec.md
8. docs/ai-agent-bot-plan.md
9. docs/development-rules.md
```

## 21. Next Implementation Order

Recommended next coding sequence:

```text
1. Create docs roadmap files
2. Add /desk/builder page
3. Add generator blueprint structs
4. Add /api/core/generator/preview
5. Add field drag/drop UI
6. Add feature options UI
7. Add JSON preview
8. Add file writer
9. Add migration planner
10. Add apply migration
```

## 22. Rules For Codex Later

When using Codex later, every prompt should include:

```text
Do not add React/Vue/npm yet.
Do not break existing /desk.
Do not rename module first-brick.
Do not create dangerous migrations.
Do not drop columns.
Do not alter core tables except planned migrations.
Keep Vanilla JS.
Run go fmt ./...
Run go test ./...
Report files changed.
```

## 23. MVP Definition

The first public MVP is complete when this works:

```text
User opens /desk/builder
Creates Module
Creates DocType
Adds fields by UI
Selects Admin CRUD + REST API + Permissions
Generates JSON
Saves JSON to module folder
Syncs metadata
Shows migration preview
Applies migration
Opens dynamic list/form
Creates records
Uses REST API
```

AI is not required for MVP, but the architecture must keep space for AI.

## 24. Final Recommended Direction

Use advanced UI-first architecture from day one, but build it step by step.

Do not build many old Laravel-style commands first.

Do not build AI before migration planner.

Do not create physical tables without preview.

Do build:

```text
Builder UI
Feature packs
Preview engine
Safe migration planner
Dynamic CRUD API
Dynamic admin UI
Permission engine
AI assistant
Bot builder
```

This will make Gogal Studio easier than Laravel generators, more flexible than Django admin, lighter than Frappe, and more attractive for Go adoption.
