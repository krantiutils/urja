"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type {
  Locale,
  StaffMember,
  StaffRole,
  CreateStaffRequest,
  UpdateStaffRequest,
} from "@/types";
import {
  UserCog,
  Plus,
  Search,
  X,
  Pencil,
  Trash2,
  Loader2,
} from "lucide-react";

const STAFF_ROLES: StaffRole[] = [
  "owner",
  "manager",
  "trainer",
  "receptionist",
];

function StatusBadge({
  status,
  label,
}: {
  status: string;
  label: string;
}) {
  const colors: Record<string, string> = {
    active: "bg-green-100 text-green-700",
    suspended: "bg-amber-100 text-amber-700",
    left: "bg-red-100 text-red-700",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase ${
        colors[status] ?? colors.active
      }`}
    >
      {label}
    </span>
  );
}

function StaffRoleBadge({
  role,
  label,
}: {
  role: string;
  label: string;
}) {
  const colors: Record<string, string> = {
    owner: "bg-amber-100 text-amber-700",
    manager: "bg-purple-100 text-purple-700",
    trainer: "bg-blue-100 text-blue-700",
    receptionist: "bg-cyan-100 text-cyan-700",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase ${
        colors[role] ?? colors.trainer
      }`}
    >
      {label}
    </span>
  );
}

