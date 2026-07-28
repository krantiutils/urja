"use client";

import { useState } from "react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * class_timetable — table | day_tabs | cards
 *
 * `day_tabs` needs client state to switch days, and the dispatcher has one
 * module per section type, so the whole file is a client component.
 *
 * A gym's timetable is the single most-read thing on its site, so every
 * variant renders the full week in the DOM: `day_tabs` hides inactive days
 * with CSS rather than unmounting them, which keeps them findable by in-page
 * search and readable when JS has not hydrated yet.
 */

interface Day {
  label: string;
  classes: Record<string, unknown>[];
}

function readDays(content: Record<string, unknown>, locale: Locale): Day[] {
  return list(content, "days")
    .map((day) => ({
      label: itemText(day, "day", locale),
      classes: list(day, "classes"),
    }))
    .filter((day) => day.label !== "" || day.classes.length > 0);
}

function ClassRow({ cls, locale }: { cls: Record<string, unknown>; locale: Locale }) {
  const time = str(cls, "time");
  const name = itemText(cls, "name", locale);
  const coach = itemText(cls, "coach", locale);
  const level = itemText(cls, "level", locale);

  return (
    <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1 py-2.5">
      <span className="w-14 shrink-0 tabular-nums text-sm text-[var(--site-fg-muted)]">{time}</span>
      <span className="font-medium">{name}</span>
      {coach ? <span className="text-sm text-[var(--site-fg-muted)]">· {coach}</span> : null}
      {level ? (
        <span className="ml-auto rounded-full border border-[var(--site-border)] px-2.5 py-0.5 text-xs text-[var(--site-fg-muted)]">
          {level}
        </span>
      ) : null}
    </div>
  );
}

export default function ClassTimetableRenderer({
  section,
  locale,
}: {
  section: Section;
  locale: Locale;
}) {
  const { content, variant } = section;
  const [activeDay, setActiveDay] = useState(0);

  const title = text(content, "title", locale);
  const note = text(content, "note", locale);
  const days = readDays(content, locale);

  if (days.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  const heading = <SectionHeading title={title} align={align} />;
  const footnote = note ? (
    <p className="mt-6 text-sm text-[var(--site-fg-muted)]">{note}</p>
  ) : null;

  if (variant === "day_tabs") {
    return (
      <>
        {heading}
        <div className="flex gap-2 overflow-x-auto pb-3 -mx-4 px-4 sm:mx-0 sm:px-0" role="tablist">
          {days.map((day, i) => (
            <button
              key={i}
              type="button"
              role="tab"
              aria-selected={i === activeDay}
              onClick={() => setActiveDay(i)}
              className={`shrink-0 rounded-[var(--site-radius)] border px-4 py-2 text-sm transition-colors ${
                i === activeDay
                  ? "border-[var(--site-accent)] bg-[var(--site-accent)] text-[var(--site-accent-fg)]"
                  : "border-[var(--site-border)] text-[var(--site-fg-muted)] hover:border-current"
              }`}
            >
              {day.label}
            </button>
          ))}
        </div>

        <div className="mt-4">
          {days.map((day, i) => (
            <div
              key={i}
              role="tabpanel"
              aria-label={day.label}
              // Hidden, not unmounted: the whole week stays in the DOM so it is
              // searchable and survives a failed hydration.
              className={i === activeDay ? "block" : "hidden"}
            >
              {day.classes.length === 0 ? (
                <p className="py-3 text-sm text-[var(--site-fg-muted)]">—</p>
              ) : (
                <div className="divide-y divide-[var(--site-border)]">
                  {day.classes.map((cls, j) => (
                    <ClassRow key={j} cls={cls} locale={locale} />
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
        {footnote}
      </>
    );
  }

  if (variant === "cards") {
    return (
      <>
        {heading}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {days.map((day, i) => (
            <div
              key={i}
              className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-5"
            >
              <h3 className="text-sm font-semibold uppercase tracking-wide text-[var(--site-accent)]">
                {day.label}
              </h3>
              {day.classes.length === 0 ? (
                <p className="mt-3 text-sm text-[var(--site-fg-muted)]">Rest day</p>
              ) : (
                <div className="mt-2 divide-y divide-[var(--site-border)]">
                  {day.classes.map((cls, j) => (
                    <ClassRow key={j} cls={cls} locale={locale} />
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
        {footnote}
      </>
    );
  }

  // table
  return (
    <>
      {heading}
      <div className="overflow-x-auto">
        <table className="w-full min-w-[34rem] border-collapse text-left">
          <tbody>
            {days.map((day, i) => (
              <tr key={i} className="border-b border-[var(--site-border)] align-top">
                <th
                  scope="row"
                  className="w-32 py-4 pr-4 text-sm font-semibold uppercase tracking-wide text-[var(--site-accent)]"
                >
                  {day.label}
                </th>
                <td className="py-2">
                  {day.classes.length === 0 ? (
                    <span className="text-sm text-[var(--site-fg-muted)]">—</span>
                  ) : (
                    <div className="divide-y divide-[var(--site-border)]">
                      {day.classes.map((cls, j) => (
                        <ClassRow key={j} cls={cls} locale={locale} />
                      ))}
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {footnote}
    </>
  );
}
