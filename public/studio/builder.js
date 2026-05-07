const doctypeList = document.getElementById("doctypeList");
const pageTitle = document.getElementById("pageTitle");
const pageSubtitle = document.getElementById("pageSubtitle");
const statusPill = document.getElementById("statusPill");
const workspace = document.getElementById("workspace");
const propertiesPanel = document.getElementById("propertiesPanel");
const showMetaFields = document.getElementById("showMetaFields");
const backBtn = document.getElementById("backBtn");
const saveDraftBtn = document.getElementById("saveDraftBtn");

const builderPolicy = window.GOGAL_BUILDER_POLICY || {};
const fieldTypes = builderPolicy.fieldTypes || [];
const layoutFieldTypes = new Set(["Section Break", "Column Break", "Tab Break"]);
const systemFieldnames = new Set(builderPolicy.systemFieldnames || []);
const hiddenForNormalBuilder = new Set([
  ...(builderPolicy.doctypeMetaOnlyFields || []),
  ...(builderPolicy.docfieldAdvancedOnlyFields || []),
  "json_hash",
  "source_path",
  "created_at",
  "updated_at",
  "status",
  "oldfieldname",
  "oldfieldtype"
]);

let activeDocType = null;
let activeFields = [];
let activeActions = [];
let activeLinks = [];
let activePermissions = [];
let activeStates = [];
let activeTab = "fields";
let selection = { type: "doctype", index: -1 };
let dragged = null;

function setStatus(message, type = "") {
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

function boolValue(value) {
  return value === true || value === 1 || value === "1" || value === "true";
}

async function apiGet(url) {
  const res = await fetch(url);

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return res.json();
}

async function apiPost(url, body) {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.error || "Request failed");
  }

  return data;
}

function encodeDocTypeName(name) {
  return encodeURIComponent(name);
}

function slugify(value) {
  return String(value || "")
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "");
}

function fieldnameFromLabel(value) {
  return slugify(value);
}

function actionNameFromLabel(value) {
  return slugify(value);
}

function normalizeWhitespace(value) {
  return String(value || "").trim().replace(/\s+/g, " ");
}

function defaultTableName(value) {
  const normalized = normalizeWhitespace(value);
  return normalized ? `tab${normalized}` : "";
}

function isSystemField(field) {
  return systemFieldnames.has(field.fieldname || "") || hiddenForNormalBuilder.has(field.fieldname || "");
}

function shouldShowField(field) {
  if (showMetaFields.checked) {
    return true;
  }

  return !isSystemField(field);
}

function isNewDraftField(field) {
  return field && field.__draft === true;
}

function visibleFieldItems() {
  return activeFields
    .map((field, index) => ({ field, index }))
    .filter((item) => shouldShowField(item.field));
}

function updateDocTypeOptionRules() {
  if (!activeDocType) {
    return;
  }

  if (boolValue(activeDocType.editable_grid)) {
    activeDocType.is_child_table = true;
    activeDocType.is_single = false;
    activeDocType.is_submittable = false;
  }

  if (boolValue(activeDocType.is_child_table)) {
    activeDocType.is_single = false;
    activeDocType.is_submittable = false;
    if (activeDocType.editable_grid === undefined) {
      activeDocType.editable_grid = true;
    }
  }

  if (boolValue(activeDocType.is_submittable)) {
    activeDocType.is_single = false;
    activeDocType.is_child_table = false;
    activeDocType.editable_grid = false;
  }

  if (boolValue(activeDocType.is_single)) {
    activeDocType.is_submittable = false;
    activeDocType.is_child_table = false;
    activeDocType.editable_grid = false;
  }
}

