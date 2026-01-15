import { useEffect, useState } from "react";
import { usersApi } from "../services/usersApi";
import type { User } from "../../../types/user";
import { useAuthStore } from "../../../store/auth/authStore";

export function useUsers() {
  const token = useAuthStore((s) => s.token);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!token) return;
    setLoading(true);
    usersApi
      .list(token)
      .then(setUsers)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [token]);

  return { users, setUsers, loading };
}
