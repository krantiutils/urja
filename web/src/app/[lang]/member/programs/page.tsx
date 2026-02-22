"use client";

import { useEffect, useState, useCallback } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type {
  Locale,
  WorkoutProgram,
  ProgramDay,
  UserProgramEnrollment,
} from "@/types";
import {
  Loader2,
  Dumbbell,
  ChevronDown,
  ChevronRight,
  ChevronLeft,
  Play,
  CheckCircle2,
  Calendar,
  Zap,
  Trophy,
  Moon,
  Filter,
  Target,
  Clock,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Muscle-group color map
// ---------------------------------------------------------------------------

const MUSCLE_COLORS: Record<string, { bg: string; text: string }> = {
  chest: { bg: "bg-rose-500/15", text: "text-rose-400" },
  shoulders: { bg: "bg-orange-500/15", text: "text-orange-400" },
  biceps: { bg: "bg-amber-500/15", text: "text-amber-400" },
  triceps: { bg: "bg-yellow-500/15", text: "text-yellow-400" },
  core: { bg: "bg-emerald-500/15", text: "text-emerald-400" },
  back: { bg: "bg-sky-500/15", text: "text-sky-400" },
  quadriceps: { bg: "bg-violet-500/15", text: "text-violet-400" },
  hamstrings: { bg: "bg-purple-500/15", text: "text-purple-400" },
  glutes: { bg: "bg-fuchsia-500/15", text: "text-fuchsia-400" },
  calves: { bg: "bg-pink-500/15", text: "text-pink-400" },
  forearms: { bg: "bg-lime-500/15", text: "text-lime-400" },
  lats: { bg: "bg-cyan-500/15", text: "text-cyan-400" },
  traps: { bg: "bg-indigo-500/15", text: "text-indigo-400" },
  obliques: { bg: "bg-teal-500/15", text: "text-teal-400" },
  lower_back: { bg: "bg-blue-500/15", text: "text-blue-400" },
  hip_flexors: { bg: "bg-red-500/15", text: "text-red-400" },
};

function muscleColor(group: string) {
  return MUSCLE_COLORS[group] ?? { bg: "bg-white/10", text: "text-fg-muted" };
}

function MuscleChip({ group }: { group: string }) {
  const c = muscleColor(group);
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-medium ${c.bg} ${c.text}`}
    >
      {group.replace(/_/g, " ")}
    </span>
  );
}

// ---------------------------------------------------------------------------
// Difficulty badge
// ---------------------------------------------------------------------------

function DifficultyBadge({
  difficulty,
  t,
}: {
  difficulty: string;
  t: ReturnType<typeof getDictionary>;
}) {
  const colorMap: Record<string, string> = {
    beginner: "bg-emerald-500/15 text-emerald-400 border-emerald-500/20",
    intermediate: "bg-amber-500/15 text-amber-400 border-amber-500/20",
    advanced: "bg-rose-500/15 text-rose-400 border-rose-500/20",
  };
  const labelMap: Record<string, string> = {
    beginner: t.exercises.beginner,
    intermediate: t.exercises.intermediate,
    advanced: t.exercises.advanced,
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${colorMap[difficulty] ?? "bg-white/10 text-fg-muted border-white/10"}`}
    >
      {labelMap[difficulty] ?? difficulty}
    </span>
  );
}

// ---------------------------------------------------------------------------
// Filter select
// ---------------------------------------------------------------------------

function FilterSelect({
  value,
  onChange,
  options,
  label,
}: {
  value: string;
  onChange: (v: string) => void;
  options: { value: string; label: string }[];
  label: string;
}) {
  return (
    <div className="relative flex-1 min-w-[140px]">
      <label className="sr-only">{label}</label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full appearance-none bg-white/[0.06] hover:bg-white/[0.1] border border-white/[0.06] rounded-xl px-3 py-2.5 text-sm text-fg cursor-pointer transition-all duration-200 focus:outline-none focus:border-accent/40 focus:ring-1 focus:ring-accent/20"
      >
        {options.map((opt) => (
          <option key={opt.value} value={opt.value} className="bg-[#1a1a2e] text-fg">
            {opt.label}
          </option>
        ))}
      </select>
      <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted pointer-events-none" />
    </div>
  );
}

