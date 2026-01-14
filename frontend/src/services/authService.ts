const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

type LoginResponse = {
    token: string;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
    const res = await fetch(`${API_URL}/api/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
        throw new Error("Invalid Credentials");
    }

    return res.json();
}