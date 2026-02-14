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

test.describe("Attendance Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockAttendanceAPI(page);
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
    await expect(page.getByPlaceholder("Member ID")).toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-checkin-modal-en.png",
      fullPage: true,
    });
  });

  test("renders attendance in Nepali", async ({ page }) => {
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockAttendanceAPI(page);
    await page.goto("/ne/dashboard/attendance");

    await expect(page.getByRole("heading", { name: "उपस्थिति" })).toBeVisible();
    await expect(page.getByRole("button", { name: "म्यानुअल चेक-इन" })).toBeVisible();

    await page.screenshot({
      path: "screenshots/attendance-ne.png",
      fullPage: true,
    });
  });
});
