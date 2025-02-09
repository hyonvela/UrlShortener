package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
	"example.com/m/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cfg := config.GetConfig()
	s := storage.New(cfg, logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL))
	defer s.Close()

	ctx := context.Background()

	t.Run("Cache test", func(t *testing.T) {
		longURL := fmt.Sprint(rand.Uint32())
		shortURL := fmt.Sprint(rand.Uint32())
		var (
			findedShort string
			findedLong  string
		)

		err := s.GetLongUrlFromCache(shortURL, &findedLong, ctx)
		require.Error(t, err)

		err = s.GetShortUrlFromCache(longURL, &findedShort, ctx)
		require.Error(t, err)

		err = s.InsertUrlsIntoCache(longURL, shortURL, ctx)
		require.NoError(t, err)

		err = s.GetLongUrlFromCache(shortURL, &findedLong, ctx)
		require.NoError(t, err)
		require.Equal(t, longURL, findedLong)

		err = s.GetShortUrlFromCache(longURL, &findedShort, ctx)
		require.NoError(t, err)
		require.Equal(t, shortURL, findedShort)
	})

}
