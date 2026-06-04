import type { Dictionary } from "@/lib/i18n";
import {
  Users,
  CalendarCheck,
  Package,
  CreditCard,
  UserCog,
  BarChart3,
} from "lucide-react";

interface FeaturesProps {
  t: Dictionary;
}

export function Features({ t }: FeaturesProps) {
  const features = [
    {
      icon: Users,
      title: t.landing.featureMembersTitle,
      desc: t.landing.featureMembersDesc,
    },
    {
      icon: CalendarCheck,
      title: t.landing.featureAttendanceTitle,
      desc: t.landing.featureAttendanceDesc,
    },
    {
      icon: Package,
      title: t.landing.featurePackagesTitle,
      desc: t.landing.featurePackagesDesc,
    },
    {
      icon: CreditCard,
      title: t.landing.featurePaymentsTitle,
      desc: t.landing.featurePaymentsDesc,
    },
    {
      icon: UserCog,
      title: t.landing.featureStaffTitle,
      desc: t.landing.featureStaffDesc,
    },
    {
      icon: BarChart3,
      title: t.landing.featureAnalyticsTitle,
      desc: t.landing.featureAnalyticsDesc,
    },
  ];

  return (
    <section id="features" className="py-20 lg:py-28">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="text-center mb-16">
          <span className="inline-block px-4 py-1.5 text-xs font-medium text-md-primary bg-md-primary-container rounded-full mb-4 uppercase tracking-wider">
            {t.landing.featuresTag}
          </span>
          <h2 className="text-3xl sm:text-4xl font-bold text-md-on-surface tracking-tight">
            {t.landing.featuresTitle}
          </h2>
          <p className="mt-4 text-lg text-md-on-surface-variant max-w-2xl mx-auto">
            {t.landing.featuresSubtitle}
          </p>
        </div>

        {/* Feature grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="bg-md-surface-container rounded-3xl p-6 shadow-md-1 hover:shadow-md-2 transition-shadow duration-300 group"
            >
              <div className="w-12 h-12 rounded-xl bg-md-primary-container flex items-center justify-center mb-4 group-hover:bg-md-primary group-hover:text-md-on-primary transition-colors">
                <feature.icon className="w-6 h-6 text-md-primary group-hover:text-md-on-primary transition-colors" />
              </div>
              <h3 className="text-lg font-semibold text-md-on-surface mb-2">
                {feature.title}
              </h3>
              <p className="text-sm text-md-on-surface-variant leading-relaxed">
                {feature.desc}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
