import { useState } from "react";
import { useAuthStore } from "../../store/auth/authStore";
import { useUsers } from "../../features/users/hooks/useUsers";
import { usersApi } from "../../features/users/services/usersApi";
import type { User } from "../../types/user";
import { Button } from "@/components/ui/button";
import { UsersTable } from "../../features/users/components/UsersTable";
import { UserFormDialog } from "../../features/users/components/UserFormDialog";

export function UsersPage() {
  const token = useAuthStore((s) => s.token);
  const { users, setUsers, loading } = useUsers();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogMode, setDialogMode] = useState<"create" | "edit">("create");
  const [selectedUser, setSelectedUser] = useState<User | undefined>();

  const openCreate = () => {
    setSelectedUser(undefined);
    setDialogMode("create");
    setDialogOpen(true);
  };

  const openEdit = (user: User) => {
    setSelectedUser(user);
    setDialogMode("edit");
    setDialogOpen(true);
  };

  const handleSubmit = async (values: {
    email: string;
    full_name: string;
    password: string;
  }) => {
    if (!token) throw new Error("No token");

    if (dialogMode === "create") {
      const created = await usersApi.create(token, values);
      setUsers((prev) => [created, ...prev]);
    } else if (dialogMode === "edit" && selectedUser) {
      const payload: Partial<{ email: string; full_name: string; password: string }> =
        {
          email: values.email,
          full_name: values.full_name,
        };
      if (values.password) {
        payload.password = values.password;
      }
      const updated = await usersApi.patch(token, selectedUser.id, payload);
      setUsers((prev) =>
        prev.map((u) => (u.id === updated.id ? updated : u))
      );
    }
  };

  const handleDelete = async (u: User) => {
    if (!token) return;
    if (!confirm(`Hapus user ${u.email}?`)) return;
    await usersApi.delete(token, u.id);
    setUsers((prev) => prev.filter((x) => x.id !== u.id));
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">User Management</h1>
        <Button onClick={openCreate}>Tambah User</Button>
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <UsersTable
          users={users}
          onEdit={openEdit}
          onDelete={handleDelete}
        />
      )}

      <UserFormDialog
        mode={dialogMode}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        initialUser={selectedUser}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
