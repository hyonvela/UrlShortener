package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"example.com/m/config"
	"example.com/m/internal/app"
	"example.com/m/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.Logs_format)

	ctx, cancel := context.WithCancel(context.Background())

	// Контекст с отменой нужен для завершения работы при сигнале ^C
	// Создается горутина ожидающая ^C, которая отменит контекст

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan
		cancel()
	}()

	log.Info("Starting app")
	app.Run(ctx, cfg, log)
}
