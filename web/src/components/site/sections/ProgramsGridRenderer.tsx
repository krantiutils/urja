import { Award, Dumbbell, Flame, Shield, Target, Users, Zap, type LucideIcon } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * programs_grid — cards | icons | list | numbered
 */

const ICONS: Record<string, LucideIcon> = {
  target: Target,
  zap: Zap,
  users: Users,
  flame: Flame,
  shield: Shield,
  award: Award,
};

function ProgramIcon({ name, className }: { name: string; className?: string }) {
  const Icon = ICONS[name] ?? Dumbbell;
  return <Icon className={className} aria-hidden="true" />;
}

export default function ProgramsGridRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const items = list(content, "items");

  if (items.length === 0 && !title && !subtitle) return null;

  const align = section.style?.align ?? "left";

  if (variant === "icons") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {items.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8 text-center">
            {items.map((item, i) => (
              <div key={i} className="flex flex-col items-center gap-3">
                <div className="flex h-14 w-14 items-center justify-center rounded-full bg-[var(--site-accent)] text-[var(--site-accent-fg)]">
                  <ProgramIcon name={str(item, "icon")} className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-medium">{itemText(item, "title", locale)}</h3>
                {itemText(item, "description", locale) ? (
                  <p className="text-sm text-[var(--site-fg-muted)] leading-relaxed">
                    {itemText(item, "description", locale)}
                  </p>
                ) : null}
              </div>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  if (variant === "list") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {items.length > 0 ? (
          <div className="divide-y divide-[var(--site-border)] border-t border-b border-[var(--site-border)]">
            {items.map((item, i) => (
              <div key={i} className="flex items-start gap-4 py-5">
                <ProgramIcon name={str(item, "icon")} className="h-5 w-5 mt-1 shrink-0 text-[var(--site-accent)]" />
                <div>
                  <h3 className="text-lg font-medium">{itemText(item, "title", locale)}</h3>
                  {itemText(item, "description", locale) ? (
                    <p className="mt-1 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                      {itemText(item, "description", locale)}
                    </p>
                  ) : null}
                </div>
              </div>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  if (variant === "numbered") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {items.length > 0 ? (
          // Two columns from md up: a single column of short entries left most
          // of the page empty, which read as unfinished rather than spacious.
          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-10 border-t border-[var(--site-border)]">
            {items.map((item, i) => (
              <div
                key={i}
                className="flex items-baseline gap-5 border-b border-[var(--site-border)] py-6"
              >
                <span
                  className="shrink-0 text-4xl sm:text-5xl leading-none text-[var(--site-accent)] opacity-90"
                  style={{ fontFamily: "var(--site-font-display)" }}
                  aria-hidden="true"
                >
                  {String(i + 1).padStart(2, "0")}
                </span>
                <div>
                  <h3 className="text-xl sm:text-2xl" style={{ fontFamily: "var(--site-font-display)" }}>
                    {itemText(item, "title", locale)}
                  </h3>
                  {itemText(item, "description", locale) ? (
                    <p className="mt-1.5 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                      {itemText(item, "description", locale)}
                    </p>
                  ) : null}
                </div>
              </div>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  // cards
  return (
    <>
      <SectionHeading title={title} subtitle={subtitle} align={align} />
      {items.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item, i) => (
            <div
              key={i}
              className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-6"
            >
              <ProgramIcon name={str(item, "icon")} className="h-6 w-6 text-[var(--site-accent)]" />
              <h3 className="mt-4 text-lg font-medium">{itemText(item, "title", locale)}</h3>
              {itemText(item, "description", locale) ? (
                <p className="mt-2 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                  {itemText(item, "description", locale)}
                </p>
              ) : null}
            </div>
          ))}
        </div>
      ) : null}
    </>
  );
}
