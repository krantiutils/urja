import { Quote, Star } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, num, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * testimonials — cards | carousel | quote
 *
 * `carousel` is a CSS scroll-snap strip rather than a JS slider: it is the
 * only interactive-looking variant of this type, and per the client-boundary
 * list only gallery's carousel needs "use client" — this one stays server
 * rendered and relies on native horizontal scroll.
 */

function Stars({ rating }: { rating: number | null }) {
  if (!rating || rating <= 0) return null;
  const count = Math.max(1, Math.min(5, Math.round(rating)));
  return (
    <div className="flex gap-0.5" aria-label={`${count} out of 5 stars`}>
      {Array.from({ length: 5 }).map((_, i) => (
        <Star
          key={i}
          className={`h-4 w-4 ${i < count ? "fill-[var(--site-accent)] text-[var(--site-accent)]" : "text-[var(--site-border)]"}`}
          aria-hidden="true"
        />
      ))}
    </div>
  );
}

export default function TestimonialsRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const items = list(content, "items").filter((item) => itemText(item, "text", locale) !== "");

  if (items.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  if (variant === "quote") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {items.length > 0 ? (
          <div className="flex flex-col gap-12">
            {items.map((item, i) => (
              <figure key={i} className="mx-auto max-w-2xl text-center">
                <Quote className="mx-auto h-8 w-8 text-[var(--site-accent)]" aria-hidden="true" />
                <blockquote className="mt-4 text-xl sm:text-2xl leading-relaxed" style={{ fontFamily: "var(--site-font-display)" }}>
                  &ldquo;{itemText(item, "text", locale)}&rdquo;
                </blockquote>
                <figcaption className="mt-4 flex flex-col items-center gap-2">
                  <Stars rating={num(item, "rating")} />
                  {itemText(item, "name", locale) ? (
                    <span className="text-sm text-[var(--site-fg-muted)]">— {itemText(item, "name", locale)}</span>
                  ) : null}
                </figcaption>
              </figure>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  if (variant === "carousel") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {items.length > 0 ? (
          <div className="flex gap-4 overflow-x-auto snap-x snap-mandatory pb-4 -mx-4 px-4 sm:mx-0 sm:px-0">
            {items.map((item, i) => (
              <figure
                key={i}
                className="shrink-0 snap-start w-[85%] sm:w-[48%] lg:w-[32%] rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-6"
              >
                <Stars rating={num(item, "rating")} />
                <blockquote className="mt-3 text-[var(--site-fg-muted)] leading-relaxed">
                  &ldquo;{itemText(item, "text", locale)}&rdquo;
                </blockquote>
                {itemText(item, "name", locale) ? (
                  <figcaption className="mt-3 text-sm font-medium">{itemText(item, "name", locale)}</figcaption>
                ) : null}
              </figure>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  // cards
  return (
    <>
      <SectionHeading title={title} align={align} />
      {items.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item, i) => (
            <figure
              key={i}
              className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-6"
            >
              <Stars rating={num(item, "rating")} />
              <blockquote className="mt-3 text-[var(--site-fg-muted)] leading-relaxed">
                &ldquo;{itemText(item, "text", locale)}&rdquo;
              </blockquote>
              {itemText(item, "name", locale) ? (
                <figcaption className="mt-3 text-sm font-medium">{itemText(item, "name", locale)}</figcaption>
              ) : null}
            </figure>
          ))}
        </div>
      ) : null}
    </>
  );
}
