package exception

import "errors"

var (
	ErrUserAlreadyActive    = errors.New("User already active")
	ErrUserInactive         = errors.New("User is inactive")
	ErrUserAlreadyInactive  = errors.New("User already inactive")
	ErrEmailAlreadyVerified = errors.New("Email already verified")

	ErrEmailRequired     = errors.New("Email is required")
	ErrEmailMinMaxLength = errors.New("Email must be between 5 and 255 characters")
	ErrEmailInvalid      = errors.New("Email format is invalid")

	ErrUsernameRequired      = errors.New("Username is required")
	ErrUsernameMinMaxLength  = errors.New("Username must be between 3 and 30 characters")
	ErrUsernameFormatInvalid = errors.New("Username format is invalid")

	ErrUserIDRequired = errors.New("UserID is required")
	ErrUserIDInvalid  = errors.New("UserID format is invalid")
)
