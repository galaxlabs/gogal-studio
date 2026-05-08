import { h } from "vue";
import StatusPill from "./StatusPill.js";

export default {
  name: "StudioTopbar",
  props: {
    title: { type: String, default: "Gogal Studio" },
    subtitle: { type: String, default: "" },
    statusMessage: { type: String, default: "Ready" },
    statusType: { type: String, default: "" }
  },
  render() {
    return h("header", { class: "gs-topbar" }, [
      h("div", [
        h("h2", this.title),
        h("p", this.subtitle)
      ]),
      h(StatusPill, { message: this.statusMessage, type: this.statusType })
    ]);
  }
};
