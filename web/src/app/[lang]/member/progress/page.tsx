"use client";

import { useEffect, useState, useCallback } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, WeightTrend, NutritionStreak } from "@/types";
import {
  Loader2,
  TrendingUp,
  TrendingDown,
  Scale,
  Flame,
  Activity,
  Target,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Period type
// ---------------------------------------------------------------------------

type Period = 30 | 90 | 180;

// ---------------------------------------------------------------------------
// SVG Weight Trend Chart
// ---------------------------------------------------------------------------

function WeightChart({
  entries,
  noDataLabel,
}: {
  entries: { date: string; weight_kg: number }[];
  noDataLabel: string;
}) {
  if (entries.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-3">
        <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
          <Scale className="w-6 h-6 text-fg-muted" />
        </div>
        <p className="text-sm text-fg-muted">{noDataLabel}</p>
      </div>
    );
  }

  // Chart dimensions
  const width = 600;
  const height = 260;
  const paddingLeft = 52;
  const paddingRight = 20;
  const paddingTop = 20;
  const paddingBottom = 40;

  const chartW = width - paddingLeft - paddingRight;
  const chartH = height - paddingTop - paddingBottom;

  // Data bounds
  const weights = entries.map((e) => e.weight_kg);
  const minW = Math.floor(Math.min(...weights) - 1);
  const maxW = Math.ceil(Math.max(...weights) + 1);
  const rangeW = maxW - minW || 1;

  // Map data to chart coordinates
  const points = entries.map((e, i) => {
    const x = paddingLeft + (i / Math.max(entries.length - 1, 1)) * chartW;
    const y =
      paddingTop + chartH - ((e.weight_kg - minW) / rangeW) * chartH;
    return { x, y, date: e.date, weight: e.weight_kg };
  });

  // Build line path
  const linePath = points
    .map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`)
    .join(" ");

  // Build area path (closed shape for gradient fill)
  const areaPath = `${linePath} L ${points[points.length - 1].x} ${paddingTop + chartH} L ${points[0].x} ${paddingTop + chartH} Z`;

  // Y-axis labels (5 ticks)
  const yTicks = 5;
  const yLabels: { value: number; y: number }[] = [];
  for (let i = 0; i <= yTicks; i++) {
    const value = minW + (rangeW * i) / yTicks;
    const y = paddingTop + chartH - (i / yTicks) * chartH;
    yLabels.push({ value: Math.round(value * 10) / 10, y });
  }

  // X-axis labels (show a subset of dates)
  const maxXLabels = 6;
  const step = Math.max(1, Math.floor(entries.length / maxXLabels));
  const xLabels: { label: string; x: number }[] = [];
  for (let i = 0; i < entries.length; i += step) {
    const d = new Date(entries[i].date);
    const label = `${d.getMonth() + 1}/${d.getDate()}`;
    xLabels.push({ label, x: points[i].x });
  }
  // Always include last entry
  if (
    entries.length > 1 &&
    xLabels[xLabels.length - 1]?.x !== points[points.length - 1].x
  ) {
    const d = new Date(entries[entries.length - 1].date);
    xLabels.push({
      label: `${d.getMonth() + 1}/${d.getDate()}`,
      x: points[points.length - 1].x,
    });
  }

  const gradientId = "weight-area-gradient";

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      className="w-full"
      preserveAspectRatio="xMidYMid meet"
    >
      <defs>
        <linearGradient id={gradientId} x1="0" x2="0" y1="0" y2="1">
          <stop offset="0%" stopColor="#a3e635" stopOpacity="0.3" />
          <stop offset="100%" stopColor="#a3e635" stopOpacity="0.02" />
        </linearGradient>
      </defs>

      {/* Grid lines */}
      {yLabels.map((yl) => (
        <line
          key={yl.value}
          x1={paddingLeft}
          x2={width - paddingRight}
          y1={yl.y}
          y2={yl.y}
          stroke="rgba(255,255,255,0.06)"
          strokeWidth="1"
        />
      ))}

      {/* Y-axis labels */}
      {yLabels.map((yl) => (
        <text
          key={`yl-${yl.value}`}
          x={paddingLeft - 8}
          y={yl.y + 4}
          textAnchor="end"
          className="fill-current text-fg-muted"
          fontSize="11"
        >
          {yl.value}
        </text>
      ))}

      {/* X-axis labels */}
      {xLabels.map((xl, i) => (
        <text
          key={`xl-${i}`}
          x={xl.x}
          y={height - 8}
          textAnchor="middle"
          className="fill-current text-fg-muted"
          fontSize="10"
        >
          {xl.label}
        </text>
      ))}

      {/* Area fill */}
      <path d={areaPath} fill={`url(#${gradientId})`} />

      {/* Line */}
      <path
        d={linePath}
        fill="none"
        stroke="#a3e635"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />

      {/* Data points */}
      {points.map((p, i) => (
        <g key={i}>
          <circle cx={p.x} cy={p.y} r="5" fill="#a3e635" />
          <circle cx={p.x} cy={p.y} r="3" fill="#1a1a2e" stroke="white" strokeWidth="1.5" />
        </g>
      ))}
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Main page component
// ---------------------------------------------------------------------------

export default function MemberProgressPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [loading, setLoading] = useState(true);
  const [weightTrend, setWeightTrend] = useState<WeightTrend | null>(null);
  const [streak, setStreak] = useState<NutritionStreak | null>(null);
  const [period, setPeriod] = useState<Period>(30);

  // Weight log form state
  const [weightInput, setWeightInput] = useState("");
  const [logging, setLogging] = useState(false);
  const [logSuccess, setLogSuccess] = useState(false);

  // -----------------------------------------------------------------------
  // Data fetching
  // -----------------------------------------------------------------------

  const loadWeightTrend = useCallback(async (days: number) => {
    try {
      const data = await api.getWeightTrend(days);
      setWeightTrend(data);
    } catch (err) {
      console.error("Failed to load weight trend:", err);
    }
  }, []);

  useEffect(() => {
    async function init() {
      setLoading(true);
      try {
        await Promise.all([
          loadWeightTrend(period),
          api
            .getNutritionStreak()
            .then((s) => setStreak(s))
            .catch((err) =>
              console.error("Failed to load nutrition streak:", err)
            ),
        ]);
      } finally {
        setLoading(false);
      }
    }
    init();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // When period changes, reload weight trend only
  useEffect(() => {
    if (!loading) {
      loadWeightTrend(period);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [period]);

  // -----------------------------------------------------------------------
  // Log weight handler
  // -----------------------------------------------------------------------

  const handleLogWeight = async () => {
    const kg = parseFloat(weightInput);
    if (isNaN(kg) || kg <= 0) return;
    setLogging(true);
    try {
      await api.quickLogWeight(kg);
      setLogSuccess(true);
      setWeightInput("");
      // Reload trend data
      await loadWeightTrend(period);
      setTimeout(() => setLogSuccess(false), 2000);
    } catch (err) {
      console.error("Failed to log weight:", err);
    } finally {
      setLogging(false);
    }
  };

  // -----------------------------------------------------------------------
  // Period label helper
  // -----------------------------------------------------------------------

  const periodLabel = (p: Period) => {
    switch (p) {
      case 30:
        return t.progress.last30Days;
      case 90:
        return t.progress.last90Days;
      case 180:
        return t.progress.last180Days;
    }
  };

  // -----------------------------------------------------------------------
  // Loading state
  // -----------------------------------------------------------------------

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const summary = weightTrend?.summary;
  const entries = weightTrend?.entries ?? [];
  const changePositive = summary && summary.change_kg > 0;
  const changeNegative = summary && summary.change_kg < 0;

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">{t.progress.title}</h1>

      {/* ----------------------------------------------------------------- */}
      {/* Weight Log Form                                                    */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
            <Scale className="w-4 h-4 text-accent" />
            {t.progress.logWeight}
          </h2>
        </div>
        <div className="p-5">
          <div className="flex gap-3">
            <input
              type="number"
              step="0.1"
              min="0"
              placeholder={t.progress.weightKg}
              value={weightInput}
              onChange={(e) => setWeightInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleLogWeight();
              }}
              className="flex-1 px-4 py-2.5 rounded-xl bg-white/[0.06] border border-white/[0.06] text-fg placeholder:text-fg-muted text-sm focus:outline-none focus:border-accent/50 transition-colors"
            />
            <button
              onClick={handleLogWeight}
              disabled={logging || !weightInput}
              className="px-5 py-2.5 rounded-xl bg-accent text-bg font-medium text-sm disabled:opacity-40 transition-opacity"
            >
              {logging ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : logSuccess ? (
                t.progress.log + " !"
              ) : (
                t.progress.log
              )}
            </button>
          </div>
        </div>
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* Summary Stats Grid                                                 */}
      {/* ----------------------------------------------------------------- */}
      {summary && summary.entries_count > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
          {/* Current Weight */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Scale className="w-4 h-4 text-accent" />
              <span className="text-xs text-fg-muted">{t.progress.current}</span>
            </div>
            <p className="text-xl font-bold text-fg">
              {summary.current_kg.toFixed(1)}
              <span className="text-xs font-normal text-fg-muted ml-1">kg</span>
            </p>
          </div>

          {/* Starting Weight */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Target className="w-4 h-4 text-fg-muted" />
              <span className="text-xs text-fg-muted">{t.progress.start}</span>
            </div>
            <p className="text-xl font-bold text-fg">
              {summary.start_kg.toFixed(1)}
              <span className="text-xs font-normal text-fg-muted ml-1">kg</span>
            </p>
          </div>

          {/* Total Change */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              {changeNegative ? (
                <TrendingDown className="w-4 h-4 text-emerald-400" />
              ) : changePositive ? (
                <TrendingUp className="w-4 h-4 text-red-400" />
              ) : (
                <Activity className="w-4 h-4 text-fg-muted" />
              )}
              <span className="text-xs text-fg-muted">{t.progress.change}</span>
            </div>
            <p
              className={`text-xl font-bold ${
                changeNegative
                  ? "text-emerald-400"
                  : changePositive
                  ? "text-red-400"
                  : "text-fg"
              }`}
            >
              {summary.change_kg > 0 ? "+" : ""}
              {summary.change_kg.toFixed(1)}
              <span className="text-xs font-normal text-fg-muted ml-1">kg</span>
            </p>
          </div>

          {/* BMI */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Activity className="w-4 h-4 text-accent" />
              <span className="text-xs text-fg-muted">{t.progress.bmi}</span>
            </div>
            <p className="text-xl font-bold text-fg">
              {summary.bmi > 0 ? summary.bmi.toFixed(1) : "--"}
            </p>
          </div>

          {/* Entries Count */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-fg-muted" />
              <span className="text-xs text-fg-muted">{t.progress.entries}</span>
            </div>
            <p className="text-xl font-bold text-fg">{summary.entries_count}</p>
          </div>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* Weight Trend Chart                                                 */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="px-5 py-4 border-b border-white/[0.06] flex items-center justify-between">
          <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
            <TrendingUp className="w-4 h-4 text-accent" />
            {t.progress.weightTrend}
          </h2>

          {/* Period Selector */}
          <div className="flex gap-1">
            {([30, 90, 180] as Period[]).map((p) => (
              <button
                key={p}
                onClick={() => setPeriod(p)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  period === p
                    ? "bg-accent text-bg"
                    : "bg-white/[0.06] text-fg-muted hover:bg-white/[0.1]"
                }`}
              >
                {periodLabel(p)}
              </button>
            ))}
          </div>
        </div>
        <div className="p-5">
          <WeightChart
            entries={entries}
            noDataLabel={t.progress.noData}
          />
          {entries.length === 0 && (
            <p className="text-xs text-fg-muted text-center mt-2">
              {t.progress.logFirst}
            </p>
          )}
        </div>
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* Nutrition Streak Card                                              */}
      {/* ----------------------------------------------------------------- */}
      {streak && (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
          <div className="px-5 py-4 border-b border-white/[0.06]">
            <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
              <Flame className="w-4 h-4 text-orange-400" />
              {t.progress.nutritionStreak}
            </h2>
          </div>
          <div className="p-5">
            <div className="grid grid-cols-2 gap-6">
              {/* Current Streak */}
              <div className="flex flex-col items-center gap-2">
                <div className="w-16 h-16 rounded-2xl bg-orange-500/10 flex items-center justify-center">
                  <Flame className="w-8 h-8 text-orange-400" />
                </div>
                <p className="text-2xl font-bold text-fg">
                  {streak.current_streak}
                </p>
                <p className="text-xs text-fg-muted">
                  {t.progress.currentStreak}{" "}
                  <span className="text-fg-muted">
                    {t.progress.streakDays}
                  </span>
                </p>
              </div>

              {/* Longest Streak */}
              <div className="flex flex-col items-center gap-2">
                <div className="w-16 h-16 rounded-2xl bg-accent/10 flex items-center justify-center">
                  <TrendingUp className="w-8 h-8 text-accent" />
                </div>
                <p className="text-2xl font-bold text-fg">
                  {streak.longest_streak}
                </p>
                <p className="text-xs text-fg-muted">
                  {t.progress.longestStreak}{" "}
                  <span className="text-fg-muted">
                    {t.progress.streakDays}
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
