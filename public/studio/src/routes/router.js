import { loadDashboard } from "../pages/dashboard.page.js";
import { showNewDocTypeForm, showEditDocTypeForm } from "../pages/doctypeForm.page.js";

export function startRouter() {
  window.addEventListener("hashchange", resolveRoute);
  resolveRoute();
}

export function goTo(path) {
  window.location.hash = path;
}

async function resolveRoute() {
  const hash = window.location.hash || "#/dashboard";

  if (hash === "#/dashboard" || hash === "#/") {
    await loadDashboard();
    return;
  }

  if (hash === "#/doctype/new") {
    showNewDocTypeForm();
    return;
  }

  if (hash.startsWith("#/doctype/")) {
    const name = decodeURIComponent(hash.replace("#/doctype/", ""));
    await showEditDocTypeForm(name);
    return;
  }

  await loadDashboard();
}