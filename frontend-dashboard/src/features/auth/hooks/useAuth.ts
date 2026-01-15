import { useState } from "react";
import { authApi } from "../services/authApi";

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await authApi.login(email, password);
      localStorage.setItem("token", res.token);
      return true;
    } catch (e) {
      setError("Login gagal");
      return false;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
  };

  return { login, logout, loading, error };
}
