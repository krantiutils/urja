"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, Notice } from "@/types";
import {
  Newspaper,
  Plus,
  X,
  Pencil,
  Trash2,
  Loader2,
  ToggleLeft,
  ToggleRight,
} from "lucide-react";

export default function StoriesPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [notices, setNotices] = useState<Notice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Form state
  const [showModal, setShowModal] = useState(false);
  const [editingNotice, setEditingNotice] = useState<Notice | null>(null);
  const [formTitle, setFormTitle] = useState("");
  const [formTitleNe, setFormTitleNe] = useState("");
  const [formContent, setFormContent] = useState("");
  const [formContentNe, setFormContentNe] = useState("");
  const [formIsActive, setFormIsActive] = useState(true);
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchNotices = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listNotices(orgId, { limit: 100 });
      setNotices(res.data ?? []);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchNotices();
  }, [fetchNotices]);

  const openCreateModal = () => {
    setEditingNotice(null);
    setFormTitle("");
    setFormTitleNe("");
    setFormContent("");
    setFormContentNe("");
    setFormIsActive(true);
    setFormError(null);
    setShowModal(true);
  };

  const openEditModal = (notice: Notice) => {
    setEditingNotice(notice);
    setFormTitle(notice.title);
    setFormTitleNe(notice.title_ne ?? "");
    setFormContent(notice.content);
    setFormContentNe(notice.content_ne ?? "");
    setFormIsActive(notice.is_active);
    setFormError(null);
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setFormLoading(true);
    setFormError(null);

    try {
      if (editingNotice) {
        await api.updateNotice(orgId, editingNotice.id, {
          title: formTitle,
          title_ne: formTitleNe || undefined,
          content: formContent,
          content_ne: formContentNe || undefined,
          is_active: formIsActive,
        });
      } else {
        await api.createNotice(orgId, {
          title: formTitle,
          title_ne: formTitleNe || undefined,
          content: formContent,
          content_ne: formContentNe || undefined,
          is_active: formIsActive,
        });
      }
      setShowModal(false);
      fetchNotices();
    } catch (err) {
      setFormError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (noticeId: string) => {
    if (!orgId) return;
    if (!confirm(t.notices.confirmDelete)) return;
    try {
      await api.deleteNotice(orgId, noticeId);
      fetchNotices();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const handleToggleActive = async (notice: Notice) => {
    if (!orgId) return;
    try {
      await api.updateNotice(orgId, notice.id, {
        is_active: !notice.is_active,
      });
      fetchNotices();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString(
      locale === "ne" ? "ne-NP" : "en-US",
      { year: "numeric", month: "short", day: "numeric" }
    );
  };

  const truncateContent = (text: string, maxLen: number = 100) => {
    if (text.length <= maxLen) return text;
    return text.slice(0, maxLen) + "...";
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Newspaper className="w-5 h-5 text-accent" />
          {t.notices.title}
        </h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.notices.addNotice}
        </button>
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Notice Cards */}
      {loading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : notices.length === 0 ? (
        <div className="text-center py-16 text-fg-muted text-sm">
          {t.notices.noNotices}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {notices.map((notice) => (
            <div
              key={notice.id}
              className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card"
            >
              {/* Card header */}
              <div className="flex items-start justify-between mb-2">
                <h3 className="text-fg font-medium text-base">
                  {notice.title}
                </h3>
                <span
                  className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
                    notice.is_active
                      ? "bg-accent/10 text-accent border-accent/20"
                      : "bg-red-500/10 text-red-400 border-red-500/20"
                  }`}
                >
                  {notice.is_active ? t.notices.active : t.notices.inactive}
                </span>
              </div>

              {/* Card content preview */}
              <p className="text-fg-muted text-sm mt-2 line-clamp-2">
                {truncateContent(notice.content)}
              </p>

              {/* Card footer */}
              <div className="flex items-center justify-between pt-4 mt-4 border-t border-white/[0.06]">
                <span className="text-xs text-fg-muted">
                  {formatDate(notice.created_at)}
                </span>
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => openEditModal(notice)}
                    className="p-1.5 rounded-lg hover:bg-surface text-fg-muted hover:text-fg transition-colors"
                    title={t.notices.editTitle}
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleToggleActive(notice)}
                    className={`p-1.5 rounded-lg transition-colors ${
                      notice.is_active
                        ? "text-accent hover:bg-accent/10"
                        : "text-fg-muted hover:bg-surface"
                    }`}
                    title={t.notices.toggleActive}
                  >
                    {notice.is_active ? (
                      <ToggleRight className="w-4 h-4" />
                    ) : (
                      <ToggleLeft className="w-4 h-4" />
                    )}
                  </button>
                  <button
                    onClick={() => handleDelete(notice.id)}
                    className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
                    title={t.common.delete}
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-lg">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {editingNotice ? t.notices.editTitle : t.notices.addTitle}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-5 space-y-4">
              {formError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {formError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.notices.noticeTitle} *
                </label>
                <input
                  type="text"
                  required
                  value={formTitle}
                  onChange={(e) => setFormTitle(e.target.value)}
                  placeholder={t.notices.titlePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.notices.titleNe}
                </label>
                <input
                  type="text"
                  value={formTitleNe}
                  onChange={(e) => setFormTitleNe(e.target.value)}
                  placeholder={t.notices.titleNePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.notices.content} *
                </label>
                <textarea
                  required
                  value={formContent}
                  onChange={(e) => setFormContent(e.target.value)}
                  placeholder={t.notices.contentPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors min-h-[80px] resize-none"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.notices.contentNe}
                </label>
                <textarea
                  value={formContentNe}
                  onChange={(e) => setFormContentNe(e.target.value)}
                  placeholder={t.notices.contentNePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors min-h-[80px] resize-none"
                />
              </div>
              <div className="flex items-center gap-3">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={formIsActive}
                    onChange={(e) => setFormIsActive(e.target.checked)}
                    className="sr-only"
                  />
                  <span
                    className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors ${
                      formIsActive ? "bg-accent" : "bg-white/[0.1]"
                    }`}
                  >
                    <span
                      className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white transition-transform ${
                        formIsActive ? "translate-x-4" : "translate-x-0.5"
                      }`}
                    />
                  </span>
                  <span className="text-sm text-fg">
                    {t.notices.active}
                  </span>
                </label>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={formLoading}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {formLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.common.save}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
