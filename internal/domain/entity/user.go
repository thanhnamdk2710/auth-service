package entity

import (
	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	vo "github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type User struct {
	ID              vo.UserID
	Username        vo.Username
	Email           vo.Email
	IsActive        bool
	IsEmailVerified bool
}

func NewUser(id vo.UserID, username vo.Username, email vo.Email) *User {
	return &User{
		ID:              id,
		Username:        username,
		Email:           email,
		IsActive:        true,
		IsEmailVerified: false,
	}
}

func (u *User) Activate() error {
	if u.IsActive {
		return exception.ErrUserAlreadyActive
	}
	u.IsActive = true
	return nil
}

func (u *User) Deactivate() error {
	if !u.IsActive {
		return exception.ErrUserAlreadyInactive
	}
	u.IsActive = false
	return nil
}

func (u *User) VerifyEmail() error {
	if !u.IsActive {
		return exception.ErrUserAlreadyInactive
	}

	if u.IsEmailVerified {
		return exception.ErrEmailAlreadyVerified
	}
	u.IsEmailVerified = true
	return nil
}
