"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type {
  Locale,
  OrgMember,
  OrgAttendance,
  ExpiringPackageEntry,
  ExpiredPackageEntry,
  PackageSummaryItem,
  AccountsSummary,
} from "@/types";
import {
  Users,
  UserCheck,
  CalendarCheck,
  Banknote,
  AlertTriangle,
  Clock,
  UserMinus,
  UserPlus,
  TrendingUp,
  ArrowUpRight,
  Loader2,
} from "lucide-react";

function formatNPR(amount: number): string {
  return `Rs. ${amount.toLocaleString("en-IN")}`;
}

function formatTime(isoStr: string): string {
  const d = new Date(isoStr);
  return d.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
}

function isToday(isoStr: string): boolean {
  const d = new Date(isoStr);
  const now = new Date();
  return (
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate()
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
  trend,
}: {
  icon: React.ElementType;
  label: string;
  value: string | number;
  trend?: string;
}) {
  return (
    <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card hover:shadow-card-hover transition-shadow duration-300">
      <div className="flex items-start justify-between">
        <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
          <Icon className="w-5 h-5 text-accent" />
        </div>
        {trend && (
          <span className="flex items-center gap-0.5 text-xs text-accent font-medium">
            <TrendingUp className="w-3 h-3" />
            {trend}
          </span>
        )}
      </div>
      <p className="mt-3 text-2xl font-semibold text-fg tracking-tight">
        {value}
      </p>
      <p className="mt-0.5 text-sm text-fg-muted">{label}</p>
    </div>
  );
}

function SectionHeader({
  icon: Icon,
  title,
  action,
}: {
  icon: React.ElementType;
  title: string;
  action?: { label: string; href: string };
}) {
  return (
    <div className="flex items-center justify-between mb-3">
      <h3 className="flex items-center gap-2 text-sm font-semibold text-fg">
        <Icon className="w-4 h-4 text-fg-muted" />
        {title}
      </h3>
      {action && (
        <Link
          href={action.href}
          className="flex items-center gap-1 text-xs text-accent hover:text-accent-bright transition-colors"
        >
          {action.label}
          <ArrowUpRight className="w-3 h-3" />
        </Link>
      )}
    </div>
  );
}

function MethodBadge({ method }: { method: string }) {
  const colors: Record<string, string> = {
    qr: "bg-blue-500/10 text-blue-400 border-blue-500/20",
    nfc: "bg-purple-500/10 text-purple-400 border-purple-500/20",
    manual: "bg-gray-500/10 text-gray-400 border-gray-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[method] ?? colors.manual
      }`}
    >
      {method}
    </span>
  );
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    active: "bg-accent/10 text-accent border-accent/20",
    expiring: "bg-amber-500/10 text-amber-400 border-amber-500/20",
    expired: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[status] ?? colors.active
      }`}
    >
      {status}
    </span>
  );
}

