package site

import (
	"encoding/json"
	"fmt"
)

// DefaultTemplate is applied to a gym that has never picked one.
const DefaultTemplate = "fight_club"

// ThemeTokens are emitted as CSS custom properties on the site root. Keeping the
// palette here rather than in Tailwind config is what lets five templates look
// genuinely different without five copies of every renderer.
type ThemeTokens struct {
	Bg          string `json:"bg"`
	Surface     string `json:"surface"`
	Fg          string `json:"fg"`
	FgMuted     string `json:"fg_muted"`
	Accent      string `json:"accent"`
	AccentFg    string `json:"accent_fg"`
	Border      string `json:"border"`
	Radius      string `json:"radius"`
	FontDisplay string `json:"font_display"`
	FontBody    string `json:"font_body"`
	// DisplayTransform lets a template shout (uppercase) or stay quiet (none)
	// without a renderer knowing which template it is rendering.
	DisplayTransform string `json:"display_transform"`
	DisplayWeight    string `json:"display_weight"`
	DisplayTracking  string `json:"display_tracking"`
}

// Template is a named starting point: a theme plus a set of seeded pages.
type Template struct {
	ID     string
	Name   string
	NameNe string
	Theme  ThemeTokens
	build  func() []Page
}

// Pages returns freshly built seed pages. Built per call so a caller mutating
// the result cannot corrupt the template registry.
func (t Template) Pages() []Page { return t.build() }

// ThemeJSON encodes the theme for storage in site_settings.theme.
func (t Template) ThemeJSON() json.RawMessage {
	b, err := json.Marshal(t.Theme)
	if err != nil {
		// The struct is fixed and always marshals; an empty object is a safe floor.
		return json.RawMessage(`{}`)
	}
	return b
}

// --- section construction helpers ---

// sec builds a section, marshalling content and assigning a stable ID. Seed
// content is authored in code, so a marshal failure is a programming error and
// the templates test catches it.
func sec(id, typ, variant string, content map[string]interface{}, style SectionStyle) Section {
	raw, err := json.Marshal(content)
	if err != nil {
		panic(fmt.Sprintf("template section %s: %v", id, err))
	}
	return Section{ID: id, Type: typ, Variant: variant, Content: raw, Style: style, Version: 1}
}

func style(bg, padding, width, align string) SectionStyle {
	return SectionStyle{Background: bg, Padding: padding, Width: width, Align: align}
}

// m is shorthand for a content map.
type m = map[string]interface{}

// --- shared seed content ---
//
// Real boxing-gym copy in both languages, not lorem ipsum: a gym admin should be
// able to publish immediately and edit later, rather than face a wall of
// placeholder text.

func seedPrograms() m {
	return m{
		"title": "What we train", "titleNe": "हामी के तालिम दिन्छौं",
		"subtitle":   "From your first jab to your first fight.",
		"subtitleNe": "तपाईंको पहिलो ज्याबदेखि पहिलो लडाइँसम्म।",
		"items": []m{
			{"icon": "target", "title": "Beginner Boxing", "titleNe": "सुरुवाती बक्सिङ",
				"description":   "Stance, footwork and the basic punches. No experience needed.",
				"descriptionNe": "स्ट्यान्स, फुटवर्क र आधारभूत प्रहारहरू। अनुभव आवश्यक छैन।"},
			{"icon": "zap", "title": "Pad Work & Technique", "titleNe": "प्याड वर्क र प्रविधि",
				"description":   "One-to-one rounds on the pads to sharpen combinations.",
				"descriptionNe": "कम्बिनेसन तिखार्न प्याडमा एक-एक राउन्ड।"},
			{"icon": "users", "title": "Sparring", "titleNe": "स्पारिङ",
				"description":   "Controlled, supervised sparring once your coach clears you.",
				"descriptionNe": "प्रशिक्षकले अनुमति दिएपछि नियन्त्रित, निगरानीमा स्पारिङ।"},
			{"icon": "flame", "title": "Boxing Fitness", "titleNe": "बक्सिङ फिटनेस",
				"description":   "Bag rounds, circuits and conditioning. Fight training, no contact.",
				"descriptionNe": "ब्याग राउन्ड, सर्किट र कन्डिसनिङ। सम्पर्करहित तालिम।"},
			{"icon": "shield", "title": "Kids Boxing", "titleNe": "बाल बक्सिङ",
				"description":   "Ages 8-14. Discipline, coordination and confidence.",
				"descriptionNe": "८-१४ वर्ष। अनुशासन, समन्वय र आत्मविश्वास।"},
			{"icon": "award", "title": "Fight Team", "titleNe": "फाइट टिम",
				"description":   "Competition training for amateur and national-level bouts.",
				"descriptionNe": "एमेच्योर र राष्ट्रिय स्तरका लडाइँका लागि प्रतिस्पर्धा तालिम।"},
		},
	}
}

