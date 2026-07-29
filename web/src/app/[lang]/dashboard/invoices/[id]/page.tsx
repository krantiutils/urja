"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Ban, Loader2, Printer, Receipt, RefreshCw, X } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import { InvoiceDocument } from "@/components/invoice/InvoiceDocument";
import "@/components/invoice/invoice-print.css";
import type { Invoice, Locale } from "@/types";

/**
 * A single bill: print it, cancel it, or revise it. There is no edit control
 * here, by design — an issued bill is a legal record, not a draft.
 */
export default function InvoiceDetailPage({
  params,
}: {
  params: { lang: string; id: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;
  const router = useRouter();
  const base = `/${locale}/dashboard`;

  const [invoice, setInvoice] = useState<Invoice | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [copyLabel, setCopyLabel] = useState<"original" | "copy" | undefined>();
  const [printPending, setPrintPending] = useState(false);
  const [printError, setPrintError] = useState<string | null>(null);

  const [cancelModalOpen, setCancelModalOpen] = useState(false);
  // True when Cancel was reached through Revise: the bill is cancelled *and*
  // the operator is sent straight to a fresh bill, because the correction
  // here is "void this one, write a correct one" — not cancellation alone.
  const [cancelThenRevise, setCancelThenRevise] = useState(false);
  const [cancelReason, setCancelReason] = useState("");
  const [cancelError, setCancelError] = useState<string | null>(null);
  const [cancelling, setCancelling] = useState(false);

  const [creditModalOpen, setCreditModalOpen] = useState(false);
  const [creditReason, setCreditReason] = useState("");
  const [creditAmount, setCreditAmount] = useState("");
  const [creditError, setCreditError] = useState<string | null>(null);
  const [crediting, setCrediting] = useState(false);

  const fetchInvoice = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const inv = await api.getInvoice(orgId, params.id);
      setInvoice(inv);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, params.id, t.common.error]);

  useEffect(() => {
    fetchInvoice();
  }, [fetchInvoice]);

  // Printing must show the label the server assigned (first copy is
  // "original", every reprint after is "copy") — window.print() only fires
  // once that label has actually reached the document via props, which is
  // one render after the state update lands.
  useEffect(() => {
    if (printPending && copyLabel) {
      window.print();
      setPrintPending(false);
    }
  }, [printPending, copyLabel]);

  async function handlePrint() {
    if (!orgId || !invoice) return;
    setPrintError(null);
    try {
      const res = await api.printInvoice(orgId, invoice.id);
      setInvoice(res.invoice);
      setCopyLabel(res.copy_label);
      setPrintPending(true);
    } catch (err) {
      setPrintError(err instanceof ApiRequestError ? err.message : t.common.error);
    }
  }

  function openCancelModal(thenRevise: boolean) {
    setCancelThenRevise(thenRevise);
    setCancelReason("");
    setCancelError(null);
    setCancelModalOpen(true);
  }

  async function confirmCancel(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId || !invoice) return;
    setCancelling(true);
    setCancelError(null);
    try {
      const updated = await api.cancelInvoice(orgId, invoice.id, cancelReason.trim());
      setInvoice(updated);
      setCancelModalOpen(false);
      if (cancelThenRevise) {
        router.push(`${base}/invoices/new`);
      }
    } catch (err) {
      setCancelError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setCancelling(false);
    }
  }

  function openCreditModal() {
    setCreditReason("");
    setCreditAmount(invoice ? String(invoice.total) : "");
    setCreditError(null);
    setCreditModalOpen(true);
  }

  async function confirmCredit(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId || !invoice) return;
    const amount = Number(creditAmount);
    if (!Number.isFinite(amount) || amount <= 0) {
      setCreditError(t.common.invalidAmount);
      return;
    }
    setCrediting(true);
    setCreditError(null);
    try {
      const note = await api.creditNote(orgId, invoice.id, {
        reason: creditReason.trim(),
        items: [
          {
            description: `${t.invoices.creditNote} — ${invoice.invoice_number}`,
            quantity: 1,
            unit_price: amount,
          },
        ],
      });
      setCreditModalOpen(false);
      router.push(`${base}/invoices/${note.id}`);
    } catch (err) {
      setCreditError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setCrediting(false);
    }
  }

  const canCancelAndReissue = invoice?.status === "issued" && invoice.print_count === 0;
  const showBillActions = invoice?.status === "issued" && invoice.doc_type === "invoice";

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex items-center justify-between flex-wrap gap-3">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Receipt className="w-5 h-5 text-accent" />
          {t.invoices.title}
        </h1>
        {invoice && (
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handlePrint}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-surface border border-white/[0.06] text-fg text-sm font-medium hover:bg-surface-hover transition-colors"
            >
              <Printer className="w-4 h-4" />
              {t.invoices.print}
            </button>
            {showBillActions && (
              <button
                type="button"
                onClick={() => openCancelModal(false)}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm font-medium hover:bg-red-500/20 transition-colors"
              >
                <Ban className="w-4 h-4" />
                {t.invoices.cancel}
              </button>
            )}
          </div>
        )}
      </div>

      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}
      {printError && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {printError}
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : (
        invoice && (
          <>
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card p-2">
              <InvoiceDocument invoice={invoice} locale={locale} t={t} copyLabel={copyLabel} />
            </div>

            {showBillActions && (
              <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card p-5 space-y-3">
                <h2 className="text-sm font-semibold text-fg flex items-center gap-2">
                  <RefreshCw className="w-4 h-4 text-fg-muted" />
                  {t.invoices.revise}
                </h2>
                <p className="text-sm text-fg-muted">
                  {canCancelAndReissue ? t.invoices.reviseExplainCancel : t.invoices.reviseExplainCredit}
                </p>
                <button
                  type="button"
                  onClick={() => (canCancelAndReissue ? openCancelModal(true) : openCreditModal())}
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity"
                >
                  {t.invoices.revise}
                </button>
              </div>
            )}
          </>
        )
      )}

      {cancelModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black/60"
            onClick={() => setCancelModalOpen(false)}
            aria-hidden="true"
          />
          <form
            onSubmit={confirmCancel}
            className="relative w-full max-w-md bg-bg-elevated border border-white/[0.06] rounded-2xl p-5 shadow-card"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base font-semibold text-fg">{t.invoices.cancel}</h2>
              <button
                type="button"
                onClick={() => setCancelModalOpen(false)}
                className="p-1.5 rounded-lg text-fg-muted hover:bg-surface transition-colors"
                aria-label={t.common.cancel}
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <p className="mb-3 text-sm text-fg-muted">{t.invoices.cancelConfirm}</p>

            {cancelError && (
              <p className="mb-3 text-sm text-red-400" role="alert">
                {cancelError}
              </p>
            )}

            <div className="mb-5">
              <label htmlFor="cancel-reason" className="block text-xs text-fg-muted mb-1.5">
                {t.invoices.cancelReason}
              </label>
              <textarea
                id="cancel-reason"
                rows={3}
                required
                value={cancelReason}
                onChange={(e) => setCancelReason(e.target.value)}
                className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
              />
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => setCancelModalOpen(false)}
                className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
              >
                {t.common.cancel}
              </button>
              <button
                type="submit"
                disabled={cancelling || !cancelReason.trim()}
                className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
              >
                {cancelling && <Loader2 className="w-4 h-4 animate-spin" />}
                {t.common.confirm}
              </button>
            </div>
          </form>
        </div>
      )}

      {creditModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black/60"
            onClick={() => setCreditModalOpen(false)}
            aria-hidden="true"
          />
          <form
            onSubmit={confirmCredit}
            className="relative w-full max-w-md bg-bg-elevated border border-white/[0.06] rounded-2xl p-5 shadow-card"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base font-semibold text-fg">{t.invoices.creditNote}</h2>
              <button
                type="button"
                onClick={() => setCreditModalOpen(false)}
                className="p-1.5 rounded-lg text-fg-muted hover:bg-surface transition-colors"
                aria-label={t.common.cancel}
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {creditError && (
              <p className="mb-3 text-sm text-red-400" role="alert">
                {creditError}
              </p>
            )}

            <div className="mb-3">
              <label htmlFor="credit-reason" className="block text-xs text-fg-muted mb-1.5">
                {t.invoices.creditReason}
              </label>
              <textarea
                id="credit-reason"
                rows={3}
                required
                value={creditReason}
                onChange={(e) => setCreditReason(e.target.value)}
                className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
              />
            </div>

            <div className="mb-5">
              <label htmlFor="credit-amount" className="block text-xs text-fg-muted mb-1.5">
                {t.invoices.creditAmount}
              </label>
              <input
                id="credit-amount"
                type="number"
                min="0.01"
                step="0.01"
                value={creditAmount}
                onChange={(e) => setCreditAmount(e.target.value)}
                className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
              />
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => setCreditModalOpen(false)}
                className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
              >
                {t.common.cancel}
              </button>
              <button
                type="submit"
                disabled={crediting || !creditReason.trim()}
                className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
              >
                {crediting && <Loader2 className="w-4 h-4 animate-spin" />}
                {t.common.confirm}
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
