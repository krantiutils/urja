package e2e

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestOrgSettings_SetPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000201", "Admin")
	orgID := createTestOrg(t, admin, "PAN Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "601234567",
		"tax_legal_name": "PAN Gym Pvt Ltd",
		"tax_address":    "Kirtipur, Kathmandu",
	}, token)
	assertStatus(t, resp, http.StatusOK)

	var got struct {
		PANNumber    string `json:"pan_number"`
		TaxLegalName string `json:"tax_legal_name"`
	}
	parseJSON(t, resp, &got)
	if got.PANNumber != "601234567" {
		t.Errorf("pan_number = %q, want %q", got.PANNumber, "601234567")
	}
	if got.TaxLegalName != "PAN Gym Pvt Ltd" {
		t.Errorf("tax_legal_name = %q, want %q", got.TaxLegalName, "PAN Gym Pvt Ltd")
	}
}

func TestOrgSettings_RejectsMalformedPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000202", "Admin")
	orgID := createTestOrg(t, admin, "Bad PAN Gym")
	token := generateTestToken(admin, "member")

	for _, bad := range []string{"12345", "abcdefghi", "6012345678"} {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
			map[string]any{"pan_number": bad}, token)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("PAN %q: status = %d, want 400", bad, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

// TestOrgSettings_ClearsBlankFields is fix-round-1 finding 2: an empty string
// for any of the three tax fields must clear it (write SQL NULL), not
// silently no-op and not 500. pan_number in particular has a CHECK requiring
// NULL or exactly 9 digits, so writing the literal empty string there is the
// specific way this used to fail.
func TestOrgSettings_ClearsBlankFields(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000203", "Admin")
	orgID := createTestOrg(t, admin, "Clearable Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "601234567",
		"tax_legal_name": "Clearable Gym Pvt Ltd",
		"tax_address":    "Baneshwor, Kathmandu",
	}, token)
	assertStatus(t, resp, http.StatusOK)
	resp.Body.Close()

	resp = doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "",
		"tax_legal_name": "",
		"tax_address":    "",
	}, token)
	assertStatus(t, resp, http.StatusOK)

	var got struct {
		PANNumber    string `json:"pan_number"`
		TaxLegalName string `json:"tax_legal_name"`
		TaxAddress   string `json:"tax_address"`
	}
	parseJSON(t, resp, &got)
	if got.PANNumber != "" {
		t.Errorf("pan_number after clearing = %q, want empty", got.PANNumber)
	}
	if got.TaxLegalName != "" {
		t.Errorf("tax_legal_name after clearing = %q, want empty", got.TaxLegalName)
	}
	if got.TaxAddress != "" {
		t.Errorf("tax_address after clearing = %q, want empty", got.TaxAddress)
	}

	// Read back directly from the database (independent of what the handler
	// just returned) and confirm the columns are actually NULL, not the
	// literal empty string — the latter would violate pan_number's format
	// CHECK the moment anything else tried to write it.
	var panInDB, legalInDB, addrInDB *string
	err := testPool.QueryRow(context.Background(),
		`SELECT pan_number, tax_legal_name, tax_address FROM organizations WHERE id = $1`, orgID,
	).Scan(&panInDB, &legalInDB, &addrInDB)
	if err != nil {
		t.Fatalf("querying organizations: %v", err)
	}
	if panInDB != nil {
		t.Errorf("pan_number in DB = %q, want NULL", *panInDB)
	}
	if legalInDB != nil {
		t.Errorf("tax_legal_name in DB = %q, want NULL", *legalInDB)
	}
	if addrInDB != nil {
		t.Errorf("tax_address in DB = %q, want NULL", *addrInDB)
	}
}

// TestOrgSettings_RejectsOversizedTaxLegalName is fix-round-1 finding 3:
// tax_legal_name is VARCHAR(255); without a handler-side check an oversized
// value fails at the database and surfaces as an opaque 500.
func TestOrgSettings_RejectsOversizedTaxLegalName(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000204", "Admin")
	orgID := createTestOrg(t, admin, "Long Name Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID,
		map[string]any{"tax_legal_name": strings.Repeat("A", 256)}, token)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
	resp.Body.Close()
}

// TestOrgPublic_DoesNotLeakTaxIdentity is fix-round-1 finding 1: pan_number,
// tax_legal_name and tax_address must never reach the anonymous gym
// directory (GET /api/v1/gyms, GET /api/v1/gyms/{id} — no Authorization
// header). Asserted on the raw response body, not a decoded struct: a struct
// with the fields simply omitted would hide a leak instead of catching it.
func TestOrgPublic_DoesNotLeakTaxIdentity(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000205", "Admin")
	orgID := createTestOrg(t, admin, "Secretive Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "609876543",
		"tax_legal_name": "Secretive Gym Pvt Ltd",
		"tax_address":    "Ason, Kathmandu",
	}, token)
	assertStatus(t, resp, http.StatusOK)
	resp.Body.Close()

	secrets := []string{"609876543", "Secretive Gym Pvt Ltd", "Ason, Kathmandu"}

	listResp := doRequest(t, http.MethodGet, "/api/v1/gyms", nil, "")
	assertStatus(t, listResp, http.StatusOK)
	listBody := readBody(t, listResp)
	for _, s := range secrets {
		if strings.Contains(listBody, s) {
			t.Errorf("GET /api/v1/gyms leaked tax identity: body contains %q\nbody: %s", s, listBody)
		}
	}

	getResp := doRequest(t, http.MethodGet, "/api/v1/gyms/"+orgID, nil, "")
	assertStatus(t, getResp, http.StatusOK)
	getBody := readBody(t, getResp)
	for _, s := range secrets {
		if strings.Contains(getBody, s) {
			t.Errorf("GET /api/v1/gyms/{id} leaked tax identity: body contains %q\nbody: %s", s, getBody)
		}
	}
}

