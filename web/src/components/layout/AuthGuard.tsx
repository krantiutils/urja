"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";
import { Loader2 } from "lucide-react";

interface AuthGuardProps {
  children: React.ReactNode;
  locale: string;
}

export function AuthGuard({ children, locale }: AuthGuardProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace(`/${locale}/login`);
    }
  }, [isAuthenticated, isLoading, router, locale]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-md-background">
        <Loader2 className="w-8 h-8 text-md-primary animate-spin" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
