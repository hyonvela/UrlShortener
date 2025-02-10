package middlewares

import (
	"net/http"
	"time"

	"example.com/m/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware собирает метрики для всех HTTP запросов
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		method := c.Request.Method
		endpoint := c.FullPath()

		metrics.RequestCount.WithLabelValues(method, endpoint).Inc()
		metrics.RequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())

		if c.Writer.Status() >= 400 {
			metrics.ErrorCount.WithLabelValues(http.StatusText(c.Writer.Status())).Inc()
		}
	}
}
