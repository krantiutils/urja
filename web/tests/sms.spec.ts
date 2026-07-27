import { test, expect, type Page } from "@playwright/test";
import { injectAuth } from "./helpers";

const MOCK_MEMBERS = {
  data: [
    {
      id: "m-001",
      phone: "9841000001",
      name: "Ram Shrestha",
      name_ne: null,
      email: null,
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
      role: "member",
      status: "active",
      joined_at: "2026-01-20T08:00:00Z",
    },
  ],
  total: 2,
};

const MOCK_BALANCE = {
  organization_id: "org-001",
  balance: 100,
  total_purchased: 200,
  total_used: 100,
  total_campaigns: 5,
};

/**
 * Mocks the member list + SMS balance/history/send APIs. `sendCalls`
 * records every body posted to /sms/send so tests can assert on exactly
 * what member_ids were sent (or that nothing was sent at all).
 */
async function mockSmsAPIs(page: Page, sendCalls: Record<string, unknown>[]) {
  await page.route("**/api/v1/orgs/*/members**", (route) => {
    if (route.request().method() === "GET") {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_MEMBERS),
      });
    }
    return route.fulfill({ status: 200, body: '{"message":"ok"}' });
  });

  await page.route("**/api/v1/orgs/*/sms/balance", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(MOCK_BALANCE),
    })
  );

  await page.route("**/api/v1/orgs/*/sms/history**", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: [], total: 0 }),
    })
  );

  await page.route("**/api/v1/orgs/*/sms/send", (route) => {
    const body = route.request().postDataJSON();
    sendCalls.push(body);
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ message: "sent" }),
    });
  });
}

test.describe("SMS Page", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/en/login");
    await injectAuth(page);
  });

  test("send is disabled with no recipients and enables once one is picked", async ({ page }) => {
    const sendCalls: Record<string, unknown>[] = [];
    await mockSmsAPIs(page, sendCalls);
    await page.goto("/en/dashboard/sms");

    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    const sendButton = page.getByRole("button", { name: "Send", exact: true });
    await page.locator("#sms-message").fill("Gym closed tomorrow");
    await expect(sendButton).toBeDisabled();

    // Clicking the member row (a <label> wrapping its checkbox) selects it.
    await page.getByText("Ram Shrestha").click();
    await expect(sendButton).toBeEnabled();

    await page.screenshot({
      path: "screenshots/sms-send-enabled-en.png",
      fullPage: true,
    });
  });

  test("select all sets member ids explicitly and confirming shows the recipient count", async ({ page }) => {
    const sendCalls: Record<string, unknown>[] = [];
    await mockSmsAPIs(page, sendCalls);
    await page.goto("/en/dashboard/sms");

    await expect(page.getByText("Ram Shrestha")).toBeVisible();
    await expect(page.getByText("Sita Maharjan")).toBeVisible();

    await page.locator("#sms-message").fill("Gym closed tomorrow");
    await page.getByRole("button", { name: "Select All" }).click();
    await expect(page.getByRole("button", { name: "Deselect All" })).toBeVisible();

    await page.getByRole("button", { name: "Send", exact: true }).click();

    await expect(page.getByTestId("sms-confirm-modal")).toBeVisible();
    await expect(page.getByText("Confirm SMS Send")).toBeVisible();
    await expect(page.getByTestId("sms-confirm-recipients-count")).toHaveText("2");
    await expect(page.getByTestId("sms-confirm-credits-count")).toHaveText("2");

    await page.getByTestId("sms-confirm-send-btn").click();

    await expect(page.getByTestId("sms-confirm-modal")).not.toBeVisible();
    expect(sendCalls).toHaveLength(1);
    const body = sendCalls[0] as { member_ids: string[]; message: string };
    expect(new Set(body.member_ids)).toEqual(new Set(["m-001", "m-002"]));
    expect(body.member_ids.length).toBe(2);

    await page.screenshot({
      path: "screenshots/sms-confirm-sent-en.png",
      fullPage: true,
    });
  });

  test("cancelling the confirmation sends nothing", async ({ page }) => {
    const sendCalls: Record<string, unknown>[] = [];
    await mockSmsAPIs(page, sendCalls);
    await page.goto("/en/dashboard/sms");

    await expect(page.getByText("Ram Shrestha")).toBeVisible();

    await page.locator("#sms-message").fill("Gym closed tomorrow");
    await page.getByText("Ram Shrestha").click();

    await page.getByRole("button", { name: "Send", exact: true }).click();
    await expect(page.getByTestId("sms-confirm-modal")).toBeVisible();
    await expect(page.getByTestId("sms-confirm-recipients-count")).toHaveText("1");

    await page.getByTestId("sms-confirm-cancel-btn").click();
    await expect(page.getByTestId("sms-confirm-modal")).not.toBeVisible();

    expect(sendCalls).toHaveLength(0);

    await page.screenshot({
      path: "screenshots/sms-confirm-cancelled-en.png",
      fullPage: true,
    });
  });
});
