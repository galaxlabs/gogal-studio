import { escapeHtml } from "../utils/escape.js";

function viewLabelFor(activeView, activeRecord) {
  if (activeView === "resource") return "Resource List";
  if (activeView === "form") return activeRecord?.name || "Form View";
  if (activeView === "fields") return "Fields";
  if (activeView === "permissions") return "Permissions";
  if (activeView === "json") return "JSON Preview";

  return activeView || "";
}

export function renderBreadcrumb(bundle, activeView, activeRecord) {
  if (!bundle?.doctype) return "";

  const moduleName = bundle.doctype.module || "";
  const doctypeName = bundle.doctype.name || "";
  const viewLabel = viewLabelFor(activeView, activeRecord);

  return `
    <div class="gs-breadcrumb-wrap">
      <div class="gs-breadcrumb">
        <span>Studio</span>
        <span class="gs-breadcrumb-sep">/</span>
        <span>${escapeHtml(moduleName)}</span>
        <span class="gs-breadcrumb-sep">/</span>
        <span>${escapeHtml(doctypeName)}</span>
        <span class="gs-breadcrumb-sep">/</span>
        <strong>${escapeHtml(viewLabel)}</strong>
      </div>
      <span class="gs-view-badge">${escapeHtml(viewLabel)}</span>
    </div>
  `;
}
