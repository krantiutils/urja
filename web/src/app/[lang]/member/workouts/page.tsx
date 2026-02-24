"use client";

import { useEffect, useState, useMemo } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type { Locale, WorkoutPlan, WorkoutTemplate, Exercise } from "@/types";
import {
  Loader2,
  Dumbbell,
  Clock,
  Search,
  ChevronRight,
  CheckCircle2,
  Sparkles,
  User,
  ClipboardList,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type Tab = "my_plan" | "browse" | "find";

type Goal = "lose_weight" | "build_muscle" | "stay_fit";
type Difficulty = "beginner" | "intermediate" | "advanced";
type DaysPerWeek = "2-3" | "4-5" | "6-7";

// ---------------------------------------------------------------------------
// Muscle-group color map
// ---------------------------------------------------------------------------

const MUSCLE_COLORS: Record<string, { bg: string; text: string; dot: string }> = {
  chest:        { bg: "bg-rose-500/15",    text: "text-rose-400",    dot: "#fb7185" },
  shoulders:    { bg: "bg-orange-500/15",  text: "text-orange-400",  dot: "#fb923c" },
  biceps:       { bg: "bg-amber-500/15",   text: "text-amber-400",   dot: "#fbbf24" },
  triceps:      { bg: "bg-yellow-500/15",  text: "text-yellow-400",  dot: "#facc15" },
  forearms:     { bg: "bg-lime-500/15",    text: "text-lime-400",    dot: "#a3e635" },
  core:         { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "#34d399" },
  obliques:     { bg: "bg-teal-500/15",    text: "text-teal-400",    dot: "#2dd4bf" },
  lats:         { bg: "bg-cyan-500/15",    text: "text-cyan-400",    dot: "#22d3ee" },
  back:         { bg: "bg-sky-500/15",     text: "text-sky-400",     dot: "#38bdf8" },
  lower_back:   { bg: "bg-blue-500/15",    text: "text-blue-400",    dot: "#60a5fa" },
  traps:        { bg: "bg-indigo-500/15",  text: "text-indigo-400",  dot: "#818cf8" },
  quadriceps:   { bg: "bg-violet-500/15",  text: "text-violet-400",  dot: "#a78bfa" },
  hamstrings:   { bg: "bg-purple-500/15",  text: "text-purple-400",  dot: "#c084fc" },
  glutes:       { bg: "bg-fuchsia-500/15", text: "text-fuchsia-400", dot: "#e879f9" },
  calves:       { bg: "bg-pink-500/15",    text: "text-pink-400",    dot: "#f472b6" },
  hip_flexors:  { bg: "bg-red-500/15",     text: "text-red-400",     dot: "#f87171" },
};

function muscleStyle(group: string) {
  return MUSCLE_COLORS[group] ?? { bg: "bg-white/10", text: "text-fg-muted", dot: "#9ca3af" };
}

// ---------------------------------------------------------------------------
// SVG Body Map
// ---------------------------------------------------------------------------

/** Anatomical positions (cx, cy) for each muscle group on the front-view body */
const MUSCLE_POSITIONS: Record<string, { cx: number; cy: number; rx: number; ry: number }> = {
  // Head area (no muscle, just reference)
  traps:        { cx: 100, cy: 62,  rx: 18, ry: 6 },
  shoulders:    { cx: 64,  cy: 72,  rx: 10, ry: 8 },
  // Mirror shoulder right
  chest:        { cx: 100, cy: 92,  rx: 22, ry: 10 },
  biceps:       { cx: 58,  cy: 102, rx: 7,  ry: 12 },
  triceps:      { cx: 142, cy: 102, rx: 7,  ry: 12 },
  forearms:     { cx: 52,  cy: 130, rx: 6,  ry: 12 },
  core:         { cx: 100, cy: 120, rx: 14, ry: 12 },
  obliques:     { cx: 80,  cy: 118, rx: 6,  ry: 10 },
  lats:         { cx: 76,  cy: 98,  rx: 8,  ry: 10 },
  back:         { cx: 100, cy: 105, rx: 12, ry: 8 },
  lower_back:   { cx: 100, cy: 138, rx: 12, ry: 6 },
  hip_flexors:  { cx: 90,  cy: 155, rx: 8,  ry: 6 },
  glutes:       { cx: 110, cy: 155, rx: 8,  ry: 6 },
  quadriceps:   { cx: 88,  cy: 185, rx: 10, ry: 18 },
  hamstrings:   { cx: 112, cy: 185, rx: 10, ry: 18 },
  calves:       { cx: 88,  cy: 230, rx: 7,  ry: 14 },
};

// Also add mirrored positions for paired muscles
const MUSCLE_POSITIONS_MIRRORED: Record<string, { cx: number; cy: number; rx: number; ry: number }> = {
  shoulders:  { cx: 136, cy: 72,  rx: 10, ry: 8 },
  forearms:   { cx: 148, cy: 130, rx: 6,  ry: 12 },
  obliques:   { cx: 120, cy: 118, rx: 6,  ry: 10 },
  lats:       { cx: 124, cy: 98,  rx: 8,  ry: 10 },
  calves:     { cx: 112, cy: 230, rx: 7,  ry: 14 },
};

function BodyMapSvg({ activeMuscles }: { activeMuscles: Set<string> }) {
  return (
    <svg
      viewBox="0 0 200 270"
      className="w-full max-w-[200px] mx-auto"
      aria-label="Body map"
    >
      {/* Human body outline */}
      <g stroke="currentColor" strokeWidth="1.2" fill="none" className="text-white/20">
        {/* Head */}
        <ellipse cx="100" cy="24" rx="16" ry="18" />
        {/* Neck */}
        <line x1="93" y1="42" x2="93" y2="52" />
        <line x1="107" y1="42" x2="107" y2="52" />
        {/* Torso */}
        <path d="M 68 52 Q 64 52 60 65 Q 55 80 55 95 L 55 100 Q 55 115 58 130 Q 60 142 72 155 L 80 160" />
        <path d="M 132 52 Q 136 52 140 65 Q 145 80 145 95 L 145 100 Q 145 115 142 130 Q 140 142 128 155 L 120 160" />
        {/* Hips */}
        <path d="M 80 160 Q 85 168 85 172 L 85 175" />
        <path d="M 120 160 Q 115 168 115 172 L 115 175" />
        <line x1="80" y1="160" x2="120" y2="160" />
        {/* Left leg */}
        <path d="M 85 175 Q 82 195 80 210 Q 78 230 80 248 L 82 262" />
        <path d="M 85 175 Q 96 195 97 210 Q 98 230 96 248 L 94 262" />
        {/* Right leg */}
        <path d="M 115 175 Q 104 195 103 210 Q 102 230 104 248 L 106 262" />
        <path d="M 115 175 Q 118 195 120 210 Q 122 230 120 248 L 118 262" />
        {/* Left arm */}
        <path d="M 60 65 Q 52 75 48 90 Q 44 105 42 120 Q 40 135 42 150" />
        <path d="M 68 68 Q 64 78 60 90 Q 58 105 56 120 Q 55 135 56 150" />
        {/* Right arm */}
        <path d="M 140 65 Q 148 75 152 90 Q 156 105 158 120 Q 160 135 158 150" />
        <path d="M 132 68 Q 136 78 140 90 Q 142 105 144 120 Q 145 135 144 150" />
        {/* Shoulders top line */}
        <path d="M 68 52 L 132 52" />
      </g>

      {/* Muscle ellipses - primary positions */}
      {Object.entries(MUSCLE_POSITIONS).map(([muscle, pos]) => {
        const isActive = activeMuscles.has(muscle);
        const color = muscleStyle(muscle).dot;
        return (
          <ellipse
            key={muscle}
            cx={pos.cx}
            cy={pos.cy}
            rx={pos.rx}
            ry={pos.ry}
            fill={isActive ? color : "#4b5563"}
            opacity={isActive ? 0.55 : 0.15}
            className="transition-all duration-300"
          />
        );
      })}

      {/* Muscle ellipses - mirrored positions */}
      {Object.entries(MUSCLE_POSITIONS_MIRRORED).map(([muscle, pos]) => {
        const isActive = activeMuscles.has(muscle);
        const color = muscleStyle(muscle).dot;
        return (
          <ellipse
            key={`${muscle}_mirror`}
            cx={pos.cx}
            cy={pos.cy}
            rx={pos.rx}
            ry={pos.ry}
            fill={isActive ? color : "#4b5563"}
            opacity={isActive ? 0.55 : 0.15}
            className="transition-all duration-300"
          />
        );
      })}
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function MuscleChip({ group }: { group: string }) {
  const s = muscleStyle(group);
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-medium ${s.bg} ${s.text}`}
    >
      {group.replace(/_/g, " ")}
    </span>
  );
}

function DifficultyBadge({
  difficulty,
  t,
}: {
  difficulty: string | null;
  t: ReturnType<typeof getDictionary>;
}) {
  if (!difficulty) return null;
  const colorMap: Record<string, string> = {
    beginner:     "bg-emerald-500/15 text-emerald-400 border-emerald-500/20",
    intermediate: "bg-amber-500/15 text-amber-400 border-amber-500/20",
    advanced:     "bg-rose-500/15 text-rose-400 border-rose-500/20",
  };
  const label =
    (t.memberPages.workout as Record<string, string>)[difficulty] ?? difficulty;
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${colorMap[difficulty] ?? "bg-white/10 text-fg-muted border-white/10"}`}
    >
      {label}
    </span>
  );
}

function AssignmentBadge({
  method,
  t,
}: {
  method: string;
  t: ReturnType<typeof getDictionary>;
}) {
  const w = t.memberPages.workout;
  const map: Record<string, { label: string; icon: React.ReactNode; cls: string }> = {
    self:          { label: w.selfAssigned, icon: <User className="w-3 h-3" />, cls: "bg-blue-500/15 text-blue-400 border-blue-500/20" },
    questionnaire: { label: w.questionnaireAssigned, icon: <Sparkles className="w-3 h-3" />, cls: "bg-purple-500/15 text-purple-400 border-purple-500/20" },
    staff:         { label: w.staffAssigned, icon: <ClipboardList className="w-3 h-3" />, cls: "bg-teal-500/15 text-teal-400 border-teal-500/20" },
  };
  const badge = map[method] ?? map.staff;
  return (
    <span
      className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium border ${badge.cls}`}
    >
      {badge.icon}
      {badge.label}
    </span>
  );
}

/** Collects all unique muscle groups from all exercises of a template */
function collectMuscleGroups(template: WorkoutTemplate | null | undefined): Set<string> {
  const groups = new Set<string>();
  if (!template?.exercises) return groups;
  for (const ex of template.exercises) {
    if (ex.muscle_groups) {
      for (const g of ex.muscle_groups) groups.add(g);
    }
  }
  return groups;
}

// ---------------------------------------------------------------------------
// Exercise list item
// ---------------------------------------------------------------------------

function ExerciseRow({
  exercise,
  locale,
  t,
}: {
  exercise: Exercise;
  locale: Locale;
  t: ReturnType<typeof getDictionary>;
}) {
  const w = t.memberPages.workout;
  const name = locale === "ne" && exercise.name_ne ? exercise.name_ne : exercise.name;

  return (
    <div className="flex flex-col gap-2 py-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-fg">{name}</span>
        <div className="flex items-center gap-3 text-xs text-fg-muted">
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/[0.06] font-mono">
            {exercise.sets} {w.sets} x {exercise.reps} {w.reps}
          </span>
          {exercise.rest_seconds > 0 && (
            <span className="inline-flex items-center gap-1">
              <Clock className="w-3 h-3" />
              {exercise.rest_seconds}{w.seconds}
            </span>
          )}
        </div>
      </div>
      {exercise.muscle_groups && exercise.muscle_groups.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {exercise.muscle_groups.map((g) => (
            <MuscleChip key={g} group={g} />
          ))}
        </div>
      )}
      {exercise.notes && (
        <p className="text-xs text-fg-muted italic">{exercise.notes}</p>
      )}
    </div>
  );
}

// ---------------------------------------------------------------------------
// Template card (for Browse tab)
// ---------------------------------------------------------------------------

function TemplateCard({
  template,
  locale,
  t,
  onChoose,
  choosing,
}: {
  template: WorkoutTemplate;
  locale: Locale;
  t: ReturnType<typeof getDictionary>;
  onChoose: (id: string) => void;
  choosing: boolean;
}) {
  const w = t.memberPages.workout;
  const name = locale === "ne" && template.name_ne ? template.name_ne : template.name;
  const desc =
    locale === "ne" && template.description_ne
      ? template.description_ne
      : template.description;
  const muscles = collectMuscleGroups(template);

  return (
    <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 flex flex-col gap-3 transition-all duration-200 hover:border-white/[0.12]">
      <div className="flex items-start justify-between gap-2">
        <h3 className="text-sm font-semibold text-fg leading-tight">{name}</h3>
        <DifficultyBadge difficulty={template.difficulty} t={t} />
      </div>

      {desc && <p className="text-xs text-fg-muted line-clamp-2">{desc}</p>}

      <div className="flex flex-wrap items-center gap-3 text-xs text-fg-muted">
        {template.duration_minutes != null && (
          <span className="inline-flex items-center gap-1">
            <Clock className="w-3.5 h-3.5" />
            {template.duration_minutes} {w.minutes}
          </span>
        )}
        <span className="inline-flex items-center gap-1">
          <Dumbbell className="w-3.5 h-3.5" />
          {template.exercises?.length ?? 0} {w.exercises}
        </span>
      </div>

      {muscles.size > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {Array.from(muscles).map((g) => (
            <MuscleChip key={g} group={g} />
          ))}
        </div>
      )}

      <button
        onClick={() => onChoose(template.id)}
        disabled={choosing}
        className="mt-auto w-full py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-sm font-medium transition-all duration-200 flex items-center justify-center gap-2"
      >
        {choosing ? (
          <Loader2 className="w-4 h-4 animate-spin" />
        ) : (
          <>
            <CheckCircle2 className="w-4 h-4" />
            {w.choosePlan}
          </>
        )}
      </button>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function MemberWorkoutsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const w = t.memberPages.workout;
  const { user } = useAuth();

  // State
  const [tab, setTab] = useState<Tab>("my_plan");
  const [loading, setLoading] = useState(true);
  const [plan, setPlan] = useState<WorkoutPlan | null>(null);

  // Browse state
  const [templates, setTemplates] = useState<WorkoutTemplate[]>([]);
  const [browseLoading, setBrowseLoading] = useState(false);
  const [goalFilter, setGoalFilter] = useState<string>("");
  const [diffFilter, setDiffFilter] = useState<string>("");
  const [choosingId, setChoosingId] = useState<string | null>(null);
  const [successBanner, setSuccessBanner] = useState(false);

  // Questionnaire state
  const [qGoal, setQGoal] = useState<Goal | "">("");
  const [qExp, setQExp] = useState<Difficulty | "">("");
  const [qDays, setQDays] = useState<DaysPerWeek | "">("");
  const [qLoading, setQLoading] = useState(false);
  const [recPlan, setRecPlan] = useState<WorkoutPlan | null>(null);
  const [qAssigning, setQAssigning] = useState(false);

  const orgId = user?.org_id ?? "";

  // ---- Load current plan on mount ----
  useEffect(() => {
    (async () => {
      try {
        const p = await api.getMyWorkoutPlan(orgId);
        setPlan(p);
      } catch {
        // no plan or error
      } finally {
        setLoading(false);
      }
    })();
  }, [orgId]);

  // ---- Load templates when browse tab is selected ----
  useEffect(() => {
    if (tab !== "browse") return;
    setBrowseLoading(true);
    (async () => {
      try {
        const res = await api.browseTemplates({
          organization_id: orgId,
          goal: goalFilter || undefined,
          difficulty: diffFilter || undefined,
        });
        setTemplates(res.data ?? []);
      } catch {
        setTemplates([]);
      } finally {
        setBrowseLoading(false);
      }
    })();
  }, [tab, orgId, goalFilter, diffFilter]);

  // ---- Derived ----
  const template = plan?.template;
  const activeMuscles = useMemo(() => collectMuscleGroups(template), [template]);

  // ---- Handlers ----
  async function handleChoosePlan(templateId: string) {
    setChoosingId(templateId);
    try {
      const newPlan = await api.selfAssignPlan({
        organization_id: orgId,
        workout_template_id: templateId,
      });
      setPlan(newPlan);
      setSuccessBanner(true);
      setTimeout(() => setSuccessBanner(false), 4000);
      setTimeout(() => setTab("my_plan"), 600);
    } catch (err) {
      console.error("Failed to assign plan:", err);
    } finally {
      setChoosingId(null);
    }
  }

  async function handleGetRecommendation() {
    if (!qGoal || !qExp || !qDays) return;
    setQLoading(true);
    try {
      const rec = await api.recommendPlan({
        organization_id: orgId,
        goal: qGoal,
        experience: qExp,
        days_per_week: qDays,
      });
      setRecPlan(rec);
    } catch (err) {
      console.error("Failed to get recommendation:", err);
    } finally {
      setQLoading(false);
    }
  }

  async function handleAssignRecommended() {
    if (!recPlan) return;
    setQAssigning(true);
    try {
      setPlan(recPlan);
      setSuccessBanner(true);
      setTimeout(() => setSuccessBanner(false), 4000);
      setTimeout(() => setTab("my_plan"), 600);
    } finally {
      setQAssigning(false);
    }
  }

  // ---- Loading screen ----
  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  // ---- Tab definition ----
  const tabs: { id: Tab; label: string }[] = [
    { id: "my_plan", label: w.myPlan },
    { id: "browse", label: w.browsePlans },
    { id: "find", label: w.findMyPlan },
  ];

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">
        {t.memberPages.workouts}
      </h1>

      {/* Success banner */}
      {successBanner && (
        <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-emerald-600/20 border border-emerald-500/30 text-emerald-400 text-sm font-medium transition-all duration-300">
          <CheckCircle2 className="w-4 h-4 shrink-0" />
          {w.planAssigned}
        </div>
      )}

      {/* Tab bar */}
      <div className="flex gap-1 p-1 rounded-xl bg-white/[0.04] border border-white/[0.06]">
        {tabs.map((t) => (
          <button
            key={t.id}
            onClick={() => setTab(t.id)}
            className={`flex-1 py-2 px-3 rounded-lg text-sm font-medium transition-all duration-200 ${
              tab === t.id
                ? "bg-white/[0.1] text-fg shadow-sm"
                : "text-fg-muted hover:text-fg hover:bg-white/[0.04]"
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {/* ========== TAB 1: My Plan ========== */}
      {tab === "my_plan" && (
        <div className="space-y-6">
          {!template ? (
            /* Empty state */
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
              <div className="flex flex-col items-center justify-center py-16 gap-4">
                <div className="w-14 h-14 rounded-2xl bg-white/[0.04] flex items-center justify-center">
                  <ClipboardList className="w-7 h-7 text-fg-muted" />
                </div>
                <p className="text-sm text-fg-muted">{w.noPlan}</p>
                <button
                  onClick={() => setTab("browse")}
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent/20 hover:bg-accent/30 text-accent text-sm font-medium transition-all duration-200"
                >
                  <Search className="w-4 h-4" />
                  {w.browsePlans}
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          ) : (
            <>
              {/* Plan header card */}
              <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
                <div className="px-5 py-4 border-b border-white/[0.06]">
                  <h2 className="text-sm font-semibold text-fg">{w.currentPlan}</h2>
                </div>
                <div className="p-5 space-y-4">
                  <div className="flex items-start justify-between gap-3">
                    <h3 className="text-lg font-semibold text-fg">
                      {locale === "ne" && template.name_ne ? template.name_ne : template.name}
                    </h3>
                    <DifficultyBadge difficulty={template.difficulty} t={t} />
                  </div>

                  {template.description && (
                    <p className="text-sm text-fg-muted">{template.description}</p>
                  )}

                  <div className="flex flex-wrap items-center gap-3 text-sm">
                    {template.duration_minutes != null && (
                      <span className="inline-flex items-center gap-1.5 text-fg-muted">
                        <Clock className="w-4 h-4" />
                        {template.duration_minutes} {w.minutes}
                      </span>
                    )}
                    {template.exercises && (
                      <span className="inline-flex items-center gap-1.5 text-fg-muted">
                        <Dumbbell className="w-4 h-4" />
                        {template.exercises.length} {w.exercises}
                      </span>
                    )}
                    {plan.assignment_method && (
                      <AssignmentBadge method={plan.assignment_method} t={t} />
                    )}
                  </div>
                </div>
              </div>

              {/* Body Map + Muscle Groups */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Body Map */}
                <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
                  <div className="px-5 py-4 border-b border-white/[0.06]">
                    <h2 className="text-sm font-semibold text-fg">{w.bodyMap}</h2>
                  </div>
                  <div className="p-5">
                    <BodyMapSvg activeMuscles={activeMuscles} />
                  </div>
                </div>

                {/* Target Muscles list */}
                <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
                  <div className="px-5 py-4 border-b border-white/[0.06]">
                    <h2 className="text-sm font-semibold text-fg">{w.targetMuscles}</h2>
                  </div>
                  <div className="p-5">
                    {activeMuscles.size === 0 ? (
                      <p className="text-sm text-fg-muted">{w.muscleGroups}</p>
                    ) : (
                      <div className="flex flex-wrap gap-2">
                        {Array.from(activeMuscles).map((g) => (
                          <MuscleChip key={g} group={g} />
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Exercises list */}
              <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
                <div className="px-5 py-4 border-b border-white/[0.06]">
                  <h2 className="text-sm font-semibold text-fg">
                    {w.exercises}{" "}
                    <span className="text-fg-muted font-normal">
                      ({template.exercises?.length ?? 0})
                    </span>
                  </h2>
                </div>
                <div className="divide-y divide-white/[0.06] px-5">
                  {template.exercises?.map((ex, i) => (
                    <ExerciseRow
                      key={`${ex.name}-${i}`}
                      exercise={ex}
                      locale={locale}
                      t={t}
                    />
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      )}

      {/* ========== TAB 2: Browse Plans ========== */}
      {tab === "browse" && (
        <div className="space-y-5">
          {/* Filter chips */}
          <div className="space-y-3">
            {/* Goal filters */}
            <div className="space-y-1.5">
              <p className="text-xs font-medium text-fg-muted uppercase tracking-wider">
                {w.filterByGoal}
              </p>
              <div className="flex flex-wrap gap-2">
                {[
                  { value: "", label: w.allGoals },
                  { value: "lose_weight", label: w.loseWeight },
                  { value: "build_muscle", label: w.buildMuscle },
                  { value: "stay_fit", label: w.stayFit },
                ].map((opt) => (
                  <button
                    key={opt.value}
                    onClick={() => setGoalFilter(opt.value)}
                    className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200 border ${
                      goalFilter === opt.value
                        ? "bg-accent/20 text-accent border-accent/30"
                        : "bg-white/[0.04] text-fg-muted border-white/[0.06] hover:bg-white/[0.08]"
                    }`}
                  >
                    {opt.label}
                  </button>
                ))}
              </div>
            </div>

            {/* Difficulty filters */}
            <div className="space-y-1.5">
              <p className="text-xs font-medium text-fg-muted uppercase tracking-wider">
                {w.filterByDifficulty}
              </p>
              <div className="flex flex-wrap gap-2">
                {[
                  { value: "", label: w.allDifficulties },
                  { value: "beginner", label: w.beginner },
                  { value: "intermediate", label: w.intermediate },
                  { value: "advanced", label: w.advanced },
                ].map((opt) => (
                  <button
                    key={opt.value}
                    onClick={() => setDiffFilter(opt.value)}
                    className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200 border ${
                      diffFilter === opt.value
                        ? "bg-purple-500/20 text-purple-400 border-purple-500/30"
                        : "bg-white/[0.04] text-fg-muted border-white/[0.06] hover:bg-white/[0.08]"
                    }`}
                  >
                    {opt.label}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Template grid */}
          {browseLoading ? (
            <div className="flex items-center justify-center py-16">
              <Loader2 className="w-6 h-6 text-accent animate-spin" />
            </div>
          ) : templates.length === 0 ? (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
              <div className="flex flex-col items-center justify-center py-16 gap-3">
                <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
                  <Dumbbell className="w-6 h-6 text-fg-muted" />
                </div>
                <p className="text-sm text-fg-muted">{t.common.noResults}</p>
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {templates.map((tmpl) => (
                <TemplateCard
                  key={tmpl.id}
                  template={tmpl}
                  locale={locale}
                  t={t}
                  onChoose={handleChoosePlan}
                  choosing={choosingId === tmpl.id}
                />
              ))}
            </div>
          )}
        </div>
      )}

      {/* ========== TAB 3: Find My Plan (Questionnaire) ========== */}
      {tab === "find" && (
        <div className="space-y-6">
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
            <div className="px-5 py-4 border-b border-white/[0.06]">
              <h2 className="text-sm font-semibold text-fg">{w.questionnaire}</h2>
            </div>
            <div className="p-5 space-y-6">
              {/* Step 1: Goal */}
              <fieldset className="space-y-2.5">
                <legend className="text-sm font-medium text-fg">
                  {w.whatIsYourGoal}
                </legend>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
                  {([
                    { value: "lose_weight" as Goal, label: w.loseWeight },
                    { value: "build_muscle" as Goal, label: w.buildMuscle },
                    { value: "stay_fit" as Goal, label: w.stayFit },
                  ]).map((opt) => (
                    <label
                      key={opt.value}
                      className={`flex items-center gap-3 px-4 py-3 rounded-xl border cursor-pointer transition-all duration-200 ${
                        qGoal === opt.value
                          ? "bg-accent/15 border-accent/40 text-fg"
                          : "bg-white/[0.02] border-white/[0.06] text-fg-muted hover:bg-white/[0.04]"
                      }`}
                    >
                      <input
                        type="radio"
                        name="goal"
                        value={opt.value}
                        checked={qGoal === opt.value}
                        onChange={() => setQGoal(opt.value)}
                        className="sr-only"
                      />
                      <div
                        className={`w-4 h-4 rounded-full border-2 flex items-center justify-center shrink-0 transition-all duration-200 ${
                          qGoal === opt.value
                            ? "border-accent bg-accent"
                            : "border-white/20"
                        }`}
                      >
                        {qGoal === opt.value && (
                          <div className="w-1.5 h-1.5 rounded-full bg-white" />
                        )}
                      </div>
                      <span className="text-sm font-medium">{opt.label}</span>
                    </label>
                  ))}
                </div>
              </fieldset>

              {/* Step 2: Experience */}
              <fieldset className="space-y-2.5">
                <legend className="text-sm font-medium text-fg">
                  {w.experienceLevel}
                </legend>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
                  {([
                    { value: "beginner" as Difficulty, label: w.beginner },
                    { value: "intermediate" as Difficulty, label: w.intermediate },
                    { value: "advanced" as Difficulty, label: w.advanced },
                  ]).map((opt) => (
                    <label
                      key={opt.value}
                      className={`flex items-center gap-3 px-4 py-3 rounded-xl border cursor-pointer transition-all duration-200 ${
                        qExp === opt.value
                          ? "bg-purple-500/15 border-purple-500/40 text-fg"
                          : "bg-white/[0.02] border-white/[0.06] text-fg-muted hover:bg-white/[0.04]"
                      }`}
                    >
                      <input
                        type="radio"
                        name="experience"
                        value={opt.value}
                        checked={qExp === opt.value}
                        onChange={() => setQExp(opt.value)}
                        className="sr-only"
                      />
                      <div
                        className={`w-4 h-4 rounded-full border-2 flex items-center justify-center shrink-0 transition-all duration-200 ${
                          qExp === opt.value
                            ? "border-purple-400 bg-purple-400"
                            : "border-white/20"
                        }`}
                      >
                        {qExp === opt.value && (
                          <div className="w-1.5 h-1.5 rounded-full bg-white" />
                        )}
                      </div>
                      <span className="text-sm font-medium">{opt.label}</span>
                    </label>
                  ))}
                </div>
              </fieldset>

              {/* Step 3: Days per week */}
              <fieldset className="space-y-2.5">
                <legend className="text-sm font-medium text-fg">
                  {w.daysPerWeek}
                </legend>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
                  {([
                    { value: "2-3" as DaysPerWeek, label: w.days23 },
                    { value: "4-5" as DaysPerWeek, label: w.days45 },
                    { value: "6-7" as DaysPerWeek, label: w.days67 },
                  ]).map((opt) => (
                    <label
                      key={opt.value}
                      className={`flex items-center gap-3 px-4 py-3 rounded-xl border cursor-pointer transition-all duration-200 ${
                        qDays === opt.value
                          ? "bg-emerald-500/15 border-emerald-500/40 text-fg"
                          : "bg-white/[0.02] border-white/[0.06] text-fg-muted hover:bg-white/[0.04]"
                      }`}
                    >
                      <input
                        type="radio"
                        name="days"
                        value={opt.value}
                        checked={qDays === opt.value}
                        onChange={() => setQDays(opt.value)}
                        className="sr-only"
                      />
                      <div
                        className={`w-4 h-4 rounded-full border-2 flex items-center justify-center shrink-0 transition-all duration-200 ${
                          qDays === opt.value
                            ? "border-emerald-400 bg-emerald-400"
                            : "border-white/20"
                        }`}
                      >
                        {qDays === opt.value && (
                          <div className="w-1.5 h-1.5 rounded-full bg-white" />
                        )}
                      </div>
                      <span className="text-sm font-medium">{opt.label}</span>
                    </label>
                  ))}
                </div>
              </fieldset>

              {/* Submit */}
              <button
                onClick={handleGetRecommendation}
                disabled={!qGoal || !qExp || !qDays || qLoading}
                className="w-full py-2.5 rounded-xl bg-accent hover:bg-accent/90 disabled:opacity-40 disabled:cursor-not-allowed text-white text-sm font-medium transition-all duration-200 flex items-center justify-center gap-2"
              >
                {qLoading ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <Sparkles className="w-4 h-4" />
                    {w.getRecommendation}
                  </>
                )}
              </button>
            </div>
          </div>

          {/* Recommended plan result */}
          {recPlan?.template && (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-emerald-500/20 rounded-2xl">
              <div className="px-5 py-4 border-b border-white/[0.06] flex items-center gap-2">
                <Sparkles className="w-4 h-4 text-emerald-400" />
                <h2 className="text-sm font-semibold text-fg">{w.recommendedPlan}</h2>
              </div>
              <div className="p-5 space-y-4">
                <div className="flex items-start justify-between gap-3">
                  <h3 className="text-lg font-semibold text-fg">
                    {locale === "ne" && recPlan.template.name_ne
                      ? recPlan.template.name_ne
                      : recPlan.template.name}
                  </h3>
                  <DifficultyBadge difficulty={recPlan.template.difficulty} t={t} />
                </div>

                {recPlan.template.description && (
                  <p className="text-sm text-fg-muted">
                    {recPlan.template.description}
                  </p>
                )}

                <div className="flex flex-wrap items-center gap-3 text-sm text-fg-muted">
                  {recPlan.template.duration_minutes != null && (
                    <span className="inline-flex items-center gap-1.5">
                      <Clock className="w-4 h-4" />
                      {recPlan.template.duration_minutes} {w.minutes}
                    </span>
                  )}
                  <span className="inline-flex items-center gap-1.5">
                    <Dumbbell className="w-4 h-4" />
                    {recPlan.template.exercises?.length ?? 0} {w.exercises}
                  </span>
                </div>

                {/* Muscle groups from recommendation */}
                {(() => {
                  const recMuscles = collectMuscleGroups(recPlan.template);
                  if (recMuscles.size === 0) return null;
                  return (
                    <div className="flex flex-wrap gap-1.5">
                      {Array.from(recMuscles).map((g) => (
                        <MuscleChip key={g} group={g} />
                      ))}
                    </div>
                  );
                })()}

                {/* Exercises preview */}
                {recPlan.template.exercises && recPlan.template.exercises.length > 0 && (
                  <div className="border-t border-white/[0.06] pt-3">
                    <p className="text-xs font-medium text-fg-muted uppercase tracking-wider mb-2">
                      {w.exercises}
                    </p>
                    <div className="divide-y divide-white/[0.04]">
                      {recPlan.template.exercises.map((ex, i) => (
                        <ExerciseRow
                          key={`rec-${ex.name}-${i}`}
                          exercise={ex}
                          locale={locale}
                          t={t}
                        />
                      ))}
                    </div>
                  </div>
                )}

                <button
                  onClick={handleAssignRecommended}
                  disabled={qAssigning}
                  className="w-full py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-sm font-medium transition-all duration-200 flex items-center justify-center gap-2"
                >
                  {qAssigning ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <>
                      <CheckCircle2 className="w-4 h-4" />
                      {w.choosePlan}
                    </>
                  )}
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
