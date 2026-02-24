"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { SuperAdminSidebar } from "@/components/layout/SuperAdminSidebar";
import { useAuth } from "@/lib/auth";
import type { Locale } from "@/types";
import { Loader2, Menu } from "lucide-react";

export default function SuperAdminLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const router = useRouter();
  const { user, isLoading, isAuthenticated } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    if (!isLoading && (!isAuthenticated || !user?.is_super_admin)) {
      router.replace(`/${locale}/login`);
    }
  }, [isLoading, isAuthenticated, user, router, locale]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-bg-base">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  if (!isAuthenticated || !user?.is_super_admin) {
    return null;
  }

  return (
    <div className="min-h-screen bg-bg-base">
      <SuperAdminSidebar
        locale={locale}
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
      />

      <div className="lg:ml-64">
        {/* Simple top bar */}
        <header className="sticky top-0 z-30 h-14 bg-bg-elevated/80 backdrop-blur-xl border-b border-white/[0.06] flex items-center px-4 lg:px-6">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="lg:hidden p-2 rounded-lg hover:bg-surface transition-colors"
            aria-label="Toggle menu"
          >
            <Menu className="w-5 h-5 text-fg-muted" />
          </button>
          <span className="ml-2 lg:ml-0 text-sm text-fg-muted">
            {user.phone}
          </span>
        </header>
        <main className="p-4 lg:p-6">{children}</main>
      </div>
    </div>
  );
}
