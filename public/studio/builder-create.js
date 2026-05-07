const createDocTypeBtn = document.getElementById("createDocTypeBtn");

let createOptions = {
  is_single: false,
  is_submittable: false,
  is_child_table: false,
  editable_grid: false
};

function getExistingModulesFromSidebar() {
  const modules = new Set();

  document.querySelectorAll(".doctype-item").forEach((item) => {
    const moduleName = item.dataset.module;
    if (moduleName) {
      modules.add(moduleName);
    }
  });

  if (modules.size === 0) {
    modules.add("Core");
  }

  return Array.from(modules).sort();
}

function openCreateDocTypePanel() {
  activeDocType = null;
  activeFields = [];
  activeActions = [];
  activeLinks = [];
  activePermissions = [];
  activeStates = [];
  selection = { type: "doctype", index: -1 };
  activeTab = "fields";
  renderTabState();

  createOptions = {
    is_single: false,
    is_submittable: false,
    is_child_table: false,
    editable_grid: false
  };

  const modules = getExistingModulesFromSidebar();

  pageTitle.textContent = "Create DocType";
  pageSubtitle.textContent = "Choose module and DocType behavior before editing fields.";

  workspace.innerHTML = `
    <div class="create-panel">
      <h3>New DocType</h3>

      <div class="create-grid">
        <div class="create-field">
          <label for="newDocTypeName">DocType Name</label>
          <input id="newDocTypeName" placeholder="Example: Customer Group" />
        </div>

        <div class="create-field">
          <label for="newModuleSelect">Module Def</label>
          <select id="newModuleSelect">
            ${modules.map((moduleName) => `<option value="${escapeHtml(moduleName)}">${escapeHtml(moduleName)}</option>`).join("")}
            <option value="__create_new__">+ Create new module...</option>
          </select>
        </div>

        <div class="create-field" id="newModuleBox" style="display:none;">
          <label for="newModuleName">New Module Name</label>
          <input id="newModuleName" placeholder="Example: Selling" />
        </div>

        <div class="option-box full-width" id="singleOption">
          <label>
            <input id="newIsSingle" type="checkbox" />
            <strong>Is Single</strong>
            <small>Settings-style DocType with one record.</small>
          </label>
        </div>

        <div class="option-box full-width" id="submittableOption">
          <label>
            <input id="newIsSubmittable" type="checkbox" />
            <strong>Is Submittable</strong>
            <small>Submitted records become locked except cancel/amend flows later.</small>
          </label>
        </div>

        <div class="option-box full-width" id="childTableOption">
          <label>
            <input id="newIsChildTable" type="checkbox" />
            <strong>Is Child Table / istable</strong>
            <small>Used inside another DocType as a table/grid.</small>
          </label>
        </div>

        <div class="option-box full-width" id="editableGridOption">
          <label>
            <input id="newEditableGrid" type="checkbox" />
            <strong>Editable Grid</strong>
            <small>Forces Child Table and defaults child rows to grid editing.</small>
          </label>
        </div>
      </div>

      <div class="create-actions">
        <button id="cancelCreateDocType" class="ghost-btn">Cancel</button>
        <button id="createDraftDocType" class="primary-btn">Create Draft</button>
      </div>
    </div>
  `;

  propertiesPanel.innerHTML = `
    <div class="notice">
      Module Def is stored on the DocType draft. It does not become a field on the canvas.
    </div>
  `;

  bindCreateDocTypeEvents();
  applyDocTypeCreateOptionRules();
  setStatus("Draft setup");
}

