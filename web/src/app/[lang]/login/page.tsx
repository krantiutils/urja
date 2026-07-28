"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import { ApiRequestError } from "@/lib/api";
import type { Locale } from "@/types";
import { Lock, Phone, ShieldCheck, ArrowRight, Loader2 } from "lucide-react";

const OTP_RESEND_SECONDS = 60;

export default function LoginPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const router = useRouter();
  const { login, verifyOtp, passwordLogin, isAuthenticated, isLoading: authLoading, user } = useAuth();

  const [step, setStep] = useState<"phone" | "otp" | "password">("phone");
  const [password, setPassword] = useState("");
  const [phone, setPhone] = useState("");
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [resendTimer, setResendTimer] = useState(0);

  // Redirect if already authenticated
  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      if (user?.is_super_admin) {
        router.replace(`/${locale}/super-admin`);
        return;
      }
      if (user && !user.onboarding_completed) {
        router.replace(`/${locale}/onboarding`);
        return;
      }
      const token = localStorage.getItem("access_token");
      if (token) {
        try {
          const payload = JSON.parse(atob(token.split(".")[1]));
          if (payload.role === "member") {
            router.replace(`/${locale}/member`);
            return;
          }
        } catch {}
      }
      router.replace(`/${locale}/dashboard`);
    }
  }, [isAuthenticated, authLoading, router, locale, user]);

  // Resend timer countdown
  useEffect(() => {
    if (resendTimer <= 0) return;
    const interval = setInterval(() => {
      setResendTimer((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(interval);
  }, [resendTimer]);

  const handleRequestOtp = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      setError("");

      // Basic client-side validation
      const cleaned = phone.replace(/\s+/g, "");
      if (!/^((\+?977)?9[6-9]\d{8})$/.test(cleaned)) {
        setError(t.auth.invalidPhone);
        return;
      }

      setIsSubmitting(true);
      try {
        await login(cleaned);
        setStep("otp");
        setResendTimer(OTP_RESEND_SECONDS);
      } catch (err) {
        if (err instanceof ApiRequestError) {
          setError(err.message);
        } else {
          setError(t.common.error);
        }
      } finally {
        setIsSubmitting(false);
      }
    },
    [phone, login, t]
  );

  /** Where a session goes once established, regardless of how it was created. */
  const routeAfterLogin = useCallback(
    (result: { onboarding_completed: boolean }) => {
      const token = localStorage.getItem("access_token");
      const claims = (() => {
        try {
          return token ? JSON.parse(atob(token.split(".")[1])) : null;
        } catch {
          return null;
        }
      })();

      if (claims?.is_super_admin) {
        router.replace(`/${locale}/super-admin`);
        return;
      }
      if (!result.onboarding_completed) {
        router.replace(`/${locale}/onboarding`);
        return;
      }
      router.replace(`/${locale}/${claims?.role === "member" ? "member" : "dashboard"}`);
    },
    [router, locale]
  );

  const handlePasswordLogin = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      setError("");

      const cleaned = phone.replace(/\s+/g, "");
      if (!/^((\+?977)?9[6-9]\d{8})$/.test(cleaned)) {
        setError(t.auth.invalidPhone);
        return;
      }

      setIsSubmitting(true);
      try {
        routeAfterLogin(await passwordLogin(cleaned, password));
      } catch (err) {
        // The API returns one message for every failure so it cannot be used
        // to discover which numbers are registered; pass it straight through.
        setError(
          err instanceof ApiRequestError ? err.message : t.auth.invalidCredentials
        );
      } finally {
        setIsSubmitting(false);
      }
    },
    [phone, password, passwordLogin, routeAfterLogin, t]
  );

  const handleVerifyOtp = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      setError("");

      if (otp.length !== 6) {
        setError(t.auth.invalidOtp);
        return;
      }

      setIsSubmitting(true);
      try {
        const cleaned = phone.replace(/\s+/g, "");
        const result = await verifyOtp(cleaned, otp);
        // Check super admin first
        const token = localStorage.getItem("access_token");
        if (token) {
          try {
            const payload = JSON.parse(atob(token.split(".")[1]));
            if (payload.is_super_admin) {
              router.replace(`/${locale}/super-admin`);
              return;
            }
          } catch {}
        }
        // New users or users who haven't completed onboarding → onboarding page
        if (!result.onboarding_completed) {
          router.replace(`/${locale}/onboarding`);
          return;
        }
        if (token) {
          try {
            const payload = JSON.parse(atob(token.split(".")[1]));
            if (payload.role === "member") {
              router.replace(`/${locale}/member`);
              return;
            }
          } catch {}
        }
        router.replace(`/${locale}/dashboard`);
      } catch (err) {
        if (err instanceof ApiRequestError) {
          setError(err.message);
        } else {
          setError(t.common.error);
        }
      } finally {
        setIsSubmitting(false);
      }
    },
    [otp, phone, verifyOtp, router, locale, t]
  );

  const handleResendOtp = useCallback(async () => {
    if (resendTimer > 0) return;
    setError("");
    setIsSubmitting(true);
    try {
      const cleaned = phone.replace(/\s+/g, "");
      await login(cleaned);
      setResendTimer(OTP_RESEND_SECONDS);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError(t.common.error);
      }
    } finally {
      setIsSubmitting(false);
    }
  }, [resendTimer, phone, login, t]);

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-bg-base">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-bg-base px-4">
      {/* Background subtle gradient */}
      <div className="fixed inset-0 bg-[radial-gradient(ellipse_at_top,#0a0a0f_0%,#050506_50%,#020203_100%)]" />

      <div className="relative w-full max-w-sm animate-fade-in">
        {/* Logo / Brand */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-semibold tracking-tight bg-gradient-to-b from-white via-white/95 to-white/70 bg-clip-text text-transparent">
            {t.common.appName}
          </h1>
          <p className="mt-2 text-fg-muted text-sm">
            {t.auth.subtitle}
          </p>
        </div>

        {/* Login Card */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-6 shadow-card">
          <h2 className="text-lg font-semibold text-fg mb-6">
            {t.auth.signIn}
          </h2>

          {step === "phone" ? (
            <form onSubmit={handleRequestOtp} className="space-y-4">
              <div>
                <label
                  htmlFor="phone"
                  className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2"
                >
                  {t.auth.phoneLabel}
                </label>
                <div className="relative">
                  <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                  <input
                    id="phone"
                    type="tel"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t.auth.phonePlaceholder}
                    className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                    autoComplete="tel"
                    autoFocus
                  />
                </div>
              </div>

              {error && (
                <p className="text-red-400 text-sm">{error}</p>
              )}

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 bg-accent hover:bg-accent-bright text-bg-base font-semibold py-2.5 rounded-lg shadow-accent-glow hover:shadow-[0_0_0_1px_rgba(132,204,22,0.6),0_4px_16px_rgba(132,204,22,0.4)] active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    {t.auth.requestOtp}
                    <ArrowRight className="w-4 h-4" />
                  </>
                )}
              </button>

              {/* Members sign in by code; staff and admins do it several times
                  a day, so a password is offered as the alternative. */}
              <button
                type="button"
                onClick={() => {
                  setError("");
                  setStep("password");
                }}
                className="w-full text-center text-sm text-fg-muted hover:text-fg transition-colors pt-1"
              >
                {t.auth.usePassword}
              </button>
            </form>
          ) : step === "password" ? (
            <form onSubmit={handlePasswordLogin} className="space-y-4">
              <div>
                <label
                  htmlFor="phone-pw"
                  className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2"
                >
                  {t.auth.phoneLabel}
                </label>
                <div className="relative">
                  <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                  <input
                    id="phone-pw"
                    type="tel"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t.auth.phonePlaceholder}
                    className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                    autoComplete="username"
                    autoFocus
                  />
                </div>
              </div>

              <div>
                <label
                  htmlFor="password"
                  className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2"
                >
                  {t.auth.passwordLabel}
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                  <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t.auth.passwordPlaceholder}
                    className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors"
                    autoComplete="current-password"
                  />
                </div>
              </div>

              {error && <p className="text-red-400 text-sm">{error}</p>}

              <button
                type="submit"
                disabled={isSubmitting || !password}
                className="w-full flex items-center justify-center gap-2 bg-accent hover:bg-accent-bright text-bg-base font-semibold py-2.5 rounded-lg shadow-accent-glow active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    {t.auth.signInWithPassword}
                    <ArrowRight className="w-4 h-4" />
                  </>
                )}
              </button>

              <button
                type="button"
                onClick={() => {
                  setError("");
                  setPassword("");
                  setStep("phone");
                }}
                className="w-full text-center text-sm text-fg-muted hover:text-fg transition-colors pt-1"
              >
                {t.auth.useOtp}
              </button>
            </form>
          ) : (
            <form onSubmit={handleVerifyOtp} className="space-y-4">
              <p className="text-sm text-fg-muted mb-2">
                {t.auth.otpSent}
              </p>

              <div>
                <label
                  htmlFor="otp"
                  className="block text-xs font-mono tracking-widest text-fg-muted uppercase mb-2"
                >
                  {t.auth.otpLabel}
                </label>
                <div className="relative">
                  <ShieldCheck className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-muted" />
                  <input
                    id="otp"
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    maxLength={6}
                    value={otp}
                    onChange={(e) =>
                      setOtp(e.target.value.replace(/\D/g, ""))
                    }
                    placeholder={t.auth.otpPlaceholder}
                    className="w-full bg-input-bg border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-fg placeholder:text-gray-500 focus:outline-none focus:border-accent focus:ring-2 focus:ring-accent/20 transition-colors tracking-[0.3em] text-center font-mono"
                    autoComplete="one-time-code"
                    autoFocus
                  />
                </div>
              </div>

              {error && (
                <p className="text-red-400 text-sm">{error}</p>
              )}

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 bg-accent hover:bg-accent-bright text-bg-base font-semibold py-2.5 rounded-lg shadow-accent-glow hover:shadow-[0_0_0_1px_rgba(132,204,22,0.6),0_4px_16px_rgba(132,204,22,0.4)] active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    {t.auth.verifyOtp}
                    <ArrowRight className="w-4 h-4" />
                  </>
                )}
              </button>

              <div className="flex items-center justify-between text-sm">
                <button
                  type="button"
                  onClick={() => {
                    setStep("phone");
                    setOtp("");
                    setError("");
                  }}
                  className="text-fg-muted hover:text-fg transition-colors"
                >
                  &larr; {t.auth.phoneLabel}
                </button>
                <button
                  type="button"
                  onClick={handleResendOtp}
                  disabled={resendTimer > 0 || isSubmitting}
                  className="text-accent hover:text-accent-bright disabled:text-fg-muted disabled:cursor-not-allowed transition-colors"
                >
                  {resendTimer > 0
                    ? `${t.auth.resendIn} ${resendTimer}s`
                    : t.auth.resendOtp}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