export default function StaffPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [staffList, setStaffList] = useState<StaffMember[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  // Modal state
  const [showModal, setShowModal] = useState(false);
  const [editingStaff, setEditingStaff] = useState<StaffMember | null>(
    null
  );
  const [formName, setFormName] = useState("");
  const [formPhone, setFormPhone] = useState("");
  const [formEmail, setFormEmail] = useState("");
  const [formStaffRole, setFormStaffRole] =
    useState<StaffRole>("trainer");
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchStaff = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listStaff(orgId, {
        search: search || undefined,
        limit: 100,
      });
      setStaffList(res.data ?? []);
      setTotal(res.total ?? 0);
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setLoading(false);
    }
  }, [orgId, search, t.common.error]);

  useEffect(() => {
    fetchStaff();
  }, [fetchStaff]);

  const openCreateModal = () => {
    setEditingStaff(null);
    setFormName("");
    setFormPhone("");
    setFormEmail("");
    setFormStaffRole("trainer");
    setFormError(null);
    setShowModal(true);
  };

  const openEditModal = (s: StaffMember) => {
    setEditingStaff(s);
    setFormName(s.name);
    setFormPhone(s.phone);
    setFormEmail(s.email ?? "");
    setFormStaffRole(s.staff_role);
    setFormError(null);
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setFormLoading(true);
    setFormError(null);

    try {
      if (editingStaff) {
        const data: UpdateStaffRequest = {
          name: formName,
          email: formEmail || undefined,
          staff_role: formStaffRole,
        };
        await api.updateStaff(orgId, editingStaff.id, data);
      } else {
        const data: CreateStaffRequest = {
          phone: formPhone,
          name: formName,
          email: formEmail || undefined,
          staff_role: formStaffRole,
        };
        await api.createStaff(orgId, data);
      }
      setShowModal(false);
      fetchStaff();
    } catch (err) {
      setFormError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (staffId: string) => {
    if (!orgId) return;
    if (!confirm(t.staff.confirmDelete)) return;
    try {
      await api.deleteStaff(orgId, staffId);
      fetchStaff();
    } catch (err) {
      setError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const statusLabel = (s: string) => {
    const map: Record<string, string> = {
      active: t.staff.active,
      suspended: t.staff.suspended,
      left: t.staff.left,
    };
    return map[s] ?? s;
  };

  const staffRoleLabel = (r: string) => {
    const map: Record<string, string> = {
      owner: t.staff.owner,
      manager: t.staff.manager,
      trainer: t.staff.trainer,
      receptionist: t.staff.receptionist,
    };
    return map[r] ?? r;
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-md-on-surface flex items-center gap-2">
          <UserCog className="w-5 h-5 text-md-primary" />
          {t.staff.title}
        </h1>
        <button
          onClick={openCreateModal}
          className="flex items-center gap-2 px-4 py-2 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors shadow-md-1"
        >
          <Plus className="w-4 h-4" />
          {t.staff.addStaff}
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-md-on-surface-variant" />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder={t.staff.searchPlaceholder}
          className="w-full pl-10 pr-4 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
        />
      </div>

      {/* Error */}
      {error && (
        <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
          {error}
        </div>
      )}

      {/* Table */}
      <div className="bg-md-surface-container border border-md-outline-variant rounded-3xl shadow-md-1 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 text-md-primary animate-spin" />
          </div>
        ) : staffList.length === 0 ? (
          <div className="text-center py-16 text-md-on-surface-variant text-sm">
            {t.staff.noStaff}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-md-outline-variant">
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.name}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.phone}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.staffRole}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.status}
                  </th>
                  <th className="text-left px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.joined}
                  </th>
                  <th className="text-right px-5 py-3 text-xs font-medium tracking-wider text-md-on-surface-variant uppercase">
                    {t.staff.actions}
                  </th>
                </tr>
              </thead>
              <tbody>
                {staffList.map((s) => (
                  <tr
                    key={s.id}
                    className="border-b border-md-outline-variant/50 last:border-0 hover:bg-md-surface-container-high transition-colors"
                  >
                    <td className="px-5 py-3">
                      <div>
                        <p className="text-md-on-surface font-medium">{s.name}</p>
                        {s.email && (
                          <p className="text-xs text-md-on-surface-variant">
                            {s.email}
                          </p>
                        )}
                      </div>
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {s.phone}
                    </td>
                    <td className="px-5 py-3">
                      <StaffRoleBadge
                        role={s.staff_role}
                        label={staffRoleLabel(s.staff_role)}
                      />
                    </td>
                    <td className="px-5 py-3">
                      <StatusBadge
                        status={s.status}
                        label={statusLabel(s.status)}
                      />
                    </td>
                    <td className="px-5 py-3 text-md-on-surface-variant font-mono text-xs">
                      {new Date(s.joined_at).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center justify-end gap-1">
                        <button
                          onClick={() => openEditModal(s)}
                          title={t.staff.editStaff}
                          className="p-1.5 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant hover:text-md-on-surface transition-colors"
                        >
                          <Pencil className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => handleDelete(s.id)}
                          title={t.common.delete}
                          className="p-1.5 rounded-lg hover:bg-red-100 text-md-on-surface-variant hover:text-red-600 transition-colors"
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
          <div className="px-5 py-3 border-t border-md-outline-variant">
            <p className="text-xs text-md-on-surface-variant">
              {total} {t.staff.title.toLowerCase()}
            </p>
          </div>
        )}
      </div>

      {/* Add/Edit Staff Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-md-surface-container border border-md-outline-variant rounded-3xl shadow-md-1 w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-md-outline-variant">
              <h2 className="text-base font-semibold text-md-on-surface">
                {editingStaff
                  ? t.staff.editStaff
                  : t.staff.addStaffTitle}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 rounded-lg hover:bg-md-surface-container-high text-md-on-surface-variant transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-5 space-y-4">
              {formError && (
                <div className="p-3 bg-md-error-container border border-md-error/20 rounded-xl text-sm text-md-error">
                  {formError}
                </div>
              )}
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.staff.name} *
                </label>
                <input
                  type="text"
                  required
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  placeholder={t.staff.namePlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.staff.phone} {editingStaff ? "" : "*"}
                </label>
                <input
                  type="tel"
                  required={!editingStaff}
                  disabled={!!editingStaff}
                  value={formPhone}
                  onChange={(e) => setFormPhone(e.target.value)}
                  placeholder={t.staff.phonePlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors disabled:opacity-50"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.staff.email}
                </label>
                <input
                  type="email"
                  value={formEmail}
                  onChange={(e) => setFormEmail(e.target.value)}
                  placeholder={t.staff.emailPlaceholder}
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface placeholder:text-md-on-surface-variant focus:outline-none focus:border-md-primary transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs text-md-on-surface-variant mb-1.5">
                  {t.staff.staffRole} *
                </label>
                <select
                  value={formStaffRole}
                  onChange={(e) =>
                    setFormStaffRole(e.target.value as StaffRole)
                  }
                  className="w-full px-3 py-2.5 bg-md-surface-container-lowest border border-md-outline-variant rounded-xl text-sm text-md-on-surface focus:outline-none focus:border-md-primary transition-colors"
                >
                  {STAFF_ROLES.map((r) => (
                    <option key={r} value={r}>
                      {staffRoleLabel(r)}
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2.5 bg-md-surface-container-high border border-md-outline-variant text-md-on-surface text-sm rounded-full hover:bg-md-surface-container-highest transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={formLoading}
                  className="flex-1 px-4 py-2.5 bg-md-primary text-md-on-primary font-medium text-sm rounded-full hover:bg-md-primary/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
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
