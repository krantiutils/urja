package attendance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const qrValidityWindow = 5 * time.Minute

// Nepal Standard Time — fixed zone, no DST.
var nptLocation = time.FixedZone("NPT", 5*3600+45*60)

type qrPayload struct {
	OrgID     string `json:"oid"`
	Timestamp int64  `json:"ts"`
}

// CheckInResult bundles the attendance record with the updated streak.
type CheckInResult struct {
	Record *Record `json:"check_in"`
	Streak *Streak `json:"streak,omitempty"`
}

// Service handles attendance business logic.
type Service struct {
	repo         *Repository
	qrSigningKey []byte
	logger       *slog.Logger
}

// NewService creates a new attendance service.
// qrSigningKey is the HMAC key used to verify QR code signatures.
func NewService(repo *Repository, qrSigningKey []byte, logger *slog.Logger) *Service {
	return &Service{repo: repo, qrSigningKey: qrSigningKey, logger: logger}
}

// SelfCheckIn handles a member checking themselves in via QR or NFC.
func (s *Service) SelfCheckIn(ctx context.Context, userID, method, qrToken, cardUID, orgID string) (*CheckInResult, error) {
	switch method {
	case "qr":
		return s.checkInQR(ctx, userID, qrToken)
	case "nfc":
		return s.checkInNFC(ctx, userID, orgID, cardUID)
	default:
		return nil, fmt.Errorf("invalid self check-in method: %s (expected qr or nfc)", method)
	}
}

// checkInQR validates a QR token and records a check-in.
// Token format: <base64url_payload>.<hex_hmac_signature>
// Payload JSON: {"oid":"<org_uuid>","ts":<unix_seconds>}
func (s *Service) checkInQR(ctx context.Context, userID, qrToken string) (*CheckInResult, error) {
	if qrToken == "" {
		return nil, fmt.Errorf("qr_token is required for QR check-in")
	}

	parts := strings.SplitN(qrToken, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid QR token format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid QR payload encoding")
	}

	sigBytes, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid QR signature encoding")
	}

	// Verify HMAC-SHA256 signature.
	mac := hmac.New(sha256.New, s.qrSigningKey)
	mac.Write(payloadBytes)
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(sigBytes, expectedMAC) {
		return nil, fmt.Errorf("invalid QR signature")
	}

	var payload qrPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("invalid QR payload")
	}

	// Reject stale or future-dated QR codes.
	elapsed := time.Since(time.Unix(payload.Timestamp, 0))
	if elapsed < -qrValidityWindow || elapsed > qrValidityWindow {
		return nil, fmt.Errorf("QR code has expired")
	}

	if payload.OrgID == "" {
		return nil, fmt.Errorf("QR payload missing organization ID")
	}

	isMember, err := s.repo.IsOrgMember(ctx, userID, payload.OrgID)
	if err != nil {
		return nil, fmt.Errorf("checking membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("not a member of this organization")
	}

	return s.RecordCheckIn(ctx, userID, payload.OrgID, "qr")
}

// checkInNFC verifies an NFC card and records a check-in.
func (s *Service) checkInNFC(ctx context.Context, userID, orgID, cardUID string) (*CheckInResult, error) {
	if orgID == "" {
		return nil, fmt.Errorf("org_id is required for NFC check-in")
	}
	if cardUID == "" {
		return nil, fmt.Errorf("card_uid is required for NFC check-in")
	}

	cardOwnerID, err := s.repo.LookupNFCCard(ctx, orgID, cardUID)
	if err != nil {
		return nil, fmt.Errorf("NFC card verification failed: %w", err)
	}

	if cardOwnerID != userID {
		return nil, fmt.Errorf("NFC card is not registered to this user")
	}

	return s.RecordCheckIn(ctx, userID, orgID, "nfc")
}

// ManualCheckIn allows staff to manually check in a member at an org.
func (s *Service) ManualCheckIn(ctx context.Context, staffUserID, memberUserID, orgID string) (*CheckInResult, error) {
	if memberUserID == "" {
		return nil, fmt.Errorf("member_id is required for manual check-in")
	}

	role, err := s.repo.GetOrgRole(ctx, staffUserID, orgID)
	if err != nil {
		return nil, fmt.Errorf("authorization check failed: %w", err)
	}
	if role != "staff" && role != "admin" {
		return nil, fmt.Errorf("insufficient role: only staff and admins can perform manual check-ins")
	}

	isMember, err := s.repo.IsOrgMember(ctx, memberUserID, orgID)
	if err != nil {
		return nil, fmt.Errorf("checking membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("target user is not a member of this organization")
	}

	return s.RecordCheckIn(ctx, memberUserID, orgID, "manual")
}

// RecordCheckIn inserts the attendance row and upserts the streak.
func (s *Service) RecordCheckIn(ctx context.Context, userID, orgID, method string) (*CheckInResult, error) {
	rec, err := s.repo.CheckIn(ctx, userID, orgID, method)
	if err != nil {
		return nil, fmt.Errorf("recording check-in: %w", err)
	}

	dateStr := rec.CheckInAt.In(nptLocation).Format("2006-01-02")

	streak, err := s.repo.UpsertStreak(ctx, userID, orgID, dateStr)
	if err != nil {
		// Streak update failed but attendance was recorded — log and continue.
		s.logger.Error("failed to update streak",
			"error", err, "user_id", userID, "org_id", orgID)
		return &CheckInResult{Record: rec}, nil
	}

	s.logger.Info("member checked in",
		"user_id", userID,
		"org_id", orgID,
		"method", method,
		"current_streak", streak.CurrentStreak,
	)

	return &CheckInResult{Record: rec, Streak: streak}, nil
}

// ListByUser lists a user's attendance history.
func (s *Service) ListByUser(ctx context.Context, userID string, limit, offset int) ([]Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

// ListByOrg lists attendance for an organization.
func (s *Service) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByOrg(ctx, orgID, limit, offset)
}

// GetStreaks returns all attendance streaks for a user across organizations.
func (s *Service) GetStreaks(ctx context.Context, userID string) ([]Streak, error) {
	return s.repo.GetStreaksByUser(ctx, userID)
}

// CalendarResponse represents a monthly attendance calendar.
type CalendarResponse struct {
	Month         string `json:"month"`
	DaysInMonth   int    `json:"days_in_month"`
	CheckInDays   []int  `json:"check_in_days"`
	TotalCheckIns int    `json:"total_check_ins"`
}

var monthRegex = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

// GetMonthlyCalendar returns the attendance calendar for a given user and month.
func (s *Service) GetMonthlyCalendar(ctx context.Context, userID, month string) (*CalendarResponse, error) {
	if !monthRegex.MatchString(month) {
		return nil, fmt.Errorf("invalid month format: expected YYYY-MM")
	}

	// Parse year and month to compute days in month.
	parts := strings.SplitN(month, "-", 2)
	year, _ := strconv.Atoi(parts[0])
	monthNum, _ := strconv.Atoi(parts[1])
	daysInMonth := time.Date(year, time.Month(monthNum)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	days, err := s.repo.GetMonthlyCalendar(ctx, userID, month)
	if err != nil {
		return nil, fmt.Errorf("getting monthly calendar: %w", err)
	}

	if days == nil {
		days = []int{}
	}

	return &CalendarResponse{
		Month:         month,
		DaysInMonth:   daysInMonth,
		CheckInDays:   days,
		TotalCheckIns: len(days),
	}, nil
}
