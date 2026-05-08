import { h } from "vue";
import { displayCellValue, displayColumnLabel } from "../utils/display.js";
import { text } from "../utils/safe.js";

function listColumns(bundle) {
  const rows = bundle?.records || [];
  const fields = bundle?.fields || [];
  const visible = fields.filter((field) => field.in_list_view && !field.hidden);

  if (visible.length) {
    return visible.map((field) => ({
      key: field.fieldname,
      label: displayColumnLabel(field.label || field.fieldname, field.fieldname)
    }));
  }

  return Object.keys(rows[0] || {}).slice(0, 8).map((key) => ({ key, label: displayColumnLabel(key, key) }));
}

export default {
  name: "ResourceListView",
  props: {
    bundle: { type: Object, required: true },
    search: { type: String, default: "" }
  },
  emits: ["update-search", "refresh", "open-record", "new-record"],
  render() {
    const rows = this.bundle?.records || [];
    const columns = listColumns(this.bundle);
    const query = this.search.trim().toLowerCase();
    const filtered = query
      ? rows.filter((row) => columns.some((column) => text(row[column.key]).toLowerCase().includes(query)))
      : rows;
    const doctype = this.bundle?.doctype?.name || "DocType";
    const resourceError = this.bundle?.resourceError || "";

    return h("div", [
      resourceError ? h("div", { class: "gs-error-box" }, [
        h("strong", `${doctype} record table is not available.`),
        h("p", resourceError),
        h("p", "Metadata tabs still work. Create the database table before using list, edit, save, or delete records for this DocType.")
      ]) : null,
      h("div", { class: "gs-list-toolbar" }, [
        h("div", { class: "gs-list-title" }, [
          h("strong", doctype),
          h("span", resourceError ? "0 loaded records; table missing" : `${filtered.length} of ${rows.length} loaded records`)
        ]),
        h("div", { class: "gs-search-wrap" }, [
          h("input", {
            id: "listSearchInput",
            class: "gs-search-input",
            type: "search",
            placeholder: "Search loaded records...",
            value: this.search,
            onInput: (event) => this.$emit("update-search", event.target.value)
          })
        ]),
        h("div", { class: "gs-list-actions" }, [
          h("button", { class: "gs-secondary-btn", type: "button", onClick: () => this.$emit("refresh") }, "Refresh"),
          h("button", { class: "gs-secondary-btn", type: "button", disabled: Boolean(resourceError), onClick: () => this.$emit("new-record") }, "New"),
          h("button", { class: "gs-disabled-btn", type: "button", disabled: true }, "Actions")
        ])
      ]),
      !resourceError && filtered.length ? h("div", { class: "gs-table-wrap" }, [
        h("table", { class: "gs-list-table" }, [
          h("thead", h("tr", [
            h("th", { class: "gs-row-checkbox" }, ""),
            ...columns.map((column) => h("th", column.label))
          ])),
          h("tbody", filtered.map((row) => h("tr", {
            class: "gs-row-clickable",
            onClick: () => this.$emit("open-record", row)
          }, [
            h("td", { class: "gs-row-checkbox" }, [
              h("input", {
                class: "gs-row-select",
                type: "checkbox",
                disabled: true,
                onClick: (event) => event.stopPropagation()
              })
            ]),
            ...columns.map((column) => h("td", displayCellValue(column.key, row[column.key])))
          ])))
        ])
      ]) : resourceError ? h("div", { class: "gs-empty" }, `No ${doctype} records can be loaded until the backing table exists.`) : h("div", { class: "gs-empty" }, "No records found.")
    ].filter(Boolean));
  }
};
