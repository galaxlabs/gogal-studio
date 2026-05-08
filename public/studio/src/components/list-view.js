import { escapeHtml } from "../utils/escape.js";
import { emptyState } from "./empty.js";

function formatCell(value) {
  if (value === null || value === undefined) return "";
  if (typeof value === "object") return JSON.stringify(value);
  return value;
}

export function renderListView({ columns, rows }) {
  if (!rows || rows.length === 0) {
    return emptyState("No records found.");
  }

  return `
    <div class="gs-table-wrap">
      <table class="gs-list-table">
        <thead>
          <tr>
            <th class="gs-row-checkbox"><input type="checkbox" disabled aria-label="Select all placeholder" /></th>
            ${columns.map((column) => `<th>${escapeHtml(column.label)}</th>`).join("")}
          </tr>
        </thead>
        <tbody>
          ${rows.map((row, index) => {
            const name = row.name || "";
            return `
              <tr class="gs-row-clickable" data-row-name="${escapeHtml(name)}" data-row-index="${escapeHtml(index)}">
                <td class="gs-row-checkbox"><input class="gs-row-select" type="checkbox" disabled aria-label="Select row placeholder" /></td>
                ${columns.map((column) => `<td>${escapeHtml(formatCell(row[column.key]))}</td>`).join("")}
              </tr>
            `;
          }).join("")}
        </tbody>
      </table>
    </div>
  `;
}
