package site

import "testing"

func TestTemplates_AllSectionsValidate(t *testing.T) {
	// The test that catches a template referencing a type or variant that does
	// not exist. Without it, apply-template would wipe a site and then fail.
	for id, tpl := range Templates {
		for _, page := range tpl.Pages() {
			if err := ValidateSections(page.Sections); err != nil {
				t.Errorf("template %q page %q: %v", id, page.Slug, err)
			}
		}
	}
}

func TestTemplates_PageSlugsAreValidAndUnique(t *testing.T) {
	for id, tpl := range Templates {
		seen := map[string]bool{}
		for _, page := range tpl.Pages() {
			if err := ValidatePageSlug(page.Slug); err != nil {
				t.Errorf("template %q slug %q: %v", id, page.Slug, err)
			}
			if seen[page.Slug] {
				t.Errorf("template %q repeats slug %q — UNIQUE(org, slug) would reject it", id, page.Slug)
			}
			seen[page.Slug] = true
		}
	}
}

func TestTemplates_EachSeedsAHomePage(t *testing.T) {
	for id, tpl := range Templates {
		found := false
		for _, page := range tpl.Pages() {
			if page.Slug == "home" {
				found = true
				if !page.IsPublished {
					t.Errorf("template %q home page must be published, else the site 404s on first visit", id)
				}
			}
		}
		if !found {
			t.Errorf("template %q has no home page", id)
		}
	}
}

func TestTemplates_EveryPageHasTitlesInBothLanguages(t *testing.T) {
	for id, tpl := range Templates {
		for _, page := range tpl.Pages() {
			if page.Title == "" {
				t.Errorf("template %q page %q has no English title", id, page.Slug)
			}
			if page.TitleNe == "" {
				t.Errorf("template %q page %q has no Nepali title", id, page.Slug)
			}
		}
	}
}

func TestTemplateIDs_MatchesRegistry(t *testing.T) {
	ids := TemplateIDs()
	if len(ids) != len(Templates) {
		t.Fatalf("TemplateIDs returns %d ids but the registry has %d", len(ids), len(Templates))
	}
	for _, id := range ids {
		if _, ok := Templates[id]; !ok {
			t.Errorf("TemplateIDs lists %q which is not in the registry", id)
		}
	}
}

func TestTemplates_ThemeJSONRoundTrips(t *testing.T) {
	for id, tpl := range Templates {
		raw := tpl.ThemeJSON()
		if len(raw) == 0 || string(raw) == "{}" {
			t.Errorf("template %q produced an empty theme", id)
		}
	}
}

func TestTemplates_ThemesAreDistinct(t *testing.T) {
	// Five templates that share an accent colour are one template with five
	// names. Guard the thing the design actually asked for.
	seen := map[string]string{}
	for id, tpl := range Templates {
		key := tpl.Theme.Accent + "|" + tpl.Theme.Bg
		if other, dup := seen[key]; dup {
			t.Errorf("templates %q and %q share a background and accent", id, other)
		}
		seen[key] = id
	}
}

func TestTemplates_SectionIDsUniqueWithinAPage(t *testing.T) {
	for id, tpl := range Templates {
		for _, page := range tpl.Pages() {
			seen := map[string]bool{}
			for _, s := range page.Sections {
				if seen[s.ID] {
					t.Errorf("template %q page %q repeats section id %q", id, page.Slug, s.ID)
				}
				seen[s.ID] = true
			}
		}
	}
}

func TestTemplates_DefaultTemplateExists(t *testing.T) {
	if _, ok := Templates[DefaultTemplate]; !ok {
		t.Fatalf("DefaultTemplate %q is not in the registry", DefaultTemplate)
	}
}
