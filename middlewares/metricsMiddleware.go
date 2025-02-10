package middlewares

import (
	"context"
	"net/http"
	"time"

	"example.com/m/pkg/metrics"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func HTTPMetricsMiddleware() gin.HandlerFunc {
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

func GRPCMetricsMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		method := info.FullMethod
		endpoint := method

		duration := time.Since(start).Seconds()

		metrics.RequestCount.WithLabelValues(method, endpoint).Inc()

		metrics.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)

		if err != nil {
			metrics.ErrorCount.WithLabelValues(status.Code(err).String()).Inc()
		}

		return resp, err
	}
}
