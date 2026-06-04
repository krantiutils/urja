import type { Dictionary } from "@/lib/i18n";
import { Dumbbell } from "lucide-react";

interface FooterProps {
  t: Dictionary;
}

export function Footer({ t }: FooterProps) {
  return (
    <footer className="bg-md-surface-container-high py-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-10">
          {/* Brand */}
          <div className="md:col-span-1">
            <div className="flex items-center gap-2.5 mb-4">
              <div className="w-9 h-9 rounded-xl bg-md-primary flex items-center justify-center">
                <Dumbbell className="w-5 h-5 text-md-on-primary" />
              </div>
              <span className="text-xl font-bold text-md-on-surface">
                {t.common.appName}
              </span>
            </div>
            <p className="text-sm text-md-on-surface-variant">
              {t.landing.footerTagline}
            </p>
          </div>

          {/* Product */}
          <div>
            <h4 className="text-sm font-semibold text-md-on-surface mb-4">
              {t.landing.footerProduct}
            </h4>
            <ul className="space-y-2.5">
              <li>
                <a
                  href="#features"
                  className="text-sm text-md-on-surface-variant hover:text-md-primary transition-colors"
                >
                  {t.landing.featuresTag}
                </a>
              </li>
              <li>
                <a
                  href="#pricing"
                  className="text-sm text-md-on-surface-variant hover:text-md-primary transition-colors"
                >
                  {t.landing.pricingTag}
                </a>
              </li>
              <li>
                <a
                  href="#faq"
                  className="text-sm text-md-on-surface-variant hover:text-md-primary transition-colors"
                >
                  {t.landing.faqTag}
                </a>
              </li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h4 className="text-sm font-semibold text-md-on-surface mb-4">
              {t.landing.footerCompany}
            </h4>
            <ul className="space-y-2.5">
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  {t.landing.footerAbout}
                </span>
              </li>
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  {t.landing.footerBlog}
                </span>
              </li>
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  {t.landing.footerCareers}
                </span>
              </li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h4 className="text-sm font-semibold text-md-on-surface mb-4">
              {t.landing.footerContact}
            </h4>
            <ul className="space-y-2.5">
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  {t.landing.footerPrivacy}
                </span>
              </li>
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  {t.landing.footerTerms}
                </span>
              </li>
              <li>
                <span className="text-sm text-md-on-surface-variant">
                  info@nepalgym.xyz
                </span>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="mt-12 pt-6 border-t border-md-outline-variant">
          <p className="text-xs text-md-on-surface-variant text-center">
            &copy; {new Date().getFullYear()} {t.common.appName}. {t.landing.footerRights}
          </p>
        </div>
      </div>
    </footer>
  );
}
