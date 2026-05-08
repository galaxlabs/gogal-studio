import { h } from "vue";
import { isSystemField } from "../utils/fieldTypes.js";
import { text } from "../utils/safe.js";

const columns = ["idx", "fieldname", "label", "fieldtype", "options", "reqd", "in_list_view"];
const detailFields = ["fieldname", "label", "fieldtype", "options", "reqd", "hidden", "read_only", "in_list_view", "idx"];

export default {
  name: "FieldDetailsView",
  props: {
    fields: { type: Array, default: () => [] },
    activeField: { type: Object, default: null }
  },
  emits: ["select-field", "delete-field"],
  render() {
    const canDeleteField = this.activeField?.name && !isSystemField(this.activeField?.fieldname);

    return h("div", { class: "gs-fields-layout" }, [
      h("div", { class: "gs-fields-table-wrap" }, this.fields.length ? h("table", { class: "gs-table" }, [
        h("thead", h("tr", columns.map((column) => h("th", column)))),
        h("tbody", this.fields.map((field) => h("tr", {
          class: ["gs-field-row", { active: this.activeField?.fieldname === field.fieldname }],
          onClick: () => this.$emit("select-field", field)
        }, columns.map((column) => h("td", text(field[column]))))))
      ]) : h("div", { class: "gs-empty" }, "No fields found.")),
      h("div", { class: "gs-field-detail-wrap" }, this.activeField ? h("aside", { class: "gs-field-detail" }, [
        h("div", { class: "gs-field-detail-header" }, [
          h("div", { class: "gs-field-detail-title" }, this.activeField.label || this.activeField.fieldname || "Field"),
          h("button", {
            class: "gs-danger-btn",
            type: "button",
            disabled: !canDeleteField,
            title: canDeleteField ? "Delete this DocField" : "System fields cannot be deleted here",
            onClick: () => this.$emit("delete-field", this.activeField)
          }, "Delete Field")
        ]),
        canDeleteField ? null : h("div", { class: "gs-form-help" }, "System fields cannot be deleted from this panel."),
        h("div", { class: "gs-field-detail-grid" }, detailFields.map((key) => h("div", { class: "gs-field-detail-row" }, [
          h("div", { class: "gs-field-detail-label" }, key),
          h("div", { class: "gs-field-detail-value" }, text(this.activeField[key]))
        ])))
      ]) : h("div", { class: "gs-field-detail-empty" }, "Click a field row to view details."))
    ]);
  }
};
