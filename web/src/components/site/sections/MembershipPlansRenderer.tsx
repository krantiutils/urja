import { Check } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { bool, list, itemText, num, str, strList, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * membership_plans — cards | table | compact
 *
 * Two data sources. "auto" (the default) renders the gym's real active
 * packages, fetched once by the page and passed in, so the website can never
 * drift from what the desk actually charges. "manual" renders whatever the
 * admin typed into the section, for gyms that price off-system.
 */

interface Plan {
  name: string;
  description: string;
  price: string;
  currency: string;
  durationLabel: string;
  features: string[];
  highlight: boolean;
}

function durationLabel(days: number | null, locale: Locale): string {
  if (!days || days <= 0) return "";
  if (days % 365 === 0) {
    const years = days / 365;
    return locale === "ne" ? `${years} वर्ष` : `${years} year${years > 1 ? "s" : ""}`;
  }
  if (days % 30 === 0) {
    const months = days / 30;
    return locale === "ne" ? `${months} महिना` : `${months} month${months > 1 ? "s" : ""}`;
  }
  return locale === "ne" ? `${days} दिन` : `${days} days`;
}

/** Trims the trailing ".00" the API sends on whole-rupee prices. */
function formatPrice(raw: string): string {
  const n = Number(raw);
  if (!Number.isFinite(n)) return raw;
  return n % 1 === 0 ? String(n) : n.toFixed(2);
}

/**
 * Packages come from the REST API, which uses `name_ne` — not the `nameNe`
 * convention section JSON uses — so this reads the translation explicitly
 * rather than going through `itemText`.
 */
function apiText(pkg: Record<string, unknown>, key: string, locale: Locale): string {
  if (locale === "ne") {
    const ne = pkg[`${key}_ne`];
    if (typeof ne === "string" && ne.trim() !== "") return ne;
  }
  return str(pkg, key);
}

function fromPackages(packages: Array<Record<string, unknown>>, locale: Locale): Plan[] {
  return packages.map((pkg) => ({
    name: apiText(pkg, "name", locale),
    description: apiText(pkg, "description", locale),
    price: formatPrice(str(pkg, "price")),
    currency: str(pkg, "currency", "NPR"),
    durationLabel: durationLabel(num(pkg, "duration_days"), locale),
    features: strList(pkg, "features"),
    highlight: false,
  }));
}

function fromContent(content: Record<string, unknown>, locale: Locale): Plan[] {
  return list(content, "plans").map((plan) => ({
    name: itemText(plan, "name", locale),
    description: itemText(plan, "description", locale),
    price: str(plan, "price"),
    currency: str(plan, "currency", "NPR"),
    durationLabel: itemText(plan, "duration", locale),
    features: strList(plan, locale === "ne" ? "featuresNe" : "features").length
      ? strList(plan, locale === "ne" ? "featuresNe" : "features")
      : strList(plan, "features"),
    highlight: bool(plan, "highlight"),
  }));
}

function PriceTag({ plan }: { plan: Plan }) {
  if (!plan.price) return null;
  return (
    <p className="flex items-baseline justify-center gap-1.5">
      <span className="text-sm text-[var(--site-fg-muted)]">{plan.currency}</span>
      <span className="text-3xl sm:text-4xl tabular-nums" style={{ fontFamily: "var(--site-font-display)" }}>
        {plan.price}
      </span>
      {plan.durationLabel ? (
        <span className="text-sm text-[var(--site-fg-muted)]">/ {plan.durationLabel}</span>
      ) : null}
    </p>
  );
}

export default function MembershipPlansRenderer({
  section,
  locale,
  packages,
}: {
  section: Section;
  locale: Locale;
  packages: Array<Record<string, unknown>>;
}) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const note = text(content, "note", locale);
  const align = section.style?.align ?? "left";

  const manual = str(content, "dataSource", "auto") === "manual";
  const plans = (manual ? fromContent(content, locale) : fromPackages(packages, locale)).filter(
    (plan) => plan.name !== ""
  );

  if (plans.length === 0) {
    // A gym mid-setup should not publish an empty pricing table, but the
    // heading alone is worth keeping so the page does not lose its anchor.
    if (!title) return null;
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {note ? <p className="text-sm text-[var(--site-fg-muted)]">{note}</p> : null}
      </>
    );
  }

  const footnote = note ? (
    <p className="mt-8 text-sm text-[var(--site-fg-muted)]">{note}</p>
  ) : null;

  if (variant === "table") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        <div className="overflow-x-auto">
          <table className="w-full min-w-[32rem] border-collapse text-left">
            <thead>
              <tr className="border-b border-[var(--site-border)] text-sm text-[var(--site-fg-muted)]">
                <th scope="col" className="py-3 pr-4 font-medium">Plan</th>
                <th scope="col" className="py-3 pr-4 font-medium">Includes</th>
                <th scope="col" className="py-3 text-right font-medium">Price</th>
              </tr>
            </thead>
            <tbody>
              {plans.map((plan, i) => (
                <tr key={i} className="border-b border-[var(--site-border)] align-top">
                  <td className="py-4 pr-4">
                    <span className="font-medium">{plan.name}</span>
                    {plan.durationLabel ? (
                      <span className="block text-sm text-[var(--site-fg-muted)]">{plan.durationLabel}</span>
                    ) : null}
                  </td>
                  <td className="py-4 pr-4 text-sm text-[var(--site-fg-muted)]">
                    {plan.features.length > 0 ? plan.features.join(", ") : plan.description}
                  </td>
                  <td className="py-4 text-right tabular-nums whitespace-nowrap">
                    {plan.currency} {plan.price}
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

  if (variant === "compact") {
    return (
      <>
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        <div className="divide-y divide-[var(--site-border)] border-t border-b border-[var(--site-border)]">
          {plans.map((plan, i) => (
            <div key={i} className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 py-4">
              <div>
                <span className="font-medium">{plan.name}</span>
                {plan.durationLabel ? (
                  <span className="ml-2 text-sm text-[var(--site-fg-muted)]">{plan.durationLabel}</span>
                ) : null}
              </div>
              <span className="tabular-nums whitespace-nowrap">
                {plan.currency} {plan.price}
              </span>
            </div>
          ))}
        </div>
        {footnote}
      </>
    );
  }

  // cards
  return (
    <>
      <SectionHeading title={title} subtitle={subtitle} align={align} />
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {plans.map((plan, i) => (
          <div
            key={i}
            className={`flex flex-col rounded-[var(--site-radius)] border bg-[var(--site-bg)] p-6 text-center ${
              plan.highlight ? "border-[var(--site-accent)] ring-1 ring-[var(--site-accent)]" : "border-[var(--site-border)]"
            }`}
          >
            <h3 className="text-lg font-medium">{plan.name}</h3>
            {plan.durationLabel ? (
              <p className="mt-1 text-sm text-[var(--site-fg-muted)]">{plan.durationLabel}</p>
            ) : null}
            <div className="my-5">
              <PriceTag plan={plan} />
            </div>
            {plan.description ? (
              <p className="text-sm text-[var(--site-fg-muted)] leading-relaxed">{plan.description}</p>
            ) : null}
            {plan.features.length > 0 ? (
              <ul className="mt-5 flex flex-col gap-2 text-left text-sm">
                {plan.features.map((feature, j) => (
                  <li key={j} className="flex items-start gap-2">
                    <Check className="mt-0.5 h-4 w-4 shrink-0 text-[var(--site-accent)]" aria-hidden="true" />
                    <span>{feature}</span>
                  </li>
                ))}
              </ul>
            ) : null}
          </div>
        ))}
      </div>
      {footnote}
    </>
  );
}
