import { renderTable } from "../components/table.js";

export function renderResourceList(bundle) {
  const rows = bundle.records || [];
  const fields = bundle.fields || [];
  const visibleFields = fields.filter((field) => field.in_list_view && !field.hidden && field.fieldname);

  if (!rows.length) {
    return renderTable([], [], "No records returned from resource API.");
  }

  const keys = visibleFields.length
    ? visibleFields.map((field) => field.fieldname)
    : Object.keys(rows[0]).slice(0, 8);

  const columns = keys.map((key) => ({ key, label: key }));

  return renderTable(columns, rows);
}
