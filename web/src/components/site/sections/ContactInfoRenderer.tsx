import { Clock, Mail, MapPin, Phone } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { text, str } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * contact_info — list | card | two_column
 */

interface Row {
  icon: typeof MapPin;
  label: string;
  href?: string;
}

function buildRows(content: Section["content"], locale: Locale): Row[] {
  const rows: Row[] = [];
  const address = text(content, "address", locale);
  const phone = str(content, "phone");
  const email = str(content, "email");
  const hoursNote = text(content, "hoursNote", locale);

  if (address) rows.push({ icon: MapPin, label: address });
  if (phone) rows.push({ icon: Phone, label: phone, href: `tel:${phone}` });
  if (email) rows.push({ icon: Mail, label: email, href: `mailto:${email}` });
  if (hoursNote) rows.push({ icon: Clock, label: hoursNote });
  return rows;
}

function RowItem({ row }: { row: Row }) {
  const Icon = row.icon;
  const inner = (
    <>
      <Icon className="h-5 w-5 shrink-0 text-[var(--site-accent)]" aria-hidden="true" />
      <span>{row.label}</span>
    </>
  );
  return (
    <li className="flex items-start gap-3">
      {row.href ? (
        <a href={row.href} className="flex items-start gap-3 hover:underline">
          {inner}
        </a>
      ) : (
        inner
      )}
    </li>
  );
}

export default function ContactInfoRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const rows = buildRows(content, locale);

  if (rows.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  if (variant === "card") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        <div className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] p-6 max-w-md">
          <ul className="flex flex-col gap-4">
            {rows.map((row, i) => (
              <RowItem key={i} row={row} />
            ))}
          </ul>
        </div>
      </>
    );
  }

  if (variant === "two_column") {
    const left = rows.filter((r) => r.icon === MapPin || r.icon === Clock);
    const right = rows.filter((r) => r.icon === Phone || r.icon === Mail);
    return (
      <>
        <SectionHeading title={title} align={align} />
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
          {left.length > 0 ? (
            <ul className="flex flex-col gap-4">
              {left.map((row, i) => (
                <RowItem key={i} row={row} />
              ))}
            </ul>
          ) : null}
          {right.length > 0 ? (
            <ul className="flex flex-col gap-4">
              {right.map((row, i) => (
                <RowItem key={i} row={row} />
              ))}
            </ul>
          ) : null}
        </div>
      </>
    );
  }

  // list
  return (
    <>
      <SectionHeading title={title} align={align} />
      <ul className="flex flex-col gap-4">
        {rows.map((row, i) => (
          <RowItem key={i} row={row} />
        ))}
      </ul>
    </>
  );
}
