"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, MemberPackage } from "@/types";
import { Loader2, Package } from "lucide-react";

const statusColors: Record<string, string> = {
  active: "bg-emerald-500/10 text-emerald-400 border-emerald-500/20",
  expired: "bg-red-500/10 text-red-400 border-red-500/20",
  cancelled: "bg-white/[0.06] text-fg-muted border-white/[0.06]",
  pending: "bg-amber-500/10 text-amber-400 border-amber-500/20",
};

export default function MemberPackagesPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await api.getMyPackages();
        setPackages(res.data ?? []);
      } catch (err) {
        console.error("Failed to load packages:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">
        {t.memberPages.packages}
      </h1>

      {packages.length === 0 ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
          <div className="flex flex-col items-center justify-center py-16 gap-3">
            <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
              <Package className="w-6 h-6 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">{t.memberPages.noPackages}</p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {packages.map((pkg) => {
            const daysRemaining = Math.max(
              0,
              Math.ceil(
                (new Date(pkg.end_date).getTime() - Date.now()) /
                  (1000 * 60 * 60 * 24)
              )
            );
            return (
              <div
                key={pkg.id}
                className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5"
              >
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="text-base font-semibold text-fg">
                      {pkg.package_name ?? "Package"}
                    </h3>
                    {pkg.org_name && (
                      <p className="text-xs text-fg-muted mt-0.5">
                        {pkg.org_name}
                      </p>
                    )}
                  </div>
                  <span
                    className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
                      statusColors[pkg.status] ?? statusColors.pending
                    }`}
                  >
                    {pkg.status}
                  </span>
                </div>

                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-fg-muted">Start</span>
                    <span className="text-fg">
                      {new Date(pkg.start_date).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-fg-muted">End</span>
                    <span className="text-fg">
                      {new Date(pkg.end_date).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-fg-muted">Amount Paid</span>
                    <span className="text-fg">NPR {pkg.amount_paid}</span>
                  </div>
                  {pkg.status === "active" && (
                    <div className="flex justify-between">
                      <span className="text-fg-muted">
                        {t.memberPages.daysRemaining}
                      </span>
                      <span className="text-accent font-medium">
                        {daysRemaining} {t.memberPages.days}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
