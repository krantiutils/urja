export type Locale = "en" | "ne";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: string;
  is_new_user: boolean;
  onboarding_completed: boolean;
}

export interface JWTPayload {
  sub: string;
  phone: string;
  role: string;
  org_id?: string;
  is_super_admin?: boolean;
  iat: number;
  exp: number;
  iss: string;
}

export type UserType = "gym_member" | "fitness_tracker" | "calorie_tracker";

export interface User {
  id: string;
  phone: string;
  /**
   * The global claim carried in the JWT. Vestigial: it is literally "member"
   * for every account. Real authority is per-organization — use `org_role`.
   */
  role: string;
  org_id?: string;
  org_name?: string;
  /**
   * The caller's role in `org_id`, read from their real memberships. This is
   * what decides whether somebody runs a gym or trains at one.
   */
  org_role?: string;
  /** The gym's subdomain label, for fetching per-gym content by slug. */
  org_slug?: string;
  user_type?: UserType;
  onboarding_completed?: boolean;
  is_super_admin?: boolean;
  /**
   * True once memberships have actually been fetched. Distinguishes "this
   * account is not an operator" from "we do not know yet", which matters
   * because a failed profile request must not eject a genuine admin.
   */
  profile_loaded?: boolean;
}

export interface ApiError {
  error: string;
  code?: string;
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
  user_type: UserType;
  onboarding_completed: boolean;
  goal_type: string | null;
  daily_water_goal_ml: number;
  daily_step_goal: number;
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
  org_slug: string;
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
  organization_id: string | null;
  exercises: unknown[];
  duration_minutes: number | null;
  notes: string;
  logged_at: string;
}

export interface WorkoutPlan {
  id: string;
  user_id: string;
  organization_id: string | null;
  workout_template_id: string;
  assigned_by: string | null;
  assignment_method: string;
  assigned_at: string;
  template: WorkoutTemplate | null;
}

