import { getDocType, getDocTypeFields, getDocTypePermissions } from "../api/core.js";
import { getResourceList, updateDocument } from "../api/resource.js";
import { renderBreadcrumb } from "../components/breadcrumb.js";
import { renderDocTypeMetaCard } from "../components/doctype-meta-card.js";
import { renderFieldDetailPanel } from "../components/field-detail-panel.js";
import { setPanelHeader } from "../components/panel.js";
import { renderPermissionDetailPanel } from "../components/permission-detail-panel.js";
import { renderDocTypes } from "../components/sidebar.js";
import { setStatus } from "../components/status.js";
import { updateActiveTab } from "../components/tabs.js";
import { setTopbar } from "../components/topbar.js";
import { renderFields } from "../renderers/field.renderer.js";
import { renderReadOnlyForm } from "../renderers/form-view.renderer.js";
import { buildDocTypeJson, renderJsonPreview } from "../renderers/json.renderer.js";
import { renderDynamicListView } from "../renderers/list-view.renderer.js";
import { getPermissionKey, renderPermissions } from "../renderers/permission.renderer.js";
import { store } from "../state/store.js";
import { getVisibleDocTypes } from "./dashboard.page.js";

export function DocTypePage(router, params = {}) {
  async function mount() {
    const name = params.name || "DocType";

    if (store.activeDocType !== name) {
      store.listSearch = "";
      store.selectedRows = new Set();
      store.activeRecord = null;
      store.activeField = null;
      store.activePermission = null;
      store.activeView = "resource";
    }

    store.activeDocType = name;
    setStatus("Loading DocType...");

    const [dtRes, fieldsRes, permsRes, recordsRes] = await Promise.all([
      getDocType(name),
      getDocTypeFields(name),
      getDocTypePermissions(name),
      getResourceList(name)
    ]);

    store.activeBundle = {
      doctype: dtRes.data,
      fields: fieldsRes.data || [],
      permissions: permsRes.data || [],
      records: recordsRes.data || []
    };

    const doctype = store.activeBundle.doctype || {};
    setTopbar(doctype.label || name, `${doctype.module || ""} / ${doctype.table_name || ""}`);
    document.getElementById("recordCount").textContent = store.activeBundle.records.length;

    renderDocTypes(getVisibleDocTypes(), (doctypeName) => {
      router.navigate("doctype", { name: doctypeName });
    });

    renderActiveView();

    setStatus("Loaded", "success");
  }

  return { mount };
}

export function renderActiveView() {
  if (!store.activeBundle) return;

  updateActiveTab();

  const mainPanel = document.getElementById("mainPanel");
  const dtName = store.activeBundle.doctype.name;
  const breadcrumb = renderBreadcrumb(store.activeBundle, store.activeView, store.activeRecord);
  const metaCard = renderDocTypeMetaCard(store.activeBundle);

  if (store.activeView === "resource") {
    setPanelHeader(`${dtName} Resource List`, `/api/resource/${dtName}?user=Administrator&limit=20`);
    mainPanel.innerHTML = breadcrumb + metaCard + renderDynamicListView(store.activeBundle, store.listSearch);
    bindListEvents();
    return;
  }

  if (store.activeView === "form") {
    const recordName = store.activeRecord?.name || "Read-only record";
    setPanelHeader(`${dtName} Form`, recordName);
    mainPanel.innerHTML = breadcrumb + metaCard + renderReadOnlyForm(store.activeBundle, store.activeRecord);
    bindFormEvents();
    return;
  }

  if (store.activeView === "fields") {
    setPanelHeader(`${dtName} Fields`, "Read-only DocField metadata");
    const activeFieldName = store.activeField?.fieldname || "";
    mainPanel.innerHTML = `
      ${breadcrumb}
      ${metaCard}
      <div class="gs-fields-layout">
        <div class="gs-fields-table-wrap">${renderFields(store.activeBundle.fields, activeFieldName)}</div>
        <div class="gs-field-detail-wrap">${renderFieldDetailPanel(store.activeField)}</div>
      </div>
    `;
    bindFieldEvents();
    return;
  }

  if (store.activeView === "permissions") {
    setPanelHeader(`${dtName} Permissions`, "Read-only DocPerm metadata");
    const activePermissionKey = store.activePermission
      ? getPermissionKey(store.activePermission, store.activeBundle.permissions.indexOf(store.activePermission))
      : "";
    mainPanel.innerHTML = `
      ${breadcrumb}
      ${metaCard}
      <div class="gs-permissions-layout">
        <div class="gs-permissions-table-wrap">${renderPermissions(store.activeBundle.permissions, activePermissionKey)}</div>
        <div class="gs-permission-detail-wrap">${renderPermissionDetailPanel(store.activePermission)}</div>
      </div>
    `;
    bindPermissionEvents();
    return;
  }

  if (store.activeView === "json") {
    setPanelHeader(`${dtName} JSON Preview`, "Generated from current API response");
    mainPanel.innerHTML = breadcrumb + metaCard + renderJsonPreview(store.activeBundle);
    bindJsonEvents();
  }
}

