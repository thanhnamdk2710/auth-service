package valueobject_test

import (
	"strings"
	"testing"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
)

func TestNewUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid username",
			input:   "john_doe",
			wantErr: nil,
		},
		{
			name:    "valid username with dots",
			input:   "john.doe",
			wantErr: nil,
		},
		{
			name:    "valid username with numbers",
			input:   "john123",
			wantErr: nil,
		},
		{
			name:    "valid username minimum length",
			input:   "abc",
			wantErr: nil,
		},
		{
			name:    "valid username maximum length",
			input:   "abcdefghijklmnopqrstuvwxyz1234",
			wantErr: nil,
		},
		{
			name:    "valid username with spaces trimmed",
			input:   "  john_doe  ",
			wantErr: nil,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: exception.ErrUsernameRequired,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: exception.ErrUsernameRequired,
		},
		{
			name:    "too short",
			input:   "ab",
			wantErr: exception.ErrUsernameMinMaxLength,
		},
		{
			name:    "too long",
			input:   "abcdefghijklmnopqrstuvwxyz12345",
			wantErr: exception.ErrUsernameMinMaxLength,
		},
		{
			name:    "invalid characters - space in middle",
			input:   "john doe",
			wantErr: exception.ErrUsernameFormatInvalid,
		},
		{
			name:    "invalid characters - special chars",
			input:   "john@doe",
			wantErr: exception.ErrUsernameFormatInvalid,
		},
		{
			name:    "invalid characters - unicode",
			input:   "john日本",
			wantErr: exception.ErrUsernameFormatInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := vo.NewUsername(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewUsername(%q) expected error %v, got nil", tt.input, tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewUsername(%q) expected error %v, got %v", tt.input, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUsername(%q) unexpected error: %v", tt.input, err)
				return
			}

			// Verify value is trimmed
			expected := strings.TrimSpace(tt.input)
			if username.String() != expected {
				t.Errorf("NewUsername(%q).String() = %q, want %q", tt.input, username.String(), expected)
			}
		})
	}
}

func TestUsername_String(t *testing.T) {
	input := "john_doe"
	username, err := vo.NewUsername(input)
	if err != nil {
		t.Fatalf("NewUsername(%q) unexpected error: %v", input, err)
	}

	if username.String() != input {
		t.Errorf("Username.String() = %q, want %q", username.String(), input)
	}
}
