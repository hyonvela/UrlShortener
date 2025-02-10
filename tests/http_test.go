package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/internal/handlers"
	"example.com/m/internal/storage"
	"example.com/m/pkg/logging"
	"example.com/m/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPHandlers(t *testing.T) {
	mockUsecase := new(mocks.MockUsecase)
	mockLogger := logging.GetLogger("text", "test")
	handler := handlers.NewHandler(mockUsecase, mockLogger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/v1")
	{
		v1.POST("url_shortener", handler.ShortenUrl)
		v1.GET("url_shortener", handler.GetLongUrl)
	}

	t.Run("ShortenUrl success", func(t *testing.T) {
		mockUsecase.ShortURL = "abc123"
		mockUsecase.Err = nil

		longURL := map[string]string{"long_url": "https://example.com"}
		jsonData, _ := json.Marshal(longURL)

		req, _ := http.NewRequest(http.MethodPost, "/v1/url_shortener", bytes.NewBuffer(jsonData))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]string
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "abc123", response["short_url"])
	})

	t.Run("ShortenUrl error", func(t *testing.T) {
		mockUsecase.Err = storage.ErrNotFound

		longURL := map[string]string{"long_url": "https://example.com"}
		jsonData, _ := json.Marshal(longURL)

		req, _ := http.NewRequest(http.MethodPost, "/v1/url_shortener", bytes.NewBuffer(jsonData))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("GetLongUrl success", func(t *testing.T) {
		mockUsecase.LongURL = "https://example.com"
		mockUsecase.Err = nil

		shortURL := map[string]string{"short_url": "abc123"}
		jsonData, _ := json.Marshal(shortURL)

		req, _ := http.NewRequest(http.MethodGet, "/v1/url_shortener", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]string
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "https://example.com", response["long_url"])
	})

	t.Run("GetLongUrl error", func(t *testing.T) {
		mockUsecase.Err = storage.ErrNotFound

		shortURL := map[string]string{"short_url": "abc123"}
		jsonData, _ := json.Marshal(shortURL)

		req, _ := http.NewRequest(http.MethodGet, "/v1/url_shortener", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
