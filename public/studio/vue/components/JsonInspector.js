import { h } from "vue";
import { jsonFileName, jsonText } from "../utils/safe.js";

function buildDocTypeJson(bundle) {
  return {
    ...(bundle?.doctype || {}),
    fields: bundle?.fields || [],
    permissions: bundle?.permissions || []
  };
}

export default {
  name: "JsonInspector",
  props: {
    bundle: { type: Object, required: true }
  },
  emits: ["status"],
  computed: {
    json() {
      return jsonText(buildDocTypeJson(this.bundle));
    }
  },
  methods: {
    async copyJson() {
      try {
        await navigator.clipboard.writeText(this.json);
        this.$emit("status", "JSON copied", "success");
      } catch (error) {
        console.error(error);
        this.$emit("status", "Copy failed", "error");
      }
    },
    downloadJson() {
      const blob = new Blob([this.json], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");

      link.href = url;
      link.download = jsonFileName(this.bundle?.doctype?.name);
      document.body.appendChild(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);

      this.$emit("status", "JSON downloaded", "success");
    }
  },
  render() {
    return h("div", [
      h("div", { class: "gs-json-toolbar" }, [
        h("div", { class: "gs-json-size" }, `${this.json.length} chars`),
        h("div", { class: "gs-json-toolbar-actions" }, [
          h("button", { id: "copyJsonBtn", class: "gs-json-btn", type: "button", onClick: this.copyJson }, "Copy JSON"),
          h("button", { id: "downloadJsonBtn", class: "gs-json-btn", type: "button", onClick: this.downloadJson }, "Download JSON")
        ])
      ]),
      h("pre", { id: "jsonPreviewText", class: "gs-json-preview" }, this.json)
    ]);
  }
};
