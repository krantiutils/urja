"use client";

import { useEffect, useState, useCallback } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type {
  Locale,
  FoodItem,
  MealTemplate,
  NutritionGoal,
  DailySummary,
  MealSummary,
  WeeklySummaryDay,
  NutritionStreak,
  WaterDailySummary,
  CreateCustomFoodInput,
} from "@/types";
import {
  Loader2,
  ChevronLeft,
  ChevronRight,
  ChevronDown,
  ChevronUp,
  Plus,
  Trash2,
  X,
  Target,
  Utensils,
  Search,
  Zap,
  BookOpen,
  Flame,
  Droplets,
  ScanBarcode,
  CookingPot,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatDateAPI(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function formatDateDisplay(d: Date): string {
  return d.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
}

function isSameDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

function getMondayOfWeek(d: Date): Date {
  const date = new Date(d);
  const day = date.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  date.setDate(date.getDate() + diff);
  date.setHours(0, 0, 0, 0);
  return date;
}

function getDayLabel(dateStr: string): string {
  const d = new Date(dateStr + "T00:00:00");
  return d.toLocaleDateString("en-US", { weekday: "short" });
}

type MealType = "breakfast" | "lunch" | "dinner" | "snack";
const MEAL_TYPES: MealType[] = ["breakfast", "lunch", "dinner", "snack"];

// ---------------------------------------------------------------------------
// Main Component
// ---------------------------------------------------------------------------

export default function MemberNutritionPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const nt = t.memberPages.nutrition;
  const { user } = useAuth();
  const orgId = user?.org_id ?? "";

  // --- State ---
  const [loading, setLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [dailySummary, setDailySummary] = useState<DailySummary | null>(null);
  const [goal, setGoal] = useState<NutritionGoal | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklySummaryDay[]>([]);
  const [templates, setTemplates] = useState<MealTemplate[]>([]);

  // Meal section toggle state
  const [expandedMeals, setExpandedMeals] = useState<Record<MealType, boolean>>(
    { breakfast: true, lunch: true, dinner: true, snack: true }
  );

  // Food search modal
  const [searchModal, setSearchModal] = useState<{
    open: boolean;
    mealType: MealType;
  }>({ open: false, mealType: "breakfast" });
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<FoodItem[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [selectedFood, setSelectedFood] = useState<FoodItem | null>(null);
  const [quantity, setQuantity] = useState("");
  const [logLoading, setLogLoading] = useState(false);

  // Barcode search state (within search modal)
  const [barcodeMode, setBarcodeMode] = useState(false);
  const [barcodeInput, setBarcodeInput] = useState("");
  const [barcodeLoading, setBarcodeLoading] = useState(false);
  const [barcodeError, setBarcodeError] = useState("");

  // Custom food modal
  const [customFoodModal, setCustomFoodModal] = useState(false);
  const [customFoodForm, setCustomFoodForm] = useState<CreateCustomFoodInput>({
    name: "",
    name_ne: "",
    category: "",
    calories_per_100g: 0,
    protein_per_100g: 0,
    carbs_per_100g: 0,
    fat_per_100g: 0,
    fiber_per_100g: 0,
    serving_size_g: 100,
    serving_label: "serving",
    barcode: "",
  });
  const [customFoodSaving, setCustomFoodSaving] = useState(false);
  const [customFoodSuccess, setCustomFoodSuccess] = useState(false);

  // Goal modal
  const [goalModal, setGoalModal] = useState(false);
  const [goalForm, setGoalForm] = useState({
    weight_kg: "",
    height_cm: "",
    age: "",
    gender: "male",
    activity_level: "moderate",
    goal_type: "maintain",
  });
  const [goalSaving, setGoalSaving] = useState(false);
  const [goalSuccess, setGoalSuccess] = useState(false);

  // Template creation modal
  const [templateModal, setTemplateModal] = useState(false);
  const [templateForm, setTemplateForm] = useState({
    name: "",
    meal_type: "breakfast" as MealType,
  });
  const [templateSaving, setTemplateSaving] = useState(false);

  // Nutrition streak
  const [streak, setStreak] = useState<NutritionStreak | null>(null);

  // Water tracking
  const [waterSummary, setWaterSummary] = useState<WaterDailySummary | null>(
    null
  );
  const [waterLoading, setWaterLoading] = useState(false);
  const [waterCustomAmount, setWaterCustomAmount] = useState("");
  const [waterShowEntries, setWaterShowEntries] = useState(false);

  // --- Data loading ---
  const loadDay = useCallback(
    async (date: Date) => {
      try {
        const dateStr = formatDateAPI(date);
        const summary = await api.getDailySummary({
          organization_id: orgId || undefined,
          date: dateStr,
        });
        setDailySummary(summary);
      } catch {
        setDailySummary(null);
      }
    },
    [orgId]
  );

  const loadWeekly = useCallback(
    async (date: Date) => {
      try {
        const monday = getMondayOfWeek(date);
        const res = await api.getWeeklySummary({
          organization_id: orgId || undefined,
          from: formatDateAPI(monday),
        });
        setWeeklyData(res.data ?? []);
      } catch {
        setWeeklyData([]);
      }
    },
    [orgId]
  );

  const loadGoal = useCallback(async () => {
    try {
      const g = await api.getNutritionGoal(orgId || undefined);
      setGoal(g);
    } catch {
      setGoal(null);
    }
  }, [orgId]);

  const loadTemplates = useCallback(async () => {
    try {
      const res = await api.getMyMealTemplates(orgId || undefined);
      setTemplates(res.data ?? []);
    } catch {
      setTemplates([]);
    }
  }, [orgId]);

  const loadStreak = useCallback(async () => {
    try {
      const s = await api.getNutritionStreak();
      setStreak(s);
    } catch {
      setStreak(null);
    }
  }, []);

  const loadWater = useCallback(
    async (date: Date) => {
      try {
        const dateStr = formatDateAPI(date);
        const ws = await api.getWaterSummary(dateStr);
        setWaterSummary(ws);
      } catch {
        setWaterSummary(null);
      }
    },
    []
  );

  // Initial load
  useEffect(() => {
    async function init() {
      setLoading(true);
      await Promise.all([
        loadDay(new Date()),
        loadGoal(),
        loadWeekly(new Date()),
        loadTemplates(),
        loadStreak(),
        loadWater(new Date()),
      ]);
      setLoading(false);
    }
    init();
  }, [loadDay, loadGoal, loadWeekly, loadTemplates, loadStreak, loadWater]);

  // When date changes
  useEffect(() => {
    loadDay(selectedDate);
    loadWeekly(selectedDate);
    loadWater(selectedDate);
  }, [selectedDate, loadDay, loadWeekly, loadWater]);

  // Food search debounce
  useEffect(() => {
    if (!searchModal.open || !searchQuery.trim()) {
      setSearchResults([]);
      return;
    }
    const timer = setTimeout(async () => {
      setSearchLoading(true);
      try {
        const res = await api.searchFoods({
          organization_id: orgId || undefined,
          q: searchQuery.trim(),
          limit: 20,
        });
        setSearchResults(res.data ?? []);
      } catch {
        setSearchResults([]);
      } finally {
        setSearchLoading(false);
      }
    }, 300);
    return () => clearTimeout(timer);
  }, [searchQuery, searchModal.open, orgId]);

  // --- Actions ---
  const navigateDate = (delta: number) => {
    const next = new Date(selectedDate);
    next.setDate(next.getDate() + delta);
    setSelectedDate(next);
  };

  const goToToday = () => setSelectedDate(new Date());

  const toggleMeal = (m: MealType) =>
    setExpandedMeals((prev) => ({ ...prev, [m]: !prev[m] }));

  const openSearch = (mealType: MealType) => {
    setSearchModal({ open: true, mealType });
    setSearchQuery("");
    setSearchResults([]);
    setSelectedFood(null);
    setQuantity("");
    setBarcodeMode(false);
    setBarcodeInput("");
    setBarcodeError("");
  };

  const closeSearch = () => {
    setSearchModal({ open: false, mealType: "breakfast" });
    setSelectedFood(null);
    setQuantity("");
    setBarcodeMode(false);
    setBarcodeInput("");
    setBarcodeError("");
  };

  const handleLogFood = async () => {
    if (!selectedFood) return;
    const qty = parseFloat(quantity);
    if (isNaN(qty) || qty <= 0) return;
    setLogLoading(true);
    try {
      await api.logFood({
        organization_id: orgId || undefined,
        food_item_id: selectedFood.id,
        meal_type: searchModal.mealType,
        quantity_grams: qty,
        logged_date: formatDateAPI(selectedDate),
      });
      closeSearch();
      await loadDay(selectedDate);
      await loadWeekly(selectedDate);
      await loadStreak();
    } catch (err) {
      console.error("Failed to log food:", err);
    } finally {
      setLogLoading(false);
    }
  };

  const handleDeleteLog = async (logId: string) => {
    try {
      await api.deleteFoodLog(logId);
      await loadDay(selectedDate);
      await loadWeekly(selectedDate);
    } catch (err) {
      console.error("Failed to delete log:", err);
    }
  };

  const handleSetGoal = async () => {
    setGoalSaving(true);
    setGoalSuccess(false);
    try {
      const g = await api.setNutritionGoal({
        organization_id: orgId || undefined,
        weight_kg: parseFloat(goalForm.weight_kg),
        height_cm: parseFloat(goalForm.height_cm),
        age: parseInt(goalForm.age, 10),
        gender: goalForm.gender,
        activity_level: goalForm.activity_level,
        goal_type: goalForm.goal_type,
      });
      setGoal(g);
      setGoalSuccess(true);
      setTimeout(() => {
        setGoalModal(false);
        setGoalSuccess(false);
      }, 1500);
    } catch (err) {
      console.error("Failed to set goal:", err);
    } finally {
      setGoalSaving(false);
    }
  };

  const handleQuickLog = async (templateId: string) => {
    try {
      await api.logMealTemplate(templateId, {
        organization_id: orgId || undefined,
        logged_date: formatDateAPI(selectedDate),
      });
      await loadDay(selectedDate);
      await loadWeekly(selectedDate);
      await loadStreak();
    } catch (err) {
      console.error("Failed to quick log template:", err);
    }
  };

  const handleCreateTemplate = async () => {
    if (!templateForm.name.trim()) return;
    setTemplateSaving(true);
    try {
      await api.createMealTemplate({
        organization_id: orgId || undefined,
        name: templateForm.name.trim(),
        meal_type: templateForm.meal_type,
        items: [],
      });
      await loadTemplates();
      setTemplateModal(false);
      setTemplateForm({ name: "", meal_type: "breakfast" });
    } catch (err) {
      console.error("Failed to create template:", err);
    } finally {
      setTemplateSaving(false);
    }
  };

  const handleDeleteTemplate = async (id: string) => {
    try {
      await api.deleteMealTemplate(id);
      await loadTemplates();
    } catch (err) {
      console.error("Failed to delete template:", err);
    }
  };

  // Barcode lookup
  const handleBarcodeLookup = async () => {
    if (!barcodeInput.trim()) return;
    setBarcodeLoading(true);
    setBarcodeError("");
    try {
      const food = await api.lookupBarcode(barcodeInput.trim());
      setSelectedFood(food);
      setQuantity(String(food.serving_size_g));
      setBarcodeMode(false);
      setBarcodeInput("");
    } catch {
      setBarcodeError(t.customFood.barcodeNotFound);
    } finally {
      setBarcodeLoading(false);
    }
  };

  // Custom food creation
  const handleCreateCustomFood = async () => {
    if (!customFoodForm.name.trim()) return;
    setCustomFoodSaving(true);
    setCustomFoodSuccess(false);
    try {
      const created = await api.createCustomFood({
        name: customFoodForm.name.trim(),
        name_ne: customFoodForm.name_ne || undefined,
        category: customFoodForm.category || undefined,
        calories_per_100g: customFoodForm.calories_per_100g,
        protein_per_100g: customFoodForm.protein_per_100g,
        carbs_per_100g: customFoodForm.carbs_per_100g,
        fat_per_100g: customFoodForm.fat_per_100g,
        fiber_per_100g: customFoodForm.fiber_per_100g || undefined,
        serving_size_g: customFoodForm.serving_size_g || undefined,
        serving_label: customFoodForm.serving_label || undefined,
        barcode: customFoodForm.barcode || undefined,
      });
      setCustomFoodSuccess(true);
      // If search modal is open, auto-select the newly created food
      if (searchModal.open) {
        setSelectedFood(created);
        setQuantity(String(created.serving_size_g));
      }
      setTimeout(() => {
        setCustomFoodModal(false);
        setCustomFoodSuccess(false);
        setCustomFoodForm({
          name: "",
          name_ne: "",
          category: "",
          calories_per_100g: 0,
          protein_per_100g: 0,
          carbs_per_100g: 0,
          fat_per_100g: 0,
          fiber_per_100g: 0,
          serving_size_g: 100,
          serving_label: "serving",
          barcode: "",
        });
      }, 1200);
    } catch (err) {
      console.error("Failed to create custom food:", err);
    } finally {
      setCustomFoodSaving(false);
    }
  };

  // Water tracking actions
  const handleLogWater = async (amountMl: number) => {
    if (amountMl <= 0) return;
    setWaterLoading(true);
    try {
      await api.logWater({
        amount_ml: amountMl,
        date: formatDateAPI(selectedDate),
      });
      await loadWater(selectedDate);
    } catch (err) {
      console.error("Failed to log water:", err);
    } finally {
      setWaterLoading(false);
    }
  };

  const handleDeleteWater = async (id: string) => {
    try {
      await api.deleteWaterLog(id);
      await loadWater(selectedDate);
    } catch (err) {
      console.error("Failed to delete water log:", err);
    }
  };

  const handleCustomWater = () => {
    const amount = parseInt(waterCustomAmount, 10);
    if (!isNaN(amount) && amount > 0) {
      handleLogWater(amount);
      setWaterCustomAmount("");
    }
  };

  // --- Derived values ---
  const consumed = dailySummary?.total_calories ?? 0;
  const calorieGoal = goal?.calorie_goal ?? 0;
  const remaining = Math.max(0, calorieGoal - consumed);
  const proteinConsumed = dailySummary?.total_protein ?? 0;
  const carbsConsumed = dailySummary?.total_carbs ?? 0;
  const fatConsumed = dailySummary?.total_fat ?? 0;
  const proteinGoal = goal?.protein_goal_g ?? 0;
  const carbsGoal = goal?.carbs_goal_g ?? 0;
  const fatGoal = goal?.fat_goal_g ?? 0;

  const meals = dailySummary?.meals ?? [];
  const getMealSummary = (type: MealType): MealSummary | undefined =>
    meals.find((m) => m.meal_type === type);

  // Calorie ring
  const ringProgress = calorieGoal > 0 ? Math.min(consumed / calorieGoal, 1.5) : 0;
  const ringColor =
    ringProgress > 1
      ? "text-red-500"
      : ringProgress > 0.85
        ? "text-amber-500"
        : "text-emerald-500";
  const ringStrokeColor =
    ringProgress > 1
      ? "#ef4444"
      : ringProgress > 0.85
        ? "#f59e0b"
        : "#10b981";

  // SVG calorie ring constants
  const ringRadius = 80;
  const ringCircumference = 2 * Math.PI * ringRadius;
  const ringDash = Math.min(ringProgress, 1) * ringCircumference;

  // Weekly chart
  const weeklyMax = Math.max(
    calorieGoal,
    ...weeklyData.map((d) => d.total_calories),
    1
  );

  // Water derived
  const waterConsumed = waterSummary?.total_ml ?? 0;
  const waterGoal = waterSummary?.goal_ml ?? 2500;
  const waterPercent = waterGoal > 0 ? Math.min((waterConsumed / waterGoal) * 100, 100) : 0;
  const waterRemaining = Math.max(0, waterGoal - waterConsumed);
  const waterEntries = waterSummary?.entries ?? [];

  // --- Loading state ---
  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  // --- Render ---
  return (
    <div className="space-y-6 pb-8">
      {/* Page Title + Streak Badge */}
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-semibold text-fg">{nt.title}</h1>
        {streak && streak.current_streak > 0 && (
          <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-orange-500/20 border border-orange-500/30">
            <Flame className="w-4 h-4 text-orange-400" />
            <span className="text-sm font-bold text-orange-400">
              {streak.current_streak}
            </span>
            <span className="text-xs text-orange-400/70">
              {t.progress.streakDays}
            </span>
          </div>
        )}
      </div>

      {/* Streak Detail (if streak exists) */}
      {streak && (
        <div className="flex items-center gap-4 text-xs text-fg-muted">
          <span>
            {t.progress.currentStreak}: <span className="text-fg font-semibold">{streak.current_streak}</span> {t.progress.streakDays}
          </span>
          <span className="text-white/[0.2]">|</span>
          <span>
            {t.progress.longestStreak}: <span className="text-fg font-semibold">{streak.longest_streak}</span> {t.progress.streakDays}
          </span>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* 1. Date Navigation Bar                                            */}
      {/* ----------------------------------------------------------------- */}
      <div className="flex items-center justify-center gap-3">
        <button
          onClick={() => navigateDate(-1)}
          className="p-2 rounded-lg bg-white/[0.06] hover:bg-white/[0.1] transition-colors"
        >
          <ChevronLeft className="w-5 h-5 text-fg-muted" />
        </button>
        <span className="text-fg font-medium min-w-[140px] text-center">
          {formatDateDisplay(selectedDate)}
        </span>
        <button
          onClick={() => navigateDate(1)}
          className="p-2 rounded-lg bg-white/[0.06] hover:bg-white/[0.1] transition-colors"
        >
          <ChevronRight className="w-5 h-5 text-fg-muted" />
        </button>
        {!isSameDay(selectedDate, new Date()) && (
          <button
            onClick={goToToday}
            className="ml-2 px-3 py-1 text-xs rounded-full bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30 transition-colors font-medium"
          >
            Today
          </button>
        )}
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* 2. Calorie Ring                                                   */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-6 flex flex-col items-center">
        <p className="text-sm text-fg-muted mb-4 font-medium">
          {nt.dailyIntake}
        </p>
        {calorieGoal > 0 ? (
          <>
            <div className="relative w-[200px] h-[200px]">
              <svg viewBox="0 0 200 200" className="w-full h-full -rotate-90">
                {/* Track */}
                <circle
                  cx="100"
                  cy="100"
                  r={ringRadius}
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="12"
                  className="text-white/[0.06]"
                />
                {/* Progress */}
                <circle
                  cx="100"
                  cy="100"
                  r={ringRadius}
                  fill="none"
                  stroke={ringStrokeColor}
                  strokeWidth="12"
                  strokeLinecap="round"
                  strokeDasharray={`${ringDash} ${ringCircumference}`}
                  className="transition-all duration-500"
                />
              </svg>
              <div className="absolute inset-0 flex flex-col items-center justify-center">
                <span className={`text-2xl font-bold ${ringColor}`}>
                  {Math.round(consumed).toLocaleString()}
                </span>
                <span className="text-xs text-fg-muted">
                  / {calorieGoal.toLocaleString()} {nt.kcal}
                </span>
              </div>
            </div>
            <p className="mt-3 text-sm text-fg-muted">
              <span className={`font-semibold ${ringColor}`}>
                {Math.round(remaining).toLocaleString()}
              </span>{" "}
              {nt.caloriesRemaining}
            </p>
          </>
        ) : (
          <div className="flex flex-col items-center gap-3 py-6">
            <Target className="w-10 h-10 text-fg-muted" />
            <p className="text-sm text-fg-muted text-center">{nt.noGoalSet}</p>
            <button
              onClick={() => setGoalModal(true)}
              className="px-4 py-2 text-sm font-medium rounded-lg bg-emerald-600 text-white hover:bg-emerald-500 transition-colors"
            >
              {nt.setGoal}
            </button>
          </div>
        )}
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* 3. Macro Progress Bars                                            */}
      {/* ----------------------------------------------------------------- */}
      {calorieGoal > 0 && (
        <div className="grid grid-cols-3 gap-3">
          {/* Protein */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-xl p-3">
            <p className="text-xs text-fg-muted font-medium mb-1">
              {nt.protein}
            </p>
            <p className="text-sm font-semibold text-fg mb-2">
              {Math.round(proteinConsumed)}/{proteinGoal}
              {nt.grams}
            </p>
            <div className="h-2 rounded-full bg-white/[0.06] overflow-hidden">
              <div
                className="h-full rounded-full bg-blue-500 transition-all duration-500"
                style={{
                  width: `${Math.min(100, proteinGoal > 0 ? (proteinConsumed / proteinGoal) * 100 : 0)}%`,
                }}
              />
            </div>
          </div>

          {/* Carbs */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-xl p-3">
            <p className="text-xs text-fg-muted font-medium mb-1">
              {nt.carbs}
            </p>
            <p className="text-sm font-semibold text-fg mb-2">
              {Math.round(carbsConsumed)}/{carbsGoal}
              {nt.grams}
            </p>
            <div className="h-2 rounded-full bg-white/[0.06] overflow-hidden">
              <div
                className="h-full rounded-full bg-amber-500 transition-all duration-500"
                style={{
                  width: `${Math.min(100, carbsGoal > 0 ? (carbsConsumed / carbsGoal) * 100 : 0)}%`,
                }}
              />
            </div>
          </div>

          {/* Fat */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-xl p-3">
            <p className="text-xs text-fg-muted font-medium mb-1">{nt.fat}</p>
            <p className="text-sm font-semibold text-fg mb-2">
              {Math.round(fatConsumed)}/{fatGoal}
              {nt.grams}
            </p>
            <div className="h-2 rounded-full bg-white/[0.06] overflow-hidden">
              <div
                className="h-full rounded-full bg-rose-500 transition-all duration-500"
                style={{
                  width: `${Math.min(100, fatGoal > 0 ? (fatConsumed / fatGoal) * 100 : 0)}%`,
                }}
              />
            </div>
          </div>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* 4. Water Tracker Section                                          */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
            <Droplets className="w-4 h-4 text-sky-400" />
            {t.water.title}
          </h2>
          {waterPercent >= 100 ? (
            <span className="text-xs font-medium text-sky-400 bg-sky-500/20 px-2.5 py-1 rounded-full">
              {t.water.completed}
            </span>
          ) : (
            <span className="text-xs text-fg-muted">
              {waterRemaining}{t.water.amount?.replace(/[^a-zA-Z.()]/g, "") || "ml"} {t.water.remaining}
            </span>
          )}
        </div>

        {/* Water progress bar */}
        <div className="mb-4">
          <div className="flex items-center justify-between text-xs text-fg-muted mb-1.5">
            <span>{waterConsumed}ml</span>
            <span>{waterGoal}ml</span>
          </div>
          <div className="h-3 rounded-full bg-white/[0.06] overflow-hidden">
            <div
              className="h-full rounded-full bg-gradient-to-r from-sky-500 to-sky-400 transition-all duration-500"
              style={{ width: `${waterPercent}%` }}
            />
          </div>
        </div>

        {/* Quick add buttons */}
        <div className="flex flex-wrap gap-2 mb-3">
          <button
            onClick={() => handleLogWater(250)}
            disabled={waterLoading}
            className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-sky-500/20 text-sky-400 hover:bg-sky-500/30 transition-colors disabled:opacity-50"
          >
            <Plus className="w-3 h-3" />
            {t.water.glassSize}
          </button>
          <button
            onClick={() => handleLogWater(500)}
            disabled={waterLoading}
            className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-sky-500/20 text-sky-400 hover:bg-sky-500/30 transition-colors disabled:opacity-50"
          >
            <Plus className="w-3 h-3" />
            {t.water.bottleSize}
          </button>
          <button
            onClick={() => handleLogWater(1000)}
            disabled={waterLoading}
            className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-sky-500/20 text-sky-400 hover:bg-sky-500/30 transition-colors disabled:opacity-50"
          >
            <Plus className="w-3 h-3" />
            {t.water.largeBottle}
          </button>
        </div>

        {/* Custom amount input */}
        <div className="flex gap-2 mb-3">
          <input
            type="number"
            value={waterCustomAmount}
            onChange={(e) => setWaterCustomAmount(e.target.value)}
            placeholder={t.water.customAmount}
            className="flex-1 px-3 py-2 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-sky-500/40"
          />
          <button
            onClick={handleCustomWater}
            disabled={waterLoading || !waterCustomAmount}
            className="px-4 py-2 text-sm font-medium rounded-lg bg-sky-600 text-white hover:bg-sky-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {waterLoading ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              t.water.addWater
            )}
          </button>
        </div>

        {/* Toggle entries list */}
        {waterEntries.length > 0 && (
          <>
            <button
              onClick={() => setWaterShowEntries(!waterShowEntries)}
              className="flex items-center gap-1 text-xs text-fg-muted hover:text-fg transition-colors font-medium"
            >
              {waterShowEntries ? (
                <ChevronUp className="w-3.5 h-3.5" />
              ) : (
                <ChevronDown className="w-3.5 h-3.5" />
              )}
              {waterEntries.length} entries
            </button>
            {waterShowEntries && (
              <ul className="mt-2 divide-y divide-white/[0.04]">
                {waterEntries.map((entry) => (
                  <li
                    key={entry.id}
                    className="flex items-center justify-between py-2"
                  >
                    <div className="flex items-center gap-2">
                      <Droplets className="w-3 h-3 text-sky-400" />
                      <span className="text-sm text-fg">{entry.amount_ml}ml</span>
                      <span className="text-xs text-fg-muted">
                        {new Date(entry.logged_at).toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </span>
                    </div>
                    <button
                      onClick={() => handleDeleteWater(entry.id)}
                      className="p-1 rounded hover:bg-white/[0.06] text-fg-muted hover:text-red-400 transition-colors"
                      title={t.water.delete}
                    >
                      <Trash2 className="w-3 h-3" />
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </>
        )}
        {waterEntries.length === 0 && (
          <p className="text-xs text-fg-muted">{t.water.noLogs}</p>
        )}
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* 5. Meal Sections                                                  */}
      {/* ----------------------------------------------------------------- */}
      <div className="space-y-3">
        {MEAL_TYPES.map((mealType) => {
          const mealSummary = getMealSummary(mealType);
          const items = mealSummary?.items ?? [];
          const mealCals = mealSummary?.calories ?? 0;
          const expanded = expandedMeals[mealType];
          const mealLabel =
            nt[mealType as keyof typeof nt] as string;

          return (
            <div
              key={mealType}
              className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden"
            >
              {/* Header */}
              <button
                onClick={() => toggleMeal(mealType)}
                className="w-full flex items-center justify-between px-5 py-4 hover:bg-white/[0.02] transition-colors"
              >
                <div className="flex items-center gap-3">
                  <Utensils className="w-4 h-4 text-fg-muted" />
                  <span className="text-sm font-semibold text-fg">
                    {mealLabel}
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-fg-muted font-medium">
                    {Math.round(mealCals)} {nt.kcal}
                  </span>
                  {expanded ? (
                    <ChevronUp className="w-4 h-4 text-fg-muted" />
                  ) : (
                    <ChevronDown className="w-4 h-4 text-fg-muted" />
                  )}
                </div>
              </button>

              {/* Items */}
              {expanded && (
                <div className="border-t border-white/[0.06]">
                  {items.length === 0 ? (
                    <p className="px-5 py-4 text-xs text-fg-muted">
                      {nt.noLogs}
                    </p>
                  ) : (
                    <ul className="divide-y divide-white/[0.04]">
                      {items.map((item) => (
                        <li
                          key={item.id}
                          className="flex items-center justify-between px-5 py-3"
                        >
                          <div className="flex-1 min-w-0">
                            <p className="text-sm text-fg truncate">
                              {item.food_item?.name ?? "Food"}
                            </p>
                            <p className="text-xs text-fg-muted">
                              {item.quantity_grams}
                              {nt.grams}
                            </p>
                          </div>
                          <div className="flex items-center gap-3">
                            <span className="text-xs text-fg-muted font-medium">
                              {Math.round(item.calories)} {nt.kcal}
                            </span>
                            <button
                              onClick={() => handleDeleteLog(item.id)}
                              className="p-1 rounded hover:bg-white/[0.06] text-fg-muted hover:text-red-400 transition-colors"
                              title={nt.deleteLog}
                            >
                              <Trash2 className="w-3.5 h-3.5" />
                            </button>
                          </div>
                        </li>
                      ))}
                    </ul>
                  )}

                  {/* Add food button */}
                  <div className="px-5 py-3 border-t border-white/[0.04]">
                    <button
                      onClick={() => openSearch(mealType)}
                      className="flex items-center gap-2 text-xs text-emerald-400 hover:text-emerald-300 transition-colors font-medium"
                    >
                      <Plus className="w-3.5 h-3.5" />
                      {nt.addFood}
                    </button>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* 6. Food Search Modal (with Barcode + Create Custom Food buttons)  */}
      {/* ----------------------------------------------------------------- */}
      {searchModal.open && (
        <div
          className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={(e) => {
            if (e.target === e.currentTarget) closeSearch();
          }}
        >
          <div className="w-full sm:max-w-md bg-bg border border-white/[0.1] rounded-t-2xl sm:rounded-2xl max-h-[85vh] flex flex-col shadow-2xl">
            {/* Modal header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
              <h3 className="text-sm font-semibold text-fg">{nt.logFood}</h3>
              <button
                onClick={closeSearch}
                className="p-1 rounded-lg hover:bg-white/[0.06] transition-colors"
              >
                <X className="w-5 h-5 text-fg-muted" />
              </button>
            </div>

            {/* Search input + action buttons */}
            <div className="px-5 py-3 border-b border-white/[0.06] space-y-3">
              {!barcodeMode ? (
                <>
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                    <input
                      type="text"
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      placeholder={nt.searchFood}
                      className="w-full pl-10 pr-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                      autoFocus
                    />
                  </div>
                  {/* Barcode + Create Custom buttons */}
                  <div className="flex gap-2">
                    <button
                      onClick={() => {
                        setBarcodeMode(true);
                        setBarcodeError("");
                      }}
                      className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-lg bg-white/[0.06] hover:bg-white/[0.1] text-fg-muted hover:text-fg transition-colors"
                    >
                      <ScanBarcode className="w-3.5 h-3.5" />
                      {t.customFood.scanBarcode}
                    </button>
                    <button
                      onClick={() => setCustomFoodModal(true)}
                      className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-lg bg-white/[0.06] hover:bg-white/[0.1] text-fg-muted hover:text-fg transition-colors"
                    >
                      <CookingPot className="w-3.5 h-3.5" />
                      {t.customFood.title}
                    </button>
                  </div>
                </>
              ) : (
                /* Barcode input mode */
                <div className="space-y-2">
                  <div className="flex gap-2">
                    <div className="relative flex-1">
                      <ScanBarcode className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                      <input
                        type="text"
                        value={barcodeInput}
                        onChange={(e) => setBarcodeInput(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") handleBarcodeLookup();
                        }}
                        placeholder={t.customFood.barcode}
                        className="w-full pl-10 pr-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                        autoFocus
                      />
                    </div>
                    <button
                      onClick={handleBarcodeLookup}
                      disabled={barcodeLoading || !barcodeInput.trim()}
                      className="px-4 py-2.5 text-sm font-medium rounded-lg bg-emerald-600 text-white hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      {barcodeLoading ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : (
                        <Search className="w-4 h-4" />
                      )}
                    </button>
                  </div>
                  {barcodeError && (
                    <p className="text-xs text-red-400">{barcodeError}</p>
                  )}
                  <button
                    onClick={() => {
                      setBarcodeMode(false);
                      setBarcodeInput("");
                      setBarcodeError("");
                    }}
                    className="text-xs text-fg-muted hover:text-fg transition-colors"
                  >
                    &larr; Back to search
                  </button>
                </div>
              )}
            </div>

            {/* Results or selected food */}
            <div className="flex-1 overflow-y-auto">
              {selectedFood ? (
                <div className="p-5 space-y-4">
                  <div>
                    <p className="text-sm font-semibold text-fg">
                      {selectedFood.name}
                    </p>
                    <p className="text-xs text-fg-muted mt-1">
                      {selectedFood.calories_per_100g} {nt.kcal} {nt.per100g} |{" "}
                      {nt.servingSize}: {selectedFood.serving_size_g}
                      {nt.grams} ({selectedFood.serving_label})
                    </p>
                  </div>
                  <div>
                    <label className="text-xs text-fg-muted font-medium">
                      {nt.quantity}
                    </label>
                    <input
                      type="number"
                      value={quantity}
                      onChange={(e) => setQuantity(e.target.value)}
                      placeholder={String(selectedFood.serving_size_g)}
                      className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    />
                  </div>
                  {quantity && parseFloat(quantity) > 0 && (
                    <div className="grid grid-cols-4 gap-2 text-center">
                      <div className="bg-white/[0.04] rounded-lg p-2">
                        <p className="text-xs text-fg-muted">{nt.calories}</p>
                        <p className="text-sm font-semibold text-emerald-400">
                          {Math.round(
                            (selectedFood.calories_per_100g *
                              parseFloat(quantity)) /
                              100
                          )}
                        </p>
                      </div>
                      <div className="bg-white/[0.04] rounded-lg p-2">
                        <p className="text-xs text-fg-muted">{nt.protein}</p>
                        <p className="text-sm font-semibold text-blue-400">
                          {Math.round(
                            (selectedFood.protein_per_100g *
                              parseFloat(quantity)) /
                              100
                          )}
                          {nt.grams}
                        </p>
                      </div>
                      <div className="bg-white/[0.04] rounded-lg p-2">
                        <p className="text-xs text-fg-muted">{nt.carbs}</p>
                        <p className="text-sm font-semibold text-amber-400">
                          {Math.round(
                            (selectedFood.carbs_per_100g *
                              parseFloat(quantity)) /
                              100
                          )}
                          {nt.grams}
                        </p>
                      </div>
                      <div className="bg-white/[0.04] rounded-lg p-2">
                        <p className="text-xs text-fg-muted">{nt.fat}</p>
                        <p className="text-sm font-semibold text-rose-400">
                          {Math.round(
                            (selectedFood.fat_per_100g *
                              parseFloat(quantity)) /
                              100
                          )}
                          {nt.grams}
                        </p>
                      </div>
                    </div>
                  )}
                  <div className="flex gap-3">
                    <button
                      onClick={() => setSelectedFood(null)}
                      className="flex-1 px-4 py-2.5 text-sm rounded-lg border border-white/[0.1] text-fg-muted hover:bg-white/[0.04] transition-colors"
                    >
                      {t.common.cancel}
                    </button>
                    <button
                      onClick={handleLogFood}
                      disabled={
                        logLoading ||
                        !quantity ||
                        parseFloat(quantity) <= 0
                      }
                      className="flex-1 px-4 py-2.5 text-sm rounded-lg bg-emerald-600 text-white font-medium hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      {logLoading ? (
                        <Loader2 className="w-4 h-4 animate-spin mx-auto" />
                      ) : (
                        nt.logFood
                      )}
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  {searchLoading ? (
                    <div className="flex justify-center py-8">
                      <Loader2 className="w-5 h-5 animate-spin text-fg-muted" />
                    </div>
                  ) : searchResults.length === 0 && searchQuery.trim() ? (
                    <p className="text-center text-sm text-fg-muted py-8">
                      {t.common.noResults}
                    </p>
                  ) : (
                    <ul className="divide-y divide-white/[0.04]">
                      {searchResults.map((food) => (
                        <li key={food.id}>
                          <button
                            onClick={() => {
                              setSelectedFood(food);
                              setQuantity(String(food.serving_size_g));
                            }}
                            className="w-full px-5 py-3 text-left hover:bg-white/[0.04] transition-colors"
                          >
                            <p className="text-sm text-fg">{food.name}</p>
                            <p className="text-xs text-fg-muted mt-0.5">
                              {food.calories_per_100g} {nt.kcal} {nt.per100g} |{" "}
                              {food.category}
                            </p>
                          </button>
                        </li>
                      ))}
                    </ul>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* 7. Custom Food Creation Modal                                     */}
      {/* ----------------------------------------------------------------- */}
      {customFoodModal && (
        <div
          className="fixed inset-0 z-[60] flex items-end sm:items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setCustomFoodModal(false);
              setCustomFoodSuccess(false);
            }
          }}
        >
          <div className="w-full sm:max-w-md bg-bg border border-white/[0.1] rounded-t-2xl sm:rounded-2xl max-h-[85vh] overflow-y-auto shadow-2xl">
            {/* Header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
              <h3 className="text-sm font-semibold text-fg">
                {t.customFood.title}
              </h3>
              <button
                onClick={() => {
                  setCustomFoodModal(false);
                  setCustomFoodSuccess(false);
                }}
                className="p-1 rounded-lg hover:bg-white/[0.06] transition-colors"
              >
                <X className="w-5 h-5 text-fg-muted" />
              </button>
            </div>

            {customFoodSuccess ? (
              <div className="p-8 flex flex-col items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-emerald-500/20 flex items-center justify-center">
                  <CookingPot className="w-6 h-6 text-emerald-400" />
                </div>
                <p className="text-sm font-medium text-emerald-400">
                  Food created!
                </p>
              </div>
            ) : (
              <div className="p-5 space-y-4">
                {/* Name */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {t.customFood.name} *
                  </label>
                  <input
                    type="text"
                    value={customFoodForm.name}
                    onChange={(e) =>
                      setCustomFoodForm((p) => ({ ...p, name: e.target.value }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    autoFocus
                  />
                </div>

                {/* Category */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {t.customFood.category}
                  </label>
                  <input
                    type="text"
                    value={customFoodForm.category}
                    onChange={(e) =>
                      setCustomFoodForm((p) => ({
                        ...p,
                        category: e.target.value,
                      }))
                    }
                    placeholder="e.g. grain, protein, vegetable"
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  />
                </div>

                {/* Nutrition values grid */}
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="text-xs text-fg-muted font-medium">
                      {t.customFood.caloriesPer100g} *
                    </label>
                    <input
                      type="number"
                      value={customFoodForm.calories_per_100g || ""}
                      onChange={(e) =>
                        setCustomFoodForm((p) => ({
                          ...p,
                          calories_per_100g: parseFloat(e.target.value) || 0,
                        }))
                      }
                      className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-fg-muted font-medium">
                      {t.customFood.proteinPer100g} *
                    </label>
                    <input
                      type="number"
                      value={customFoodForm.protein_per_100g || ""}
                      onChange={(e) =>
                        setCustomFoodForm((p) => ({
                          ...p,
                          protein_per_100g: parseFloat(e.target.value) || 0,
                        }))
                      }
                      className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-fg-muted font-medium">
                      {t.customFood.carbsPer100g} *
                    </label>
                    <input
                      type="number"
                      value={customFoodForm.carbs_per_100g || ""}
                      onChange={(e) =>
                        setCustomFoodForm((p) => ({
                          ...p,
                          carbs_per_100g: parseFloat(e.target.value) || 0,
                        }))
                      }
                      className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-fg-muted font-medium">
                      {t.customFood.fatPer100g} *
                    </label>
                    <input
                      type="number"
                      value={customFoodForm.fat_per_100g || ""}
                      onChange={(e) =>
                        setCustomFoodForm((p) => ({
                          ...p,
                          fat_per_100g: parseFloat(e.target.value) || 0,
                        }))
                      }
                      className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                    />
                  </div>
                </div>

                {/* Barcode */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {t.customFood.barcode}
                  </label>
                  <input
                    type="text"
                    value={customFoodForm.barcode}
                    onChange={(e) =>
                      setCustomFoodForm((p) => ({
                        ...p,
                        barcode: e.target.value,
                      }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  />
                </div>

                {/* Submit */}
                <button
                  onClick={handleCreateCustomFood}
                  disabled={
                    customFoodSaving ||
                    !customFoodForm.name.trim() ||
                    customFoodForm.calories_per_100g <= 0
                  }
                  className="w-full px-4 py-2.5 text-sm rounded-lg bg-emerald-600 text-white font-medium hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {customFoodSaving ? (
                    <Loader2 className="w-4 h-4 animate-spin mx-auto" />
                  ) : (
                    t.customFood.create
                  )}
                </button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* 8. Weekly Progress Chart                                          */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
        <h2 className="text-sm font-semibold text-fg mb-4">
          {nt.weeklyProgress}
        </h2>
        {weeklyData.length > 0 ? (
          <div className="relative">
            {/* Goal line */}
            {calorieGoal > 0 && (
              <div
                className="absolute left-0 right-0 border-t-2 border-dashed border-emerald-500/40 pointer-events-none"
                style={{
                  bottom: `${(calorieGoal / weeklyMax) * 100}%`,
                }}
              >
                <span className="absolute -top-4 right-0 text-[10px] text-emerald-400 font-medium">
                  {nt.calorieGoal}
                </span>
              </div>
            )}
            <div className="flex items-end justify-between gap-2 h-40">
              {weeklyData.map((day) => {
                const pct =
                  weeklyMax > 0
                    ? (day.total_calories / weeklyMax) * 100
                    : 0;
                const isOver =
                  calorieGoal > 0 && day.total_calories > calorieGoal;
                return (
                  <div
                    key={day.date}
                    className="flex-1 flex flex-col items-center gap-1"
                  >
                    <span className="text-[10px] text-fg-muted font-medium">
                      {day.total_calories > 0
                        ? Math.round(day.total_calories)
                        : ""}
                    </span>
                    <div className="w-full flex justify-center">
                      <div
                        className={`w-full max-w-[32px] rounded-t-md transition-all duration-300 ${
                          isOver ? "bg-red-500/70" : "bg-emerald-500/70"
                        }`}
                        style={{
                          height: `${Math.max(pct, 2)}%`,
                          minHeight: day.total_calories > 0 ? "8px" : "2px",
                        }}
                      />
                    </div>
                    <span className="text-[10px] text-fg-muted">
                      {getDayLabel(day.date)}
                    </span>
                  </div>
                );
              })}
            </div>
          </div>
        ) : (
          <p className="text-sm text-fg-muted text-center py-6">{nt.noLogs}</p>
        )}
      </div>

      {/* ----------------------------------------------------------------- */}
      {/* 9. TDEE Goal Setup                                                */}
      {/* ----------------------------------------------------------------- */}
      {goal ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
              <Target className="w-4 h-4 text-emerald-400" />
              {nt.dailyCalories}
            </h2>
            <button
              onClick={() => {
                setGoalForm({
                  weight_kg: String(goal.weight_kg),
                  height_cm: String(goal.height_cm),
                  age: String(goal.age),
                  gender: goal.gender,
                  activity_level: goal.activity_level,
                  goal_type: goal.goal_type,
                });
                setGoalModal(true);
              }}
              className="text-xs text-emerald-400 hover:text-emerald-300 transition-colors font-medium"
            >
              {nt.editGoal}
            </button>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 text-sm">
            <div>
              <p className="text-xs text-fg-muted">{nt.calories}</p>
              <p className="font-semibold text-fg">
                {goal.calorie_goal} {nt.kcal}
              </p>
            </div>
            <div>
              <p className="text-xs text-fg-muted">{nt.protein}</p>
              <p className="font-semibold text-blue-400">
                {goal.protein_goal_g}
                {nt.grams}
              </p>
            </div>
            <div>
              <p className="text-xs text-fg-muted">{nt.carbs}</p>
              <p className="font-semibold text-amber-400">
                {goal.carbs_goal_g}
                {nt.grams}
              </p>
            </div>
            <div>
              <p className="text-xs text-fg-muted">{nt.fat}</p>
              <p className="font-semibold text-rose-400">
                {goal.fat_goal_g}
                {nt.grams}
              </p>
            </div>
          </div>
        </div>
      ) : (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex flex-col items-center gap-3 py-4">
            <Target className="w-10 h-10 text-fg-muted" />
            <p className="text-sm text-fg-muted text-center">{nt.noGoalSet}</p>
            <button
              onClick={() => setGoalModal(true)}
              className="px-4 py-2 text-sm font-medium rounded-lg bg-emerald-600 text-white hover:bg-emerald-500 transition-colors"
            >
              {nt.setGoal}
            </button>
          </div>
        </div>
      )}

      {/* Goal Modal */}
      {goalModal && (
        <div
          className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setGoalModal(false);
              setGoalSuccess(false);
            }
          }}
        >
          <div className="w-full sm:max-w-md bg-bg border border-white/[0.1] rounded-t-2xl sm:rounded-2xl max-h-[85vh] overflow-y-auto shadow-2xl">
            <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
              <h3 className="text-sm font-semibold text-fg">{nt.setGoal}</h3>
              <button
                onClick={() => {
                  setGoalModal(false);
                  setGoalSuccess(false);
                }}
                className="p-1 rounded-lg hover:bg-white/[0.06] transition-colors"
              >
                <X className="w-5 h-5 text-fg-muted" />
              </button>
            </div>

            {goalSuccess ? (
              <div className="p-8 flex flex-col items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-emerald-500/20 flex items-center justify-center">
                  <Target className="w-6 h-6 text-emerald-400" />
                </div>
                <p className="text-sm font-medium text-emerald-400">
                  {nt.goalSet}
                </p>
              </div>
            ) : (
              <div className="p-5 space-y-4">
                {/* Weight */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {nt.weight}
                  </label>
                  <input
                    type="number"
                    value={goalForm.weight_kg}
                    onChange={(e) =>
                      setGoalForm((p) => ({ ...p, weight_kg: e.target.value }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  />
                </div>

                {/* Height */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {nt.height}
                  </label>
                  <input
                    type="number"
                    value={goalForm.height_cm}
                    onChange={(e) =>
                      setGoalForm((p) => ({ ...p, height_cm: e.target.value }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  />
                </div>

                {/* Age */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {nt.age}
                  </label>
                  <input
                    type="number"
                    value={goalForm.age}
                    onChange={(e) =>
                      setGoalForm((p) => ({ ...p, age: e.target.value }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  />
                </div>

                {/* Gender */}
                <div>
                  <label className="text-xs text-fg-muted font-medium mb-2 block">
                    {nt.gender}
                  </label>
                  <div className="flex gap-3">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="gender"
                        value="male"
                        checked={goalForm.gender === "male"}
                        onChange={(e) =>
                          setGoalForm((p) => ({
                            ...p,
                            gender: e.target.value,
                          }))
                        }
                        className="accent-emerald-500"
                      />
                      <span className="text-sm text-fg">{nt.male}</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="gender"
                        value="female"
                        checked={goalForm.gender === "female"}
                        onChange={(e) =>
                          setGoalForm((p) => ({
                            ...p,
                            gender: e.target.value,
                          }))
                        }
                        className="accent-emerald-500"
                      />
                      <span className="text-sm text-fg">{nt.female}</span>
                    </label>
                  </div>
                </div>

                {/* Activity Level */}
                <div>
                  <label className="text-xs text-fg-muted font-medium">
                    {nt.activityLevel}
                  </label>
                  <select
                    value={goalForm.activity_level}
                    onChange={(e) =>
                      setGoalForm((p) => ({
                        ...p,
                        activity_level: e.target.value,
                      }))
                    }
                    className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  >
                    <option value="sedentary">{nt.sedentary}</option>
                    <option value="light">{nt.light}</option>
                    <option value="moderate">{nt.moderate}</option>
                    <option value="active">{nt.active}</option>
                    <option value="very_active">{nt.veryActive}</option>
                  </select>
                </div>

                {/* Goal Type */}
                <div>
                  <label className="text-xs text-fg-muted font-medium mb-2 block">
                    {nt.goalType}
                  </label>
                  <div className="flex flex-col gap-2">
                    {[
                      { value: "lose_weight", label: nt.loseWeight },
                      { value: "maintain", label: nt.maintain },
                      { value: "build_muscle", label: nt.buildMuscle },
                    ].map((opt) => (
                      <label
                        key={opt.value}
                        className="flex items-center gap-2 cursor-pointer"
                      >
                        <input
                          type="radio"
                          name="goal_type"
                          value={opt.value}
                          checked={goalForm.goal_type === opt.value}
                          onChange={(e) =>
                            setGoalForm((p) => ({
                              ...p,
                              goal_type: e.target.value,
                            }))
                          }
                          className="accent-emerald-500"
                        />
                        <span className="text-sm text-fg">{opt.label}</span>
                      </label>
                    ))}
                  </div>
                </div>

                {/* Submit */}
                <button
                  onClick={handleSetGoal}
                  disabled={
                    goalSaving ||
                    !goalForm.weight_kg ||
                    !goalForm.height_cm ||
                    !goalForm.age
                  }
                  className="w-full px-4 py-2.5 text-sm rounded-lg bg-emerald-600 text-white font-medium hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {goalSaving ? (
                    <Loader2 className="w-4 h-4 animate-spin mx-auto" />
                  ) : (
                    nt.calculateGoal
                  )}
                </button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* ----------------------------------------------------------------- */}
      {/* 10. Meal Templates Section                                        */}
      {/* ----------------------------------------------------------------- */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
            <BookOpen className="w-4 h-4 text-fg-muted" />
            {nt.mealTemplates}
          </h2>
          <button
            onClick={() => setTemplateModal(true)}
            className="flex items-center gap-1 text-xs text-emerald-400 hover:text-emerald-300 transition-colors font-medium"
          >
            <Plus className="w-3.5 h-3.5" />
            {nt.createTemplate}
          </button>
        </div>

        {templates.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-10 gap-2">
            <BookOpen className="w-8 h-8 text-fg-muted" />
            <p className="text-xs text-fg-muted">{nt.noTemplates}</p>
          </div>
        ) : (
          <ul className="divide-y divide-white/[0.04]">
            {templates.map((tpl) => (
              <li
                key={tpl.id}
                className="flex items-center justify-between px-5 py-3"
              >
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-fg truncate">{tpl.name}</p>
                  <p className="text-xs text-fg-muted capitalize">
                    {nt[tpl.meal_type as keyof typeof nt] as string} |{" "}
                    {tpl.items?.length ?? 0} items
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handleQuickLog(tpl.id)}
                    className="flex items-center gap-1 px-3 py-1.5 text-xs font-medium rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30 transition-colors"
                  >
                    <Zap className="w-3 h-3" />
                    {nt.quickLog}
                  </button>
                  <button
                    onClick={() => handleDeleteTemplate(tpl.id)}
                    className="p-1.5 rounded hover:bg-white/[0.06] text-fg-muted hover:text-red-400 transition-colors"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Template Creation Modal */}
      {templateModal && (
        <div
          className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={(e) => {
            if (e.target === e.currentTarget) setTemplateModal(false);
          }}
        >
          <div className="w-full sm:max-w-md bg-bg border border-white/[0.1] rounded-t-2xl sm:rounded-2xl shadow-2xl">
            <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
              <h3 className="text-sm font-semibold text-fg">
                {nt.createTemplate}
              </h3>
              <button
                onClick={() => setTemplateModal(false)}
                className="p-1 rounded-lg hover:bg-white/[0.06] transition-colors"
              >
                <X className="w-5 h-5 text-fg-muted" />
              </button>
            </div>
            <div className="p-5 space-y-4">
              <div>
                <label className="text-xs text-fg-muted font-medium">
                  {nt.templateName}
                </label>
                <input
                  type="text"
                  value={templateForm.name}
                  onChange={(e) =>
                    setTemplateForm((p) => ({ ...p, name: e.target.value }))
                  }
                  className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg placeholder:text-fg-muted focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                  autoFocus
                />
              </div>
              <div>
                <label className="text-xs text-fg-muted font-medium">
                  Meal Type
                </label>
                <select
                  value={templateForm.meal_type}
                  onChange={(e) =>
                    setTemplateForm((p) => ({
                      ...p,
                      meal_type: e.target.value as MealType,
                    }))
                  }
                  className="w-full mt-1 px-4 py-2.5 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-fg focus:outline-none focus:ring-2 focus:ring-emerald-500/40"
                >
                  {MEAL_TYPES.map((m) => (
                    <option key={m} value={m}>
                      {nt[m as keyof typeof nt] as string}
                    </option>
                  ))}
                </select>
              </div>
              <button
                onClick={handleCreateTemplate}
                disabled={templateSaving || !templateForm.name.trim()}
                className="w-full px-4 py-2.5 text-sm rounded-lg bg-emerald-600 text-white font-medium hover:bg-emerald-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {templateSaving ? (
                  <Loader2 className="w-4 h-4 animate-spin mx-auto" />
                ) : (
                  nt.createTemplate
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
