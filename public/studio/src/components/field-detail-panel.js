import { escapeHtml } from "../utils/escape.js";

const detailFields = [
  "fieldname",
  "label",
  "fieldtype",
  "options",
  "reqd",
  "hidden",
  "read_only",
  "in_list_view",
  "idx"
];

export function renderFieldDetailPanel(field) {
  if (!field) {
    return `<div class="gs-field-detail-empty">Click a field row to view details.</div>`;
  }

  return `
    <aside class="gs-field-detail" aria-label="Field details">
      <div class="gs-field-detail-title">${escapeHtml(field.label || field.fieldname || "Field")}</div>
      <div class="gs-field-detail-grid">
        ${detailFields.map((key) => `
          <div class="gs-field-detail-row">
            <div class="gs-field-detail-label">${escapeHtml(key)}</div>
            <div class="gs-field-detail-value">${escapeHtml(field[key] ?? "")}</div>
          </div>
        `).join("")}
      </div>
    </aside>
  `;
}
