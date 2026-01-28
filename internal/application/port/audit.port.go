package port

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type AuditLogger interface {
	Log(ctx context.Context, log *entity.AuditLog)
	Start()
	Stop()
}
