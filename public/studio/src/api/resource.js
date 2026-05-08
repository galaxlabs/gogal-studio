import { apiGet, apiRequest } from "./client.js";

export function getResourceList(doctype, user = "Administrator", limit = 20) {
  const encoded = encodeURIComponent(doctype);
  return apiGet(`/api/resource/${encoded}?user=${encodeURIComponent(user)}&limit=${limit}`);
}

export function updateDocument(doctype, name, updates, user = "Administrator") {
  const encodedDoctype = encodeURIComponent(doctype);
  const encodedName = encodeURIComponent(name);

  return apiRequest(`/api/core/documents/${encodedDoctype}/${encodedName}?user=${encodeURIComponent(user)}`, {
    method: "PUT",
    body: JSON.stringify(updates)
  });
}
