import { h } from "vue";

export default {
  name: "StatusPill",
  props: {
    message: { type: String, default: "Ready" },
    type: { type: String, default: "" }
  },
  render() {
    return h("span", {
      class: ["gs-status-pill", this.type]
    }, this.message);
  }
};
