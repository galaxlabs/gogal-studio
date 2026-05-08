import { escapeHtml } from "../utils/escape.js";

export function renderSearchBox(value = "") {
  return `
    <div class="gs-search-wrap">
      <input
        id="listSearchInput"
        class="gs-search-input"
        type="search"
        placeholder="Search loaded records..."
        value="${escapeHtml(value)}"
      />
    </div>
  `;
}
