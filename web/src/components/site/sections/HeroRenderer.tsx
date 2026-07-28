import type { CSSProperties } from "react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, str, list, itemText, buttons } from "@/lib/site/content";
import { SiteButton } from "@/components/site/primitives/SiteButton";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * hero — centered | split | fullbleed | minimal | reel
 *
 * The page's primary heading, so the title renders as an <h1> rather than
 * going through SectionHeading (which is for secondary, in-page headings).
 */

import { HeroReel } from "@/components/site/primitives/HeroReel";

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
  const video = str(content, "video");
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

  if (variant === "reel") {
    const eyebrow = text(content, "eyebrow", locale);
    const facts = list(content, "facts");

    return (
      <div className="relative isolate min-h-[88svh] flex flex-col justify-end overflow-hidden">
        {video ? (
          <HeroReel
            src={video}
            poster={image}
            className="absolute inset-0 -z-20 h-full w-full object-cover"
          />
        ) : image ? (
          // eslint-disable-next-line @next/next/no-img-element -- tenant-uploaded, arbitrary host
          <img src={image} alt="" className="absolute inset-0 -z-20 h-full w-full object-cover" />
        ) : null}

        {/* Two layers rather than one flat scrim: a bottom-weighted gradient so
            the copy sits on near-solid colour, and a light overall wash so the
            footage never fights the text at the top of the frame. */}
        <div
          className="absolute inset-0 -z-10"
          style={{
            background:
              "linear-gradient(to top, var(--site-bg) 4%, color-mix(in srgb, var(--site-bg) 82%, transparent) 34%, color-mix(in srgb, var(--site-bg) 30%, transparent) 100%)",
          }}
          aria-hidden="true"
        />

        <div className="relative w-full max-w-6xl mx-auto px-5 sm:px-8 pb-12 sm:pb-20">
          {eyebrow ? (
            <p className="mb-5 inline-flex items-center gap-2.5 text-[11px] sm:text-xs font-mono uppercase tracking-[0.32em]">
              <span className="h-px w-8 bg-[var(--site-accent)]" aria-hidden="true" />
              {/* On its own the accent was unreadable against moving footage;
                  the rule carries the colour and the label stays legible. */}
              <span className="text-[var(--site-fg)]">{eyebrow}</span>
            </p>
          ) : null}

          {title ? (
            <h1
              className="max-w-4xl leading-[0.88] text-[clamp(2.75rem,11vw,7.5rem)]"
              style={displayStyle}
            >
              {title}
            </h1>
          ) : null}

          {subtitle ? (
            <p className="mt-6 max-w-xl text-base sm:text-lg text-[var(--site-fg-muted)] leading-relaxed">
              {subtitle}
            </p>
          ) : null}

          {btns.length > 0 && (
            <div className="mt-8 flex flex-wrap gap-3">
              {btns.map((b, i) => (
                <SiteButton key={i} {...b} />
              ))}
            </div>
          )}

          {facts.length > 0 ? (
            <dl
              className="mt-10 grid grid-cols-2 sm:grid-cols-4 gap-px overflow-hidden rounded-[var(--site-radius)] border border-[var(--site-border)]"
              style={{ background: "var(--site-border)" }}
            >
              {facts.map((f, i) => (
                <div
                  key={i}
                  className="px-4 py-3.5 backdrop-blur-md"
                  style={{ background: "color-mix(in srgb, var(--site-bg) 55%, transparent)" }}
                >
                  <dt className="text-[10px] font-mono uppercase tracking-[0.18em] text-[var(--site-fg-muted)]">
                    {itemText(f, "label", locale)}
                  </dt>
                  <dd className="mt-1 text-sm sm:text-base font-medium">
                    {itemText(f, "value", locale)}
                  </dd>
                </div>
              ))}
            </dl>
          ) : null}
        </div>
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
