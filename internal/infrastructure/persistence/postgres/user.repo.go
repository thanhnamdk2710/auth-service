package postgres

import (
	"context"
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
)

type PostgreUserRepo struct {
	db *DB
}

func NewPostgreUserRepo(db *DB) repository.UserRepository {
	return &PostgreUserRepo{db: db}
}

func (r *PostgreUserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, is_active, is_email_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(),
		user.Username.String(),
		user.Email.String(),
		"", // password_hash - will be added later
		user.IsActive,
		user.IsEmailVerified,
	)

	return err
}

func (r *PostgreUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, username, email, is_active, is_email_verified
		FROM users WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanUser(row)
}

func (r *PostgreUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, username, email, is_active, is_email_verified
		FROM users WHERE username = $1
	`

	row := r.db.QueryRowContext(ctx, query, username)
	return scanUser(row)
}

func (r *PostgreUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, username, email, is_active, is_email_verified
		FROM users WHERE email = $1
	`

	row := r.db.QueryRowContext(ctx, query, email)
	return scanUser(row)
}

func (r *PostgreUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (r *PostgreUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *PostgreUserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, is_active = $4, is_email_verified = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(),
		user.Username.String(),
		user.Email.String(),
		user.IsActive,
		user.IsEmailVerified,
	)

	return err
}

func scanUser(row *sql.Row) (*entity.User, error) {
	var id, username, email string
	var isActive, isEmailVerified bool

	err := row.Scan(&id, &username, &email, &isActive, &isEmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	userID, err := vo.NewUserID(id)
	if err != nil {
		return nil, err
	}

	userUsername, err := vo.NewUsername(username)
	if err != nil {
		return nil, err
	}

	userEmail, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:              userID,
		Username:        *userUsername,
		Email:           userEmail,
		IsActive:        isActive,
		IsEmailVerified: isEmailVerified,
	}, nil
}
