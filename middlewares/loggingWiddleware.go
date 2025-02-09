// middlewares/logging.go
package middlewares

import (
	"time"

	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{
			Latency:      time.Since(start),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		entry := logger.WithFields(logrus.Fields{
			"status":   param.StatusCode,
			"method":   param.Method,
			"path":     param.Path,
			"latency":  param.Latency,
			"clientIP": param.ClientIP,
		})

		if param.ErrorMessage != "" {
			entry = entry.WithField("error", param.ErrorMessage)
		}

		if param.StatusCode >= 500 {
			entry.Error()
		} else if param.StatusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
