package e2e

// This file covers the authorization holes described in
// docs/superpowers/plans/2026-07-27-phase0-bug-fixes.md, Tasks 2, 3, 4, and the
// routes portion of Task 8:
//
//   - Task 2: org-scoped route groups (staff, packages, accounts, absentees,
//     feedbacks, workout-templates, sms, notices, dues) must gate on the
//     per-org role resolved by OrgScope (middleware.RequireOrgRole), never on
//     the global JWT role claim (middleware.RequireRole), which historically
//     came from an arbitrary organization_members row unrelated to the org
//     being accessed.
//   - Task 3: subscription package-assignment routes had no role gate at all.
//   - Task 4: member role updates must not allow privilege self-escalation.
//   - Task 8 (routes only): NFC management routes had no role gate at all.
//
// generateTestToken(userID, role) signs a JWT with an explicit "role" claim
// that does NOT come from the database — this lets these tests directly
// simulate the historical vulnerability, where the JWT's global role claim
// could reflect a user's role in a *different* org than the one being
// accessed (the arbitrary-row bug). Setting the claim to "admin" here
// represents the worst case that bug could produce.

import (
	"net/http"
	"testing"
)

// TestCrossOrgRoleEscalation is the core regression test for Task 2 and 3: a user
// who is a genuine admin of Org A and only a plain member of Org B must not be able
// to use Org-A-flavored admin powers against Org B, no matter what their JWT's
// global role claim says. Every case here must return 403. Before the fix (routes
// gated by middleware.RequireRole, which trusts the JWT claim) this test fails with
// 2xx/201 on most cases; after the fix (middleware.RequireOrgRole, which trusts only
// the per-org role OrgScope resolved from the database for Org B) it must all be 403.
func TestCrossOrgRoleEscalation(t *testing.T) {
	cleanupTables(t)

	// crossUser is a genuine admin of Org A, and only a plain "member" of Org B.
	crossUser := createTestUser(t, "9806000001", "Cross User")
	_ = createTestOrg(t, crossUser, "AuthZ Org A") // makes crossUser admin of Org A

	orgBAdmin := createTestUser(t, "9806000002", "Org B Admin")
	orgB := createTestOrg(t, orgBAdmin, "AuthZ Org B") // makes orgBAdmin the real admin of Org B
	createTestOrgMember(t, crossUser, orgB, "member")

	memberB := createTestUser(t, "9806000003", "Org B Target Member")
	createTestOrgMember(t, memberB, orgB, "member")

	pkgB := createTestPackage(t, orgB, "Monthly", 30, 1000)

	// crossToken's "role" claim says "admin" — simulating the historical bug where
	// the arbitrary-row COALESCE query could return crossUser's Org A admin role
	// even when the request targets Org B, where crossUser is only a member.
	crossToken := generateTestToken(crossUser, "admin")

	// NOTE: training-guides is intentionally NOT in this list. internal/guide/ is
	// out of scope for this change (see TestTrainingGuidesCrossOrgEscalation_KnownGap
	// below for why, and what remains open).
	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"staff list", http.MethodGet, "/api/v1/orgs/" + orgB + "/staff"},
		{"packages create", http.MethodPost, "/api/v1/orgs/" + orgB + "/packages"},
		{"accounts list", http.MethodGet, "/api/v1/orgs/" + orgB + "/accounts"},
		{"absentees list", http.MethodGet, "/api/v1/orgs/" + orgB + "/absentees"},
		{"feedbacks list", http.MethodGet, "/api/v1/orgs/" + orgB + "/feedbacks"},
		{"workout-templates create", http.MethodPost, "/api/v1/orgs/" + orgB + "/workout-templates"},
		{"sms balance", http.MethodGet, "/api/v1/orgs/" + orgB + "/sms/balance"},
		{"notices create", http.MethodPost, "/api/v1/orgs/" + orgB + "/notices"},
		{"dues list", http.MethodGet, "/api/v1/orgs/" + orgB + "/dues"},
		{"subscription assign", http.MethodPost, "/api/v1/orgs/" + orgB + "/members/" + memberB + "/packages/assign"},
		{"subscription renew", http.MethodPost, "/api/v1/orgs/" + orgB + "/members/" + memberB + "/packages/renew"},
		{"nfc device registration", http.MethodPost, "/api/v1/orgs/" + orgB + "/nfc-devices"},
		{"nfc card registration", http.MethodPost, "/api/v1/orgs/" + orgB + "/nfc-cards"},
		{"nfc card assignment", http.MethodPut, "/api/v1/orgs/" + orgB + "/nfc-cards/no-such-card/assign"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := doRequest(t, tc.method, tc.path, nil, crossToken)
			if resp.StatusCode != http.StatusForbidden {
				body := readBody(t, resp)
				t.Fatalf("expected 403, got %d — cross-org role escalation is open: %s", resp.StatusCode, body)
			}
		})
	}

	// Positive controls: a genuine admin of Org B must still be able to use these
	// same routes. Without this, a naive fix (e.g. blocking the routes entirely)
	// would also make the negative cases above pass while breaking the product.
	orgBAdminToken := generateTestToken(orgBAdmin, "admin")

	positiveCases := []struct {
		name   string
		method string
		path   string
		body   interface{}
		want   int
	}{
		{"staff list", http.MethodGet, "/api/v1/orgs/" + orgB + "/staff", nil, http.StatusOK},
		{"packages list", http.MethodGet, "/api/v1/orgs/" + orgB + "/packages", nil, http.StatusOK},
		{"accounts list", http.MethodGet, "/api/v1/orgs/" + orgB + "/accounts", nil, http.StatusOK},
		{"absentees list", http.MethodGet, "/api/v1/orgs/" + orgB + "/absentees", nil, http.StatusOK},
		{"feedbacks list", http.MethodGet, "/api/v1/orgs/" + orgB + "/feedbacks", nil, http.StatusOK},
		{"workout-templates create", http.MethodPost, "/api/v1/orgs/" + orgB + "/workout-templates",
			map[string]string{"name": "Admin Template"}, http.StatusCreated},
		{"sms balance", http.MethodGet, "/api/v1/orgs/" + orgB + "/sms/balance", nil, http.StatusOK},
		{"notices create", http.MethodPost, "/api/v1/orgs/" + orgB + "/notices",
			map[string]string{"title": "Admin Notice", "content": "Body"}, http.StatusCreated},
		{"dues list", http.MethodGet, "/api/v1/orgs/" + orgB + "/dues", nil, http.StatusOK},
		{"subscription assign", http.MethodPost, "/api/v1/orgs/" + orgB + "/members/" + memberB + "/packages/assign",
			map[string]interface{}{
				"package_id":     pkgB,
				"start_date":     "2026-03-01",
				"payment_method": "cash",
				"amount_paid":    1000.0,
			}, http.StatusCreated},
		// "subscription renew" is registered in the same RegisterMemberRoutes group,
		// behind the exact same RequireOrgRole("staff", "admin") middleware instance
		// as "payments" — proving the gate opens here proves it opens for renew too,
		// without needing to fabricate a member_package to renew.
		{"subscription payments (same gate as renew)", http.MethodGet,
			"/api/v1/orgs/" + orgB + "/members/" + memberB + "/payments", nil, http.StatusOK},
		{"nfc devices list", http.MethodGet, "/api/v1/orgs/" + orgB + "/nfc-devices", nil, http.StatusOK},
		{"nfc cards list", http.MethodGet, "/api/v1/orgs/" + orgB + "/nfc-cards", nil, http.StatusOK},
	}

	for _, tc := range positiveCases {
		t.Run("admin: "+tc.name, func(t *testing.T) {
			resp := doRequest(t, tc.method, tc.path, tc.body, orgBAdminToken)
			if resp.StatusCode != tc.want {
				body := readBody(t, resp)
				t.Fatalf("expected %d, got %d — a genuine admin should not be locked out: %s", tc.want, resp.StatusCode, body)
			}
		})
	}
}

