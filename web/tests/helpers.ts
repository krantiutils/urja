/**
 * Build a fake JWT accepted by the client-side auth code.
 *
 * The app only base64-decodes the payload section (`atob(parts[1])`)
 * and checks `exp` for token validity.  Header and signature are
 * ignored beyond the three-part split, so we can use stubs.
 */
export function fakeJWT(overrides: Record<string, unknown> = {}): string {
  const header = Buffer.from(JSON.stringify({ alg: "HS256", typ: "JWT" })).toString("base64");
  const payload = Buffer.from(
    JSON.stringify({
      sub: "test-user-1",
      phone: "9841234567",
      role: "admin",
      org_id: "test-org-1",
      iat: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + 86400, // 24 hours from now
      iss: "urja-test",
      ...overrides,
    })
  ).toString("base64");
  const signature = Buffer.from("fake-sig").toString("base64");
  return `${header}.${payload}.${signature}`;
}

/**
 * Inject auth tokens into localStorage so AuthGuard lets us through.
 *
 * Must be called via `page.addInitScript` *before* navigating to a
 * protected route, because the AuthProvider reads localStorage on mount.
 */
export function injectAuthScript(token: string): string {
  return `
    localStorage.setItem('access_token', '${token}');
    localStorage.setItem('refresh_token', 'fake-refresh-token');
  `;
}
