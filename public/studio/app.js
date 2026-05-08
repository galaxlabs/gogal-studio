(function () {
  const els = {
    moduleNav: document.getElementById("moduleNav"),
    doctypeList: document.getElementById("doctypeList"),
    refreshBtn: document.getElementById("refreshBtn"),
    statusPill: document.getElementById("statusPill"),
    pageTitle: document.getElementById("pageTitle"),
    pageSubtitle: document.getElementById("pageSubtitle"),
    appCount: document.getElementById("appCount"),
    moduleCount: document.getElementById("moduleCount"),
    doctypeCount: document.getElementById("doctypeCount"),
    recordCount: document.getElementById("recordCount"),
    panelTitle: document.getElementById("panelTitle"),
    panelSubtitle: document.getElementById("panelSubtitle"),
    mainPanel: document.getElementById("mainPanel")
  };

  const state = {
    apps: [],
    modules: [],
    doctypes: [],
    selectedModule: "All",
    selectedDoctypeName: "",
    selectedDoctype: null,
    fields: [],
    permissions: [],
    records: [],
    activeTab: "resources"
  };

  function escapeHtml(value) {
    return String(value ?? "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#39;");
  }

  function boolValue(value) {
    return value === true || value === 1 || value === "1" || value === "true";
  }

  function boolBadge(value) {
    const on = boolValue(value);
    return `<span class="gs-badge ${on ? "on" : "off"}">${on ? "Yes" : "No"}</span>`;
  }

  function setStatus(message, type) {
    els.statusPill.textContent = message;
    els.statusPill.className = "gs-status-pill";
    if (type) {
      els.statusPill.classList.add(type);
    }
  }

  async function apiGet(url) {
    const res = await fetch(url, { headers: { Accept: "application/json" } });
    if (!res.ok) {
      const detail = await res.text();
      throw new Error(detail || `${res.status} ${res.statusText}`);
    }
    return res.json();
  }

  function dataOf(response) {
    return Array.isArray(response?.data) ? response.data : [];
  }

  function encodedName(name) {
    return encodeURIComponent(name);
  }

  function moduleLabel(module) {
    return module?.module_name || module?.name || "";
  }

  function doctypeSubtitle(dt) {
    const moduleName = dt?.module || "No module";
    const tableName = dt?.table_name || "No table";
    return `${moduleName} / ${tableName}`;
  }

  async function loadDashboard() {
    setStatus("Loading...");
    els.mainPanel.innerHTML = `<div class="gs-empty">Loading Studio metadata...</div>`;

    try {
      const [appsRes, modulesRes, doctypesRes] = await Promise.all([
        apiGet("/api/core/installed-apps"),
        apiGet("/api/core/modules"),
        apiGet("/api/core/doctypes")
      ]);

      state.apps = dataOf(appsRes);
      state.modules = dataOf(modulesRes);
      state.doctypes = dataOf(doctypesRes);

      els.appCount.textContent = state.apps.length;
      els.moduleCount.textContent = state.modules.length;
      els.doctypeCount.textContent = state.doctypes.length;

      renderModules();
      renderDocTypes();

      const defaultDocType = state.doctypes.find((dt) => dt.name === "DocType") || state.doctypes[0];
      if (defaultDocType) {
        await loadDocType(defaultDocType.name);
      } else {
        setStatus("No DocTypes", "error");
        els.mainPanel.innerHTML = `<div class="gs-empty">No DocTypes found.</div>`;
      }
    } catch (err) {
      console.error(err);
      setStatus("Load failed", "error");
      els.mainPanel.innerHTML = `<div class="gs-empty">${escapeHtml(err.message)}</div>`;
    }
  }

  function renderModules() {
    const modules = [{ name: "All", module_name: "All Modules" }, ...state.modules];
    els.moduleNav.innerHTML = modules.map((module) => {
      const name = module.name || module.module_name;
      const label = moduleLabel(module);
      const active = state.selectedModule === name ? " active" : "";
      return `<button class="gs-nav-item${active}" type="button" data-module="${escapeHtml(name)}">${escapeHtml(label)}</button>`;
    }).join("");

    els.moduleNav.querySelectorAll("[data-module]").forEach((button) => {
      button.addEventListener("click", () => {
        state.selectedModule = button.dataset.module || "All";
        renderModules();
        renderDocTypes();
      });
    });
  }

  function filteredDocTypes() {
    if (state.selectedModule === "All") {
      return state.doctypes;
    }

    const selected = state.modules.find((module) => module.name === state.selectedModule);
    const names = new Set([state.selectedModule, moduleLabel(selected)].filter(Boolean));
    return state.doctypes.filter((dt) => names.has(dt.module));
  }

  function renderDocTypes() {
    const items = filteredDocTypes();
    if (!items.length) {
      els.doctypeList.innerHTML = `<div class="gs-empty">No DocTypes in this module.</div>`;
      return;
    }

    els.doctypeList.innerHTML = items.map((dt) => {
      const active = dt.name === state.selectedDoctypeName ? " active" : "";
      return `
        <button class="gs-doctype-item${active}" type="button" data-doctype="${escapeHtml(dt.name)}">
          <strong>${escapeHtml(dt.label || dt.name)}</strong>
          <small>${escapeHtml(doctypeSubtitle(dt))}</small>
        </button>
      `;
    }).join("");

    els.doctypeList.querySelectorAll("[data-doctype]").forEach((button) => {
      button.addEventListener("click", () => loadDocType(button.dataset.doctype));
    });
  }

  async function loadDocType(name) {
    if (!name) {
      return;
    }

    setStatus("Loading DocType...");
    state.selectedDoctypeName = name;
    renderDocTypes();

    try {
      const encoded = encodedName(name);
      const [doctypeRes, fieldsRes, permissionsRes, resourcesRes] = await Promise.all([
        apiGet(`/api/core/doctypes/${encoded}`),
        apiGet(`/api/core/doctypes/${encoded}/fields`),
        apiGet(`/api/core/doctypes/${encoded}/permissions`),
        apiGet(`/api/resource/${encoded}?user=Administrator&limit=20`)
      ]);

      state.selectedDoctype = doctypeRes.data || {};
      state.fields = dataOf(fieldsRes);
      state.permissions = dataOf(permissionsRes);
      state.records = dataOf(resourcesRes);

      els.pageTitle.textContent = state.selectedDoctype.label || state.selectedDoctype.name || name;
      els.pageSubtitle.textContent = doctypeSubtitle(state.selectedDoctype);
      els.recordCount.textContent = state.records.length;
      els.panelTitle.textContent = state.selectedDoctype.name || name;
      els.panelSubtitle.textContent = `${state.fields.length} fields, ${state.permissions.length} permissions, ${state.records.length} records loaded`;

      renderDocTypes();
      renderActiveTab();
      setStatus("Loaded", "success");
    } catch (err) {
      console.error(err);
      setStatus("Load failed", "error");
      els.panelTitle.textContent = name;
      els.panelSubtitle.textContent = "Could not load this DocType.";
      els.mainPanel.innerHTML = `<div class="gs-empty">${escapeHtml(err.message)}</div>`;
    }
  }

  function renderActiveTab() {
    document.querySelectorAll(".gs-tab").forEach((tab) => {
      tab.classList.toggle("active", tab.dataset.tab === state.activeTab);
    });

    if (state.activeTab === "fields") {
      renderFields();
      return;
    }
    if (state.activeTab === "permissions") {
      renderPermissions();
      return;
    }
    if (state.activeTab === "json") {
      renderJson();
      return;
    }
    renderResources();
  }

  function listViewColumns() {
    const fromFields = state.fields
      .filter((field) => boolValue(field.in_list_view) && field.fieldname)
      .map((field) => field.fieldname);

    if (fromFields.length) {
      return fromFields;
    }

    const firstRecord = state.records[0] || {};
    return Object.keys(firstRecord).slice(0, 8);
  }

  function renderResources() {
    if (!state.records.length) {
      els.mainPanel.innerHTML = `<div class="gs-empty">No records found for this DocType.</div>`;
      return;
    }

    const columns = listViewColumns();
    if (!columns.length) {
      els.mainPanel.innerHTML = `<div class="gs-empty">Records loaded, but no display columns were found.</div>`;
      return;
    }

    els.mainPanel.innerHTML = `
      <div class="gs-meta-strip">
        <span class="gs-badge">${escapeHtml(state.selectedDoctypeName)}</span>
        <span class="gs-badge">${state.records.length} records</span>
      </div>
      <div class="gs-table-wrap">
        <table class="gs-table">
          <thead>
            <tr>${columns.map((column) => `<th>${escapeHtml(column)}</th>`).join("")}</tr>
          </thead>
          <tbody>
            ${state.records.map((record) => `
              <tr>
                ${columns.map((column) => `<td>${escapeHtml(formatCell(record[column]))}</td>`).join("")}
              </tr>
            `).join("")}
          </tbody>
        </table>
      </div>
    `;
  }

  function renderFields() {
    if (!state.fields.length) {
      els.mainPanel.innerHTML = `<div class="gs-empty">No fields found.</div>`;
      return;
    }

    els.mainPanel.innerHTML = `
      <div class="gs-table-wrap">
        <table class="gs-table">
          <thead>
            <tr>
              <th>Idx</th>
              <th>Fieldname</th>
              <th>Label</th>
              <th>Fieldtype</th>
              <th>Options</th>
              <th>Reqd</th>
              <th>In List View</th>
            </tr>
          </thead>
          <tbody>
            ${state.fields.map((field) => `
              <tr>
                <td>${escapeHtml(field.idx)}</td>
                <td><strong>${escapeHtml(field.fieldname)}</strong></td>
                <td>${escapeHtml(field.label)}</td>
                <td><span class="gs-badge">${escapeHtml(field.fieldtype)}</span></td>
                <td>${escapeHtml(field.options)}</td>
                <td>${boolBadge(field.reqd)}</td>
                <td>${boolBadge(field.in_list_view)}</td>
              </tr>
            `).join("")}
          </tbody>
        </table>
      </div>
    `;
  }

  function renderPermissions() {
    if (!state.permissions.length) {
      els.mainPanel.innerHTML = `<div class="gs-empty">No permissions found.</div>`;
      return;
    }

    els.mainPanel.innerHTML = `
      <div class="gs-table-wrap">
        <table class="gs-table">
          <thead>
            <tr>
              <th>Role</th>
              <th>Read</th>
              <th>Write</th>
              <th>Create</th>
              <th>Delete</th>
            </tr>
          </thead>
          <tbody>
            ${state.permissions.map((permission) => {
              const read = permission.read ?? permission.read_perm ?? false;
              const write = permission.write ?? permission.write_perm ?? false;
              const create = permission.create ?? permission.create_perm ?? false;
              const del = permission.delete ?? permission.delete_perm ?? false;
              return `
                <tr>
                  <td>${escapeHtml(permission.role)}</td>
                  <td>${boolBadge(read)}</td>
                  <td>${boolBadge(write)}</td>
                  <td>${boolBadge(create)}</td>
                  <td>${boolBadge(del)}</td>
                </tr>
              `;
            }).join("")}
          </tbody>
        </table>
      </div>
    `;
  }

  function renderJson() {
    const preview = {
      ...(state.selectedDoctype || {}),
      fields: state.fields,
      permissions: state.permissions
    };

    els.mainPanel.innerHTML = `<pre class="gs-json-preview">${escapeHtml(JSON.stringify(preview, null, 2))}</pre>`;
  }

  function formatCell(value) {
    if (value === null || value === undefined) {
      return "";
    }
    if (typeof value === "object") {
      return JSON.stringify(value);
    }
    return value;
  }

  document.querySelectorAll(".gs-tab").forEach((tab) => {
    tab.addEventListener("click", () => {
      state.activeTab = tab.dataset.tab || "resources";
      renderActiveTab();
    });
  });

  els.refreshBtn.addEventListener("click", loadDashboard);
  loadDashboard();
})();
