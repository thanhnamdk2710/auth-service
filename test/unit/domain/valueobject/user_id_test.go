package valueobject_test

import (
	"strings"
	"testing"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
)

func TestNewUserID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid UUID v7",
			input:   "0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f",
			wantErr: nil,
		},
		{
			name:    "valid UUID v7 uppercase",
			input:   "0190A5B0-7E1C-7B3D-8F4E-9A1B2C3D4E5F",
			wantErr: nil,
		},
		{
			name:    "valid UUID v7 with spaces",
			input:   "  0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f  ",
			wantErr: nil,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: exception.ErrUserIDRequired,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: exception.ErrUserIDRequired,
		},
		{
			name:    "invalid format - no dashes",
			input:   "0190a5b07e1c7b3d8f4e9a1b2c3d4e5f",
			wantErr: exception.ErrUserIDInvalid,
		},
		{
			name:    "invalid format - wrong version (v4 instead of v7)",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: exception.ErrUserIDInvalid,
		},
		{
			name:    "invalid format - too short",
			input:   "0190a5b0-7e1c-7b3d",
			wantErr: exception.ErrUserIDInvalid,
		},
		{
			name:    "invalid format - invalid characters",
			input:   "0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4xyz",
			wantErr: exception.ErrUserIDInvalid,
		},
		{
			name:    "invalid variant bits",
			input:   "0190a5b0-7e1c-7b3d-0f4e-9a1b2c3d4e5f",
			wantErr: exception.ErrUserIDInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := vo.NewUserID(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewUserID(%q) expected error %v, got nil", tt.input, tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewUserID(%q) expected error %v, got %v", tt.input, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUserID(%q) unexpected error: %v", tt.input, err)
				return
			}

			// Verify value is lowercased and trimmed
			expected := strings.TrimSpace(strings.ToLower(tt.input))
			if userID.String() != expected {
				t.Errorf("NewUserID(%q).String() = %q, want %q", tt.input, userID.String(), expected)
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	input := "0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f"
	userID, err := vo.NewUserID(input)
	if err != nil {
		t.Fatalf("NewUserID(%q) unexpected error: %v", input, err)
	}

	if userID.String() != input {
		t.Errorf("UserID.String() = %q, want %q", userID.String(), input)
	}
}
