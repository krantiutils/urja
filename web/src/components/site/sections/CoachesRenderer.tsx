import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, strList, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";
import { SiteImage } from "@/components/site/primitives/SiteImage";

/**
 * coaches — cards | list | spotlight
 */

function CredentialTags({ credentials }: { credentials: string[] }) {
  if (credentials.length === 0) return null;
  return (
    <ul className="mt-3 flex flex-wrap gap-2">
      {credentials.map((c, i) => (
        <li
          key={i}
          className="rounded-[var(--site-radius)] border border-[var(--site-border)] px-2 py-1 text-xs text-[var(--site-fg-muted)]"
        >
          {c}
        </li>
      ))}
    </ul>
  );
}

export default function CoachesRenderer({ section, locale }: { section: Section; locale: Locale }) {
  const { content, variant } = section;
  const title = text(content, "title", locale);
  const coaches = list(content, "coaches");

  if (coaches.length === 0 && !title) return null;

  const align = section.style?.align ?? "left";

  if (variant === "list") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {coaches.length > 0 ? (
          <div className="flex flex-col divide-y divide-[var(--site-border)] border-t border-b border-[var(--site-border)]">
            {coaches.map((coach, i) => (
              <div key={i} className="flex items-center gap-4 py-5">
                <SiteImage
                  src={str(coach, "photo")}
                  alt={itemText(coach, "name", locale)}
                  ratio="aspect-square"
                  className="w-16 h-16 shrink-0 rounded-full"
                />
                <div className="min-w-0">
                  <h3 className="font-medium">{itemText(coach, "name", locale)}</h3>
                  <p className="text-sm text-[var(--site-fg-muted)]">{itemText(coach, "role", locale)}</p>
                </div>
                {str(coach, "record") ? (
                  <span className="ml-auto shrink-0 text-sm text-[var(--site-accent)]">{str(coach, "record")}</span>
                ) : null}
              </div>
            ))}
          </div>
        ) : null}
      </>
    );
  }

  if (variant === "spotlight") {
    return (
      <>
        <SectionHeading title={title} align={align} />
        {coaches.length > 0 ? (
          <div className="flex flex-col gap-12">
            {coaches.map((coach, i) => (
              <div key={i} className="grid gap-8 lg:grid-cols-2 lg:items-center">
                <SiteImage
                  src={str(coach, "photo")}
                  alt={itemText(coach, "name", locale)}
                  ratio="aspect-[4/3]"
                  className={i % 2 === 1 ? "lg:order-last" : ""}
                />
                <div>
                  <h3 className="text-2xl font-medium">{itemText(coach, "name", locale)}</h3>
                  <p className="mt-1 text-[var(--site-accent)]">{itemText(coach, "role", locale)}</p>
                  {str(coach, "record") ? (
                    <p className="mt-1 text-sm text-[var(--site-fg-muted)]">Record: {str(coach, "record")}</p>
                  ) : null}
                  {itemText(coach, "bio", locale) ? (
                    <p className="mt-4 leading-relaxed text-[var(--site-fg-muted)]">{itemText(coach, "bio", locale)}</p>
                  ) : null}
                  <CredentialTags credentials={strList(coach, "credentials")} />
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
      <SectionHeading title={title} align={align} />
      {coaches.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {coaches.map((coach, i) => (
            <div
              key={i}
              className="rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-surface)] overflow-hidden"
            >
              <SiteImage src={str(coach, "photo")} alt={itemText(coach, "name", locale)} ratio="aspect-square" />
              <div className="p-5">
                <h3 className="font-medium">{itemText(coach, "name", locale)}</h3>
                <p className="text-sm text-[var(--site-accent)]">{itemText(coach, "role", locale)}</p>
                {itemText(coach, "bio", locale) ? (
                  <p className="mt-3 text-sm text-[var(--site-fg-muted)] leading-relaxed">
                    {itemText(coach, "bio", locale)}
                  </p>
                ) : null}
                <CredentialTags credentials={strList(coach, "credentials")} />
              </div>
            </div>
          ))}
        </div>
      ) : null}
    </>
  );
}
