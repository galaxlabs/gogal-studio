import { bindTabs } from "./components/tabs.js";
import { mountDeskShell } from "./boot/desk.js";
import { renderActiveView } from "./pages/doctype.page.js";
import { createRouter } from "./routes/router.js";

document.addEventListener("DOMContentLoaded", () => {
  const router = createRouter();

  mountDeskShell();
  bindTabs(renderActiveView);

  document.getElementById("refreshBtn")?.addEventListener("click", () => {
    router.navigate("dashboard");
  });

  router.start();
});
