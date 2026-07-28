package site

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Section is one block of a page. Pages store an ordered array of these in a
// JSONB column. The shape is deliberately open — Content varies per Type — so
// ValidateSections is the only thing standing between a request body and the
// public renderer. Keep it strict.
type Section struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Variant string          `json:"variant"`
	Content json.RawMessage `json:"content"`
	Style   SectionStyle    `json:"style"`
	Hidden  bool            `json:"hidden,omitempty"`
	// Version lets a future migration helper distinguish section shapes.
	// Section JSON has no migration story of its own, so renderers must also
	// tolerate missing fields regardless of this value.
	Version int `json:"version"`
}

// SectionStyle holds the presentation controls every section shares. Individual
// renderers must not reimplement these — the shared shell component applies them.
type SectionStyle struct {
	Background string `json:"background"` // base | surface | accent | none
	Padding    string `json:"padding"`    // none | sm | md | lg
	Width      string `json:"width"`      // full | contained | narrow
	Align      string `json:"align"`      // left | center | right
}

// SectionSpec declares what a section type is allowed to look like.
type SectionSpec struct {
	Variants []string `json:"variants"`
	Label    string   `json:"label"`
	LabelNe  string   `json:"label_ne"`
	Category string   `json:"category"`
}

// SectionSpecs is the authoritative registry of section types. The TypeScript
// registry in web/src/lib/site/section-specs.ts mirrors this, and a test
// compares the two — if you add a type here, add it there.
var SectionSpecs = map[string]SectionSpec{
	"hero": {
		Variants: []string{"centered", "split", "fullbleed", "minimal", "gloves"},
		Label:    "Hero", LabelNe: "मुख्य ब्यानर", Category: "header",
	},
	"stats_bar": {
		Variants: []string{"inline", "cards", "bordered"},
		Label:    "Stats Bar", LabelNe: "तथ्याङ्क", Category: "social_proof",
	},
	"class_timetable": {
		Variants: []string{"table", "day_tabs", "cards"},
		Label:    "Class Timetable", LabelNe: "कक्षा तालिका", Category: "gym",
	},
	"coaches": {
		Variants: []string{"cards", "list", "spotlight"},
		Label:    "Coaches", LabelNe: "प्रशिक्षकहरू", Category: "gym",
	},
	"programs_grid": {
		Variants: []string{"cards", "icons", "list", "numbered"},
		Label:    "Programs", LabelNe: "कार्यक्रमहरू", Category: "gym",
	},
	"membership_plans": {
		Variants: []string{"cards", "table", "compact"},
		Label:    "Membership Plans", LabelNe: "सदस्यता योजना", Category: "gym",
	},
	"gallery": {
		Variants: []string{"grid", "masonry", "carousel"},
		Label:    "Gallery", LabelNe: "ग्यालरी", Category: "media",
	},
	"testimonials": {
		Variants: []string{"cards", "carousel", "quote"},
		Label:    "Testimonials", LabelNe: "प्रशंसापत्र", Category: "social_proof",
	},
	"faq": {
		Variants: []string{"accordion", "list", "two_column"},
		Label:    "FAQ", LabelNe: "प्रश्नोत्तर", Category: "content",
	},
	"cta_banner": {
		Variants: []string{"solid", "gradient", "image", "split"},
		Label:    "Call to Action", LabelNe: "कार्य आह्वान", Category: "conversion",
	},
	"contact_info": {
		Variants: []string{"list", "card", "two_column"},
		Label:    "Contact Info", LabelNe: "सम्पर्क जानकारी", Category: "contact",
	},
	"map_embed": {
		Variants: []string{"standard", "with_info", "full_width"},
		Label:    "Map", LabelNe: "नक्सा", Category: "contact",
	},
	"lead_form": {
		Variants: []string{"inline", "card", "split"},
		Label:    "Enquiry Form", LabelNe: "सोधपुछ फारम", Category: "conversion",
	},
	"opening_hours": {
		Variants: []string{"table", "list", "compact"},
		Label:    "Opening Hours", LabelNe: "खुल्ने समय", Category: "contact",
	},
	"logo_strip": {
		Variants: []string{"grid", "marquee", "simple"},
		Label:    "Logo Strip", LabelNe: "लोगो पट्टी", Category: "social_proof",
	},
	"fight_record": {
		Variants: []string{"timeline", "table", "cards"},
		Label:    "Fight Record", LabelNe: "लडाइँ रेकर्ड", Category: "gym",
	},
	"rich_text": {
		Variants: []string{"standard", "two_column", "centered"},
		Label:    "Text", LabelNe: "पाठ", Category: "content",
	},
	"media": {
		Variants: []string{"image", "video"},
		Label:    "Image or Video", LabelNe: "तस्बिर वा भिडियो", Category: "media",
	},
	"divider": {
		Variants: []string{"line", "dots", "space"},
		Label:    "Divider", LabelNe: "विभाजक", Category: "layout",
	},
}

var (
	validBackgrounds = map[string]bool{"base": true, "surface": true, "accent": true, "none": true}
	validPaddings    = map[string]bool{"none": true, "sm": true, "md": true, "lg": true}
	validWidths      = map[string]bool{"full": true, "contained": true, "narrow": true}
	validAligns      = map[string]bool{"left": true, "center": true, "right": true}
)

// maxSectionsPerPage bounds how much work a single page can force the renderer
// to do, and bounds the JSONB row size.
const maxSectionsPerPage = 60

// maxContentBytes bounds a single section's content payload.
const maxContentBytes = 256 * 1024

// scriptPattern catches HTML injection attempts in Markdown fields. Rich text is
// stored as Markdown and sanitized at render; this is defense in depth so a
// payload never reaches storage in the first place.
var scriptPattern = regexp.MustCompile(`(?i)<\s*(script|iframe|object|embed|style)\b|javascript\s*:|\bon[a-z]+\s*=`)

