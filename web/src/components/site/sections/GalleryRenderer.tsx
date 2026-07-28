import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * gallery — grid | masonry | carousel
 *
 * `carousel` is a CSS scroll-snap strip, not a JS slider, so this module stays
 * server rendered. A gym that has not uploaded photos yet sees its
 * `emptyHint` rather than a row of placeholder boxes.
 */
export default function GalleryRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);

  // Tolerates both shapes: a list of objects, and a bare list of URLs from
  // hand-edited JSON.
  const raw = content?.images;
  const images = Array.isArray(raw)
    ? raw
        .map((entry) =>
          typeof entry === "string"
            ? { src: entry, alt: "" }
            : {
                src: str(entry as Record<string, unknown>, "src") ||
                  str(entry as Record<string, unknown>, "url"),
                alt: itemText(entry as Record<string, unknown>, "alt", locale),
              }
        )
        .filter((img) => img.src !== "")
    : [];

  const align = section.style?.align ?? "left";

  if (images.length === 0) {
    const hint = text(content, "emptyHint", locale);
    if (!title && !hint) return null;
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {hint ? (
          <p className="text-sm text-[var(--site-fg-muted)] leading-relaxed">{hint}</p>
        ) : null}
      </>
    );
  }

  if (variant === "masonry") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {/* CSS columns give a masonry flow without measuring anything in JS. */}
        <div className="columns-2 lg:columns-3 gap-4 [column-fill:_balance]">
          {images.map((img, i) => (
            <div key={i} className="mb-4 break-inside-avoid">
              <SiteImage src={img.src} alt={img.alt} ratio="aspect-auto" />
            </div>
          ))}
        </div>
      </>
    );
  }

  if (variant === "carousel") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        <div className="flex gap-4 overflow-x-auto snap-x snap-mandatory pb-4 -mx-4 px-4 sm:mx-0 sm:px-0">
          {images.map((img, i) => (
            <div key={i} className="shrink-0 snap-start w-[85%] sm:w-[48%] lg:w-[32%]">
              <SiteImage src={img.src} alt={img.alt} />
            </div>
          ))}
        </div>
      </>
    );
  }

  // grid
  return (
    <>
      <SectionHeading title={title} subtitle={subtitle} align={align} />
      <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
        {images.map((img, i) => (
          <SiteImage key={i} src={img.src} alt={img.alt} ratio="aspect-square" />
        ))}
      </div>
    </>
  );
}
