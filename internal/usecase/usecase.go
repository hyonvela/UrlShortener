package usecase

import (
	"context"

	"example.com/m/internal/storage"
	"example.com/m/pkg/logging"
)

type Usecase interface {
	ShortenUrl(long string, ctx context.Context) (string, error)
	GetLongUrl(short string, ctx context.Context) (string, error)
}

func New(storage storage.Storage, logger *logging.Logger) URLusecase {
	return URLusecase{storage, logger}
}
