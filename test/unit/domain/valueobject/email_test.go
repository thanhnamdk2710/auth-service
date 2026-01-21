package valueobject_test

import (
	"strings"
	"testing"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	vo "github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid email",
			input:   "test@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with subdomain",
			input:   "test@mail.example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with plus",
			input:   "test+tag@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with dots",
			input:   "first.last@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email uppercase converted to lowercase",
			input:   "TEST@EXAMPLE.COM",
			wantErr: nil,
		},
		{
			name:    "valid email with spaces trimmed",
			input:   "  test@example.com  ",
			wantErr: nil,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: exception.ErrEmailRequired,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: exception.ErrEmailRequired,
		},
		{
			name:    "too short",
			input:   "a@b",
			wantErr: exception.ErrEmailMinMaxLength,
		},
		{
			name:    "missing @",
			input:   "testexample.com",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "missing domain",
			input:   "test@",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "missing local part",
			input:   "@example.com",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "missing TLD",
			input:   "test@example",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "double @",
			input:   "test@@example.com",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "spaces in middle",
			input:   "test @example.com",
			wantErr: exception.ErrEmailInvalid,
		},
		{
			name:    "invalid TLD - single char",
			input:   "test@example.c",
			wantErr: exception.ErrEmailInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewEmail(%q) expected error %v, got nil", tt.input, tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewEmail(%q) expected error %v, got %v", tt.input, tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("NewEmail(%q) unexpected error: %v", tt.input, err)
				return
			}

			// Verify value is lowercased and trimmed
			expected := strings.TrimSpace(strings.ToLower(tt.input))
			if email.String() != expected {
				t.Errorf("NewEmail(%q).String() = %q, want %q", tt.input, email.String(), expected)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	input := "test@example.com"
	email, err := vo.NewEmail(input)
	if err != nil {
		t.Fatalf("NewEmail(%q) unexpected error: %v", input, err)
	}

	if email.String() != input {
		t.Errorf("Email.String() = %q, want %q", email.String(), input)
	}
}
