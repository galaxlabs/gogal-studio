import { escapeHtml } from "../utils/escape.js";

export function badge(value, tone = "") {
  const toneClass = tone ? ` ${escapeHtml(tone)}` : "";
  return `<span class="gs-badge${toneClass}">${escapeHtml(value)}</span>`;
}

export function boolBadge(value) {
  const on = value === true || value === 1 || value === "1" || value === "true";
  return badge(on ? "Yes" : "No", on ? "on" : "off");
}