// TestTrainingGuidesCrossOrgScope covers the training-guides half of the
// cross-org escalation fix. Guides now carry an organization_id (migration
// 000049) and RegisterOrgRoutes gates on the per-org role, so a gym's staff can
// neither reach another gym's guides nor mutate the platform-wide presets that
// carry a NULL organization_id.
func TestTrainingGuidesCrossOrgScope(t *testing.T) {
	cleanupTables(t)

	adminA := createTestUser(t, "9806300001", "Guide Admin A")
	orgA := createTestOrg(t, adminA, "Guide Gym A")
	tokenA := generateTestToken(adminA, "admin")

	adminB := createTestUser(t, "9806300002", "Guide Admin B")
	orgB := createTestOrg(t, adminB, "Guide Gym B")
	tokenB := generateTestToken(adminB, "admin")

	// Org A creates a guide.
	resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgA+"/training-guides",
		map[string]interface{}{
			"title": "Org A Guide", "content": "Private to org A", "category": "strength",
		}, tokenA)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("creating guide in org A: expected 201, got %d: %s", resp.StatusCode, readBody(t, resp))
	}

	var created map[string]interface{}
	parseJSON(t, resp, &created)
	guideID, _ := created["id"].(string)
	if guideID == "" {
		t.Fatal("created guide has no id")
	}

	// Org B must not be able to read, edit, publish or delete it.
	t.Run("org B cannot read org A's guide", func(t *testing.T) {
		r := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgB+"/training-guides/"+guideID, nil, tokenB)
		assertStatus(t, r, http.StatusNotFound)
	})

	t.Run("org B cannot update org A's guide", func(t *testing.T) {
		r := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgB+"/training-guides/"+guideID,
			map[string]interface{}{"title": "Hijacked", "content": "Hijacked", "category": "strength"}, tokenB)
		if r.StatusCode == http.StatusOK {
			t.Fatal("org B updated org A's guide — cross-org tampering is open")
		}
	})

	t.Run("org B cannot delete org A's guide", func(t *testing.T) {
		r := doRequest(t, http.MethodDelete, "/api/v1/orgs/"+orgB+"/training-guides/"+guideID, nil, tokenB)
		if r.StatusCode == http.StatusOK {
			t.Fatal("org B deleted org A's guide — the platform library was destroyable")
		}
	})

	// The guide must still be there afterwards.
	t.Run("org A's guide survived", func(t *testing.T) {
		r := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgA+"/training-guides/"+guideID, nil, tokenA)
		assertStatus(t, r, http.StatusOK)
	})

	// And a plain member of org A still cannot manage guides at all.
	t.Run("plain member is forbidden", func(t *testing.T) {
		memberID := createTestUser(t, "9806300003", "Guide Member")
		createTestOrgMember(t, memberID, orgA, "member")
		r := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgA+"/training-guides",
			nil, generateTestToken(memberID, "admin"))
		assertStatus(t, r, http.StatusForbidden)
	})
}

