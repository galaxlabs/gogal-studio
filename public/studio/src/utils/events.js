export function on(elOrSelector, eventName, handler) {
  const el = typeof elOrSelector === "string" ? document.querySelector(elOrSelector) : elOrSelector;
  if (!el) return () => {};

  el.addEventListener(eventName, handler);
  return () => el.removeEventListener(eventName, handler);
}

export function onAll(selector, eventName, handler) {
  const cleanups = Array.from(document.querySelectorAll(selector)).map((el) => on(el, eventName, handler));
  return () => cleanups.forEach((cleanup) => cleanup());
}
