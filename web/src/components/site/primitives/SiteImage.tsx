import { ImageOff } from "lucide-react";

/**
 * Defensive image wrapper for tenant-supplied media.
 *
 * Site content URLs come from arbitrary uploads on arbitrary domains, and
 * next/image requires every remote host to be allow-listed ahead of time —
 * not something a renderer can do per-gym. A plain <img> avoids that, and an
 * empty `src` renders a muted placeholder instead of a broken-image icon.
 */
export function SiteImage({
  src,
  alt,
  className = "",
  ratio = "aspect-video",
}: {
  src?: string;
  alt: string;
  className?: string;
  ratio?: string;
}) {
  if (!src) {
    return (
      <div
        role="img"
        aria-label={alt || "Image not available"}
        className={`${ratio} w-full flex items-center justify-center bg-[var(--site-surface)] text-[var(--site-fg-muted)] rounded-[var(--site-radius)] ${className}`}
      >
        <ImageOff className="h-8 w-8 opacity-40" aria-hidden="true" />
      </div>
    );
  }

  return (
    // eslint-disable-next-line @next/next/no-img-element -- remote gym-uploaded URLs, arbitrary hosts
    <img
      src={src}
      alt={alt}
      loading="lazy"
      className={`${ratio} w-full object-cover rounded-[var(--site-radius)] ${className}`}
    />
  );
}
