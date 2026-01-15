import { api } from "../../../services/apiService";
import type { AuthUser } from "../../../store/auth/authStore";

export type LoginResponse = {
  token: string;
  user: AuthUser;
};

export const authApi = {
  login(email: string, password: string) {
    return api.post<LoginResponse>("/login", { email, password });
  },
};