// TestSubscriptionSelfAssignForbidden covers the second half of Task 3's required
// coverage: even within their OWN org, a plain member must never be able to assign
// themselves a package (which would let them pick their own amount_paid).
func TestSubscriptionSelfAssignForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestUser(t, "9806100001", "Sub Admin")
	orgID := createTestOrg(t, adminID, "AuthZ Sub Gym")

	memberID := createTestUser(t, "9806100002", "Sub Member")
	createTestOrgMember(t, memberID, orgID, "member")

	pkgID := createTestPackage(t, orgID, "Monthly", 30, 1000)
	token := generateTestToken(memberID, "member")

	resp := doRequest(t, http.MethodPost,
		"/api/v1/orgs/"+orgID+"/members/"+memberID+"/packages/assign",
		map[string]interface{}{
			"package_id":     pkgID,
			"start_date":     "2026-03-01",
			"payment_method": "cash",
			"amount_paid":    1000.0,
		}, token)
	assertStatus(t, resp, http.StatusForbidden)
}

// TestMemberRoleEscalationGuards covers Task 4: a staff user must not be able to
// promote anyone (including themselves) to a higher role, but must still be able to
// edit non-role fields. A genuine admin remains able to change roles, except their
// own.
func TestMemberRoleEscalationGuards(t *testing.T) {
	cleanupTables(t)

	adminID := createTestUser(t, "9806200001", "Role Admin")
	orgID := createTestOrg(t, adminID, "AuthZ Role Gym")

	staffID := createTestUser(t, "9806200002", "Role Staff")
	createTestOrgMember(t, staffID, orgID, "staff")

	targetID := createTestUser(t, "9806200003", "Role Target")
	createTestOrgMember(t, targetID, orgID, "member")

	staffToken := generateTestToken(staffID, "staff")
	adminToken := generateTestToken(adminID, "admin")

	t.Run("staff cannot promote another member to admin", func(t *testing.T) {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/members/"+targetID,
			map[string]string{"role": "admin"}, staffToken)
		assertStatus(t, resp, http.StatusForbidden)
	})

	t.Run("staff cannot change their own role", func(t *testing.T) {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/members/"+staffID,
			map[string]string{"role": "admin"}, staffToken)
		assertStatus(t, resp, http.StatusForbidden)
	})

	t.Run("staff can still edit non-role fields", func(t *testing.T) {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/members/"+targetID,
			map[string]string{"status": "suspended"}, staffToken)
		assertStatus(t, resp, http.StatusOK)
	})

	t.Run("admin can change another member's role (positive control)", func(t *testing.T) {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/members/"+targetID,
			map[string]string{"role": "staff"}, adminToken)
		assertStatus(t, resp, http.StatusOK)
	})

	t.Run("admin cannot change their own role", func(t *testing.T) {
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/members/"+adminID,
			map[string]string{"role": "staff"}, adminToken)
		assertStatus(t, resp, http.StatusForbidden)
	})
}

