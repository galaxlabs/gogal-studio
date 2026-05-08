import { escapeHtml } from "../utils/escape.js";

const columns = [
  { key: "role", label: "Role" },
  { key: "permlevel", label: "Level" },
  { key: "read", label: "Read" },
  { key: "write", label: "Write" },
  { key: "create", label: "Create" },
  { key: "delete", label: "Delete" }
];

function permissionKey(permission, index) {
  const role = permission.role || "";
  const permlevel = permission.permlevel ?? 0;
  const idx = permission.idx ?? index;
  return `${role}-${permlevel}-${idx}`;
}

function normalizePermission(permission, index) {
  return {
    key: permissionKey(permission, index),
    role: permission.role || "",
    permlevel: permission.permlevel ?? 0,
    read: permission.read ?? permission.read_perm ?? false,
    write: permission.write ?? permission.write_perm ?? false,
    create: permission.create ?? permission.create_perm ?? false,
    delete: permission.delete ?? permission.delete_perm ?? false
  };
}

export function getPermissionKey(permission, index = 0) {
  return permissionKey(permission || {}, index);
}

export function renderPermissions(perms, activePermissionKey = "") {
  const rows = (perms || []).map((permission, index) => normalizePermission(permission, index));

  if (!rows.length) {
    return `<div class="gs-empty">No permissions found.</div>`;
  }

  return `
    <table class="gs-table">
      <thead>
        <tr>
          ${columns.map((column) => `<th>${escapeHtml(column.label)}</th>`).join("")}
        </tr>
      </thead>
      <tbody>
        ${rows.map((row) => {
          const activeClass = activePermissionKey && row.key === activePermissionKey ? " active" : "";

          return `
            <tr class="gs-permission-row${activeClass}" data-permission-key="${escapeHtml(row.key)}">
              ${columns.map((column) => `<td>${escapeHtml(row[column.key])}</td>`).join("")}
            </tr>
          `;
        }).join("")}
      </tbody>
    </table>
  `;
}
