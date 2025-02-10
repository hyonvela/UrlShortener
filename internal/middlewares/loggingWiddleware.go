// middlewares/logging.go
package middlewares

import (
	"context"
	"time"

	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func HTTPLoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
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

func GRPCLoggingMiddleware(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		clientIP := getClientIP(ctx)

		statusCode := status.Code(err)

		logger.WithFields(logrus.Fields{
			"method":     info.FullMethod,
			"took":       time.Since(start),
			"clientIP":   clientIP,
			"statusCode": statusCode,
		}).Info("Handled request")

		if err != nil {
			logger.WithField("error", err.Error()).Error("Request failed")
		}

		return resp, err
	}
}

func getClientIP(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ip := md["client-ip"]; len(ip) > 0 {
			return ip[0]
		}
	}
	return ""
}
