package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thanhnamdk2710/auth-service/internal/pkg/metrics"
)

func Metrics(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		m.HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		m.HTTPRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
	}
}
