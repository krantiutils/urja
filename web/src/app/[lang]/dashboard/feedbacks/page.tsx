"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, Feedback } from "@/types";
import { MessageCircle, Trash2, Loader2, User } from "lucide-react";

export default function FeedbacksPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const orgId = user?.org_id;

  const fetchFeedbacks = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listFeedbacks(orgId, { limit: 100 });
      setFeedbacks(res.data ?? []);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchFeedbacks();
  }, [fetchFeedbacks]);

  const handleDelete = async (id: string) => {
    if (!orgId) return;
    if (!confirm(t.feedbacks.confirmDelete)) return;
    try {
      await api.deleteFeedback(orgId, id);
      fetchFeedbacks();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <MessageCircle className="w-5 h-5 text-accent" />
          {t.feedbacks.title}
        </h1>
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Feedback List */}
      {loading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : feedbacks.length === 0 ? (
        <div className="text-center py-16 text-fg-muted text-sm">
          {t.feedbacks.noFeedbacks}
        </div>
      ) : (
        <div className="space-y-4">
          {feedbacks.map((fb) => (
            <div
              key={fb.id}
              className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card"
            >
              <div className="flex items-start gap-4">
                {/* Avatar */}
                {fb.avatar_url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={fb.avatar_url}
                    alt={fb.member_name}
                    className="w-10 h-10 rounded-full object-cover"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0">
                    <User className="w-5 h-5 text-accent" />
                  </div>
                )}

                {/* Content */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between gap-3">
                    <p className="text-fg font-medium text-sm">
                      {fb.member_name}
                    </p>
                    <p className="text-xs text-fg-muted font-mono">
                      {new Date(fb.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <p className="text-fg-muted text-sm mt-2">{fb.message}</p>
                </div>

                {/* Delete button */}
                <button
                  onClick={() => handleDelete(fb.id)}
                  title={t.common.delete}
                  className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