// ---------------------------------------------------------------------------
// Program card (list view)
// ---------------------------------------------------------------------------

function ProgramCard({
  program,
  locale,
  t,
  onView,
}: {
  program: WorkoutProgram;
  locale: Locale;
  t: ReturnType<typeof getDictionary>;
  onView: () => void;
}) {
  const name = locale === "ne" && program.name_ne ? program.name_ne : program.name;

  return (
    <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 flex flex-col gap-3 transition-all duration-200 hover:border-white/[0.12]">
      {/* Title + difficulty */}
      <div className="flex items-start justify-between gap-2">
        <h3 className="text-sm font-semibold text-fg leading-tight">{name}</h3>
        <DifficultyBadge difficulty={program.difficulty} t={t} />
      </div>

      {/* Description */}
      {program.description && (
        <p className="text-xs text-fg-muted line-clamp-2">{program.description}</p>
      )}

      {/* Meta: weeks, equipment, goal, type */}
      <div className="flex flex-wrap items-center gap-2 text-xs text-fg-muted">
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-white/[0.06]">
          <Calendar className="w-3 h-3" />
          {program.duration_weeks} {t.programs.weeks}
        </span>
        {program.program_type && (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-sky-500/10 text-sky-400">
            <Zap className="w-3 h-3" />
            {program.program_type.replace(/_/g, " ")}
          </span>
        )}
        {program.equipment && program.equipment !== "none" && (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-white/[0.06]">
            <Dumbbell className="w-3 h-3" />
            {program.equipment.replace(/_/g, " ")}
          </span>
        )}
        {program.goal && (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-accent/10 text-accent">
            <Target className="w-3 h-3" />
            {program.goal.replace(/_/g, " ")}
          </span>
        )}
      </div>

      {/* View details button */}
      <button
        onClick={onView}
        className="mt-auto w-full py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] text-fg text-sm font-medium transition-all duration-200 flex items-center justify-center gap-2"
      >
        <ChevronRight className="w-4 h-4" />
        {t.programs.programDetail}
      </button>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Day card (program detail)
// ---------------------------------------------------------------------------

function DayCard({
  day,
  isToday,
  locale,
}: {
  day: ProgramDay;
  isToday: boolean;
  locale: Locale;
}) {
  const exercises = (day.exercises ?? []) as Array<{
    name?: string;
    name_ne?: string;
    muscle_groups?: string[];
    sets?: number;
    reps?: number;
  }>;

  return (
    <div
      className={`bg-gradient-to-b from-white/[0.08] to-white/[0.02] border rounded-2xl p-4 space-y-3 transition-all duration-200 ${
        isToday
          ? "border-accent/40 ring-1 ring-accent/20"
          : "border-white/[0.06]"
      }`}
    >
      {/* Day header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="text-sm font-semibold text-fg">
            {day.day_name || `Day ${day.day_number}`}
          </span>
          {isToday && (
            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-accent/15 text-accent">
              Today
            </span>
          )}
        </div>
        {day.rest_day && (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium bg-purple-500/15 text-purple-400">
            <Moon className="w-3 h-3" />
            Rest Day
          </span>
        )}
      </div>

      {/* Exercises list or rest day */}
      {day.rest_day ? (
        <div className="flex items-center justify-center py-4">
          <div className="flex flex-col items-center gap-2">
            <Moon className="w-8 h-8 text-purple-400/50" />
            <p className="text-xs text-fg-muted">Recovery &amp; rest</p>
          </div>
        </div>
      ) : exercises.length > 0 ? (
        <div className="space-y-2">
          {exercises.map((ex, i) => {
            const name =
              locale === "ne" && ex.name_ne ? ex.name_ne : ex.name ?? `Exercise ${i + 1}`;
            return (
              <div
                key={i}
                className="flex items-center justify-between py-1.5 px-2 rounded-lg bg-white/[0.03]"
              >
                <div className="flex items-center gap-2">
                  <div className="w-5 h-5 rounded-md bg-white/[0.06] flex items-center justify-center text-[10px] text-fg-muted font-mono">
                    {i + 1}
                  </div>
                  <span className="text-xs text-fg font-medium">{name}</span>
                </div>
                {ex.sets && ex.reps && (
                  <span className="text-[11px] text-fg-muted font-mono">
                    {ex.sets}x{ex.reps}
                  </span>
                )}
              </div>
            );
          })}
          {/* Muscle chips from exercises */}
          {(() => {
            const muscles = new Set<string>();
            exercises.forEach((ex) => {
              ex.muscle_groups?.forEach((g) => muscles.add(g));
            });
            if (muscles.size === 0) return null;
            return (
              <div className="flex flex-wrap gap-1 pt-1">
                {Array.from(muscles).map((g) => (
                  <MuscleChip key={g} group={g} />
                ))}
              </div>
            );
          })()}
        </div>
      ) : (
        <p className="text-xs text-fg-muted py-2">No exercises listed</p>
      )}
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

type View = "list" | "detail";

export default function MemberProgramsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  // View state
  const [view, setView] = useState<View>("list");
  const [selectedProgramId, setSelectedProgramId] = useState<string | null>(null);

  // List data
  const [programs, setPrograms] = useState<WorkoutProgram[]>([]);
  const [loading, setLoading] = useState(true);

  // Filters
  const [typeFilter, setTypeFilter] = useState("");
  const [difficultyFilter, setDifficultyFilter] = useState("");
  const [equipmentFilter, setEquipmentFilter] = useState("");

  // Detail data
  const [detailProgram, setDetailProgram] = useState<WorkoutProgram | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // Enrollment
  const [enrollment, setEnrollment] = useState<UserProgramEnrollment | null>(null);
  const [enrollmentLoading, setEnrollmentLoading] = useState(true);
  const [enrolling, setEnrolling] = useState(false);
  const [completing, setCompleting] = useState(false);
  const [completedBanner, setCompletedBanner] = useState(false);
  const [enrolledBanner, setEnrolledBanner] = useState(false);

  // Load current enrollment
  useEffect(() => {
    (async () => {
      try {
        const current = await api.getCurrentProgram();
        setEnrollment(current);
      } catch {
        // No active enrollment
        setEnrollment(null);
      } finally {
        setEnrollmentLoading(false);
      }
    })();
  }, []);

  // Load programs list
  const loadPrograms = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.listPrograms({
        type: typeFilter || undefined,
        difficulty: difficultyFilter || undefined,
        equipment: equipmentFilter || undefined,
        limit: 50,
      });
      setPrograms(res.data ?? []);
    } catch (err) {
      console.error("Failed to load programs:", err);
      setPrograms([]);
    } finally {
      setLoading(false);
    }
  }, [typeFilter, difficultyFilter, equipmentFilter]);

  useEffect(() => {
    loadPrograms();
  }, [loadPrograms]);

  // View program detail
  async function handleViewDetail(programId: string) {
    setSelectedProgramId(programId);
    setView("detail");
    setDetailLoading(true);
    try {
      const detail = await api.getProgram(programId);
      setDetailProgram(detail);
    } catch (err) {
      console.error("Failed to load program detail:", err);
      setDetailProgram(null);
    } finally {
      setDetailLoading(false);
    }
  }

  // Enroll in program
  async function handleEnroll(programId: string) {
    setEnrolling(true);
    try {
      const result = await api.enrollInProgram(programId);
      setEnrollment(result);
      setEnrolledBanner(true);
      setTimeout(() => setEnrolledBanner(false), 4000);
    } catch (err) {
      console.error("Failed to enroll:", err);
    } finally {
      setEnrolling(false);
    }
  }

  // Complete today's workout
  async function handleCompleteDay() {
    setCompleting(true);
    try {
      const result = await api.completeProgramDay();
      setEnrollment(result);
      setCompletedBanner(true);
      setTimeout(() => setCompletedBanner(false), 4000);
    } catch (err) {
      console.error("Failed to complete day:", err);
    } finally {
      setCompleting(false);
    }
  }

  // Back to list
  function handleBackToList() {
    setView("list");
    setSelectedProgramId(null);
    setDetailProgram(null);
  }

  // Check if current program matches selected
  const isEnrolledInSelected =
    enrollment && selectedProgramId && enrollment.program_id === selectedProgramId;

  // Filter options
  const typeOptions = [
    { value: "", label: t.programs.allTypes },
    { value: "strength", label: t.exercises.strength },
    { value: "cardio", label: t.exercises.cardio },
    { value: "flexibility", label: t.exercises.flexibility },
    { value: "hiit", label: t.exercises.hiit },
    { value: "bodyweight", label: t.programs.bodyweight },
    { value: "mixed", label: t.programs.mixed },
  ];

  const difficultyOptions = [
    { value: "", label: t.exercises.allTypes },
    { value: "beginner", label: t.exercises.beginner },
    { value: "intermediate", label: t.exercises.intermediate },
    { value: "advanced", label: t.exercises.advanced },
  ];

  const equipmentOptions = [
    { value: "", label: t.exercises.allEquipment },
    { value: "none", label: t.exercises.none },
    { value: "dumbbells", label: t.exercises.dumbbells },
    { value: "barbell", label: t.exercises.barbell },
    { value: "machine", label: t.exercises.machine },
    { value: "kettlebell", label: t.exercises.kettlebell },
  ];

  // =========================================================================
  // DETAIL VIEW
  // =========================================================================

  if (view === "detail") {
    return (
      <div className="space-y-6">
        {/* Back button */}
        <button
          onClick={handleBackToList}
          className="inline-flex items-center gap-1.5 text-sm text-fg-muted hover:text-fg transition-colors"
        >
          <ChevronLeft className="w-4 h-4" />
          {t.programs.browsePrograms}
        </button>

        {/* Banners */}
        {enrolledBanner && (
          <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-emerald-600/20 border border-emerald-500/30 text-emerald-400 text-sm font-medium transition-all duration-300">
            <CheckCircle2 className="w-4 h-4 shrink-0" />
            {t.programs.enrolled}
          </div>
        )}
        {completedBanner && (
          <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-accent/20 border border-accent/30 text-accent text-sm font-medium transition-all duration-300">
            <Trophy className="w-4 h-4 shrink-0" />
            {t.programs.dayCompleted}
          </div>
        )}

        {detailLoading ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 className="w-8 h-8 text-accent animate-spin" />
          </div>
        ) : !detailProgram ? (
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
            <div className="flex flex-col items-center justify-center py-20 gap-3">
              <p className="text-sm text-fg-muted">{t.programs.noPrograms}</p>
            </div>
          </div>
        ) : (
          <>
            {/* Program header */}
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
              <div className="p-5 space-y-4">
                <div className="flex items-start justify-between gap-3">
                  <h1 className="text-lg font-semibold text-fg">
                    {locale === "ne" && detailProgram.name_ne
                      ? detailProgram.name_ne
                      : detailProgram.name}
                  </h1>
                  <DifficultyBadge difficulty={detailProgram.difficulty} t={t} />
                </div>

                {detailProgram.description && (
                  <p className="text-sm text-fg-muted leading-relaxed">
                    {detailProgram.description}
                  </p>
                )}

                {/* Meta */}
                <div className="flex flex-wrap items-center gap-2 text-xs text-fg-muted">
                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-white/[0.06]">
                    <Calendar className="w-3 h-3" />
                    {detailProgram.duration_weeks} {t.programs.weeks}
                  </span>
                  {detailProgram.program_type && (
                    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-sky-500/10 text-sky-400">
                      <Zap className="w-3 h-3" />
                      {detailProgram.program_type.replace(/_/g, " ")}
                    </span>
                  )}
                  {detailProgram.equipment && detailProgram.equipment !== "none" && (
                    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-white/[0.06]">
                      <Dumbbell className="w-3 h-3" />
                      {detailProgram.equipment.replace(/_/g, " ")}
                    </span>
                  )}
                  {detailProgram.goal && (
                    <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-accent/10 text-accent">
                      <Target className="w-3 h-3" />
                      {detailProgram.goal.replace(/_/g, " ")}
                    </span>
                  )}
                </div>

                {/* Progress indicator if enrolled */}
                {isEnrolledInSelected && enrollment && (
                  <div className="flex items-center gap-3 py-2 px-3 rounded-xl bg-accent/10 border border-accent/20">
                    <Trophy className="w-4 h-4 text-accent shrink-0" />
                    <div className="flex-1">
                      <p className="text-xs font-medium text-accent">
                        {t.programs.week} {enrollment.current_week}, {t.programs.day}{" "}
                        {enrollment.current_day}
                      </p>
                      {/* Progress bar */}
                      <div className="mt-1.5 h-1.5 rounded-full bg-white/[0.06] overflow-hidden">
                        <div
                          className="h-full rounded-full bg-accent transition-all duration-500"
                          style={{
                            width: `${Math.min(
                              100,
                              detailProgram.duration_weeks > 0
                                ? ((enrollment.current_week - 1) * 7 +
                                    enrollment.current_day) /
                                    (detailProgram.duration_weeks * 7) *
                                    100
                                : 0
                            )}%`,
                          }}
                        />
                      </div>
                    </div>
                  </div>
                )}

                {/* Enroll / Complete day button */}
                {isEnrolledInSelected ? (
                  <button
                    onClick={handleCompleteDay}
                    disabled={completing}
                    className="w-full py-2.5 rounded-xl bg-accent text-bg text-sm font-medium transition-all duration-200 hover:bg-accent/90 disabled:opacity-50 flex items-center justify-center gap-2"
                  >
                    {completing ? (
                      <Loader2 className="w-4 h-4 animate-spin" />
                    ) : (
                      <>
                        <CheckCircle2 className="w-4 h-4" />
                        {t.programs.completeDay}
                      </>
                    )}
                  </button>
                ) : (
                  <button
                    onClick={() => handleEnroll(detailProgram.id)}
                    disabled={enrolling}
                    className="w-full py-2.5 rounded-xl bg-accent text-bg text-sm font-medium transition-all duration-200 hover:bg-accent/90 disabled:opacity-50 flex items-center justify-center gap-2"
                  >
                    {enrolling ? (
                      <Loader2 className="w-4 h-4 animate-spin" />
                    ) : (
                      <>
                        <Play className="w-4 h-4" />
                        {t.programs.enrollNow}
                      </>
                    )}
                  </button>
                )}
              </div>
            </div>

            {/* Week-by-week schedule */}
            {detailProgram.days && detailProgram.days.length > 0 && (
              <div className="space-y-6">
                {/* Group days by week */}
                {(() => {
                  const weekMap = new Map<number, ProgramDay[]>();
                  for (const day of detailProgram.days) {
                    const week = day.week_number;
                    if (!weekMap.has(week)) weekMap.set(week, []);
                    weekMap.get(week)!.push(day);
                  }
                  // Sort weeks
                  const weeks = Array.from(weekMap.entries()).sort(
                    ([a], [b]) => a - b
                  );

                  return weeks.map(([weekNum, days]) => {
                    // Sort days within week
                    const sortedDays = [...days].sort(
                      (a, b) => a.day_number - b.day_number
                    );

                    // Check if this is the current week for enrolled program
                    const isCurrentWeek =
                      isEnrolledInSelected &&
                      enrollment &&
                      enrollment.current_week === weekNum;

                    return (
                      <div key={weekNum} className="space-y-3">
                        <div className="flex items-center gap-2">
                          <h2 className="text-sm font-semibold text-fg">
                            {t.programs.week} {weekNum}
                          </h2>
                          {isCurrentWeek && (
                            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-accent/15 text-accent">
                              Current
                            </span>
                          )}
                        </div>
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                          {sortedDays.map((day) => {
                            const isTodayDay =
                              isCurrentWeek &&
                              enrollment &&
                              enrollment.current_day === day.day_number;
                            return (
                              <DayCard
                                key={day.id}
                                day={day}
                                isToday={!!isTodayDay}
                                locale={locale}
                              />
                            );
                          })}
                        </div>
                      </div>
                    );
                  });
                })()}
              </div>
            )}
          </>
        )}
      </div>
    );
  }

  // =========================================================================
  // LIST VIEW
  // =========================================================================

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">{t.programs.title}</h1>

      {/* Banners */}
      {enrolledBanner && (
        <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-emerald-600/20 border border-emerald-500/30 text-emerald-400 text-sm font-medium transition-all duration-300">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          {t.programs.enrolled}
        </div>
      )}
      {completedBanner && (
        <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-accent/20 border border-accent/30 text-accent text-sm font-medium transition-all duration-300">
          <Trophy className="w-4 h-4 shrink-0" />
          {t.programs.dayCompleted}
        </div>
      )}

      {/* Current program section */}
      {!enrollmentLoading && enrollment && enrollment.program && (
        <div className="bg-gradient-to-b from-accent/[0.08] to-accent/[0.02] border border-accent/20 rounded-2xl">
          <div className="px-5 py-4 border-b border-accent/10 flex items-center gap-2">
            <Trophy className="w-4 h-4 text-accent" />
            <h2 className="text-sm font-semibold text-fg">{t.programs.currentProgram}</h2>
          </div>
          <div className="p-5 space-y-4">
            {/* Program info */}
            <div className="flex items-start justify-between gap-3">
              <div>
                <h3 className="text-base font-semibold text-fg">
                  {locale === "ne" && enrollment.program.name_ne
                    ? enrollment.program.name_ne
                    : enrollment.program.name}
                </h3>
                <p className="text-xs text-fg-muted mt-1">
                  {enrollment.program.program_type?.replace(/_/g, " ")}
                  {enrollment.program.equipment && enrollment.program.equipment !== "none"
                    ? ` / ${enrollment.program.equipment.replace(/_/g, " ")}`
                    : ""}
                </p>
              </div>
              <DifficultyBadge difficulty={enrollment.program.difficulty} t={t} />
            </div>

            {/* Progress */}
            <div className="flex items-center gap-3 py-2 px-3 rounded-xl bg-white/[0.04]">
              <div className="flex-1">
                <p className="text-xs font-medium text-accent">
                  {t.programs.week} {enrollment.current_week}, {t.programs.day}{" "}
                  {enrollment.current_day}
                </p>
                <div className="mt-1.5 h-1.5 rounded-full bg-white/[0.06] overflow-hidden">
                  <div
                    className="h-full rounded-full bg-accent transition-all duration-500"
                    style={{
                      width: `${Math.min(
                        100,
                        enrollment.program.duration_weeks > 0
                          ? ((enrollment.current_week - 1) * 7 +
                              enrollment.current_day) /
                              (enrollment.program.duration_weeks * 7) *
                              100
                          : 0
                      )}%`,
                    }}
                  />
                </div>
              </div>
            </div>

            {/* Today's workout */}
            {enrollment.today_workout && (
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Clock className="w-3.5 h-3.5 text-fg-muted" />
                  <span className="text-xs font-medium text-fg-muted uppercase tracking-wider">
                    {t.programs.todayWorkout}
                  </span>
                </div>
                {enrollment.today_workout.rest_day ? (
                  <div className="flex items-center gap-2 py-3 px-4 rounded-xl bg-purple-500/10 border border-purple-500/20">
                    <Moon className="w-4 h-4 text-purple-400" />
                    <span className="text-sm font-medium text-purple-400">
                      {t.programs.restDay}
                    </span>
                  </div>
                ) : (
                  <div className="space-y-1.5">
                    {(
                      (enrollment.today_workout.exercises ?? []) as Array<{
                        name?: string;
                        name_ne?: string;
                        sets?: number;
                        reps?: number;
                      }>
                    ).map((ex, i) => (
                      <div
                        key={i}
                        className="flex items-center justify-between py-1.5 px-3 rounded-lg bg-white/[0.03]"
                      >
                        <span className="text-xs text-fg font-medium">
                          {locale === "ne" && ex.name_ne
                            ? ex.name_ne
                            : ex.name ?? `Exercise ${i + 1}`}
                        </span>
                        {ex.sets && ex.reps && (
                          <span className="text-[11px] text-fg-muted font-mono">
                            {ex.sets}x{ex.reps}
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {/* Action buttons */}
            <div className="flex gap-3">
              <button
                onClick={() => handleViewDetail(enrollment.program_id)}
                className="flex-1 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] text-fg text-sm font-medium transition-all duration-200 flex items-center justify-center gap-2"
              >
                <ChevronRight className="w-4 h-4" />
                {t.programs.programDetail}
              </button>
              {enrollment.today_workout && !enrollment.today_workout.rest_day && (
                <button
                  onClick={handleCompleteDay}
                  disabled={completing}
                  className="flex-1 py-2 rounded-xl bg-accent text-bg text-sm font-medium transition-all duration-200 hover:bg-accent/90 disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {completing ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <>
                      <CheckCircle2 className="w-4 h-4" />
                      {t.programs.completeDay}
                    </>
                  )}
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
        <div className="flex items-center gap-2 mb-3">
          <Filter className="w-4 h-4 text-fg-muted" />
          <span className="text-xs font-medium text-fg-muted uppercase tracking-wider">
            {t.programs.browsePrograms}
          </span>
        </div>
        <div className="flex flex-wrap gap-3">
          <FilterSelect
            value={typeFilter}
            onChange={setTypeFilter}
            options={typeOptions}
            label={t.programs.allTypes}
          />
          <FilterSelect
            value={difficultyFilter}
            onChange={setDifficultyFilter}
            options={difficultyOptions}
            label={t.exercises.allTypes}
          />
          <FilterSelect
            value={equipmentFilter}
            onChange={setEquipmentFilter}
            options={equipmentOptions}
            label={t.exercises.allEquipment}
          />
        </div>
      </div>

      {/* Programs grid */}
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="w-8 h-8 text-accent animate-spin" />
        </div>
      ) : programs.length === 0 ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
          <div className="flex flex-col items-center justify-center py-20 gap-4">
            <div className="w-14 h-14 rounded-2xl bg-white/[0.04] flex items-center justify-center">
              <Dumbbell className="w-7 h-7 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">{t.programs.noPrograms}</p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {programs.map((program) => (
            <ProgramCard
              key={program.id}
              program={program}
              locale={locale}
              t={t}
              onView={() => handleViewDetail(program.id)}
            />
          ))}
        </div>
      )}

      {/* Result count */}
      {!loading && programs.length > 0 && (
        <p className="text-center text-xs text-fg-muted">
          {programs.length} {programs.length === 1 ? "program" : "programs"}
        </p>
      )}
    </div>
  );
}
