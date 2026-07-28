"use client";

import { useState } from "react";
import { Sidebar } from "@/components/layout/Sidebar";
import { TopBar } from "@/components/layout/TopBar";
import { AuthGuard } from "@/components/layout/AuthGuard";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";

export default function DashboardLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <AuthGuard locale={locale} requireOperator>
      <div className="min-h-screen bg-bg-base">
        <Sidebar
          t={t}
          locale={locale}
          isOpen={sidebarOpen}
          onClose={() => setSidebarOpen(false)}
          orgName={user?.org_name}
        />

        {/* Main content area offset by sidebar width on desktop */}
        <div className="lg:ml-64">
          <TopBar
            t={t}
            locale={locale}
            onMenuToggle={() => setSidebarOpen(!sidebarOpen)}
          />
          <main className="p-4 lg:p-6">{children}</main>
        </div>
      </div>
    </AuthGuard>
  );
}
