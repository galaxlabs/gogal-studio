import { escapeHtml } from "../utils/escape.js";
import { emptyState } from "./empty.js";

function formatCell(value) {
  if (value === null || value === undefined) return "";
  if (typeof value === "object") return JSON.stringify(value);
  return value;
}

export function renderTable(columns, rows, emptyMessage = "No records found.") {
  if (!rows || rows.length === 0) {
    return emptyState(emptyMessage);
  }

  return `
    <div class="gs-table-wrap">
      <table class="gs-table">
        <thead>
          <tr>
            ${columns.map((col) => `<th>${escapeHtml(col.label)}</th>`).join("")}
          </tr>
        </thead>
        <tbody>
          ${rows.map((row) => `
            <tr>
              ${columns.map((col) => `<td>${escapeHtml(formatCell(row[col.key]))}</td>`).join("")}
            </tr>
          `).join("")}
        </tbody>
      </table>
    </div>
  `;
}
