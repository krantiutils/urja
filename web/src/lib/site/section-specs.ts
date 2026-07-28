/**
 * Section registry — the TypeScript mirror of SectionSpecs in
 * internal/site/models.go.
 *
 * The builder offers exactly what this declares, and the API accepts exactly
 * what the Go registry declares. If the two drift, an admin gets to add a
 * section the server then rejects. `registry-parity.spec.ts` diffs them.
 */

import type { SectionType, SectionStyle } from "@/types/site";

export type SectionCategory =
  | "header"
  | "gym"
  | "content"
  | "media"
  | "social_proof"
  | "conversion"
  | "contact"
  | "layout";

export interface SectionSpec {
  variants: string[];
  label: string;
  labelNe: string;
  category: SectionCategory;
  /** Seed content used when an admin adds this section from the panel. */
  defaultContent: Record<string, unknown>;
  defaultStyle: SectionStyle;
}

const contained: SectionStyle = {
  background: "base",
  padding: "lg",
  width: "contained",
  align: "left",
};

export const SECTION_SPECS: Record<SectionType, SectionSpec> = {
  hero: {
    // "reel" plays gym footage behind the headline; see HeroReel.
    variants: ["centered", "split", "fullbleed", "minimal", "reel"],
    label: "Hero",
    labelNe: "मुख्य ब्यानर",
    category: "header",
    defaultContent: {
      title: "Your gym name",
      titleNe: "तपाईंको जिमको नाम",
      subtitle: "A line about what you do.",
      subtitleNe: "तपाईं के गर्नुहुन्छ भन्ने एक लाइन।",
      image: "",
      buttons: [
        { label: "Get started", labelNe: "सुरु गर्नुहोस्", href: "/contact", style: "solid" },
      ],
    },
    defaultStyle: { background: "base", padding: "lg", width: "full", align: "center" },
  },
  stats_bar: {
    variants: ["inline", "cards", "bordered"],
    label: "Stats Bar",
    labelNe: "तथ्याङ्क",
    category: "social_proof",
    defaultContent: {
      items: [
        { value: "10+", label: "Years", labelNe: "वर्ष" },
        { value: "300+", label: "Members", labelNe: "सदस्य" },
        { value: "6", label: "Days a week", labelNe: "हप्ताको दिन" },
      ],
    },
    defaultStyle: { background: "surface", padding: "md", width: "contained", align: "center" },
  },
  class_timetable: {
    variants: ["table", "day_tabs", "cards"],
    label: "Class Timetable",
    labelNe: "कक्षा तालिका",
    category: "gym",
    defaultContent: {
      title: "Weekly timetable",
      titleNe: "साप्ताहिक तालिका",
      note: "",
      noteNe: "",
      days: [
        {
          day: "Sunday",
          dayNe: "आइतबार",
          classes: [
            { time: "06:00", name: "Morning session", nameNe: "बिहानको सत्र", coach: "", level: "All" },
          ],
        },
      ],
    },
    defaultStyle: contained,
  },
  coaches: {
    variants: ["cards", "list", "spotlight"],
    label: "Coaches",
    labelNe: "प्रशिक्षकहरू",
    category: "gym",
    defaultContent: {
      title: "Your coaches",
      titleNe: "तपाईंका प्रशिक्षकहरू",
      coaches: [
        {
          name: "Coach name",
          nameNe: "प्रशिक्षकको नाम",
          role: "Head Coach",
          roleNe: "प्रमुख प्रशिक्षक",
          bio: "",
          bioNe: "",
          record: "",
          photo: "",
          credentials: [],
        },
      ],
    },
    defaultStyle: contained,
  },
  programs_grid: {
    variants: ["cards", "icons", "list", "numbered"],
    label: "Programs",
    labelNe: "कार्यक्रमहरू",
    category: "gym",
    defaultContent: {
      title: "What we train",
      titleNe: "हामी के तालिम दिन्छौं",
      subtitle: "",
      subtitleNe: "",
      items: [
        {
          icon: "target",
          title: "Programme name",
          titleNe: "कार्यक्रमको नाम",
          description: "",
          descriptionNe: "",
        },
      ],
    },
    defaultStyle: contained,
  },
  membership_plans: {
    variants: ["cards", "table", "compact"],
    label: "Membership Plans",
    labelNe: "सदस्यता योजना",
    category: "gym",
    defaultContent: {
      title: "Membership",
      titleNe: "सदस्यता",
      // "auto" reads the gym's active packages; "manual" uses `plans` below.
      dataSource: "auto",
      note: "",
      noteNe: "",
      plans: [],
    },
    defaultStyle: { background: "surface", padding: "lg", width: "contained", align: "center" },
  },
  gallery: {
    variants: ["grid", "masonry", "carousel"],
    label: "Gallery",
    labelNe: "ग्यालरी",
    category: "media",
    defaultContent: {
      title: "Inside the gym",
      titleNe: "जिम भित्र",
      images: [],
    },
    defaultStyle: contained,
  },
  testimonials: {
    variants: ["cards", "carousel", "quote"],
    label: "Testimonials",
    labelNe: "प्रशंसापत्र",
    category: "social_proof",
    defaultContent: {
      title: "What members say",
      titleNe: "सदस्यहरू के भन्छन्",
      items: [{ name: "", text: "", textNe: "", rating: 5 }],
    },
    defaultStyle: contained,
  },
  faq: {
    variants: ["accordion", "list", "two_column"],
    label: "FAQ",
    labelNe: "प्रश्नोत्तर",
    category: "content",
    defaultContent: {
      title: "Common questions",
      titleNe: "सामान्य प्रश्नहरू",
      items: [{ question: "", questionNe: "", answer: "", answerNe: "" }],
    },
    defaultStyle: { background: "base", padding: "lg", width: "narrow", align: "left" },
  },
  cta_banner: {
    // "solid" is the default variant and leans on the shell's accent; gradient
    // and image paint their own and want background "none".
    variants: ["solid", "gradient", "image", "split"],
    label: "Call to Action",
    labelNe: "कार्य आह्वान",
    category: "conversion",
    defaultContent: {
      title: "Ready to start?",
      titleNe: "सुरु गर्न तयार?",
      subtitle: "",
      subtitleNe: "",
      image: "",
      buttons: [
        { label: "Contact us", labelNe: "सम्पर्क गर्नुहोस्", href: "/contact", style: "solid" },
      ],
    },
    defaultStyle: { background: "accent", padding: "lg", width: "full", align: "center" },
  },
  contact_info: {
    variants: ["list", "card", "two_column"],
    label: "Contact Info",
    labelNe: "सम्पर्क जानकारी",
    category: "contact",
    defaultContent: {
      title: "Find us",
      titleNe: "हामीलाई भेट्नुहोस्",
      address: "",
      addressNe: "",
      phone: "",
      email: "",
      hoursNote: "",
      hoursNoteNe: "",
    },
    defaultStyle: contained,
  },
  map_embed: {
    variants: ["standard", "with_info", "full_width"],
    label: "Map",
    labelNe: "नक्सा",
    category: "contact",
    defaultContent: {
      title: "Find us",
      titleNe: "हामीलाई भेट्नुहोस्",
      address: "",
      addressNe: "",
      latitude: null,
      longitude: null,
    },
    defaultStyle: { background: "base", padding: "none", width: "full", align: "left" },
  },
  lead_form: {
    variants: ["inline", "card", "split"],
    label: "Enquiry Form",
    labelNe: "सोधपुछ फारम",
    category: "conversion",
    defaultContent: {
      title: "Book a free trial",
      titleNe: "निःशुल्क परीक्षण बुक गर्नुहोस्",
      subtitle: "",
      subtitleNe: "",
      submitLabel: "Send",
      submitLabelNe: "पठाउनुहोस्",
      successMessage: "Thanks — we will be in touch.",
      successMessageNe: "धन्यवाद — हामी सम्पर्क गर्नेछौं।",
      interests: [],
    },
    defaultStyle: { background: "surface", padding: "lg", width: "narrow", align: "left" },
  },
  opening_hours: {
    variants: ["table", "list", "compact"],
    label: "Opening Hours",
    labelNe: "खुल्ने समय",
    category: "contact",
    defaultContent: {
      title: "Opening hours",
      titleNe: "खुल्ने समय",
      days: [{ day: "Sunday – Friday", dayNe: "आइतबार – शुक्रबार", hours: "06:00 – 20:00" }],
    },
    defaultStyle: contained,
  },
  logo_strip: {
    variants: ["grid", "marquee", "simple"],
    label: "Logo Strip",
    labelNe: "लोगो पट्टी",
    category: "social_proof",
    defaultContent: {
      title: "Affiliations",
      titleNe: "सम्बद्धता",
      logos: [],
    },
    defaultStyle: { background: "surface", padding: "sm", width: "contained", align: "center" },
  },
  fight_record: {
    variants: ["timeline", "table", "cards"],
    label: "Fight Record",
    labelNe: "लडाइँ रेकर्ड",
    category: "gym",
    defaultContent: {
      title: "Competition record",
      titleNe: "प्रतिस्पर्धा रेकर्ड",
      bouts: [],
    },
    defaultStyle: contained,
  },
  rich_text: {
    variants: ["standard", "two_column", "centered"],
    label: "Text",
    labelNe: "पाठ",
    category: "content",
    defaultContent: {
      title: "",
      titleNe: "",
      body: "Write something about your gym.",
      bodyNe: "आफ्नो जिमको बारेमा केही लेख्नुहोस्।",
    },
    defaultStyle: { background: "base", padding: "md", width: "narrow", align: "left" },
  },
  reel_wall: {
    variants: ["grid", "strip"],
    label: "Reel Wall",
    labelNe: "रिल वाल",
    category: "media",
    defaultContent: {
      title: "Inside the gym",
      titleNe: "जिम भित्र",
      subtitle: "",
      subtitleNe: "",
      items: [],
    },
    defaultStyle: contained,
  },

  media: {
    variants: ["image", "video"],
    label: "Image or Video",
    labelNe: "तस्बिर वा भिडियो",
    category: "media",
    defaultContent: { url: "", caption: "", captionNe: "", alt: "" },
    defaultStyle: contained,
  },
  divider: {
    variants: ["line", "dots", "space"],
    label: "Divider",
    labelNe: "विभाजक",
    category: "layout",
    defaultContent: {},
    defaultStyle: { background: "base", padding: "none", width: "narrow", align: "center" },
  },
};

/** Stable display order for the add-section panel. */
export const SECTION_CATEGORY_ORDER: SectionCategory[] = [
  "header",
  "gym",
  "content",
  "media",
  "social_proof",
  "conversion",
  "contact",
  "layout",
];

export const CATEGORY_LABELS: Record<SectionCategory, { en: string; ne: string }> = {
  header: { en: "Header", ne: "हेडर" },
  gym: { en: "Gym", ne: "जिम" },
  content: { en: "Content", ne: "सामग्री" },
  media: { en: "Media", ne: "मिडिया" },
  social_proof: { en: "Social proof", ne: "प्रमाण" },
  conversion: { en: "Conversion", ne: "रूपान्तरण" },
  contact: { en: "Contact", ne: "सम्पर्क" },
  layout: { en: "Layout", ne: "लेआउट" },
};

/** All section types in a stable order. */
export const SECTION_TYPES = Object.keys(SECTION_SPECS) as SectionType[];

/** Reports whether a type/variant pair is one the API will accept. */
export function isValidVariant(type: SectionType, variant: string): boolean {
  return SECTION_SPECS[type]?.variants.includes(variant) ?? false;
}
