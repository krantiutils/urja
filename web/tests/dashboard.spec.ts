import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

test.describe("Dashboard", () => {
  test.beforeEach(async ({ page }) => {
    // Need to navigate to the origin first so we can set localStorage
    await page.goto("/en/login");
    await injectAuth(page);
  });

  test("shows stat cards and data tables", async ({ page }) => {
    await page.goto("/en/dashboard");

    // Stat cards should be visible with mock data
    await expect(page.getByText("342")).toBeVisible();
    await expect(page.getByText("287")).toBeVisible();
    await expect(page.getByText("89", { exact: true })).toBeVisible();
    await expect(page.getByText("Rs. 4,85,000")).toBeVisible();

    // Section headers
    await expect(page.getByText("Today's Activity")).toBeVisible();
    await expect(page.getByText("Expiring Packages")).toBeVisible();

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
    await page.goto("/ne/dashboard");

    // Nepali text
    await expect(page.getByText("कुल सदस्यहरू")).toBeVisible();
    await expect(page.getByText("आजको गतिविधि")).toBeVisible();

    await page.screenshot({ path: "screenshots/dashboard-ne.png", fullPage: true });
  });
});
