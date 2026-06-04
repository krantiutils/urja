"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, OrgMember, CreateMemberRequest } from "@/types";
import {
  Users,
  Plus,
  Search,
  ChevronLeft,
  ChevronRight,
  X,
  UserMinus,
  UserCheck,
  Loader2,
} from "lucide-react";

const PAGE_SIZE = 20;

function StatusBadge({
  status,
  label,
}: {
  status: string;
  label: string;
}) {
  const colors: Record<string, string> = {
    active: "bg-green-100 text-green-700",
    suspended: "bg-amber-100 text-amber-700",
    left: "bg-red-100 text-red-700",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase ${
        colors[status] ?? colors.active
      }`}
    >
      {label}
    </span>
  );
}

function RoleBadge({ role, label }: { role: string; label: string }) {
  const colors: Record<string, string> = {
    admin: "bg-purple-100 text-purple-700",
    staff: "bg-blue-100 text-blue-700",
    member: "bg-gray-100 text-gray-600",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase ${
        colors[role] ?? colors.member
      }`}
    >
      {label}
    </span>
  );
}

export default function MembersPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [members, setMembers] = useState<OrgMember[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal state
  const [showModal, setShowModal] = useState(false);
  const [formName, setFormName] = useState("");
  const [formPhone, setFormPhone] = useState("");
  const [formEmail, setFormEmail] = useState("");
  const [formRole, setFormRole] = useState<"member" | "staff" | "admin">("member");
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchMembers = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listMembers(orgId, {
        limit: PAGE_SIZE,
        offset,
      });
      setMembers(res.data ?? []);
      setTotal(res.total ?? 0);
    } catch (err) {
      setError(
        err instanceof ApiRequestError
          ? err.message
          : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, offset, t.common.error]);

  useEffect(() => {
    fetchMembers();
  }, [fetchMembers]);

  const filteredMembers = search
    ? members.filter(
        (m) =>
          m.name.toLowerCase().includes(search.toLowerCase()) ||
          m.phone.includes(search)
      )
    : members;

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setFormLoading(true);
    setFormError(null);
    try {
      const data: CreateMemberRequest = {
        phone: formPhone,
        name: formName,
        role: formRole,
      };
      if (formEmail) data.email = formEmail;
      await api.createMember(orgId, data);
      setShowModal(false);
      setFormName("");
      setFormPhone("");
      setFormEmail("");
      setFormRole("member");
      fetchMembers();
    } catch (err) {
      setFormError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setFormLoading(false);
    }
  };

  const handleStatusChange = async (
    memberId: string,
    newStatus: "active" | "suspended"
  ) => {
    if (!orgId) return;
    try {
      await api.updateMember(orgId, memberId, { status: newStatus });
      fetchMembers();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const handleRemove = async (memberId: string) => {
    if (!orgId) return;
    if (!confirm(t.members.confirmDelete)) return;
    try {
      await api.deleteMember(orgId, memberId);
      fetchMembers();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

  const statusLabel = (s: string) => {
    const map: Record<string, string> = {
      active: t.members.active,
      suspended: t.members.suspended,
      left: t.members.left,
    };
    return map[s] ?? s;
  };

  const roleLabel = (r: string) => {
    const map: Record<string, string> = {
      member: t.members.member,
      staff: t.members.staff,
      admin: t.members.admin,
    };
    return map[r] ?? r;
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-md-on-surface flex items-center gap-2">
          <Users className="w-5 h-5 text-md-primary" />
          {t.members.title}
        </h1>
        <button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors shadow-md-1"
        >
          <Plus className="w-4 h-4" />
          {t.members.addMember}
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-md-on-surface-variant" />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t.members.searchPlaceholder}
          className="w-full pl-10 pr-4 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
        />
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
          {error}
        </div>
      )}

      {/* Table */}
      <div className="bg-md-surface-container border border-md-outline-variant rounded-3xl shadow-md-1 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-md-primary animate-spin" />
          </div>
        ) : filteredMembers.length === 0 ? (
          <div className="text-center py-16 text-md-on-surface-variant text-sm">
            {t.members.noMembers}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.name}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.phone}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.role}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.status}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.joined}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.members.actions}
                  </th>
                </tr>
              </thead>
              <tbody>
                {filteredMembers.map((member) => (
                  <tr
                    key={member.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="px-5 py-3">
                      <div>
                        <p className="text-md-on-surface font-medium">{member.name}</p>
                        {member.email && (
                          <p className="text-xs text-md-on-surface-variant">
                            {member.email}
                          </p>
                        )}
                      </div>
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {member.phone}
                    </td>
                    <td className="px-5 py-3">
                      <RoleBadge role={member.role} label={roleLabel(member.role)} />
                    </td>
                    <td className="px-5 py-3">
                      <StatusBadge
                        status={member.status}
                        label={statusLabel(member.status)}
                      />
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {new Date(member.joined_at).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center justify-end gap-1">
                        {member.status === "active" ? (
                          <button
                            onClick={() =>
                              handleStatusChange(member.id, "suspended")
                            }
                            title={t.members.suspend}
                            className="p-1.5 rounded-lg hover:bg-amber-100 text-md-on-surface-variant hover:text-amber-600 transition-colors"
                          >
                            <UserMinus className="w-4 h-4" />
                          </button>
                        ) : member.status === "suspended" ? (
                          <button
                            onClick={() =>
                              handleStatusChange(member.id, "active")
                            }
                            title={t.members.activate}
                            className="p-1.5 rounded-lg hover:bg-md-primary-container text-md-on-surface-variant hover:text-md-primary transition-colors"
                          >
                            <UserCheck className="w-4 h-4" />
                          </button>
                        ) : null}
                        <button
                          onClick={() => handleRemove(member.id)}
                          title={t.members.remove}
                          className="p-1.5 rounded-lg hover:bg-red-100 text-md-on-surface-variant hover:text-red-600 transition-colors"
                        >
                          <X className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Pagination */}
        {total > PAGE_SIZE && (
          <div className="flex items-center justify-between px-5 py-3 border-t border-md-outline-variant">
            <p className="text-xs text-md-on-surface-variant">
              {t.members.showing} {offset + 1}–
              {Math.min(offset + PAGE_SIZE, total)} {t.members.of} {total}
            </p>
            <div className="flex items-center gap-1">
              <button
                onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))}
                disabled={offset === 0}
                className="p-1.5 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant disabled:opacity-30 transition-colors"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <span className="text-xs text-md-on-surface-variant px-2">
                {currentPage} / {totalPages}
              </span>
              <button
                onClick={() =>
                  setOffset(
                    offset + PAGE_SIZE < total ? offset + PAGE_SIZE : offset
                  )
                }
                disabled={offset + PAGE_SIZE >= total}
                className="p-1.5 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant disabled:opacity-30 transition-colors"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Add Member Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-md-surface-container border border-md-outline-variant rounded-3xl shadow-md-1 w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-md-outline-variant">
              <h2 className="text-base font-semibold text-md-on-surface">
                {t.members.addMemberTitle}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleAddMember} className="p-5 space-y-4">
              {formError && (
                <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
                  {formError}
                </div>
              )}
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.members.name} *
                </label>
                <input
                  type="text"
                  required
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  placeholder={t.members.namePlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.members.phone} *
                </label>
                <input
                  type="tel"
                  required
                  value={formPhone}
                  onChange={(e) => setFormPhone(e.target.value)}
                  placeholder={t.members.phonePlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.members.email}
                </label>
                <input
                  type="email"
                  value={formEmail}
                  onChange={(e) => setFormEmail(e.target.value)}
                  placeholder={t.members.emailPlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.members.role}
                </label>
                <select
                  value={formRole}
                  onChange={(e) =>
                    setFormRole(
                      e.target.value as "member" | "staff" | "admin"
                    )
                  }
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface focus:outline-none focus:border-md-primary transition-colors"
                >
                  <option value="member">{t.members.member}</option>
                  <option value="staff">{t.members.staff}</option>
                  <option value="admin">{t.members.admin}</option>
                </select>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2.5 bg-md-surface-container-high border border-md-outline-variant text-md-on-surface text-sm rounded-full hover:bg-md-surface-container-highest transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={formLoading}
                  className="flex-1 px-4 py-2.5 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {formLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.common.save}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
