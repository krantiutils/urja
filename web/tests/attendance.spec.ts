import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

function todayISO(): string {
  return new Date().toISOString().split("T")[0];
}

function makeMockAttendance() {
  const today = todayISO();
  return {
    data: [
      {
        id: "att-001",
        user_id: "m-001",
        org_id: "org-001",
        check_in_at: `${today}T06:15:00Z`,
        method: "qr",
        member_name: "Ram Shrestha",
      },
      {
        id: "att-002",
        user_id: "m-002",
        org_id: "org-001",
        check_in_at: `${today}T06:42:00Z`,
        method: "nfc",
        member_name: "Sita Maharjan",
      },
      {
        id: "att-003",
        user_id: "m-003",
        org_id: "org-001",
        check_in_at: `${today}T07:01:00Z`,
        method: "manual",
        member_name: "Hari Tamang",
      },
      {
        id: "att-004",
        user_id: "m-004",
        org_id: "org-001",
        check_in_at: "2026-02-13T08:00:00Z",
        method: "qr",
        member_name: "Deepa Thapa",
      },
    ],
  };
}

function mockAttendanceAPI(page: import("@playwright/test").Page) {
  const mock = makeMockAttendance();
  return page.route("**/api/v1/orgs/*/attendance*", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(mock),
      });
    }
    if (route.request().method() === "POST") {
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify({
          check_in: mock.data[0],
          streak: {
            id: "s-001",
            member_id: "m-001",
            org_id: "org-001",
            current_streak: 5,
            longest_streak: 10,
            last_check_in: todayISO(),
            updated_at: new Date().toISOString(),
          },
        }),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

// AttendancePage.fetchData() runs listAttendance and listMembers in the same
// Promise.all — the members endpoint must also be mocked or the whole fetch
// rejects (no backend at localhost:8080) and the attendance mock never
// renders. Member ids here must line up with the attendance mock's user_id
// values since the table resolves names via memberNameMap, not
// record.member_name.
function makeMockMembers() {
  return {
    data: [
      {
        id: "m-001",
        phone: "9841000001",
        name: "Ram Shrestha",
        name_ne: "राम श्रेष्ठ",
        email: "ram@example.com",
        avatar_url: null,
        role: "member",
        status: "active",
        joined_at: "2026-01-15T10:30:00Z",
      },
      {
        id: "m-002",
        phone: "9841000002",
        name: "Sita Maharjan",
        name_ne: null,
        email: null,
        avatar_url: null,
        role: "member",
        status: "active",
        joined_at: "2026-01-20T08:00:00Z",
      },
      {
        id: "m-003",
        phone: "9841000003",
        name: "Hari Tamang",
        name_ne: null,
        email: "hari@example.com",
        avatar_url: null,
        role: "member",
        status: "active",
        joined_at: "2025-12-01T09:00:00Z",
      },
      {
        id: "m-004",
        phone: "9841000004",
        name: "Deepa Thapa",
        name_ne: null,
        email: null,
        avatar_url: null,
        role: "member",
        status: "active",
        joined_at: "2025-11-10T09:00:00Z",
      },
    ],
    total: 4,
  };
}

function mockMembersAPI(page: import("@playwright/test").Page) {
  const mock = makeMockMembers();
  return page.route("**/api/v1/orgs/*/members*", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(mock),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Attendance Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockAttendanceAPI(page);
    await mockMembersAPI(page);
  });

  test("shows today attendance records", async ({ page }) => {
    await page.goto("/en/dashboard/attendance");

    // Title
    await expect(
      page.getByRole("heading", { name: "Attendance" })
    ).toBeVisible();

    // Today's check-ins stat card shows 3 (today's records)
    await expect(page.getByText("Today's Check-ins")).toBeVisible();

    // Member names from today
    await expect(page.getByText("Ram Shrestha")).toBeVisible();
    await expect(page.getByText("Sita Maharjan")).toBeVisible();
    await expect(page.getByText("Hari Tamang")).toBeVisible();

    // Method badges (use locator within table to avoid sidebar NFC link)
    const table = page.locator("table");
    await expect(table.getByText("qr").first()).toBeVisible();
    await expect(table.getByText("nfc")).toBeVisible();
    await expect(table.getByText("manual")).toBeVisible();

    // Manual Check-in button
    await expect(page.getByText("Manual Check-in")).toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-en.png",
      fullPage: true,
    });
  });

  test("method filter works", async ({ page }) => {
    await page.goto("/en/dashboard/attendance");

    // All records visible
    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    // Click QR filter
    await page
      .locator("button")
      .filter({ hasText: /^QR$/i })
      .click();

    // Only QR records visible (Ram is QR)
    await expect(page.getByText("Ram Shrestha")).toBeVisible();
    // NFC records hidden (Sita)
    await expect(page.getByText("Sita Maharjan")).not.toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-filtered-en.png",
      fullPage: true,
    });
  });

  test("manual check-in modal opens", async ({ page }) => {
    await page.goto("/en/dashboard/attendance");

    await page.getByText("Manual Check-in").click();
    await expect(page.getByText("Check in a member")).toBeVisible();

    // The modal takes a searchable member picker, not a raw id: staff know a
    // member by name, not by UUID.
    const search = page.getByPlaceholder("Search by name or phone...");
    await expect(search).toBeVisible();

    await search.fill("Sita");
    const option = page.getByRole("button", { name: /Sita Maharjan/ });
    await expect(option).toBeVisible();

    await option.click();
    // Picking a member replaces the search box with the chosen name. Scoped to
    // the form because the same name is also in the table behind the modal.
    await expect(search).toBeHidden();
    await expect(page.locator("form").getByText("Sita Maharjan")).toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-checkin-modal-en.png",
      fullPage: true,
    });
  });

  test("renders attendance in Nepali", async ({ page }) => {
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockAttendanceAPI(page);
    await mockMembersAPI(page);
    await page.goto("/ne/dashboard/attendance");

    await expect(page.getByRole("heading", { name: "उपस्थिति" })).toBeVisible();
    await expect(page.getByRole("button", { name: "म्यानुअल चेक-इन" })).toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-ne.png",
      fullPage: true,
    });
  });
});
