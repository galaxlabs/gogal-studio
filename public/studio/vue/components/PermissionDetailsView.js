import { h } from "vue";
import { boolLabel, text } from "../utils/safe.js";

const booleanKeys = ["read", "write", "create", "delete", "submit", "cancel", "amend", "print", "email", "export", "import", "share", "report"];
const detailFields = ["role", "permlevel", ...booleanKeys, "idx"];

export function permissionKey(permission, index = 0) {
  return `${permission?.role || ""}-${permission?.permlevel ?? 0}-${permission?.idx ?? index}`;
}

function normalized(permission = {}) {
  return {
    role: permission.role || "",
    permlevel: permission.permlevel ?? 0,
    read: permission.read ?? permission.read_perm ?? false,
    write: permission.write ?? permission.write_perm ?? false,
    create: permission.create ?? permission.create_perm ?? false,
    delete: permission.delete ?? permission.delete_perm ?? false,
    submit: permission.submit ?? permission.submit_perm ?? false,
    cancel: permission.cancel ?? permission.cancel_perm ?? false,
    amend: permission.amend ?? permission.amend_perm ?? false,
    print: permission.print ?? permission.print_perm ?? false,
    email: permission.email ?? permission.email_perm ?? false,
    export: permission.export ?? permission.export_perm ?? false,
    import: permission.import ?? permission.import_perm ?? false,
    share: permission.share ?? permission.share_perm ?? false,
    report: permission.report ?? permission.report_perm ?? false,
    idx: permission.idx ?? ""
  };
}

function badge(value) {
  const isOn = Boolean(value);
  return h("span", { class: ["gs-permission-badge", isOn ? "true" : "false"] }, boolLabel(isOn));
}

export default {
  name: "PermissionDetailsView",
  props: {
    permissions: { type: Array, default: () => [] },
    activePermission: { type: Object, default: null }
  },
  emits: ["select-permission"],
  render() {
    const rows = this.permissions.map((permission, index) => ({
      raw: permission,
      key: permissionKey(permission, index),
      ...normalized(permission)
    }));
    const activeKey = this.activePermission ? permissionKey(this.activePermission, this.permissions.indexOf(this.activePermission)) : "";
    const active = this.activePermission ? normalized(this.activePermission) : null;

    return h("div", { class: "gs-permissions-layout" }, [
      h("div", { class: "gs-permissions-table-wrap" }, rows.length ? h("table", { class: "gs-table" }, [
        h("thead", h("tr", ["Role", "Level", "Read", "Write", "Create", "Delete"].map((label) => h("th", label)))),
        h("tbody", rows.map((row) => h("tr", {
          class: ["gs-permission-row", { active: row.key === activeKey }],
          onClick: () => this.$emit("select-permission", row.raw)
        }, [
          h("td", row.role),
          h("td", text(row.permlevel)),
          h("td", boolLabel(row.read)),
          h("td", boolLabel(row.write)),
          h("td", boolLabel(row.create)),
          h("td", boolLabel(row.delete))
        ])))
      ]) : h("div", { class: "gs-empty" }, "No permissions found.")),
      h("div", { class: "gs-permission-detail-wrap" }, active ? h("aside", { class: "gs-permission-detail" }, [
        h("div", { class: "gs-permission-detail-title" }, active.role || "Permission"),
        h("div", { class: "gs-permission-detail-grid" }, detailFields.map((key) => h("div", { class: "gs-permission-detail-row" }, [
          h("div", { class: "gs-permission-detail-label" }, key),
          h("div", { class: "gs-permission-detail-value" }, booleanKeys.includes(key) ? [badge(active[key])] : text(active[key]))
        ])))
      ]) : h("div", { class: "gs-permission-detail-empty" }, "Click a permission row to view details."))
    ]);
  }
};
