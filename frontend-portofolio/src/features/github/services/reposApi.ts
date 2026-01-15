import { api } from "@/services/apiService";

export type Repo = {
  name: string;
  description: string;
  html_url: string;
  private: boolean;
  created_at: string;
  updated_at: string;
};

export const reposApi = {
  listMyRepos() {
    return api.get<Repo[]>("/github/me/repos");
  },
};
