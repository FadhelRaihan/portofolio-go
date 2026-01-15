import { httpRequest } from "./httpService";

export const api = {
  get: <T>(path: string, token?: string | null) =>
    httpRequest<T>(path, {}, token),
  post: <T>(path: string, body: unknown, token?: string | null) =>
    httpRequest<T>(
      path,
      { method: "POST", body: JSON.stringify(body) },
      token
    ),
  patch: <T>(path: string, body: unknown, token?: string | null) =>
    httpRequest<T>(
      path,
      { method: "PATCH", body: JSON.stringify(body) },
      token
    ),
  del: <T>(path: string, token?: string | null) =>
    httpRequest<T>(path, { method: "DELETE" }, token),
};
