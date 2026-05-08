export function byId(id) {
  return document.getElementById(id);
}

export function setHtml(elOrId, html) {
  const el = typeof elOrId === "string" ? byId(elOrId) : elOrId;
  if (!el) return;
  el.innerHTML = html;
}
