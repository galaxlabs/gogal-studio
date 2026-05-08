export function text(value) {
  return value === null || value === undefined ? "" : String(value);
}

export function boolLabel(value) {
  return Boolean(value) ? "true" : "false";
}

export function jsonText(value) {
  return JSON.stringify(value, null, 2);
}

export function jsonFileName(name) {
  return `${text(name || "doctype").trim().toLowerCase().replace(/\s+/g, "_")}.json`;
}
