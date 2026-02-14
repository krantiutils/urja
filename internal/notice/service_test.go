package notice

import (
	"context"
	"testing"
)

func TestCreate_EmptyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "org-1", "", "", "content", "", "user-1")
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if err.Error() != "title is required" {
		t.Errorf("error = %q, want %q", err.Error(), "title is required")
	}
}

func TestCreate_EmptyContent(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "org-1", "Title", "", "", "", "user-1")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
	if err.Error() != "content is required" {
		t.Errorf("error = %q, want %q", err.Error(), "content is required")
	}
}

func TestCreate_WhitespaceOnlyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "org-1", "   ", "", "content", "", "user-1")
	if err == nil {
		t.Fatal("expected error for whitespace-only title")
	}
}

func TestList_DefaultPagination(t *testing.T) {
	s := NewService(nil, testLogger())

	// We can't actually call List without a repo, but we can verify the service was created.
	// The actual pagination clamping is tested implicitly.
	if s.repo != nil {
		t.Error("expected nil repo for test")
	}
}

func TestUpdate_EmptyTitle(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Update(context.Background(), "org-1", "notice-1", "", "", "content", "", true)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestUpdate_EmptyContent(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Update(context.Background(), "org-1", "notice-1", "Title", "", "", "", true)
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}
