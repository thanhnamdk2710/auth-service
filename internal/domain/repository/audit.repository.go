package repository

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type AuditRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error)
	FindByCorrelationID(ctx context.Context, correlationID string) ([]*entity.AuditLog, error)
}
