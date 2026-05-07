# Gogal Studio

Gogal Studio is a Go-powered visual app builder by Galaxy Labs.

Goal:
- Build business apps visually
- Generate DocType-style metadata
- Preview and apply migrations
- Auto-create CRUD APIs
- Auto-create admin list/form UI
- Support feature packs
- Support AI agents and bots later

Current stage:
Foundation setup.

## Current Focus: DocType Builder UI

Builder page path:
- `/studio-assets/builder.html`

What works now:
- `/studio` remains the read-only metadata viewer.
- `/studio-assets/builder.html` is the current DocType Builder UI.
- The builder supports creating a DocType draft, selecting an existing Module Def or entering a new module name, mutually exclusive DocType options, field add/edit/reorder/remove, action add/edit/reorder/remove, and formatted JSON preview with copy.
- Builder save writes `modules/{module}/doctype/{doctype}/{doctype}.json` and then syncs metadata from that file.
- Module Def is stored as a DocType property and is not shown as a field card.
- Normal canvas mode hides meta/internal fields, with a checkbox to inspect them.

Draft-only:
- Builder edits stay in UI draft state until metadata save.
- No physical database table is created from this UI.
- No migration is applied from this UI.

How to run:
- `go run .\cmd\server`
- Open `http://127.0.0.1:8080/studio`
- Open `http://127.0.0.1:8080/studio-assets/builder.html`

How to verify:
- Confirm `/studio` still loads as the metadata viewer.
- In the builder, create a DocType draft, select or create a Module Def, test Single/Submittable/Child Table/Editable Grid option rules, add and reorder fields, edit field properties, add and reorder actions, and confirm JSON Preview updates.
- Save the draft and confirm a JSON file is created under `modules/{module}/doctype/{doctype}/{doctype}.json`.
- Confirm no database migration or physical table creation happens.

Next steps:
- Add a safer persisted draft/history layer.
- Add migration preview before any schema apply.
- Add module/history validation and conflict handling around file writes.
