"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import { getDictionary } from "@/lib/i18n";
import { ApiRequestError } from "@/lib/api";
import type { Locale } from "@/types";
import { Phone, ShieldCheck, ArrowRight, Loader2 } from "lucide-react";

const OTP_RESEND_SECONDS = 60;

export default function LoginPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const router = useRouter();
  const { login, verifyOtp, isAuthenticated, isLoading: authLoading } = useAuth();

  const [step, setStep] = useState<"phone" | "otp">("phone");
  const [phone, setPhone] = useState("");
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [resendTimer, setResendTimer] = useState(0);

  // Redirect if already authenticated
  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      router.replace(`/${locale}/dashboard`);
    }
  }, [isAuthenticated, authLoading, router, locale]);

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
        await verifyOtp(cleaned, otp);
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
      <div className="min-h-screen flex items-center justify-center bg-md-background">
        <Loader2 className="w-8 h-8 text-md-primary animate-spin" />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-md-background px-4">
      {/* Background organic shapes */}
      <div className="fixed top-20 -left-32 w-[500px] h-[500px] rounded-full bg-md-primary-container blur-[150px] opacity-30 pointer-events-none" />
      <div className="fixed bottom-20 -right-20 w-[400px] h-[400px] rounded-full bg-md-tertiary-container blur-[120px] opacity-25 pointer-events-none" />

      <div className="relative w-full max-w-sm animate-fade-in">
        {/* Logo / Brand */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold tracking-tight text-md-on-surface">
            {t.common.appName}
          </h1>
          <p className="mt-2 text-md-on-surface-variant text-sm">
            {t.auth.subtitle}
          </p>
        </div>

        {/* Login Card */}
        <div className="bg-md-surface-container border border-md-outline-variant rounded-3xl p-6 shadow-md-1">
          <h2 className="text-lg font-semibold text-md-on-surface mb-6">
            {t.auth.signIn}
          </h2>

          {step === "phone" ? (
            <form onSubmit={handleRequestOtp} className="space-y-4">
              <div>
                <label
                  htmlFor="phone"
                  className="block text-xs font-medium tracking-wider text-md-on-surface-variant uppercase mb-2"
                >
                  {t.auth.phoneLabel}
                </label>
                <div className="relative">
                  <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-md-on-surface-variant" />
                  <input
                    id="phone"
                    type="tel"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t.auth.phonePlaceholder}
                    className="w-full bg-md-surface-container-lowest border border-md-outline rounded-xl py-2.5 pl-10 pr-4 text-md-on-surface placeholder:text-md-on-surface-variant/50 focus:outline-none focus:border-md-primary focus:ring-2 focus:ring-md-primary/20 transition-colors"
                    autoComplete="tel"
                    autoFocus
                  />
                </div>
              </div>

              {error && (
                <p className="text-md-error text-sm">{error}</p>
              )}

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 bg-md-primary hover:bg-md-primary/90 text-md-on-primary font-medium py-2.5 rounded-full shadow-md-1 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
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
            </form>
          ) : (
            <form onSubmit={handleVerifyOtp} className="space-y-4">
              <p className="text-sm text-md-on-surface-variant mb-2">
                {t.auth.otpSent}
              </p>

              <div>
                <label
                  htmlFor="otp"
                  className="block text-xs font-medium tracking-wider text-md-on-surface-variant uppercase mb-2"
                >
                  {t.auth.otpLabel}
                </label>
                <div className="relative">
                  <ShieldCheck className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-md-on-surface-variant" />
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
                    className="w-full bg-md-surface-container-lowest border border-md-outline rounded-xl py-2.5 pl-10 pr-4 text-md-on-surface placeholder:text-md-on-surface-variant/50 focus:outline-none focus:border-md-primary focus:ring-2 focus:ring-md-primary/20 transition-colors tracking-[0.3em] text-center"
                    autoComplete="one-time-code"
                    autoFocus
                  />
                </div>
              </div>

              {error && (
                <p className="text-md-error text-sm">{error}</p>
              )}

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 bg-md-primary hover:bg-md-primary/90 text-md-on-primary font-medium py-2.5 rounded-full shadow-md-1 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
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
                  className="text-md-on-surface-variant hover:text-md-on-surface transition-colors"
                >
                  &larr; {t.auth.phoneLabel}
                </button>
                <button
                  type="button"
                  onClick={handleResendOtp}
                  disabled={resendTimer > 0 || isSubmitting}
                  className="text-md-primary hover:text-md-primary/80 disabled:text-md-on-surface-variant disabled:cursor-not-allowed transition-colors"
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
