import { store } from "../state/store.js";

export function bindTabs(onTabChange) {
  document.querySelectorAll(".gs-tab").forEach((tab) => {
    tab.addEventListener("click", () => {
      store.activeRecord = null;
      store.activeView = tab.dataset.view;
      onTabChange();
    });
  });
}

export function updateActiveTab() {
  document.querySelectorAll(".gs-tab").forEach((tab) => {
    tab.classList.toggle("active", tab.dataset.view === store.activeView);
  });
}
