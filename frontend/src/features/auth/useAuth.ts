import { useState } from "react";
import { authApi } from "./services/authApi";
import { useAuthStore } from "../../store/auth/authStore";

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const setAuth = useAuthStore((s) => s.setAuth);

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await authApi.login(email, password);
      setAuth(res.token, res.user);
    } catch (e: any) {
      setError("Login gagal");
      throw e;
    } finally {
      setLoading(false);
    }
  };

  return { login, loading, error };
}
