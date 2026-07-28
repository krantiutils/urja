/**
 * Types for the authenticated side of the tenant website: the page builder,
 * the lead inbox, and member boxing profiles.
 *
 * Public site types live in `site.ts`; these are the shapes only an admin or a
 * signed-in member ever sees.
 */

import type { Section, ThemeTokens, TemplateId } from "@/types/site";

export interface SiteTemplateOption {
  id: TemplateId;
  name: string;
  name_ne?: string;
  theme: ThemeTokens;
}

/** Body accepted by create/update page. Every field is sent on every save. */
export interface SitePageInput {
  slug: string;
  title: string;
  title_ne?: string;
  sections: Section[];
  seo_description?: string;
  is_published: boolean;
  show_in_nav: boolean;
  sort_order: number;
}

// --- Boxing ---

export interface BoxingRecord {
  wins: number;
  losses: number;
  draws: number;
  no_contests: number;
  total: number;
}

export interface Bout {
  id: string;
  user_id: string;
  organization_id: string;
  bout_date: string;
  opponent?: string;
  event_name?: string;
  result: "win" | "loss" | "draw" | "no_contest";
  method?: string;
  rounds?: number | null;
  weight_class?: string;
  notes?: string;
  created_at: string;
}

export interface BoxingProfileView {
  id: string;
  user_id: string;
  organization_id: string;
  stance?: string;
  weight_class?: string;
  skill_level?: string;
  sparring_cleared: boolean;
  sparring_cleared_at?: string | null;
  sparring_cleared_by?: string;
  reach_cm?: number | null;
  notes?: string;
  record: BoxingRecord;
  bouts: Bout[];
  /** Derived from the member's latest logged weight; absent if never logged. */
  suggested_weight_class?: string;
}

/**
 * Deliberately has no `sparring_cleared`. Clearance is a coach's safety
 * decision, granted through a separate staff-only endpoint — a field a member
 * can never set should not exist in the shape they submit.
 */
export interface BoxingProfileInput {
  stance: string;
  weight_class: string;
  skill_level: string;
  reach_cm?: number | null;
  notes: string;
}

export interface BoutInput {
  bout_date: string;
  opponent: string;
  event_name: string;
  result: string;
  method: string;
  rounds?: number | null;
  weight_class: string;
  notes: string;
}
