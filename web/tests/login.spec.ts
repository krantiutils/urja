import { test, expect } from "@playwright/test";

test.describe("Login page", () => {
  test("renders phone input form in English", async ({ page }) => {
    await page.goto("/en/login");

    // Brand heading
    await expect(page.getByRole("heading", { name: "Urja" })).toBeVisible();

    // Subtitle
    await expect(page.getByText("Gym management, simplified")).toBeVisible();

    // "Sign In" card heading
    await expect(page.getByRole("heading", { name: "Sign In" })).toBeVisible();

    // Phone input visible and focused
    const phoneInput = page.locator("#phone");
    await expect(phoneInput).toBeVisible();
    await expect(phoneInput).toHaveAttribute("placeholder", "98XXXXXXXX");

    // Submit button
    await expect(page.getByRole("button", { name: /Send OTP/i })).toBeVisible();

    await page.screenshot({ path: "screenshots/login-en-phone.png", fullPage: true });
  });

  test("shows validation error for invalid phone number", async ({ page }) => {
    await page.goto("/en/login");

    // Type an invalid number and submit
    await page.locator("#phone").fill("1234");
    await page.getByRole("button", { name: /Send OTP/i }).click();

    // Validation error should appear
    await expect(page.getByText("Enter a valid Nepali mobile number")).toBeVisible();

    await page.screenshot({ path: "screenshots/login-en-phone-error.png", fullPage: true });
  });

  test("transitions to OTP step after submitting valid phone", async ({ page }) => {
    // Mock the login API endpoint
    await page.route("**/api/v1/auth/login", (route) =>
      route.fulfill({ status: 200, contentType: "application/json", body: "{}" })
    );

    await page.goto("/en/login");

    await page.locator("#phone").fill("9841234567");
    await page.getByRole("button", { name: /Send OTP/i }).click();

    // Wait for OTP form to appear
    await expect(page.getByText("OTP sent to your phone")).toBeVisible();
    await expect(page.locator("#otp")).toBeVisible();
    await expect(page.getByRole("button", { name: /Verify/i })).toBeVisible();

    // Back button and resend timer should be visible
    await expect(page.getByText(/Mobile Number/i)).toBeVisible();
    await expect(page.getByText(/Resend in/i)).toBeVisible();

    await page.screenshot({ path: "screenshots/login-en-otp.png", fullPage: true });
  });

  test("renders login page in Nepali locale", async ({ page }) => {
    await page.goto("/ne/login");

    // Nepali brand name
    await expect(page.getByRole("heading", { name: "ऊर्जा" })).toBeVisible();

    // Nepali subtitle
    await expect(page.getByText("जिम व्यवस्थापन, सरलीकृत")).toBeVisible();

    // Nepali sign-in heading
    await expect(page.getByRole("heading", { name: "साइन इन" })).toBeVisible();

    // Phone input
    await expect(page.locator("#phone")).toBeVisible();

    // Nepali OTP button text
    await expect(
      page.getByRole("button", { name: /OTP पठाउनुहोस्/i })
    ).toBeVisible();

    await page.screenshot({ path: "screenshots/login-ne.png", fullPage: true });
  });
});
