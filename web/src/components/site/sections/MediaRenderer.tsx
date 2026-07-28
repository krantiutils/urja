import { PlayCircle } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, str } from "@/lib/site/content";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * media — image | video
 */

function youTubeEmbed(url: string): string | null {
  const patterns = [/youtu\.be\/([\w-]{11})/, /[?&]v=([\w-]{11})/, /youtube\.com\/embed\/([\w-]{11})/];
  for (const p of patterns) {
    const m = url.match(p);
    if (m) return `https://www.youtube.com/embed/${m[1]}`;
  }
  return null;
}

function vimeoEmbed(url: string): string | null {
  const m = url.match(/vimeo\.com\/(\d+)/);
  return m ? `https://player.vimeo.com/video/${m[1]}` : null;
}

export default function MediaRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const url = str(content, "url");
  const poster = str(content, "poster");
  const caption = text(content, "caption", locale);
  const alt = text(content, "alt", locale) || caption;

  if (!url && !caption) return null;

  if (variant === "video") {
    const embed = url ? youTubeEmbed(url) ?? vimeoEmbed(url) : null;
    return (
      <figure>
        {!url ? (
          <div className="aspect-video w-full flex items-center justify-center rounded-[var(--site-radius)] bg-[var(--site-surface)] text-[var(--site-fg-muted)]">
            <PlayCircle className="h-10 w-10 opacity-40" aria-hidden="true" />
          </div>
        ) : embed ? (
          <div className="aspect-video w-full overflow-hidden rounded-[var(--site-radius)]">
            <iframe
              src={embed}
              title={caption || alt || "Video"}
              className="h-full w-full"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            />
          </div>
        ) : (
          // eslint-disable-next-line jsx-a11y/media-has-caption -- captions aren't part of the source content model
          // preload="metadata" so a phone on mobile data does not pull the whole
          // clip before anyone presses play; the poster carries the first
          // impression. max-h keeps a 9:16 reel from swallowing a desktop page.
          <video
            controls
            playsInline
            preload="metadata"
            poster={poster || undefined}
            className="mx-auto max-h-[70vh] w-auto max-w-full rounded-[var(--site-radius)]"
            aria-label={caption || alt || "Video"}
          >
            <source src={url} />
          </video>
        )}
        {caption ? (
          <figcaption className="mt-2 text-sm text-[var(--site-fg-muted)] text-center">{caption}</figcaption>
        ) : null}
      </figure>
    );
  }

  // image
  return (
    <figure>
      <SiteImage src={url} alt={alt || "Gym photo"} ratio="aspect-video" />
      {caption ? (
        <figcaption className="mt-2 text-sm text-[var(--site-fg-muted)] text-center">{caption}</figcaption>
      ) : null}
    </figure>
  );
}
