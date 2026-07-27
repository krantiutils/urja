package member

import (
	"context"
	"errors"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func TestValidateProfileUpdate(t *testing.T) {
	tests := []struct {
		name    string
		update  *ProfileUpdate
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty update is valid",
			update:  &ProfileUpdate{},
			wantErr: false,
		},
		{
			name:    "valid name update",
			update:  &ProfileUpdate{Name: strPtr("Ram Sharma")},
			wantErr: false,
		},
		{
			name:    "empty name rejected",
			update:  &ProfileUpdate{Name: strPtr("")},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name:    "valid email",
			update:  &ProfileUpdate{Email: strPtr("ram@example.com")},
			wantErr: false,
		},
		{
			name:    "invalid email",
			update:  &ProfileUpdate{Email: strPtr("not-an-email")},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "empty email is valid (clearing)",
			update:  &ProfileUpdate{Email: strPtr("")},
			wantErr: false,
		},
		{
			name:    "valid gender male",
			update:  &ProfileUpdate{Gender: strPtr("male")},
			wantErr: false,
		},
		{
			name:    "valid gender female",
			update:  &ProfileUpdate{Gender: strPtr("female")},
			wantErr: false,
		},
		{
			name:    "valid gender other",
			update:  &ProfileUpdate{Gender: strPtr("other")},
			wantErr: false,
		},
		{
			name:    "valid gender prefer_not_say",
			update:  &ProfileUpdate{Gender: strPtr("prefer_not_say")},
			wantErr: false,
		},
		{
			name:    "invalid gender",
			update:  &ProfileUpdate{Gender: strPtr("invalid")},
			wantErr: true,
			errMsg:  "invalid gender: must be male, female, other, or prefer_not_say",
		},
		{
			name:    "valid emergency phone",
			update:  &ProfileUpdate{EmergencyContactPhone: strPtr("9841234567")},
			wantErr: false,
		},
		{
			name:    "invalid emergency phone too short",
			update:  &ProfileUpdate{EmergencyContactPhone: strPtr("98412345")},
			wantErr: true,
			errMsg:  "invalid emergency contact phone number",
		},
		{
			name:    "invalid emergency phone wrong prefix",
			update:  &ProfileUpdate{EmergencyContactPhone: strPtr("1234567890")},
			wantErr: true,
			errMsg:  "invalid emergency contact phone number",
		},
		{
			name: "multiple valid fields",
			update: &ProfileUpdate{
				Name:      strPtr("Ram Sharma"),
				NameNe:    strPtr("राम शर्मा"),
				Email:     strPtr("ram@example.com"),
				Gender:    strPtr("male"),
				AvatarURL: strPtr("https://example.com/avatar.jpg"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProfileUpdate(tt.update)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Fatalf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateCreateMemberInput(t *testing.T) {
	tests := []struct {
		name    string
		input   *CreateMemberInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid minimal input",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "Ram Sharma",
			},
			wantErr: false,
		},
		{
			name: "valid full input",
			input: &CreateMemberInput{
				Phone:  "9841234567",
				Name:   "Ram Sharma",
				NameNe: "राम शर्मा",
				Email:  "ram@example.com",
				Role:   "staff",
			},
			wantErr: false,
		},
		{
			name: "invalid phone",
			input: &CreateMemberInput{
				Phone: "1234567890",
				Name:  "Ram",
			},
			wantErr: true,
			errMsg:  "invalid phone number: must be 10 digits starting with 9, 8, 7, or 6",
		},
		{
			name: "empty phone",
			input: &CreateMemberInput{
				Phone: "",
				Name:  "Ram",
			},
			wantErr: true,
			errMsg:  "invalid phone number: must be 10 digits starting with 9, 8, 7, or 6",
		},
		{
			name: "empty name",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid email",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "Ram",
				Email: "bad-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid role",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "Ram",
				Role:  "superadmin",
			},
			wantErr: true,
			errMsg:  "invalid role: must be member, staff, or admin",
		},
		{
			name: "valid role member",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "Ram",
				Role:  "member",
			},
			wantErr: false,
		},
		{
			name: "valid role admin",
			input: &CreateMemberInput{
				Phone: "9841234567",
				Name:  "Ram",
				Role:  "admin",
			},
			wantErr: false,
		},
		{
			name: "phone starting with 6",
			input: &CreateMemberInput{
				Phone: "6841234567",
				Name:  "Ram",
			},
			wantErr: false,
		},
		{
			name: "phone starting with 7",
			input: &CreateMemberInput{
				Phone: "7841234567",
				Name:  "Ram",
			},
			wantErr: false,
		},
		{
			name: "phone starting with 8",
			input: &CreateMemberInput{
				Phone: "8841234567",
				Name:  "Ram",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateMemberInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Fatalf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateUpdateOrgMemberInput(t *testing.T) {
	tests := []struct {
		name    string
		input   *UpdateOrgMemberInput
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty update rejected",
			input:   &UpdateOrgMemberInput{},
			wantErr: true,
			errMsg:  "at least one of role or status must be provided",
		},
		{
			name:    "valid role update",
			input:   &UpdateOrgMemberInput{Role: strPtr("staff")},
			wantErr: false,
		},
		{
			name:    "valid status update",
			input:   &UpdateOrgMemberInput{Status: strPtr("suspended")},
			wantErr: false,
		},
		{
			name: "valid both fields",
			input: &UpdateOrgMemberInput{
				Role:   strPtr("admin"),
				Status: strPtr("active"),
			},
			wantErr: false,
		},
		{
			name:    "invalid role",
			input:   &UpdateOrgMemberInput{Role: strPtr("owner")},
			wantErr: true,
			errMsg:  "invalid role: must be member, staff, or admin",
		},
		{
			name:    "invalid status",
			input:   &UpdateOrgMemberInput{Status: strPtr("banned")},
			wantErr: true,
			errMsg:  "invalid status: must be active, suspended, or left",
		},
		{
			name:    "valid status left",
			input:   &UpdateOrgMemberInput{Status: strPtr("left")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateOrgMemberInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Fatalf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestUpdateOrgMember_RoleChangeGuards exercises the privilege-escalation guards in
// UpdateOrgMember. All of these cases are rejected before the service ever reaches
// the repository (repo stays nil), so a zero-value Service is sufficient here — no
// database required. Non-role edits and the "staff can edit non-role fields" case
// require a live repository and are covered instead by tests/e2e/authz_test.go.
func TestUpdateOrgMember_RoleChangeGuards(t *testing.T) {
	tests := []struct {
		name       string
		callerID   string
		callerRole string
		memberID   string
	}{
		{
			name:       "staff cannot change another member's role",
			callerID:   "user-1",
			callerRole: "staff",
			memberID:   "user-2",
		},
		{
			name:       "staff cannot change their own role",
			callerID:   "user-1",
			callerRole: "staff",
			memberID:   "user-1",
		},
		{
			name:       "member cannot change another member's role",
			callerID:   "user-1",
			callerRole: "member",
			memberID:   "user-2",
		},
		{
			name:       "admin cannot change their own role",
			callerID:   "user-1",
			callerRole: "admin",
			memberID:   "user-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			err := svc.UpdateOrgMember(context.Background(), tt.callerID, tt.callerRole, "org-1", tt.memberID,
				&UpdateOrgMemberInput{Role: strPtr("admin")})
			if err == nil {
				t.Fatal("expected forbidden error, got nil")
			}
			if !errors.Is(err, ErrForbidden) {
				t.Fatalf("expected ErrForbidden, got %v", err)
			}
		})
	}
}

