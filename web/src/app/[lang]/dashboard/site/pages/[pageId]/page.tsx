"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeft,
  ChevronDown,
  ChevronUp,
  Eye,
  EyeOff,
  Loader2,
  Plus,
  Trash2,
} from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import type { Section, SectionType, SitePage } from "@/types/site";
import { SECTION_SPECS } from "@/lib/site/section-specs";

/**
 * Section editor for one page.
 *
 * The whole page is edited locally and saved in one PUT: the API validates the
 * section array as a unit, so a partial save would be rejected anyway, and
 * batching means reordering ten sections is one request rather than ten.
 *
 * `new` as the page id creates a page instead of loading one — the two forms
 * are otherwise identical, and duplicating the screen to add four fields is
 * how they drift apart.
 */

const BACKGROUNDS = ["base", "surface", "accent", "none"] as const;
const PADDINGS = ["none", "sm", "md", "lg"] as const;
const WIDTHS = ["full", "contained", "narrow"] as const;
const ALIGNS = ["left", "center", "right"] as const;

/** Section ids only need to be unique within a page. */
function newSectionId(existing: Section[]): string {
  let n = existing.length + 1;
  const taken = new Set(existing.map((s) => s.id));
  while (taken.has(`s${n}`)) n += 1;
  return `s${n}`;
}

function blankSection(type: SectionType, existing: Section[]): Section {
  const spec = SECTION_SPECS[type];
  return {
    id: newSectionId(existing),
    type,
    variant: spec.variants[0],
    content: structuredClone(spec.defaultContent),
    style: { ...spec.defaultStyle },
  };
}

