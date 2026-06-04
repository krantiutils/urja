import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import { Check } from "lucide-react";

interface PricingProps {
  t: Dictionary;
  locale: Locale;
}

export function Pricing({ t, locale }: PricingProps) {
  const tiers = [
    {
      name: t.landing.pricingBasicName,
      price: t.landing.pricingBasicPrice,
      period: t.landing.pricingBasicPeriod,
      features: [
        t.landing.pricingBasicF1,
        t.landing.pricingBasicF2,
        t.landing.pricingBasicF3,
        t.landing.pricingBasicF4,
      ],
      cta: t.landing.pricingBasicCta,
      featured: false,
    },
    {
      name: t.landing.pricingProName,
      price: t.landing.pricingProPrice,
      period: t.landing.pricingProPeriod,
      features: [
        t.landing.pricingProF1,
        t.landing.pricingProF2,
        t.landing.pricingProF3,
        t.landing.pricingProF4,
        t.landing.pricingProF5,
      ],
      cta: t.landing.pricingProCta,
      featured: true,
    },
    {
      name: t.landing.pricingEnterpriseName,
      price: t.landing.pricingEnterprisePrice,
      period: t.landing.pricingEnterprisePeriod,
      features: [
        t.landing.pricingEnterpriseF1,
        t.landing.pricingEnterpriseF2,
        t.landing.pricingEnterpriseF3,
        t.landing.pricingEnterpriseF4,
        t.landing.pricingEnterpriseF5,
      ],
      cta: t.landing.pricingEnterpriseCta,
      featured: false,
    },
  ];

  return (
    <section id="pricing" className="py-20 lg:py-28">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="text-center mb-16">
          <span className="inline-block px-4 py-1.5 text-xs font-medium text-md-primary bg-md-primary-container rounded-full mb-4 uppercase tracking-wider">
            {t.landing.pricingTag}
          </span>
          <h2 className="text-3xl sm:text-4xl font-bold text-md-on-surface tracking-tight">
            {t.landing.pricingTitle}
          </h2>
          <p className="mt-4 text-lg text-md-on-surface-variant max-w-xl mx-auto">
            {t.landing.pricingSubtitle}
          </p>
        </div>

        {/* Pricing cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {tiers.map((tier) => (
            <div
              key={tier.name}
              className={`rounded-3xl p-8 flex flex-col ${
                tier.featured
                  ? "bg-md-primary text-md-on-primary shadow-md-3 scale-105"
                  : "bg-md-surface-container shadow-md-1"
              }`}
            >
              <h3
                className={`text-lg font-semibold mb-2 ${
                  tier.featured
                    ? "text-md-on-primary"
                    : "text-md-on-surface"
                }`}
              >
                {tier.name}
              </h3>
              <div className="mb-6">
                <span
                  className={`text-4xl font-bold ${
                    tier.featured
                      ? "text-md-on-primary"
                      : "text-md-on-surface"
                  }`}
                >
                  {tier.price}
                </span>
                {tier.period && (
                  <span
                    className={`text-sm ml-1 ${
                      tier.featured
                        ? "text-md-on-primary/70"
                        : "text-md-on-surface-variant"
                    }`}
                  >
                    {tier.period}
                  </span>
                )}
              </div>

              <ul className="space-y-3 mb-8 flex-1">
                {tier.features.map((f) => (
                  <li key={f} className="flex items-start gap-3">
                    <Check
                      className={`w-5 h-5 flex-shrink-0 mt-0.5 ${
                        tier.featured
                          ? "text-md-on-primary/80"
                          : "text-md-primary"
                      }`}
                    />
                    <span
                      className={`text-sm ${
                        tier.featured
                          ? "text-md-on-primary/90"
                          : "text-md-on-surface-variant"
                      }`}
                    >
                      {f}
                    </span>
                  </li>
                ))}
              </ul>

              <Link
                href={`/${locale}/login`}
                className={`block w-full text-center px-6 py-3 text-sm font-medium rounded-full transition-all ${
                  tier.featured
                    ? "bg-md-on-primary text-md-primary hover:bg-md-on-primary/90 shadow-md-1"
                    : "bg-md-primary text-md-on-primary hover:bg-md-primary/90 shadow-md-1"
                }`}
              >
                {tier.cta}
              </Link>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
