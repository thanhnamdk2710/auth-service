package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/thanhnamdk2710/auth-service/internal/config"
)

type DB struct {
	*sql.DB
}

func NewConnection(ctx context.Context, cfg *config.DBConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Stats() sql.DBStats {
	return db.DB.Stats()
}
