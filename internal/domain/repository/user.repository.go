package repository

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *entity.User) error
}
