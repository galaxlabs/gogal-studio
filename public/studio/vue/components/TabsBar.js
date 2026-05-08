import { h } from "vue";

const tabs = [
  ["resource", "Resource List"],
  ["fields", "Fields"],
  ["permissions", "Permissions"],
  ["json", "JSON Preview"]
];

export default {
  name: "TabsBar",
  props: {
    activeView: { type: String, default: "resource" }
  },
  emits: ["change-view"],
  render() {
    return h("div", { class: "gs-tabs", role: "tablist", "aria-label": "DocType panels" },
      tabs.map(([view, label]) => h("button", {
        class: ["gs-tab", { active: this.activeView === view }],
        type: "button",
        onClick: () => this.$emit("change-view", view)
      }, label))
    );
  }
};
