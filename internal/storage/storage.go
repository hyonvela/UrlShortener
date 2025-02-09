package storage

import (
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

	return &Storage{db, redis}
}

func (s *Storage) GetShortUrl(id uint32, short *string) error {
	return s.GetShortUrlFromDB(id, short)
}

func (s *Storage) InsertShortUrl(id uint32, url string, short string) error {
	return s.InsertShortUrlIntoDB(id, url, short)
}

func (s *Storage) GetLongUrl(short string, long *string) error {
	return s.GetLongUrlFromDB(short, long)
}

func (r *Storage) Close() error {
	err := r.cache.Close()
	if err != nil {
		return err
	}

	return r.db.Close()
}
