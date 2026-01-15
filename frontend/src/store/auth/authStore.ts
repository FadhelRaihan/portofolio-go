import { create } from "zustand";

export type AuthUser = {
  id: string;
  email: string;
  full_name: string;
};

type AuthState = {
  token: string | null;
  user: AuthUser | null;
  setAuth: (token: string, user: AuthUser) => void;
  clearAuth: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem("token"),
  user: null,
  setAuth: (token, user) => {
    set({ token, user });
    localStorage.setItem("token", token);
  },
  clearAuth: () => {
    set({ token: null, user: null });
    localStorage.removeItem("token");
  },
}));
