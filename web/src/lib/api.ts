import type {
  AuthTokens,
  ApiError,
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
} from "@/types";

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
      throw new ApiRequestError(body.error, response.status);
    }

    return response.json() as Promise<T>;
  }

  async login(phone: string): Promise<{ message: string }> {
    return this.request("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ phone }),
    });
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
}

export class ApiRequestError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
  }
}

export const api = new ApiClient(API_BASE);
