import { ExternalLink, MapPin } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { num, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * map_embed — standard | with_info | full_width
 *
 * Coordinates come from the gym, so the embed URL is built here from two
 * numbers rather than accepting a tenant-supplied iframe src — the API already
 * rejects <iframe> in content, and this keeps the only framed origin a fixed
 * one. A gym that has given an address but no coordinates gets a directions
 * link instead of an empty frame.
 */

/** OpenStreetMap needs no API key, so a gym can publish a map with no setup. */
function embedUrl(lat: number, lon: number): string {
  const d = 0.006;
  const bbox = [lon - d, lat - d / 2, lon + d, lat + d / 2].join("%2C");
  return `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}&layer=mapnik&marker=${lat}%2C${lon}`;
}

function directionsUrl(lat: number | null, lon: number | null, address: string): string {
  if (lat !== null && lon !== null) {
    return `https://www.google.com/maps/search/?api=1&query=${lat},${lon}`;
  }
  return `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(address)}`;
}

export default function MapEmbedRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const address = text(content, "address", locale);
  const lat = num(content, "latitude");
  const lon = num(content, "longitude");
  const hasPin = lat !== null && lon !== null;

  if (!hasPin && !address && !title) return null;

  const align = section.style?.align ?? "left";
  const directionsLabel = locale === "ne" ? "दिशा हेर्नुहोस्" : "Get directions";

  const frame = hasPin ? (
    <iframe
      src={embedUrl(lat, lon)}
      title={title || address || "Map"}
      loading="lazy"
      referrerPolicy="no-referrer-when-downgrade"
      className={`w-full border-0 ${variant === "full_width" ? "h-[24rem] sm:h-[30rem]" : "aspect-[16/9] rounded-[var(--site-radius)]"}`}
    />
  ) : null;

  const details = (
    <div className="flex flex-col gap-3">
      {address ? (
        <p className="flex items-start gap-2 text-[var(--site-fg-muted)]">
          <MapPin className="mt-0.5 h-5 w-5 shrink-0 text-[var(--site-accent)]" aria-hidden="true" />
          <span>{address}</span>
        </p>
      ) : null}
      {address || hasPin ? (
        <a
          href={directionsUrl(lat, lon, address)}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1.5 text-sm font-medium text-[var(--site-accent)] hover:underline"
        >
          {directionsLabel}
          <ExternalLink className="h-4 w-4" aria-hidden="true" />
        </a>
      ) : null}
    </div>
  );

  if (variant === "full_width") {
    return (
      <>
        {frame}
        {!frame ? <div className="px-4 sm:px-6">{details}</div> : null}
      </>
    );
  }

  if (variant === "with_info") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
          <div className="lg:col-span-2">{frame ?? details}</div>
          {frame ? <div>{details}</div> : null}
        </div>
      </>
    );
  }

  // standard
  return (
    <>
      <SectionHeading title={title} align={align} />
      {frame}
      <div className={frame ? "mt-4" : ""}>{details}</div>
    </>
  );
}
