import type { Dictionary } from "@/lib/i18n";

interface HowItWorksProps {
  t: Dictionary;
}

export function HowItWorks({ t }: HowItWorksProps) {
  const steps = [
    { num: "1", title: t.landing.howStep1Title, desc: t.landing.howStep1Desc },
    { num: "2", title: t.landing.howStep2Title, desc: t.landing.howStep2Desc },
    { num: "3", title: t.landing.howStep3Title, desc: t.landing.howStep3Desc },
  ];

  return (
    <section className="py-20 lg:py-28 bg-md-surface-container-low">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="text-center mb-16">
          <span className="inline-block px-4 py-1.5 text-xs font-medium text-md-primary bg-md-primary-container rounded-full mb-4 uppercase tracking-wider">
            {t.landing.howTag}
          </span>
          <h2 className="text-3xl sm:text-4xl font-bold text-md-on-surface tracking-tight">
            {t.landing.howTitle}
          </h2>
        </div>

        {/* Steps */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12">
          {steps.map((step, idx) => (
            <div key={step.num} className="relative text-center">
              {/* Connecting line (desktop only) */}
              {idx < steps.length - 1 && (
                <div className="hidden md:block absolute top-8 left-[60%] w-[80%] h-0.5 bg-md-outline-variant" />
              )}

              {/* Number circle */}
              <div className="relative inline-flex items-center justify-center w-16 h-16 rounded-full bg-md-primary text-md-on-primary text-2xl font-bold mb-6 shadow-md-2">
                {step.num}
              </div>

              <h3 className="text-xl font-semibold text-md-on-surface mb-3">
                {step.title}
              </h3>
              <p className="text-sm text-md-on-surface-variant leading-relaxed max-w-xs mx-auto">
                {step.desc}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
