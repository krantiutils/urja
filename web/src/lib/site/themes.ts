/**
 * Theme tokens → CSS custom properties.
 *
 * Five visually distinct templates share one set of renderers. That only works
 * if no renderer knows which template it is inside: they read `var(--site-*)`
 * and the layout supplies the values. Keeping the palette out of the Tailwind
 * config is also what lets an admin override individual tokens at runtime.
 */

import type { ThemeTokens, TemplateId } from "@/types/site";

/**
 * Fallback theme. Used when a gym's stored theme is empty or partial — every
 * token must always resolve, or a renderer paints transparent text on
 * transparent background.
 */
export const FALLBACK_THEME: ThemeTokens = {
  bg: "#08080a",
  surface: "#121216",
  fg: "#f5f5f7",
  fg_muted: "#9a9aa4",
  accent: "#dc2626",
  accent_fg: "#ffffff",
  border: "rgba(255,255,255,0.08)",
  radius: "2px",
  font_display: '"Oswald", "Arial Narrow", sans-serif',
  font_body: '"Inter", system-ui, sans-serif',
  display_transform: "uppercase",
  display_weight: "700",
  display_tracking: "0.02em",
};

/**
 * Template themes, mirroring internal/site/templates.go.
 *
 * Used for the builder's template picker previews. At render time the site
 * layout prefers the theme stored on the gym's settings row, because an admin
 * may have customised individual tokens.
 */
export const TEMPLATE_THEMES: Record<TemplateId, ThemeTokens> = {
  fight_club: FALLBACK_THEME,
  iron_sweat: {
    bg: "#f4f4f0",
    surface: "#e6e6e0",
    fg: "#16161a",
    fg_muted: "#55555f",
    accent: "#facc15",
    accent_fg: "#16161a",
    border: "#16161a",
    radius: "0px",
    font_display: '"JetBrains Mono", ui-monospace, monospace',
    font_body: '"Inter", system-ui, sans-serif',
    display_transform: "uppercase",
    display_weight: "800",
    display_tracking: "-0.02em",
  },
  champion: {
    bg: "#fffdf8",
    surface: "#f6f1e7",
    fg: "#1c1917",
    fg_muted: "#78716c",
    accent: "#b45309",
    accent_fg: "#ffffff",
    border: "rgba(28,25,23,0.12)",
    radius: "10px",
    font_display: '"Playfair Display", Georgia, serif',
    font_body: '"Inter", system-ui, sans-serif',
    display_transform: "none",
    display_weight: "600",
    display_tracking: "-0.01em",
  },
  community: {
    bg: "#fffbf7",
    surface: "#fff1e6",
    fg: "#2d1b12",
    fg_muted: "#7c6155",
    accent: "#ea580c",
    accent_fg: "#ffffff",
    border: "rgba(45,27,18,0.10)",
    radius: "18px",
    font_display: '"Inter", system-ui, sans-serif',
    font_body: '"Inter", system-ui, sans-serif',
    display_transform: "none",
    display_weight: "700",
    display_tracking: "-0.02em",
  },
  minimal_pro: {
    bg: "#ffffff",
    surface: "#fafafa",
    fg: "#0a0a0a",
    fg_muted: "#737373",
    accent: "#0a0a0a",
    accent_fg: "#ffffff",
    border: "rgba(10,10,10,0.10)",
    radius: "4px",
    font_display: '"Inter", system-ui, sans-serif',
    font_body: '"Inter", system-ui, sans-serif',
    display_transform: "none",
    display_weight: "500",
    display_tracking: "-0.03em",
  },
};

export const TEMPLATE_LABELS: Record<TemplateId, { en: string; ne: string }> = {
  fight_club: { en: "Fight Club", ne: "फाइट क्लब" },
  iron_sweat: { en: "Iron & Sweat", ne: "आइरन एन्ड स्वेट" },
  champion: { en: "Champion", ne: "च्याम्पियन" },
  community: { en: "Community", ne: "सामुदायिक" },
  minimal_pro: { en: "Minimal Pro", ne: "मिनिमल प्रो" },
};

export const TEMPLATE_IDS: TemplateId[] = [
  "fight_club",
  "iron_sweat",
  "champion",
  "community",
  "minimal_pro",
];

/**
 * Merges a gym's stored (possibly partial) theme over the template default,
 * so a missing or newly-added token always resolves to something sensible.
 */
export function resolveTheme(
  template: TemplateId | undefined,
  stored: Partial<ThemeTokens> | undefined | null
): ThemeTokens {
  const base = (template && TEMPLATE_THEMES[template]) || FALLBACK_THEME;
  if (!stored) return base;

  const merged = { ...base };
  for (const key of Object.keys(base) as (keyof ThemeTokens)[]) {
    const value = stored[key];
    if (typeof value === "string" && value.trim() !== "") {
      merged[key] = value;
    }
  }
  return merged;
}

/** Converts resolved tokens into the inline CSS custom properties the layout sets. */
export function themeToCssVars(theme: ThemeTokens): Record<string, string> {
  return {
    "--site-bg": theme.bg,
    "--site-surface": theme.surface,
    "--site-fg": theme.fg,
    "--site-fg-muted": theme.fg_muted,
    "--site-accent": theme.accent,
    "--site-accent-fg": theme.accent_fg,
    "--site-border": theme.border,
    "--site-radius": theme.radius,
    "--site-font-display": theme.font_display,
    "--site-font-body": theme.font_body,
    "--site-display-transform": theme.display_transform,
    "--site-display-weight": theme.display_weight,
    "--site-display-tracking": theme.display_tracking,
  };
}