func seedTimetable() m {
	return m{
		"title": "Weekly timetable", "titleNe": "साप्ताहिक तालिका",
		"note": "Walk in 10 minutes early for your first session.", "noteNe": "पहिलो सत्रका लागि १० मिनेट अगाडि आउनुहोस्।",
		"days": []m{
			{"day": "Sunday", "dayNe": "आइतबार", "classes": []m{
				{"time": "06:00", "name": "Boxing Fitness", "nameNe": "बक्सिङ फिटनेस", "coach": "Suman", "level": "All"},
				{"time": "17:30", "name": "Beginner Boxing", "nameNe": "सुरुवाती बक्सिङ", "coach": "Rajan", "level": "Beginner"},
				{"time": "19:00", "name": "Fight Team", "nameNe": "फाइट टिम", "coach": "Rajan", "level": "Advanced"},
			}},
			{"day": "Monday", "dayNe": "सोमबार", "classes": []m{
				{"time": "06:00", "name": "Bag Work", "nameNe": "ब्याग वर्क", "coach": "Suman", "level": "All"},
				{"time": "16:00", "name": "Kids Boxing", "nameNe": "बाल बक्सिङ", "coach": "Anita", "level": "Kids"},
				{"time": "18:00", "name": "Technique", "nameNe": "प्रविधि", "coach": "Rajan", "level": "Intermediate"},
			}},
			{"day": "Tuesday", "dayNe": "मंगलबार", "classes": []m{
				{"time": "06:00", "name": "Boxing Fitness", "nameNe": "बक्सिङ फिटनेस", "coach": "Suman", "level": "All"},
				{"time": "18:00", "name": "Beginner Boxing", "nameNe": "सुरुवाती बक्सिङ", "coach": "Anita", "level": "Beginner"},
				{"time": "19:30", "name": "Sparring", "nameNe": "स्पारिङ", "coach": "Rajan", "level": "Cleared only"},
			}},
			{"day": "Wednesday", "dayNe": "बुधबार", "classes": []m{
				{"time": "06:00", "name": "Bag Work", "nameNe": "ब्याग वर्क", "coach": "Suman", "level": "All"},
				{"time": "16:00", "name": "Kids Boxing", "nameNe": "बाल बक्सिङ", "coach": "Anita", "level": "Kids"},
				{"time": "18:00", "name": "Technique", "nameNe": "प्रविधि", "coach": "Rajan", "level": "Intermediate"},
			}},
			{"day": "Thursday", "dayNe": "बिहिबार", "classes": []m{
				{"time": "06:00", "name": "Boxing Fitness", "nameNe": "बक्सिङ फिटनेस", "coach": "Suman", "level": "All"},
				{"time": "18:00", "name": "Beginner Boxing", "nameNe": "सुरुवाती बक्सिङ", "coach": "Anita", "level": "Beginner"},
				{"time": "19:30", "name": "Sparring", "nameNe": "स्पारिङ", "coach": "Rajan", "level": "Cleared only"},
			}},
			{"day": "Friday", "dayNe": "शुक्रबार", "classes": []m{
				{"time": "06:00", "name": "Open Gym", "nameNe": "खुला जिम", "coach": "—", "level": "All"},
				{"time": "18:00", "name": "Fight Team", "nameNe": "फाइट टिम", "coach": "Rajan", "level": "Advanced"},
			}},
			{"day": "Saturday", "dayNe": "शनिबार", "classes": []m{
				{"time": "07:00", "name": "Conditioning", "nameNe": "कन्डिसनिङ", "coach": "Suman", "level": "All"},
				{"time": "10:00", "name": "Open Mat & Pads", "nameNe": "खुला म्याट र प्याड", "coach": "All coaches", "level": "All"},
			}},
		},
	}
}

