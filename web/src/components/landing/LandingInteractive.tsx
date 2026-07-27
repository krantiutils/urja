"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, X, ChevronDown, Globe } from "lucide-react";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";

interface MobileNavProps {
  t: Dictionary;
  locale: Locale;
}

export function MobileNav({ t, locale }: MobileNavProps) {
  const [open, setOpen] = useState(false);
  const otherLocale = locale === "en" ? "ne" : "en";

  return (
    <>
      <button
        onClick={() => setOpen(!open)}
        className="lg:hidden p-2 rounded-lg hover:bg-surface transition-colors"
        aria-label="Toggle menu"
      >
        {open ? (
          <X className="w-5 h-5 text-fg" />
        ) : (
          <Menu className="w-5 h-5 text-fg" />
        )}
      </button>

      {open && (
        <div className="absolute top-full left-0 right-0 bg-bg-elevated/95 backdrop-blur-xl border-b border-white/[0.06] lg:hidden">
          <nav className="flex flex-col p-4 gap-1">
            <a
              href="#features"
              onClick={() => setOpen(false)}
              className="px-4 py-3 text-fg-muted hover:text-fg hover:bg-surface rounded-lg transition-colors"
            >
              {t.landing.navFeatures}
            </a>
            <a
              href="#faq"
              onClick={() => setOpen(false)}
              className="px-4 py-3 text-fg-muted hover:text-fg hover:bg-surface rounded-lg transition-colors"
            >
              {t.landing.navFaq}
            </a>
            <Link
              href={`/${otherLocale}`}
              className="flex items-center gap-2 px-4 py-3 text-fg-muted hover:text-fg hover:bg-surface rounded-lg transition-colors"
            >
              <Globe className="w-4 h-4" />
              {otherLocale === "en" ? "English" : "नेपाली"}
            </Link>
            <Link
              href={`/${locale}/login`}
              className="mt-2 px-4 py-3 text-center bg-accent text-bg-deep font-semibold rounded-lg hover:bg-accent-bright transition-colors"
            >
              {t.landing.login}
            </Link>
          </nav>
        </div>
      )}
    </>
  );
}

interface FaqAccordionProps {
  items: { question: string; answer: string }[];
}

export function FaqAccordion({ items }: FaqAccordionProps) {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  return (
    <div className="space-y-3 max-w-2xl mx-auto">
      {items.map((item, i) => (
        <div
          key={i}
          className="rounded-xl border border-white/[0.06] bg-bg-elevated/50 backdrop-blur-sm overflow-hidden"
        >
          <button
            onClick={() => setOpenIndex(openIndex === i ? null : i)}
            className="w-full flex items-center justify-between px-6 py-4 text-left hover:bg-surface transition-colors"
          >
            <span className="text-fg font-medium pr-4">{item.question}</span>
            <ChevronDown
              className={`w-5 h-5 text-fg-muted flex-shrink-0 transition-transform duration-200 ${
                openIndex === i ? "rotate-180" : ""
              }`}
            />
          </button>
          {openIndex === i && (
            <div className="px-6 pb-4 text-fg-muted leading-relaxed">
              {item.answer}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
