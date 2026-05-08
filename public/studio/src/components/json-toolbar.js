import { escapeHtml } from "../utils/escape.js";

export function renderJsonToolbar(jsonText, doctypeName) {
  return `
    <div class="gs-json-toolbar">
      <div class="gs-json-size">${escapeHtml(jsonText.length)} chars</div>
      <div class="gs-json-toolbar-actions" aria-label="${escapeHtml(doctypeName)} JSON actions">
        <button id="copyJsonBtn" class="gs-json-btn" type="button">Copy JSON</button>
        <button id="downloadJsonBtn" class="gs-json-btn" type="button">Download JSON</button>
      </div>
    </div>
  `;
}