func seedCoaches() m {
	return m{
		"title": "Your coaches", "titleNe": "तपाईंका प्रशिक्षकहरू",
		"coaches": []m{
			{"name": "Rajan Tamang", "nameNe": "राजन तामाङ",
				"role": "Head Coach", "roleNe": "प्रमुख प्रशिक्षक",
				"bio":    "National-level amateur, coaching in Kirtipur since 2015. Specialises in southpaw counter-work.",
				"bioNe":  "राष्ट्रिय स्तरका एमेच्योर, २०१५ देखि कीर्तिपुरमा प्रशिक्षण। साउथपा काउन्टरमा विशेषज्ञ।",
				"record": "18-4-1", "photo": "",
				"credentials": []string{"NBA Level 2 Coach", "National Championship Silver 2014"}},
			{"name": "Anita Gurung", "nameNe": "अनिता गुरुङ",
				"role": "Coach — Beginners & Kids", "roleNe": "प्रशिक्षक — सुरुवाती र बालबालिका",
				"bio":    "Builds first-timers from zero. Runs the kids programme and the women's beginner block.",
				"bioNe":  "शून्यबाट सुरु गर्नेहरूलाई तयार पार्छिन्। बाल कार्यक्रम र महिला सुरुवाती ब्लक चलाउँछिन्।",
				"record": "9-2-0", "photo": "",
				"credentials": []string{"Level 1 Coach", "First Aid certified"}},
			{"name": "Suman Shrestha", "nameNe": "सुमन श्रेष्ठ",
				"role": "Strength & Conditioning", "roleNe": "स्ट्रेन्थ र कन्डिसनिङ",
				"bio":    "Handles the morning conditioning block and fight-camp physical prep.",
				"bioNe":  "बिहानको कन्डिसनिङ ब्लक र फाइट क्याम्पको शारीरिक तयारी सम्हाल्छन्।",
				"record": "", "photo": "",
				"credentials": []string{"CSCS", "Sports nutrition diploma"}},
		},
	}
}

func seedFAQ() m {
	return m{
		"title": "Before you come", "titleNe": "आउनु अघि",
		"items": []m{
			{"question": "I have never boxed. Can I join?", "questionNe": "मैले कहिल्यै बक्सिङ गरेको छैन। सामेल हुन सक्छु?",
				"answer":   "Yes — most people who walk in have never trained. Start in Beginner Boxing; nobody will put you in with a fighter.",
				"answerNe": "सक्नुहुन्छ — आउने अधिकांशले कहिल्यै तालिम लिएका हुँदैनन्। सुरुवाती बक्सिङबाट सुरु गर्नुहोस्।"},
			{"question": "What should I bring?", "questionNe": "के ल्याउनु पर्छ?",
				"answer":   "Training clothes, a water bottle and a towel. We lend gloves and wraps for your first two weeks.",
				"answerNe": "तालिमका लुगा, पानीको बोतल र तौलिया। पहिलो दुई हप्ता हामी दस्ताना र र्‍याप दिन्छौं।"},
			{"question": "When can I start sparring?", "questionNe": "कहिलेदेखि स्पारिङ गर्न सक्छु?",
				"answer":   "When your coach clears you — usually after three to four months of consistent technique work. It is never rushed.",
				"answerNe": "प्रशिक्षकले अनुमति दिएपछि — सामान्यतया तीनदेखि चार महिनाको निरन्तर प्रविधि अभ्यासपछि।"},
			{"question": "Do you train women?", "questionNe": "के तपाईं महिलालाई तालिम दिनुहुन्छ?",
				"answer":   "Yes. Anita runs a women's beginner block, and every class is open to everyone.",
				"answerNe": "हो। अनिताले महिला सुरुवाती ब्लक चलाउँछिन्, र हरेक कक्षा सबैका लागि खुला छ।"},
			{"question": "Is there a trial session?", "questionNe": "के परीक्षण सत्र छ?",
				"answer":   "Your first session is free. Fill in the form below or just turn up ten minutes before a beginner class.",
				"answerNe": "तपाईंको पहिलो सत्र निःशुल्क छ। तलको फारम भर्नुहोस् वा सुरुवाती कक्षाभन्दा दस मिनेट अगाडि आउनुहोस्।"},
		},
	}
}

