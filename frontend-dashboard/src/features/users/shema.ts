import { z } from "zod";

export const userCreateSchema = z.object({
  email: z
    .string()
    .min(1, { message: "Email wajib diisi" })
    .email({ message: "Format email tidak valid" }),
  full_name: z.string().min(2, { message: "Nama minimal 2 karakter" }),
  password: z.string().min(6, { message: "Password minimal 6 karakter" }),
});

export type UserCreateSchema = z.infer<typeof userCreateSchema>;

export const userEditSchema = userCreateSchema.partial().extend({
  email: z
    .string()
    .min(1, { message: "Email wajib diisi" })
    .email({ message: "Format email tidak valid" })
    .optional(),
  full_name: z
    .string()
    .min(2, { message: "Nama minimal 2 karakter" })
    .optional(),
  password: z
    .string()
    .min(6, { message: "Password minimal 6 karakter" })
    .optional(),
});

export type UserEditSchema = z.infer<typeof userEditSchema>;
