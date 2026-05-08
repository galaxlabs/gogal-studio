async function responseError(res) {
  const body = await res.text();

  if (!body) {
    return new Error(`${res.status} ${res.statusText}`);
  }

  try {
    const parsed = JSON.parse(body);
    if (typeof parsed.error === "string") {
      return new Error(parsed.detail ? `${parsed.error}: ${parsed.detail}` : parsed.error);
    }
    if (parsed.error && typeof parsed.error === "object") {
      const message = parsed.error.message || parsed.error.code || body;
      return new Error(parsed.error.detail ? `${message}: ${parsed.error.detail}` : message);
    }
    return new Error(parsed.message || body);
  } catch {
    return new Error(body);
  }
}

export async function apiGet(url) {
  const res = await fetch(url);

  if (!res.ok) {
    throw await responseError(res);
  }

  return res.json();
}

export async function apiPut(url, payload) {
  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) {
    throw await responseError(res);
  }

  return res.json();
}

export async function apiPost(url, payload) {
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) {
    throw await responseError(res);
  }

  return res.json();
}

export async function apiDelete(url) {
  const res = await fetch(url, {
    method: "DELETE"
  });

  if (!res.ok) {
    throw await responseError(res);
  }

  return res.json();
}
