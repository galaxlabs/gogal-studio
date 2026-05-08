import { renderActiveView } from "./doctype.page.js";

export function ResourceListPage(router, params = {}) {
  async function mount() {
    await router.navigate("doctype", { name: params.name || "DocType" });
    renderActiveView();
  }

  return { mount };
}
