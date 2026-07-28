import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, str, buttons } from "@/lib/site/content";
import { SiteButton } from "@/components/site/primitives/SiteButton";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * cta_banner — solid | gradient | image | split
 *
 * SectionShell already paints the section's own background (usually
 * `accent` for this type). `solid` leans on that as-is; the other three
 * variants layer their own treatment on top of it.
 */
export default function CtaBannerRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const image = str(content, "image");
  const btns = buttons(content, locale);

  if (!title && !subtitle && btns.length === 0) return null;

  const buttonRow = btns.length > 0 && (
    <div className="flex flex-wrap gap-4">
      {btns.map((b, i) => (
        <SiteButton key={i} {...b} />
      ))}
    </div>
  );

  if (variant === "split") {
    return (
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6">
        <div>
          {title ? (
            <h2 className="text-2xl sm:text-3xl" style={{ fontFamily: "var(--site-font-display)" }}>
              {title}
            </h2>
          ) : null}
          {subtitle ? <p className="mt-2 text-[var(--site-fg-muted)] max-w-xl">{subtitle}</p> : null}
        </div>
        {btns.length > 0 && (
          <div className="flex flex-wrap gap-4 sm:shrink-0">
            {btns.map((b, i) => (
              <SiteButton key={i} {...b} />
            ))}
          </div>
        )}
      </div>
    );
  }

  if (variant === "image" && image) {
    return (
      <div className="relative overflow-hidden rounded-[var(--site-radius)]">
        <div className="absolute inset-0">
          <SiteImage src={image} alt="" ratio="" className="h-full" />
          <div className="absolute inset-0 bg-[var(--site-bg)] opacity-60" />
        </div>
        <div className="relative z-10 flex flex-col items-center gap-4 text-center px-6 py-16">
          {title ? (
            <h2 className="text-3xl sm:text-4xl" style={{ fontFamily: "var(--site-font-display)" }}>
              {title}
            </h2>
          ) : null}
          {subtitle ? <p className="text-[var(--site-fg-muted)] max-w-xl">{subtitle}</p> : null}
          {buttonRow}
        </div>
      </div>
    );
  }

  if (variant === "gradient") {
    return (
      // Negative margins cancel SectionShell's horizontal padding so the
      // gradient reaches the edges. Rounded and inset, it painted a second
      // coloured card on top of a section already filled with the accent —
      // a box inside a box.
      <div className="-mx-4 sm:-mx-6 bg-gradient-to-br from-[var(--site-accent)] to-[var(--site-surface)] px-6 py-16">
        <div className="flex flex-col items-center gap-4 text-center text-[var(--site-accent-fg)]">
          {title ? (
            <h2 className="text-3xl sm:text-4xl" style={{ fontFamily: "var(--site-font-display)" }}>
              {title}
            </h2>
          ) : null}
          {subtitle ? <p className="max-w-xl opacity-90">{subtitle}</p> : null}
          {buttonRow}
        </div>
      </div>
    );
  }

  // solid (also the fallback for "image" with no image URL)
  return (
    <div className="flex flex-col items-center gap-4 text-center">
      {title ? (
        <h2 className="text-3xl sm:text-4xl" style={{ fontFamily: "var(--site-font-display)" }}>
          {title}
        </h2>
      ) : null}
      {subtitle ? <p className="text-[var(--site-fg-muted)] max-w-xl">{subtitle}</p> : null}
      {buttonRow}
    </div>
  );
}
