package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	cfg := config.GetConfig()
	s := storage.New(cfg, logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL))
	defer s.Close()

	t.Run("DB test", func(t *testing.T) {
		longURL := fmt.Sprint(rand.Uint32())
		shortURL := fmt.Sprint(rand.Uint32())
		id := usecase.GetID(longURL)
		var (
			findedShort string
			findedLong  string
		)

		err := s.GetLongUrlFromDB(shortURL, &findedLong)
		require.Error(t, err)

		err = s.GetShortUrlFromDB(id, &findedShort)
		require.Error(t, err)

		err = s.InsertShortUrlIntoDB(id, longURL, shortURL)
		require.NoError(t, err)

		err = s.GetLongUrlFromDB(shortURL, &findedLong)
		require.NoError(t, err)
		require.Equal(t, longURL, findedLong)

		err = s.GetShortUrlFromDB(id, &findedShort)
		require.NoError(t, err)
		require.Equal(t, shortURL, findedShort)
	})
}
