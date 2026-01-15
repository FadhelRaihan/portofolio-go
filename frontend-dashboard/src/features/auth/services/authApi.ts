import { api } from "@/services/apiService";

export type LoginResponse = {
  token: string;
};

export const authApi = {
  login(email: string, password: string) {
    return api.post<LoginResponse>("/login", { email, password });
    // request → POST /api/login
  },
};
