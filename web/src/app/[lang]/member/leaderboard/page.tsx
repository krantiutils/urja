"use client";

import { useCallback, useEffect, useState } from "react";
import { Loader2, Trophy } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { LeaderboardEntry, Locale } from "@/types";

/**
 * The gym leaderboard, as a member sees it.
 *
 * The endpoint existed with no screen, so the one bit of gamification in the
 * product was unreachable. Members who have opted out of the board through
 * their privacy settings simply do not appear — that is the API's decision,
 * not this screen's, so a member missing from the list is told why rather than
 * left to wonder.
 */

const PERIODS = ["weekly", "monthly", "alltime"] as const;

export default function LeaderboardPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [period, setPeriod] = useState<(typeof PERIODS)[number]>("monthly");
  const [rankings, setRankings] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchBoard = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.getMyLeaderboard({ period, limit: 50 });
      setRankings(res.rankings ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [period, t.common.error]);

  useEffect(() => {
    fetchBoard();
  }, [fetchBoard]);

  const label = {
    weekly: t.leaderboard.thisWeek,
    monthly: t.leaderboard.thisMonth,
    alltime: t.leaderboard.allTime,
  };
  const onBoard = rankings.some((r) => r.member_id === user?.id);

  return (
    <div className="space-y-6 max-w-2xl">
      <div>
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Trophy className="w-5 h-5 text-accent" />
          {t.leaderboard.title}
        </h1>
        <p className="text-sm text-fg-muted mt-1">{t.leaderboard.subtitle}</p>
      </div>

      <div className="flex items-center gap-1">
        {PERIODS.map((p) => (
          <button
            key={p}
            type="button"
            onClick={() => setPeriod(p)}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              period === p ? "bg-accent/10 text-accent" : "text-fg-muted hover:bg-surface"
            }`}
          >
            {label[p]}
          </button>
        ))}
      </div>

      {error && (
        <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        {loading ? (
          <div className="flex justify-center py-10">
            <Loader2 className="w-5 h-5 animate-spin text-fg-muted" />
          </div>
        ) : rankings.length === 0 ? (
          <p className="py-10 text-center text-sm text-fg-muted">{t.leaderboard.empty}</p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {rankings.map((r) => {
              const isYou = r.member_id === user?.id;
              return (
                <div
                  key={r.member_id}
                  className={`py-3 flex items-center gap-4 ${isYou ? "text-accent" : ""}`}
                >
                  <span
                    className={`w-8 shrink-0 text-center font-mono ${
                      r.rank <= 3 ? "text-accent font-bold" : "text-fg-muted"
                    }`}
                  >
                    {r.rank}
                  </span>
                  <span className="min-w-0 flex-1 truncate text-sm">
                    {r.name}
                    {isYou && (
                      <span className="ml-2 text-xs font-mono uppercase tracking-widest">
                        {t.leaderboard.you}
                      </span>
                    )}
                  </span>
                  <span className="shrink-0 text-sm tabular-nums">{r.value}</span>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* A member who opted out of the board would otherwise just not find
          themselves and assume it was broken. */}
      {!loading && rankings.length > 0 && !onBoard && (
        <p className="text-xs text-fg-muted">{t.leaderboard.hiddenNote}</p>
      )}
    </div>
  );
}
