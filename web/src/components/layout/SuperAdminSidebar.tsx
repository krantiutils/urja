"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import type { Locale } from "@/types";
import {
  Building2,
  LogOut,
  Dumbbell,
  X,
} from "lucide-react";

interface SuperAdminSidebarProps {
  locale: Locale;
  isOpen: boolean;
  onClose: () => void;
}

export function SuperAdminSidebar({ locale, isOpen, onClose }: SuperAdminSidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { logout } = useAuth();

  const base = `/${locale}/super-admin`;

  const navItems = [
    { label: "Organizations", href: base, icon: <Building2 className="w-4 h-4" /> },
  ];

  const isActive = (href: string) => pathname === href;

  const handleLogout = async () => {
    await logout();
    router.replace(`/${locale}/login`);
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
              <h2 className="text-sm font-semibold text-fg">Urja</h2>
              <p className="text-xs text-fg-muted">Super Admin</p>
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
        <nav className="flex-1 overflow-y-auto py-3 px-3">
          <div className="space-y-0.5">
            {navItems.map((item) => {
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
          </div>
        </nav>

        {/* Logout at bottom */}
        <div className="p-3 border-t border-white/[0.06]">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-fg-muted hover:bg-surface hover:text-fg transition-all duration-200 w-full"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      </aside>
    </>
  );
}
