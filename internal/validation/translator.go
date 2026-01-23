package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateAll(err error) []string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	out := make([]string, 0, len(errs))
	for _, e := range errs {
		out = append(out, translate(e))
	}

	return out
}

func translate(err validator.FieldError) string {
	field := err.Field()

	msg, ok := tagMessages[err.Tag()]
	if !ok {
		return fmt.Sprintf("%s is invalid", field)
	}

	switch err.Tag() {
	case "min", "max", "gte", "lte":
		return fmt.Sprintf("%s %s", field, fmt.Sprintf(msg, err.Param()))
	case "eqfield":
		return fmt.Sprintf("%s %s", field, fmt.Sprintf(msg, strings.ToLower(err.Param())))
	default:
		return fmt.Sprintf("%s %s", field, msg)
	}
}
