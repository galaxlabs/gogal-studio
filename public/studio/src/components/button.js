import { escapeHtml } from "../utils/escape.js";

export function button(label, options = {}) {
  const type = options.type || "button";
  const className = options.className || "gs-primary-btn";
  const attrs = options.attrs || {};
  const extraAttrs = Object.entries(attrs)
    .map(([key, value]) => `${escapeHtml(key)}="${escapeHtml(value)}"`)
    .join(" ");

  return `<button class="${escapeHtml(className)}" type="${escapeHtml(type)}" ${extraAttrs}>${escapeHtml(label)}</button>`;
}
