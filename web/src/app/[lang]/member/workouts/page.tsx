"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type { Locale, WorkoutLog, WorkoutPlan } from "@/types";
import { Loader2, Dumbbell, ClipboardList, Clock } from "lucide-react";

export default function MemberWorkoutsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [plan, setPlan] = useState<WorkoutPlan | null>(null);
  const [logs, setLogs] = useState<WorkoutLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const orgId = user?.org_id;
        const results = await Promise.allSettled([
          orgId ? api.getMyWorkoutPlan(orgId) : Promise.reject("no org"),
          api.getMyWorkoutLogs({
            organization_id: orgId,
            limit: 20,
          }),
        ]);

        if (results[0].status === "fulfilled") {
          setPlan(results[0].value);
        }
        if (results[1].status === "fulfilled") {
          setLogs(results[1].value.data ?? []);
        }
      } catch (err) {
        console.error("Failed to load workouts:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [user?.org_id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const template = plan?.template;

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">
        {t.memberPages.workouts}
      </h1>

      {/* Assigned plan */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg">
            {t.memberPages.workoutPlan}
          </h2>
        </div>

        {!template ? (
          <div className="flex flex-col items-center justify-center py-16 gap-3">
            <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
              <ClipboardList className="w-6 h-6 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">{t.memberPages.noPlan}</p>
          </div>
        ) : (
          <div className="p-5 space-y-3">
            <h3 className="text-lg font-semibold text-fg">{template.name}</h3>
            {template.description && (
              <p className="text-sm text-fg-muted">{template.description}</p>
            )}
            <div className="flex flex-wrap gap-3 text-sm">
              {template.category && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border bg-accent/10 text-accent border-accent/20">
                  {template.category}
                </span>
              )}
              {template.difficulty && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border bg-purple-500/10 text-purple-400 border-purple-500/20">
                  {template.difficulty}
                </span>
              )}
              {template.duration_minutes != null && (
                <span className="inline-flex items-center gap-1 text-fg-muted">
                  <Clock className="w-3.5 h-3.5" />
                  {template.duration_minutes} {t.memberPages.minutes}
                </span>
              )}
            </div>
          </div>
        )}
      </div>

      {/* Recent workout logs */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg">
            {t.memberPages.recentWorkouts}
          </h2>
        </div>

        {logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 gap-3">
            <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
              <Dumbbell className="w-6 h-6 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">
              {t.memberPages.noWorkouts}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {logs.map((log) => {
              const date = new Date(log.logged_at);
              return (
                <div key={log.id} className="px-5 py-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-fg">
                        {date.toLocaleDateString()}
                      </p>
                      {log.notes && (
                        <p className="text-xs text-fg-muted mt-0.5 line-clamp-1">
                          {log.notes}
                        </p>
                      )}
                    </div>
                    {log.duration_minutes != null && (
                      <span className="inline-flex items-center gap-1 text-xs text-fg-muted">
                        <Clock className="w-3.5 h-3.5" />
                        {log.duration_minutes} {t.memberPages.minutes}
                      </span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
