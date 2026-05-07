import {
  getDocType,
  getDocTypeFields,
  getDocTypePermissions
} from "../api/core.js";

import {
  renderDocTypeDetails,
  renderDocTypeFields,
  renderDocTypePermissions
} from "../renderers/doctype.renderer.js";

import { renderDynamicFormPreview } from "../renderers/form.renderer.js";

export function DocTypePage(params) {
  async function mount() {
    const name = params.name;

    if (!name) {
      setStatus("No DocType selected", "error");
      return;
    }

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
      document.getElementById("pageSubtitle").textContent = `${dt.module} · ${dt.app_name}`;

      document.getElementById("fieldCount").textContent = fields.length;
      document.getElementById("permCount").textContent = perms.length;

      document.getElementById("detailsPanel").innerHTML = renderDocTypeDetails(dt);
      document.getElementById("fieldsPanel").innerHTML = renderDocTypeFields(fields);
      document.getElementById("permsPanel").innerHTML = renderDocTypePermissions(perms);

      const previewPanel = document.getElementById("formPreviewPanel");
      if (previewPanel) {
        previewPanel.innerHTML = renderDynamicFormPreview(dt, fields);
      }

      markActive(dt.name);
      setStatus("Loaded", "success");
    } catch (err) {
      console.error(err);
      setStatus("Load failed", "error");

      const detailsPanel = document.getElementById("detailsPanel");
      if (detailsPanel) {
        detailsPanel.innerHTML = `<div class="muted">${escapeHtml(err.message)}</div>`;
      }
    }
  }

  return {
    mount
  };
}

function markActive(name) {
  document.querySelectorAll(".doctype-item").forEach((item) => {
    item.classList.toggle("active", item.dataset.doctype === name);
  });
}

function setStatus(message, type = "") {
  const statusPill = document.getElementById("statusPill");

  if (!statusPill) {
    return;
  }

  statusPill.textContent = message;
  statusPill.className = "status-pill";

  if (type) {
    statusPill.classList.add(type);
  }
}

function escapeHtml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}