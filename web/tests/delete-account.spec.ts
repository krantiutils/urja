import { test, expect, type Page } from "@playwright/test";

const ACCESS_TOKEN = "fake-access-token-abc123";
const REFRESH_TOKEN = "fake-refresh-token-xyz789";

/**
 * Mocks the account-deletion flow's three calls: request OTP, verify OTP
 * (real route is /auth/verify-otp, returning {access_token, ...} — not
 * {token} from a nonexistent /auth/verify), and the authenticated DELETE.
 * Returns the list of Authorization headers seen on the DELETE call so
 * tests can assert the token from verify-otp is what actually gets sent.
 */
async function mockDeleteFlowAPIs(page: Page): Promise<string[]> {
  await page.route("**/api/v1/auth/login", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ message: "OTP sent" }),
    })
  );

  await page.route("**/api/v1/auth/verify-otp", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        access_token: ACCESS_TOKEN,
        refresh_token: REFRESH_TOKEN,
        token_type: "Bearer",
        is_new_user: false,
        onboarding_completed: true,
      }),
    })
  );

  const deleteAuthHeaders: string[] = [];
  await page.route("**/api/v1/members/me", (route) => {
    deleteAuthHeaders.push(route.request().headers()["authorization"] ?? "");
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ message: "Account deleted" }),
    });
  });

  return deleteAuthHeaders;
}

test.describe("Delete Account Page", () => {
  test("account deletion completes end to end and reaches the success state", async ({ page }) => {
    const deleteAuthHeaders = await mockDeleteFlowAPIs(page);

    await page.goto("/en/delete");
    await expect(page.getByRole("heading", { name: "Delete Account" })).toBeVisible();

    await page.locator("#phone").fill("9841000000");
    await page.getByRole("button", { name: "Send Verification Code" }).click();

    await expect(page.locator("#otp")).toBeVisible();

    await page.locator("#otp").fill("123456");
    await page.getByRole("button", { name: "Delete My Account" }).click();

    await expect(
      page.getByText("Your account has been successfully deleted.")
    ).toBeVisible();

    // The DELETE must be authorized with the real access_token from
    // verify-otp, not an undefined token from a nonexistent {token} field.
    expect(deleteAuthHeaders).toHaveLength(1);
    expect(deleteAuthHeaders[0]).toBe(`Bearer ${ACCESS_TOKEN}`);

    await page.screenshot({
      path: "screenshots/delete-account-success-en.png",
      fullPage: true,
    });
  });

  test("never calls the nonexistent /auth/verify endpoint and never sends a Bearer undefined token", async ({ page }) => {
    const deleteAuthHeaders = await mockDeleteFlowAPIs(page);

    let hitOldVerifyEndpoint = false;
    await page.route("**/api/v1/auth/verify", (route) => {
      hitOldVerifyEndpoint = true;
      return route.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "not found" }),
      });
    });

    await page.goto("/en/delete");
    await page.locator("#phone").fill("9841000000");
    await page.getByRole("button", { name: "Send Verification Code" }).click();
    await expect(page.locator("#otp")).toBeVisible();

    await page.locator("#otp").fill("123456");
    await page.getByRole("button", { name: "Delete My Account" }).click();

    await expect(
      page.getByText("Your account has been successfully deleted.")
    ).toBeVisible();

    expect(hitOldVerifyEndpoint).toBe(false);
    expect(deleteAuthHeaders.every((h) => !h.includes("undefined"))).toBe(true);
  });

  test("shows an error and does not delete when OTP verification fails", async ({ page }) => {
    await page.route("**/api/v1/auth/login", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ message: "OTP sent" }),
      })
    );
    await page.route("**/api/v1/auth/verify-otp", (route) =>
      route.fulfill({
        status: 400,
        contentType: "application/json",
        body: JSON.stringify({ error: "Invalid or expired OTP" }),
      })
    );
    let deleteWasCalled = false;
    await page.route("**/api/v1/members/me", (route) => {
      deleteWasCalled = true;
      return route.fulfill({ status: 200, body: '{"message":"ok"}' });
    });

    await page.goto("/en/delete");
    await page.locator("#phone").fill("9841000000");
    await page.getByRole("button", { name: "Send Verification Code" }).click();
    await expect(page.locator("#otp")).toBeVisible();

    await page.locator("#otp").fill("000000");
    await page.getByRole("button", { name: "Delete My Account" }).click();

    await expect(page.getByText("Invalid or expired OTP")).toBeVisible();
    expect(deleteWasCalled).toBe(false);
  });
});
