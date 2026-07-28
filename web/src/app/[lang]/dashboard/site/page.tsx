"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import {
  AlertTriangle,
  ExternalLink,
  Eye,
  EyeOff,
  GripVertical,
  Loader2,
  Plus,
  Trash2,
} from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import type { PageSummary, SiteSettings } from "@/types/site";
import type { SiteTemplateOption } from "@/types/site-admin";
import { BASE_DOMAIN } from "@/lib/subdomain";

/**
 * Website overview: publish state, template, and the page list.
 *
 * The section editor lives one level down, per page. This screen is the thing
 * an owner opens to answer "is my site up, and what is on it".
 */
export default function SitePage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;

  const [settings, setSettings] = useState<SiteSettings | null>(null);
  const [pages, setPages] = useState<PageSummary[]>([]);
  const [templates, setTemplates] = useState<SiteTemplateOption[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const [templatePicker, setTemplatePicker] = useState(false);
  const [pendingTemplate, setPendingTemplate] = useState<string | null>(null);

  const fetchAll = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const [s, p, tpl] = await Promise.all([
        api.getSiteSettings(orgId),
        api.listSitePages(orgId),
        api.listSiteTemplates(),
      ]);
      setSettings(s);
      setPages(p.data ?? []);
      setTemplates(tpl.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  const siteUrl = settings?.slug
    ? `https://${settings.slug}.${BASE_DOMAIN}`
    : null;

  async function toggleLive() {
    if (!orgId || !settings) return;
    setBusy(true);
    setError(null);
    try {
      setSettings(await api.updateSiteSettings(orgId, { is_live: !settings.is_live }));
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setBusy(false);
    }
  }

  async function applyTemplate(templateId: string) {
    if (!orgId) return;
    setBusy(true);
    setError(null);
    try {
      const res = await api.applySiteTemplate(orgId, templateId);
      setPages(res.data ?? []);
      setSettings(await api.getSiteSettings(orgId));
      setTemplatePicker(false);
      setPendingTemplate(null);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setBusy(false);
    }
  }

  async function deletePage(page: PageSummary) {
    if (!orgId) return;
    if (!window.confirm(t.site.deletePageConfirm)) return;
    setBusy(true);
    try {
      await api.deleteSitePage(orgId, page.id);
      setPages((prev) => prev.filter((p) => p.id !== page.id));
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setBusy(false);
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Loader2 className="w-6 h-6 animate-spin text-fg-muted" />
      </div>
    );
  }

  return (
    <div className="max-w-5xl">
      <div className="flex flex-wrap items-start justify-between gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-fg">{t.site.title}</h1>
          <p className="text-sm text-fg-muted mt-1">{t.site.subtitle}</p>
        </div>
        {siteUrl && settings?.is_live ? (
          <a
            href={siteUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-white/[0.06] text-sm text-fg hover:bg-surface transition-colors"
          >
            {t.site.livePreview}
            <ExternalLink className="w-3.5 h-3.5" />
          </a>
        ) : null}
      </div>

      {error && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Publish state */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card mb-5">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2.5">
              <span
                className={`w-2 h-2 rounded-full ${settings?.is_live ? "bg-emerald-400" : "bg-amber-400"}`}
              />
              <span className="text-sm font-medium text-fg">
                {settings?.is_live ? t.site.live : t.site.draft}
              </span>
            </div>
            <p className="text-xs text-fg-muted mt-1.5">
              {settings?.is_live ? t.site.liveHint : t.site.draftHint}
            </p>
            {siteUrl && (
              <p className="text-xs font-mono text-fg-muted mt-1">{siteUrl}</p>
            )}
          </div>
          <button
            type="button"
            onClick={toggleLive}
            disabled={busy}
            className={`inline-flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-colors disabled:opacity-50 ${
              settings?.is_live
                ? "border border-white/[0.06] text-fg hover:bg-surface"
                : "bg-accent text-accent-fg hover:opacity-90"
            }`}
          >
            {settings?.is_live ? (
              <>
                <EyeOff className="w-4 h-4" />
                {t.site.takeOffline}
              </>
            ) : (
              <>
                <Eye className="w-4 h-4" />
                {t.site.goLive}
              </>
            )}
          </button>
        </div>
      </div>

      {/* Template */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card mb-5">
        <div className="flex items-center justify-between gap-4 mb-4">
          <div>
            <h2 className="text-sm font-medium text-fg">{t.site.template}</h2>
            <p className="text-xs text-fg-muted mt-1">
              {templates.find((tpl) => tpl.id === settings?.template)?.name ??
                settings?.template}
            </p>
          </div>
          <button
            type="button"
            onClick={() => setTemplatePicker((v) => !v)}
            className="px-4 py-2 rounded-xl border border-white/[0.06] text-sm text-fg hover:bg-surface transition-colors"
          >
            {t.site.changeTemplate}
          </button>
        </div>

        {templatePicker && (
          <>
            <div className="flex items-start gap-2.5 px-4 py-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-xs text-amber-300 mb-4">
              <AlertTriangle className="w-4 h-4 shrink-0 mt-0.5" />
              <span>{t.site.templateWarning}</span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {templates.map((tpl) => (
                <div
                  key={tpl.id}
                  className={`rounded-xl border p-4 transition-colors ${
                    pendingTemplate === tpl.id
                      ? "border-accent"
                      : "border-white/[0.06] hover:border-white/20"
                  }`}
                >
                  {/* Swatch built from the template's own tokens, so the choice
                      is made on how it looks rather than on its name. */}
                  <div
                    className="h-16 rounded-lg mb-3 flex items-end p-2 gap-1.5"
                    style={{ background: tpl.theme.bg, border: `1px solid ${tpl.theme.border}` }}
                  >
                    <span
                      className="h-2.5 w-10 rounded-full"
                      style={{ background: tpl.theme.accent }}
                    />
                    <span
                      className="h-2.5 w-6 rounded-full"
                      style={{ background: tpl.theme.fg_muted }}
                    />
                  </div>
                  <p className="text-sm text-fg">
                    {locale === "ne" && tpl.name_ne ? tpl.name_ne : tpl.name}
                  </p>
                  <button
                    type="button"
                    onClick={() =>
                      pendingTemplate === tpl.id
                        ? applyTemplate(tpl.id)
                        : setPendingTemplate(tpl.id)
                    }
                    disabled={busy}
                    className={`mt-3 w-full px-3 py-2 rounded-lg text-xs font-medium transition-colors disabled:opacity-50 ${
                      pendingTemplate === tpl.id
                        ? "bg-accent text-accent-fg"
                        : "border border-white/[0.06] text-fg-muted hover:text-fg"
                    }`}
                  >
                    {pendingTemplate === tpl.id ? t.site.applyTemplate : t.site.changeTemplate}
                  </button>
                </div>
              ))}
            </div>
          </>
        )}
      </div>

      {/* Pages */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        <div className="flex items-center justify-between gap-4 mb-4">
          <h2 className="text-sm font-medium text-fg">{t.site.pages}</h2>
          <Link
            href={`/${locale}/dashboard/site/pages/new`}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-accent-fg text-sm font-medium hover:opacity-90 transition-opacity"
          >
            <Plus className="w-4 h-4" />
            {t.site.addPage}
          </Link>
        </div>

        {pages.length === 0 ? (
          <p className="py-8 text-center text-sm text-fg-muted">{t.site.noPages}</p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {[...pages]
              .sort((a, b) => a.sort_order - b.sort_order)
              .map((page) => (
                <div key={page.id} className="flex items-center gap-3 py-3">
                  <GripVertical className="w-4 h-4 text-fg-muted shrink-0" />
                  <Link
                    href={`/${locale}/dashboard/site/pages/${page.id}`}
                    className="min-w-0 flex-1 group"
                  >
                    <p className="text-sm text-fg truncate group-hover:text-accent transition-colors">
                      {locale === "ne" && page.title_ne ? page.title_ne : page.title}
                    </p>
                    <p className="text-xs font-mono text-fg-muted truncate">
                      /{page.slug === "home" ? "" : page.slug}
                    </p>
                  </Link>
                  {!page.is_published && (
                    <span className="shrink-0 px-2 py-0.5 rounded-full text-[10px] font-mono uppercase bg-amber-500/10 text-amber-400 border border-amber-500/20">
                      {t.site.draft}
                    </span>
                  )}
                  <button
                    type="button"
                    onClick={() => deletePage(page)}
                    disabled={busy}
                    aria-label={t.site.deletePage}
                    className="shrink-0 p-2 rounded-lg text-fg-muted hover:text-red-400 hover:bg-surface transition-colors disabled:opacity-50"
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
