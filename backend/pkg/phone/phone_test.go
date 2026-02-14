package phone

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"plain valid", "9841234567", "9841234567", false},
		{"with +977", "+9779841234567", "9841234567", false},
		{"with 977 prefix", "9779841234567", "9841234567", false},
		{"with spaces", "984 123 4567", "9841234567", false},
		{"with dashes", "984-123-4567", "9841234567", false},
		{"starts with 6", "6012345678", "6012345678", false},
		{"starts with 7", "7012345678", "7012345678", false},
		{"starts with 8", "8012345678", "8012345678", false},

		// Invalid
		{"too short", "984123456", "", true},
		{"too long", "98412345678", "", true},
		{"starts with 1", "1234567890", "", true},
		{"starts with 0", "0234567890", "", true},
		{"starts with 5", "5234567890", "", true},
		{"empty", "", "", true},
		{"letters", "98abc41234", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid 9", "9841234567", false},
		{"valid 6", "6012345678", false},
		{"valid 7", "7012345678", false},
		{"valid 8", "8012345678", false},
		{"invalid start 1", "1234567890", true},
		{"invalid start 0", "0234567890", true},
		{"too short", "98412", true},
		{"too long", "98412345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
