"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, OrgMember, OrgPackage, CreateMemberRequest } from "@/types";
import {
  Users,
  CreditCard,
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
    active: "bg-accent/10 text-accent border-accent/20",
    suspended: "bg-amber-500/10 text-amber-400 border-amber-500/20",
    left: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[status] ?? colors.active
      }`}
    >
      {label}
    </span>
  );
}

function RoleBadge({ role, label }: { role: string; label: string }) {
  const colors: Record<string, string> = {
    admin: "bg-purple-500/10 text-purple-400 border-purple-500/20",
    staff: "bg-blue-500/10 text-blue-400 border-blue-500/20",
    member: "bg-gray-500/10 text-gray-400 border-gray-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
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

  // Selling a package to a member. The gym could define packages and add
  // members but had no way to put the two together — the central transaction
  // of the business was unreachable from the product.
  const [sellFor, setSellFor] = useState<OrgMember | null>(null);
  const [packages, setPackages] = useState<OrgPackage[]>([]);
  const [sellPackageId, setSellPackageId] = useState("");
  const [sellStart, setSellStart] = useState("");
  const [sellPaid, setSellPaid] = useState("");
  const [sellDiscount, setSellDiscount] = useState("");
  const [sellMethod, setSellMethod] = useState("cash");
  const [sellError, setSellError] = useState<string | null>(null);
  const [sellLoading, setSellLoading] = useState(false);
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

  const openSell = useCallback(
    async (member: OrgMember) => {
      setSellFor(member);
      setSellPackageId("");
      setSellStart(new Date().toISOString().slice(0, 10));
      setSellPaid("");
      setSellDiscount("");
      setSellMethod("cash");
      setSellError(null);

      if (!orgId || packages.length > 0) return;
      try {
        const res = await api.listPackages(orgId, { limit: 100 });
        setPackages((res.data ?? []).filter((p) => p.is_active));
      } catch {
        setSellError(t.common.error);
      }
    },
    [orgId, packages.length, t.common.error]
  );

  const handleSell = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!orgId || !sellFor || !sellPackageId) return;

      setSellLoading(true);
      setSellError(null);
      try {
        await api.assignPackage(orgId, sellFor.id, {
          package_id: sellPackageId,
          start_date: sellStart,
          payment_method: sellMethod,
          amount_paid: parseFloat(sellPaid) || 0,
          discount: parseFloat(sellDiscount) || 0,
        });
        setSellFor(null);
        await fetchMembers();
      } catch (err) {
        setSellError(err instanceof ApiRequestError ? err.message : t.common.error);
      } finally {
        setSellLoading(false);
      }
    },
    [orgId, sellFor, sellPackageId, sellStart, sellMethod, sellPaid, sellDiscount, fetchMembers, t]
  );

  const selectedPackage = packages.find((p) => p.id === sellPackageId);
  const sellBalance = selectedPackage
    ? Math.max(
        0,
        parseFloat(selectedPackage.price) -
          (parseFloat(sellDiscount) || 0) -
          (parseFloat(sellPaid) || 0)
      )
    : 0;

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
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Users className="w-5 h-5 text-accent" />
          {t.members.title}
        </h1>
        <button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.members.addMember}
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t.members.searchPlaceholder}
          className="w-full pl-10 pr-4 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
        />
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
        ) : filteredMembers.length === 0 ? (
          <div className="text-center py-16 text-fg-muted text-sm">
            {t.members.noMembers}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.name}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.phone}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.role}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.status}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.joined}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.members.actions}
                  </th>
                </tr>
              </thead>
              <tbody>
                {filteredMembers.map((member) => (
                  <tr
                    key={member.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="px-5 py-3">
                      <div>
                        <p className="text-fg font-medium">{member.name}</p>
                        {member.email && (
                          <p className="text-xs text-fg-muted">
                            {member.email}
                          </p>
                        )}
                      </div>
                    </td>
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
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
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                      {new Date(member.joined_at).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center justify-end gap-1">
                        <button
                          onClick={() => openSell(member)}
                          title={t.members.sellPackage}
                          className="p-1.5 rounded-lg hover:bg-accent/10 text-fg-muted hover:text-accent transition-colors"
                        >
                          <CreditCard className="w-4 h-4" />
                        </button>
                        {member.status === "active" ? (
                          <button
                            onClick={() =>
                              handleStatusChange(member.id, "suspended")
                            }
                            title={t.members.suspend}
                            className="p-1.5 rounded-lg hover:bg-amber-500/10 text-fg-muted hover:text-amber-400 transition-colors"
                          >
                            <UserMinus className="w-4 h-4" />
                          </button>
                        ) : member.status === "suspended" ? (
                          <button
                            onClick={() =>
                              handleStatusChange(member.id, "active")
                            }
                            title={t.members.activate}
                            className="p-1.5 rounded-lg hover:bg-accent/10 text-fg-muted hover:text-accent transition-colors"
                          >
                            <UserCheck className="w-4 h-4" />
                          </button>
                        ) : null}
                        <button
                          onClick={() => handleRemove(member.id)}
                          title={t.members.remove}
                          className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
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
          <div className="flex items-center justify-between px-5 py-3 border-t border-white/[0.06]">
            <p className="text-xs text-fg-muted">
              {t.members.showing} {offset + 1}–
              {Math.min(offset + PAGE_SIZE, total)} {t.members.of} {total}
            </p>
            <div className="flex items-center gap-1">
              <button
                onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))}
                disabled={offset === 0}
                className="p-1.5 rounded-lg hover:bg-surface text-fg-muted disabled:opacity-30 transition-colors"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <span className="text-xs text-fg-muted px-2">
                {currentPage} / {totalPages}
              </span>
              <button
                onClick={() =>
                  setOffset(
                    offset + PAGE_SIZE < total ? offset + PAGE_SIZE : offset
                  )
                }
                disabled={offset + PAGE_SIZE >= total}
                className="p-1.5 rounded-lg hover:bg-surface text-fg-muted disabled:opacity-30 transition-colors"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Add Member Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {t.members.addMemberTitle}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleAddMember} className="p-5 space-y-4">
              {formError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {formError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.members.name} *
                </label>
                <input
                  type="text"
                  required
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  placeholder={t.members.namePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.members.phone} *
                </label>
                <input
                  type="tel"
                  required
                  value={formPhone}
                  onChange={(e) => setFormPhone(e.target.value)}
                  placeholder={t.members.phonePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.members.email}
                </label>
                <input
                  type="email"
                  value={formEmail}
                  onChange={(e) => setFormEmail(e.target.value)}
                  placeholder={t.members.emailPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.members.role}
                </label>
                <select
                  value={formRole}
                  onChange={(e) =>
                    setFormRole(
                      e.target.value as "member" | "staff" | "admin"
                    )
                  }
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50 transition-colors"
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
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={formLoading}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
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

      {/* Sell a package to a member */}
      {sellFor && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black/60"
            onClick={() => setSellFor(null)}
            aria-hidden="true"
          />
          <form
            onSubmit={handleSell}
            className="relative w-full max-w-md bg-bg-elevated border border-white/[0.06] rounded-2xl p-5 shadow-card"
          >
            <div className="flex items-center justify-between mb-1">
              <h2 className="text-base font-semibold text-fg">{t.members.assignPackage}</h2>
              <button
                type="button"
                onClick={() => setSellFor(null)}
                className="p-1.5 rounded-lg text-fg-muted hover:bg-surface transition-colors"
                aria-label={t.common.cancel}
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <p className="text-sm text-fg-muted mb-4">{sellFor.name}</p>

            {sellError && (
              <p className="mb-3 text-sm text-red-400" role="alert">
                {sellError}
              </p>
            )}

            <div className="mb-4">
              <label htmlFor="sell-pkg" className="block text-xs text-fg-muted mb-1.5">
                {t.members.choosePackage}
              </label>
              <select
                id="sell-pkg"
                value={sellPackageId}
                onChange={(e) => setSellPackageId(e.target.value)}
                className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
              >
                <option value="">—</option>
                {packages.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name} · {p.currency} {p.price} · {p.duration_days}d
                  </option>
                ))}
              </select>
            </div>

            <div className="grid grid-cols-2 gap-3 mb-4">
              <div>
                <label htmlFor="sell-start" className="block text-xs text-fg-muted mb-1.5">
                  {t.members.startDate}
                </label>
                <input
                  id="sell-start"
                  type="date"
                  value={sellStart}
                  onChange={(e) => setSellStart(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                />
              </div>
              <div>
                <label htmlFor="sell-method" className="block text-xs text-fg-muted mb-1.5">
                  {t.members.paymentMethod}
                </label>
                <select
                  id="sell-method"
                  value={sellMethod}
                  onChange={(e) => setSellMethod(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                >
                  <option value="cash">{t.dues.cash}</option>
                  <option value="esewa">{t.dues.esewa}</option>
                  <option value="bank">{t.dues.bank}</option>
                </select>
              </div>
              <div>
                <label htmlFor="sell-paid" className="block text-xs text-fg-muted mb-1.5">
                  {t.members.amountPaid}
                </label>
                <input
                  id="sell-paid"
                  type="number"
                  min="0"
                  step="0.01"
                  value={sellPaid}
                  onChange={(e) => setSellPaid(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                />
              </div>
              <div>
                <label htmlFor="sell-discount" className="block text-xs text-fg-muted mb-1.5">
                  {t.members.discount}
                </label>
                <input
                  id="sell-discount"
                  type="number"
                  min="0"
                  step="0.01"
                  value={sellDiscount}
                  onChange={(e) => setSellDiscount(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                />
              </div>
            </div>

            {/* The balance is the thing a gym gets wrong on paper, so it is
                shown before saving rather than discovered later on the dues
                screen. */}
            {selectedPackage && (
              <div className="mb-5 px-3 py-2.5 rounded-xl bg-surface text-sm">
                <div className="flex justify-between text-fg-muted">
                  <span>{selectedPackage.currency} {selectedPackage.price}</span>
                  <span>
                    {t.members.amountPaid}: {parseFloat(sellPaid) || 0}
                  </span>
                </div>
                {sellBalance > 0 && (
                  <p className="mt-1.5 text-amber-400">
                    {t.dues.title}: {selectedPackage.currency} {sellBalance.toLocaleString()} — {t.members.balanceNote}
                  </p>
                )}
              </div>
            )}

            <button
              type="submit"
              disabled={sellLoading || !sellPackageId}
              className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {sellLoading && <Loader2 className="w-4 h-4 animate-spin" />}
              {t.members.assignPackage}
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
