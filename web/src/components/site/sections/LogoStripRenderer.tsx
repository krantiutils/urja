"use client";

import { useState } from "react";
import { Pause, Play } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * logo_strip — grid | marquee | simple
 *
 * Only `marquee` is genuinely interactive (a pause control, for users who
 * find continuously-moving content distracting per WCAG 2.2.2), but the
 * dispatcher only has one file per section type, so the whole module carries
 * "use client" per the client-boundary list in the brief.
 */
export default function LogoStripRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const logos = list(content, "logos").filter((logo) => str(logo, "url") !== "");
  const [paused, setPaused] = useState(false);

  if (logos.length === 0 && !title) return null;

  const align = section.style?.align ?? "center";

  if (variant === "marquee") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {logos.length > 0 ? (
          <div className="relative">
            <style>{`
              @keyframes site-logo-marquee {
                from { transform: translateX(0); }
                to { transform: translateX(-50%); }
              }
            `}</style>
            <div className="overflow-hidden">
              <div
                className="flex w-max items-center gap-12 animate-[site-logo-marquee_28s_linear_infinite]"
                style={{ animationPlayState: paused ? "paused" : "running" }}
              >
                {[...logos, ...logos].map((logo, i) => (
                  <div key={i} className="h-10 w-28 shrink-0 opacity-70">
                    <SiteImage
                      src={str(logo, "url")}
                      alt={itemText(logo, "name", locale)}
                      ratio=""
                      className="h-full object-contain"
                    />
                  </div>
                ))}
              </div>
            </div>
            <button
              type="button"
              onClick={() => setPaused((p) => !p)}
              aria-pressed={paused}
              className="mt-4 inline-flex items-center gap-2 mx-auto text-xs text-[var(--site-fg-muted)] hover:text-[var(--site-fg)]"
            >
              {paused ? <Play className="h-3.5 w-3.5" aria-hidden="true" /> : <Pause className="h-3.5 w-3.5" aria-hidden="true" />}
              {paused ? (locale === "ne" ? "जारी राख्नुहोस्" : "Resume") : locale === "ne" ? "रोक्नुहोस्" : "Pause"}
            </button>
          </div>
        ) : null}
      </>
    );
  }

  if (variant === "simple") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {logos.length > 0 ? (
          <div className="flex flex-wrap items-center justify-center gap-8">
            {logos.map((logo, i) => (
              <div key={i} className="h-8 w-24 opacity-70 hover:opacity-100 transition-opacity">
                <SiteImage
                  src={str(logo, "url")}
                  alt={itemText(logo, "name", locale)}
                  ratio=""
                  className="h-full object-contain"
                />
              </div>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  // grid
  return (
    <>
      <SectionHeading title={title} align={align} />
      {logos.length > 0 ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-8 items-center justify-items-center">
          {logos.map((logo, i) => (
            <div key={i} className="h-12 w-32 grayscale hover:grayscale-0 transition-[filter] opacity-80 hover:opacity-100">
              <SiteImage
                src={str(logo, "url")}
                alt={itemText(logo, "name", locale)}
                ratio=""
                className="h-full object-contain"
              />
            </div>
          ))}
        </div>
      ) : null}
    </>
  );
}
