"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type {
  Locale,
  MemberProfile,
  MemberStreak,
  MemberPackage,
  AttendanceCalendar,
  LeaderboardResponse,
  LeaderboardEntry,
} from "@/types";
import {
  Loader2,
  Flame,
  Trophy,
  Crown,
  Dumbbell,
  MessageSquare,
  X,
  Check,
} from "lucide-react";

export default function MemberDashboardPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [profile, setProfile] = useState<MemberProfile | null>(null);
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [calendar, setCalendar] = useState<AttendanceCalendar | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [feedbackOpen, setFeedbackOpen] = useState(false);
  const [feedbackText, setFeedbackText] = useState("");
  const [feedbackSent, setFeedbackSent] = useState(false);

  useEffect(() => {
    async function load() {
      try {
        const [profileRes, streaksRes, packagesRes, calendarRes, leaderboardRes] =
          await Promise.allSettled([
            api.getMyProfile(),
            api.getMyStreaks(),
            api.getMyPackages(),
            api.getMyAttendanceCalendar(),
            api.getMyLeaderboard({ period: "monthly", limit: 3 }),
          ]);
        if (profileRes.status === "fulfilled") setProfile(profileRes.value);
        if (streaksRes.status === "fulfilled") setStreaks(streaksRes.value.data ?? []);
        if (packagesRes.status === "fulfilled") setPackages(packagesRes.value.data ?? []);
        if (calendarRes.status === "fulfilled") setCalendar(calendarRes.value);
        if (leaderboardRes.status === "fulfilled") setLeaderboard(leaderboardRes.value);
      } catch (err) {
        console.error("Failed to load member dashboard:", err);
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

  const displayName = profile?.name || profile?.phone || "";
  const firstName = displayName.includes(" ")
    ? displayName.split(" ")[0]
    : displayName;
  const orgName =
    profile?.organizations?.find((o) => o.org_id === user?.org_id)?.org_name ?? "";
  const orgId = user?.org_id ?? profile?.organizations?.[0]?.org_id ?? "";

  const activePackage = packages.find(
    (p) => p.status?.toLowerCase() === "active"
  );
  const daysRemaining = activePackage
    ? Math.max(
        0,
        Math.ceil(
          (new Date(activePackage.end_date).getTime() - Date.now()) /
            (1000 * 60 * 60 * 24)
        )
      )
    : 0;
  const totalDays = activePackage
    ? Math.max(
        1,
        Math.ceil(
          (new Date(activePackage.end_date).getTime() -
            new Date(activePackage.start_date).getTime()) /
            (1000 * 60 * 60 * 24)
        )
      )
    : 1;

  const streak = streaks[0];
  const currentStreak = streak?.current_streak ?? 0;

  const now = new Date();
  const monthNames = [
    "January", "February", "March", "April", "May", "June",
    "July", "August", "September", "October", "November", "December",
  ];
  const monthLabel = monthNames[now.getMonth()];
  const daysInMonth = calendar?.days_in_month ?? new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
  const checkInDays = new Set(calendar?.check_in_days ?? []);
  const totalCheckIns = calendar?.total_check_ins ?? 0;
  const pct = daysInMonth > 0 ? Math.round((totalCheckIns / daysInMonth) * 100) : 0;

  // Calendar grid
  const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
  const firstWeekday = firstDayOfMonth.getDay(); // 0=Sun
  const today = now.getDate();
  const dayLabels = ["S", "M", "T", "W", "T", "F", "S"];

  const calendarRows: (number | null)[][] = [];
  let dayNum = 1;
  const totalCells = firstWeekday + daysInMonth;
  const numRows = Math.ceil(totalCells / 7);
  for (let row = 0; row < numRows; row++) {
    const cells: (number | null)[] = [];
    for (let col = 0; col < 7; col++) {
      const idx = row * 7 + col;
      if (idx < firstWeekday || dayNum > daysInMonth) {
        cells.push(null);
      } else {
        cells.push(dayNum);
        dayNum++;
      }
    }
    calendarRows.push(cells);
  }

  // Leaderboard
  const rankings = leaderboard?.rankings ?? [];
  let first: LeaderboardEntry | undefined,
    second: LeaderboardEntry | undefined,
    third: LeaderboardEntry | undefined;
  for (const r of rankings) {
    if (r.rank === 1) first = r;
    if (r.rank === 2) second = r;
    if (r.rank === 3) third = r;
  }

  async function handleFeedbackSubmit() {
    if (!feedbackText.trim() || !orgId) return;
    try {
      await api.submitFeedback(orgId, { message: feedbackText.trim() });
      setFeedbackSent(true);
      setFeedbackText("");
      setTimeout(() => {
        setFeedbackOpen(false);
        setFeedbackSent(false);
      }, 2000);
    } catch (err) {
      console.error("Feedback submit failed:", err);
    }
  }

  function PodiumEntry({
    entry,
    height,
    isFirst,
  }: {
    entry: LeaderboardEntry;
    height: string;
    isFirst: boolean;
  }) {
    const name = entry.name.includes(" ")
      ? entry.name.split(" ")[0]
      : entry.name;
    const initial = name[0]?.toUpperCase() ?? "?";
    return (
      <div className="flex flex-col items-center">
        {isFirst && <Crown className="w-5 h-5 text-yellow-400 mb-1" />}
        <div
          className={`rounded-full flex items-center justify-center font-bold ${
            isFirst
              ? "w-14 h-14 bg-yellow-500/20 text-yellow-400 text-lg"
              : "w-11 h-11 bg-accent/10 text-accent"
          }`}
        >
          {entry.avatar_url ? (
            <img
              src={entry.avatar_url}
              alt={name}
              className="w-full h-full rounded-full object-cover"
            />
          ) : (
            initial
          )}
        </div>
        <span
          className={`mt-1 truncate max-w-[72px] text-center ${
            isFirst
              ? "text-sm font-semibold text-fg"
              : "text-xs font-medium text-fg-muted"
          }`}
        >
          {name}
        </span>
        <span
          className={`font-bold ${
            isFirst ? "text-lg text-yellow-400" : "text-base text-fg-muted"
          }`}
        >
          {entry.value}
        </span>
        <div
          className={`mt-1 rounded-t-lg flex items-center justify-center font-bold text-sm ${
            isFirst
              ? "w-16 bg-yellow-500/20 text-yellow-400"
              : "w-13 bg-accent/10 text-fg-muted"
          }`}
          style={{ height }}
        >
          #{entry.rank}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Greeting header */}
      <div className="flex items-center gap-4">
        <div className="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center text-accent text-xl font-bold">
          {profile?.avatar_url ? (
            <img
              src={profile.avatar_url}
              alt={firstName}
              className="w-full h-full rounded-full object-cover"
            />
          ) : (
            firstName[0]?.toUpperCase() ?? "?"
          )}
        </div>
        <div>
          <h1 className="text-2xl font-semibold text-fg">
            {t.memberPages.hello}, {firstName}
          </h1>
          <p className="text-sm text-fg-muted">{t.memberPages.letsGetMoving}</p>
        </div>
      </div>

      {/* Monthly streak card */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
        <div className="flex items-center justify-between mb-4">
          <span className="text-lg font-semibold text-fg">{monthLabel}</span>
          <div className="flex items-center gap-1.5">
            <Flame className="w-4 h-4 text-orange-400" />
            <span className="text-sm font-semibold text-orange-400">
              {currentStreak} {t.memberPages.daysStreak}
            </span>
          </div>
        </div>

        {/* Day labels */}
        <div className="grid grid-cols-7 gap-1 mb-1">
          {dayLabels.map((d, i) => (
            <div key={i} className="text-center text-xs text-fg-muted font-medium py-1">
              {d}
            </div>
          ))}
        </div>

        {/* Calendar grid */}
        {calendarRows.map((row, ri) => (
          <div key={ri} className="grid grid-cols-7 gap-1">
            {row.map((day, ci) => {
              if (day === null)
                return <div key={ci} className="w-full aspect-square" />;
              const isChecked = checkInDays.has(day);
              const isToday = day === today;
              const isPast = day < today;
              return (
                <div
                  key={ci}
                  className={`w-full aspect-square rounded-full flex items-center justify-center text-xs font-medium ${
                    isChecked
                      ? "bg-accent/20 text-accent"
                      : isToday
                        ? "ring-2 ring-accent text-accent font-bold"
                        : isPast
                          ? "bg-red-500/10 text-red-400/60"
                          : "text-fg-muted"
                  }`}
                >
                  {isChecked ? (
                    <Check className="w-3.5 h-3.5" />
                  ) : (
                    day
                  )}
                </div>
              );
            })}
          </div>
        ))}

        {/* Progress bar */}
        <div className="mt-4">
          <div className="h-1.5 bg-white/[0.06] rounded-full overflow-hidden">
            <div
              className="h-full bg-accent rounded-full transition-all"
              style={{ width: `${pct}%` }}
            />
          </div>
          <div className="flex justify-between mt-2 text-xs">
            <span className="text-fg-muted">
              {totalCheckIns} / {daysInMonth} {t.memberPages.days}
            </span>
            <span className="text-accent font-semibold">{pct}%</span>
          </div>
        </div>
      </div>

      {/* Two-column: Leaderboard + Subscription */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Leaderboard podium */}
        {rankings.length > 0 && (
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
            <div className="flex items-center gap-2 mb-5">
              <Trophy className="w-4 h-4 text-yellow-400" />
              <span className="text-sm font-semibold text-fg">
                {t.memberPages.top3CheckIns} &bull; {monthLabel}
              </span>
            </div>
            <div className="flex items-end justify-center gap-4">
              {second && (
                <PodiumEntry entry={second} height="48px" isFirst={false} />
              )}
              {first && (
                <PodiumEntry entry={first} height="64px" isFirst={true} />
              )}
              {third && (
                <PodiumEntry entry={third} height="40px" isFirst={false} />
              )}
            </div>
          </div>
        )}

        {/* Subscription card */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center justify-between mb-4">
            <span className="text-sm text-fg-muted">
              {t.memberPages.mySubscription}
            </span>
            {activePackage && (
              <span className="px-2 py-0.5 rounded-full bg-accent/20 text-accent text-[10px] font-bold">
                {t.memberPages.subscriptionActive}
              </span>
            )}
          </div>
          {activePackage ? (
            <div className="flex items-center gap-5">
              {/* Circular progress */}
              <div className="relative w-20 h-20 flex-shrink-0">
                <svg className="w-20 h-20 -rotate-90" viewBox="0 0 80 80">
                  <circle
                    cx="40"
                    cy="40"
                    r="34"
                    fill="none"
                    stroke="rgba(255,255,255,0.06)"
                    strokeWidth="6"
                  />
                  <circle
                    cx="40"
                    cy="40"
                    r="34"
                    fill="none"
                    stroke={daysRemaining <= 7 ? "#f59e0b" : "#84cc16"}
                    strokeWidth="6"
                    strokeLinecap="round"
                    strokeDasharray={`${2 * Math.PI * 34}`}
                    strokeDashoffset={`${2 * Math.PI * 34 * (1 - (totalDays - daysRemaining) / totalDays)}`}
                  />
                </svg>
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                  <span
                    className={`text-xl font-bold ${
                      daysRemaining <= 7 ? "text-yellow-400" : "text-accent"
                    }`}
                  >
                    {daysRemaining}
                  </span>
                  <span className="text-[10px] text-fg-muted">
                    {t.memberPages.daysLeft}
                  </span>
                </div>
              </div>
              <div>
                <p className="text-base font-semibold text-fg">
                  {activePackage.package_name}
                </p>
                {activePackage.amount_paid && (
                  <p className="text-sm text-fg-muted mt-1">
                    NPR {activePackage.amount_paid}
                  </p>
                )}
              </div>
            </div>
          ) : (
            <p className="text-sm text-fg-muted">
              {t.memberPages.noActivePackage}
            </p>
          )}
        </div>
      </div>

      {/* Gym info card */}
      {orgName && (
        <div className="bg-[#1a1a2e] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center">
              <Dumbbell className="w-6 h-6 text-accent" />
            </div>
            <div>
              <p className="text-xs text-fg-muted">
                {t.memberPages.myHomeClub}
              </p>
              <p className="text-lg font-semibold text-fg">{orgName}</p>
            </div>
          </div>
        </div>
      )}

      {/* Feedback button */}
      <button
        onClick={() => setFeedbackOpen(true)}
        className="w-full py-3 rounded-xl border border-accent/30 text-accent text-sm font-medium flex items-center justify-center gap-2 hover:bg-accent/5 transition-colors"
      >
        <MessageSquare className="w-4 h-4" />
        {t.memberPages.sendFeedback}
      </button>

      {/* Feedback modal */}
      {feedbackOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
          <div className="bg-elevated border border-white/[0.08] rounded-2xl p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-fg">
                {t.memberPages.sendFeedback}
              </h3>
              <button
                onClick={() => {
                  setFeedbackOpen(false);
                  setFeedbackSent(false);
                }}
                className="text-fg-muted hover:text-fg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            {feedbackSent ? (
              <p className="text-accent text-center py-8">
                {t.memberPages.feedbackSent}
              </p>
            ) : (
              <>
                <textarea
                  value={feedbackText}
                  onChange={(e) => setFeedbackText(e.target.value)}
                  placeholder={t.memberPages.feedbackHint}
                  rows={4}
                  className="w-full bg-surface border border-white/[0.06] rounded-xl p-3 text-fg text-sm placeholder:text-fg-muted resize-none focus:outline-none focus:ring-1 focus:ring-accent"
                />
                <button
                  onClick={handleFeedbackSubmit}
                  disabled={!feedbackText.trim()}
                  className="mt-3 w-full py-2.5 rounded-xl bg-accent text-bg font-medium text-sm disabled:opacity-40 hover:bg-accent-bright transition-colors"
                >
                  {t.common.save}
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
