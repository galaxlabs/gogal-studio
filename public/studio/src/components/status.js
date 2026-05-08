export function setStatus(message, type = "") {
  const statusPill = document.getElementById("statusPill");

  if (!statusPill) return;

  statusPill.textContent = message;
  statusPill.className = "gs-status-pill";

  if (type) {
    statusPill.classList.add(type);
  }
}
