"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, X, Globe } from "lucide-react";
import type { Locale } from "@/types";
import type { PublicSite } from "@/types/site";
import { text } from "@/lib/site/content";

/**
 * Tenant site navigation.
 *
 * Links are relative to the subdomain root: on ibckirtipur.nepalgym.xyz a page
 * lives at /classes (or /ne/classes), never at /site/ibckirtipur/classes — the
 * middleware rewrite is invisible to the visitor and must stay that way.
 */

function pageHref(locale: Locale, pageSlug: string): string {
  const path = pageSlug === "home" ? "" : `/${pageSlug}`;
  return locale === "en" ? path || "/" : `/ne${path}`;
}

export function SiteHeader({
  site,
  locale,
  currentPage,
  apexUrl,
}: {
  site: PublicSite;
  locale: Locale;
  currentPage: string;
  apexUrl: string;
}) {
  const [open, setOpen] = useState(false);

  const navPages = site.pages
    .filter((p) => p.show_in_nav)
    .sort((a, b) => a.sort_order - b.sort_order);

  const orgName =
    locale === "ne" && site.org_name_ne ? site.org_name_ne : site.org_name;

  const extraLinks = site.nav?.links ?? [];
  const ctaLabel = text(site.nav as Record<string, unknown>, "ctaLabel", locale);
  const ctaHref = site.nav?.ctaHref;

  const otherLocale: Locale = locale === "en" ? "ne" : "en";
  const otherLocaleHref = pageHref(otherLocale, currentPage);

  return (
    <header className="sticky top-0 z-50 bg-[var(--site-bg)]/90 backdrop-blur-md border-b border-[var(--site-border)]">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between gap-4">
        <Link
          href={pageHref(locale, "home")}
          className="text-lg sm:text-xl truncate"
          style={{
            fontFamily: "var(--site-font-display)",
            textTransform: "var(--site-display-transform)" as never,
            fontWeight: "var(--site-display-weight)" as never,
            letterSpacing: "var(--site-display-tracking)",
          }}
        >
          {orgName}
        </Link>

        <nav className="hidden md:flex items-center gap-6">
          {navPages.map((p) => (
            <Link
              key={p.id}
              href={pageHref(locale, p.slug)}
              aria-current={p.slug === currentPage ? "page" : undefined}
              className={`text-sm transition-opacity hover:opacity-100 ${
                p.slug === currentPage ? "opacity-100" : "opacity-70"
              }`}
            >
              {locale === "ne" && p.title_ne ? p.title_ne : p.title}
            </Link>
          ))}
          {extraLinks.map((l, i) => (
            <a
              key={`${l.href}-${i}`}
              href={l.href}
              className="text-sm opacity-70 hover:opacity-100 transition-opacity"
            >
              {locale === "ne" && l.labelNe ? l.labelNe : l.label}
            </a>
          ))}
        </nav>

        <div className="hidden md:flex items-center gap-3">
          <Link
            href={otherLocaleHref}
            className="flex items-center gap-1.5 text-sm opacity-70 hover:opacity-100 transition-opacity"
            aria-label={otherLocale === "en" ? "Switch to English" : "नेपालीमा हेर्नुहोस्"}
          >
            <Globe className="w-4 h-4" />
            {otherLocale === "en" ? "EN" : "ने"}
          </Link>
          {ctaLabel && ctaHref ? (
            <a
              href={ctaHref}
              className="px-4 py-2 text-sm bg-[var(--site-accent)] text-[var(--site-accent-fg)] rounded-[var(--site-radius)]"
            >
              {ctaLabel}
            </a>
          ) : null}
          {/* The authenticated app lives on the apex, never on a tenant host. */}
          <a href={apexUrl} className="text-sm opacity-70 hover:opacity-100 transition-opacity">
            {locale === "ne" ? "सदस्य लगइन" : "Member Login"}
          </a>
        </div>

        <button
          type="button"
          className="md:hidden p-2 -mr-2"
          onClick={() => setOpen((v) => !v)}
          aria-expanded={open}
          aria-label={open ? "Close menu" : "Open menu"}
        >
          {open ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
        </button>
      </div>

      {open ? (
        <div className="md:hidden border-t border-[var(--site-border)] bg-[var(--site-bg)]">
          <nav className="px-4 py-4 flex flex-col gap-1">
            {navPages.map((p) => (
              <Link
                key={p.id}
                href={pageHref(locale, p.slug)}
                onClick={() => setOpen(false)}
                aria-current={p.slug === currentPage ? "page" : undefined}
                className="py-2.5 text-sm"
              >
                {locale === "ne" && p.title_ne ? p.title_ne : p.title}
              </Link>
            ))}
            {extraLinks.map((l, i) => (
              <a key={`${l.href}-${i}`} href={l.href} className="py-2.5 text-sm">
                {locale === "ne" && l.labelNe ? l.labelNe : l.label}
              </a>
            ))}
            <div className="mt-2 pt-3 border-t border-[var(--site-border)] flex items-center gap-4">
              <Link href={otherLocaleHref} className="py-2 text-sm flex items-center gap-1.5">
                <Globe className="w-4 h-4" />
                {otherLocale === "en" ? "English" : "नेपाली"}
              </Link>
              <a href={apexUrl} className="py-2 text-sm opacity-70">
                {locale === "ne" ? "सदस्य लगइन" : "Member Login"}
              </a>
            </div>
          </nav>
        </div>
      ) : null}
    </header>
  );
}
