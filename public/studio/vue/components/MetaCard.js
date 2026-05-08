import { h } from "vue";
import { displayAppName, displayTableName } from "../utils/display.js";
import { boolLabel } from "../utils/safe.js";

function boolBadge(label, value) {
  const isOn = Boolean(value);
  return h("span", { class: ["gs-meta-badge", isOn ? "success" : "muted"] }, `${label}: ${boolLabel(isOn)}`);
}

export default {
  name: "MetaCard",
  props: {
    bundle: { type: Object, default: null },
    canDelete: { type: Boolean, default: false }
  },
  emits: ["delete-doctype"],
  render() {
    const doctype = this.bundle?.doctype || {};
    const records = this.bundle?.records || [];

    return h("section", { class: "gs-meta-card", "aria-label": "DocType metadata" }, [
      h("div", [
        h("div", { class: "gs-meta-title" }, doctype.name || "DocType"),
        h("div", { class: "gs-meta-subtitle" }, `${doctype.module || ""} / ${displayTableName(doctype.table_name)}`)
      ]),
      h("div", { class: "gs-meta-side" }, [
        h("div", { class: "gs-meta-badges" }, [
          h("span", { class: "gs-meta-badge" }, `App: ${displayAppName(doctype.app_name)}`),
          h("span", { class: "gs-meta-badge" }, `Table: ${displayTableName(doctype.table_name)}`),
          boolBadge("Single", doctype.is_single),
          boolBadge("Child", doctype.is_child_table),
          boolBadge("Submittable", doctype.is_submittable),
          boolBadge("Tree", doctype.is_tree),
          h("span", { class: "gs-meta-badge" }, `Records: ${records.length}`)
        ]),
        h("button", {
          class: "gs-danger-btn",
          type: "button",
          disabled: !this.canDelete,
          title: this.canDelete ? "Delete this DocType metadata record" : "Core Studio DocTypes are protected",
          onClick: () => this.$emit("delete-doctype", doctype)
        }, "Delete DocType")
      ])
    ]);
  }
};
