import {
  getInstalledApps,
  getModules,
  getDocTypes,
  getDocType,
  getDocTypeFields,
  getDocTypePermissions
} from "../api/core.js";

import { renderDocTypeDetails } from "../renderers/doctypeRenderer.js";
import { renderFields } from "../renderers/fieldRenderer.js";
import { renderPermissions } from "../renderers/permissionRenderer.js";
import { goTo } from "../routes/router.js";
let dashboardBound = false;

export async function loadDashboard() {
  const statusPill = document.getElementById("statusPill");
  const doctypeList = document.getElementById("doctypeList");

  setStatus("Loading...");

  try {
    const [appsRes, modulesRes, doctypesRes] = await Promise.all([
      getInstalledApps(),
      getModules(),
      getDocTypes()
    ]);

    const apps = appsRes.data || [];
    const modules = modulesRes.data || [];
    const doctypes = doctypesRes.data || [];

    document.getElementById("appCount").textContent = apps.length;
    document.getElementById("moduleCount").textContent = modules.length;
    document.getElementById("doctypeCount").textContent = doctypes.length;

    doctypeList.innerHTML = doctypes.map((dt) => `
      <button class="doctype-item" data-doctype="${dt.name}">
        <strong>${dt.name}</strong>
        <small>${dt.module} · ${dt.table_name}</small>
      </button>
    `).join("");

    doctypeList.querySelectorAll(".doctype-item").forEach((btn) => {
      btn.addEventListener("click", () => {
  goTo(`/doctype/${encodeURIComponent(btn.dataset.doctype)}`);
});
    });

    if (!dashboardBound) {
      document.getElementById("refreshBtn").addEventListener("click", loadDashboard);

      const newDocTypeBtn = document.getElementById("newDocTypeBtn");

      if (newDocTypeBtn && !newDocTypeBtn.dataset.bound) {
        newDocTypeBtn.addEventListener("click", () => goTo("/doctype/new"));
        newDocTypeBtn.dataset.bound = "1";
      }

      dashboardBound = true;
    }

    statusPill.className = "status-pill success";
    statusPill.textContent = "Loaded";
  } catch (err) {
    console.error(err);
    statusPill.className = "status-pill error";
    statusPill.textContent = "Load failed";
    doctypeList.innerHTML = `<div class="muted">${err.message}</div>`;
  }
}

async function loadDocType(name) {
  setStatus("Loading DocType...");

  try {
    const [dtRes, fieldsRes, permsRes] = await Promise.all([
      getDocType(name),
      getDocTypeFields(name),
      getDocTypePermissions(name)
    ]);

    const dt = dtRes.data;
    const fields = fieldsRes.data || [];
    const perms = permsRes.data || [];

    document.getElementById("pageTitle").textContent = dt.name;
    document.getElementById("pageSubtitle").textContent = `${dt.module} · ${dt.table_name}`;

    document.getElementById("fieldCount").textContent = fields.length;
    document.getElementById("permCount").textContent = perms.length;

    document.querySelectorAll(".doctype-item").forEach((item) => {
      item.classList.toggle("active", item.dataset.doctype === name);
    });

    document.getElementById("detailsPanel").innerHTML = renderDocTypeDetails(dt);
    document.getElementById("fieldsPanel").innerHTML = renderFields(fields);
    document.getElementById("permsPanel").innerHTML = renderPermissions(perms);

    setStatus("Loaded", "success");
  } catch (err) {
    console.error(err);
    setStatus("Load failed", "error");
  }
}

function setStatus(message, type = "") {
  const statusPill = document.getElementById("statusPill");

  statusPill.textContent = message;
  statusPill.className = "status-pill";

  if (type) {
    statusPill.classList.add(type);
  }
}