import { byId } from "../utils/dom.js";

export function mountDeskSummary(store) {
  byId("appCount").textContent = store.apps.length;
  byId("moduleCount").textContent = store.modules.length;
  byId("doctypeCount").textContent = store.doctypes.length;
}

export function mountDeskShell() {
  document.body.classList.add("gs-desk-ready");
}