export interface WorkoutTemplate {
  id: string;
  name: string;
  name_ne: string | null;
  description: string | null;
  description_ne?: string | null;
  category: string | null;
  difficulty: string | null;
  duration_minutes: number | null;
  exercises: Exercise[];
  is_preset?: boolean;
  organization_id?: string | null;
  goal?: string | null;
  days_per_week?: string | null;
  created_at?: string;
  updated_at?: string;
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

// --- Dues (from GET /api/v1/orgs/{orgId}/dues) ---

export interface Due {
  id: string;
  organization_id: string;
  user_id: string;
  member_name: string;
  member_phone: string;
  amount: string;
  due_date: string;
  description: string;
  status: "unpaid" | "paid" | "waived";
  paid_at: string | null;
  paid_amount: string | null;
  payment_method: string | null;
  payment_reference: string | null;
  created_at: string;
}

export interface DueList {
  data: Due[];
  total: number;
}

// --- Accounts / Transactions (from GET /api/v1/orgs/{orgId}/accounts) ---

export interface Transaction {
  id: string;
  organization_id: string;
  category: string;
  description: string;
  transaction_date: string;
  transaction_type: "income" | "expense";
  amount: string;
  payment_type: string;
  reference: string;
  entry_by: string;
  entry_by_name: string;
  created_at: string;
  updated_at: string;
}

export interface TransactionList {
  data: Transaction[];
  total: number;
}

export interface AccountsSummary {
  total_income: string;
  total_expenses: string;
  gross_profit: string;
  profit_percent: number;
}

// --- Notices (from GET /api/v1/orgs/{orgId}/notices) ---

export interface Notice {
  id: string;
  organization_id: string;
  title: string;
  title_ne: string | null;
  content: string;
  content_ne: string | null;
  created_by: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// --- Feedbacks (from GET /api/v1/orgs/{orgId}/feedbacks) ---

export interface Feedback {
  id: string;
  organization_id: string;
  user_id: string;
  member_name: string;
  avatar_url: string | null;
  message: string;
  created_at: string;
}

// --- NFC (from GET /api/v1/orgs/{orgId}/nfc-cards) ---

export interface NfcCard {
  id: string;
  org_id: string;
  card_number: string;
  assigned_member: string | null;
  full_name: string | null;
  status: "available" | "assigned" | "disabled";
  assigned_at: string | null;
  last_updated: string;
}

export interface NfcDevice {
  id: string;
  org_id: string;
  name: string;
  device_identifier: string;
  status: "online" | "offline";
  door_state: "locked" | "unlocked";
  last_sync_at: string | null;
  uptime_seconds: number;
  created_at: string;
  updated_at: string;
}

// --- Profile Update (PUT /api/v1/members/me) ---

export interface ProfileUpdateRequest {
  name?: string;
  name_ne?: string;
  email?: string;
  date_of_birth?: string;
  gender?: string;
  avatar_url?: string;
  emergency_contact_name?: string;
  emergency_contact_phone?: string;
}

export interface PrivacySettingsUpdate {
  show_email: boolean;
  show_phone: boolean;
  show_profile: boolean;
  show_attendance: boolean;
  show_on_leaderboard: boolean;
}

// --- Organization (from GET /api/v1/gyms/{id}) ---

export interface Organization {
  id: string;
  name: string;
  name_ne: string;
  slug: string;
  description: string;
  description_ne: string;
  logo_url: string;
  address: string;
  address_ne: string;
  phone: string;
  email: string;
  latitude: number | null;
  longitude: number | null;
  timezone: string;
  settings: Record<string, unknown>;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  pan_number?: string;
  is_vat_registered?: boolean;
  tax_legal_name?: string;
  tax_address?: string;
}

export interface CreateOrganizationRequest {
  name: string;
  slug?: string;
  address?: string;
  phone?: string;
  email?: string;
}

export interface UpdateOrganizationRequest {
  name?: string;
  name_ne?: string;
  description?: string;
  description_ne?: string;
  logo_url?: string;
  address?: string;
  address_ne?: string;
  phone?: string;
  email?: string;
  latitude?: number;
  longitude?: number;
  settings?: Record<string, unknown>;
  pan_number?: string;
  is_vat_registered?: boolean;
  tax_legal_name?: string;
  tax_address?: string;
}

// --- SMS (from GET /api/v1/orgs/{orgId}/sms) ---

export interface SmsBalance {
  organization_id: string;
  balance: number;
  total_purchased: number;
  total_used: number;
  total_campaigns: number;
}

export interface SmsPurchase {
  id: string;
  organization_id: string;
  quantity: number;
  rate: string;
  amount: string;
  payment_method: string;
  payment_status: string;
  purchased_by: string;
  created_at: string;
}

// --- Dashboard Summary Endpoints ---

export interface PackageSummaryItem {
  id: string;
  name: string;
  member_count: number;
  price: string;
  duration_days: number;
}

export interface ExpiringPackageEntry {
  member_name: string;
  member_phone: string;
  package_name: string;
  expiry_date: string;
  expires_in: number;
}

export interface ExpiredPackageEntry {
  member_name: string;
  member_phone: string;
  package_name: string;
  expiry_date: string;
  expired_ago: number;
}

// --- Exercise (typed, with muscle groups) ---

export interface Exercise {
  name: string;
  name_ne?: string;
  sets: number;
  reps: number;
  rest_seconds: number;
  notes?: string;
  muscle_groups?: string[];
}

// --- Workout Planner Additions ---

export interface WorkoutQuestionnaireInput {
  organization_id: string;
  goal: "lose_weight" | "build_muscle" | "stay_fit";
  experience: "beginner" | "intermediate" | "advanced";
  days_per_week: "2-3" | "4-5" | "6-7";
}

// --- Nutrition Types ---

export interface FoodItem {
  id: string;
  organization_id: string | null;
  name: string;
  name_ne: string;
  category: string;
  calories_per_100g: number;
  protein_per_100g: number;
  carbs_per_100g: number;
  fat_per_100g: number;
  fiber_per_100g: number;
  serving_size_g: number;
  serving_label: string;
  serving_label_ne: string;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface FoodLog {
  id: string;
  user_id: string;
  organization_id: string | null;
  food_item_id: string;
  food_item?: FoodItem;
  meal_type: "breakfast" | "lunch" | "dinner" | "snack";
  quantity_grams: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  logged_date: string;
  logged_at: string;
  notes: string;
}

export interface MealTemplate {
  id: string;
  user_id: string;
  organization_id: string | null;
  name: string;
  name_ne: string;
  meal_type: "breakfast" | "lunch" | "dinner" | "snack";
  items: { food_item_id: string; quantity_grams: number }[];
  created_at: string;
}

export interface NutritionGoal {
  id: string;
  user_id: string;
  organization_id: string | null;
  calorie_goal: number;
  protein_goal_g: number;
  carbs_goal_g: number;
  fat_goal_g: number;
  weight_kg: number;
  height_cm: number;
  age: number;
  gender: string;
  activity_level: string;
  goal_type: string;
  created_at: string;
  updated_at: string;
}

export interface MealSummary {
  meal_type: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  items: FoodLog[];
}

export interface DailySummary {
  date: string;
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
  meals: MealSummary[];
}

export interface WeeklySummaryDay {
  date: string;
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
}

// --- Water Tracking ---

export interface WaterLog {
  id: string;
  user_id: string;
  amount_ml: number;
  logged_date: string;
  logged_at: string;
}

export interface WaterDailySummary {
  date: string;
  total_ml: number;
  goal_ml: number;
  entries: WaterLog[];
}

// --- Onboarding ---

export interface OnboardingInput {
  user_type: UserType;
  name: string;
  goal_type: "lose_weight" | "build_muscle" | "stay_fit" | "general_health";
  weight_kg?: number;
  height_cm?: number;
  age?: number;
  gender?: string;
  activity_level?: string;
}

// --- Exercise Library ---

export interface ExerciseItem {
  id: string;
  name: string;
  name_ne: string;
  description: string;
  muscle_groups: string[];
  equipment: string;
  exercise_type: string;
  difficulty: string;
  calories_per_minute: number;
  instructions: string;
  image_url: string;
  created_at: string;
}

export interface WorkoutProgram {
  id: string;
  name: string;
  name_ne: string;
  description: string;
  program_type: string;
  difficulty: string;
  duration_weeks: number;
  equipment: string;
  goal: string;
  image_url: string;
  created_at: string;
  days?: ProgramDay[];
}

export interface ProgramDay {
  id: string;
  program_id: string;
  week_number: number;
  day_number: number;
  day_name: string;
  exercises: unknown[];
  rest_day: boolean;
}

export interface UserProgramEnrollment {
  id: string;
  user_id: string;
  program_id: string;
  started_at: string;
  current_week: number;
  current_day: number;
  status: string;
  program?: WorkoutProgram;
  today_workout?: ProgramDay;
}

// --- Nutrition Streak ---

export interface NutritionStreak {
  user_id: string;
  current_streak: number;
  longest_streak: number;
  last_log_date: string | null;
  updated_at: string;
}

// --- Daily Dashboard ---

export interface DailyDashboard {
  date: string;
  calories: { consumed: number; goal: number; remaining: number };
  macros: {
    protein: { g: number; goal: number };
    carbs: { g: number; goal: number };
    fat: { g: number; goal: number };
  };
  water: { ml: number; goal: number };
  meals: MealSummary[];
  streak: { current: number; longest: number };
}

// --- Weight Trend ---

export interface WeightEntry {
  date: string;
  weight_kg: number;
}

export interface WeightTrendSummary {
  current_kg: number;
  start_kg: number;
  change_kg: number;
  bmi: number;
  entries_count: number;
}

export interface WeightTrend {
  entries: WeightEntry[];
  summary: WeightTrendSummary;
}

// --- Custom Food ---

export interface CreateCustomFoodInput {
  name: string;
  name_ne?: string;
  category?: string;
  calories_per_100g: number;
  protein_per_100g: number;
  carbs_per_100g: number;
  fat_per_100g: number;
  fiber_per_100g?: number;
  serving_size_g?: number;
  serving_label?: string;
  barcode?: string;
}


/** A package a member currently holds. Mirrors internal/subscription. */
export interface MemberSubscription {
  id: string;
  package_id: string;
  package_name: string;
  start_date: string;
  end_date: string;
  status: string;
  days_remaining: number;
  amount_paid: number;
  payment_method?: string;
  created_at: string;
}

export interface SubscriptionPayment {
  id: string;
  member_package_id: string;
  particular: string;
  total_amount: number;
  discount: number;
  paid_amount: number;
  due: number;
  payment_method?: string;
  payment_reference?: string;
  payment_date: string;
}

export interface AssignPackageRequest {
  package_id: string;
  start_date: string;
  payment_method: string;
  amount_paid: number;
  discount?: number;
  payment_reference?: string;
}


/** A member who has stopped turning up. Mirrors internal/absentee. */
export interface Absentee {
  user_id: string;
  name: string;
  avatar_url?: string;
  phone: string;
  absent_days: number;
  joined_date: string;
}


/** A training guide a gym writes for its members. Mirrors internal/guide. */
export interface TrainingGuide {
  id: string;
  title: string;
  title_ne?: string;
  content: string;
  content_ne?: string;
  category?: string;
  cover_image_url?: string;
  author_id?: string;
  is_published: boolean;
  created_at: string;
  updated_at: string;
}

export interface TrainingGuideInput {
  title: string;
  title_ne?: string;
  content: string;
  content_ne?: string;
  category?: string;
  cover_image_url?: string;
}

/** A line item on a PAN tax invoice. Mirrors internal/invoice. */
export interface InvoiceItem {
  line_no: number;
  description: string;
  description_ne?: string;
  quantity: number;
  unit_price: number;
  amount: number;
}

export interface Invoice {
  id: string;
  organization_id: string;
  fiscal_year: string;
  sequence: number;
  invoice_number: string;
  doc_type: "invoice" | "credit_note";
  credit_note_for?: string;
  /** Parent's invoice_number, snapshotted at credit time — the id above is a
   *  UUID the customer never sees; the printed document needs this. */
  credit_note_for_number?: string;

  seller_name: string;
  seller_pan: string;
  seller_address?: string;
  seller_vat_registered: boolean;

  customer_user_id?: string;
  customer_name: string;
  customer_pan?: string;
  customer_address?: string;
  customer_phone?: string;

  issued_date: string;
  issued_date_bs: string;

  subtotal: number;
  discount: number;
  taxable_amount: number;
  vat_rate: number;
  vat_amount: number;
  total: number;
  amount_in_words: string;

  payment_method?: string;
  status: "issued" | "cancelled";
  cancelled_at?: string;
  cancellation_reason?: string;

  transaction_id?: string;
  member_package_id?: string;

  issued_by: string;
  print_count: number;
  created_at: string;

  items?: InvoiceItem[];
}

export interface InvoiceItemInput {
  description: string;
  description_ne?: string;
  quantity: number;
  unit_price: number;
}

export interface IssueInvoiceInput {
  customer_user_id?: string;
  customer_name: string;
  customer_pan?: string;
  customer_address?: string;
  customer_phone?: string;
  payment_method?: string;
  transaction_id?: string;
  member_package_id?: string;
  discount?: number;
  items: InvoiceItemInput[];
}