func seedLeadForm() m {
	return m{
		"title": "Book a free trial", "titleNe": "निःशुल्क परीक्षण बुक गर्नुहोस्",
		"subtitle":    "Leave your number and we will call you back the same day.",
		"subtitleNe":  "आफ्नो नम्बर छोड्नुहोस्, हामी सोही दिन फोन गर्नेछौं।",
		"submitLabel": "Request my trial", "submitLabelNe": "मेरो परीक्षण अनुरोध गर्नुहोस्",
		"successMessage":   "Thanks — we will call you shortly.",
		"successMessageNe": "धन्यवाद — हामी चाँडै फोन गर्नेछौं।",
		"interests": []m{
			{"value": "beginner", "label": "Beginner Boxing", "labelNe": "सुरुवाती बक्सिङ"},
			{"value": "fitness", "label": "Boxing Fitness", "labelNe": "बक्सिङ फिटनेस"},
			{"value": "kids", "label": "Kids Boxing", "labelNe": "बाल बक्सिङ"},
			{"value": "fight_team", "label": "Fight Team", "labelNe": "फाइट टिम"},
		},
	}
}

func seedContact() m {
	return m{
		"title": "Find us", "titleNe": "हामीलाई भेट्नुहोस्",
		"address": "Kirtipur, Kathmandu", "addressNe": "कीर्तिपुर, काठमाडौं",
		"phone": "", "email": "",
		"hoursNote":   "Open six days a week. Closed Saturday evening.",
		"hoursNoteNe": "हप्ताको छ दिन खुला। शनिबार साँझ बन्द।",
	}
}

func seedStats() m {
	return m{
		"items": []m{
			{"value": "10+", "label": "Years coaching", "labelNe": "वर्षको प्रशिक्षण"},
			{"value": "300+", "label": "Members trained", "labelNe": "तालिम प्राप्त सदस्य"},
			{"value": "24", "label": "Amateur bouts won", "labelNe": "एमेच्योर लडाइँ जितिएको"},
			{"value": "6", "label": "Days a week", "labelNe": "हप्ताको दिन"},
		},
	}
}

func seedHours() m {
	return m{
		"title": "Opening hours", "titleNe": "खुल्ने समय",
		"days": []m{
			{"day": "Sunday – Thursday", "dayNe": "आइतबार – बिहिबार", "hours": "06:00 – 08:00, 16:00 – 21:00"},
			{"day": "Friday", "dayNe": "शुक्रबार", "hours": "06:00 – 08:00, 18:00 – 20:00"},
			{"day": "Saturday", "dayNe": "शनिबार", "hours": "07:00 – 12:00"},
		},
	}
}

func seedTestimonials() m {
	return m{
		"title": "From the gym floor", "titleNe": "जिमबाट",
		"items": []m{
			{"name": "Bibek K.", "text": "Walked in unable to skip for thirty seconds. Six months later I had my first bout.",
				"textNe": "तीस सेकेन्ड पनि स्किप गर्न नसक्ने भएर आएँ। छ महिनापछि मेरो पहिलो लडाइँ भयो।", "rating": 5},
			{"name": "Sarita M.", "text": "The coaches actually correct you. I learned more in a month here than a year on YouTube.",
				"textNe": "प्रशिक्षकहरूले साँच्चै सच्याउँछन्। युट्युबमा एक वर्षभन्दा यहाँ एक महिनामा बढी सिकें।", "rating": 5},
			{"name": "Prakash T.", "text": "Hard training, no ego. They will not let you spar until you are genuinely ready.",
				"textNe": "कडा तालिम, अहंकार छैन। साँच्चै तयार नभएसम्म स्पारिङ गर्न दिँदैनन्।", "rating": 5},
		},
	}
}

func seedGallery() m {
	return m{
		"title": "Inside the gym", "titleNe": "जिम भित्र",
		"images":      []m{},
		"emptyHint":   "Add photos of your gym floor, bags, ring and team.",
		"emptyHintNe": "आफ्नो जिम, ब्याग, रिङ र टिमका तस्बिरहरू थप्नुहोस्।",
	}
}

// --- Templates ---

