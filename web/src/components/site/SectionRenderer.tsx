import type { Section } from "@/types/site";
import type { Locale } from "@/types";
import { SectionShell } from "@/components/site/primitives/SectionShell";

import ClassTimetableRenderer from "@/components/site/sections/ClassTimetableRenderer";
import CoachesRenderer from "@/components/site/sections/CoachesRenderer";
import ContactInfoRenderer from "@/components/site/sections/ContactInfoRenderer";
import CtaBannerRenderer from "@/components/site/sections/CtaBannerRenderer";
import DividerRenderer from "@/components/site/sections/DividerRenderer";
import FaqRenderer from "@/components/site/sections/FaqRenderer";
import FightRecordRenderer from "@/components/site/sections/FightRecordRenderer";
import GalleryRenderer from "@/components/site/sections/GalleryRenderer";
import HeroRenderer from "@/components/site/sections/HeroRenderer";
import LeadFormRenderer from "@/components/site/sections/LeadFormRenderer";
import LogoStripRenderer from "@/components/site/sections/LogoStripRenderer";
import MapEmbedRenderer from "@/components/site/sections/MapEmbedRenderer";
import MediaRenderer from "@/components/site/sections/MediaRenderer";
import MembershipPlansRenderer from "@/components/site/sections/MembershipPlansRenderer";
import OpeningHoursRenderer from "@/components/site/sections/OpeningHoursRenderer";
import ProgramsGridRenderer from "@/components/site/sections/ProgramsGridRenderer";
import RichTextRenderer from "@/components/site/sections/RichTextRenderer";
import StatsBarRenderer from "@/components/site/sections/StatsBarRenderer";
import TestimonialsRenderer from "@/components/site/sections/TestimonialsRenderer";

/**
 * Dispatches a section to its renderer.
 *
 * Every renderer returns only its inner content; SectionShell applies the
 * background, padding, width and alignment once, here — so the inspector can
 * restyle any section without each renderer reimplementing those four
 * controls.
 *
 * An unknown type renders nothing rather than throwing. The API validates
 * types on write, so this can only happen if a section type is retired while
 * pages still reference it, and one stale section must not take down the
 * whole page.
 */
export interface SectionRendererProps {
  section: Section;
  locale: Locale;
  siteSlug: string;
  packages: Array<Record<string, unknown>>;
  /** Page slug, recorded against leads so a gym knows which page converted. */
  sourcePage?: string;
}

export function SectionRenderer({
  section,
  locale,
  siteSlug,
  packages,
  sourcePage,
}: SectionRendererProps) {
  if (section.hidden) return null;

  const body = renderBody({ section, locale, siteSlug, packages, sourcePage });
  if (body === null) return null;

  return (
    <SectionShell style={section.style} sectionId={section.id} sectionType={section.type}>
      {body}
    </SectionShell>
  );
}

function renderBody({
  section,
  locale,
  siteSlug,
  packages,
  sourcePage,
}: SectionRendererProps) {
  switch (section.type) {
    case "hero":
      return <HeroRenderer section={section} locale={locale} />;
    case "stats_bar":
      return <StatsBarRenderer section={section} locale={locale} />;
    case "class_timetable":
      return <ClassTimetableRenderer section={section} locale={locale} />;
    case "coaches":
      return <CoachesRenderer section={section} locale={locale} />;
    case "programs_grid":
      return <ProgramsGridRenderer section={section} locale={locale} />;
    case "membership_plans":
      return <MembershipPlansRenderer section={section} locale={locale} packages={packages} />;
    case "gallery":
      return <GalleryRenderer section={section} locale={locale} />;
    case "testimonials":
      return <TestimonialsRenderer section={section} locale={locale} />;
    case "faq":
      return <FaqRenderer section={section} locale={locale} />;
    case "cta_banner":
      return <CtaBannerRenderer section={section} locale={locale} />;
    case "contact_info":
      return <ContactInfoRenderer section={section} locale={locale} />;
    case "map_embed":
      return <MapEmbedRenderer section={section} locale={locale} />;
    case "lead_form":
      return (
        <LeadFormRenderer
          section={section}
          locale={locale}
          siteSlug={siteSlug}
          sourcePage={sourcePage}
        />
      );
    case "opening_hours":
      return <OpeningHoursRenderer section={section} locale={locale} />;
    case "logo_strip":
      return <LogoStripRenderer section={section} locale={locale} />;
    case "fight_record":
      return <FightRecordRenderer section={section} locale={locale} />;
    case "rich_text":
      return <RichTextRenderer section={section} locale={locale} />;
    case "media":
      return <MediaRenderer section={section} locale={locale} />;
    case "divider":
      return <DividerRenderer section={section} locale={locale} />;
    default:
      return null;
  }
}
