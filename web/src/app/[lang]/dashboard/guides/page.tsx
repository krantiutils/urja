"use client";

import { useCallback, useEffect, useState } from "react";
import { BookOpen, Eye, EyeOff, Loader2, Pencil, Plus, Trash2, X } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale, TrainingGuide } from "@/types";

/**
 * Training guides a gym writes for its members.
 *
 * Eight endpoints existed with no screen at all — staff had no way to write a
 * guide and members no way to read one. Drafts are visible here and nowhere
 * else, so a guide can be worked on before anybody sees it.
 */

// Mirrors validCategories in internal/guide/service.go; the API rejects anything else.
const CATEGORIES = ["strength", "cardio", "flexibility", "nutrition", "recovery", "beginner"];

export default function GuidesPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;

  const [guides, setGuides] = useState<TrainingGuide[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState<string | null>(null);

  const [editing, setEditing] = useState<TrainingGuide | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [title, setTitle] = useState("");
  const [titleNe, setTitleNe] = useState("");
  const [content, setContent] = useState("");
  const [contentNe, setContentNe] = useState("");
  const [category, setCategory] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const fetchGuides = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listGuides(orgId, { limit: 100 });
      setGuides(res.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchGuides();
  }, [fetchGuides]);

  function openForm(guide: TrainingGuide | null) {
    setEditing(guide);
    setTitle(guide?.title ?? "");
    setTitleNe(guide?.title_ne ?? "");
    setContent(guide?.content ?? "");
    setContentNe(guide?.content_ne ?? "");
    setCategory(guide?.category ?? "");
    setFormError(null);
    setShowForm(true);
  }

  async function save(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId) return;
    setSaving(true);
    setFormError(null);
    try {
      const body = {
        title: title.trim(),
        title_ne: titleNe.trim(),
        content,
        content_ne: contentNe,
        category,
      };
      if (editing) {
        await api.updateGuide(orgId, editing.id, body);
      } else {
        await api.createGuide(orgId, body);
      }
      setShowForm(false);
      await fetchGuides();
    } catch (err) {
      setFormError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setSaving(false);
    }
  }

  async function togglePublish(guide: TrainingGuide) {
    if (!orgId) return;
    setBusy(guide.id);
    try {
      await api.publishGuide(orgId, guide.id, !guide.is_published);
      await fetchGuides();
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setBusy(null);
    }
  }

  async function remove(guide: TrainingGuide) {
    if (!orgId || !window.confirm(t.guides.deleteConfirm)) return;
    setBusy(guide.id);
    try {
      await api.deleteGuide(orgId, guide.id);
      setGuides((prev) => prev.filter((g) => g.id !== guide.id));
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setBusy(null);
    }
  }

  const inputClass =
    "w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors";

  return (
    <div className="space-y-6 animate-fade-in max-w-5xl">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
            <BookOpen className="w-5 h-5 text-accent" />
            {t.guides.title}
          </h1>
          <p className="text-sm text-fg-muted mt-1">{t.guides.subtitle}</p>
        </div>
        <button
          type="button"
          onClick={() => openForm(null)}
          className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity"
        >
          <Plus className="w-4 h-4" />
          {t.guides.newGuide}
        </button>
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
        ) : guides.length === 0 ? (
          <p className="py-10 text-center text-sm text-fg-muted">{t.guides.empty}</p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {guides.map((g) => (
              <div key={g.id} className="py-3 flex flex-wrap items-center gap-3">
                <div className="min-w-0 flex-1">
                  <p className="text-sm text-fg">{g.title}</p>
                  <p className="text-xs text-fg-muted">
                    {[g.category, new Date(g.updated_at).toLocaleDateString()]
                      .filter(Boolean)
                      .join(" · ")}
                  </p>
                </div>

                <span
                  className={`shrink-0 px-2 py-0.5 rounded-full text-[10px] font-mono uppercase border ${
                    g.is_published
                      ? "bg-accent/10 text-accent border-accent/20"
                      : "bg-amber-500/10 text-amber-400 border-amber-500/20"
                  }`}
                >
                  {g.is_published ? t.guides.published : t.guides.draft}
                </span>

                <div className="shrink-0 flex items-center gap-1">
                  <button
                    type="button"
                    onClick={() => togglePublish(g)}
                    disabled={busy === g.id}
                    title={g.is_published ? t.guides.unpublish : t.guides.publish}
                    className="p-1.5 rounded-lg text-fg-muted hover:text-accent hover:bg-surface transition-colors disabled:opacity-50"
                  >
                    {g.is_published ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                  <button
                    type="button"
                    onClick={() => openForm(g)}
                    title={t.common.save}
                    className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors"
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    type="button"
                    onClick={() => remove(g)}
                    disabled={busy === g.id}
                    title={t.common.delete}
                    className="p-1.5 rounded-lg text-fg-muted hover:text-red-400 hover:bg-surface transition-colors disabled:opacity-50"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black/60"
            onClick={() => setShowForm(false)}
            aria-hidden="true"
          />
          <form
            onSubmit={save}
            className="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto bg-bg-elevated border border-white/[0.06] rounded-2xl p-5 shadow-card"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base font-semibold text-fg">
                {editing ? t.guides.editGuide : t.guides.newGuide}
              </h2>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="p-1.5 rounded-lg text-fg-muted hover:bg-surface transition-colors"
                aria-label={t.common.cancel}
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {formError && (
              <p className="mb-3 text-sm text-red-400" role="alert">
                {formError}
              </p>
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-3">
              <div>
                <label htmlFor="g-title" className="block text-xs text-fg-muted mb-1.5">
                  {t.guides.guideTitle}
                </label>
                <input
                  id="g-title"
                  className={inputClass}
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                />
              </div>
              <div>
                <label htmlFor="g-title-ne" className="block text-xs text-fg-muted mb-1.5">
                  {t.guides.guideTitleNe}
                </label>
                <input
                  id="g-title-ne"
                  className={inputClass}
                  value={titleNe}
                  onChange={(e) => setTitleNe(e.target.value)}
                />
              </div>
            </div>

            <div className="mb-3">
              <label htmlFor="g-category" className="block text-xs text-fg-muted mb-1.5">
                {t.guides.category}
              </label>
              <select
                id="g-category"
                className={inputClass}
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              >
                <option value="">—</option>
                {CATEGORIES.map((c) => (
                  <option key={c} value={c}>
                    {c}
                  </option>
                ))}
              </select>
            </div>

            <div className="mb-3">
              <label htmlFor="g-content" className="block text-xs text-fg-muted mb-1.5">
                {t.guides.content}
              </label>
              <textarea
                id="g-content"
                rows={10}
                className={inputClass}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                required
              />
            </div>

            <div className="mb-5">
              <label htmlFor="g-content-ne" className="block text-xs text-fg-muted mb-1.5">
                {t.guides.contentNe}
              </label>
              <textarea
                id="g-content-ne"
                rows={6}
                className={inputClass}
                value={contentNe}
                onChange={(e) => setContentNe(e.target.value)}
              />
            </div>

            <button
              type="submit"
              disabled={saving || !title.trim() || !content.trim()}
              className="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {saving && <Loader2 className="w-4 h-4 animate-spin" />}
              {t.common.save}
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
