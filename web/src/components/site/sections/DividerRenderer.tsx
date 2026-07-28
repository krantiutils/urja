import type { Section } from "@/types/site";
import type { Locale } from "@/types";

/**
 * divider — line | dots | space
 *
 * Carries no content — SECTION_SPECS.divider.defaultContent is `{}` — so
 * nothing here reads `section.content` at all.
 */
export default function DividerRenderer({ section }: { section: Section; locale: Locale }) {
  if (section.variant === "dots") {
    return (
      <div className="flex items-center justify-center gap-2" role="separator" aria-hidden="true">
        <span className="h-1.5 w-1.5 rounded-full bg-[var(--site-accent)]" />
        <span className="h-1.5 w-1.5 rounded-full bg-[var(--site-accent)]" />
        <span className="h-1.5 w-1.5 rounded-full bg-[var(--site-accent)]" />
      </div>
    );
  }

  if (section.variant === "space") {
    return <div className="h-8 sm:h-12" role="presentation" />;
  }

  // line
  return <hr className="border-t border-[var(--site-border)]" />;
}
