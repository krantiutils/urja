import { test, expect, type Page, type Route } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_SUMMARY = {
  total_income: "10000.00",
  total_expenses: "2000.00",
  gross_profit: "8000.00",
  profit_percent: 80,
};

function baseTransaction(overrides: Record<string, unknown> = {}) {
  return {
    id: "tx-001",
    organization_id: "org-001",
    category: "membership",
    description: "Existing seed transaction",
    transaction_date: "2026-07-01T00:00:00Z",
    transaction_type: "income",
    amount: "500.00",
    payment_type: "cash",
    reference: "",
    entry_by: "u-001",
    entry_by_name: "Test Admin",
    created_at: "2026-07-01T00:00:00Z",
    updated_at: "2026-07-01T00:00:00Z",
    ...overrides,
  };
}

/**
 * Mocks the accounts API. The POST handler mirrors the real Go handler's
 * bug (`amount` decodes into a bare float64, no `,string` tag): if the
 * request body's `amount` isn't a JSON number, it 400s exactly like the
 * live backend did before the fix. This makes the "persists" test fail
 * against the old string-amount code and pass against the fix.
 */
function mockAccountsAPI(page: Page) {
  const transactions = [baseTransaction()];
  return page.route("**/api/v1/orgs/*/accounts**", (route: Route) => {
    const url = new URL(route.request().url());
    const method = route.request().method();

    if (url.pathname.endsWith("/summary")) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_SUMMARY),
      });
    }

    if (method === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: transactions, total: transactions.length }),
      });
    }

    if (method === "POST") {
      const body = route.request().postDataJSON();
      if (typeof body.amount !== "number") {
        return route.fulfill({
          status: 400,
          contentType: "application/json",
          body: JSON.stringify({ error: "json: cannot unmarshal string into Go struct field .amount of type float64" }),
        });
      }
      const created = baseTransaction({
        id: `tx-${transactions.length + 1}`,
        category: body.category,
        description: body.description,
        transaction_date: body.transaction_date,
        transaction_type: body.transaction_type,
        amount: String(body.amount),
        payment_type: body.payment_type,
        reference: body.reference ?? "",
      });
      transactions.push(created);
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(created),
      });
    }

    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });
}

test.describe("Accounts Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockAccountsAPI(page);
  });

  test("adding a transaction with amount 1500.50 persists and appears in the list", async ({ page }) => {
    await page.goto("/en/dashboard/accounts");

    await expect(page.getByText("Existing seed transaction")).toBeVisible();

    await page.getByRole("button", { name: "Add Transaction" }).click();
    await page.locator("#tx-category").fill("membership");
    await page.locator("#tx-description").fill("Playwright test entry");
    await page.locator("#tx-amount").fill("1500.50");
    await page.getByRole("button", { name: "Save" }).click();

    // Modal closes (its "Add Transaction" heading unmounts) and the new
    // row shows up after the post-create refetch.
    await expect(page.getByRole("heading", { name: "Add Transaction" })).not.toBeVisible();
    await expect(page.getByText("Playwright test entry")).toBeVisible();

    await page.screenshot({
      path: "screenshots/accounts-add-success-en.png",
      fullPage: true,
    });
  });

  test("amount field rejects empty, non-numeric, zero and negative input", async ({ page }) => {
    await page.goto("/en/dashboard/accounts");

    let postCount = 0;
    page.on("request", (req) => {
      if (req.method() === "POST" && /\/accounts$/.test(new URL(req.url()).pathname)) {
        postCount++;
      }
    });

    await page.getByRole("button", { name: "Add Transaction" }).click();
    await page.locator("#tx-category").fill("membership");
    await page.locator("#tx-description").fill("Bad amount test");

    for (const bad of ["", "abc", "0", "-50"]) {
      await page.locator("#tx-amount").fill(bad);
      await page.getByRole("button", { name: "Save" }).click();
      await expect(page.getByText("Enter a valid amount greater than 0")).toBeVisible();
    }

    expect(postCount).toBe(0);

    await page.screenshot({
      path: "screenshots/accounts-invalid-amount-en.png",
      fullPage: true,
    });
  });
});
