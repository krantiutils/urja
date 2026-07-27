/**
 * Server-side data access for public tenant sites.
 *
 * These run in React Server Components, so they call the API directly rather
 * than going through `lib/api.ts` (which reads a bearer token from
 * localStorage and only exists in the browser). Public site endpoints need no
 * authentication.
 */

import type { PublicSite, SitePage } from "@/types/site";

/**
 * Base URL for server-to-server calls. Inside Docker the web container reaches
 * the API over the internal network, so this is separate from the browser's
 * NEXT_PUBLIC_API_URL.
 */
function apiBase(): string {
  return (
    process.env.INTERNAL_API_URL?.trim() ||
    process.env.NEXT_PUBLIC_API_URL?.trim() ||
    "http://localhost:8080"
  );
}

/** How long a rendered tenant page stays cached before revalidating. */
export const SITE_REVALIDATE_SECONDS = 60;

/**
 * Fetches a live site's settings, nav and published page index.
 * Returns null when the gym does not exist or has not gone live — the caller
 * renders not-found, so an unfinished site is never discoverable.
 */
export async function fetchSite(slug: string): Promise<PublicSite | null> {
  try {
    const res = await fetch(`${apiBase()}/api/v1/sites/${encodeURIComponent(slug)}`, {
      next: { revalidate: SITE_REVALIDATE_SECONDS },
    });
    if (!res.ok) return null;
    return (await res.json()) as PublicSite;
  } catch {
    // A site that cannot be reached renders as not-found rather than a crash:
    // an API blip should not take down every gym's website with a 500.
    return null;
  }
}

/** Fetches one published page. Returns null for unknown or unpublished pages. */
export async function fetchPage(
  slug: string,
  pageSlug: string
): Promise<SitePage | null> {
  try {
    const res = await fetch(
      `${apiBase()}/api/v1/sites/${encodeURIComponent(slug)}/pages/${encodeURIComponent(pageSlug)}`,
      { next: { revalidate: SITE_REVALIDATE_SECONDS } }
    );
    if (!res.ok) return null;
    return (await res.json()) as SitePage;
  } catch {
    return null;
  }
}

/**
 * Fetches a gym's active membership packages, for membership_plans sections
 * using dataSource "auto". Failure yields an empty list so the rest of the page
 * still renders.
 */
export async function fetchPackages(
  orgSlug: string
): Promise<Array<Record<string, unknown>>> {
  try {
    const res = await fetch(
      `${apiBase()}/api/v1/packages?gym=${encodeURIComponent(orgSlug)}`,
      { next: { revalidate: SITE_REVALIDATE_SECONDS } }
    );
    if (!res.ok) return [];
    const body = await res.json();
    return Array.isArray(body?.data) ? body.data : [];
  } catch {
    return [];
  }
}
