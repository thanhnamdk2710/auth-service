package vo

import (
	"regexp"
	"strings"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+$`)
)

const (
	UsernameMinLength = 3
	UsernameMaxLength = 30
)

type Username struct {
	value string
}

func NewUsername(value string) (*Username, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return nil, exception.ErrUsernameRequired
	}

	if len(trimmed) < UsernameMinLength || len(trimmed) > UsernameMaxLength {
		return nil, exception.ErrUsernameMinMaxLength
	}

	if !usernameRegex.MatchString(trimmed) {
		return nil, exception.ErrUsernameFormatInvalid
	}

	return &Username{
		value: trimmed,
	}, nil
}

func (e Username) String() string {
	return e.value
}
