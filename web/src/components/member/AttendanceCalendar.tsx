"use client";

import { useCallback, useEffect, useState } from "react";
import { ChevronLeft, ChevronRight, Loader2 } from "lucide-react";
import { api } from "@/lib/api";
import type { AttendanceCalendar as CalendarData, Locale } from "@/types";
import type { Dictionary } from "@/lib/i18n";

/**
 * A month grid of the member's check-ins.
 *
 * The endpoint existed with nothing calling it, and the attendance page showed
 * a list of timestamps — which answers "when did I last come" but not "am I
 * actually turning up", which is the question a member has. A filled square per
 * day answers it at a glance.
 */

function monthKey(d: Date): string {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
}

export function AttendanceCalendar({
  t,
  locale,
}: {
  t: Dictionary;
  locale: Locale;
}) {
  const [cursor, setCursor] = useState(() => new Date());
  const [data, setData] = useState<CalendarData | null>(null);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      setData(await api.getMyAttendanceCalendar({ month: monthKey(cursor) }));
    } catch {
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [cursor]);

  useEffect(() => {
    load();
  }, [load]);

  const shift = (months: number) =>
    setCursor((c) => new Date(c.getFullYear(), c.getMonth() + months, 1));

  const days = data?.days_in_month ?? 0;
  const checkedIn = new Set(data?.check_in_days ?? []);
  // Which weekday the month starts on, so the grid lines up with the headings.
  const firstWeekday = new Date(cursor.getFullYear(), cursor.getMonth(), 1).getDay();

  const monthLabel = cursor.toLocaleDateString(locale === "ne" ? "ne-NP" : "en-GB", {
    month: "long",
    year: "numeric",
  });

  return (
    <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-sm font-medium text-fg">{t.attendanceCalendar.title}</h2>
          <p className="text-xs text-fg-muted mt-0.5">
            {monthLabel}
            {data ? ` · ${data.total_check_ins} ${t.attendanceCalendar.checkIns}` : ""}
          </p>
        </div>
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => shift(-1)}
            aria-label={t.attendanceCalendar.prev}
            className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <button
            type="button"
            onClick={() => shift(1)}
            aria-label={t.attendanceCalendar.next}
            className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <Loader2 className="w-5 h-5 animate-spin text-fg-muted" />
        </div>
      ) : (
        <div className="grid grid-cols-7 gap-1.5">
          {["S", "M", "T", "W", "T", "F", "S"].map((d, i) => (
            <div
              key={i}
              className="text-center text-[10px] font-mono uppercase text-fg-muted pb-1"
            >
              {d}
            </div>
          ))}

          {Array.from({ length: firstWeekday }).map((_, i) => (
            <div key={`pad-${i}`} />
          ))}

          {Array.from({ length: days }).map((_, i) => {
            const day = i + 1;
            const came = checkedIn.has(day);
            return (
              <div
                key={day}
                title={came ? `${day}` : undefined}
                className={`aspect-square rounded-md flex items-center justify-center text-xs transition-colors ${
                  came
                    ? "bg-accent/80 text-bg-base font-medium"
                    : "bg-white/[0.04] text-fg-muted"
                }`}
              >
                {day}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
