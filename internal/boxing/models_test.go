package boxing

import "testing"

func TestWeightClassFor_Divisions(t *testing.T) {
	cases := []struct {
		kg   float64
		want string
	}{
		// Well inside each division.
		{46, "light_flyweight"},
		{50, "flyweight"},
		{53, "bantamweight"},
		{56, "featherweight"},
		{59, "lightweight"},
		{62, "light_welterweight"},
		{65, "welterweight"},
		{70, "light_middleweight"},
		{74, "middleweight"},
		{79, "light_heavyweight"},
		{85, "cruiserweight"},
		{90, "heavyweight"},
		{120, SuperHeavyweight},
	}
	for _, tc := range cases {
		if got := WeightClassFor(tc.kg); got != tc.want {
			t.Errorf("WeightClassFor(%.1f) = %q, want %q", tc.kg, got, tc.want)
		}
	}
}

func TestWeightClassFor_ExactBoundaries(t *testing.T) {
	// The limit belongs to the lighter division — a boxer who makes 60.0kg
	// exactly is a lightweight, not a light welterweight. Getting this wrong
	// puts someone in the wrong division on fight night.
	cases := []struct {
		kg   float64
		want string
	}{
		{48, "light_flyweight"},
		{48.1, "flyweight"},
		{51, "flyweight"},
		{51.1, "bantamweight"},
		{60, "lightweight"},
		{60.1, "light_welterweight"},
		{63.5, "light_welterweight"},
		{63.6, "welterweight"},
		{92, "heavyweight"},
		{92.1, SuperHeavyweight},
	}
	for _, tc := range cases {
		if got := WeightClassFor(tc.kg); got != tc.want {
			t.Errorf("WeightClassFor(%.2f) = %q, want %q", tc.kg, got, tc.want)
		}
	}
}

func TestWeightClassFor_NonPositiveWeight(t *testing.T) {
	// A member who has never logged a weight must not be silently filed as a
	// light flyweight.
	for _, kg := range []float64{0, -1, -70} {
		if got := WeightClassFor(kg); got != "" {
			t.Errorf("WeightClassFor(%.1f) = %q, want empty", kg, got)
		}
	}
}

func TestTallyBouts(t *testing.T) {
	bouts := []Bout{
		{Result: "win"}, {Result: "win"}, {Result: "win"},
		{Result: "loss"},
		{Result: "draw"},
		{Result: "no_contest"},
	}
	rec := TallyBouts(bouts)

	if rec.Wins != 3 || rec.Losses != 1 || rec.Draws != 1 || rec.NoContests != 1 {
		t.Fatalf("unexpected tally: %+v", rec)
	}
	if rec.Total != 6 {
		t.Errorf("Total = %d, want 6", rec.Total)
	}
	// No-contests are counted but conventionally excluded from the W-L-D string.
	if rec.String() != "3-1-1" {
		t.Errorf("String() = %q, want %q", rec.String(), "3-1-1")
	}
}

func TestTallyBouts_Empty(t *testing.T) {
	rec := TallyBouts(nil)
	if rec.Total != 0 || rec.String() != "0-0-0" {
		t.Errorf("empty tally = %+v / %q", rec, rec.String())
	}
}

func TestValidators_AllowEmptyForOptionalFields(t *testing.T) {
	// Stance, skill level and method are optional: a member who has not been
	// assessed yet must still be saveable.
	if err := ValidateStance(""); err != nil {
		t.Errorf("empty stance should be allowed: %v", err)
	}
	if err := ValidateSkillLevel(""); err != nil {
		t.Errorf("empty skill level should be allowed: %v", err)
	}
	if err := ValidateMethod(""); err != nil {
		t.Errorf("empty method should be allowed: %v", err)
	}
	// Result is required, so empty must fail.
	if err := ValidateResult(""); err == nil {
		t.Error("empty bout result must be rejected")
	}
}

func TestValidators_RejectUnknownValues(t *testing.T) {
	if err := ValidateStance("sideways"); err == nil {
		t.Error("unknown stance must be rejected")
	}
	if err := ValidateSkillLevel("legendary"); err == nil {
		t.Error("unknown skill level must be rejected")
	}
	if err := ValidateResult("sort_of_won"); err == nil {
		t.Error("unknown result must be rejected")
	}
	if err := ValidateMethod("vibes"); err == nil {
		t.Error("unknown method must be rejected")
	}
}

func TestValidators_AcceptKnownValues(t *testing.T) {
	for _, s := range []string{"orthodox", "southpaw", "switch"} {
		if err := ValidateStance(s); err != nil {
			t.Errorf("stance %q should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"beginner", "intermediate", "amateur", "pro"} {
		if err := ValidateSkillLevel(s); err != nil {
			t.Errorf("skill level %q should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"win", "loss", "draw", "no_contest"} {
		if err := ValidateResult(s); err != nil {
			t.Errorf("result %q should be valid: %v", s, err)
		}
	}
}
