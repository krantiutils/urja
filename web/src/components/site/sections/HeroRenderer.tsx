import type { CSSProperties } from "react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, str, buttons } from "@/lib/site/content";
import { SiteButton } from "@/components/site/primitives/SiteButton";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * hero — centered | split | fullbleed | minimal
 *
 * The page's primary heading, so the title renders as an <h1> rather than
 * going through SectionHeading (which is for secondary, in-page headings).
 */

const displayStyle: CSSProperties = {
  fontFamily: "var(--site-font-display)",
  textTransform: "var(--site-display-transform)" as never,
  fontWeight: "var(--site-display-weight)" as never,
  letterSpacing: "var(--site-display-tracking)",
};

function alignClasses(align: string | undefined) {
  if (align === "center") return { items: "items-center", justify: "justify-center" };
  if (align === "right") return { items: "items-end", justify: "justify-end" };
  return { items: "items-start", justify: "justify-start" };
}

export default function HeroRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const image = str(content, "image");
  const btns = buttons(content, locale);

  if (!title && !subtitle && !image && btns.length === 0) return null;

  const buttonRow = btns.length > 0 && (
    <div className="flex flex-wrap gap-4">
      {btns.map((b, i) => (
        <SiteButton key={i} {...b} />
      ))}
    </div>
  );

  if (variant === "fullbleed") {
    return (
      <div className="relative w-full overflow-hidden rounded-[var(--site-radius)] min-h-[60vh] sm:min-h-[70vh] flex items-center justify-center">
        {image ? (
          <div className="absolute inset-0">
            <SiteImage src={image} alt="" ratio="" className="h-full min-h-[60vh] sm:min-h-[70vh]" />
            <div className="absolute inset-0 bg-[var(--site-bg)] opacity-60" />
          </div>
        ) : null}
        <div className="relative z-10 flex flex-col items-center gap-6 max-w-3xl text-center px-4 py-16">
          {title ? (
            <h1 className="text-4xl sm:text-5xl lg:text-6xl" style={displayStyle}>
              {title}
            </h1>
          ) : null}
          {subtitle ? (
            <p className="text-lg sm:text-xl text-[var(--site-fg-muted)] leading-relaxed">{subtitle}</p>
          ) : null}
          {btns.length > 0 && (
            <div className="flex flex-wrap gap-4 justify-center">
              {btns.map((b, i) => (
                <SiteButton key={i} {...b} />
              ))}
            </div>
          )}
        </div>
      </div>
    );
  }

  if (variant === "split") {
    const { items } = alignClasses(section.style?.align);
    return (
      <div className={`grid gap-10 lg:items-center ${image ? "lg:grid-cols-2" : ""}`}>
        <div className={`flex flex-col ${items} gap-6`}>
          {title ? (
            <h1 className="text-3xl sm:text-4xl lg:text-5xl" style={displayStyle}>
              {title}
            </h1>
          ) : null}
          {subtitle ? (
            <p className="text-base sm:text-lg text-[var(--site-fg-muted)] leading-relaxed">{subtitle}</p>
          ) : null}
          {buttonRow}
        </div>
        {image ? <SiteImage src={image} alt={title} ratio="aspect-[4/3]" /> : null}
      </div>
    );
  }

  if (variant === "minimal") {
    const { items } = alignClasses(section.style?.align);
    return (
      <div className={`flex flex-col ${items} gap-4`}>
        {title ? (
          <h1 className="text-2xl sm:text-3xl lg:text-4xl" style={displayStyle}>
            {title}
          </h1>
        ) : null}
        {subtitle ? (
          <p className="text-base sm:text-lg text-[var(--site-fg-muted)] leading-relaxed">{subtitle}</p>
        ) : null}
        {buttonRow}
      </div>
    );
  }

  // centered (default / fallback for an unrecognised variant string)
  return (
    <div className="flex flex-col items-center gap-6 max-w-3xl mx-auto text-center">
      {title ? (
        <h1 className="text-4xl sm:text-5xl lg:text-6xl" style={displayStyle}>
          {title}
        </h1>
      ) : null}
      {subtitle ? (
        <p className="text-lg sm:text-xl text-[var(--site-fg-muted)] leading-relaxed">{subtitle}</p>
      ) : null}
      {btns.length > 0 && (
        <div className="flex flex-wrap gap-4 justify-center">
          {btns.map((b, i) => (
            <SiteButton key={i} {...b} />
          ))}
        </div>
      )}
      {image ? <SiteImage src={image} alt={title} ratio="aspect-video" className="mt-4 max-w-4xl" /> : null}
    </div>
  );
}