async function loadDocTypes() {
  setStatus("Loading...");

  try {
    const result = await apiGet("/api/core/doctypes");
    const doctypes = result.data || [];

    doctypeList.innerHTML = "";

    if (doctypes.length === 0) {
      doctypeList.innerHTML = `<div class="muted">No DocTypes found.</div>`;
      setStatus("Ready");
      return;
    }

    doctypes.forEach((dt) => {
      const btn = document.createElement("button");
      btn.className = "doctype-item";
      btn.dataset.name = dt.name;
      btn.dataset.module = dt.module || "";
      btn.innerHTML = `
        <strong>${escapeHtml(dt.label || dt.name)}</strong>
        <small>${escapeHtml(dt.name)} - ${escapeHtml(dt.module)}</small>
      `;
      btn.addEventListener("click", () => loadDocType(dt.name));
      doctypeList.appendChild(btn);
    });

    setStatus("Ready", "success");
  } catch (err) {
    console.error(err);
    doctypeList.innerHTML = `<div class="muted">${escapeHtml(err.message)}</div>`;
    setStatus("Load failed", "error");
  }
}

async function loadDocType(name) {
  setStatus("Loading DocType...");
  selection = { type: "doctype", index: -1 };

  try {
    const encoded = encodeDocTypeName(name);
    const [docResult, fieldsResult, actionsResult, linksResult, permissionsResult, statesResult] = await Promise.all([
      apiGet(`/api/core/doctypes/${encoded}`),
      apiGet(`/api/core/doctypes/${encoded}/fields`),
      apiGet(`/api/core/doctypes/${encoded}/actions`),
      apiGet(`/api/core/doctypes/${encoded}/links`),
      apiGet(`/api/core/doctypes/${encoded}/permissions`),
      apiGet(`/api/core/doctypes/${encoded}/states`)
    ]);

    activeDocType = docResult.data || {};
    activeFields = (fieldsResult.data || []).map((field) => ({ ...field, __draft: false }));
    activeActions = (actionsResult.data || []).map((action) => ({ ...action, __draft: false }));
    activeLinks = linksResult.data || [];
    activePermissions = permissionsResult.data || [];
    activeStates = statesResult.data || [];

    pageTitle.textContent = activeDocType.label || activeDocType.name;
    pageSubtitle.textContent = `${activeDocType.name} - ${activeDocType.table_name || defaultTableName(activeDocType.name)}`;
    updateActiveSidebar(name);
    renderWorkspace();
    renderProperties();
    setStatus("Loaded", "success");
  } catch (err) {
    console.error(err);
    setStatus("Load failed", "error");
    workspace.innerHTML = `<div class="empty-state">${escapeHtml(err.message)}</div>`;
  }
}

function updateActiveSidebar(name) {
  document.querySelectorAll(".doctype-item").forEach((item) => {
    item.classList.toggle("active", item.dataset.name === name);
  });
}

function setActiveTab(tab) {
  activeTab = tab;
  document.querySelectorAll(".tab").forEach((tabButton) => {
    tabButton.classList.toggle("active", tabButton.dataset.tab === activeTab);
  });
  selection = selection.type === "doctype" ? selection : { type: activeTab === "actions" ? "action" : "field", index: selection.index };
  renderWorkspace();
  renderProperties();
}

function renderWorkspace() {
  if (!activeDocType) {
    workspace.innerHTML = `<div class="empty-state">Select a DocType or create a new draft.</div>`;
    return;
  }

  if (activeTab === "actions") {
    renderActions();
    return;
  }

  if (activeTab === "json") {
    renderJSONPreview();
    return;
  }

  renderFields();
}