function bindCreateDocTypeEvents() {
  const nameInput = document.getElementById("newDocTypeName");
  const moduleSelect = document.getElementById("newModuleSelect");
  const newModuleBox = document.getElementById("newModuleBox");

  moduleSelect.addEventListener("change", () => {
    newModuleBox.style.display = moduleSelect.value === "__create_new__" ? "flex" : "none";
  });

  document.getElementById("newIsSingle").addEventListener("change", (event) => {
    createOptions.is_single = event.target.checked;
    if (createOptions.is_single) {
      createOptions.is_submittable = false;
      createOptions.is_child_table = false;
      createOptions.editable_grid = false;
    }
    applyDocTypeCreateOptionRules();
  });

  document.getElementById("newIsSubmittable").addEventListener("change", (event) => {
    createOptions.is_submittable = event.target.checked;
    if (createOptions.is_submittable) {
      createOptions.is_single = false;
      createOptions.is_child_table = false;
      createOptions.editable_grid = false;
    }
    applyDocTypeCreateOptionRules();
  });

  document.getElementById("newIsChildTable").addEventListener("change", (event) => {
    createOptions.is_child_table = event.target.checked;
    if (createOptions.is_child_table) {
      createOptions.is_single = false;
      createOptions.is_submittable = false;
      createOptions.editable_grid = true;
    } else {
      createOptions.editable_grid = false;
    }
    applyDocTypeCreateOptionRules();
  });

  document.getElementById("newEditableGrid").addEventListener("change", (event) => {
    createOptions.editable_grid = event.target.checked;
    if (createOptions.editable_grid) {
      createOptions.is_child_table = true;
      createOptions.is_single = false;
      createOptions.is_submittable = false;
    }
    applyDocTypeCreateOptionRules();
  });

  document.getElementById("cancelCreateDocType").addEventListener("click", () => {
    workspace.innerHTML = `<div class="empty-state">Select a DocType or create a new draft.</div>`;
    propertiesPanel.innerHTML = `<div class="muted">No draft selected.</div>`;
    setStatus("Ready");
  });

  document.getElementById("createDraftDocType").addEventListener("click", createDocTypeDraftFromForm);
}

function applyDocTypeCreateOptionRules() {
  const isSingle = document.getElementById("newIsSingle");
  const isSubmittable = document.getElementById("newIsSubmittable");
  const isChildTable = document.getElementById("newIsChildTable");
  const editableGrid = document.getElementById("newEditableGrid");
  const singleBox = document.getElementById("singleOption");
  const submittableBox = document.getElementById("submittableOption");
  const childBox = document.getElementById("childTableOption");
  const gridBox = document.getElementById("editableGridOption");

  isSingle.checked = createOptions.is_single;
  isSubmittable.checked = createOptions.is_submittable;
  isChildTable.checked = createOptions.is_child_table;
  editableGrid.checked = createOptions.editable_grid;

  [singleBox, submittableBox, childBox, gridBox].forEach((box) => box.classList.remove("disabled"));
  [isSingle, isSubmittable, isChildTable, editableGrid].forEach((input) => input.disabled = false);

  if (createOptions.is_single) {
    submittableBox.classList.add("disabled");
    childBox.classList.add("disabled");
    gridBox.classList.add("disabled");
    isSubmittable.disabled = true;
    isChildTable.disabled = true;
    editableGrid.disabled = true;
  }

  if (createOptions.is_submittable) {
    singleBox.classList.add("disabled");
    childBox.classList.add("disabled");
    gridBox.classList.add("disabled");
    isSingle.disabled = true;
    isChildTable.disabled = true;
    editableGrid.disabled = true;
  }

  if (createOptions.is_child_table || createOptions.editable_grid) {
    singleBox.classList.add("disabled");
    submittableBox.classList.add("disabled");
    isSingle.disabled = true;
    isSubmittable.disabled = true;
  }
}

function createDocTypeDraftFromForm() {
  const name = document.getElementById("newDocTypeName").value.trim();
  const moduleSelect = document.getElementById("newModuleSelect").value;
  const newModuleName = document.getElementById("newModuleName")?.value.trim();

  if (!name) {
    alert("DocType Name is required.");
    return;
  }

  let moduleName = moduleSelect;
  if (moduleSelect === "__create_new__") {
    if (!newModuleName) {
      alert("New Module Name is required.");
      return;
    }
    moduleName = newModuleName;
  }

  activeDocType = {
    name,
    module: moduleName,
    label: name,
    table_name: defaultTableName(name),
    is_core: false,
    is_single: createOptions.is_single,
    is_submittable: createOptions.is_submittable,
    is_child_table: createOptions.is_child_table,
    editable_grid: createOptions.editable_grid,
    quick_entry: true,
    allow_import: true,
    allow_export: true,
    track_changes: true,
    naming_rule: "autoname",
    title_field: "",
    sort_field: "created_at",
    sort_order: "DESC"
  };

  activeFields = [];
  activeActions = [];
  activeLinks = [];
  activePermissions = [
    {
      role: "System Manager",
      permlevel: 0,
      read: true,
      write: true,
      create: true,
      delete: true,
      export: true,
      import: true,
      idx: 1
    }
  ];
  activeStates = [];
  selection = { type: "doctype", index: -1 };

  pageTitle.textContent = activeDocType.label;
  pageSubtitle.textContent = `${activeDocType.name} - ${activeDocType.table_name}`;
  renderWorkspace();
  renderProperties();
  setStatus("Draft created", "success");
}

createDocTypeBtn.addEventListener("click", openCreateDocTypePanel);
