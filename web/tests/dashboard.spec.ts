import { test, expect } from "@playwright/test";
import { fakeJWT, injectAuthScript } from "./helpers";

test.describe("Dashboard page", () => {
  const token = fakeJWT();

  test.beforeEach(async ({ page }) => {
    // Inject auth tokens before the app reads localStorage on mount
    await page.addInitScript(injectAuthScript(token));
  });

  test("renders stat cards with mock data", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Wait for dashboard content to load (stat cards)
    await expect(page.getByText("342")).toBeVisible();
    await expect(page.getByText("287")).toBeVisible();
    await expect(page.getByText("89", { exact: true })).toBeVisible();
    await expect(page.getByText("Rs. 4,85,000")).toBeVisible();

    // Stat labels
    await expect(page.getByText("Total Members")).toBeVisible();
    await expect(page.getByText("Active Members")).toBeVisible();
    await expect(page.getByText("Today's Attendance")).toBeVisible();
    await expect(page.getByText("Monthly Revenue")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-full.png", fullPage: true });
  });

  test("renders attendance table with mock records", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Attendance records — names appear in multiple tables, scope to first match
    await expect(page.getByText("Ram Shrestha").first()).toBeVisible();
    await expect(page.getByText("06:15 AM")).toBeVisible();
    await expect(page.getByText("Sita Maharjan").first()).toBeVisible();

    // Method badges
    const qrBadges = page.locator("span", { hasText: /^qr$/i });
    await expect(qrBadges.first()).toBeVisible();
  });

  test("renders expiring and expired package tables", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Expiring packages section
    await expect(page.getByText("Expiring Packages")).toBeVisible();
    await expect(page.getByText("3d left")).toBeVisible();
    await expect(page.getByText("1d left")).toBeVisible();

    // Expired packages section
    await expect(page.getByText("Expired Packages")).toBeVisible();
    await expect(page.getByText("Bikash Gurung")).toBeVisible();
  });

  test("renders package summary table with totals", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Package summary
    await expect(page.getByText("Package Summary")).toBeVisible();
    // Package names appear in multiple tables, so use .first()
    await expect(page.getByText("1 Month Basic").first()).toBeVisible();
    await expect(page.getByText("3 Month Premium").first()).toBeVisible();
    await expect(page.getByText("6 Month Standard").first()).toBeVisible();
    await expect(page.getByText("1 Year Gold").first()).toBeVisible();

    // Total row
    await expect(page.getByText("329")).toBeVisible(); // 98+124+65+42
  });

  test("sidebar is visible on desktop with navigation links", async ({ page }) => {
    await page.goto("/en/dashboard");

    const sidebar = page.locator("aside");
    await expect(sidebar).toBeVisible();

    // Sidebar header
    await expect(sidebar.getByText("Gym Name")).toBeVisible();

    // Navigation sections
    await expect(sidebar.getByText("Manage Members")).toBeVisible();
    await expect(sidebar.getByText("Run Operations")).toBeVisible();
    await expect(sidebar.getByText("Engage With Members")).toBeVisible();
    await expect(sidebar.getByText("Access Control")).toBeVisible();

    // Navigation items
    await expect(sidebar.getByRole("link", { name: "Dashboard" })).toBeVisible();
    await expect(sidebar.getByRole("link", { name: "Members" })).toBeVisible();
    await expect(sidebar.getByRole("link", { name: "Packages" })).toBeVisible();
    await expect(sidebar.getByRole("link", { name: "Settings" })).toBeVisible();

    // Dashboard link should be active (highlighted)
    const dashboardLink = sidebar.getByRole("link", { name: "Dashboard" });
    await expect(dashboardLink).toHaveClass(/text-accent/);

    await page.screenshot({ path: "screenshots/dashboard-sidebar.png", fullPage: true });
  });

  test("sidebar collapses section on click", async ({ page }) => {
    await page.goto("/en/dashboard");

    const sidebar = page.locator("aside");

    // "Manage Members" section items should be visible initially
    await expect(sidebar.getByRole("link", { name: "Members" })).toBeVisible();

    // Click the section header to collapse
    await sidebar.getByText("Manage Members").click();

    // Items should be hidden after collapse
    await expect(sidebar.getByRole("link", { name: "Members" })).toBeHidden();

    await page.screenshot({
      path: "screenshots/dashboard-sidebar-collapsed.png",
      fullPage: true,
    });
  });

  test("mobile viewport hides sidebar and shows menu button", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto("/en/dashboard");

    // Sidebar should be off-screen (translated)
    const sidebar = page.locator("aside");
    await expect(sidebar).toHaveClass(/-translate-x-full/);

    // Menu toggle button should be visible
    const menuButton = page.locator("button").filter({ has: page.locator("svg") }).first();
    await expect(menuButton).toBeVisible();

    // Stat cards should still render
    await expect(page.getByText("342")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-mobile.png", fullPage: true });
  });

  test("mobile viewport can open sidebar via menu button", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto("/en/dashboard");

    // Click the menu toggle (first button in top bar with an SVG icon)
    // TopBar has the menu toggle as the first button with lg:hidden
    const menuButton = page.locator("header button").first();
    await menuButton.click();

    // Sidebar should now be visible (translate-x-0)
    const sidebar = page.locator("aside");
    await expect(sidebar).toHaveClass(/translate-x-0/);

    // Overlay backdrop should appear
    await expect(page.locator(".fixed.inset-0.bg-black\\/60")).toBeVisible();

    await page.screenshot({
      path: "screenshots/dashboard-mobile-sidebar-open.png",
      fullPage: true,
    });
  });

  test("top bar shows language toggle and profile", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Language toggle — when on EN locale, button shows "NE" (target locale)
    const header = page.locator("header");
    await expect(header.getByText("NE", { exact: true })).toBeVisible();

    // QR button
    await expect(header.getByText("QR Code")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-topbar.png" });
  });
});