function renderFields() {
  const items = visibleFieldItems();
  const body = items.map(({ field, index }) => {
    const layout = layoutFieldTypes.has(field.fieldtype || "");
    const system = isSystemField(field);
    return `
      <div class="field-card ${selection.type === "field" && selection.index === index ? "active" : ""} ${system ? "system" : ""} ${layout ? "layout" : ""}"
        draggable="true"
        data-kind="field"
        data-index="${index}">
        <div class="field-row-top">
          <div class="field-title">
            <strong>${escapeHtml(field.label || field.fieldname || "Untitled Field")}</strong>
            <small>${escapeHtml(field.fieldname || "")}</small>
          </div>
          <span class="badge ${layout ? "layout" : ""}">${escapeHtml(field.fieldtype || "Data")}</span>
        </div>
        <div class="field-flags">
          ${boolValue(field.required || field.reqd) ? `<span class="small-badge green">Required</span>` : ""}
          ${boolValue(field.unique || field.unique_field) ? `<span class="small-badge">Unique</span>` : ""}
          ${boolValue(field.in_list_view) ? `<span class="small-badge">List</span>` : ""}
          ${boolValue(field.read_only) ? `<span class="small-badge">Read Only</span>` : ""}
          ${boolValue(field.hidden) ? `<span class="small-badge red">Hidden</span>` : ""}
          ${system ? `<span class="small-badge">Meta</span>` : ""}
        </div>
      </div>
    `;
  }).join("");

  workspace.innerHTML = `
    <div class="toolbar">
      <div>
        <strong>${items.length} visible fields</strong>
        <span class="helper">Only business/user fields are shown unless meta/internal fields are enabled.</span>
      </div>
      <button id="addFieldBtn" class="primary-btn">+ Add Field</button>
    </div>
    ${body || `<div class="empty-state">No visible fields. Enable meta/internal fields to inspect hidden metadata.</div>`}
  `;

  document.getElementById("addFieldBtn").addEventListener("click", addField);
  bindSortableCards("field");
}

function renderActions() {
  const body = activeActions.map((action, index) => {
    return `
      <div class="action-card ${selection.type === "action" && selection.index === index ? "active" : ""}"
        draggable="true"
        data-kind="action"
        data-index="${index}">
        <div class="action-row-top">
          <div class="action-title">
            <strong>${escapeHtml(action.label || action.action_name || "Untitled Action")}</strong>
            <small>${escapeHtml(action.action_name || "")}</small>
          </div>
          <span class="badge">${escapeHtml(action.action_type || "server")}</span>
        </div>
        <div class="field-flags">
          <span class="small-badge">${escapeHtml(action.method || "POST")}</span>
          ${action.permission ? `<span class="small-badge">${escapeHtml(action.permission)}</span>` : ""}
          ${boolValue(action.enabled) ? `<span class="small-badge green">Enabled</span>` : `<span class="small-badge red">Disabled</span>`}
        </div>
      </div>
    `;
  }).join("");

  workspace.innerHTML = `
    <div class="toolbar">
      <div>
        <strong>${activeActions.length} actions</strong>
        <span class="helper">Actions are stored separately under draft JSON actions: [].</span>
      </div>
      <button id="addActionBtn" class="primary-btn">+ Add Action</button>
    </div>
    ${body || `<div class="empty-state">No actions yet.</div>`}
  `;

  document.getElementById("addActionBtn").addEventListener("click", addAction);
  bindSortableCards("action");
}

function renderJSONPreview() {
  const json = buildFullDocTypeJSON();
  workspace.innerHTML = `
    <div class="toolbar">
      <div>
        <strong>DocType JSON Preview</strong>
        <span class="helper">Preview only. No physical database migration is applied here.</span>
      </div>
      <button id="copyJsonBtn" class="ghost-btn">Copy JSON</button>
    </div>
    <pre id="jsonPreview" class="json-preview">${escapeHtml(JSON.stringify(json, null, 2))}</pre>
  `;

  document.getElementById("copyJsonBtn").addEventListener("click", copyJSONPreview);
}

function bindSortableCards(kind) {
  document.querySelectorAll(`[data-kind="${kind}"]`).forEach((card) => {
    card.addEventListener("click", () => {
      selection = { type: kind, index: Number(card.dataset.index) };
      renderWorkspace();
      renderProperties();
    });

    card.addEventListener("dragstart", () => {
      dragged = { kind, index: Number(card.dataset.index) };
      card.classList.add("dragging");
    });

    card.addEventListener("dragend", () => {
      card.classList.remove("dragging");
      dragged = null;
    });

    card.addEventListener("dragover", (event) => event.preventDefault());

    card.addEventListener("drop", (event) => {
      event.preventDefault();
      const targetIndex = Number(card.dataset.index);
      if (!dragged || dragged.kind !== kind || dragged.index === targetIndex) {
        return;
      }
      reorderCollection(kind, dragged.index, targetIndex);
    });
  });
}

