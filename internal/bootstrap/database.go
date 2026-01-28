package bootstrap

import (
	"context"
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
)

type Database struct {
	conn *postgres.DB
}

func NewDatabase(ctx context.Context, cfg *config.DBConfig) (*Database, error) {
	conn, err := postgres.NewConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (d *Database) Conn() *postgres.DB {
	return d.conn
}

func (d *Database) SQL() *sql.DB {
	return d.conn.DB
}

func (d *Database) Close() error {
	return d.conn.Close()
}
