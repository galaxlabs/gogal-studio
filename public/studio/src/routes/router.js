import { routes } from "./routes.js";

export function createRouter() {
  let currentPage = null;

  async function navigate(routeName, params = {}) {
    const route = routes.find((item) => item.name === routeName);

    if (!route) {
      throw new Error(`Route not found: ${routeName}`);
    }

    if (currentPage && typeof currentPage.destroy === "function") {
      currentPage.destroy();
    }

    currentPage = route.page(router, params);

    if (typeof currentPage.mount === "function") {
      await currentPage.mount();
    }
  }

  const router = {
    start() {
      window.gogalRouter = router;
      return navigate("dashboard");
    },
    navigate
  };

  return router;
}