// Templates is the registry of available starting points.
var Templates = map[string]Template{
	"fight_club":  fightClub,
	"iron_sweat":  ironSweat,
	"champion":    champion,
	"community":   community,
	"minimal_pro": minimalPro,
}

// TemplateIDs returns the registry keys in a stable, presentable order.
func TemplateIDs() []string {
	return []string{"fight_club", "iron_sweat", "champion", "community", "minimal_pro"}
}

var fightClub = Template{
	ID: "fight_club", Name: "Fight Club", NameNe: "फाइट क्लब",
	Theme: ThemeTokens{
		Bg: "#08080a", Surface: "#121216", Fg: "#f5f5f7", FgMuted: "#9a9aa4",
		Accent: "#dc2626", AccentFg: "#ffffff", Border: "rgba(255,255,255,0.08)",
		Radius: "2px", FontDisplay: `"Oswald", "Arial Narrow", sans-serif`,
		FontBody:         `"Inter", system-ui, sans-serif`,
		DisplayTransform: "uppercase", DisplayWeight: "700", DisplayTracking: "0.02em",
	},
	build: func() []Page {
		return []Page{
			{Slug: "home", Title: "Home", TitleNe: "गृहपृष्ठ", SortOrder: 0, ShowInNav: true, IsPublished: true,
				Sections: []Section{
					sec("h1", "hero", "fullbleed", m{
						"title": "Train like you mean it", "titleNe": "साँच्चै तालिम गर्नुहोस्",
						"subtitle":   "Boxing in Kirtipur. Beginners welcome, fighters made.",
						"subtitleNe": "कीर्तिपुरमा बक्सिङ। सुरुवातीहरूलाई स्वागत, लडाकुहरू तयार।",
						"image":      "",
						"buttons": []m{
							{"label": "Book a free trial", "labelNe": "निःशुल्क परीक्षण", "href": "/contact", "style": "solid"},
							{"label": "See the timetable", "labelNe": "तालिका हेर्नुहोस्", "href": "/classes", "style": "outline"},
						},
					}, style("base", "lg", "full", "center")),
					sec("s1", "stats_bar", "bordered", seedStats(), style("surface", "md", "contained", "center")),
					sec("p1", "programs_grid", "cards", seedPrograms(), style("base", "lg", "contained", "left")),
					sec("t1", "class_timetable", "day_tabs", seedTimetable(), style("surface", "lg", "contained", "left")),
					sec("c1", "coaches", "cards", seedCoaches(), style("base", "lg", "contained", "left")),
					sec("g1", "gallery", "masonry", seedGallery(), style("base", "lg", "contained", "left")),
					sec("r1", "testimonials", "cards", seedTestimonials(), style("surface", "lg", "contained", "left")),
					sec("cta1", "cta_banner", "image", m{
						"title": "First session is free", "titleNe": "पहिलो सत्र निःशुल्क",
						"subtitle":   "Turn up ten minutes early. Bring water.",
						"subtitleNe": "दस मिनेट अगाडि आउनुहोस्। पानी ल्याउनुहोस्।",
						"buttons":    []m{{"label": "Get started", "labelNe": "सुरु गर्नुहोस्", "href": "/contact", "style": "solid"}},
					}, style("accent", "lg", "full", "center")),
				}},
			classesPage(), coachesPage(), contactPage(),
		}
	},
}

