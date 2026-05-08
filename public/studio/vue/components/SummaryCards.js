import { h } from "vue";

export default {
  name: "SummaryCards",
  props: {
    appsCount: { type: Number, default: 0 },
    modulesCount: { type: Number, default: 0 },
    doctypesCount: { type: Number, default: 0 },
    recordCount: { type: Number, default: 0 }
  },
  render() {
    const cards = [
      ["Apps", this.appsCount],
      ["Modules", this.modulesCount],
      ["DocTypes", this.doctypesCount],
      ["Records", this.recordCount]
    ];

    return h("section", { class: "gs-summary-grid", "aria-label": "Studio summary" },
      cards.map(([label, value]) => h("div", { class: "gs-summary-card" }, [
        h("span", label),
        h("strong", value)
      ]))
    );
  }
};
