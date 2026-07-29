import type { Locale } from "@/types";

/**
 * Schema.org structured data for a gym's public site.
 *
 * Two audiences read this and neither reads the prose. Google uses it for the
 * local-business panel — name, hours, phone, map pin. Assistants use it as the
 * factual spine of an answer: asked "where can I learn boxing in Kirtipur", a
 * model that has crawled the page will quote whatever it can extract with
 * confidence, and a JSON-LD block is far more quotable than a sentence inside a
 * hero image caption.
 *
 * Everything here is derived from what the gym already entered in the page
 * builder. Nothing is invented — a fabricated opening time or a guessed
 * coordinate is worse than an absent one, because both surfaces present it as
 * fact.
 */

interface Section {
  type: string;
  content?: Record<string, unknown>;
}

interface SiteLike {
  slug: string;
  org_name: string;
  org_name_ne?: string;
  /** Loose on purpose: SiteSocials is a fixed-key interface with no index
   *  signature, and all this needs is "every value that looks like a URL". */
  socials?: object | null;
}

function str(v: unknown): string | undefined {
  return typeof v === "string" && v.trim() ? v.trim() : undefined;
}

function findSection(sections: Section[], type: string): Section | undefined {
  return sections.find((s) => s.type === type);
}

/**
 * Splits "Kirtipur 44618, Kathmandu" into the parts schema.org expects.
 * Deliberately conservative: if the shape is unfamiliar the whole string goes
 * in streetAddress rather than being guessed into the wrong fields.
 */
function postalAddress(address: string | undefined) {
  if (!address) return undefined;

  const parts = address.split(",").map((p) => p.trim()).filter(Boolean);
  const base = { "@type": "PostalAddress", addressCountry: "NP" };

  if (parts.length >= 2) {
    const street = parts[0];
    const postal = street.match(/\b(\d{5})\b/)?.[1];
    return {
      ...base,
      streetAddress: postal ? street.replace(postal, "").trim() : street,
      addressLocality: parts[parts.length - 1],
      ...(postal ? { postalCode: postal } : {}),
    };
  }
  return { ...base, streetAddress: address };
}

export interface GymSchemaInput {
  site: SiteLike;
  sections: Section[];
  locale: Locale;
  /** Absolute origin of this gym's site, e.g. https://ibckirtipur.nepalgym.xyz */
  origin: string;
}

/**
 * Builds a SportsActivityLocation node, or null when there is too little to
 * say. An entry with only a name is not worth emitting — it gives neither
 * surface anything it could not already read off the page.
 */
export function buildGymSchema({
  site,
  sections,
  locale,
  origin,
}: GymSchemaInput): Record<string, unknown> | null {
  const contact = findSection(sections, "contact_info")?.content ?? {};
  const hero = findSection(sections, "hero")?.content ?? {};

  const name =
    (locale === "ne" ? site.org_name_ne : undefined) || site.org_name;

  const phone = str(contact.phone);
  const address = str(locale === "ne" ? contact.addressNe : contact.address) ?? str(contact.address);
  const hours = str(contact.hoursNote) ?? str(contact.hoursNoteNe);
  const description =
    str(locale === "ne" ? hero.subtitleNe : hero.subtitle) ?? str(hero.subtitle);

  // A bare name tells neither Google nor an assistant anything useful.
  if (!phone && !address) return null;

  const sameAs = Object.values(site.socials ?? {}).filter(
    (v): v is string => typeof v === "string" && v.startsWith("http")
  );

  return {
    "@context": "https://schema.org",
    "@type": "SportsActivityLocation",
    "@id": `${origin}#gym`,
    name,
    url: origin,
    ...(description ? { description } : {}),
    ...(phone ? { telephone: phone } : {}),
    ...(address ? { address: postalAddress(address) } : {}),
    // Free-form rather than openingHoursSpecification: the gym writes hours as
    // a sentence, and parsing "mornings and 5:00 - 7:00 PM" into structured
    // day/time pairs would mean guessing which days the mornings cover.
    ...(hours ? { openingHours: hours } : {}),
    ...(sameAs.length ? { sameAs } : {}),
    ...(locale === "ne" ? { inLanguage: "ne-NP" } : { inLanguage: "en" }),
  };
}

/** Breadcrumbs help Google render the path under a result and cost nothing. */
export function buildBreadcrumbs(
  origin: string,
  pageSlug: string,
  pageTitle: string
): Record<string, unknown> | null {
  if (pageSlug === "home") return null;
  return {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      { "@type": "ListItem", position: 1, name: "Home", item: origin },
      { "@type": "ListItem", position: 2, name: pageTitle, item: `${origin}/${pageSlug}` },
    ],
  };
}
