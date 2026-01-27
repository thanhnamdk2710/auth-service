package app

import (
	"github.com/thanhnamdk2710/auth-service/internal/application/service"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type Services struct {
	audit service.AuditService
}

func NewServices(db *Database, log *logger.Logger) *Services {
	auditRepo := postgres.NewAuditRepo(db.Conn())
	auditService := service.NewAuditService(
		auditRepo,
		log,
		service.DefaultAuditServiceConfig(),
	)

	return &Services{
		audit: auditService,
	}
}

func (s *Services) Audit() service.AuditService {
	return s.audit
}

func (s *Services) Start() {
	s.audit.Start()
}

func (s *Services) Stop() {
	s.audit.Stop()
}
