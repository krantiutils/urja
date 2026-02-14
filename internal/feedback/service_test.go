package feedback

import (
	"context"
	"testing"
)

func TestCreate_EmptyMessage(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "org-1", "user-1", "")
	if err == nil {
		t.Fatal("expected error for empty message")
	}
	if err.Error() != "message is required" {
		t.Errorf("error = %q, want %q", err.Error(), "message is required")
	}
}

func TestCreate_WhitespaceOnlyMessage(t *testing.T) {
	s := NewService(nil, testLogger())

	_, err := s.Create(context.Background(), "org-1", "user-1", "   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only message")
	}
}
