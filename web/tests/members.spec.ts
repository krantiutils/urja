import { test, expect } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_MEMBERS = {
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
      role: "staff",
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
      status: "suspended",
      joined_at: "2025-12-01T09:00:00Z",
    },
  ],
  total: 3,
};

function mockMembersAPI(page: import("@playwright/test").Page) {
  return page.route("**/api/v1/orgs/*/members*", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_MEMBERS),
      });
    }
    if (route.request().method() === "POST") {
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(MOCK_MEMBERS.data[0]),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Members Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockMembersAPI(page);
  });

  test("shows members table with data", async ({ page }) => {
    await page.goto("/en/dashboard/members");

    // Title visible
    await expect(page.getByRole("heading", { name: "Members" })).toBeVisible();

    // Member names
    await expect(page.getByText("Ram Shrestha")).toBeVisible();
    await expect(page.getByText("Sita Maharjan")).toBeVisible();
    await expect(page.getByText("Hari Tamang")).toBeVisible();

    // Phone numbers
    await expect(page.getByText("9841000001")).toBeVisible();

    // Status badges
    await expect(page.getByText("active").first()).toBeVisible();
    await expect(page.getByText("suspended")).toBeVisible();

    // Add Member button
    await expect(page.getByText("Add Member")).toBeVisible();

    await page.screenshot({
      path: "screenshots/members-en.png",
      fullPage: true,
    });
  });

  test("search filters members", async ({ page }) => {
    await page.goto("/en/dashboard/members");

    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    // Type in search
    await page.getByPlaceholder("Search by name or phone...").fill("Sita");

    // Ram should be hidden, Sita visible
    await expect(page.getByText("Ram Shrestha")).not.toBeVisible();
    await expect(page.getByText("Sita Maharjan")).toBeVisible();

    await page.screenshot({
      path: "screenshots/members-search-en.png",
      fullPage: true,
    });
  });

  test("add member modal opens and closes", async ({ page }) => {
    await page.goto("/en/dashboard/members");

    // Open modal
    await page.getByText("Add Member").click();
    await expect(page.getByText("Add New Member")).toBeVisible();

    // Form fields visible
    await expect(page.getByPlaceholder("Full name")).toBeVisible();
    await expect(page.getByPlaceholder("98XXXXXXXX")).toBeVisible();

    await page.screenshot({
      path: "screenshots/members-add-modal-en.png",
      fullPage: true,
    });

    // Close modal
    await page.getByText("Cancel").click();
    await expect(page.getByText("Add New Member")).not.toBeVisible();
  });

  test("renders members page in Nepali", async ({ page }) => {
    await page.goto("/ne/login");
    await injectAuth(page);
    await mockMembersAPI(page);
    await page.goto("/ne/dashboard/members");

    await expect(page.getByRole("heading", { name: "सदस्यहरू" })).toBeVisible();
    await expect(page.getByRole("button", { name: "सदस्य थप्नुहोस्" })).toBeVisible();

    await page.screenshot({
      path: "screenshots/members-ne.png",
      fullPage: true,
    });
  });
});
