package bootstrap

import (
	"github.com/thanhnamdk2710/auth-service/internal/application/port"
	"github.com/thanhnamdk2710/auth-service/internal/application/usecase"
	infralogger "github.com/thanhnamdk2710/auth-service/internal/infrastructure/logger"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/uuid"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/handler"
)

type Handlers struct {
	Auth *handler.AuthHandler
}

func NewHandlers(db *Database, auditLogger port.AuditLogger, log *logger.Logger) *Handlers {
	// Infrastructure layer
	userRepo := postgres.NewPostgreUserRepo(db.Conn())
	uuidGenerator := uuid.NewGenerator()
	logAdapter := infralogger.NewAdapter(log)

	// Application layer
	registerUC := usecase.NewRegisterUsecase(userRepo, auditLogger, logAdapter, uuidGenerator)

	// Presentation layer
	authHandler := handler.NewAuthHandler(registerUC, logAdapter)

	return &Handlers{
		Auth: authHandler,
	}
}