// ValidateSections checks an entire page's section array. It returns the first
// problem found, identifying the offending index so the builder UI can point at it.
func ValidateSections(sections []Section) error {
	if len(sections) > maxSectionsPerPage {
		return fmt.Errorf("too many sections: %d (max %d)", len(sections), maxSectionsPerPage)
	}

	seenIDs := make(map[string]bool, len(sections))
	for i, s := range sections {
		if err := validateSection(s); err != nil {
			return fmt.Errorf("section %d (%s): %w", i, s.Type, err)
		}
		if seenIDs[s.ID] {
			return fmt.Errorf("section %d: duplicate id %q", i, s.ID)
		}
		seenIDs[s.ID] = true
	}
	return nil
}

func validateSection(s Section) error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("id is required")
	}

	spec, ok := SectionSpecs[s.Type]
	if !ok {
		return fmt.Errorf("unknown section type")
	}

	variantOK := false
	for _, v := range spec.Variants {
		if v == s.Variant {
			variantOK = true
			break
		}
	}
	if !variantOK {
		return fmt.Errorf("variant %q is not valid for this type (allowed: %s)",
			s.Variant, strings.Join(spec.Variants, ", "))
	}

	if err := validateStyle(s.Style); err != nil {
		return err
	}

	if len(s.Content) > maxContentBytes {
		return fmt.Errorf("content exceeds %d bytes", maxContentBytes)
	}
	// Content must be valid JSON and must be an object, not a bare scalar or
	// array — every renderer reads named fields.
	if len(s.Content) > 0 {
		var probe map[string]interface{}
		if err := json.Unmarshal(s.Content, &probe); err != nil {
			return fmt.Errorf("content must be a JSON object: %w", err)
		}
		if scriptPattern.Match(s.Content) {
			return fmt.Errorf("content contains disallowed HTML or script markup")
		}
	}

	return nil
}

func validateStyle(st SectionStyle) error {
	// An empty style is legal — the renderer falls back to its defaults.
	if st.Background != "" && !validBackgrounds[st.Background] {
		return fmt.Errorf("invalid background %q", st.Background)
	}
	if st.Padding != "" && !validPaddings[st.Padding] {
		return fmt.Errorf("invalid padding %q", st.Padding)
	}
	if st.Width != "" && !validWidths[st.Width] {
		return fmt.Errorf("invalid width %q", st.Width)
	}
	if st.Align != "" && !validAligns[st.Align] {
		return fmt.Errorf("invalid align %q", st.Align)
	}
	return nil
}

// reservedSlugs can never be used as a page slug: they either collide with
// Next.js internals, with the API, or with the subdomain labels the middleware
// refuses to treat as a tenant.
var reservedSlugs = map[string]bool{
	"api": true, "www": true, "admin": true, "_next": true,
	"dashboard": true, "login": true, "member": true, "static": true,
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// NormalizeSlug lowercases and trims a candidate slug. Callers normalize first,
// then validate — the validator deliberately does not normalize, so a slug that
// passes validation is byte-identical to what gets stored and later looked up.
func NormalizeSlug(slug string) string {
	return strings.TrimSpace(strings.ToLower(slug))
}

// ValidatePageSlug enforces a URL-safe, non-reserved page slug. It rejects
// rather than repairs: see NormalizeSlug.
func ValidatePageSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug is required")
	}
	if len(slug) > 120 {
		return fmt.Errorf("slug is too long (max 120 characters)")
	}
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("slug may contain only lowercase letters, numbers and single hyphens")
	}
	if reservedSlugs[slug] {
		return fmt.Errorf("slug %q is reserved", slug)
	}
	return nil
}

// --- Persistence models ---

// Page is a single page of a gym's public site.
type Page struct {
	ID             string    `json:"id"`
	OrgID          string    `json:"organization_id"`
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	TitleNe        string    `json:"title_ne,omitempty"`
	Sections       []Section `json:"sections"`
	SEODescription string    `json:"seo_description,omitempty"`
	IsPublished    bool      `json:"is_published"`
	ShowInNav      bool      `json:"show_in_nav"`
	SortOrder      int       `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PageSummary is the lightweight form used for nav and listings — no sections.
type PageSummary struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	TitleNe     string `json:"title_ne,omitempty"`
	IsPublished bool   `json:"is_published"`
	ShowInNav   bool   `json:"show_in_nav"`
	SortOrder   int    `json:"sort_order"`
}

// Settings holds a gym's site-wide configuration.
type Settings struct {
	OrgID string `json:"organization_id"`
	// Slug is the gym's subdomain label. Read-only and joined from
	// organizations: the builder has to be able to tell an owner the address
	// their site is published at.
	Slug      string          `json:"slug"`
	Template  string          `json:"template"`
	Theme     json.RawMessage `json:"theme"`
	Nav       json.RawMessage `json:"nav"`
	Footer    json.RawMessage `json:"footer"`
	Socials   json.RawMessage `json:"socials"`
	IsLive    bool            `json:"is_live"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Lead is an enquiry submitted through a public site form.
type Lead struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"organization_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email,omitempty"`
	Message    string    `json:"message,omitempty"`
	Interest   string    `json:"interest,omitempty"`
	SourcePage string    `json:"source_page,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// validLeadStatuses bounds the lead pipeline.
var validLeadStatuses = map[string]bool{
	"new": true, "contacted": true, "trial_booked": true, "joined": true, "lost": true,
}

// ValidateLeadStatus reports whether a status transition target is known.
func ValidateLeadStatus(status string) error {
	if !validLeadStatuses[status] {
		return fmt.Errorf("invalid lead status %q", status)
	}
	return nil
}
