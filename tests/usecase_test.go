package tests

import (
	"context"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/storage"
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
			longURL := RandString()
			// Проверяем, что два вызова для одного URL возвращают одинаковый short
			short1, err := uc.ShortenUrl(longURL, ctx)
			require.NoError(t, err)
			short2, err := uc.ShortenUrl(longURL, ctx)
			require.NoError(t, err)
			require.Equal(t, short1, short2)

			// Проверяем, что short корректно восстанавливается в long
			var long string
			long, err = uc.GetLongUrl(short1, ctx)
			require.NoError(t, err)
			require.Equal(t, longURL, long)
		}
	})
}
