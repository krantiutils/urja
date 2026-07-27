import { test, expect, type Page, type Route } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_DUES = {
  data: [
    {
      id: "due-001",
      organization_id: "org-001",
      user_id: "m-001",
      member_name: "Ram Shrestha",
      member_phone: "9841000001",
      amount: "1500.00",
      due_date: "2026-07-01T00:00:00Z",
      description: "Monthly membership",
      status: "unpaid",
      paid_at: null,
      paid_amount: null,
      payment_method: null,
      payment_reference: null,
      created_at: "2026-06-01T00:00:00Z",
    },
  ],
  total: 1,
};

/**
 * Mocks the dues API. The `/pay` POST handler mirrors the real Go
 * handler's bug (`amount` decodes into a bare float64): a string amount
 * 400s just like the live backend did before the fix.
 */
function mockDuesAPI(page: Page, onPay?: (body: Record<string, unknown>) => void) {
  return page.route("**/api/v1/orgs/*/dues**", (route: Route) => {
    const url = new URL(route.request().url());
    const method = route.request().method();

    if (method === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_DUES),
      });
    }

    if (method === "POST" && url.pathname.endsWith("/pay")) {
      const body = route.request().postDataJSON();
      onPay?.(body);
      if (typeof body.amount !== "number") {
        return route.fulfill({
          status: 400,
          contentType: "application/json",
          body: JSON.stringify({ error: "json: cannot unmarshal string into Go struct field .amount of type float64" }),
        });
      }
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ message: "Payment recorded" }),
      });
    }

    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Due Payments Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
  });

  test("recording a payment sends amount as a number and succeeds", async ({ page }) => {
    const captured: { body: Record<string, unknown> | null } = { body: null };
    await mockDuesAPI(page, (body) => {
      captured.body = body;
    });

    await page.goto("/en/dashboard/due-payments");
    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    await page.getByRole("button", { name: "Record Payment" }).click();
    await page.locator("#due-pay-amount").fill("");
    await page.locator("#due-pay-amount").fill("1500.50");
    await page.getByRole("button", { name: "Confirm Payment" }).click();

    // Modal closes on success (no payError shown, dialog unmounts).
    await expect(page.getByText("Confirm Payment")).not.toBeVisible();

    expect(captured.body).not.toBeNull();
    expect(typeof captured.body?.amount).toBe("number");
    expect(captured.body?.amount).toBe(1500.5);

    await page.screenshot({
      path: "screenshots/due-payments-pay-success-en.png",
      fullPage: true,
    });
  });

  test("payment amount field rejects empty, non-numeric, zero and negative input", async ({ page }) => {
    await mockDuesAPI(page);
    await page.goto("/en/dashboard/due-payments");
    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    let payPostCount = 0;
    page.on("request", (req) => {
      if (req.method() === "POST" && req.url().includes("/pay")) payPostCount++;
    });

    await page.getByRole("button", { name: "Record Payment" }).click();

    for (const bad of ["", "abc", "0", "-20"]) {
      await page.locator("#due-pay-amount").fill(bad);
      await page.getByRole("button", { name: "Confirm Payment" }).click();
      await expect(page.getByText("Enter a valid amount greater than 0")).toBeVisible();
    }

    expect(payPostCount).toBe(0);

    await page.screenshot({
      path: "screenshots/due-payments-invalid-amount-en.png",
      fullPage: true,
    });
  });
});
