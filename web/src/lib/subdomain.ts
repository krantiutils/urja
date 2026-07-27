/**
 * Subdomain resolution for tenant gym sites.
 *
 * A gym's organization `slug` is its subdomain: the gym with slug
 * `ibckirtipur` is served at `ibckirtipur.nepalgym.xyz`. There is no separate
 * identity concept — the slug column already carries a UNIQUE constraint.
 *
 * These are pure functions so they can be unit tested without a request.
 */

import { locales, defaultLocale } from "@/lib/i18n";

/** Base domain, overridable so staging and local development work unchanged. */
export const BASE_DOMAIN =
  process.env.NEXT_PUBLIC_BASE_DOMAIN?.trim() || "nepalgym.xyz";

/**
 * Labels that are never a tenant. `www` and the apex serve the marketing site
 * and the authenticated app; `api` and `admin` are reserved so a gym can never
 * claim a slug that would shadow infrastructure.
 */
const RESERVED_LABELS = new Set(["www", "api", "admin", "static", "assets"]);

/** Strips the port and lowercases a Host header value. */
function normalizeHost(host: string): string {
  return host.split(":")[0].trim().toLowerCase();
}

/**
 * Extracts the tenant slug from a Host header, or null when the host is the
 * apex, a reserved label, or an unrelated domain.
 *
 * Accepts `<slug>.localhost` too, so `ibckirtipur.localhost:3000` works in
 * development without touching /etc/hosts on most systems.
 */
export function extractSubdomain(host: string | null | undefined): string | null {
  if (!host) return null;

  const h = normalizeHost(host);
  if (!h) return null;

  let label: string | null = null;

  if (h.endsWith(`.${BASE_DOMAIN}`)) {
    label = h.slice(0, -(BASE_DOMAIN.length + 1));
  } else if (h.endsWith(".localhost")) {
    label = h.slice(0, -".localhost".length);
  } else {
    return null;
  }

  // A multi-level label (`a.b.nepalgym.xyz`) is not a tenant — we serve exactly
  // one level of subdomain, and accepting more would make slugs ambiguous.
  if (!label || label.includes(".")) return null;
  if (RESERVED_LABELS.has(label)) return null;

  // Must look like a slug. Anything else cannot match an organization anyway,
  // and rejecting here keeps junk out of the rewrite path.
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(label)) return null;

  return label;
}

export interface SiteRewrite {
  slug: string;
  locale: string;
  /** Page slug with no leading slash; "home" when the path is empty. */
  pagePath: string;
}

/**
 * Paths that must never be rewritten to a tenant page: the authenticated app
 * lives on the apex host only. A visitor who lands on one of these via a
 * subdomain is redirected to the apex instead.
 */
const APP_ONLY_PREFIXES = ["dashboard", "member", "super-admin", "login", "onboarding"];

/** Reports whether a path belongs to the authenticated app rather than a site. */
export function isAppOnlyPath(pathname: string): boolean {
  const segments = pathname.split("/").filter(Boolean);
  if (segments.length === 0) return false;
  // Skip a leading locale segment if present.
  const first = (locales as readonly string[]).includes(segments[0])
    ? segments[1]
    : segments[0];
  return !!first && APP_ONLY_PREFIXES.includes(first);
}

/**
 * Resolves a tenant request to the internal Next.js route parameters.
 *
 * `ibckirtipur.nepalgym.xyz/`            -> { slug, locale: "en", pagePath: "home" }
 * `ibckirtipur.nepalgym.xyz/ne/coaches`  -> { slug, locale: "ne", pagePath: "coaches" }
 * `nepalgym.xyz/en/dashboard`            -> null (not a tenant host)
 */
export function resolveSiteRewrite(
  host: string | null | undefined,
  pathname: string
): SiteRewrite | null {
  const slug = extractSubdomain(host);
  if (!slug) return null;

  const segments = pathname.split("/").filter(Boolean);

  let locale = defaultLocale as string;
  let rest = segments;
  if (segments.length > 0 && (locales as readonly string[]).includes(segments[0])) {
    locale = segments[0];
    rest = segments.slice(1);
  }

  const pagePath = rest.length > 0 ? rest.join("/") : "home";

  return { slug, locale, pagePath };
}

/** Builds the internal route a tenant request rewrites to. */
export function siteRewritePath(r: SiteRewrite): string {
  return `/${r.locale}/site/${r.slug}/${r.pagePath}`;
}

/**
 * Builds the absolute apex URL used to bounce app paths off a tenant host.
 * The protocol is taken from the incoming request rather than assumed, so a
 * local http:// session is not redirected to an https:// URL that will not
 * resolve.
 */
export function apexUrl(pathname: string, protocol: string): string {
  const scheme = protocol.replace(/:$/, "");
  return `${scheme}://${BASE_DOMAIN}${pathname}`;
}
