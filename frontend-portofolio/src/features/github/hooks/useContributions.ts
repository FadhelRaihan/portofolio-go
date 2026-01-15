import { useEffect, useState, type SetStateAction } from "react";
import {
  contributionsApi,
  type ContributionDay,
} from "@/features/github/services/contributionsApi";

export function useContributions(username: string, year?: number) {
  const [days, setDays] = useState<ContributionDay[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!username) return;

    setLoading(true);
    setError(null);

    contributionsApi
      .list(username,year)
      .then((res: { days: SetStateAction<ContributionDay[]> }) =>
        setDays(res.days)
      )
      .catch((e) => {
        console.error(e);
        setError("Gagal mengambil data kontribusi");
      })
      .finally(() => setLoading(false));
  }, [username, year]);

  return { days, loading, error };
}
