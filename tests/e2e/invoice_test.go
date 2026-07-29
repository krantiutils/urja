package e2e

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// setPAN gives an org a PAN so it can issue bills.
func setPAN(t *testing.T, orgID, pan string) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`UPDATE organizations SET pan_number = $2 WHERE id = $1`, orgID, pan)
	if err != nil {
		t.Fatalf("setPAN: %v", err)
	}
}

func issueBody(customer string, qty, price float64) map[string]any {
	return map[string]any{
		"customer_name":  customer,
		"payment_method": "cash",
		"items": []map[string]any{
			{"description": "Monthly Boxing", "quantity": qty, "unit_price": price},
		},
	}
}

func TestInvoice_IssueRequiresPAN(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000301", "Admin")
	orgID := createTestOrg(t, admin, "No PAN Gym")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram Bahadur", 1, 3000), token)
	assertStatus(t, resp, http.StatusBadRequest)

	var body struct{ Code string }
	parseJSON(t, resp, &body)
	if body.Code != "pan_not_configured" {
		t.Errorf("code = %q, want %q", body.Code, "pan_not_configured")
	}
}

func TestInvoice_IssueComputesTotalsAndNumber(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000302", "Admin")
	orgID := createTestOrg(t, admin, "Billing Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	body := issueBody("Ram Bahadur", 2, 1500)
	body["discount"] = 500

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusCreated)

	var inv struct {
		InvoiceNumber string  `json:"invoice_number"`
		Sequence      int     `json:"sequence"`
		Subtotal      float64 `json:"subtotal"`
		Discount      float64 `json:"discount"`
		Total         float64 `json:"total"`
		VATAmount     float64 `json:"vat_amount"`
		AmountInWords string  `json:"amount_in_words"`
		SellerPAN     string  `json:"seller_pan"`
		Status        string  `json:"status"`
		Items         []struct {
			Amount float64 `json:"amount"`
		} `json:"items"`
	}
	parseJSON(t, resp, &inv)

	if inv.Sequence != 1 {
		t.Errorf("sequence = %d, want 1", inv.Sequence)
	}
	if inv.Subtotal != 3000 || inv.Discount != 500 || inv.Total != 2500 {
		t.Errorf("totals = subtotal %.2f discount %.2f total %.2f, want 3000/500/2500",
			inv.Subtotal, inv.Discount, inv.Total)
	}
	if inv.VATAmount != 0 {
		t.Errorf("vat_amount = %.2f, want 0 on a PAN-only bill", inv.VATAmount)
	}
	if inv.SellerPAN != "601234567" {
		t.Errorf("seller_pan = %q, want the org's PAN snapshotted", inv.SellerPAN)
	}
	if inv.Status != "issued" {
		t.Errorf("status = %q, want issued", inv.Status)
	}
	if len(inv.Items) != 1 || inv.Items[0].Amount != 3000 {
		t.Errorf("items = %+v, want one line of 3000", inv.Items)
	}
	if inv.AmountInWords == "" {
		t.Error("amount_in_words is empty")
	}
}

// The core guarantee: numbers must run 1..N with no gap and no duplicate even
// when several staff issue at the same moment.
func TestInvoice_ConcurrentIssuesAreGapless(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000303", "Admin")
	orgID := createTestOrg(t, admin, "Concurrent Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	const n = 12
	seqs := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
				issueBody(fmt.Sprintf("Customer %d", i), 1, 1000), token)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				return
			}
			var inv struct {
				Sequence int `json:"sequence"`
			}
			parseJSON(t, resp, &inv)
			seqs[i] = inv.Sequence
		}(i)
	}
	wg.Wait()

	seen := map[int]bool{}
	for _, s := range seqs {
		if s == 0 {
			t.Fatal("an issue failed; all 12 should succeed")
		}
		if seen[s] {
			t.Fatalf("sequence %d was allocated twice", s)
		}
		seen[s] = true
	}
	for want := 1; want <= n; want++ {
		if !seen[want] {
			t.Errorf("sequence %d is missing — the run has a gap", want)
		}
	}
}

func TestInvoice_CrossTenantIsInvisible(t *testing.T) {
	cleanupTables(t)
	adminA := createTestUser(t, "9800000304", "Admin A")
	adminB := createTestUser(t, "9800000305", "Admin B")
	orgA := createTestOrg(t, adminA, "Tenant A Gym")
	orgB := createTestOrg(t, adminB, "Tenant B Gym")
	setPAN(t, orgA, "601234567")
	tokenA := generateTestToken(adminA, "member")
	tokenB := generateTestToken(adminB, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgA+"/invoices",
		issueBody("Ram", 1, 1000), tokenA)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	// B asking for A's invoice through B's own org must 404.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgB+"/invoices/"+inv.ID, nil, tokenB)
	assertStatus(t, resp, http.StatusNotFound)

	// And B must not reach it through A's org either.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgA+"/invoices/"+inv.ID, nil, tokenB)
	if resp.StatusCode == http.StatusOK {
		t.Error("org B read org A's invoice through org A's path")
	}
	resp.Body.Close()
}

func TestInvoice_MemberCannotIssue(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000306", "Admin")
	memberID := createTestUser(t, "9800000307", "Member")
	orgID := createTestOrg(t, admin, "Gated Gym")
	createTestOrgMember(t, memberID, orgID, "member")
	setPAN(t, orgID, "601234567")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), generateTestToken(memberID, "member"))
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403 — billing is staff and admin only", resp.StatusCode)
	}
	resp.Body.Close()
}
