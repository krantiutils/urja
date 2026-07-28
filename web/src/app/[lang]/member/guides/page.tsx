"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft, BookOpen, Loader2 } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import { parseMarkdownBlocks } from "@/lib/site/content";
import type { Locale, TrainingGuide } from "@/types";

/**
 * Training guides, as a member reads them.
 *
 * Guides are per-gym content, so the list is fetched for the member's own gym.
 * Content is rendered through the same restricted Markdown reader the tenant
 * site uses rather than interpolated as HTML — a guide is staff-authored, but
 * staff are not a trusted source of markup.
 */
export default function MemberGuidesPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [guides, setGuides] = useState<TrainingGuide[]>([]);
  const [open, setOpen] = useState<TrainingGuide | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const gym = user?.org_slug;

  const fetchGuides = useCallback(async () => {
    if (!gym) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const res = await api.listPublishedGuides(gym, { limit: 100 });
      setGuides(res.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [gym, t.common.error]);

  useEffect(() => {
    fetchGuides();
  }, [fetchGuides]);

  const title = (g: TrainingGuide) =>
    locale === "ne" && g.title_ne ? g.title_ne : g.title;
  const body = (g: TrainingGuide) =>
    locale === "ne" && g.content_ne ? g.content_ne : g.content;

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Loader2 className="w-6 h-6 animate-spin text-fg-muted" />
      </div>
    );
  }

  if (open) {
    return (
      <div className="max-w-2xl">
        <button
          type="button"
          onClick={() => setOpen(null)}
          className="inline-flex items-center gap-2 text-sm text-fg-muted hover:text-fg transition-colors mb-4"
        >
          <ArrowLeft className="w-4 h-4" />
          {t.guides.title}
        </button>

        <h1 className="text-2xl font-semibold text-fg">{title(open)}</h1>
        {open.category && (
          <p className="mt-1 text-xs font-mono uppercase tracking-widest text-accent">
            {open.category}
          </p>
        )}

        <div className="mt-6 flex flex-col gap-3">
          {parseMarkdownBlocks(body(open)).map((block, i) => {
            if (block.kind === "heading") {
              return block.level === 2 ? (
                <h2 key={i} className="mt-4 text-lg font-semibold text-fg">
                  {block.text}
                </h2>
              ) : (
                <h3 key={i} className="mt-3 font-medium text-fg">
                  {block.text}
                </h3>
              );
            }
            if (block.kind === "listItem") {
              return (
                <p key={i} className="pl-4 text-sm text-fg-muted leading-relaxed">
                  • {block.text}
                </p>
              );
            }
            return (
              <p key={i} className="text-sm text-fg-muted leading-relaxed">
                {block.text}
              </p>
            );
          })}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-2xl">
      <div>
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <BookOpen className="w-5 h-5 text-accent" />
          {t.guides.title}
        </h1>
        <p className="text-sm text-fg-muted mt-1">{t.guides.memberSubtitle}</p>
      </div>

      {error && (
        <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      {guides.length === 0 ? (
        <p className="py-10 text-center text-sm text-fg-muted">{t.guides.noneYet}</p>
      ) : (
        <div className="flex flex-col gap-2">
          {guides.map((g) => (
            <button
              key={g.id}
              type="button"
              onClick={() => setOpen(g)}
              className="text-left p-4 rounded-2xl bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] hover:border-accent/30 transition-colors"
            >
              <p className="text-sm text-fg">{title(g)}</p>
              {g.category && (
                <p className="mt-1 text-xs font-mono uppercase tracking-widest text-fg-muted">
                  {g.category}
                </p>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
