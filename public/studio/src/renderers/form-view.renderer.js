import { boolBadge } from "../components/badge.js";
import { renderFormShell } from "../components/form-view.js";
import { escapeHtml } from "../utils/escape.js";

const systemFields = [
  "name",
  "owner",
  "created_by",
  "modified_by",
  "created_at",
  "updated_at",
  "creation",
  "modified",
  "docstatus",
  "idx"
];

const layoutFieldTypes = new Set(["Section", "Column", "Section Break", "Column Break", "Table"]);

function valueFor(record, fieldname) {
  const value = record?.[fieldname];
  return value === null || value === undefined ? "" : value;
}

function readonlyValue(field, value) {
  const fieldtype = field.fieldtype || "";

  if (fieldtype === "Check" || typeof value === "boolean") {
    return boolBadge(value);
  }

  if (fieldtype === "Link") {
    return `<span class="gs-value-badge">${escapeHtml(value)}</span>`;
  }

  if (fieldtype === "Table") {
    return `<span class="gs-form-value gs-muted-value">Child table editing not implemented yet</span>`;
  }

  if (fieldtype === "JSON" || fieldtype === "Code") {
    const rendered = typeof value === "object" ? JSON.stringify(value, null, 2) : value;
    return `<pre class="gs-form-pre">${escapeHtml(rendered)}</pre>`;
  }

  return `<span class="gs-form-value">${escapeHtml(value)}</span>`;
}

function editableControl(field, value) {
  const fieldname = field.fieldname || "";
  const fieldtype = field.fieldtype || "Data";
  const required = field.reqd ? " required" : "";
  const baseAttrs = `data-fieldname="${escapeHtml(fieldname)}" data-fieldtype="${escapeHtml(fieldtype)}"${required}`;

  if (fieldtype === "Check") {
    return `
      <label class="gs-check-control">
        <input class="gs-form-input" type="checkbox" ${baseAttrs} ${value ? "checked" : ""} />
        <span>${value ? "True" : "False"}</span>
      </label>
    `;
  }

  if (["Text", "Small Text", "Long Text", "Text Editor", "JSON", "Code"].includes(fieldtype)) {
    const textValue = typeof value === "object" ? JSON.stringify(value, null, 2) : value;
    return `<textarea class="gs-form-input gs-form-textarea" ${baseAttrs}>${escapeHtml(textValue)}</textarea>`;
  }

  if (fieldtype === "Select") {
    const options = String(field.options || "")
      .split(/\r?\n/)
      .map((item) => item.trim())
      .filter(Boolean);

    return `
      <select class="gs-form-input" ${baseAttrs}>
        <option value=""></option>
        ${options.map((option) => `
          <option value="${escapeHtml(option)}" ${String(value) === option ? "selected" : ""}>${escapeHtml(option)}</option>
        `).join("")}
      </select>
    `;
  }

  const typeMap = {
    Int: "number",
    Float: "number",
    Currency: "number",
    Date: "date",
    Datetime: "datetime-local",
    Time: "time"
  };
  const inputType = typeMap[fieldtype] || "text";
  const step = ["Float", "Currency"].includes(fieldtype) ? ` step="any"` : "";

  return `<input class="gs-form-input" type="${escapeHtml(inputType)}"${step} value="${escapeHtml(value)}" ${baseAttrs} />`;
}

function renderFieldRows(fields, record) {
  const visibleFields = (fields || []).filter((field) => !field.hidden && field.fieldname);

  if (!visibleFields.length) {
    return `<div class="gs-empty">No visible fields found for this DocType.</div>`;
  }

  return `
    <div class="gs-form-grid">
      ${visibleFields.map((field) => {
        const label = field.label || field.fieldname;
        const value = valueFor(record, field.fieldname);
        const isEditable = !field.read_only && !layoutFieldTypes.has(field.fieldtype || "");

        return `
          <div class="gs-form-row ${isEditable ? "editable" : "readonly"}">
            <div class="gs-form-label">${escapeHtml(label)}${field.reqd ? `<span class="gs-required">*</span>` : ""}</div>
            <div>${isEditable ? editableControl(field, value) : readonlyValue(field, value)}</div>
          </div>
        `;
      }).join("")}
    </div>
  `;
}

function renderSystemRows(record) {
  const rows = systemFields
    .filter((fieldname) => Object.prototype.hasOwnProperty.call(record || {}, fieldname))
    .map((fieldname) => ({
      label: fieldname,
      value: valueFor(record, fieldname)
    }));

  if (!rows.length) {
    return "";
  }

  return `
    <details class="gs-form-section gs-system-section">
      <summary class="gs-form-section-title">System Fields</summary>
      <div class="gs-form-grid">
        ${rows.map((row) => `
          <div class="gs-form-row readonly">
            <div class="gs-form-label">${escapeHtml(row.label)}</div>
            <div class="gs-form-value">${escapeHtml(row.value)}</div>
          </div>
        `).join("")}
      </div>
    </details>
  `;
}

export function renderEditableForm(bundle, record) {
  const doctype = bundle.doctype?.name || "DocType";
  const recordName = record?.name || "Editable record";
  const fields = bundle.fields || [];

  const bodyHtml = `
    <section class="gs-form-section">
      <div class="gs-form-section-title">Fields</div>
      ${renderFieldRows(fields, record)}
    </section>
    ${renderSystemRows(record)}
  `;

  return renderFormShell({ doctype, recordName, bodyHtml });
}

export const renderReadOnlyForm = renderEditableForm;
