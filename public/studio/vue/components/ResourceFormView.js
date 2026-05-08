import { h } from "vue";
import { getInputType, isEditableField, normalizeValueForSave } from "../utils/fieldTypes.js";
import { jsonText, text } from "../utils/safe.js";

const systemFields = [
  "name",
  "owner",
  "created_by",
  "modified_by",
  "created_at",
  "updated_at",
  "creation",
  "modified",
  "docstatus",
  "idx"
];

export default {
  name: "ResourceFormView",
  props: {
    bundle: { type: Object, required: true },
    record: { type: Object, default: null },
    saveEnabled: { type: Boolean, default: true },
    saving: { type: Boolean, default: false },
    error: { type: String, default: "" },
    isNew: { type: Boolean, default: false }
  },
  emits: ["back", "save", "delete"],
  data() {
    return { draft: {} };
  },
  watch: {
    record: {
      immediate: true,
      handler(record) {
        this.draft = { ...(record || {}) };
      }
    }
  },
  methods: {
    updateValue(field, value) {
      this.draft = {
        ...this.draft,
        [field.fieldname]: value
      };
    },
    save() {
      const payload = {};
      for (const field of this.bundle.fields || []) {
        if (!isEditableField(field)) continue;
        payload[field.fieldname] = normalizeValueForSave(field, this.draft[field.fieldname]);
      }
      this.$emit("save", payload);
    },
    resetDraft() {
      this.draft = { ...(this.record || {}) };
    },
    renderInput(field) {
      const value = this.draft[field.fieldname];
      const type = getInputType(field.fieldtype);

      if (field.fieldtype === "Table") {
        return h("div", { class: "gs-form-help" }, "Child table editing is not implemented yet.");
      }

      if (field.fieldtype === "Long Text" || field.fieldtype === "Small Text" || field.fieldtype === "Text" || field.fieldtype === "JSON" || field.fieldtype === "Code") {
        return h("textarea", {
          class: "gs-form-textarea",
          value: typeof value === "object" ? jsonText(value) : text(value),
          onInput: (event) => this.updateValue(field, event.target.value)
        });
      }

      if (field.fieldtype === "Select" && field.options) {
        const options = text(field.options).split("\n").map((item) => item.trim()).filter(Boolean);
        return h("select", {
          class: "gs-form-input",
          value: text(value),
          onChange: (event) => this.updateValue(field, event.target.value)
        }, [
          h("option", { value: "" }, ""),
          ...options.map((option) => h("option", { value: option }, option))
        ]);
      }

      if (type === "checkbox") {
        return h("label", { class: "gs-form-checkbox" }, [
          h("input", {
            type: "checkbox",
            checked: Boolean(value),
            onChange: (event) => this.updateValue(field, event.target.checked)
          }),
          h("span", Boolean(value) ? "Checked" : "Unchecked")
        ]);
      }

      return h("input", {
        class: "gs-form-input",
        type,
        value: text(value),
        onInput: (event) => this.updateValue(field, event.target.value)
      });
    }
  },
  render() {
    const fields = (this.bundle.fields || []).filter(isEditableField);
    const presentSystemFields = systemFields.filter((fieldname) => Object.prototype.hasOwnProperty.call(this.record || {}, fieldname));

    return h("div", { class: "gs-form-edit" }, [
      h("div", { class: "gs-form-header" }, [
        h("button", { class: "gs-back-btn", type: "button", onClick: () => this.$emit("back") }, "Back to List"),
        h("div", [
          h("h3", { class: "gs-form-title" }, this.isNew ? `New ${this.bundle.doctype?.name || "Record"}` : this.record?.name || "Record"),
          h("p", { class: "gs-form-subtitle" }, this.isNew ? "Create mode" : this.bundle.doctype?.name || "DocType")
        ])
      ]),
      this.error ? h("div", { class: "gs-error-box" }, this.error) : null,
      fields.length ? h("section", { class: "gs-form-section" }, [
        h("div", { class: "gs-form-section-title" }, "Editable Fields"),
        h("div", { class: "gs-form-grid" }, fields.map((field) => h("div", { class: "gs-form-field" }, [
          h("label", { class: "gs-form-label" }, [
            field.label || field.fieldname,
            field.reqd ? h("span", { class: "gs-required" }, "*") : null
          ]),
          this.renderInput(field)
        ])))
      ]) : h("div", { class: "gs-empty" }, "No editable fields found for this DocType."),
      h("div", { class: "gs-form-actions" }, [
        h("button", {
          class: "gs-secondary-btn",
          type: "button",
          disabled: this.saving,
          onClick: this.resetDraft
        }, "Reset"),
        h("button", {
          id: "saveRecordBtn",
          class: "gs-save-btn",
          type: "button",
          disabled: !this.saveEnabled || this.saving,
          onClick: this.save
        }, this.saving ? "Saving..." : this.isNew ? "Save New Record" : "Save"),
        this.isNew ? null : h("button", {
          class: "gs-danger-btn",
          type: "button",
          disabled: this.saving || !this.record?.name,
          onClick: () => this.$emit("delete", this.record)
        }, "Delete Record")
      ]),
      presentSystemFields.length ? h("details", { class: "gs-form-section gs-form-system" }, [
        h("summary", { class: "gs-form-section-title" }, "System Fields"),
        h("div", { class: "gs-form-system-grid" }, presentSystemFields.map((fieldname) => h("div", { class: "gs-form-row readonly" }, [
          h("div", { class: "gs-form-label" }, fieldname),
          h("div", { class: "gs-form-value" }, text(this.record?.[fieldname]))
        ])))
      ]) : null
    ].filter(Boolean));
  }
};
