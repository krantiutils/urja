"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth, isOperator } from "@/lib/auth";
import { Loader2 } from "lucide-react";

interface AuthGuardProps {
  children: React.ReactNode;
  locale: string;
  /**
   * Require that the caller runs a gym rather than trains at one.
   *
   * The API already refuses org-scoped requests from a plain member, so this is
   * not what keeps them out — it is what stops them being shown a gym
   * management console in which every panel fails to load. Wrong-place is a
   * better message than broken.
   */
  requireOperator?: boolean;
}

export function AuthGuard({ children, locale, requireOperator }: AuthGuardProps) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const router = useRouter();

  // Only redirect once memberships are known. If the profile request failed we
  // cannot tell an admin from a member, and throwing a real admin out of their
  // own dashboard over a network blip is worse than showing it briefly — the
  // API refuses their requests either way.
  const wrongPlace =
    requireOperator && isAuthenticated && user?.profile_loaded === true && !isOperator(user);

  useEffect(() => {
    if (isLoading) return;
    if (!isAuthenticated) {
      router.replace(`/${locale}/login`);
      return;
    }
    if (wrongPlace) {
      router.replace(`/${locale}/member`);
    }
  }, [isAuthenticated, isLoading, wrongPlace, router, locale]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-bg-base">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  if (!isAuthenticated || wrongPlace) {
    return null;
  }

  return <>{children}</>;
}
