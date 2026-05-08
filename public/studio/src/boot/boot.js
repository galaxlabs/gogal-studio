import { getDocTypes, getInstalledApps, getModules } from "../api/core.js";
import { store } from "../state/store.js";

export async function loadBootData() {
  const [appsRes, modulesRes, doctypesRes] = await Promise.all([
    getInstalledApps(),
    getModules(),
    getDocTypes()
  ]);

  store.apps = appsRes.data || [];
  store.modules = modulesRes.data || [];
  store.doctypes = doctypesRes.data || [];

  return {
    apps: store.apps,
    modules: store.modules,
    doctypes: store.doctypes
  };
}