function bindListEvents() {
  const searchInput = document.getElementById("listSearchInput");
  const refreshBtn = document.getElementById("listRefreshBtn");

  if (searchInput) {
    searchInput.addEventListener("input", (event) => {
      const cursorStart = event.target.selectionStart;
      const cursorEnd = event.target.selectionEnd;

      store.listSearch = event.target.value;
      renderActiveView();

      const nextSearchInput = document.getElementById("listSearchInput");
      if (nextSearchInput) {
        nextSearchInput.focus();

        if (cursorStart !== null && cursorEnd !== null) {
          nextSearchInput.setSelectionRange(cursorStart, cursorEnd);
        }
      }
    });
  }

  if (refreshBtn) {
    refreshBtn.addEventListener("click", () => {
      window.gogalRouter.navigate("doctype", { name: store.activeDocType });
    });
  }

  document.querySelectorAll(".gs-row-select").forEach((checkbox) => {
    checkbox.addEventListener("click", (event) => {
      event.stopPropagation();
    });
  });

  document.querySelectorAll(".gs-row-clickable").forEach((row) => {
    row.addEventListener("click", (event) => {
      if (event.target instanceof HTMLInputElement) return;
      const name = row.dataset.rowName || "";
      const index = Number(row.dataset.rowIndex || "-1");
      const record = store.activeBundle.records.find((item) => item.name === name) || store.activeBundle.records[index];

      console.log("Open form placeholder:", name);

      store.activeRecord = record || null;
      store.previousView = store.activeView;
      store.activeView = "form";
      renderActiveView();
    });
  });
}

function bindFormEvents() {
  const backBtn = document.getElementById("backToListBtn");
  const saveBtn = document.getElementById("saveRecordBtn");

  if (backBtn) {
    backBtn.addEventListener("click", () => {
      store.activeRecord = null;
      store.activeView = "resource";
      renderActiveView();
    });
  }

  if (saveBtn) {
    saveBtn.addEventListener("click", saveActiveRecord);
  }
}

function normalizeInputValue(input) {
  const fieldtype = input.dataset.fieldtype || "Data";

  if (input.type === "checkbox") {
    return input.checked;
  }

  if (["Int"].includes(fieldtype)) {
    return input.value === "" ? null : Number.parseInt(input.value, 10);
  }

  if (["Float", "Currency"].includes(fieldtype)) {
    return input.value === "" ? null : Number.parseFloat(input.value);
  }

  if (["JSON"].includes(fieldtype)) {
    if (input.value.trim() === "") return null;

    try {
      return JSON.parse(input.value);
    } catch {
      return input.value;
    }
  }

  return input.value;
}

function collectFormUpdates() {
  const updates = {};

  document.querySelectorAll(".gs-form-input[data-fieldname]").forEach((input) => {
    const fieldname = input.dataset.fieldname;

    if (!fieldname) return;

    updates[fieldname] = normalizeInputValue(input);
  });

  return updates;
}

async function saveActiveRecord() {
  if (!store.activeBundle || !store.activeRecord?.name) {
    setStatus("No record selected", "error");
    return;
  }

  const doctype = store.activeBundle.doctype.name;
  const name = store.activeRecord.name;
  const saveBtn = document.getElementById("saveRecordBtn");

  try {
    if (saveBtn) saveBtn.disabled = true;
    setStatus("Saving...");

    const res = await updateDocument(doctype, name, collectFormUpdates());
    const saved = res.data || {};
    const existingIndex = store.activeBundle.records.findIndex((row) => row.name === name);

    store.activeRecord = saved;

    if (existingIndex >= 0) {
      store.activeBundle.records[existingIndex] = saved;
    }

    setStatus("Saved", "success");
    renderActiveView();
  } catch (error) {
    console.error(error);
    setStatus("Save failed", "error");
  } finally {
    const nextSaveBtn = document.getElementById("saveRecordBtn");
    if (nextSaveBtn) nextSaveBtn.disabled = false;
  }
}

function bindFieldEvents() {
  document.querySelectorAll(".gs-field-row").forEach((row) => {
    row.addEventListener("click", () => {
      const fieldname = row.dataset.fieldname || "";
      store.activeField = store.activeBundle.fields.find((field) => field.fieldname === fieldname) || null;
      renderActiveView();
    });
  });
}

function bindPermissionEvents() {
  document.querySelectorAll(".gs-permission-row").forEach((row) => {
    row.addEventListener("click", () => {
      const permissionKey = row.dataset.permissionKey || "";
      store.activePermission = store.activeBundle.permissions.find((permission, index) => (
        getPermissionKey(permission, index) === permissionKey
      )) || null;
      renderActiveView();
    });
  });
}

function safeJsonFileName(name) {
  return `${String(name || "doctype").trim().toLowerCase().replace(/\s+/g, "_")}.json`;
}

function currentJsonText() {
  return JSON.stringify(buildDocTypeJson(store.activeBundle), null, 2);
}

function bindJsonEvents() {
  const copyBtn = document.getElementById("copyJsonBtn");
  const downloadBtn = document.getElementById("downloadJsonBtn");

  if (copyBtn) {
    copyBtn.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText(currentJsonText());
        setStatus("JSON copied", "success");
      } catch (error) {
        console.error(error);
        setStatus("Copy failed", "error");
      }
    });
  }

  if (downloadBtn) {
    downloadBtn.addEventListener("click", () => {
      const doctypeName = store.activeBundle?.doctype?.name || store.activeDocType;
      const blob = new Blob([currentJsonText()], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");

      link.href = url;
      link.download = safeJsonFileName(doctypeName);
      document.body.appendChild(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);

      setStatus("JSON downloaded", "success");
    });
  }
}
