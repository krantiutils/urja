package e2e

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/healthz", nil, "")
	assertStatus(t, resp, http.StatusOK)
}
