import { api } from "../../../services/apiService";
import type { User } from "../../..//types/user";

export const usersApi = {
  list() {
    return api.get<User[]>("/users");
  },
  create(payload: { email: string; full_name: string; password: string }) {
    return api.post<User>("/users", payload);
  },
  patch(
    id: string,
    payload: Partial<{ email: string; full_name: string; password: string }>
  ) {
    return api.patch<User>(`/users/${id}`, payload);
  },
  delete(id: string) {
    return api.del<void>(`/users/${id}`);
  },
};
