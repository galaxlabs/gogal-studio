# Gogal Core and Desk Scope

## Naming

- Framework/app name: `gogal`
- CLI: `gogal`
- Browser product name: `Gogal Studio`
- Main UI module: `Desk`
- Default site: `gogal.dev`

Current temporary installed app name may still be `gogal_studio`.
It will be renamed to `gogal` later through a safe migration.

---

## Current focus

Before building full dynamic CRUD/List/Form generation, complete these two modules:

1. Core
2. Desk

No business modules yet.

Do not create Customer, Supplier, Item, Invoice, ERP, Transport, or app-specific DocTypes in this phase.

---

## Core module purpose

Core owns the framework engine metadata.

Core must support:

- app registry
- module registry
- DocType registry
- DocField registry
- DocPerm registry
- user registry
- role registry
- user-role assignment
- installed app/module tracking

## Core DocTypes for first completion

Required now:

- Installed App
- Installed Module
- Module Def
- DocType
- DocField
- DocPerm
- User
- Role
- Has Role

Later Core additions:

- File
- Comment
- Activity Log
- Error Log
- Patch Log
- Language
- Report
- Page
- System Settings

---

## Desk module purpose

Desk owns the Gogal Studio UI metadata and navigation layer.

Desk is not a separate app.
Desk is a module inside Gogal.

Desk supports:

- workspace UI
- navigation/sidebar
- dashboard blocks
- list view settings
- route history
- notifications
- todo/task widgets

## Desk DocTypes for first completion

Required first:

- Workspace
- Workspace Link
- Workspace Shortcut
- Workspace Sidebar
- Workspace Sidebar Item
- List View Settings
- Route History
- ToDo
- Notification Log
- Dashboard
- Number Card

Later Desk additions:

- Kanban Board
- Calendar View
- Dashboard Chart
- Custom HTML Block
- Onboarding
- System Console
- System Health Report

---

## Design rule

Do not copy Frappe directly.

Use Frappe Core/Desk only as reference for ideas:

- metadata
- workspace
- route history
- list settings
- dashboard
- notification
- permissions

Use Sponge only as reference for:

- definition-driven generation
- generator templates
- clean output separation

Use Encore template-engine only as reference for:

- backend route serving UI
- static asset separation
- route-to-template pattern

Gogal must remain runtime metadata-driven first.

---

## Immediate implementation order

1. Finish Core APIs.
2. Seed minimum Desk metadata tables/DocTypes.
3. Add Desk API:
   - GET /api/desk/workspaces
   - GET /api/desk/sidebar
   - GET /api/desk/dashboard
4. Improve `/studio` UI to read Desk metadata.
5. Add modular frontend build structure.
6. Then start dynamic DocType UI generation.