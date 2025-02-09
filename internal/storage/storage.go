package storage

import (
	"context"

	"example.com/m/config"
	"example.com/m/pkg/logging"
	postgersDB "example.com/m/pkg/postgres"
	redisCache "example.com/m/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db    *sqlx.DB
	cache *redis.Client
	cfg   *config.Config
	log   *logging.Logger
}

func New(cfg *config.Config, log *logging.Logger) *Storage {
	log.Info("Connecting to database")
	dsn := cfg.GetDSN()
	log.Info(dsn)
	db, err := postgersDB.New(dsn)
	if err != nil {
		log.Fatalf("Faild to connect to database. error: %v", err)
	}

	log.Info("Connecting to redis client")
	redis := redisCache.New(cfg.Redis.RedisHost, cfg.Redis.RedisPort, cfg.Redis.RedisDB)

	return &Storage{db, redis, cfg, log}
}

func (s *Storage) GetShortUrl(id uint32, long string, short *string, ctx context.Context) error {
	err := s.GetShortUrlFromCache(long, short, ctx)

	if err != nil {
		err = s.GetShortUrlFromDB(id, short)
		s.log.Info(s.InsertUrlsIntoCache(long, *short, ctx))
		return err
	}

	return nil
}

func (s *Storage) InsertShortUrl(id uint32, long string, short string, ctx context.Context) error {
	err := s.InsertUrlsIntoCache(long, short, ctx)
	if err != nil {
		s.log.Errorf("Cache error: %s", err.Error())
	}

	return s.InsertShortUrlIntoDB(id, long, short)
}

func (s *Storage) GetLongUrl(short string, long *string, ctx context.Context) error {
	err := s.GetLongUrlFromCache(short, long, ctx)

	if err != nil {
		err = s.GetLongUrlFromDB(short, long)
		s.log.Info(s.InsertUrlsIntoCache(*long, short, ctx))
		return err
	}

	return nil
}

func (r *Storage) Close() error {
	err := r.cache.Close()
	if err != nil {
		return err
	}

	return r.db.Close()
}
