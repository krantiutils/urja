"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, Due, OrgMember } from "@/types";
import {
  CreditCard,
  Plus,
  Search,
  X,
  Loader2,
} from "lucide-react";

type DueStatus = "unpaid" | "paid" | "waived";

function StatusBadge({
  status,
  label,
}: {
  status: string;
  label: string;
}) {
  const colors: Record<string, string> = {
    paid: "bg-accent/10 text-accent border-accent/20",
    unpaid: "bg-amber-500/10 text-amber-400 border-amber-500/20",
    waived: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[status] ?? colors.unpaid
      }`}
    >
      {label}
    </span>
  );
}

export default function DuePaymentsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [dues, setDues] = useState<Due[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<DueStatus | "all">("all");

  // Pay modal state
  const [showPayModal, setShowPayModal] = useState(false);
  const [payingDue, setPayingDue] = useState<Due | null>(null);
  const [payAmount, setPayAmount] = useState("");
  const [payMethod, setPayMethod] = useState<"cash" | "esewa" | "bank">("cash");
  const [payReference, setPayReference] = useState("");
  const [payError, setPayError] = useState<string | null>(null);
  const [payLoading, setPayLoading] = useState(false);

  // Recording a new due
  const [showNewModal, setShowNewModal] = useState(false);
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [memberSearch, setMemberSearch] = useState("");
  const [newMember, setNewMember] = useState<OrgMember | null>(null);
  const [newAmount, setNewAmount] = useState("");
  const [newDate, setNewDate] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const [newError, setNewError] = useState<string | null>(null);
  const [newLoading, setNewLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchDues = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listDues(orgId, {
        status: statusFilter === "all" ? undefined : statusFilter,
        search: search || undefined,
        limit: 100,
      });
      setDues(res.data ?? []);
      setTotal(res.total ?? 0);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, statusFilter, search, t.common.error]);

  useEffect(() => {
    fetchDues();
  }, [fetchDues]);

  const openPayModal = (due: Due) => {
    setPayingDue(due);
    setPayAmount(due.amount);
    setPayMethod("cash");
    setPayReference("");
    setPayError(null);
    setShowPayModal(true);
  };

  const handlePay = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId || !payingDue) return;
    setPayError(null);

    const parsedAmount = Number(payAmount);
    if (!Number.isFinite(parsedAmount) || parsedAmount <= 0) {
      setPayError(t.common.invalidAmount);
      return;
    }

    setPayLoading(true);

    try {
      await api.payDue(orgId, payingDue.id, {
        amount: parsedAmount,
        payment_method: payMethod,
        payment_reference: payReference || undefined,
      });
      setShowPayModal(false);
      fetchDues();
    } catch (err) {
      setPayError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setPayLoading(false);
    }
  };

  const statusLabel = (s: string) => {
    const map: Record<string, string> = {
      unpaid: t.dues.unpaid,
      paid: t.dues.paid,
      waived: t.dues.waived,
    };
    return map[s] ?? s;
  };

  /**
   * Opens the form and loads members to pick from. Fetched on open rather than
   * with the page: most visits here are to take a payment, not to raise a new
   * due, and a gym can have hundreds of members.
   */
  const openNewDue = useCallback(async () => {
    setNewMember(null);
    setNewAmount("");
    setNewDescription("");
    // Defaults to today, which is what a gym recording an unpaid month wants.
    setNewDate(new Date().toISOString().slice(0, 10));
    setNewError(null);
    setShowNewModal(true);

    if (!orgId || members.length > 0) return;
    try {
      const res = await api.listMembers(orgId, { limit: 500 });
      setMembers(res.data ?? []);
    } catch {
      setNewError(t.common.error);
    }
  }, [orgId, members.length, t.common.error]);

  const handleCreateDue = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!orgId || !newMember) return;

      const amount = parseFloat(newAmount);
      if (!Number.isFinite(amount) || amount <= 0) {
        setNewError(t.common.invalidAmount);
        return;
      }

      setNewLoading(true);
      setNewError(null);
      try {
        await api.createDue(orgId, {
          user_id: newMember.id,
          amount,
          due_date: newDate,
          description: newDescription.trim() || undefined,
        });
        setShowNewModal(false);
        await fetchDues();
      } catch (err) {
        setNewError(err instanceof ApiRequestError ? err.message : t.common.error);
      } finally {
        setNewLoading(false);
      }
    },
    [orgId, newMember, newAmount, newDate, newDescription, fetchDues, t]
  );

  const filteredMembers = memberSearch
    ? members.filter((m) =>
        `${m.name} ${m.phone}`.toLowerCase().includes(memberSearch.toLowerCase())
      )
    : members;

  const filterTabs: { key: DueStatus | "all"; label: string }[] = [
    { key: "all", label: t.dues.all },
    { key: "unpaid", label: t.dues.unpaid },
    { key: "paid", label: t.dues.paid },
    { key: "waived", label: t.dues.waived },
  ];

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <CreditCard className="w-5 h-5 text-accent" />
          {t.dues.title}
        </h1>
        <button
          type="button"
          onClick={openNewDue}
          className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity"
        >
          <Plus className="w-4 h-4" />
          {t.dues.recordDue}
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t.dues.search}
          className="w-full pl-10 pr-4 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
        />
      </div>

      {/* Status filter tabs */}
      <div className="flex items-center gap-1">
        {filterTabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setStatusFilter(tab.key)}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              statusFilter === tab.key
                ? "bg-accent/10 text-accent"
                : "text-fg-muted hover:bg-surface"
            }`}
          >
            {tab.label}
          </button>
        ))}
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
        ) : dues.length === 0 ? (
          <div className="text-center py-16 text-fg-muted text-sm">
            {t.dues.noDues}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.member}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.amount}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.dueDate}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.description}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.status}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.dues.actions}
                  </th>
                </tr>
              </thead>
              <tbody>
                {dues.map((due) => (
                  <tr
                    key={due.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="px-5 py-3">
                      <div>
                        <p className="text-fg font-medium">
                          {due.member_name}
                        </p>
                        <p className="text-xs text-fg-muted">
                          {due.member_phone}
                        </p>
                      </div>
                    </td>
                    <td className="px-5 py-3 text-fg font-mono text-xs">
                      NPR {due.amount}
                    </td>
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                      {new Date(due.due_date).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-3 text-fg-muted text-xs">
                      {due.description}
                    </td>
                    <td className="px-5 py-3">
                      <StatusBadge
                        status={due.status}
                        label={statusLabel(due.status)}
                      />
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center justify-end gap-1">
                        {due.status === "unpaid" && (
                          <button
                            onClick={() => openPayModal(due)}
                            className="px-3 py-1.5 bg-accent text-bg-deep font-medium text-xs rounded-lg hover:bg-accent-bright transition-colors"
                          >
                            {t.dues.payDue}
                          </button>
                        )}
                        {due.status === "paid" && due.paid_at && (
                          <span className="text-xs text-fg-muted font-mono">
                            {new Date(due.paid_at).toLocaleDateString()}
                          </span>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Count footer */}
        {total > 0 && (
          <div className="px-5 py-3 border-t border-white/[0.06]">
            <p className="text-xs text-fg-muted">
              {total} {t.dues.title.toLowerCase()}
            </p>
          </div>
        )}
      </div>

      {/* Pay Modal */}
      {showPayModal && payingDue && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {t.dues.payDue}
              </h2>
              <button
                onClick={() => setShowPayModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handlePay} className="p-5 space-y-4">
              {payError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {payError}
                </div>
              )}

              {/* Due info summary */}
              <div className="p-3 bg-surface rounded-xl space-y-1">
                <p className="text-sm text-fg font-medium">
                  {payingDue.member_name}
                </p>
                <p className="text-xs text-fg-muted">
                  {payingDue.description}
                </p>
                <p className="text-xs text-fg-muted font-mono">
                  {t.dues.dueDate}: {new Date(payingDue.due_date).toLocaleDateString()}
                </p>
              </div>

              <div>
                <label htmlFor="due-pay-amount" className="block text-xs text-fg-muted mb-1.5">
                  {t.dues.payAmount} *
                </label>
                <input
                  id="due-pay-amount"
                  type="text"
                  inputMode="decimal"
                  value={payAmount}
                  onChange={(e) => setPayAmount(e.target.value)}
                  placeholder="0.00"
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label htmlFor="due-pay-method" className="block text-xs text-fg-muted mb-1.5">
                  {t.dues.paymentMethod} *
                </label>
                <select
                  id="due-pay-method"
                  value={payMethod}
                  onChange={(e) =>
                    setPayMethod(e.target.value as "cash" | "esewa" | "bank")
                  }
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50 transition-colors"
                >
                  <option value="cash">{t.dues.cash}</option>
                  <option value="esewa">{t.dues.esewa}</option>
                  <option value="bank">{t.dues.bank}</option>
                </select>
              </div>
              <div>
                <label htmlFor="due-pay-reference" className="block text-xs text-fg-muted mb-1.5">
                  {t.dues.paymentReference}
                </label>
                <input
                  id="due-pay-reference"
                  type="text"
                  value={payReference}
                  onChange={(e) => setPayReference(e.target.value)}
                  placeholder={t.dues.paymentReference}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowPayModal(false)}
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={payLoading}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {payLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.dues.confirmPay}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Record a new due */}
      {showNewModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black/60"
            onClick={() => setShowNewModal(false)}
            aria-hidden="true"
          />
          <form
            onSubmit={handleCreateDue}
            className="relative w-full max-w-md bg-bg-elevated border border-white/[0.06] rounded-2xl p-5 shadow-card"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base font-semibold text-fg">{t.dues.newDue}</h2>
              <button
                type="button"
                onClick={() => setShowNewModal(false)}
                className="p-1.5 rounded-lg text-fg-muted hover:bg-surface transition-colors"
                aria-label={t.common.cancel}
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {newError && (
              <p className="mb-3 text-sm text-red-400" role="alert">
                {newError}
              </p>
            )}

            <div className="mb-4">
              <label className="block text-xs text-fg-muted mb-1.5">
                {t.dues.selectMember}
              </label>
              {newMember ? (
                <div className="flex items-center justify-between px-3 py-2.5 bg-input-bg border border-accent/30 rounded-xl text-sm text-fg">
                  <span>
                    {newMember.name}
                    <span className="ml-2 font-mono text-xs text-fg-muted">
                      {newMember.phone}
                    </span>
                  </span>
                  <button
                    type="button"
                    onClick={() => setNewMember(null)}
                    className="p-0.5 rounded hover:bg-surface text-fg-muted"
                    aria-label={t.common.cancel}
                  >
                    <X className="w-3.5 h-3.5" />
                  </button>
                </div>
              ) : (
                <>
                  <input
                    type="text"
                    value={memberSearch}
                    onChange={(e) => setMemberSearch(e.target.value)}
                    placeholder={t.dues.searchMember}
                    className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50"
                  />
                  <div className="mt-1 max-h-40 overflow-y-auto rounded-xl border border-white/[0.06]">
                    {filteredMembers.length === 0 ? (
                      <p className="px-3 py-3 text-sm text-fg-muted text-center">
                        {t.dues.noMemberSelected}
                      </p>
                    ) : (
                      filteredMembers.slice(0, 50).map((m) => (
                        <button
                          key={m.id}
                          type="button"
                          onClick={() => setNewMember(m)}
                          className="w-full text-left px-3 py-2.5 text-sm hover:bg-surface transition-colors border-b border-white/[0.03] last:border-0"
                        >
                          <span className="text-fg">{m.name}</span>
                          <span className="ml-2 font-mono text-xs text-fg-muted">
                            {m.phone}
                          </span>
                        </button>
                      ))
                    )}
                  </div>
                </>
              )}
            </div>

            <div className="grid grid-cols-2 gap-3 mb-4">
              <div>
                <label htmlFor="due-amount" className="block text-xs text-fg-muted mb-1.5">
                  {t.dues.amount}
                </label>
                <input
                  id="due-amount"
                  type="number"
                  min="1"
                  step="0.01"
                  value={newAmount}
                  onChange={(e) => setNewAmount(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                />
              </div>
              <div>
                <label htmlFor="due-date" className="block text-xs text-fg-muted mb-1.5">
                  {t.dues.dueDate}
                </label>
                <input
                  id="due-date"
                  type="date"
                  value={newDate}
                  onChange={(e) => setNewDate(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
                />
              </div>
            </div>

            <div className="mb-5">
              <label htmlFor="due-desc" className="block text-xs text-fg-muted mb-1.5">
                {t.dues.description}
              </label>
              <input
                id="due-desc"
                type="text"
                value={newDescription}
                onChange={(e) => setNewDescription(e.target.value)}
                maxLength={500}
                className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50"
              />
            </div>

            <button
              type="submit"
              disabled={newLoading || !newMember || !newAmount}
              className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {newLoading && <Loader2 className="w-4 h-4 animate-spin" />}
              {t.dues.saveDue}
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
