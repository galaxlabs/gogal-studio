# Gogal Studio Framework Packages

Gogal Studio is a Go-powered full-stack business application framework by **Galaxy Labs**.

It is not only a UI builder. It is a complete framework engine for building, generating, migrating, securing, and running business applications.

Official naming:

```text
Brand/company: Galaxy Labs
Product: Gogal Studio
Engine: Gogal Engine
CLI: gogal
Go module: github.com/galaxylabs/gogal-studio
Local path: E:\dev\gogal-studio
```

Main product promise:

```text
Describe app → Choose features → Preview → Apply → Use app
```

---

## 1. Framework Direction

Gogal Studio should combine the productivity ideas of Frappe, Django Admin, Laravel Nova/Filament, and modern low-code builders, but implemented in a Go-native way.

Primary direction:

```text
UI-first builder
Feature packs
Metadata-driven DocTypes
Migration preview before apply
Dynamic CRUD API
Dynamic Admin UI
PostgreSQL backend
Future multi-site/multi-database support
Future AI agents and bots
Secure from day one
```

---

## 2. Main Design Rule

Do not make the main workflow command-heavy.

Avoid making this the primary workflow:

```text
gogal make:model
gogal make:controller
gogal make:crud
gogal make:api
gogal make:admin
gogal make:migration
```

Preferred workflow:

```text
Open Gogal Studio
Choose or create Module
Create DocType
Add fields
Choose feature packs
Preview generated JSON
Preview migration
Apply migration
Use generated app
```

CLI still exists, but mainly for developer/admin operations.

---

## 3. Package Naming Rule

Packages are internal framework areas, not necessarily separate public repositories at the beginning.

Initial internal package style:

```text
internal/config
internal/db
internal/http
internal/core/doctype
internal/core/module
internal/core/migration
internal/core/generator
internal/core/permissions
internal/core/api
internal/studio
internal/auth
internal/users
```

Future public package names can be:

```text
gogal-core
gogal-cli
gogal-doctype
gogal-module
gogal-migration
gogal-generator
gogal-crud
gogal-admin
gogal-auth
gogal-permissions
gogal-api
gogal-reports
gogal-jobs
gogal-events
gogal-files
gogal-email
gogal-notifications
gogal-realtime
gogal-ai
gogal-bots
```

---

## 4. Package Map

## 4.1 gogal-config

Path:

```text
internal/config
```

Purpose:

Application and environment configuration.

Responsibilities:

```text
APP_NAME
APP_ENV
APP_PORT
DATABASE_URL
Future CONTROL_DATABASE_URL
Future REDIS_URL
Future SITE_MODE
Future LOG_LEVEL
```

Current status:

```text
Started
```

---

## 4.2 gogal-db

Path:

```text
internal/db
```

Purpose:

Database connection and database utilities.

Responsibilities:

```text
PostgreSQL pgxpool connection
Database health check
Transaction helper
Future site database resolver
Future control database resolver
Future query helper
```

Current status:

```text
Started
```

Important rule:

The current `door_app` database is only a development site database. Later Gogal Studio should support one control database and many site databases.

---

## 4.3 gogal-http

Path:

```text
internal/http
```

Purpose:

HTTP server, routes, middleware, and static UI serving.

Responsibilities:

```text
Fiber server setup
Health route
API root
Version route
Studio UI route
Static file serving
Route grouping
Middleware registration
Error handling
```

Current status:

```text
Started
```

---

## 4.4 gogal-core

Path:

```text
internal/core
```

Purpose:

Core engine foundation.

Responsibilities:

```text
DocType metadata
Module system
Migration planner
Permission engine
Dynamic CRUD
Core APIs
Generator engine
Feature packs
```

Current status:

```text
Started
```

---

## 4.5 gogal-doctype

Path:

```text
internal/core/doctype
```

Purpose:

Metadata-driven DocType system.

This is the backbone of Gogal Studio.

Responsibilities:

```text
DocType JSON schema
DocField JSON schema
DocType Action schema
DocType Link schema
DocPerm schema
DocType State schema
JSON loader
JSON validator
JSON hash calculation
Metadata sync engine
Metadata migration log
```

Required bootstrap tables:

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

Required Core JSON files:

