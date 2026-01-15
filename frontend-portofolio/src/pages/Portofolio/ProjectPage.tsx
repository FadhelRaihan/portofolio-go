import { useState } from "react";
import { useRepos } from "@/features/github/hooks/useRepos";
import { RepoList } from "@/features/github/components/RepoList";
import { Button } from "@/components/ui/button";
import { ContributionSection } from "@/features/github/components/ContributionsSection";

export function ProjectPage() {
  const [visibility, setVisibility] = useState<"all" | "public" | "private">(
    "all"
  );
  const { repos, loading, error } = useRepos(visibility);

  return (
    <div className="max-w-5xl mx-auto py-10 px-4 space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Repositories</h1>
          <p className="text-sm text-slate-500">
            Daftar repository dari GitHub kamu (public & private).
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            size="sm"
            variant={visibility === "all" ? "default" : "outline"}
            onClick={() => setVisibility("all")}
          >
            All
          </Button>
          <Button
            size="sm"
            variant={visibility === "public" ? "default" : "outline"}
            onClick={() => setVisibility("public")}
          >
            Public
          </Button>
          <Button
            size="sm"
            variant={visibility === "private" ? "default" : "outline"}
            onClick={() => setVisibility("private")}
          >
            Private
          </Button>
        </div>
      </div>

      {loading && <p>Loading repos...</p>}
      {error && <p className="text-sm text-red-500">{error}</p>}

      {!loading && !error && <RepoList repos={repos} />}

      <ContributionSection />
    </div>
  );
}