export default function DashboardHomePage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Data state
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [totalMembers, setTotalMembers] = useState(0);
  const [todayAttendance, setTodayAttendance] = useState<
    (OrgAttendance & { resolvedName?: string })[]
  >([]);
  const [expiringPackages, setExpiringPackages] = useState<
    ExpiringPackageEntry[]
  >([]);
  const [expiredPackages, setExpiredPackages] = useState<
    ExpiredPackageEntry[]
  >([]);
  const [packageSummary, setPackageSummary] = useState<PackageSummaryItem[]>(
    []
  );
  const [accountsSummary, setAccountsSummary] =
    useState<AccountsSummary | null>(null);

  const orgId = user?.org_id;

  const fetchDashboardData = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);

    const results = await Promise.allSettled([
      api.listMembers(orgId, { limit: 500 }),
      api.listAttendance(orgId, { limit: 100 }),
      api.getExpiringPackages(orgId, { days: 7, limit: 10 }),
      api.getExpiredPackages(orgId, { limit: 10 }),
      api.getPackageSummary(orgId),
      api.getAccountsSummary(orgId),
    ]);

    const [
      membersResult,
      attendanceResult,
      expiringResult,
      expiredResult,
      pkgSummaryResult,
      accountsResult,
    ] = results;

    // Track if any fetch failed
    const errors: string[] = [];

    // Members
    let membersList: OrgMember[] = [];
    if (membersResult.status === "fulfilled") {
      membersList = membersResult.value.data ?? [];
      setMembers(membersList);
      setTotalMembers(membersResult.value.total ?? membersList.length);
    } else {
      errors.push("members");
    }

    // Build member lookup map: user_id -> name
    const memberLookup = new Map<string, string>();
    for (const m of membersList) {
      memberLookup.set(m.id, m.name);
    }

    // Attendance - filter today's records and join names
    if (attendanceResult.status === "fulfilled") {
      const allAttendance = attendanceResult.value.data ?? [];
      const todayRecords = allAttendance
        .filter((r) => isToday(r.check_in_at))
        .map((r) => ({
          ...r,
          resolvedName:
            r.member_name || memberLookup.get(r.user_id) || "Unknown",
        }));
      setTodayAttendance(todayRecords);
    } else {
      errors.push("attendance");
    }

    // Expiring packages
    if (expiringResult.status === "fulfilled") {
      setExpiringPackages(expiringResult.value.data ?? []);
    } else {
      errors.push("expiring packages");
    }

    // Expired packages
    if (expiredResult.status === "fulfilled") {
      setExpiredPackages(expiredResult.value.data ?? []);
    } else {
      errors.push("expired packages");
    }

    // Package summary
    if (pkgSummaryResult.status === "fulfilled") {
      setPackageSummary(
        Array.isArray(pkgSummaryResult.value)
          ? pkgSummaryResult.value
          : []
      );
    } else {
      errors.push("package summary");
    }

    // Accounts summary
    if (accountsResult.status === "fulfilled") {
      setAccountsSummary(accountsResult.value);
    } else {
      errors.push("accounts");
    }

    if (errors.length > 0 && errors.length === results.length) {
      // All failed
      setError(t.common.error);
    } else if (errors.length > 0) {
      // Partial failure - show what we have but note the error
      setError(`Failed to load: ${errors.join(", ")}`);
    }

    setLoading(false);
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchDashboardData();
  }, [fetchDashboardData]);

  // Derived stats
  const activeMembers = members.filter((m) => m.status === "active").length;
  const monthlyRevenue = accountsSummary
    ? parseFloat(accountsSummary.total_income) || 0
    : 0;

  // Recently joined: members sorted by joined_at descending, top 5
  const recentlyJoined = [...members]
    .sort(
      (a, b) =>
        new Date(b.joined_at).getTime() - new Date(a.joined_at).getTime()
    )
    .slice(0, 5);

  // Loading state while orgId resolves
  if (!orgId) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  // Loading state while data fetches
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Error banner (partial failure) */}
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl px-4 py-3 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          icon={Users}
          label={t.dashboard.totalMembers}
          value={totalMembers}
        />
        <StatCard
          icon={UserCheck}
          label={t.dashboard.activeMembers}
          value={activeMembers}
        />
        <StatCard
          icon={CalendarCheck}
          label={t.dashboard.todayAttendance}
          value={todayAttendance.length}
        />
        <StatCard
          icon={Banknote}
          label={t.dashboard.monthlyRevenue}
          value={formatNPR(monthlyRevenue)}
        />
      </div>

      {/* Two-column layout for tables */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Today's Attendance */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={CalendarCheck}
            title={t.dashboard.todayActivity}
            action={{ label: t.dashboard.viewAll, href: `/${locale}/dashboard/attendance` }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.checkIn}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.method}
                  </th>
                </tr>
              </thead>
              <tbody>
                {todayAttendance.length === 0 ? (
                  <tr>
                    <td
                      colSpan={3}
                      className="py-6 text-center text-fg-muted text-sm"
                    >
                      No check-ins today
                    </td>
                  </tr>
                ) : (
                  todayAttendance.map((record) => (
                    <tr
                      key={record.id}
                      className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                    >
                      <td className="py-2.5 text-fg">
                        {record.resolvedName}
                      </td>
                      <td className="py-2.5 text-fg-muted font-mono text-xs">
                        {formatTime(record.check_in_at)}
                      </td>
                      <td className="py-2.5">
                        <MethodBadge method={record.method} />
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Expiring Packages */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={AlertTriangle}
            title={t.dashboard.expiringPackages}
            action={{ label: t.dashboard.viewAll, href: `/${locale}/dashboard/packages` }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.package}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.status}
                  </th>
                </tr>
              </thead>
              <tbody>
                {expiringPackages.length === 0 ? (
                  <tr>
                    <td
                      colSpan={3}
                      className="py-6 text-center text-fg-muted text-sm"
                    >
                      No expiring packages
                    </td>
                  </tr>
                ) : (
                  expiringPackages.map((entry, idx) => (
                    <tr
                      key={`${entry.member_phone}-${entry.package_name}-${idx}`}
                      className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                    >
                      <td className="py-2.5">
                        <div>
                          <p className="text-fg">{entry.member_name}</p>
                          <p className="text-xs text-fg-muted font-mono">
                            {entry.member_phone}
                          </p>
                        </div>
                      </td>
                      <td className="py-2.5 text-fg-muted">
                        {entry.package_name}
                      </td>
                      <td className="py-2.5">
                        <span className="text-amber-400 text-xs font-mono">
                          {entry.expires_in}d left
                        </span>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Expired Packages */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={Clock}
            title={t.dashboard.expiredPackages}
            action={{ label: t.dashboard.viewAll, href: `/${locale}/dashboard/packages` }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.package}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.status}
                  </th>
                </tr>
              </thead>
              <tbody>
                {expiredPackages.length === 0 ? (
                  <tr>
                    <td
                      colSpan={3}
                      className="py-6 text-center text-fg-muted text-sm"
                    >
                      No expired packages
                    </td>
                  </tr>
                ) : (
                  expiredPackages.map((entry, idx) => (
                    <tr
                      key={`${entry.member_phone}-${entry.package_name}-${idx}`}
                      className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                    >
                      <td className="py-2.5">
                        <div>
                          <p className="text-fg">{entry.member_name}</p>
                          <p className="text-xs text-fg-muted font-mono">
                            {entry.member_phone}
                          </p>
                        </div>
                      </td>
                      <td className="py-2.5 text-fg-muted">
                        {entry.package_name}
                      </td>
                      <td className="py-2.5">
                        <StatusBadge status="expired" />
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Recently Joined */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={UserPlus}
            title={t.dashboard.recentlyJoined}
            action={{ label: t.dashboard.viewAll, href: `/${locale}/dashboard/members` }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.status}
                  </th>
                  <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dashboard.joined}
                  </th>
                </tr>
              </thead>
              <tbody>
                {recentlyJoined.length === 0 ? (
                  <tr>
                    <td
                      colSpan={3}
                      className="py-6 text-center text-fg-muted text-sm"
                    >
                      No members yet
                    </td>
                  </tr>
                ) : (
                  recentlyJoined.map((member) => (
                    <tr
                      key={member.id}
                      className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                    >
                      <td className="py-2.5">
                        <div>
                          <p className="text-fg">{member.name}</p>
                          <p className="text-xs text-fg-muted font-mono">
                            {member.phone}
                          </p>
                        </div>
                      </td>
                      <td className="py-2.5">
                        <StatusBadge status={member.status} />
                      </td>
                      <td className="py-2.5 text-fg-muted font-mono text-xs">
                        {new Date(member.joined_at)
                          .toISOString()
                          .split("T")[0]}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Package Summary — full width */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        <SectionHeader
          icon={UserMinus}
          title={t.dashboard.packageSummary}
          action={{ label: t.dashboard.viewAll, href: `/${locale}/dashboard/packages` }}
        />
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-white/[0.06]">
                <th className="text-left py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                  {t.dashboard.name}
                </th>
                <th className="text-right py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                  {t.dashboard.count}
                </th>
                <th className="text-right py-2 text-xs font-mono tracking-widest text-fg-muted uppercase">
                  {t.dashboard.revenue}
                </th>
              </tr>
            </thead>
            <tbody>
              {packageSummary.length === 0 ? (
                <tr>
                  <td
                    colSpan={3}
                    className="py-6 text-center text-fg-muted text-sm"
                  >
                    No packages
                  </td>
                </tr>
              ) : (
                <>
                  {packageSummary.map((pkg) => (
                    <tr
                      key={pkg.id}
                      className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                    >
                      <td className="py-2.5 text-fg">{pkg.name}</td>
                      <td className="py-2.5 text-fg-muted text-right font-mono">
                        {pkg.member_count}
                      </td>
                      <td className="py-2.5 text-fg-muted text-right font-mono">
                        {formatNPR(pkg.member_count * parseFloat(pkg.price))}
                      </td>
                    </tr>
                  ))}
                  <tr className="border-t border-white/[0.08]">
                    <td className="py-2.5 text-fg font-semibold">Total</td>
                    <td className="py-2.5 text-fg text-right font-mono font-semibold">
                      {packageSummary.reduce(
                        (a, b) => a + b.member_count,
                        0
                      )}
                    </td>
                    <td className="py-2.5 text-fg text-right font-mono font-semibold">
                      {formatNPR(
                        packageSummary.reduce(
                          (a, b) =>
                            a + b.member_count * parseFloat(b.price),
                          0
                        )
                      )}
                    </td>
                  </tr>
                </>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
