"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import type { User, UserType, JWTPayload, AuthTokens } from "@/types";
import { api, ApiRequestError } from "./api";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

interface AuthContextValue extends AuthState {
  login: (phone: string) => Promise<void>;
  verifyOtp: (phone: string, otp: string) => Promise<{ is_new_user: boolean; onboarding_completed: boolean }>;
  passwordLogin: (phone: string, password: string) => Promise<{ is_new_user: boolean; onboarding_completed: boolean }>;
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
    is_super_admin: payload.is_super_admin,
  };
}

async function resolveOrgDetails(user: User): Promise<User> {
  try {
    const profile = await api.getMyProfile();
    const orgs = profile.organizations ?? [];
    const resolved: User = {
      ...user,
      user_type: profile.user_type as UserType | undefined,
      onboarding_completed: profile.onboarding_completed,
    };
    if (orgs.length > 0) {
      resolved.org_id = user.org_id || orgs[0].org_id;
      resolved.org_name = orgs[0].org_name;
    }
    return resolved;
  } catch {
    // If profile fetch fails, return user as-is
  }
  return user;
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
        .then(async (tokens) => {
          localStorage.setItem("access_token", tokens.access_token);
          localStorage.setItem("refresh_token", tokens.refresh_token);
          const user = getUserFromToken(tokens.access_token);
          if (!user) {
            setState({ user: null, isLoading: false, isAuthenticated: false });
            return;
          }
          const resolved = await resolveOrgDetails(user);
          setState({ user: resolved, isLoading: false, isAuthenticated: true });
        })
        .catch(() => {
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          setState({ user: null, isLoading: false, isAuthenticated: false });
        });
    } else {
      const user = getUserFromToken(accessToken);
      if (!user) {
        setState({ user: null, isLoading: false, isAuthenticated: false });
        return;
      }
      resolveOrgDetails(user).then((resolved) => {
        setState({ user: resolved, isLoading: false, isAuthenticated: true });
      });
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
        if (user) {
          const resolved = await resolveOrgDetails(user);
          setState({ user: resolved, isLoading: false, isAuthenticated: true });
        } else {
          setState({ user: null, isLoading: false, isAuthenticated: false });
        }
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

  /** Establishes the session from a token pair, whichever route issued it. */
  const adoptTokens = useCallback(async (tokens: AuthTokens) => {
    localStorage.setItem("access_token", tokens.access_token);
    localStorage.setItem("refresh_token", tokens.refresh_token);
    const user = getUserFromToken(tokens.access_token);
    if (!user) {
      setState({ user: null, isLoading: false, isAuthenticated: false });
      return { is_new_user: false, onboarding_completed: false };
    }
    // Enrich user with onboarding status from token response
    const enrichedUser: User = {
      ...user,
      onboarding_completed: tokens.onboarding_completed,
    };
    const resolved = await resolveOrgDetails(enrichedUser);
    setState({ user: resolved, isLoading: false, isAuthenticated: true });
    return {
      is_new_user: tokens.is_new_user,
      onboarding_completed: tokens.onboarding_completed,
    };
  }, []);

  const verifyOtp = useCallback(
    async (phone: string, otp: string) => adoptTokens(await api.verifyOtp(phone, otp)),
    [adoptTokens]
  );

  const passwordLogin = useCallback(
    async (phone: string, password: string) =>
      adoptTokens(await api.passwordLogin(phone, password)),
    [adoptTokens]
  );

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
    () => ({ ...state, login, verifyOtp, passwordLogin, logout }),
    [state, login, verifyOtp, passwordLogin, logout]
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
