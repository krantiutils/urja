import { test, expect, type Page, type Route } from "@playwright/test";
import { injectAdminAuth, injectAuth } from "./helpers";

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
    // The Tax section is gated the same way the sibling org cards are
    // (canEditOrg): an "admin" identity is required to see the editable form
    // rather than the read-only fallback.
    await injectAdminAuth(page);
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

  test("loads a previously saved PAN from the authenticated org endpoint", async ({ page }) => {
    // A saved PAN must actually show up on reload. This mocks only
    // /api/v1/orgs/{orgId} (authenticated, carries the tax identity) — not
    // the public /api/v1/gyms/{id} directory, which never does (see the
    // backend's orgPublicColumns). If api.getOrg regresses back to the
    // /gyms/ endpoint, this request goes unmocked, the fetch fails, and the
    // field renders blank instead of pre-filled — that's the failure this
    // test exists to catch.
    await page.route("**/api/v1/orgs/org-001", (route: Route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          id: "org-001",
          name: "Test Gym",
          slug: "test-gym",
          pan_number: "609999999",
          tax_legal_name: "Existing Legal Name Pvt Ltd",
          tax_address: "Existing Address, Kathmandu",
        }),
      })
    );

    await page.goto("/en/dashboard/settings");

    await expect(page.getByLabel(/PAN/i)).toHaveValue("609999999");
    await expect(page.getByLabel(/registered legal name/i)).toHaveValue(
      "Existing Legal Name Pvt Ltd"
    );
    await expect(page.getByLabel(/registered address/i)).toHaveValue(
      "Existing Address, Kathmandu"
    );
  });
});

function baseInvoice(overrides: Record<string, unknown> = {}) {
  return {
    id: "inv-001",
    organization_id: "org-001",
    fiscal_year: "2082-83",
    sequence: 1,
    invoice_number: "2082-83/000001",
    doc_type: "invoice",
    seller_name: "Test Gym",
    seller_pan: "601234567",
    seller_vat_registered: false,
    customer_name: "Ram Bahadur",
    issued_date: "2025-07-20",
    issued_date_bs: "2082-04-04",
    subtotal: 3000,
    discount: 0,
    taxable_amount: 3000,
    vat_rate: 0,
    vat_amount: 0,
    total: 3000,
    amount_in_words: "Three thousand rupees only",
    status: "issued",
    issued_by: "test-user-001",
    print_count: 0,
    created_at: "2025-07-20T04:00:00Z",
    items: [
      {
        line_no: 1,
        description: "Monthly Boxing",
        quantity: 1,
        unit_price: 3000,
        amount: 3000,
      },
    ],
    ...overrides,
  };
}

function mockInvoiceAPI(page: Page) {
  let invoice = baseInvoice();
  return page.route("**/api/v1/orgs/*/invoices**", (route: Route) => {
    const url = new URL(route.request().url());
    const method = route.request().method();

    if (url.pathname.endsWith("/next-number")) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ invoice_number: "2082-83/000001" }),
      });
    }

    if (url.pathname.endsWith("/cancel")) {
      const { reason } = route.request().postDataJSON();
      if (invoice.status === "cancelled") {
        return route.fulfill({
          status: 409,
          contentType: "application/json",
          body: JSON.stringify({ error: "already cancelled", code: "already_cancelled" }),
        });
      }
      invoice = baseInvoice({ status: "cancelled", cancellation_reason: reason });
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(invoice),
      });
    }

    if (method === "POST" && url.pathname.endsWith("/invoices")) {
      const body = route.request().postDataJSON();
      invoice = baseInvoice({ customer_name: body.customer_name });
      return route.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(invoice),
      });
    }

    if (method === "GET" && url.pathname.endsWith("/invoices")) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ data: [invoice], total: 1 }),
      });
    }

    // GET one
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(invoice),
    });
  });
}

test.describe("Bills", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
    await mockInvoiceAPI(page);
  });

  test("issuing a bill lands on its document", async ({ page }) => {
    await page.goto("/en/dashboard/invoices/new");

    await page.getByLabel(/customer name/i).fill("Ram Bahadur");
    await page.getByLabel(/description/i).first().fill("Monthly Boxing");
    await page.getByLabel(/rate/i).first().fill("3000");
    await page.getByRole("button", { name: /issue/i }).click();

    await expect(page).toHaveURL(/\/dashboard\/invoices\/inv-001/);
    await expect(page.getByText("Ram Bahadur")).toBeVisible();
    await expect(page.getByText("2082-83/000001")).toBeVisible();
  });

  test("a cancelled bill is marked and offers no second cancel", async ({ page }) => {
    await page.goto("/en/dashboard/invoices/inv-001");

    await page.getByRole("button", { name: /cancel bill/i }).click();
    await page.getByLabel(/why is this being cancelled/i).fill("wrong customer");
    await page.getByRole("button", { name: /confirm/i }).click();

    await expect(page.getByText(/cancelled/i).first()).toBeVisible();
    await expect(page.getByRole("button", { name: /cancel bill/i })).toHaveCount(0);
  });

  test("the list shows the bill number and customer", async ({ page }) => {
    await page.goto("/en/dashboard/invoices");
    await expect(page.getByText("2082-83/000001")).toBeVisible();
    await expect(page.getByText("Ram Bahadur")).toBeVisible();
  });
});
