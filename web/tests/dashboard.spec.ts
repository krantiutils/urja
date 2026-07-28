import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

function todayISO(): string {
  return new Date().toISOString().split("T")[0];
}

/**
 * The dashboard derives every stat card from live API data:
 *   total members    <- members.total
 *   active members   <- members whose status is "active"
 *   today's activity <- attendance records whose check_in_at is today
 *   monthly revenue  <- accounts summary total_income
 *
 * So the mocks below are chosen to make each card a distinct number, and the
 * assertions read those numbers back rather than hardcoding figures that no
 * longer come from anywhere.
 */
function member(id: string, name: string, status: string, joined: string) {
  return {
    id,
    name,
    phone: `98410000${id.slice(-2)}`,
    name_ne: null,
    email: null,
    avatar_url: null,
    role: "member",
    status,
    joined_at: joined,
  };
}

const MOCK_MEMBERS = {
  // 4 returned, 3 of them active, but a total of 342 across all pages.
  data: [
    member("m-01", "Ram Shrestha", "active", "2026-01-15T10:30:00Z"),
    member("m-02", "Sita Maharjan", "active", "2026-01-20T08:00:00Z"),
    member("m-03", "Hari Tamang", "active", "2025-12-01T09:00:00Z"),
    member("m-04", "Deepa Thapa", "inactive", "2025-11-10T09:00:00Z"),
  ],
  total: 342,
};

function mockDashboardAPI(page: import("@playwright/test").Page) {
  const today = todayISO();

  const routes: Array<[string, unknown]> = [
    ["**/api/v1/orgs/*/members*", MOCK_MEMBERS],
    [
      "**/api/v1/orgs/*/attendance*",
      {
        data: [
          { id: "a1", user_id: "m-01", check_in_at: `${today}T06:15:00Z`, method: "qr" },
          { id: "a2", user_id: "m-02", check_in_at: `${today}T06:42:00Z`, method: "nfc" },
          // Yesterday — must not count toward today's activity.
          { id: "a3", user_id: "m-03", check_in_at: "2026-02-13T08:00:00Z", method: "qr" },
        ],
      },
    ],
    ["**/api/v1/orgs/*/packages/expiring*", { data: [] }],
    ["**/api/v1/orgs/*/packages/expired*", { data: [] }],
    ["**/api/v1/orgs/*/packages/summary*", []],
    ["**/api/v1/orgs/*/accounts/summary*", { total_income: "485000", total_expense: "0" }],
  ];

  return Promise.all(
    routes.map(([pattern, body]) =>
      page.route(pattern as string, (route) =>
        route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(body),
        })
      )
    )
  );
}

test.describe("Dashboard", () => {
  test.beforeEach(async ({ page }) => {
    // Need to navigate to the origin first so we can set localStorage
    await page.goto("/en/login");
    await injectAuth(page);
    await mockDashboardAPI(page);
  });

  test("shows stat cards and data tables", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Stat cards, each derived from a different part of the mock.
    await expect(page.getByText("342")).toBeVisible(); // total members
    await expect(page.getByText("3", { exact: true }).first()).toBeVisible(); // active
    await expect(page.getByText("Rs. 4,85,000")).toBeVisible(); // monthly revenue

    // Section headers, matched by role: the headings carry an icon alongside
    // the text, and the sidebar links use some of the same words.
    await expect(page.getByRole("heading", { name: "Today's Activity" })).toBeVisible();
    await expect(page.getByRole("heading", { name: "Expiring Packages" })).toBeVisible();

    // Member names from mock data
    await expect(page.getByText("Ram Shrestha").first()).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-full-en.png", fullPage: true });
  });

  test("sidebar is visible on desktop", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Sidebar should be present (fixed left on desktop, lg:translate-x-0)
    const sidebar = page.locator("aside");
    await expect(sidebar).toBeVisible();

    // Navigation sections should be visible
    await expect(page.getByText("Manage Members")).toBeVisible();
    await expect(page.getByText("Run Operations")).toBeVisible();

    // Dashboard nav item should be active
    const dashLink = sidebar.getByRole("link", { name: "Dashboard" });
    await expect(dashLink).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-sidebar-en.png", fullPage: true });
  });

  test("responsive mobile layout", async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto("/en/dashboard");

    // On mobile the sidebar is hidden (translated off-screen)
    const sidebar = page.locator("aside");
    // The sidebar exists in DOM but is translated left
    await expect(sidebar).toHaveCSS("transform", /matrix/);

    // Stat cards should still be visible
    await expect(page.getByText("342")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-mobile-en.png", fullPage: true });
  });

  test("mobile sidebar opens via menu button", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto("/en/dashboard");

    // Click the hamburger/menu button in the top bar
    const menuButton = page.locator("button").filter({ has: page.locator("svg.lucide-menu") });
    await menuButton.click();

    // Sidebar should now be visible (translate-x-0)
    const sidebar = page.locator("aside");
    await expect(sidebar).toBeVisible();

    // Overlay backdrop should be visible
    await expect(page.locator(".fixed.inset-0.bg-black\\/60")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-mobile-sidebar-en.png", fullPage: true });
  });

  test("renders dashboard in Nepali", async ({ page }) => {
    // Inject auth again since we need it for the Nepali route too
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockDashboardAPI(page);
    await page.goto("/ne/dashboard");

    // Nepali text
    await expect(page.getByText("कुल सदस्यहरू")).toBeVisible();
    await expect(page.getByText("आजको गतिविधि")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-ne.png", fullPage: true });
  });
});
