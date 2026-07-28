"use client";

import { useCallback, useEffect, useState } from "react";
import { Loader2, MessageSquare, Phone, UserX } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Absentee, Locale } from "@/types";

/**
 * Members who have stopped turning up.
 *
 * The endpoints existed with no screen at all, which meant the one retention
 * tool the product has could not be used. A gym notices somebody has drifted
 * off weeks late, if ever; this is the list that makes it obvious, and the
 * message goes out in one click.
 */

const RANGES = [7, 14, 30, 60];

export default function AbsenteesPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;

  const [absentees, setAbsentees] = useState<Absentee[]>([]);
  const [days, setDays] = useState(14);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notifying, setNotifying] = useState<string | null>(null);
  // Who has already been texted this session, so staff do not send twice.
  const [notified, setNotified] = useState<Set<string>>(new Set());

  const fetchAbsentees = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listAbsentees(orgId, { days, limit: 200 });
      setAbsentees(res.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, days, t.common.error]);

  useEffect(() => {
    fetchAbsentees();
  }, [fetchAbsentees]);

  async function notify(member: Absentee) {
    if (!orgId) return;
    setNotifying(member.user_id);
    setError(null);
    try {
      await api.notifyAbsentee(orgId, member.user_id);
      setNotified((prev) => new Set(prev).add(member.user_id));
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setNotifying(null);
    }
  }

  return (
    <div className="space-y-6 animate-fade-in max-w-5xl">
      <div>
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <UserX className="w-5 h-5 text-accent" />
          {t.absentees.title}
        </h1>
        <p className="text-sm text-fg-muted mt-1">{t.absentees.subtitle}</p>
      </div>

      {error && (
        <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      <div className="flex items-center gap-1">
        {RANGES.map((d) => (
          <button
            key={d}
            type="button"
            onClick={() => setDays(d)}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              days === d ? "bg-accent/10 text-accent" : "text-fg-muted hover:bg-surface"
            }`}
          >
            {d}+ {t.absentees.days}
          </button>
        ))}
      </div>

      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        {loading ? (
          <div className="flex justify-center py-10">
            <Loader2 className="w-5 h-5 animate-spin text-fg-muted" />
          </div>
        ) : absentees.length === 0 ? (
          <p className="py-10 text-center text-sm text-fg-muted">{t.absentees.nobody}</p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {absentees.map((m) => (
              <div key={m.user_id} className="py-3 flex flex-wrap items-center gap-3">
                <div className="min-w-0 flex-1">
                  <p className="text-sm text-fg">{m.name}</p>
                  <a
                    href={`tel:${m.phone}`}
                    className="inline-flex items-center gap-1.5 text-xs font-mono text-fg-muted hover:text-accent transition-colors"
                  >
                    <Phone className="w-3 h-3" />
                    {m.phone}
                  </a>
                </div>

                <span
                  className={`shrink-0 px-2.5 py-1 rounded-full text-xs font-mono ${
                    m.absent_days >= 30
                      ? "bg-red-500/10 text-red-400"
                      : "bg-amber-500/10 text-amber-400"
                  }`}
                >
                  {m.absent_days} {t.absentees.daysAway}
                </span>

                <button
                  type="button"
                  onClick={() => notify(m)}
                  disabled={notifying === m.user_id || notified.has(m.user_id)}
                  className="shrink-0 inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-white/[0.06] text-sm text-fg hover:bg-surface transition-colors disabled:opacity-50"
                >
                  {notifying === m.user_id ? (
                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  ) : (
                    <MessageSquare className="w-3.5 h-3.5" />
                  )}
                  {notified.has(m.user_id) ? t.absentees.sent : t.absentees.sendMessage}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
