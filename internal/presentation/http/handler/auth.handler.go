package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/port"
	"github.com/thanhnamdk2710/auth-service/internal/presentation/http/request"
	"github.com/thanhnamdk2710/auth-service/internal/validation"
)

type AuthHandler struct {
	registerUC port.RegisterUseCase
	logger     port.Logger
}

func NewAuthHandler(registerUC port.RegisterUseCase, logger port.Logger) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		logger:     logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		if errs := validation.TranslateAll(err); errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "validation failed",
				"errors":  errs,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.registerUC.Execute(ctx, input.RegisterInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: c.ClientIP(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": result.UserID,
		"message": result.Message,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login API",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Forgot Password API",
	})
}
