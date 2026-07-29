package invoice

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_Issue_RejectsEmptyItems(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		CustomerName: "Ram",
		Items:        nil,
	})
	if !errors.Is(err, ErrNoItems) {
		t.Fatalf("error = %v, want ErrNoItems", err)
	}
}

func TestService_Issue_RejectsMissingCustomerName(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		Items: []ItemInput{{Description: "Boxing", Quantity: 1, UnitPrice: 1000}},
	})
	if err == nil {
		t.Fatal("expected an error for a bill with no customer name, got nil")
	}
}

func TestService_Issue_RejectsBadLineValues(t *testing.T) {
	tests := []struct {
		name string
		item ItemInput
	}{
		{"zero quantity", ItemInput{Description: "Boxing", Quantity: 0, UnitPrice: 1000}},
		{"negative quantity", ItemInput{Description: "Boxing", Quantity: -1, UnitPrice: 1000}},
		{"negative price", ItemInput{Description: "Boxing", Quantity: 1, UnitPrice: -5}},
		{"blank description", ItemInput{Description: "  ", Quantity: 1, UnitPrice: 1000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{logger: testLogger()}
			_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
				CustomerName: "Ram",
				Items:        []ItemInput{tt.item},
			})
			if err == nil {
				t.Fatalf("expected an error for %s, got nil", tt.name)
			}
		})
	}
}

func TestService_Issue_RejectsDiscountAboveSubtotal(t *testing.T) {
	s := &Service{logger: testLogger()}
	_, err := s.Issue(context.Background(), "org-1", "user-1", IssueInput{
		CustomerName: "Ram",
		Discount:     2000,
		Items:        []ItemInput{{Description: "Boxing", Quantity: 1, UnitPrice: 1000}},
	})
	if !errors.Is(err, ErrInvalidDiscount) {
		t.Fatalf("error = %v, want ErrInvalidDiscount", err)
	}
}
