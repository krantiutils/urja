"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, Transaction, AccountsSummary } from "@/types";
import {
  BookOpen,
  Plus,
  X,
  Pencil,
  Trash2,
  Loader2,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Percent,
} from "lucide-react";

type TransactionType = "income" | "expense";
type PaymentType = "cash" | "esewa" | "bank";
type FilterType = "all" | "income" | "expense";

function TypeBadge({ type, label }: { type: string; label: string }) {
  const colors: Record<string, string> = {
    income: "bg-accent/10 text-accent border-accent/20",
    expense: "bg-red-500/10 text-red-400 border-red-500/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[type] ?? colors.income
      }`}
    >
      {label}
    </span>
  );
}

export default function AccountsPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterType, setFilterType] = useState<FilterType>("all");

  // Summary
  const [summary, setSummary] = useState<AccountsSummary | null>(null);
  const [summaryLoading, setSummaryLoading] = useState(true);

  // Modal state
  const [showModal, setShowModal] = useState(false);
  const [editingTransaction, setEditingTransaction] =
    useState<Transaction | null>(null);
  const [formCategory, setFormCategory] = useState("");
  const [formDescription, setFormDescription] = useState("");
  const [formDate, setFormDate] = useState("");
  const [formType, setFormType] = useState<TransactionType>("income");
  const [formAmount, setFormAmount] = useState("");
  const [formPaymentType, setFormPaymentType] = useState<PaymentType>("cash");
  const [formReference, setFormReference] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchTransactions = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listTransactions(orgId, {
        transaction_type: filterType === "all" ? undefined : filterType,
        limit: 100,
      });
      setTransactions(res.data ?? []);
      setTotal(res.total ?? 0);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, filterType, t.common.error]);

  const fetchSummary = useCallback(async () => {
    if (!orgId) return;
    setSummaryLoading(true);
    try {
      const res = await api.getAccountsSummary(orgId);
      setSummary(res);
    } catch {
      // Summary fetch failure is non-blocking
    } finally {
      setSummaryLoading(false);
    }
  }, [orgId]);

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  useEffect(() => {
    fetchSummary();
  }, [fetchSummary]);

  const openCreateModal = () => {
    setEditingTransaction(null);
    setFormCategory("");
    setFormDescription("");
    setFormDate(new Date().toISOString().split("T")[0]);
    setFormType("income");
    setFormAmount("");
    setFormPaymentType("cash");
    setFormReference("");
    setFormError(null);
    setShowModal(true);
  };

  const openEditModal = (tx: Transaction) => {
    setEditingTransaction(tx);
    setFormCategory(tx.category);
    setFormDescription(tx.description);
    setFormDate(tx.transaction_date.split("T")[0]);
    setFormType(tx.transaction_type);
    setFormAmount(tx.amount);
    setFormPaymentType(tx.payment_type as PaymentType);
    setFormReference(tx.reference ?? "");
    setFormError(null);
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setFormLoading(true);
    setFormError(null);

    try {
      if (editingTransaction) {
        await api.updateTransaction(orgId, editingTransaction.id, {
          category: formCategory,
          description: formDescription,
          transaction_date: formDate,
          transaction_type: formType,
          amount: formAmount,
          payment_type: formPaymentType,
          reference: formReference || undefined,
        });
      } else {
        await api.createTransaction(orgId, {
          category: formCategory,
          description: formDescription,
          transaction_date: formDate,
          transaction_type: formType,
          amount: formAmount,
          payment_type: formPaymentType,
          reference: formReference || undefined,
        });
      }
      setShowModal(false);
      fetchTransactions();
      fetchSummary();
    } catch (err) {
      setFormError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!orgId) return;
    if (!confirm(t.accounts.confirmDelete)) return;
    try {
      await api.deleteTransaction(orgId, id);
      fetchTransactions();
      fetchSummary();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const typeLabel = (type: string) => {
    const map: Record<string, string> = {
      income: t.accounts.income,
      expense: t.accounts.expense,
    };
    return map[type] ?? type;
  };

  const filterTabs: { key: FilterType; label: string }[] = [
    { key: "all", label: t.accounts.all },
    { key: "income", label: t.accounts.income },
    { key: "expense", label: t.accounts.expense },
  ];

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <BookOpen className="w-5 h-5 text-accent" />
          {t.accounts.title}
        </h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.accounts.addTransaction}
        </button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-4 h-4 text-accent" />
            <span className="text-xs text-fg-muted font-mono uppercase tracking-widest">
              {t.accounts.totalIncome}
            </span>
          </div>
          <p className="text-2xl font-bold text-accent">
            {summaryLoading ? "..." : summary?.total_income ?? "0"}
          </p>
        </div>
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
          <div className="flex items-center gap-2 mb-2">
            <TrendingDown className="w-4 h-4 text-red-400" />
            <span className="text-xs text-fg-muted font-mono uppercase tracking-widest">
              {t.accounts.totalExpenses}
            </span>
          </div>
          <p className="text-2xl font-bold text-red-400">
            {summaryLoading ? "..." : summary?.total_expenses ?? "0"}
          </p>
        </div>
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
          <div className="flex items-center gap-2 mb-2">
            <DollarSign className="w-4 h-4 text-fg-muted" />
            <span className="text-xs text-fg-muted font-mono uppercase tracking-widest">
              {t.accounts.grossProfit}
            </span>
          </div>
          <p className="text-2xl font-bold text-fg">
            {summaryLoading ? "..." : summary?.gross_profit ?? "0"}
          </p>
        </div>
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-4 shadow-card">
          <div className="flex items-center gap-2 mb-2">
            <Percent className="w-4 h-4 text-fg-muted" />
            <span className="text-xs text-fg-muted font-mono uppercase tracking-widest">
              {t.accounts.profitPercent}
            </span>
          </div>
          <p className="text-2xl font-bold text-fg">
            {summaryLoading
              ? "..."
              : summary?.profit_percent != null
                ? `${summary.profit_percent.toFixed(1)}%`
                : "0%"}
          </p>
        </div>
      </div>

      {/* Filter Tabs */}
      <div className="flex items-center gap-1">
        {filterTabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setFilterType(tab.key)}
            className={`px-4 py-1.5 text-sm rounded-lg font-medium transition-colors ${
              filterType === tab.key
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
        ) : transactions.length === 0 ? (
          <div className="text-center py-16 text-fg-muted text-sm">
            {t.accounts.noTransactions}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.date}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.category}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.description}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.type}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.amount}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.paymentType}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.reference}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.entryBy}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                    {t.accounts.actions}
                  </th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr
                    key={tx.id}
                    className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                  >
                    <td className="px-5 py-3 text-fg-muted font-mono text-xs whitespace-nowrap">
                      {new Date(tx.transaction_date).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-3 text-fg font-medium">
                      {tx.category}
                    </td>
                    <td className="px-5 py-3 text-fg-muted text-xs max-w-[200px] truncate">
                      {tx.description}
                    </td>
                    <td className="px-5 py-3">
                      <TypeBadge
                        type={tx.transaction_type}
                        label={typeLabel(tx.transaction_type)}
                      />
                    </td>
                    <td
                      className={`px-5 py-3 text-right font-mono font-medium ${
                        tx.transaction_type === "income"
                          ? "text-accent"
                          : "text-red-400"
                      }`}
                    >
                      {tx.transaction_type === "income" ? "+" : "-"}
                      {tx.amount}
                    </td>
                    <td className="px-5 py-3 text-fg-muted text-xs uppercase">
                      {tx.payment_type}
                    </td>
                    <td className="px-5 py-3 text-fg-muted text-xs">
                      {tx.reference || "-"}
                    </td>
                    <td className="px-5 py-3 text-fg-muted text-xs">
                      {tx.entry_by_name || "-"}
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center justify-end gap-1">
                        <button
                          onClick={() => openEditModal(tx)}
                          title={t.accounts.editTitle}
                          className="p-1.5 rounded-lg hover:bg-surface text-fg-muted hover:text-fg transition-colors"
                        >
                          <Pencil className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => handleDelete(tx.id)}
                          title={t.common.delete}
                          className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
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
              {total} {t.accounts.title.toLowerCase()}
            </p>
          </div>
        )}
      </div>

      {/* Add/Edit Transaction Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {editingTransaction
                  ? t.accounts.editTitle
                  : t.accounts.addTitle}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-5 space-y-4">
              {formError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {formError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.category} *
                </label>
                <input
                  type="text"
                  required
                  value={formCategory}
                  onChange={(e) => setFormCategory(e.target.value)}
                  placeholder={t.accounts.categoryPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.description} *
                </label>
                <input
                  type="text"
                  required
                  value={formDescription}
                  onChange={(e) => setFormDescription(e.target.value)}
                  placeholder={t.accounts.descriptionPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.date} *
                </label>
                <input
                  type="date"
                  required
                  value={formDate}
                  onChange={(e) => setFormDate(e.target.value)}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.type} *
                </label>
                <select
                  value={formType}
                  onChange={(e) =>
                    setFormType(e.target.value as TransactionType)
                  }
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50 transition-colors"
                >
                  <option value="income">{t.accounts.income}</option>
                  <option value="expense">{t.accounts.expense}</option>
                </select>
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.amount} *
                </label>
                <input
                  type="text"
                  required
                  value={formAmount}
                  onChange={(e) => setFormAmount(e.target.value)}
                  placeholder={t.accounts.amountPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.paymentType} *
                </label>
                <select
                  value={formPaymentType}
                  onChange={(e) =>
                    setFormPaymentType(e.target.value as PaymentType)
                  }
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg focus:outline-none focus:border-accent/50 transition-colors"
                >
                  <option value="cash">Cash</option>
                  <option value="esewa">eSewa</option>
                  <option value="bank">Bank</option>
                </select>
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.accounts.reference}
                </label>
                <input
                  type="text"
                  value={formReference}
                  onChange={(e) => setFormReference(e.target.value)}
                  placeholder={t.accounts.referencePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
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
    </div>
  );
}
