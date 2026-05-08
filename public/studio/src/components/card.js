import { escapeHtml } from "../utils/escape.js";

export function summaryCard(label, value) {
  return `
    <div class="gs-summary-card">
      <span>${escapeHtml(label)}</span>
      <strong>${escapeHtml(value)}</strong>
    </div>
  `;
}
