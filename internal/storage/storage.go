package storage

import (
	"context"

	"example.com/m/config"
	"example.com/m/pkg/logging"
)

type Storage interface {
	GetShortURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, shortURL string) (string, error)
	GetLongURLByID(ctx context.Context, id uint32) (string, error)
	SaveURL(ctx context.Context, id uint32, longURL, shortURL string) error
	Close() error
}

func New(cfg *config.Config, log *logging.Logger) Storage {
	switch cfg.StorageType {
	case "postgres":
		return NewPostgresStorage(cfg, log)
	case "redis":
		return NewRedisStorage(cfg, log)
	default:
		panic("unknown storage type")
	}
}
