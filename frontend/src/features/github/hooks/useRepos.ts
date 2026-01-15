import { useEffect, useState } from "react";
import { reposApi, type Repo } from "../services/reposApi";

export function useRepos(visibility: "all" | "public" | "private" = "all") {
    const [repos, setRepos] = useState<Repo[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);

        reposApi.listMyRepos(visibility).then(setRepos).catch((e) => {
            console.error(e);
            setError("Gagal mengambil repositories")
        }).finally(() => setLoading(false));
    }, [visibility]);

    return { repos, loading, error };
}