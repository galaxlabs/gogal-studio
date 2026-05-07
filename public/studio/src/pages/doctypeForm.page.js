import {
  getDocType,
  getDocTypeFields,
  getDocTypePermissions,
  saveDocType
} 
from "../api/core.js";

function setStatus(message, type = "") {
  const statusPill = document.getElementById("statusPill");

  statusPill.textContent = message;
  statusPill.className = "status-pill";

  if (type) {
    statusPill.classList.add(type);
  }
}

function defaultDocTypeDraft() {
  return {
    name: "",
    module: "Core",
    app_name: "gogal_studio",
    table_name: "",
    is_single: false,
    is_submittable: false,
    is_child_table: false,
    is_tree: false,
    fields: [
      {
        idx: 1,
        fieldname: "section_break_1",
        label: "Section",
        fieldtype: "Section Break",
        options: "",
        reqd: false,
        in_list_view: false
      }
    ],
    permissions: [
      {
        role: "System Manager",
        permlevel: 0,
        read: true,
        write: true,
        create: true,
        delete: true
      }
    ]
  };
}

export function showNewDocTypeForm() {
  const draft = defaultDocTypeDraft();

  document.getElementById("pageTitle").textContent = "New DocType";
  document.getElementById("pageSubtitle").textContent = "Create metadata draft. Physical table migration comes later.";

  renderDocTypeForm(draft, "new");
  setStatus("Draft", "success");
}

export async function showEditDocTypeForm(name) {
  setStatus("Loading DocType...");

  try {
    const [dtRes, fieldsRes, permsRes] = await Promise.all([
      getDocType(name),
      getDocTypeFields(name),
      getDocTypePermissions(name)
    ]);

    const doc = {
      ...dtRes.data,
      fields: fieldsRes.data || [],
      permissions: permsRes.data || []
    };

    document.getElementById("pageTitle").textContent = `Edit ${doc.name}`;
    document.getElementById("pageSubtitle").textContent = `${doc.module} · ${doc.table_name}`;

    renderDocTypeForm(doc, "edit");
    setStatus("Loaded", "success");
  } catch (err) {
    console.error(err);
    setStatus("Load failed", "error");
  }
}

function renderDocTypeForm(doc, mode) {
  const detailsPanel = document.getElementById("detailsPanel");
  const fieldsPanel = document.getElementById("fieldsPanel");
  const permsPanel = document.getElementById("permsPanel");

  detailsPanel.innerHTML = `
    <div class="form-grid">
      ${inputField("name", "DocType Name", doc.name, mode === "edit")}
      ${inputField("module", "Module", doc.module)}
      ${inputField("app_name", "App Name", doc.app_name)}
      ${inputField("table_name", "Table Name", doc.table_name)}

      ${checkField("is_single", "Is Single", doc.is_single)}
      ${checkField("is_child_table", "Is Child Table", doc.is_child_table)}
      ${checkField("is_submittable", "Is Submittable", doc.is_submittable)}
      ${checkField("is_tree", "Is Tree", doc.is_tree)}
    </div>

    <div class="topbar-actions" style="margin-top:16px;">
      <button id="addFieldBtn" class="secondary-btn">+ Add Field</button>
      <button id="previewJsonBtn" class="secondary-btn">Preview JSON</button>
      <button id="saveDocTypeBtn" class="primary-btn">Save Metadata</button>
    </div>
  `;

  fieldsPanel.innerHTML = renderEditableFields(doc.fields || []);
  permsPanel.innerHTML = renderEditablePermissions(doc.permissions || []);

  document.getElementById("addFieldBtn").addEventListener("click", () => {
    addFieldRow();
  });

  document.getElementById("previewJsonBtn").addEventListener("click", () => {
    showJSONPreview(collectDocTypeForm());
  });

document.getElementById("saveDocTypeBtn").addEventListener("click", async () => {
  const payload = collectDocTypeForm();
  setStatus("Saving...");

  try {
    const result = await saveDocType(payload);
    showJSONPreview({
      result,
      payload,
      note: "Metadata saved only. Physical migration not applied yet."
    });
    setStatus("Saved", "success");
  } catch (err) {
    console.error(err);
    showJSONPreview({ error: err.message, payload });
    setStatus("Save failed", "error");
  }
});

  bindDynamicFormRules();
}

function inputField(name, label, value = "", readonly = false) {
  return `
    <label class="form-field">
      <span>${label}</span>
      <input name="${name}" value="${value || ""}" ${readonly ? "readonly" : ""} />
    </label>
  `;
}

