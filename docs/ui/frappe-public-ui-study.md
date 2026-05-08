# Frappe Public UI Study for Gogal Studio

## 1. What Was Inspected

Inspected `E:\dev\frappe-version-16\frappe\public` with emphasis on:

- `public/js`: desk bootstrap, router, views, list view, form layout, controls, sidebar, toolbar, workspace, reports, and bundled JS entry files.
- `public/css`: local Bootstrap, local fonts, icon font assets, and small standalone CSS files.
- `public/scss`: desk, common, website, print, report, login, and bundle entry SCSS files.
- `public/html`: print template assets.
- `public/js/frappe/form/templates`: form sidebar, dashboard, footer, links, sharing, contact, address, and timeline templates.
- `public/js/frappe/views`: view factories and page/list/form/report/tree/workspace-oriented view modules.
- `public/js/frappe/list`: list factory, base list, list view, filters, settings, bulk operations, and sidebar grouping.
- `public/js/frappe/form`: form, layout, grid, toolbar, sidebar, dashboard, save, formatters, controls, and timeline modules.
- `public/js/frappe/ui`: page shell, sidebar, toolbar/navbar/search, filters, dialogs, messages, tags, trees, and datatable adapters.

`public/build.json`, `public/app.html`, and `public/desk` were not present in this checkout.

## 2. Useful Ideas to Borrow

- Keep Desk as a shell that boots once, then swaps routed page content.
- Separate startup/boot logic from page rendering.
- Keep route definitions small and independent from page implementations.
- Treat List, Form, Report, and Workspace as distinct view/page types.
- Keep metadata-driven rendering behind renderer modules rather than embedding table/form logic in route code.
- Keep sidebar, topbar, tabs, tables, badges, empty states, cards, and panels as small reusable UI components.
- Keep source files modular and bundle to a single browser-loaded asset.
- Keep style bundles and vendor assets local.

## 3. What NOT to Copy

- Do not copy Frappe source code, templates, CSS, or bundled JS.
- Do not copy Frappe globals such as `frappe.*` as an application pattern.
- Do not copy the large jQuery application runtime.
- Do not copy Vue-based builders or Frappe-specific build conventions.
- Do not copy permission, form-save, workflow, report, or realtime behavior before Gogal APIs support them.
- Do not copy Frappe asset paths, app routes, or naming.

## 4. Recommended Gogal Studio UI Structure

Gogal Studio should stay Vanilla JS first:

```text
public/studio/
  index.html
  styles.css
  app.js
  dist/studio.js
  src/
    main.js
    boot/
    api/
    routes/
    pages/
    components/
    renderers/
    layout/
    state/
    utils/
```

The structure should keep API calls, state, route/page orchestration, components, and renderers separate while remaining small enough to understand without a framework.

## 5. Asset Loading Approach

- Serve everything from the Go server on one port.
- Keep `/studio` as the HTML entry.
- Serve Studio source output from `/studio-assets`.
- Serve vendor assets from `/vendor`.
- Load only local assets in final HTML:
  - `/vendor/tailwind/tailwind.js`
  - `/vendor/bootstrap/bootstrap.min.css`
  - `/vendor/jquery/jquery.min.js`
  - `/vendor/bootstrap/bootstrap.bundle.min.js`
  - `/studio-assets/dist/studio.js`
- Do not use CDN references.

## 6. Route/Page System Approach

Use a tiny router with named routes:

- `dashboard`: loads installed apps, modules, and DocTypes.
- `doctype`: loads DocType metadata, fields, permissions, and read-only resource rows.
- `resource-list`: reserved for a future standalone resource list page.

Routes should create page objects with `mount()` and optional `destroy()` methods. The current UI can remain single-screen while the source shape supports future page expansion.

## 7. Component System Approach

Use small functional components that return HTML strings or update known DOM nodes:

- `topbar`
- `sidebar`
- `tabs`
- `status`
- `card`
- `panel`
- `table`
- `badge`
- `empty`

Components should not know about API endpoints. They should receive data and callbacks.

## 8. Form/List Renderer Approach

Keep renderers metadata-driven:

- Resource lists choose `in_list_view` fields first, then fall back to keys from the first returned row.
- Field metadata renders as a read-only DocField table.
- Permission metadata normalizes `read/read_perm`, `write/write_perm`, `create/create_perm`, and `delete/delete_perm`.
- JSON preview renders the current API bundle.

Form rendering should come later and should be built from DocField metadata, not hand-coded per DocType.

## 9. Future Build Pipeline Notes

- Keep esbuild for now.
- Keep `npm run studio:build`, `studio:watch`, and `studio:prod`.
- Avoid React, Vite, Vue, Angular, Svelte, and Next until the product need is proven.
- Later, split CSS into source files only when the design system grows enough to justify it.
- Later, add a development watch task that runs esbuild and the Go server together, but keep production output as local static files served by Go.
