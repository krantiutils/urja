package site

import (
	"encoding/json"
	"strings"
	"testing"
)

// okSection builds a minimally valid section of the given type, picking that
// type's first declared variant.
func okSection(id, typ string) Section {
	return Section{
		ID:      id,
		Type:    typ,
		Variant: SectionSpecs[typ].Variants[0],
		Content: json.RawMessage(`{}`),
		Style:   SectionStyle{Background: "base", Padding: "md", Width: "contained", Align: "left"},
	}
}

func TestValidateSections_AcceptsEveryDeclaredTypeAndVariant(t *testing.T) {
	// This is the test that catches the registry and a renderer disagreeing
	// about what exists. Every declared combination must validate.
	for typ, spec := range SectionSpecs {
		for _, variant := range spec.Variants {
			s := okSection("sec-1", typ)
			s.Variant = variant
			if err := ValidateSections([]Section{s}); err != nil {
				t.Errorf("type %q variant %q should be valid: %v", typ, variant, err)
			}
		}
	}
}

func TestSectionSpecs_EveryTypeDeclaresAVariantAndLabels(t *testing.T) {
	for typ, spec := range SectionSpecs {
		if len(spec.Variants) == 0 {
			t.Errorf("type %q declares no variants", typ)
		}
		if spec.Label == "" || spec.LabelNe == "" {
			t.Errorf("type %q is missing an English or Nepali label", typ)
		}
		if spec.Category == "" {
			t.Errorf("type %q has no category, so the add panel cannot group it", typ)
		}
		seen := map[string]bool{}
		for _, v := range spec.Variants {
			if seen[v] {
				t.Errorf("type %q lists variant %q twice", typ, v)
			}
			seen[v] = true
		}
	}
}

func TestValidateSections_RejectsUnknownType(t *testing.T) {
	s := Section{ID: "a", Type: "arbitrary_code", Variant: "x", Content: json.RawMessage(`{}`)}
	err := ValidateSections([]Section{s})
	if err == nil {
		t.Fatal("unknown section type must be rejected")
	}
	if !strings.Contains(err.Error(), "unknown section type") {
		t.Errorf("error should name the problem, got: %v", err)
	}
}

func TestValidateSections_RejectsUnknownVariantForType(t *testing.T) {
	s := okSection("a", "hero")
	s.Variant = "nonsense"
	if err := ValidateSections([]Section{s}); err == nil {
		t.Fatal("a variant not declared for the type must be rejected")
	}
}

func TestValidateSections_RejectsVariantBorrowedFromAnotherType(t *testing.T) {
	// "accordion" is valid for faq but not for hero. Cross-type variant leakage
	// is the most likely real-world mistake, so pin it explicitly.
	s := okSection("a", "hero")
	s.Variant = "accordion"
	if err := ValidateSections([]Section{s}); err == nil {
		t.Fatal("hero must not accept the faq variant 'accordion'")
	}
}

func TestValidateSections_RejectsMissingID(t *testing.T) {
	s := okSection("", "hero")
	if err := ValidateSections([]Section{s}); err == nil {
		t.Fatal("a section without an id must be rejected")
	}
}

func TestValidateSections_RejectsDuplicateIDs(t *testing.T) {
	a := okSection("same", "hero")
	b := okSection("same", "faq")
	if err := ValidateSections([]Section{a, b}); err == nil {
		t.Fatal("duplicate section ids must be rejected — reorder and delete key on id")
	}
}

func TestValidateSections_RejectsInvalidStyleTokens(t *testing.T) {
	cases := []struct {
		name  string
		style SectionStyle
	}{
		{"background", SectionStyle{Background: "chartreuse"}},
		{"padding", SectionStyle{Padding: "enormous"}},
		{"width", SectionStyle{Width: "wide"}},
		{"align", SectionStyle{Align: "justify"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := okSection("a", "hero")
			s.Style = tc.style
			if err := ValidateSections([]Section{s}); err == nil {
				t.Fatalf("invalid %s token must be rejected", tc.name)
			}
		})
	}
}

func TestValidateSections_AcceptsEmptyStyle(t *testing.T) {
	// An unset style is legal: renderers fall back to their own defaults.
	s := okSection("a", "hero")
	s.Style = SectionStyle{}
	if err := ValidateSections([]Section{s}); err != nil {
		t.Fatalf("an empty style should be allowed: %v", err)
	}
}

