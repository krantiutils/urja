import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * opening_hours — table | list | compact
 *
 * Hours are read bilingually rather than as a plain string: "Closed" is a word,
 * not a time, and a gym that shuts on Saturday needs to say so in both
 * languages.
 */
export default function OpeningHoursRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const days = list(content, "days");
  const note = text(content, "note", locale);

  if (days.length === 0 && !title && !note) return null;

  const align = section.style?.align ?? "left";

  // Sessions within the opening window are the thing members actually plan
  // around, so a note has to survive to the page.
  const footnote = note ? (
    <p className="mt-4 text-sm text-[var(--site-fg-muted)] leading-relaxed">{note}</p>
  ) : null;

  if (variant === "table") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-[var(--site-border)]">
                <th scope="col" className="py-2 pr-4 font-medium">
                  {locale === "ne" ? "दिन" : "Day"}
                </th>
                <th scope="col" className="py-2 font-medium">
                  {locale === "ne" ? "समय" : "Hours"}
                </th>
              </tr>
            </thead>
            <tbody>
              {days.map((d, i) => (
                <tr key={i} className="border-b border-[var(--site-border)] last:border-b-0">
                  <th scope="row" className="py-3 pr-4 font-normal text-[var(--site-fg-muted)]">
                    {itemText(d, "day", locale)}
                  </th>
                  <td className="py-3">{itemText(d, "hours", locale)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {footnote}
      </>
    );
  }

  if (variant === "compact") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        <dl className="grid grid-cols-[auto_1fr] gap-x-6 gap-y-1 text-sm">
          {days.map((d, i) => (
            <div key={i} className="contents">
              <dt className="text-[var(--site-fg-muted)]">{itemText(d, "day", locale)}</dt>
              <dd>{itemText(d, "hours", locale)}</dd>
            </div>
          ))}
        </dl>
        {footnote}
      </>
    );
  }

  // list
  return (
    <>
      <SectionHeading title={title} align={align} />
      <ul className="flex flex-col divide-y divide-[var(--site-border)] border-t border-b border-[var(--site-border)]">
        {days.map((d, i) => (
          <li key={i} className="flex items-center justify-between py-3">
            <span className="font-medium">{itemText(d, "day", locale)}</span>
            <span className="text-[var(--site-fg-muted)]">{itemText(d, "hours", locale)}</span>
          </li>
        ))}
      </ul>
      {footnote}
    </>
  );
}
