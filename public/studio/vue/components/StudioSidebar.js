import { h } from "vue";

function moduleName(item) {
  return item?.name || item?.module_name || "";
}

export default {
  name: "StudioSidebar",
  props: {
    modules: { type: Array, default: () => [] },
    doctypes: { type: Array, default: () => [] },
    activeModule: { type: String, default: "All" },
    activeDocType: { type: String, default: "DocType" }
  },
  emits: ["select-module", "select-doctype", "refresh"],
  render() {
    const moduleItems = [{ name: "All", module_name: "All Modules", app_name: "All apps" }, ...this.modules];

    return [
      h("div", { class: "gs-brand" }, [
        h("div", { class: "gs-brand-logo" }, "G"),
        h("div", [
          h("h1", "Gogal Studio"),
          h("p", "Gogal App Builder")
        ])
      ]),
      h("button", {
        class: "gs-primary-btn",
        type: "button",
        onClick: () => this.$emit("refresh")
      }, "Refresh"),
      h("div", { class: "gs-sidebar-section" }, [
        h("div", { class: "gs-section-title" }, "Modules"),
        h("div", { class: "gs-module-nav", "aria-label": "Modules" },
          moduleItems.map((item) => {
            const name = moduleName(item);
            return h("button", {
              class: ["gs-nav-item", { active: this.activeModule === name }],
              type: "button",
              onClick: () => this.$emit("select-module", name)
            }, [
              h("strong", item.module_name || name)
            ]);
          })
        )
      ]),
      h("div", { class: "gs-sidebar-section gs-sidebar-grow" }, [
        h("div", { class: "gs-section-title" }, "DocTypes"),
        h("div", { class: "gs-selected-doctype" }, [
          h("span", "Selected"),
          h("strong", this.activeDocType || "None")
        ]),
        h("div", { class: "gs-doctype-list", "aria-label": "DocTypes" },
          this.doctypes.length
            ? this.doctypes.map((dt) => h("button", {
              class: ["gs-doctype-item", { active: dt.name === this.activeDocType }],
              type: "button",
              onClick: () => this.$emit("select-doctype", dt.name)
            }, [
              h("strong", dt.label || dt.name)
            ]))
            : [h("div", { class: "gs-empty-sidebar" }, "No DocTypes found.")]
        )
      ])
    ];
  }
};
