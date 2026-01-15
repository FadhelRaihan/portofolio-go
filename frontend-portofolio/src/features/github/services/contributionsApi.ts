import { api } from "@/services/apiService";

export type ContributionDay = {
  date: string;
  count: number;
};

export type ContributionsResponse = {
  days: ContributionDay[];
};

export const contributionsApi = {
  list(username: string, year?: number) {
    const qp = year ? `?year=${year}` : "";
    return api.get<ContributionsResponse>(
      `/github/${username}/contributions${qp}`,
    );
  },
};
