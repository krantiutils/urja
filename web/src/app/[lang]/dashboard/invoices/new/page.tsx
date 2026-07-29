"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Plus, Receipt, Trash2, X } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { InvoiceItemInput, IssueInvoiceInput, Locale, OrgMember } from "@/types";

type CustomerMode = "walkIn" | "existing";

interface ItemRow {
  description: string;
  quantity: string;
  unit_price: string;
}

function emptyRow(): ItemRow {
  return { description: "", quantity: "1", unit_price: "" };
}

const labelClass = "block text-xs text-fg-muted mb-1.5";
const inputClass =
  "w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors";

/**
 * Issue a bill.
 *
 * Totals shown here are computed client-side purely for feedback while
 * typing — the server recomputes subtotal/taxable/total from the items and
 * discount it actually receives and is the only authority on the number that
 * ends up on the printed document. Nothing computed here is sent; only the
 * raw items and discount are.
 */
export default function NewInvoicePage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;
  const router = useRouter();
  const base = `/${locale}/dashboard`;

  const [mode, setMode] = useState<CustomerMode>("walkIn");
  const [customerUserId, setCustomerUserId] = useState<string | undefined>();
  const [customerName, setCustomerName] = useState("");
  const [customerPhone, setCustomerPhone] = useState<string | undefined>();
  const [customerPan, setCustomerPan] = useState("");
  const [customerAddress, setCustomerAddress] = useState("");
  const [paymentMethod, setPaymentMethod] = useState("");
  const [discount, setDiscount] = useState("");
  const [items, setItems] = useState<ItemRow[]>([emptyRow()]);

  const [members, setMembers] = useState<OrgMember[]>([]);
  const [membersLoaded, setMembersLoaded] = useState(false);
  const [memberSearch, setMemberSearch] = useState("");

  const [previewNumber, setPreviewNumber] = useState<string | null>(null);
  const [panMissing, setPanMissing] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // A preview, not a reservation: fetched purely to give the operator a hint
  // of what's coming next. The server assigns the real number at issue time,
  // so a bill issued by someone else a moment later does not desync this.
  useEffect(() => {
    if (!orgId) return;
    api
      .nextInvoiceNumber(orgId)
      .then((res) => setPreviewNumber(res.invoice_number))
      .catch((err) => {
        if (err instanceof ApiRequestError && err.code === "pan_not_configured") {
          setPanMissing(true);
        }
      });
  }, [orgId]);

  const switchToExisting = useCallback(async () => {
    setMode("existing");
    setCustomerUserId(undefined);
    setCustomerName("");
    setCustomerPhone(undefined);
    if (!orgId || membersLoaded) return;
    try {
      const res = await api.listMembers(orgId, { limit: 500 });
      setMembers(res.data ?? []);
    } catch {
      // Non-blocking: the picker just shows no results to search.
    } finally {
      setMembersLoaded(true);
    }
  }, [orgId, membersLoaded]);

  const switchToWalkIn = useCallback(() => {
    setMode("walkIn");
    setCustomerUserId(undefined);
    setCustomerName("");
    setCustomerPhone(undefined);
  }, []);

  function pickMember(m: OrgMember) {
    setCustomerUserId(m.id);
    setCustomerName(m.name);
    setCustomerPhone(m.phone);
  }

  function clearMember() {
    setCustomerUserId(undefined);
    setCustomerName("");
    setCustomerPhone(undefined);
  }

  function updateItem(idx: number, field: keyof ItemRow, value: string) {
    setItems((prev) => prev.map((row, i) => (i === idx ? { ...row, [field]: value } : row)));
  }

  function addLine() {
    setItems((prev) => [...prev, emptyRow()]);
  }

  function removeLine(idx: number) {
    setItems((prev) => (prev.length === 1 ? prev : prev.filter((_, i) => i !== idx)));
  }

  const lineAmount = (row: ItemRow) => (Number(row.quantity) || 0) * (Number(row.unit_price) || 0);
  const subtotal = useMemo(() => items.reduce((sum, row) => sum + lineAmount(row), 0), [items]);
  const discountNum = Number(discount) || 0;
  const total = Math.max(subtotal - discountNum, 0);

  const filteredMembers = memberSearch
    ? members.filter((m) =>
        `${m.name} ${m.phone}`.toLowerCase().includes(memberSearch.toLowerCase())
      )
    : members;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!orgId) return;
    setFormError(null);

    const lineItems: InvoiceItemInput[] = items
      .filter((row) => row.description.trim() && Number(row.unit_price) > 0)
      .map((row) => ({
        description: row.description.trim(),
        quantity: Number(row.quantity) || 1,
        unit_price: Number(row.unit_price) || 0,
      }));

    if (lineItems.length === 0) {
      setFormError(t.common.invalidAmount);
      return;
    }

    setSubmitting(true);
    try {
      const payload: IssueInvoiceInput = {
        customer_user_id: mode === "existing" ? customerUserId : undefined,
        customer_name: customerName.trim(),
        customer_pan: customerPan.trim() || undefined,
        customer_address: customerAddress.trim() || undefined,
        customer_phone: customerPhone || undefined,
        payment_method: paymentMethod || undefined,
        discount: discountNum || undefined,
        items: lineItems,
      };
      const invoice = await api.issueInvoice(orgId, payload);
      router.push(`${base}/invoices/${invoice.id}`);
    } catch (err) {
      if (err instanceof ApiRequestError && err.code === "pan_not_configured") {
        setPanMissing(true);
      } else {
        setFormError(err instanceof ApiRequestError ? err.message : t.common.error);
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="space-y-6 animate-fade-in max-w-3xl">
      <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
        <Receipt className="w-5 h-5 text-accent" />
        {t.invoices.newBill}
      </h1>

      {panMissing && (
        <div className="p-4 bg-amber-500/10 border border-amber-500/20 rounded-xl text-sm text-amber-300 flex flex-wrap items-center justify-between gap-3">
          <span>{t.invoices.panMissing}</span>
          <a
            href={`${base}/settings`}
            className="px-3 py-1.5 rounded-lg bg-amber-500/20 text-amber-200 font-medium hover:bg-amber-500/30 transition-colors shrink-0"
          >
            {t.invoices.panMissingCta}
          </a>
        </div>
      )}

      {previewNumber && (
        <div className="p-3 rounded-xl border border-white/[0.06] bg-surface text-xs">
          <span className="text-fg-muted">{t.invoices.numberPreview}: </span>
          <span className="font-mono text-fg">{previewNumber}</span>
          <p className="mt-1 text-fg-muted">{t.invoices.numberPreviewHint}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5">
        {formError && (
          <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400" role="alert">
            {formError}
          </div>
        )}

        <div className="inline-flex rounded-xl border border-white/[0.06] p-1 bg-surface">
          <button
            type="button"
            onClick={switchToWalkIn}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              mode === "walkIn" ? "bg-accent/10 text-accent font-medium" : "text-fg-muted hover:text-fg"
            }`}
          >
            {t.invoices.walkIn}
          </button>
          <button
            type="button"
            onClick={switchToExisting}
            className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
              mode === "existing" ? "bg-accent/10 text-accent font-medium" : "text-fg-muted hover:text-fg"
            }`}
          >
            {t.invoices.existingMember}
          </button>
        </div>

        {mode === "walkIn" ? (
          <div>
            <label htmlFor="inv-customer-name" className={labelClass}>
              {t.invoices.customerName}
            </label>
            <input
              id="inv-customer-name"
              type="text"
              required
              value={customerName}
              onChange={(e) => setCustomerName(e.target.value)}
              className={inputClass}
            />
          </div>
        ) : customerUserId ? (
          <div className="flex items-center justify-between px-3 py-2.5 bg-input-bg border border-accent/30 rounded-xl text-sm text-fg">
            <span>
              {customerName}
              {customerPhone && (
                <span className="ml-2 font-mono text-xs text-fg-muted">{customerPhone}</span>
              )}
            </span>
            <button
              type="button"
              onClick={clearMember}
              className="p-0.5 rounded hover:bg-surface text-fg-muted"
              aria-label={t.common.cancel}
            >
              <X className="w-3.5 h-3.5" />
            </button>
          </div>
        ) : (
          <div>
            <label htmlFor="inv-member-search" className={labelClass}>
              {t.invoices.existingMember}
            </label>
            <input
              id="inv-member-search"
              type="text"
              value={memberSearch}
              onChange={(e) => setMemberSearch(e.target.value)}
              placeholder={t.dues.searchMember}
              className={inputClass}
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
                    onClick={() => pickMember(m)}
                    className="w-full text-left px-3 py-2.5 text-sm hover:bg-surface transition-colors border-b border-white/[0.03] last:border-0"
                  >
                    <span className="text-fg">{m.name}</span>
                    <span className="ml-2 font-mono text-xs text-fg-muted">{m.phone}</span>
                  </button>
                ))
              )}
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label htmlFor="inv-customer-pan" className={labelClass}>
              {t.invoices.customerPan}
            </label>
            <input
              id="inv-customer-pan"
              type="text"
              inputMode="numeric"
              value={customerPan}
              onChange={(e) => setCustomerPan(e.target.value)}
              className={inputClass}
            />
          </div>
          <div>
            <label htmlFor="inv-customer-address" className={labelClass}>
              {t.invoices.customerAddress}
            </label>
            <input
              id="inv-customer-address"
              type="text"
              value={customerAddress}
              onChange={(e) => setCustomerAddress(e.target.value)}
              className={inputClass}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label htmlFor="inv-payment-method" className={labelClass}>
              {t.invoices.paymentMethod}
            </label>
            <select
              id="inv-payment-method"
              value={paymentMethod}
              onChange={(e) => setPaymentMethod(e.target.value)}
              className={inputClass}
            >
              <option value="">—</option>
              <option value="cash">{t.dues.cash}</option>
              <option value="esewa">{t.dues.esewa}</option>
              <option value="bank">{t.dues.bank}</option>
            </select>
          </div>
          <div>
            <label htmlFor="inv-discount" className={labelClass}>
              {t.invoices.discount}
            </label>
            <input
              id="inv-discount"
              type="number"
              min="0"
              step="0.01"
              value={discount}
              onChange={(e) => setDiscount(e.target.value)}
              placeholder="0.00"
              className={inputClass}
            />
          </div>
        </div>

        <div className="space-y-3">
          {items.map((row, idx) => (
            <div
              key={idx}
              className="grid grid-cols-12 gap-2 items-end p-3 rounded-xl bg-surface border border-white/[0.06]"
            >
              <div className="col-span-12 sm:col-span-5">
                <label htmlFor={`item-desc-${idx}`} className={labelClass}>
                  {t.invoices.description}
                </label>
                <input
                  id={`item-desc-${idx}`}
                  type="text"
                  required
                  value={row.description}
                  onChange={(e) => updateItem(idx, "description", e.target.value)}
                  className={inputClass}
                />
              </div>
              <div className="col-span-4 sm:col-span-2">
                <label htmlFor={`item-qty-${idx}`} className={labelClass}>
                  {t.invoices.quantity}
                </label>
                <input
                  id={`item-qty-${idx}`}
                  type="number"
                  min="1"
                  step="1"
                  value={row.quantity}
                  onChange={(e) => updateItem(idx, "quantity", e.target.value)}
                  className={inputClass}
                />
              </div>
              <div className="col-span-4 sm:col-span-2">
                <label htmlFor={`item-rate-${idx}`} className={labelClass}>
                  {t.invoices.unitPrice}
                </label>
                <input
                  id={`item-rate-${idx}`}
                  type="number"
                  min="0"
                  step="0.01"
                  value={row.unit_price}
                  onChange={(e) => updateItem(idx, "unit_price", e.target.value)}
                  placeholder="0.00"
                  className={inputClass}
                />
              </div>
              <div className="col-span-3 sm:col-span-2">
                <span className={labelClass}>{t.invoices.lineTotal}</span>
                <p className="px-3 py-2.5 text-sm text-fg font-mono">{lineAmount(row).toFixed(2)}</p>
              </div>
              <div className="col-span-1 flex justify-end">
                <button
                  type="button"
                  onClick={() => removeLine(idx)}
                  disabled={items.length === 1}
                  aria-label={t.invoices.removeLine}
                  className="p-1.5 rounded-lg text-fg-muted hover:text-red-400 hover:bg-surface-hover transition-colors disabled:opacity-30"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}

          <button
            type="button"
            onClick={addLine}
            className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm text-accent hover:bg-accent/10 transition-colors"
          >
            <Plus className="w-4 h-4" />
            {t.invoices.addLine}
          </button>
        </div>

        <div className="ml-auto w-full sm:w-64 space-y-1 text-sm">
          <div className="flex justify-between">
            <span className="text-fg-muted">{t.invoices.subtotal}</span>
            <span className="font-mono text-fg">{subtotal.toFixed(2)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-fg-muted">{t.invoices.discount}</span>
            <span className="font-mono text-fg">{discountNum.toFixed(2)}</span>
          </div>
          <div className="flex justify-between border-t border-white/[0.06] pt-1 font-semibold">
            <span className="text-fg">{t.invoices.total}</span>
            <span className="font-mono text-accent">{total.toFixed(2)}</span>
          </div>
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-bg-base text-sm font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {submitting && <Loader2 className="w-4 h-4 animate-spin" />}
          {t.invoices.issueBill}
        </button>
      </form>
    </div>
  );
}
