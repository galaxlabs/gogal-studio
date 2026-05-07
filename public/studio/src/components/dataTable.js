export function renderTable(columns, rows) {
  if (!rows || rows.length === 0) {
    return `<div class="muted">No records found.</div>`;
  }

  return `
    <table class="table">
      <thead>
        <tr>
          ${columns.map((col) => `<th>${col.label}</th>`).join("")}
        </tr>
      </thead>
      <tbody>
        ${rows.map((row) => `
          <tr>
            ${columns.map((col) => {
              const value = col.render ? col.render(row) : row[col.key];
              return `<td>${value ?? ""}</td>`;
            }).join("")}
          </tr>
        `).join("")}
      </tbody>
    </table>
  `;
}