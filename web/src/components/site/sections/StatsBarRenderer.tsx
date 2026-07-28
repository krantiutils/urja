import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str } from "@/lib/site/content";

/**
 * stats_bar — inline | cards | bordered
 */
export default function StatsBarRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const items = list(section.content, "items");
  if (items.length === 0) return null;

  if (section.variant === "cards") {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {items.map((item, i) => (
          <div
            key={i}
            className="rounded-[var(--site-radius)] bg-[var(--site-surface)] px-4 py-6 text-center"
          >
            <div className="text-2xl sm:text-3xl font-semibold text-[var(--site-accent)]">
              {str(item, "value")}
            </div>
            <div className="mt-1 text-sm text-[var(--site-fg-muted)]">{itemText(item, "label", locale)}</div>
          </div>
        ))}
      </div>
    );
  }

  if (section.variant === "bordered") {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-4 divide-x divide-y sm:divide-y-0 divide-[var(--site-border)] border border-[var(--site-border)] rounded-[var(--site-radius)]">
        {items.map((item, i) => (
          <div key={i} className="px-4 py-6 text-center">
            <div className="text-2xl sm:text-3xl font-semibold">{str(item, "value")}</div>
            <div className="mt-1 text-sm text-[var(--site-fg-muted)]">{itemText(item, "label", locale)}</div>
          </div>
        ))}
      </div>
    );
  }

  // inline
  return (
    <div className="flex flex-wrap justify-center gap-x-10 gap-y-6">
      {items.map((item, i) => (
        <div key={i} className="flex flex-col items-center px-2">
          <div className="text-2xl sm:text-3xl font-semibold">{str(item, "value")}</div>
          <div className="mt-1 text-sm text-[var(--site-fg-muted)]">{itemText(item, "label", locale)}</div>
        </div>
      ))}
    </div>
  );
}
