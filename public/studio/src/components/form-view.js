import { escapeHtml } from "../utils/escape.js";

export function renderFormShell({ doctype, recordName, bodyHtml }) {
  return `
    <div class="gs-form-shell">
      <div class="gs-form-header">
        <button id="backToListBtn" class="gs-back-btn" type="button">Back to List</button>
        <div>
          <h3 class="gs-form-title">${escapeHtml(recordName || "Untitled")}</h3>
          <p class="gs-form-subtitle">${escapeHtml(doctype || "DocType")} / editable record</p>
        </div>
        <button id="saveRecordBtn" class="gs-save-btn" type="button">Save</button>
      </div>
      <div class="gs-form-body">
        ${bodyHtml}
      </div>
    </div>
  `;
}