// TestOrgSettings_GetAuthenticated_RoundTrip is fix-round-2 finding 1: a saved
// PAN must actually be readable back by the org that saved it. Before this
// fix there was no authenticated GET at all — the settings page's initial
// fetch hit the unauthenticated GET /api/v1/gyms/{id} directory, which never
// carries pan_number/tax_legal_name/tax_address by design (see
// TestOrgPublic_DoesNotLeakTaxIdentity above), so a saved PAN would reload as
// blank forever even though the database had it.
func TestOrgSettings_GetAuthenticated_RoundTrip(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000206", "Admin")
	orgID := createTestOrg(t, admin, "Round Trip Gym")
	token := generateTestToken(admin, "member")

	putResp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "601111111",
		"tax_legal_name": "Round Trip Gym Pvt Ltd",
		"tax_address":    "Patan, Lalitpur",
	}, token)
	assertStatus(t, putResp, http.StatusOK)
	putResp.Body.Close()

	getResp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID, nil, token)
	assertStatus(t, getResp, http.StatusOK)

	var got struct {
		PANNumber    string `json:"pan_number"`
		TaxLegalName string `json:"tax_legal_name"`
		TaxAddress   string `json:"tax_address"`
	}
	parseJSON(t, getResp, &got)
	if got.PANNumber != "601111111" {
		t.Errorf("pan_number = %q, want %q", got.PANNumber, "601111111")
	}
	if got.TaxLegalName != "Round Trip Gym Pvt Ltd" {
		t.Errorf("tax_legal_name = %q, want %q", got.TaxLegalName, "Round Trip Gym Pvt Ltd")
	}
	if got.TaxAddress != "Patan, Lalitpur" {
		t.Errorf("tax_address = %q, want %q", got.TaxAddress, "Patan, Lalitpur")
	}
}

// TestOrgSettings_GetAuthenticated_CrossTenant404 is fix-round-2 finding 1's
// other half: an admin of one gym must not be able to read another gym's tax
// identity by guessing its org ID, and the response must not even confirm the
// org exists — 404, not the 403 that OrgScope's shared "not a member" check
// would otherwise give away (which is why this route is deliberately not
// gated by OrgScope; see cmd/api/main.go).
func TestOrgSettings_GetAuthenticated_CrossTenant404(t *testing.T) {
	cleanupTables(t)
	ownerA := createTestUser(t, "9800000207", "Owner A")
	orgA := createTestOrg(t, ownerA, "Org A")

	ownerB := createTestUser(t, "9800000208", "Owner B")
	createTestOrg(t, ownerB, "Org B") // ownerB legitimately admins a different org
	tokenB := generateTestToken(ownerB, "member")

	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgA, nil, tokenB)
	assertStatus(t, resp, http.StatusNotFound)
}

// TestOrgSettings_GetAuthenticated_SuperAdminNoMembership is fix-round-3
// finding 1: a genuine platform super_admin has no organization_members row
// for a gym they didn't found — nothing seeds one, not Repository.Create, not
// any trigger — so this is the normal case for a super_admin browsing a gym's
// settings, not an edge case. GetAuthenticated must check IsSuperAdmin
// independently of GetOrgRole; checking it only as a fallback reached after a
// membership row already exists (the order Update uses, safe there only
// because OrgScope requires membership before Update ever runs) would 404 a
// real super_admin on the one route that was deliberately moved outside
// OrgScope so a super_admin could reach it.
func TestOrgSettings_GetAuthenticated_SuperAdminNoMembership(t *testing.T) {
	cleanupTables(t)
	owner := createTestUser(t, "9800000209", "Owner")
	orgID := createTestOrg(t, owner, "Platform-Viewed Gym")
	ownerToken := generateTestToken(owner, "member")

	putResp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID, map[string]any{
		"pan_number":     "602222222",
		"tax_legal_name": "Platform-Viewed Gym Pvt Ltd",
		"tax_address":    "Boudha, Kathmandu",
	}, ownerToken)
	assertStatus(t, putResp, http.StatusOK)
	putResp.Body.Close()

	// No createTestOrg/createTestOrgMember call for this user against orgID —
	// deliberately no organization_members row at all, matching a real
	// super_admin who never joined this gym.
	superAdmin := createTestSuperAdmin(t, "9800000210", "Platform Admin")
	superToken := generateTestToken(superAdmin, "member")

	getResp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID, nil, superToken)
	assertStatus(t, getResp, http.StatusOK)

	var got struct {
		PANNumber string `json:"pan_number"`
	}
	parseJSON(t, getResp, &got)
	if got.PANNumber != "602222222" {
		t.Errorf("pan_number = %q, want %q", got.PANNumber, "602222222")
	}
}
