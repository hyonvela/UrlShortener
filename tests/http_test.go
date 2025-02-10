package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/config"
	"example.com/m/internal/app"
	"example.com/m/internal/entity"
	"example.com/m/internal/storage"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {
	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.LogsFormat, cfg.LogsLVL)
	s := storage.New(cfg, log)
	defer s.Close()

	uc := usecase.New(s, log)

	gin.SetMode(gin.TestMode)
	r := app.SetupRouter(uc, log)

	type testCase struct {
		short string
		long  string
	}

	testCases := []testCase{
		{"2lb2PU0000", "QUmhyWlX47a7zwU"},
		{"1JLGSP0000", "WbOvG8lS1ssWaYJyhqpi8lNdzLAPH6Z8i2FPMTqd3B0jlPwpWyHglMN1AEMijuqef2aRjfoItahzG1YPvQDt17eOsqqd7t4iWQAKoyMMhnJP1CpAEjMxoAjhwufmPisXdV77LeT08sTEb6ViU9Dmw1xMZkKUkDZe1C"},
		{"37poLS0000", "qbk0UibRb5gxKh8LP5evBjAbqg3zoRr9An8GtD_beFPtA6aMo2leXKbr0Ne188enpvjDm1yd_mr0CSf5NAILumfAZDbP60HpG2on4NJlUxIbVk64F6R96S1D8qFVgxqYEiuvaLWVBWZezRvqPE7JkeW8f_wO762AZfg_NzfBgIpWnC0nuzpuoP5Zj8lyPXfrJnQbFJ1UvPFVnuQsBly0_aZ"},
		{"QZItq00000", "7jgbBvpvAw14LITaDUWSYG4bas3i3a6b5ypIAt_ePxY1CHms3WUCtNJ0K5"},
		{"1bYJtg0000", "jx1PcuBMCRAGt_vO3rpTAC7W74"},
		{"3TICds0000", "21w0"},
		{"3lAGY80000", "TyRNIvPFGRfCa"},
		{"3ih_v00000", "uDpvPkxiIf1ozkE6_mHYJWmia0zpqitVdEA7S4Pt6VxKYnu4zP6vGj3zQvh4Zliv"},
		{"48zo8_0000", "p"},
		{"2XJW6E0000", "pGD"},
	}

	t.Run("ShortenUrl test", func(t *testing.T) {
		for _, testcase := range testCases {
			m := entity.LongUrl{LongUrl: testcase.long}
			jsonData, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("Ошибка при маршализации JSON: %v", err)
			}

			req, _ := http.NewRequest(http.MethodPost, "http://0.0.0.0:8080/v1/url_shortener", bytes.NewBuffer(jsonData))
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

			var answer entity.ShortUrl
			json.Unmarshal(resp.Body.Bytes(), &answer)
			assert.Equal(t, testcase.short, answer.ShortUrl)
		}

	})

	t.Run("GetLongUrl test", func(t *testing.T) {
		for _, testcase := range testCases {
			m := entity.ShortUrl{ShortUrl: testcase.short}
			jsonData, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("Ошибка при маршализации JSON: %v", err)
			}

			req, _ := http.NewRequest(http.MethodGet, "http://0.0.0.0:8080/v1/url_shortener", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

			var answer entity.LongUrl
			json.Unmarshal(resp.Body.Bytes(), &answer)
			assert.Equal(t, testcase.long, answer.LongUrl)
		}
	})
}
