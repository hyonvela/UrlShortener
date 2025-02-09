package tests

import (
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
	s := storage.New(cfg, logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL))
	defer s.Close()

	t.Run("DB test", func(t *testing.T) {
		longURL := RandString()
		id := usecase.GetID(longURL)
		shortURL := urlshortener.Shorten(id)

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