function reorderCollection(kind, fromIndex, toIndex) {
  const list = kind === "field" ? activeFields : activeActions;
  const [item] = list.splice(fromIndex, 1);
  list.splice(toIndex, 0, item);
  list.forEach((entry, index) => {
    entry.idx = index + 1;
  });
  selection = { type: kind, index: toIndex };
  renderWorkspace();
  renderProperties();
  setStatus("Reordered", "success");
}

function addField() {
  const next = activeFields.length + 1;
  const field = {
    fieldname: `new_field_${next}`,
    label: `New Field ${next}`,
    fieldtype: "Data",
    options: "",
    required: false,
    reqd: false,
    unique: false,
    unique_field: false,
    hidden: false,
    read_only: false,
    in_list_view: false,
    default: "",
    description: "",
    depends_on: "",
    idx: next,
    __draft: true
  };

  activeFields.push(field);
  selection = { type: "field", index: activeFields.length - 1 };
  activeTab = "fields";
  renderTabState();
  renderWorkspace();
  renderProperties();
  setStatus("Field added", "success");
}

function addAction() {
  const next = activeActions.length + 1;
  const action = {
    action_name: `new_action_${next}`,
    label: `New Action ${next}`,
    action_type: "server",
    method: "POST",
    handler: "",
    route: "",
    permission: "System Manager",
    visible_when: "",
    enabled: true,
    idx: next,
    __draft: true
  };

  activeActions.push(action);
  selection = { type: "action", index: activeActions.length - 1 };
  activeTab = "actions";
  renderTabState();
  renderWorkspace();
  renderProperties();
  setStatus("Action added", "success");
}

function renderProperties() {
  if (!activeDocType) {
    propertiesPanel.innerHTML = `<div class="muted">No draft selected.</div>`;
    return;
  }

  if (selection.type === "field" && activeFields[selection.index]) {
    renderFieldProperties(activeFields[selection.index]);
    return;
  }

  if (selection.type === "action" && activeActions[selection.index]) {
    renderActionProperties(activeActions[selection.index]);
    return;
  }

  renderDocTypeProperties();
}

function renderDocTypeProperties() {
  updateDocTypeOptionRules();
  propertiesPanel.innerHTML = `
    <div class="notice">DocType settings update the draft JSON only. Module is a DocType property and never appears as a field card.</div>
    ${textInput("dtLabel", "Label", activeDocType.label || "")}
    ${textInput("dtName", "DocType Name", activeDocType.name || "", true)}
    ${textInput("dtModule", "Module Def", activeDocType.module || "")}
    <div class="prop-grid">
      ${checkInput("dtSingle", "Is Single", activeDocType.is_single, boolValue(activeDocType.is_submittable) || boolValue(activeDocType.is_child_table) || boolValue(activeDocType.editable_grid))}
      ${checkInput("dtSubmittable", "Is Submittable", activeDocType.is_submittable, boolValue(activeDocType.is_single) || boolValue(activeDocType.is_child_table) || boolValue(activeDocType.editable_grid))}
      ${checkInput("dtChildTable", "Is Child Table", activeDocType.is_child_table, boolValue(activeDocType.is_single) || boolValue(activeDocType.is_submittable))}
      ${checkInput("dtEditableGrid", "Editable Grid", activeDocType.editable_grid, boolValue(activeDocType.is_single) || boolValue(activeDocType.is_submittable))}
      ${checkInput("dtQuickEntry", "Quick Entry", activeDocType.quick_entry)}
      ${checkInput("dtTrackChanges", "Track Changes", activeDocType.track_changes)}
      ${checkInput("dtAllowImport", "Allow Import", activeDocType.allow_import)}
      ${checkInput("dtAllowExport", "Allow Export", activeDocType.allow_export)}
    </div>
    ${textInput("dtNamingRule", "Naming Rule", activeDocType.naming_rule || "autoname")}
    ${textInput("dtTitleField", "Title Field", activeDocType.title_field || "")}
    ${textInput("dtSortField", "Sort Field", activeDocType.sort_field || "created_at")}
    ${selectInput("dtSortOrder", "Sort Order", activeDocType.sort_order || "DESC", ["ASC", "DESC"])}
  `;

  bindValue("dtLabel", (value) => {
    activeDocType.label = value;
    pageTitle.textContent = value || activeDocType.name;
  });
  bindValue("dtModule", (value) => activeDocType.module = value);
  bindCheck("dtSingle", (value) => {
    activeDocType.is_single = value;
    updateDocTypeOptionRules();
    renderProperties();
  });
  bindCheck("dtSubmittable", (value) => {
    activeDocType.is_submittable = value;
    updateDocTypeOptionRules();
    renderProperties();
  });
  bindCheck("dtChildTable", (value) => {
    activeDocType.is_child_table = value;
    if (value) activeDocType.editable_grid = true;
    updateDocTypeOptionRules();
    renderProperties();
  });
  bindCheck("dtEditableGrid", (value) => {
    activeDocType.editable_grid = value;
    updateDocTypeOptionRules();
    renderProperties();
  });
  bindCheck("dtQuickEntry", (value) => activeDocType.quick_entry = value);
  bindCheck("dtTrackChanges", (value) => activeDocType.track_changes = value);
  bindCheck("dtAllowImport", (value) => activeDocType.allow_import = value);
  bindCheck("dtAllowExport", (value) => activeDocType.allow_export = value);
  bindValue("dtNamingRule", (value) => activeDocType.naming_rule = value);
  bindValue("dtTitleField", (value) => activeDocType.title_field = value);
  bindValue("dtSortField", (value) => activeDocType.sort_field = value);
  bindValue("dtSortOrder", (value) => activeDocType.sort_order = value);
}

