import type { MetadataRoute } from "next";
import { headers } from "next/headers";
import { BASE_DOMAIN } from "@/lib/subdomain";

/**
 * robots.txt, resolved per host.
 *
 * Two things matter here beyond the usual.
 *
 * First, the authenticated app must never be crawled. `/dashboard`, `/member`
 * and the rest are behind a login, so a crawler only ever gets a redirect —
 * but indexing the URLs still surfaces them in results and wastes crawl budget
 * that should go to the gym's public pages.
 *
 * Second, assistant crawlers are allowed on purpose. A gym's discoverability
 * increasingly runs through answers rather than links: somebody asking an
 * assistant "where can I learn boxing in Kirtipur" should be able to find this
 * gym. Those crawlers are separate user-agents from Googlebot and several
 * default to *not* crawling unless named, so silence here reads as refusal.
 * Naming them is the whole point — the marketing site and the gym sites are
 * public information we want quoted.
 */

const ASSISTANT_CRAWLERS = [
  "GPTBot", // OpenAI — ChatGPT browsing and training
  "OAI-SearchBot", // OpenAI — ChatGPT search results
  "ChatGPT-User", // OpenAI — user-initiated fetches
  "ClaudeBot", // Anthropic
  "Claude-Web",
  "anthropic-ai",
  "PerplexityBot", // Perplexity
  "Perplexity-User",
  "Google-Extended", // Google — Gemini / AI Overviews grounding
  "Applebot-Extended", // Apple Intelligence
  "cohere-ai",
  "Bytespider",
];

/** Paths behind authentication: no crawler has anything to gain from them. */
const PRIVATE_PATHS = [
  "/api/",
  "/dashboard",
  "/member",
  "/super-admin",
  "/login",
  "/onboarding",
  // Locale-prefixed variants of the same, since the app serves /en/... and /ne/...
  "/en/dashboard",
  "/en/member",
  "/en/super-admin",
  "/en/login",
  "/en/onboarding",
  "/ne/dashboard",
  "/ne/member",
  "/ne/super-admin",
  "/ne/login",
  "/ne/onboarding",
];

export default function robots(): MetadataRoute.Robots {
  const host = (headers().get("host") ?? BASE_DOMAIN).split(":")[0].toLowerCase();
  const proto = headers().get("x-forwarded-proto") ?? "https";
  const origin = `${proto}://${host}`;
  const isApex = host === BASE_DOMAIN || host === `www.${BASE_DOMAIN}`;

  return {
    rules: [
      { userAgent: "*", allow: "/", disallow: PRIVATE_PATHS },
      ...ASSISTANT_CRAWLERS.map((userAgent) => ({
        userAgent,
        allow: "/",
        disallow: PRIVATE_PATHS,
      })),
    ],
    // Each host advertises its own sitemap; the apex additionally advertises
    // the index, which is how a crawler arriving at nepalgym.xyz discovers
    // that gym subdomains exist at all.
    sitemap: isApex
      ? [`${origin}/sitemap.xml`, `${origin}/sitemap-index.xml`]
      : `${origin}/sitemap.xml`,
  };
}
