"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth, isOperator } from "@/lib/auth";
import type { Locale, UserType, OnboardingInput, NutritionGoal } from "@/types";
import {
  Loader2,
  Dumbbell,
  Activity,
  Utensils,
  Target,
  Heart,
  Flame,
  Sparkles,
  ArrowRight,
  ArrowLeft,
  Check,
  Building2,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type GoalType = OnboardingInput["goal_type"];

interface OnboardingState {
  user_type: UserType | "";
  org_id: string;
  goal_type: GoalType | "";
  name: string;
  weight_kg: string;
  height_cm: string;
  age: string;
  gender: string;
  activity_level: string;
}

// ---------------------------------------------------------------------------
// Step indicator
// ---------------------------------------------------------------------------

function StepIndicator({
  current,
  total,
  t,
}: {
  current: number;
  total: number;
  t: ReturnType<typeof getDictionary>;
}) {
  return (
    <div className="flex flex-col items-center gap-3">
      <p className="text-xs text-fg-muted font-medium tracking-wide uppercase">
        {t.onboarding.step} {current} {t.onboarding.of} {total}
      </p>
      <div className="flex items-center gap-2">
        {Array.from({ length: total }, (_, i) => (
          <div
            key={i}
            className={`h-1.5 rounded-full transition-all duration-500 ${
              i + 1 < current
                ? "w-8 bg-accent"
                : i + 1 === current
                  ? "w-10 bg-accent"
                  : "w-6 bg-white/[0.1]"
            }`}
          />
        ))}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function OnboardingPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const ob = t.onboarding;
  const router = useRouter();
  const { isLoading: authLoading, user } = useAuth();

  // Form state
  const [form, setForm] = useState<OnboardingState>({
    user_type: "",
    org_id: "",
    goal_type: "",
    name: "",
    weight_kg: "",
    height_cm: "",
    age: "",
    gender: "",
    activity_level: "moderate",
  });

  // Somebody who runs a gym has no business being asked which gym they would
  // like to join. Reachable directly by URL, or by an operator whose
  // onboarding_completed was never set, so it is guarded here too rather than
  // only at the login redirect.
  useEffect(() => {
    if (!authLoading && isOperator(user)) {
      router.replace(`/${locale}/dashboard`);
    }
  }, [authLoading, user, router, locale]);

  // Step management
  const [step, setStep] = useState(1);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [joiningGym, setJoiningGym] = useState(false);
  const [result, setResult] = useState<NutritionGoal | null>(null);

  // Direction for animation
  const [direction, setDirection] = useState<"forward" | "backward">("forward");

  // Calculate total steps: if gym_member, include gym join step
  const isGymMember = form.user_type === "gym_member";
  const totalSteps = isGymMember ? 5 : 4;

  // Map logical step to display step (accounting for optional gym step)
  function getDisplayStep(): number {
    if (isGymMember) return step;
    // If not gym member, steps are: 1 (type), 2 (goal), 3 (about), 4 (done)
    // We skip the gym step (which would be step 2 for gym members)
    if (step >= 2) return step;
    return step;
  }

  // Navigation
  function goNext() {
    setDirection("forward");
    setError("");
    // If step 1 and not gym member, skip gym step
    if (step === 1 && !isGymMember) {
      setStep(2);
    } else if (step === 1 && isGymMember) {
      setStep(2); // gym join step
    } else {
      setStep((s) => s + 1);
    }
  }

  function goBack() {
    setDirection("backward");
    setError("");
    if (step === 2 && !isGymMember) {
      setStep(1);
    } else {
      setStep((s) => s - 1);
    }
  }

  // For gym members, the steps are:
  //   1: user type selection
  //   2: gym join
  //   3: goal
  //   4: about you
  //   5: all set
  //
  // For non-gym members:
  //   1: user type selection
  //   2: goal
  //   3: about you
  //   4: all set

  // Determine which "logical content" to show per step
  function getStepContent(): "user_type" | "gym" | "goal" | "about" | "done" {
    if (isGymMember) {
      switch (step) {
        case 1: return "user_type";
        case 2: return "gym";
        case 3: return "goal";
        case 4: return "about";
        case 5: return "done";
        default: return "user_type";
      }
    } else {
      switch (step) {
        case 1: return "user_type";
        case 2: return "goal";
        case 3: return "about";
        case 4: return "done";
        default: return "user_type";
      }
    }
  }

  // Submit onboarding
  async function handleSubmit() {
    if (!form.user_type || !form.goal_type || !form.name) return;
    setSubmitting(true);
    setError("");

    try {
      const payload: OnboardingInput = {
        user_type: form.user_type as UserType,
        name: form.name,
        goal_type: form.goal_type as GoalType,
        weight_kg: form.weight_kg ? parseFloat(form.weight_kg) : undefined,
        height_cm: form.height_cm ? parseFloat(form.height_cm) : undefined,
        age: form.age ? parseInt(form.age, 10) : undefined,
        gender: form.gender || undefined,
        activity_level: form.activity_level || undefined,
      };

      await api.submitOnboarding(payload);

      // If weight/height/age provided, try to get nutrition goal for the results screen
      if (form.weight_kg && form.height_cm && form.age) {
        try {
          const goalData = await api.setNutritionGoal({
            weight_kg: parseFloat(form.weight_kg),
            height_cm: parseFloat(form.height_cm),
            age: parseInt(form.age, 10),
            gender: form.gender || "male",
            activity_level: form.activity_level || "moderate",
            goal_type: form.goal_type as string,
          });
          setResult(goalData);
        } catch {
          // Non-critical, continue anyway
        }
      }

      setDirection("forward");
      setStep(isGymMember ? 5 : 4);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setSubmitting(false);
    }
  }

  // Join gym
  async function handleJoinGym() {
    if (!form.org_id.trim()) return;
    setJoiningGym(true);
    setError("");
    try {
      await api.joinGym(form.org_id.trim());
      goNext();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to join gym");
    } finally {
      setJoiningGym(false);
    }
  }

  // Go to member home
  function handleLetsGo() {
    router.push(`/${locale}/member`);
  }

  // Can proceed to next step
  function canProceed(): boolean {
    const content = getStepContent();
    switch (content) {
      case "user_type":
        return !!form.user_type;
      case "gym":
        return true; // can skip
      case "goal":
        return !!form.goal_type;
      case "about":
        return !!form.name.trim();
      default:
        return true;
    }
  }

  // Loading state
  if (authLoading) {
    return (
      <div className="min-h-screen bg-bg-deep flex items-center justify-center">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const content = getStepContent();
  const displayStep = getDisplayStep();

  return (
    <div className="min-h-screen bg-bg-deep flex flex-col">
      {/* Ambient glow */}
      <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] rounded-full bg-accent/[0.04] blur-[150px] pointer-events-none" />

      {/* Header */}
      <header className="relative z-10 pt-8 pb-4 text-center">
        <div className="flex items-center justify-center gap-2.5 mb-6">
          <div className="w-10 h-10 rounded-xl bg-accent/20 flex items-center justify-center">
            <Dumbbell className="w-5 h-5 text-accent" />
          </div>
          <span className="text-xl font-bold tracking-tight text-fg">
            {t.common.appName}
          </span>
        </div>
        {content !== "done" && (
          <StepIndicator current={displayStep} total={totalSteps} t={t} />
        )}
      </header>

      {/* Content area */}
      <main className="relative z-10 flex-1 flex items-start justify-center px-4 pb-12">
        <div
          className={`w-full max-w-lg transition-all duration-400 ${
            direction === "forward"
              ? "animate-fade-in"
              : "animate-fade-in"
          }`}
          key={step}
        >
          {/* ============================================================= */}
          {/* Step: User Type Selection                                      */}
          {/* ============================================================= */}
          {content === "user_type" && (
            <div className="space-y-6">
              <div className="text-center">
                <h1 className="text-2xl sm:text-3xl font-bold text-fg tracking-tight">
                  {ob.whatBringsYou}
                </h1>
              </div>

              <div className="space-y-3">
                {/* Gym Member */}
                <button
                  onClick={() => setForm((f) => ({ ...f, user_type: "gym_member" }))}
                  className={`w-full flex items-start gap-4 p-5 rounded-2xl border transition-all duration-200 text-left ${
                    form.user_type === "gym_member"
                      ? "bg-accent/[0.08] border-accent/30 shadow-card-hover"
                      : "bg-white/[0.03] border-white/[0.06] hover:bg-white/[0.05] hover:border-white/[0.1]"
                  }`}
                >
                  <div
                    className={`w-12 h-12 rounded-xl flex items-center justify-center shrink-0 transition-colors ${
                      form.user_type === "gym_member"
                        ? "bg-accent/20"
                        : "bg-white/[0.06]"
                    }`}
                  >
                    <Building2
                      className={`w-6 h-6 ${
                        form.user_type === "gym_member"
                          ? "text-accent"
                          : "text-fg-muted"
                      }`}
                    />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-fg">{ob.gymMember}</p>
                    <p className="text-xs text-fg-muted mt-1 leading-relaxed">
                      {ob.gymMemberDesc}
                    </p>
                  </div>
                  {form.user_type === "gym_member" && (
                    <div className="w-5 h-5 rounded-full bg-accent flex items-center justify-center shrink-0 mt-0.5">
                      <Check className="w-3 h-3 text-bg-deep" />
                    </div>
                  )}
                </button>

                {/* Fitness Tracker */}
                <button
                  onClick={() => setForm((f) => ({ ...f, user_type: "fitness_tracker" }))}
                  className={`w-full flex items-start gap-4 p-5 rounded-2xl border transition-all duration-200 text-left ${
                    form.user_type === "fitness_tracker"
                      ? "bg-accent/[0.08] border-accent/30 shadow-card-hover"
                      : "bg-white/[0.03] border-white/[0.06] hover:bg-white/[0.05] hover:border-white/[0.1]"
                  }`}
                >
                  <div
                    className={`w-12 h-12 rounded-xl flex items-center justify-center shrink-0 transition-colors ${
                      form.user_type === "fitness_tracker"
                        ? "bg-accent/20"
                        : "bg-white/[0.06]"
                    }`}
                  >
                    <Activity
                      className={`w-6 h-6 ${
                        form.user_type === "fitness_tracker"
                          ? "text-accent"
                          : "text-fg-muted"
                      }`}
                    />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-fg">
                      {ob.fitnessTracker}
                    </p>
                    <p className="text-xs text-fg-muted mt-1 leading-relaxed">
                      {ob.fitnessTrackerDesc}
                    </p>
                  </div>
                  {form.user_type === "fitness_tracker" && (
                    <div className="w-5 h-5 rounded-full bg-accent flex items-center justify-center shrink-0 mt-0.5">
                      <Check className="w-3 h-3 text-bg-deep" />
                    </div>
                  )}
                </button>

                {/* Calorie Tracker */}
                <button
                  onClick={() => setForm((f) => ({ ...f, user_type: "calorie_tracker" }))}
                  className={`w-full flex items-start gap-4 p-5 rounded-2xl border transition-all duration-200 text-left ${
                    form.user_type === "calorie_tracker"
                      ? "bg-accent/[0.08] border-accent/30 shadow-card-hover"
                      : "bg-white/[0.03] border-white/[0.06] hover:bg-white/[0.05] hover:border-white/[0.1]"
                  }`}
                >
                  <div
                    className={`w-12 h-12 rounded-xl flex items-center justify-center shrink-0 transition-colors ${
                      form.user_type === "calorie_tracker"
                        ? "bg-accent/20"
                        : "bg-white/[0.06]"
                    }`}
                  >
                    <Utensils
                      className={`w-6 h-6 ${
                        form.user_type === "calorie_tracker"
                          ? "text-accent"
                          : "text-fg-muted"
                      }`}
                    />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-fg">
                      {ob.calorieTracker}
                    </p>
                    <p className="text-xs text-fg-muted mt-1 leading-relaxed">
                      {ob.calorieTrackerDesc}
                    </p>
                  </div>
                  {form.user_type === "calorie_tracker" && (
                    <div className="w-5 h-5 rounded-full bg-accent flex items-center justify-center shrink-0 mt-0.5">
                      <Check className="w-3 h-3 text-bg-deep" />
                    </div>
                  )}
                </button>
              </div>

              {/* Next button */}
              <button
                onClick={goNext}
                disabled={!canProceed()}
                className="w-full py-3.5 rounded-xl bg-accent text-bg-deep font-semibold text-sm hover:bg-accent-bright disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 shadow-accent-glow flex items-center justify-center gap-2"
              >
                {ob.next}
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          )}

          {/* ============================================================= */}
          {/* Step: Gym Join (only for gym_member)                           */}
          {/* ============================================================= */}
          {content === "gym" && (
            <div className="space-y-6">
              <div className="text-center">
                <h1 className="text-2xl sm:text-3xl font-bold text-fg tracking-tight">
                  {ob.joinGym}
                </h1>
                <p className="text-sm text-fg-muted mt-2">
                  {ob.scanGymQr}
                </p>
              </div>

              <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-6 space-y-4">
                <div>
                  <label className="text-xs text-fg-muted font-medium block mb-1.5">
                    Gym ID
                  </label>
                  <input
                    type="text"
                    value={form.org_id}
                    onChange={(e) =>
                      setForm((f) => ({ ...f, org_id: e.target.value }))
                    }
                    placeholder="Enter gym organization ID"
                    className="w-full px-4 py-3 text-sm bg-white/[0.04] border border-white/[0.08] rounded-xl text-fg placeholder:text-fg-muted/60 focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent/30 transition-all"
                  />
                </div>

                {error && (
                  <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                    {error}
                  </p>
                )}

                <button
                  onClick={handleJoinGym}
                  disabled={!form.org_id.trim() || joiningGym}
                  className="w-full py-3 rounded-xl bg-accent text-bg-deep font-semibold text-sm hover:bg-accent-bright disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 shadow-accent-glow flex items-center justify-center gap-2"
                >
                  {joiningGym ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <>
                      <Building2 className="w-4 h-4" />
                      {ob.joinGym}
                    </>
                  )}
                </button>
              </div>

              {/* Navigation */}
              <div className="flex items-center gap-3">
                <button
                  onClick={goBack}
                  className="flex-1 py-3 rounded-xl border border-white/[0.08] text-fg-muted font-medium text-sm hover:bg-white/[0.04] transition-all duration-200 flex items-center justify-center gap-2"
                >
                  <ArrowLeft className="w-4 h-4" />
                  {ob.back}
                </button>
                <button
                  onClick={goNext}
                  className="flex-1 py-3 rounded-xl border border-white/[0.08] text-fg-muted font-medium text-sm hover:bg-white/[0.04] transition-all duration-200 flex items-center justify-center gap-2"
                >
                  {ob.skipForNow}
                  <ArrowRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          )}

          {/* ============================================================= */}
          {/* Step: Goal Selection                                           */}
          {/* ============================================================= */}
          {content === "goal" && (
            <div className="space-y-6">
              <div className="text-center">
                <h1 className="text-2xl sm:text-3xl font-bold text-fg tracking-tight">
                  {ob.whatsYourGoal}
                </h1>
              </div>

              <div className="grid grid-cols-2 gap-3">
                {([
                  {
                    value: "lose_weight" as GoalType,
                    label: ob.loseWeight,
                    icon: Flame,
                    color: "text-orange-400",
                    activeBg: "bg-orange-500/[0.08]",
                    activeBorder: "border-orange-500/30",
                  },
                  {
                    value: "build_muscle" as GoalType,
                    label: ob.buildMuscle,
                    icon: Dumbbell,
                    color: "text-blue-400",
                    activeBg: "bg-blue-500/[0.08]",
                    activeBorder: "border-blue-500/30",
                  },
                  {
                    value: "stay_fit" as GoalType,
                    label: ob.stayFit,
                    icon: Activity,
                    color: "text-emerald-400",
                    activeBg: "bg-emerald-500/[0.08]",
                    activeBorder: "border-emerald-500/30",
                  },
                  {
                    value: "general_health" as GoalType,
                    label: ob.generalHealth,
                    icon: Heart,
                    color: "text-pink-400",
                    activeBg: "bg-pink-500/[0.08]",
                    activeBorder: "border-pink-500/30",
                  },
                ]).map((goal) => {
                  const isActive = form.goal_type === goal.value;
                  return (
                    <button
                      key={goal.value}
                      onClick={() =>
                        setForm((f) => ({ ...f, goal_type: goal.value }))
                      }
                      className={`flex flex-col items-center gap-3 p-5 rounded-2xl border transition-all duration-200 ${
                        isActive
                          ? `${goal.activeBg} ${goal.activeBorder} shadow-card-hover`
                          : "bg-white/[0.03] border-white/[0.06] hover:bg-white/[0.05] hover:border-white/[0.1]"
                      }`}
                    >
                      <div
                        className={`w-12 h-12 rounded-xl flex items-center justify-center transition-colors ${
                          isActive ? "bg-white/[0.1]" : "bg-white/[0.06]"
                        }`}
                      >
                        <goal.icon
                          className={`w-6 h-6 ${
                            isActive ? goal.color : "text-fg-muted"
                          }`}
                        />
                      </div>
                      <span
                        className={`text-sm font-medium ${
                          isActive ? "text-fg" : "text-fg-muted"
                        }`}
                      >
                        {goal.label}
                      </span>
                      {isActive && (
                        <div className="w-5 h-5 rounded-full bg-accent flex items-center justify-center">
                          <Check className="w-3 h-3 text-bg-deep" />
                        </div>
                      )}
                    </button>
                  );
                })}
              </div>

              {/* Navigation */}
              <div className="flex items-center gap-3">
                <button
                  onClick={goBack}
                  className="py-3 px-5 rounded-xl border border-white/[0.08] text-fg-muted font-medium text-sm hover:bg-white/[0.04] transition-all duration-200 flex items-center justify-center gap-2"
                >
                  <ArrowLeft className="w-4 h-4" />
                  {ob.back}
                </button>
                <button
                  onClick={goNext}
                  disabled={!canProceed()}
                  className="flex-1 py-3 rounded-xl bg-accent text-bg-deep font-semibold text-sm hover:bg-accent-bright disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 shadow-accent-glow flex items-center justify-center gap-2"
                >
                  {ob.next}
                  <ArrowRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          )}

          {/* ============================================================= */}
          {/* Step: About You                                                */}
          {/* ============================================================= */}
          {content === "about" && (
            <div className="space-y-6">
              <div className="text-center">
                <h1 className="text-2xl sm:text-3xl font-bold text-fg tracking-tight">
                  {ob.aboutYou}
                </h1>
              </div>

              <div className="space-y-4">
                {/* Name */}
                <div>
                  <label className="text-xs text-fg-muted font-medium block mb-1.5">
                    {ob.name}
                  </label>
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) =>
                      setForm((f) => ({ ...f, name: e.target.value }))
                    }
                    placeholder={ob.namePlaceholder}
                    className="w-full px-4 py-3 text-sm bg-white/[0.04] border border-white/[0.08] rounded-xl text-fg placeholder:text-fg-muted/60 focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent/30 transition-all"
                  />
                </div>

                {/* Weight & Height side by side */}
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="text-xs text-fg-muted font-medium block mb-1.5">
                      {ob.weight}
                    </label>
                    <input
                      type="number"
                      value={form.weight_kg}
                      onChange={(e) =>
                        setForm((f) => ({ ...f, weight_kg: e.target.value }))
                      }
                      placeholder="70"
                      className="w-full px-4 py-3 text-sm bg-white/[0.04] border border-white/[0.08] rounded-xl text-fg placeholder:text-fg-muted/60 focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent/30 transition-all"
                    />
                  </div>
                  <div>
                    <label className="text-xs text-fg-muted font-medium block mb-1.5">
                      {ob.height}
                    </label>
                    <input
                      type="number"
                      value={form.height_cm}
                      onChange={(e) =>
                        setForm((f) => ({ ...f, height_cm: e.target.value }))
                      }
                      placeholder="170"
                      className="w-full px-4 py-3 text-sm bg-white/[0.04] border border-white/[0.08] rounded-xl text-fg placeholder:text-fg-muted/60 focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent/30 transition-all"
                    />
                  </div>
                </div>

                {/* Age */}
                <div>
                  <label className="text-xs text-fg-muted font-medium block mb-1.5">
                    {ob.age}
                  </label>
                  <input
                    type="number"
                    value={form.age}
                    onChange={(e) =>
                      setForm((f) => ({ ...f, age: e.target.value }))
                    }
                    placeholder="25"
                    className="w-full px-4 py-3 text-sm bg-white/[0.04] border border-white/[0.08] rounded-xl text-fg placeholder:text-fg-muted/60 focus:outline-none focus:ring-2 focus:ring-accent/40 focus:border-accent/30 transition-all"
                  />
                </div>

                {/* Gender */}
                <div>
                  <label className="text-xs text-fg-muted font-medium block mb-2">
                    {ob.gender}
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {([
                      { value: "male", label: ob.male },
                      { value: "female", label: ob.female },
                      { value: "other", label: ob.other },
                      { value: "prefer_not_say", label: ob.preferNotSay },
                    ]).map((g) => (
                      <button
                        key={g.value}
                        onClick={() =>
                          setForm((f) => ({ ...f, gender: g.value }))
                        }
                        className={`px-4 py-2 rounded-xl text-sm font-medium border transition-all duration-200 ${
                          form.gender === g.value
                            ? "bg-accent/[0.15] border-accent/30 text-fg"
                            : "bg-white/[0.03] border-white/[0.06] text-fg-muted hover:bg-white/[0.05]"
                        }`}
                      >
                        {g.label}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Activity Level */}
                <div>
                  <label className="text-xs text-fg-muted font-medium block mb-2">
                    {ob.activityLevel}
                  </label>
                  <div className="space-y-2">
                    {([
                      { value: "sedentary", label: ob.sedentary },
                      { value: "light", label: ob.light },
                      { value: "moderate", label: ob.moderate },
                      { value: "active", label: ob.active },
                      { value: "very_active", label: ob.veryActive },
                    ]).map((a) => (
                      <button
                        key={a.value}
                        onClick={() =>
                          setForm((f) => ({ ...f, activity_level: a.value }))
                        }
                        className={`w-full flex items-center justify-between px-4 py-3 rounded-xl text-sm border transition-all duration-200 ${
                          form.activity_level === a.value
                            ? "bg-accent/[0.08] border-accent/30 text-fg"
                            : "bg-white/[0.03] border-white/[0.06] text-fg-muted hover:bg-white/[0.05]"
                        }`}
                      >
                        <span className="font-medium">{a.label}</span>
                        {form.activity_level === a.value && (
                          <div className="w-4 h-4 rounded-full bg-accent flex items-center justify-center">
                            <Check className="w-2.5 h-2.5 text-bg-deep" />
                          </div>
                        )}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              {error && (
                <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                  {error}
                </p>
              )}

              {/* Navigation */}
              <div className="flex items-center gap-3">
                <button
                  onClick={goBack}
                  className="py-3 px-5 rounded-xl border border-white/[0.08] text-fg-muted font-medium text-sm hover:bg-white/[0.04] transition-all duration-200 flex items-center justify-center gap-2"
                >
                  <ArrowLeft className="w-4 h-4" />
                  {ob.back}
                </button>
                <button
                  onClick={handleSubmit}
                  disabled={!canProceed() || submitting}
                  className="flex-1 py-3 rounded-xl bg-accent text-bg-deep font-semibold text-sm hover:bg-accent-bright disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 shadow-accent-glow flex items-center justify-center gap-2"
                >
                  {submitting ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <>
                      {ob.next}
                      <ArrowRight className="w-4 h-4" />
                    </>
                  )}
                </button>
              </div>
            </div>
          )}

          {/* ============================================================= */}
          {/* Step: All Set                                                   */}
          {/* ============================================================= */}
          {content === "done" && (
            <div className="space-y-8 text-center">
              {/* Celebration */}
              <div className="flex flex-col items-center gap-4">
                <div className="relative">
                  <div className="w-20 h-20 rounded-3xl bg-accent/20 flex items-center justify-center">
                    <Sparkles className="w-10 h-10 text-accent" />
                  </div>
                  <div className="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-accent flex items-center justify-center">
                    <Check className="w-3.5 h-3.5 text-bg-deep" />
                  </div>
                </div>
                <h1 className="text-2xl sm:text-3xl font-bold text-fg tracking-tight">
                  {ob.allSet}
                </h1>
              </div>

              {/* Calorie / macro targets */}
              {result && (
                <div className="bg-white/[0.03] border border-white/[0.06] rounded-2xl p-6 space-y-5">
                  {/* Daily calories */}
                  <div className="flex flex-col items-center gap-2">
                    <div className="w-14 h-14 rounded-2xl bg-accent/[0.12] flex items-center justify-center">
                      <Target className="w-7 h-7 text-accent" />
                    </div>
                    <p className="text-xs text-fg-muted font-medium uppercase tracking-wider">
                      {ob.dailyCalories}
                    </p>
                    <p className="text-4xl font-bold text-fg tracking-tight">
                      {result.calorie_goal.toLocaleString()}
                    </p>
                    <p className="text-xs text-fg-muted">kcal / day</p>
                  </div>

                  {/* Macro breakdown */}
                  <div className="grid grid-cols-3 gap-3">
                    {/* Protein */}
                    <div className="bg-blue-500/[0.06] border border-blue-500/[0.12] rounded-xl p-4 flex flex-col items-center gap-1">
                      <p className="text-xs text-blue-400 font-medium">
                        {ob.protein}
                      </p>
                      <p className="text-xl font-bold text-fg">
                        {result.protein_goal_g}
                      </p>
                      <p className="text-[10px] text-fg-muted">g / day</p>
                    </div>

                    {/* Carbs */}
                    <div className="bg-amber-500/[0.06] border border-amber-500/[0.12] rounded-xl p-4 flex flex-col items-center gap-1">
                      <p className="text-xs text-amber-400 font-medium">
                        {ob.carbs}
                      </p>
                      <p className="text-xl font-bold text-fg">
                        {result.carbs_goal_g}
                      </p>
                      <p className="text-[10px] text-fg-muted">g / day</p>
                    </div>

                    {/* Fat */}
                    <div className="bg-rose-500/[0.06] border border-rose-500/[0.12] rounded-xl p-4 flex flex-col items-center gap-1">
                      <p className="text-xs text-rose-400 font-medium">
                        {ob.fat}
                      </p>
                      <p className="text-xl font-bold text-fg">
                        {result.fat_goal_g}
                      </p>
                      <p className="text-[10px] text-fg-muted">g / day</p>
                    </div>
                  </div>
                </div>
              )}

              {/* Let's Go button */}
              <button
                onClick={handleLetsGo}
                className="w-full py-3.5 rounded-xl bg-accent text-bg-deep font-semibold text-sm hover:bg-accent-bright transition-all duration-200 shadow-accent-glow flex items-center justify-center gap-2"
              >
                {ob.letsGo}
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
