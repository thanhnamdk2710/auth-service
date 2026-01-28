package postgres

import (
	"context"
	"database/sql"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
)

type AuditRepo struct {
	db *DB
}

func NewAuditRepo(db *DB) repository.AuditRepository {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, timestamp, user_id, action, details, ip_address, correlation_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.Timestamp,
		log.UserID,
		log.Action,
		log.Details,
		sql.NullString{String: log.IPAddress, Valid: log.IPAddress != ""},
		log.CorrelationID,
	)

	return err
}

func (r *AuditRepo) CreateBatch(ctx context.Context, logs []*entity.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO audit_logs (id, timestamp, user_id, action, details, ip_address, correlation_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, log := range logs {
		_, err := stmt.ExecContext(ctx,
			log.ID,
			log.Timestamp,
			log.UserID,
			log.Action,
			log.Details,
			sql.NullString{String: log.IPAddress, Valid: log.IPAddress != ""},
			log.CorrelationID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AuditRepo) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error) {
	query := `
		SELECT id, timestamp, user_id, action, details, ip_address, correlation_id
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAuditLogs(rows)
}

func (r *AuditRepo) FindByCorrelationID(ctx context.Context, correlationID string) ([]*entity.AuditLog, error) {
	query := `
		SELECT id, timestamp, user_id, action, details, ip_address, correlation_id
		FROM audit_logs
		WHERE correlation_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, correlationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAuditLogs(rows)
}

func scanAuditLogs(rows *sql.Rows) ([]*entity.AuditLog, error) {
	var logs []*entity.AuditLog

	for rows.Next() {
		var log entity.AuditLog
		var userID sql.NullString
		var ipAddress sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&userID,
			&log.Action,
			&log.Details,
			&ipAddress,
			&log.CorrelationID,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			log.UserID = &userID.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}
