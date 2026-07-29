import type {
  AuthTokens,
  ApiError,
  Organization,
  CreateOrganizationRequest,
  UpdateOrganizationRequest,
  OrgMemberList,
  OrgMember,
  CreateMemberRequest,
  UpdateMemberRequest,
  OrgPackageList,
  OrgPackage,
  CreatePackageRequest,
  UpdatePackageRequest,
  OrgAttendanceList,
  CheckInResponse,
  ManualCheckInRequest,
  StaffList,
  StaffMember,
  CreateStaffRequest,
  UpdateStaffRequest,
  ProfileUpdateRequest,
  PrivacySettingsUpdate,
  MemberProfile,
  MemberAttendanceRecord,
  MemberStreak,
  MemberPackage,
  HealthMetric,
  WorkoutLog,
  WorkoutPlan,
  WorkoutTemplate,
  WorkoutQuestionnaireInput,
  FoodItem,
  FoodLog,
  MealTemplate,
  NutritionGoal,
  DailySummary,
  WeeklySummaryDay,
  DueList,
  Due,
  TransactionList,
  AccountsSummary,
  Notice,
  Feedback,
  NfcCard,
  NfcDevice,
  SmsBalance,
  PackageSummaryItem,
  ExpiringPackageEntry,
  ExpiredPackageEntry,
  AttendanceCalendar,
  LeaderboardResponse,
  WaterLog,
  WaterDailySummary,
  OnboardingInput,
  ExerciseItem,
  WorkoutProgram,
  UserProgramEnrollment,
  NutritionStreak,
  DailyDashboard,
  WeightTrend,
  CreateCustomFoodInput,
  Absentee,
  TrainingGuide,
  TrainingGuideInput,
  MemberSubscription,
  SubscriptionPayment,
  AssignPackageRequest,
  Invoice,
  InvoiceItemInput,
  IssueInvoiceInput,
} from "@/types";
import type {
  SiteSettings,
  SitePage,
  PageSummary,
  SiteLead,
  LeadStatus,
} from "@/types/site";
import type {
  SiteTemplateOption,
  SitePageInput,
  BoxingProfileView,
  BoxingProfileInput,
  Bout,
  BoutInput,
} from "@/types/site-admin";

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${path}`;
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...((options.headers as Record<string, string>) ?? {}),
    };

    const accessToken =
      typeof window !== "undefined"
        ? localStorage.getItem("access_token")
        : null;
    if (accessToken) {
      headers["Authorization"] = `Bearer ${accessToken}`;
    }

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
      const body = (await response.json().catch(() => ({
        error: `HTTP ${response.status}`,
      }))) as ApiError;
      throw new ApiRequestError(body.error, response.status, body.code);
    }

    return response.json() as Promise<T>;
  }

  async login(phone: string): Promise<{ message: string }> {
    return this.request("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ phone }),
    });
  }

  /** Password sign-in, for staff and admins who log in several times a day. */
  async passwordLogin(phone: string, password: string): Promise<AuthTokens> {
    return this.request("/api/v1/auth/password-login", {
      method: "POST",
      body: JSON.stringify({ phone, password }),
    });
  }

  async setPassword(password: string): Promise<{ message: string }> {
    return this.request("/api/v1/auth/password", {
      method: "POST",
      body: JSON.stringify({ password }),
    });
  }

  async passwordStatus(): Promise<{ password_set: boolean }> {
    return this.request("/api/v1/auth/password");
  }

  async verifyOtp(phone: string, otp: string): Promise<AuthTokens> {
    return this.request("/api/v1/auth/verify-otp", {
      method: "POST",
      body: JSON.stringify({ phone, otp }),
    });
  }

  async refreshToken(refreshToken: string): Promise<AuthTokens> {
    return this.request("/api/v1/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  async logout(refreshToken: string): Promise<{ message: string }> {
    return this.request("/api/v1/auth/logout", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  // --- Organization ---

  async getOrg(orgId: string): Promise<Organization> {
    return this.request(`/api/v1/gyms/${orgId}`);
  }

  async updateOrg(
    orgId: string,
    data: UpdateOrganizationRequest
  ): Promise<Organization> {
    return this.request(`/api/v1/orgs/${orgId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  // --- Super Admin ---

  async listOrgs(
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: Organization[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs${qs ? `?${qs}` : ""}`);
  }

  async createOrg(data: CreateOrganizationRequest): Promise<Organization> {
    return this.request("/api/v1/orgs", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // --- Members ---

  async listMembers(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<OrgMemberList> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/members${qs ? `?${qs}` : ""}`);
  }

  async createMember(
    orgId: string,
    data: CreateMemberRequest
  ): Promise<OrgMember> {
    return this.request(`/api/v1/orgs/${orgId}/members`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateMember(
    orgId: string,
    memberId: string,
    data: UpdateMemberRequest
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteMember(
    orgId: string,
    memberId: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}`, {
      method: "DELETE",
    });
  }

  // --- Packages ---

  async listPackages(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<OrgPackageList> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/packages${qs ? `?${qs}` : ""}`);
  }

  async createPackage(
    orgId: string,
    data: CreatePackageRequest
  ): Promise<OrgPackage> {
    return this.request(`/api/v1/orgs/${orgId}/packages`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updatePackage(
    orgId: string,
    packageId: string,
    data: UpdatePackageRequest
  ): Promise<OrgPackage> {
    return this.request(`/api/v1/orgs/${orgId}/packages/${packageId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deletePackage(
    orgId: string,
    packageId: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/packages/${packageId}`, {
      method: "DELETE",
    });
  }

  // --- Attendance ---

  async listAttendance(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<OrgAttendanceList> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(
      `/api/v1/orgs/${orgId}/attendance${qs ? `?${qs}` : ""}`
    );
  }

  async manualCheckIn(
    orgId: string,
    data: ManualCheckInRequest
  ): Promise<CheckInResponse> {
    return this.request(`/api/v1/orgs/${orgId}/attendance/check-in`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // --- Staff ---

  async listStaff(
    orgId: string,
    params: { search?: string; limit?: number; offset?: number } = {}
  ): Promise<StaffList> {
    const q = new URLSearchParams();
    if (params.search) q.set("search", params.search);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/staff${qs ? `?${qs}` : ""}`);
  }

  async createStaff(
    orgId: string,
    data: CreateStaffRequest
  ): Promise<StaffMember> {
    return this.request(`/api/v1/orgs/${orgId}/staff`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getStaff(orgId: string, staffId: string): Promise<StaffMember> {
    return this.request(`/api/v1/orgs/${orgId}/staff/${staffId}`);
  }

  async updateStaff(
    orgId: string,
    staffId: string,
    data: UpdateStaffRequest
  ): Promise<StaffMember> {
    return this.request(`/api/v1/orgs/${orgId}/staff/${staffId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteStaff(
    orgId: string,
    staffId: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/staff/${staffId}`, {
      method: "DELETE",
    });
  }

  // --- Member (self) ---

  async getMyProfile(): Promise<MemberProfile> {
    return this.request("/api/v1/members/me");
  }

  async updateMyProfile(
    data: ProfileUpdateRequest
  ): Promise<{ message: string }> {
    return this.request("/api/v1/members/me", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async updateMyPrivacy(
    data: PrivacySettingsUpdate
  ): Promise<{ message: string }> {
    return this.request("/api/v1/members/me/privacy", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteMyAccount(): Promise<{ message: string }> {
    return this.request("/api/v1/members/me", {
      method: "DELETE",
    });
  }

  async getMyAttendance(
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: MemberAttendanceRecord[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/attendance${qs ? `?${qs}` : ""}`);
  }

  async getMyStreaks(): Promise<{ data: MemberStreak[] }> {
    return this.request("/api/v1/members/me/streaks");
  }

  async getMyPackages(
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: MemberPackage[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/packages${qs ? `?${qs}` : ""}`);
  }

  async getMyHealth(
    params: { type?: string; limit?: number } = {}
  ): Promise<{ data: HealthMetric[] }> {
    const q = new URLSearchParams();
    if (params.type) q.set("type", params.type);
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/health${qs ? `?${qs}` : ""}`);
  }

  async getMyWorkoutLogs(
    params: { organization_id?: string; limit?: number; offset?: number } = {}
  ): Promise<{ data: WorkoutLog[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/workout-logs${qs ? `?${qs}` : ""}`);
  }

  async getMyWorkoutPlan(orgId?: string): Promise<WorkoutPlan> {
    const q = new URLSearchParams();
    if (orgId) q.set("organization_id", orgId);
    const qs = q.toString();
    return this.request(`/api/v1/members/me/workout-plan${qs ? `?${qs}` : ""}`);
  }

  // --- Workout Planner ---

  async browseTemplates(params: {
    organization_id?: string;
    goal?: string;
    difficulty?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ data: WorkoutTemplate[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    if (params.goal) q.set("goal", params.goal);
    if (params.difficulty) q.set("difficulty", params.difficulty);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    return this.request(`/api/v1/members/me/workout-templates?${q.toString()}`);
  }

  async selfAssignPlan(data: {
    organization_id?: string;
    workout_template_id: string;
  }): Promise<WorkoutPlan> {
    return this.request("/api/v1/members/me/workout-plan", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async recommendPlan(data: WorkoutQuestionnaireInput): Promise<WorkoutPlan> {
    return this.request("/api/v1/members/me/workout-plan/recommend", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // --- Nutrition ---

  async searchFoods(params: {
    organization_id?: string;
    q?: string;
    category?: string;
    limit?: number;
  }): Promise<{ data: FoodItem[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    if (params.q) q.set("q", params.q);
    if (params.category) q.set("category", params.category);
    if (params.limit) q.set("limit", String(params.limit));
    return this.request(`/api/v1/members/me/foods?${q.toString()}`);
  }

  async getFoodItem(id: string): Promise<FoodItem> {
    return this.request(`/api/v1/members/me/foods/${id}`);
  }

  async logFood(data: {
    organization_id?: string;
    food_item_id: string;
    meal_type: string;
    quantity_grams: number;
    logged_date: string;
    notes?: string;
  }): Promise<FoodLog> {
    return this.request("/api/v1/members/me/food-logs", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getMyFoodLogs(params: {
    organization_id?: string;
    date?: string;
    meal_type?: string;
  }): Promise<{ data: FoodLog[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    if (params.date) q.set("date", params.date);
    if (params.meal_type) q.set("meal_type", params.meal_type);
    return this.request(`/api/v1/members/me/food-logs?${q.toString()}`);
  }

  async deleteFoodLog(id: string): Promise<{ message: string }> {
    return this.request(`/api/v1/members/me/food-logs/${id}`, {
      method: "DELETE",
    });
  }

  async getDailySummary(params: {
    organization_id?: string;
    date: string;
  }): Promise<DailySummary> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    q.set("date", params.date);
    return this.request(`/api/v1/members/me/nutrition/summary?${q.toString()}`);
  }

  async getWeeklySummary(params: {
    organization_id?: string;
    from: string;
  }): Promise<{ data: WeeklySummaryDay[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    q.set("from", params.from);
    return this.request(`/api/v1/members/me/nutrition/weekly?${q.toString()}`);
  }

  async setNutritionGoal(data: {
    organization_id?: string;
    weight_kg: number;
    height_cm: number;
    age: number;
    gender: string;
    activity_level: string;
    goal_type: string;
  }): Promise<NutritionGoal> {
    return this.request("/api/v1/members/me/nutrition/goal", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getNutritionGoal(orgId?: string): Promise<NutritionGoal> {
    const q = new URLSearchParams();
    if (orgId) q.set("organization_id", orgId);
    const qs = q.toString();
    return this.request(`/api/v1/members/me/nutrition/goal${qs ? `?${qs}` : ""}`);
  }

  async createMealTemplate(data: {
    organization_id?: string;
    name: string;
    name_ne?: string;
    meal_type: string;
    items: { food_item_id: string; quantity_grams: number }[];
  }): Promise<MealTemplate> {
    return this.request("/api/v1/members/me/meal-templates", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getMyMealTemplates(orgId?: string): Promise<{ data: MealTemplate[] }> {
    const q = new URLSearchParams();
    if (orgId) q.set("organization_id", orgId);
    const qs = q.toString();
    return this.request(`/api/v1/members/me/meal-templates${qs ? `?${qs}` : ""}`);
  }

  async deleteMealTemplate(id: string): Promise<{ message: string }> {
    return this.request(`/api/v1/members/me/meal-templates/${id}`, {
      method: "DELETE",
    });
  }

  async logMealTemplate(
    id: string,
    data: { organization_id?: string; logged_date: string }
  ): Promise<{ data: FoodLog[] }> {
    return this.request(`/api/v1/members/me/meal-templates/${id}/log`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getMyAttendanceCalendar(
    params: { month?: string } = {}
  ): Promise<AttendanceCalendar> {
    const q = new URLSearchParams();
    if (params.month) q.set("month", params.month);
    const qs = q.toString();
    return this.request(`/api/v1/members/me/attendance/calendar${qs ? `?${qs}` : ""}`);
  }

  async getMyLeaderboard(
    params: { period?: string; limit?: number } = {}
  ): Promise<LeaderboardResponse> {
    const q = new URLSearchParams();
    if (params.period) q.set("period", params.period);
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/leaderboard${qs ? `?${qs}` : ""}`);
  }

  async submitFeedback(
    orgId: string,
    data: { message: string }
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/feedbacks`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // --- Dues ---

  async listDues(
    orgId: string,
    params: { status?: string; search?: string; limit?: number; offset?: number } = {}
  ): Promise<DueList> {
    const q = new URLSearchParams();
    if (params.status) q.set("status", params.status);
    if (params.search) q.set("search", params.search);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/dues${qs ? `?${qs}` : ""}`);
  }

  /** Records that a member owes money. */
  async createDue(
    orgId: string,
    data: { user_id: string; amount: number; due_date: string; description?: string }
  ): Promise<Due> {
    return this.request(`/api/v1/orgs/${orgId}/dues`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async payDue(
    orgId: string,
    dueId: string,
    data: { amount: number; payment_method: string; payment_reference?: string }
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/dues/${dueId}/pay`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // --- Accounts / Transactions ---

  async listTransactions(
    orgId: string,
    params: { category?: string; transaction_type?: string; from?: string; to?: string; limit?: number; offset?: number } = {}
  ): Promise<TransactionList> {
    const q = new URLSearchParams();
    if (params.category) q.set("category", params.category);
    if (params.transaction_type) q.set("transaction_type", params.transaction_type);
    if (params.from) q.set("from", params.from);
    if (params.to) q.set("to", params.to);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/accounts${qs ? `?${qs}` : ""}`);
  }

  async createTransaction(
    orgId: string,
    data: { category: string; description: string; transaction_date: string; transaction_type: string; amount: number; payment_type: string; reference?: string }
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/accounts`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateTransaction(
    orgId: string,
    id: string,
    data: { category?: string; description?: string; transaction_date?: string; transaction_type?: string; amount?: number; payment_type?: string; reference?: string }
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/accounts/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteTransaction(
    orgId: string,
    id: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/accounts/${id}`, {
      method: "DELETE",
    });
  }

  async getAccountsSummary(
    orgId: string,
    params: { from?: string; to?: string } = {}
  ): Promise<AccountsSummary> {
    const q = new URLSearchParams();
    if (params.from) q.set("from", params.from);
    if (params.to) q.set("to", params.to);
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/accounts/summary${qs ? `?${qs}` : ""}`);
  }

  // --- Notices ---

  async listNotices(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: Notice[]; total: number }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/notices${qs ? `?${qs}` : ""}`);
  }

  async createNotice(
    orgId: string,
    data: { title: string; title_ne?: string; content: string; content_ne?: string; is_active?: boolean }
  ): Promise<Notice> {
    return this.request(`/api/v1/orgs/${orgId}/notices`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateNotice(
    orgId: string,
    id: string,
    data: { title?: string; title_ne?: string; content?: string; content_ne?: string; is_active?: boolean }
  ): Promise<Notice> {
    return this.request(`/api/v1/orgs/${orgId}/notices/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteNotice(
    orgId: string,
    id: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/notices/${id}`, {
      method: "DELETE",
    });
  }

  // --- Feedbacks ---

  async listFeedbacks(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: Feedback[]; total: number }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/feedbacks${qs ? `?${qs}` : ""}`);
  }

  async deleteFeedback(
    orgId: string,
    id: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/feedbacks/${id}`, {
      method: "DELETE",
    });
  }

  // --- NFC ---

  async listNfcCards(
    orgId: string,
    params: { status?: string; limit?: number; offset?: number } = {}
  ): Promise<{ data: NfcCard[]; total: number }> {
    const q = new URLSearchParams();
    if (params.status) q.set("status", params.status);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/nfc-cards${qs ? `?${qs}` : ""}`);
  }

  async registerNfcCard(
    orgId: string,
    data: { card_number: string }
  ): Promise<NfcCard> {
    return this.request(`/api/v1/orgs/${orgId}/nfc-cards`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async assignNfcCard(
    orgId: string,
    cardId: string,
    userId: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/nfc-cards/${cardId}/assign`, {
      method: "PUT",
      body: JSON.stringify({ user_id: userId }),
    });
  }

  async unassignNfcCard(
    orgId: string,
    cardId: string
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/nfc-cards/${cardId}/unassign`, {
      method: "PUT",
    });
  }

  async listNfcDevices(
    orgId: string
  ): Promise<{ data: NfcDevice[] }> {
    return this.request(`/api/v1/orgs/${orgId}/nfc-devices`);
  }

  // --- SMS ---

  async getSmsBalance(orgId: string): Promise<SmsBalance> {
    return this.request(`/api/v1/orgs/${orgId}/sms/balance`);
  }

  async sendSms(
    orgId: string,
    data: { message: string; member_ids: string[] }
  ): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/sms/send`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getSmsHistory(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: unknown[]; total: number }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/sms/history${qs ? `?${qs}` : ""}`);
  }

  // --- Dashboard Summary ---

  async getPackageSummary(
    orgId: string
  ): Promise<PackageSummaryItem[]> {
    return this.request(`/api/v1/orgs/${orgId}/packages/summary`);
  }

  async getExpiringPackages(
    orgId: string,
    params: { days?: number; limit?: number; offset?: number } = {}
  ): Promise<{ data: ExpiringPackageEntry[] }> {
    const q = new URLSearchParams();
    if (params.days) q.set("days", String(params.days));
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/packages/expiring${qs ? `?${qs}` : ""}`);
  }

  async getExpiredPackages(
    orgId: string,
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: ExpiredPackageEntry[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/packages/expired${qs ? `?${qs}` : ""}`);
  }

  // --- Onboarding ---

  async submitOnboarding(data: OnboardingInput): Promise<MemberProfile> {
    return this.request("/api/v1/members/me/onboarding", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async joinGym(orgId: string): Promise<{ message: string }> {
    return this.request("/api/v1/members/me/join-gym", {
      method: "POST",
      body: JSON.stringify({ org_id: orgId }),
    });
  }

  async leaveGym(orgId: string): Promise<{ message: string }> {
    return this.request("/api/v1/members/me/leave-gym", {
      method: "POST",
      body: JSON.stringify({ org_id: orgId }),
    });
  }

  // --- Water Tracking ---

  async logWater(data: {
    amount_ml: number;
    date?: string;
  }): Promise<WaterLog> {
    return this.request("/api/v1/members/me/water-logs", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getWaterLogs(date: string): Promise<{ data: WaterLog[] }> {
    return this.request(`/api/v1/members/me/water-logs?date=${date}`);
  }

  async deleteWaterLog(id: string): Promise<{ message: string }> {
    return this.request(`/api/v1/members/me/water-logs/${id}`, {
      method: "DELETE",
    });
  }

  async getWaterSummary(date: string): Promise<WaterDailySummary> {
    return this.request(`/api/v1/members/me/water-logs/summary?date=${date}`);
  }

  // --- Exercise Library ---

  async listExercises(params?: {
    type?: string;
    equipment?: string;
    muscle_group?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ data: ExerciseItem[] }> {
    const q = new URLSearchParams();
    if (params?.type) q.set("type", params.type);
    if (params?.equipment) q.set("equipment", params.equipment);
    if (params?.muscle_group) q.set("muscle_group", params.muscle_group);
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/exercises${qs ? `?${qs}` : ""}`);
  }

  async getExercise(id: string): Promise<ExerciseItem> {
    return this.request(`/api/v1/members/me/exercises/${id}`);
  }

  async listPrograms(params?: {
    type?: string;
    difficulty?: string;
    equipment?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ data: WorkoutProgram[] }> {
    const q = new URLSearchParams();
    if (params?.type) q.set("type", params.type);
    if (params?.difficulty) q.set("difficulty", params.difficulty);
    if (params?.equipment) q.set("equipment", params.equipment);
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/programs${qs ? `?${qs}` : ""}`);
  }

  async getProgram(id: string): Promise<WorkoutProgram> {
    return this.request(`/api/v1/members/me/programs/${id}`);
  }

  async enrollInProgram(programId: string): Promise<UserProgramEnrollment> {
    return this.request(`/api/v1/members/me/programs/${programId}/enroll`, {
      method: "POST",
    });
  }

  async getCurrentProgram(): Promise<UserProgramEnrollment> {
    return this.request(`/api/v1/members/me/programs/current`);
  }

  async completeProgramDay(): Promise<UserProgramEnrollment> {
    return this.request(`/api/v1/members/me/programs/current/complete-day`, {
      method: "POST",
    });
  }

  // --- Custom Foods + Barcode ---

  async createCustomFood(input: CreateCustomFoodInput): Promise<FoodItem> {
    return this.request(`/api/v1/members/me/foods`, {
      method: "POST",
      body: JSON.stringify(input),
    });
  }

  async lookupBarcode(barcode: string): Promise<FoodItem> {
    return this.request(`/api/v1/members/me/foods/barcode/${barcode}`);
  }

  // --- Nutrition Streak ---

  async getNutritionStreak(): Promise<NutritionStreak> {
    return this.request(`/api/v1/members/me/nutrition/streak`);
  }

  // --- Daily Dashboard ---

  async getDailyDashboard(params?: {
    organization_id?: string;
    date?: string;
  }): Promise<DailyDashboard> {
    const q = new URLSearchParams();
    if (params?.organization_id) q.set("organization_id", params.organization_id);
    if (params?.date) q.set("date", params.date);
    const qs = q.toString();
    return this.request(`/api/v1/members/me/nutrition/daily-dashboard${qs ? `?${qs}` : ""}`);
  }

  // --- Weight Trend ---

  async quickLogWeight(weightKg: number): Promise<HealthMetric> {
    return this.request(`/api/v1/members/me/health/weight`, {
      method: "POST",
      body: JSON.stringify({ weight_kg: weightKg }),
    });
  }

  async getWeightTrend(days?: number): Promise<WeightTrend> {
    const q = days ? `?days=${days}` : "";
    return this.request(`/api/v1/members/me/health/weight-trend${q}`);
  }

  // --- QR Code ---

  getQrCodeUrl(orgId: string): string {
    return `${this.baseUrl}/api/v1/orgs/${orgId}/qr-code`;
  }

  /** Tops up SMS credits. The server prices the purchase; only the count is ours to choose. */
  async buySmsCredits(
    orgId: string,
    quantity: number,
    paymentMethod: string
  ): Promise<{ quantity: number; amount: number; rate: number }> {
    return this.request(`/api/v1/orgs/${orgId}/sms/buy`, {
      method: "POST",
      body: JSON.stringify({ quantity, payment_method: paymentMethod }),
    });
  }

  // --- Training guides ---

  /** Published guides for one gym. The gym is required — guides are per-gym content. */
  async listPublishedGuides(
    gym: string,
    params: { category?: string; search?: string; limit?: number } = {}
  ): Promise<{ data: TrainingGuide[]; total: number }> {
    const q = new URLSearchParams({ gym });
    if (params.category) q.set("category", params.category);
    if (params.search) q.set("search", params.search);
    if (params.limit) q.set("limit", String(params.limit));
    return this.request(`/api/v1/training-guides?${q.toString()}`);
  }

  async getPublishedGuide(id: string): Promise<TrainingGuide> {
    return this.request(`/api/v1/training-guides/${id}`);
  }

  async listGuides(
    orgId: string,
    params: { category?: string; search?: string; limit?: number } = {}
  ): Promise<{ data: TrainingGuide[]; total: number }> {
    const q = new URLSearchParams();
    if (params.category) q.set("category", params.category);
    if (params.search) q.set("search", params.search);
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/training-guides${qs ? `?${qs}` : ""}`);
  }

  async createGuide(orgId: string, data: TrainingGuideInput): Promise<TrainingGuide> {
    return this.request(`/api/v1/orgs/${orgId}/training-guides`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateGuide(
    orgId: string,
    id: string,
    data: TrainingGuideInput
  ): Promise<TrainingGuide> {
    return this.request(`/api/v1/orgs/${orgId}/training-guides/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async publishGuide(
    orgId: string,
    id: string,
    isPublished: boolean
  ): Promise<TrainingGuide> {
    return this.request(`/api/v1/orgs/${orgId}/training-guides/${id}/publish`, {
      method: "PATCH",
      body: JSON.stringify({ is_published: isPublished }),
    });
  }

  async deleteGuide(orgId: string, id: string): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/training-guides/${id}`, {
      method: "DELETE",
    });
  }

  // --- Invoices (PAN tax billing) ---

  async listInvoices(
    orgId: string,
    params: {
      status?: string;
      fiscal_year?: string;
      q?: string;
      limit?: number;
      offset?: number;
    } = {}
  ): Promise<{ data: Invoice[]; total: number }> {
    const q = new URLSearchParams();
    if (params.status) q.set("status", params.status);
    if (params.fiscal_year) q.set("fiscal_year", params.fiscal_year);
    if (params.q) q.set("q", params.q);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/invoices${qs ? `?${qs}` : ""}`);
  }

  async getInvoice(orgId: string, id: string): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}`);
  }

  async issueInvoice(orgId: string, data: IssueInvoiceInput): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async cancelInvoice(orgId: string, id: string, reason: string): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/cancel`, {
      method: "POST",
      body: JSON.stringify({ reason }),
    });
  }

  async creditNote(
    orgId: string,
    id: string,
    data: { reason: string; items: InvoiceItemInput[] }
  ): Promise<Invoice> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/credit-note`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async printInvoice(
    orgId: string,
    id: string
  ): Promise<{ invoice: Invoice; copy_label: "original" | "copy" }> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/${id}/print`, {
      method: "POST",
    });
  }

  async nextInvoiceNumber(orgId: string): Promise<{ invoice_number: string }> {
    return this.request(`/api/v1/orgs/${orgId}/invoices/next-number`);
  }

  // --- Absentees (members who have stopped coming) ---

  async listAbsentees(
    orgId: string,
    params: { days?: number; limit?: number } = {}
  ): Promise<{ data: Absentee[] }> {
    const q = new URLSearchParams();
    if (params.days) q.set("days", String(params.days));
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/absentees${qs ? `?${qs}` : ""}`);
  }

  async notifyAbsentee(orgId: string, memberId: string): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/absentees/${memberId}/notify`, {
      method: "POST",
    });
  }

  // --- Memberships (selling a package to a member) ---

  async assignPackage(
    orgId: string,
    memberId: string,
    data: AssignPackageRequest
  ): Promise<{ subscription: MemberSubscription; payment: SubscriptionPayment }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}/packages/assign`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async renewPackage(
    orgId: string,
    memberId: string,
    data: {
      member_package_id: string;
      payment_method: string;
      amount_paid: number;
      discount?: number;
      payment_reference?: string;
    }
  ): Promise<{ subscription: MemberSubscription; payment: SubscriptionPayment }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}/packages/renew`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async listMemberSubscriptions(
    orgId: string,
    memberId: string
  ): Promise<{ data: MemberSubscription[] }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}/subscriptions`);
  }

  async listMemberPayments(
    orgId: string,
    memberId: string
  ): Promise<{ data: SubscriptionPayment[] }> {
    return this.request(`/api/v1/orgs/${orgId}/members/${memberId}/payments`);
  }

  // --- Tenant website ---

  async listSiteTemplates(): Promise<{ data: SiteTemplateOption[] }> {
    return this.request(`/api/v1/sites/templates`);
  }

  async getSiteSettings(orgId: string): Promise<SiteSettings> {
    return this.request(`/api/v1/orgs/${orgId}/site/settings`);
  }

  async updateSiteSettings(
    orgId: string,
    data: Partial<
      Pick<SiteSettings, "template" | "theme" | "nav" | "footer" | "socials" | "is_live">
    >
  ): Promise<SiteSettings> {
    return this.request(`/api/v1/orgs/${orgId}/site/settings`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  /**
   * Replaces every page with the template's. Destructive, so the API requires
   * `confirm` explicitly rather than inferring intent from the call.
   */
  async applySiteTemplate(
    orgId: string,
    template: string
  ): Promise<{ data: PageSummary[] }> {
    return this.request(`/api/v1/orgs/${orgId}/site/apply-template`, {
      method: "POST",
      body: JSON.stringify({ template, confirm: true }),
    });
  }

  async listSitePages(orgId: string): Promise<{ data: PageSummary[] }> {
    return this.request(`/api/v1/orgs/${orgId}/site/pages`);
  }

  async getSitePage(orgId: string, pageId: string): Promise<SitePage> {
    return this.request(`/api/v1/orgs/${orgId}/site/pages/${pageId}`);
  }

  async createSitePage(orgId: string, data: SitePageInput): Promise<SitePage> {
    return this.request(`/api/v1/orgs/${orgId}/site/pages`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateSitePage(
    orgId: string,
    pageId: string,
    data: SitePageInput
  ): Promise<SitePage> {
    return this.request(`/api/v1/orgs/${orgId}/site/pages/${pageId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteSitePage(orgId: string, pageId: string): Promise<{ message: string }> {
    return this.request(`/api/v1/orgs/${orgId}/site/pages/${pageId}`, {
      method: "DELETE",
    });
  }

  async listSiteLeads(
    orgId: string,
    params: { status?: string; limit?: number } = {}
  ): Promise<{ data: SiteLead[] }> {
    const q = new URLSearchParams();
    if (params.status) q.set("status", params.status);
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/orgs/${orgId}/site/leads${qs ? `?${qs}` : ""}`);
  }

  async updateSiteLead(
    orgId: string,
    leadId: string,
    status: LeadStatus
  ): Promise<SiteLead> {
    return this.request(`/api/v1/orgs/${orgId}/site/leads/${leadId}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
  }

  // --- Boxing ---

  async getMyBoxingProfile(orgId: string): Promise<BoxingProfileView> {
    return this.request(`/api/v1/members/me/boxing?organization_id=${orgId}`);
  }

  async updateMyBoxingProfile(
    orgId: string,
    data: BoxingProfileInput
  ): Promise<BoxingProfileView> {
    return this.request(`/api/v1/members/me/boxing?organization_id=${orgId}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async createMyBout(orgId: string, data: BoutInput): Promise<Bout> {
    return this.request(`/api/v1/members/me/bouts?organization_id=${orgId}`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async deleteMyBout(orgId: string, boutId: string): Promise<{ message: string }> {
    return this.request(
      `/api/v1/members/me/bouts/${boutId}?organization_id=${orgId}`,
      { method: "DELETE" }
    );
  }
}

export class ApiRequestError extends Error {
  status: number;
  code?: string;

  constructor(message: string, status: number, code?: string) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
    this.code = code;
  }
}

export const api = new ApiClient(API_BASE);
