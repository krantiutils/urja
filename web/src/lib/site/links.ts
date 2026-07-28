/**
 * Link construction for tenant sites.
 *
 * Links are relative to the subdomain root: on ibckirtipur.nepalgym.xyz a page
 * lives at /classes (or /ne/classes), never at /site/ibckirtipur/classes — the
 * middleware rewrite is invisible to the visitor and must stay that way.
 */

import type { Locale } from "@/types";

/** Path to a page of the current site, in the given locale. */
export function pageHref(locale: Locale, pageSlug: string): string {
  const path = pageSlug === "home" ? "" : `/${pageSlug}`;
  return locale === "en" ? path || "/" : `/ne${path}`;
}

/**
 * True for links that leave the site, which must never be rewritten:
 * absolute URLs, mail/phone/SMS handoffs, and in-page anchors.
 */
function isExternal(href: string): boolean {
  return /^([a-z][a-z0-9+.-]*:|\/\/|#)/i.test(href);
}

/**
 * Localizes an internal href written by an admin or stored in section content.
 *
 * Content is authored once and read in both languages, so a CTA saved as
 * "/contact" would otherwise send every Nepali visitor to the English page.
 * Anything already carrying a locale prefix, and anything pointing off-site,
 * is returned untouched.
 */
export function siteHref(locale: Locale, href: string): string {
  const target = href.trim();
  if (target === "") return "#";
  if (isExternal(target)) return target;

  // Relative paths ("contact") are not something the builder produces, but a
  // hand-edited page can contain one; treat it as site-root relative.
  const path = target.startsWith("/") ? target : `/${target}`;

  if (locale === "en") return path;
  if (path === "/ne" || path.startsWith("/ne/")) return path;
  return path === "/" ? "/ne" : `/ne${path}`;
}
