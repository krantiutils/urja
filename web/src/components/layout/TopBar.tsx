"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Menu,
  QrCode,
  ChevronDown,
  User,
  LogOut,
  Globe,
} from "lucide-react";

interface TopBarProps {
  t: Dictionary;
  locale: Locale;
  onMenuToggle: () => void;
}

export function TopBar({ t, locale, onMenuToggle }: TopBarProps) {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [profileOpen, setProfileOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node)
      ) {
        setProfileOpen(false);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  const handleLogout = async () => {
    setProfileOpen(false);
    await logout();
    router.replace(`/${locale}/login`);
  };

  const toggleLocale = () => {
    const newLocale = locale === "en" ? "ne" : "en";
    const currentPath = window.location.pathname;
    const newPath = currentPath.replace(`/${locale}/`, `/${newLocale}/`);
    router.push(newPath);
  };

  return (
    <header className="sticky top-0 z-30 h-14 bg-bg-elevated/80 backdrop-blur-xl border-b border-white/[0.06] flex items-center justify-between px-4 lg:px-6">
      {/* Left: menu toggle (mobile) + page title area */}
      <div className="flex items-center gap-3">
        <button
          onClick={onMenuToggle}
          className="lg:hidden p-2 rounded-lg hover:bg-surface transition-colors"
          aria-label="Toggle menu"
        >
          <Menu className="w-5 h-5 text-fg-muted" />
        </button>
      </div>

      {/* Right: actions */}
      <div className="flex items-center gap-2">
        {/* Language toggle */}
        <button
          onClick={toggleLocale}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm text-fg-muted hover:bg-surface hover:text-fg transition-colors"
          title={locale === "en" ? "नेपालीमा हेर्नुहोस्" : "View in English"}
        >
          <Globe className="w-4 h-4" />
          <span className="text-xs font-mono uppercase">
            {locale === "en" ? "NE" : "EN"}
          </span>
        </button>

        {/* QR Code button */}
        <button
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm text-fg-muted hover:bg-surface hover:text-fg transition-colors"
          title={t.topbar.orgQr}
        >
          <QrCode className="w-4 h-4" />
          <span className="hidden sm:inline">{t.topbar.orgQr}</span>
        </button>

        {/* Profile dropdown */}
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setProfileOpen(!profileOpen)}
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-surface transition-colors"
          >
            <div className="w-7 h-7 rounded-full bg-accent/20 flex items-center justify-center">
              <User className="w-3.5 h-3.5 text-accent" />
            </div>
            <span className="hidden sm:inline text-sm text-fg-muted">
              {user?.phone ?? ""}
            </span>
            <ChevronDown
              className={`w-3 h-3 text-fg-muted transition-transform duration-200 ${
                profileOpen ? "rotate-180" : ""
              }`}
            />
          </button>

          {profileOpen && (
            <div className="absolute right-0 top-full mt-1 w-48 bg-bg-elevated border border-white/[0.06] rounded-xl shadow-card overflow-hidden animate-fade-in">
              <div className="px-3 py-2 border-b border-white/[0.06]">
                <p className="text-sm text-fg font-medium truncate">
                  {user?.phone}
                </p>
                <p className="text-xs text-fg-muted capitalize">
                  {user?.role}
                </p>
              </div>
              <div className="py-1">
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-400 hover:bg-surface transition-colors"
                >
                  <LogOut className="w-4 h-4" />
                  {t.topbar.logout}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
