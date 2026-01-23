package validation_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

func TestTranslateAll_NonValidationError(t *testing.T) {
	result := validation.TranslateAll(errors.New("some other error"))
	if result != nil {
		t.Errorf("expected nil for non-validation error, got %v", result)
	}
}

func TestTranslateAll_Required(t *testing.T) {
	type req struct {
		Name string `validate:"required" json:"name"`
	}

	errs := validateStruct(t, req{Name: ""})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "Name is required")
}

func TestTranslateAll_Gte(t *testing.T) {
	type req struct {
		Name string `validate:"gte=3" json:"name"`
	}

	errs := validateStruct(t, req{Name: "ab"})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "Name must be at least 3 characters")
}

func TestTranslateAll_Lte(t *testing.T) {
	type req struct {
		Name string `validate:"lte=5" json:"name"`
	}

	errs := validateStruct(t, req{Name: "abcdef"})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "Name must be at most 5 characters")
}

func TestTranslateAll_Email(t *testing.T) {
	type req struct {
		Email string `validate:"email" json:"email"`
	}

	errs := validateStruct(t, req{Email: "not-an-email"})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "Email must be a valid email address")
}

func TestTranslateAll_Eqfield(t *testing.T) {
	type req struct {
		Password             string `validate:"required" json:"password"`
		PasswordConfirmation string `validate:"eqfield=Password" json:"password_confirmation"`
	}

	errs := validateStruct(t, req{Password: "secret123", PasswordConfirmation: "mismatch"})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "PasswordConfirmation must match the password field")
}

func TestTranslateAll_UnknownTag(t *testing.T) {
	type req struct {
		URL string `validate:"url" json:"url"`
	}

	errs := validateStruct(t, req{URL: "not-a-url"})
	msgs := validation.TranslateAll(errs)

	assertContains(t, msgs, "URL is invalid")
}

func TestTranslateAll_MultipleErrors(t *testing.T) {
	type req struct {
		Name  string `validate:"required" json:"name"`
		Email string `validate:"required" json:"email"`
	}

	errs := validateStruct(t, req{})
	msgs := validation.TranslateAll(errs)

	if len(msgs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(msgs), msgs)
	}
}

// validateStruct validates the given struct and returns the error.
// Fails the test if validation passes unexpectedly.
func validateStruct(t *testing.T, s interface{}) error {
	t.Helper()
	v := validator.New()
	err := v.Struct(s)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	return err
}

// assertContains checks that the expected message exists in the slice.
func assertContains(t *testing.T, msgs []string, expected string) {
	t.Helper()
	for _, msg := range msgs {
		if msg == expected {
			return
		}
	}
	t.Errorf("expected %q in messages, got %v", expected, msgs)
}
