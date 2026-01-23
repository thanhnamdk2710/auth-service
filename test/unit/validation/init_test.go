package validation_test

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

func TestInit_UsesJSONTagNames(t *testing.T) {
	validation.Init()

	type req struct {
		FullName string `json:"full_name" validate:"required" binding:"required"`
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("failed to get validator engine")
	}

	err := v.Struct(req{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatal("expected ValidationErrors type")
	}

	field := errs[0].Field()
	if field != "full_name" {
		t.Errorf("expected field name %q from json tag, got %q", "full_name", field)
	}
}
