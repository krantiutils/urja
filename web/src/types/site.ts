/**
 * Tenant site types.
 *
 * These mirror internal/site/models.go. The API validates every section on
 * write, so anything reaching a renderer has already passed the Go allow-lists —
 * but renderers must still tolerate missing content fields, because section
 * JSON has no migration story and old rows keep their original shape.
 */

export type SectionType =
  | "hero"
  | "stats_bar"
  | "class_timetable"
  | "coaches"
  | "programs_grid"
  | "membership_plans"
  | "gallery"
  | "testimonials"
  | "faq"
  | "cta_banner"
  | "contact_info"
  | "map_embed"
  | "lead_form"
  | "opening_hours"
  | "logo_strip"
  | "fight_record"
  | "rich_text"
  | "media"
  | "reel_wall"
  | "divider";

export type SectionBackground = "base" | "surface" | "accent" | "none";
export type SectionPadding = "none" | "sm" | "md" | "lg";
export type SectionWidth = "full" | "contained" | "narrow";
export type SectionAlign = "left" | "center" | "right";

export interface SectionStyle {
  background?: SectionBackground;
  padding?: SectionPadding;
  width?: SectionWidth;
  align?: SectionAlign;
}

export interface Section {
  id: string;
  type: SectionType;
  variant: string;
  /** Shape varies per type. Renderers read named fields defensively. */
  content: Record<string, unknown>;
  style: SectionStyle;
  hidden?: boolean;
  version?: number;
}

export interface PageSummary {
  id: string;
  slug: string;
  title: string;
  title_ne?: string;
  is_published: boolean;
  show_in_nav: boolean;
  sort_order: number;
}

export interface SitePage extends PageSummary {
  organization_id: string;
  sections: Section[];
  seo_description?: string;
  created_at: string;
  updated_at: string;
}

export interface ThemeTokens {
  bg: string;
  surface: string;
  fg: string;
  fg_muted: string;
  accent: string;
  accent_fg: string;
  border: string;
  radius: string;
  font_display: string;
  font_body: string;
  display_transform: string;
  display_weight: string;
  display_tracking: string;
}

export interface SiteNav {
  /** Extra links beyond the pages flagged show_in_nav. */
  links?: Array<{ label: string; labelNe?: string; href: string }>;
  ctaLabel?: string;
  ctaLabelNe?: string;
  ctaHref?: string;
}

export interface SiteFooter {
  tagline?: string;
  taglineNe?: string;
  note?: string;
  noteNe?: string;
}

export interface SiteSocials {
  facebook?: string;
  instagram?: string;
  youtube?: string;
  tiktok?: string;
  whatsapp?: string;
}

export interface PublicSite {
  slug: string;
  org_name: string;
  org_name_ne?: string;
  template: TemplateId;
  theme: Partial<ThemeTokens>;
  nav: SiteNav;
  footer: SiteFooter;
  socials: SiteSocials;
  pages: PageSummary[];
  /** The gym's public contact detail, read from its own record so every page
   *  — the home page above all — can emit it as structured data. */
  contact?: SiteContact | null;
}

export interface SiteContact {
  address?: string;
  address_ne?: string;
  phone?: string;
  email?: string;
  latitude?: number;
  longitude?: number;
}

export type TemplateId =
  | "fight_club"
  | "iron_sweat"
  | "champion"
  | "community"
  | "minimal_pro";

export interface SiteSettings {
  organization_id: string;
  /** The gym's subdomain label. Read-only; set when the organization is created. */
  slug: string;
  template: TemplateId;
  theme: Partial<ThemeTokens>;
  nav: SiteNav;
  footer: SiteFooter;
  socials: SiteSocials;
  is_live: boolean;
  updated_at: string;
}

export type LeadStatus = "new" | "contacted" | "trial_booked" | "joined" | "lost";

export interface SiteLead {
  id: string;
  organization_id: string;
  name: string;
  phone: string;
  email?: string;
  message?: string;
  interest?: string;
  source_page?: string;
  status: LeadStatus;
  created_at: string;
}
