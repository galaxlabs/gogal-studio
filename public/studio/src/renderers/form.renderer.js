import { getFieldTypeConfig } from "../core/fieldTypes.js";

export function renderDynamicFormPreview(doctype, fields) {
  const visibleFields = (fields || []).filter((field) => !field.hidden);

  if (!visibleFields.length) {
    return `<div class="muted">No visible fields found for ${doctype.name}.</div>`;
  }

  return `
    <div class="dynamic-form">
      ${visibleFields.map((field) => renderFieldPreview(field)).join("")}
    </div>
  `;
}

function renderFieldPreview(field) {
  const config = getFieldTypeConfig(field.fieldtype);
  const label = field.label || field.fieldname;
  const required = field.reqd ? `<span class="required">*</span>` : "";
  const readonly = field.read_only ? "readonly" : "";
  const columns = config.defaultColumns || 6;

  if (config.control === "layout") {
    if (config.layoutType === "section") {
      return `<div class="form-section-title">${label}</div>`;
    }

    return `<div class="form-column-break"></div>`;
  }

  return `
    <div class="form-field col-${columns}">
      <label>${escapeHtml(label)} ${required}</label>
      ${renderControl(field, config, readonly)}
      ${renderFieldHelp(field, config)}
    </div>
  `;
}

function renderControl(field, config, readonly) {
  const value = field.default_value || "";

  switch (config.control) {
    case "textarea":
      return `<textarea rows="${config.rows || 3}" ${readonly} placeholder="${escapeHtml(config.placeholder || "")}">${escapeHtml(value)}</textarea>`;

    case "editor":
      return `
        <div class="editor-preview" contenteditable="${readonly ? "false" : "true"}">
          ${escapeHtml(value || config.placeholder || "Rich text")}
        </div>
      `;

    case "checkbox":
      return `<input type="checkbox" ${value === "1" || value === true ? "checked" : ""} ${readonly} />`;

    case "select":
      return renderSelect(field, readonly);

    case "link":
      return `<input type="text" ${readonly} placeholder="Link: ${escapeHtml(field.options || "Select DocType")}" />`;

    case "table":
      return `
        <div class="table-preview">
          Child Table: <strong>${escapeHtml(field.options || "Select child DocType")}</strong>
        </div>
      `;

    case "file":
      return `<input type="file" ${readonly} />`;

    case "image":
      return `<input type="file" accept="image/*" ${readonly} />`;

    case "code":
      return `<textarea rows="6" ${readonly} class="code-input" placeholder="${escapeHtml(config.language || "code")}">${escapeHtml(value)}</textarea>`;

    default:
      return `<input type="${config.inputType || "text"}" step="${config.step || ""}" ${readonly} placeholder="${escapeHtml(config.placeholder || "")}" value="${escapeHtml(value)}" />`;
  }
}

function renderSelect(field, readonly) {
  const options = String(field.options || "")
    .split("\n")
    .map((item) => item.trim())
    .filter(Boolean);

  return `
    <select ${readonly ? "disabled" : ""}>
      <option value="">Select</option>
      ${options.map((option) => `<option value="${escapeHtml(option)}">${escapeHtml(option)}</option>`).join("")}
    </select>
  `;
}

function renderFieldHelp(field, config) {
  if (config.optionsMode === "none") {
    return "";
  }

  if (config.optionsMode === "lines") {
    return `<small class="field-help">Options: one value per line</small>`;
  }

  if (config.optionsMode === "doctype") {
    return `<small class="field-help">Options should be target DocType name</small>`;
  }

  if (config.optionsMode === "child_doctype") {
    return `<small class="field-help">Options should be child table DocType name</small>`;
  }

  return "";
}

function escapeHtml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}