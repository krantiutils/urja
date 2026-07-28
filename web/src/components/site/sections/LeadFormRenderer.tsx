"use client";

import { useState } from "react";
import { Check, Loader2 } from "lucide-react";
import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { list, itemText, str, text } from "@/lib/site/content";
import { SectionHeading } from "@/components/site/primitives/SectionShell";

/**
 * lead_form — inline | card | split
 *
 * Posts to the public leads endpoint, which applies its own per-phone rate
 * limit. The `website` field is the honeypot the API expects: it is hidden
 * from people and from assistive technology, so anything that fills it is a
 * bot. The server accepts those submissions with a 201 and discards them, so
 * this component deliberately shows the same success state either way —
 * telling a bot it was caught only teaches it to avoid the trap.
 */

type Status = "idle" | "sending" | "sent" | "error";

export default function LeadFormRenderer({
  section,
  locale,
  siteSlug,
  sourcePage,
}: {
  section: Section;
  locale: Locale;
  siteSlug: string;
  sourcePage?: string;
}) {
  const { content, variant } = section;
  const [status, setStatus] = useState<Status>("idle");
  const [error, setError] = useState("");

  const title = text(content, "title", locale);
  const subtitle = text(content, "subtitle", locale);
  const submitLabel = text(content, "submitLabel", locale) || (locale === "ne" ? "पठाउनुहोस्" : "Send");
  const successMessage =
    text(content, "successMessage", locale) ||
    (locale === "ne" ? "धन्यवाद — हामी सम्पर्क गर्नेछौं।" : "Thanks — we will be in touch.");
  const interests = list(content, "interests");

  const labels =
    locale === "ne"
      ? { name: "नाम", phone: "फोन", email: "इमेल", message: "सन्देश", interest: "रुचि", optional: "वैकल्पिक" }
      : { name: "Name", phone: "Phone", email: "Email", message: "Message", interest: "Interest", optional: "optional" };

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    setStatus("sending");
    setError("");

    try {
      // Same-origin proxy rather than the API directly: see
      // app/api/site-leads/route.ts.
      const res = await fetch("/api/site-leads", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          site: siteSlug,
          name: form.get("name"),
          phone: form.get("phone"),
          email: form.get("email") || "",
          message: form.get("message") || "",
          interest: form.get("interest") || "",
          source_page: sourcePage || "",
          website: form.get("website") || "",
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        setError(typeof body?.error === "string" ? body.error : "Something went wrong.");
        setStatus("error");
        return;
      }
      setStatus("sent");
    } catch {
      setError(locale === "ne" ? "पठाउन सकिएन।" : "Could not send. Please try again.");
      setStatus("error");
    }
  }

  const align = section.style?.align ?? "left";

  if (status === "sent") {
    return (
      <div className="flex flex-col items-center gap-3 py-8 text-center">
        <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--site-accent)] text-[var(--site-accent-fg)]">
          <Check className="h-6 w-6" aria-hidden="true" />
        </span>
        <p className="text-lg" role="status">
          {successMessage}
        </p>
      </div>
    );
  }

  const fieldClasses =
    "w-full rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-bg)] px-4 py-3 text-[var(--site-fg)] placeholder:text-[var(--site-fg-muted)] focus:border-[var(--site-accent)] focus:outline-none";

  const form = (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 text-left">
      {/* Honeypot. aria-hidden + tabIndex keep it away from real users. */}
      <div className="absolute h-px w-px overflow-hidden opacity-0" aria-hidden="true">
        <label htmlFor={`website-${section.id}`}>Website</label>
        <input id={`website-${section.id}`} name="website" type="text" tabIndex={-1} autoComplete="off" />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label htmlFor={`name-${section.id}`} className="mb-1.5 block text-sm font-medium">
            {labels.name}
          </label>
          <input id={`name-${section.id}`} name="name" type="text" required maxLength={120} className={fieldClasses} />
        </div>
        <div>
          <label htmlFor={`phone-${section.id}`} className="mb-1.5 block text-sm font-medium">
            {labels.phone}
          </label>
          <input
            id={`phone-${section.id}`}
            name="phone"
            type="tel"
            inputMode="numeric"
            required
            maxLength={20}
            className={fieldClasses}
          />
        </div>
      </div>

      {interests.length > 0 ? (
        <div>
          <label htmlFor={`interest-${section.id}`} className="mb-1.5 block text-sm font-medium">
            {labels.interest}
          </label>
          <select id={`interest-${section.id}`} name="interest" defaultValue="" className={fieldClasses}>
            <option value="">—</option>
            {interests.map((item, i) => (
              <option key={i} value={str(item, "value")}>
                {itemText(item, "label", locale)}
              </option>
            ))}
          </select>
        </div>
      ) : null}

      <div>
        <label htmlFor={`message-${section.id}`} className="mb-1.5 block text-sm font-medium">
          {labels.message}{" "}
          <span className="font-normal text-[var(--site-fg-muted)]">({labels.optional})</span>
        </label>
        <textarea id={`message-${section.id}`} name="message" rows={3} maxLength={1000} className={fieldClasses} />
      </div>

      {status === "error" && error ? (
        <p role="alert" className="text-sm text-[var(--site-fg)]">
          {error}
        </p>
      ) : null}

      <button
        type="submit"
        disabled={status === "sending"}
        className="inline-flex items-center justify-center gap-2 rounded-[var(--site-radius)] bg-[var(--site-accent)] px-6 py-3 font-medium text-[var(--site-accent-fg)] transition-opacity hover:opacity-90 disabled:opacity-60"
      >
        {status === "sending" ? <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" /> : null}
        {submitLabel}
      </button>
    </form>
  );

  if (variant === "split") {
    return (
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-16 items-start">
        <div>
          <SectionHeading title={title} subtitle={subtitle} align={align} />
        </div>
        <div className="relative">{form}</div>
      </div>
    );
  }

  if (variant === "card") {
    return (
      <div className="relative mx-auto max-w-xl rounded-[var(--site-radius)] border border-[var(--site-border)] bg-[var(--site-bg)] p-6 sm:p-8">
        <SectionHeading title={title} subtitle={subtitle} align={align} />
        {form}
      </div>
    );
  }

  // inline
  return (
    <div className="relative">
      <SectionHeading title={title} subtitle={subtitle} align={align} />
      {form}
    </div>
  );
}