function renderFieldProperties(field) {
  const fieldType = field.fieldtype || "Data";
  const system = isSystemField(field);
  const layout = layoutFieldTypes.has(fieldType);
  const optionsLabel = fieldType === "Link"
    ? "Link Target DocType"
    : fieldType === "Table"
      ? "Child Table DocType"
      : "Options";

  propertiesPanel.innerHTML = `
    ${layout ? `<div class="notice">Layout fields structure the form and are not normal data inputs.</div>` : ""}
    ${fieldType === "Button" ? `<div class="notice">Button handler wiring belongs in the Actions panel later.</div>` : ""}
    ${textInput("fieldLabel", "Label", field.label || "")}
    ${textInput("fieldName", "Field Name", field.fieldname || "", !isNewDraftField(field))}
    ${selectInput("fieldType", "Field Type", fieldType, fieldTypes)}
    ${fieldType === "Select"
      ? textArea("fieldOptions", "Options (one per line)", field.options || "", 5)
      : textArea("fieldOptions", optionsLabel, field.options || "", 3)}
    <div class="prop-grid">
      ${checkInput("fieldRequired", "Required", field.required || field.reqd)}
      ${checkInput("fieldUnique", "Unique", field.unique || field.unique_field)}
      ${checkInput("fieldHidden", "Hidden", field.hidden)}
      ${checkInput("fieldReadOnly", "Read Only", field.read_only)}
      ${checkInput("fieldListView", "In List View", field.in_list_view)}
    </div>
    ${textInput("fieldDefault", "Default", field.default || "")}
    ${textArea("fieldDescription", "Description", field.description || "", 3)}
    ${textInput("fieldDependsOn", "Depends On", field.depends_on || "")}
    <button id="selectDocTypeProps" class="ghost-btn">Edit DocType Properties</button>
    ${!system ? `<button id="removeFieldBtn" class="danger-btn">Remove Field</button>` : `<div class="notice">System/meta fields can be inspected but not removed.</div>`}
  `;

  bindValue("fieldLabel", (value) => {
    field.label = value;
    if (isNewDraftField(field) && !field.__fieldnameTouched) {
      field.fieldname = fieldnameFromLabel(value);
      const fieldNameInput = document.getElementById("fieldName");
      if (fieldNameInput) {
        fieldNameInput.value = field.fieldname;
      }
    }
    renderWorkspace();
  });
  bindValue("fieldName", (value) => {
    field.fieldname = fieldnameFromLabel(value);
    field.__fieldnameTouched = true;
    renderWorkspace();
  });
  bindValue("fieldType", (value) => {
    field.fieldtype = value;
    renderWorkspace();
    renderProperties();
  });
  bindValue("fieldOptions", (value) => field.options = value);
  bindCheck("fieldRequired", (value) => {
    field.required = value;
    field.reqd = value;
    renderWorkspace();
  });
  bindCheck("fieldUnique", (value) => {
    field.unique = value;
    field.unique_field = value;
    renderWorkspace();
  });
  bindCheck("fieldHidden", (value) => {
    field.hidden = value;
    renderWorkspace();
  });
  bindCheck("fieldReadOnly", (value) => {
    field.read_only = value;
    renderWorkspace();
  });
  bindCheck("fieldListView", (value) => {
    field.in_list_view = value;
    renderWorkspace();
  });
  bindValue("fieldDefault", (value) => field.default = value);
  bindValue("fieldDescription", (value) => field.description = value);
  bindValue("fieldDependsOn", (value) => field.depends_on = value);

  document.getElementById("selectDocTypeProps").addEventListener("click", () => {
    selection = { type: "doctype", index: -1 };
    renderWorkspace();
    renderProperties();
  });

  const removeBtn = document.getElementById("removeFieldBtn");
  if (removeBtn) {
    removeBtn.addEventListener("click", () => {
      activeFields.splice(selection.index, 1);
      activeFields.forEach((entry, index) => entry.idx = index + 1);
      selection = { type: "doctype", index: -1 };
      renderWorkspace();
      renderProperties();
      setStatus("Field removed", "success");
    });
  }
}

