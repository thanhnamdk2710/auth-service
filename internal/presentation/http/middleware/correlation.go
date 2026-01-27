package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/thanhnamdk2710/auth-service/internal/pkg/correlationid"
)

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		corrID := c.GetHeader(correlationid.HeaderName)
		if corrID == "" {
			corrID = correlationid.New()
		}

		ctx := correlationid.WithContext(c.Request.Context(), corrID)
		c.Request = c.Request.WithContext(ctx)

		c.Header(correlationid.HeaderName, corrID)

		c.Next()
	}
}
