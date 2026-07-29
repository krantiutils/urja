import type { MetadataRoute } from "next";
import { headers } from "next/headers";
import { BASE_DOMAIN, extractSubdomain } from "@/lib/subdomain";
import { locales, defaultLocale } from "@/lib/i18n";

/**
 * Sitemap, resolved per host.
 *
 * This app serves two different things from one deployment: the apex is the
 * product's marketing site, and every `<slug>.nepalgym.xyz` is a gym's own
 * public website. A single static sitemap would be wrong for both — a gym's
 * sitemap must list that gym's pages at that gym's hostname, or Google treats
 * the URLs as belonging to a site it wasn't asked to crawl and drops them.
 *
 * Reading the Host header makes this route dynamic, which is correct: page
 * lists change whenever a gym publishes, and a build-time snapshot would go
 * stale the moment they do.
 */

const API_BASE = process.env.NEXT_PUBLIC_API_URL || `https://${BASE_DOMAIN}`;

interface SitePage {
  slug: string;
  is_published: boolean;
}

async function fetchPages(gymSlug: string): Promise<SitePage[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/sites/${gymSlug}`, {
      // Fresh enough that a newly published page appears within the hour,
      // cheap enough that a crawler cannot hammer the API.
      next: { revalidate: 3600 },
    });
    if (!res.ok) return [];
    const site = (await res.json()) as { pages?: SitePage[] };
    return (site.pages ?? []).filter((p) => p.is_published);
  } catch {
    return [];
  }
}

/** Every locale variant of one path, cross-linked so Google knows they are translations. */
function localeEntries(origin: string, path: string, lastModified: Date) {
  const href = (locale: string) =>
    locale === defaultLocale
      ? `${origin}${path === "" ? "" : `/${path}`}`
      : `${origin}/${locale}${path === "" ? "" : `/${path}`}`;

  const languages = Object.fromEntries(locales.map((l) => [l, href(l)]));

  return locales.map((locale) => ({
    url: href(locale),
    lastModified,
    changeFrequency: "weekly" as const,
    priority: path === "" ? 1 : 0.7,
    alternates: { languages },
  }));
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const host = headers().get("host");
  const gymSlug = extractSubdomain(host);
  const proto = headers().get("x-forwarded-proto") ?? "https";
  const origin = `${proto}://${host ?? BASE_DOMAIN}`;
  const now = new Date();

  // Apex: the product's own pages. A gym's pages do not belong here — they are
  // published under their own hostname and listed by that host's sitemap.
  if (!gymSlug) {
    return ["", "privacy"].flatMap((p) => localeEntries(origin, p, now));
  }

  const pages = await fetchPages(gymSlug);
  return pages.flatMap((p) =>
    localeEntries(origin, p.slug === "home" ? "" : p.slug, now)
  );
}
