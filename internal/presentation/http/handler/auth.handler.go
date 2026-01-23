package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/usecase"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/request"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

type AuthHandler struct {
	registerUC usecase.RegisterUseCase
}

func NewAuthHandler(registerUC usecase.RegisterUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
	}
}

func (h AuthHandler) Register(ctx *gin.Context) {
	var req request.RegisterRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		if errs := validation.TranslateAll(err); errs != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "validation failed",
				"errors":  errs,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	uc, err := h.registerUC.Execute(input.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": uc,
	})
}

func (h AuthHandler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login API",
	})
}

func (h AuthHandler) ForgotPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Forgot Password API",
	})
}
