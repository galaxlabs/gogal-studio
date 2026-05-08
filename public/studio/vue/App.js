import { h } from "vue";
import { getDocType, getDocTypeFields, getDocTypePermissions, getDocTypes, getInstalledApps, getModules } from "./api/core.js";
import { createResourceDoc, deleteResourceDoc, getResourceDoc, getResourceList, updateResourceDoc } from "./api/resource.js";
import Breadcrumb from "./components/Breadcrumb.js";
import FieldDetailsView from "./components/FieldDetailsView.js";
import JsonInspector from "./components/JsonInspector.js";
import MetaCard from "./components/MetaCard.js";
import PermissionDetailsView from "./components/PermissionDetailsView.js";
import ResourceFormView from "./components/ResourceFormView.js";
import ResourceListView from "./components/ResourceListView.js";
import StudioShell from "./components/StudioShell.js";
import StudioSidebar from "./components/StudioSidebar.js";
import StudioTopbar from "./components/StudioTopbar.js";
import SummaryCards from "./components/SummaryCards.js";
import TabsBar from "./components/TabsBar.js";
import { displayTableName } from "./utils/display.js";

function viewTitle(view) {
  if (view === "resource") return "Resource List";
  if (view === "form") return "Form View";
  if (view === "fields") return "Fields";
  if (view === "permissions") return "Permissions";
  if (view === "json") return "JSON Preview";
  return view;
}

