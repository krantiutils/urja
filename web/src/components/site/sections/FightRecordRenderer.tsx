import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, num, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * fight_record — timeline | table | cards
 *
 * Mirrors the shape of internal/boxing bouts so a gym can publish the same
 * record it tracks in the dashboard, but reads defensively: this is section
 * JSON, not a typed API response.
 */

const RESULT_LABELS: Record<string, { short: string; classes: string }> = {
  win: { short: "W", classes: "bg-[var(--site-accent)] text-[var(--site-accent-fg)]" },
  loss: { short: "L", classes: "border border-[var(--site-border)] text-[var(--site-fg-muted)]" },
  draw: { short: "D", classes: "border border-[var(--site-border)] text-[var(--site-fg-muted)]" },
  no_contest: { short: "NC", classes: "border border-[var(--site-border)] text-[var(--site-fg-muted)]" },
};

function ResultBadge({ result }: { result: string }) {
  const spec = RESULT_LABELS[result];
  if (!spec) return null;
  return (
    <span
      className={`inline-flex h-7 min-w-7 items-center justify-center rounded-full px-2 text-xs font-bold ${spec.classes}`}
      aria-label={result.replace("_", " ")}
    >
      {spec.short}
    </span>
  );
}

/** Tally shown above the list, mirroring boxing.TallyBouts. */
function tally(bouts: Record<string, unknown>[]) {
  const counts = { win: 0, loss: 0, draw: 0, no_contest: 0 };
  for (const bout of bouts) {
    const result = str(bout, "result");
    if (result in counts) counts[result as keyof typeof counts] += 1;
  }
  return counts;
}

function boutMeta(bout: Record<string, unknown>, locale: Locale): string {
  const rounds = num(bout, "rounds");
  return [
    itemText(bout, "method", locale),
    rounds ? `R${rounds}` : "",
    itemText(bout, "weight_class", locale) || itemText(bout, "weightClass", locale),
  ]
    .filter(Boolean)
    .join(" · ");
}

export default function FightRecordRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const bouts = list(content, "bouts");
  const align = section.style?.align ?? "left";

  if (bouts.length === 0) {
    const hint = text(content, "emptyHint", locale);
    if (!title && !hint) return null;
    return (
      <>
        <SectionHeading title={title} align={align} />
        {hint ? (
          <p className="text-sm text-[var(--site-fg-muted)] leading-relaxed">{hint}</p>
        ) : null}
      </>
    );
  }

  const counts = tally(bouts);
  const summary = (
    <p className="mb-8 text-3xl sm:text-4xl tabular-nums" style={{ fontFamily: "var(--site-font-display)" }}>
      {counts.win}
      <span className="text-[var(--site-fg-muted)]">–</span>
      {counts.loss}
      <span className="text-[var(--site-fg-muted)]">–</span>
      {counts.draw}
      {counts.no_contest > 0 ? (
        <span className="ml-2 text-base text-[var(--site-fg-muted)]">({counts.no_contest} NC)</span>
      ) : null}
    </p>
  );

  if (variant === "table") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {summary}
        <div className="overflow-x-auto">
          <table className="w-full min-w-[32rem] border-collapse text-left text-sm">
            <thead>
              <tr className="border-b border-[var(--site-border)] text-[var(--site-fg-muted)]">
                <th scope="col" className="py-2 pr-4 font-medium">Date</th>
                <th scope="col" className="py-2 pr-4 font-medium">Opponent</th>
                <th scope="col" className="py-2 pr-4 font-medium">Event</th>
                <th scope="col" className="py-2 font-medium">Result</th>
              </tr>
            </thead>
            <tbody>
              {bouts.map((bout, i) => (
                <tr key={i} className="border-b border-[var(--site-border)]">
                  <td className="py-3 pr-4 tabular-nums text-[var(--site-fg-muted)]">
                    {str(bout, "date") || str(bout, "bout_date")}
                  </td>
                  <td className="py-3 pr-4 font-medium">{itemText(bout, "opponent", locale)}</td>
                  <td className="py-3 pr-4 text-[var(--site-fg-muted)]">
                    {itemText(bout, "event", locale) || itemText(bout, "event_name", locale)}
                  </td>
                  <td className="py-3">
                    <ResultBadge result={str(bout, "result")} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </>
    );
  }

  if (variant === "cards") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {summary}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {bouts.map((bout, i) => (
            <div
              key={i}
              className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-5"
            >
              <div className="flex items-center justify-between gap-3">
                <ResultBadge result={str(bout, "result")} />
                <span className="text-xs tabular-nums text-[var(--site-fg-muted)]">
                  {str(bout, "date") || str(bout, "bout_date")}
                </span>
              </div>
              <h3 className="mt-3 font-medium">{itemText(bout, "opponent", locale)}</h3>
              {boutMeta(bout, locale) ? (
                <p className="mt-1 text-sm text-[var(--site-fg-muted)]">{boutMeta(bout, locale)}</p>
              ) : null}
            </div>
          ))}
        </div>
      </>
    );
  }

  // timeline
  return (
    <>
      <SectionHeading title={title} align={align} />
      {summary}
      <ol className="border-l border-[var(--site-border)] pl-6">
        {bouts.map((bout, i) => (
          <li key={i} className="relative pb-8 last:pb-0">
            <span
              className="absolute -left-[1.8rem] top-1 h-3 w-3 rounded-full bg-[var(--site-accent)] ring-4 ring-[var(--site-bg)]"
              aria-hidden="true"
            />
            <div className="flex flex-wrap items-center gap-3">
              <ResultBadge result={str(bout, "result")} />
              <span className="font-medium">{itemText(bout, "opponent", locale)}</span>
              <span className="text-xs tabular-nums text-[var(--site-fg-muted)]">
                {str(bout, "date") || str(bout, "bout_date")}
              </span>
            </div>
            {itemText(bout, "event", locale) || itemText(bout, "event_name", locale) ? (
              <p className="mt-1 text-sm">{itemText(bout, "event", locale) || itemText(bout, "event_name", locale)}</p>
            ) : null}
            {boutMeta(bout, locale) ? (
              <p className="mt-0.5 text-sm text-[var(--site-fg-muted)]">{boutMeta(bout, locale)}</p>
            ) : null}
          </li>
        ))}
      </ol>
    </>
  );
}
