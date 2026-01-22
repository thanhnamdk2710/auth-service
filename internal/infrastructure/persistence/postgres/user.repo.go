package postgres

import (
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type PostgreUserRepo struct {
	db *sql.DB
}

func NewPostgreUserRepo() *PostgreUserRepo {
	return &PostgreUserRepo{}
}

func (r *PostgreUserRepo) Create(user entity.User) entity.User {
	return user
}

func (r *PostgreUserRepo) FindByUsername(username string) entity.User {
	return entity.User{}
}

func (r *PostgreUserRepo) FindByEmail(email string) entity.User {
	return entity.User{}
}
