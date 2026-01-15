const BASE_URL = import.meta.env.VITE_API_URL;

async function request<T>(input: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${input}`, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
    ...init,
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }

  return res.json() as Promise<T>;
}

export const api = {
  get<T>(path: string) {
    return request<T>(path);
  },
};
