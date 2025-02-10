package tests

import (
	"context"
	"testing"

	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"example.com/m/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUsecase(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockLogger := logging.GetLogger("text", "test")
	uc := usecase.New(mockStorage, mockLogger)

	t.Run("ShortenUrl returns short URL", func(t *testing.T) {
		longURL := "https://example.com"
		shortURL := "abc123"
		mockStorage.ShortURL = shortURL
		mockStorage.Err = nil

		result, err := uc.ShortenUrl(longURL, context.Background())
		assert.NoError(t, err)
		assert.Equal(t, shortURL, result)
	})

	t.Run("GetLongUrl returns long URL", func(t *testing.T) {
		longURL := "https://example.com"
		shortURL := "abc123"
		mockStorage.LongURL = longURL
		mockStorage.Err = nil

		result, err := uc.GetLongUrl(shortURL, context.Background())
		assert.NoError(t, err)
		assert.Equal(t, longURL, result)
	})

}