var ironSweat = Template{
	ID: "iron_sweat", Name: "Iron & Sweat", NameNe: "आइरन एन्ड स्वेट",
	Theme: ThemeTokens{
		Bg: "#f4f4f0", Surface: "#e6e6e0", Fg: "#16161a", FgMuted: "#55555f",
		Accent: "#facc15", AccentFg: "#16161a", Border: "#16161a",
		Radius: "0px", FontDisplay: `"JetBrains Mono", ui-monospace, monospace`,
		FontBody:         `"Inter", system-ui, sans-serif`,
		DisplayTransform: "uppercase", DisplayWeight: "800", DisplayTracking: "-0.02em",
	},
	build: func() []Page {
		return []Page{
			{Slug: "home", Title: "Home", TitleNe: "गृहपृष्ठ", SortOrder: 0, ShowInNav: true, IsPublished: true,
				Sections: []Section{
					sec("h1", "hero", "split", m{
						"title": "No shortcuts. Just rounds.", "titleNe": "कुनै सर्टकट छैन। केवल राउन्ड।",
						"subtitle":   "A working boxing gym in Kirtipur since 2015.",
						"subtitleNe": "२०१५ देखि कीर्तिपुरमा एक कार्यरत बक्सिङ जिम।",
						"image":      "",
						"buttons":    []m{{"label": "Start training", "labelNe": "तालिम सुरु गर्नुहोस्", "href": "/contact", "style": "solid"}},
					}, style("base", "lg", "contained", "left")),
					sec("l1", "logo_strip", "simple", m{
						"title": "Affiliations", "titleNe": "सम्बद्धता", "logos": []m{},
					}, style("surface", "sm", "contained", "center")),
					sec("p1", "programs_grid", "numbered", seedPrograms(), style("base", "lg", "contained", "left")),
					sec("t1", "class_timetable", "table", seedTimetable(), style("surface", "lg", "contained", "left")),
					sec("c1", "coaches", "list", seedCoaches(), style("base", "lg", "contained", "left")),
					sec("f1", "fight_record", "timeline", m{
						"title": "Gym record", "titleNe": "जिम रेकर्ड",
						"bouts":       []m{},
						"emptyHint":   "Add your team's bouts as they happen.",
						"emptyHintNe": "टिमका लडाइँहरू भएसँगै थप्नुहोस्।",
					}, style("surface", "lg", "contained", "left")),
					sec("q1", "faq", "two_column", seedFAQ(), style("base", "lg", "contained", "left")),
					sec("lf1", "lead_form", "split", seedLeadForm(), style("accent", "lg", "contained", "left")),
				}},
			classesPage(), coachesPage(), contactPage(),
		}
	},
}

var champion = Template{
	ID: "champion", Name: "Champion", NameNe: "च्याम्पियन",
	Theme: ThemeTokens{
		Bg: "#fffdf8", Surface: "#f6f1e7", Fg: "#1c1917", FgMuted: "#78716c",
		Accent: "#b45309", AccentFg: "#ffffff", Border: "rgba(28,25,23,0.12)",
		Radius: "10px", FontDisplay: `"Playfair Display", Georgia, serif`,
		FontBody:         `"Inter", system-ui, sans-serif`,
		DisplayTransform: "none", DisplayWeight: "600", DisplayTracking: "-0.01em",
	},
	build: func() []Page {
		return []Page{
			{Slug: "home", Title: "Home", TitleNe: "गृहपृष्ठ", SortOrder: 0, ShowInNav: true, IsPublished: true,
				Sections: []Section{
					sec("h1", "hero", "centered", m{
						"title": "The sweet science, taught properly", "titleNe": "मीठो विज्ञान, राम्ररी सिकाइएको",
						"subtitle":   "Technique first. Everything else follows.",
						"subtitleNe": "पहिले प्रविधि। बाँकी सबै त्यसपछि।",
						"image":      "",
						"buttons":    []m{{"label": "Visit the gym", "labelNe": "जिम भ्रमण गर्नुहोस्", "href": "/contact", "style": "solid"}},
					}, style("base", "lg", "narrow", "center")),
					sec("rt1", "rich_text", "centered", m{
						"body":   "We are a small boxing gym in Kirtipur. We teach the fundamentals slowly and properly, and we do not put anyone in the ring before they are ready.",
						"bodyNe": "हामी कीर्तिपुरको एउटा सानो बक्सिङ जिम हौं। हामी आधारभूत कुरा बिस्तारै र राम्ररी सिकाउँछौं, र कोही तयार नभएसम्म रिङमा पठाउँदैनौं।",
					}, style("base", "md", "narrow", "center")),
					sec("c1", "coaches", "spotlight", seedCoaches(), style("surface", "lg", "contained", "left")),
					sec("g1", "gallery", "grid", seedGallery(), style("base", "lg", "contained", "center")),
					sec("mp1", "membership_plans", "cards", m{
						"title": "Membership", "titleNe": "सदस्यता", "dataSource": "auto",
						"note":   "Ask at the desk about student and family rates.",
						"noteNe": "विद्यार्थी र पारिवारिक दरका लागि डेस्कमा सोध्नुहोस्।",
					}, style("surface", "lg", "contained", "center")),
					sec("r1", "testimonials", "quote", seedTestimonials(), style("base", "lg", "narrow", "center")),
					sec("cta1", "cta_banner", "gradient", m{
						"title": "Come and watch a session", "titleNe": "एउटा सत्र हेर्न आउनुहोस्",
						"subtitle": "No pressure, no sign-up.", "subtitleNe": "कुनै दबाब छैन, दर्ता छैन।",
						"buttons": []m{{"label": "Contact us", "labelNe": "सम्पर्क गर्नुहोस्", "href": "/contact", "style": "solid"}},
					}, style("accent", "lg", "full", "center")),
				}},
			classesPage(), coachesPage(), contactPage(),
		}
	},
}

