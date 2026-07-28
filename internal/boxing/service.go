package boxing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ErrForbidden is returned when a caller may not act on the target member.
var ErrForbidden = errors.New("forbidden")

// Service handles boxing profile business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new boxing service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ProfileView is a profile plus its derived bout record.
type ProfileView struct {
	*Profile
	Record Record `json:"record"`
	Bouts  []Bout `json:"bouts"`
	// SuggestedWeightClass is derived from the member's latest logged weight.
	// Empty when they have never logged one.
	SuggestedWeightClass string `json:"suggested_weight_class,omitempty"`
}

// GetProfile returns a member's profile, bouts and derived record. A member who
// has no profile row yet gets an empty one rather than a 404, so the UI has
// something to render on first visit.
func (s *Service) GetProfile(ctx context.Context, orgID, userID string) (*ProfileView, error) {
	profile, err := s.repo.GetProfile(ctx, orgID, userID)
	if errors.Is(err, ErrNotFound) {
		profile = &Profile{UserID: userID, OrgID: orgID}
	} else if err != nil {
		return nil, err
	}

	bouts, err := s.repo.ListBouts(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	view := &ProfileView{Profile: profile, Record: TallyBouts(bouts), Bouts: bouts}

	// Suggest a division from the latest logged weight. A member who has never
	// logged one simply gets no suggestion — not a bogus light flyweight.
	if kg, err := s.repo.LatestWeightKg(ctx, userID); err == nil {
		view.SuggestedWeightClass = WeightClassFor(kg)
	}

	return view, nil
}

// UpdateProfileInput carries the member-editable fields.
type UpdateProfileInput struct {
	Stance      string
	WeightClass string
	SkillLevel  string
	ReachCm     *float64
	Notes       string
}

// UpdateProfile validates and saves the member-editable fields. It cannot
// change sparring clearance: that is a staff safety decision, and the field is
// absent from both this input and the repository's update statement.
func (s *Service) UpdateProfile(ctx context.Context, orgID, userID string, in UpdateProfileInput) (*ProfileView, error) {
	stance := strings.ToLower(strings.TrimSpace(in.Stance))
	skill := strings.ToLower(strings.TrimSpace(in.SkillLevel))
	weightClass := strings.ToLower(strings.TrimSpace(in.WeightClass))

	if err := ValidateStance(stance); err != nil {
		return nil, err
	}
	if err := ValidateSkillLevel(skill); err != nil {
		return nil, err
	}
	if in.ReachCm != nil && (*in.ReachCm <= 0 || *in.ReachCm > 300) {
		return nil, fmt.Errorf("reach must be between 0 and 300 cm")
	}
	if len(in.Notes) > 2000 {
		return nil, fmt.Errorf("notes are too long")
	}

	if _, err := s.repo.UpsertProfile(ctx, orgID, userID, ProfileUpdate{
		Stance: stance, WeightClass: weightClass, SkillLevel: skill,
		ReachCm: in.ReachCm, Notes: strings.TrimSpace(in.Notes),
	}); err != nil {
		return nil, err
	}

	s.logger.Info("boxing profile updated", "user_id", userID, "org_id", orgID)
	return s.GetProfile(ctx, orgID, userID)
}

// SetSparringClearance grants or revokes clearance. Staff-only; the handler
// enforces the role, and this method additionally refuses to let a caller clear
// themselves so a coach cannot self-certify.
func (s *Service) SetSparringClearance(ctx context.Context, orgID, callerID, targetID string, cleared bool) (*ProfileView, error) {
	if callerID == targetID {
		return nil, fmt.Errorf("%w: sparring clearance must be granted by another member of staff", ErrForbidden)
	}

	isMember, err := s.repo.IsOrgMember(ctx, orgID, targetID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("%w: target is not an active member of this organization", ErrForbidden)
	}

	if _, err := s.repo.SetSparringClearance(ctx, orgID, targetID, callerID, cleared); err != nil {
		return nil, err
	}

	s.logger.Info("sparring clearance changed", "target", targetID, "org_id", orgID,
		"cleared", cleared, "by", callerID)
	return s.GetProfile(ctx, orgID, targetID)
}

// CreateBoutInput carries a new bout record.
type CreateBoutInput struct {
	BoutDate    string
	Opponent    string
	EventName   string
	Result      string
	Method      string
	Rounds      *int
	WeightClass string
	Notes       string
}

// CreateBout validates and records a bout.
func (s *Service) CreateBout(ctx context.Context, orgID, userID string, in CreateBoutInput) (*Bout, error) {
	date, err := time.Parse("2006-01-02", strings.TrimSpace(in.BoutDate))
	if err != nil {
		return nil, fmt.Errorf("bout_date must be YYYY-MM-DD")
	}
	// A bout in the future is a data-entry slip, not a record.
	if date.After(time.Now().AddDate(0, 0, 1)) {
		return nil, fmt.Errorf("bout_date cannot be in the future")
	}

	result := strings.ToLower(strings.TrimSpace(in.Result))
	method := strings.ToLower(strings.TrimSpace(in.Method))
	if err := ValidateResult(result); err != nil {
		return nil, err
	}
	if err := ValidateMethod(method); err != nil {
		return nil, err
	}
	if in.Rounds != nil && (*in.Rounds < 1 || *in.Rounds > 15) {
		return nil, fmt.Errorf("rounds must be between 1 and 15")
	}

	bout, err := s.repo.CreateBout(ctx, orgID, userID, &Bout{
		BoutDate:    date,
		Opponent:    strings.TrimSpace(in.Opponent),
		EventName:   strings.TrimSpace(in.EventName),
		Result:      result,
		Method:      method,
		Rounds:      in.Rounds,
		WeightClass: strings.ToLower(strings.TrimSpace(in.WeightClass)),
		Notes:       strings.TrimSpace(in.Notes),
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("bout recorded", "bout_id", bout.ID, "user_id", userID, "org_id", orgID, "result", result)
	return bout, nil
}

// DeleteBout removes a bout.
func (s *Service) DeleteBout(ctx context.Context, orgID, userID, boutID string) error {
	return s.repo.DeleteBout(ctx, orgID, userID, boutID)
}

// EnsureOrgMember reports an error unless the target is an active member.
func (s *Service) EnsureOrgMember(ctx context.Context, orgID, userID string) error {
	isMember, err := s.repo.IsOrgMember(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("%w: not an active member of this organization", ErrForbidden)
	}
	return nil
}
