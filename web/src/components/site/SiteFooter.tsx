import Link from "next/link";
import type { Locale } from "@/types";
import type { PublicSite } from "@/types/site";
import { text } from "@/lib/site/content";

/** Social platforms rendered in the footer, in a stable order. */
const SOCIAL_KEYS = [
  ["facebook", "Facebook"],
  ["instagram", "Instagram"],
  ["youtube", "YouTube"],
  ["tiktok", "TikTok"],
  ["whatsapp", "WhatsApp"],
] as const;

function pageHref(locale: Locale, pageSlug: string): string {
  const path = pageSlug === "home" ? "" : `/${pageSlug}`;
  return locale === "en" ? path || "/" : `/ne${path}`;
}

export function SiteFooter({
  site,
  locale,
}: {
  site: PublicSite;
  locale: Locale;
}) {
  const orgName =
    locale === "ne" && site.org_name_ne ? site.org_name_ne : site.org_name;

  const tagline = text(site.footer as Record<string, unknown>, "tagline", locale);
  const note = text(site.footer as Record<string, unknown>, "note", locale);

  const socials = SOCIAL_KEYS.filter(
    ([key]) => typeof site.socials?.[key] === "string" && site.socials[key] !== ""
  );

  const navPages = site.pages
    .filter((p) => p.show_in_nav)
    .sort((a, b) => a.sort_order - b.sort_order);

  return (
    <footer className="border-t border-[var(--site-border)] bg-[var(--site-bg)] text-[var(--site-fg)]">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 py-12">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-8">
          <div>
            <div
              className="text-lg mb-3"
              style={{
                fontFamily: "var(--site-font-display)",
                textTransform: "var(--site-display-transform)" as never,
                fontWeight: "var(--site-display-weight)" as never,
              }}
            >
              {orgName}
            </div>
            {tagline ? (
              <p className="text-sm text-[var(--site-fg-muted)] leading-relaxed">
                {tagline}
              </p>
            ) : null}
          </div>

          <div>
            <nav className="flex flex-col gap-2">
              {navPages.map((p) => (
                <Link
                  key={p.id}
                  href={pageHref(locale, p.slug)}
                  className="text-sm text-[var(--site-fg-muted)] hover:text-[var(--site-fg)] transition-colors"
                >
                  {locale === "ne" && p.title_ne ? p.title_ne : p.title}
                </Link>
              ))}
            </nav>
          </div>

          <div>
            {socials.length > 0 ? (
              <div className="flex flex-wrap gap-4">
                {socials.map(([key, label]) => (
                  <a
                    key={key}
                    href={site.socials[key]}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-[var(--site-fg-muted)] hover:text-[var(--site-fg)] transition-colors"
                  >
                    {label}
                  </a>
                ))}
              </div>
            ) : null}
            {note ? (
              <p className="mt-4 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                {note}
              </p>
            ) : null}
          </div>
        </div>

        <div className="mt-10 pt-6 border-t border-[var(--site-border)] flex flex-col sm:flex-row gap-2 sm:items-center sm:justify-between">
          <p className="text-xs text-[var(--site-fg-muted)]">
            &copy; {new Date().getFullYear()} {orgName}
          </p>
          <p className="text-xs text-[var(--site-fg-muted)]">
            {locale === "ne" ? "Urja द्वारा संचालित" : "Powered by Urja"}
          </p>
        </div>
      </div>
    </footer>
  );
}
