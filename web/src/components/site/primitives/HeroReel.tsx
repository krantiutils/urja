"use client";

import { useEffect, useRef, useState } from "react";

/**
 * Background footage for the reel hero.
 *
 * Autoplaying video is a bandwidth tax on a phone in Kathmandu, so this is
 * deliberately conservative: the poster is what paints first and is all a
 * reduced-motion visitor or a save-data connection ever gets, the clip only
 * starts once the hero is actually on screen, and it pauses the moment it is
 * scrolled away or the tab is hidden.
 *
 * The video is decoration behind text. If any of it fails, the poster remains
 * and the hero reads exactly the same.
 */
export function HeroReel({
  src,
  poster,
  className = "",
}: {
  src: string;
  poster: string;
  className?: string;
}) {
  const ref = useRef<HTMLVideoElement>(null);
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    if (window.matchMedia?.("(prefers-reduced-motion: reduce)").matches) return;

    // Chromium exposes the browser/OS data saver here; honouring it costs
    // nothing and the poster already carries the design.
    const conn = (navigator as { connection?: { saveData?: boolean; effectiveType?: string } })
      .connection;
    if (conn?.saveData) return;
    if (conn?.effectiveType && /^(slow-)?2g$/.test(conn.effectiveType)) return;

    setEnabled(true);
  }, []);

  useEffect(() => {
    const video = ref.current;
    if (!video || !enabled) return;

    const play = () => void video.play().catch(() => {});
    const io = new IntersectionObserver(
      ([e]) => (e.isIntersecting && !document.hidden ? play() : video.pause()),
      { threshold: 0.1 }
    );
    io.observe(video);

    const onVisibility = () => (document.hidden ? video.pause() : play());
    document.addEventListener("visibilitychange", onVisibility);

    return () => {
      io.disconnect();
      document.removeEventListener("visibilitychange", onVisibility);
    };
  }, [enabled]);

  return (
    <video
      ref={ref}
      className={className}
      poster={poster}
      muted
      loop
      playsInline
      // preload="none" until we have decided to play: the poster is a 50 KB
      // image, the clip is 800 KB, and most visitors judge the page from the
      // first paint alone.
      preload={enabled ? "auto" : "none"}
      aria-hidden="true"
      tabIndex={-1}
    >
      {enabled ? <source src={src} type="video/mp4" /> : null}
    </video>
  );
}
