package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

func (s *Storage) GetShortUrlFromCache(long string, short *string, ctx context.Context) error {
	data, err := s.cache.Get(ctx, long).Result()
	if err != nil {
		return err
	}
	if data == "null" {
		return redis.Nil
	}

	*short = data
	return nil
}

func (s *Storage) GetLongUrlFromCache(short string, long *string, ctx context.Context) error {
	data, err := s.cache.Get(ctx, short).Result()
	if err != nil {
		return err
	}
	if data == "null" {
		return redis.Nil
	}

	*long = data
	return nil
}

func (s *Storage) InsertUrlsIntoCache(long string, short string, ctx context.Context) error {
	err := s.cache.Set(ctx, short, long, time.Duration(s.cfg.Redis.LifeTime)*time.Second).Err()
	if err != nil {
		return err
	}
	err = s.cache.Set(ctx, long, short, time.Duration(s.cfg.Redis.LifeTime)*time.Second).Err()
	return err
}
