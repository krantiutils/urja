"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Home,
  Users,
  Package,
  CreditCard,
  CalendarCheck,
  UserCog,
  BookOpen,
  MessageSquare,
  Mail,
  MessageCircle,
  NfcIcon,
  Settings,
  ChevronDown,
  Dumbbell,
  X,
} from "lucide-react";
import { useState } from "react";

interface NavSection {
  label: string;
  items: NavItem[];
}

interface NavItem {
  label: string;
  href: string;
  icon: React.ReactNode;
}

function buildNavSections(t: Dictionary, locale: Locale): NavSection[] {
  const base = `/${locale}/dashboard`;
  return [
    {
      label: t.nav.manageMembers,
      items: [
        { label: t.nav.dashboard, href: base, icon: <Home className="w-4 h-4" /> },
        { label: t.nav.members, href: `${base}/members`, icon: <Users className="w-4 h-4" /> },
        { label: t.nav.packages, href: `${base}/packages`, icon: <Package className="w-4 h-4" /> },
        { label: t.nav.duePayments, href: `${base}/due-payments`, icon: <CreditCard className="w-4 h-4" /> },
        { label: t.nav.attendance, href: `${base}/attendance`, icon: <CalendarCheck className="w-4 h-4" /> },
      ],
    },
    {
      label: t.nav.runOperations,
      items: [
        { label: t.nav.staff, href: `${base}/staff`, icon: <UserCog className="w-4 h-4" /> },
        { label: t.nav.accounts, href: `${base}/accounts`, icon: <BookOpen className="w-4 h-4" /> },
      ],
    },
    {
      label: t.nav.engageMembers,
      items: [
        { label: t.nav.stories, href: `${base}/stories`, icon: <MessageSquare className="w-4 h-4" /> },
        { label: t.nav.sms, href: `${base}/sms`, icon: <Mail className="w-4 h-4" /> },
        { label: t.nav.feedbacks, href: `${base}/feedbacks`, icon: <MessageCircle className="w-4 h-4" /> },
      ],
    },
    {
      label: t.nav.accessControl,
      items: [
        { label: t.nav.nfcAccess, href: `${base}/nfc`, icon: <NfcIcon className="w-4 h-4" /> },
      ],
    },
  ];
}

interface SidebarProps {
  t: Dictionary;
  locale: Locale;
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({ t, locale, isOpen, onClose }: SidebarProps) {
  const pathname = usePathname();
  const sections = buildNavSections(t, locale);
  const [collapsed, setCollapsed] = useState<Record<string, boolean>>({});

  const toggleSection = (label: string) => {
    setCollapsed((prev) => ({ ...prev, [label]: !prev[label] }));
  };

  const isActive = (href: string) => {
    if (href === `/${locale}/dashboard`) {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/40 backdrop-blur-sm z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={`fixed top-0 left-0 h-full w-64 bg-md-surface-container border-r border-md-outline-variant z-50 flex flex-col transition-transform duration-300 lg:translate-x-0 ${
          isOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-md-outline-variant">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-xl bg-md-primary flex items-center justify-center">
              <Dumbbell className="w-4 h-4 text-md-on-primary" />
            </div>
            <div>
              <h2 className="text-sm font-semibold text-md-on-surface truncate">
                {t.common.appName}
              </h2>
              <p className="text-xs text-md-on-surface-variant truncate">Gym Name</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="lg:hidden p-1 rounded-full hover:bg-md-surface-container-high transition-colors"
          >
            <X className="w-4 h-4 text-md-on-surface-variant" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto py-3 px-3 space-y-1">
          {sections.map((section) => (
            <div key={section.label} className="mb-2">
              <button
                onClick={() => toggleSection(section.label)}
                className="flex items-center justify-between w-full px-2 py-1.5 text-[11px] font-medium tracking-widest uppercase text-md-on-surface-variant hover:text-md-on-surface transition-colors"
              >
                {section.label}
                <ChevronDown
                  className={`w-3 h-3 transition-transform duration-200 ${
                    collapsed[section.label] ? "-rotate-90" : ""
                  }`}
                />
              </button>

              {!collapsed[section.label] && (
                <div className="mt-0.5 space-y-0.5">
                  {section.items.map((item) => {
                    const active = isActive(item.href);
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={onClose}
                        className={`flex items-center gap-3 px-3 py-2 rounded-full text-sm transition-all duration-200 ${
                          active
                            ? "bg-md-primary-container text-md-on-primary-container font-medium"
                            : "text-md-on-surface-variant hover:bg-md-surface-container-high hover:text-md-on-surface"
                        }`}
                      >
                        <span className={active ? "text-md-primary" : ""}>
                          {item.icon}
                        </span>
                        {item.label}
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          ))}
        </nav>

        {/* Settings at bottom */}
        <div className="p-3 border-t border-md-outline-variant">
          <Link
            href={`/${locale}/dashboard/settings`}
            className={`flex items-center gap-3 px-3 py-2 rounded-full text-sm transition-all duration-200 ${
              isActive(`/${locale}/dashboard/settings`)
                ? "bg-md-primary-container text-md-on-primary-container font-medium"
                : "text-md-on-surface-variant hover:bg-md-surface-container-high hover:text-md-on-surface"
            }`}
          >
            <Settings className="w-4 h-4" />
            {t.nav.settings}
          </Link>
        </div>
      </aside>
    </>
  );
}
