"use client";

import { useCallback, useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api, ApiRequestError } from "@/lib/api";
import type { Locale, NfcCard, NfcDevice } from "@/types";
import {
  CreditCard,
  Plus,
  X,
  Loader2,
  Wifi,
  WifiOff,
  Lock,
  Unlock,
  UserPlus,
  UserMinus,
} from "lucide-react";

type Tab = "cards" | "devices";

function StatusBadge({
  status,
  label,
}: {
  status: string;
  label: string;
}) {
  const colors: Record<string, string> = {
    available: "bg-accent/10 text-accent border-accent/20",
    assigned: "bg-blue-500/10 text-blue-400 border-blue-500/20",
    disabled: "bg-red-500/10 text-red-400 border-red-500/20",
    online: "bg-accent/10 text-accent border-accent/20",
    offline: "bg-red-500/10 text-red-400 border-red-500/20",
    locked: "bg-amber-500/10 text-amber-400 border-amber-500/20",
    unlocked: "bg-accent/10 text-accent border-accent/20",
  };
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${
        colors[status] ?? "bg-accent/10 text-accent border-accent/20"
      }`}
    >
      {label}
    </span>
  );
}

function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
}

export default function NfcPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [activeTab, setActiveTab] = useState<Tab>("cards");

  // Cards state
  const [cards, setCards] = useState<NfcCard[]>([]);
  const [cardsLoading, setCardsLoading] = useState(true);
  const [cardsError, setCardsError] = useState<string | null>(null);

  // Devices state
  const [devices, setDevices] = useState<NfcDevice[]>([]);
  const [devicesLoading, setDevicesLoading] = useState(false);
  const [devicesError, setDevicesError] = useState<string | null>(null);

  // Register card modal
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const [registerCardNumber, setRegisterCardNumber] = useState("");
  const [registerError, setRegisterError] = useState<string | null>(null);
  const [registerLoading, setRegisterLoading] = useState(false);

  // Assign modal
  const [showAssignModal, setShowAssignModal] = useState(false);
  const [assignCardId, setAssignCardId] = useState<string | null>(null);
  const [assignMemberId, setAssignMemberId] = useState("");
  const [assignError, setAssignError] = useState<string | null>(null);
  const [assignLoading, setAssignLoading] = useState(false);

  const orgId = user?.org_id;

  const fetchCards = useCallback(async () => {
    if (!orgId) return;
    setCardsLoading(true);
    setCardsError(null);
    try {
      const res = await api.listNfcCards(orgId, { limit: 100 });
      setCards(res.data ?? []);
    } catch (err) {
      setCardsError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setCardsLoading(false);
    }
  }, [orgId, t.common.error]);

  const fetchDevices = useCallback(async () => {
    if (!orgId) return;
    setDevicesLoading(true);
    setDevicesError(null);
    try {
      const res = await api.listNfcDevices(orgId);
      setDevices(res.data ?? []);
    } catch (err) {
      setDevicesError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setDevicesLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    if (activeTab === "cards") {
      fetchCards();
    } else {
      fetchDevices();
    }
  }, [activeTab, fetchCards, fetchDevices]);

  const handleRegisterCard = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setRegisterLoading(true);
    setRegisterError(null);
    try {
      await api.registerNfcCard(orgId, { card_number: registerCardNumber });
      setShowRegisterModal(false);
      setRegisterCardNumber("");
      fetchCards();
    } catch (err) {
      setRegisterError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setRegisterLoading(false);
    }
  };

  const openAssignModal = (cardId: string) => {
    setAssignCardId(cardId);
    setAssignMemberId("");
    setAssignError(null);
    setShowAssignModal(true);
  };

  const handleAssign = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId || !assignCardId) return;
    setAssignLoading(true);
    setAssignError(null);
    try {
      await api.assignNfcCard(orgId, assignCardId, assignMemberId);
      setShowAssignModal(false);
      setAssignCardId(null);
      setAssignMemberId("");
      fetchCards();
    } catch (err) {
      setAssignError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    } finally {
      setAssignLoading(false);
    }
  };

  const handleUnassign = async (cardId: string) => {
    if (!orgId) return;
    if (!confirm(t.nfc.confirmUnassign)) return;
    try {
      await api.unassignNfcCard(orgId, cardId);
      fetchCards();
    } catch (err) {
      setCardsError(
        err instanceof ApiRequestError ? err.message : t.common.error
      );
    }
  };

  const statusLabel = (s: string) => {
    const map: Record<string, string> = {
      available: t.nfc.available,
      assigned: t.nfc.assigned,
      disabled: t.nfc.disabled,
    };
    return map[s] ?? s;
  };

  const deviceStatusLabel = (s: string) => {
    const map: Record<string, string> = {
      online: t.nfc.online,
      offline: t.nfc.offline,
    };
    return map[s] ?? s;
  };

  const doorStateLabel = (s: string) => {
    const map: Record<string, string> = {
      locked: t.nfc.locked,
      unlocked: t.nfc.unlocked,
    };
    return map[s] ?? s;
  };

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-fg flex items-center gap-2">
          <CreditCard className="w-5 h-5 text-accent" />
          {t.nfc.title}
        </h1>
        <button
          onClick={() => {
            setRegisterCardNumber("");
            setRegisterError(null);
            setShowRegisterModal(true);
          }}
          className="flex items-center gap-2 px-4 py-2 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors shadow-accent-glow"
        >
          <Plus className="w-4 h-4" />
          {t.nfc.registerCard}
        </button>
      </div>

      {/* Tabs */}
      <div className="flex gap-4 border-b border-white/[0.06] pb-0">
        <button
          onClick={() => setActiveTab("cards")}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "cards"
              ? "bg-accent/10 text-accent border-b-2 border-accent"
              : "text-fg-muted hover:text-fg"
          }`}
        >
          {t.nfc.cards}
        </button>
        <button
          onClick={() => setActiveTab("devices")}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "devices"
              ? "bg-accent/10 text-accent border-b-2 border-accent"
              : "text-fg-muted hover:text-fg"
          }`}
        >
          {t.nfc.devices}
        </button>
      </div>

      {/* Cards Tab */}
      {activeTab === "cards" && (
        <>
          {/* Error */}
          {cardsError && (
            <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
              {cardsError}
            </div>
          )}

          {/* Cards Table */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            {cardsLoading ? (
              <div className="flex items-center justify-center py-16">
                <Loader2 className="w-6 h-6 text-accent animate-spin" />
              </div>
            ) : cards.length === 0 ? (
              <div className="text-center py-16 text-fg-muted text-sm">
                {t.nfc.noCards}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-white/[0.06]">
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.cardNumber}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.assignedTo}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.status}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.assignedAt}
                      </th>
                      <th className="text-right px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.actions}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {cards.map((card) => (
                      <tr
                        key={card.id}
                        className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                      >
                        <td className="px-5 py-3 text-fg font-mono text-xs">
                          {card.card_number}
                        </td>
                        <td className="px-5 py-3 text-fg-muted text-sm">
                          {card.full_name ?? "\u2014"}
                        </td>
                        <td className="px-5 py-3">
                          <StatusBadge
                            status={card.status}
                            label={statusLabel(card.status)}
                          />
                        </td>
                        <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                          {card.assigned_at
                            ? new Date(card.assigned_at).toLocaleDateString()
                            : "\u2014"}
                        </td>
                        <td className="px-5 py-3">
                          <div className="flex items-center justify-end gap-1">
                            {card.status === "available" && (
                              <button
                                onClick={() => openAssignModal(card.id)}
                                title={t.nfc.assign}
                                className="p-1.5 rounded-lg hover:bg-surface text-fg-muted hover:text-fg transition-colors"
                              >
                                <UserPlus className="w-4 h-4" />
                              </button>
                            )}
                            {card.status === "assigned" && (
                              <button
                                onClick={() => handleUnassign(card.id)}
                                title={t.nfc.unassign}
                                className="p-1.5 rounded-lg hover:bg-red-500/10 text-fg-muted hover:text-red-400 transition-colors"
                              >
                                <UserMinus className="w-4 h-4" />
                              </button>
                            )}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </>
      )}

      {/* Devices Tab */}
      {activeTab === "devices" && (
        <>
          {/* Error */}
          {devicesError && (
            <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
              {devicesError}
            </div>
          )}

          {/* Devices Table */}
          <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl shadow-card overflow-hidden">
            {devicesLoading ? (
              <div className="flex items-center justify-center py-16">
                <Loader2 className="w-6 h-6 text-accent animate-spin" />
              </div>
            ) : devices.length === 0 ? (
              <div className="text-center py-16 text-fg-muted text-sm">
                {t.nfc.noDevices}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-white/[0.06]">
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.deviceName}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.identifier}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.status}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.doorState}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.lastSync}
                      </th>
                      <th className="text-left px-5 py-3 text-xs font-mono tracking-widest text-fg-muted uppercase">
                        {t.nfc.uptime}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {devices.map((device) => (
                      <tr
                        key={device.id}
                        className="border-b border-white/[0.03] last:border-0 hover:bg-surface transition-colors"
                      >
                        <td className="px-5 py-3">
                          <div className="flex items-center gap-2">
                            {device.status === "online" ? (
                              <Wifi className="w-4 h-4 text-accent" />
                            ) : (
                              <WifiOff className="w-4 h-4 text-red-400" />
                            )}
                            <span className="text-fg font-medium">
                              {device.name}
                            </span>
                          </div>
                        </td>
                        <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                          {device.device_identifier}
                        </td>
                        <td className="px-5 py-3">
                          <StatusBadge
                            status={device.status}
                            label={deviceStatusLabel(device.status)}
                          />
                        </td>
                        <td className="px-5 py-3">
                          <div className="flex items-center gap-1.5">
                            {device.door_state === "locked" ? (
                              <Lock className="w-3.5 h-3.5 text-amber-400" />
                            ) : (
                              <Unlock className="w-3.5 h-3.5 text-accent" />
                            )}
                            <StatusBadge
                              status={device.door_state}
                              label={doorStateLabel(device.door_state)}
                            />
                          </div>
                        </td>
                        <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                          {device.last_sync_at
                            ? new Date(device.last_sync_at).toLocaleString()
                            : "\u2014"}
                        </td>
                        <td className="px-5 py-3 text-fg-muted font-mono text-xs">
                          {formatUptime(device.uptime_seconds)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </>
      )}

      {/* Register Card Modal */}
      {showRegisterModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {t.nfc.registerCard}
              </h2>
              <button
                onClick={() => setShowRegisterModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleRegisterCard} className="p-5 space-y-4">
              {registerError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {registerError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.nfc.cardNumber} *
                </label>
                <input
                  type="text"
                  required
                  value={registerCardNumber}
                  onChange={(e) => setRegisterCardNumber(e.target.value)}
                  placeholder={t.nfc.cardNumberPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowRegisterModal(false)}
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={registerLoading}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {registerLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.common.save}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Assign Card Modal */}
      {showAssignModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-bg-elevated border border-white/[0.06] rounded-2xl shadow-card w-full max-w-md">
            <div className="flex items-center justify-between p-5 border-b border-white/[0.06]">
              <h2 className="text-base font-semibold text-fg">
                {t.nfc.assign}
              </h2>
              <button
                onClick={() => setShowAssignModal(false)}
                className="p-1 rounded-lg hover:bg-surface text-fg-muted transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <form onSubmit={handleAssign} className="p-5 space-y-4">
              {assignError && (
                <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm text-red-400">
                  {assignError}
                </div>
              )}
              <div>
                <label className="block text-xs text-fg-muted mb-1.5">
                  {t.nfc.assignedTo} *
                </label>
                <input
                  type="text"
                  required
                  value={assignMemberId}
                  onChange={(e) => setAssignMemberId(e.target.value)}
                  placeholder={t.nfc.memberIdPlaceholder}
                  className="w-full px-3 py-2.5 bg-input-bg border border-white/[0.06] rounded-xl text-sm text-fg placeholder:text-fg-muted focus:outline-none focus:border-accent/50 transition-colors"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowAssignModal(false)}
                  className="flex-1 px-4 py-2.5 bg-surface border border-white/[0.06] text-fg text-sm rounded-xl hover:bg-surface-hover transition-colors"
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  disabled={assignLoading}
                  className="flex-1 px-4 py-2.5 bg-accent text-bg-deep font-medium text-sm rounded-xl hover:bg-accent-bright transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {assignLoading && (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  )}
                  {t.nfc.assign}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
