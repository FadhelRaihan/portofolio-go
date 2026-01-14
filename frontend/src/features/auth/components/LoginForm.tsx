import { type FormEvent, useState } from "react";
import { useLogin } from "../hooks/useLogin";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function LoginForm() {
  const { handleLogin, loading, error, success } = useLogin();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    await handleLogin(email, password);
  };

  return (
    <form onSubmit={onSubmit} className="max-w-sm mx-auto mt-24 space-y-4">
      <h1 className="text-2xl font-semibold text-center">Login</h1>
      {error && <p className="text-red-500 text-sm">{error}</p>}
      {success && <p className="text-green-500 text-sm">{success}</p>}

      <div className="space-y-1">
        <label className="block text-sm font-medium">Email</label>
        <Input
          type="email"
          className="w-full border rounded px-3 py-2"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </div>

      <div className="space-y-1">
        <label className="block text-sm font-medium">Password</label>
        <Input
          type="password"
          className="w-full border rounded px-3 py-2"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>

      <Button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-600 text-white py-2 rounded disabled:opacity-60"
      >
        {loading ? "Loading..." : "Login"}
      </Button>
    </form>
  );
}