```text
modules/core/module.json
modules/core/doctype/module_def/module_def.json
modules/core/doctype/doctype/doctype.json
modules/core/doctype/docfield/docfield.json
modules/core/doctype/doctype_action/doctype_action.json
modules/core/doctype/doctype_link/doctype_link.json
modules/core/doctype/docperm/docperm.json
modules/core/doctype/doctype_state/doctype_state.json
```

Current status:

```text
Next main package to complete
```

---

## 4.6 gogal-module

Path:

```text
internal/core/module
```

Purpose:

Module/app package system.

Responsibilities:

```text
module.json support
Module discovery
Module install
Module uninstall
Module enable/disable
Module dependency checking
Module versioning
Installed module records
```

Target structure:

```text
modules/
  core/
    module.json
    doctype/
  selling/
    module.json
    doctype/
  buying/
    module.json
    doctype/
  inventory/
    module.json
    doctype/
  transport/
    module.json
    doctype/
```

Current status:

```text
Planned
```

---

## 4.7 gogal-migration

Path:

```text
internal/core/migration
```

Purpose:

Schema migration planner and applier.

Responsibilities:

```text
Read DocType metadata
Compare metadata with PostgreSQL schema
Create migration preview
Create table plan
Add column plan
Alter column plan
Index plan
Constraint plan
Apply migration
Log migration result
Future rollback support
```

Important rule:

Do not blindly alter business tables on every DocType save.

Safe flow:

```text
DocType saved
↓
Migration preview generated
↓
User reviews changes
↓
User clicks Apply
↓
Schema changes applied
↓
Migration log written
```

Current status:

```text
Planned after Core metadata bootstrap
```

---

## 4.8 gogal-generator

Path:

```text
internal/core/generator
```

Purpose:

Feature-based generator engine.

This replaces many Laravel-style generator commands.

Responsibilities:

```text
Generate DocType JSON
Generate Module JSON
Generate feature options
Generate CRUD metadata
Generate Admin List metadata
Generate Admin Form metadata
Generate API metadata
Generate permission defaults
Generate migration plan input
Generate route/menu metadata
```

Feature options:

```text
Database Table
Admin CRUD
REST API
Permissions
Import/Export
List View
Form View
Workflow/States
Dashboard Cards
Public API
TypeScript Client
AI Assistant
Bot
Webhook
Background Jobs
```

Current status:

```text
Planned
```

---
---

### gogal-orm

Path:

