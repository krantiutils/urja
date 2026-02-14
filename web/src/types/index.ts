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

// --- Organization Member (from GET /api/v1/orgs/{orgId}/members) ---

export interface OrgMember {
  id: string;
  phone: string;
  name: string;
  name_ne: string | null;
  email: string | null;
  avatar_url: string | null;
  role: "member" | "staff" | "admin";
  status: "active" | "suspended" | "left";
  joined_at: string;
}

export interface OrgMemberList {
  data: OrgMember[];
  total: number;
}

export interface CreateMemberRequest {
  phone: string;
  name: string;
  name_ne?: string;
  email?: string;
  role?: "member" | "staff" | "admin";
}

export interface UpdateMemberRequest {
  role?: "member" | "staff" | "admin";
  status?: "active" | "suspended" | "left";
}

// --- Package (from GET /api/v1/orgs/{orgId}/packages) ---

export interface OrgPackage {
  id: string;
  organization_id: string;
  name: string;
  name_ne: string | null;
  description: string | null;
  description_ne: string | null;
  duration_days: number;
  price: string;
  currency: string;
  max_members: number | null;
  features: string[] | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface OrgPackageList {
  data: OrgPackage[];
}

export interface CreatePackageRequest {
  name: string;
  name_ne?: string;
  description?: string;
  description_ne?: string;
  duration_days: number;
  price: string;
  max_members?: number;
  features?: string[];
}

export interface UpdatePackageRequest {
  name?: string;
  name_ne?: string;
  description?: string;
  description_ne?: string;
  duration_days?: number;
  price?: string;
  max_members?: number;
  features?: string[];
  is_active?: boolean;
}

// --- Attendance (from GET /api/v1/orgs/{orgId}/attendance) ---

export interface OrgAttendance {
  id: string;
  user_id: string;
  org_id: string;
  check_in_at: string;
  method: "qr" | "nfc" | "manual";
  member_name?: string;
}

export interface OrgAttendanceList {
  data: OrgAttendance[];
}

export interface ManualCheckInRequest {
  member_id: string;
}

export interface CheckInResponse {
  check_in: OrgAttendance;
  streak: {
    id: string;
    member_id: string;
    org_id: string;
    current_streak: number;
    longest_streak: number;
    last_check_in: string;
    updated_at: string;
  };
}

// --- Staff (from GET /api/v1/orgs/{orgId}/staff) ---

export type StaffRole = "owner" | "manager" | "trainer" | "receptionist";

export interface StaffMember {
  id: string;
  user_id: string;
  name: string;
  phone: string;
  email: string | null;
  avatar_url: string | null;
  role: "staff" | "admin";
  staff_role: StaffRole;
  status: "active" | "suspended" | "left";
  joined_at: string;
}

export interface StaffList {
  data: StaffMember[];
  total: number;
}

export interface CreateStaffRequest {
  phone: string;
  name: string;
  email?: string;
  avatar_url?: string;
  staff_role: StaffRole;
}

export interface UpdateStaffRequest {
  name?: string;
  email?: string;
  avatar_url?: string;
  staff_role?: StaffRole;
}
