import type { ReactNode } from "react";
import type { SectionStyle } from "@/types/site";

/**
 * Shared wrapper applying a section's background, padding, width and alignment.
 *
 * Every renderer wraps its content in this. No renderer reimplements these four
 * controls, so the inspector can change any section's layout without knowing
 * what kind of section it is.
 */

const BACKGROUNDS: Record<string, string> = {
  base: "bg-[var(--site-bg)] text-[var(--site-fg)]",
  surface: "bg-[var(--site-surface)] text-[var(--site-fg)]",
  accent: "bg-[var(--site-accent)] text-[var(--site-accent-fg)]",
  none: "",
};

const PADDINGS: Record<string, string> = {
  none: "py-0",
  sm: "py-6 sm:py-8",
  md: "py-10 sm:py-14",
  lg: "py-16 sm:py-24",
};

const WIDTHS: Record<string, string> = {
  full: "w-full px-4 sm:px-6",
  contained: "w-full max-w-6xl mx-auto px-4 sm:px-6",
  narrow: "w-full max-w-3xl mx-auto px-4 sm:px-6",
};

const ALIGNS: Record<string, string> = {
  left: "text-left",
  center: "text-center",
  right: "text-right",
};

export interface SectionShellProps {
  style?: SectionStyle;
  children: ReactNode;
  /** Set on the outer element so tests and the builder can target a section. */
  sectionId?: string;
  sectionType?: string;
}

export function SectionShell({
  style,
  children,
  sectionId,
  sectionType,
}: SectionShellProps) {
  // Defaults matter: section JSON has no migrations, so a row written before a
  // token existed arrives with it undefined.
  const background = BACKGROUNDS[style?.background ?? "base"] ?? BACKGROUNDS.base;
  const padding = PADDINGS[style?.padding ?? "md"] ?? PADDINGS.md;
  const width = WIDTHS[style?.width ?? "contained"] ?? WIDTHS.contained;
  const align = ALIGNS[style?.align ?? "left"] ?? ALIGNS.left;

  return (
    <section
      className={`${background} ${padding}`}
      data-section-id={sectionId}
      data-section-type={sectionType}
    >
      <div className={`${width} ${align}`}>{children}</div>
    </section>
  );
}

/**
 * Section heading rendered in the template's display font. Returns null for an
 * empty title so a section without one does not leave a gap.
 */
export function SectionHeading({
  title,
  subtitle,
  align = "left",
}: {
  title?: string;
  subtitle?: string;
  align?: string;
}) {
  if (!title && !subtitle) return null;

  return (
    <div className={`mb-8 sm:mb-12 ${align === "center" ? "mx-auto max-w-2xl" : ""}`}>
      {title ? (
        <h2
          className="text-2xl sm:text-3xl lg:text-4xl"
          style={{
            fontFamily: "var(--site-font-display)",
            textTransform: "var(--site-display-transform)" as never,
            fontWeight: "var(--site-display-weight)" as never,
            letterSpacing: "var(--site-display-tracking)",
          }}
        >
          {title}
        </h2>
      ) : null}
      {subtitle ? (
        <p className="mt-3 text-base sm:text-lg text-[var(--site-fg-muted)] leading-relaxed">
          {subtitle}
        </p>
      ) : null}
    </div>
  );
}
