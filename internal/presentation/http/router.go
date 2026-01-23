package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/application/usecase"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/handler"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

func NewRouter() *gin.Engine {
	validation.Init()

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	api := r.Group("api/v1")
	{
		auth := api.Group("auth")
		{
			userRepo := postgres.NewPostgreUserRepo()
			registerUC := usecase.NewRegisterUsecase(userRepo)
			authHandler := handler.NewAuthHandler(registerUC)

			auth.POST("register", authHandler.Register)
			auth.GET("login", authHandler.Login)
			auth.GET("forgot-password", authHandler.ForgotPassword)
		}
	}

	return r
}