export default function SitePageEditor({
  params,
}: {
  params: { lang: string; pageId: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const router = useRouter();
  const { user } = useAuth();
  const orgId = user?.org_id;

  const isNew = params.pageId === "new";

  const [loading, setLoading] = useState(!isNew);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [dirty, setDirty] = useState(false);
  const [justSaved, setJustSaved] = useState(false);

  const [slug, setSlug] = useState("");
  const [title, setTitle] = useState("");
  const [titleNe, setTitleNe] = useState("");
  const [seoDescription, setSeoDescription] = useState("");
  const [isPublished, setIsPublished] = useState(true);
  const [showInNav, setShowInNav] = useState(true);
  const [sortOrder, setSortOrder] = useState(0);
  const [sections, setSections] = useState<Section[]>([]);

  const [openSection, setOpenSection] = useState<string | null>(null);
  const [addOpen, setAddOpen] = useState(false);
  /** Raw JSON text per section, so a half-typed edit is not thrown away. */
  const [contentDraft, setContentDraft] = useState<Record<string, string>>({});
  const [contentError, setContentError] = useState<Record<string, boolean>>({});

  const load = useCallback(async () => {
    if (!orgId || isNew) return;
    setLoading(true);
    try {
      const page: SitePage = await api.getSitePage(orgId, params.pageId);
      setSlug(page.slug);
      setTitle(page.title);
      setTitleNe(page.title_ne ?? "");
      setSeoDescription(page.seo_description ?? "");
      setIsPublished(page.is_published);
      setShowInNav(page.show_in_nav);
      setSortOrder(page.sort_order);
      setSections(page.sections ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, params.pageId, isNew, t.common.error]);

  useEffect(() => {
    load();
  }, [load]);

  const touch = useCallback(() => {
    setDirty(true);
    setJustSaved(false);
  }, []);

  function mutateSection(id: string, patch: Partial<Section>) {
    setSections((prev) => prev.map((s) => (s.id === id ? { ...s, ...patch } : s)));
    touch();
  }

  function moveSection(index: number, delta: number) {
    const target = index + delta;
    if (target < 0 || target >= sections.length) return;
    setSections((prev) => {
      const next = [...prev];
      [next[index], next[target]] = [next[target], next[index]];
      return next;
    });
    touch();
  }

  function removeSection(id: string) {
    setSections((prev) => prev.filter((s) => s.id !== id));
    touch();
  }

  function addSection(type: SectionType) {
    setSections((prev) => {
      const next = [...prev, blankSection(type, prev)];
      setOpenSection(next[next.length - 1].id);
      return next;
    });
    setAddOpen(false);
    touch();
  }

  /**
   * Commits a section's JSON draft. Invalid JSON is flagged but kept in the
   * textarea rather than reverted, so a typo does not discard the edit.
   */
  function commitContent(id: string, text: string) {
    setContentDraft((prev) => ({ ...prev, [id]: text }));
    try {
      const parsed = JSON.parse(text);
      if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
        throw new Error("not an object");
      }
      setContentError((prev) => ({ ...prev, [id]: false }));
      mutateSection(id, { content: parsed });
    } catch {
      setContentError((prev) => ({ ...prev, [id]: true }));
      setDirty(true);
    }
  }

  const hasContentError = useMemo(
    () => Object.values(contentError).some(Boolean),
    [contentError]
  );

  async function save() {
    if (!orgId) return;
    setSaving(true);
    setError(null);
    try {
      const body = {
        slug: slug.trim(),
        title: title.trim(),
        title_ne: titleNe.trim(),
        sections,
        seo_description: seoDescription.trim(),
        is_published: isPublished,
        show_in_nav: showInNav,
        sort_order: sortOrder,
      };

      if (isNew) {
        const created = await api.createSitePage(orgId, body);
        router.replace(`/${locale}/dashboard/site/pages/${created.id}`);
      } else {
        await api.updateSitePage(orgId, params.pageId, body);
      }
      setDirty(false);
      setJustSaved(true);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setSaving(false);
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

  return (
    <div className="max-w-4xl">
      <Link
        href={`/${locale}/dashboard/site`}
        className="inline-flex items-center gap-2 text-sm text-fg-muted hover:text-fg transition-colors mb-4"
      >
        <ArrowLeft className="w-4 h-4" />
        {t.site.title}
      </Link>

      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <h1 className="text-2xl font-semibold text-fg">
          {isNew ? t.site.newPage : t.site.editPage}
        </h1>
        <div className="flex items-center gap-3">
          {dirty && !hasContentError && (
            <span className="text-xs text-amber-400">{t.site.unsaved}</span>
          )}
          {justSaved && <span className="text-xs text-emerald-400">{t.site.saved}</span>}
          <button
            type="button"
            onClick={save}
            disabled={saving || hasContentError || !title.trim() || !slug.trim()}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-accent-fg text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
          >
            {saving && <Loader2 className="w-4 h-4 animate-spin" />}
            {t.site.save}
          </button>
        </div>
      </div>

      {error && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Page settings */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card mb-5">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label htmlFor="page-title" className={labelClass}>
              {t.site.pageTitle}
            </label>
            <input
              id="page-title"
              className={inputClass}
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                touch();
              }}
            />
          </div>
          <div>
            <label htmlFor="page-title-ne" className={labelClass}>
              {t.site.pageTitleNe}
            </label>
            <input
              id="page-title-ne"
              className={inputClass}
              value={titleNe}
              onChange={(e) => {
                setTitleNe(e.target.value);
                touch();
              }}
            />
          </div>
          <div>
            <label htmlFor="page-slug" className={labelClass}>
              {t.site.pageSlug}
            </label>
            <input
              id="page-slug"
              className={`${inputClass} font-mono`}
              value={slug}
              onChange={(e) => {
                setSlug(e.target.value);
                touch();
              }}
            />
            <p className="text-[11px] text-fg-muted mt-1">{t.site.slugHint}</p>
          </div>
          <div>
            <label htmlFor="page-seo" className={labelClass}>
              {t.site.seoDescription}
            </label>
            <input
              id="page-seo"
              className={inputClass}
              value={seoDescription}
              onChange={(e) => {
                setSeoDescription(e.target.value);
                touch();
              }}
            />
          </div>
        </div>

        <div className="flex flex-wrap gap-5 mt-4">
          <label className="flex items-center gap-2 text-sm text-fg cursor-pointer">
            <input
              type="checkbox"
              checked={isPublished}
              onChange={(e) => {
                setIsPublished(e.target.checked);
                touch();
              }}
              className="accent-accent"
            />
            {t.site.published}
          </label>
          <label className="flex items-center gap-2 text-sm text-fg cursor-pointer">
            <input
              type="checkbox"
              checked={showInNav}
              onChange={(e) => {
                setShowInNav(e.target.checked);
                touch();
              }}
              className="accent-accent"
            />
            {t.site.showInNav}
          </label>
        </div>
      </div>

      {/* Sections */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        <div className="flex items-center justify-between gap-4 mb-4">
          <h2 className="text-sm font-medium text-fg">{t.site.sections}</h2>
          <button
            type="button"
            onClick={() => setAddOpen((v) => !v)}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-white/[0.06] text-sm text-fg hover:bg-surface transition-colors"
          >
            <Plus className="w-4 h-4" />
            {t.site.addSection}
          </button>
        </div>

        {addOpen && (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2 mb-5 p-3 rounded-xl bg-black/20 border border-white/[0.06]">
            {(Object.keys(SECTION_SPECS) as SectionType[]).map((type) => (
              <button
                key={type}
                type="button"
                onClick={() => addSection(type)}
                className="px-3 py-2 rounded-lg text-left text-xs text-fg-muted hover:text-fg hover:bg-surface transition-colors"
              >
                {locale === "ne"
                  ? SECTION_SPECS[type].labelNe
                  : SECTION_SPECS[type].label}
              </button>
            ))}
          </div>
        )}

        {sections.length === 0 ? (
          <p className="py-8 text-center text-sm text-fg-muted">{t.site.noSections}</p>
        ) : (
          <div className="flex flex-col gap-2">
            {sections.map((section, index) => {
              const spec = SECTION_SPECS[section.type];
              const isOpen = openSection === section.id;
              return (
                <div
                  key={section.id}
                  className="rounded-xl border border-white/[0.06] overflow-hidden"
                >
                  <div className="flex items-center gap-2 px-3 py-2.5 bg-black/20">
                    <button
                      type="button"
                      onClick={() => setOpenSection(isOpen ? null : section.id)}
                      className="min-w-0 flex-1 flex items-center gap-2 text-left"
                      aria-expanded={isOpen}
                    >
                      {isOpen ? (
                        <ChevronUp className="w-4 h-4 text-fg-muted shrink-0" />
                      ) : (
                        <ChevronDown className="w-4 h-4 text-fg-muted shrink-0" />
                      )}
                      <span className="text-sm text-fg truncate">
                        {locale === "ne" ? spec?.labelNe : spec?.label ?? section.type}
                      </span>
                      <span className="text-[11px] font-mono text-fg-muted truncate">
                        {section.variant}
                      </span>
                      {section.hidden && (
                        <span className="px-2 py-0.5 rounded-full text-[10px] font-mono uppercase bg-white/5 text-fg-muted">
                          {t.site.hidden}
                        </span>
                      )}
                    </button>

                    <button
                      type="button"
                      onClick={() => moveSection(index, -1)}
                      disabled={index === 0}
                      aria-label={t.site.moveUp}
                      className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors disabled:opacity-30"
                    >
                      <ChevronUp className="w-4 h-4" />
                    </button>
                    <button
                      type="button"
                      onClick={() => moveSection(index, 1)}
                      disabled={index === sections.length - 1}
                      aria-label={t.site.moveDown}
                      className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors disabled:opacity-30"
                    >
                      <ChevronDown className="w-4 h-4" />
                    </button>
                    <button
                      type="button"
                      onClick={() => mutateSection(section.id, { hidden: !section.hidden })}
                      aria-label={section.hidden ? t.site.unhide : t.site.hide}
                      className="p-1.5 rounded-lg text-fg-muted hover:text-fg hover:bg-surface transition-colors"
                    >
                      {section.hidden ? (
                        <EyeOff className="w-4 h-4" />
                      ) : (
                        <Eye className="w-4 h-4" />
                      )}
                    </button>
                    <button
                      type="button"
                      onClick={() => removeSection(section.id)}
                      aria-label={t.site.removeSection}
                      className="p-1.5 rounded-lg text-fg-muted hover:text-red-400 hover:bg-surface transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>

                  {isOpen && (
                    <div className="p-4 flex flex-col gap-4">
                      <div className="grid grid-cols-2 lg:grid-cols-5 gap-3">
                        <div>
                          <label className={labelClass}>{t.site.variant}</label>
                          <select
                            value={section.variant}
                            onChange={(e) =>
                              mutateSection(section.id, { variant: e.target.value })
                            }
                            className={inputClass}
                          >
                            {spec?.variants.map((v) => (
                              <option key={v} value={v}>
                                {v}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <label className={labelClass}>{t.site.background}</label>
                          <select
                            value={section.style?.background ?? "base"}
                            onChange={(e) =>
                              mutateSection(section.id, {
                                style: { ...section.style, background: e.target.value as never },
                              })
                            }
                            className={inputClass}
                          >
                            {BACKGROUNDS.map((v) => (
                              <option key={v} value={v}>
                                {v}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <label className={labelClass}>{t.site.padding}</label>
                          <select
                            value={section.style?.padding ?? "md"}
                            onChange={(e) =>
                              mutateSection(section.id, {
                                style: { ...section.style, padding: e.target.value as never },
                              })
                            }
                            className={inputClass}
                          >
                            {PADDINGS.map((v) => (
                              <option key={v} value={v}>
                                {v}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <label className={labelClass}>{t.site.width}</label>
                          <select
                            value={section.style?.width ?? "contained"}
                            onChange={(e) =>
                              mutateSection(section.id, {
                                style: { ...section.style, width: e.target.value as never },
                              })
                            }
                            className={inputClass}
                          >
                            {WIDTHS.map((v) => (
                              <option key={v} value={v}>
                                {v}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <label className={labelClass}>{t.site.align}</label>
                          <select
                            value={section.style?.align ?? "left"}
                            onChange={(e) =>
                              mutateSection(section.id, {
                                style: { ...section.style, align: e.target.value as never },
                              })
                            }
                            className={inputClass}
                          >
                            {ALIGNS.map((v) => (
                              <option key={v} value={v}>
                                {v}
                              </option>
                            ))}
                          </select>
                        </div>
                      </div>

                      <div>
                        <label
                          htmlFor={`content-${section.id}`}
                          className={labelClass}
                        >
                          {t.site.content}
                        </label>
                        <textarea
                          id={`content-${section.id}`}
                          rows={12}
                          spellCheck={false}
                          value={
                            contentDraft[section.id] ??
                            JSON.stringify(section.content, null, 2)
                          }
                          onChange={(e) => commitContent(section.id, e.target.value)}
                          className={`${inputClass} font-mono text-xs leading-relaxed ${
                            contentError[section.id] ? "border-red-500/50" : ""
                          }`}
                        />
                        <p className="text-[11px] text-fg-muted mt-1">
                          {contentError[section.id]
                            ? t.site.invalidJson
                            : t.site.contentHint}
                        </p>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
