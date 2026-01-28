package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/handler"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/middleware"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

type RouterDeps struct {
	Logger      *logger.Logger
	Metrics     *metrics.Metrics
	AuthHandler *handler.AuthHandler
}

func New(deps RouterDeps) *gin.Engine {
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

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimitDefault())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", deps.AuthHandler.Register)
			auth.POST("/login", deps.AuthHandler.Login)
			auth.POST("/forgot-password", deps.AuthHandler.ForgotPassword)
		}
	}

	return r
}
