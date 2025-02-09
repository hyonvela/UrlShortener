package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
	urlshortener "example.com/m/internal/url_shortener"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestUsecase(t *testing.T) {
	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL)
	s := storage.New(cfg, log)
	defer s.Close()

	ctx := context.Background()
	uc := usecase.New(s, log)

	t.Run("Shortener test", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			longURL := fmt.Sprint(rand.Uint32())
			id := usecase.GetID(longURL)
			shortURL := urlshortener.Shorten(id)

			newShort, err := uc.ShortenUrl(longURL, ctx)
			require.NoError(t, err)
			require.Equal(t, shortURL, newShort)
		}
	})
}
