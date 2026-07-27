/**
 * Defensive readers for section content.
 *
 * Section JSON has no migration story: a row written months ago keeps its
 * original shape, and a renderer that assumes a field exists will throw on old
 * data. Everything here returns a safe default rather than throwing, and
 * handles the `field` / `fieldNe` bilingual convention in one place.
 */

import type { Locale } from "@/types";

type Content = Record<string, unknown> | undefined | null;

/**
 * Reads a bilingual string field. For `ne` it prefers `<key>Ne`, falling back to
 * the base field when the translation is missing or blank — a half-translated
 * page should show English, not an empty element.
 */
export function text(content: Content, key: string, locale: Locale): string {
  if (!content) return "";

  if (locale === "ne") {
    const ne = content[`${key}Ne`];
    if (typeof ne === "string" && ne.trim() !== "") return ne;
  }

  const base = content[key];
  return typeof base === "string" ? base : "";
}

/** Reads a plain (non-translated) string field, such as a URL or a time. */
export function str(content: Content, key: string, fallback = ""): string {
  if (!content) return fallback;
  const v = content[key];
  return typeof v === "string" ? v : fallback;
}

/** Reads a number field, tolerating numeric strings from hand-edited JSON. */
export function num(content: Content, key: string): number | null {
  if (!content) return null;
  const v = content[key];
  if (typeof v === "number" && Number.isFinite(v)) return v;
  if (typeof v === "string" && v.trim() !== "") {
    const parsed = Number(v);
    if (Number.isFinite(parsed)) return parsed;
  }
  return null;
}

/** Reads a boolean field. */
export function bool(content: Content, key: string, fallback = false): boolean {
  if (!content) return fallback;
  const v = content[key];
  return typeof v === "boolean" ? v : fallback;
}

/**
 * Reads an array of objects, filtering out anything that is not an object.
 * Always returns an array, so `.map` is safe without a guard at every call site.
 */
export function list(content: Content, key: string): Record<string, unknown>[] {
  if (!content) return [];
  const v = content[key];
  if (!Array.isArray(v)) return [];
  return v.filter(
    (item): item is Record<string, unknown> =>
      typeof item === "object" && item !== null && !Array.isArray(item)
  );
}

/** Reads an array of plain strings. */
export function strList(content: Content, key: string): string[] {
  if (!content) return [];
  const v = content[key];
  if (!Array.isArray(v)) return [];
  return v.filter((item): item is string => typeof item === "string" && item !== "");
}

/** Bilingual read against an arbitrary item rather than the section content. */
export function itemText(
  item: Record<string, unknown> | undefined | null,
  key: string,
  locale: Locale
): string {
  return text(item as Content, key, locale);
}

export interface ButtonSpec {
  label: string;
  href: string;
  style: "solid" | "outline" | "pill";
}

/**
 * Reads a button array, dropping entries without a usable label and forcing the
 * style to a known value.
 */
export function buttons(content: Content, locale: Locale): ButtonSpec[] {
  return list(content, "buttons")
    .map((b) => {
      const style = str(b, "style", "solid");
      return {
        label: itemText(b, "label", locale),
        href: str(b, "href", "#"),
        style: (["solid", "outline", "pill"].includes(style) ? style : "solid") as
          | "solid"
          | "outline"
          | "pill",
      };
    })
    .filter((b) => b.label !== "");
}

/**
 * Renders a restricted subset of Markdown to safe HTML fragments.
 *
 * The API already rejects script and HTML markup on write, but rich text is
 * still user input, so nothing here interpolates raw HTML: the caller receives
 * structured blocks and renders them as React elements.
 */
export type MarkdownBlock =
  | { kind: "heading"; level: 2 | 3; text: string }
  | { kind: "paragraph"; text: string }
  | { kind: "listItem"; text: string };

export function parseMarkdownBlocks(source: string): MarkdownBlock[] {
  if (!source) return [];

  const blocks: MarkdownBlock[] = [];
  for (const rawLine of source.split("\n")) {
    const line = rawLine.trim();
    if (line === "") continue;

    if (line.startsWith("### ")) {
      blocks.push({ kind: "heading", level: 3, text: line.slice(4) });
    } else if (line.startsWith("## ")) {
      blocks.push({ kind: "heading", level: 2, text: line.slice(3) });
    } else if (line.startsWith("- ") || line.startsWith("* ")) {
      blocks.push({ kind: "listItem", text: line.slice(2) });
    } else {
      blocks.push({ kind: "paragraph", text: line });
    }
  }
  return blocks;
}
