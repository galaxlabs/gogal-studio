import { text } from "./safe.js";

export function displayAppName(value) {
  const raw = text(value).trim();
  if (!raw) return "";
  if (raw === "gogal_studio") return "Gogal Studio";

  return raw
    .replace(/[_-]+/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

export function displayTableName(value) {
  const raw = text(value).trim();
  if (!raw) return "";

  return raw.startsWith("tab") ? raw.slice(3) : raw;
}

export function displayColumnLabel(label, key = "") {
  if (key === "app_name") return "App";
  if (key === "table_name") return "Table";
  if (label === "app_name") return "App";
  if (label === "table_name") return "Table";

  return text(label || key);
}

export function displayCellValue(key, value) {
  if (key === "app_name") return displayAppName(value);
  if (key === "table_name") return displayTableName(value);

  return text(value);
}
