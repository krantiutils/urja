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
    <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1 hover:shadow-md-2 transition-shadow duration-300">
      <div className="flex items-start justify-between">
        <div className="w-10 h-10 rounded-xl bg-md-primary-container flex items-center justify-center">
          <Icon className="w-5 h-5 text-md-primary" />
        </div>
        {trend && (
          <span className="flex items-center gap-0.5 text-xs text-md-primary font-medium">
            <TrendingUp className="w-3 h-3" />
            {trend}
          </span>
        )}
      </div>
      <p className="mt-3 text-2xl font-semibold text-md-on-surface tracking-tight">
        {value}
      </p>
      <p className="mt-0.5 text-sm text-md-on-surface-variant">{label}</p>
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
      <h3 className="flex items-center gap-2 text-sm font-semibold text-md-on-surface">
        <Icon className="w-4 h-4 text-md-on-surface-variant" />
        {title}
      </h3>
      {action && (
        <button className="flex items-center gap-1 text-xs text-md-primary hover:text-md-primary/80 transition-colors">
          {action.label}
          <ArrowUpRight className="w-3 h-3" />
        </button>
      )}
    </div>
  );
}

function MethodBadge({ method }: { method: string }) {
  const colors: Record<string, string> = {
    qr: "bg-blue-100 text-blue-700",
    nfc: "bg-purple-100 text-purple-700",
    manual: "bg-gray-100 text-gray-600",
  };
  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-medium uppercase ${
        colors[method] ?? colors.manual
      }`}
    >
      {method}
    </span>
  );
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    active: "bg-green-100 text-green-700",
    expiring: "bg-amber-100 text-amber-700",
    expired: "bg-red-100 text-red-700",
  };
  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-medium uppercase ${
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
        <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1">
          <SectionHeader
            icon={CalendarCheck}
            title={t.dashboard.todayActivity}
            action={{ label: t.dashboard.viewAll }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.checkIn}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.method}
                  </th>
                </tr>
              </thead>
              <tbody>
                {MOCK_TODAY_ATTENDANCE.map((record) => (
                  <tr
                    key={record.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="py-2.5 text-md-on-surface">{record.memberName}</td>
                    <td className="py-2.5 text-md-on-surface-variant text-xs">
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
        <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1">
          <SectionHeader
            icon={AlertTriangle}
            title={t.dashboard.expiringPackages}
            action={{ label: t.dashboard.viewAll }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.package}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.status}
                  </th>
                </tr>
              </thead>
              <tbody>
                {MOCK_EXPIRING.map((member) => (
                  <tr
                    key={member.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="py-2.5">
                      <div>
                        <p className="text-md-on-surface">{member.name}</p>
                        <p className="text-xs text-md-on-surface-variant">
                          {member.phone}
                        </p>
                      </div>
                    </td>
                    <td className="py-2.5 text-md-on-surface-variant">{member.package}</td>
                    <td className="py-2.5">
                      <span className="text-amber-600 text-xs font-medium">
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
        <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1">
          <SectionHeader
            icon={Clock}
            title={t.dashboard.expiredPackages}
            action={{ label: t.dashboard.viewAll }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.package}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.status}
                  </th>
                </tr>
              </thead>
              <tbody>
                {MOCK_EXPIRED.map((member) => (
                  <tr
                    key={member.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="py-2.5">
                      <div>
                        <p className="text-md-on-surface">{member.name}</p>
                        <p className="text-xs text-md-on-surface-variant">
                          {member.phone}
                        </p>
                      </div>
                    </td>
                    <td className="py-2.5 text-md-on-surface-variant">{member.package}</td>
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
        <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1">
          <SectionHeader
            icon={UserPlus}
            title={t.dashboard.recentlyJoined}
            action={{ label: t.dashboard.viewAll }}
          />
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.member}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.package}
                  </th>
                  <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.dashboard.joined}
                  </th>
                </tr>
              </thead>
              <tbody>
                {MOCK_RECENTLY_JOINED.map((member) => (
                  <tr
                    key={member.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="py-2.5">
                      <div>
                        <p className="text-md-on-surface">{member.name}</p>
                        <p className="text-xs text-md-on-surface-variant">
                          {member.phone}
                        </p>
                      </div>
                    </td>
                    <td className="py-2.5 text-md-on-surface-variant">{member.package}</td>
                    <td className="py-2.5 text-md-on-surface-variant text-xs">
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
      <div className="bg-md-surface-container rounded-3xl p-5 shadow-md-1">
        <SectionHeader
          icon={UserMinus}
          title={t.dashboard.packageSummary}
          action={{ label: t.dashboard.viewAll }}
        />
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-md-outline-variant">
                <th className="text-left py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                  {t.dashboard.name}
                </th>
                <th className="text-right py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                  {t.dashboard.count}
                </th>
                <th className="text-right py-2 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                  {t.dashboard.revenue}
                </th>
              </tr>
            </thead>
            <tbody>
              {MOCK_PACKAGE_SUMMARY.map((pkg) => (
                <tr
                  key={pkg.name}
                  className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                >
                  <td className="py-2.5 text-md-on-surface">{pkg.name}</td>
                  <td className="py-2.5 text-md-on-surface-variant text-right">
                    {pkg.count}
                  </td>
                  <td className="py-2.5 text-md-on-surface-variant text-right">
                    {formatNPR(pkg.revenue)}
                  </td>
                </tr>
              ))}
              <tr className="border-t border-md-outline-variant">
                <td className="py-2.5 text-md-on-surface font-semibold">Total</td>
                <td className="py-2.5 text-md-on-surface text-right font-semibold">
                  {MOCK_PACKAGE_SUMMARY.reduce((a, b) => a + b.count, 0)}
                </td>
                <td className="py-2.5 text-md-on-surface text-right font-semibold">
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
