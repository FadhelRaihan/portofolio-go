import { Outlet, Link, useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/auth/authStore";
import { Button } from "@/components/ui/button";

export function AdminLayout() {
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const navigate = useNavigate();

  const logout = () => {
    clearAuth();
    navigate("/login", { replace: true });
  };

  return (
    <div className="min-h-screen flex">
      <aside className="w-56 bg-slate-900 text-slate-100 p-4">
        <h2 className="text-lg font-semibold mb-4">Admin</h2>
        <nav className="space-y-2">
          <Link to="/dashboard" className="block text-sm">
            Dashboard
          </Link>
          <Link to="/users" className="block text-sm">
            Users
          </Link>
        </nav>
        <Button
          size="sm"
          variant="outline"
          className="mt-6"
          onClick={logout}
        >
          Logout
        </Button>
      </aside>

      <main className="flex-1 bg-slate-50 p-6">
        <Outlet />
      </main>
    </div>
  );
}
