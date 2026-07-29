import { BASE_DOMAIN } from "@/lib/subdomain";
import { headers } from "next/headers";

/**
 * Sitemap index for the whole network, served from the apex.
 *
 * Each gym lives on its own subdomain and publishes its own `/sitemap.xml`.
 * This lists them all in one place so a single submission in Search Console
 * covers every gym, and so a crawler that finds nepalgym.xyz can discover the
 * tenant sites without being told each subdomain by hand.
 *
 * Cross-host entries in a sitemap index are only honoured when the submitting
 * property covers those hosts — so this must be submitted under the *domain*
 * property `nepalgym.xyz`, not the URL-prefix property `https://nepalgym.xyz/`.
 * A URL-prefix property does not include subdomains and Google will reject the
 * subdomain entries as out of scope.
 *
 * Served on the apex only: a gym's own subdomain has no business advertising
 * its competitors' sitemaps.
 */

export const revalidate = 3600;

const API_BASE = process.env.NEXT_PUBLIC_API_URL || `https://${BASE_DOMAIN}`;

interface LiveSite {
  slug: string;
  updated_at: string;
}

async function fetchLiveSites(): Promise<LiveSite[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/sites`, {
      next: { revalidate: 3600 },
    });
    if (!res.ok) return [];
    const body = (await res.json()) as { data?: LiveSite[] };
    return body.data ?? [];
  } catch {
    return [];
  }
}

function entry(loc: string, lastmod?: string): string {
  const when = lastmod ? `\n    <lastmod>${lastmod}</lastmod>` : "";
  return `  <sitemap>\n    <loc>${loc}</loc>${when}\n  </sitemap>`;
}

export async function GET() {
  const host = headers().get("host") ?? BASE_DOMAIN;
  const proto = headers().get("x-forwarded-proto") ?? "https";

  // Only the apex serves the index. A tenant host asking for it gets a 404
  // rather than a list of every other gym on the platform.
  const isApex =
    host.split(":")[0].toLowerCase() === BASE_DOMAIN ||
    host.split(":")[0].toLowerCase() === `www.${BASE_DOMAIN}`;
  if (!isApex) {
    return new Response("Not found", { status: 404 });
  }

  const sites = await fetchLiveSites();

  const body =
    `<?xml version="1.0" encoding="UTF-8"?>\n` +
    `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
    [
      entry(`${proto}://${BASE_DOMAIN}/sitemap.xml`),
      ...sites.map((s) =>
        entry(
          `${proto}://${s.slug}.${BASE_DOMAIN}/sitemap.xml`,
          new Date(s.updated_at).toISOString()
        )
      ),
    ].join("\n") +
    `\n</sitemapindex>\n`;

  return new Response(body, {
    headers: {
      "Content-Type": "application/xml; charset=utf-8",
      "Cache-Control": "public, max-age=3600, s-maxage=3600",
    },
  });
}
