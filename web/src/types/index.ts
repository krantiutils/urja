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

// --- Member Profile (from GET /members/me) ---

export interface MemberProfile {
  id: string;
  phone: string;
  name: string;
  name_ne: string | null;
  email: string | null;
  date_of_birth: string | null;
  gender: string | null;
  avatar_url: string | null;
  emergency_contact_name: string | null;
  emergency_contact_phone: string | null;
  privacy_settings: PrivacySettings;
  organizations: OrgMembership[];
  created_at: string;
  updated_at: string;
}

export interface PrivacySettings {
  show_email: boolean;
  show_phone: boolean;
  show_profile: boolean;
  show_attendance: boolean;
  show_on_leaderboard: boolean;
}

export interface OrgMembership {
  org_id: string;
  org_name: string;
  role: string;
  status: string;
  joined_at: string;
}

export interface MemberAttendanceRecord {
  id: string;
  user_id: string;
  org_id: string;
  check_in_at: string;
  method: "qr" | "nfc" | "manual";
}

export interface MemberStreak {
  id: string;
  member_id: string;
  org_id: string;
  current_streak: number;
  longest_streak: number;
  last_check_in: string | null;
  updated_at: string;
}

export interface MemberPackage {
  id: string;
  user_id: string;
  package_id: string;
  organization_id: string;
  start_date: string;
  end_date: string;
  payment_method: string | null;
  payment_reference: string | null;
  amount_paid: string;
  status: "pending" | "active" | "expired" | "cancelled";
  created_at: string;
  package_name: string | null;
  org_name: string | null;
}

export interface HealthMetric {
  id: string;
  member_id: string;
  metric_type: string;
  value: Record<string, unknown>;
  recorded_at: string;
}

export interface WorkoutLog {
  id: string;
  user_id: string;
  workout_template_id: string | null;
  organization_id: string;
  exercises: unknown[];
  duration_minutes: number | null;
  notes: string;
  logged_at: string;
}

export interface WorkoutPlan {
  id: string;
  user_id: string;
  organization_id: string;
  workout_template_id: string;
  assigned_by: string | null;
  assigned_at: string;
  template: WorkoutTemplate | null;
}

export interface WorkoutTemplate {
  id: string;
  name: string;
  name_ne: string | null;
  description: string | null;
  category: string | null;
  difficulty: string | null;
  duration_minutes: number | null;
  exercises: unknown[];
}

// --- Attendance Calendar (from GET /api/v1/members/me/attendance/calendar) ---

export interface AttendanceCalendar {
  month: string;
  days_in_month: number;
  check_in_days: number[];
  total_check_ins: number;
}

// --- Leaderboard (from GET /api/v1/members/me/leaderboard) ---

export interface LeaderboardEntry {
  rank: number;
  member_id: string;
  name: string;
  avatar_url: string | null;
  value: number;
  metric: string;
}

export interface LeaderboardResponse {
  period: string;
  rankings: LeaderboardEntry[];
}
