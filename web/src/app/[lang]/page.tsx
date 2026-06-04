import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import { Navbar } from "@/components/landing/Navbar";
import { Hero } from "@/components/landing/Hero";
import { Features } from "@/components/landing/Features";
import { HowItWorks } from "@/components/landing/HowItWorks";
import { Pricing } from "@/components/landing/Pricing";
import { FAQ } from "@/components/landing/FAQ";
import { Footer } from "@/components/landing/Footer";

export default function LandingPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  return (
    <main className="min-h-screen bg-md-background">
      <Navbar t={t} locale={locale} />
      <Hero t={t} locale={locale} />
      <Features t={t} />
      <HowItWorks t={t} />
      <Pricing t={t} locale={locale} />
      <FAQ t={t} />
      <Footer t={t} />
    </main>
  );
}
