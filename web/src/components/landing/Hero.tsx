"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  ArrowRight,
  Dumbbell,
  Users,
  CalendarCheck,
  TrendingUp,
  Zap,
} from "lucide-react";

interface HeroProps {
  t: Dictionary;
  locale: Locale;
}

/* Animated counter that counts up from 0 to target */
function AnimatedNumber({ target, suffix = "" }: { target: number; suffix?: string }) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    let frame: number;
    const duration = 2000;
    const start = performance.now();

    function tick(now: number) {
      const elapsed = now - start;
      const progress = Math.min(elapsed / duration, 1);
      // ease-out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      setCount(Math.floor(eased * target));
      if (progress < 1) frame = requestAnimationFrame(tick);
    }

    frame = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(frame);
  }, [target]);

  return (
    <span>
      {count.toLocaleString()}
      {suffix}
    </span>
  );
}

export function Hero({ t, locale }: HeroProps) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setVisible(true), 100);
    return () => clearTimeout(timer);
  }, []);

  return (
    <section className="relative overflow-hidden min-h-[90vh] flex items-center">
      {/* Animated gradient mesh background */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-[30%] -left-[20%] w-[70%] h-[70%] rounded-full bg-md-primary-container opacity-40 blur-[120px] animate-[drift_20s_ease-in-out_infinite]" />
        <div className="absolute -bottom-[20%] -right-[15%] w-[60%] h-[60%] rounded-full bg-md-tertiary-container opacity-30 blur-[100px] animate-[drift_25s_ease-in-out_infinite_reverse]" />
        <div className="absolute top-[20%] right-[20%] w-[40%] h-[40%] rounded-full bg-md-secondary-container opacity-25 blur-[80px] animate-[drift_18s_ease-in-out_infinite_2s]" />
        {/* Subtle grid overlay */}
        <div
          className="absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage:
              "linear-gradient(var(--md-outline) 1px, transparent 1px), linear-gradient(90deg, var(--md-outline) 1px, transparent 1px)",
            backgroundSize: "60px 60px",
          }}
        />
      </div>

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-32 pb-20 lg:pt-40 lg:pb-28 w-full">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
          {/* Left: Copy */}
          <div
            className={`transition-all duration-1000 ${
              visible
                ? "opacity-100 translate-y-0"
                : "opacity-0 translate-y-8"
            }`}
          >
            {/* Social proof pill */}
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-md-surface-container rounded-full mb-8 shadow-md-1 border border-md-outline-variant">
              <div className="flex -space-x-1.5">
                {[
                  "bg-md-primary",
                  "bg-md-tertiary",
                  "bg-md-secondary",
                ].map((bg, i) => (
                  <div
                    key={i}
                    className={`w-5 h-5 rounded-full ${bg} border-2 border-md-surface-container`}
                  />
                ))}
              </div>
              <span className="text-sm font-medium text-md-on-surface-variant">
                {t.landing.socialProof}
              </span>
            </div>

            {/* Headline */}
            <h1 className="text-4xl sm:text-5xl lg:text-[3.5rem] xl:text-6xl font-bold text-md-on-surface tracking-tight leading-[1.1] text-balance">
              {t.landing.heroTitle.split(" ").map((word, i, arr) => {
                // Highlight the last two words in primary color
                if (i >= arr.length - 2) {
                  return (
                    <span key={i} className="text-md-primary">
                      {word}{" "}
                    </span>
                  );
                }
                return <span key={i}>{word} </span>;
              })}
            </h1>

            {/* Subtitle */}
            <p className="mt-6 text-lg sm:text-xl text-md-on-surface-variant max-w-lg leading-relaxed">
              {t.landing.heroSubtitle}
            </p>

            {/* CTAs */}
            <div className="mt-10 flex flex-col sm:flex-row items-start gap-4">
              <Link
                href={`/${locale}/login`}
                className="group inline-flex items-center gap-2 px-8 py-4 text-base font-medium text-md-on-primary bg-md-primary rounded-full shadow-md-2 hover:shadow-md-3 active:scale-[0.98] transition-all"
              >
                {t.landing.ctaGetStarted}
                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </Link>
              <a
                href="#features"
                className="inline-flex items-center gap-2 px-8 py-4 text-base font-medium text-md-primary border border-md-outline rounded-full hover:bg-md-primary-container/30 transition-colors"
              >
                {t.landing.ctaLearnMore}
              </a>
            </div>

            {/* Stats row */}
            <div className="mt-12 flex gap-8">
              {[
                { value: 50, suffix: "+", label: "Gyms" },
                { value: 5000, suffix: "+", label: "Members" },
                { value: 99, suffix: "%", label: "Uptime" },
              ].map((stat) => (
                <div key={stat.label}>
                  <p className="text-2xl sm:text-3xl font-bold text-md-on-surface">
                    <AnimatedNumber target={stat.value} suffix={stat.suffix} />
                  </p>
                  <p className="text-sm text-md-on-surface-variant mt-0.5">
                    {stat.label}
                  </p>
                </div>
              ))}
            </div>
          </div>

          {/* Right: Dashboard preview / floating cards */}
          <div
            className={`relative hidden lg:block transition-all duration-1000 delay-300 ${
              visible
                ? "opacity-100 translate-y-0"
                : "opacity-0 translate-y-12"
            }`}
          >
            {/* Main card - simulated dashboard */}
            <div className="relative bg-md-surface-container rounded-3xl shadow-md-3 p-6 border border-md-outline-variant">
              {/* Fake top bar */}
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-xl bg-md-primary flex items-center justify-center">
                    <Dumbbell className="w-4 h-4 text-md-on-primary" />
                  </div>
                  <span className="text-sm font-semibold text-md-on-surface">
                    Dashboard
                  </span>
                </div>
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-md-error/60" />
                  <div className="w-3 h-3 rounded-full bg-amber-400/60" />
                  <div className="w-3 h-3 rounded-full bg-green-400/60" />
                </div>
              </div>

              {/* Mini stat cards */}
              <div className="grid grid-cols-2 gap-3 mb-4">
                {[
                  {
                    icon: Users,
                    label: "Members",
                    value: "342",
                    color: "bg-md-primary-container",
                    iconColor: "text-md-primary",
                  },
                  {
                    icon: CalendarCheck,
                    label: "Today",
                    value: "89",
                    color: "bg-md-tertiary-container",
                    iconColor: "text-md-tertiary",
                  },
                  {
                    icon: TrendingUp,
                    label: "Revenue",
                    value: "Rs. 4.8L",
                    color: "bg-md-secondary-container",
                    iconColor: "text-md-secondary",
                  },
                  {
                    icon: Zap,
                    label: "Active",
                    value: "287",
                    color: "bg-green-100 dark:bg-green-900/30",
                    iconColor: "text-green-600 dark:text-green-400",
                  },
                ].map((card) => (
                  <div
                    key={card.label}
                    className="bg-md-surface-container-low rounded-2xl p-3 border border-md-outline-variant/50"
                  >
                    <div
                      className={`w-8 h-8 rounded-lg ${card.color} flex items-center justify-center mb-2`}
                    >
                      <card.icon className={`w-4 h-4 ${card.iconColor}`} />
                    </div>
                    <p className="text-xs text-md-on-surface-variant">
                      {card.label}
                    </p>
                    <p className="text-lg font-bold text-md-on-surface">
                      {card.value}
                    </p>
                  </div>
                ))}
              </div>

              {/* Fake chart area */}
              <div className="bg-md-surface-container-low rounded-2xl border border-md-outline-variant/50 p-4 h-32 flex items-end gap-1.5">
                {[40, 65, 50, 80, 60, 90, 75, 95, 70, 85, 65, 100].map(
                  (h, i) => (
                    <div
                      key={i}
                      className="flex-1 bg-md-primary/20 rounded-t-sm relative overflow-hidden"
                      style={{ height: `${h}%` }}
                    >
                      <div
                        className="absolute bottom-0 w-full bg-md-primary rounded-t-sm transition-all duration-1000"
                        style={{
                          height: visible ? "100%" : "0%",
                          transitionDelay: `${800 + i * 100}ms`,
                        }}
                      />
                    </div>
                  )
                )}
              </div>
            </div>

            {/* Floating notification card */}
            <div
              className={`absolute -left-8 top-1/3 bg-md-surface-container-lowest rounded-2xl shadow-md-2 p-4 border border-md-outline-variant w-52 transition-all duration-700 delay-[1200ms] ${
                visible
                  ? "opacity-100 -translate-x-2"
                  : "opacity-0 translate-x-4"
              }`}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center flex-shrink-0">
                  <CalendarCheck className="w-5 h-5 text-green-600 dark:text-green-400" />
                </div>
                <div>
                  <p className="text-xs font-medium text-md-on-surface">
                    New check-in
                  </p>
                  <p className="text-[11px] text-md-on-surface-variant">
                    Ram B. via QR code
                  </p>
                </div>
              </div>
            </div>

            {/* Floating streak card */}
            <div
              className={`absolute -right-4 bottom-8 bg-md-surface-container-lowest rounded-2xl shadow-md-2 p-4 border border-md-outline-variant transition-all duration-700 delay-[1500ms] ${
                visible
                  ? "opacity-100 translate-x-2"
                  : "opacity-0 -translate-x-4"
              }`}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-md-primary-container flex items-center justify-center flex-shrink-0">
                  <Zap className="w-5 h-5 text-md-primary" />
                </div>
                <div>
                  <p className="text-xs font-medium text-md-on-surface">
                    12-day streak!
                  </p>
                  <p className="text-[11px] text-md-on-surface-variant">
                    Sita K. is on fire
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <style jsx>{`
        @keyframes drift {
          0%,
          100% {
            transform: translate(0, 0);
          }
          25% {
            transform: translate(30px, -20px);
          }
          50% {
            transform: translate(-20px, 30px);
          }
          75% {
            transform: translate(20px, 20px);
          }
        }
      `}</style>
    </section>
  );
}
