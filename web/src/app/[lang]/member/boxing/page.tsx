"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { ArrowLeft, Loader2, Plus, ShieldCheck, ShieldAlert, Trash2 } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import type { BoxingProfileView } from "@/types/site-admin";

/**
 * A member's boxing profile and fight record.
 *
 * Sparring clearance is displayed but never editable here: it is a coach's
 * safety decision, granted through a staff-only endpoint. Showing it read-only
 * is the point — a member needs to know whether they are cleared.
 */

const STANCES = ["", "orthodox", "southpaw", "switch"];
const SKILL_LEVELS = ["", "beginner", "intermediate", "amateur", "pro"];
const RESULTS = ["win", "loss", "draw", "no_contest"];
const METHODS = ["", "decision", "ko", "tko", "submission", "dq", "walkover"];

function titleCase(value: string): string {
  return value
    .split("_")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

export default function MemberBoxingPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;

  const [profile, setProfile] = useState<BoxingProfileView | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [stance, setStance] = useState("");
  const [skillLevel, setSkillLevel] = useState("");
  const [weightClass, setWeightClass] = useState("");
  const [reachCm, setReachCm] = useState("");
  const [notes, setNotes] = useState("");

  const [boutOpen, setBoutOpen] = useState(false);
  const [boutDate, setBoutDate] = useState("");
  const [opponent, setOpponent] = useState("");
  const [eventName, setEventName] = useState("");
  const [result, setResult] = useState("win");
  const [method, setMethod] = useState("");
  const [rounds, setRounds] = useState("");
  const [boutError, setBoutError] = useState<string | null>(null);

  const applyProfile = useCallback((p: BoxingProfileView) => {
    setProfile(p);
    setStance(p.stance ?? "");
    setSkillLevel(p.skill_level ?? "");
    setWeightClass(p.weight_class ?? "");
    setReachCm(p.reach_cm ? String(p.reach_cm) : "");
    setNotes(p.notes ?? "");
  }, []);

  const load = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    try {
      applyProfile(await api.getMyBoxingProfile(orgId));
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, applyProfile, t.common.error]);

  useEffect(() => {
    load();
  }, [load]);

  async function saveProfile(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId) return;
    setSaving(true);
    setError(null);
    try {
      applyProfile(
        await api.updateMyBoxingProfile(orgId, {
          stance,
          skill_level: skillLevel,
          weight_class: weightClass,
          reach_cm: reachCm ? Number(reachCm) : null,
          notes,
        })
      );
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setSaving(false);
    }
  }

  async function addBout(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId) return;
    setBoutError(null);
    try {
      await api.createMyBout(orgId, {
        bout_date: boutDate,
        opponent,
        event_name: eventName,
        result,
        method,
        rounds: rounds ? Number(rounds) : null,
        weight_class: weightClass,
        notes: "",
      });
      setBoutOpen(false);
      setBoutDate("");
      setOpponent("");
      setEventName("");
      setMethod("");
      setRounds("");
      await load();
    } catch (err) {
      setBoutError(err instanceof ApiRequestError ? err.message : t.common.error);
    }
  }

  async function removeBout(boutId: string) {
    if (!orgId) return;
    try {
      await api.deleteMyBout(orgId, boutId);
      await load();
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Loader2 className="w-6 h-6 animate-spin text-fg-muted" />
      </div>
    );
  }

  const inputClass =
    "w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors";
  const labelClass = "block text-xs text-fg-muted mb-1.5";
  const record = profile?.record;

  return (
    <div className="max-w-3xl">
      <Link
        href={`/${locale}/member`}
        className="inline-flex items-center gap-2 text-sm text-fg-muted hover:text-fg transition-colors mb-4"
      >
        <ArrowLeft className="w-4 h-4" />
        {t.nav.dashboard}
      </Link>

      <h1 className="text-2xl font-semibold text-fg mb-6">Boxing</h1>

      {error && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Record + clearance */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card mb-5">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs text-fg-muted mb-1">Record</p>
            <p className="text-3xl font-semibold tabular-nums text-fg">
              {record?.wins ?? 0}
              <span className="text-fg-muted">–</span>
              {record?.losses ?? 0}
              <span className="text-fg-muted">–</span>
              {record?.draws ?? 0}
              {record && record.no_contests > 0 && (
                <span className="ml-2 text-sm text-fg-muted">
                  ({record.no_contests} NC)
                </span>
              )}
            </p>
          </div>

          <div
            className={`inline-flex items-center gap-2 px-3.5 py-2 rounded-xl border text-sm ${
              profile?.sparring_cleared
                ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                : "bg-white/[0.03] text-fg-muted border-white/[0.06]"
            }`}
          >
            {profile?.sparring_cleared ? (
              <ShieldCheck className="w-4 h-4" />
            ) : (
              <ShieldAlert className="w-4 h-4" />
            )}
            {profile?.sparring_cleared ? "Cleared to spar" : "Not cleared to spar"}
          </div>
        </div>
        {!profile?.sparring_cleared && (
          <p className="text-xs text-fg-muted mt-3">
            Your coach clears you when your technique is ready. It is never rushed.
          </p>
        )}
      </div>

      {/* Profile */}
      <form
        onSubmit={saveProfile}
        className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card mb-5"
      >
        <h2 className="text-sm font-medium text-fg mb-4">Profile</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label htmlFor="stance" className={labelClass}>
              Stance
            </label>
            <select
              id="stance"
              className={inputClass}
              value={stance}
              onChange={(e) => setStance(e.target.value)}
            >
              {STANCES.map((s) => (
                <option key={s} value={s}>
                  {s ? titleCase(s) : "—"}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="skill" className={labelClass}>
              Level
            </label>
            <select
              id="skill"
              className={inputClass}
              value={skillLevel}
              onChange={(e) => setSkillLevel(e.target.value)}
            >
              {SKILL_LEVELS.map((s) => (
                <option key={s} value={s}>
                  {s ? titleCase(s) : "—"}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="weight-class" className={labelClass}>
              Weight class
            </label>
            <input
              id="weight-class"
              className={inputClass}
              value={weightClass}
              onChange={(e) => setWeightClass(e.target.value)}
              placeholder={profile?.suggested_weight_class ?? "e.g. welterweight"}
            />
            {profile?.suggested_weight_class && (
              // Derived from the member's own logged weight, so they do not have
              // to work out which division they are in.
              <button
                type="button"
                onClick={() => setWeightClass(profile.suggested_weight_class!)}
                className="mt-1.5 text-[11px] text-accent hover:underline"
              >
                Use {titleCase(profile.suggested_weight_class)} (from your latest weight)
              </button>
            )}
          </div>
          <div>
            <label htmlFor="reach" className={labelClass}>
              Reach (cm)
            </label>
            <input
              id="reach"
              type="number"
              inputMode="decimal"
              step="0.5"
              className={inputClass}
              value={reachCm}
              onChange={(e) => setReachCm(e.target.value)}
            />
          </div>
        </div>

        <div className="mt-4">
          <label htmlFor="notes" className={labelClass}>
            Notes
          </label>
          <textarea
            id="notes"
            rows={3}
            className={inputClass}
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />
        </div>

        <button
          type="submit"
          disabled={saving}
          className="mt-4 inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-accent-fg text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {saving && <Loader2 className="w-4 h-4 animate-spin" />}
          {t.common.save}
        </button>
      </form>

      {/* Bouts */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        <div className="flex items-center justify-between gap-4 mb-4">
          <h2 className="text-sm font-medium text-fg">Bouts</h2>
          <button
            type="button"
            onClick={() => setBoutOpen((v) => !v)}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-white/[0.06] text-sm text-fg hover:bg-surface transition-colors"
          >
            <Plus className="w-4 h-4" />
            Add bout
          </button>
        </div>

        {boutOpen && (
          <form
            onSubmit={addBout}
            className="mb-5 p-4 rounded-xl bg-black/20 border border-white/[0.06]"
          >
            {boutError && (
              <p className="mb-3 text-xs text-red-400" role="alert">
                {boutError}
              </p>
            )}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label htmlFor="bout-date" className={labelClass}>
                  Date
                </label>
                <input
                  id="bout-date"
                  type="date"
                  required
                  className={inputClass}
                  value={boutDate}
                  onChange={(e) => setBoutDate(e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="opponent" className={labelClass}>
                  Opponent
                </label>
                <input
                  id="opponent"
                  className={inputClass}
                  value={opponent}
                  onChange={(e) => setOpponent(e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="event" className={labelClass}>
                  Event
                </label>
                <input
                  id="event"
                  className={inputClass}
                  value={eventName}
                  onChange={(e) => setEventName(e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="result" className={labelClass}>
                  Result
                </label>
                <select
                  id="result"
                  className={inputClass}
                  value={result}
                  onChange={(e) => setResult(e.target.value)}
                >
                  {RESULTS.map((r) => (
                    <option key={r} value={r}>
                      {titleCase(r)}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor="method" className={labelClass}>
                  Method
                </label>
                <select
                  id="method"
                  className={inputClass}
                  value={method}
                  onChange={(e) => setMethod(e.target.value)}
                >
                  {METHODS.map((m) => (
                    <option key={m} value={m}>
                      {m ? m.toUpperCase() : "—"}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor="rounds" className={labelClass}>
                  Rounds
                </label>
                <input
                  id="rounds"
                  type="number"
                  min={1}
                  max={15}
                  className={inputClass}
                  value={rounds}
                  onChange={(e) => setRounds(e.target.value)}
                />
              </div>
            </div>
            <button
              type="submit"
              className="mt-4 px-4 py-2 rounded-xl bg-accent text-accent-fg text-sm font-medium hover:opacity-90 transition-opacity"
            >
              {t.common.save}
            </button>
          </form>
        )}

        {!profile?.bouts?.length ? (
          <p className="py-8 text-center text-sm text-fg-muted">
            No bouts recorded yet.
          </p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {profile.bouts.map((bout) => (
              <div key={bout.id} className="flex items-center gap-3 py-3">
                <span
                  className={`shrink-0 inline-flex h-7 min-w-7 items-center justify-center rounded-full px-2 text-xs font-bold ${
                    bout.result === "win"
                      ? "bg-emerald-500/15 text-emerald-400"
                      : "bg-white/5 text-fg-muted"
                  }`}
                >
                  {bout.result === "no_contest" ? "NC" : bout.result.charAt(0).toUpperCase()}
                </span>
                <div className="min-w-0 flex-1">
                  <p className="text-sm text-fg truncate">
                    {bout.opponent || "Unnamed opponent"}
                  </p>
                  <p className="text-xs text-fg-muted truncate">
                    {[
                      bout.bout_date.slice(0, 10),
                      bout.event_name,
                      bout.method?.toUpperCase(),
                      bout.rounds ? `R${bout.rounds}` : "",
                    ]
                      .filter(Boolean)
                      .join(" · ")}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => removeBout(bout.id)}
                  aria-label={t.common.delete}
                  className="shrink-0 p-2 rounded-lg text-fg-muted hover:text-red-400 hover:bg-surface transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
