import type { Page } from "@playwright/test";

/**
 * Build a fake JWT token with the given payload.
 * The signature is garbage but the app only decodes the payload (no server verification).
 */
function fakeJWT(payload: Record<string, unknown>): string {
  const header = Buffer.from(JSON.stringify({ alg: "HS256", typ: "JWT" })).toString("base64url");
  const body = Buffer.from(JSON.stringify(payload)).toString("base64url");
  const sig = "fake-sig";
  return `${header}.${body}.${sig}`;
}

async function injectSession(page: Page, role: string): Promise<void> {
  const now = Math.floor(Date.now() / 1000);
  const accessToken = fakeJWT({
    sub: "test-user-001",
    phone: "9841000000",
    role,
    org_id: "org-001",
    iat: now,
    exp: now + 3600, // 1 hour from now
    iss: "urja",
  });
  const refreshToken = fakeJWT({
    sub: "test-user-001",
    type: "refresh",
    iat: now,
    exp: now + 86400,
    iss: "urja",
  });

  await page.evaluate(
    ({ access, refresh }) => {
      localStorage.setItem("access_token", access);
      localStorage.setItem("refresh_token", refresh);
    },
    { access: accessToken, refresh: refreshToken },
  );
}

/**
 * Inject a fake authenticated session into localStorage so AuthGuard lets us through.
 */
export async function injectAuth(page: Page): Promise<void> {
  return injectSession(page, "super_admin");
}

/**
 * Like injectAuth, but with role "admin" instead of "super_admin". Needed for
 * screens gated by a client-side check on `user.role === "admin" || "staff"`
 * (e.g. settings' canEditOrg) — that check reads the JWT's role claim
 * directly and does not treat "super_admin" as satisfying it.
 */
export async function injectAdminAuth(page: Page): Promise<void> {
  return injectSession(page, "admin");
}
