import { useState } from "react";
import { login } from "@/services/authService";
import { useAuthStore } from "@/store/authStore";

export function useLogin() {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const setAuth = useAuthStore((s) => s.setAuth);

    const handleLogin = async (email: string, password: string) => {
        setLoading(true);
        setError(null);

        try {
            const res = await login(email, password);
            setAuth(res.token, {email});
            setSuccess("Success Login")
        } catch (e: any) {
            setError("Login Failed");
        } finally {
            setLoading(false);
        }
    };

    return { handleLogin, loading, error, success };
}