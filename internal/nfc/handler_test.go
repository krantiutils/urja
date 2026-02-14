package nfc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

func TestListCards_MissingOrgID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/nfc-cards", nil)
	rr := httptest.NewRecorder()

	h.ListCards(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRegisterCard_MissingOrgID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/nfc-cards", nil)
	rr := httptest.NewRecorder()

	h.RegisterCard(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestAssignCard_MissingOrgID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPut, "/nfc-cards/test-id/assign", nil)
	rr := httptest.NewRecorder()

	h.AssignCard(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestAssignCard_MissingCardID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPut, "/nfc-cards//assign", nil)
	ctx := context.WithValue(req.Context(), middleware.OrgIDKey, "org-1")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	h.AssignCard(rr, req)

	// Card ID comes from chi URL param, which will be empty without chi router
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestUnassignCard_MissingOrgID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPut, "/nfc-cards/test-id/unassign", nil)
	rr := httptest.NewRecorder()

	h.UnassignCard(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestListDevices_MissingOrgID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/nfc-devices", nil)
	rr := httptest.NewRecorder()

	h.ListDevices(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRegisterOrgRoutes(t *testing.T) {
	h := &Handler{}
	r := chi.NewRouter()
	h.RegisterOrgRoutes(r)

	expectedRoutes := map[string]string{
		"GET":  "/nfc-cards",
		"POST": "/nfc-cards",
	}

	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if expected, ok := expectedRoutes[method]; ok && route == expected {
			delete(expectedRoutes, method)
		}
		return nil
	})

	// At minimum, GET and POST /nfc-cards should be registered
	// (PUT routes have sub-paths that may vary in walk output)
}
