package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	cfg := config.GetConfig()
	s := storage.New(cfg, logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL))
	defer s.Close()

	ctx := context.Background()

	t.Run("Get test", func(t *testing.T) {
		longURL := fmt.Sprint(rand.Uint32())
		id := usecase.GetID(longURL)
		var shortURL string
		var result string

		err := s.GetShortUrl(id, longURL, &shortURL, ctx)
		require.Error(t, err, "sql: no rows in result set")

		err = s.GetLongUrl(shortURL, &result, ctx)
		require.NoError(t, err, "sql: no rows in result set")
		require.Equal(t, longURL, result)
	})

	t.Run("Insert test", func(t *testing.T) {
		longURL := fmt.Sprint(rand.Uint32())
		shortURL := fmt.Sprint(rand.Uint32())
		id := usecase.GetID(longURL)
		var (
			findedShort string
			findedLong  string
		)

		err := s.InsertShortUrl(id, longURL, shortURL, ctx)
		require.NoError(t, err)

		err = s.GetShortUrl(id, longURL, &findedShort, ctx)
		require.NoError(t, err)
		require.Equal(t, shortURL, findedShort)

		err = s.GetLongUrl(shortURL, &findedLong, ctx)
		require.NoError(t, err)
		require.Equal(t, longURL, findedLong)

	})
}
