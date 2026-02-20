"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type { Locale, MemberProfile, MemberStreak, MemberPackage } from "@/types";
import { Loader2, Package, Flame, Trophy, User } from "lucide-react";

export default function MemberDashboardPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [profile, setProfile] = useState<MemberProfile | null>(null);
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [profileRes, streaksRes, packagesRes] = await Promise.all([
          api.getMyProfile(),
          api.getMyStreaks(),
          api.getMyPackages(),
        ]);
        setProfile(profileRes);
        setStreaks(streaksRes.data ?? []);
        setPackages(packagesRes.data ?? []);
      } catch (err) {
        console.error("Failed to load member dashboard:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const orgName =
    profile?.organizations?.find((o) => o.org_id === user?.org_id)?.org_name ??
    "";

  const activePackage = packages.find((p) => p.status === "active");
  const daysRemaining = activePackage
    ? Math.max(
        0,
        Math.ceil(
          (new Date(activePackage.end_date).getTime() - Date.now()) /
            (1000 * 60 * 60 * 24)
        )
      )
    : 0;

  const streak = streaks[0];
  const currentStreak = streak?.current_streak ?? 0;
  const longestStreak = streak?.longest_streak ?? 0;

  return (
    <div className="space-y-6">
      {/* Welcome header */}
      <div>
        <h1 className="text-2xl font-semibold text-fg">
          {t.memberPages.welcomeBack}, {profile?.name ?? ""}
        </h1>
        {orgName && (
          <p className="text-sm text-fg-muted mt-1">{orgName}</p>
        )}
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Current Package */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
              <Package className="w-5 h-5 text-accent" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberPages.currentPackage}
            </span>
          </div>
          {activePackage ? (
            <>
              <p className="text-lg font-semibold text-fg">
                {activePackage.package_name}
              </p>
              <p className="text-xs text-fg-muted mt-1">
                {daysRemaining} {t.memberPages.daysRemaining}
              </p>
            </>
          ) : (
            <p className="text-sm text-fg-muted">
              {t.memberPages.noActivePackage}
            </p>
          )}
        </div>

        {/* Current Streak */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-orange-500/10 flex items-center justify-center">
              <Flame className="w-5 h-5 text-orange-400" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberPages.currentStreak}
            </span>
          </div>
          <p className="text-lg font-semibold text-fg">
            {currentStreak} {t.memberPages.days}
          </p>
        </div>

        {/* Longest Streak */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-yellow-500/10 flex items-center justify-center">
              <Trophy className="w-5 h-5 text-yellow-400" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberPages.longestStreak}
            </span>
          </div>
          <p className="text-lg font-semibold text-fg">
            {longestStreak} {t.memberPages.days}
          </p>
        </div>

        {/* Profile info */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center">
              <User className="w-5 h-5 text-blue-400" />
            </div>
            <span className="text-sm text-fg-muted">
              {t.memberNav.profile}
            </span>
          </div>
          <p className="text-lg font-semibold text-fg">{profile?.phone}</p>
          {profile?.created_at && (
            <p className="text-xs text-fg-muted mt-1">
              {t.memberPages.memberSince}{" "}
              {new Date(profile.created_at).toLocaleDateString()}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
