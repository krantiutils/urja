"use client";

import { useEffect, useState, useCallback } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, ExerciseItem } from "@/types";
import {
  Loader2,
  Dumbbell,
  Zap,
  Filter,
  ChevronDown,
  ChevronUp,
  Flame,
  Search,
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

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

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

function DifficultyBadge({ difficulty, t }: { difficulty: string; t: ReturnType<typeof getDictionary> }) {
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

function TypeBadge({ exerciseType, t }: { exerciseType: string; t: ReturnType<typeof getDictionary> }) {
  const colorMap: Record<string, string> = {
    cardio: "bg-sky-500/15 text-sky-400",
    strength: "bg-orange-500/15 text-orange-400",
    flexibility: "bg-teal-500/15 text-teal-400",
    hiit: "bg-rose-500/15 text-rose-400",
  };
  const labelMap: Record<string, string> = {
    cardio: t.exercises.cardio,
    strength: t.exercises.strength,
    flexibility: t.exercises.flexibility,
    hiit: t.exercises.hiit,
  };
  return (
    <span
      className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium ${colorMap[exerciseType] ?? "bg-white/10 text-fg-muted"}`}
    >
      <Zap className="w-3 h-3" />
      {labelMap[exerciseType] ?? exerciseType}
    </span>
  );
}

// ---------------------------------------------------------------------------
// Filter select component
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
// Exercise card
// ---------------------------------------------------------------------------

function ExerciseCard({
  exercise,
  locale,
  t,
  expanded,
  onToggle,
}: {
  exercise: ExerciseItem;
  locale: Locale;
  t: ReturnType<typeof getDictionary>;
  expanded: boolean;
  onToggle: () => void;
}) {
  const name = locale === "ne" && exercise.name_ne ? exercise.name_ne : exercise.name;

  return (
    <div
      className={`bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl transition-all duration-200 hover:border-white/[0.12] ${
        expanded ? "ring-1 ring-accent/20" : ""
      }`}
    >
      {/* Card header - clickable */}
      <button
        onClick={onToggle}
        className="w-full text-left p-5 space-y-3"
      >
        {/* Top row: name + difficulty */}
        <div className="flex items-start justify-between gap-2">
          <h3 className="text-sm font-semibold text-fg leading-tight">{name}</h3>
          <DifficultyBadge difficulty={exercise.difficulty} t={t} />
        </div>

        {/* Type + equipment + calories */}
        <div className="flex flex-wrap items-center gap-2 text-xs text-fg-muted">
          <TypeBadge exerciseType={exercise.exercise_type} t={t} />
          {exercise.equipment && exercise.equipment !== "none" && (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-white/[0.06]">
              <Dumbbell className="w-3 h-3" />
              {exercise.equipment.replace(/_/g, " ")}
            </span>
          )}
          {exercise.calories_per_minute > 0 && (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-orange-500/10 text-orange-400">
              <Flame className="w-3 h-3" />
              {exercise.calories_per_minute} {t.exercises.calPerMin}
            </span>
          )}
        </div>

        {/* Muscle group chips */}
        {exercise.muscle_groups && exercise.muscle_groups.length > 0 && (
          <div className="flex flex-wrap gap-1.5">
            {exercise.muscle_groups.map((g) => (
              <MuscleChip key={g} group={g} />
            ))}
          </div>
        )}

        {/* Expand indicator */}
        <div className="flex items-center justify-center pt-1">
          {expanded ? (
            <ChevronUp className="w-4 h-4 text-fg-muted" />
          ) : (
            <ChevronDown className="w-4 h-4 text-fg-muted" />
          )}
        </div>
      </button>

      {/* Expanded detail */}
      {expanded && (
        <div className="px-5 pb-5 space-y-4 border-t border-white/[0.06] pt-4">
          {/* Description */}
          {exercise.description && (
            <p className="text-sm text-fg-muted leading-relaxed">
              {exercise.description}
            </p>
          )}

          {/* Instructions */}
          {exercise.instructions && (
            <div className="space-y-2">
              <h4 className="text-xs font-semibold text-fg uppercase tracking-wider">
                {t.exercises.instructions}
              </h4>
              <div className="text-sm text-fg-muted leading-relaxed whitespace-pre-line bg-white/[0.03] rounded-xl p-4 border border-white/[0.04]">
                {exercise.instructions}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

// ---------------------------------------------------------------------------
// Muscle group options
// ---------------------------------------------------------------------------

const MUSCLE_GROUPS = [
  "chest",
  "shoulders",
  "biceps",
  "triceps",
  "core",
  "back",
  "quadriceps",
  "hamstrings",
  "glutes",
  "calves",
  "forearms",
  "lats",
  "traps",
  "obliques",
  "lower_back",
  "hip_flexors",
];

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function MemberExercisesPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  // Filters
  const [typeFilter, setTypeFilter] = useState("");
  const [equipmentFilter, setEquipmentFilter] = useState("");
  const [muscleFilter, setMuscleFilter] = useState("");

  // Data
  const [exercises, setExercises] = useState<ExerciseItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [detailCache, setDetailCache] = useState<Record<string, ExerciseItem>>({});
  const [detailLoading, setDetailLoading] = useState<string | null>(null);

  // Load exercises
  const loadExercises = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.listExercises({
        type: typeFilter || undefined,
        equipment: equipmentFilter || undefined,
        muscle_group: muscleFilter || undefined,
        limit: 50,
      });
      setExercises(res.data ?? []);
    } catch (err) {
      console.error("Failed to load exercises:", err);
      setExercises([]);
    } finally {
      setLoading(false);
    }
  }, [typeFilter, equipmentFilter, muscleFilter]);

  useEffect(() => {
    loadExercises();
  }, [loadExercises]);

  // Toggle expand with detail fetch
  async function handleToggle(id: string) {
    if (expandedId === id) {
      setExpandedId(null);
      return;
    }

    setExpandedId(id);

    // Fetch detail if not cached
    if (!detailCache[id]) {
      setDetailLoading(id);
      try {
        const detail = await api.getExercise(id);
        setDetailCache((prev) => ({ ...prev, [id]: detail }));
        // Update exercise in list with full detail
        setExercises((prev) =>
          prev.map((ex) => (ex.id === id ? { ...ex, ...detail } : ex))
        );
      } catch (err) {
        console.error("Failed to load exercise detail:", err);
      } finally {
        setDetailLoading(null);
      }
    }
  }

  // Filter options
  const typeOptions = [
    { value: "", label: t.exercises.allTypes },
    { value: "cardio", label: t.exercises.cardio },
    { value: "strength", label: t.exercises.strength },
    { value: "flexibility", label: t.exercises.flexibility },
    { value: "hiit", label: t.exercises.hiit },
  ];

  const equipmentOptions = [
    { value: "", label: t.exercises.allEquipment },
    { value: "none", label: t.exercises.none },
    { value: "dumbbells", label: t.exercises.dumbbells },
    { value: "barbell", label: t.exercises.barbell },
    { value: "resistance_band", label: t.exercises.resistanceBand },
    { value: "pull_up_bar", label: t.exercises.pullUpBar },
    { value: "machine", label: t.exercises.machine },
    { value: "kettlebell", label: t.exercises.kettlebell },
  ];

  const muscleOptions = [
    { value: "", label: t.exercises.allMuscles },
    ...MUSCLE_GROUPS.map((g) => ({
      value: g,
      label: g.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase()),
    })),
  ];

  const hasActiveFilters = typeFilter || equipmentFilter || muscleFilter;

  return (
    <div className="space-y-6">
      {/* Page title */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-fg">{t.exercises.title}</h1>
        {hasActiveFilters && (
          <button
            onClick={() => {
              setTypeFilter("");
              setEquipmentFilter("");
              setMuscleFilter("");
            }}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium text-fg-muted hover:text-fg bg-white/[0.04] hover:bg-white/[0.08] border border-white/[0.06] transition-all duration-200"
          >
            <Filter className="w-3 h-3" />
            Clear filters
          </button>
        )}
      </div>

      {/* Filter bar */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
        <div className="flex items-center gap-2 mb-3">
          <Filter className="w-4 h-4 text-fg-muted" />
          <span className="text-xs font-medium text-fg-muted uppercase tracking-wider">
            {t.exercises.browseExercises}
          </span>
        </div>
        <div className="flex flex-wrap gap-3">
          <FilterSelect
            value={typeFilter}
            onChange={setTypeFilter}
            options={typeOptions}
            label={t.exercises.allTypes}
          />
          <FilterSelect
            value={equipmentFilter}
            onChange={setEquipmentFilter}
            options={equipmentOptions}
            label={t.exercises.allEquipment}
          />
          <FilterSelect
            value={muscleFilter}
            onChange={setMuscleFilter}
            options={muscleOptions}
            label={t.exercises.allMuscles}
          />
        </div>
      </div>

      {/* Exercise grid */}
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="w-8 h-8 text-accent animate-spin" />
        </div>
      ) : exercises.length === 0 ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
          <div className="flex flex-col items-center justify-center py-20 gap-4">
            <div className="w-14 h-14 rounded-2xl bg-white/[0.04] flex items-center justify-center">
              <Search className="w-7 h-7 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">{t.exercises.noExercises}</p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {exercises.map((exercise) => (
            <ExerciseCard
              key={exercise.id}
              exercise={
                detailCache[exercise.id]
                  ? { ...exercise, ...detailCache[exercise.id] }
                  : exercise
              }
              locale={locale}
              t={t}
              expanded={expandedId === exercise.id}
              onToggle={() => handleToggle(exercise.id)}
            />
          ))}
        </div>
      )}

      {/* Result count */}
      {!loading && exercises.length > 0 && (
        <p className="text-center text-xs text-fg-muted">
          {exercises.length} {exercises.length === 1 ? "exercise" : "exercises"}
          {hasActiveFilters ? " (filtered)" : ""}
        </p>
      )}
    </div>
  );
}
