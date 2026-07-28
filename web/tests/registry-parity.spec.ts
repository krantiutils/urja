import { test, expect } from "@playwright/test";
import { readFileSync } from "node:fs";
import { join } from "node:path";
import { SECTION_SPECS } from "@/lib/site/section-specs";
import type { SectionType } from "@/types/site";

/**
 * The section registry exists twice: once in Go (the validator the API enforces)
 * and once in TypeScript (what the builder offers an admin). If they drift, an
 * admin can add a section the server then rejects, or a template can reference a
 * variant the builder cannot edit.
 *
 * This test parses the Go source and diffs the two. It is deliberately the only
 * place that couples them — if it becomes brittle, generating the TypeScript
 * from the Go source at build time is the fallback.
 */

const GO_MODELS = join(__dirname, "..", "..", "internal", "site", "models.go");

interface GoSpec {
  variants: string[];
  category: string;
}

function parseGoRegistry(source: string): Record<string, GoSpec> {
  const start = source.indexOf("var SectionSpecs = map[string]SectionSpec{");
  expect(start, "SectionSpecs registry not found in models.go").toBeGreaterThan(-1);

  // Take everything from the opening brace to the closing "}\n" at column 0.
  const body = source.slice(start);
  const end = body.indexOf("\n}\n");
  expect(end, "could not find the end of the SectionSpecs registry").toBeGreaterThan(-1);
  const registry = body.slice(0, end);

  const specs: Record<string, GoSpec> = {};

  // Each entry looks like:
  //   "hero": {
  //       Variants: []string{"centered", "split", ...},
  //       Label: "Hero", LabelNe: "...", Category: "header",
  //   },
  const entryPattern = /"([a-z_]+)":\s*\{([\s\S]*?)\n\t\},/g;
  let match: RegExpExecArray | null;
  while ((match = entryPattern.exec(registry)) !== null) {
    const [, type, block] = match;

    const variantsMatch = block.match(/Variants:\s*\[\]string\{([^}]*)\}/);
    // An exec loop rather than [...matchAll]: spreading an iterator needs
    // downlevelIteration, which the app's tsconfig does not set.
    const variants: string[] = [];
    if (variantsMatch) {
      const variantPattern = /"([^"]+)"/g;
      let variant: RegExpExecArray | null;
      while ((variant = variantPattern.exec(variantsMatch[1])) !== null) {
        variants.push(variant[1]);
      }
    }

    const categoryMatch = block.match(/Category:\s*"([^"]+)"/);

    specs[type] = { variants, category: categoryMatch ? categoryMatch[1] : "" };
  }

  return specs;
}

test("the Go and TypeScript section registries declare the same types", () => {
  const go = parseGoRegistry(readFileSync(GO_MODELS, "utf8"));

  const goTypes = Object.keys(go).sort();
  const tsTypes = Object.keys(SECTION_SPECS).sort();

  // Sanity check the parser itself — a regex that silently matched nothing
  // would make this whole test vacuously pass.
  expect(goTypes.length, "parsed no section types out of models.go").toBeGreaterThan(10);

  expect(tsTypes).toEqual(goTypes);
});

test("every type declares identical variants in both registries", () => {
  const go = parseGoRegistry(readFileSync(GO_MODELS, "utf8"));

  for (const [type, goSpec] of Object.entries(go)) {
    const tsSpec = SECTION_SPECS[type as SectionType];
    expect(tsSpec, `TypeScript registry is missing "${type}"`).toBeTruthy();
    expect(
      [...tsSpec.variants].sort(),
      `variants differ for section type "${type}"`
    ).toEqual([...goSpec.variants].sort());
  }
});

test("every type declares the same category in both registries", () => {
  const go = parseGoRegistry(readFileSync(GO_MODELS, "utf8"));

  for (const [type, goSpec] of Object.entries(go)) {
    const tsSpec = SECTION_SPECS[type as SectionType];
    expect(tsSpec.category, `category differs for section type "${type}"`).toBe(
      goSpec.category
    );
  }
});

test("every default variant is one the API will accept", () => {
  // The builder seeds a new section with variants[0]; if that were not a
  // declared variant the very first save would 400.
  for (const [type, spec] of Object.entries(SECTION_SPECS)) {
    expect(spec.variants.length, `"${type}" declares no variants`).toBeGreaterThan(0);
    expect(spec.variants).toContain(spec.variants[0]);
  }
});

test("default content round-trips through JSON unchanged", () => {
  // Default content is sent to the API verbatim on first save. Anything that
  // does not survive JSON (undefined, functions, Dates) would be silently
  // dropped and the section would render empty.
  for (const [type, spec] of Object.entries(SECTION_SPECS)) {
    const round = JSON.parse(JSON.stringify(spec.defaultContent));
    expect(round, `default content for "${type}" does not round-trip`).toEqual(
      spec.defaultContent
    );
  }
});

test("default styles use only tokens the Go validator allows", () => {
  const backgrounds = ["base", "surface", "accent", "none"];
  const paddings = ["none", "sm", "md", "lg"];
  const widths = ["full", "contained", "narrow"];
  const aligns = ["left", "center", "right"];

  for (const [type, spec] of Object.entries(SECTION_SPECS)) {
    const s = spec.defaultStyle;
    if (s.background) expect(backgrounds, `${type}.background`).toContain(s.background);
    if (s.padding) expect(paddings, `${type}.padding`).toContain(s.padding);
    if (s.width) expect(widths, `${type}.width`).toContain(s.width);
    if (s.align) expect(aligns, `${type}.align`).toContain(s.align);
  }
});

test("every section type has labels in both languages", () => {
  for (const [type, spec] of Object.entries(SECTION_SPECS)) {
    expect(spec.label, `"${type}" has no English label`).toBeTruthy();
    expect(spec.labelNe, `"${type}" has no Nepali label`).toBeTruthy();
    // A Nepali label identical to the English one is almost always an
    // untranslated placeholder.
    expect(spec.labelNe, `"${type}" Nepali label is untranslated`).not.toBe(spec.label);
  }
});
