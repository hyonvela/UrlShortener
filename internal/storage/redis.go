package storage

import (
	"context"
	"fmt"

	"example.com/m/config"
	"example.com/m/pkg/logging"
	redisCache "example.com/m/pkg/redis"
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	db  *redis.Client
	log *logging.Logger
}

func NewRedisStorage(cfg *config.Config, log *logging.Logger) *RedisStorage {
	r := redisCache.New(cfg.Redis.RedisHost, cfg.Redis.RedisPort, cfg.Redis.RedisDB)
	return &RedisStorage{r, log}
}

func (s *RedisStorage) GetShortURL(ctx context.Context, longURL string) (string, error) {
	shortURL, err := s.db.Get(ctx, longURL).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return shortURL, err
}

func (s *RedisStorage) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := s.db.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return longURL, err
}

func (s *RedisStorage) GetLongURLByID(ctx context.Context, id uint32) (string, error) {
	longURL, err := s.db.Get(ctx, fmt.Sprint(id)).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return longURL, err
}

func (s *RedisStorage) SaveURL(ctx context.Context, id uint32, longURL, shortURL string) error {
	err := s.db.Set(ctx, longURL, shortURL, 0).Err()
	if err != nil {
		return err
	}
	err = s.db.Set(ctx, shortURL, longURL, 0).Err()
	if err != nil {
		return err
	}
	return s.db.Set(ctx, fmt.Sprint(id), longURL, 0).Err()
}

func (s *RedisStorage) Close() error {
	return s.db.Close()
}
