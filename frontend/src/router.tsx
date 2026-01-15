import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { SignInPage } from "./pages/Auth/SignInPage";
import { UsersPage } from "./pages/Dashboard/UsersPage";
import { AdminLayout } from "./layouts/AdminLayout";
import { useAuthStore } from "./store/auth/authStore";
import { ProjectPage } from "./pages/Portofolio/ProjectPage";

function ProtectedRoute() {
  const token = useAuthStore((s) => s.token);
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return <AdminLayout />;
}

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ProjectPage />}/>
        {/* <Route path="/login" element={<SignInPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<UsersPage />} />
          <Route path="/users" element={<UsersPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} /> */}
      </Routes>
    </BrowserRouter>
  );
}
