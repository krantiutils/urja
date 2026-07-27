import Link from "next/link";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Dumbbell,
  Globe,
  Users,
  CalendarCheck,
  Package,
  UserCog,
  Heart,
  Activity,
  NfcIcon,
  Mail,
  ArrowRight,
} from "lucide-react";
import { MobileNav, FaqAccordion } from "@/components/landing/LandingInteractive";

export default function LandingPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const otherLocale = locale === "en" ? "ne" : "en";

  const features = [
    { icon: Users, title: t.landing.featureMembers, desc: t.landing.featureMembersDesc },
    { icon: CalendarCheck, title: t.landing.featureAttendance, desc: t.landing.featureAttendanceDesc },
    { icon: Package, title: t.landing.featurePackages, desc: t.landing.featurePackagesDesc },
    { icon: UserCog, title: t.landing.featureStaff, desc: t.landing.featureStaffDesc },
    { icon: Heart, title: t.landing.featureHealth, desc: t.landing.featureHealthDesc },
    { icon: Activity, title: t.landing.featureWorkouts, desc: t.landing.featureWorkoutsDesc },
    { icon: NfcIcon, title: t.landing.featureNfc, desc: t.landing.featureNfcDesc },
    { icon: Mail, title: t.landing.featureSms, desc: t.landing.featureSmsDesc },
  ];

  const steps = [
    { num: "01", title: t.landing.step1Title, desc: t.landing.step1Desc },
    { num: "02", title: t.landing.step2Title, desc: t.landing.step2Desc },
    { num: "03", title: t.landing.step3Title, desc: t.landing.step3Desc },
  ];

  const faqItems = [
    { question: t.landing.faq1Q, answer: t.landing.faq1A },
    { question: t.landing.faq2Q, answer: t.landing.faq2A },
    { question: t.landing.faq3Q, answer: t.landing.faq3A },
    { question: t.landing.faq4Q, answer: t.landing.faq4A },
    { question: t.landing.faq5Q, answer: t.landing.faq5A },
  ];


  return (
    <div className="min-h-screen bg-bg-deep text-fg">
      {/* ── Navbar ── */}
      <header className="sticky top-0 z-50 bg-bg-deep/80 backdrop-blur-xl border-b border-white/[0.06]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between relative">
          {/* Logo */}
          <Link href={`/${locale}`} className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-accent/20 flex items-center justify-center">
              <Dumbbell className="w-4 h-4 text-accent" />
            </div>
            <span className="text-lg font-bold tracking-tight">{t.common.appName}</span>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden lg:flex items-center gap-8">
            <a href="#features" className="text-sm text-fg-muted hover:text-fg transition-colors">
              {t.landing.navFeatures}
            </a>
            <a href="#faq" className="text-sm text-fg-muted hover:text-fg transition-colors">
              {t.landing.navFaq}
            </a>
          </nav>

          {/* Desktop actions */}
          <div className="hidden lg:flex items-center gap-3">
            <Link
              href={`/${otherLocale}`}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-fg-muted hover:text-fg rounded-lg hover:bg-surface transition-colors"
            >
              <Globe className="w-4 h-4" />
              {otherLocale === "en" ? "EN" : "ने"}
            </Link>
            <Link
              href={`/${locale}/login`}
              className="px-4 py-2 text-sm font-semibold bg-accent text-bg-deep rounded-lg hover:bg-accent-bright transition-colors shadow-accent-glow"
            >
              {t.landing.login}
            </Link>
          </div>

          {/* Mobile nav */}
          <MobileNav t={t} locale={locale} />
        </div>
      </header>

      {/* ── Hero ── */}
      <section className="relative overflow-hidden">
        {/* Ambient glow */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] rounded-full bg-accent/[0.07] blur-[120px] pointer-events-none" />

        <div className="relative max-w-6xl mx-auto px-4 sm:px-6 pt-24 pb-20 sm:pt-32 sm:pb-28 text-center">
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight leading-[1.1] animate-fade-in">
            <span className="bg-gradient-to-r from-fg via-accent-bright to-accent bg-clip-text text-transparent">
              {t.landing.heroTitle}
            </span>
          </h1>
          <p className="mt-6 text-lg sm:text-xl text-fg-muted max-w-2xl mx-auto leading-relaxed animate-fade-in-up">
            {t.landing.heroSubtitle}
          </p>
          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4 animate-fade-in-up">
            <Link
              href={`/${locale}/login`}
              className="w-full sm:w-auto px-8 py-3.5 text-sm font-semibold bg-accent text-bg-deep rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
            >
              {t.landing.getStarted}
            </Link>
            <a
              href="#features"
              className="w-full sm:w-auto px-8 py-3.5 text-sm font-semibold text-fg border border-white/[0.1] rounded-xl hover:bg-surface hover:border-white/[0.15] transition-colors"
            >
              {t.landing.learnMore}
            </a>
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section id="features" className="py-20 sm:py-28">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight">
              {t.landing.featuresTitle}
            </h2>
            <p className="mt-4 text-fg-muted text-lg max-w-xl mx-auto">
              {t.landing.featuresSubtitle}
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {features.map((f) => (
              <div
                key={f.title}
                className="group p-6 rounded-2xl border border-white/[0.06] bg-bg-elevated/50 backdrop-blur-sm hover:border-white/[0.1] hover:shadow-card-hover transition-all duration-300"
              >
                <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center mb-4 group-hover:bg-accent/20 transition-colors">
                  <f.icon className="w-5 h-5 text-accent" />
                </div>
                <h3 className="text-fg font-semibold mb-2">{f.title}</h3>
                <p className="text-sm text-fg-muted leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── How it works ── */}
      <section className="py-20 sm:py-28 border-t border-white/[0.04]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight">
              {t.landing.howItWorksTitle}
            </h2>
            <p className="mt-4 text-fg-muted text-lg max-w-xl mx-auto">
              {t.landing.howItWorksSubtitle}
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {steps.map((step, i) => (
              <div key={step.num} className="relative text-center">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-accent/10 border border-accent/20 mb-6">
                  <span className="text-2xl font-bold font-mono text-accent">{step.num}</span>
                </div>
                <h3 className="text-lg font-semibold mb-2">{step.title}</h3>
                <p className="text-sm text-fg-muted leading-relaxed">{step.desc}</p>
                {i < steps.length - 1 && (
                  <ArrowRight className="hidden md:block absolute top-8 -right-4 w-6 h-6 text-white/[0.1]" />
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── FAQ ── */}
      <section id="faq" className="py-20 sm:py-28 border-t border-white/[0.04]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight">
              {t.landing.faqTitle}
            </h2>
            <p className="mt-4 text-fg-muted text-lg max-w-xl mx-auto">
              {t.landing.faqSubtitle}
            </p>
          </div>

          <FaqAccordion items={faqItems} />
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="border-t border-white/[0.06] py-16">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {/* Brand */}
            <div className="col-span-2 md:col-span-1">
              <div className="flex items-center gap-2.5 mb-4">
                <div className="w-8 h-8 rounded-xl bg-accent/20 flex items-center justify-center">
                  <Dumbbell className="w-4 h-4 text-accent" />
                </div>
                <span className="text-lg font-bold">{t.common.appName}</span>
              </div>
              <p className="text-sm text-fg-muted leading-relaxed">
                {t.landing.heroSubtitle.slice(0, 80)}...
              </p>
            </div>

            {/* Product */}
            <div>
              <h4 className="text-sm font-semibold mb-4">{t.landing.footerProduct}</h4>
              <ul className="space-y-2.5">
                <li><a href="#features" className="text-sm text-fg-muted hover:text-fg transition-colors">{t.landing.navFeatures}</a></li>
                <li><a href="#faq" className="text-sm text-fg-muted hover:text-fg transition-colors">{t.landing.navFaq}</a></li>
              </ul>
            </div>

            {/* Company */}
            <div>
              <h4 className="text-sm font-semibold mb-4">{t.landing.footerCompany}</h4>
              <ul className="space-y-2.5">
                <li><span className="text-sm text-fg-muted">{t.landing.footerAbout}</span></li>
                <li><span className="text-sm text-fg-muted">{t.landing.footerCareers}</span></li>
                <li><span className="text-sm text-fg-muted">{t.landing.footerBlog}</span></li>
              </ul>
            </div>

            {/* Contact */}
            <div>
              <h4 className="text-sm font-semibold mb-4">{t.landing.footerContact}</h4>
              <ul className="space-y-2.5">
                <li><span className="text-sm text-fg-muted">{t.landing.footerEmail}</span></li>
                <li><span className="text-sm text-fg-muted">{t.landing.footerPhone}</span></li>
              </ul>
            </div>
          </div>

          <div className="mt-12 pt-8 border-t border-white/[0.06] text-center text-sm text-fg-muted">
            &copy; {new Date().getFullYear()} {t.common.appName}. {t.landing.footerRights}
          </div>
        </div>
      </footer>
    </div>
  );
}
