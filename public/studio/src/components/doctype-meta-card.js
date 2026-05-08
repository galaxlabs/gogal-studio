import { escapeHtml } from "../utils/escape.js";

function booleanBadge(label, value) {
  const isOn = Boolean(value);
  const stateClass = isOn ? "success" : "muted";

  return `
    <span class="gs-meta-badge ${stateClass}">
      ${escapeHtml(label)}: ${escapeHtml(isOn)}
    </span>
  `;
}

function valueBadge(label, value) {
  return `
    <span class="gs-meta-badge">
      ${escapeHtml(label)}: ${escapeHtml(value ?? "")}
    </span>
  `;
}

export function renderDocTypeMetaCard(bundle) {
  const doctype = bundle?.doctype || {};
  const records = bundle?.records || [];
  const name = doctype.name || "DocType";
  const moduleName = doctype.module || "";
  const tableName = doctype.table_name || "";

  return `
    <section class="gs-meta-card" aria-label="DocType metadata">
      <div>
        <div class="gs-meta-title">${escapeHtml(name)}</div>
        <div class="gs-meta-subtitle">${escapeHtml(moduleName)} &middot; ${escapeHtml(tableName)}</div>
      </div>
      <div class="gs-meta-badges">
        ${valueBadge("App", doctype.app_name)}
        ${booleanBadge("Single", doctype.is_single)}
        ${booleanBadge("Child", doctype.is_child_table)}
        ${booleanBadge("Submittable", doctype.is_submittable)}
        ${booleanBadge("Tree", doctype.is_tree)}
        ${valueBadge("Records", records.length)}
      </div>
    </section>
  `;
}
