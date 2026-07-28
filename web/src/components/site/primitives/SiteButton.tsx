import type { ButtonSpec } from "@/lib/site/content";

/**
 * Renders a single CTA button in one of the three spec-defined styles.
 *
 * `outline` uses `border-current` / `text-current` rather than a theme
 * variable so it always matches whatever text colour SectionShell already put
 * on the ambient section (base fg, or accent-fg on an accent background)
 * without needing to know which background it is sitting on.
 */
export function SiteButton({ label, href, style }: ButtonSpec) {
  if (!label) return null;

  const base =
    "inline-flex items-center justify-center px-6 py-3 text-sm sm:text-base font-medium transition-opacity hover:opacity-90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-current";

  const variants: Record<ButtonSpec["style"], string> = {
    solid: "bg-[var(--site-accent)] text-[var(--site-accent-fg)] rounded-[var(--site-radius)]",
    outline: "bg-transparent border border-current rounded-[var(--site-radius)]",
    pill: "bg-[var(--site-accent)] text-[var(--site-accent-fg)] rounded-full",
  };

  return (
    <a href={href} className={`${base} ${variants[style]}`}>
      {label}
    </a>
  );
}
