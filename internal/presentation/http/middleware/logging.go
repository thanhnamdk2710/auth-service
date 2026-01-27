package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

func Logging(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		ctx := c.Request.Context()

		log.InfoCtx(ctx, "Request started",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.InfoCtx(ctx, "Request completed",
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}