export default {
  name: "GogalStudioApp",
  data() {
    return {
      apps: [],
      modules: [],
      doctypes: [],
      activeModule: "All",
      doctypeSearch: "",
      activeDocType: "DocType",
      activeView: "resource",
      activeBundle: null,
      activeRecord: null,
      activeField: null,
      activePermission: null,
      listSearch: "",
      statusMessage: "Ready",
      statusType: "",
      loading: false,
      error: "",
      saveError: "",
      saving: false,
      saveEnabled: true,
      isNewRecord: false
    };
  },
  computed: {
    visibleDocTypes() {
      let rows = this.doctypes;

      if (this.activeModule !== "All") {
        rows = rows.filter((dt) => dt.module === this.activeModule);
      }

      return rows;
    },
    topbarTitle() {
      return this.activeBundle?.doctype?.label || this.activeDocType || "Gogal Studio";
    },
    topbarSubtitle() {
      const dt = this.activeBundle?.doctype;
      return dt ? `${dt.module || ""} / ${displayTableName(dt.table_name)}` : "Loading core metadata...";
    },
    panelTitle() {
      return `${this.activeDocType} ${viewTitle(this.activeView)}`;
    },
    panelSubtitle() {
      if (this.activeView === "resource") return `/api/resource/${this.activeDocType}?user=Administrator&limit=20`;
      if (this.activeView === "form") return this.activeRecord?.name || "Editable record";
      if (this.activeView === "fields") return "DocField metadata";
      if (this.activeView === "permissions") return "DocPerm metadata";
      return "Generated from current API response";
    }
  },
  async mounted() {
    await this.boot();
  },
  methods: {
    setStatus(message, type = "") {
      this.statusMessage = message;
      this.statusType = type;
    },
    async boot() {
      this.loading = true;
      this.error = "";
      this.setStatus("Loading...");

      try {
        const [appsRes, modulesRes, doctypesRes] = await Promise.all([
          getInstalledApps(),
          getModules(),
          getDocTypes()
        ]);

        this.apps = appsRes.data || [];
        this.modules = modulesRes.data || [];
        this.doctypes = doctypesRes.data || [];
        await this.loadDocType(this.activeDocType || "DocType");
        this.setStatus("Loaded", "success");
      } catch (error) {
        console.error(error);
        this.error = error.message;
        this.setStatus("Load failed", "error");
      } finally {
        this.loading = false;
      }
    },
    selectModule(moduleName) {
      this.activeModule = moduleName || "All";
    },
    selectDocType(name) {
      return this.loadDocType(name);
    },
    async loadDocType(name) {
      this.activeDocType = name || "DocType";
      this.activeRecord = null;
      this.activeField = null;
      this.activePermission = null;
      this.isNewRecord = false;
      this.listSearch = "";
      this.activeView = "resource";
      this.saveError = "";
      this.setStatus("Loading DocType...");

      const [dtRes, fieldsRes, permsRes] = await Promise.all([
        getDocType(this.activeDocType),
        getDocTypeFields(this.activeDocType),
        getDocTypePermissions(this.activeDocType)
      ]);

      let records = [];
      let resourceError = "";

      try {
        const recordsRes = await getResourceList(this.activeDocType);
        records = recordsRes.data || [];
      } catch (error) {
        console.error(error);
        resourceError = error.message;
      }

      this.activeBundle = {
        doctype: dtRes.data,
        fields: fieldsRes.data || [],
        permissions: permsRes.data || [],
        records,
        resourceError
      };
      this.setStatus(resourceError ? "Metadata loaded; table missing" : "Loaded", resourceError ? "error" : "success");
    },
    async openRecord(row) {
      this.setStatus("Loading record...");
      this.saveError = "";

      try {
        const res = await getResourceDoc(this.activeDocType, row.name);
        this.activeRecord = res.data || row;
      } catch (error) {
        console.error(error);
        this.activeRecord = row;
      }

      this.activeView = "form";
      this.isNewRecord = false;
      this.setStatus("Record loaded", "success");
    },
    openNewRecord() {
      this.activeRecord = {};
      this.isNewRecord = true;
      this.saveError = "";
      this.activeView = "form";
      this.setStatus("New record", "success");
    },
    backToList() {
      this.activeRecord = null;
      this.isNewRecord = false;
      this.activeView = "resource";
      this.saveError = "";
    },
    setView(view) {
      this.activeView = view;
      if (view !== "form") {
        this.activeRecord = null;
        this.isNewRecord = false;
      }
    },
    refreshActive() {
      return this.loadDocType(this.activeDocType);
    },
    async saveRecord(payload) {
      if (!this.isNewRecord && !this.activeRecord?.name) return;

      this.saving = true;
      this.saveError = "";
      this.setStatus(this.isNewRecord ? "Creating..." : "Saving...");
      const wasNew = this.isNewRecord;

      try {
        const res = wasNew
          ? await createResourceDoc(this.activeDocType, payload)
          : await updateResourceDoc(this.activeDocType, this.activeRecord.name, payload);
        const saved = res.data || {};
        const records = this.activeBundle.records || [];
        const index = records.findIndex((row) => row.name === saved.name);

        if (index >= 0) records.splice(index, 1, saved);
        else records.unshift(saved);
        this.activeRecord = saved;
        this.isNewRecord = false;
        this.activeBundle = {
          ...this.activeBundle,
          records
        };
        this.setStatus(wasNew ? "Created" : "Saved", "success");
      } catch (error) {
        console.error(error);
        this.saveError = error.message;
        this.setStatus("Save failed", "error");
      } finally {
        this.saving = false;
      }
    },
    async deleteRecord(record) {
      if (!record?.name) return;

      const ok = window.confirm(`Delete record "${record.name}" from ${this.activeDocType}?`);
      if (!ok) return;

      this.setStatus("Deleting record...");

      try {
        await deleteResourceDoc(this.activeDocType, record.name);
        const records = (this.activeBundle.records || []).filter((row) => row.name !== record.name);
        this.activeBundle = {
          ...this.activeBundle,
          records
        };
        this.activeRecord = null;
        this.isNewRecord = false;
        this.activeView = "resource";
        this.setStatus("Record deleted", "success");
      } catch (error) {
        console.error(error);
        this.saveError = error.message;
        this.setStatus("Delete failed", "error");
      }
    },
    async deleteField(field) {
      if (!field?.name) return;

      const label = field.label || field.fieldname || field.name;
      const ok = window.confirm(`Delete field "${label}" from ${this.activeDocType}? This removes the DocField metadata row only.`);
      if (!ok) return;

      this.setStatus("Deleting field...");

      try {
        await deleteResourceDoc("DocField", field.name);
        const fields = (this.activeBundle.fields || []).filter((item) => item.name !== field.name);

        this.activeBundle = {
          ...this.activeBundle,
          fields
        };
        this.activeField = null;
        this.setStatus("Field deleted", "success");
      } catch (error) {
        console.error(error);
        this.setStatus("Delete failed", "error");
        this.error = error.message;
      }
    },
    canDeleteDocType(doctype) {
      const protectedDocTypes = new Set([
        "DocType",
        "DocField",
        "DocPerm",
        "Module Def",
        "Installed App",
        "Installed Module",
        "User",
        "Role",
        "Has Role",
        "Naming Series"
      ]);

      return Boolean(doctype?.name) && !protectedDocTypes.has(doctype.name);
    },
    async deleteDocType(doctype) {
      if (!this.canDeleteDocType(doctype)) return;

      const ok = window.confirm(`Delete DocType "${doctype.name}"? This removes the DocType metadata record.`);
      if (!ok) return;

      this.setStatus("Deleting DocType...");

      try {
        await deleteResourceDoc("DocType", doctype.name);
        const doctypes = this.doctypes.filter((item) => item.name !== doctype.name);
        this.doctypes = doctypes;

        const next = doctypes.find((item) => item.name === "DocType") || doctypes[0];
        this.activeBundle = null;
        this.activeRecord = null;
        this.activeField = null;
        this.activePermission = null;
        this.isNewRecord = false;

        if (next?.name) {
          await this.loadDocType(next.name);
        }

        this.setStatus("DocType deleted", "success");
      } catch (error) {
        console.error(error);
        this.error = error.message;
        this.setStatus("Delete failed", "error");
      }
    },
    renderActiveViewContent() {
      if (!this.activeBundle) {
        return h("div", { class: "gs-empty" }, this.error || "Loading Studio...");
      }

      const shared = [
        h(Breadcrumb, { bundle: this.activeBundle, activeView: this.activeView, activeRecord: this.activeRecord }),
        h(MetaCard, {
          bundle: this.activeBundle,
          canDelete: this.canDeleteDocType(this.activeBundle.doctype),
          onDeleteDoctype: this.deleteDocType
        })
      ];

      if (this.activeView === "resource") {
        return [
          ...shared,
          h(ResourceListView, {
            bundle: this.activeBundle,
            search: this.listSearch,
            onUpdateSearch: (value) => { this.listSearch = value; },
            onRefresh: this.refreshActive,
            onOpenRecord: this.openRecord,
            onNewRecord: this.openNewRecord
          })
        ];
      }

      if (this.activeView === "form") {
        return [
          ...shared,
          h(ResourceFormView, {
            bundle: this.activeBundle,
            record: this.activeRecord,
            saveEnabled: this.saveEnabled,
            saving: this.saving,
            error: this.saveError,
            isNew: this.isNewRecord,
            onBack: this.backToList,
            onSave: this.saveRecord,
            onDelete: this.deleteRecord
          })
        ];
      }

      if (this.activeView === "fields") {
        return [
          ...shared,
          h(FieldDetailsView, {
            fields: this.activeBundle.fields || [],
            activeField: this.activeField,
            onSelectField: (field) => { this.activeField = field; },
            onDeleteField: this.deleteField
          })
        ];
      }

      if (this.activeView === "permissions") {
        return [
          ...shared,
          h(PermissionDetailsView, {
            permissions: this.activeBundle.permissions || [],
            activePermission: this.activePermission,
            onSelectPermission: (permission) => { this.activePermission = permission; }
          })
        ];
      }

      return [
        ...shared,
        h(JsonInspector, {
          bundle: this.activeBundle,
          onStatus: this.setStatus
        })
      ];
    }
  },
  render() {
    return h(StudioShell, null, {
      sidebar: () => h(StudioSidebar, {
        modules: this.modules,
        doctypes: this.visibleDocTypes,
        activeModule: this.activeModule,
        activeDocType: this.activeDocType,
        onSelectModule: this.selectModule,
        onSelectDoctype: this.selectDocType,
        onRefresh: this.refreshActive
      }),
      default: () => [
        h(StudioTopbar, {
          title: this.topbarTitle,
          subtitle: this.topbarSubtitle,
          statusMessage: this.statusMessage,
          statusType: this.statusType
        }),
        h(SummaryCards, {
          appsCount: this.apps.length,
          modulesCount: this.modules.length,
          doctypesCount: this.doctypes.length,
          recordCount: this.activeBundle?.records?.length || 0
        }),
        h("section", { class: "gs-panel" }, [
          h("div", { class: "gs-panel-header" }, [
            h("div", [
              h("h3", this.panelTitle),
              h("p", this.panelSubtitle)
            ])
          ]),
          h(TabsBar, {
          activeView: this.activeView,
            onChangeView: this.setView
          }),
          h("div", { id: "mainPanel", class: "gs-panel-body" }, this.renderActiveViewContent())
        ])
      ]
    });
  }
};
