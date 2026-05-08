import { escapeHtml } from "../utils/escape.js";

const booleanKeys = [
  "read",
  "write",
  "create",
  "delete",
  "submit",
  "cancel",
  "amend",
  "print",
  "email",
  "export",
  "import",
  "share",
  "report"
];

function normalizedPermission(permission = {}) {
  return {
    role: permission.role || "",
    permlevel: permission.permlevel ?? 0,
    read: permission.read ?? permission.read_perm ?? false,
    write: permission.write ?? permission.write_perm ?? false,
    create: permission.create ?? permission.create_perm ?? false,
    delete: permission.delete ?? permission.delete_perm ?? false,
    submit: permission.submit ?? permission.submit_perm ?? false,
    cancel: permission.cancel ?? permission.cancel_perm ?? false,
    amend: permission.amend ?? permission.amend_perm ?? false,
    print: permission.print ?? permission.print_perm ?? false,
    email: permission.email ?? permission.email_perm ?? false,
    export: permission.export ?? permission.export_perm ?? false,
    import: permission.import ?? permission.import_perm ?? false,
    share: permission.share ?? permission.share_perm ?? false,
    report: permission.report ?? permission.report_perm ?? false,
    idx: permission.idx ?? ""
  };
}

function renderValue(key, value) {
  if (booleanKeys.includes(key)) {
    const isOn = Boolean(value);
    return `<span class="gs-permission-badge ${isOn ? "true" : "false"}">${escapeHtml(isOn)}</span>`;
  }

  return escapeHtml(value ?? "");
}

export function renderPermissionDetailPanel(permission) {
  if (!permission) {
    return `<div class="gs-permission-detail-empty">Click a permission row to view details.</div>`;
  }

  const details = normalizedPermission(permission);
  const rows = [
    "role",
    "permlevel",
    ...booleanKeys,
    "idx"
  ];

  return `
    <aside class="gs-permission-detail" aria-label="Permission details">
      <div class="gs-permission-detail-title">${escapeHtml(details.role || "Permission")}</div>
      <div class="gs-permission-detail-grid">
        ${rows.map((key) => `
          <div class="gs-permission-detail-row">
            <div class="gs-permission-detail-label">${escapeHtml(key)}</div>
            <div class="gs-permission-detail-value">${renderValue(key, details[key])}</div>
          </div>
        `).join("")}
      </div>
    </aside>
  `;
}
