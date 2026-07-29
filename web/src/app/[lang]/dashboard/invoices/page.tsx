"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Loader2, Plus, Receipt } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Invoice, Locale } from "@/types";

/**
 * Bills issued to customers under the gym's PAN.
 *
 * There is no edit action anywhere in this feature — a bill, once issued, is
 * a legal record. Corrections happen through Cancel or a credit note on the
 * detail screen, never by mutating what was printed.
 */
export default function InvoicesPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;
  const base = `/${locale}/dashboard`;

  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  // pan_not_configured is a distinct error code precisely so this screen can
  // link the operator to settings instead of dead-ending on a bare message.
  const [panMissing, setPanMissing] = useState(false);

  const fetchInvoices = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    setPanMissing(false);
    try {
      const res = await api.listInvoices(orgId, { limit: 100 });
      setInvoices(res.data ?? []);
      setTotal(res.total ?? 0);
    } catch (err) {
      if (err instanceof ApiRequestError && err.code === "pan_not_configured") {
        setPanMissing(true);
      } else {
        setError(err instanceof ApiRequestError ? err.message : t.common.error);
      }
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchInvoices();
  }, [fetchInvoices]);

  const statusLabel = (inv: Invoice) => {
    if (inv.doc_type === "credit_note") return t.invoices.creditNote;
    return inv.status === "cancelled" ? t.invoices.cancelled : t.invoices.issued;
  };

  const statusColor = (inv: Invoice) => {
    if (inv.doc_type === "credit_note") return "bg-amber-500/10 text-amber-400 border-amber-500/20";
    return inv.status === "cancelled"
      ? "bg-red-500/10 text-red-400 border-red-500/20"
      : "bg-accent/10 text-accent border-accent/20";
  };

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
            <Receipt className="w-5 h-5 text-accent" />
            {t.invoices.title}
          </h1>
          <p className="text-sm text-fg-muted mt-1">{t.invoices.subtitle}</p>
        </div>
        <Link
          href={`${base}/invoices/new`}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.invoices.newBill}
        </Link>
      </div>

      {panMissing && (
        <div className="p-4 bg-amber-500/10 border border-amber-500/20 rounded-xl text-sm text-amber-300 flex flex-wrap items-center justify-between gap-3">
          <span>{t.invoices.panMissing}</span>
          <Link
            href={`${base}/settings`}
            className="px-3 py-1.5 rounded-lg bg-amber-500/20 text-amber-200 font-medium hover:bg-amber-500/30 transition-colors shrink-0"
          >
            {t.invoices.panMissingCta}
          </Link>
        </div>
      )}

      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-accent animate-spin" />
          </div>
        ) : invoices.length === 0 ? (
          <div className="text-center py-16 text-fg-muted text-sm">{t.invoices.empty}</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.invoices.number}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.invoices.customer}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.invoices.date}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.invoices.amount}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.invoices.status}
                  </th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((inv) => (
                  <tr
                    key={inv.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="px-5 py-3">
                      <Link
                        href={`${base}/invoices/${inv.id}`}
                        className="font-mono text-xs text-fg hover:text-accent transition-colors"
                      >
                        {inv.invoice_number}
                      </Link>
                    </td>
                    <td className="px-5 py-3 text-fg font-medium">{inv.customer_name}</td>
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                      {inv.issued_date_bs}
                    </td>
                    <td className="px-5 py-3 text-right font-mono text-fg">
                      {inv.total.toLocaleString()}
                    </td>
                    <td className="px-5 py-3">
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${statusColor(inv)}`}
                      >
                        {statusLabel(inv)}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {total > 0 && (
          <div className="px-5 py-3 border-t border-white/[0.06]">
            <p className="text-xs text-fg-muted">
              {total} {t.invoices.title.toLowerCase()}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
