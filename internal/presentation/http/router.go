package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/thanhnamdk2710/auth-service/internal/application/service"
	"github.com/thanhnamdk2710/auth-service/internal/application/usecase"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/handler"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/middleware"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

type Deps struct {
	Logger       *logger.Logger
	Metrics      *metrics.Metrics
	DB           *postgres.DB
	AuditService service.AuditService
}

func New(deps Deps) *gin.Engine {
	validation.Init()

	r := gin.New()

	r.Use(middleware.CorrelationID())
	r.Use(middleware.Recovery(deps.Logger))
	r.Use(middleware.Logging(deps.Logger))
	r.Use(middleware.Metrics(deps.Metrics))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	userRepo := postgres.NewPostgreUserRepo(deps.DB)
	registerUC := usecase.NewRegisterUsecase(userRepo, deps.AuditService, deps.Logger)
	authHandler := handler.NewAuthHandler(registerUC, deps.Logger)

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimitDefault())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.GET("/login", authHandler.Login)
			auth.GET("/forgot-password", authHandler.ForgotPassword)
		}
	}

	return r
}
