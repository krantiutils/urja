"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, MemberAttendanceRecord, MemberStreak } from "@/types";
import { Loader2, Flame, Trophy, CalendarCheck } from "lucide-react";
import { AttendanceCalendar } from "@/components/member/AttendanceCalendar";

const methodColors: Record<string, string> = {
  qr: "bg-blue-500/10 text-blue-400 border-blue-500/20",
  nfc: "bg-purple-500/10 text-purple-400 border-purple-500/20",
  manual: "bg-white/[0.06] text-fg-muted border-white/[0.06]",
};

export default function MemberAttendancePage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [records, setRecords] = useState<MemberAttendanceRecord[]>([]);
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [attendanceRes, streaksRes] = await Promise.all([
          api.getMyAttendance({ limit: 50 }),
          api.getMyStreaks(),
        ]);
        setRecords(attendanceRes.data ?? []);
        setStreaks(streaksRes.data ?? []);
      } catch (err) {
        console.error("Failed to load attendance:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const streak = streaks[0];
  const currentStreak = streak?.current_streak ?? 0;
  const longestStreak = streak?.longest_streak ?? 0;

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">
        {t.memberPages.attendance}
      </h1>

      {/* A month at a glance answers "am I actually turning up", which a list
          of timestamps does not. */}
      <AttendanceCalendar t={t} locale={locale} />

      {/* Streak cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-orange-500/10 flex items-center justify-center">
              <Flame className="w-5 h-5 text-orange-400" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberPages.currentStreak}
            </span>
          </div>
          <p className="text-2xl font-semibold text-fg">
            {currentStreak} {t.memberPages.days}
          </p>
        </div>

        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center">
              <Trophy className="w-5 h-5 text-yellow-400" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberPages.longestStreak}
            </span>
          </div>
          <p className="text-2xl font-semibold text-fg">
            {longestStreak} {t.memberPages.days}
          </p>
        </div>
      </div>

      {/* Check-in history */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg">
            {t.memberPages.recentCheckIns}
          </h2>
        </div>

        {records.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 gap-3">
            <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
              <CalendarCheck className="w-6 h-6 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">{t.memberPages.noCheckIns}</p>
          </div>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {records.map((record) => {
              const date = new Date(record.check_in_at);
              return (
                <div
                  key={record.id}
                  className="px-5 py-3 flex items-center justify-between"
                >
                  <div className="flex items-center gap-4">
                    <div>
                      <p className="text-sm font-medium text-fg">
                        {date.toLocaleDateString()}
                      </p>
                      <p className="text-xs text-fg-muted">
                        {date.toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </p>
                    </div>
                  </div>
                  <span
                    className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
                      methodColors[record.method] ?? methodColors.manual
                    }`}
                  >
                    {record.method}
                  </span>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
