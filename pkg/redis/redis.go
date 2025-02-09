package redisCache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func New(host string, port string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   db,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("%s:%s", host, port)
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	return rdb
}
