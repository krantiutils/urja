"use client";

import { useState, useEffect, useCallback } from "react";
import { api, ApiRequestError } from "@/lib/api";
import type { Organization } from "@/types";
import {
  Building2,
  Plus,
  Loader2,
  X,
  MapPin,
  Phone,
  Mail,
} from "lucide-react";

export default function SuperAdminPage() {
  const [orgs, setOrgs] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showCreate, setShowCreate] = useState(false);

  const fetchOrgs = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.listOrgs({ limit: 50 });
      setOrgs(res.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : "Failed to load organizations");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchOrgs();
  }, [fetchOrgs]);

  return (
    <div className="max-w-5xl mx-auto">
      {/* Page header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-fg">Organizations</h1>
          <p className="text-sm text-fg-muted mt-1">
            Manage all gyms on the platform
          </p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="flex items-center gap-2 px-4 py-2 bg-accent hover:bg-accent-bright text-bg-base font-medium text-sm rounded-lg transition-colors"
        >
          <Plus className="w-4 h-4" />
          Create Organization
        </button>
      </div>

      {/* Error state */}
      {error && (
        <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
          {error}
        </div>
      )}

      {/* Loading state */}
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="w-6 h-6 text-accent animate-spin" />
        </div>
      ) : orgs.length === 0 ? (
        /* Empty state */
        <div className="flex flex-col items-center justify-center py-20 text-center">
          <div className="w-12 h-12 rounded-2xl bg-surface flex items-center justify-center mb-4">
            <Building2 className="w-6 h-6 text-fg-muted" />
          </div>
          <h3 className="text-fg font-medium mb-1">No organizations yet</h3>
          <p className="text-sm text-fg-muted mb-4">
            Create your first organization to get started
          </p>
          <button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2 px-4 py-2 bg-accent hover:bg-accent-bright text-bg-base font-medium text-sm rounded-lg transition-colors"
          >
            <Plus className="w-4 h-4" />
            Create Organization
          </button>
        </div>
      ) : (
        /* Orgs table */
        <div className="bg-bg-elevated border border-white/[0.06] rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/[0.06]">
                  <th className="text-left text-xs font-mono tracking-widest text-fg-muted uppercase px-4 py-3">
                    Name
                  </th>
                  <th className="text-left text-xs font-mono tracking-widest text-fg-muted uppercase px-4 py-3 hidden sm:table-cell">
                    Phone
                  </th>
                  <th className="text-left text-xs font-mono tracking-widest text-fg-muted uppercase px-4 py-3 hidden md:table-cell">
                    Address
                  </th>
                  <th className="text-left text-xs font-mono tracking-widest text-fg-muted uppercase px-4 py-3 hidden lg:table-cell">
                    Created
                  </th>
                  <th className="text-left text-xs font-mono tracking-widest text-fg-muted uppercase px-4 py-3">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody>
                {orgs.map((org) => (
                  <tr
                    key={org.id}
                    className="border-b border-white/[0.04] last:border-0 hover:bg-surface/50 transition-colors"
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                          <Building2 className="w-4 h-4 text-accent" />
                        </div>
                        <div className="min-w-0">
                          <p className="text-sm font-medium text-fg truncate">
                            {org.name}
                          </p>
                          <p className="text-xs text-fg-muted truncate">
                            {org.slug}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3 hidden sm:table-cell">
                      <span className="text-sm text-fg-muted">
                        {org.phone || "—"}
                      </span>
                    </td>
                    <td className="px-4 py-3 hidden md:table-cell">
                      <span className="text-sm text-fg-muted truncate block max-w-[200px]">
                        {org.address || "—"}
                      </span>
                    </td>
                    <td className="px-4 py-3 hidden lg:table-cell">
                      <span className="text-sm text-fg-muted">
                        {new Date(org.created_at).toLocaleDateString()}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                          org.is_active
                            ? "bg-green-500/10 text-green-400"
                            : "bg-red-500/10 text-red-400"
                        }`}
                      >
                        {org.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Create Organization Modal */}
      {showCreate && (
        <CreateOrgModal
          onClose={() => setShowCreate(false)}
          onCreated={() => {
            setShowCreate(false);
            fetchOrgs();
          }}
        />
      )}
    </div>
  );
}

function CreateOrgModal({
  onClose,
  onCreated,
}: {
  onClose: () => void;
  onCreated: () => void;
}) {
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError("Name is required");
      return;
    }

    setError("");
    setSubmitting(true);
    try {
      await api.createOrg({
        name: name.trim(),
        address: address.trim() || undefined,
        phone: phone.trim() || undefined,
        email: email.trim() || undefined,
      });
      onCreated();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : "Failed to create organization"
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div
        className="bg-bg-elevated border border-white/[0.06] rounded-2xl w-full max-w-md shadow-card animate-fade-in"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
          <h2 className="text-lg font-semibold text-fg">
            Create Organization
          </h2>
          <button
            onClick={onClose}
            className="p-1 rounded-lg hover:bg-surface transition-colors"
          >
            <X className="w-4 h-4 text-fg-muted" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-5 space-y-4">
          <div>
            <label className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2">
              Name *
            </label>
            <div className="relative">
              <Building2 className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Gym name"
                className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                autoFocus
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2">
              Address
            </label>
            <div className="relative">
              <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="Gym address"
                className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2">
                Phone
              </label>
              <div className="relative">
                <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                <input
                  type="tel"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="98XXXXXXXX"
                  className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2">
                Email
              </label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="gym@email.com"
                  className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                />
              </div>
            </div>
          </div>

          {error && (
            <p className="text-red-400 text-sm">{error}</p>
          )}

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm font-medium rounded-lg hover:bg-surface-hover transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-accent hover:bg-accent-bright text-bg-base font-medium text-sm rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {submitting ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <>
                  <Plus className="w-4 h-4" />
                  Create
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