var community = Template{
	ID: "community", Name: "Community", NameNe: "सामुदायिक",
	Theme: ThemeTokens{
		Bg: "#fffbf7", Surface: "#fff1e6", Fg: "#2d1b12", FgMuted: "#7c6155",
		Accent: "#ea580c", AccentFg: "#ffffff", Border: "rgba(45,27,18,0.10)",
		Radius: "18px", FontDisplay: `"Inter", system-ui, sans-serif`,
		FontBody:         `"Inter", system-ui, sans-serif`,
		DisplayTransform: "none", DisplayWeight: "700", DisplayTracking: "-0.02em",
	},
	build: func() []Page {
		return []Page{
			{Slug: "home", Title: "Home", TitleNe: "गृहपृष्ठ", SortOrder: 0, ShowInNav: true, IsPublished: true,
				Sections: []Section{
					sec("h1", "hero", "split", m{
						"title": "Everyone starts somewhere", "titleNe": "सबै कतैबाट सुरु गर्छन्",
						"subtitle":   "A friendly boxing gym for Kirtipur — kids, beginners and fighters.",
						"subtitleNe": "कीर्तिपुरको मैत्रीपूर्ण बक्सिङ जिम — बालबालिका, सुरुवाती र लडाकुहरू।",
						"image":      "",
						"buttons": []m{
							{"label": "Try a free class", "labelNe": "निःशुल्क कक्षा", "href": "/contact", "style": "pill"},
						},
					}, style("base", "lg", "contained", "left")),
					sec("s1", "stats_bar", "cards", seedStats(), style("base", "md", "contained", "center")),
					sec("p1", "programs_grid", "icons", seedPrograms(), style("surface", "lg", "contained", "center")),
					sec("oh1", "opening_hours", "compact", seedHours(), style("base", "md", "contained", "left")),
					sec("t1", "class_timetable", "cards", seedTimetable(), style("base", "lg", "contained", "left")),
					sec("c1", "coaches", "cards", seedCoaches(), style("surface", "lg", "contained", "center")),
					sec("q1", "faq", "accordion", seedFAQ(), style("base", "lg", "narrow", "left")),
					sec("lf1", "lead_form", "card", seedLeadForm(), style("surface", "lg", "narrow", "center")),
					sec("mp1", "map_embed", "with_info", seedContact(), style("base", "none", "full", "left")),
				}},
			classesPage(), coachesPage(), contactPage(),
		}
	},
}

var minimalPro = Template{
	ID: "minimal_pro", Name: "Minimal Pro", NameNe: "मिनिमल प्रो",
	Theme: ThemeTokens{
		Bg: "#ffffff", Surface: "#fafafa", Fg: "#0a0a0a", FgMuted: "#737373",
		Accent: "#0a0a0a", AccentFg: "#ffffff", Border: "rgba(10,10,10,0.10)",
		Radius: "4px", FontDisplay: `"Inter", system-ui, sans-serif`,
		FontBody:         `"Inter", system-ui, sans-serif`,
		DisplayTransform: "none", DisplayWeight: "500", DisplayTracking: "-0.03em",
	},
	build: func() []Page {
		return []Page{
			{Slug: "home", Title: "Home", TitleNe: "गृहपृष्ठ", SortOrder: 0, ShowInNav: true, IsPublished: true,
				Sections: []Section{
					sec("h1", "hero", "minimal", m{
						"title": "Boxing. Kirtipur.", "titleNe": "बक्सिङ। कीर्तिपुर।",
						"subtitle":   "Six days a week, from six in the morning.",
						"subtitleNe": "हप्ताको छ दिन, बिहान छ बजेदेखि।",
						"buttons":    []m{{"label": "Get in touch", "labelNe": "सम्पर्क", "href": "/contact", "style": "outline"}},
					}, style("base", "lg", "narrow", "left")),
					sec("d1", "divider", "line", m{}, style("base", "none", "narrow", "center")),
					sec("p1", "programs_grid", "list", seedPrograms(), style("base", "lg", "narrow", "left")),
					sec("t1", "class_timetable", "table", seedTimetable(), style("base", "lg", "contained", "left")),
					sec("c1", "coaches", "list", seedCoaches(), style("surface", "lg", "narrow", "left")),
					sec("mp1", "membership_plans", "table", m{
						"title": "Membership", "titleNe": "सदस्यता", "dataSource": "auto",
					}, style("base", "lg", "narrow", "left")),
					sec("q1", "faq", "list", seedFAQ(), style("base", "lg", "narrow", "left")),
					sec("ci1", "contact_info", "two_column", seedContact(), style("surface", "lg", "narrow", "left")),
				}},
			classesPage(), coachesPage(), contactPage(),
		}
	},
}