function checkField(name, label, checked = false) {
  return `
    <label class="check-field">
      <input type="checkbox" name="${name}" ${checked ? "checked" : ""} />
      <span>${label}</span>
    </label>
  `;
}

function renderEditableFields(fields) {
  return `
    <div class="topbar-actions" style="margin-bottom:12px;">
      <strong>Fields</strong>
    </div>

    <table class="table" id="fieldsTable">
      <thead>
        <tr>
          <th>Idx</th>
          <th>Fieldname</th>
          <th>Label</th>
          <th>Type</th>
          <th>Options</th>
          <th>Reqd</th>
          <th>List</th>
        </tr>
      </thead>
      <tbody>
        ${(fields || []).map((field, index) => fieldRow(field, index + 1)).join("")}
      </tbody>
    </table>
  `;
}

function fieldRow(field = {}, idx = 1) {
  return `
    <tr class="field-row">
      <td><input name="idx" value="${field.idx || idx}" /></td>
      <td><input name="fieldname" value="${field.fieldname || ""}" /></td>
      <td><input name="label" value="${field.label || ""}" data-auto-fieldname="1" /></td>
      <td>
        <select name="fieldtype">
          ${fieldTypeOptions(field.fieldtype || "Data")}
        </select>
      </td>
      <td><input name="options" value="${field.options || ""}" /></td>
      <td><input type="checkbox" name="reqd" ${field.reqd ? "checked" : ""} /></td>
      <td><input type="checkbox" name="in_list_view" ${field.in_list_view ? "checked" : ""} /></td>
    </tr>
  `;
}

function fieldTypeOptions(selected) {
  const types = [
    "Data",
    "Text",
    "Small Text",
    "Long Text",
    "Int",
    "Float",
    "Currency",
    "Check",
    "Date",
    "Datetime",
    "Select",
    "Link",
    "Table",
    "Attach",
    "JSON",
    "Code",
    "Section Break",
    "Column Break",
    "Tab Break",
    "Button",
    "HTML",
    "Read Only"
  ];

  return types.map((type) => `
    <option value="${type}" ${type === selected ? "selected" : ""}>${type}</option>
  `).join("");
}

function renderEditablePermissions(perms) {
  return `
    <div class="topbar-actions" style="margin-bottom:12px;">
      <strong>Permissions</strong>
    </div>

    <table class="table" id="permissionsTable">
      <thead>
        <tr>
          <th>Role</th>
          <th>Level</th>
          <th>Read</th>
          <th>Write</th>
          <th>Create</th>
          <th>Delete</th>
        </tr>
      </thead>
      <tbody>
        ${(perms || []).map((perm) => `
          <tr class="perm-row">
            <td><input name="role" value="${perm.role || "System Manager"}" /></td>
            <td><input name="permlevel" value="${perm.permlevel || 0}" /></td>
            <td><input type="checkbox" name="read" ${readBool(perm, "read") ? "checked" : ""} /></td>
            <td><input type="checkbox" name="write" ${readBool(perm, "write") ? "checked" : ""} /></td>
            <td><input type="checkbox" name="create" ${readBool(perm, "create") ? "checked" : ""} /></td>
            <td><input type="checkbox" name="delete" ${readBool(perm, "delete") ? "checked" : ""} /></td>
          </tr>
        `).join("")}
      </tbody>
    </table>
  `;
}

function readBool(row, key) {
  return row[key] ?? row[`${key}_perm`] ?? false;
}

function addFieldRow() {
  const tbody = document.querySelector("#fieldsTable tbody");
  const idx = tbody.querySelectorAll("tr").length + 1;

  tbody.insertAdjacentHTML("beforeend", fieldRow({}, idx));
  bindDynamicFormRules();
}

function bindDynamicFormRules() {
  document.querySelectorAll('#fieldsTable input[name="label"]').forEach((input) => {
    if (input.dataset.bound) return;

    input.addEventListener("input", () => {
      const row = input.closest("tr");
      const fieldnameInput = row.querySelector('[name="fieldname"]');

      if (!fieldnameInput.value || fieldnameInput.dataset.auto === "1") {
        fieldnameInput.value = slugFieldname(input.value);
        fieldnameInput.dataset.auto = "1";
      }
    });

    input.dataset.bound = "1";
  });

  document.querySelectorAll('#fieldsTable input[name="fieldname"]').forEach((input) => {
    if (input.dataset.bound) return;

    input.addEventListener("input", () => {
      input.dataset.auto = "0";
      input.value = slugFieldname(input.value);
    });

    input.dataset.bound = "1";
  });

  ["is_single", "is_child_table", "is_submittable", "is_tree"].forEach((name) => {
    const input = document.querySelector(`#detailsPanel [name="${name}"]`);
    if (!input || input.dataset.bound) return;

    input.addEventListener("change", applyDocTypeOptionRules);
    input.dataset.bound = "1";
  });

  applyDocTypeOptionRules();
}

