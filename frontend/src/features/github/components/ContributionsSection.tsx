import { useState } from "react";
import { useContributions } from "../hooks/useContributions";
import { ContributionHeatMap } from "./ContributionHeatMap";
import { Button } from "@/components/ui/button";

const USERNAME = "FadhelRaihan";

const YEARS = [2026, 2025, 2024]; // Tambahkan sesuai kebutuhan

export function ContributionSection() {
  const [year, setYear] = useState<number | undefined>(undefined); // unidifined = last 12 months
  const { days, loading, error } = useContributions(USERNAME, year);

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold">GitHub Contributions</h2>
          <p className="text-sm text-slate-500">
            Aktivitas kontribusi dalam bentuk heatmap.
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            size="sm"
            variant={year === undefined ? "default" : "outline"}
            onClick={() => setYear(undefined)}
          >
            Last 12 months
          </Button>
          {YEARS.map((y) => (
            <Button
              key={y}
              size="sm"
              variant={year === y ? "default" : "outline"}
              onClick={() => setYear(y)}
            >
              {y}
            </Button>
          ))}
        </div>
      </div>

      {loading && <p>Loading contributions...</p>}
      {error && <p className="text-sm text-red-500">{error}</p>}
      
      {!loading && !error && <ContributionHeatMap days={days} year={year} />}
    </section>
  );
}