// --- shared inner pages ---
//
// The four seeded pages are the same set for every template; only the home page
// composition and the theme differ. This keeps a gym's navigation predictable
// regardless of which look they pick.

func classesPage() Page {
	return Page{Slug: "classes", Title: "Classes", TitleNe: "कक्षाहरू", SortOrder: 1, ShowInNav: true, IsPublished: true,
		Sections: []Section{
			sec("ch1", "hero", "minimal", m{
				"title": "Classes & timetable", "titleNe": "कक्षा र तालिका",
				"subtitle":   "Drop into any class marked for your level.",
				"subtitleNe": "आफ्नो स्तरका जुनसुकै कक्षामा आउन सक्नुहुन्छ।",
			}, style("base", "md", "contained", "left")),
			sec("cp1", "programs_grid", "cards", seedPrograms(), style("base", "lg", "contained", "left")),
			sec("ct1", "class_timetable", "table", seedTimetable(), style("surface", "lg", "contained", "left")),
			sec("cq1", "faq", "accordion", seedFAQ(), style("base", "lg", "narrow", "left")),
			sec("clf1", "lead_form", "inline", seedLeadForm(), style("accent", "lg", "contained", "center")),
		}}
}

func coachesPage() Page {
	return Page{Slug: "coaches", Title: "Coaches", TitleNe: "प्रशिक्षकहरू", SortOrder: 2, ShowInNav: true, IsPublished: true,
		Sections: []Section{
			sec("kh1", "hero", "minimal", m{
				"title": "The people who will train you", "titleNe": "तपाईंलाई तालिम दिने व्यक्तिहरू",
			}, style("base", "md", "contained", "left")),
			sec("kc1", "coaches", "spotlight", seedCoaches(), style("base", "lg", "contained", "left")),
			sec("kf1", "fight_record", "table", m{
				"title": "Competition record", "titleNe": "प्रतिस्पर्धा रेकर्ड",
				"bouts":       []m{},
				"emptyHint":   "Add bouts to show your team's competition history.",
				"emptyHintNe": "टिमको प्रतिस्पर्धा इतिहास देखाउन लडाइँहरू थप्नुहोस्।",
			}, style("surface", "lg", "contained", "left")),
		}}
}

func contactPage() Page {
	return Page{Slug: "contact", Title: "Contact", TitleNe: "सम्पर्क", SortOrder: 3, ShowInNav: true, IsPublished: true,
		Sections: []Section{
			sec("nh1", "hero", "minimal", m{
				"title": "Come and train", "titleNe": "आएर तालिम गर्नुहोस्",
				"subtitle":   "Your first session is free.",
				"subtitleNe": "तपाईंको पहिलो सत्र निःशुल्क छ।",
			}, style("base", "md", "contained", "left")),
			sec("nlf1", "lead_form", "card", seedLeadForm(), style("base", "lg", "narrow", "left")),
			sec("nci1", "contact_info", "two_column", seedContact(), style("surface", "lg", "contained", "left")),
			sec("noh1", "opening_hours", "table", seedHours(), style("base", "md", "contained", "left")),
			sec("nm1", "map_embed", "full_width", seedContact(), style("base", "none", "full", "left")),
		}}
}
