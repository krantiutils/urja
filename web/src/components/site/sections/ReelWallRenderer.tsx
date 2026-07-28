"use client";

import { useState } from "react";
import { Play } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * reel_wall — grid | strip
 *
 * A wall of vertical clips, which is the shape gym footage actually comes in.
 * A single 9:16 reel dropped into a page section renders as a narrow box
 * marooned in empty space; several of them in a 9:16 grid read as intentional,
 * and match what a visitor already recognises from the gym's Instagram.
 *
 * Nothing autoplays and nothing preloads. Each tile is a poster image until it
 * is clicked, so a page with six clips costs six thumbnails rather than sixty
 * megabytes — the clips here run 3–10 MB each.
 */
export default function ReelWallRenderer({
  section,
  locale,
}: {
  section: Section;
  locale: Locale;
}) {
  const { content, variant } = section;
  const [playing, setPlaying] = useState<number | null>(null);

  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const items = list(content, "items").filter((i) => str(i, "url") !== "");

  if (items.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  const grid =
    variant === "strip"
      ? "flex gap-4 overflow-x-auto snap-x snap-mandatory pb-4 -mx-4 px-4 sm:mx-0 sm:px-0"
      : "grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4";

  const tile =
    variant === "strip"
      ? "shrink-0 snap-start w-[62%] sm:w-[38%] lg:w-[23%]"
      : "";

  return (
    <>
      <SectionHeading title={title} subtitle={subtitle} align={align} />

      <div className={grid}>
        {items.map((item, i) => {
          const url = str(item, "url");
          const poster = str(item, "poster");
          const caption = itemText(item, "caption", locale);
          const isPlaying = playing === i;

          return (
            <figure key={i} className={tile}>
              <div className="relative aspect-[9/16] overflow-hidden rounded-[var(--site-radius)] bg-[var(--site-surface)]">
                {isPlaying ? (
                  // eslint-disable-next-line jsx-a11y/media-has-caption -- captions are not part of the content model
                  <video
                    src={url}
                    poster={poster || undefined}
                    controls
                    autoPlay
                    playsInline
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <button
                    type="button"
                    onClick={() => setPlaying(i)}
                    className="group absolute inset-0 h-full w-full"
                    aria-label={caption ? `Play: ${caption}` : "Play clip"}
                  >
                    {poster ? (
                      // eslint-disable-next-line @next/next/no-img-element -- tenant-uploaded, arbitrary host
                      <img
                        src={poster}
                        alt=""
                        loading="lazy"
                        className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                      />
                    ) : null}
                    <span
                      className="absolute inset-0"
                      style={{
                        background:
                          "linear-gradient(to top, color-mix(in srgb, var(--site-bg) 75%, transparent) 0%, transparent 55%)",
                      }}
                    />
                    <span className="absolute inset-0 flex items-center justify-center">
                      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--site-accent)] text-[var(--site-accent-fg)] shadow-lg transition-transform duration-300 group-hover:scale-110">
                        <Play className="h-5 w-5 translate-x-[1px]" fill="currentColor" aria-hidden="true" />
                      </span>
                    </span>
                  </button>
                )}
              </div>

              {caption ? (
                <figcaption className="mt-2 text-xs sm:text-sm text-[var(--site-fg-muted)] leading-snug">
                  {caption}
                </figcaption>
              ) : null}
            </figure>
          );
        })}
      </div>
    </>
  );
}
