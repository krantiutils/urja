export type Locale = "en" | "ne";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: string;
}

export interface JWTPayload {
  sub: string;
  phone: string;
  role: string;
  org_id?: string;
  iat: number;
  exp: number;
  iss: string;
}

export interface User {
  id: string;
  phone: string;
  role: string;
  org_id?: string;
}

export interface ApiError {
  error: string;
}

export interface DashboardStats {
  totalMembers: number;
  activeMembers: number;
  todayAttendance: number;
  revenue: number;
}

export interface MemberSummary {
  id: string;
  name: string;
  phone: string;
  package: string;
  joinedAt: string;
  status: "active" | "expired" | "expiring";
}

export interface AttendanceRecord {
  id: string;
  memberName: string;
  checkInTime: string;
  method: "qr" | "nfc" | "manual";
}

export interface PackageSummary {
  name: string;
  count: number;
  revenue: number;
}
