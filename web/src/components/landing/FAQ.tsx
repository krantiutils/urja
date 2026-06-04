"use client";

import { useState } from "react";
import type { Dictionary } from "@/lib/i18n";
import { ChevronDown } from "lucide-react";

interface FAQProps {
  t: Dictionary;
}

export function FAQ({ t }: FAQProps) {
  const items = [
    { q: t.landing.faq1Q, a: t.landing.faq1A },
    { q: t.landing.faq2Q, a: t.landing.faq2A },
    { q: t.landing.faq3Q, a: t.landing.faq3A },
    { q: t.landing.faq4Q, a: t.landing.faq4A },
    { q: t.landing.faq5Q, a: t.landing.faq5A },
  ];

  const [openIndex, setOpenIndex] = useState<number | null>(null);

  return (
    <section id="faq" className="py-20 lg:py-28 bg-md-surface-container-low">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="text-center mb-12">
          <span className="inline-block px-4 py-1.5 text-xs font-medium text-md-primary bg-md-primary-container rounded-full mb-4 uppercase tracking-wider">
            {t.landing.faqTag}
          </span>
          <h2 className="text-3xl sm:text-4xl font-bold text-md-on-surface tracking-tight">
            {t.landing.faqTitle}
          </h2>
        </div>

        {/* Accordion */}
        <div className="space-y-3">
          {items.map((item, idx) => {
            const isOpen = openIndex === idx;
            return (
              <div
                key={idx}
                className="bg-md-surface-container rounded-2xl overflow-hidden shadow-md-1"
              >
                <button
                  onClick={() => setOpenIndex(isOpen ? null : idx)}
                  className="w-full flex items-center justify-between px-6 py-4 text-left"
                >
                  <span className="text-sm font-medium text-md-on-surface pr-4">
                    {item.q}
                  </span>
                  <ChevronDown
                    className={`w-5 h-5 text-md-on-surface-variant flex-shrink-0 transition-transform duration-200 ${
                      isOpen ? "rotate-180" : ""
                    }`}
                  />
                </button>
                {isOpen && (
                  <div className="px-6 pb-4">
                    <p className="text-sm text-md-on-surface-variant leading-relaxed">
                      {item.a}
                    </p>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
