package activitylog

import (
	"testing"
	"time"
)

func TestGroupByDate_Empty(t *testing.T) {
	groups := groupByDate(nil)
	if len(groups) != 0 {
		t.Errorf("groupByDate(nil) returned %d groups, want 0", len(groups))
	}

	groups = groupByDate([]Entry{})
	if len(groups) != 0 {
		t.Errorf("groupByDate([]) returned %d groups, want 0", len(groups))
	}
}

func TestGroupByDate_SingleDay(t *testing.T) {
	base := time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC)
	entries := []Entry{
		{ID: "1", EventType: "check_in", Description: "Member checked in", CreatedAt: base},
		{ID: "2", EventType: "check_in", Description: "Another check-in", CreatedAt: base.Add(1 * time.Hour)},
		{ID: "3", EventType: "membership_approved", Description: "Membership approved", CreatedAt: base.Add(2 * time.Hour)},
	}

	groups := groupByDate(entries)
	if len(groups) != 1 {
		t.Fatalf("groupByDate() returned %d groups, want 1", len(groups))
	}
	if groups[0].Date != "2026-02-14" {
		t.Errorf("Date = %q, want %q", groups[0].Date, "2026-02-14")
	}
	if len(groups[0].Entries) != 3 {
		t.Errorf("entries count = %d, want 3", len(groups[0].Entries))
	}
}

func TestGroupByDate_MultipleDays(t *testing.T) {
	day1 := time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 2, 13, 15, 0, 0, 0, time.UTC)
	day3 := time.Date(2026, 2, 12, 8, 0, 0, 0, time.UTC)

	// Entries should be in descending order (newest first)
	entries := []Entry{
		{ID: "1", EventType: "check_in", CreatedAt: day1},
		{ID: "2", EventType: "check_in", CreatedAt: day1.Add(-1 * time.Hour)},
		{ID: "3", EventType: "package_changed", CreatedAt: day2},
		{ID: "4", EventType: "member_joined", CreatedAt: day3},
		{ID: "5", EventType: "check_in", CreatedAt: day3.Add(-2 * time.Hour)},
	}

	groups := groupByDate(entries)
	if len(groups) != 3 {
		t.Fatalf("groupByDate() returned %d groups, want 3", len(groups))
	}

	if groups[0].Date != "2026-02-14" {
		t.Errorf("groups[0].Date = %q, want %q", groups[0].Date, "2026-02-14")
	}
	if len(groups[0].Entries) != 2 {
		t.Errorf("groups[0] entries = %d, want 2", len(groups[0].Entries))
	}

	if groups[1].Date != "2026-02-13" {
		t.Errorf("groups[1].Date = %q, want %q", groups[1].Date, "2026-02-13")
	}
	if len(groups[1].Entries) != 1 {
		t.Errorf("groups[1] entries = %d, want 1", len(groups[1].Entries))
	}

	if groups[2].Date != "2026-02-12" {
		t.Errorf("groups[2].Date = %q, want %q", groups[2].Date, "2026-02-12")
	}
	if len(groups[2].Entries) != 2 {
		t.Errorf("groups[2] entries = %d, want 2", len(groups[2].Entries))
	}
}

func TestEventTypeConstants(t *testing.T) {
	// Verify event type constants match the CHECK constraint in the migration
	validTypes := map[EventType]bool{
		EventCheckIn:              true,
		EventSubscriptionApproved: true,
		EventMembershipRequest:    true,
		EventMembershipApproved:   true,
		EventPackageChanged:       true,
		EventNFCCardAssigned:      true,
		EventNFCCardUnassigned:    true,
		EventNFCCardRegistered:    true,
		EventMemberJoined:         true,
		EventMemberSuspended:      true,
	}

	expectedStrings := []string{
		"check_in",
		"subscription_approved",
		"membership_request",
		"membership_approved",
		"package_changed",
		"nfc_card_assigned",
		"nfc_card_unassigned",
		"nfc_card_registered",
		"member_joined",
		"member_suspended",
	}

	for _, s := range expectedStrings {
		if !validTypes[EventType(s)] {
			t.Errorf("event type %q not found in constants", s)
		}
	}

	if len(validTypes) != len(expectedStrings) {
		t.Errorf("constant count = %d, expected strings count = %d", len(validTypes), len(expectedStrings))
	}
}
