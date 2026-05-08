import { escapeHtml } from "../utils/escape.js";

export function emptyState(message) {
  return `<div class="gs-empty">${escapeHtml(message)}</div>`;
}
