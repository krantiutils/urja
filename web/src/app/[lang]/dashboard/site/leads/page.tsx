"use client";

import { useCallback, useEffect, useState } from "react";
import { Loader2, Phone } from "lucide-react";
import { api, ApiRequestError } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import type { LeadStatus, SiteLead } from "@/types/site";

/**
 * Enquiries captured by lead_form sections on the gym's public site.
 *
 * A gym answers these by phone, so the number is a tel: link and the status is
 * a single select rather than a workflow — the point is to know who has been
 * called back, not to model a CRM.
 */

const STATUSES: LeadStatus[] = ["new", "contacted", "trial_booked", "joined", "lost"];

const STATUS_CLASSES: Record<LeadStatus, string> = {
  new: "bg-accent/10 text-accent border-accent/20",
  contacted: "bg-sky-500/10 text-sky-400 border-sky-500/20",
  trial_booked: "bg-amber-500/10 text-amber-400 border-amber-500/20",
  joined: "bg-emerald-500/10 text-emerald-400 border-emerald-500/20",
  lost: "bg-gray-500/10 text-gray-400 border-gray-500/20",
};

export default function SiteLeadsPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();
  const orgId = user?.org_id;

  const [leads, setLeads] = useState<SiteLead[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<LeadStatus | "all">("all");

  const statusLabel: Record<LeadStatus, string> = {
    new: t.site.statusNew,
    contacted: t.site.statusContacted,
    trial_booked: t.site.statusTrialBooked,
    joined: t.site.statusJoined,
    lost: t.site.statusLost,
  };

  const fetchLeads = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listSiteLeads(orgId, { limit: 200 });
      setLeads(res.data ?? []);
    } catch (err) {
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    } finally {
      setLoading(false);
    }
  }, [orgId, t.common.error]);

  useEffect(() => {
    fetchLeads();
  }, [fetchLeads]);

  async function setStatus(lead: SiteLead, status: LeadStatus) {
    if (!orgId) return;
    // Optimistic: the select should not freeze while the request is in flight.
    const previous = lead.status;
    setLeads((prev) => prev.map((l) => (l.id === lead.id ? { ...l, status } : l)));
    try {
      await api.updateSiteLead(orgId, lead.id, status);
    } catch (err) {
      setLeads((prev) =>
        prev.map((l) => (l.id === lead.id ? { ...l, status: previous } : l))
      );
      setError(err instanceof ApiRequestError ? err.message : t.common.error);
    }
  }

  const visible = filter === "all" ? leads : leads.filter((l) => l.status === filter);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24">
        <Loader2 className="w-6 h-6 animate-spin text-fg-muted" />
      </div>
    );
  }

  return (
    <div className="max-w-5xl">
      <h1 className="text-2xl font-semibold text-fg">{t.site.enquiries}</h1>
      <p className="text-sm text-fg-muted mt-1 mb-6">{t.site.enquiriesSubtitle}</p>

      {error && (
        <div className="mb-5 px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      <div className="flex flex-wrap gap-2 mb-5">
        {(["all", ...STATUSES] as const).map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setFilter(s)}
            className={`px-3.5 py-1.5 rounded-full text-xs transition-colors ${
              filter === s
                ? "bg-accent text-accent-fg"
                : "border border-white/[0.06] text-fg-muted hover:text-fg"
            }`}
          >
            {s === "all" ? t.site.all : statusLabel[s]}
            {s !== "all" && (
              <span className="ml-1.5 opacity-60">
                {leads.filter((l) => l.status === s).length}
              </span>
            )}
          </button>
        ))}
      </div>

      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5 shadow-card">
        {visible.length === 0 ? (
          <p className="py-10 text-center text-sm text-fg-muted">{t.site.noEnquiries}</p>
        ) : (
          <div className="divide-y divide-white/[0.06]">
            {visible.map((lead) => (
              <div key={lead.id} className="py-4 flex flex-wrap items-start gap-4">
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2.5">
                    <span className="text-sm font-medium text-fg">{lead.name}</span>
                    <a
                      href={`tel:${lead.phone}`}
                      className="inline-flex items-center gap-1.5 text-xs font-mono text-accent hover:underline"
                    >
                      <Phone className="w-3 h-3" />
                      {lead.phone}
                    </a>
                    <span
                      className={`px-2 py-0.5 rounded-full text-[10px] font-mono uppercase border ${STATUS_CLASSES[lead.status]}`}
                    >
                      {statusLabel[lead.status]}
                    </span>
                  </div>

                  {lead.message && (
                    <p className="mt-1.5 text-sm text-fg-muted leading-relaxed">
                      {lead.message}
                    </p>
                  )}

                  <p className="mt-1.5 text-[11px] text-fg-muted">
                    {[
                      lead.interest && `${t.site.leadInterest}: ${lead.interest}`,
                      lead.source_page && `${t.site.leadPage}: /${lead.source_page}`,
                      new Date(lead.created_at).toLocaleDateString(
                        locale === "ne" ? "ne-NP" : "en-GB",
                        { day: "numeric", month: "short", year: "numeric" }
                      ),
                    ]
                      .filter(Boolean)
                      .join(" · ")}
                  </p>
                </div>

                <select
                  value={lead.status}
                  onChange={(e) => setStatus(lead, e.target.value as LeadStatus)}
                  aria-label={t.site.leadStatus}
                  className="shrink-0 px-3 py-2 bg-input-bg border border-white/[0.06] rounded-xl text-xs text-fg focus:outline-none focus:border-accent/50"
                >
                  {STATUSES.map((s) => (
                    <option key={s} value={s}>
                      {statusLabel[s]}
                    </option>
                  ))}
                </select>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