function renderActionProperties(action) {
  propertiesPanel.innerHTML = `
    ${textInput("actionLabel", "Label", action.label || "")}
    ${textInput("actionName", "Action Name", action.action_name || "")}
    ${selectInput("actionType", "Action Type", action.action_type || "server", ["server", "client", "route", "modal", "external"])}
    ${selectInput("actionMethod", "Method", action.method || "POST", ["GET", "POST", "PUT", "PATCH", "DELETE"])}
    ${textInput("actionHandler", "Handler", action.handler || "")}
    ${textInput("actionRoute", "Route", action.route || "")}
    ${textInput("actionPermission", "Permission", action.permission || "")}
    ${textInput("actionVisibleWhen", "Visible When", action.visible_when || "")}
    ${checkInput("actionEnabled", "Enabled", action.enabled !== false)}
    <button id="selectDocTypeProps" class="ghost-btn">Edit DocType Properties</button>
    <button id="removeActionBtn" class="danger-btn">Remove Action</button>
  `;

  bindValue("actionLabel", (value) => {
    action.label = value;
    if (!action.action_name || action.__draft) {
      action.action_name = actionNameFromLabel(value);
    }
    renderWorkspace();
  });
  bindValue("actionName", (value) => {
    action.action_name = actionNameFromLabel(value);
    renderWorkspace();
  });
  bindValue("actionType", (value) => {
    action.action_type = value;
    renderWorkspace();
  });
  bindValue("actionMethod", (value) => action.method = value);
  bindValue("actionHandler", (value) => action.handler = value);
  bindValue("actionRoute", (value) => action.route = value);
  bindValue("actionPermission", (value) => action.permission = value);
  bindValue("actionVisibleWhen", (value) => action.visible_when = value);
  bindCheck("actionEnabled", (value) => {
    action.enabled = value;
    renderWorkspace();
  });

  document.getElementById("selectDocTypeProps").addEventListener("click", () => {
    selection = { type: "doctype", index: -1 };
    renderWorkspace();
    renderProperties();
  });

  document.getElementById("removeActionBtn").addEventListener("click", () => {
    activeActions.splice(selection.index, 1);
    activeActions.forEach((entry, index) => entry.idx = index + 1);
    selection = { type: "doctype", index: -1 };
    renderWorkspace();
    renderProperties();
    setStatus("Action removed", "success");
  });
}

function textInput(id, label, value, disabled = false) {
  return `
    <div class="prop-group">
      <label for="${id}">${label}</label>
      <input id="${id}" value="${escapeHtml(value)}" ${disabled ? "disabled" : ""} />
    </div>
  `;
}

