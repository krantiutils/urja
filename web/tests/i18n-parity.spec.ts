import { test, expect } from "@playwright/test";
import { getDictionary, locales } from "@/lib/i18n";

/**
 * The en and ne dictionaries must stay structurally identical.
 *
 * A key present in one but not the other renders as `undefined` to users in
 * that locale — silently, with no type error, because the dictionary is a plain
 * object literal. This is the guard.
 */

type Leaf = string;

function flatten(obj: unknown, prefix = ""): Record<string, Leaf> {
  const out: Record<string, Leaf> = {};
  if (obj === null || typeof obj !== "object") return out;

  for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
    const path = prefix ? `${prefix}.${key}` : key;
    if (value !== null && typeof value === "object" && !Array.isArray(value)) {
      Object.assign(out, flatten(value, path));
    } else {
      out[path] = String(value);
    }
  }
  return out;
}

test("every locale dictionary has exactly the same keys", () => {
  const flattened = locales.map((l) => ({ locale: l, keys: flatten(getDictionary(l)) }));

  const [reference, ...rest] = flattened;
  const referenceKeys = Object.keys(reference.keys).sort();

  expect(referenceKeys.length, "dictionary looks empty — is the export shape right?")
    .toBeGreaterThan(100);

  for (const other of rest) {
    const otherKeys = Object.keys(other.keys).sort();

    const missing = referenceKeys.filter((k) => !(k in other.keys));
    const extra = otherKeys.filter((k) => !(k in reference.keys));

    expect(
      missing,
      `keys present in "${reference.locale}" but missing from "${other.locale}" — these render as undefined`
    ).toEqual([]);
    expect(
      extra,
      `keys present in "${other.locale}" but missing from "${reference.locale}"`
    ).toEqual([]);
  }
});

test("no dictionary value is empty", () => {
  for (const locale of locales) {
    const flat = flatten(getDictionary(locale));
    const blank = Object.entries(flat)
      .filter(([, v]) => v.trim() === "")
      .map(([k]) => k);
    expect(blank, `empty strings in the "${locale}" dictionary`).toEqual([]);
  }
});

test("the marketing landing page no longer references SaaS pricing", () => {
  // Phase 5 removed the pricing section. Its dictionary keys went with it —
  // orphaned keys are how a deleted feature quietly comes back.
  for (const locale of locales) {
    const flat = flatten(getDictionary(locale));
    const pricingKeys = Object.keys(flat).filter((k) =>
      /^landing\.(pricing|navPricing)/.test(k)
    );
    expect(pricingKeys, `orphaned SaaS pricing keys in "${locale}"`).toEqual([]);
  }
});
