"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type {
  Locale,
  UserType,
  MemberProfile,
  MemberStreak,
  MemberPackage,
  AttendanceCalendar,
  LeaderboardResponse,
  LeaderboardEntry,
  DailySummary,
  NutritionGoal,
  WaterDailySummary,
  WorkoutPlan,
  HealthMetric,
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
  Droplets,
  Utensils,
  TrendingDown,
  TrendingUp,
  Activity,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const today = new Date().toISOString().split("T")[0];

function CalorieRing({
  consumed,
  goal,
  size = 160,
  strokeWidth = 14,
}: {
  consumed: number;
  goal: number;
  size?: number;
  strokeWidth?: number;
}) {
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const pct = goal > 0 ? Math.min(consumed / goal, 1) : 0;
  const offset = circumference * (1 - pct);
  const remaining = Math.max(0, goal - consumed);

  return (
    <div className="relative" style={{ width: size, height: size }}>
      <svg
        className="-rotate-90"
        width={size}
        height={size}
        viewBox={`0 0 ${size} ${size}`}
      >
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="rgba(255,255,255,0.06)"
          strokeWidth={strokeWidth}
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="#84cc16"
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          className="transition-all duration-700"
        />
      </svg>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-white">{remaining}</span>
        <span className="text-xs text-gray-400">kcal left</span>
      </div>
    </div>
  );
}

function WaterProgressBar({
  current,
  goal,
  t,
}: {
  current: number;
  goal: number;
  t: ReturnType<typeof getDictionary>;
}) {
  const pct = goal > 0 ? Math.min((current / goal) * 100, 100) : 0;
  const remaining = Math.max(0, goal - current);
  const completed = current >= goal;

  return (
    <div className="bg-gray-800 rounded-2xl p-6">
      <div className="flex items-center gap-2 mb-4">
        <Droplets className="w-5 h-5 text-cyan-400" />
        <span className="text-base font-semibold text-white">
          {t.water.title}
        </span>
      </div>
      <div className="flex items-center justify-between text-sm mb-2">
        <span className="text-gray-400">
          {t.water.todayIntake}: {(current / 1000).toFixed(1)}L
        </span>
        <span className="text-gray-400">
          {t.water.goal}: {(goal / 1000).toFixed(1)}L
        </span>
      </div>
      <div className="h-3 bg-white/[0.06] rounded-full overflow-hidden">
        <div
          className="h-full rounded-full transition-all duration-500 bg-cyan-400"
          style={{ width: `${pct}%` }}
        />
      </div>
      <p className="mt-2 text-xs text-center font-medium">
        {completed ? (
          <span className="text-cyan-400">{t.water.completed}</span>
        ) : (
          <span className="text-gray-400">
            {t.water.remaining}: {remaining}ml
          </span>
        )}
      </p>
    </div>
  );
}

function MacroPill({
  label,
  grams,
  goal,
  color,
}: {
  label: string;
  grams: number;
  goal: number;
  color: string;
}) {
  const pct = goal > 0 ? Math.min((grams / goal) * 100, 100) : 0;
  return (
    <div className="flex-1 min-w-0">
      <div className="flex items-center justify-between mb-1">
        <span className="text-xs text-gray-400">{label}</span>
        <span className="text-xs text-gray-400">
          {Math.round(grams)}/{goal}g
        </span>
      </div>
      <div className="h-1.5 bg-white/[0.06] rounded-full overflow-hidden">
        <div
          className="h-full rounded-full transition-all duration-500"
          style={{ width: `${pct}%`, backgroundColor: color }}
        />
      </div>
    </div>
  );
}

