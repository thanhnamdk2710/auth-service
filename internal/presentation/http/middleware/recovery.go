package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()
				log.ErrorCtx(ctx, "Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
