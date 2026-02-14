package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// insertTestStreak inserts a streak row directly for testing leaderboard queries.
func insertTestStreak(t *testing.T, memberID, orgID string, currentStreak, longestStreak int) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO streaks (member_id, organization_id, current_streak, longest_streak, last_check_in)
		 VALUES ($1, $2, $3, $4, CURRENT_DATE)
		 ON CONFLICT (member_id, organization_id) DO UPDATE
		 SET current_streak = $3, longest_streak = $4, last_check_in = CURRENT_DATE`,
		memberID, orgID, currentStreak, longestStreak,
	)
	if err != nil {
		t.Fatalf("insertTestStreak: %v", err)
	}
}

// insertTestAttendance inserts an attendance row directly for testing monthly leaderboard.
func insertTestAttendance(t *testing.T, userID, orgID string) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO attendance (user_id, organization_id, check_in_method) VALUES ($1, $2, 'manual')`,
		userID, orgID,
	)
	if err != nil {
		t.Fatalf("insertTestAttendance: %v", err)
	}
}

func TestLeaderboard_Weekly(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Leaderboard Gym")

	member1 := createTestUser(t, "9800500002", "Alice")
	member2 := createTestUser(t, "9800500003", "Bob")
	member3 := createTestUser(t, "9800500004", "Charlie")

	createTestOrgMember(t, member1, orgID, "member")
	createTestOrgMember(t, member2, orgID, "member")
	createTestOrgMember(t, member3, orgID, "member")

	insertTestStreak(t, member1, orgID, 10, 15)
	insertTestStreak(t, member2, orgID, 5, 20)
	insertTestStreak(t, member3, orgID, 8, 8)

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=weekly", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Period   string `json:"period"`
		Rankings []struct {
			Rank     int    `json:"rank"`
			MemberID string `json:"member_id"`
			Name     string `json:"name"`
			Value    int    `json:"value"`
			Metric   string `json:"metric"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if body.Period != "weekly" {
		t.Errorf("expected period=weekly, got %s", body.Period)
	}
	if len(body.Rankings) != 3 {
		t.Fatalf("expected 3 rankings, got %d", len(body.Rankings))
	}
	// Alice (10) > Charlie (8) > Bob (5)
	if body.Rankings[0].Name != "Alice" || body.Rankings[0].Value != 10 || body.Rankings[0].Rank != 1 {
		t.Errorf("rank 1: expected Alice with value 10, got %s with %d", body.Rankings[0].Name, body.Rankings[0].Value)
	}
	if body.Rankings[1].Name != "Charlie" || body.Rankings[1].Value != 8 {
		t.Errorf("rank 2: expected Charlie with value 8, got %s with %d", body.Rankings[1].Name, body.Rankings[1].Value)
	}
	if body.Rankings[2].Name != "Bob" || body.Rankings[2].Value != 5 {
		t.Errorf("rank 3: expected Bob with value 5, got %s with %d", body.Rankings[2].Name, body.Rankings[2].Value)
	}
	if body.Rankings[0].Metric != "current_streak" {
		t.Errorf("expected metric=current_streak, got %s", body.Rankings[0].Metric)
	}
}

func TestLeaderboard_Monthly(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Monthly Gym")

	member1 := createTestUser(t, "9800500002", "Alice")
	member2 := createTestUser(t, "9800500003", "Bob")

	createTestOrgMember(t, member1, orgID, "member")
	createTestOrgMember(t, member2, orgID, "member")

	// Alice: 3 check-ins this month, Bob: 1
	for i := 0; i < 3; i++ {
		insertTestAttendance(t, member1, orgID)
	}
	insertTestAttendance(t, member2, orgID)

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=monthly", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Period   string `json:"period"`
		Rankings []struct {
			Rank     int    `json:"rank"`
			MemberID string `json:"member_id"`
			Name     string `json:"name"`
			Value    int    `json:"value"`
			Metric   string `json:"metric"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if body.Period != "monthly" {
		t.Errorf("expected period=monthly, got %s", body.Period)
	}
	if len(body.Rankings) != 2 {
		t.Fatalf("expected 2 rankings, got %d", len(body.Rankings))
	}
	if body.Rankings[0].Name != "Alice" || body.Rankings[0].Value != 3 {
		t.Errorf("rank 1: expected Alice with 3 check-ins, got %s with %d", body.Rankings[0].Name, body.Rankings[0].Value)
	}
	if body.Rankings[1].Name != "Bob" || body.Rankings[1].Value != 1 {
		t.Errorf("rank 2: expected Bob with 1 check-in, got %s with %d", body.Rankings[1].Name, body.Rankings[1].Value)
	}
	if body.Rankings[0].Metric != "check_ins" {
		t.Errorf("expected metric=check_ins, got %s", body.Rankings[0].Metric)
	}
}

func TestLeaderboard_AllTime(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "AllTime Gym")

	member1 := createTestUser(t, "9800500002", "Alice")
	member2 := createTestUser(t, "9800500003", "Bob")

	createTestOrgMember(t, member1, orgID, "member")
	createTestOrgMember(t, member2, orgID, "member")

	insertTestStreak(t, member1, orgID, 3, 15)
	insertTestStreak(t, member2, orgID, 7, 25)

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=alltime", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Period   string `json:"period"`
		Rankings []struct {
			Rank     int    `json:"rank"`
			MemberID string `json:"member_id"`
			Name     string `json:"name"`
			Value    int    `json:"value"`
			Metric   string `json:"metric"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if body.Period != "alltime" {
		t.Errorf("expected period=alltime, got %s", body.Period)
	}
	if len(body.Rankings) != 2 {
		t.Fatalf("expected 2 rankings, got %d", len(body.Rankings))
	}
	// Bob (longest=25) > Alice (longest=15)
	if body.Rankings[0].Name != "Bob" || body.Rankings[0].Value != 25 {
		t.Errorf("rank 1: expected Bob with 25, got %s with %d", body.Rankings[0].Name, body.Rankings[0].Value)
	}
	if body.Rankings[1].Name != "Alice" || body.Rankings[1].Value != 15 {
		t.Errorf("rank 2: expected Alice with 15, got %s with %d", body.Rankings[1].Name, body.Rankings[1].Value)
	}
	if body.Rankings[0].Metric != "longest_streak" {
		t.Errorf("expected metric=longest_streak, got %s", body.Rankings[0].Metric)
	}
}

func TestLeaderboard_DefaultPeriod(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Default Gym")
	token := generateTestToken(adminID, "admin")

	// No period param — should default to "weekly"
	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Period   string        `json:"period"`
		Rankings []interface{} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if body.Period != "weekly" {
		t.Errorf("expected default period=weekly, got %s", body.Period)
	}
	if body.Rankings == nil {
		t.Error("expected rankings to be non-nil (empty array)")
	}
}

func TestLeaderboard_InvalidPeriod(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Invalid Gym")
	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=daily", nil, token)
	assertStatus(t, resp, http.StatusBadRequest)
}

func TestLeaderboard_PrivacyOptOut(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Privacy Gym")

	member1 := createTestUser(t, "9800500002", "Visible")
	member2 := createTestUser(t, "9800500003", "Hidden")

	createTestOrgMember(t, member1, orgID, "member")
	createTestOrgMember(t, member2, orgID, "member")

	insertTestStreak(t, member1, orgID, 10, 10)
	insertTestStreak(t, member2, orgID, 20, 20)

	// Opt out member2 from leaderboard
	_, err := testPool.Exec(context.Background(),
		`UPDATE users SET privacy_settings = jsonb_set(
			privacy_settings, '{show_on_leaderboard}', 'false'
		) WHERE id = $1`, member2)
	if err != nil {
		t.Fatalf("updating privacy settings: %v", err)
	}

	token := generateTestToken(adminID, "admin")

	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=weekly", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Rankings []struct {
			MemberID string `json:"member_id"`
			Name     string `json:"name"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if len(body.Rankings) != 1 {
		t.Fatalf("expected 1 ranking (privacy opt-out), got %d", len(body.Rankings))
	}
	if body.Rankings[0].Name != "Visible" {
		t.Errorf("expected Visible member, got %s", body.Rankings[0].Name)
	}
}

func TestLeaderboard_OrgScoping(t *testing.T) {
	cleanupTables(t)

	admin1 := createTestSuperAdmin(t, "9800500001", "Admin1")
	admin2 := createTestSuperAdmin(t, "9800500005", "Admin2")
	org1 := createTestOrg(t, admin1, "Gym One")
	org2 := createTestOrg(t, admin2, "Gym Two")

	member1 := createTestUser(t, "9800500002", "Org1Member")
	member2 := createTestUser(t, "9800500003", "Org2Member")

	createTestOrgMember(t, member1, org1, "member")
	createTestOrgMember(t, member2, org2, "member")

	insertTestStreak(t, member1, org1, 10, 10)
	insertTestStreak(t, member2, org2, 20, 20)

	// Query org1's leaderboard — should only see member1
	token := generateTestToken(admin1, "admin")
	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+org1+"/leaderboard?period=weekly", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Rankings []struct {
			Name string `json:"name"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &body)

	if len(body.Rankings) != 1 {
		t.Fatalf("expected 1 ranking for org1, got %d", len(body.Rankings))
	}
	if body.Rankings[0].Name != "Org1Member" {
		t.Errorf("expected Org1Member, got %s", body.Rankings[0].Name)
	}
}

func TestLeaderboard_Pagination(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800500001", "Admin")
	orgID := createTestOrg(t, adminID, "Pagination Gym")

	// Create 5 members with different streaks
	for i := 1; i <= 5; i++ {
		phone := fmt.Sprintf("98005001%02d", i)
		name := fmt.Sprintf("Member%d", i)
		memberID := createTestUser(t, phone, name)
		createTestOrgMember(t, memberID, orgID, "member")
		insertTestStreak(t, memberID, orgID, i*2, i*3)
	}

	token := generateTestToken(adminID, "admin")

	// Request limit=2, offset=0 — should get top 2
	resp := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=weekly&limit=2&offset=0", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var page1 struct {
		Rankings []struct {
			Rank  int `json:"rank"`
			Value int `json:"value"`
		} `json:"rankings"`
	}
	parseJSON(t, resp, &page1)

	if len(page1.Rankings) != 2 {
		t.Fatalf("expected 2 results for page 1, got %d", len(page1.Rankings))
	}
	if page1.Rankings[0].Rank != 1 {
		t.Errorf("expected rank 1, got %d", page1.Rankings[0].Rank)
	}

	// Request limit=2, offset=2 — should get next 2 with correct ranks
	resp2 := doRequest(t, http.MethodGet,
		"/api/v1/orgs/"+orgID+"/leaderboard?period=weekly&limit=2&offset=2", nil, token)
	assertStatus(t, resp2, http.StatusOK)

	var page2 struct {
		Rankings []struct {
			Rank  int `json:"rank"`
			Value int `json:"value"`
		} `json:"rankings"`
	}
	parseJSON(t, resp2, &page2)

	if len(page2.Rankings) != 2 {
		t.Fatalf("expected 2 results for page 2, got %d", len(page2.Rankings))
	}
	if page2.Rankings[0].Rank != 3 {
		t.Errorf("expected rank 3 for first item on page 2, got %d", page2.Rankings[0].Rank)
	}
}

func TestLeaderboard_NoAuth(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/v1/orgs/some-org-id/leaderboard", nil, "")
	assertStatus(t, resp, http.StatusUnauthorized)
}
