import { test, expect } from "@playwright/test";
import { pageHref, siteHref } from "@/lib/site/links";
import { buttons } from "@/lib/site/content";

// Pure unit tests for tenant site link construction. Content is authored once
// and read in both languages, so an internal href that keeps its English path
// silently drops Nepali visitors onto English pages.

test.describe("pageHref", () => {
  test("maps the home slug to the site root", () => {
    expect(pageHref("en", "home")).toBe("/");
    expect(pageHref("ne", "home")).toBe("/ne");
  });

  test("prefixes the locale for other pages", () => {
    expect(pageHref("en", "classes")).toBe("/classes");
    expect(pageHref("ne", "classes")).toBe("/ne/classes");
  });
});

test.describe("siteHref", () => {
  test("leaves English paths alone", () => {
    expect(siteHref("en", "/contact")).toBe("/contact");
    expect(siteHref("en", "/")).toBe("/");
  });

  test("prefixes internal paths in Nepali", () => {
    expect(siteHref("ne", "/contact")).toBe("/ne/contact");
    expect(siteHref("ne", "/")).toBe("/ne");
  });

  test("does not double-prefix an href that already carries the locale", () => {
    expect(siteHref("ne", "/ne")).toBe("/ne");
    expect(siteHref("ne", "/ne/contact")).toBe("/ne/contact");
  });

  test("never rewrites links that leave the site", () => {
    for (const href of [
      "https://example.com/x",
      "http://example.com",
      "//cdn.example.com/a.png",
      "mailto:hi@example.com",
      "tel:+9779800000000",
      "#timetable",
    ]) {
      expect(siteHref("ne", href)).toBe(href);
    }
  });

  test("treats a blank href as a dead link rather than the site root", () => {
    // "" would otherwise become "/ne", silently sending a visitor home.
    expect(siteHref("ne", "")).toBe("#");
    expect(siteHref("en", "   ")).toBe("#");
  });
});

test.describe("buttons", () => {
  const content = {
    buttons: [
      { label: "Book a free trial", labelNe: "निःशुल्क परीक्षण", href: "/contact", style: "solid" },
      { label: "Call us", href: "tel:+9779800000000", style: "outline" },
    ],
  };

  test("localizes internal CTA hrefs", () => {
    expect(buttons(content, "ne").map((b) => b.href)).toEqual([
      "/ne/contact",
      "tel:+9779800000000",
    ]);
    expect(buttons(content, "en").map((b) => b.href)).toEqual([
      "/contact",
      "tel:+9779800000000",
    ]);
  });

  test("uses the Nepali label when present", () => {
    expect(buttons(content, "ne")[0].label).toBe("निःशुल्क परीक्षण");
    // No labelNe on the second button, so it falls back to English rather than
    // rendering an empty button.
    expect(buttons(content, "ne")[1].label).toBe("Call us");
  });
});
