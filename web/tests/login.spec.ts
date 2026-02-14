import { test, expect } from "@playwright/test";

test.describe("Login page", () => {
  test("renders phone input form in English", async ({ page }) => {
    await page.goto("/en/login");

    // The phone input should be visible
    const phoneInput = page.locator("#phone");
    await expect(phoneInput).toBeVisible();
    await expect(phoneInput).toHaveAttribute("placeholder", "98XXXXXXXX");

    // Brand name should be visible
    await expect(page.getByText("Urja")).toBeVisible();

    // "Send OTP" button
    await expect(page.getByRole("button", { name: /Send OTP/i })).toBeVisible();

    await page.screenshot({ path: "screenshots/login-en.png", fullPage: true });
  });

  test("shows OTP step after submitting phone", async ({ page }) => {
    // Mock the login API to succeed
    await page.route("**/api/v1/auth/login", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ message: "OTP sent" }),
      }),
    );

    await page.goto("/en/login");

    // Fill a valid Nepali phone number and submit
    await page.locator("#phone").fill("9841000000");
    await page.getByRole("button", { name: /Send OTP/i }).click();

    // Wait for OTP form to appear
    const otpInput = page.locator("#otp");
    await expect(otpInput).toBeVisible({ timeout: 5000 });
    await expect(page.getByText("OTP sent to your phone")).toBeVisible();

    // Resend and back buttons should be present
    await expect(page.getByRole("button", { name: /Verify/i })).toBeVisible();

    await page.screenshot({ path: "screenshots/login-otp-en.png", fullPage: true });
  });

  test("shows validation error for invalid phone", async ({ page }) => {
    await page.goto("/en/login");

    // Enter an invalid number
    await page.locator("#phone").fill("12345");
    await page.getByRole("button", { name: /Send OTP/i }).click();

    // Should show the validation error
    await expect(page.getByText("Enter a valid Nepali mobile number")).toBeVisible();

    await page.screenshot({ path: "screenshots/login-invalid-phone-en.png", fullPage: true });
  });

  test("renders login page in Nepali", async ({ page }) => {
    await page.goto("/ne/login");

    // Nepali brand
    await expect(page.getByText("ऊर्जा")).toBeVisible();

    // Nepali sign in text
    await expect(page.getByText("साइन इन")).toBeVisible();

    // Phone input still present
    await expect(page.locator("#phone")).toBeVisible();

    await page.screenshot({ path: "screenshots/login-ne.png", fullPage: true });
  });
});
