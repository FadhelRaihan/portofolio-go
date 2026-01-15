import * as React from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {
  userCreateSchema,
  userEditSchema,
  type UserCreateSchema,
  type UserEditSchema,
} from "../shema";
import type { User } from "../../../types/user";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";

type BaseValues = {
  email: string;
  full_name: string;
  password: string;
};

type Props = {
  mode: "create" | "edit";
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialUser?: User;
  onSubmit: (values: BaseValues) => Promise<void>;
};

export function UserFormDialog({
  mode,
  open,
  onOpenChange,
  initialUser,
  onSubmit,
}: Props) {
  const title = mode === "create" ? "Tambah User" : "Edit User";

  const form = useForm<UserCreateSchema | UserEditSchema>({
    resolver: zodResolver(
      mode === "create" ? userCreateSchema : userEditSchema
    ),
    defaultValues: {
      email: initialUser?.email ?? "",
      full_name: initialUser?.full_name ?? "",
      password: "",
    },
  });

  // Reset nilai ketika dialog dibuka / user berubah
  React.useEffect(() => {
    form.reset({
      email: initialUser?.email ?? "",
      full_name: initialUser?.full_name ?? "",
      password: "",
    });
  }, [initialUser, open, form, mode]);

  const handleSubmit = async (values: any) => {
    // Normalisasi payload untuk service
    const payload: BaseValues = {
      email: values.email ?? initialUser?.email ?? "",
      full_name: values.full_name ?? initialUser?.full_name ?? "",
      password: values.password ?? "",
    };

    await onSubmit(payload);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <DialogHeader>
              <DialogTitle>{title}</DialogTitle>
            </DialogHeader>

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="[email protected]"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="full_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Full Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Nama user" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    Password{" "}
                    {mode === "edit" && (
                      <span className="text-xs text-slate-500">
                        (kosongkan jika tidak diganti)
                      </span>
                    )}
                  </FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="••••••" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="flex gap-2 justify-end">
              <DialogClose asChild>
                <Button type="button" variant="outline">
                  Batal
                </Button>
              </DialogClose>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting ? "Menyimpan..." : "Simpan"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
