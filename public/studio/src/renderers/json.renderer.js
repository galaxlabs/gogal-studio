import { renderJsonToolbar } from "../components/json-toolbar.js";
import { escapeHtml } from "../utils/escape.js";

export function buildDocTypeJson(bundle) {
  return {
    ...(bundle?.doctype || {}),
    fields: bundle?.fields || [],
    permissions: bundle?.permissions || []
  };
}

export function renderJsonPreview(bundle) {
  const preview = buildDocTypeJson(bundle);
  const doctypeName = preview.name || "doctype";
  const jsonText = JSON.stringify(preview, null, 2);

  return `
    ${renderJsonToolbar(jsonText, doctypeName)}
    <pre class="gs-json-preview" id="jsonPreviewText">${escapeHtml(jsonText)}</pre>
  `;
}
