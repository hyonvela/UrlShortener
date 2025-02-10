package mocks

import (
	"context"
)

type MockUsecase struct {
	ShortURL string
	LongURL  string
	Err      error
}

func (m *MockUsecase) ShortenUrl(long string, ctx context.Context) (string, error) {
	return m.ShortURL, m.Err
}

func (m *MockUsecase) GetLongUrl(short string, ctx context.Context) (string, error) {
	return m.LongURL, m.Err
}