function textArea(id, label, value, rows = 3) {
  return `
    <div class="prop-group">
      <label for="${id}">${label}</label>
      <textarea id="${id}" rows="${rows}">${escapeHtml(value)}</textarea>
    </div>
  `;
}

function selectInput(id, label, value, options) {
  return `
    <div class="prop-group">
      <label for="${id}">${label}</label>
      <select id="${id}">
        ${options.map((option) => `<option value="${escapeHtml(option)}" ${option === value ? "selected" : ""}>${escapeHtml(option)}</option>`).join("")}
      </select>
    </div>
  `;
}

function checkInput(id, label, value, disabled = false) {
  return `
    <label class="check-row">
      <input id="${id}" type="checkbox" ${boolValue(value) ? "checked" : ""} ${disabled ? "disabled" : ""} />
      ${label}
    </label>
  `;
}

function bindValue(id, callback) {
  const input = document.getElementById(id);
  if (!input) return;
  input.addEventListener("input", (event) => callback(event.target.value));
  input.addEventListener("change", (event) => callback(event.target.value));
}

function bindCheck(id, callback) {
  const input = document.getElementById(id);
  if (!input) return;
  input.addEventListener("change", (event) => callback(event.target.checked));
}

function buildFullDocTypeJSON() {
  updateDocTypeOptionRules();
  return {
    name: activeDocType.name,
    module: activeDocType.module,
    label: activeDocType.label || activeDocType.name,
    table_name: activeDocType.table_name || defaultTableName(activeDocType.name),
    is_core: boolValue(activeDocType.is_core),
    is_single: boolValue(activeDocType.is_single),
    is_submittable: boolValue(activeDocType.is_submittable),
    is_child_table: boolValue(activeDocType.is_child_table),
    editable_grid: boolValue(activeDocType.editable_grid),
    quick_entry: boolValue(activeDocType.quick_entry),
    allow_import: boolValue(activeDocType.allow_import),
    allow_export: boolValue(activeDocType.allow_export),
    track_changes: boolValue(activeDocType.track_changes),
    naming_rule: activeDocType.naming_rule || "autoname",
    title_field: activeDocType.title_field || "",
    sort_field: activeDocType.sort_field || "created_at",
    sort_order: activeDocType.sort_order || "DESC",
    fields: activeFields.map(cleanDraftMarkers),
    actions: activeActions.map((action, index) => cleanDraftMarkers({ ...action, idx: index + 1 })),
    links: activeLinks || [],
    permissions: activePermissions || [],
    states: activeStates || []
  };
}

function cleanDraftMarkers(value) {
  const clone = { ...value };
  delete clone.__draft;
  delete clone.__fieldnameTouched;
  return clone;
}

async function copyJSONPreview() {
  const json = JSON.stringify(buildFullDocTypeJSON(), null, 2);
  await navigator.clipboard.writeText(json);
  setStatus("JSON copied", "success");
}

async function saveDraft() {
  if (!activeDocType) {
    alert("Select or create a DocType first.");
    return;
  }

  try {
    setStatus("Saving...");
    const json = buildFullDocTypeJSON();

    if (!json.name || !json.module || !json.table_name) {
      throw new Error("name, module, and table_name are required.");
    }

    const result = await apiPost("/api/core/doctypes", json);
    setStatus("Saved", "success");
    alert(`DocType JSON saved to ${result.file_path}. No physical migration was applied.`);
    await loadDocTypes();
    updateActiveSidebar(json.name);
  } catch (err) {
    console.error(err);
    setStatus("Save failed", "error");
    alert(err.message);
  }
}

function renderTabState() {
  document.querySelectorAll(".tab").forEach((tabButton) => {
    tabButton.classList.toggle("active", tabButton.dataset.tab === activeTab);
  });
}

document.querySelectorAll(".tab").forEach((tabButton) => {
  tabButton.addEventListener("click", () => setActiveTab(tabButton.dataset.tab));
});

showMetaFields.addEventListener("change", () => {
  selection = { type: "doctype", index: -1 };
  renderWorkspace();
  renderProperties();
});

backBtn.addEventListener("click", () => {
  window.location.href = "/studio";
});

saveDraftBtn.addEventListener("click", saveDraft);

loadDocTypes();
