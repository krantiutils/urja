package e2e

import (
	"context"
	"net/http"
	"testing"
)

// The ledger invariant this whole feature rests on: it always equals the sum
// of non-cancelled documents. These tests assert the transactions rows
// directly, not just the HTTP response, because that invariant lives in what
// got written to the ledger, not in what the API echoed back.

func TestInvoiceLedger_FromScratchBillRecordsIncome(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000701", "Admin")
	orgID := createTestOrg(t, admin, "Ledger Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct {
		ID            string `json:"id"`
		TransactionID string `json:"transaction_id"`
	}
	parseJSON(t, resp, &inv)

	var count int
	var txnID string
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*), MIN(id::text) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'income'
		    AND category = 'Sales' AND amount = 1000`, orgID,
	).Scan(&count, &txnID)
	if err != nil {
		t.Fatalf("counting income rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("income rows = %d, want exactly 1", count)
	}
	if inv.TransactionID != txnID {
		t.Errorf("invoice transaction_id = %q, want it to point at the income row %q", inv.TransactionID, txnID)
	}

	var ownsTransaction bool
	if err := testPool.QueryRow(context.Background(),
		`SELECT owns_transaction FROM invoices WHERE id = $1`, inv.ID,
	).Scan(&ownsTransaction); err != nil {
		t.Fatalf("reading owns_transaction: %v", err)
	}
	if !ownsTransaction {
		t.Error("owns_transaction = false, want true for a bill raised from scratch")
	}
}

func TestInvoiceLedger_PackageLinkedBillRecordsNothingNew(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000702", "Admin")
	orgID := createTestOrg(t, admin, "Package Linked Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	txnID := insertTestTransaction(t, orgID, admin)

	var before int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions WHERE organization_id = $1`, orgID,
	).Scan(&before); err != nil {
		t.Fatalf("counting transactions before issue: %v", err)
	}

	body := issueBody("Ram", 1, 1000)
	body["transaction_id"] = txnID
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct {
		ID            string `json:"id"`
		TransactionID string `json:"transaction_id"`
	}
	parseJSON(t, resp, &inv)

	var after int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions WHERE organization_id = $1`, orgID,
	).Scan(&after); err != nil {
		t.Fatalf("counting transactions after issue: %v", err)
	}
	if after != before {
		t.Errorf("ledger row count = %d, want unchanged at %d — a package-linked bill must write nothing new", after, before)
	}
	if inv.TransactionID != txnID {
		t.Errorf("transaction_id = %q, want the supplied %q", inv.TransactionID, txnID)
	}

	var ownsTransaction bool
	if err := testPool.QueryRow(context.Background(),
		`SELECT owns_transaction FROM invoices WHERE id = $1`, inv.ID,
	).Scan(&ownsTransaction); err != nil {
		t.Fatalf("reading owns_transaction: %v", err)
	}
	if ownsTransaction {
		t.Error("owns_transaction = true, want false for a bill that only links a package sale's row")
	}
}

func TestInvoiceLedger_CancellingFromScratchBillNetsToZero(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000703", "Admin")
	orgID := createTestOrg(t, admin, "Cancel Nets Zero Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	var incomeCount, expenseCount int
	var incomeAmount, expenseAmount float64
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'income' AND category = 'Sales'`,
		orgID,
	).Scan(&incomeCount, &incomeAmount); err != nil {
		t.Fatalf("counting income rows: %v", err)
	}
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM transactions
		  WHERE organization_id = $1 AND transaction_type = 'expense' AND category = 'Sales reversal'`,
		orgID,
	).Scan(&expenseCount, &expenseAmount); err != nil {
		t.Fatalf("counting expense reversal rows: %v", err)
	}

	if incomeCount != 1 {
		t.Errorf("income/Sales rows = %d, want exactly 1", incomeCount)
	}
	if expenseCount != 1 {
		t.Errorf("expense/Sales reversal rows = %d, want exactly 1", expenseCount)
	}
	if incomeAmount != 1000 || expenseAmount != 1000 {
		t.Errorf("income %.2f / reversal %.2f, want both 1000", incomeAmount, expenseAmount)
	}
	if incomeAmount-expenseAmount != 0 {
		t.Errorf("net = %.2f, want 0", incomeAmount-expenseAmount)
	}
}

func TestInvoiceLedger_CancellingPackageLinkedBillLeavesLedgerUntouched(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000704", "Admin")
	orgID := createTestOrg(t, admin, "Cancel Linked Untouched Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	txnID := insertTestTransaction(t, orgID, admin)

	body := issueBody("Ram", 1, 1000)
	body["transaction_id"] = txnID
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	var beforeCount int
	var beforeTotal float64
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM transactions WHERE organization_id = $1`, orgID,
	).Scan(&beforeCount, &beforeTotal); err != nil {
		t.Fatalf("counting transactions before cancel: %v", err)
	}

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+inv.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	var afterCount int
	var afterTotal float64
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM transactions WHERE organization_id = $1`, orgID,
	).Scan(&afterCount, &afterTotal); err != nil {
		t.Fatalf("counting transactions after cancel: %v", err)
	}

	if afterCount != beforeCount {
		t.Errorf("ledger row count after cancel = %d, want unchanged at %d", afterCount, beforeCount)
	}
	if afterTotal != beforeTotal {
		t.Errorf("ledger total after cancel = %.2f, want unchanged at %.2f", afterTotal, beforeTotal)
	}
}

func TestInvoiceLedger_ZeroTotalBillWritesNoIncomeRow(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000705", "Admin")
	orgID := createTestOrg(t, admin, "Zero Total Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	// A free line item (unit_price 0) is a legitimate bill — e.g. a
	// complimentary item — but there is no money to record as income.
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 0), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	var count int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM transactions WHERE organization_id = $1`, orgID,
	).Scan(&count); err != nil {
		t.Fatalf("counting transactions: %v", err)
	}
	if count != 0 {
		t.Errorf("ledger rows for a zero-total bill = %d, want 0", count)
	}

	var ownsTransaction bool
	if err := testPool.QueryRow(context.Background(),
		`SELECT owns_transaction FROM invoices WHERE id = $1`, inv.ID,
	).Scan(&ownsTransaction); err != nil {
		t.Fatalf("reading owns_transaction: %v", err)
	}
	if ownsTransaction {
		t.Error("owns_transaction = true, want false — a zero-total bill created no income row to own")
	}
}

