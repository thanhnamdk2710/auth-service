package repository

import "github.com/thanhnamdk2710/auth-service/internal/domain/entity"

type UserRepository interface {
	Create(entity.User) entity.User
	FindByUsername(string) entity.User
	FindByEmail(string) entity.User
}
