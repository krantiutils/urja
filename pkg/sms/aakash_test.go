package sms

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// The response shapes below are the ones Aakash actually returns. `error` is a
// boolean; it was decoded as a string, so every send — including the
// successful ones — failed to parse and was reported as an error. No SMS this
// product sent was ever recorded as delivered.
func TestSend_AcceptsRealResponseShapes(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr bool
	}{
		{
			name:   "success — error is boolean false",
			status: http.StatusOK,
			body:   `{"error":false,"message":"Message sent successfully","data":{"valid":1}}`,
		},
		{
			name:    "rejection — error is boolean true",
			status:  http.StatusOK,
			body:    `{"error":true,"message":"Invalid auth token"}`,
			wantErr: true,
		},
		{
			name:   "message arrives as an object rather than a string",
			status: http.StatusOK,
			body:   `{"error":false,"message":{"detail":"queued"}}`,
		},
		{
			name:   "unparseable body on a 200 is treated as sent",
			status: http.StatusOK,
			body:   `not json at all`,
		},
		{
			name:    "non-200 is an error regardless of body",
			status:  http.StatusUnauthorized,
			body:    `{"error":true,"message":"unauthorized"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := NewClient("token", srv.URL, testLogger())
			err := c.SendOTP(context.Background(), "9812345678", "123456")

			if tt.wantErr && err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected success, got %v", err)
			}
		})
	}
}

// A 200 that cannot be parsed must not be reported as a failure: the message
// was accepted, and telling the caller it failed makes them request a second
// code while the first is already in flight.
func TestSend_UnparseableSuccessDoesNotDoubleSend(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html>maintenance</html>`))
	}))
	defer srv.Close()

	c := NewClient("token", srv.URL, testLogger())
	if err := c.SendOTP(context.Background(), "9812345678", "123456"); err != nil {
		t.Fatalf("a 200 should be treated as sent, got %v", err)
	}
	if calls != 1 {
		t.Errorf("sent %d times, want exactly 1", calls)
	}
}

func TestSend_PostsTheFieldsAakashExpects(t *testing.T) {
	var gotToken, gotTo, gotText string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotToken = r.PostFormValue("auth_token")
		gotTo = r.PostFormValue("to")
		gotText = r.PostFormValue("text")
		_, _ = w.Write([]byte(`{"error":false,"message":"ok"}`))
	}))
	defer srv.Close()

	c := NewClient("secret-token", srv.URL, testLogger())
	if err := c.SendOTP(context.Background(), "9812345678", "654321"); err != nil {
		t.Fatalf("send: %v", err)
	}

	if gotToken != "secret-token" {
		t.Errorf("auth_token = %q", gotToken)
	}
	if gotTo != "9812345678" {
		t.Errorf("to = %q", gotTo)
	}
	// The code must reach the recipient; a template change that drops it would
	// send a message nobody can act on.
	if gotText == "" || !contains(gotText, "654321") {
		t.Errorf("text %q does not carry the code", gotText)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle ||
			len(needle) == 0 ||
			indexOf(haystack, needle) >= 0)
}

func indexOf(h, n string) int {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}
