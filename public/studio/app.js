const doctypeList = document.getElementById("doctypeList");
const refreshBtn = document.getElementById("refreshBtn");
const statusPill = document.getElementById("statusPill");

const pageTitle = document.getElementById("pageTitle");
const pageSubtitle = document.getElementById("pageSubtitle");

const appCount = document.getElementById("appCount");
const moduleCount = document.getElementById("moduleCount");
const doctypeCount = document.getElementById("doctypeCount");
const fieldCount = document.getElementById("fieldCount");
const permCount = document.getElementById("permCount");

const detailsPanel = document.getElementById("detailsPanel");
const fieldsPanel = document.getElementById("fieldsPanel");
const permsPanel = document.getElementById("permsPanel");

let doctypes = [];

function setStatus(message, type = "") {
  statusPill.textContent = message;
  statusPill.className = "status-pill";

  if (type) {
    statusPill.classList.add(type);
  }
}

async function apiGet(url) {
  const res = await fetch(url);

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return res.json();
}

function encodeName(name) {
  return encodeURIComponent(name);
}

async function loadDashboard() {
  setStatus("Loading...");

  try {
    const [appsRes, modulesRes, doctypesRes] = await Promise.all([
      apiGet("/api/core/installed-apps"),
      apiGet("/api/core/modules"),
      apiGet("/api/core/doctypes")
    ]);

    const apps = appsRes.data || [];
    const modules = modulesRes.data || [];
    doctypes = doctypesRes.data || [];

    appCount.textContent = apps.length;
    moduleCount.textContent = modules.length;
    doctypeCount.textContent = doctypes.length;

    renderDocTypeList(doctypes);

    setStatus("Loaded", "success");
  } catch (err) {
    console.error(err);
    setStatus("Load failed", "error");
    doctypeList.innerHTML = `<div class="muted">${err.message}</div>`;
  }
}

function renderDocTypeList(items) {
  doctypeList.innerHTML = "";

  items.forEach((dt) => {
    const btn = document.createElement("button");
    btn.className = "doctype-item";
    btn.dataset.name = dt.name;

    btn.innerHTML = `
      <strong>${dt.name}</strong>
      <small>${dt.module} · ${dt.table_name}</small>
    `;

    btn.addEventListener("click", () => loadDocType(dt.name));

    doctypeList.appendChild(btn);
  });
}

async function loadDocType(name) {
  setStatus("Loading DocType...");

  try {
    const encoded = encodeName(name);

    const [dtRes, fieldsRes, permsRes] = await Promise.all([
      apiGet(`/api/core/doctypes/${encoded}`),
      apiGet(`/api/core/doctypes/${encoded}/fields`),
      apiGet(`/api/core/doctypes/${encoded}/permissions`)
    ]);

    const dt = dtRes.data;
    const fields = fieldsRes.data || [];
    const perms = permsRes.data || [];

    pageTitle.textContent = dt.name;
    pageSubtitle.textContent = `${dt.module} · ${dt.table_name}`;

    fieldCount.textContent = fields.length;
    permCount.textContent = perms.length;

    markActive(name);
    renderDetails(dt);
    renderFields(fields);
    renderPerms(perms);

    setStatus("Loaded", "success");
  } catch (err) {
    console.error(err);
    setStatus("Load failed", "error");
  }
}

function markActive(name) {
  document.querySelectorAll(".doctype-item").forEach((item) => {
    item.classList.toggle("active", item.dataset.name === name);
  });
}

function renderDetails(dt) {
  detailsPanel.innerHTML = `
    <div class="kv"><span>Name</span><strong>${dt.name}</strong></div>
    <div class="kv"><span>Module</span><strong>${dt.module}</strong></div>
    <div class="kv"><span>App</span><strong>${dt.app_name}</strong></div>
    <div class="kv"><span>Table</span><strong>${dt.table_name}</strong></div>
    <div class="kv"><span>Single</span><strong>${dt.is_single}</strong></div>
    <div class="kv"><span>Child Table</span><strong>${dt.is_child_table}</strong></div>
    <div class="kv"><span>Submittable</span><strong>${dt.is_submittable}</strong></div>
    <div class="kv"><span>Tree</span><strong>${dt.is_tree}</strong></div>
  `;
}

function renderFields(fields) {
  if (!fields.length) {
    fieldsPanel.innerHTML = `<div class="muted">No fields found.</div>`;
    return;
  }

  fieldsPanel.innerHTML = `
    <table class="table">
      <thead>
        <tr>
          <th>Idx</th>
          <th>Field</th>
          <th>Type</th>
          <th>Options</th>
          <th>Reqd</th>
        </tr>
      </thead>
      <tbody>
        ${fields.map((f) => `
          <tr>
            <td>${f.idx}</td>
            <td><strong>${f.fieldname}</strong><br><small>${f.label}</small></td>
            <td><span class="badge">${f.fieldtype}</span></td>
            <td>${f.options || ""}</td>
            <td>${f.reqd}</td>
          </tr>
        `).join("")}
      </tbody>
    </table>
  `;
}

function renderPerms(perms) {
  if (!perms.length) {
    permsPanel.innerHTML = `<div class="muted">No permissions found.</div>`;
    return;
  }

  permsPanel.innerHTML = `
    <table class="table">
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
        ${perms.map((p) => `
          <tr>
            <td>${p.role}</td>
            <td>${p.read}</td>
            <td>${p.write}</td>
            <td>${p.create ?? p.create_perm ?? false}</td>
            <td>${p.delete ?? p.delete_perm ?? false}</td>
          </tr>
        `).join("")}
      </tbody>
    </table>
  `;
}

refreshBtn.addEventListener("click", loadDashboard);

loadDashboard();