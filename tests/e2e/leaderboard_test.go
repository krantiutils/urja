package e2e

import (
	"net/http"
	"testing"
)

// A period outside the accepted vocabulary is the caller's mistake, not a
// server fault. It returned 500, which made a typo in a query string look like
// an outage — and did exactly that: the member leaderboard screen sent
// "month" instead of "monthly" and reported an internal server error.
func TestLeaderboard_InvalidPeriodIsBadRequest(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Board Gym")
	memberID := createTestUser(t, "9800600002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	for _, period := range []string{"month", "week", "all", "nonsense"} {
		resp := doRequest(t, http.MethodGet,
			"/api/v1/members/me/leaderboard?period="+period, nil, token)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("period %q returned %d, want 400", period, resp.StatusCode)
		}
	}
}

func TestLeaderboard_AcceptedPeriods(t *testing.T) {
	cleanupTables(t)

	adminID := createTestSuperAdmin(t, "9800600001", "Admin")
	orgID := createTestOrg(t, adminID, "Board Gym")
	memberID := createTestUser(t, "9800600002", "Member")
	createTestOrgMember(t, memberID, orgID, "member")
	token := generateTestToken(memberID, "member")

	for _, period := range []string{"weekly", "monthly", "alltime"} {
		resp := doRequest(t, http.MethodGet,
			"/api/v1/members/me/leaderboard?period="+period, nil, token)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("period %q returned %d, want 200", period, resp.StatusCode)
		}
	}
}
