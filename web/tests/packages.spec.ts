import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_PACKAGES = {
  data: [
    {
      id: "pkg-001",
      organization_id: "org-001",
      name: "1 Month Basic",
      name_ne: "१ महिना बेसिक",
      description: "Basic gym access",
      description_ne: null,
      duration_days: 30,
      price: "2000.00",
      currency: "NPR",
      max_members: null,
      features: ["gym_access"],
      is_active: true,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    },
    {
      id: "pkg-002",
      organization_id: "org-001",
      name: "3 Month Premium",
      name_ne: null,
      description: "Full gym + trainer sessions",
      description_ne: null,
      duration_days: 90,
      price: "5000.00",
      currency: "NPR",
      max_members: 50,
      features: ["gym_access", "trainer_sessions", "pool_access"],
      is_active: true,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    },
    {
      id: "pkg-003",
      organization_id: "org-001",
      name: "Expired Offer",
      name_ne: null,
      description: "No longer available",
      description_ne: null,
      duration_days: 30,
      price: "1000.00",
      currency: "NPR",
      max_members: null,
      features: [],
      is_active: false,
      created_at: "2025-06-01T00:00:00Z",
      updated_at: "2025-12-01T00:00:00Z",
    },
  ],
};

function mockPackagesAPI(page: import("@playwright/test").Page) {
  return page.route("**/api/v1/orgs/*/packages*", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_PACKAGES),
      });
    }
    if (route.request().method() === "POST") {
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(MOCK_PACKAGES.data[0]),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Packages Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockPackagesAPI(page);
  });

  test("shows package cards with data", async ({ page }) => {
    await page.goto("/en/dashboard/packages");

    // Title
    await expect(page.getByRole("heading", { name: "Packages" })).toBeVisible();

    // Package names
    await expect(page.getByText("1 Month Basic")).toBeVisible();
    await expect(page.getByText("3 Month Premium")).toBeVisible();
    await expect(page.getByText("Expired Offer")).toBeVisible();

    // Prices
    await expect(page.getByText("Rs. 2,000")).toBeVisible();
    await expect(page.getByText("Rs. 5,000")).toBeVisible();

    // Durations
    await expect(page.getByText("30 days").first()).toBeVisible();
    await expect(page.getByText("90 days")).toBeVisible();

    // Features
    await expect(page.getByText("gym_access").first()).toBeVisible();
    await expect(page.getByText("trainer_sessions")).toBeVisible();

    // Status badges
    const activeBadges = page.getByText("Active", { exact: true });
    expect(await activeBadges.count()).toBeGreaterThanOrEqual(1);

    await page.screenshot({
      path: "screenshots/packages-en.png",
      fullPage: true,
    });
  });

  test("filter toggles work", async ({ page }) => {
    await page.goto("/en/dashboard/packages");

    // All packages visible by default
    await expect(page.getByText("1 Month Basic")).toBeVisible();
    await expect(page.getByText("Expired Offer")).toBeVisible();

    // Click Active filter
    await page.locator("button").filter({ hasText: /^Active$/i }).click();
    await expect(page.getByText("1 Month Basic")).toBeVisible();
    await expect(page.getByText("Expired Offer")).not.toBeVisible();

    // Click Inactive filter
    await page.locator("button").filter({ hasText: /^Inactive$/i }).click();
    await expect(page.getByText("Expired Offer")).toBeVisible();
    await expect(page.getByText("1 Month Basic")).not.toBeVisible();

    await page.screenshot({
      path: "screenshots/packages-filtered-en.png",
      fullPage: true,
    });
  });

  test("add package modal opens", async ({ page }) => {
    await page.goto("/en/dashboard/packages");

    await page.getByText("Add Package").click();
    await expect(page.getByText("Add New Package")).toBeVisible();

    // Form fields
    await expect(page.getByPlaceholder("Package name")).toBeVisible();
    await expect(page.getByPlaceholder("0.00")).toBeVisible();

    await page.screenshot({
      path: "screenshots/packages-add-modal-en.png",
      fullPage: true,
    });
  });

  test("renders packages page in Nepali", async ({ page }) => {
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockPackagesAPI(page);
    await page.goto("/ne/dashboard/packages");

    await expect(page.getByRole("heading", { name: "प्याकेजहरू" })).toBeVisible();
    await expect(page.getByRole("button", { name: "प्याकेज थप्नुहोस्" })).toBeVisible();

    await page.screenshot({
      path: "screenshots/packages-ne.png",
      fullPage: true,
    });
  });
});
