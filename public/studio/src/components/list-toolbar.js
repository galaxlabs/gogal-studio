import { escapeHtml } from "../utils/escape.js";
import { renderSearchBox } from "./search.js";

export function renderListToolbar({ doctype, count, search }) {
  return `
    <div class="gs-list-toolbar">
      <div class="gs-list-title">
        <strong>${escapeHtml(doctype)}</strong>
        <span>${escapeHtml(count)} records</span>
      </div>
      ${renderSearchBox(search)}
      <div class="gs-list-actions">
        <button id="listRefreshBtn" class="gs-secondary-btn" type="button">Refresh</button>
        <button class="gs-disabled-btn" type="button" disabled>Actions</button>
      </div>
    </div>
  `;
}
