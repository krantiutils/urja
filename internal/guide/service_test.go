package guide

import (
	"context"
	"testing"
)

func TestCreate_EmptyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "", "", "content", "", "strength", "", "user-1")
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if err.Error() != "title is required" {
		t.Errorf("error = %q, want %q", err.Error(), "title is required")
	}
}

func TestCreate_EmptyContent(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "Title", "", "", "", "strength", "", "user-1")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
	if err.Error() != "content is required" {
		t.Errorf("error = %q, want %q", err.Error(), "content is required")
	}
}

func TestCreate_WhitespaceOnlyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "   ", "", "content", "", "", "", "user-1")
	if err == nil {
		t.Fatal("expected error for whitespace-only title")
	}
}

func TestCreate_InvalidCategory(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "Title", "", "content", "", "invalid_cat", "", "user-1")
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
	if err.Error() != "invalid category: invalid_cat" {
		t.Errorf("error = %q, want %q", err.Error(), "invalid category: invalid_cat")
	}
}

func TestUpdate_EmptyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Update(context.Background(), "guide-1", "", "", "content", "", "strength", "")
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestUpdate_EmptyContent(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Update(context.Background(), "guide-1", "Title", "", "", "", "strength", "")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestUpdate_InvalidCategory(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Update(context.Background(), "guide-1", "Title", "", "content", "", "bogus", "")
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestValidCategories(t *testing.T) {
	expected := []string{"strength", "cardio", "flexibility", "nutrition", "recovery", "beginner"}
	for _, cat := range expected {
		if !validCategories[cat] {
			t.Errorf("expected %q to be a valid category", cat)
		}
	}
}
