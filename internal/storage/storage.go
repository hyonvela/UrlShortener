package storage

import (
	"example.com/m/config"
	"example.com/m/pkg/logging"
	postgersDB "example.com/m/pkg/postgres"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
	// cache    *redis.Client
	// producer *kafka.Producer
}

func New(cfg *config.Config, logger *logging.Logger) *Storage {
	logger.Info("Connecting to database")
	dsn := cfg.GetDSN()
	logger.Info(dsn)
	db, err := postgersDB.New(dsn)
	if err != nil {
		logger.Fatalf("Faild to connect to database. error: %v", err)
	}

	// logger.Info("Connecting to redis client")
	// redis := redisCache.New(cfg.Redis.RedisHost, cfg.Redis.RedisPort, cfg.Redis.RedisDB)

	// logger.Info("Connecting to kafka")
	// producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	// if err != nil {
	// 	logger.Errorf("Faild to connect to kafka. error: %v", err)
	// }

	return &Storage{db}
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
	// r.producer.Close()

	// err := r.cache.Close()
	// if err != nil {
	// 	return err
	// }

	return r.db.Close()
}
