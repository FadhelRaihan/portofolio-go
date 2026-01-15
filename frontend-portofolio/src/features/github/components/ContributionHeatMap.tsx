import type { ContributionDay } from "../services/contributionsApi";

type Props = {
  days: ContributionDay[];
  year?: number; // undefined = last 12 months
};

function getColorLevel(count: number): number {
  if (count === 0) return 0;
  if (count < 3) return 1;
  if (count < 7) return 2;
  if (count < 12) return 3;
  return 4;
}

export function ContributionHeatMap({ days, year }: Props) {
  if (!days.length) {
    return <p className="text-sm text-slate-500">Tidak ada kontribusi.</p>;
  }

  // Map tanggal -> count
  const map = new Map<string, number>();
  days.forEach((d) => {
    map.set(d.date, d.count);
  });

  // Tentukan range tanggal dari year (atau last 12 months)
  let start: Date;
  let end: Date;

  if (year) {
    start = new Date(year, 0, 1); // 1 Jan
    end = new Date(year, 11, 31); // 31 Dec
  } else {
    const today = new Date();
    end = today;
    start = new Date();
    start.setFullYear(today.getFullYear() - 1);
  }

  // Normalisasi ke awal/akhir minggu
  const normalizeStart = new Date(start);
  normalizeStart.setDate(start.getDate() - start.getDay());

  const normalizeEnd = new Date(end);
  normalizeEnd.setDate(end.getDate() + (6 - end.getDay()));

  const weeks: Date[][] = [];
  let current = new Date(normalizeStart);

  while (current <= normalizeEnd) {
    const week: Date[] = [];
    for (let i = 0; i < 7; i++) {
      week.push(new Date(current));
      current.setDate(current.getDate() + 1);
    }
    weeks.push(week);
  }

  console.log("weeks length", weeks.length); // sekarang harus ~52

  return (
    <div className="flex gap-1 overflow-x-auto">
      <div className="flex flex-col justify-between py-1 mr-1 text-[10px] text-slate-400">
        <span>Mon</span>
        <span>Wed</span>
        <span>Fri</span>
      </div>

      <div className="flex gap-[2px]">
        {weeks.map((week, wi) => (
          <div key={wi} className="flex flex-col gap-[2px]">
            {week.map((day) => {
              const iso = day.toISOString().slice(0, 10);
              const count = map.get(iso) ?? 0;
              const level = getColorLevel(count);

              const baseClass = "w-3 h-3 rounded-[3px]";
              const colorClass =
                level === 0
                  ? "bg-slate-200 dark:bg-slate-800"
                  : level === 1
                  ? "bg-emerald-200"
                  : level === 2
                  ? "bg-emerald-400"
                  : level === 3
                  ? "bg-emerald-600"
                  : "bg-emerald-800";

              return (
                <div
                  key={iso}
                  className={`${baseClass} ${colorClass}`}
                  title={`${count} contributions on ${iso}`}
                />
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}