```text
internal/core/orm


## 4.9 gogal-crud

Path:

```text
internal/core/crud
```

Purpose:

Dynamic CRUD API for all DocTypes.

Responsibilities:

```text
List records
Create record
Read record
Update record
Delete record
Validate fields from DocField metadata
Apply naming rules
Enforce permissions
Support filters
Support sorting
Support pagination
Support owner tracking
```

Future routes:

```text
GET    /api/resource/:doctype
POST   /api/resource/:doctype
GET    /api/resource/:doctype/:name
PUT    /api/resource/:doctype/:name
DELETE /api/resource/:doctype/:name
```

Current status:

```text
Planned after migration planner
```

---

## 4.10 gogal-admin

Path:

```text
internal/studio/admin
public/studio/admin
```

Purpose:

Dynamic admin/list/form engine.

Responsibilities:

```text
Dynamic list view
Dynamic form view
Child table grid
Quick create
Filters
Search
Bulk actions
Row actions
Status pills
Column settings
Permission-aware buttons
```

Current status:

```text
Planned
```

---

## 4.11 gogal-studio

Path:

```text
internal/studio
public/studio
```

Purpose:

Main visual builder and admin UI.

Responsibilities:

```text
Studio shell
Builder UI
Module builder
DocType builder
Field drag/drop
Feature options
JSON preview
Migration preview
Apply migration
Dynamic admin list/form
```

Current status:

```text
Planned after Core APIs
```

---

## 4.12 gogal-auth

Path:

```text
internal/auth
```

Purpose:

Authentication system.

Responsibilities:

```text
Signup
Login
Logout
JWT
Refresh tokens
Sessions
Password hashing
Password reset
Email verification
Google login later
OTP login later
```

Current status:

```text
Planned
```

---

## 4.13 gogal-users

Path:

```text
internal/users
```

Purpose:

User, profile, and account management.

Responsibilities:

```text
User profile
Avatar
User status
User settings
User roles
Owner tracking
Future User DocType integration
```

Current status:

```text
Planned
```

---

## 4.14 gogal-permissions

Path:

```text
internal/core/permissions
```

Purpose:

Permission engine.

Responsibilities:

```text
Role permissions
DocPerm rules
Owner permissions
Field-level permissions
Row-level permissions
Action permissions
API permissions
Permission-aware UI buttons
```

Current status:

```text
Planned
```

---

## 4.15 gogal-api

Path:

```text
internal/core/api
```

Purpose:

Core API and generated API system.

Responsibilities:

```text
Core metadata APIs
Dynamic resource APIs
Custom API registry
OpenAPI generation later
TypeScript client generation later
Public/private API controls
Webhook API support
```

Current status:

```text
Partially planned
```

---

## 4.16 gogal-reports

Path:

```text
internal/core/reports
```

Purpose:

Reports and dashboards.

Responsibilities:

```text
Query reports
Saved reports
Dashboard cards
Charts
Filters
Export
Future report builder
```

Current status:

```text
Planned
```

---

## 4.17 gogal-jobs

Path:

```text
internal/jobs
```

Purpose:

Background jobs and scheduler.

Responsibilities:

```text
Queue jobs
Scheduled jobs
Email queue
Import/export jobs
Retry failed jobs
Job dashboard
Worker process
```

Future dependency:

```text
Redis
Asynq or own queue layer
```

Current status:

```text
Planned
```

---

## 4.18 gogal-events

Path:

```text
internal/events
```

Purpose:

Framework hooks/events.

Responsibilities:

```text
before_insert
after_insert
before_save
after_save
before_submit
after_submit
before_delete
after_delete
Custom module hooks
```

Current status:

```text
Planned
```

---

## 4.19 gogal-files

Path:

```text
internal/files
```

Purpose:

File manager and attachments.

Responsibilities:

```text
File upload
Private files
Public files
File metadata
Attach field support
Storage drivers later
```

Current status:

```text
Planned
```

---

## 4.20 gogal-email

Path:

```text
internal/email
```

Purpose:

Email system.

Responsibilities:

```text
SMTP settings
Email templates
Email queue
Password reset email
System notifications
Document notifications
```

Current status:

```text
Planned
```

---

## 4.21 gogal-notifications

Path:

```text
internal/notifications
```

Purpose:

Notification system.

Responsibilities:

```text
In-app notifications
Email notifications
Webhook notifications
Realtime notifications later
User notification preferences
```

Current status:

```text
Planned
```

---

## 4.22 gogal-realtime

Path:

```text
internal/realtime
```

Purpose:

Realtime/WebSocket layer.

Responsibilities:

```text
Live notifications
Live form updates
Queue updates
Chat/bot messages later
Realtime dashboard updates
```

Current status:

```text
Planned
```

---

## 4.23 gogal-importer

Path:

```text
internal/importer
```

Purpose:

Import/export system.

Responsibilities:

```text
CSV import
Excel import later
Data validation
Import preview
Error rows
Export CSV
Export Excel later
```

Current status:

```text
Planned
```

---

## 4.24 gogal-ai

Path:

```text
internal/core/ai
```

Purpose:

AI provider and app-generation assistance.

Responsibilities:

```text
AI provider settings
OpenAI/OpenRouter/Gemini/Ollama support later
Prompt templates
Generate DocTypes
Suggest fields
Suggest validations
Suggest permissions
Explain migration plans
Generate reports
```

Current status:

```text
Planned after migration planner
```

---

## 4.25 gogal-bots

Path:

```text
internal/core/bots
```

Purpose:

Bot/agent runtime.

Responsibilities:

```text
Bot registry
Bot tools
Permission-aware data access
App assistant
Report assistant
CRUD assistant
Future chat interface
```

Current status:

```text
Planned after AI package
```

---

## 5. External Dependencies Policy

Use fewer external dependencies at the beginning.

Allowed foundation dependencies:

```text
Fiber       → HTTP server
pgxpool     → PostgreSQL connection
Goose       → migrations
JWT         → auth token later
bcrypt      → password hashing later
Cobra       → CLI later
```

Later dependencies:

```text
Redis/go-redis
Asynq
Validator
Zap/Zerolog
OpenAPI generator
QR code package
PDF service/package
```

Rule:

```text
Do not add packages before they are needed.
```

---

## 6. Feature Pack Concept

Users should choose feature packs from Studio UI.

Example packs:

```text
Admin CRUD Pack
API Pack
Permission Pack
Import/Export Pack
Workflow Pack
Report Pack
Dashboard Pack
AI Assistant Pack
Bot Pack
Mobile Client Pack
```

Each pack should add required metadata automatically.

Admin CRUD Pack creates:

```text
List View
Form View
Create/Edit/Delete actions
Filters
Search
Bulk actions
Menu entry
Permissions
```

API Pack creates:

```text
REST resource routes
Permission checks
OpenAPI metadata later
API documentation metadata
```

Permission Pack creates:

```text
DocPerm records
Role defaults
Owner rules
Field-level rules later
```

Workflow Pack creates:

```text
States
Transitions
Allowed actions
Status field
Submit/cancel rules
```

AI Pack creates:

```text
AI assistant metadata
Prompt templates
Tool access plan
Permission boundaries
```

---

## 7. Site and Database Architecture

Gogal Studio must be designed for future multi-site/multi-database use.

There are two levels:

```text
Control database
Site database
```

Control database stores:

```text
Sites
Installations
Site database names
Site status
Installed modules per site
Global license/settings later
```

Each site database stores:

```text
Core metadata tables
Installed module metadata
Business tables
Site migration logs
Site users/roles later
```

Current development database:

```text
door_app
```

This is only the first development site database, not the final global control database.

Future site creation flow:

```text
Create site record in control DB
Create new PostgreSQL database
Run bootstrap migrations
Sync Core module JSON
Install selected modules
Record installation
Open Studio for that site
```

---

## 8. Core Metadata Bootstrap Rule

Core metadata tables are not normal business tables.

They are engine tables.

Required engine tables:

```text
core_module_def
core_doctype
core_docfield
core_doctype_action
core_doctype_link
core_docperm
core_doctype_state
core_migration_log
core_site_installation
core_installed_module
```

Core DocTypes still exist as metadata:

```text
Core → Module Def
Core → DocType
Core → DocField
Core → DocType Action
Core → DocType Link
Core → DocPerm
Core → DocType State
```

But their physical storage tables are fixed bootstrap tables:

```text
core_module_def
core_doctype
core_docfield
core_doctype_action
core_doctype_link
core_docperm
core_doctype_state
```

---

## 9. JSON Storage Rule

Every DocType must have its own folder.

Correct structure:

```text
modules/{module_slug}/doctype/{doctype_slug}/{doctype_slug}.json
```

Examples:

```text
modules/core/doctype/docfield/docfield.json
modules/selling/doctype/customer/customer.json
modules/transport/doctype/builty/builty.json
```

Avoid all-in-one model files.

Each DocType JSON should be independent, portable, syncable, and hashable.

---

## 10. Migration Safety Rule

Metadata sync and physical schema migration are separate.

Metadata sync:

```text
JSON file
↓
core_doctype/core_docfield/etc.
```

Physical migration:

```text
core_doctype/core_docfield
↓
Migration planner
↓
Preview
↓
Apply
↓
Actual PostgreSQL business table
```

Do not auto-create or auto-alter business tables blindly when user edits a DocType.

---

## 11. Development Order

Build packages in this order:

```text
1. gogal-config
2. gogal-db
3. gogal-http
4. gogal-core
5. gogal-doctype
6. gogal-module
7. gogal-migration
8. gogal-generator
9. gogal-crud
10. gogal-admin
11. gogal-auth
12. gogal-users
13. gogal-permissions
14. gogal-studio
15. gogal-api
16. gogal-reports
17. gogal-jobs
18. gogal-events
19. gogal-files
20. gogal-email
21. gogal-notifications
22. gogal-realtime
23. gogal-importer
24. gogal-ai
25. gogal-bots
```

Immediate next implementation:

```text
Core metadata bootstrap migrations
```

---

## 12. Current Next Chapter

Next chapter:

```text
Core metadata bootstrap migrations for Gogal Studio
```

Goal:

```text
Create the fixed engine tables needed before JSON sync, builder UI, migration planner, CRUD API, and dynamic admin UI.
```

Required tables:

```text
core_module_def
core_doctype
core_docfield
core_doctype_action
core_doctype_link
core_docperm
core_doctype_state
core_migration_log
core_site_installation
core_installed_module
```

After this:

```text
Create Core module JSON files
Create DocType JSON schema structs
Create JSON sync engine
Expose Core metadata APIs
Start Studio UI shell
```
