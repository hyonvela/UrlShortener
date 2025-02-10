package mocks

import (
	"context"
)

type MockStorage struct {
	ShortURL string
	LongURL  string
	Err      error
}

func (m *MockStorage) GetShortURL(ctx context.Context, longURL string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.ShortURL, nil
}

func (m *MockStorage) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.LongURL, nil
}

func (m *MockStorage) GetLongURLByID(ctx context.Context, id uint32) (string, error) {
	return "", m.Err
}

func (m *MockStorage) SaveURL(ctx context.Context, id uint32, longURL, shortURL string) error {
	return m.Err
}

func (m *MockStorage) Close() error {
	return nil
}
