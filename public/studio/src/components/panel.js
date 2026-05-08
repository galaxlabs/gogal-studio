import { escapeHtml } from "../utils/escape.js";

export function setPanelHeader(title, subtitle = "") {
  const panelTitle = document.getElementById("panelTitle");
  const panelSubtitle = document.getElementById("panelSubtitle");

  if (panelTitle) panelTitle.textContent = title;
  if (panelSubtitle) panelSubtitle.textContent = subtitle;
}

export function panel(title, body, subtitle = "") {
  return `
    <section class="gs-panel">
      <div class="gs-panel-header">
        <div>
          <h3>${escapeHtml(title)}</h3>
          <p>${escapeHtml(subtitle)}</p>
        </div>
      </div>
      <div class="gs-panel-body">${body}</div>
    </section>
  `;
}
