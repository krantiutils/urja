"use client";

import type { Dictionary } from "@/lib/i18n";
import type { Invoice, Locale } from "@/types";

/**
 * The printed bill.
 *
 * Every value is read from the invoice's own snapshot rather than from the
 * org or the member: a bill reprinted next year must show what it showed the
 * day it was issued, even if the gym has since changed its name or PAN.
 */
export function InvoiceDocument({
  invoice,
  locale,
  t,
  copyLabel,
}: {
  invoice: Invoice;
  locale: Locale;
  t: Dictionary;
  copyLabel?: "original" | "copy";
}) {
  const isCredit = invoice.doc_type === "credit_note";
  // numberingSystem is pinned to "latn": Node's ICU defaults "ne-NP" to
  // Devanagari digits while not every browser ships that locale's data, so
  // leaving it unset renders differently server- vs client-side and trips a
  // hydration mismatch. Pinning it also keeps money in the same digits as
  // invoice_number and issued_date_bs, which are plain strings and always
  // Arabic numerals regardless of locale — a bill mixing digit scripts within
  // itself would look like a mistake, not a translation.
  const money = (n: number) =>
    n.toLocaleString(locale === "ne" ? "ne-NP" : "en-IN", {
      numberingSystem: "latn",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

  return (
    <article className="invoice-doc bg-white text-black p-8 max-w-[210mm] mx-auto">
      <header className="text-center border-b border-black/20 pb-4">
        <h1 className="text-lg font-bold">{invoice.seller_name}</h1>
        {invoice.seller_address && <p className="text-sm">{invoice.seller_address}</p>}
        <p className="text-sm font-mono">PAN: {invoice.seller_pan}</p>
        <p className="mt-2 text-sm font-semibold uppercase tracking-widest">
          {isCredit ? t.invoices.creditNote : t.invoices.title}
          {copyLabel && (
            <span className="ml-2 font-normal">
              ({copyLabel === "original" ? t.invoices.original : t.invoices.copy})
            </span>
          )}
        </p>
      </header>

      <section className="grid grid-cols-2 gap-4 py-4 text-sm">
        <div>
          <p>
            <span className="text-black/60">{t.invoices.customer}: </span>
            {invoice.customer_name}
          </p>
          {invoice.customer_pan && (
            <p>
              <span className="text-black/60">PAN: </span>
              {invoice.customer_pan}
            </p>
          )}
          {invoice.customer_address && <p>{invoice.customer_address}</p>}
        </div>
        <div className="text-right">
          <p className="font-mono">{invoice.invoice_number}</p>
          <p>{invoice.issued_date_bs} (BS)</p>
          <p className="text-black/60">{invoice.issued_date} (AD)</p>
        </div>
      </section>

      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-y border-black/20">
            <th className="text-left py-1 w-8">#</th>
            <th className="text-left py-1">{t.invoices.description}</th>
            <th className="text-right py-1 w-16">{t.invoices.quantity}</th>
            <th className="text-right py-1 w-24">{t.invoices.unitPrice}</th>
            <th className="text-right py-1 w-28">{t.invoices.lineTotal}</th>
          </tr>
        </thead>
        <tbody>
          {(invoice.items ?? []).map((it) => (
            <tr key={it.line_no} className="border-b border-black/10">
              <td className="py-1">{it.line_no}</td>
              <td className="py-1">
                {locale === "ne" && it.description_ne ? it.description_ne : it.description}
              </td>
              <td className="py-1 text-right tabular-nums">{it.quantity}</td>
              <td className="py-1 text-right tabular-nums">{money(it.unit_price)}</td>
              <td className="py-1 text-right tabular-nums">{money(it.amount)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <section className="mt-4 ml-auto w-64 text-sm">
        <Row label={t.invoices.subtotal} value={money(invoice.subtotal)} />
        {invoice.discount > 0 && (
          <Row label={t.invoices.discount} value={`- ${money(invoice.discount)}`} />
        )}
        {/* VAT is deferred; the row appears only if a bill ever carries one. */}
        {invoice.vat_amount > 0 && (
          <Row label={`VAT ${invoice.vat_rate}%`} value={money(invoice.vat_amount)} />
        )}
        <div className="flex justify-between border-t border-black/30 pt-1 font-semibold">
          <span>{t.invoices.total}</span>
          <span className="tabular-nums">Rs {money(invoice.total)}</span>
        </div>
      </section>

      <p className="mt-3 text-sm">
        <span className="text-black/60">{t.invoices.amountInWords}: </span>
        {invoice.amount_in_words}
      </p>

      {invoice.status === "cancelled" && (
        <p className="mt-4 border border-black px-3 py-2 text-sm font-semibold uppercase">
          {t.invoices.cancelled}
          {invoice.cancellation_reason ? ` — ${invoice.cancellation_reason}` : ""}
        </p>
      )}
    </article>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between">
      <span className="text-black/60">{label}</span>
      <span className="tabular-nums">{value}</span>
    </div>
  );
}
