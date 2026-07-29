import { test, expect, type Page, type Route } from "@playwright/test";
import { injectAuth } from "./helpers";

const ORG = {
  id: "org-001",
  name: "Test Gym",
  slug: "test-gym",
  pan_number: "",
  tax_legal_name: "",
  tax_address: "",
};

function mockOrgAPI(page: Page) {
  let org = { ...ORG };
  return page.route("**/api/v1/orgs/org-001", (route: Route) => {
    if (route.request().method() === "PUT") {
      org = { ...org, ...route.request().postDataJSON() };
    }
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(org),
    });
  });
}

test.describe("Tax settings", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockOrgAPI(page);
  });

  test("rejects a PAN that is not nine digits", async ({ page }) => {
    await page.goto("/en/dashboard/settings");

    const pan = page.getByLabel(/PAN/i);
    await expect(pan).toBeVisible();

    await pan.fill("12345");
    // Scoped to the Tax form: the page also has "Save Changes"/"Save password"
    // buttons on other settings sections, so an unscoped role query is ambiguous.
    await page.locator("form").getByRole("button", { name: /save/i }).click();
    await expect(page.getByText(/9 digits/i)).toBeVisible();
  });

  test("accepts a nine-digit PAN", async ({ page }) => {
    await page.goto("/en/dashboard/settings");

    await page.getByLabel(/PAN/i).fill("601234567");
    await page.locator("form").getByRole("button", { name: /save/i }).click();
    await expect(page.getByText(/9 digits/i)).toHaveCount(0);
  });

  test("clearing the PAN is accepted, not rejected client-side", async ({ page }) => {
    await page.goto("/en/dashboard/settings");

    const pan = page.getByLabel(/PAN/i);
    const saveButton = page.locator("form").getByRole("button", { name: /save/i });

    // A gym that saved a wrong PAN removes it by blanking the box. The API
    // treats an empty string as "clear the field", not "leave unchanged", so
    // the client must let an empty submission through rather than complain.
    await pan.fill("601234567");
    await saveButton.click();
    await expect(page.getByText(/9 digits/i)).toHaveCount(0);

    await pan.fill("");
    await saveButton.click();
    await expect(page.getByText(/9 digits/i)).toHaveCount(0);
  });
});
