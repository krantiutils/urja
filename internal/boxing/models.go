package boxing

import (
	"fmt"
	"time"
)

// Stance options.
const (
	StanceOrthodox = "orthodox"
	StanceSouthpaw = "southpaw"
	StanceSwitch   = "switch"
)

var validStances = map[string]bool{
	StanceOrthodox: true, StanceSouthpaw: true, StanceSwitch: true,
}

var validSkillLevels = map[string]bool{
	"beginner": true, "intermediate": true, "amateur": true, "pro": true,
}

var validResults = map[string]bool{
	"win": true, "loss": true, "draw": true, "no_contest": true,
}

var validMethods = map[string]bool{
	"ko": true, "tko": true, "decision": true, "split_decision": true,
	"unanimous_decision": true, "rsc": true, "dq": true, "walkover": true,
}

// Profile is a member's combat-sports profile within one organization.
type Profile struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	OrgID             string     `json:"organization_id"`
	Stance            string     `json:"stance,omitempty"`
	WeightClass       string     `json:"weight_class,omitempty"`
	SkillLevel        string     `json:"skill_level,omitempty"`
	SparringCleared   bool       `json:"sparring_cleared"`
	SparringClearedAt *time.Time `json:"sparring_cleared_at,omitempty"`
	SparringClearedBy string     `json:"sparring_cleared_by,omitempty"`
	ReachCm           *float64   `json:"reach_cm,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Bout is one recorded competitive fight.
type Bout struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	OrgID       string    `json:"organization_id"`
	BoutDate    time.Time `json:"bout_date"`
	Opponent    string    `json:"opponent,omitempty"`
	EventName   string    `json:"event_name,omitempty"`
	Result      string    `json:"result"`
	Method      string    `json:"method,omitempty"`
	Rounds      *int      `json:"rounds,omitempty"`
	WeightClass string    `json:"weight_class,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Record is a win-loss-draw tally derived from a member's bouts.
type Record struct {
	Wins       int `json:"wins"`
	Losses     int `json:"losses"`
	Draws      int `json:"draws"`
	NoContests int `json:"no_contests"`
	Total      int `json:"total"`
}

// String renders the conventional "18-4-1" form.
func (r Record) String() string {
	return fmt.Sprintf("%d-%d-%d", r.Wins, r.Losses, r.Draws)
}

// TallyBouts computes a record from a bout list.
func TallyBouts(bouts []Bout) Record {
	var rec Record
	for _, b := range bouts {
		switch b.Result {
		case "win":
			rec.Wins++
		case "loss":
			rec.Losses++
		case "draw":
			rec.Draws++
		case "no_contest":
			rec.NoContests++
		}
	}
	rec.Total = len(bouts)
	return rec
}

// weightDivision is one amateur boxing division, keyed by its upper bound.
type weightDivision struct {
	maxKg float64
	name  string
}

// amateurDivisions are the standard men's amateur boxing weight classes, in
// ascending order. A boxer competes in the lightest division whose upper bound
// they do not exceed, so the comparison below is inclusive.
var amateurDivisions = []weightDivision{
	{48, "light_flyweight"},
	{51, "flyweight"},
	{54, "bantamweight"},
	{57, "featherweight"},
	{60, "lightweight"},
	{63.5, "light_welterweight"},
	{67, "welterweight"},
	{71, "light_middleweight"},
	{75, "middleweight"},
	{81, "light_heavyweight"},
	{86, "cruiserweight"},
	{92, "heavyweight"},
}

// SuperHeavyweight is the open-ended top division.
const SuperHeavyweight = "super_heavyweight"

// WeightClassFor maps a bodyweight in kilograms to an amateur boxing division.
// A non-positive weight yields an empty string rather than a bogus division.
func WeightClassFor(weightKg float64) string {
	if weightKg <= 0 {
		return ""
	}
	for _, d := range amateurDivisions {
		if weightKg <= d.maxKg {
			return d.name
		}
	}
	return SuperHeavyweight
}

// ValidateStance checks a stance value, allowing empty (unset).
func ValidateStance(s string) error {
	if s == "" || validStances[s] {
		return nil
	}
	return fmt.Errorf("invalid stance %q", s)
}

// ValidateSkillLevel checks a skill level, allowing empty (unset).
func ValidateSkillLevel(s string) error {
	if s == "" || validSkillLevels[s] {
		return nil
	}
	return fmt.Errorf("invalid skill level %q", s)
}

// ValidateResult checks a bout result. Unlike the others this is required.
func ValidateResult(s string) error {
	if validResults[s] {
		return nil
	}
	return fmt.Errorf("invalid bout result %q", s)
}

// ValidateMethod checks a bout method, allowing empty (unset).
func ValidateMethod(s string) error {
	if s == "" || validMethods[s] {
		return nil
	}
	return fmt.Errorf("invalid bout method %q", s)
}
