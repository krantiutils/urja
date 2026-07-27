"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, OrgMember, SmsBalance, SmsPurchase } from "@/types";
import {
  MessageSquare,
  Send,
  Loader2,
  Mail,
  ShoppingCart,
  BarChart3,
  CreditCard,
  Search,
  X,
} from "lucide-react";

// The members list endpoint caps `limit` at 100 per page (anything higher
// resets to the default of 20), so paging is required to build a full
// recipient picker. MAX_MEMBER_PAGES is a sane upper bound (~5,000 members)
// to guarantee this loop always terminates.
const MEMBER_FETCH_LIMIT = 100;
const MAX_MEMBER_PAGES = 50;

function formatDate(isoStr: string): string {
  return new Date(isoStr).toLocaleDateString();
}

function PaymentStatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    completed: "bg-accent/10 text-accent border-accent/20",
    paid: "bg-accent/10 text-accent border-accent/20",
    pending: "bg-amber-500/10 text-amber-400 border-amber-500/20",
    failed: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[status.toLowerCase()] ?? "bg-gray-500/10 text-gray-400 border-gray-500/20"
      }`}
    >
      {status}
    </span>
  );
}

export default function SmsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const orgId = user?.org_id;

  const [balance, setBalance] = useState<SmsBalance | null>(null);
  const [purchases, setPurchases] = useState<SmsPurchase[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Send SMS state
  const [message, setMessage] = useState("");
  const [sending, setSending] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const [sendSuccess, setSendSuccess] = useState(false);

  // Recipient picker state
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [membersLoading, setMembersLoading] = useState(true);
  const [memberSearch, setMemberSearch] = useState("");
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [showConfirm, setShowConfirm] = useState(false);

  const fetchData = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const [balanceRes, historyRes] = await Promise.all([
        api.getSmsBalance(orgId),
        api.getSmsHistory(orgId, { limit: 50 }),
      ]);
      setBalance(balanceRes);
      setPurchases((historyRes.data ?? []) as SmsPurchase[]);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  const fetchAllMembers = useCallback(async () => {
    if (!orgId) return;
    setMembersLoading(true);
    try {
      let offset = 0;
      let all: OrgMember[] = [];
      for (let page = 0; page < MAX_MEMBER_PAGES; page++) {
        const res = await api.listMembers(orgId, {
          limit: MEMBER_FETCH_LIMIT,
          offset,
        });
        const batch = res.data ?? [];
        all = all.concat(batch);
        const total = res.total ?? all.length;
        if (batch.length < MEMBER_FETCH_LIMIT || all.length >= total) break;
        offset += MEMBER_FETCH_LIMIT;
      }
      setMembers(all);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setMembersLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  useEffect(() => {
    fetchAllMembers();
  }, [fetchAllMembers]);

  const filteredMembers = useMemo(() => {
    const q = memberSearch.trim().toLowerCase();
    if (!q) return members;
    return members.filter(
      (m) => m.name.toLowerCase().includes(q) || m.phone.includes(q)
    );
  }, [members, memberSearch]);

  const allFilteredSelected =
    filteredMembers.length > 0 &&
    filteredMembers.every((m) => selectedIds.has(m.id));

  const toggleMember = (id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const handleSelectAllToggle = () => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (allFilteredSelected) {
        filteredMembers.forEach((m) => next.delete(m.id));
      } else {
        filteredMembers.forEach((m) => next.add(m.id));
      }
      return next;
    });
  };

  const insufficientBalance =
    balance != null && selectedIds.size > balance.balance;

  const handleOpenConfirm = (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId || !message.trim() || selectedIds.size === 0) return;
    setSendError(null);
    setSendSuccess(false);
    setShowConfirm(true);
  };

  const handleConfirmSend = async () => {
    if (!orgId) return;
    setSending(true);
    setSendError(null);
    setSendSuccess(false);
    try {
      // member_ids is always sent explicitly and non-empty here — Send is
      // disabled until at least one recipient is picked, and "Select all"
      // fills this array rather than leaving it empty, since the backend
      // treats an empty list as "every active member in the org".
      await api.sendSms(orgId, {
        message: message.trim(),
        member_ids: Array.from(selectedIds),
      });
      setSendSuccess(true);
      setMessage("");
      setSelectedIds(new Set());
      setShowConfirm(false);
      // Refresh balance after sending
      const balanceRes = await api.getSmsBalance(orgId);
      setBalance(balanceRes);
    } catch (err) {
      setSendError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
      setShowConfirm(false);
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <MessageSquare className="w-5 h-5 text-accent" />
          {t.sms.title}
        </h1>
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Balance Cards */}
      {loading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Remaining Balance */}
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
              <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
                <Mail className="w-5 h-5 text-accent" />
              </div>
              <p className="text-2xl font-bold text-fg mt-2">
                {balance?.balance ?? 0}
              </p>
              <p className="text-xs text-fg-muted font-mono uppercase tracking-widest mt-1">
                {t.sms.remaining}
              </p>
            </div>

            {/* Total Purchased */}
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
              <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
                <ShoppingCart className="w-5 h-5 text-accent" />
              </div>
              <p className="text-2xl font-bold text-fg mt-2">
                {balance?.total_purchased ?? 0}
              </p>
              <p className="text-xs text-fg-muted font-mono uppercase tracking-widest mt-1">
                {t.sms.totalPurchased}
              </p>
            </div>

            {/* Total Used */}
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
              <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
                <Send className="w-5 h-5 text-accent" />
              </div>
              <p className="text-2xl font-bold text-fg mt-2">
                {balance?.total_used ?? 0}
              </p>
              <p className="text-xs text-fg-muted font-mono uppercase tracking-widest mt-1">
                {t.sms.totalUsed}
              </p>
            </div>

            {/* Total Campaigns */}
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
              <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
                <BarChart3 className="w-5 h-5 text-accent" />
              </div>
              <p className="text-2xl font-bold text-fg mt-2">
                {balance?.total_campaigns ?? 0}
              </p>
              <p className="text-xs text-fg-muted font-mono uppercase tracking-widest mt-1">
                {t.sms.totalCampaigns}
              </p>
            </div>
          </div>

          {/* Send SMS Section */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
            <h2 className="text-sm font-semibold text-fg flex items-center gap-2 mb-4">
              <Send className="w-4 h-4 text-fg-muted" />
              {t.sms.sendSms}
            </h2>

            {sendError && (
              <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                {sendError}
              </div>
            )}

            {sendSuccess && (
              <div className="mb-4 p-3 bg-accent/10 border border-accent/20 rounded-xl text-sm text-accent">
                {t.sms.sendSuccessMessage}
              </div>
            )}

            <form onSubmit={handleOpenConfirm} className="space-y-4">
              <div>
                <label htmlFor="sms-message" className="block text-xs text-fg-muted mb-1.5">
                  {t.sms.message}
                </label>
                <textarea
                  id="sms-message"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder={t.sms.messagePlaceholder}
                  required
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors min-h-[100px] resize-none"
                />
              </div>

              <div>
                <div className="flex items-center justify-between mb-1.5">
                  <label className="block text-xs text-fg-muted">
                    {t.sms.selectMembers} *
                  </label>
                  <span className="text-xs text-fg-muted font-mono">
                    {selectedIds.size} {t.sms.of} {members.length}{" "}
                    {t.sms.selected}
                  </span>
                </div>

                <div className="relative mb-2">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                  <input
                    type="text"
                    value={memberSearch}
                    onChange={(e) => setMemberSearch(e.target.value)}
                    placeholder={t.sms.searchMembersPlaceholder}
                    className="w-full pl-10 pr-4 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                  />
                </div>

                <div className="border border-white/[0.06] rounded-xl overflow-hidden">
                  <button
                    type="button"
                    onClick={handleSelectAllToggle}
                    disabled={filteredMembers.length === 0}
                    className="w-full flex items-center gap-2 px-3 py-2 text-xs font-medium text-accent hover:bg-surface transition-colors border-b border-white/[0.06] disabled:opacity-40"
                  >
                    <input
                      type="checkbox"
                      readOnly
                      checked={allFilteredSelected}
                      className="w-3.5 h-3.5 accent-accent pointer-events-none"
                    />
                    {allFilteredSelected ? t.sms.deselectAll : t.sms.selectAll}
                  </button>

                  <div className="max-h-56 overflow-y-auto divide-y divide-white/[0.03]">
                    {membersLoading ? (
                      <div className="flex items-center justify-center py-6 text-xs text-fg-muted gap-2">
                        <Loader2 className="w-4 h-4 text-accent animate-spin" />
                        {t.sms.loadingMembers}
                      </div>
                    ) : filteredMembers.length === 0 ? (
                      <div className="text-center py-6 text-fg-muted text-xs">
                        {t.sms.noMembersFound}
                      </div>
                    ) : (
                      filteredMembers.map((m) => (
                        <label
                          key={m.id}
                          className="flex items-center gap-3 px-3 py-2 text-sm hover:bg-surface transition-colors cursor-pointer"
                        >
                          <input
                            type="checkbox"
                            checked={selectedIds.has(m.id)}
                            onChange={() => toggleMember(m.id)}
                            className="w-3.5 h-3.5 accent-accent"
                          />
                          <span className="flex-1 min-w-0">
                            <span className="block text-fg truncate">
                              {m.name}
                            </span>
                            <span className="block text-xs text-fg-muted font-mono">
                              {m.phone}
                            </span>
                          </span>
                        </label>
                      ))
                    )}
                  </div>
                </div>
              </div>

              <div className="flex justify-end">
                <button
                  type="submit"
                  disabled={sending || !message.trim() || selectedIds.size === 0}
                  className="flex items-center gap-2 px-5 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow disabled:opacity-50"
                >
                  {sending ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Send className="w-4 h-4" />
                  )}
                  {t.sms.send}
                </button>
              </div>
            </form>
          </div>

          {/* Confirm Send Modal */}
          {showConfirm && (
            <div
              className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4"
              data-testid="sms-confirm-modal"
            >
              <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-sm">
                <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
                  <h2 className="text-base font-semibold text-fg">
                    {t.sms.confirmSendTitle}
                  </h2>
                  <button
                    onClick={() => setShowConfirm(false)}
                    className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
                <div className="p-5 space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-fg-muted">
                      {t.sms.confirmRecipientsLabel}
                    </span>
                    <span
                      className="text-fg font-mono font-semibold"
                      data-testid="sms-confirm-recipients-count"
                    >
                      {selectedIds.size}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-fg-muted">
                      {t.sms.confirmCreditsLabel}
                    </span>
                    <span
                      className="text-fg font-mono font-semibold"
                      data-testid="sms-confirm-credits-count"
                    >
                      {selectedIds.size}
                    </span>
                  </div>
                  {insufficientBalance && (
                    <p className="text-xs text-red-400">
                      {t.sms.insufficientBalance}
                    </p>
                  )}
                  <p className="text-xs text-fg-muted">
                    {t.sms.confirmSendWarning}
                  </p>
                  <div className="flex gap-3 pt-2">
                    <button
                      type="button"
                      onClick={() => setShowConfirm(false)}
                      data-testid="sms-confirm-cancel-btn"
                      className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                    >
                      {t.common.cancel}
                    </button>
                    <button
                      type="button"
                      onClick={handleConfirmSend}
                      disabled={sending || insufficientBalance}
                      data-testid="sms-confirm-send-btn"
                      className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                    >
                      {sending && <Loader2 className="w-4 h-4 animate-spin" />}
                      {t.sms.confirmSend}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Purchase History */}
          <div>
            <h2 className="text-sm font-semibold text-fg flex items-center gap-2 mb-3">
              <CreditCard className="w-4 h-4 text-fg-muted" />
              {t.sms.history}
            </h2>

            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
              {purchases.length === 0 ? (
                <div className="text-center py-16 text-fg-muted text-sm">
                  {t.sms.noHistory}
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-white/[0.06]">
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.quantity}
                        </th>
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.rate}
                        </th>
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.amount}
                        </th>
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.paymentMethod}
                        </th>
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.paymentStatus}
                        </th>
                        <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                          {t.sms.date}
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {purchases.map((purchase) => (
                        <tr
                          key={purchase.id}
                          className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                        >
                          <td className="px-5 py-3 text-fg font-mono">
                            {purchase.quantity}
                          </td>
                          <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                            {purchase.rate}
                          </td>
                          <td className="px-5 py-3 text-fg font-mono">
                            {purchase.amount}
                          </td>
                          <td className="px-5 py-3 text-fg-muted">
                            {purchase.payment_method}
                          </td>
                          <td className="px-5 py-3">
                            <PaymentStatusBadge status={purchase.payment_status} />
                          </td>
                          <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                            {formatDate(purchase.created_at)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
