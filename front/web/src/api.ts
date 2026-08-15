const API_BASE = import.meta.env.VITE_API_BASE || "/api/v1";

export async function api<T>(path: string, options?: RequestInit): Promise<T> {
  const request: RequestInit = { ...options };
  if (!request.method || request.method.toUpperCase() === "GET") {
    request.cache = "no-store";
  }
  const response = await fetch(`${API_BASE}${path}`, request);
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(payload.error || "Request failed");
  return payload as T;
}

export function submitForm(path: string, data: FormData) {
  return api<{ message: string; id: number }>(path, { method: "POST", body: data });
}

export function externalLink(url: string) {
  if (!url) return "#";
  return url.replace(/^https?:\/\/www\.hbfittings\.net/i, "");
}
