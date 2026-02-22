"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, SmsBalance, SmsPurchase } from "@/types";
import {
  MessageSquare,
  Send,
  Loader2,
  Mail,
  ShoppingCart,
  BarChart3,
  CreditCard,
} from "lucide-react";

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

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleSendSms = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId || !message.trim()) return;
    setSending(true);
    setSendError(null);
    setSendSuccess(false);
    try {
      await api.sendSms(orgId, { message: message.trim(), member_ids: [] });
      setSendSuccess(true);
      setMessage("");
      // Refresh balance after sending
      const balanceRes = await api.getSmsBalance(orgId);
      setBalance(balanceRes);
    } catch (err) {
      setSendError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
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
                SMS sent successfully
              </div>
            )}

            <form onSubmit={handleSendSms} className="space-y-4">
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.sms.message}
                </label>
                <textarea
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder={t.sms.messagePlaceholder}
                  required
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors min-h-[100px] resize-none"
                />
              </div>

              <p className="text-xs text-fg-muted">
                {t.sms.selectMembers}: Select members from the Members page
              </p>

              <div className="flex justify-end">
                <button
                  type="submit"
                  disabled={sending || !message.trim()}
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
