package e2e

import (
	"context"
	"net/http"
	"sync"
	"testing"
)

func TestInvoicePrint_FirstIsOriginalThenCopy(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000601", "Admin")
	orgID := createTestOrg(t, admin, "Print Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	var out struct {
		CopyLabel string `json:"copy_label"`
		Invoice   struct {
			PrintCount int `json:"print_count"`
		} `json:"invoice"`
	}

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/print", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &out)
	if out.CopyLabel != "original" {
		t.Errorf("first print label = %q, want original", out.CopyLabel)
	}
	if out.Invoice.PrintCount != 1 {
		t.Errorf("print_count = %d, want 1", out.Invoice.PrintCount)
	}

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/print", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &out)
	if out.CopyLabel != "copy" {
		t.Errorf("second print label = %q, want copy", out.CopyLabel)
	}
	if out.Invoice.PrintCount != 2 {
		t.Errorf("print_count = %d, want 2", out.Invoice.PrintCount)
	}
}

func TestInvoiceNextNumber_PreviewsWithoutConsuming(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000602", "Admin")
	orgID := createTestOrg(t, admin, "Preview Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	var preview struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/invoices/next-number", nil, token)
	assertStatus(t, resp, http.StatusOK)
	parseJSON(t, resp, &preview)

	// Previewing twice must return the same number — it reserves nothing.
	resp = doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/invoices/next-number", nil, token)
	var second struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	parseJSON(t, resp, &second)
	if preview.InvoiceNumber != second.InvoiceNumber {
		t.Errorf("preview consumed a number: %q then %q", preview.InvoiceNumber, second.InvoiceNumber)
	}

	// And the bill actually issued must take that number.
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	var inv struct {
		InvoiceNumber string `json:"invoice_number"`
	}
	parseJSON(t, resp, &inv)
	if inv.InvoiceNumber != preview.InvoiceNumber {
		t.Errorf("issued %q but preview promised %q", inv.InvoiceNumber, preview.InvoiceNumber)
	}
}

// Printing through another org's path must 404, not 403 — a 403 would
// confirm the invoice exists in a gym the caller has no business seeing —
// and unlike a read, a print has side effects: it must not log a print or
// bump the count on a document org A was never allowed to touch.
func TestInvoicePrint_CrossTenantDoesNotPrint(t *testing.T) {
	cleanupTables(t)
	adminA := createTestUser(t, "9800000603", "Admin A")
	adminB := createTestUser(t, "9800000604", "Admin B")
	orgA := createTestOrg(t, adminA, "Tenant A Print")
	orgB := createTestOrg(t, adminB, "Tenant B Print")
	setPAN(t, orgB, "601234567")
	tokenA := generateTestToken(adminA, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgB+"/invoices",
		issueBody("Ram", 1, 1000), generateTestToken(adminB, "member"))
	assertStatus(t, resp, http.StatusCreated)
	var invB struct{ ID string }
	parseJSON(t, resp, &invB)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgA+"/invoices/"+invB.ID+"/print", nil, tokenA)
	assertStatus(t, resp, http.StatusNotFound)

	var printCount int
	if err := testPool.QueryRow(context.Background(),
		`SELECT print_count FROM invoices WHERE id = $1`, invB.ID).Scan(&printCount); err != nil {
		t.Fatalf("reading print_count: %v", err)
	}
	if printCount != 0 {
		t.Errorf("print_count = %d, want 0 — org A must not have printed org B's invoice", printCount)
	}

	var logged int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM invoice_prints WHERE invoice_id = $1`, invB.ID).Scan(&logged); err != nil {
		t.Fatalf("counting invoice_prints: %v", err)
	}
	if logged != 0 {
		t.Errorf("invoice_prints rows = %d, want 0 — the cross-tenant attempt must leave no audit trail", logged)
	}
}

// Two staff hitting print at the same moment must not lose one of the two
// log rows, and print_count must land on 2 — not on 1 because one increment
// overwrote the other. A naive "read count, then write count+1" done without
// a lock is exactly the kind of race that produces a lost update; this test
// exists to catch it rather than trust the increment by inspection.
func TestInvoicePrint_ConcurrentPrintsNoLostUpdate(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000605", "Admin")
	orgID := createTestOrg(t, admin, "Concurrent Print Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	const n = 2
	var wg sync.WaitGroup
	statuses := make([]int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp := doRequest(t, http.MethodPost,
				"/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/print", nil, token)
			defer resp.Body.Close()
			statuses[i] = resp.StatusCode
		}(i)
	}
	wg.Wait()

	for _, s := range statuses {
		if s != http.StatusOK {
			t.Fatalf("a concurrent print returned %d, want 200", s)
		}
	}

	var printCount int
	if err := testPool.QueryRow(context.Background(),
		`SELECT print_count FROM invoices WHERE id = $1`, inv.ID).Scan(&printCount); err != nil {
		t.Fatalf("reading print_count: %v", err)
	}
	if printCount != n {
		t.Errorf("print_count = %d, want %d — the increment needs a lock", printCount, n)
	}

	var logged int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM invoice_prints WHERE invoice_id = $1`, inv.ID).Scan(&logged); err != nil {
		t.Fatalf("counting invoice_prints: %v", err)
	}
	if logged != n {
		t.Errorf("invoice_prints rows = %d, want %d — a concurrent print was lost", logged, n)
	}
}
