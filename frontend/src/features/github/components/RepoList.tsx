import { type Repo } from "../services/reposApi";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

type Props = {
  repos: Repo[];
};

export function RepoList({ repos }: Props) {
  if (!repos.length) {
    return <p className="text-sm text-slate-500">Belum ada repository.</p>;
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      {repos.map((repo) => (
        <Card key={repo.html_url} className="flex flex-col">
          <CardHeader>
            <div className="flex items-start justify-between gap-2">
              <div>
                <CardTitle className="text-base">{repo.name}</CardTitle>
                {repo.description && (
                  <CardDescription className="mt-1 line-clamp-2">
                    {repo.description}
                  </CardDescription>
                )}
              </div>
              <Badge variant={repo.private ? "secondary" : "outline"}>
                {repo.private ? "Private" : "Public"}
              </Badge>
            </div>
          </CardHeader>
          <CardContent className="text-xs text-slate-500 space-y-1">
            <p>Dibuat: {new Date(repo.created_at).toLocaleDateString()}</p>
            <p>Update: {new Date(repo.updated_at).toLocaleDateString()}</p>
          </CardContent>
          <CardFooter className="mt-auto">
            <Button asChild variant="outline" size="sm" className="w-full">
              <a href={repo.html_url} target="_blank" rel="noreferrer">
                Lihat di GitHub
              </a>
            </Button>
          </CardFooter>
        </Card>
      ))}
    </div>
  );
}
