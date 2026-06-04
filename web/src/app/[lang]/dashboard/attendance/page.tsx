"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, OrgAttendance } from "@/types";
import {
  CalendarCheck,
  UserPlus,
  X,
  Loader2,
  Calendar,
} from "lucide-react";

type MethodFilter = "all" | "qr" | "nfc" | "manual";

function MethodBadge({ method }: { method: string }) {
  const colors: Record<string, string> = {
    qr: "bg-blue-100 text-blue-700",
    nfc: "bg-purple-100 text-purple-700",
    manual: "bg-gray-100 text-gray-600",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase ${
        colors[method] ?? colors.manual
      }`}
    >
      {method}
    </span>
  );
}

function formatTime(isoStr: string): string {
  const d = new Date(isoStr);
  return d.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
}

function formatDate(isoStr: string): string {
  return new Date(isoStr).toLocaleDateString();
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

export default function AttendancePage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [records, setRecords] = useState<OrgAttendance[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [methodFilter, setMethodFilter] = useState<MethodFilter>("all");
  const [selectedDate, setSelectedDate] = useState<string>(
    new Date().toISOString().split("T")[0]
  );

  // Manual check-in modal
  const [showCheckIn, setShowCheckIn] = useState(false);
  const [checkInMemberId, setCheckInMemberId] = useState("");
  const [checkInLoading, setCheckInLoading] = useState(false);
  const [checkInError, setCheckInError] = useState<string | null>(null);

  const orgId = user?.org_id;

  const fetchAttendance = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listAttendance(orgId, { limit: 100 });
      setRecords(res.data ?? []);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchAttendance();
  }, [fetchAttendance]);

  // Filter by method and date
  const filtered = records.filter((r) => {
    const matchesMethod =
      methodFilter === "all" || r.method === methodFilter;
    const recordDate = new Date(r.check_in_at)
      .toISOString()
      .split("T")[0];
    const matchesDate = recordDate === selectedDate;
    return matchesMethod && matchesDate;
  });

  const todayCount = records.filter((r) => isToday(r.check_in_at)).length;

  const handleManualCheckIn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setCheckInLoading(true);
    setCheckInError(null);
    try {
      await api.manualCheckIn(orgId, { member_id: checkInMemberId });
      setShowCheckIn(false);
      setCheckInMemberId("");
      fetchAttendance();
    } catch (err) {
      setCheckInError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setCheckInLoading(false);
    }
  };

  const isSelectedToday =
    selectedDate === new Date().toISOString().split("T")[0];

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-md-on-surface flex items-center gap-2">
          <CalendarCheck className="w-5 h-5 text-md-primary" />
          {t.attendance.title}
        </h1>
        <button
          onClick={() => setShowCheckIn(true)}
          className="flex items-center gap-2 px-4 py-2 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors shadow-md-1"
        >
          <UserPlus className="w-4 h-4" />
          {t.attendance.manualCheckIn}
        </button>
      </div>

      {/* Stats Card */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="bg-md-surface-container border border-md-outline-variant rounded-2xl p-5 shadow-md-1">
          <div className="w-10 h-10 rounded-xl bg-md-primary-container flex items-center justify-center mb-3">
            <CalendarCheck className="w-5 h-5 text-md-primary" />
          </div>
          <p className="text-2xl font-semibold text-md-on-surface tracking-tight">
            {todayCount}
          </p>
          <p className="mt-0.5 text-sm text-md-on-surface-variant">
            {t.attendance.todayCheckIns}
          </p>
        </div>
        <div className="bg-md-surface-container border border-md-outline-variant rounded-2xl p-5 shadow-md-1">
          <div className="w-10 h-10 rounded-xl bg-md-primary-container flex items-center justify-center mb-3">
            <Calendar className="w-5 h-5 text-md-primary" />
          </div>
          <p className="text-2xl font-semibold text-md-on-surface tracking-tight">
            {filtered.length}
          </p>
          <p className="mt-0.5 text-sm text-md-on-surface-variant">
            {t.attendance.totalCheckIns} ({isSelectedToday ? t.attendance.today : selectedDate})
          </p>
        </div>
      </div>

      {/* Filters Row */}
      <div className="flex flex-wrap items-center gap-3">
        {/* Date Picker */}
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 text-md-on-surface-variant" />
          <input
            type="date"
            value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            className="px-3 py-1.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-lg text-sm text-md-on-surface focus:outline-none focus:border-md-primary transition-colors"
          />
        </div>

        {/* Method Filter */}
        <div className="flex gap-2">
          {(["all", "qr", "nfc", "manual"] as MethodFilter[]).map((m) => {
            const label =
              m === "all"
                ? t.attendance.filterAll
                : m === "qr"
                  ? t.attendance.qr
                  : m === "nfc"
                    ? t.attendance.nfc
                    : t.attendance.manual;
            return (
              <button
                key={m}
                onClick={() => setMethodFilter(m)}
                className={`px-3 py-1.5 text-xs font-mono uppercase rounded-full border transition-colors ${
                  methodFilter === m
                    ? "bg-md-primary-container text-md-on-primary-container border-transparent"
                    : "text-md-on-surface-variant border-md-outline-variant hover:bg-md-surface-container-high"
                }`}
              >
                {label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
          {error}
        </div>
      )}

      {/* Table */}
      <div className="bg-md-surface-container border border-md-outline-variant rounded-2xl shadow-md-1 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-md-primary animate-spin" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="text-center py-16 text-md-on-surface-variant text-sm">
            {t.attendance.noRecords}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.attendance.member}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.attendance.date}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.attendance.time}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.attendance.method}
                  </th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((record) => (
                  <tr
                    key={record.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="px-5 py-3 text-md-on-surface">
                      {record.member_name ?? record.user_id}
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {formatDate(record.check_in_at)}
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {formatTime(record.check_in_at)}
                    </td>
                    <td className="px-5 py-3">
                      <MethodBadge method={record.method} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Manual Check-in Modal */}
      {showCheckIn && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-md-surface-container border border-md-outline-variant rounded-2xl shadow-md-1 w-full max-w-sm">
            <div className="flex items-center justify-between p-5 border-b border-md-outline-variant">
              <h2 className="text-base font-semibold text-md-on-surface">
                {t.attendance.manualCheckIn}
              </h2>
              <button
                onClick={() => setShowCheckIn(false)}
                className="p-1 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleManualCheckIn} className="p-5 space-y-4">
              {checkInError && (
                <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
                  {checkInError}
                </div>
              )}
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.attendance.checkInMember}
                </label>
                <input
                  type="text"
                  required
                  value={checkInMemberId}
                  onChange={(e) => setCheckInMemberId(e.target.value)}
                  placeholder={t.attendance.memberIdPlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowCheckIn(false)}
                  className="flex-1 px-4 py-2.5 bg-md-surface-container-high border border-md-outline-variant text-md-on-surface text-sm rounded-full hover:bg-md-surface-container-highest transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={checkInLoading}
                  className="flex-1 px-4 py-2.5 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {checkInLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.attendance.manualCheckIn}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
