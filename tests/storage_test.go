package tests

import (
	"context"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
	urlshortener "example.com/m/internal/url_shortener"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	cfg := config.GetConfig()
	ctx := context.Background()
	storages := [2]string{"postgres", "redis"}

	for _, i := range storages {
		cfg.StorageType = i
		s := storage.New(cfg, logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL))

		t.Run("Get test", func(t *testing.T) {
			longURL := RandString()
			id := usecase.GetID(longURL)

			shortURL, err := s.GetShortURL(ctx, longURL)
			require.Error(t, err, storage.ErrNotFound)

			_, err = s.GetLongURL(ctx, shortURL)
			require.Error(t, err, storage.ErrNotFound)

			_, err = s.GetLongURLByID(ctx, id)
			require.Error(t, err, storage.ErrNotFound)
		})

		t.Run("Insert test", func(t *testing.T) {
			longURL := RandString()
			id := usecase.GetID(longURL)
			shortURL := urlshortener.Shorten(id)

			err := s.SaveURL(ctx, id, longURL, shortURL)
			require.NoError(t, err)

			findedShort, err := s.GetShortURL(ctx, longURL)
			require.NoError(t, err)
			require.Equal(t, shortURL, findedShort)

			findedLong, err := s.GetLongURL(ctx, shortURL)
			require.NoError(t, err)
			require.Equal(t, longURL, findedLong)

			findedLong, err = s.GetLongURLByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, longURL, findedLong)

		})

		s.Close()
	}

}
