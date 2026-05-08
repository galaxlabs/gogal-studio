import { h } from "vue";

export default {
  name: "StudioShell",
  render() {
    return h("div", { class: "gs-shell" }, [
      h("aside", { class: "gs-sidebar" }, this.$slots.sidebar?.()),
      h("main", { class: "gs-main" }, this.$slots.default?.())
    ]);
  }
};
