"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Home,
  CalendarCheck,
  Package,
  Heart,
  Dumbbell,
  Apple,
  TrendingUp,
  BookOpen,
  ClipboardList,
  X,
} from "lucide-react";

interface NavItem {
  label: string;
  href: string;
  icon: React.ReactNode;
}

function buildNavItems(t: Dictionary, locale: Locale): NavItem[] {
  const base = `/${locale}/member`;
  return [
    { label: t.memberNav.myDashboard, href: base, icon: <Home className="w-4 h-4" /> },
    { label: t.memberNav.myAttendance, href: `${base}/attendance`, icon: <CalendarCheck className="w-4 h-4" /> },
    { label: t.memberNav.myPackages, href: `${base}/packages`, icon: <Package className="w-4 h-4" /> },
    { label: t.memberNav.myHealth, href: `${base}/health`, icon: <Heart className="w-4 h-4" /> },
    { label: t.memberNav.myWorkouts, href: `${base}/workouts`, icon: <Dumbbell className="w-4 h-4" /> },
    { label: t.memberNav.myNutrition, href: `${base}/nutrition`, icon: <Apple className="w-4 h-4" /> },
    { label: t.memberNav.myExercises, href: `${base}/exercises`, icon: <BookOpen className="w-4 h-4" /> },
    { label: t.memberNav.myPrograms, href: `${base}/programs`, icon: <ClipboardList className="w-4 h-4" /> },
    { label: t.memberNav.myProgress, href: `${base}/progress`, icon: <TrendingUp className="w-4 h-4" /> },
  ];
}

interface MemberSidebarProps {
  t: Dictionary;
  locale: Locale;
  isOpen: boolean;
  onClose: () => void;
}

export function MemberSidebar({ t, locale, isOpen, onClose }: MemberSidebarProps) {
  const pathname = usePathname();
  const items = buildNavItems(t, locale);

  const isActive = (href: string) => {
    if (href === `/${locale}/member`) {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={`fixed top-0 left-0 h-full w-64 bg-bg-elevated border-r border-white/[0.06] z-50 flex flex-col transition-transform duration-300 lg:translate-x-0 ${
          isOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-white/[0.06]">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-xl bg-accent/20 flex items-center justify-center">
              <Dumbbell className="w-4 h-4 text-accent" />
            </div>
            <div>
              <h2 className="text-sm font-semibold text-fg truncate">
                {t.common.appName}
              </h2>
              <p className="text-xs text-fg-muted truncate">Member Portal</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="lg:hidden p-1 rounded-lg hover:bg-surface transition-colors"
          >
            <X className="w-4 h-4 text-fg-muted" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto py-3 px-3 space-y-0.5">
          {items.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={onClose}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-200 ${
                  active
                    ? "bg-accent/10 text-accent font-medium"
                    : "text-fg-muted hover:bg-surface hover:text-fg"
                }`}
              >
                <span className={active ? "text-accent" : ""}>
                  {item.icon}
                </span>
                {item.label}
              </Link>
            );
          })}
        </nav>
      </aside>
    </>
  );
}
