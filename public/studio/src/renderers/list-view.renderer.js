import { renderListToolbar } from "../components/list-toolbar.js";
import { renderListView } from "../components/list-view.js";

function getColumns(bundle) {
  const rows = bundle.records || [];
  const fields = bundle.fields || [];
  const visibleFields = fields.filter((field) => field.in_list_view && !field.hidden && field.fieldname);

  if (visibleFields.length) {
    return visibleFields.map((field) => ({
      key: field.fieldname,
      label: field.label || field.fieldname
    }));
  }

  const firstRow = rows[0] || {};
  return Object.keys(firstRow).slice(0, 8).map((key) => ({ key, label: key }));
}

function filterRows(rows, columns, searchText) {
  const query = String(searchText || "").trim().toLowerCase();
  if (!query) return rows;

  return rows.filter((row) => columns.some((column) => {
    const value = row[column.key];
    if (value === null || value === undefined) return false;
    return String(value).toLowerCase().includes(query);
  }));
}

export function renderDynamicListView(bundle, searchText = "") {
  const rows = bundle.records || [];
  const doctype = bundle.doctype?.name || "Resource";
  const columns = getColumns(bundle);
  const filteredRows = filterRows(rows, columns, searchText);

  return `
    ${renderListToolbar({ doctype, count: filteredRows.length, search: searchText })}
    <div class="gs-list-status">
      <span>Read-only list view</span>
      <span>Bulk actions placeholder</span>
    </div>
    ${renderListView({ columns, rows: filteredRows })}
  `;
}
