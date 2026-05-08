import { h } from "vue";

function viewLabel(activeView, activeRecord) {
  if (activeView === "resource") return "Resource List";
  if (activeView === "form") return activeRecord?.name || "Form View";
  if (activeView === "fields") return "Fields";
  if (activeView === "permissions") return "Permissions";
  if (activeView === "json") return "JSON Preview";
  return activeView || "";
}

export default {
  name: "Breadcrumb",
  props: {
    bundle: { type: Object, default: null },
    activeView: { type: String, default: "resource" },
    activeRecord: { type: Object, default: null }
  },
  render() {
    if (!this.bundle?.doctype) return null;
    const dt = this.bundle.doctype;
    const label = viewLabel(this.activeView, this.activeRecord);
    const parts = ["Studio", dt.module || "", dt.name || ""];

    return h("div", { class: "gs-breadcrumb-wrap" }, [
      h("div", { class: "gs-breadcrumb" }, [
        ...parts.flatMap((part, index) => [
          index ? h("span", { class: "gs-breadcrumb-sep" }, "/") : null,
          h("span", part)
        ]).filter(Boolean),
        h("span", { class: "gs-breadcrumb-sep" }, "/"),
        h("strong", label)
      ]),
      h("span", { class: "gs-view-badge" }, label)
    ]);
  }
};
