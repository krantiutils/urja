"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import { Dumbbell, Menu, X, Sun, Moon } from "lucide-react";
import { useTheme } from "@/lib/theme";

interface NavbarProps {
  t: Dictionary;
  locale: Locale;
}

export function Navbar({ t, locale }: NavbarProps) {
  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const { resolved, setTheme } = useTheme();

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  const toggleTheme = () => {
    setTheme(resolved === "dark" ? "light" : "dark");
  };

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled
          ? "bg-md-surface/95 backdrop-blur-xl shadow-md-1"
          : "bg-transparent"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link
            href={`/${locale}`}
            className="flex items-center gap-2.5"
          >
            <div className="w-9 h-9 rounded-xl bg-md-primary flex items-center justify-center">
              <Dumbbell className="w-5 h-5 text-md-on-primary" />
            </div>
            <span className="text-xl font-bold text-md-on-surface">
              {t.common.appName}
            </span>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden md:flex items-center gap-1">
            <a
              href="#features"
              className="px-4 py-2 text-sm font-medium text-md-on-surface-variant hover:text-md-on-surface rounded-full hover:bg-md-surface-container transition-colors"
            >
              {t.landing.featuresTag}
            </a>
            <a
              href="#pricing"
              className="px-4 py-2 text-sm font-medium text-md-on-surface-variant hover:text-md-on-surface rounded-full hover:bg-md-surface-container transition-colors"
            >
              {t.landing.pricingTag}
            </a>
            <a
              href="#faq"
              className="px-4 py-2 text-sm font-medium text-md-on-surface-variant hover:text-md-on-surface rounded-full hover:bg-md-surface-container transition-colors"
            >
              {t.landing.faqTag}
            </a>
          </nav>

          {/* Desktop CTAs */}
          <div className="hidden md:flex items-center gap-3">
            <button
              onClick={toggleTheme}
              className="p-2 rounded-full hover:bg-md-surface-container text-md-on-surface-variant hover:text-md-on-surface transition-colors"
              aria-label="Toggle theme"
            >
              {resolved === "dark" ? (
                <Sun className="w-5 h-5" />
              ) : (
                <Moon className="w-5 h-5" />
              )}
            </button>
            <Link
              href={`/${locale}/login`}
              className="px-5 py-2.5 text-sm font-medium text-md-primary border border-md-outline rounded-full hover:bg-md-primary-container/30 transition-colors"
            >
              {t.auth.signIn}
            </Link>
            <Link
              href={`/${locale}/login`}
              className="px-5 py-2.5 text-sm font-medium text-md-on-primary bg-md-primary rounded-full hover:bg-md-primary/90 shadow-md-1 transition-all"
            >
              {t.landing.ctaGetStarted}
            </Link>
          </div>

          {/* Mobile menu button */}
          <div className="flex md:hidden items-center gap-2">
            <button
              onClick={toggleTheme}
              className="p-2 rounded-full hover:bg-md-surface-container text-md-on-surface-variant transition-colors"
              aria-label="Toggle theme"
            >
              {resolved === "dark" ? (
                <Sun className="w-5 h-5" />
              ) : (
                <Moon className="w-5 h-5" />
              )}
            </button>
            <button
              onClick={() => setMobileOpen(!mobileOpen)}
              className="p-2 rounded-full hover:bg-md-surface-container transition-colors"
            >
              {mobileOpen ? (
                <X className="w-5 h-5 text-md-on-surface" />
              ) : (
                <Menu className="w-5 h-5 text-md-on-surface" />
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileOpen && (
        <div className="md:hidden bg-md-surface border-t border-md-outline-variant animate-fade-in">
          <div className="px-4 py-4 space-y-2">
            <a
              href="#features"
              onClick={() => setMobileOpen(false)}
              className="block px-4 py-3 text-sm font-medium text-md-on-surface-variant hover:bg-md-surface-container rounded-2xl transition-colors"
            >
              {t.landing.featuresTag}
            </a>
            <a
              href="#pricing"
              onClick={() => setMobileOpen(false)}
              className="block px-4 py-3 text-sm font-medium text-md-on-surface-variant hover:bg-md-surface-container rounded-2xl transition-colors"
            >
              {t.landing.pricingTag}
            </a>
            <a
              href="#faq"
              onClick={() => setMobileOpen(false)}
              className="block px-4 py-3 text-sm font-medium text-md-on-surface-variant hover:bg-md-surface-container rounded-2xl transition-colors"
            >
              {t.landing.faqTag}
            </a>
            <div className="pt-2 space-y-2">
              <Link
                href={`/${locale}/login`}
                className="block w-full text-center px-5 py-3 text-sm font-medium text-md-primary border border-md-outline rounded-full hover:bg-md-primary-container/30 transition-colors"
              >
                {t.auth.signIn}
              </Link>
              <Link
                href={`/${locale}/login`}
                className="block w-full text-center px-5 py-3 text-sm font-medium text-md-on-primary bg-md-primary rounded-full shadow-md-1 transition-all"
              >
                {t.landing.ctaGetStarted}
              </Link>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}
