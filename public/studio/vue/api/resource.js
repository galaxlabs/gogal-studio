import { apiDelete, apiGet, apiPost, apiPut } from "./client.js";

export function getResourceList(doctype, user = "Administrator", limit = 20) {
  return apiGet(`/api/resource/${encodeURIComponent(doctype)}?user=${encodeURIComponent(user)}&limit=${limit}`);
}

export function getResourceDoc(doctype, name, user = "Administrator") {
  return apiGet(`/api/resource/${encodeURIComponent(doctype)}/${encodeURIComponent(name)}?user=${encodeURIComponent(user)}`);
}

export function createResourceDoc(doctype, payload, user = "Administrator") {
  return apiPost(`/api/resource/${encodeURIComponent(doctype)}?user=${encodeURIComponent(user)}`, payload);
}

export function updateResourceDoc(doctype, name, payload, user = "Administrator") {
  return apiPut(`/api/resource/${encodeURIComponent(doctype)}/${encodeURIComponent(name)}?user=${encodeURIComponent(user)}`, payload);
}

export function deleteResourceDoc(doctype, name, user = "Administrator") {
  return apiDelete(`/api/resource/${encodeURIComponent(doctype)}/${encodeURIComponent(name)}?user=${encodeURIComponent(user)}`);
}