function applyDocTypeOptionRules() {
  const detailsPanel = document.getElementById("detailsPanel");

  const isSingle = detailsPanel.querySelector('[name="is_single"]');
  const isChildTable = detailsPanel.querySelector('[name="is_child_table"]');
  const isSubmittable = detailsPanel.querySelector('[name="is_submittable"]');
  const isTree = detailsPanel.querySelector('[name="is_tree"]');

  if (!isSingle || !isChildTable || !isSubmittable || !isTree) return;

  if (isChildTable.checked) {
    isSingle.checked = false;
    isSubmittable.checked = false;
    isTree.checked = false;
  }

  if (isSingle.checked) {
    isChildTable.checked = false;
    isSubmittable.checked = false;
  }

  if (isSubmittable.checked) {
    isChildTable.checked = false;
    isSingle.checked = false;
  }

  toggleFieldWrapper(isChildTable, !isSingle.checked && !isSubmittable.checked);
  toggleFieldWrapper(isSubmittable, !isChildTable.checked && !isSingle.checked);
  toggleFieldWrapper(isSingle, !isChildTable.checked && !isSubmittable.checked);
  toggleFieldWrapper(isTree, !isChildTable.checked);
}

function toggleFieldWrapper(input, show) {
  const wrapper = input.closest(".check-field");
  if (!wrapper) return;

  wrapper.style.display = show ? "flex" : "none";
}

function slugFieldname(value) {
  return (value || "")
    .trim()
    .toLowerCase()
    .replaceAll("-", "_")
    .replace(/\s+/g, "_")
    .replace(/[^a-z0-9_]/g, "")
    .replace(/_+/g, "_")
    .replace(/^_+|_+$/g, "");
}

function collectDocTypeForm() {
  const detailsPanel = document.getElementById("detailsPanel");

  const getValue = (name) => detailsPanel.querySelector(`[name="${name}"]`)?.value || "";
  const getChecked = (name) => detailsPanel.querySelector(`[name="${name}"]`)?.checked || false;

  const fields = Array.from(document.querySelectorAll("#fieldsTable tbody tr")).map((row) => {
    return {
      idx: Number(row.querySelector(`[name="idx"]`)?.value || 0),
      fieldname: row.querySelector(`[name="fieldname"]`)?.value || "",
      label: row.querySelector(`[name="label"]`)?.value || "",
      fieldtype: row.querySelector(`[name="fieldtype"]`)?.value || "Data",
      options: row.querySelector(`[name="options"]`)?.value || "",
      reqd: row.querySelector(`[name="reqd"]`)?.checked || false,
      in_list_view: row.querySelector(`[name="in_list_view"]`)?.checked || false
    };
  });

  const permissions = Array.from(document.querySelectorAll("#permissionsTable tbody tr")).map((row) => {
    return {
      role: row.querySelector(`[name="role"]`)?.value || "System Manager",
      permlevel: Number(row.querySelector(`[name="permlevel"]`)?.value || 0),
      read: row.querySelector(`[name="read"]`)?.checked || false,
      write: row.querySelector(`[name="write"]`)?.checked || false,
      create: row.querySelector(`[name="create"]`)?.checked || false,
      delete: row.querySelector(`[name="delete"]`)?.checked || false
    };
  });

  const name = getValue("name");

  return {
    name,
    module: getValue("module"),
    app_name: getValue("app_name"),
    table_name: getValue("table_name") || `tab${name}`,
    is_single: getChecked("is_single"),
    is_child_table: getChecked("is_child_table"),
    is_submittable: getChecked("is_submittable"),
    is_tree: getChecked("is_tree"),
    fields,
    permissions
  };
}

function showJSONPreview(payload) {
  document.getElementById("permsPanel").innerHTML = `
    <pre class="json-preview">${escapeHtml(JSON.stringify(payload, null, 2))}</pre>
  `;
}

function escapeHtml(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}