// TestNFCMemberForbidden covers the NFC half of Task 8's required coverage
// literally: a plain member of their OWN org (not just a cross-org admin, which
// TestCrossOrgRoleEscalation already covers) must be rejected from device
// registration and card registration/assignment. Positive-control coverage for
// admin success on POST/PUT already exists in tests/e2e/nfc_test.go
// (TestNFC_RegisterCard, TestNFC_AssignCard, TestNFC_ListDevices, etc.) using an
// admin token that is unaffected by this change; a lightweight GET is included here
// too so this file is self-contained.
func TestNFCMemberForbidden(t *testing.T) {
	cleanupTables(t)

	adminID := createTestUser(t, "9806300001", "NFC Admin")
	orgID := createTestOrg(t, adminID, "AuthZ NFC Gym")

	memberID := createTestUser(t, "9806300002", "NFC Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	t.Run("device registration forbidden", func(t *testing.T) {
		resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/nfc-devices",
			map[string]string{"name": "Front Door", "device_identifier": "dev-1", "device_secret": "s3cr3t"}, token)
		assertStatus(t, resp, http.StatusForbidden)
	})

	t.Run("card registration forbidden", func(t *testing.T) {
		resp := doRequest(t, http.MethodPost, "/api/v1/orgs/"+orgID+"/nfc-cards",
			map[string]string{"card_hex": "DEADBEEF"}, token)
		assertStatus(t, resp, http.StatusForbidden)
	})

	t.Run("card assignment forbidden", func(t *testing.T) {
		cardID := createTestNFCCard(t, orgID, "CAFEBABE")
		resp := doRequest(t, http.MethodPut, "/api/v1/orgs/"+orgID+"/nfc-cards/"+cardID+"/assign",
			map[string]string{"user_id": memberID}, token)
		assertStatus(t, resp, http.StatusForbidden)
	})

	t.Run("admin can list devices (positive control)", func(t *testing.T) {
		adminToken := generateTestToken(adminID, "admin")
		resp := doRequest(t, http.MethodGet, "/api/v1/orgs/"+orgID+"/nfc-devices", nil, adminToken)
		assertStatus(t, resp, http.StatusOK)
	})
}