func TestInvoiceLedger_ImmutabilityTriggerProtectsOwnsTransaction(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000706", "Admin")
	orgID := createTestOrg(t, admin, "Owns Transaction Immutable Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var inv struct{ ID string }
	parseJSON(t, resp, &inv)

	_, err := testPool.Exec(context.Background(),
		`UPDATE invoices SET owns_transaction = NOT owns_transaction WHERE id = $1`, inv.ID)
	if err == nil {
		t.Fatal("expected the immutability trigger to reject an owns_transaction change, got nil")
	}
}

// The mirror image of the bug this task fixed: cancelling a from-scratch
// bill reverses its income but does not — cannot, the transaction_id column
// is immutable — release the transaction_id row it created.
// idx_invoices_one_per_transaction only blocks a second *issued* invoice from
// linking that row, so without this check a brand-new bill could quote the
// same transaction_id, come back with owns_transaction = false, and leave
// the ledger net at zero while the new bill counts as 1000 of live sales.
func TestInvoiceLedger_CannotRelinkAReversedTransaction(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000707", "Admin")
	orgID := createTestOrg(t, admin, "No Relink Reversed Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var a struct {
		ID            string `json:"id"`
		TransactionID string `json:"transaction_id"`
	}
	parseJSON(t, resp, &a)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+a.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	body := issueBody("Shyam", 1, 1000)
	body["transaction_id"] = a.TransactionID
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusBadRequest)

	var respBody struct{ Code string }
	parseJSON(t, resp, &respBody)
	if respBody.Code != "transaction_already_owned" {
		t.Errorf("code = %q, want transaction_already_owned", respBody.Code)
	}

	var count int
	if err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM invoices WHERE transaction_id = $1`, a.TransactionID,
	).Scan(&count); err != nil {
		t.Fatalf("counting invoices against the transaction: %v", err)
	}
	if count != 1 {
		t.Errorf("invoices linking the reversed transaction = %d, want exactly 1 (the cancelled original)", count)
	}
}

// The must-keep-working case: a package-flow transaction is legitimately
// reusable across cancel-and-reissue — it is never owned by any invoice, so
// the new ownership check must leave this path alone.
func TestInvoiceLedger_PackageFlowTransactionStaysRelinkableAfterCancel(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000708", "Admin")
	orgID := createTestOrg(t, admin, "Package Relink Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	txnID := insertTestTransaction(t, orgID, admin)

	body := issueBody("Ram", 1, 1000)
	body["transaction_id"] = txnID
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusCreated)
	var a struct{ ID string }
	parseJSON(t, resp, &a)

	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices/"+a.ID+"/cancel",
		map[string]any{"reason": "wrong customer"}, token)
	assertStatus(t, resp, http.StatusOK)

	body2 := issueBody("Shyam", 1, 1000)
	body2["transaction_id"] = txnID
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body2, token)
	assertStatus(t, resp, http.StatusCreated)

	var b struct {
		TransactionID string `json:"transaction_id"`
	}
	parseJSON(t, resp, &b)
	if b.TransactionID != txnID {
		t.Errorf("transaction_id = %q, want the reused %q", b.TransactionID, txnID)
	}

	var ownsTransaction bool
	if err := testPool.QueryRow(context.Background(),
		`SELECT owns_transaction FROM invoices WHERE transaction_id = $1 AND status = 'issued'`, txnID,
	).Scan(&ownsTransaction); err != nil {
		t.Fatalf("reading owns_transaction: %v", err)
	}
	if ownsTransaction {
		t.Error("owns_transaction = true, want false — this bill only linked a package sale's row")
	}
}

// Before this fix, quoting an owned transaction while its owning invoice was
// still issued was already rejected — but only because
// idx_invoices_one_per_transaction collided, surfacing as already_billed.
// verifyOrgReferences must now catch it earlier and return the clearer code.
func TestInvoiceLedger_OwnedTransactionCannotBeLinkedWhileIssued(t *testing.T) {
	cleanupTables(t)
	admin := createTestUser(t, "9800000709", "Admin")
	orgID := createTestOrg(t, admin, "Owned While Issued Gym")
	setPAN(t, orgID, "601234567")
	token := generateTestToken(admin, "member")

	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices",
		issueBody("Ram", 1, 1000), token)
	assertStatus(t, resp, http.StatusCreated)
	var a struct {
		TransactionID string `json:"transaction_id"`
	}
	parseJSON(t, resp, &a)

	body := issueBody("Shyam", 1, 1000)
	body["transaction_id"] = a.TransactionID
	resp = doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/invoices", body, token)
	assertStatus(t, resp, http.StatusBadRequest)

	var respBody struct{ Code string }
	parseJSON(t, resp, &respBody)
	if respBody.Code != "transaction_already_owned" {
		t.Errorf("code = %q, want transaction_already_owned (not already_billed — "+
			"the ownership check must catch this before the unique index does)", respBody.Code)
	}
}