func TestValidateSections_RejectsScriptMarkupInContent(t *testing.T) {
	payloads := []string{
		`{"body":"<script>alert(1)</script>"}`,
		`{"body":"<iframe src=\"evil\"></iframe>"}`,
		`{"href":"javascript:alert(1)"}`,
		`{"body":"<img src=x onerror=alert(1)>"}`,
	}
	for _, p := range payloads {
		s := okSection("a", "rich_text")
		s.Content = json.RawMessage(p)
		if err := ValidateSections([]Section{s}); err == nil {
			t.Errorf("payload should have been rejected: %s", p)
		}
	}
}

func TestValidateSections_AcceptsOrdinaryMarkdown(t *testing.T) {
	s := okSection("a", "rich_text")
	s.Content = json.RawMessage(`{"body":"## Train with us\n\n**Boxing** since 2015 — [join](/contact)."}`)
	if err := ValidateSections([]Section{s}); err != nil {
		t.Fatalf("ordinary Markdown must be accepted: %v", err)
	}
}

func TestValidateSections_RejectsNonObjectContent(t *testing.T) {
	for _, bad := range []string{`"a string"`, `42`, `[1,2,3]`, `{oops`} {
		s := okSection("a", "hero")
		s.Content = json.RawMessage(bad)
		if err := ValidateSections([]Section{s}); err == nil {
			t.Errorf("content %s must be rejected — renderers read named fields", bad)
		}
	}
}

func TestValidateSections_RejectsTooManySections(t *testing.T) {
	many := make([]Section, maxSectionsPerPage+1)
	for i := range many {
		many[i] = okSection(string(rune('a'+i%26))+string(rune('0'+i/26)), "divider")
	}
	if err := ValidateSections(many); err == nil {
		t.Fatal("a page over the section cap must be rejected")
	}
}

func TestValidateSections_ErrorNamesTheOffendingIndex(t *testing.T) {
	good := okSection("a", "hero")
	bad := okSection("b", "faq")
	bad.Variant = "nope"
	err := ValidateSections([]Section{good, bad})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "section 1") {
		t.Errorf("error must identify the offending index so the builder can point at it, got: %v", err)
	}
}

func TestValidatePageSlug(t *testing.T) {
	valid := []string{"home", "classes", "our-coaches", "class-timetable-2026"}
	for _, s := range valid {
		if err := ValidatePageSlug(s); err != nil {
			t.Errorf("slug %q should be valid: %v", s, err)
		}
	}

	invalid := []string{
		"",                       // empty
		"Home",                   // uppercase
		"our coaches",            // space
		"-leading",               // leading hyphen
		"trailing-",              // trailing hyphen
		"double--dash",           // consecutive hyphens
		"under_score",            // underscore
		"api",                    // reserved
		"_next",                  // reserved
		"dashboard",              // reserved
		strings.Repeat("a", 121), // too long
	}
	for _, s := range invalid {
		if err := ValidatePageSlug(s); err == nil {
			t.Errorf("slug %q should be rejected", s)
		}
	}
}

func TestValidateLeadStatus(t *testing.T) {
	for _, s := range []string{"new", "contacted", "trial_booked", "joined", "lost"} {
		if err := ValidateLeadStatus(s); err != nil {
			t.Errorf("status %q should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"", "NEW", "deleted", "spam"} {
		if err := ValidateLeadStatus(s); err == nil {
			t.Errorf("status %q should be rejected", s)
		}
	}
}

func TestNormalizeSlug(t *testing.T) {
	// Normalization is a separate, explicit step so that a slug which passes
	// validation is byte-identical to what is stored and later looked up.
	cases := map[string]string{
		"Home":        "home",
		"  Classes  ": "classes",
		"OUR-COACHES": "our-coaches",
		"contact":     "contact",
	}
	for in, want := range cases {
		if got := NormalizeSlug(in); got != want {
			t.Errorf("NormalizeSlug(%q) = %q, want %q", in, got, want)
		}
	}

	// Normalize then validate is the intended pipeline and must succeed.
	if err := ValidatePageSlug(NormalizeSlug("  Our-Coaches ")); err != nil {
		t.Errorf("normalize-then-validate should succeed: %v", err)
	}
}
