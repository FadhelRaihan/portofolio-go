import { api } from "../../../services/apiService";
import type { User } from "../../..//types/user";

export const usersApi = {
  list(token: string | null) {
    return api.get<User[]>("/users", token);
  },
  create(token: string | null, payload: { email: string; full_name: string; password: string }) {
    return api.post<User>("/users", payload, token);
  },
  patch(
    token: string | null,
    id: string,
    payload: Partial<{ email: string; full_name: string; password: string }>
  ) {
    return api.patch<User>(`/users/${id}`, payload, token);
  },
  delete(token: string | null, id: string) {
    return api.del<void>(`/users/${id}`, token);
  },
};
