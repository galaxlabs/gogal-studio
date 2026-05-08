import { apiGet } from "./client.js";

export function getInstalledApps() {
  return apiGet("/api/core/installed-apps");
}

export function getModules() {
  return apiGet("/api/core/modules");
}

export function getDocTypes() {
  return apiGet("/api/core/doctypes");
}

export function getDocType(name) {
  return apiGet(`/api/core/doctypes/${encodeURIComponent(name)}`);
}

export function getDocTypeFields(name) {
  return apiGet(`/api/core/doctypes/${encodeURIComponent(name)}/fields`);
}

export function getDocTypePermissions(name) {
  return apiGet(`/api/core/doctypes/${encodeURIComponent(name)}/permissions`);
}
