"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, OrgAttendance, OrgMember } from "@/types";
import {
  CalendarCheck,
  UserPlus,
  X,
  Loader2,
  Calendar,
  Search,
  ChevronDown,
} from "lucide-react";

type MethodFilter = "all" | "qr" | "nfc" | "manual";

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
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [memberNameMap, setMemberNameMap] = useState<Map<string, string>>(
    new Map()
  );
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

  // Member picker state
  const [memberSearch, setMemberSearch] = useState("");
  const [showMemberDropdown, setShowMemberDropdown] = useState(false);
  const [selectedMemberName, setSelectedMemberName] = useState("");
  const dropdownRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const orgId = user?.org_id;

  const fetchData = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const [attendanceRes, membersRes] = await Promise.all([
        api.listAttendance(orgId, { limit: 100 }),
        api.listMembers(orgId, { limit: 200 }),
      ]);
      setRecords(attendanceRes.data ?? []);

      const membersList = membersRes.data ?? [];
      setMembers(membersList);

      const nameMap = new Map<string, string>();
      for (const m of membersList) {
        nameMap.set(m.id, m.name);
      }
      setMemberNameMap(nameMap);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(e.target as Node)
      ) {
        setShowMemberDropdown(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

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

  // Filtered members for the picker dropdown
  const filteredMembers = memberSearch
    ? members.filter(
        (m) =>
          m.name.toLowerCase().includes(memberSearch.toLowerCase()) ||
          m.phone.includes(memberSearch)
      )
    : members;

  const handleSelectMember = (member: OrgMember) => {
    setCheckInMemberId(member.id);
    setSelectedMemberName(member.name);
    setMemberSearch("");
    setShowMemberDropdown(false);
  };

  const handleClearMember = () => {
    setCheckInMemberId("");
    setSelectedMemberName("");
    setMemberSearch("");
    setShowMemberDropdown(false);
  };

  const handleManualCheckIn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId || !checkInMemberId) return;
    setCheckInLoading(true);
    setCheckInError(null);
    try {
      await api.manualCheckIn(orgId, { member_id: checkInMemberId });
      setShowCheckIn(false);
      setCheckInMemberId("");
      setSelectedMemberName("");
      setMemberSearch("");
      fetchData();
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
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <CalendarCheck className="w-5 h-5 text-accent" />
          {t.attendance.title}
        </h1>
        <button
          onClick={() => setShowCheckIn(true)}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <UserPlus className="w-4 h-4" />
          {t.attendance.manualCheckIn}
        </button>
      </div>

      {/* Stats Card */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center mb-3">
            <CalendarCheck className="w-5 h-5 text-accent" />
          </div>
          <p className="text-2xl font-semibold text-fg tracking-tight">
            {todayCount}
          </p>
          <p className="mt-0.5 text-sm text-fg-muted">
            {t.attendance.todayCheckIns}
          </p>
        </div>
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
          <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center mb-3">
            <Calendar className="w-5 h-5 text-accent" />
          </div>
          <p className="text-2xl font-semibold text-fg tracking-tight">
            {filtered.length}
          </p>
          <p className="mt-0.5 text-sm text-fg-muted">
            {t.attendance.totalCheckIns} ({isSelectedToday ? t.attendance.today : selectedDate})
          </p>
        </div>
      </div>

      {/* Filters Row */}
      <div className="flex flex-wrap items-center gap-3">
        {/* Date Picker */}
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 text-fg-muted" />
          <input
            type="date"
            value={selectedDate}
            onChange={(e) => setSelectedDate(e.target.value)}
            className="px-3 py-1.5 bg-input-bg border border-white/[0.06] rounded-lg text-sm text-fg focus:outline-none focus:border-accent/50 transition-colors"
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
                className={`px-3 py-1.5 text-xs font-mono uppercase rounded-lg border transition-colors ${
                  methodFilter === m
                    ? "bg-accent/10 text-accent border-accent/20"
                    : "text-fg-muted border-white/[0.06] hover:bg-surface"
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
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Table */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-accent animate-spin" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="text-center py-16 text-fg-muted text-sm">
            {t.attendance.noRecords}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.attendance.member}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.attendance.date}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.attendance.time}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.attendance.method}
                  </th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((record) => (
                  <tr
                    key={record.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="px-5 py-3 text-fg">
                      {memberNameMap.get(record.user_id) ?? record.user_id}
                    </td>
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                      {formatDate(record.check_in_at)}
                    </td>
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
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
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-sm">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {t.attendance.manualCheckIn}
              </h2>
              <button
                onClick={() => {
                  setShowCheckIn(false);
                  handleClearMember();
                  setCheckInError(null);
                }}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleManualCheckIn} className="p-5 space-y-4">
              {checkInError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {checkInError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.attendance.checkInMember}
                </label>

                {/* Searchable member picker */}
                <div ref={dropdownRef} className="relative">
                  {selectedMemberName ? (
                    /* Selected member display */
                    <div className="w-full px-3 py-2.5 bg-input-bg border border-accent/30 rounded-xl text-sm text-fg flex items-center justify-between">
                      <span>{selectedMemberName}</span>
                      <button
                        type="button"
                        onClick={handleClearMember}
                        className="p-0.5 rounded hover:bg-surface text-fg-muted transition-colors"
                      >
                        <X className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  ) : (
                    /* Search input */
                    <div className="relative">
                      <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                      <input
                        ref={inputRef}
                        type="text"
                        value={memberSearch}
                        onChange={(e) => {
                          setMemberSearch(e.target.value);
                          setShowMemberDropdown(true);
                        }}
                        onFocus={() => setShowMemberDropdown(true)}
                        placeholder={t.members?.searchPlaceholder ?? "Search by name or phone..."}
                        className="w-full pl-9 pr-8 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                      />
                      <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted pointer-events-none" />
                    </div>
                  )}

                  {/* Dropdown */}
                  {showMemberDropdown && !selectedMemberName && (
                    <div className="absolute z-10 mt-1 w-full max-h-48 overflow-y-auto bg-bg-elevated border border-white/[0.06] rounded-xl shadow-card">
                      {filteredMembers.length === 0 ? (
                        <div className="px-3 py-3 text-sm text-fg-muted text-center">
                          {t.members?.noMembers ?? "No members found"}
                        </div>
                      ) : (
                        filteredMembers.map((member) => (
                          <button
                            key={member.id}
                            type="button"
                            onClick={() => handleSelectMember(member)}
                            className="w-full text-left px-3 py-2.5 hover:bg-surface transition-colors flex items-center justify-between gap-2 border-b border-white/[0.03] last:border-0"
                          >
                            <div className="min-w-0">
                              <p className="text-sm text-fg truncate">
                                {member.name}
                              </p>
                              <p className="text-xs text-fg-muted font-mono">
                                {member.phone}
                              </p>
                            </div>
                            <span
                              className={`shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-mono uppercase border ${
                                member.status === "active"
                                  ? "bg-accent/10 text-accent border-accent/20"
                                  : "bg-gray-500/10 text-gray-400 border-gray-500/20"
                              }`}
                            >
                              {member.status}
                            </span>
                          </button>
                        ))
                      )}
                    </div>
                  )}
                </div>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => {
                    setShowCheckIn(false);
                    handleClearMember();
                    setCheckInError(null);
                  }}
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={checkInLoading || !checkInMemberId}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
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
