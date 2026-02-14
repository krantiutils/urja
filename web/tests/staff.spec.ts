import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_STAFF = {
  data: [
    {
      id: "s-001",
      user_id: "u-001",
      name: "Bikash Gurung",
      phone: "9841100001",
      email: "bikash@fitgym.com",
      avatar_url: null,
      role: "admin",
      staff_role: "owner",
      status: "active",
      joined_at: "2025-06-01T00:00:00Z",
    },
    {
      id: "s-002",
      user_id: "u-002",
      name: "Anita Rai",
      phone: "9841100002",
      email: "anita@fitgym.com",
      avatar_url: null,
      role: "staff",
      staff_role: "trainer",
      status: "active",
      joined_at: "2025-08-15T00:00:00Z",
    },
    {
      id: "s-003",
      user_id: "u-003",
      name: "Deepak Thapa",
      phone: "9841100003",
      email: null,
      avatar_url: null,
      role: "staff",
      staff_role: "receptionist",
      status: "suspended",
      joined_at: "2025-10-01T00:00:00Z",
    },
  ],
  total: 3,
};

function mockStaffAPI(page: import("@playwright/test").Page) {
  return page.route("**/api/v1/orgs/*/staff*", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_STAFF),
      });
    }
    if (route.request().method() === "POST") {
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(MOCK_STAFF.data[1]),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Staff Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockStaffAPI(page);
  });

  test("shows staff table with data", async ({ page }) => {
    await page.goto("/en/dashboard/staff");

    // Title
    await expect(page.getByRole("heading", { name: "Staff" })).toBeVisible();

    // Staff names
    await expect(page.getByText("Bikash Gurung")).toBeVisible();
    await expect(page.getByText("Anita Rai")).toBeVisible();
    await expect(page.getByText("Deepak Thapa")).toBeVisible();

    // Phone
    await expect(page.getByText("9841100001")).toBeVisible();

    // Staff roles
    await expect(page.getByText("owner")).toBeVisible();
    await expect(page.getByText("trainer")).toBeVisible();
    await expect(page.getByText("receptionist")).toBeVisible();

    // Status
    await expect(page.getByText("active").first()).toBeVisible();
    await expect(page.getByText("suspended")).toBeVisible();

    // Count footer
    await expect(page.getByText("3 staff")).toBeVisible();

    await page.screenshot({
      path: "screenshots/staff-en.png",
      fullPage: true,
    });
  });

  test("add staff modal opens with role selector", async ({ page }) => {
    await page.goto("/en/dashboard/staff");

    await page.getByText("Add Staff").click();
    await expect(page.getByText("Add New Staff")).toBeVisible();

    // Form fields
    await expect(page.getByPlaceholder("Full name")).toBeVisible();
    await expect(page.getByPlaceholder("98XXXXXXXX")).toBeVisible();

    // Staff role selector
    const roleSelect = page.locator("select");
    await expect(roleSelect).toBeVisible();

    await page.screenshot({
      path: "screenshots/staff-add-modal-en.png",
      fullPage: true,
    });
  });

  test("search filters staff", async ({ page }) => {
    await page.goto("/en/dashboard/staff");

    await expect(page.getByText("Bikash Gurung")).toBeVisible();

    // The search triggers a new API call (server-side search), so re-mock with filtered data
    await page.route("**/api/v1/orgs/*/staff*", (route) => {
      const url = new URL(route.request().url());
      const searchParam = url.searchParams.get("search");
      if (searchParam && searchParam.includes("Anita")) {
        return route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            data: [MOCK_STAFF.data[1]],
            total: 1,
          }),
        });
      }
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_STAFF),
      });
    });

    await page.getByPlaceholder("Search by name or phone...").fill("Anita");

    // Wait for debounced refetch
    await page.waitForTimeout(500);
    await expect(page.getByText("Anita Rai")).toBeVisible();

    await page.screenshot({
      path: "screenshots/staff-search-en.png",
      fullPage: true,
    });
  });

  test("renders staff page in Nepali", async ({ page }) => {
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockStaffAPI(page);
    await page.goto("/ne/dashboard/staff");

    await expect(page.getByRole("heading", { name: "कर्मचारी" })).toBeVisible();
    await expect(page.getByRole("button", { name: "कर्मचारी थप्नुहोस्" })).toBeVisible();

    await page.screenshot({
      path: "screenshots/staff-ne.png",
      fullPage: true,
    });
  });
});
