import { useNavigate } from "react-router-dom";
import { LoginForm } from "../../features/auth/components/LoginForm";

export function SignInPage() {
  const navigate = useNavigate();
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50">
      <LoginForm onSuccess={() => navigate("/users", { replace: true })} />
    </div>
  );
}
