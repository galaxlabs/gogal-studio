import { escapeHtml } from "../utils/escape.js";

const columns = [
  { key: "idx", label: "Idx" },
  { key: "fieldname", label: "Fieldname" },
  { key: "label", label: "Label" },
  { key: "fieldtype", label: "Fieldtype" },
  { key: "options", label: "Options" },
  { key: "reqd", label: "Reqd" },
  { key: "in_list_view", label: "In List View" }
];

export function renderFields(fields, activeFieldName = "") {
  const rows = fields || [];

  if (!rows.length) {
    return `<div class="gs-empty">No fields found.</div>`;
  }

  return `
    <table class="gs-table">
      <thead>
        <tr>
          ${columns.map((column) => `<th>${escapeHtml(column.label)}</th>`).join("")}
        </tr>
      </thead>
      <tbody>
        ${rows.map((field) => {
          const fieldname = field.fieldname || "";
          const activeClass = activeFieldName && fieldname === activeFieldName ? " active" : "";

          return `
            <tr class="gs-field-row${activeClass}" data-fieldname="${escapeHtml(fieldname)}">
              ${columns.map((column) => `<td>${escapeHtml(field[column.key])}</td>`).join("")}
            </tr>
          `;
        }).join("")}
      </tbody>
    </table>
  `;
}
