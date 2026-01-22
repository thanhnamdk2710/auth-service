package vo

import (
	"regexp"
	"strings"

	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
)

var (
	// UUID v7 format: xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-7[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
)

type UserID struct {
	value string
}

func NewUserID(value string) (UserID, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))

	if trimmed == "" {
		return UserID{}, exception.ErrUserIDRequired
	}

	if !uuidRegex.MatchString(trimmed) {
		return UserID{}, exception.ErrUserIDInvalid
	}

	return UserID{
		value: trimmed,
	}, nil
}

func (e UserID) String() string {
	return e.value
}
