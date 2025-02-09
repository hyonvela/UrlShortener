package app

import (
	"context"
	"net/http"
	"time"

	"example.com/m/config"
	"example.com/m/internal/handlers"
	"example.com/m/internal/storage"
	"example.com/m/internal/usecase"
	"example.com/m/middlewares"
	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
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
		Addr:         cfg.GetAdress(),
		Handler:      r.Handler(),
		ReadTimeout:  time.Duration(cfg.Listen.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Listen.WriteTimeout) * time.Second,
	}

	// Запуск http сервиса
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	// Поступил сигнал ^C
	// Запуск логики graceful shutdown

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func SetupRouter(uc *usecase.Usecase, logger *logging.Logger) *gin.Engine {
	// Настройка роутера gin

	// В качестве логгера выбран logrus, для него написсана обертка,
	// чтобы использовать logrus в качестве логгера в gin
	gin.DefaultWriter = logging.NewGinLogrusWriter(logger.Entry)
	gin.DefaultErrorWriter = logging.NewGinLogrusWriter(logger.Entry)

	router := gin.New()
	// Создание логера, подключение хэндлеров и мидлварей

	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
	}))

	s := handlers.NewHandler(uc, logger)

	v1 := router.Group("/v1")
	{
		v1.POST("url_shortener", s.GetLongUrl)
		v1.GET("url_shortener/:url", s.ShortenUrl)
	}

	return router
}
