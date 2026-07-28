import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * faq — accordion | list | two_column
 *
 * `accordion` uses native <details>/<summary> rather than JS state: it is
 * keyboard accessible and open-by-search out of the box, and it keeps this
 * module server-rendered.
 */
export default function FaqRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);

  // An entry with no question is a half-filled builder row, not content.
  const items = list(content, "items").filter(
    (item) => itemText(item, "question", locale) !== ""
  );

  if (items.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  if (variant === "list") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        <dl className="flex flex-col gap-8">
          {items.map((item, i) => (
            <div key={i}>
              <dt className="text-lg font-medium">{itemText(item, "question", locale)}</dt>
              {itemText(item, "answer", locale) ? (
                <dd className="mt-2 text-[var(--site-fg-muted)] leading-relaxed">
                  {itemText(item, "answer", locale)}
                </dd>
              ) : null}
            </div>
          ))}
        </dl>
      </>
    );
  }

  if (variant === "two_column") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-x-10 gap-y-8">
          {items.map((item, i) => (
            <div key={i}>
              <dt className="text-lg font-medium">{itemText(item, "question", locale)}</dt>
              {itemText(item, "answer", locale) ? (
                <dd className="mt-2 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                  {itemText(item, "answer", locale)}
                </dd>
              ) : null}
            </div>
          ))}
        </dl>
      </>
    );
  }

  // accordion
  return (
    <>
      <SectionHeading title={title} subtitle={subtitle} align={align} />
      <div className="divide-y divide-[var(--site-border)] border-t border-b border-[var(--site-border)]">
        {items.map((item, i) => (
          <details key={i} className="group py-4">
            <summary className="flex cursor-pointer list-none items-center justify-between gap-4 text-left text-lg font-medium marker:hidden">
              {itemText(item, "question", locale)}
              <span
                className="shrink-0 text-[var(--site-accent)] transition-transform group-open:rotate-45"
                aria-hidden="true"
              >
                +
              </span>
            </summary>
            {itemText(item, "answer", locale) ? (
              <p className="mt-3 text-[var(--site-fg-muted)] leading-relaxed">
                {itemText(item, "answer", locale)}
              </p>
            ) : null}
          </details>
        ))}
      </div>
    </>
  );
}
