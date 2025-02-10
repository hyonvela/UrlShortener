package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"example.com/m/config"
	"example.com/m/internal/handlers"
	"example.com/m/internal/middlewares"
	"example.com/m/internal/storage"
	"example.com/m/internal/usecase"
	grpcV1 "example.com/m/pkg/grpc.v1"
	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config, log *logging.Logger) {
	log.Info("Storage initialization")
	s := storage.New(cfg, log)

	log.Info("Storage initialization")
	uc := usecase.New(s, log)

	log.Info("Router initialization")
	gin.SetMode(gin.DebugMode)
	r := SetupRouter(uc, log)

	log.Info("REST initialization")
	server := &http.Server{
		Addr:         cfg.GetHTTPAdress(),
		Handler:      r.Handler(),
		ReadTimeout:  time.Duration(cfg.Listen.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Listen.WriteTimeout) * time.Second,
	}

	// Запуск http сервиса
	go func() {
		log.Infof("HTTP server listening on %s", cfg.Listen.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Info("GRPC initialization")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middlewares.GRPCLoggingMiddleware(log),
				middlewares.GRPCMetricsMiddleware(),
			),
		),
	)
	grpcHandler := handlers.NewGRPCHandler(uc, log)
	grpcV1.RegisterUrlShortenerServiceServer(grpcServer, grpcHandler)

	grpcListener, err := net.Listen("tcp", cfg.Listen.BindIp+":"+cfg.Listen.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Запуск grpc сервиса
	go func() {
		log.Infof("GRPC server listening on %s", cfg.Listen.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()

	<-ctx.Done()
	// Поступил сигнал ^C
	// Запуск логику graceful shutdown

	log.Println("Shutdown Servers ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	grpcServer.GracefulStop()

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func SetupRouter(uc usecase.Usecase, logger *logging.Logger) *gin.Engine {
	// Настройка роутера gin

	// В качестве логгера выбран logrus, для него написсана обертка,
	// чтобы использовать logrus в качестве логгера в gin
	gin.DefaultWriter = logging.NewGinLogrusWriter(logger.Entry)
	gin.DefaultErrorWriter = logging.NewGinLogrusWriter(logger.Entry)

	router := gin.New()
	// Создание логера, подключение хэндлеров и мидлварей

	router.Use(middlewares.HTTPLoggingMiddleware(logger))
	router.Use(middlewares.HTTPMetricsMiddleware())
	router.Use(gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
	}))

	s := handlers.NewHTTPHandler(uc, logger)

	v1 := router.Group("/v1")
	{
		v1.POST("url_shortener", s.ShortenUrl)
		v1.GET("url_shortener", s.GetLongUrl)
		v1.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	return router
}