function WeightTrendCard({
  metrics,
  t,
}: {
  metrics: HealthMetric[];
  t: ReturnType<typeof getDictionary>;
}) {
  const weightMetrics = metrics
    .filter((m) => m.metric_type === "weight")
    .sort(
      (a, b) =>
        new Date(b.recorded_at).getTime() - new Date(a.recorded_at).getTime()
    );

  if (weightMetrics.length === 0) return null;

  const latest = weightMetrics[0];
  const previous = weightMetrics.length > 1 ? weightMetrics[1] : null;
  const currentWeight = Number(
    (latest.value as Record<string, unknown>).weight_kg ??
      (latest.value as Record<string, unknown>).value ??
      0
  );
  const prevWeight = previous
    ? Number(
        (previous.value as Record<string, unknown>).weight_kg ??
          (previous.value as Record<string, unknown>).value ??
          0
      )
    : null;
  const diff = prevWeight !== null ? currentWeight - prevWeight : null;

  return (
    <div className="bg-gray-800 rounded-2xl p-6">
      <div className="flex items-center gap-2 mb-3">
        <Activity className="w-5 h-5 text-lime-500" />
        <span className="text-base font-semibold text-white">
          {t.memberPages.health.weight}
        </span>
      </div>
      <div className="flex items-end gap-2">
        <span className="text-3xl font-bold text-white">
          {currentWeight.toFixed(1)}
        </span>
        <span className="text-sm text-gray-400 mb-1">kg</span>
        {diff !== null && diff !== 0 && (
          <span
            className={`flex items-center text-sm font-medium ml-auto ${
              diff < 0 ? "text-green-400" : "text-red-400"
            }`}
          >
            {diff < 0 ? (
              <TrendingDown className="w-4 h-4 mr-0.5" />
            ) : (
              <TrendingUp className="w-4 h-4 mr-0.5" />
            )}
            {Math.abs(diff).toFixed(1)}kg
          </span>
        )}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main page component
// ---------------------------------------------------------------------------

export default function MemberDashboardPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  // --- shared state ---
  const [profile, setProfile] = useState<MemberProfile | null>(null);
  const [loading, setLoading] = useState(true);

  // --- gym_member state ---
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [calendar, setCalendar] = useState<AttendanceCalendar | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardResponse | null>(
    null
  );
  const [feedbackOpen, setFeedbackOpen] = useState(false);
  const [feedbackText, setFeedbackText] = useState("");
  const [feedbackSent, setFeedbackSent] = useState(false);

  // --- nutrition / water state (gym_member optional, fitness_tracker, calorie_tracker) ---
  const [dailySummary, setDailySummary] = useState<DailySummary | null>(null);
  const [nutritionGoal, setNutritionGoal] = useState<NutritionGoal | null>(
    null
  );
  const [waterSummary, setWaterSummary] = useState<WaterDailySummary | null>(
    null
  );

  // --- fitness_tracker state ---
  const [workoutPlan, setWorkoutPlan] = useState<WorkoutPlan | null>(null);

  // --- weight trend state ---
  const [healthMetrics, setHealthMetrics] = useState<HealthMetric[]>([]);

  const userType: UserType =
    user?.user_type ?? profile?.user_type ?? "gym_member";

  useEffect(() => {
    async function load() {
      try {
        // Always fetch profile first to determine user type
        const profileRes = await api.getMyProfile().catch(() => null);
        if (profileRes) setProfile(profileRes);

        const effectiveType: UserType =
          user?.user_type ?? profileRes?.user_type ?? "gym_member";

        // Build parallel fetch list based on user type
        const fetches: Promise<void>[] = [];

        if (effectiveType === "gym_member") {
          fetches.push(
            api
              .getMyStreaks()
              .then((r) => setStreaks(r.data ?? []))
              .catch(() => {}),
            api
              .getMyPackages()
              .then((r) => setPackages(r.data ?? []))
              .catch(() => {}),
            api
              .getMyAttendanceCalendar()
              .then((r) => setCalendar(r))
              .catch(() => {}),
            api
              .getMyLeaderboard({ period: "monthly", limit: 3 })
              .then((r) => setLeaderboard(r))
              .catch(() => {}),
            // Optional nutrition summary for gym members
            api
              .getDailySummary({ date: today })
              .then((r) => setDailySummary(r))
              .catch(() => {}),
            api
              .getNutritionGoal()
              .then((r) => setNutritionGoal(r))
              .catch(() => {})
          );
        }

        if (effectiveType === "fitness_tracker") {
          fetches.push(
            api
              .getMyWorkoutPlan()
              .then((r) => setWorkoutPlan(r))
              .catch(() => {}),
            api
              .getDailySummary({ date: today })
              .then((r) => setDailySummary(r))
              .catch(() => {}),
            api
              .getNutritionGoal()
              .then((r) => setNutritionGoal(r))
              .catch(() => {}),
            api
              .getWaterSummary(today)
              .then((r) => setWaterSummary(r))
              .catch(() => {}),
            api
              .getMyHealth({ type: "weight", limit: 5 })
              .then((r) => setHealthMetrics(r.data ?? []))
              .catch(() => {})
          );
        }

        if (effectiveType === "calorie_tracker") {
          fetches.push(
            api
              .getDailySummary({ date: today })
              .then((r) => setDailySummary(r))
              .catch(() => {}),
            api
              .getNutritionGoal()
              .then((r) => setNutritionGoal(r))
              .catch(() => {}),
            api
              .getWaterSummary(today)
              .then((r) => setWaterSummary(r))
              .catch(() => {}),
            api
              .getMyHealth({ type: "weight", limit: 5 })
              .then((r) => setHealthMetrics(r.data ?? []))
              .catch(() => {})
          );
        }

        await Promise.allSettled(fetches);
      } catch (err) {
        console.error("Failed to load member dashboard:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // ------------------------------------------------------------------
  // Loading
  // ------------------------------------------------------------------

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-lime-500 animate-spin" />
      </div>
    );
  }

  // ------------------------------------------------------------------
  // Shared derived values
  // ------------------------------------------------------------------

  const displayName = profile?.name || profile?.phone || "";
  const firstName = displayName.includes(" ")
    ? displayName.split(" ")[0]
    : displayName;

  const dateLabel = new Date().toLocaleDateString(
    locale === "ne" ? "ne-NP" : "en-US",
    { weekday: "long", month: "long", day: "numeric" }
  );

  // ------------------------------------------------------------------
  // Greeting header (shared by all user types)
  // ------------------------------------------------------------------

  const GreetingHeader = (
    <div className="flex items-center gap-4">
      <div className="w-12 h-12 rounded-full bg-lime-500/20 flex items-center justify-center text-lime-500 text-xl font-bold">
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
        <h1 className="text-2xl font-semibold text-white">
          {t.memberPages.hello}, {firstName}
        </h1>
        <p className="text-sm text-gray-400">{dateLabel}</p>
      </div>
    </div>
  );

  // ------------------------------------------------------------------
  // GYM MEMBER VIEW
  // ------------------------------------------------------------------

  if (userType === "gym_member") {
    return <GymMemberView />;
  }

  // ------------------------------------------------------------------
  // FITNESS TRACKER VIEW
  // ------------------------------------------------------------------

  if (userType === "fitness_tracker") {
    return <FitnessTrackerView />;
  }

  // ------------------------------------------------------------------
  // CALORIE TRACKER VIEW
  // ------------------------------------------------------------------

  return <CalorieTrackerView />;

  // ===================================================================
  // GYM MEMBER (the existing home page + optional nutrition card)
  // ===================================================================

  function GymMemberView() {
    const orgName =
      profile?.organizations?.find((o) => o.org_id === user?.org_id)
        ?.org_name ?? "";
    const orgId =
      user?.org_id ?? profile?.organizations?.[0]?.org_id ?? "";

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
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ];
    const monthLabel = monthNames[now.getMonth()];
    const daysInMonth =
      calendar?.days_in_month ??
      new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
    const checkInDays = new Set(calendar?.check_in_days ?? []);
    const totalCheckIns = calendar?.total_check_ins ?? 0;
    const pct =
      daysInMonth > 0
        ? Math.round((totalCheckIns / daysInMonth) * 100)
        : 0;

    // Calendar grid
    const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    const firstWeekday = firstDayOfMonth.getDay();
    const todayDate = now.getDate();
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

    // Nutrition quick summary (optional for gym members)
    const hasNutritionData =
      nutritionGoal && dailySummary && dailySummary.total_calories > 0;

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
                : "w-11 h-11 bg-lime-500/10 text-lime-500"
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
                ? "text-sm font-semibold text-white"
                : "text-xs font-medium text-gray-400"
            }`}
          >
            {name}
          </span>
          <span
            className={`font-bold ${
              isFirst ? "text-lg text-yellow-400" : "text-base text-gray-400"
            }`}
          >
            {entry.value}
          </span>
          <div
            className={`mt-1 rounded-t-lg flex items-center justify-center font-bold text-sm ${
              isFirst
                ? "w-16 bg-yellow-500/20 text-yellow-400"
                : "w-13 bg-lime-500/10 text-gray-400"
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
        {GreetingHeader}

        {/* Monthly streak card */}
        <div className="bg-gray-800 border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center justify-between mb-4">
            <span className="text-lg font-semibold text-white">
              {monthLabel}
            </span>
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
              <div
                key={i}
                className="text-center text-xs text-gray-400 font-medium py-1"
              >
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
                const isToday = day === todayDate;
                const isPast = day < todayDate;
                return (
                  <div
                    key={ci}
                    className={`w-full aspect-square rounded-full flex items-center justify-center text-xs font-medium ${
                      isChecked
                        ? "bg-lime-500/20 text-lime-500"
                        : isToday
                          ? "ring-2 ring-lime-500 text-lime-500 font-bold"
                          : isPast
                            ? "bg-red-500/10 text-red-400/60"
                            : "text-gray-400"
                    }`}
                  >
                    {isChecked ? <Check className="w-3.5 h-3.5" /> : day}
                  </div>
                );
              })}
            </div>
          ))}

          {/* Progress bar */}
          <div className="mt-4">
            <div className="h-1.5 bg-white/[0.06] rounded-full overflow-hidden">
              <div
                className="h-full bg-lime-500 rounded-full transition-all"
                style={{ width: `${pct}%` }}
              />
            </div>
            <div className="flex justify-between mt-2 text-xs">
              <span className="text-gray-400">
                {totalCheckIns} / {daysInMonth} {t.memberPages.days}
              </span>
              <span className="text-lime-500 font-semibold">{pct}%</span>
            </div>
          </div>
        </div>

        {/* Two-column: Leaderboard + Subscription */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Leaderboard podium */}
          {rankings.length > 0 && (
            <div className="bg-gray-800 border border-white/[0.06] rounded-2xl p-5">
              <div className="flex items-center gap-2 mb-5">
                <Trophy className="w-4 h-4 text-yellow-400" />
                <span className="text-sm font-semibold text-white">
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
          <div className="bg-gray-800 border border-white/[0.06] rounded-2xl p-5">
            <div className="flex items-center justify-between mb-4">
              <span className="text-sm text-gray-400">
                {t.memberPages.mySubscription}
              </span>
              {activePackage && (
                <span className="px-2 py-0.5 rounded-full bg-lime-500/20 text-lime-500 text-[10px] font-bold">
                  {t.memberPages.subscriptionActive}
                </span>
              )}
            </div>
            {activePackage ? (
              <div className="flex items-center gap-5">
                {/* Circular progress */}
                <div className="relative w-20 h-20 flex-shrink-0">
                  <svg
                    className="w-20 h-20 -rotate-90"
                    viewBox="0 0 80 80"
                  >
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
                        daysRemaining <= 7 ? "text-yellow-400" : "text-lime-500"
                      }`}
                    >
                      {daysRemaining}
                    </span>
                    <span className="text-[10px] text-gray-400">
                      {t.memberPages.daysLeft}
                    </span>
                  </div>
                </div>
                <div>
                  <p className="text-base font-semibold text-white">
                    {activePackage.package_name}
                  </p>
                  {activePackage.amount_paid && (
                    <p className="text-sm text-gray-400 mt-1">
                      NPR {activePackage.amount_paid}
                    </p>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-400">
                {t.memberPages.noActivePackage}
              </p>
            )}
          </div>
        </div>

        {/* Quick nutrition summary card (optional, shown only if user has data) */}
        {hasNutritionData && nutritionGoal && dailySummary && (
          <div className="bg-gray-800 border border-white/[0.06] rounded-2xl p-5">
            <div className="flex items-center gap-2 mb-4">
              <Utensils className="w-4 h-4 text-lime-500" />
              <span className="text-sm font-semibold text-white">
                {t.memberPages.nutrition.dailyIntake}
              </span>
            </div>
            <div className="flex items-center gap-6">
              <CalorieRing
                consumed={dailySummary.total_calories}
                goal={nutritionGoal.calorie_goal}
                size={80}
                strokeWidth={8}
              />
              <div className="flex-1 space-y-2">
                <MacroPill
                  label={t.memberPages.nutrition.protein}
                  grams={dailySummary.total_protein}
                  goal={nutritionGoal.protein_goal_g}
                  color="#84cc16"
                />
                <MacroPill
                  label={t.memberPages.nutrition.carbs}
                  grams={dailySummary.total_carbs}
                  goal={nutritionGoal.carbs_goal_g}
                  color="#facc15"
                />
                <MacroPill
                  label={t.memberPages.nutrition.fat}
                  grams={dailySummary.total_fat}
                  goal={nutritionGoal.fat_goal_g}
                  color="#f97316"
                />
              </div>
            </div>
          </div>
        )}

        {/* Gym info card */}
        {orgName && (
          <div className="bg-gray-800 border border-white/[0.06] rounded-2xl p-5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-lime-500/10 flex items-center justify-center">
                <Dumbbell className="w-6 h-6 text-lime-500" />
              </div>
              <div>
                <p className="text-xs text-gray-400">
                  {t.memberPages.myHomeClub}
                </p>
                <p className="text-lg font-semibold text-white">{orgName}</p>
              </div>
            </div>
          </div>
        )}

        {/* Feedback button */}
        <button
          onClick={() => setFeedbackOpen(true)}
          className="w-full py-3 rounded-xl border border-lime-500/30 text-lime-500 text-sm font-medium flex items-center justify-center gap-2 hover:bg-lime-500/5 transition-colors"
        >
          <MessageSquare className="w-4 h-4" />
          {t.memberPages.sendFeedback}
        </button>

        {/* Feedback modal */}
        {feedbackOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 border border-white/[0.08] rounded-2xl p-6 w-full max-w-md">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-white">
                  {t.memberPages.sendFeedback}
                </h3>
                <button
                  onClick={() => {
                    setFeedbackOpen(false);
                    setFeedbackSent(false);
                  }}
                  className="text-gray-400 hover:text-white"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>
              {feedbackSent ? (
                <p className="text-lime-500 text-center py-8">
                  {t.memberPages.feedbackSent}
                </p>
              ) : (
                <>
                  <textarea
                    value={feedbackText}
                    onChange={(e) => setFeedbackText(e.target.value)}
                    placeholder={t.memberPages.feedbackHint}
                    rows={4}
                    className="w-full bg-gray-900 border border-white/[0.06] rounded-xl p-3 text-white text-sm placeholder:text-gray-400 resize-none focus:outline-none focus:ring-1 focus:ring-lime-500"
                  />
                  <button
                    onClick={handleFeedbackSubmit}
                    disabled={!feedbackText.trim()}
                    className="mt-3 w-full py-2.5 rounded-xl bg-lime-500 text-gray-900 font-medium text-sm disabled:opacity-40 hover:bg-lime-400 transition-colors"
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

  // ===================================================================
  // FITNESS TRACKER VIEW
  // ===================================================================

  function FitnessTrackerView() {
    const template = workoutPlan?.template;
    const exercises = template?.exercises ?? [];

    return (
      <div className="space-y-6">
        {GreetingHeader}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Today's Workout Plan Card */}
          <div className="bg-gray-800 rounded-2xl p-6">
            <div className="flex items-center gap-2 mb-4">
              <Dumbbell className="w-5 h-5 text-lime-500" />
              <span className="text-base font-semibold text-white">
                {t.memberPages.workout.currentPlan}
              </span>
            </div>
            {template ? (
              <>
                <p className="text-lg font-bold text-white mb-1">
                  {template.name}
                </p>
                {template.duration_minutes && (
                  <p className="text-sm text-gray-400 mb-3">
                    {template.duration_minutes} {t.memberPages.workout.minutes}
                    {template.difficulty && ` - ${template.difficulty}`}
                  </p>
                )}
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {exercises.slice(0, 5).map((ex, i) => (
                    <div
                      key={i}
                      className="flex items-center justify-between py-1.5 border-b border-white/[0.04] last:border-0"
                    >
                      <span className="text-sm text-white">{ex.name}</span>
                      <span className="text-xs text-gray-400">
                        {ex.sets} {t.memberPages.workout.sets} x {ex.reps}{" "}
                        {t.memberPages.workout.reps}
                      </span>
                    </div>
                  ))}
                  {exercises.length > 5 && (
                    <p className="text-xs text-gray-400 text-center pt-1">
                      +{exercises.length - 5} more
                    </p>
                  )}
                </div>
              </>
            ) : (
              <p className="text-sm text-gray-400 py-4">
                {t.memberPages.workout.noPlan}
              </p>
            )}
          </div>

          {/* Nutrition Donut Chart */}
          <div className="bg-gray-800 rounded-2xl p-6">
            <div className="flex items-center gap-2 mb-4">
              <Utensils className="w-5 h-5 text-lime-500" />
              <span className="text-base font-semibold text-white">
                {t.memberPages.nutrition.dailyIntake}
              </span>
            </div>
            {nutritionGoal ? (
              <div className="flex flex-col items-center">
                <CalorieRing
                  consumed={dailySummary?.total_calories ?? 0}
                  goal={nutritionGoal.calorie_goal}
                  size={140}
                  strokeWidth={12}
                />
                <div className="flex gap-4 mt-4 w-full">
                  <MacroPill
                    label={t.memberPages.nutrition.protein}
                    grams={dailySummary?.total_protein ?? 0}
                    goal={nutritionGoal.protein_goal_g}
                    color="#84cc16"
                  />
                  <MacroPill
                    label={t.memberPages.nutrition.carbs}
                    grams={dailySummary?.total_carbs ?? 0}
                    goal={nutritionGoal.carbs_goal_g}
                    color="#facc15"
                  />
                  <MacroPill
                    label={t.memberPages.nutrition.fat}
                    grams={dailySummary?.total_fat ?? 0}
                    goal={nutritionGoal.fat_goal_g}
                    color="#f97316"
                  />
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-400 py-4 text-center">
                {t.memberPages.nutrition.noGoalSet}
              </p>
            )}
          </div>

          {/* Water Tracker */}
          <WaterProgressBar
            current={waterSummary?.total_ml ?? 0}
            goal={waterSummary?.goal_ml ?? profile?.daily_water_goal_ml ?? 2500}
            t={t}
          />

          {/* Weight Trend */}
          <WeightTrendCard metrics={healthMetrics} t={t} />
        </div>
      </div>
    );
  }

  // ===================================================================
  // CALORIE TRACKER VIEW
  // ===================================================================

  function CalorieTrackerView() {
    const consumed = dailySummary?.total_calories ?? 0;
    const goal = nutritionGoal?.calorie_goal ?? 0;
    const meals = dailySummary?.meals ?? [];

    const mealTypes = ["breakfast", "lunch", "dinner", "snack"] as const;
    const mealLabels: Record<string, string> = {
      breakfast: t.memberPages.nutrition.breakfast,
      lunch: t.memberPages.nutrition.lunch,
      dinner: t.memberPages.nutrition.dinner,
      snack: t.memberPages.nutrition.snack,
    };
    const mealColors: Record<string, string> = {
      breakfast: "#facc15",
      lunch: "#84cc16",
      dinner: "#f97316",
      snack: "#a78bfa",
    };

    return (
      <div className="space-y-6">
        {GreetingHeader}

        {/* Large Calorie Budget Ring */}
        <div className="bg-gray-800 rounded-2xl p-6 flex flex-col items-center">
          {nutritionGoal ? (
            <>
              <CalorieRing
                consumed={consumed}
                goal={goal}
                size={200}
                strokeWidth={16}
              />
              <div className="flex gap-8 mt-4">
                <div className="text-center">
                  <p className="text-lg font-bold text-white">{consumed}</p>
                  <p className="text-xs text-gray-400">
                    {t.memberPages.nutrition.caloriesConsumed}
                  </p>
                </div>
                <div className="text-center">
                  <p className="text-lg font-bold text-lime-500">{goal}</p>
                  <p className="text-xs text-gray-400">
                    {t.memberPages.nutrition.calorieGoal}
                  </p>
                </div>
              </div>
              {/* Macro bars */}
              <div className="flex gap-4 mt-4 w-full max-w-sm">
                <MacroPill
                  label={t.memberPages.nutrition.protein}
                  grams={dailySummary?.total_protein ?? 0}
                  goal={nutritionGoal.protein_goal_g}
                  color="#84cc16"
                />
                <MacroPill
                  label={t.memberPages.nutrition.carbs}
                  grams={dailySummary?.total_carbs ?? 0}
                  goal={nutritionGoal.carbs_goal_g}
                  color="#facc15"
                />
                <MacroPill
                  label={t.memberPages.nutrition.fat}
                  grams={dailySummary?.total_fat ?? 0}
                  goal={nutritionGoal.fat_goal_g}
                  color="#f97316"
                />
              </div>
            </>
          ) : (
            <p className="text-sm text-gray-400 py-8">
              {t.memberPages.nutrition.noGoalSet}
            </p>
          )}
        </div>

        {/* Meal Summaries */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {mealTypes.map((mealType) => {
            const meal = meals.find((m) => m.meal_type === mealType);
            const cal = meal?.calories ?? 0;
            const itemCount = meal?.items?.length ?? 0;
            return (
              <div
                key={mealType}
                className="bg-gray-800 rounded-2xl p-5 flex items-center gap-4"
              >
                <div
                  className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{ backgroundColor: mealColors[mealType] + "20" }}
                >
                  <Utensils
                    className="w-5 h-5"
                    style={{ color: mealColors[mealType] }}
                  />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-semibold text-white">
                    {mealLabels[mealType]}
                  </p>
                  <p className="text-xs text-gray-400">
                    {itemCount > 0
                      ? `${itemCount} item${itemCount > 1 ? "s" : ""}`
                      : t.memberPages.nutrition.noLogs}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-lg font-bold text-white">{cal}</p>
                  <p className="text-[10px] text-gray-400">
                    {t.memberPages.nutrition.kcal}
                  </p>
                </div>
              </div>
            );
          })}
        </div>

        {/* Water + Weight row */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <WaterProgressBar
            current={waterSummary?.total_ml ?? 0}
            goal={
              waterSummary?.goal_ml ?? profile?.daily_water_goal_ml ?? 2500
            }
            t={t}
          />
          <WeightTrendCard metrics={healthMetrics} t={t} />
        </div>
      </div>
    );
  }
}
