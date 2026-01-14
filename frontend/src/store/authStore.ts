import { create } from "zustand";

type User = {
    email: string;
}

type AuthState = {
    token: string | null;
    user: User | null;
    setAuth: (token: string, user: User) => void;
    clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set: any) => ({
    token: null,
    user: null,
    setAuth: (token: any, user: any) => {
        set({ token, user });
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));
    },
    clearAuth: () => {
        set({ token: null, user: null });
        localStorage.removeItem("token");
        localStorage.removeItem("user");
    }
}))
