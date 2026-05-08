import { store } from "../state/store.js";
import { escapeHtml } from "../utils/escape.js";

function moduleLabel(module) {
  return module?.module_name || module?.name || "";
}

export function renderModules(modules, onModuleClick) {
  const moduleNav = document.getElementById("moduleNav");

  if (!moduleNav) return;

  const items = [{ name: "All", module_name: "All Modules", app_name: "All apps" }, ...(modules || [])];

  moduleNav.innerHTML = items.map((m) => {
    const name = m.name || m.module_name;
    const active = store.activeModule === name ? " active" : "";
    return `
      <button class="gs-nav-item${active}" data-module="${escapeHtml(name)}" type="button">
        <strong>${escapeHtml(moduleLabel(m))}</strong>
        <small>${escapeHtml(m.app_name)}</small>
      </button>
    `;
  }).join("");

  moduleNav.querySelectorAll(".gs-nav-item").forEach((btn) => {
    btn.addEventListener("click", () => onModuleClick(btn.dataset.module));
  });
}

export function renderDocTypes(doctypes, onDocTypeClick) {
  const doctypeList = document.getElementById("doctypeList");

  if (!doctypeList) return;

  if (!doctypes || doctypes.length === 0) {
    doctypeList.innerHTML = `<div class="gs-empty-sidebar">No DocTypes found.</div>`;
    return;
  }

  doctypeList.innerHTML = doctypes.map((dt) => `
    <button class="gs-doctype-item ${dt.name === store.activeDocType ? "active" : ""}" data-doctype="${escapeHtml(dt.name)}" type="button">
      <strong>${escapeHtml(dt.name)}</strong>
      <small>${escapeHtml(dt.module)} / ${escapeHtml(dt.table_name)}</small>
    </button>
  `).join("");

  doctypeList.querySelectorAll(".gs-doctype-item").forEach((btn) => {
    btn.addEventListener("click", () => onDocTypeClick(btn.dataset.doctype));
  });
}

export function bindDocTypeSearch(onSearch) {
  const input = document.getElementById("doctypeSearchInput");

  if (!input) return;

  input.value = store.doctypeSearch || "";
  input.oninput = (event) => {
    onSearch(event.target.value);
  };
}
