"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import type { User, JWTPayload } from "@/types";
import { api, ApiRequestError } from "./api";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

interface AuthContextValue extends AuthState {
  login: (phone: string) => Promise<void>;
  verifyOtp: (phone: string, otp: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function decodeJWT(token: string): JWTPayload | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;
    const payload = JSON.parse(atob(parts[1]));
    return payload as JWTPayload;
  } catch {
    return null;
  }
}

function isTokenExpired(token: string): boolean {
  const payload = decodeJWT(token);
  if (!payload) return true;
  // 30 second buffer before actual expiry
  return Date.now() >= (payload.exp - 30) * 1000;
}

function getUserFromToken(token: string): User | null {
  const payload = decodeJWT(token);
  if (!payload) return null;
  return {
    id: payload.sub,
    phone: payload.phone,
    role: payload.role,
    org_id: payload.org_id,
  };
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    isLoading: true,
    isAuthenticated: false,
  });

  // Initialize from stored tokens
  useEffect(() => {
    const accessToken = localStorage.getItem("access_token");
    const refreshToken = localStorage.getItem("refresh_token");

    if (!accessToken || !refreshToken) {
      setState({ user: null, isLoading: false, isAuthenticated: false });
      return;
    }

    if (isTokenExpired(accessToken)) {
      // Try refreshing
      api
        .refreshToken(refreshToken)
        .then((tokens) => {
          localStorage.setItem("access_token", tokens.access_token);
          localStorage.setItem("refresh_token", tokens.refresh_token);
          const user = getUserFromToken(tokens.access_token);
          setState({ user, isLoading: false, isAuthenticated: !!user });
        })
        .catch(() => {
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          setState({ user: null, isLoading: false, isAuthenticated: false });
        });
    } else {
      const user = getUserFromToken(accessToken);
      setState({ user, isLoading: false, isAuthenticated: !!user });
    }
  }, []);

  // Auto-refresh before expiry
  useEffect(() => {
    if (!state.isAuthenticated) return;

    const accessToken = localStorage.getItem("access_token");
    if (!accessToken) return;

    const payload = decodeJWT(accessToken);
    if (!payload) return;

    // Refresh 60 seconds before expiry
    const msUntilRefresh = (payload.exp - 60) * 1000 - Date.now();
    if (msUntilRefresh <= 0) return;

    const timer = setTimeout(async () => {
      const refreshToken = localStorage.getItem("refresh_token");
      if (!refreshToken) return;

      try {
        const tokens = await api.refreshToken(refreshToken);
        localStorage.setItem("access_token", tokens.access_token);
        localStorage.setItem("refresh_token", tokens.refresh_token);
        const user = getUserFromToken(tokens.access_token);
        setState({ user, isLoading: false, isAuthenticated: !!user });
      } catch {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        setState({ user: null, isLoading: false, isAuthenticated: false });
      }
    }, msUntilRefresh);

    return () => clearTimeout(timer);
  }, [state.isAuthenticated, state.user]);

  const login = useCallback(async (phone: string) => {
    await api.login(phone);
  }, []);

  const verifyOtp = useCallback(async (phone: string, otp: string) => {
    const tokens = await api.verifyOtp(phone, otp);
    localStorage.setItem("access_token", tokens.access_token);
    localStorage.setItem("refresh_token", tokens.refresh_token);
    const user = getUserFromToken(tokens.access_token);
    setState({ user, isLoading: false, isAuthenticated: !!user });
  }, []);

  const logout = useCallback(async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      try {
        await api.logout(refreshToken);
      } catch (err) {
        // Logout even if server call fails (token might already be invalid)
        if (!(err instanceof ApiRequestError && err.status === 401)) {
          console.error("Logout error:", err);
        }
      }
    }
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    setState({ user: null, isLoading: false, isAuthenticated: false });
  }, []);

  const value = useMemo(
    () => ({ ...state, login, verifyOtp, logout }),
    [state, login, verifyOtp, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
