"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type {
  Locale,
  OrgPackage,
  CreatePackageRequest,
  UpdatePackageRequest,
} from "@/types";
import {
  Package,
  Plus,
  X,
  Pencil,
  Trash2,
  Loader2,
  Check,
  XCircle,
} from "lucide-react";

type Filter = "all" | "active" | "inactive";

function formatNPR(amount: string | number): string {
  const num = typeof amount === "string" ? parseFloat(amount) : amount;
  return `Rs. ${num.toLocaleString("en-IN")}`;
}

export default function PackagesPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [packages, setPackages] = useState<OrgPackage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<Filter>("all");

  // Form state
  const [showModal, setShowModal] = useState(false);
  const [editingPkg, setEditingPkg] = useState<OrgPackage | null>(null);
  const [formName, setFormName] = useState("");
  const [formDescription, setFormDescription] = useState("");
  const [formPrice, setFormPrice] = useState("");
  const [formDuration, setFormDuration] = useState("");
  const [formFeatures, setFormFeatures] = useState("");
  const [formMaxMembers, setFormMaxMembers] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchPackages = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listPackages(orgId, { limit: 100 });
      setPackages(res.data ?? []);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchPackages();
  }, [fetchPackages]);

  const filtered =
    filter === "all"
      ? packages
      : packages.filter((p) =>
          filter === "active" ? p.is_active : !p.is_active
        );

  const openCreateModal = () => {
    setEditingPkg(null);
    setFormName("");
    setFormDescription("");
    setFormPrice("");
    setFormDuration("");
    setFormFeatures("");
    setFormMaxMembers("");
    setFormError(null);
    setShowModal(true);
  };

  const openEditModal = (pkg: OrgPackage) => {
    setEditingPkg(pkg);
    setFormName(pkg.name);
    setFormDescription(pkg.description ?? "");
    setFormPrice(pkg.price);
    setFormDuration(String(pkg.duration_days));
    setFormFeatures(pkg.features?.join(", ") ?? "");
    setFormMaxMembers(pkg.max_members ? String(pkg.max_members) : "");
    setFormError(null);
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setFormLoading(true);
    setFormError(null);

    const features = formFeatures
      .split(",")
      .map((f) => f.trim())
      .filter(Boolean);

    try {
      if (editingPkg) {
        const data: UpdatePackageRequest = {
          name: formName,
          description: formDescription || undefined,
          price: formPrice,
          duration_days: parseInt(formDuration),
          features: features.length > 0 ? features : undefined,
          max_members: formMaxMembers
            ? parseInt(formMaxMembers)
            : undefined,
        };
        await api.updatePackage(orgId, editingPkg.id, data);
      } else {
        const data: CreatePackageRequest = {
          name: formName,
          description: formDescription || undefined,
          price: formPrice,
          duration_days: parseInt(formDuration),
          features: features.length > 0 ? features : undefined,
          max_members: formMaxMembers
            ? parseInt(formMaxMembers)
            : undefined,
        };
        await api.createPackage(orgId, data);
      }
      setShowModal(false);
      fetchPackages();
    } catch (err) {
      setFormError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (pkgId: string) => {
    if (!orgId) return;
    try {
      await api.deletePackage(orgId, pkgId);
      fetchPackages();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const handleToggleActive = async (pkg: OrgPackage) => {
    if (!orgId) return;
    try {
      await api.updatePackage(orgId, pkg.id, {
        is_active: !pkg.is_active,
      });
      fetchPackages();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <Package className="w-5 h-5 text-accent" />
          {t.packages.title}
        </h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.packages.addPackage}
        </button>
      </div>

      {/* Filters */}
      <div className="flex gap-2">
        {(["all", "active", "inactive"] as Filter[]).map((f) => {
          const label =
            f === "all"
              ? t.packages.filterAll
              : f === "active"
                ? t.packages.filterActive
                : t.packages.filterInactive;
          return (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className={`px-3 py-1.5 text-xs font-mono uppercase rounded-lg border transition-colors ${
                filter === f
                  ? "bg-accent/10 text-accent border-accent/20"
                  : "text-fg-muted border-white/[0.06] hover:bg-surface"
              }`}
            >
              {label}
            </button>
          );
        })}
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Package Cards */}
      {loading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-16 text-fg-muted text-sm">
          {t.packages.noPackages}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((pkg) => (
            <div
              key={pkg.id}
              className={`bg-gradient-to-b from-white/[0.08] to-white/[0.02] border rounded-2xl p-5 shadow-card hover:shadow-card-hover transition-shadow duration-300 ${
                pkg.is_active
                  ? "border-white/[0.06]"
                  : "border-red-500/20 opacity-60"
              }`}
            >
              <div className="flex items-start justify-between mb-3">
                <div>
                  <h3 className="text-base font-semibold text-fg">
                    {pkg.name}
                  </h3>
                  {pkg.description && (
                    <p className="text-xs text-fg-muted mt-0.5">
                      {pkg.description}
                    </p>
                  )}
                </div>
                <span
                  className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
                    pkg.is_active
                      ? "bg-accent/10 text-accent border-accent/20"
                      : "bg-red-500/10 text-red-400 border-red-500/20"
                  }`}
                >
                  {pkg.is_active ? t.packages.active : t.packages.inactive}
                </span>
              </div>

              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-fg-muted">{t.packages.price}</span>
                  <span className="text-fg font-mono font-medium">
                    {formatNPR(pkg.price)}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-fg-muted">
                    {t.packages.duration}
                  </span>
                  <span className="text-fg font-mono">
                    {pkg.duration_days} {t.packages.days}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-fg-muted">
                    {t.packages.maxMembers}
                  </span>
                  <span className="text-fg font-mono">
                    {pkg.max_members ?? t.packages.unlimited}
                  </span>
                </div>
              </div>

              {pkg.features && pkg.features.length > 0 && (
                <div className="flex flex-wrap gap-1 mb-4">
                  {pkg.features.map((f) => (
                    <span
                      key={f}
                      className="px-2 py-0.5 bg-surface border border-white/[0.06] rounded-lg text-[11px] text-fg-muted font-mono"
                    >
                      {f}
                    </span>
                  ))}
                </div>
              )}

              <div className="flex items-center gap-1 pt-3 border-t border-white/[0.06]">
                <button
                  onClick={() => openEditModal(pkg)}
                  className="p-1.5 rounded-lg hover:bg-surface text-fg-muted hover:text-fg transition-colors"
                  title={t.packages.editPackage}
                >
                  <Pencil className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleToggleActive(pkg)}
                  className={`p-1.5 rounded-lg transition-colors ${
                    pkg.is_active
                      ? "hover:bg-red-500/10 text-fg-muted hover:text-red-400"
                      : "hover:bg-accent/10 text-fg-muted hover:text-accent"
                  }`}
                  title={
                    pkg.is_active
                      ? t.packages.inactive
                      : t.packages.active
                  }
                >
                  {pkg.is_active ? (
                    <XCircle className="w-4 h-4" />
                  ) : (
                    <Check className="w-4 h-4" />
                  )}
                </button>
                <button
                  onClick={() => handleDelete(pkg.id)}
                  className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
                  title={t.common.delete}
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {editingPkg
                  ? t.packages.editPackage
                  : t.packages.addPackageTitle}
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
                  {t.packages.name} *
                </label>
                <input
                  type="text"
                  required
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  placeholder={t.packages.namePlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.packages.description}
                </label>
                <input
                  type="text"
                  value={formDescription}
                  onChange={(e) => setFormDescription(e.target.value)}
                  placeholder={t.packages.descriptionPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs text-fg-muted mb-1.5">
                    {t.packages.price} ({t.packages.currency}) *
                  </label>
                  <input
                    type="number"
                    required
                    min="0"
                    step="0.01"
                    value={formPrice}
                    onChange={(e) => setFormPrice(e.target.value)}
                    placeholder={t.packages.pricePlaceholder}
                    className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                  />
                </div>
                <div>
                  <label className="block text-xs text-fg-muted mb-1.5">
                    {t.packages.duration} ({t.packages.days}) *
                  </label>
                  <input
                    type="number"
                    required
                    min="1"
                    value={formDuration}
                    onChange={(e) => setFormDuration(e.target.value)}
                    placeholder={t.packages.durationPlaceholder}
                    className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.packages.maxMembers}
                </label>
                <input
                  type="number"
                  min="1"
                  value={formMaxMembers}
                  onChange={(e) => setFormMaxMembers(e.target.value)}
                  placeholder={t.packages.unlimited}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.packages.features}
                </label>
                <input
                  type="text"
                  value={formFeatures}
                  onChange={(e) => setFormFeatures(e.target.value)}
                  placeholder={t.packages.featuresPlaceholder}
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
