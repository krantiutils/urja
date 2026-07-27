import { test, expect } from "@playwright/test";
import {
  extractSubdomain,
  resolveSiteRewrite,
  siteRewritePath,
  isAppOnlyPath,
  apexUrl,
} from "@/lib/subdomain";

// Pure unit tests for tenant subdomain resolution. No browser needed — these
// guard the routing rules that decide whether a request is a gym's public site
// or the authenticated app.

test.describe("extractSubdomain", () => {
  test("pulls the slug from a tenant host", () => {
    expect(extractSubdomain("ibckirtipur.nepalgym.xyz")).toBe("ibckirtipur");
    expect(extractSubdomain("pimbahal-gym.nepalgym.xyz")).toBe("pimbahal-gym");
  });

  test("returns null for the apex and reserved labels", () => {
    expect(extractSubdomain("nepalgym.xyz")).toBeNull();
    expect(extractSubdomain("www.nepalgym.xyz")).toBeNull();
    expect(extractSubdomain("api.nepalgym.xyz")).toBeNull();
    expect(extractSubdomain("admin.nepalgym.xyz")).toBeNull();
  });

  test("strips the port", () => {
    expect(extractSubdomain("ibckirtipur.nepalgym.xyz:443")).toBe("ibckirtipur");
    expect(extractSubdomain("ibckirtipur.localhost:3000")).toBe("ibckirtipur");
  });

  test("is case insensitive", () => {
    expect(extractSubdomain("IBCKirtipur.NepalGym.xyz")).toBe("ibckirtipur");
  });

  test("rejects multi-level labels", () => {
    // Accepting these would make a slug ambiguous.
    expect(extractSubdomain("a.b.nepalgym.xyz")).toBeNull();
  });

  test("rejects unrelated domains and junk", () => {
    expect(extractSubdomain("example.com")).toBeNull();
    expect(extractSubdomain("evil.nepalgym.xyz.attacker.com")).toBeNull();
    expect(extractSubdomain("under_score.nepalgym.xyz")).toBeNull();
    expect(extractSubdomain("-leading.nepalgym.xyz")).toBeNull();
    expect(extractSubdomain("")).toBeNull();
    expect(extractSubdomain(null)).toBeNull();
    expect(extractSubdomain(undefined)).toBeNull();
  });
});

test.describe("resolveSiteRewrite", () => {
  test("defaults the root path to the home page in the default locale", () => {
    expect(resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/")).toEqual({
      slug: "ibckirtipur",
      locale: "en",
      pagePath: "home",
    });
  });

  test("honours an explicit locale prefix", () => {
    expect(resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/ne/coaches")).toEqual({
      slug: "ibckirtipur",
      locale: "ne",
      pagePath: "coaches",
    });
  });

  test("treats a bare path as the default locale", () => {
    expect(resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/classes")).toEqual({
      slug: "ibckirtipur",
      locale: "en",
      pagePath: "classes",
    });
  });

  test("returns null on the apex host", () => {
    expect(resolveSiteRewrite("nepalgym.xyz", "/en/dashboard")).toBeNull();
    expect(resolveSiteRewrite("www.nepalgym.xyz", "/")).toBeNull();
  });

  test("builds the internal rewrite path", () => {
    const r = resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/ne/contact")!;
    expect(siteRewritePath(r)).toBe("/ne/site/ibckirtipur/contact");
  });

  test("root rewrites to the home page route", () => {
    const r = resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/")!;
    expect(siteRewritePath(r)).toBe("/en/site/ibckirtipur/home");
  });
});

test.describe("isAppOnlyPath", () => {
  test("recognises authenticated app paths with and without a locale", () => {
    expect(isAppOnlyPath("/dashboard")).toBe(true);
    expect(isAppOnlyPath("/en/dashboard")).toBe(true);
    expect(isAppOnlyPath("/ne/dashboard/members")).toBe(true);
    expect(isAppOnlyPath("/en/member/workouts")).toBe(true);
    expect(isAppOnlyPath("/en/login")).toBe(true);
  });

  test("treats public site paths as not app-only", () => {
    expect(isAppOnlyPath("/")).toBe(false);
    expect(isAppOnlyPath("/en")).toBe(false);
    expect(isAppOnlyPath("/coaches")).toBe(false);
    expect(isAppOnlyPath("/ne/classes")).toBe(false);
    // A gym could legitimately name a page "membership" — only the exact
    // "member" segment is the app.
    expect(isAppOnlyPath("/en/membership")).toBe(false);
  });
});

test.describe("apexUrl", () => {
  test("preserves the incoming protocol", () => {
    expect(apexUrl("/en/dashboard", "https:")).toBe("https://nepalgym.xyz/en/dashboard");
    expect(apexUrl("/en/dashboard", "http")).toBe("http://nepalgym.xyz/en/dashboard");
  });
});
