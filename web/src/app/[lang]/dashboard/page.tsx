"use client";

import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
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
} from "lucide-react";

// Mock data — will be replaced with API calls once backend endpoints exist
const MOCK_STATS = {
  totalMembers: 342,
  activeMembers: 287,
  todayAttendance: 89,
  monthlyRevenue: 485000,
};

const MOCK_EXPIRING = [
  { id: "1", name: "Ram Shrestha", phone: "9841234567", package: "3 Month Premium", daysLeft: 3, status: "expiring" as const },
  { id: "2", name: "Sita Maharjan", phone: "9851234567", package: "1 Month Basic", daysLeft: 1, status: "expiring" as const },
  { id: "3", name: "Hari Tamang", phone: "9861234567", package: "6 Month Standard", daysLeft: 5, status: "expiring" as const },
];

const MOCK_EXPIRED = [
  { id: "4", name: "Bikash Gurung", phone: "9801234567", package: "1 Month Basic", expiredDays: 2, status: "expired" as const },
  { id: "5", name: "Anita Rai", phone: "9811234567", package: "3 Month Premium", expiredDays: 7, status: "expired" as const },
];

const MOCK_TODAY_ATTENDANCE = [
  { id: "1", memberName: "Ram Shrestha", checkInTime: "06:15 AM", method: "qr" as const },
  { id: "2", memberName: "Sita Maharjan", checkInTime: "06:42 AM", method: "nfc" as const },
  { id: "3", memberName: "Hari Tamang", checkInTime: "07:01 AM", method: "manual" as const },
  { id: "4", memberName: "Deepa Thapa", checkInTime: "07:23 AM", method: "qr" as const },
  { id: "5", memberName: "Sujan KC", checkInTime: "07:45 AM", method: "nfc" as const },
];

const MOCK_RECENTLY_JOINED = [
  { id: "1", name: "Priya Adhikari", phone: "9871234567", package: "3 Month Premium", joinedAt: "2026-02-13" },
  { id: "2", name: "Rajesh Bhandari", phone: "9881234567", package: "1 Month Basic", joinedAt: "2026-02-12" },
  { id: "3", name: "Mina Poudel", phone: "9891234567", package: "6 Month Standard", joinedAt: "2026-02-11" },
];

const MOCK_PACKAGE_SUMMARY = [
  { name: "1 Month Basic", count: 98, revenue: 78400 },
  { name: "3 Month Premium", count: 124, revenue: 297600 },
  { name: "6 Month Standard", count: 65, revenue: 195000 },
  { name: "1 Year Gold", count: 42, revenue: 252000 },
];

function formatNPR(amount: number): string {
  return `Rs. ${amount.toLocaleString("en-IN")}`;
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
  action?: { label: string; onClick?: () => void };
}) {
  return (
    <div className="flex items-center justify-between mb-3">
      <h3 className="flex items-center gap-2 text-sm font-semibold text-fg">
        <Icon className="w-4 h-4 text-fg-muted" />
        {title}
      </h3>
      {action && (
        <button className="flex items-center gap-1 text-xs text-accent hover:text-accent-bright transition-colors">
          {action.label}
          <ArrowUpRight className="w-3 h-3" />
        </button>
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

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          icon={Users}
          label={t.dashboard.totalMembers}
          value={MOCK_STATS.totalMembers}
          trend="+12"
        />
        <StatCard
          icon={UserCheck}
          label={t.dashboard.activeMembers}
          value={MOCK_STATS.activeMembers}
        />
        <StatCard
          icon={CalendarCheck}
          label={t.dashboard.todayAttendance}
          value={MOCK_STATS.todayAttendance}
          trend="+5%"
        />
        <StatCard
          icon={Banknote}
          label={t.dashboard.monthlyRevenue}
          value={formatNPR(MOCK_STATS.monthlyRevenue)}
        />
      </div>

      {/* Two-column layout for tables */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Today's Attendance */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={CalendarCheck}
            title={t.dashboard.todayActivity}
            action={{ label: t.dashboard.viewAll }}
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
                {MOCK_TODAY_ATTENDANCE.map((record) => (
                  <tr
                    key={record.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="py-2.5 text-fg">{record.memberName}</td>
                    <td className="py-2.5 text-fg-muted font-mono text-xs">
                      {record.checkInTime}
                    </td>
                    <td className="py-2.5">
                      <MethodBadge method={record.method} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Expiring Packages */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={AlertTriangle}
            title={t.dashboard.expiringPackages}
            action={{ label: t.dashboard.viewAll }}
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
                {MOCK_EXPIRING.map((member) => (
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
                    <td className="py-2.5 text-fg-muted">{member.package}</td>
                    <td className="py-2.5">
                      <span className="text-amber-400 text-xs font-mono">
                        {member.daysLeft}d left
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Expired Packages */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={Clock}
            title={t.dashboard.expiredPackages}
            action={{ label: t.dashboard.viewAll }}
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
                {MOCK_EXPIRED.map((member) => (
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
                    <td className="py-2.5 text-fg-muted">{member.package}</td>
                    <td className="py-2.5">
                      <StatusBadge status="expired" />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Recently Joined */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <SectionHeader
            icon={UserPlus}
            title={t.dashboard.recentlyJoined}
            action={{ label: t.dashboard.viewAll }}
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
                    {t.dashboard.joined}
                  </th>
                </tr>
              </thead>
              <tbody>
                {MOCK_RECENTLY_JOINED.map((member) => (
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
                    <td className="py-2.5 text-fg-muted">{member.package}</td>
                    <td className="py-2.5 text-fg-muted font-mono text-xs">
                      {member.joinedAt}
                    </td>
                  </tr>
                ))}
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
          action={{ label: t.dashboard.viewAll }}
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
              {MOCK_PACKAGE_SUMMARY.map((pkg) => (
                <tr
                  key={pkg.name}
                  className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                >
                  <td className="py-2.5 text-fg">{pkg.name}</td>
                  <td className="py-2.5 text-fg-muted text-right font-mono">
                    {pkg.count}
                  </td>
                  <td className="py-2.5 text-fg-muted text-right font-mono">
                    {formatNPR(pkg.revenue)}
                  </td>
                </tr>
              ))}
              <tr className="border-t border-white/[0.08]">
                <td className="py-2.5 text-fg font-semibold">Total</td>
                <td className="py-2.5 text-fg text-right font-mono font-semibold">
                  {MOCK_PACKAGE_SUMMARY.reduce((a, b) => a + b.count, 0)}
                </td>
                <td className="py-2.5 text-fg text-right font-mono font-semibold">
                  {formatNPR(
                    MOCK_PACKAGE_SUMMARY.reduce((a, b) => a + b.revenue, 0)
                  )}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
