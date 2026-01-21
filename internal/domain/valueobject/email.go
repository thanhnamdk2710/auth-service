package vo

import (
	"regexp"
	"strings"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

const (
	EmailMinLength = 5
	EmailMaxLength = 255
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))

	if trimmed == "" {
		return Email{}, exception.ErrEmailRequired
	}

	if len(trimmed) < EmailMinLength || len(trimmed) > EmailMaxLength {
		return Email{}, exception.ErrEmailMinMaxLength
	}

	if !emailRegex.MatchString(trimmed) {
		return Email{}, exception.ErrEmailInvalid
	}

	return Email{
		value: trimmed,
	}, nil
}

func (e Email) String() string {
	return e.value
}
