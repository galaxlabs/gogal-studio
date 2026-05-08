import { loadBootData } from "../boot/boot.js";
import { mountDeskSummary } from "../boot/desk.js";
import { bindDocTypeSearch, renderDocTypes, renderModules } from "../components/sidebar.js";
import { setStatus } from "../components/status.js";
import { store } from "../state/store.js";

export function getVisibleDocTypes() {
  const searchText = (store.doctypeSearch || "").trim().toLowerCase();
  let doctypes = store.doctypes || [];

  if (store.activeModule && store.activeModule !== "All") {
    doctypes = doctypes.filter((dt) => dt.module === store.activeModule);
  }

  if (searchText) {
    doctypes = doctypes.filter((dt) => {
      const haystack = `${dt.name || ""} ${dt.module || ""} ${dt.table_name || ""}`.toLowerCase();
      return haystack.includes(searchText);
    });
  }

  return doctypes;
}

export function DashboardPage(router) {
  async function mount() {
    setStatus("Loading...");

    await loadBootData();
    store.activeModule = store.activeModule || "All";
    store.doctypeSearch = store.doctypeSearch || "";

    mountDeskSummary(store);

    const renderVisibleDocTypes = () => {
      renderDocTypes(getVisibleDocTypes(), (doctypeName) => {
        router.navigate("doctype", { name: doctypeName });
      });
    };

    const handleModuleClick = (moduleName) => {
      store.activeModule = moduleName || "All";
      renderModules(store.modules, handleModuleClick);
      renderVisibleDocTypes();
    };

    renderModules(store.modules, handleModuleClick);
    bindDocTypeSearch((searchText) => {
      store.doctypeSearch = searchText;
      renderVisibleDocTypes();
    });

    renderVisibleDocTypes();

    setStatus("Loaded", "success");

    await router.navigate("doctype", { name: store.activeDocType || "DocType" });
  }

  return { mount };
}
