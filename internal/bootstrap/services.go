package bootstrap

import (
	"github.com/thanhnamdk2710/auth-service/internal/application/port"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/audit"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type Services struct {
	audit port.AuditLogger
}

func NewServices(db *Database, log *logger.Logger) *Services {
	auditRepo := postgres.NewAuditRepo(db.Conn())
	auditLogger := audit.NewAsyncLogger(
		auditRepo,
		log,
		audit.DefaultConfig(),
	)

	return &Services{
		audit: auditLogger,
	}
}

func (s *Services) Audit() port.AuditLogger {
	return s.audit
}

func (s *Services) Start() {
	s.audit.Start()
}

func (s *Services) Stop() {
	s.audit.Stop()
}
